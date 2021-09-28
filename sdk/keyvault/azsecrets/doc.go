// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
/*

Package azsecrets can be used to access Azure KeyVault Secrets instance.

Azure KeyVault helps securely store and control access to tokens, passwords, certificates, API
keys, and other secrets.

A secret consists of a secret value and its associated metadata and management information. This
library library handles secret values as strings, but Azure Key Vault does not store them
as such. For more information about secrets about secrets and how Key Vault stores and manages them,
check out the Key Vault documentation (https://docs.microsoft.com/azure/key-vault/general/about-keys-secrets-certificates).

*/

package azsecrets
