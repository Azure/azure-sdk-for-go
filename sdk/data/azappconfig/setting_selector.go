// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/generated"
)

// SettingFilter to select configuration setting entities.
type SettingFilter struct {
	// Key filter that will be used to select a set of configuration setting entities.
	KeyFilter *string

	// Label filter that will be used to select a set of configuration setting entities.
	LabelFilter *string
}

// SettingSelector is a set of options that allows selecting a filtered set of configuration setting entities
// from the configuration store, and optionally allows indicating which fields of each setting to retrieve.
type SettingSelector struct {
	// Key filter that will be used to select a set of configuration setting entities.
	KeyFilter *string

	// Label filter that will be used to select a set of configuration setting entities.
	LabelFilter *string

	// Tags filter that will be used to select a set of configuration setting entities.
	// This is a list of tag filters in the format {tagName=tagValue}. For more information about filtering by tags, see:
	// https://aka.ms/azconfig/docs/keyvaluefiltering
	TagsFilter []string

	// Indicates the point in time in the revision history of the selected configuration setting entities to retrieve.
	// If set, all properties of the configuration setting entities in the returned group will be exactly what they were at this time.
	AcceptDateTime *time.Time

	// The fields of the configuration setting to retrieve for each setting in the retrieved group.
	Fields []SettingFields
}

// AllSettingFields returns a collection of all setting fields to use in SettingSelector.
func AllSettingFields() []SettingFields {
	return []SettingFields{
		SettingFieldsKey,
		SettingFieldsLabel,
		SettingFieldsValue,
		SettingFieldsContentType,
		SettingFieldsETag,
		SettingFieldsLastModified,
		SettingFieldsIsReadOnly,
		SettingFieldsTags,
	}
}

func (sc SettingSelector) toGeneratedGetRevisions() *generated.AzureAppConfigurationClientGetRevisionsOptions {
	var dt *string
	if sc.AcceptDateTime != nil {
		str := sc.AcceptDateTime.Format(timeFormat)
		dt = &str
	}

	sf := make([]SettingFields, len(sc.Fields))
	for i := range sc.Fields {
		sf[i] = SettingFields(sc.Fields[i])
	}

	return &generated.AzureAppConfigurationClientGetRevisionsOptions{
		After:  dt,
		Key:    sc.KeyFilter,
		Label:  sc.LabelFilter,
		Select: sf,
		Tags:   sc.TagsFilter,
	}
}

func (sc SettingSelector) toGeneratedGetKeyValues() *generated.AzureAppConfigurationClientGetKeyValuesOptions {
	var dt *string
	if sc.AcceptDateTime != nil {
		str := sc.AcceptDateTime.Format(timeFormat)
		dt = &str
	}

	sf := make([]SettingFields, len(sc.Fields))
	for i := range sc.Fields {
		sf[i] = SettingFields(sc.Fields[i])
	}

	return &generated.AzureAppConfigurationClientGetKeyValuesOptions{
		After:  dt,
		Key:    sc.KeyFilter,
		Label:  sc.LabelFilter,
		Select: sf,
		Tags:   sc.TagsFilter,
	}
}
