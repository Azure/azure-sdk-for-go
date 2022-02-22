//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/appconfiguration/azappconfiguration/internal/generated"
)

// Fields to retrieve from a configuration setting.
type SettingFields string

const (
	// The primary identifier of a configuration setting.
	SettingFieldsKey = SettingFields(generated.Enum6Key)

	// A label used to group configuration settings.
	SettingFieldsLabel = SettingFields(generated.Enum6Label)

	// The value of the configuration setting.
	SettingFieldsValue = SettingFields(generated.Enum6Value)

	// The content type of the configuration setting's value.
	SettingFieldsContentType = SettingFields(generated.Enum6ContentType)

	// An ETag indicating the version of a configuration setting within a configuration store.
	SettingFieldsETag = SettingFields(generated.Enum6Etag)

	// The last time a modifying operation was performed on the given configuration setting.
	SettingFieldsLastModified = SettingFields(generated.Enum6LastModified)

	// A value indicating whether the configuration setting is read-only.
	SettingFieldsIsReadOnly = SettingFields(generated.Enum6Locked)

	// A list of tags that can help identify what a configuration setting may be applicable for.
	SettingFieldsTags = SettingFields(generated.Enum6Tags)
)

// SettingSelector is a set of options that allows selecting a filtered set of configuration setting entities
// from the configuration store, and optionally allows indicating which fields of each setting to retrieve.
type SettingSelector struct {
	keyFilter      *string
	labelFilter    *string
	acceptDateTime *time.Time
	fields         []SettingFields
}

var allSettingFields []SettingFields = []SettingFields{
	SettingFieldsKey,
	SettingFieldsLabel,
	SettingFieldsValue,
	SettingFieldsContentType,
	SettingFieldsETag,
	SettingFieldsLastModified,
	SettingFieldsIsReadOnly,
	SettingFieldsTags,
}

// NewSettingSelector creates a new setting selector.
func NewSettingSelector() SettingSelector {
	return SettingSelector{fields: allSettingFields}
}

// WithKeyFilter creates a copy of the setting selector with the key filter provided.
func (ss SettingSelector) WithKeyFilter(keyFilter string) SettingSelector {
	result := ss
	result.keyFilter = &keyFilter
	return result
}

// WithLabelFilter creates a copy of the setting selector with the label filter provided.
func (ss SettingSelector) WithLabelFilter(labelFilter string) SettingSelector {
	result := ss
	result.labelFilter = &labelFilter
	return result
}

// WithAcceptDateTime creates a copy of the setting selector with the time
// that specifies a point in time in the revision history of the selected entities to retrieve.
func (ss SettingSelector) WithAcceptDateTime(acceptDateTime time.Time) SettingSelector {
	result := ss
	result.acceptDateTime = &acceptDateTime
	return result
}

// WithAcceptDateTime creates a copy of the setting selector with the fields
// to retrieve for each setting in the retrieved group.
func (ss SettingSelector) WithFields(fields []SettingFields) SettingSelector {
	result := ss
	result.fields = fields
	return result
}

func (sc SettingSelector) toGenerated() *generated.AzureAppConfigurationClientGetRevisionsOptions {
	var dt *string
	if sc.acceptDateTime != nil {
		str := sc.acceptDateTime.Format(timeFormat)
		dt = &str
	}

	sf := make([]generated.Enum6Tags, len(sc.fields))
	for i := range sc.fields {
		sf[i] = (generated.Enum6Tags)(sc.fields[i])
	}

	return &generated.AzureAppConfigurationClientGetRevisionsOptions{
		AcceptDateTime: dt,
		Key:            sc.keyFilter,
		Label:          sc.labelFilter,
		Select:         sf,
	}
}
