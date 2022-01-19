//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var ctx = context.Background()

const HSMTEST = "HSM"
const REGULARTEST = "NON-HSM"
const INVALIDKEYNAME = "key!@#$%"

var testTypes = []string{REGULARTEST, HSMTEST}

func TestConstructor(t *testing.T) {
	client, err := NewClient("https://fakekvurl.vault.azure.net/", &FakeCredential{}, nil)
	require.NoError(t, err)
	require.NotNil(t, client.kvClient)
}

func TestCreateKeyRSA(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)

			resp, err := client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)
			require.NotNil(t, resp.Key)

			resp2, err := client.CreateRSAKey(ctx, key+"hsm", &CreateRSAKeyOptions{HardwareProtected: true})
			require.NoError(t, err)
			require.NotNil(t, resp2.Key)

			cleanUpKey(t, client, key)
			cleanUpKey(t, client, key+"hsm")

			invalid, err := client.CreateRSAKey(ctx, "invalidName!@#$", nil)
			require.Error(t, err)
			require.Nil(t, invalid.Attributes)
		})
	}
}
func TestCreateKeyRSATags(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t, REGULARTEST)
	require.NoError(t, err)

	key, err := createRandomName(t, "key")
	require.NoError(t, err)

	resp, err := client.CreateRSAKey(ctx, key, &CreateRSAKeyOptions{
		Tags: map[string]string{
			"Tag1": "Val1",
		},
	})
	defer cleanUpKey(t, client, key)
	require.NoError(t, err)
	require.NotNil(t, resp.Key)
	require.Equal(t, 1, len(resp.Tags))

	// Remove the tag
	resp2, err := client.UpdateKeyProperties(ctx, key, &UpdateKeyPropertiesOptions{
		Tags: map[string]string{},
	})
	require.NoError(t, err)
	require.Equal(t, 0, len(resp2.Tags))
}

func TestCreateECKey(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)

			resp, err := client.CreateECKey(ctx, key, nil)
			require.NoError(t, err)
			require.NotNil(t, resp.Key)

			invalid, err := client.CreateECKey(ctx, "key!@#$", nil)
			require.Error(t, err)
			require.Nil(t, invalid.Key)

			cleanUpKey(t, client, key)
		})
	}
}

func TestCreateOCTKey(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			if testType == REGULARTEST {
				t.Skip("OCT Key is HSM only")
			}
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)

			resp, err := client.CreateOCTKey(ctx, key, &CreateOCTKeyOptions{KeySize: to.Int32Ptr(256), HardwareProtected: true})
			require.NoError(t, err)
			require.NotNil(t, resp.Key)

			cleanUpKey(t, client, key)
		})
	}
}

func TestListKeys(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			for i := 0; i < 4; i++ {
				key, err := createRandomName(t, fmt.Sprintf("key-%d", i))
				require.NoError(t, err)

				_, err = client.CreateKey(ctx, key, RSA, nil)
				require.NoError(t, err)
			}

			pager := client.ListKeys(nil)
			count := 0
			for pager.NextPage(ctx) {
				count += len(pager.PageResponse().Keys)
				for _, key := range pager.PageResponse().Keys {
					require.NotNil(t, key)
				}
			}

			require.NoError(t, pager.Err())
			require.GreaterOrEqual(t, count, 4)

			for i := 0; i < 4; i++ {
				key, err := createRandomName(t, fmt.Sprintf("key-%d", i))
				require.NoError(t, err)
				cleanUpKey(t, client, key)
			}
		})
	}
}

func TestGetKey(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)

			_, err = client.CreateKey(ctx, key, RSA, nil)
			require.NoError(t, err)

			resp, err := client.GetKey(ctx, key, nil)
			require.NoError(t, err)
			require.NotNil(t, resp.Key)

			invalid, err := client.CreateKey(ctx, "invalidkey[]()", RSA, nil)
			require.Error(t, err)
			require.Nil(t, invalid.Attributes)
		})
	}
}

