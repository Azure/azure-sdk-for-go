// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

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
		accountKey = recording.GetEnvVariable("PREMIUM_STORAGE_ACCOUNT_Key", "premfakestorageaccountkeykey")
	case testAccountSecondary:
		accountName = recording.GetEnvVariable("SECONDARY_STORAGE_ACCOUNT_NAME", "secondaryfakestorageaccount")
		accountKey = recording.GetEnvVariable("SECONDARY_STORAGE_ACCOUNT_Key", "secondaryfakestorageaccountkeykey")
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
