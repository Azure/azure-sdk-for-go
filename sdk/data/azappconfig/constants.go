//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import "github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/internal/generated"

// SettingFields are fields to retrieve from a configuration setting.
type SettingFields = generated.SettingFields

const (
	// The primary identifier of a configuration setting.
	SettingFieldsKey SettingFields = generated.SettingFieldsKey

	// A label used to group configuration settings.
	SettingFieldsLabel SettingFields = generated.SettingFieldsLabel

	// The value of the configuration setting.
	SettingFieldsValue SettingFields = generated.SettingFieldsValue

	// The content type of the configuration setting's value.
	SettingFieldsContentType SettingFields = generated.SettingFieldsContentType

	// An ETag indicating the version of a configuration setting within a configuration store.
	SettingFieldsETag SettingFields = generated.SettingFieldsEtag

	// The last time a modifying operation was performed on the given configuration setting.
	SettingFieldsLastModified SettingFields = generated.SettingFieldsLastModified

	// A value indicating whether the configuration setting is read-only.
	SettingFieldsIsReadOnly SettingFields = generated.SettingFieldsLocked

	// A list of tags that can help identify what a configuration setting may be applicable for.
	SettingFieldsTags SettingFields = generated.SettingFieldsTags
)
