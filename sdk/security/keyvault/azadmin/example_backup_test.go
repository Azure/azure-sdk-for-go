// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azadmin_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azadmin"
)

func ExampleNewBackupClient() {
	vaultURL := "https://<TODO: your vault name>.managedhsm.azure.net/"
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}

	client, err := azadmin.NewBackupClient(vaultURL, cred, nil)
	if err != nil {
		// TODO: handle error
	}

	_ = client
}
