//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azappconfig

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	key := "key"
	label := "label"
	value := "value"
	client, err := NewClientFromConnectionString(os.Getenv("APPCONFIGURATION_CONNECTION_STRING"), nil)
	require.NoError(t, err)
	require.NotEmpty(t, client)

	addResp, err2 := client.AddSetting(context.TODO(), Setting{Key: &key, Label: &label, Value: &value}, nil)
	require.NoError(t, err2)
	require.NotEmpty(t, addResp)
	require.NotNil(t, addResp.Key)
	require.NotNil(t, addResp.Label)
	require.NotNil(t, addResp.Value)
	require.Equal(t, key, *addResp.Key)
	require.Equal(t, label, *addResp.Label)
	require.Equal(t, value, *addResp.Value)

	getResp, err3 := client.GetSetting(context.TODO(), Setting{Key: &key, Label: &label}, nil)
	require.NoError(t, err3)
	require.NotEmpty(t, getResp)
	require.NotNil(t, getResp.Key)
	require.NotNil(t, getResp.Label)
	require.NotNil(t, getResp.Value)
	require.Equal(t, key, *getResp.Key)
	require.Equal(t, label, *getResp.Label)
	require.Equal(t, value, *getResp.Value)

	value = "value2"
	setResp, err4 := client.SetSetting(context.TODO(), Setting{Key: &key, Label: &label, Value: &value}, nil)
	require.NoError(t, err4)
	require.NotEmpty(t, setResp)
	require.NotNil(t, setResp.Key)
	require.NotNil(t, setResp.Label)
	require.NotNil(t, setResp.Value)
	require.Equal(t, key, *setResp.Key)
	require.Equal(t, label, *setResp.Label)
	require.Equal(t, value, *setResp.Value)

	roResp, err5 := client.SetReadOnly(context.TODO(), Setting{Key: &key, Label: &label}, true, nil)
	require.NoError(t, err5)
	require.NotEmpty(t, roResp)
	require.NotNil(t, roResp.Key)
	require.NotNil(t, roResp.Label)
	require.NotNil(t, roResp.Value)
	require.NotNil(t, roResp.IsReadOnly)
	require.Equal(t, key, *roResp.Key)
	require.Equal(t, label, *roResp.Label)
	require.Equal(t, value, *roResp.Value)
	require.True(t, *roResp.IsReadOnly)

	roResp2, err6 := client.SetReadOnly(context.TODO(), Setting{Key: &key, Label: &label}, false, nil)
	require.NoError(t, err6)
	require.NotEmpty(t, roResp2)
	require.NotNil(t, roResp2.Key)
	require.NotNil(t, roResp2.Label)
	require.NotNil(t, roResp2.Value)
	require.NotNil(t, roResp2.IsReadOnly)
	require.Equal(t, key, *roResp2.Key)
	require.Equal(t, label, *roResp2.Label)
	require.Equal(t, value, *roResp2.Value)
	require.False(t, *roResp2.IsReadOnly)

	any := "*"
	revPgr := client.ListRevisions(SettingSelector{KeyFilter: &any, LabelFilter: &any, Fields: AllSettingFields()}, nil)
	require.NotEmpty(t, revPgr)
	revHasPage := revPgr.NextPage(context.TODO())
	require.True(t, revHasPage)
	revResp := revPgr.PageResponse()
	require.NotEmpty(t, revResp)
	require.Equal(t, key, *revResp.Settings[0].Key)
	require.Equal(t, label, *revResp.Settings[0].Label)
	require.Equal(t, value, *revResp.Settings[0].Value)

	delResp, err7 := client.DeleteSetting(context.TODO(), Setting{Key: &key, Label: &label}, nil)
	require.NoError(t, err7)
	require.NotEmpty(t, delResp)
	require.NotNil(t, delResp.Key)
	require.NotNil(t, delResp.Label)
	require.NotNil(t, delResp.Value)
	require.Equal(t, key, *delResp.Key)
	require.Equal(t, label, *delResp.Label)
	require.Equal(t, value, *delResp.Value)
}
