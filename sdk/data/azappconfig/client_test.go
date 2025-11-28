// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azappconfig_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

// testId will be used for local testing as a unique identifier. The test proxy will record the base
// request for snapshots. The deletion time for a snapshot is minimum 1 hour. For quicker
// local iteration we will use a unique suffix for each test run.
// to use: switch the testId being used
//
//	 Snapshot Name: `// + string(testId)`
//		KeyValue Prefix: `/*testId +*/`

// Record Mode
var testId = "120823uid"

// // Local Testing Mode
// var currTime = time.Now().Unix()
// var testId = strconv.FormatInt(currTime, 10)[len(strconv.FormatInt(currTime, 10))-6:]

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

func TestSettingWithEscaping(t *testing.T) {
	const (
		key         = ".appconfig.featureflag/TestSettingWithEscaping"
		contentType = "application/vnd.microsoft.appconfig.ff+json;charset=utf-8"
	)
	client := NewClientFromConnectionString(t)

	addResp, err := client.AddSetting(context.Background(), key, nil, &azappconfig.AddSettingOptions{
		ContentType: to.Ptr(contentType),
	})
	require.NoError(t, err)
	require.NotZero(t, addResp)

	getResp, err := client.GetSetting(context.Background(), key, nil)
	require.NoError(t, err)
	require.NotNil(t, getResp.Key)
	require.EqualValues(t, key, *getResp.Key)

	resp, err := client.DeleteSetting(context.Background(), key, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Key)
	require.EqualValues(t, key, *resp.Key)
}

func TestSnapshotListConfigurationSettings(t *testing.T) {
	snapshotName := "listConfigurationsSnapshotTest" + string(testId)
	client := NewClientFromConnectionString(t)

	type VL struct {
		Value string
		Label string
	}

	Settings := []azappconfig.Setting{
		{
			Value: to.Ptr("value3"),
			Label: to.Ptr("label"),
		},
		{
			Value: to.Ptr("Val1"),
			Label: to.Ptr("Label1"),
		},
		{
			Label: to.Ptr("Label1"),
		},
		{
			Value: to.Ptr("Val1"),
		},
		{
			Label: to.Ptr("Label2"),
		},
		{},
	}

	Keys := []string{
		"Key",
		"Key1",
		"Key2",
		"KeyNoLabel",
		"KeyNoVal",
		"NoValNoLabelKey",
	}

	require.Equal(t, len(Settings), len(Keys))

	for i, key := range Keys {
		Settings[i].Key = to.Ptr(testId + key)
	}

	settingMap := make(map[string][]VL)

	for _, setting := range Settings {

		key := *setting.Key
		value := setting.Value
		label := setting.Label

		// Add setting to Map
		mapV := VL{}

		if value != nil {
			mapV.Value = *value
		}

		if label != nil {
			mapV.Label = *label
		}

		settingMap[key] = append(settingMap[key], mapV)

		_, err := client.AddSetting(context.Background(), key, value, nil)

		require.NoError(t, err)
	}

	keyFilter := testId + "*"
	sf := []azappconfig.SettingFilter{
		{
			KeyFilter: &keyFilter,
		},
	}

	_, err := CreateSnapshot(client, snapshotName, sf)
	require.NoError(t, err)

	respPgr := client.NewListSettingsForSnapshotPager(snapshotName, nil)
	require.NotEmpty(t, respPgr)

	settingsAdded := 0

	for respPgr.More() {
		page, err := respPgr.NextPage(context.Background())

		require.NoError(t, err)
		require.NotEmpty(t, page)

		for _, setting := range page.Settings {
			require.NotNil(t, setting.Key)
			found := false

			// Check if setting is in the map
			for _, configuration := range settingMap[*setting.Key] {
				if setting.Value != nil {
					if *setting.Value != configuration.Value {
						continue
					}
				}

				if setting.Label != nil {
					if *setting.Label != configuration.Label {
						continue
					}
				}

				found = true
				settingsAdded++
				break
			}

			// Check that the key follows the filtering pattern
			if !found {
				require.True(t, strings.HasPrefix(*setting.Key, keyFilter[:len(keyFilter)-1]))
			}
		}
	}

	require.Equal(t, len(settingMap), settingsAdded)

	// Cleanup Settings
	for _, setting := range Settings {
		_, err = client.DeleteSetting(context.Background(), *setting.Key, nil)
		require.NoError(t, err)
	}

	// Cleanup Snapshots
	require.NoError(t, CleanupSnapshot(client, snapshotName))
}

