//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

/*

Package azkeys can be used to access an Azure Key Vault Keys instance.

Azure KeyVault helps securely store and control access to tokens, passwords, certificates, API
keys, and other secrets. Azure Key Vault provides two types of resources to store and manage
cryptographic keys. Vaults support software-protected and HSM-protected (Hardware Security Module)
keys. Managed HSMs only support HSM-protected keys.

* Vaults - Vaults provide a low-cost, easy to deploy, multi-tenant, zone-resilient
(where available), highly available key management solution suitable for most common
cloud application scenarios.
* Managed HSMs - Managed HSM provides single-tenant, zone-resilient (where available),
highly available HSMs to store and manage your cryptographic keys. Most suitable for
applications and usage scenarios that handle high value keys. Also helps to meet most
stringent security, compliance, and regulatory requirements.

*/

package azkeys
