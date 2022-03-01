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
	switch recording.GetRecordMode() {
	case recording.PlaybackMode:
		err := recording.SetDefaultMatcher(nil, &recording.SetDefaultMatcherOptions{
			ExcludedHeaders: []string{"x-ms-tags", "x-ms-copy-source"},
		})
		if err != nil {
			panic(err)
		}
	case recording.RecordingMode:
		vals := [][]string{
			{"STORAGE_ACCOUNT_NAME", "fakestorageaccount"},
			{"PREMIUM_STORAGE_ACCOUNT_NAME", "premfakestorageaccount"},
			{"SECONDARY_STORAGE_ACCOUNT_NAME", "secondaryfakestorageaccount"},
		}
		for _, val := range vals {
			account, ok := os.LookupEnv(val[0])
			if !ok {
				fmt.Printf("Could not find environment variable: %s", val)
			} else {

				err := recording.AddGeneralRegexSanitizer(val[1], account, nil)
				if err != nil {
					panic(err)
				}
			}
		}

	}
	// Run tests
	exitVal := m.Run()

	// 3. Reset
	if recording.GetRecordMode() == "record" {
		err := recording.ResetProxy(nil)
		if err != nil {
			panic(err)
		}
	}

	// 4. Error out if applicable
	os.Exit(exitVal)
}

var pathToPackage = "sdk/storage/azblob/testdata"

func getCredential(accountType testAccountType) (*SharedKeyCredential, error) {
	var accountName string
	var accountKey string
	switch accountType {
	case testAccountDefault:
		accountName = recording.GetEnvVariable("STORAGE_ACCOUNT_NAME", "fakestorageaccount")
		accountKey = recording.GetEnvVariable("STORAGE_ACCOUNT_Key", "fakestorageaccountkeykey")
	case testAccountPremium:
		accountName = recording.GetEnvVariable("PREMIUM_STORAGE_ACCOUNT_NAME", "premfakestorageaccount")
		accountKey = recording.GetEnvVariable("PREMIUM_STORAGE_ACCOUNT_Key", "fakestorageaccountkeykey")
	case testAccountSecondary:
		accountName = recording.GetEnvVariable("SECONDARY_STORAGE_ACCOUNT_NAME", "secondaryfakestorageaccount")
		accountKey = recording.GetEnvVariable("SECONDARY_STORAGE_ACCOUNT_Key", "fakestorageaccountkeykey")
	default:
		return nil, fmt.Errorf("invalid test account type: %s", accountType)
	}
	return NewSharedKeyCredential(accountName, accountKey)
}

func start(t *testing.T) func() {
	err := recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)

	return func() {
		err = recording.Stop(t, nil)
		require.NoError(t, err)
	}
}

func createServiceClient(t *testing.T, accountType testAccountType) (ServiceClient, error) {
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	cred, err := getCredential(accountType)
	if err != nil {
		return ServiceClient{}, err
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", cred.AccountName())
	return NewServiceClientWithSharedKey(serviceURL, cred, &ClientOptions{
		Transporter: transport,
	})
}

func createServiceClientFromConnectionString(t *testing.T, accountType testAccountType) (ServiceClient, error) {
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)

	var accountName string
	var accountKey string
	fakeaccountkey := "fakestorageaccountkeykey"
	switch accountType {
	case testAccountDefault:
		accountName = recording.GetEnvVariable("STORAGE_ACCOUNT_NAME", "fakestorageaccount")
		accountKey = recording.GetEnvVariable("STORAGE_ACCOUNT_Key", fakeaccountkey)
	case testAccountPremium:
		accountName = recording.GetEnvVariable("PREMIUM_STORAGE_ACCOUNT_NAME", "premfakestorageaccount")
		accountKey = recording.GetEnvVariable("PREMIUM_STORAGE_ACCOUNT_Key", fakeaccountkey)
	case testAccountSecondary:
		accountName = recording.GetEnvVariable("SECONDARY_STORAGE_ACCOUNT_NAME", "secondaryfakestorageaccount")
		accountKey = recording.GetEnvVariable("SECONDARY_STORAGE_ACCOUNT_Key", fakeaccountkey)
	default:
		return ServiceClient{}, fmt.Errorf("invalid test account type: %s", accountType)
	}

	connectionString := fmt.Sprintf("DefaultEndpointsProtocol=https;AccountName=%s;AccountKey=%s;EndpointSuffix=core.windows.net/", accountName, accountKey)
	primaryURL, cred, err := parseConnectionString(connectionString)
	require.NoError(t, err)

	return NewServiceClientWithSharedKey(primaryURL, cred, &ClientOptions{
		Transporter: transport,
	})
}