func TestGetSnapshots(t *testing.T) {
	snapshotName := "getSnapshotsTest" + string(testId)

	const (
		ssCreateCount = 5
	)

	client := NewClientFromConnectionString(t)

	for i := 0; i < ssCreateCount; i++ {
		createSSName := snapshotName + fmt.Sprintf("%d", i)

		_, err := client.GetSnapshot(context.Background(), createSSName, nil)

		if err != nil {
			_, err = CreateSnapshot(client, createSSName, nil)
			require.NoError(t, err)
		}
	}

	// Get Snapshots
	ssPgr := client.NewListSnapshotsPager(nil)

	require.NotEmpty(t, ssPgr)

	snapshotCount := 0

	for ssPgr.More() {
		page, err := ssPgr.NextPage(context.Background())

		require.NoError(t, err)
		require.NotEmpty(t, page)

		for _, snapshot := range page.Snapshots {
			if strings.HasPrefix(*snapshot.Name, snapshotName) {
				snapshotCount++
			}
		}
	}

	require.Equal(t, ssCreateCount, snapshotCount)

	// Cleanup Snapshots
	for i := 0; i < ssCreateCount; i++ {
		cleanSSName := snapshotName + fmt.Sprintf("%d", i)
		require.NoError(t, CleanupSnapshot(client, cleanSSName))
	}
}

func TestSnapshotArchive(t *testing.T) {
	snapshotName := "archiveSnapshotsTest" + string(testId)

	client := NewClientFromConnectionString(t)

	snapshot, err := CreateSnapshot(client, snapshotName, nil)
	require.NoError(t, err)

	// Snapshot must exist
	_, err = client.GetSnapshot(context.Background(), snapshotName, nil)
	require.NoError(t, err)
	require.Equal(t, azappconfig.SnapshotStatusReady, *snapshot.Status)

	// Archive the snapshot
	archiveSnapshot, err := client.ArchiveSnapshot(context.Background(), snapshotName, nil)
	require.NoError(t, err)
	require.Equal(t, azappconfig.SnapshotStatusArchived, *archiveSnapshot.Snapshot.Status)

	// Best effort snapshot cleanup
	require.NoError(t, CleanupSnapshot(client, snapshotName))
}

func TestSnapshotRecover(t *testing.T) {
	snapshotName := "recoverSnapshotsTest" + string(testId)

	client := NewClientFromConnectionString(t)

	snapshot, err := CreateSnapshot(client, snapshotName, nil)
	require.NoError(t, err)

	_, err = client.GetSnapshot(context.Background(), snapshotName, nil)
	require.NoError(t, err)

	_, err = client.ArchiveSnapshot(context.Background(), snapshotName, nil)
	require.NoError(t, err)

	// Check that snapshot is archived
	archivedSnapshot, err := client.GetSnapshot(context.Background(), *snapshot.Name, nil)
	require.NoError(t, err)
	require.Equal(t, azappconfig.SnapshotStatusArchived, *archivedSnapshot.Snapshot.Status)

	// Recover the snapshot
	readySnapshot, err := client.RecoverSnapshot(context.Background(), *snapshot.Name, nil)
	require.NoError(t, err)
	require.Equal(t, azappconfig.SnapshotStatusReady, *readySnapshot.Snapshot.Status)

	// Best effort snapshot cleanup
	require.NoError(t, CleanupSnapshot(client, snapshotName))
}

func TestSnapshotCreate(t *testing.T) {
	snapshotName := "createSnapshotsTest" + string(testId)

	client := NewClientFromConnectionString(t)

	// Create a snapshot
	snapshot, err := CreateSnapshot(client, snapshotName, nil)

	require.NoError(t, err)
	require.Equal(t, snapshotName, *snapshot.Name)

	// Best effort cleanup snapshot
	require.NoError(t, CleanupSnapshot(client, snapshotName))
}

func createMultipleKeys(t *testing.T, client *azappconfig.Client, batchKey string, count int) string {
	resp, err := client.GetSetting(context.Background(), batchKey, nil)
	if err == nil {
		return *resp.Value
	}

	key, err := recording.GenerateAlphaNumericID(t, "key-", 10, true)
	require.NoError(t, err)

	for i := 0; i < count; i++ {
		_, err = client.AddSetting(context.Background(), key, to.Ptr("test_value"), &azappconfig.AddSettingOptions{
			Label: to.Ptr(fmt.Sprintf("%d", i)),
		})
		require.NoError(t, err)
	}
	_, err = client.SetSetting(context.Background(), batchKey, &key, nil)
	require.NoError(t, err)
	return key
}