func TestDeleteKey(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			_, err = client.CreateKey(ctx, key, RSA, nil)
			require.NoError(t, err)

			resp, err := client.BeginDeleteKey(ctx, key, nil)
			require.NoError(t, err)
			_, err = resp.PollUntilDone(ctx, delay())
			require.Nil(t, err)

			_, err = client.GetKey(ctx, key, nil)
			require.Error(t, err)

			_, err = client.PurgeDeletedKey(ctx, key, nil)
			require.NoError(t, err)

			for i := 0; i < 5; i++ {
				_, err = client.GetDeletedKey(ctx, key, nil)
				if err != nil {
					break
				}
				require.NoError(t, err)
				recording.Sleep(time.Second * 2)
			}

			_, err = client.GetDeletedKey(ctx, key, nil)
			require.Error(t, err)

			_, err = resp.Poller.FinalResponse(ctx)
			require.NoError(t, err)

			invalidResp, err := client.BeginDeleteKey(ctx, "nonexistent", nil)
			require.Error(t, err)
			require.Nil(t, invalidResp.Poller)
		})
	}
}

func TestBackupKey(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "backup-key")
			require.NoError(t, err)

			_, err = client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)

			defer cleanUpKey(t, client, key)

			backupResp, err := client.BackupKey(ctx, key, nil)
			require.NoError(t, err)
			require.Greater(t, len(backupResp.Value), 0)

			respPoller, err := client.BeginDeleteKey(ctx, key, nil)
			require.NoError(t, err)
			_, err = respPoller.PollUntilDone(ctx, delay())
			require.NoError(t, err)

			_, err = client.PurgeDeletedKey(ctx, key, nil)
			require.NoError(t, err)

			_, err = client.GetKey(ctx, key, nil)
			var httpErr *azcore.ResponseError
			require.True(t, errors.As(err, &httpErr))
			require.Equal(t, httpErr.RawResponse.StatusCode, http.StatusNotFound)

			_, err = client.GetDeletedKey(ctx, key, nil)
			require.True(t, errors.As(err, &httpErr))
			require.Equal(t, httpErr.RawResponse.StatusCode, http.StatusNotFound)

			time.Sleep(30 * delay())
			// Poll this operation manually
			var restoreResp RestoreKeyBackupResponse
			var i int
			for i = 0; i < 10; i++ {
				restoreResp, err = client.RestoreKeyBackup(ctx, backupResp.Value, nil)
				if err == nil {
					break
				}
				time.Sleep(delay())
			}
			require.NoError(t, err)
			require.NotNil(t, restoreResp.Key)

			// Now the Key should be Get-able
			_, err = client.GetKey(ctx, key, nil)
			require.NoError(t, err)

			// confirm invalid response
			invalidResp, err := client.BackupKey(ctx, INVALIDKEYNAME, nil)
			require.Error(t, err)
			require.Equal(t, 0, len(invalidResp.Value))

			// confirm invalid restore key backup
			invalidResp2, err := client.RestoreKeyBackup(ctx, []byte("doesnotexist"), nil)
			require.Error(t, err)
			require.Nil(t, invalidResp2.RawResponse)
		})
	}
}

func TestRecoverDeletedKey(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)

			_, err = client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)

			defer cleanUpKey(t, client, key)

			pollerResp, err := client.BeginDeleteKey(ctx, key, nil)
			require.NoError(t, err)

			_, err = pollerResp.PollUntilDone(ctx, delay())
			require.NoError(t, err)

			resp, err := client.BeginRecoverDeletedKey(ctx, key, nil)
			require.NoError(t, err)

			_, err = resp.PollUntilDone(ctx, delay())
			require.NoError(t, err)

			getResp, err := client.GetKey(ctx, key, nil)
			require.NoError(t, err)
			require.NotNil(t, getResp.Key)

			invalidResp, err := client.BeginRecoverDeletedKey(ctx, "INVALIDKEYNAME", nil)
			require.Error(t, err)
			require.Nil(t, invalidResp.Poller)
		})
	}
}

