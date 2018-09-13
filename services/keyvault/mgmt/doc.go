// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

/*
	Azure Key Vault management allows you to create and manage the 'vaults' which
	store secrets in Azure. These vaults are backed by secure hardware modules
	to keep your information safe.

	You can learn more about Azure Key Vault on the Microsoft documentation site:

	    - What is Azure Key Vault: https://docs.microsoft.com/en-us/azure/key-vault/key-vault-overview
	    - About keys, secrets, and certificates: https://docs.microsoft.com/en-us/azure/key-vault/about-keys-secrets-and-certificates
	    - Secure your key vault: https://docs.microsoft.com/en-us/azure/key-vault/key-vault-secure-your-key-vault

	Example - Create a new vault

	The following example shows how to create a new vault which can be managed by the provided user,
	Some helper functions are used to elide boilerplate details of authentication and configuring a
	connection to the management plane of key vault.

	    func CreateKeyVault(ctx context.Context, vaultName, userID string) (vault keyvault.Vault, err error) {
		vaultsClient := getVaultsClient()

		tenantID, err := uuid.FromString(iam.TenantID())
		if err != nil {
		    return
		}

		apList := []keyvault.AccessPolicyEntry{
		    keyvault.AccessPolicyEntry{
				ObjectID: to.StringPtr(userID),
				TenantID: &tenantID,
				Permissions: &keyvault.Permissions{
					Keys: &[]keyvault.KeyPermissions{
						keyvault.KeyPermissionsCreate,
					},
					Secrets: &[]keyvault.SecretPermissions{
						keyvault.SecretPermissionsSet,
					},
				},
		    },
		}

		return vaultsClient.CreateOrUpdate(
		    ctx,
		    helpers.ResourceGroupName(),
		    vaultName,
		    keyvault.VaultCreateOrUpdateParameters{
			Location: to.StringPtr(helpers.Location()),
			Properties: &keyvault.VaultProperties{
			    AccessPolicies:           &apList,
			    EnabledForDiskEncryption: to.BoolPtr(true),
			    Sku: &keyvault.Sku{
					Family: to.StringPtr("A"),
					Name:   keyvault.Standard,
			    },
			    TenantID: &tenantID,
			},
		    })
	    }

	Samples

	You can see all of the available samples for Azure Key Vault by browsing our samples repository. All samples are runnable through go test.
	    https://github.com/Azure-Samples/azure-sdk-for-go-samples/tree/master/keyvault

	Contribute

	We encourage and welcome contributions to the Azure SDK for Go. The best ways to help us out are:

	    - File an issue on github: https://github.com/azure/azure-sdk-for-go/issues
	    - Contribute a pull request: https://github.com/Azure/azure-sdk-for-go/blob/master/CONTRIBUTING.md

	All Microsoft open source projects follow a standard Code of Conduct. You can learn more about our CoC at
	https://opensource.microsoft.com/codeofconduct/
*/
package keyvault
