// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestGetAccountInfo(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	// Ensure the call succeeded. Don't test for specific account properties because we can't/don't want to set account properties.
	sAccInfo, err := svcClient.GetAccountInfo(context.Background())
	require.NoError(t, err)
	require.NotEqualValues(t, sAccInfo, ServiceGetAccountInfoResponse{})
}

func TestServiceClientFromConnectionString(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()

	svcClient, err := createServiceClientFromConnectionString(t, testAccountDefault)
	require.NoError(t, err)
	containerClient := createNewContainer(t, generateContainerName(testName), svcClient)
	defer deleteContainer(t, containerClient)
}

func TestListContainersBasic(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	md := map[string]string{
		"foo": "foovalue",
		"bar": "barvalue",
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)
	_, err = containerClient.Create(ctx, &CreateContainerOptions{Metadata: md})
	defer func(containerClient ContainerClient, ctx context.Context, options *DeleteContainerOptions) {
		_, err := containerClient.Delete(ctx, options)
		if err != nil {
			require.NoError(t, err)
		}
	}(containerClient, ctx, nil)
	require.NoError(t, err)
	prefix := containerPrefix
	listOptions := ListContainersOptions{Prefix: &prefix, Include: ListContainersDetail{Metadata: true}}
	pager := svcClient.ListContainers(&listOptions)

	count := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, container := range resp.ServiceListContainersSegmentResult.ContainerItems {
			require.NotNil(t, container.Name)

			if *container.Name == containerName {
				require.NotNil(t, container.Properties)
				require.NotNil(t, container.Properties.LastModified)
				require.NotNil(t, container.Properties.Etag)
				require.Equal(t, *container.Properties.LeaseStatus, LeaseStatusTypeUnlocked)
				require.Equal(t, *container.Properties.LeaseState, LeaseStateTypeAvailable)
				require.Nil(t, container.Properties.LeaseDuration)
				require.Nil(t, container.Properties.PublicAccess)
				require.NotNil(t, container.Metadata)

				unwrappedMeta := map[string]string{}
				for k, v := range container.Metadata {
					if v != nil {
						unwrappedMeta[k] = *v
					}
				}

				require.EqualValues(t, unwrappedMeta, md)
			}
		}
	}

	require.Nil(t, pager.Err())
	require.NoError(t, err)
	require.GreaterOrEqual(t, count, 0)
}

func TestListContainersBasicUsingConnectionString(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClientFromConnectionString(t, testAccountDefault)
	require.NoError(t, err)

	md := map[string]string{
		"foo": "foovalue",
		"bar": "barvalue",
	}

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)
	_, err = containerClient.Create(ctx, &CreateContainerOptions{Metadata: md})
	defer func(containerClient ContainerClient, ctx context.Context, options *DeleteContainerOptions) {
		_, err := containerClient.Delete(ctx, options)
		require.NoError(t, err)
	}(containerClient, ctx, nil)
	require.NoError(t, err)
	prefix := containerPrefix
	listOptions := ListContainersOptions{Prefix: &prefix, Include: ListContainersDetail{Metadata: true}}
	pager := svcClient.ListContainers(&listOptions)

	count := 0
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, container := range resp.ServiceListContainersSegmentResult.ContainerItems {
			require.NotNil(t, container.Name)

			if *container.Name == containerName {
				require.NotNil(t, container.Properties)
				require.NotNil(t, container.Properties.LastModified)
				require.NotNil(t, container.Properties.Etag)
				require.Equal(t, *container.Properties.LeaseStatus, LeaseStatusTypeUnlocked)
				require.Equal(t, *container.Properties.LeaseState, LeaseStateTypeAvailable)
				require.Nil(t, container.Properties.LeaseDuration)
				require.Nil(t, container.Properties.PublicAccess)
				require.NotNil(t, container.Metadata)

				unwrappedMeta := map[string]string{}
				for k, v := range container.Metadata {
					if v != nil {
						unwrappedMeta[k] = *v
					}
				}

				require.EqualValues(t, unwrappedMeta, md)
			}
		}
	}

	require.NoError(t, pager.Err())
	require.NoError(t, err)
	require.GreaterOrEqual(t, count, 0)
}