func TestListSettingsPagerWithETagUnmodifiedPage(t *testing.T) {
	client := NewClientFromConnectionString(t)

	key := createMultipleKeys(t, client, "TestListSettingsPagerWithETagUnmodifiedPage", 105)

	selector := azappconfig.SettingSelector{
		KeyFilter: &key,
	}

	// get all page ETags
	pager := client.NewListSettingsPager(selector, nil)
	matchConditions := []azcore.MatchConditions{}
	countPages := 0
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		matchConditions = append(matchConditions, azcore.MatchConditions{
			IfNoneMatch: page.ETag,
		})
		countPages++
	}
	require.EqualValues(t, 2, countPages)

	// validate all pages are not modified and returns an empty list of settings
	countPages = 0
	pager = client.NewListSettingsPager(selector, &azappconfig.ListSettingsOptions{
		MatchConditions: matchConditions,
	})
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		require.Empty(t, page.Settings)
		countPages++
	}
	require.EqualValues(t, 2, countPages)
}

func TestListSettingsPagerWithETagModifiedPage(t *testing.T) {
	client := NewClientFromConnectionString(t)

	key := createMultipleKeys(t, client, "TestListSettingsPagerWithETagModifiedPage", 105)

	selector := azappconfig.SettingSelector{
		KeyFilter: &key,
	}

	// get all page ETags
	var lastSetting azappconfig.Setting
	pager := client.NewListSettingsPager(selector, nil)
	matchConditions := []azcore.MatchConditions{}
	countPages := 0
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		for _, setting := range page.Settings {
			lastSetting = setting
		}
		matchConditions = append(matchConditions, azcore.MatchConditions{
			IfNoneMatch: page.ETag,
		})
		countPages++
	}
	require.EqualValues(t, 2, countPages)

	// modify the last setting
	require.NotNil(t, lastSetting.Key)
	require.NotNil(t, lastSetting.Value)
	lastSetting.Value = to.Ptr(fmt.Sprintf("%s-1", *lastSetting.Value))
	_, err := client.SetSetting(context.Background(), *lastSetting.Key, lastSetting.Value, &azappconfig.SetSettingOptions{
		Label: lastSetting.Label,
	})
	require.NoError(t, err)

	// validate second page is modified
	countPages = 0
	pager = client.NewListSettingsPager(selector, &azappconfig.ListSettingsOptions{
		MatchConditions: matchConditions,
	})
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		require.NoError(t, err)
		if countPages == 0 {
			require.Empty(t, page.Settings)
		} else {
			require.NotEmpty(t, page.Settings)
		}
		countPages++
	}
	require.EqualValues(t, 2, countPages)
}

func CreateSnapshot(c *azappconfig.Client, snapshotName string, sf []azappconfig.SettingFilter) (azappconfig.CreateSnapshotResponse, error) {
	if sf == nil {
		all := "*"
		sf = []azappconfig.SettingFilter{
			{
				KeyFilter: &all,
			},
		}
	}

	opts := &azappconfig.BeginCreateSnapshotOptions{
		RetentionPeriod: to.Ptr[int64](3600),
	}

	// Create a snapshot
	resp, err := c.BeginCreateSnapshot(context.Background(), snapshotName, sf, opts)

	if err != nil {
		return azappconfig.CreateSnapshotResponse{}, err
	}

	snapshot, err := resp.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{
		Frequency: 1 * time.Second,
	})

	if err != nil {
		return azappconfig.CreateSnapshotResponse{}, err
	}

	// Check if snapshot exists. If not fail the test
	_, err = c.GetSnapshot(context.Background(), snapshotName, nil)

	if err != nil {
		return azappconfig.CreateSnapshotResponse{}, err
	}

	if snapshotName != *snapshot.Name {
		return azappconfig.CreateSnapshotResponse{}, fmt.Errorf("Snapshot name does not match")
	}

	return snapshot, nil
}

func CleanupSnapshot(client *azappconfig.Client, snapshotName string) error {
	_, err := client.ArchiveSnapshot(context.Background(), snapshotName, nil)

	if err != nil {
		return err
	}

	// Check if snapshot exists
	snapshot, err := client.GetSnapshot(context.Background(), snapshotName, nil)

	if err != nil || *snapshot.Status != azappconfig.SnapshotStatusArchived {
		return fmt.Errorf("Snapshot still exists")
	}

	return nil
}

