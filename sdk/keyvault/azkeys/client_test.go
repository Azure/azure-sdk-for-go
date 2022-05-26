//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

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
			validateKey(t, to.Ptr(resp.Key))

			resp2, err := client.CreateRSAKey(ctx, key+"hsm", &CreateRSAKeyOptions{HardwareProtected: to.Ptr(true)})
			require.NoError(t, err)
			validateKey(t, to.Ptr(resp2.Key))

			cleanUpKey(t, client, key)
			cleanUpKey(t, client, key+"hsm")

			invalid, err := client.CreateRSAKey(ctx, "invalidName!@#$", nil)
			require.Error(t, err)
			require.Nil(t, invalid.Key.Properties)
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
		Tags: map[string]*string{
			"Tag1": to.Ptr("Val1"),
		},
	})
	defer cleanUpKey(t, client, key)
	require.NoError(t, err)
	validateKey(t, to.Ptr(resp.Key))
	require.Equal(t, 1, len(resp.Key.Properties.Tags))

	resp.Key.Properties.Tags = map[string]*string{}
	// Remove the tag
	resp2, err := client.UpdateKeyProperties(ctx, resp.Key, nil)
	require.NoError(t, err)
	require.Equal(t, 0, len(resp2.Properties.Tags))
	validateKey(t, &resp2.Key)
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
			validateKey(t, to.Ptr(resp.Key))

			invalid, err := client.CreateECKey(ctx, "key!@#$", nil)
			require.Error(t, err)
			require.Nil(t, invalid.Key.Properties)

			cleanUpKey(t, client, key)
		})
	}
}

func TestCreateOCTKey(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)

			resp, err := client.CreateOctKey(ctx, key, &CreateOctKeyOptions{
				Size:              to.Ptr(int32(256)),
				HardwareProtected: to.Ptr(true)},
			)

			if testType == REGULARTEST {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				validateKey(t, to.Ptr(resp.Key))

				cleanUpKey(t, client, key)
			}
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

				_, err = client.CreateKey(ctx, key, KeyTypeRSA, nil)
				require.NoError(t, err)
			}

			pager := client.NewListPropertiesOfKeysPager(nil)
			count := 0
			for pager.More() {
				resp, err := pager.NextPage(ctx)
				require.NoError(t, err)
				count += len(resp.Keys)
				for _, key := range resp.Keys {
					require.NotNil(t, key)
				}
			}
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

			r, err := client.CreateKey(ctx, key, KeyTypeRSA, nil)
			require.NoError(t, err)
			validateKey(t, to.Ptr(r.Key))

			resp, err := client.GetKey(ctx, key, nil)
			require.NoError(t, err)
			validateKey(t, to.Ptr(resp.Key))

			invalid, err := client.CreateKey(ctx, "invalidkey[]()", KeyTypeRSA, nil)
			require.Error(t, err)
			require.Nil(t, invalid.Key.Properties)
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

			r, err := client.CreateKey(ctx, key, KeyTypeRSA, nil)
			require.NoError(t, err)
			validateKey(t, to.Ptr(r.Key))

			poller, err := client.BeginDeleteKey(ctx, key, nil)
			require.NoError(t, err)
			deleteResp, err := poller.PollUntilDone(ctx, delay())
			require.NoError(t, err)
			require.NotNil(t, deleteResp.Key)
			require.NotNil(t, deleteResp.Key.ID)

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

			_, err = poller.Result(ctx)
			require.NoError(t, err)

			_, err = client.BeginDeleteKey(ctx, "nonexistent", nil)
			require.Error(t, err)
		})
	}
}

