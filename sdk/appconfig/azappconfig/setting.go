//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"

	"github.com/Azure/azure-sdk-for-go/sdk/appconfig/azappconfig/internal/generated"
)

// Setting is a setting, defined by a unique combination of a Key and Label.
type Setting struct {
	// The primary identifier of the configuration setting.
	// A Key is used together with a Label to uniquely identify a configuration setting.
	Key *string

	// The configuration setting's value.
	Value *string

	// A value used to group configuration settings.
	// A Label is used together with a Key to uniquely identify a configuration setting.
	Label *string

	// The content type of the configuration setting's value.
	// Providing a proper content-type can enable transformations of values when they are retrieved by applications.
	ContentType *string

	// An ETag indicating the state of a configuration setting within a configuration store.
	ETag *azcore.ETag

	// A dictionary of tags used to assign additional properties to a configuration setting.
	// These can be used to indicate how a configuration setting may be applied.
	Tags map[string]string

	// The last time a modifying operation was performed on the given configuration setting.
	LastModified *time.Time

	// A value indicating whether the configuration setting is read only.
	// A read only configuration setting may not be modified until it is made writable.
	IsReadOnly *bool
}

func settingFromGenerated(kv generated.KeyValue) Setting {
	var tags map[string]string
	for k, v := range kv.Tags {
		if v != nil {
			tags[k] = *v
		}
	}

	return Setting{
		Key:          kv.Key,
		Value:        kv.Value,
		Label:        kv.Label,
		ContentType:  kv.ContentType,
		ETag:         (*azcore.ETag)(kv.Etag),
		Tags:         tags,
		LastModified: kv.LastModified,
		IsReadOnly:   kv.Locked,
	}
}

func (cs Setting) toGenerated() *generated.KeyValue {
	var tags map[string]*string
	for k, v := range cs.Tags {
		tags[k] = &v
	}

	return &generated.KeyValue{
		ContentType:  cs.ContentType,
		Etag:         (*string)(cs.ETag),
		Key:          cs.Key,
		Label:        cs.Label,
		LastModified: cs.LastModified,
		Locked:       cs.IsReadOnly,
		Tags:         tags,
		Value:        cs.Value,
	}
}
