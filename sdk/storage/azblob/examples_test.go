// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azstorage"
)

const (
	endpoint = "https://<endpoint>.blob.core.windows.net/"
)

const (
	tenantID     = "<tenant>"
	clientID     = "<client>"
	clientSecret = "<secret>"
)

const (
	accountName = "<storageaccount>"
	accountKey  = "<accountkey>"
)

func clientSecretCredential() azcore.Credential {
	secret := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	// if err != nil {
	// 	panic(err)
	// }
	return secret
}

func sharedKeyCredential() azcore.Credential {
	sharedKey, err := azstorage.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		panic(err)
	}
	return sharedKey
}

func defaultCredential() azcore.Credential {
	cred, err := azidentity.NewDefaultTokenCredential(nil)
	if err != nil {
		panic(err)
	}
	return cred
}

func ExampleAnonymousCredential() {
	client, err := NewServiceClient(endpoint,
		// switch out with other credential functions above
		azcore.AnonymousCredential(),
		azcore.PipelineOptions{})
	if err != nil {
		panic(err)
	}
	iter := client.ListContainers(nil)
	for iter.NextItem(context.Background()) {
		fmt.Println(iter.Item().Name)
	}
	if iter.Err() != nil {
		panic(iter.Err())
	}
}
