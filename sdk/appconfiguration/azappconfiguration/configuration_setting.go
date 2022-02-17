//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"time"

	"sdk/appconfiguration/sdk/appconfiguration/azappconfiguration/internal/generated"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type ConfigurationSetting struct {
	Key         *string
	Value       *string
	Label       *string
	ContentType *string

	etag         *azcore.ETag
	tags         map[string]*string
	lastModified *time.Time
	isReadOnly   *bool
}

func (cs ConfigurationSetting) ETag() *azcore.ETag {
	return cs.etag
}

func (cs ConfigurationSetting) Tags() map[string]*string {
	return cs.tags
}

func (cs ConfigurationSetting) LastModified() *time.Time {
	return cs.lastModified
}

func (cs ConfigurationSetting) IsReadOnly() *bool {
	return cs.isReadOnly
}

func configurationSettingFromGenerated(kv generated.KeyValue) ConfigurationSetting {
	return ConfigurationSetting{
		Key:          kv.Key,
		Value:        kv.Value,
		Label:        kv.Label,
		ContentType:  kv.ContentType,
		etag:         (*azcore.ETag)(kv.Etag),
		tags:         kv.Tags,
		lastModified: kv.LastModified,
		isReadOnly:   kv.Locked,
	}
}

func (cs ConfigurationSetting) toGenerated() *generated.KeyValue {
	return &generated.KeyValue{
		ContentType:  cs.ContentType,
		Etag:         (*string)(cs.etag),
		Key:          cs.Key,
		Label:        cs.Label,
		LastModified: cs.lastModified,
		Locked:       cs.isReadOnly,
		Tags:         cs.tags,
		Value:        cs.Value,
	}
}
