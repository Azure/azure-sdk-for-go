// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys_test

import (
	"context"
	"crypto/sha3"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys/fake"
	"github.com/stretchr/testify/require"
)

// mlDSAExternalMu computes the ML-DSA message representative mu for pure (non-prehash) signing as defined
// in FIPS 204 Algorithm 2 and Algorithm 7: mu = SHAKE256(SHAKE256(pk, 64) || 0 || len(ctx) || ctx || msg, 64).
func mlDSAExternalMu(t *testing.T, publicKey, mlDSAContext, msg []byte) []byte {
	require.LessOrEqual(t, len(mlDSAContext), 255, "ML-DSA context must be at most 255 bytes")
	tr := sha3.SumSHAKE256(publicKey, 64)
	input := make([]byte, 0, len(tr)+2+len(mlDSAContext)+len(msg))
	input = append(input, tr...)
	input = append(input, 0, byte(len(mlDSAContext)))
	input = append(input, mlDSAContext...)
	input = append(input, msg...)
	return sha3.SumSHAKE256(input, 64)
}

// TestMLDSAModelsSerde verifies the wire format of the ML-DSA additions against the 2026-07-01-preview REST API examples
func TestMLDSAModelsSerde(t *testing.T) {
	b64 := base64.RawURLEncoding.EncodeToString

	t.Run("CreateKeyParameters", func(t *testing.T) {
		params := azkeys.CreateKeyParameters{
			Alg:           to.Ptr(azkeys.AKPAlgorithmMLDSA65),
			KeyAttributes: &azkeys.KeyAttributes{},
			KeyOps:        to.SliceOfPtrs(azkeys.KeyOperationSign, azkeys.KeyOperationVerify),
			Kty:           to.Ptr(azkeys.KeyTypeAKP),
			Tags:          map[string]*string{"purpose": to.Ptr("unit test"), "test name": to.Ptr("CreateMLDSAKeyTest")},
		}
		data, err := json.Marshal(params)
		require.NoError(t, err)
		require.JSONEq(t, `{
			"kty": "AKP",
			"alg": "ML-DSA-65",
			"key_ops": ["sign", "verify"],
			"attributes": {},
			"tags": {"purpose": "unit test", "test name": "CreateMLDSAKeyTest"}
		}`, string(data))
		testSerde(t, &params)
	})

	t.Run("KeyBundle", func(t *testing.T) {
		var bundle azkeys.KeyBundle
		require.NoError(t, json.Unmarshal([]byte(`{
			"key": {
				"kid": "https://myvault.vault.azure.net/keys/CreateMLDSAKeyTest/a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4",
				"kty": "AKP-HSM",
				"alg": "ML-DSA-65",
				"key_ops": ["sign", "verify"],
				"pub": "base64url-encoded-public-key-bytes"
			},
			"attributes": {"enabled": true, "created": 1706000000, "updated": 1706000000, "recoveryLevel": "Recoverable+Purgeable"}
		}`), &bundle))
		require.NotNil(t, bundle.Key)
		require.Equal(t, azkeys.KeyTypeAKPHSM, *bundle.Key.Kty)
		require.Equal(t, azkeys.AKPAlgorithmMLDSA65, *bundle.Key.Alg)
		require.Equal(t, to.SliceOfPtrs(azkeys.KeyOperationSign, azkeys.KeyOperationVerify), bundle.Key.KeyOps)
		expectedPub, err := base64.RawURLEncoding.DecodeString("base64url-encoded-public-key-bytes")
		require.NoError(t, err)
		require.Equal(t, expectedPub, bundle.Key.Pub)
		require.Equal(t, "CreateMLDSAKeyTest", bundle.Key.KID.Name())
		testSerde(t, bundle.Key)
	})

	msg := []byte("EA7788548B49F1E7E6702C8C4C02BC90516295362948A76CA9E9354063FF2962")
	mlDSAContext := []byte("application-context-v1")
	signature := []byte("ML-DSA signature")
	mu := make([]byte, 64)
	for i := range mu {
		mu[i] = byte(i)
	}

	t.Run("SignParameters", func(t *testing.T) {
		// pure ML-DSA signing sends the message without an algorithm; the key's algorithm is implied
		params := azkeys.SignParameters{Context: mlDSAContext, Value: msg}
		data, err := json.Marshal(params)
		require.NoError(t, err)
		require.JSONEq(t, `{
			"value": "RUE3Nzg4NTQ4QjQ5RjFFN0U2NzAyQzhDNEMwMkJDOTA1MTYyOTUzNjI5NDhBNzZDQTlFOTM1NDA2M0ZGMjk2Mg",
			"context": "YXBwbGljYXRpb24tY29udGV4dC12MQ"
		}`, string(data))
		testSerde(t, &params)

		// external mu mode sends only the precomputed mu
		params = azkeys.SignParameters{ExternalMu: mu}
		data, err = json.Marshal(params)
		require.NoError(t, err)
		require.JSONEq(t, `{"external_mu": "`+b64(mu)+`"}`, string(data))
		testSerde(t, &params)

		var roundTripped azkeys.SignParameters
		require.NoError(t, json.Unmarshal(data, &roundTripped))
		require.Equal(t, mu, roundTripped.ExternalMu)
		require.Nil(t, roundTripped.Algorithm)
		require.Nil(t, roundTripped.Value)
	})

	t.Run("VerifyParameters", func(t *testing.T) {
		params := azkeys.VerifyParameters{Context: mlDSAContext, Digest: msg, Signature: signature}
		data, err := json.Marshal(params)
		require.NoError(t, err)
		require.JSONEq(t, `{
			"digest": "RUE3Nzg4NTQ4QjQ5RjFFN0U2NzAyQzhDNEMwMkJDOTA1MTYyOTUzNjI5NDhBNzZDQTlFOTM1NDA2M0ZGMjk2Mg",
			"value": "`+b64(signature)+`",
			"context": "YXBwbGljYXRpb24tY29udGV4dC12MQ"
		}`, string(data))
		testSerde(t, &params)

		params = azkeys.VerifyParameters{ExternalMu: mu, Signature: signature}
		data, err = json.Marshal(params)
		require.NoError(t, err)
		require.JSONEq(t, `{"external_mu": "`+b64(mu)+`", "value": "`+b64(signature)+`"}`, string(data))
		testSerde(t, &params)
	})
}

