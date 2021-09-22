// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"
)

func TestGetSDKInternalPolicies(t *testing.T) {
	cosmosClientOptions := &CosmosClientOptions{}
	policies := cosmosClientOptions.getSDKInternalPolicies()
	if policies == nil {
		t.Error("Expected policies to be not nil")
	}

	if len(policies) == 0 {
		t.Error("Expected policies to have more than 0 items ")
	}
}

func Test_newCosmosClientConnection(t *testing.T) {
	cred, _ := NewSharedKeyCredential("someKey")
	connection := newCosmosClientConnection("https://test.com", cred, &CosmosClientOptions{})
	if connection == nil {
		t.Error("Expected connection to be not nil")
	}
}
