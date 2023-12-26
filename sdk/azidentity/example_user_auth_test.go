//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity_test

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

// This example shows how to authenticate a user with [InteractiveBrowserCredential], enabling persistent
// caching so that the user doesn't need to authenticate interactively the next time the application runs.
func Example_userAuthentication() {
	cred, err := azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		// By default, credentials begin interactive authentication whenever necessary. To instead control when
		// a credential prompts for user interaction, set this option true. The credential will then return
		// azidentity.ErrAuthenticationRequired instead of prompting for authentication. The application
		// can then call the credential's Authenticate method when it's convenient to prompt the user.
		DisableAutomaticAuthentication: true,

		// By default, credentials cache in memory. Set TokenCachePersistenceOptions to enable persistent caching.
		TokenCachePersistenceOptions: &azidentity.TokenCachePersistenceOptions{
			// optionally set Name to isolate this credential's cache from other applications
			Name: "myapp",
		},
	})
	if err != nil {
		// TODO: handle error
	}

	// The Authenticate method begins interactive authentication. Call it whenever it's convenient for
	// your application to authenticate a user. If Authenticate succeeds, the credential is ready for
	// use with a client.
	record, err := cred.Authenticate(context.TODO(), nil)
	if err != nil {
		// TODO: handle error
	}

	// The record contains no authentication secrets. You can marshal it for storage.
	b, err := json.Marshal(record)
	if err != nil {
		// TODO: handle error
	}
	// TODO: store bytes
	_ = b

	// An authentication record stored by your application enables other credentials to access data from
	// past authentications. If the cache contains sufficient data, your application won't need to prompt
	// for authentication.
	var unmarshaled azidentity.AuthenticationRecord
	err = json.Unmarshal(b, &unmarshaled)
	if err != nil {
		// TODO: handle error
	}

	// this credential will be able to access authentication data cached by cred above, even in another process
	newCred, err := azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		AuthenticationRecord:           unmarshaled,
		DisableAutomaticAuthentication: true,
		TokenCachePersistenceOptions: &azidentity.TokenCachePersistenceOptions{
			Name: "myapp",
		},
	})
	if err != nil {
		// TODO: handle error
	}
	_ = newCred
}
