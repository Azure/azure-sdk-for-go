// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

// tests general functionality
func TestBlobListWrapper(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(t.Name())
	containerClient := getContainerClient(containerName, svcClient)

	_, err = containerClient.Create(ctx, nil)
	require.NoError(t, err)
	defer deleteContainer(t, containerClient)

	files := []string{"a123", "b234", "c345"}

	createNewBlobs(t, files, containerClient)

	pager := containerClient.ListBlobsFlat(nil)

	found := make([]string, 0)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blob := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			found = append(found, *blob.Name)
		}
	}
	require.Nil(t, pager.Err())

	sort.Strings(files)
	sort.Strings(found)

	require.EqualValues(t, found, files)
}

// tests that the buffer filling isn't a problem
func TestBlobListWrapperFullBuffer(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := getContainerClient(generateContainerName(t.Name()), svcClient)

	_, err = containerClient.Create(ctx, nil)
	require.NoError(t, err)
	defer deleteContainer(t, containerClient)

	files := []string{"a123", "b234", "c345"}

	createNewBlobs(t, files, containerClient)

	pager := containerClient.ListBlobsFlat(nil)

	found := make([]string, 0)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blob := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			found = append(found, *blob.Name)
		}
	}
	require.Nil(t, pager.Err())

	sort.Strings(files)
	sort.Strings(found)

	require.EqualValues(t, files, found)
}

func TestBlobListWrapperListingError(t *testing.T) {
	stop := start(t)
	defer stop()

	svcClient, err := createServiceClientWithSharedKeyForRecording(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := getContainerClient(generateContainerName(t.Name()), svcClient)

	pager := containerClient.ListBlobsFlat(nil)

	require.Equal(t, pager.NextPage(ctx), false)
	require.NotNil(t, pager.Err())
}