func TestBeginDeleteKeyRehydrate(t *testing.T) {
	stop := startTest(t)
	defer stop()

	client, err := createClient(t, testTypes[0])
	require.NoError(t, err)

	key, err := createRandomName(t, "rehydrate-poller")
	require.NoError(t, err)

	_, err = client.CreateRSAKey(ctx, key, nil)
	require.NoError(t, err)

	defer cleanUpKey(t, client, key)

	deletePoller, err := client.BeginDeleteKey(ctx, key, nil)
	require.NoError(t, err)

	rt, err := deletePoller.ResumeToken()
	require.NoError(t, err)

	rehydrated, err := client.BeginDeleteKey(ctx, key, &BeginDeleteKeyOptions{ResumeToken: rt})
	require.NoError(t, err)

	_, err = rehydrated.PollUntilDone(ctx, delay())
	require.NoError(t, err)

	// Validate key is not get-able
	_, err = client.GetKey(ctx, key, nil)
	require.Error(t, err)

	// Recover deleted
	recover, err := client.BeginRecoverDeletedKey(ctx, key, nil)
	require.NoError(t, err)

	rt, err = recover.ResumeToken()
	require.NoError(t, err)

	rehydratedRecover, err := client.BeginRecoverDeletedKey(ctx, key, &BeginRecoverDeletedKeyOptions{ResumeToken: rt})
	require.NoError(t, err)

	_, err = rehydratedRecover.PollUntilDone(ctx, delay())
	require.NoError(t, err)

	_, err = client.GetKey(ctx, key, nil)
	require.NoError(t, err)
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

			const retries = 5
			for i := 0; i < retries; i++ {
				// unfortunately purging a deleted key is non-deterministic so we
				// need to retry until we either succeed or hit the retry cap.
				restoreResp, err := client.RestoreKeyBackup(ctx, backupResp.Value, nil)
				if err != nil && i+1 == retries {
					t.Fatal("retry limit reached")
				} else if err != nil {
					if recording.GetRecordMode() != recording.PlaybackMode {
						time.Sleep(time.Minute)
					}
					continue
				}
				require.NoError(t, err)
				require.NotNil(t, restoreResp.Key)
				break
			}

			// Now the Key should be Get-able
			_, err = client.GetKey(ctx, key, nil)
			require.NoError(t, err)

			// confirm invalid response
			invalidResp, err := client.BackupKey(ctx, INVALIDKEYNAME, nil)
			require.Error(t, err)
			require.Equal(t, 0, len(invalidResp.Value))

			// confirm invalid restore key backup
			_, err = client.RestoreKeyBackup(ctx, []byte("doesnotexist"), nil)
			require.Error(t, err)
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

			_, err = client.GetDeletedKey(ctx, key, nil)
			require.NoError(t, err)

			resp, err := client.BeginRecoverDeletedKey(ctx, key, nil)
			require.NoError(t, err)

			_, err = resp.PollUntilDone(ctx, delay())
			require.NoError(t, err)

			getResp, err := client.GetKey(ctx, key, nil)
			require.NoError(t, err)
			require.NotNil(t, getResp.Key)

			_, err = client.BeginRecoverDeletedKey(ctx, "INVALIDKEYNAME", nil)
			require.Error(t, err)
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

			createResp, err := client.CreateRSAKey(ctx, key, &CreateRSAKeyOptions{})
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			createResp.Key.Properties.Tags = map[string]*string{
				"Tag1": to.Ptr("Val1"),
			}
			createResp.Key.Properties.ExpiresOn = to.Ptr(time.Now().AddDate(1, 0, 0))

			resp, err := client.UpdateKeyProperties(ctx, createResp.Key, nil)
			require.NoError(t, err)
			require.NotNil(t, resp.Properties)
			require.Equal(t, *resp.Properties.Tags["Tag1"], "Val1")
			require.NotNil(t, resp.Properties.ExpiresOn)

			createResp.Key.Properties.Name = to.Ptr("doesnotexist")
			invalid, err := client.UpdateKeyProperties(ctx, createResp.Key, nil)
			require.Error(t, err)
			require.Nil(t, invalid.Properties)
		})
	}
}

func TestUpdateKeyPropertiesImmutable(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			if testType == HSMTEST {
				t.Skip("HSM does not recognize immutable yet.")
			}
			stop := startTest(t)
			defer stop()
			err := recording.SetBodilessMatcher(t, nil)
			require.NoError(t, err)

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "immuta")
			require.NoError(t, err)

			marshalledPolicy, err := json.Marshal(map[string]interface{}{
				"anyOf": []map[string]interface{}{
					{
						"anyOf": []map[string]interface{}{
							{
								"claim":  "sdk-test",
								"equals": "true",
							}},
						"authority": os.Getenv("AZURE_KEYVAULT_ATTESTATION_URL"),
					},
				},
				"version": "1.0.0",
			})
			require.NoError(t, err)

			// retry creating the release policy because Key Vault sometimes can't reach
			// the fake attestation service we use in CI for several minutes after deployment
			var createResp CreateRSAKeyResponse
			for i := 0; i < 5; i++ {
				createResp, err = client.CreateRSAKey(ctx, key, &CreateRSAKeyOptions{
					HardwareProtected: to.Ptr(true),
					Properties: &Properties{
						Exportable: to.Ptr(true),
					},
					ReleasePolicy: &ReleasePolicy{
						Immutable:     to.Ptr(true),
						EncodedPolicy: marshalledPolicy,
					},
					Operations: []*Operation{to.Ptr(OperationEncrypt), to.Ptr(OperationDecrypt)},
				})
				if err == nil {
					break
				}
				if recording.GetRecordMode() != recording.PlaybackMode {
					time.Sleep(time.Minute)
				}
			}
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			newMarshalledPolicy, err := json.Marshal(map[string]interface{}{
				"anyOf": []map[string]interface{}{
					{
						"anyOf": []map[string]interface{}{
							{
								"claim":  "sdk-test",
								"equals": "false",
							}},
						"authority": os.Getenv("AZURE_KEYVAULT_ATTESTATION_URL"),
					},
				},
				"version": "1.0.0",
			})
			require.NoError(t, err)

			createResp.Key.ReleasePolicy = &ReleasePolicy{
				Immutable:     to.Ptr(true),
				EncodedPolicy: newMarshalledPolicy,
			}

			_, err = client.UpdateKeyProperties(ctx, createResp.Key, nil)
			require.Error(t, err)
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

			pager := client.NewListDeletedKeysPager(nil)
			count := 0
			for pager.More() {
				resp, err := pager.NextPage(ctx)
				require.NoError(t, err)
				count += len(resp.DeletedKeys)
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

			pager := client.NewListPropertiesOfKeyVersionsPager(key, nil)
			count := 0
			for pager.More() {
				resp, err := pager.NextPage(ctx)
				require.NoError(t, err)
				count += len(resp.Keys)
			}
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

			jwk := JSONWebKey{
				KeyType: to.Ptr(KeyTypeRSA),
				KeyOps:  to.SliceOfPtrs(OperationEncrypt, OperationDecrypt, OperationSign, OperationVerify, OperationWrapKey, OperationUnwrapKey),
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
			require.Nil(t, invalid.Properties)
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

			resp, err := client.GetRandomBytes(ctx, to.Ptr(int32(100)), nil)
			require.NoError(t, err)
			require.Equal(t, 100, len(resp.Value))

			_, err = client.GetRandomBytes(ctx, to.Ptr(int32(-1)), nil)
			require.Error(t, err)
		})
	}
}

