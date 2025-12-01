// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
)

// This example shows how to cache service principal authentication data persistently to make it accessible
// to multiple processes. The example uses [ClientCertificateCredential], however the pattern is the same
// for all service principal credential types having a Cache field in their options. The key steps are:
//
//  1. Call [github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache.New] to construct a persistent cache
//  2. Set the Cache field in the credential's options
//
// Credentials that authenticate users such as [InteractiveBrowserCredential] have a different pattern; see
// the [persistent user authentication example].
//
// [persistent user authentication example]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#example-package-PersistentUserAuthentication
func Example_persistentServicePrincipalAuthentication() {
	c, err := cache.New(nil)
	if err != nil {
		// TODO: handle error. An error here means persistent
		// caching is impossible in the runtime environment.
	}
	certs, key, err := azidentity.ParseCertificates([]byte("cert data"), nil)
	if err != nil {
		// TODO: handle error
	}
	opts := &azidentity.ClientCertificateCredentialOptions{
		// Credentials cache in memory by default. Setting Cache with a
		// nonzero value from cache.New() enables persistent caching.
		Cache: c,
	}
	cred, err := azidentity.NewClientCertificateCredential("tenant ID", "client ID", certs, key, opts)
	if err != nil {
		// TODO: handle error
	}
	// TODO: pass the credential to an Azure SDK client constructor
	_ = cred
}