func TestUpdateKeyProperties(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()
			err := recording.SetBodilessMatcher(t, nil)
			require.NoError(t, err)

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)

			_, err = client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			resp, err := client.UpdateKeyProperties(ctx, key, &UpdateKeyPropertiesOptions{
				Tags: map[string]string{
					"Tag1": "Val1",
				},
				KeyAttributes: &KeyAttributes{
					Attributes: Attributes{
						Expires: to.TimePtr(time.Now().AddDate(1, 0, 0)),
					},
				},
			})
			require.NoError(t, err)
			require.NotNil(t, resp.Attributes)
			require.Equal(t, resp.Tags["Tag1"], "Val1")
			require.NotNil(t, resp.Attributes.Updated)

			invalid, err := client.UpdateKeyProperties(ctx, "doesnotexist", nil)
			require.Error(t, err)
			require.Nil(t, invalid.Attributes)
		})
	}
}

func TestListDeletedKeys(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "list-del-key0")
			require.NoError(t, err)
			_, err = client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			pollerResp, err := client.BeginDeleteKey(ctx, key, nil)
			require.NoError(t, err)
			_, err = pollerResp.PollUntilDone(ctx, delay())
			require.NoError(t, err)

			key, err = createRandomName(t, "list-del-key1")
			require.NoError(t, err)
			_, err = client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			pollerResp, err = client.BeginDeleteKey(ctx, key, nil)
			require.NoError(t, err)
			_, err = pollerResp.PollUntilDone(ctx, delay())
			require.NoError(t, err)

			key, err = createRandomName(t, "list-del-key2")
			require.NoError(t, err)
			_, err = client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			pollerResp, err = client.BeginDeleteKey(ctx, key, nil)
			require.NoError(t, err)
			_, err = pollerResp.PollUntilDone(ctx, delay())
			require.NoError(t, err)

			pager := client.ListDeletedKeys(nil)
			count := 0
			for pager.NextPage(ctx) {
				count += len(pager.PageResponse().DeletedKeys)
			}

			require.GreaterOrEqual(t, count, 3)
		})
	}
}

func TestListKeyVersions(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)
			_, err = client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			for i := 0; i < 5; i++ {
				_, err = client.CreateRSAKey(ctx, key, nil)
				require.NoError(t, err)
			}

			pager := client.ListKeyVersions(key, nil)
			count := 0
			for pager.NextPage(ctx) {
				count += len(pager.PageResponse().Keys)
			}
			require.NoError(t, pager.Err())
			require.GreaterOrEqual(t, count, 6)
		})
	}
}

