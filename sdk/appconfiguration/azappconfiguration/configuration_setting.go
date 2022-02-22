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

// ConfigurationSetting is a setting, defined by a unique combination of a key and label.
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

// NewConfigurationSetting creates a new ConfigurationSetting.
func NewConfigurationSetting(key string) ConfigurationSetting {
	return ConfigurationSetting{
		key: &key,
	}
}

// WithKey creates a copy of the configuration setting with the key provided.
func (cs ConfigurationSetting) WithKey(key string) ConfigurationSetting {
	result := cs
	result.key = &key
	return result
}

// GetKey gets the configuration setting key.
func (cs ConfigurationSetting) GetKey() *string {
	return cs.key
}

// WithValue creates a copy of the configuration setting with the value provided.
func (cs ConfigurationSetting) WithValue(value string) ConfigurationSetting {
	result := cs
	result.value = &value
	return result
}

// GetValue gets the configuration setting value.
func (cs ConfigurationSetting) GetValue() *string {
	return cs.value
}

// WithLabel creates a copy of the configuration setting with the label provided.
func (cs ConfigurationSetting) WithLabel(label string) ConfigurationSetting {
	result := cs
	result.label = &label
	return result
}

// GetLabel gets the configuration setting label.
func (cs ConfigurationSetting) GetLabel() *string {
	return cs.label
}

// WithContentType creates a copy of the configuration setting with the content type provided.
func (cs ConfigurationSetting) WithContentType(contentType string) ConfigurationSetting {
	result := cs
	result.contentType = &contentType
	return result
}

// GetContentType gets the configuration setting content type.
func (cs ConfigurationSetting) GetContentType() *string {
	return cs.contentType
}

// GetETag gets the configuration setting content ETag.
func (cs ConfigurationSetting) GetETag() *azcore.ETag {
	return cs.etag
}

// GetTags gets the list of the configuration setting tags.
func (cs ConfigurationSetting) GetTags() map[string]*string {
	return cs.tags
}

// GetLastModified gets the configuration setting last modified timestamp.
func (cs ConfigurationSetting) GetLastModified() *time.Time {
	return cs.lastModified
}

// IsReadOnly gets the read-only status of the configuration setting.
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
