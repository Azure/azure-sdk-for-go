//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"time"

	"sdk/appconfiguration/azappconfiguration/internal/generated"
)

type SettingFields string

const (
	SettingFieldsKey          = SettingFields(generated.Enum6Key)
	SettingFieldsLabel        = SettingFields(generated.Enum6Label)
	SettingFieldsValue        = SettingFields(generated.Enum6Value)
	SettingFieldsContentType  = SettingFields(generated.Enum6ContentType)
	SettingFieldsETag         = SettingFields(generated.Enum6Etag)
	SettingFieldsLastModified = SettingFields(generated.Enum6LastModified)
	SettingFieldsIsReadOnly   = SettingFields(generated.Enum6Locked)
	SettingFieldsTags         = SettingFields(generated.Enum6Tags)
)

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

func NewSettingSelector(key string) SettingSelector {
	return SettingSelector{fields: allSettingFields}
}

func (ss SettingSelector) WithKeyFilter(keyFilter string) SettingSelector {
	result := ss
	result.keyFilter = &keyFilter
	return result
}

func (ss SettingSelector) WittLabelFilter(labelFilter string) SettingSelector {
	result := ss
	result.labelFilter = &labelFilter
	return result
}

func (ss SettingSelector) WithAcceptDateTime(acceptDateTime time.Time) SettingSelector {
	result := ss
	result.acceptDateTime = &acceptDateTime
	return result
}

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
