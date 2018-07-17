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
	Azure Key Vault allows you to create and manage secrets and certificates
	on the cloud in a secure fashion. These secrets are stored in 'vaults', backed by secure hardware
	modules.

	To create and manage vaults themselves, use the mgmt subpackage to interact with the key vault
	management API.

	You can learn more about Azure Key vault on the Microsoft documentation site:

	    - What is Azure Key Vault: https://docs.microsoft.com/en-us/azure/key-vault/key-vault-overview
	    - About keys, secrets, and certificates: https://docs.microsoft.com/en-us/azure/key-vault/about-keys-secrets-and-certificates
	    - Storage Account Keys: https://docs.microsoft.com/en-us/azure/key-vault/key-vault-ovw-storage-keys

	Example - Create a new key

	The following example shows how to create a new key, given a vault name. Some helper functions are included
	to elide boilerplate details of authentication and configuring a connection to the management plane of key vault.

	    func CreateKeyBundle(ctx context.Context, vaultName string) (key keyvault.KeyBundle, err error) {
		vaultsClient := getVaultsClient()
		vault, err := vaultsClient.Get(ctx, helpers.ResourceGroupName(), vaultName)
		if err != nil {
		    return
		}
		vaultURL := *vault.Properties.VaultURI

		keyClient := getKeysClient()
		return keyClient.CreateKey(
		    ctx,
		    vaultURL,
		    keyName,
		    keyvault.KeyCreateParameters{
			KeyAttributes: &keyvault.KeyAttributes{
			    Enabled: to.BoolPtr(true),
			},
			KeySize: to.Int32Ptr(2048), // As of writing this sample, 2048 is the only supported KeySize.
			KeyOps: &[]keyvault.JSONWebKeyOperation{
			    keyvault.Encrypt,
			    keyvault.Decrypt,
			},
			Kty: keyvault.RSA,
		    })
	    }

	Samples

	You can see all of the available samples for key vault by browsing our samples repository. All samples are runnable through go test.
	    https://github.com/Azure-Samples/azure-sdk-for-go-samples/tree/master/keyvault

	Contribute

	We encourage and welcome contributions to the Azure SDK for Go. The best ways to help us out are:

	    - File an issue on github: https://github.com/azure/azure-sdk-for-go/issues
	    - Contribute a pull request: https://github.com/Azure/azure-sdk-for-go/blob/master/CONTRIBUTING.md

	All Microsoft open source projects follow a standard Code of Conduct. You can learn more about our CoC at
	https://opensource.microsoft.com/codeofconduct/
*/
package keyvault
