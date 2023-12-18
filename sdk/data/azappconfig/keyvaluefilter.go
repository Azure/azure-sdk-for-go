//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

// KeyValueFilter contains filters to retrieve key-values from a configuration store.
type KeyValueFilter struct {
	// REQUIRED; Filters key-values by their key field.
	Key *string

	// Filters key-values by their label field.
	Label *string
}