func TestAKPConstants(t *testing.T) {
	require.Equal(t, []azkeys.AKPAlgorithm{
		azkeys.AKPAlgorithmMLDSA44,
		azkeys.AKPAlgorithmMLDSA65,
		azkeys.AKPAlgorithmMLDSA87,
	}, azkeys.PossibleAKPAlgorithmValues())
	require.Equal(t, azkeys.AKPAlgorithm("ML-DSA-44"), azkeys.AKPAlgorithmMLDSA44)
	require.Equal(t, azkeys.AKPAlgorithm("ML-DSA-65"), azkeys.AKPAlgorithmMLDSA65)
	require.Equal(t, azkeys.AKPAlgorithm("ML-DSA-87"), azkeys.AKPAlgorithmMLDSA87)

	require.Contains(t, azkeys.PossibleKeyTypeValues(), azkeys.KeyTypeAKP)
	require.Contains(t, azkeys.PossibleKeyTypeValues(), azkeys.KeyTypeAKPHSM)
	require.Equal(t, azkeys.KeyType("AKP"), azkeys.KeyTypeAKP)
	require.Equal(t, azkeys.KeyType("AKP-HSM"), azkeys.KeyTypeAKPHSM)
}

func TestMLDSAExternalMu(t *testing.T) {
	pub := []byte("public key")
	msg := []byte("message")

	mu := mlDSAExternalMu(t, pub, nil, msg)
	require.Len(t, mu, 64)
	require.Equal(t, mu, mlDSAExternalMu(t, pub, []byte{}, msg), "nil and empty contexts are equivalent")
	require.NotEqual(t, mu, mlDSAExternalMu(t, pub, []byte("ctx"), msg), "mu must depend on the context")
	require.NotEqual(t, mu, mlDSAExternalMu(t, []byte("other key"), nil, msg), "mu must depend on the public key")
}

