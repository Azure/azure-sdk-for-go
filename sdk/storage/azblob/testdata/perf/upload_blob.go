// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type uploadPerfTest struct {
	perf.PerfTestOptions
	containerName string
	blobName      string
	blobClient    azblob.BlockBlobClient
	data          io.ReadSeekCloser
}

func (m *uploadPerfTest) GlobalSetup(ctx context.Context) error {
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, m.containerName, nil)
	if err != nil {
		fmt.Println("Error creating the container client: ")
		return err
	}
	_, err = containerClient.Create(context.Background(), nil)
	if err != nil {
		fmt.Printf("Error creating the container: '%s'\n", m.containerName)
		return err
	}

	return nil
}

func (m *uploadPerfTest) Setup(ctx context.Context) error {
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, m.containerName, nil)
	if err != nil {
		return err
	}

	m.blobClient = containerClient.NewBlockBlobClient(m.blobName)
	return nil
}

func (m *uploadPerfTest) Run(ctx context.Context) error {
	_, err := m.blobClient.Upload(ctx, m.data, &azblob.UploadBlockBlobOptions{})
	if err != nil {
		return err
	}
	_, err = m.data.Seek(0, io.SeekStart) // rewind to the beginning
	return err
}

func (m *uploadPerfTest) Cleanup(ctx context.Context) error {
	return nil
}

func (m *uploadPerfTest) GlobalCleanup(ctx context.Context) error {
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

func (m *uploadPerfTest) GetMetadata() perf.PerfTestOptions {
	return m.PerfTestOptions
}

func NewUploadTest(options *perf.PerfTestOptions) perf.PerfTest {
	if options == nil {
		options = &perf.PerfTestOptions{}
	}
	options.Name = "BlobUploadTest"

	if size == nil {
		size = to.Int64Ptr(10240)
	}

	if count == nil {
		count = to.Int64Ptr(100)
	}
	data, err := perf.NewRandomStream(int(*size))
	if err != nil {
		panic(err)
	}
	return &uploadPerfTest{
		PerfTestOptions: *options,
		blobName:        "uploadtest",
		containerName:   "uploadcontainer",
		data:            data,
	}
}
