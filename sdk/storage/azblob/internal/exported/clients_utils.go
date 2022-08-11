//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ClientOptions contains the optional parameters when creating a Client.
// NOTE: all clients use this options type.
type ClientOptions struct {
	azcore.ClientOptions
}

func GetClientOptions(options *ClientOptions) *ClientOptions {
	if options == nil {
		options = &ClientOptions{}
	}
	return options
}
