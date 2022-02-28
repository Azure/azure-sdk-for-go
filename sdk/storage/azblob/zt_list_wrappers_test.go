// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// tests general functionality
func TestBlobListWrapper(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, svcClient)

	_, err = containerClient.Create(ctx, nil)
	require.NoError(t, err)
	defer deleteContainer(_assert, containerClient)

	files := []string{"a123", "b234", "c345"}

	createNewBlobs(_assert, files, containerClient)

	pager := containerClient.ListBlobsFlat(nil)

	found := make([]string, 0)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blob := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			found = append(found, *blob.Name)
		}
	}
	require.NoError(t, pager.Err())

	sort.Strings(files)
	sort.Strings(found)

	require.EqualValues(t, found, files)
}

// tests that the buffer filling isn't a problem
func TestBlobListWrapperFullBuffer(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := getContainerClient(generateContainerName(testName), svcClient)

	_, err = containerClient.Create(ctx, nil)
	_assert.NoError(err)
	defer deleteContainer(_assert, containerClient)

	files := []string{"a123", "b234", "c345"}

	createNewBlobs(_assert, files, containerClient)

	pager := containerClient.ListBlobsFlat(nil)

	found := make([]string, 0)

	for pager.NextPage(ctx) {
		resp := pager.PageResponse()

		for _, blob := range resp.ContainerListBlobFlatSegmentResult.Segment.BlobItems {
			found = append(found, *blob.Name)
		}
	}
	require.NoError(t, pager.Err())

	sort.Strings(files)
	sort.Strings(found)

	require.EqualValues(t, files, found)
}

func TestBlobListWrapperListingError(t *testing.T) {
	stop := start(t)
	defer stop()

	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerClient := getContainerClient(generateContainerName(testName), svcClient)

	pager := containerClient.ListBlobsFlat(nil)

	require.False(t, pager.NextPage(ctx))
	require.Error(t, pager.Err())
}