func TestGetDeletedKey(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
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

			resp, err := client.GetDeletedKey(ctx, key, nil)
			require.NoError(t, err)
			require.Contains(t, *resp.Key.ID, key)

			_, err = client.PurgeDeletedKey(ctx, key, nil)
			require.NoError(t, err)
		})
	}
}

func TestRotateKey(t *testing.T) {
	for _, testType := range testTypes {
		t.Run(fmt.Sprintf("%s_%s", t.Name(), testType), func(t *testing.T) {
			skipHSM(t, testType)
			stop := startTest(t)
			defer stop()

			client, err := createClient(t, testType)
			require.NoError(t, err)

			key, err := createRandomName(t, "key")
			require.NoError(t, err)
			createResp, err := client.CreateRSAKey(ctx, key, nil)
			require.NoError(t, err)
			defer cleanUpKey(t, client, key)

			if testType == HSMTEST {
				// MHSM keys don't have a default rotation policy
				_, err = client.UpdateKeyRotationPolicy(ctx, key, RotationPolicy{Attributes: &RotationPolicyAttributes{ExpiresIn: to.Ptr("P30D")}}, nil)
				require.NoError(t, err)
			}
			resp, err := client.RotateKey(ctx, key, nil)
			require.NoError(t, err)

			require.NotEqual(t, *createResp.Key.JSONWebKey.ID, *resp.JSONWebKey.ID)
			require.NotEqual(t, createResp.Key.JSONWebKey.N, resp.JSONWebKey.N)

			invalid, err := client.RotateKey(ctx, "keynonexistent", nil)
			require.Error(t, err)
			require.Zero(t, invalid.Key)
		})
	}
}

func TestGetKeyRotationPolicy(t *testing.T) {
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

			_, err = client.GetKeyRotationPolicy(ctx, key, nil)
			require.NoError(t, err)
		})
	}
}

func TestReleaseKey(t *testing.T) {
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

			// Get attestation token from service
			attestationURL := recording.GetEnvVariable("AZURE_KEYVAULT_ATTESTATION_URL", "https://fakewebsite.net/")
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/generate-test-token", attestationURL), nil)
			require.NoError(t, err)

			if recording.GetRecordMode() == recording.PlaybackMode {
				t.Skip("Skipping test in playback")
			}

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, resp.StatusCode, http.StatusOK)
			defer resp.Body.Close()

			type targetResponse struct {
				Token string `json:"token"`
			}

			var tR targetResponse
			err = json.NewDecoder(resp.Body).Decode(&tR)
			require.NoError(t, err)

			_, err = client.ReleaseKey(ctx, key, "target", nil)
			require.Error(t, err)
		})
	}
}

func TestUpdateKeyRotationPolicy(t *testing.T) {
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

			get, err := client.GetKeyRotationPolicy(ctx, key, nil)
			require.NoError(t, err)
			get.Attributes = &RotationPolicyAttributes{ExpiresIn: to.Ptr("P90D")}
			get.LifetimeActions = []*LifetimeActions{
				{
					Action: &LifetimeActionsType{
						Type: to.Ptr(RotationActionRotate),
					},
					Trigger: &LifetimeActionsTrigger{
						TimeBeforeExpiry: to.Ptr("P30D"),
					},
				},
			}

			_, err = client.UpdateKeyRotationPolicy(ctx, key, get.RotationPolicy, nil)
			require.NoError(t, err)
		})
	}
}
