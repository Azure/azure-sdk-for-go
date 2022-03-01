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

// Setting is a setting, defined by a unique combination of a Key and Label.
type Setting struct {
	Key          *string
	Value        *string
	Label        *string
	ContentType  *string
	ETag         *azcore.ETag
	Tags         map[string]*string
	LastModified *time.Time
	IsReadOnly   *bool
}

func configurationSettingFromGenerated(kv generated.KeyValue) Setting {
	return Setting{
		Key:          kv.Key,
		Value:        kv.Value,
		Label:        kv.Label,
		ContentType:  kv.ContentType,
		ETag:         (*azcore.ETag)(kv.Etag),
		Tags:         kv.Tags,
		LastModified: kv.LastModified,
		IsReadOnly:   kv.Locked,
	}
}

func (cs Setting) toGenerated() *generated.KeyValue {
	return &generated.KeyValue{
		ContentType:  cs.ContentType,
		Etag:         (*string)(cs.ETag),
		Key:          cs.Key,
		Label:        cs.Label,
		LastModified: cs.LastModified,
		Locked:       cs.IsReadOnly,
		Tags:         cs.Tags,
		Value:        cs.Value,
	}
}
