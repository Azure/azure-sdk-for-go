//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/generated"
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
	Tags map[string]*string

	// The last time a modifying operation was performed on the given configuration setting.
	LastModified *time.Time

	// A value indicating whether the configuration setting is read only.
	// A read only configuration setting may not be modified until it is made writable.
	IsReadOnly *bool
}

func settingFromGenerated(kv generated.KeyValue) Setting {
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

func toGeneratedETagString(etag *azcore.ETag) *string {
	if etag == nil || *etag == azcore.ETagAny {
		return (*string)(etag)
	}

	str := "\"" + (string)(*etag) + "\""
	return &str
}

func (cs Setting) toGenerated() generated.KeyValue {
	return generated.KeyValue{
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

func (cs Setting) toGeneratedDeleteLockOptions(ifMatch *azcore.ETag) *generated.AzureAppConfigurationClientDeleteLockOptions {
	return &generated.AzureAppConfigurationClientDeleteLockOptions{
		IfMatch: toGeneratedETagString(ifMatch),
		Label:   cs.Label,
	}
}

func (cs Setting) toGeneratedDeleteOptions(ifMatch *azcore.ETag) *generated.AzureAppConfigurationClientDeleteKeyValueOptions {
	return &generated.AzureAppConfigurationClientDeleteKeyValueOptions{
		IfMatch: toGeneratedETagString(ifMatch),
		Label:   cs.Label,
	}
}

func (cs Setting) toGeneratedGetOptions(ifNoneMatch *azcore.ETag, acceptDateTime *time.Time) *generated.AzureAppConfigurationClientGetKeyValueOptions {
	var dt *string
	if acceptDateTime != nil {
		str := acceptDateTime.Format(timeFormat)
		dt = &str
	}

	return &generated.AzureAppConfigurationClientGetKeyValueOptions{
		AcceptDatetime: dt,
		IfNoneMatch:    toGeneratedETagString(ifNoneMatch),
		Label:          cs.Label,
	}
}

func (cs Setting) toGeneratedPutLockOptions(ifMatch *azcore.ETag) *generated.AzureAppConfigurationClientPutLockOptions {
	return &generated.AzureAppConfigurationClientPutLockOptions{
		IfMatch: toGeneratedETagString(ifMatch),
		Label:   cs.Label,
	}
}

func (cs Setting) toGeneratedPutOptions(ifMatch *azcore.ETag, ifNoneMatch *azcore.ETag) (generated.KeyValue, generated.AzureAppConfigurationClientPutKeyValueOptions) {
	return cs.toGenerated(), generated.AzureAppConfigurationClientPutKeyValueOptions{
		IfMatch:     toGeneratedETagString(ifMatch),
		IfNoneMatch: toGeneratedETagString(ifNoneMatch),
		Label:       cs.Label,
	}
}