func TestAddSettingWithTags(t *testing.T) {
	const (
		key   = "key-TestAddSettingWithTags"
		label = "label"
	)

	value := "value"
	tags := map[string]*string{
		"environment": to.Ptr("test"),
		"region":      to.Ptr("eastus"),
		"component":   to.Ptr("api"),
	}
	client := NewClientFromConnectionString(t)

	// Test AddSetting with tags
	addResp, err := client.AddSetting(context.Background(), key, &value, &azappconfig.AddSettingOptions{
		Label: to.Ptr(label),
		Tags:  tags,
	})
	require.NoError(t, err)
	require.NotEmpty(t, addResp)
	require.NotNil(t, addResp.Key)
	require.NotNil(t, addResp.Label)
	require.NotNil(t, addResp.Value)
	require.Equal(t, key, *addResp.Key)
	require.Equal(t, label, *addResp.Label)
	require.Equal(t, value, *addResp.Value)
	require.Equal(t, tags, addResp.Tags)

	// Verify tags are persisted
	getResp, err := client.GetSetting(context.Background(), key, &azappconfig.GetSettingOptions{
		Label: to.Ptr(label),
	})
	require.NoError(t, err)
	require.Equal(t, tags, getResp.Tags)

	// Clean up
	_, err = client.DeleteSetting(context.Background(), key, &azappconfig.DeleteSettingOptions{
		Label: to.Ptr(label),
	})
	require.NoError(t, err)
}

func TestSetSettingWithTags(t *testing.T) {
	const (
		key   = "key-TestSetSettingWithTags1"
		label = "label"
	)

	value := "value1"
	tags := map[string]*string{
		"environment": to.Ptr("prod"),
		"region":      to.Ptr("westus"),
		"version":     to.Ptr("1.0"),
	}
	client := NewClientFromConnectionString(t)

	// Test SetSetting with tags
	setResp, err := client.SetSetting(context.Background(), key, &value, &azappconfig.SetSettingOptions{
		Label: to.Ptr(label),
		Tags:  tags,
	})
	require.NoError(t, err)
	require.NotEmpty(t, setResp)
	require.NotNil(t, setResp.Key)
	require.NotNil(t, setResp.Label)
	require.NotNil(t, setResp.Value)
	require.Equal(t, key, *setResp.Key)
	require.Equal(t, label, *setResp.Label)
	require.Equal(t, value, *setResp.Value)
	require.Equal(t, tags, setResp.Tags)

	// Update with different tags
	newTags := map[string]*string{
		"environment": to.Ptr("staging"),
		"region":      to.Ptr("centralus"),
		"version":     to.Ptr("2.0"),
		"feature":     to.Ptr("beta"),
	}
	newValue := "updated-value"

	setResp2, err := client.SetSetting(context.Background(), key, &newValue, &azappconfig.SetSettingOptions{
		Label: to.Ptr(label),
		Tags:  newTags,
	})
	require.NoError(t, err)
	require.Equal(t, newTags, setResp2.Tags)
	require.Equal(t, newValue, *setResp2.Value)

	// Clean up
	_, err = client.DeleteSetting(context.Background(), key, &azappconfig.DeleteSettingOptions{
		Label: to.Ptr(label),
	})
	require.NoError(t, err)
}

