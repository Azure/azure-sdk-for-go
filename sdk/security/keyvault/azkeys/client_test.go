// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stretchr/testify/require"
)

// pollStatus calls a function until it stops returning a response error with the given status code.
// If this takes more than 2 minutes, it fails the test.
func pollStatus(t *testing.T, expectedStatus int, fn func() error) {
	var err error
	for i := 0; i < 12; i++ {
		err = fn()
		var respErr *azcore.ResponseError
		if !(errors.As(err, &respErr) && respErr.StatusCode == expectedStatus) {
			break
		}
		if i < 11 {
			recording.Sleep(10 * time.Second)
		}
	}
	require.NoError(t, err)
}

func requireEqualAttributes(t *testing.T, a, b *azkeys.KeyAttributes) {
	if a == nil || b == nil {
		require.Equal(t, a, b)
		return
	}
	require.Equal(t, a.Created, b.Created)
	require.Equal(t, a.Enabled, b.Enabled)
	require.Equal(t, a.Expires, b.Expires)
	require.Equal(t, a.Exportable, b.Exportable)
	require.Equal(t, a.NotBefore, b.NotBefore)
	require.Equal(t, a.RecoverableDays, b.RecoverableDays)
	require.Equal(t, a.RecoveryLevel, b.RecoveryLevel)
	require.Equal(t, a.Updated, b.Updated)
}

func TestBackupRestore(t *testing.T) {
	name := "KV"
	for _, mhsm := range []bool{false, true} {
		if mhsm {
			name = "MHSM"
		}
		t.Run(name, func(t *testing.T) {
			client := startTest(t, mhsm)

			keyName := createRandomName(t, "testbackuprestore")
			createResp, err := client.CreateKey(context.Background(), keyName, azkeys.CreateKeyParameters{Kty: to.Ptr(azkeys.KeyTypeRSA)}, nil)
			require.NoError(t, err)
			require.Equal(t, keyName, createResp.Key.KID.Name())
			require.NotEmpty(t, createResp.Key.KID.Version())
			require.NotNil(t, createResp.Attributes)
			require.NotNil(t, createResp.Key)

			backupResp, err := client.BackupKey(context.Background(), keyName, nil)
			require.NoError(t, err)
			require.NotEmpty(t, backupResp.Value)

			deleteResp, err := client.DeleteKey(context.Background(), keyName, nil)
			require.NoError(t, err)
			require.Equal(t, createResp.Key.KID.Name(), deleteResp.Key.KID.Name())
			require.Equal(t, createResp.Key.KID.Version(), deleteResp.Key.KID.Version())
			requireEqualAttributes(t, createResp.Attributes, deleteResp.Attributes)
			require.NotNil(t, deleteResp.Key)
			require.NotEmpty(t, deleteResp.RecoveryID)
			require.NotEmpty(t, deleteResp.ScheduledPurgeDate)
			pollStatus(t, 404, func() error {
				_, err := client.GetDeletedKey(context.Background(), keyName, nil)
				return err
			})

			_, err = client.PurgeDeletedKey(context.Background(), keyName, nil)
			require.NoError(t, err)

			var restoreResp azkeys.RestoreKeyResponse
			restoreParams := azkeys.RestoreKeyParameters{KeyBackup: backupResp.Value}
			pollStatus(t, 409, func() error {
				restoreResp, err = client.RestoreKey(context.Background(), restoreParams, nil)
				return err
			})
			require.NoError(t, err)
			defer cleanUpKey(t, client, restoreResp.Key.KID)
			require.NotNil(t, restoreResp.Key)
			testSerde(t, &restoreParams)

			getResp, err := client.GetKey(context.Background(), keyName, "", nil)
			require.NoError(t, err)
			require.Equal(t, restoreResp.Attributes, getResp.Attributes)
			require.Equal(t, createResp.Key.KID.Name(), getResp.Key.KID.Name())
			require.Equal(t, createResp.Key.KID.Version(), getResp.Key.KID.Version())
		})
	}
}

