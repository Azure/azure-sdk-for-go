// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package packages

// VerifyFunc ...
type VerifyFunc func(p Package) error

// Verifier could verify a SDK package
type Verifier interface {
	// Verify verifies the given package
	Verify(pkg Package) []error
}
