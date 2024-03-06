//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity_test

import (
	"context"
	"encoding/json"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	// importing the cache module registers the cache implementation for the current platform
	_ "github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
)

// this example shows file storage but any form of byte storage would work
func retrieveRecord() (azidentity.AuthenticationRecord, error) {
	record := azidentity.AuthenticationRecord{}
	b, err := os.ReadFile(authRecordPath)
	if err == nil {
		err = json.Unmarshal(b, &record)
	}
	return record, err
}

func storeRecord(record azidentity.AuthenticationRecord) error {
	b, err := json.Marshal(record)
	if err == nil {
		err = os.WriteFile(authRecordPath, b, 0600)
	}
	return err
}

// This example shows how to cache authentication data persistently so a user doesn't need to authenticate
// interactively every time the application runs. The example uses [InteractiveBrowserCredential], however
// [DeviceCodeCredential] has the same API. The key steps are:
//
//  1. Enable persistent caching by importing "github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache" and
//     setting [TokenCachePersistenceOptions]
//  2. Call Authenticate to acquire an [AuthenticationRecord] and store that for future use. An [AuthenticationRecord]
//     enables credentials to access data in the persistent cache. The record contains no authentication secrets.
//  3. Add the [AuthenticationRecord] to the credential's options
func Example_persistentUserAuthentication() {
	record, err := retrieveRecord()
	if err != nil {
		// TODO: handle error
	}
	cred, err := azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		AuthenticationRecord: record,
		// Credentials cache in memory by default. Set TokenCachePersistenceOptions to enable persistent caching.
		TokenCachePersistenceOptions: &azidentity.TokenCachePersistenceOptions{},
	})
	if err != nil {
		// TODO: handle error
	}

	if record == (azidentity.AuthenticationRecord{}) {
		// No stored record; call Authenticate to acquire one
		record, err = cred.Authenticate(context.TODO(), nil)
		if err != nil {
			// TODO: handle error
		}
		err = storeRecord(record)
		if err != nil {
			// TODO: handle error
		}
	}
}