func TestCRUD(t *testing.T) {
	attributes := &azkeys.KeyAttributes{
		Expires:   to.Ptr(time.Date(2050, 1, 1, 1, 1, 1, 0, time.UTC)),
		NotBefore: to.Ptr(time.Date(2040, 1, 1, 1, 1, 1, 0, time.UTC)),
	}
	tags := map[string]*string{"key": to.Ptr("value")}
	for _, mhsm := range []bool{false, true} {
		for _, params := range []azkeys.CreateKeyParameters{
			{
				Kty:           to.Ptr(azkeys.KeyTypeEC),
				Curve:         to.Ptr(azkeys.CurveNameP256K),
				KeyAttributes: attributes,
				Tags:          tags,
			},
			{
				Kty:            to.Ptr(azkeys.KeyTypeRSA),
				KeyAttributes:  attributes,
				KeySize:        to.Ptr(int32(2048)),
				PublicExponent: to.Ptr(int32(65537)),
				Tags:           tags,
			},
		} {
			testSerde(t, &params)
			name := string(*params.Kty)
			if mhsm {
				name += "_MHSM"
			}
			t.Run(name, func(t *testing.T) {
				client := startTest(t, mhsm)

				keyName := createRandomName(t, "testcrud")
				createResp, err := client.CreateKey(context.Background(), keyName, params, nil)
				require.NoError(t, err)
				require.Equal(t, keyName, createResp.Key.KID.Name())
				require.NotEmpty(t, createResp.Key.KID.Version())
				require.NotNil(t, createResp.Attributes)
				require.NotNil(t, createResp.Key)
				require.True(t, *createResp.Attributes.Enabled)

				getResp, err := client.GetKey(context.Background(), keyName, "", nil)
				require.NoError(t, err)
				requireEqualAttributes(t, createResp.Attributes, getResp.Attributes)
				require.Equal(t, createResp.Key.KID.Name(), getResp.Key.KID.Name())
				require.Equal(t, createResp.Key.KID.Version(), getResp.Key.KID.Version())
				testSerde(t, &getResp.KeyBundle)

				if mhsm {
					t.Skip("Skipping MHSM attestation until it supports 2025-07-01")

					getAttResp, err := client.GetKeyAttestation(context.Background(), keyName, "", nil)
					require.NoError(t, err)
					require.Equal(t, createResp.Key.KID.Name(), getAttResp.Key.KID.Name())
					require.Equal(t, createResp.Key.KID.Version(), getAttResp.Key.KID.Version())
					require.NotEmpty(t, getAttResp.Attributes.Attestation.CertificatePEMFile)
					require.NotEmpty(t, getAttResp.Attributes.Attestation.PrivateKeyAttestation)
					testSerde(t, &getAttResp.KeyBundle)
				}

				updateParams := azkeys.UpdateKeyParameters{
					KeyAttributes: &azkeys.KeyAttributes{
						Enabled: to.Ptr(false),
					},
				}
				testSerde(t, &updateParams)
				updateResp, err := client.UpdateKey(context.Background(), keyName, createResp.Key.KID.Version(), updateParams, nil)
				require.NoError(t, err)
				require.Equal(t, createResp.Key.KID.Name(), updateResp.Key.KID.Name())
				require.Equal(t, createResp.Key.KID.Version(), updateResp.Key.KID.Version())
				require.False(t, *updateResp.Attributes.Enabled)

				deleteResp, err := client.DeleteKey(context.Background(), keyName, nil)
				require.NoError(t, err)
				require.Equal(t, createResp.Key.KID.Name(), deleteResp.Key.KID.Name())
				require.Equal(t, createResp.Key.KID.Version(), deleteResp.Key.KID.Version())
				requireEqualAttributes(t, updateResp.Attributes, deleteResp.Attributes)
				testSerde(t, &deleteResp.DeletedKey)
				pollStatus(t, 404, func() error {
					_, err := client.GetDeletedKey(context.Background(), keyName, nil)
					return err
				})

				_, err = client.PurgeDeletedKey(context.Background(), keyName, nil)
				require.NoError(t, err)
			})
		}
	}
}

