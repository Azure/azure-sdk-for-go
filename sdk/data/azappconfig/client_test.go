//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azappconfig_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	const (
		key   = "key-TestClient"
		label = "label"
	)

	contentType := "content-type"
	value := "value"
	client := NewClientFromConnectionString(t)

	addResp, err2 := client.AddSetting(context.Background(), key, &value, &azappconfig.AddSettingOptions{
		Label:       to.Ptr(label),
		ContentType: &contentType,
	})
	require.NoError(t, err2)
	require.NotEmpty(t, addResp)
	require.NotNil(t, addResp.Key)
	require.NotNil(t, addResp.Label)
	require.NotNil(t, addResp.ContentType)
	require.NotNil(t, addResp.Value)
	require.Equal(t, key, *addResp.Key)
	require.Equal(t, label, *addResp.Label)
	require.Equal(t, contentType, *addResp.ContentType)
	require.Equal(t, value, *addResp.Value)

	getResp, err3 := client.GetSetting(context.Background(), key, &azappconfig.GetSettingOptions{
		Label: to.Ptr(label),
	})
	require.NoError(t, err3)
	require.NotEmpty(t, getResp)
	require.NotNil(t, getResp.Key)
	require.NotNil(t, getResp.Label)
	require.NotNil(t, getResp.ContentType)
	require.NotNil(t, getResp.Value)
	require.Equal(t, key, *getResp.Key)
	require.Equal(t, label, *getResp.Label)
	require.Equal(t, contentType, *getResp.ContentType)
	require.Equal(t, value, *getResp.Value)

	etag := getResp.ETag
	getResp2, err4 := client.GetSetting(context.Background(), key, &azappconfig.GetSettingOptions{
		Label:         to.Ptr(label),
		OnlyIfChanged: etag,
	})
	require.Error(t, err4)
	require.Empty(t, getResp2)

	value = "value2"
	contentType = "content-type2"
	setResp, err5 := client.SetSetting(context.Background(), key, &value, &azappconfig.SetSettingOptions{
		Label:       to.Ptr(label),
		ContentType: &contentType,
	})
	require.NoError(t, err5)
	require.NotEmpty(t, setResp)
	require.NotNil(t, setResp.Key)
	require.NotNil(t, setResp.Label)
	require.NotNil(t, setResp.ContentType)
	require.NotNil(t, setResp.Value)
	require.Equal(t, key, *setResp.Key)
	require.Equal(t, label, *setResp.Label)
	require.Equal(t, contentType, *setResp.ContentType)
	require.Equal(t, value, *setResp.Value)
	require.NotNil(t, setResp.SyncToken)

	// after changing the setting, update the sync token so we don't get cached, stale data
	client.UpdateSyncToken(*setResp.SyncToken)

	getResp3, err6 := client.GetSetting(context.Background(), key, &azappconfig.GetSettingOptions{
		Label:         to.Ptr(label),
		OnlyIfChanged: etag,
	})
	require.NoError(t, err6)
	require.NotEmpty(t, getResp3)
	require.NotNil(t, getResp3.Key)
	require.NotNil(t, getResp3.Label)
	require.NotNil(t, getResp3.ContentType)
	require.NotNil(t, getResp3.Value)
	require.Equal(t, key, *getResp3.Key)
	require.Equal(t, label, *getResp3.Label)
	require.Equal(t, contentType, *getResp3.ContentType)
	require.Equal(t, value, *getResp3.Value)

	etag = getResp3.ETag
	value = "value3"
	setResp2, err7 := client.SetSetting(context.Background(), key, &value, &azappconfig.SetSettingOptions{
		Label:           to.Ptr(label),
		OnlyIfUnchanged: etag,
	})
	require.NoError(t, err7)
	require.NotEmpty(t, setResp2)
	require.NotNil(t, setResp2.Key)
	require.NotNil(t, setResp2.Label)
	require.Nil(t, setResp2.ContentType)
	require.NotNil(t, setResp2.Value)
	require.Equal(t, key, *setResp2.Key)
	require.Equal(t, label, *setResp2.Label)
	require.Equal(t, value, *setResp2.Value)
	require.NotNil(t, setResp.SyncToken)

	// after changing the setting, update the sync token so we don't get cached, stale data
	client.UpdateSyncToken(*setResp.SyncToken)

	setResp3, err8 := client.SetSetting(context.Background(), key, &value, &azappconfig.SetSettingOptions{
		Label:           to.Ptr(label),
		OnlyIfUnchanged: etag,
	})
	require.Error(t, err8)
	require.Empty(t, setResp3)

	roResp, err9 := client.SetReadOnly(context.Background(), key, true, &azappconfig.SetReadOnlyOptions{
		Label: to.Ptr(label),
	})
	require.NoError(t, err9)
	require.NotEmpty(t, roResp)
	require.NotNil(t, roResp.Key)
	require.NotNil(t, roResp.Label)
	require.Nil(t, roResp.ContentType)
	require.NotNil(t, roResp.Value)
	require.NotNil(t, roResp.IsReadOnly)
	require.Equal(t, key, *roResp.Key)
	require.Equal(t, label, *roResp.Label)
	require.Equal(t, value, *roResp.Value)
	require.True(t, *roResp.IsReadOnly)
	require.NotNil(t, setResp.SyncToken)

	// after changing the setting, update the sync token so we don't get cached, stale data
	client.UpdateSyncToken(*setResp.SyncToken)

	roResp2, err10 := client.SetReadOnly(context.Background(), key, false, &azappconfig.SetReadOnlyOptions{
		Label: to.Ptr(label),
	})
	require.NoError(t, err10)
	require.NotEmpty(t, roResp2)
	require.NotNil(t, roResp2.Key)
	require.NotNil(t, roResp2.Label)
	require.Nil(t, roResp2.ContentType)
	require.NotNil(t, roResp2.Value)
	require.NotNil(t, roResp2.IsReadOnly)
	require.Equal(t, key, *roResp2.Key)
	require.Equal(t, label, *roResp2.Label)
	require.Equal(t, value, *roResp2.Value)
	require.False(t, *roResp2.IsReadOnly)
	require.NotNil(t, setResp.SyncToken)

	// after changing the setting, update the sync token so we don't get cached, stale data
	client.UpdateSyncToken(*setResp.SyncToken)

	roResp3, err11 := client.SetReadOnly(context.Background(), key, true, &azappconfig.SetReadOnlyOptions{
		Label:           to.Ptr(label),
		OnlyIfUnchanged: etag,
	})
	require.Error(t, err11)
	require.Empty(t, roResp3)

	etag = roResp2.ETag
	roResp4, err12 := client.SetReadOnly(context.Background(), key, true, &azappconfig.SetReadOnlyOptions{
		Label:           to.Ptr(label),
		OnlyIfUnchanged: etag,
	})
	require.NoError(t, err12)
	require.NotEmpty(t, roResp4)
	require.NotNil(t, roResp4.Key)
	require.NotNil(t, roResp4.Label)
	require.Nil(t, roResp4.ContentType)
	require.NotNil(t, roResp4.Value)
	require.NotNil(t, roResp4.IsReadOnly)
	require.Equal(t, key, *roResp4.Key)
	require.Equal(t, label, *roResp4.Label)
	require.Equal(t, value, *roResp4.Value)
	require.True(t, *roResp4.IsReadOnly)
	require.NotNil(t, setResp.SyncToken)

	// after changing the setting, update the sync token so we don't get cached, stale data
	client.UpdateSyncToken(*setResp.SyncToken)

	roResp5, err13 := client.SetReadOnly(context.Background(), key, false, &azappconfig.SetReadOnlyOptions{
		Label:           to.Ptr(label),
		OnlyIfUnchanged: etag,
	})
	require.Error(t, err13)
	require.Empty(t, roResp5)

	etag = roResp4.ETag
	roResp6, err14 := client.SetReadOnly(context.Background(), key, false, &azappconfig.SetReadOnlyOptions{
		Label:           to.Ptr(label),
		OnlyIfUnchanged: etag,
	})
	require.NoError(t, err14)
	require.NotEmpty(t, roResp6)
	require.NotNil(t, roResp6.Key)
	require.NotNil(t, roResp6.Label)
	require.Nil(t, roResp6.ContentType)
	require.NotNil(t, roResp6.Value)
	require.NotNil(t, roResp6.IsReadOnly)
	require.Equal(t, key, *roResp6.Key)
	require.Equal(t, label, *roResp6.Label)
	require.Equal(t, value, *roResp6.Value)
	require.False(t, *roResp6.IsReadOnly)
	require.NotNil(t, setResp.SyncToken)

	// after changing the setting, update the sync token so we don't get cached, stale data
	client.UpdateSyncToken(*setResp.SyncToken)

	any := "*"
	revPgr := client.NewListRevisionsPager(azappconfig.SettingSelector{
		KeyFilter:   &any,
		LabelFilter: &any,
		Fields:      azappconfig.AllSettingFields(),
	}, nil)
	require.NotEmpty(t, revPgr)
	hasMoreRevs := revPgr.More()
	require.True(t, hasMoreRevs)
	revResp, err15 := revPgr.NextPage(context.Background())
	require.NoError(t, err15)
	require.NotEmpty(t, revResp)
	require.Equal(t, key, *revResp.Settings[0].Key)
	require.Equal(t, label, *revResp.Settings[0].Label)

	settsPgr := client.NewListSettingsPager(azappconfig.SettingSelector{
		KeyFilter:   &any,
		LabelFilter: &any,
		Fields:      azappconfig.AllSettingFields(),
	}, nil)
	require.NotEmpty(t, settsPgr)
	hasMoreSetts := settsPgr.More()
	require.True(t, hasMoreSetts)
	settsResp, err16 := settsPgr.NextPage(context.Background())
	require.NoError(t, err16)
	require.NotEmpty(t, settsResp)
	require.Equal(t, key, *settsResp.Settings[0].Key)
	require.Equal(t, label, *settsResp.Settings[0].Label)
	require.Equal(t, value, *settsResp.Settings[0].Value)
	require.False(t, *settsResp.Settings[0].IsReadOnly)

	delResp, err17 := client.DeleteSetting(context.Background(), key, &azappconfig.DeleteSettingOptions{
		Label: to.Ptr(label),
	})
	require.NoError(t, err17)
	require.NotEmpty(t, delResp)
	require.NotNil(t, delResp.Key)
	require.NotNil(t, delResp.Label)
	require.Nil(t, delResp.ContentType)
	require.NotNil(t, delResp.Value)
	require.Equal(t, key, *delResp.Key)
	require.Equal(t, label, *delResp.Label)
	require.Equal(t, value, *delResp.Value)
	require.NotNil(t, setResp.SyncToken)

	// after changing the setting, update the sync token so we don't get cached, stale data
	client.UpdateSyncToken(*setResp.SyncToken)

	addResp2, err18 := client.AddSetting(context.Background(), key, &value, &azappconfig.AddSettingOptions{
		Label:       to.Ptr(label),
		ContentType: &contentType,
	})
	require.NoError(t, err18)
	require.NotEmpty(t, addResp2)
	require.NotNil(t, addResp2.Key)
	require.NotNil(t, addResp2.Label)
	require.NotNil(t, addResp2.ContentType)
	require.NotNil(t, addResp2.Value)
	require.Equal(t, key, *addResp2.Key)
	require.Equal(t, label, *addResp2.Label)
	require.Equal(t, contentType, *addResp2.ContentType)
	require.Equal(t, value, *addResp2.Value)

	delResp2, err19 := client.DeleteSetting(context.Background(), key, &azappconfig.DeleteSettingOptions{
		Label:           to.Ptr(label),
		OnlyIfUnchanged: etag,
	})
	require.Error(t, err19)
	require.Empty(t, delResp2)

	etag = addResp2.ETag
	delResp3, err20 := client.DeleteSetting(context.Background(), key, &azappconfig.DeleteSettingOptions{
		Label:           to.Ptr(label),
		OnlyIfUnchanged: etag,
	})
	require.NoError(t, err20)
	require.NotEmpty(t, delResp3)
	require.NotNil(t, delResp3.Key)
	require.NotNil(t, delResp3.Label)
	require.NotNil(t, delResp3.ContentType)
	require.NotNil(t, delResp3.Value)
	require.Equal(t, key, *delResp3.Key)
	require.Equal(t, label, *delResp3.Label)
	require.Equal(t, contentType, *delResp3.ContentType)
	require.Equal(t, value, *delResp3.Value)
}

func TestSettingNilValue(t *testing.T) {
	const (
		key         = "key-TestSettingNilValue"
		contentType = "content-type"
	)
	client := NewClientFromConnectionString(t)

	addResp, err := client.AddSetting(context.Background(), key, nil, &azappconfig.AddSettingOptions{
		ContentType: to.Ptr(contentType),
	})
	require.NoError(t, err)
	require.NotZero(t, addResp)

	resp, err := client.DeleteSetting(context.Background(), key, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Key)
	require.EqualValues(t, key, *resp.Key)
}