func TestAccountListContainersEmptyPrefix(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient1 := createNewContainer(t, generateContainerName(testName)+"1", svcClient)
	defer deleteContainer(t, containerClient1)
	containerClient2 := createNewContainer(t, generateContainerName(testName)+"2", svcClient)
	defer deleteContainer(t, containerClient2)

	count := 0
	pager := svcClient.ListContainers(nil)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, container := range resp.ServiceListContainersSegmentResult.ContainerItems {
			count++
			require.NotNil(t, container.Name)
		}
	}
	require.NoError(t, pager.Err())
	require.GreaterOrEqual(t, count, 2)
}

func TestAccountDeleteRetentionPolicy(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	days := to.Int32Ptr(5)
	enabled := to.BoolPtr(true)
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: enabled, Days: days}})
	require.NoError(t, err)

	// From FE, 30 seconds is guaranteed to be enough.
	recording.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(ctx)
	require.NoError(t, err)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	disabled := false
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &disabled}})
	require.NoError(t, err)

	// From FE, 30 seconds is guaranteed to be enough.
	recording.Sleep(time.Second * 30)

	resp, err = svcClient.GetProperties(ctx)
	require.NoError(t, err)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, false)
	require.Nil(t, resp.StorageServiceProperties.DeleteRetentionPolicy.Days)
}

func TestAccountDeleteRetentionPolicyEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	days := to.Int32Ptr(5)
	enabled := to.BoolPtr(true)
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: enabled, Days: days}})
	require.NoError(t, err)

	// From FE, 30 seconds is guaranteed to be enough.
	recording.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(ctx)
	require.NoError(t, err)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	// Empty retention policy causes an error, this is different from track 1.5
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{}})
	require.Error(t, err)
}

func TestAccountDeleteRetentionPolicyNil(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	days := to.Int32Ptr(5)
	enabled := to.BoolPtr(true)
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: enabled, Days: days}})
	require.NoError(t, err)

	// From FE, 30 seconds is guaranteed to be enough.
	recording.Sleep(time.Second * 30)

	resp, err := svcClient.GetProperties(ctx)
	require.NoError(t, err)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{})
	require.NoError(t, err)

	// From FE, 30 seconds is guaranteed to be enough.
	recording.Sleep(time.Second * 30)

	// If an element of service properties is not passed, the service keeps the current settings.
	resp, err = svcClient.GetProperties(ctx)
	require.NoError(t, err)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Enabled, *enabled)
	require.EqualValues(t, *resp.StorageServiceProperties.DeleteRetentionPolicy.Days, *days)

	// Disable for other tests
	enabled = to.BoolPtr(false)
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: enabled}})
	require.NoError(t, err)
}

func TestAccountDeleteRetentionPolicyDaysTooSmall(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	days := int32(0) // Minimum days is 1. Validated on the client.
	enabled := true
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
	require.Error(t, err)
}

func TestAccountDeleteRetentionPolicyDaysTooLarge(t *testing.T) {
	stop := start(t)
	defer stop()

	var svcClient ServiceClient
	var err error
	for i := 1; i <= 2; i++ {
		if i == 1 {
			svcClient, err = createServiceClient(t, testAccountDefault)
		} else {
			svcClient, err = createServiceClientFromConnectionString(t, testAccountDefault)
		}
		require.NoError(t, err)

		days := int32(366) // Max days is 365. Left to the service for validation.
		enabled := true
		_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled, Days: &days}})
		require.Error(t, err)

		validateStorageError(t, err, StorageErrorCodeInvalidXMLDocument)
	}
}

func TestAccountDeleteRetentionPolicyDaysOmitted(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	// Days is required if enabled is true.
	enabled := true
	_, err = svcClient.SetProperties(ctx, StorageServiceProperties{DeleteRetentionPolicy: &RetentionPolicy{Enabled: &enabled}})
	require.Error(t, err)

	validateStorageError(t, err, StorageErrorCodeInvalidXMLDocument)
}