func TestDisableChallengeResourceVerification(t *testing.T) {
	authResource := `"Bearer authorization="https://login.microsoftonline.com/tenant", resource="%s""`
	authScope := `"Bearer authorization="https://login.microsoftonline.com/tenant", scope="%s""`
	vaultURL := "https://fakevault.vault.azure.net"
	for _, test := range []struct {
		challenge, resource string
		disableVerify, err  bool
	}{
		// happy path: resource matches requested vault's host (vault.azure.net)
		{challenge: authResource, resource: "https://vault.azure.net"},
		{challenge: authScope, resource: "https://vault.azure.net/.default"},
		{challenge: authResource, resource: "https://vault.azure.net", disableVerify: true},
		{challenge: authScope, resource: "https://vault.azure.net/.default", disableVerify: true},

		// error cases: resource/scope doesn't match the requested vault's host (vault.azure.net)
		{challenge: authResource, resource: "https://vault.azure.cn", err: true},
		{challenge: authResource, resource: "https://myvault.azure.net", err: true},
		{challenge: authScope, resource: "https://vault.azure.cn/.default", err: true},
		{challenge: authScope, resource: "https://myvault.azure.net/.default", err: true},

		// the policy shouldn't return errors for the above error cases when verification is disabled
		{challenge: authResource, resource: "https://vault.azure.cn", disableVerify: true},
		{challenge: authResource, resource: "https://myvault.azure.net", disableVerify: true},
		{challenge: authScope, resource: "https://vault.azure.cn/.default", disableVerify: true},
		{challenge: authScope, resource: "https://myvault.azure.net/.default", disableVerify: true},
	} {
		t.Run("", func(t *testing.T) {
			srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
			defer close()
			srv.AppendResponse(mock.WithStatusCode(401), mock.WithHeader("WWW-Authenticate", fmt.Sprintf(test.challenge, test.resource)))
			srv.AppendResponse(mock.WithStatusCode(200), mock.WithBody([]byte(`{"value":[]}`)))
			options := &azkeys.ClientOptions{
				ClientOptions: policy.ClientOptions{
					Transport: srv,
				},
				DisableChallengeResourceVerification: test.disableVerify,
			}
			client, err := azkeys.NewClient(vaultURL, &azcred.Fake{}, options)
			require.NoError(t, err)
			pager := client.NewListKeyPropertiesPager(nil)
			_, err = pager.NextPage(context.Background())
			if test.err {
				require.Error(t, err)
				require.Contains(t, err.Error(), "challenge resource")
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	for _, mhsm := range []bool{false, true} {
		name := "KV"
		if mhsm {
			name = "MHSM"
		}
		t.Run(name, func(t *testing.T) {
			client := startTest(t, mhsm)

			keyName := createRandomName(t, "key")
			createParams := azkeys.CreateKeyParameters{
				Kty:    to.Ptr(azkeys.KeyTypeRSAHSM),
				KeyOps: to.SliceOfPtrs(azkeys.KeyOperationEncrypt, azkeys.KeyOperationDecrypt),
			}
			createResp, err := client.CreateKey(context.Background(), keyName, createParams, nil)
			require.NoError(t, err)

			encryptParams := azkeys.KeyOperationParameters{
				Algorithm: to.Ptr(azkeys.EncryptionAlgorithmRSAOAEP256),
				Value:     []byte("plaintext"),
			}
			testSerde(t, &encryptParams)
			encryptResponse, err := client.Encrypt(context.Background(), keyName, createResp.Key.KID.Version(), encryptParams, nil)
			require.NoError(t, err)
			require.NotEmpty(t, encryptResponse.Result)
			testSerde(t, &encryptResponse.KeyOperationResult)

			decryptParams := azkeys.KeyOperationParameters{
				Algorithm: encryptParams.Algorithm,
				Value:     encryptResponse.Result,
			}
			testSerde(t, &decryptParams)
			decryptResponse, err := client.Decrypt(context.Background(), keyName, "", decryptParams, nil)
			require.NoError(t, err)
			require.Equal(t, decryptResponse.Result, encryptParams.Value)
			testSerde(t, &encryptResponse.KeyOperationResult)
		})
	}
}

func TestEncryptDecryptSymmetric(t *testing.T) {
	client := startTest(t, true)

	keyName := createRandomName(t, "key")
	createParams := azkeys.CreateKeyParameters{
		Kty:     to.Ptr(azkeys.KeyTypeOct),
		KeyOps:  to.SliceOfPtrs(azkeys.KeyOperationEncrypt, azkeys.KeyOperationDecrypt),
		KeySize: to.Ptr(int32(256)),
	}
	createResp, err := client.CreateKey(context.Background(), keyName, createParams, nil)
	require.NoError(t, err)

	encryptParams := azkeys.KeyOperationParameters{
		Algorithm: to.Ptr(azkeys.EncryptionAlgorithmA256CBCPAD),
		// IV must be random in real usage. This value is static only to ensure it matches in playback.
		IV:    []byte("0123456789ABCDEF"),
		Value: []byte("plaintext"),
	}
	testSerde(t, &encryptParams)
	encryptResponse, err := client.Encrypt(context.Background(), keyName, createResp.Key.KID.Version(), encryptParams, nil)
	require.NoError(t, err)
	require.NotEmpty(t, encryptResponse.Result)

	decryptParams := azkeys.KeyOperationParameters{
		Algorithm: encryptParams.Algorithm,
		IV:        encryptResponse.IV,
		Value:     encryptResponse.Result,
	}
	testSerde(t, &decryptParams)
	decryptResponse, err := client.Decrypt(context.Background(), keyName, "", decryptParams, nil)
	require.NoError(t, err)
	require.Equal(t, decryptResponse.Result, encryptParams.Value)
}

func TestGetRandomBytes(t *testing.T) {
	client := startTest(t, true)
	req := azkeys.GetRandomBytesParameters{Count: to.Ptr(int32(100))}
	testSerde(t, &req)
	resp, err := client.GetRandomBytes(context.Background(), req, nil)
	require.NoError(t, err)
	require.Equal(t, 100, len(resp.Value))
	testSerde(t, &resp)
}

func TestID(t *testing.T) {
	for _, test := range []struct{ ID, name, version string }{
		{"https://foo.vault.azure.net/keys/name/version", "name", "version"},
		{"https://foo.vault.azure.net/keys/name", "name", ""},
	} {
		t.Run(test.ID, func(t *testing.T) {
			ID := azkeys.ID(test.ID)
			require.Equal(t, test.name, ID.Name())
			require.Equal(t, test.version, ID.Version())
		})
	}
}

func TestImportKey(t *testing.T) {
	for _, mhsm := range []bool{false, true} {
		name := "KV"
		if mhsm {
			name = "MHSM"
		}
		t.Run(name, func(t *testing.T) {
			client := startTest(t, mhsm)
			jwk := &azkeys.JSONWebKey{
				KeyOps: to.SliceOfPtrs(azkeys.KeyOperationEncrypt),
				Kty:    to.Ptr(azkeys.KeyTypeRSA),
				N:      toBytes("a0914d00234ac683b21b4c15d5bed887bdc959c2e57af54ae734e8f00720d775d275e455207e3784ceeb60a50a4655dd72a7a94d271e8ee8f7959a669ca6e775bf0e23badae991b4529d978528b4bd90521d32dd2656796ba82b6bbfc7668c8f5eeb5053747fd199319d29a8440d08f4412d527ff9311eda71825920b47b1c46b11ab3e91d7316407e89c7f340f7b85a34042ce51743b27d4718403d34c7b438af6181be05e4d11eb985d38253d7fe9bf53fc2f1b002d22d2d793fa79a504b6ab42d0492804d7071d727a06cf3a8893aa542b1503f832b296371b6707d4dc6e372f8fe67d8ded1c908fde45ce03bc086a71487fa75e43aa0e0679aa0d20efe35", t),
				E:      toBytes("10001", t),
				D:      toBytes("627c7d24668148fe2252c7fa649ea8a5a9ed44d75c766cda42b29b660e99404f0e862d4561a6c95af6a83d213e0a2244b03cd28576473215073785fb067f015da19084ade9f475e08b040a9a2c7ba00253bb8125508c9df140b75161d266be347a5e0f6900fe1d8bbf78ccc25eeb37e0c9d188d6e1fc15169ba4fe12276193d77790d2326928bd60d0d01d6ead8d6ac4861abadceec95358fd6689c50a1671a4a936d2376440a41445501da4e74bfb98f823bd19c45b94eb01d98fc0d2f284507f018ebd929b8180dbe6381fdd434bffb7800aaabdd973d55f9eaf9bb88a6ea7b28c2a80231e72de1ad244826d665582c2362761019de2e9f10cb8bcc2625649", t),
				P:      toBytes("00d1deac8d68ddd2c1fd52d5999655b2cf1565260de5269e43fd2a85f39280e1708ffff0682166cb6106ee5ea5e9ffd9f98d0becc9ff2cda2febc97259215ad84b9051e563e14a051dce438bc6541a24ac4f014cf9732d36ebfc1e61a00d82cbe412090f7793cfbd4b7605be133dfc3991f7e1bed5786f337de5036fc1e2df4cf3", t),
				Q:      toBytes("00c3dc66b641a9b73cd833bc439cd34fc6574465ab5b7e8a92d32595a224d56d911e74624225b48c15a670282a51c40d1dad4bc2e9a3c8dab0c76f10052dfb053bc6ed42c65288a8e8bace7a8881184323f94d7db17ea6dfba651218f931a93b8f738f3d8fd3f6ba218d35b96861a0f584b0ab88ddcf446b9815f4d287d83a3237", t),
				DP:     toBytes("00c9a159be7265cbbabc9afcc4967eb74fe58a4c4945431902d1142da599b760e03838f8cbd26b64324fea6bdc9338503f459793636e59b5361d1e6951e08ddb089e1b507be952a81fbeaf7e76890ea4f536e25505c3f648b1e88377dfc19b4c304e738dfca07211b792286a392a704d0f444c0a802539110b7f1f121c00cff0a9", t),
				DQ:     toBytes("00a0bd4c0a3d9f64436a082374b5caf2488bac1568696153a6a5e4cd85d186db31e2f58f024c617d29f37b4e6b54c97a1e25efec59c4d1fd3061ac33509ce8cae5c11f4cd2e83f41a8264f785e78dc0996076ee23dfdfc43d67c463afaa0180c4a718357f9a6f270d542479a0f213870e661fb950abca4a14ca290570ba7983347", t),
				QI:     toBytes("009fe7ae42e92bc04fcd5780464bd21d0c8ac0c599f9af020fde6ab0a7e7d1d39902f5d8fb6c614184c4c1b103fb46e94cd10a6c8a40f9991a1f28269f326435b6c50276fda6493353c650a833f724d80c7d522ba16c79f0eb61f672736b68fb8be3243d10943c4ab7028d09e76cfb5892222e38bc4d35585bf35a88cd68c73b07", t),
			}
			params := azkeys.ImportKeyParameters{HSM: to.Ptr(true), Key: jwk}
			testSerde(t, &params)
			resp, err := client.ImportKey(context.Background(), createRandomName(t, "testimport"), params, nil)
			require.NoError(t, err)
			defer cleanUpKey(t, client, resp.Key.KID)
			require.Equal(t, jwk.KeyOps, resp.Key.KeyOps)
			require.Equal(t, jwk.N, resp.Key.N)
			require.Equal(t, jwk.E, resp.Key.E)
		})
	}
}

func TestListDeletedKeys(t *testing.T) {
	for _, mhsm := range []bool{false, true} {
		name := "KV"
		if mhsm {
			name = "MHSM"
		}
		t.Run(name, func(t *testing.T) {
			client := startTest(t, mhsm)
			count := 4
			keyNames := make([]string, count)
			createParams := azkeys.CreateKeyParameters{
				Kty:  to.Ptr(azkeys.KeyTypeRSA),
				Tags: map[string]*string{"count-this-key": to.Ptr("yes")},
			}
			for i := 0; i < len(keyNames); i++ {
				n := createRandomName(t, fmt.Sprintf("listdeletedkeys%d", i))
				keyNames[i] = n
				createResp, err := client.CreateKey(context.Background(), n, createParams, nil)
				require.NoError(t, err)
				cleanUpKey(t, client, createResp.Key.KID)
			}
			for i := 0; i < len(keyNames); i++ {
				pollStatus(t, 404, func() error {
					_, err := client.GetDeletedKey(context.Background(), keyNames[i], nil)
					return err
				})
			}
			pager := client.NewListDeletedKeyPropertiesPager(nil)
			for pager.More() {
				resp, err := pager.NextPage(context.Background())
				require.NoError(t, err)
				testSerde(t, &resp.DeletedKeyPropertiesListResult)
				for _, key := range resp.Value {
					require.NotEmpty(t, key.Attributes)
					require.NotNil(t, key.DeletedDate)
					require.NotEmpty(t, key.KID.Name())
					require.NotNil(t, key.RecoveryID)
					require.NotNil(t, key.ScheduledPurgeDate)
					if strings.HasPrefix(key.KID.Name(), "listdeletedkeys") {
						require.NotEmpty(t, key.Tags)
						if *key.Tags["count-this-key"] == "yes" {
							count--
						}
						testSerde(t, key)
					}
				}
			}
			require.Equal(t, count, 0)
		})
	}
}

func TestListKeys(t *testing.T) {
	name := "KV"
	for _, mhsm := range []bool{false, true} {
		if mhsm {
			name = "MHSM"
		}
		t.Run(name, func(t *testing.T) {
			client := startTest(t, mhsm)
			count := 0
			keyNamePrefix := "testlistkeys"
			for i := 0; i < 4; i++ {
				n := createRandomName(t, fmt.Sprintf("%s-%d", keyNamePrefix, i))
				resp, err := client.CreateKey(context.Background(), n, azkeys.CreateKeyParameters{Kty: to.Ptr(azkeys.KeyTypeRSA)}, nil)
				require.NoError(t, err)
				defer cleanUpKey(t, client, resp.Key.KID)
				count++
			}

			pager := client.NewListKeyPropertiesPager(nil)
			for pager.More() {
				resp, err := pager.NextPage(context.Background())
				require.NoError(t, err)
				testSerde(t, &resp.KeyPropertiesListResult)
				for _, key := range resp.Value {
					require.NotNil(t, key)
					require.NotNil(t, key.Attributes)
					require.NotNil(t, key.KID)
					if strings.HasPrefix(key.KID.Name(), keyNamePrefix) {
						count--
					}
					testSerde(t, key)
				}
			}
			require.Equal(t, count, 0)
		})
	}
}

func TestListKeyVersions(t *testing.T) {
	for _, mhsm := range []bool{false, true} {
		name := "KV"
		if mhsm {
			name = "MHSM"
		}
		t.Run(name, func(t *testing.T) {
			client := startTest(t, mhsm)

			var createResp azkeys.CreateKeyResponse
			var err error
			keyName := createRandomName(t, "listkeyversions")
			expectedVersions := make(map[string]struct{}, 4)
			for i := 0; i < 4; i++ {
				createResp, err = client.CreateKey(context.Background(), keyName, azkeys.CreateKeyParameters{Kty: to.Ptr(azkeys.KeyTypeRSA)}, nil)
				expectedVersions[createResp.Key.KID.Version()] = struct{}{}
				require.NoError(t, err)
			}
			defer cleanUpKey(t, client, createResp.Key.KID)

			pager := client.NewListKeyPropertiesVersionsPager(keyName, nil)
			for pager.More() {
				resp, err := pager.NextPage(context.Background())
				require.NoError(t, err)
				testSerde(t, &resp.KeyPropertiesListResult)
				for _, key := range resp.Value {
					testSerde(t, key)
					require.NotNil(t, key)
					require.NotNil(t, key.Attributes)
					require.NotNil(t, key.KID)
					require.Equal(t, keyName, key.KID.Name())
					version := key.KID.Version()
					require.NotEmpty(t, keyName, version)
					require.Contains(t, expectedVersions, version)
					delete(expectedVersions, version)
				}
			}
			require.Empty(t, expectedVersions)
		})
	}
}

func TestRecoverDeletedKey(t *testing.T) {
	for _, mhsm := range []bool{false, true} {
		name := "KV"
		if mhsm {
			name = "MHSM"
		}
		t.Run(name, func(t *testing.T) {
			client := startTest(t, mhsm)

			key := createRandomName(t, "key")
			createResp, err := client.CreateKey(context.Background(), key, azkeys.CreateKeyParameters{Kty: to.Ptr(azkeys.KeyTypeEC)}, nil)
			require.NoError(t, err)

			_, err = client.DeleteKey(context.Background(), key, nil)
			require.NoError(t, err)
			pollStatus(t, 404, func() error {
				_, err := client.GetDeletedKey(context.Background(), key, nil)
				return err
			})

			recoverResp, err := client.RecoverDeletedKey(context.Background(), key, nil)
			require.NoError(t, err)
			pollStatus(t, 404, func() error {
				_, err := client.GetKey(context.Background(), key, createResp.Key.KID.Version(), nil)
				return err
			})
			cleanUpKey(t, client, createResp.Key.KID)
			require.Equal(t, createResp.Key.KID, recoverResp.Key.KID)
			require.NotNil(t, recoverResp.Attributes)
		})
	}
}

func TestReleaseKey(t *testing.T) {
	for _, mhsm := range []bool{false, true} {
		name := "KV"
		if mhsm {
			name = "MHSM"
		}
		t.Run(name, func(t *testing.T) {
			client := startTest(t, false)
			key := createRandomName(t, "testreleasekey")

			// retry creating the key because Key Vault sometimes can't reach the fake
			// attestation service we use in CI for several minutes after deployment
			var createResp azkeys.CreateKeyResponse
			var err error
			for i := 0; i < 5; i++ {
				params := azkeys.CreateKeyParameters{
					Curve: to.Ptr(azkeys.CurveNameP256K),
					KeyAttributes: &azkeys.KeyAttributes{
						Exportable: to.Ptr(true),
					},
					Kty: to.Ptr(azkeys.KeyTypeECHSM),
					ReleasePolicy: &azkeys.KeyReleasePolicy{
						EncodedPolicy: getMarshalledReleasePolicy(attestationURL),
						Immutable:     to.Ptr(true),
					},
				}
				createResp, err = client.CreateKey(context.Background(), key, params, nil)
				if err == nil {
					break
				}
				if i < 4 {
					recording.Sleep(30 * time.Second)
				}
			}
			require.NoError(t, err)
			require.NotNil(t, createResp.Key.KID)
			defer cleanUpKey(t, client, createResp.Key.KID)

			attestationClient, err := recording.NewRecordingHTTPClient(t, nil)
			require.NoError(t, err)
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/generate-test-token", attestationURL), nil)
			require.NoError(t, err)
			resp, err := attestationClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			defer resp.Body.Close()

			var tR struct {
				Token *string `json:"token"`
			}
			err = json.NewDecoder(resp.Body).Decode(&tR)
			require.NoError(t, err)

			params := azkeys.ReleaseParameters{TargetAttestationToken: tR.Token}
			testSerde(t, &params)
			releaseResp, err := client.Release(context.Background(), key, "", params, nil)
			if err != nil && strings.Contains(err.Error(), "Target environment attestation statement cannot be verified.") {
				t.Skip("test encountered a transient service fault; see https://github.com/Azure/azure-sdk-for-net/issues/27957")
			}
			require.NoError(t, err)
			require.NotEmpty(t, releaseResp.KeyReleaseResult.Value)
			testSerde(t, &releaseResp.KeyReleaseResult)
		})
	}

}

func TestRotateKey(t *testing.T) {
	for _, mhsm := range []bool{false, true} {
		name := "KV"
		if mhsm {
			name = "MHSM"
		}
		t.Run(name, func(t *testing.T) {
			client := startTest(t, mhsm)
			key := createRandomName(t, "testrotatekey")
			createResp, err := client.CreateKey(context.Background(), key, azkeys.CreateKeyParameters{Kty: to.Ptr(azkeys.KeyTypeECHSM)}, nil)
			require.NoError(t, err)
			defer cleanUpKey(t, client, createResp.Key.KID)

			timeAfterCreate := to.Ptr("P30D")
			policy := azkeys.KeyRotationPolicy{
				Attributes: &azkeys.KeyRotationPolicyAttributes{
					ExpiryTime: to.Ptr("P90D"),
				},
				LifetimeActions: []*azkeys.LifetimeAction{
					{
						Action: &azkeys.LifetimeActionType{
							Type: to.Ptr(azkeys.KeyRotationPolicyActionRotate),
						},
						Trigger: &azkeys.LifetimeActionTrigger{
							TimeAfterCreate: timeAfterCreate,
						},
					},
				}}
			updateResp, err := client.UpdateKeyRotationPolicy(context.Background(), key, policy, nil)
			require.NoError(t, err)
			require.Equal(t, policy.Attributes.ExpiryTime, updateResp.Attributes.ExpiryTime)
			require.NotEmpty(t, updateResp.LifetimeActions)

			getResp, err := client.GetKeyRotationPolicy(context.Background(), key, nil)
			require.NoError(t, err)
			require.Equal(t, updateResp.Attributes.ExpiryTime, getResp.Attributes.ExpiryTime)
			require.NotEmpty(t, getResp.LifetimeActions)
			require.Condition(t, func() bool {
				for _, action := range getResp.LifetimeActions {
					if strings.EqualFold(string(*action.Action.Type), string(azkeys.KeyRotationPolicyActionRotate)) && strings.EqualFold(string(*action.Trigger.TimeAfterCreate), *timeAfterCreate) {
						return true
					}
				}
				return false
			}, "GetKeyRotationPolicy returned a policy missing the updated action")

			rotateResp, err := client.RotateKey(context.Background(), key, nil)
			require.NoError(t, err)
			require.NotNil(t, rotateResp.Key.KID)

			invalid, err := client.RotateKey(context.Background(), "keynonexistent", nil)
			require.Error(t, err)
			require.Zero(t, invalid.Key)
		})
	}
}

func TestSignVerify(t *testing.T) {
	for _, mhsm := range []bool{false, true} {
		name := "KV"
		if mhsm {
			name = "MHSM"
		}
		t.Run(name, func(t *testing.T) {
			client := startTest(t, mhsm)

			keyName := createRandomName(t, "key")

			createParams := azkeys.CreateKeyParameters{
				Curve:  to.Ptr(azkeys.CurveNameP256K),
				KeyOps: to.SliceOfPtrs(azkeys.KeyOperationSign, azkeys.KeyOperationVerify),
				Kty:    to.Ptr(azkeys.KeyTypeEC),
			}
			_, err := client.CreateKey(context.Background(), keyName, createParams, nil)
			require.NoError(t, err)

			hasher := sha256.New()
			_, err = hasher.Write([]byte("plaintext"))
			require.NoError(t, err)
			digest := hasher.Sum(nil)

			signParams := azkeys.SignParameters{Algorithm: to.Ptr(azkeys.SignatureAlgorithmES256K), Value: digest}
			testSerde(t, &signParams)
			signResponse, err := client.Sign(context.Background(), keyName, "", signParams, nil)
			require.NoError(t, err)
			testSerde(t, &signResponse.KeyOperationResult)

			verifyParams := azkeys.VerifyParameters{Algorithm: signParams.Algorithm, Digest: digest, Signature: signResponse.Result}
			testSerde(t, &verifyParams)
			verifyResponse, err := client.Verify(context.Background(), keyName, "", verifyParams, nil)
			require.NoError(t, err)
			require.True(t, *verifyResponse.Value)
			testSerde(t, &verifyResponse.KeyVerifyResult)
		})
	}
}

func TestWrapUnwrap(t *testing.T) {
	for _, mhsm := range []bool{false, true} {
		name := "KV"
		if mhsm {
			name = "MHSM"
		}
		t.Run(name, func(t *testing.T) {
			client := startTest(t, mhsm)

			keyName := createRandomName(t, "key")

			createParams := azkeys.CreateKeyParameters{
				KeyOps: to.SliceOfPtrs(azkeys.KeyOperationWrapKey, azkeys.KeyOperationUnwrapKey),
				Kty:    to.Ptr(azkeys.KeyTypeRSA),
			}
			_, err := client.CreateKey(context.Background(), keyName, createParams, nil)
			require.NoError(t, err)

			keyBytes := []byte("5063e6aaa845f150200547944fd199679c98ed6f99da0a0b2dafeaf1f4684496fd532c1c229968cb9dee44957fcef7ccef59ceda0b362e56bcd78fd3faee5781c623c0bb22b35beabde0664fd30e0e824aba3dd1b0afffc4a3d955ede20cf6a854d52cfd")

			wrapParams := azkeys.KeyOperationParameters{Algorithm: to.Ptr(azkeys.EncryptionAlgorithmRSAOAEP), Value: keyBytes}
			wrapResp, err := client.WrapKey(context.Background(), keyName, "", wrapParams, nil)
			require.NoError(t, err)

			unwrapResp, err := client.UnwrapKey(context.Background(), keyName, "", azkeys.KeyOperationParameters{Algorithm: wrapParams.Algorithm, Value: wrapResp.Result}, nil)
			require.NoError(t, err)
			require.Equal(t, keyBytes, unwrapResp.Result)
		})
	}
}

func TestAPIVersion(t *testing.T) {
	apiVersion := "7.3"
	var requireVersion = func(req *http.Request) bool {
		version := req.URL.Query().Get("api-version")
		require.Equal(t, version, apiVersion)
		return true
	}
	srv, close := mock.NewServer(mock.WithTransformAllRequestsToTestServerUrl())
	defer close()
	srv.AppendResponse(
		mock.WithStatusCode(200),
		mock.WithPredicate(requireVersion),
	)
	srv.AppendResponse(mock.WithStatusCode(http.StatusInternalServerError))

	opts := &azkeys.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport:  srv,
			APIVersion: apiVersion,
		},
	}
	client, err := azkeys.NewClient(vaultURL, &azcred.Fake{}, opts)
	require.NoError(t, err)

	_, err = client.GetKey(context.Background(), "name", "", nil)
	require.NoError(t, err)
}
