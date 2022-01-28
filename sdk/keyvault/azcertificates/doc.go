//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

/*

Key Vault certificates support provides for management of your x509 certificates and the following behaviors:

* Allows a certificate owner to create a certificate through a Key Vault creation process or through the import of an existing certificate. Includes both self-signed and Certificate Authority generated certificates.
* Allows a Key Vault certificate owner to implement secure storage and management of X509 certificates without interaction with private key material.
* Allows a certificate owner to create a policy that directs Key Vault to manage the life-cycle of a certificate.
* Allows certificate owners to provide contact information for notification about life-cycle events of expiration and renewal of certificate.
* Supports automatic renewal with selected issuers - Key Vault partner X509 certificate providers / certificate authorities.

*/

package azcertificates
