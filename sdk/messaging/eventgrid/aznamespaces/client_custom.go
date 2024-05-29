//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznamespaces

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ClientOptions contains optional settings for [Client]
type ClientOptions struct {
	azcore.ClientOptions
}

var authScope = "https://eventgrid.azure.net/.default"
