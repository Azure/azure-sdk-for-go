// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztemplate

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func TestOutput(t *testing.T) {
	cred, _ := azidentity.NewDefaultAzureCredential(nil)

	client, _ := NewClient(cred, nil)

	_, _ = client.SomeServiceAction(context.TODO(), nil)

	client.NewListValuesPager(nil)

	_, _ = client.BeginLongRunningOperation(context.TODO(), nil)
}
