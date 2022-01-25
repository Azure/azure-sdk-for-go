// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type listBlobPerfTest struct {
	perf.PerfTestOptions
	containerName   string
	blobName        string
	containerClient azblob.ContainerClient
}

func (m *listBlobPerfTest) GlobalSetup(ctx context.Context) error {
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, m.containerName, nil)
	if err != nil {
		return err
	}
	_, err = containerClient.Create(context.Background(), nil)
	if err != nil {
		return err
	}

	for i := 0; i < 100; i++ {
		blobClient := containerClient.NewBlockBlobClient(fmt.Sprintf("%s%d", m.blobName, i))
		_, err = blobClient.Upload(
			context.Background(),
			NopCloser(bytes.NewReader([]byte(""))),
			nil,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *listBlobPerfTest) Setup(ctx context.Context) error {
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, m.containerName, nil)
	m.containerClient = containerClient
	return err
}

func (m *listBlobPerfTest) Run(ctx context.Context) error {
	pager := m.containerClient.ListBlobsFlat(&azblob.ContainerListBlobFlatSegmentOptions{Maxresults: to.Int32Ptr(int32(*count))})
	for pager.NextPage(context.Background()) {
	}
	return pager.Err()
}

func (m *listBlobPerfTest) Cleanup(ctx context.Context) error {
	return nil
}

func (m *listBlobPerfTest) GlobalCleanup(ctx context.Context) error {

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, m.containerName, nil)
	if err != nil {
		return err
	}

	_, err = containerClient.Delete(context.Background(), nil)
	return err
}

func (m *listBlobPerfTest) GetMetadata() perf.PerfTestOptions {
	return m.PerfTestOptions
}

func (l *listBlobPerfTest) RegisterArguments() error {
	return nil
}

func NewListTest(options *perf.PerfTestOptions) perf.PerfTest {
	if options == nil {
		options = &perf.PerfTestOptions{}
	}
	options.Name = "BlobListTest"
	return &listBlobPerfTest{
		PerfTestOptions: *options,
		blobName:        "listTest",
		containerName:   "listtest",
	}
}
