//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"time"

	"sdk/appconfiguration/sdk/appconfiguration/azappconfiguration/internal/generated"
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
	KeyFilter      *string
	LabelFilter    *string
	AcceptDateTime *time.Time
	Fields         []SettingFields
}

const SettingSelectorFilterAny string = "*"
