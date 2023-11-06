//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import "github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/internal/generated"

// KeyValueFields are fields to retrieve from a configuration setting.
type KeyValueFields = generated.KeyValueFields

const (
	// The primary identifier of a configuration setting.
	KeyValueFieldsKey KeyValueFields = generated.KeyValueFieldsKey

	// A label used to group configuration settings.
	KeyValueFieldsLabel KeyValueFields = generated.KeyValueFieldsLabel

	// The value of the configuration setting.
	KeyValueFieldsValue KeyValueFields = generated.KeyValueFieldsValue

	// The content type of the configuration setting's value.
	KeyValueFieldsContentType KeyValueFields = generated.KeyValueFieldsContentType

	// An ETag indicating the version of a configuration setting within a configuration store.
	KeyValueFieldsETag KeyValueFields = generated.KeyValueFieldsEtag

	// The last time a modifying operation was performed on the given configuration setting.
	KeyValueFieldsLastModified KeyValueFields = generated.KeyValueFieldsLastModified

	// A value indicating whether the configuration setting is read-only.
	KeyValueFieldsIsReadOnly KeyValueFields = generated.KeyValueFieldsLocked

	// A list of tags that can help identify what a configuration setting may be applicable for.
	KeyValueFieldsTags KeyValueFields = generated.KeyValueFieldsTags
)
