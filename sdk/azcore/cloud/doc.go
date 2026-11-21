// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

/*
Package cloud implements a configuration API for applications deployed to sovereign or private Azure clouds.

Azure SDK client configuration defaults are appropriate for Azure Public Cloud (sometimes referred to as
"Azure Commercial" or simply "Microsoft Azure"). This package enables applications deployed to other
Azure Clouds to configure clients appropriately.

This package contains predefined configuration for well-known sovereign clouds such as Azure Government and
Azure China. Azure SDK clients accept this configuration via the Cloud field of azcore.ClientOptions. For
example, configuring a credential and ARM client for Azure Government:

	opts := azcore.ClientOptions{Cloud: cloud.AzureGovernment}
	cred, err := azidentity.NewDefaultAzureCredential(
		&azidentity.DefaultAzureCredentialOptions{ClientOptions: opts},
	)
	handle(err)

	client, err := armsubscription.NewClient(
		cred, &arm.ClientOptions{ClientOptions: opts},
	)
	handle(err)

Azure Resource Manager clients whose options don't specify a cloud use the well-known cloud named by the
environment variable AZURE_ENVIRONMENT, when it's set. Valid values are "AzureChina", "AzureGovernment",
and "AzurePublic". Values are case insensitive, surrounding whitespace is ignored, and the names
earlier Azure SDKs and the Azure CLI use for these clouds, such as "AzureUSGovernmentCloud", are also
accepted. Constructors of these clients return an error when the variable has any other value.
Configuration set in code always takes precedence over the environment variable. Clients for other
services don't read the variable; continue to configure them in code.

Azure Resource Manager clients don't read AZURE_AUTHORITY_HOST, which configures only azidentity
credentials. An application that sets AZURE_AUTHORITY_HOST for a sovereign cloud must also set
AZURE_ENVIRONMENT, or specify the cloud in code, to direct its Resource Manager clients to that cloud.

Applications deployed to a private cloud such as Azure Stack create a Configuration object with
appropriate values:

	c := cloud.Configuration{
		ActiveDirectoryAuthorityHost: "https://...",
		Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
			cloud.ResourceManager: {
				Audience: "...",
				Endpoint: "https://...",
			},
		},
	}
	opts := azcore.ClientOptions{Cloud: c}

	cred, err := azidentity.NewDefaultAzureCredential(
		&azidentity.DefaultAzureCredentialOptions{ClientOptions: opts},
	)
	handle(err)

	client, err := armsubscription.NewClient(
		cred, &arm.ClientOptions{ClientOptions: opts},
	)
	handle(err)
*/
package cloud