// TestSignVerifyMLDSA_Fake verifies the client sends ML-DSA parameters and parses the results using the generated fake server
func TestSignVerifyMLDSA_Fake(t *testing.T) {
	const keyName, keyVersion = "mldsa-key", "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
	kid := azkeys.ID("https://fakevault.vault.azure.net/keys/" + keyName + "/" + keyVersion)
	pub := []byte("ML-DSA-65 public key")
	msg := []byte("message")
	mlDSAContext := []byte("application-context-v1")
	signature := []byte("ML-DSA signature")
	mu := mlDSAExternalMu(t, pub, nil, msg)

	// The fake server runs these callbacks on its own goroutine, where a failing require would call
	// FailNow and leave the client call waiting forever. So the callbacks only record their arguments,
	// and the test goroutine asserts on them after each call returns.
	var (
		gotName, gotVersion string
		gotCreate           azkeys.CreateKeyParameters
		gotSign             azkeys.SignParameters
		gotVerify           azkeys.VerifyParameters
	)
	server := fake.Server{
		CreateKey: func(ctx context.Context, name string, parameters azkeys.CreateKeyParameters, options *azkeys.CreateKeyOptions) (resp azfake.Responder[azkeys.CreateKeyResponse], errResp azfake.ErrorResponder) {
			gotName, gotCreate = name, parameters
			resp.SetResponse(http.StatusOK, azkeys.CreateKeyResponse{
				KeyBundle: azkeys.KeyBundle{
					Key: &azkeys.JSONWebKey{
						Alg:    parameters.Alg,
						KeyOps: parameters.KeyOps,
						KID:    to.Ptr(kid),
						Kty:    to.Ptr(azkeys.KeyTypeAKPHSM),
						Pub:    pub,
					},
				},
			}, nil)
			return
		},
		Sign: func(ctx context.Context, name string, version string, parameters azkeys.SignParameters, options *azkeys.SignOptions) (resp azfake.Responder[azkeys.SignResponse], errResp azfake.ErrorResponder) {
			gotName, gotVersion, gotSign = name, version, parameters
			resp.SetResponse(http.StatusOK, azkeys.SignResponse{
				KeyOperationResult: azkeys.KeyOperationResult{KID: to.Ptr(kid), Result: signature},
			}, nil)
			return
		},
		Verify: func(ctx context.Context, name string, version string, parameters azkeys.VerifyParameters, options *azkeys.VerifyOptions) (resp azfake.Responder[azkeys.VerifyResponse], errResp azfake.ErrorResponder) {
			gotName, gotVersion, gotVerify = name, version, parameters
			resp.SetResponse(http.StatusOK, azkeys.VerifyResponse{
				KeyVerifyResult: azkeys.KeyVerifyResult{Value: to.Ptr(true)},
			}, nil)
			return
		},
	}

	client, err := azkeys.NewClient(vaultURL, &azcred.Fake{}, &azkeys.ClientOptions{
		ClientOptions: azcore.ClientOptions{Transport: fake.NewServerTransport(&server)},
	})
	require.NoError(t, err)

	createResp, err := client.CreateKey(context.Background(), keyName, azkeys.CreateKeyParameters{
		Alg:    to.Ptr(azkeys.AKPAlgorithmMLDSA65),
		KeyOps: to.SliceOfPtrs(azkeys.KeyOperationSign, azkeys.KeyOperationVerify),
		Kty:    to.Ptr(azkeys.KeyTypeAKP),
	}, nil)
	require.NoError(t, err)
	require.Equal(t, keyName, gotName)
	require.Equal(t, to.Ptr(azkeys.KeyTypeAKP), gotCreate.Kty)
	require.Equal(t, to.Ptr(azkeys.AKPAlgorithmMLDSA65), gotCreate.Alg)
	require.Equal(t, to.SliceOfPtrs(azkeys.KeyOperationSign, azkeys.KeyOperationVerify), gotCreate.KeyOps)
	require.Equal(t, azkeys.KeyTypeAKPHSM, *createResp.Key.Kty)
	require.Equal(t, azkeys.AKPAlgorithmMLDSA65, *createResp.Key.Alg)
	require.Equal(t, pub, createResp.Key.Pub)
	require.Equal(t, keyName, createResp.Key.KID.Name())
	require.Equal(t, keyVersion, createResp.Key.KID.Version())

	requireKey := func() {
		t.Helper()
		require.Equal(t, keyName, gotName)
		require.Equal(t, keyVersion, gotVersion)
	}

	signResp, err := client.Sign(context.Background(), keyName, keyVersion, azkeys.SignParameters{Context: mlDSAContext, Value: msg}, nil)
	require.NoError(t, err)
	requireKey()
	require.Nil(t, gotSign.Algorithm, "ML-DSA signing doesn't take an algorithm")
	require.Equal(t, msg, gotSign.Value)
	require.Equal(t, mlDSAContext, gotSign.Context)
	require.Nil(t, gotSign.ExternalMu)
	require.Equal(t, signature, signResp.Result)

	verifyResp, err := client.Verify(context.Background(), keyName, keyVersion, azkeys.VerifyParameters{Context: mlDSAContext, Digest: msg, Signature: signResp.Result}, nil)
	require.NoError(t, err)
	requireKey()
	require.Nil(t, gotVerify.Algorithm)
	require.Equal(t, signature, gotVerify.Signature)
	require.Equal(t, msg, gotVerify.Digest)
	require.Equal(t, mlDSAContext, gotVerify.Context)
	require.Nil(t, gotVerify.ExternalMu)
	require.True(t, *verifyResp.Value)

	signResp, err = client.Sign(context.Background(), keyName, keyVersion, azkeys.SignParameters{ExternalMu: mu}, nil)
	require.NoError(t, err)
	requireKey()
	require.Nil(t, gotSign.Algorithm, "ML-DSA signing doesn't take an algorithm")
	require.Equal(t, mu, gotSign.ExternalMu)
	require.Nil(t, gotSign.Value)
	require.Nil(t, gotSign.Context)
	require.Equal(t, signature, signResp.Result)

	verifyResp, err = client.Verify(context.Background(), keyName, keyVersion, azkeys.VerifyParameters{ExternalMu: mu, Signature: signResp.Result}, nil)
	require.NoError(t, err)
	requireKey()
	require.Nil(t, gotVerify.Algorithm)
	require.Equal(t, signature, gotVerify.Signature)
	require.Equal(t, mu, gotVerify.ExternalMu)
	require.Nil(t, gotVerify.Digest)
	require.True(t, *verifyResp.Value)
}
