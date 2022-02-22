//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"

	"github.com/Azure/azure-sdk-for-go/sdk/appconfiguration/azappconfiguration/internal/generated"
)

type ConfigurationSetting struct {
	key          *string
	value        *string
	label        *string
	contentType  *string
	etag         *azcore.ETag
	tags         map[string]*string
	lastModified *time.Time
	isReadOnly   *bool
}

func NewConfigurationSetting(key string) ConfigurationSetting {
	return ConfigurationSetting{
		key: &key,
	}
}

func (cs ConfigurationSetting) WithKey(key string) ConfigurationSetting {
	result := cs
	result.key = &key
	return result
}

func (cs ConfigurationSetting) GetKey() *string {
	return cs.key
}

func (cs ConfigurationSetting) WithValue(value string) ConfigurationSetting {
	result := cs
	result.value = &value
	return result
}

func (cs ConfigurationSetting) GetValue() *string {
	return cs.value
}

func (cs ConfigurationSetting) WithLabel(label string) ConfigurationSetting {
	result := cs
	result.label = &label
	return result
}

func (cs ConfigurationSetting) GetLabel() *string {
	return cs.label
}

func (cs ConfigurationSetting) WithContentType(contentType string) ConfigurationSetting {
	result := cs
	result.contentType = &contentType
	return result
}

func (cs ConfigurationSetting) GetContentType() *string {
	return cs.contentType
}

func (cs ConfigurationSetting) GetETag() *azcore.ETag {
	return cs.etag
}

func (cs ConfigurationSetting) GetTags() map[string]*string {
	return cs.tags
}

func (cs ConfigurationSetting) GetLastModified() *time.Time {
	return cs.lastModified
}

func (cs ConfigurationSetting) IsReadOnly() *bool {
	return cs.isReadOnly
}

func configurationSettingFromGenerated(kv generated.KeyValue) ConfigurationSetting {
	return ConfigurationSetting{
		key:          kv.Key,
		value:        kv.Value,
		label:        kv.Label,
		contentType:  kv.ContentType,
		etag:         (*azcore.ETag)(kv.Etag),
		tags:         kv.Tags,
		lastModified: kv.LastModified,
		isReadOnly:   kv.Locked,
	}
}

func (cs ConfigurationSetting) toGenerated() *generated.KeyValue {
	return &generated.KeyValue{
		ContentType:  cs.contentType,
		Etag:         (*string)(cs.etag),
		Key:          cs.key,
		Label:        cs.label,
		LastModified: cs.lastModified,
		Locked:       cs.isReadOnly,
		Tags:         cs.tags,
		Value:        cs.value,
	}
}
