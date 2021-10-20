//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztemplate

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func TestOutput(t *testing.T) {

	//options := &ClientOptions{}

	cred, _ := azidentity.NewDefaultAzureCredential(nil)

	client, _ := NewClient(cred, nil)

	client.SomeServiceAction()

}
