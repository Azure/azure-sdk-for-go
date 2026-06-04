// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

// this file contains handwritten additions to the generated enums in constants.go

// The secureWrapKey and secureUnwrapKey KeyOperation values are accepted by the
// Key Vault service but are not yet part of the published KeyOperation enum in
// the service's OpenAPI/TypeSpec definition. They are declared here so callers
// can use them via the strongly-typed KeyOperation type.
const (
	KeyOperationSecureWrapKey   KeyOperation = "secureWrapKey"
	KeyOperationSecureUnwrapKey KeyOperation = "secureUnwrapKey"
)
