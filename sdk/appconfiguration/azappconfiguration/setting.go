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

// Setting is a setting, defined by a unique combination of a key and label.
type Setting struct {
	key          *string
	value        *string
	label        *string
	contentType  *string
	etag         *azcore.ETag
	tags         map[string]*string
	lastModified *time.Time
	isReadOnly   *bool
}

// NewConfigurationSetting creates a new Setting.
func NewConfigurationSetting(key string) Setting {
	return Setting{
		key: &key,
	}
}

// WithKey creates a copy of the configuration setting with the key provided.
func (cs Setting) WithKey(key string) Setting {
	result := cs
	result.key = &key
	return result
}

// GetKey gets the configuration setting key.
func (cs Setting) GetKey() *string {
	return cs.key
}

// WithValue creates a copy of the configuration setting with the value provided.
func (cs Setting) WithValue(value string) Setting {
	result := cs
	result.value = &value
	return result
}

// GetValue gets the configuration setting value.
func (cs Setting) GetValue() *string {
	return cs.value
}

// WithLabel creates a copy of the configuration setting with the label provided.
func (cs Setting) WithLabel(label string) Setting {
	result := cs
	result.label = &label
	return result
}

// GetLabel gets the configuration setting label.
func (cs Setting) GetLabel() *string {
	return cs.label
}

// WithContentType creates a copy of the configuration setting with the content type provided.
func (cs Setting) WithContentType(contentType string) Setting {
	result := cs
	result.contentType = &contentType
	return result
}

// GetContentType gets the configuration setting content type.
func (cs Setting) GetContentType() *string {
	return cs.contentType
}

// GetETag gets the configuration setting content ETag.
func (cs Setting) GetETag() *azcore.ETag {
	return cs.etag
}

// GetTags gets the list of the configuration setting tags.
func (cs Setting) GetTags() map[string]*string {
	return cs.tags
}

// GetLastModified gets the configuration setting last modified timestamp.
func (cs Setting) GetLastModified() *time.Time {
	return cs.lastModified
}

// IsReadOnly gets the read-only status of the configuration setting.
func (cs Setting) IsReadOnly() *bool {
	return cs.isReadOnly
}

func configurationSettingFromGenerated(kv generated.KeyValue) Setting {
	return Setting{
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

func (cs Setting) toGenerated() *generated.KeyValue {
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