func TestImportKey(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			r := RSA
			jwk := JSONWebKey{
				KeyType: &r,
				KeyOps:  to.StringPtrArray("encrypt", "decrypt", "sign", "verify", "wrapKey", "unwrapKey"),
				N:       toBytes("00a0914d00234ac683b21b4c15d5bed887bdc959c2e57af54ae734e8f00720d775d275e455207e3784ceeb60a50a4655dd72a7a94d271e8ee8f7959a669ca6e775bf0e23badae991b4529d978528b4bd90521d32dd2656796ba82b6bbfc7668c8f5eeb5053747fd199319d29a8440d08f4412d527ff9311eda71825920b47b1c46b11ab3e91d7316407e89c7f340f7b85a34042ce51743b27d4718403d34c7b438af6181be05e4d11eb985d38253d7fe9bf53fc2f1b002d22d2d793fa79a504b6ab42d0492804d7071d727a06cf3a8893aa542b1503f832b296371b6707d4dc6e372f8fe67d8ded1c908fde45ce03bc086a71487fa75e43aa0e0679aa0d20efe35", t),
				E:       toBytes("10001", t),
				D:       toBytes("627c7d24668148fe2252c7fa649ea8a5a9ed44d75c766cda42b29b660e99404f0e862d4561a6c95af6a83d213e0a2244b03cd28576473215073785fb067f015da19084ade9f475e08b040a9a2c7ba00253bb8125508c9df140b75161d266be347a5e0f6900fe1d8bbf78ccc25eeb37e0c9d188d6e1fc15169ba4fe12276193d77790d2326928bd60d0d01d6ead8d6ac4861abadceec95358fd6689c50a1671a4a936d2376440a41445501da4e74bfb98f823bd19c45b94eb01d98fc0d2f284507f018ebd929b8180dbe6381fdd434bffb7800aaabdd973d55f9eaf9bb88a6ea7b28c2a80231e72de1ad244826d665582c2362761019de2e9f10cb8bcc2625649", t),
				P:       toBytes("00d1deac8d68ddd2c1fd52d5999655b2cf1565260de5269e43fd2a85f39280e1708ffff0682166cb6106ee5ea5e9ffd9f98d0becc9ff2cda2febc97259215ad84b9051e563e14a051dce438bc6541a24ac4f014cf9732d36ebfc1e61a00d82cbe412090f7793cfbd4b7605be133dfc3991f7e1bed5786f337de5036fc1e2df4cf3", t),
				Q:       toBytes("00c3dc66b641a9b73cd833bc439cd34fc6574465ab5b7e8a92d32595a224d56d911e74624225b48c15a670282a51c40d1dad4bc2e9a3c8dab0c76f10052dfb053bc6ed42c65288a8e8bace7a8881184323f94d7db17ea6dfba651218f931a93b8f738f3d8fd3f6ba218d35b96861a0f584b0ab88ddcf446b9815f4d287d83a3237", t),
				DP:      toBytes("00c9a159be7265cbbabc9afcc4967eb74fe58a4c4945431902d1142da599b760e03838f8cbd26b64324fea6bdc9338503f459793636e59b5361d1e6951e08ddb089e1b507be952a81fbeaf7e76890ea4f536e25505c3f648b1e88377dfc19b4c304e738dfca07211b792286a392a704d0f444c0a802539110b7f1f121c00cff0a9", t),
				DQ:      toBytes("00a0bd4c0a3d9f64436a082374b5caf2488bac1568696153a6a5e4cd85d186db31e2f58f024c617d29f37b4e6b54c97a1e25efec59c4d1fd3061ac33509ce8cae5c11f4cd2e83f41a8264f785e78dc0996076ee23dfdfc43d67c463afaa0180c4a718357f9a6f270d542479a0f213870e661fb950abca4a14ca290570ba7983347", t),
				QI:      toBytes("009fe7ae42e92bc04fcd5780464bd21d0c8ac0c599f9af020fde6ab0a7e7d1d39902f5d8fb6c614184c4c1b103fb46e94cd10a6c8a40f9991a1f28269f326435b6c50276fda6493353c650a833f724d80c7d522ba16c79f0eb61f672736b68fb8be3243d10943c4ab7028d09e76cfb5892222e38bc4d35585bf35a88cd68c73b07", t),
			}

			resp, err := client.ImportKey(ctx, "importedKey", jwk, nil)
			require.NoError(t, err)
			require.NotNil(t, resp.Key)

			invalid, err := client.ImportKey(ctx, "invalid", JSONWebKey{}, nil)
			require.Error(t, err)
			require.Nil(t, invalid.Attributes)
		})
	}
}

func TestGetRandomBytes(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			if testType == REGULARTEST {
				t.Skip("Managed HSM Only")
			}
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			resp, err := client.GetRandomBytes(ctx, to.Int32Ptr(100), nil)
			require.NoError(t, err)
			require.Equal(t, 100, len(resp.Value))

			invalid, err := client.GetRandomBytes(ctx, to.Int32Ptr(-1), nil)
			require.Error(t, err)
			require.Nil(t, invalid.RawResponse)
		})
	}
}

