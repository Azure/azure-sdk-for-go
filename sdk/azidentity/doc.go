//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

/*
azidentity provides Azure Active Directory token authentication for Azure SDK clients.

Azure SDK clients supporting token authentication can use any azidentity credential.
For example, authenticating a subscription client with DefaultAzureCredential:

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	...
	client := armsubscription.NewClient(cred, nil)
)

Different credential types implement different authentication flows and so have different
options. Each credential's documentation describes its particular options and how it
authenticates.

Credentials which authenticate service and user principals default to configuration
appropriate for Azure Public Cloud (sometimes called "Azure Commercial"). Applications
accessing resources in a private cloud,or a sovereign cloud such as Azure Government
can override this default. For example:

	import "github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"

	opts := azcore.ClientOptions{Cloud: cloud.AzureGovernment}
	cred, err := azidentity.NewDefaultAzureCredential(
		&azidentity.DefaultAzureCredentialOptions{ClientOptions: opts}
	)
*/

package azidentity
