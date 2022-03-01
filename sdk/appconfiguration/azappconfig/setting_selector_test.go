//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azappconfig

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSettingSelectorDefault(t *testing.T) {
	ss := NewSettingSelector()

	require.Empty(t, ss.keyFilter)
	require.Empty(t, ss.labelFilter)
	require.Empty(t, ss.acceptDateTime)

	require.Len(t, ss.fields, 8)
	require.Equal(t, ss.fields[0], SettingFieldsKey)
	require.Equal(t, ss.fields[1], SettingFieldsLabel)
	require.Equal(t, ss.fields[2], SettingFieldsValue)
	require.Equal(t, ss.fields[3], SettingFieldsContentType)
	require.Equal(t, ss.fields[4], SettingFieldsETag)
	require.Equal(t, ss.fields[5], SettingFieldsLastModified)
	require.Equal(t, ss.fields[6], SettingFieldsIsReadOnly)
	require.Equal(t, ss.fields[7], SettingFieldsTags)
}

func TestSettingSelectorWithKeyFilter(t *testing.T) {
	ss := NewSettingSelector().WithKeyFilter("kf")

	require.NotEmpty(t, ss.keyFilter)
	require.Equal(t, *ss.keyFilter, "kf")

	require.Empty(t, ss.labelFilter)
	require.Empty(t, ss.acceptDateTime)

	require.Len(t, ss.fields, 8)
	require.Equal(t, ss.fields[0], SettingFieldsKey)
	require.Equal(t, ss.fields[1], SettingFieldsLabel)
	require.Equal(t, ss.fields[2], SettingFieldsValue)
	require.Equal(t, ss.fields[3], SettingFieldsContentType)
	require.Equal(t, ss.fields[4], SettingFieldsETag)
	require.Equal(t, ss.fields[5], SettingFieldsLastModified)
	require.Equal(t, ss.fields[6], SettingFieldsIsReadOnly)
	require.Equal(t, ss.fields[7], SettingFieldsTags)
}

func TestSettingSelectorWithLabelFilter(t *testing.T) {
	ss := NewSettingSelector().WithLabelFilter("lf")

	require.Empty(t, ss.keyFilter)

	require.NotEmpty(t, ss.labelFilter)
	require.Equal(t, *ss.labelFilter, "lf")

	require.Empty(t, ss.acceptDateTime)

	require.Len(t, ss.fields, 8)
	require.Equal(t, ss.fields[0], SettingFieldsKey)
	require.Equal(t, ss.fields[1], SettingFieldsLabel)
	require.Equal(t, ss.fields[2], SettingFieldsValue)
	require.Equal(t, ss.fields[3], SettingFieldsContentType)
	require.Equal(t, ss.fields[4], SettingFieldsETag)
	require.Equal(t, ss.fields[5], SettingFieldsLastModified)
	require.Equal(t, ss.fields[6], SettingFieldsIsReadOnly)
	require.Equal(t, ss.fields[7], SettingFieldsTags)
}

func TestSettingSelectorWithAcceptDateTime(t *testing.T) {
	ss := NewSettingSelector().WithAcceptDateTime(time.Parse("2006-01-02", "2022-02-22"))

	require.Empty(t, ss.keyFilter)
	require.Empty(t, ss.labelFilter)

	require.NotEmpty(t, ss.acceptDateTime)
	require.Equal(t, *ss.acceptDateTime, time.Parse("2006-01-02", "2022-02-22"))

	require.Len(t, ss.fields, 8)
	require.Equal(t, ss.fields[0], SettingFieldsKey)
	require.Equal(t, ss.fields[1], SettingFieldsLabel)
	require.Equal(t, ss.fields[2], SettingFieldsValue)
	require.Equal(t, ss.fields[3], SettingFieldsContentType)
	require.Equal(t, ss.fields[4], SettingFieldsETag)
	require.Equal(t, ss.fields[5], SettingFieldsLastModified)
	require.Equal(t, ss.fields[6], SettingFieldsIsReadOnly)
	require.Equal(t, ss.fields[7], SettingFieldsTags)
}

func TestSettingSelectorIncrementalBuilding(t *testing.T) {
	ss := NewSettingSelector().WithLabelFilter("lf").WithKeyFilter("kf1")

	require.NotEmpty(t, ss.keyFilter)
	require.Equal(t, *ss.keyFilter, "kf1")

	require.NotEmpty(t, ss.labelFilter)
	require.Equal(t, *ss.labelFilter, "lf")

	require.Empty(t, ss.acceptDateTime)

	require.Len(t, ss.fields, 8)
	require.Equal(t, ss.fields[0], SettingFieldsKey)
	require.Equal(t, ss.fields[1], SettingFieldsLabel)
	require.Equal(t, ss.fields[2], SettingFieldsValue)
	require.Equal(t, ss.fields[3], SettingFieldsContentType)
	require.Equal(t, ss.fields[4], SettingFieldsETag)
	require.Equal(t, ss.fields[5], SettingFieldsLastModified)
	require.Equal(t, ss.fields[6], SettingFieldsIsReadOnly)
	require.Equal(t, ss.fields[7], SettingFieldsTags)

	ss = ss.WithKeyFilter("kf2").WithAcceptDateTime(time.Parse("2006-01-02", "2022-02-22")).WithFields([]SettingFields{})

	require.NotEmpty(t, ss.keyFilter)
	require.Equal(t, *ss.keyFilter, "kf2")

	require.NotEmpty(t, ss.labelFilter)
	require.Equal(t, *ss.labelFilter, "lf")

	require.NotEmpty(t, ss.acceptDateTime)
	require.Equal(t, *ss.acceptDateTime, time.Parse("2006-01-02", "2022-02-22"))

	require.Len(t, ss.fields, 0)
}