func TestGetDeletedKey(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			if testType == HSMTEST {
				t.Skip()
			}
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "keyName")
			require.NoError(t, err)
			_, err = client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			poller, err := client.BeginDeleteKey(ctx, key, nil)
			require.NoError(t, err)
			_, err = poller.PollUntilDone(ctx, delay())
			require.NoError(t, err)

			time.Sleep(10 * delay())

			resp, err := client.GetDeletedKey(ctx, key, nil)
			require.NoError(t, err)
			require.Contains(t, *resp.Key.ID, key)

			_, err = client.PurgeDeletedKey(ctx, key, nil)
			require.NoError(t, err)
		})
	}
}

func TestRotateKey(t *testing.T) {
	t.Skipf("Skipping while service disabled feature")
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			alwaysSkipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)
			createResp, err := client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			resp, err := client.RotateKey(ctx, key, nil)
			require.NoError(t, err)

			require.NotEqual(t, *createResp.Key.ID, *resp.Key.ID)
			require.NotEqual(t, createResp.Key.N, resp.Key.N)

			invalid, err := client.RotateKey(ctx, "keynonexistent", nil)
			require.Error(t, err)
			require.Nil(t, invalid.Key)
			require.Nil(t, invalid.Key)
		})
	}
}

func TestGetKeyRotationPolicy(t *testing.T) {
	t.Skipf("Skipping while service disabled feature")
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			alwaysSkipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)
			_, err = client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			_, err = client.GetKeyRotationPolicy(ctx, key, nil)
			require.NoError(t, err)
		})
	}
}

// This test is not ready, it will be ready in the 7.4 swagger, leaving this test for once that change is made.
func TestReleaseKey(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			alwaysSkipHSM(t, testType)
			// t.Skip("Release is not currently not enabled in API Version 7.3-preview")
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)
			_, err = client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			// Get attestation token from service
			attestationURL := recording.GetEnvVariable("AZURE_KEYVAULT_ATTESTATION_URL", "https://fakewebsite.net/")
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/generate-test-token", attestationURL), nil)
			require.NoError(t, err)

			if recording.GetRecordMode() == recording.PlaybackMode {
				t.Skip("Skipping test in playback")
			}

			// Issue when deploying HSM as well
			if _, ok := os.LookupEnv("AZURE_MANAGEDHSM_URL"); !ok {
				_, err = http.DefaultClient.Do(req)
				require.Error(t, err) // This URL doesn't exist so this should fail, will pass after 7.4-preview release
				// require.Equal(t, resp.StatusCode, http.StatusOK)
				// defer resp.Body.Close()

				// type targetResponse struct {
				// 	Token string `json:"token"`
				// }

				// var tR targetResponse
				// err = json.NewDecoder(resp.Body).Decode(&tR)
				// require.NoError(t, err)

				_, err = client.ReleaseKey(ctx, key, "target", nil)
				require.Error(t, err)
			}
		})
	}
}

func TestUpdateKeyRotationPolicy(t *testing.T) {
	t.Skipf("Skipping while service disabled feature")
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			alwaysSkipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)
			_, err = client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			_, err = client.UpdateKeyRotationPolicy(ctx, key, &UpdateKeyRotationPolicyOptions{
				Attributes: &KeyRotationPolicyAttributes{
					ExpiryTime: to.StringPtr("P90D"),
				},
				LifetimeActions: []*LifetimeActions{
					{
						Action: &LifetimeActionsType{
							Type: ActionTypeNotify.ToPtr(),
						},
						Trigger: &LifetimeActionsTrigger{
							TimeBeforeExpiry: to.StringPtr("P30D"),
						},
					},
				},
			})
			require.NoError(t, err)
		})
	}
}