func TestListSettingsWithTagsFilter(t *testing.T) {
	const (
		key1 = "key1-TestListSettingsWithTagsFilter"
		key2 = "key2-TestListSettingsWithTagsFilter"
		key3 = "key3-TestListSettingsWithTagsFilter"
	)

	client := NewClientFromConnectionString(t)

	// Create settings with different tags
	settings := []struct {
		key   string
		value string
		tags  map[string]*string
	}{
		{
			key:   key1,
			value: "value1",
			tags: map[string]*string{
				"environment": to.Ptr("prod"),
				"region":      to.Ptr("eastus"),
			},
		},
		{
			key:   key2,
			value: "value2",
			tags: map[string]*string{
				"environment": to.Ptr("test"),
				"region":      to.Ptr("eastus"),
			},
		},
		{
			key:   key3,
			value: "value3",
			tags: map[string]*string{
				"environment": to.Ptr("prod"),
				"region":      to.Ptr("westus"),
			},
		},
	}

	// Add all settings
	for _, setting := range settings {
		_, err := client.SetSetting(context.Background(), setting.key, &setting.value, &azappconfig.SetSettingOptions{
			Tags:  setting.tags,
			Label: to.Ptr("tag-filter-test"),
		})
		require.NoError(t, err)
	}

	// Test filtering by single tag
	selector := azappconfig.SettingSelector{
		TagsFilter: []string{"environment=prod"},
	}

	var foundKeys []string
	settsPgr := client.NewListSettingsPager(selector, nil)
	for settsPgr.More() {
		page, err := settsPgr.NextPage(context.Background())
		require.NoError(t, err)

		for _, setting := range page.Settings {
			if *setting.Key == key1 || *setting.Key == key2 || *setting.Key == key3 {
				foundKeys = append(foundKeys, *setting.Key)
				// Verify the setting has the expected tag
				require.Equal(t, "prod", *setting.Tags["environment"])
			}
		}
	}

	// Should find key1 and key3 (both have environment=prod)
	require.Contains(t, foundKeys, key1)
	require.Contains(t, foundKeys, key3)
	require.NotContains(t, foundKeys, key2) // key2 has environment=test

	// Test filtering by multiple tags (AND operation)
	selectorMultiple := azappconfig.SettingSelector{
		TagsFilter: []string{"environment=prod", "region=eastus"},
	}

	foundKeysMultiple := []string{}
	settsPgrMultiple := client.NewListSettingsPager(selectorMultiple, nil)
	for settsPgrMultiple.More() {
		page, err := settsPgrMultiple.NextPage(context.Background())
		require.NoError(t, err)

		for _, setting := range page.Settings {
			if *setting.Key == key1 || *setting.Key == key2 || *setting.Key == key3 {
				foundKeysMultiple = append(foundKeysMultiple, *setting.Key)
				// Verify the setting has both expected tags
				require.Equal(t, "prod", *setting.Tags["environment"])
				require.Equal(t, "eastus", *setting.Tags["region"])
			}
		}
	}

	// Should find only key1 (has both environment=prod AND region=eastus)
	require.Contains(t, foundKeysMultiple, key1)
	require.NotContains(t, foundKeysMultiple, key2) // key2 has environment=test
	require.NotContains(t, foundKeysMultiple, key3) // key3 has region=westus

	// Clean up
	for _, setting := range settings {
		_, err := client.DeleteSetting(context.Background(), setting.key, &azappconfig.DeleteSettingOptions{
			Label: to.Ptr("tag-filter-test"),
		})
		require.NoError(t, err)
	}
}

func TestSettingTagsWithSpecialCharacters(t *testing.T) {
	const (
		key = "key-TestSettingTagsWithSpecialCharacters"
	)

	client := NewClientFromConnectionString(t)

	// Test tags with special characters
	tags := map[string]*string{
		"app&name":    to.Ptr("my-app"),
		"team.owner":  to.Ptr("backend-team"),
		"cost_center": to.Ptr("engineering"),
		"env?type":    to.Ptr("staging-test"),
	}
	value := "value"

	addResp, err := client.SetSetting(context.Background(), key, &value, &azappconfig.SetSettingOptions{
		Tags:  tags,
		Label: to.Ptr("special-char-test"),
	})
	require.NoError(t, err)
	require.Equal(t, tags, addResp.Tags)

	// Test filtering with special characters
	selector := azappconfig.SettingSelector{
		TagsFilter: []string{"app&name=my-app"},
	}

	found := false
	settsPgr := client.NewListSettingsPager(selector, nil)
	for settsPgr.More() {
		page, err := settsPgr.NextPage(context.Background())
		require.NoError(t, err)

		for _, setting := range page.Settings {
			if *setting.Key == key {
				require.Equal(t, tags, setting.Tags)
				found = true
			}
		}
	}
	require.True(t, found, "Setting with special character tags not found")

	// Clean up
	_, err = client.DeleteSetting(context.Background(), key, &azappconfig.DeleteSettingOptions{
		Label: to.Ptr("special-char-test"),
	})
	require.NoError(t, err)
}

func TestEmptyTagsHandling(t *testing.T) {
	const (
		key = "key-TestEmptyTagsHandling"
	)

	client := NewClientFromConnectionString(t)
	value := "value"

	// Test setting with nil tags
	addResp, err := client.SetSetting(context.Background(), key, &value, &azappconfig.SetSettingOptions{
		Tags: nil,
	})
	require.NoError(t, err)
	require.Empty(t, addResp.Tags)

	// Test setting with empty tags map
	setResp, err := client.SetSetting(context.Background(), key, &value, &azappconfig.SetSettingOptions{
		Tags: make(map[string]*string),
	})
	require.NoError(t, err)
	require.Empty(t, setResp.Tags)

	// Clean up
	_, err = client.DeleteSetting(context.Background(), key, nil)
	require.NoError(t, err)
}

