//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azappconfiguration

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func TestConfigurationSettingWithKey(t *testing.T) {
	cs := NewConfigurationSetting("key")

	require.NotEmpty(t, cs.GetKey())
	require.Equal(t, *cs.GetKey(), "key")

	require.Empty(t, cs.GetValue())
	require.Empty(t, cs.GetLabel())
	require.Empty(t, cs.GetContentType())
	require.Empty(t, cs.GetETag())
	require.Empty(t, cs.GetTags())
	require.Empty(t, cs.GetLastModified())
	require.Empty(t, cs.IsReadOnly())
}

func TestConfigurationSettingWithKeyValue(t *testing.T) {
	cs := NewConfigurationSetting("key").WithValue("value")

	require.NotEmpty(t, cs.GetKey())
	require.Equal(t, *cs.GetKey(), "key")

	require.NotEmpty(t, cs.GetValue())
	require.Equal(t, *cs.GetValue(), "value")

	require.Empty(t, cs.GetLabel())
	require.Empty(t, cs.GetContentType())
	require.Empty(t, cs.GetETag())
	require.Empty(t, cs.GetTags())
	require.Empty(t, cs.GetLastModified())
	require.Empty(t, cs.IsReadOnly())
}

func TestConfigurationSettingWithKeyValueLabel(t *testing.T) {
	cs := NewConfigurationSetting("key").WithValue("value").WithLabel("label")

	require.NotEmpty(t, cs.GetKey())
	require.Equal(t, *cs.GetKey(), "key")

	require.NotEmpty(t, cs.GetValue())
	require.Equal(t, *cs.GetValue(), "value")

	require.NotEmpty(t, cs.GetLabel())
	require.Equal(t, *cs.GetLabel(), "label")

	require.Empty(t, cs.GetContentType())
	require.Empty(t, cs.GetETag())
	require.Empty(t, cs.GetTags())
	require.Empty(t, cs.GetLastModified())
	require.Empty(t, cs.IsReadOnly())
}

func TestConfigurationSettingWithKeyValueContentType(t *testing.T) {
	cs := NewConfigurationSetting("key").WithValue("value").WithContentType("contentType")

	require.NotEmpty(t, cs.GetKey())
	require.Equal(t, *cs.GetKey(), "key")

	require.NotEmpty(t, cs.GetValue())
	require.Equal(t, *cs.GetValue(), "value")

	require.Empty(t, cs.GetLabel())

	require.NotEmpty(t, cs.GetContentType())
	require.Equal(t, *cs.GetContentType(), "contentType")

	require.Empty(t, cs.GetETag())
	require.Empty(t, cs.GetTags())
	require.Empty(t, cs.GetLastModified())
	require.Empty(t, cs.IsReadOnly())
}

func TestConfigurationSettingIncrementalBuilding(t *testing.T) {
	cs := NewConfigurationSetting("key")

	require.NotEmpty(t, cs.GetKey())
	require.Equal(t, *cs.GetKey(), "key")

	require.Empty(t, cs.GetValue())
	require.Empty(t, cs.GetLabel())
	require.Empty(t, cs.GetContentType())
	require.Empty(t, cs.GetETag())
	require.Empty(t, cs.GetTags())
	require.Empty(t, cs.GetLastModified())
	require.Empty(t, cs.IsReadOnly())

	cs = cs.WithValue("value1")

	require.NotEmpty(t, cs.GetKey())
	require.Equal(t, *cs.GetKey(), "key")

	require.NotEmpty(t, cs.GetValue())
	require.Equal(t, *cs.GetValue(), "value1")

	require.Empty(t, cs.GetLabel())
	require.Empty(t, cs.GetContentType())
	require.Empty(t, cs.GetETag())
	require.Empty(t, cs.GetTags())
	require.Empty(t, cs.GetLastModified())
	require.Empty(t, cs.IsReadOnly())

	cs = cs.WithLabel("label")

	require.NotEmpty(t, cs.GetKey())
	require.Equal(t, *cs.GetKey(), "key")

	require.NotEmpty(t, cs.GetValue())
	require.Equal(t, *cs.GetValue(), "value1")

	require.NotEmpty(t, cs.GetLabel())
	require.Equal(t, *cs.GetLabel(), "label")

	require.Empty(t, cs.GetContentType())
	require.Empty(t, cs.GetETag())
	require.Empty(t, cs.GetTags())
	require.Empty(t, cs.GetLastModified())
	require.Empty(t, cs.IsReadOnly())

	cs = cs.WithValue("value2")

	require.NotEmpty(t, cs.GetKey())
	require.Equal(t, *cs.GetKey(), "key")

	require.NotEmpty(t, cs.GetValue())
	require.Equal(t, *cs.GetValue(), "value2")

	require.NotEmpty(t, cs.GetLabel())
	require.Equal(t, *cs.GetLabel(), "label")

	require.Empty(t, cs.GetContentType())
	require.Empty(t, cs.GetETag())
	require.Empty(t, cs.GetTags())
	require.Empty(t, cs.GetLastModified())
	require.Empty(t, cs.IsReadOnly())
}

func TestConfigurationSettingPrivateFieldGetters(t *testing.T) {
	key := "key"
	value := "value"
	label := "label"
	contentType := "contentType"
	etag := azcore.ETagAny

	var tags map[string]*string
	tagValue := "tagValue"
	tags["tagKey"] = &tagvalue

	isReadOnly := true

	cs := Setting{
		key:         &key,
		value:       &value,
		label:       &label,
		contentType: &contentType,
		etag:        &etag,
		tags:        tags,
		isReadOnly:  &isReadOnly,
	}

	require.Equal(t, cs.GetKey(), &key)
	require.Equal(t, cs.GetValue(), &value)
	require.Equal(t, cs.GetLabel(), &label)
	require.Equal(t, cs.GetContentType(), &contentType)
	require.Equal(t, cs.GetETag(), &etag)
	require.Equal(t, cs.IsReadOnly, &isReadOnly)

	tv, ok := cs.GetTags()["tagKey"]
	require.Equal(t, ok, true)
	require.Equal(t, tv, &tagValue)
	require.Len(t, cs.GetTags(), 1)
}
