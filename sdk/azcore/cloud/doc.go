//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

/*
Package cloud implements a configuration API for applications deployed to sovereign or private Azure clouds.

Azure SDK client configuration defaults are appropriate for Azure Public Cloud (sometimes referred to as
"Azure Commercial"). This package enables applications deployed to other Azure Clouds to configure clients
appropriately.

Applications deployed to a sovereign cloud such as Azure US Government can use WellKnownClouds to configure
Azure SDK clients, for example:

	opts := azcore.ClientOptions{Cloud: cloud.AzureGovernment}
	cred, err := azidentity.NewDefaultAzureCredential(
		&azidentity.DefaultAzureCredentialOptions{ClientOptions: opts}
	)
	handle(err)

	client, err := armsubscription.NewClient(
		cred, &arm.ClientOptions{ClientOptions: opts}
	)
	handle(err)

Applications deployed to a private cloud such as Azure Stack create a Configuration object with
appropriate values:

	c = cloud.Configuration{
		LoginEndpoint: "https://...",
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {
				Audience: "...",
				Endpoint: "https://...",
			},
		},
	}
	opts := azcore.ClientOptions{Cloud: c}

	cred, err := azidentity.NewDefaultAzureCredential(
		&azidentity.DefaultAzureCredentialOptions{ClientOptions: opts}
	)
	handle(err)

	client, err := armsubscription.NewClient(
		cred, &arm.ClientOptions{ClientOptions: opts}
	)
	handle(err)
*/

package cloud