func TestTagsFilteringEndToEnd(t *testing.T) {
	client := NewClientFromConnectionString(t)

	// Create a comprehensive test scenario with multiple settings and different tag combinations
	testData := []struct {
		key   string
		value string
		label string
		tags  map[string]*string
	}{
		{
			key:   "app-config-1",
			value: "value1",
			label: "production",
			tags: map[string]*string{
				"environment": to.Ptr("production"),
				"service":     to.Ptr("api"),
				"version":     to.Ptr("1.0"),
			},
		},
		{
			key:   "app-config-2",
			value: "value2",
			label: "staging",
			tags: map[string]*string{
				"environment": to.Ptr("staging"),
				"service":     to.Ptr("api"),
				"version":     to.Ptr("1.1"),
			},
		},
		{
			key:   "app-config-3",
			value: "value3",
			label: "production",
			tags: map[string]*string{
				"environment": to.Ptr("production"),
				"service":     to.Ptr("web"),
				"version":     to.Ptr("2.0"),
			},
		},
		{
			key:   "app-config-4",
			value: "value4",
			label: "development",
			tags: map[string]*string{
				"environment": to.Ptr("development"),
				"service":     to.Ptr("api"),
				"version":     to.Ptr("1.0"),
			},
		},
	}

	// Add all test settings
	for _, data := range testData {
		_, err := client.SetSetting(context.Background(), data.key, &data.value, &azappconfig.SetSettingOptions{
			Label: to.Ptr(data.label),
			Tags:  data.tags,
		})
		require.NoError(t, err)
	}

	// Test 1: Filter by single tag (environment=production)
	t.Run("SingleTagFilter", func(t *testing.T) {
		selector := azappconfig.SettingSelector{
			TagsFilter: []string{"environment=production"},
		}

		var foundKeys []string
		pager := client.NewListSettingsPager(selector, nil)
		for pager.More() {
			page, err := pager.NextPage(context.Background())
			require.NoError(t, err)

			for _, setting := range page.Settings {
				if strings.HasPrefix(*setting.Key, "app-config-") {
					foundKeys = append(foundKeys, *setting.Key)
				}
			}
		}

		// Should find app-config-1 and app-config-3
		require.Contains(t, foundKeys, "app-config-1")
		require.Contains(t, foundKeys, "app-config-3")
		require.NotContains(t, foundKeys, "app-config-2")
		require.NotContains(t, foundKeys, "app-config-4")
	})

	// Test 2: Filter by multiple tags (service=api AND version=1.0)
	t.Run("MultipleTagsFilter", func(t *testing.T) {
		selector := azappconfig.SettingSelector{
			TagsFilter: []string{"service=api", "version=1.0"},
		}

		var foundKeys []string
		pager := client.NewListSettingsPager(selector, nil)
		for pager.More() {
			page, err := pager.NextPage(context.Background())
			require.NoError(t, err)

			for _, setting := range page.Settings {
				if strings.HasPrefix(*setting.Key, "app-config-") {
					foundKeys = append(foundKeys, *setting.Key)
				}
			}
		}

		// Should find only app-config-1 and app-config-4
		require.Contains(t, foundKeys, "app-config-1")
		require.Contains(t, foundKeys, "app-config-4")
		require.NotContains(t, foundKeys, "app-config-2")
		require.NotContains(t, foundKeys, "app-config-3")
	})

	// Test 3: Combine TagsFilter with other filters
	t.Run("CombinedFilters", func(t *testing.T) {
		productionLabel := "production"
		selector := azappconfig.SettingSelector{
			LabelFilter: &productionLabel,
			TagsFilter:  []string{"service=api"},
		}

		var foundKeys []string
		pager := client.NewListSettingsPager(selector, nil)
		for pager.More() {
			page, err := pager.NextPage(context.Background())
			require.NoError(t, err)

			for _, setting := range page.Settings {
				if strings.HasPrefix(*setting.Key, "app-config-") {
					foundKeys = append(foundKeys, *setting.Key)
					// Verify the setting matches both filters
					require.Equal(t, "production", *setting.Label)
					require.Equal(t, "api", *setting.Tags["service"])
				}
			}
		}

		// Should find only app-config-1 (has both label=production AND service=api)
		require.Contains(t, foundKeys, "app-config-1")
		require.NotContains(t, foundKeys, "app-config-2") // staging label
		require.NotContains(t, foundKeys, "app-config-3") // service=web
		require.NotContains(t, foundKeys, "app-config-4") // development label
	})

	// Test 4: Test with revisions pager
	t.Run("RevisionsWithTagsFilter", func(t *testing.T) {
		// Update one of the settings to create a revision
		newValue := "updated-value"
		_, err := client.SetSetting(context.Background(), "app-config-1", &newValue, &azappconfig.SetSettingOptions{
			Label: to.Ptr("production"),
			Tags: map[string]*string{
				"environment": to.Ptr("production"),
				"service":     to.Ptr("api"),
				"version":     to.Ptr("1.1"), // Changed version
			},
		})
		require.NoError(t, err)

		// Filter revisions by tags
		keyFilter := "app-config-1"
		selector := azappconfig.SettingSelector{
			KeyFilter:  &keyFilter,
			TagsFilter: []string{"service=api"},
		}

		var revisions []azappconfig.Setting
		pager := client.NewListRevisionsPager(selector, nil)
		for pager.More() {
			page, err := pager.NextPage(context.Background())
			require.NoError(t, err)

			revisions = append(revisions, page.Settings...)
		}

		// Should find both revisions (original and updated) with service=api
		require.GreaterOrEqual(t, len(revisions), 2)
		for _, revision := range revisions {
			if *revision.Key == "app-config-1" {
				require.Equal(t, "api", *revision.Tags["service"])
			}
		}
	})

	// Clean up all test settings
	for _, data := range testData {
		_, err := client.DeleteSetting(context.Background(), data.key, &azappconfig.DeleteSettingOptions{
			Label: to.Ptr(data.label),
		})
		require.NoError(t, err)
	}
}

