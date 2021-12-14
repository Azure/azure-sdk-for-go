// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// 1. Set up session level sanitizers
	if recording.GetRecordMode() == "record" {
		err := recording.ResetProxy(nil)
		if err != nil {
			panic(err)
		}

		account, ok := os.LookupEnv(AccountNameEnvVar)
		if !ok {
			fmt.Printf("Could not find environment variable: %s", AccountNameEnvVar)
			os.Exit(1)
		}

		err = recording.AddURISanitizer("fakeaccount", account, nil)
		if err != nil {
			panic(err)
		}

	}

	// Run tests
	exitVal := m.Run()

	// 3. Reset
	// TODO: Add after sanitizer PR
	if recording.GetRecordMode() != "live" {
		err := recording.ResetProxy(nil)
		if err != nil {
			panic(err)
		}
	}

	// 4. Error out if applicable
	os.Exit(exitVal)
}

var pathToPackage = "sdk/storage/azblob/testdata"

func createServiceClientWithSharedKeyForRecording(t *testing.T, accountType testAccountType) (ServiceClient, error) {
	cred, err := getRecordingCredential(t, accountType)
	require.NoError(t, err)

	transporter, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{
		Transporter: transporter,
	}
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", cred.AccountName())
	return NewServiceClientWithSharedKey(serviceURL, cred, options)
}

func createServiceClientWithConnStrForRecording(t *testing.T, accountType testAccountType) (ServiceClient, error) {
	transporter, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	options := &ClientOptions{
		Transporter: transporter,
	}

	accountName, accountKey := getAccountNameKey(t, accountType)

	connectionString := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net/", accountName, accountKey)

	return NewServiceClientFromConnectionString(connectionString, options)
}

func getRecordingCredential(t *testing.T, accountType testAccountType) (*SharedKeyCredential, error) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		return NewSharedKeyCredential("fakeAccount", "daaaaaaaaaabbbbbbbbbbcccccccccccccccccccdddddddddddddddddddeeeeeeeeeeefffffffffffggggg==")
	}
	accountName, accountKey := getAccountNameKey(t, accountType)
	return NewSharedKeyCredential(accountName, accountKey)
}

func getAccountNameKey(t *testing.T, accountType testAccountType) (string, string) {
	var accountName string
	var accountKey string
	var ok bool

	if accountType == testAccountDefault {
		accountName, ok = os.LookupEnv("STORAGE_ACCOUNT_NAME")
		require.True(t, ok)
		accountKey, ok = os.LookupEnv("STORAGE_ACCOUNT_KEY")
		require.True(t, ok)
	} else if accountType == testAccountSecondary {
		accountName, ok = os.LookupEnv("SECONDARY_STORAGE_ACCOUNT_NAME")
		require.True(t, ok)
		accountKey, ok = os.LookupEnv("SECONDARY_STORAGE_ACCOUNT_KEY")
		require.True(t, ok)
	} else if accountType == testAccountPremium {
		accountName, ok = os.LookupEnv("PREMIUM_STORAGE_ACCOUNT_NAME")
		require.True(t, ok)
		accountKey, ok = os.LookupEnv("PREMIUM_STORAGE_ACCOUNT_KEY")
		require.True(t, ok)
	}

	return accountName, accountKey
}

func start(t *testing.T) func() {
	err := recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)
	return func() {
		err = recording.Stop(t, nil)
		require.NoError(t, err)
	}
}
