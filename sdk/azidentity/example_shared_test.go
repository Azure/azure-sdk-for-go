//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity_test

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// Helpers, variables, fakes to keep the examples tidy

const (
	authRecordPath = "fake/path"
	certPath       = "testdata/certificate.pem"
	clientID       = "fake-client-id"
	tenantID       = "fake-tenant"
)

func handleError(err error) {
	if err != nil {
		log.Panicf("example failed: %v", err)
	}
}

var cred azcore.TokenCredential
var err error