func TestTagsFilterMultipleAndNilValues(t *testing.T) {
	client := NewClientFromConnectionString(t)

	// Test data with various tag scenarios including nil and empty values
	testData := []struct {
		key   string
		value string
		label string
		tags  map[string]*string
	}{
		{
			key:   "setting-with-normal-tags",
			value: "value1",
			label: "prod",
			tags: map[string]*string{
				"env":     to.Ptr("production"),
				"service": to.Ptr("api"),
				"version": to.Ptr("1.0"),
			},
		},
		{
			key:   "setting-with-nil-tag-value",
			value: "value2",
			label: "prod",
			tags: map[string]*string{
				"env":     to.Ptr("production"),
				"service": nil, // nil tag value
				"version": to.Ptr("1.0"),
			},
		},
		{
			key:   "setting-with-empty-tag-value",
			value: "value3",
			label: "prod",
			tags: map[string]*string{
				"env":     to.Ptr("production"),
				"service": to.Ptr(""), // empty string tag value
				"version": to.Ptr("1.0"),
			},
		},
		{
			key:   "setting-mixed-tags",
			value: "value4",
			label: "staging",
			tags: map[string]*string{
				"env":       to.Ptr("staging"),
				"service":   to.Ptr("web"),
				"version":   to.Ptr("2.0"),
				"feature":   nil,                // nil value
				"debug":     to.Ptr(""),         // empty value
				"component": to.Ptr("frontend"), // normal value
			},
		},
		{
			key:   "setting-only-nil-tags",
			value: "value5",
			label: "test",
			tags: map[string]*string{
				"tag1": nil,
				"tag2": nil,
				"tag3": nil,
			},
		},
		{
			key:   "setting-only-empty-tags",
			value: "value6",
			label: "test",
			tags: map[string]*string{
				"tag1": to.Ptr(""),
				"tag2": to.Ptr(""),
				"tag3": to.Ptr(""),
			},
		},
	}

	// Create all test settings
	for _, data := range testData {
		_, err := client.AddSetting(context.Background(), data.key, &data.value, &azappconfig.AddSettingOptions{
			Label: to.Ptr(data.label),
			Tags:  data.tags,
		})
		require.NoError(t, err)
	}

	// Test 1: Filter by multiple specific tag values
	t.Run("MultipleTagsFilter", func(t *testing.T) {
		selector := azappconfig.SettingSelector{
			TagsFilter: []string{"env=production", "version=1.0"},
		}

		var foundKeys []string
		pager := client.NewListSettingsPager(selector, nil)
		for pager.More() {
			page, err := pager.NextPage(context.Background())
			require.NoError(t, err)

			for _, setting := range page.Settings {
				if strings.HasPrefix(*setting.Key, "setting-") {
					foundKeys = append(foundKeys, *setting.Key)
				}
			}
		}

		// Should find settings with both env=production AND version=1.0
		require.Contains(t, foundKeys, "setting-with-normal-tags")
		require.Contains(t, foundKeys, "setting-with-nil-tag-value")
		require.Contains(t, foundKeys, "setting-with-empty-tag-value")
		require.NotContains(t, foundKeys, "setting-mixed-tags") // has env=staging
	})

	// Test 2: Filter by exact tag key-value match
	t.Run("TagKeyExistsFilter", func(t *testing.T) {
		selector := azappconfig.SettingSelector{
			TagsFilter: []string{"service=api"},
		}

		var foundKeys []string
		pager := client.NewListSettingsPager(selector, nil)
		for pager.More() {
			page, err := pager.NextPage(context.Background())
			require.NoError(t, err)

			for _, setting := range page.Settings {
				if strings.HasPrefix(*setting.Key, "setting-") {
					foundKeys = append(foundKeys, *setting.Key)
				}
			}
		}

		// Should find only the setting with service=api
		require.Len(t, foundKeys, 1)
		require.Contains(t, foundKeys, "setting-with-normal-tags")
		require.NotContains(t, foundKeys, "setting-with-nil-tag-value")   // service is nil
		require.NotContains(t, foundKeys, "setting-with-empty-tag-value") // service is empty
	})

	// Test 3: Filter by empty string tag value
	t.Run("EmptyStringTagValueFilter", func(t *testing.T) {
		selector := azappconfig.SettingSelector{
			TagsFilter: []string{"service="},
		}

		var foundKeys []string
		pager := client.NewListSettingsPager(selector, nil)
		for pager.More() {
			page, err := pager.NextPage(context.Background())
			require.NoError(t, err)

			for _, setting := range page.Settings {
				if strings.HasPrefix(*setting.Key, "setting-") {
					foundKeys = append(foundKeys, *setting.Key)
				}
			}
		}

		// Should find settings where service tag has empty string value
		require.Contains(t, foundKeys, "setting-with-empty-tag-value")
		require.NotContains(t, foundKeys, "setting-with-normal-tags")   // service=api
		require.NotContains(t, foundKeys, "setting-with-nil-tag-value") // service is nil
	})

	// Test 4: Complex filter with multiple conditions
	t.Run("ComplexMultipleTagsFilter", func(t *testing.T) {
		selector := azappconfig.SettingSelector{
			TagsFilter: []string{"env=staging", "service=web", "version=2.0"},
		}

		var foundKeys []string
		pager := client.NewListSettingsPager(selector, nil)
		for pager.More() {
			page, err := pager.NextPage(context.Background())
			require.NoError(t, err)

			for _, setting := range page.Settings {
				if strings.HasPrefix(*setting.Key, "setting-") {
					foundKeys = append(foundKeys, *setting.Key)
				}
			}
		}

		// Should find only setting-mixed-tags that matches all conditions
		require.Contains(t, foundKeys, "setting-mixed-tags")
		require.Len(t, foundKeys, 1) // Only one setting should match all conditions
	})

	// Test 6: Filter with non-existent tag
	t.Run("NonExistentTagFilter", func(t *testing.T) {
		selector := azappconfig.SettingSelector{
			TagsFilter: []string{"nonexistent=value"},
		}

		var foundKeys []string
		pager := client.NewListSettingsPager(selector, nil)
		for pager.More() {
			page, err := pager.NextPage(context.Background())
			require.NoError(t, err)

			for _, setting := range page.Settings {
				if strings.HasPrefix(*setting.Key, "setting-") {
					foundKeys = append(foundKeys, *setting.Key)
				}
			}
		}

		// Should find no settings with non-existent tag
		require.Empty(t, foundKeys)
	})

	// Test 7: Filter by tag key that exists with nil values
	t.Run("TagKeyExistsWithNilValuesFilter", func(t *testing.T) {
		selector := azappconfig.SettingSelector{
			TagsFilter: []string{"tag1=\x00"}, // tag1 exists with nil value
		}
		var foundKeys []string
		pager := client.NewListSettingsPager(selector, nil)
		for pager.More() {
			page, err := pager.NextPage(context.Background())
			require.NoError(t, err)
			for _, setting := range page.Settings {
				if strings.HasPrefix(*setting.Key, "setting-") {
					foundKeys = append(foundKeys, *setting.Key)
				}
			}
		}

		// Should find settings that have tag1 with nil values
		require.Contains(t, foundKeys, "setting-only-nil-tags")
		require.Len(t, foundKeys, 1) // Only two settings should match
	})

	// Clean up all test settings
	for _, data := range testData {
		_, err := client.DeleteSetting(context.Background(), data.key, &azappconfig.DeleteSettingOptions{
			Label: to.Ptr(data.label),
		})
		require.NoError(t, err)
	}
}
