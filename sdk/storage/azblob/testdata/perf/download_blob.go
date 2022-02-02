// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type downloadPerfTest struct {
	perf.PerfTestOptions
	blobName   string
	blobClient azblob.BlockBlobClient
}

func (d *downloadPerfTest) GlobalSetup(ctx context.Context) error {
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, d.blobName, nil)
	if err != nil {
		return err
	}
	_, err = containerClient.Create(context.Background(), nil)
	if err != nil {
		return err
	}

	blobClient := containerClient.NewBlockBlobClient(d.blobName)

	data, err := perf.NewRandomStream(int(*size))
	if err != nil {
		return err
	}

	_, err = blobClient.Upload(context.Background(), data, nil)

	return err
}

func (d *downloadPerfTest) Setup(ctx context.Context) error {
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, d.blobName, &azblob.ClientOptions{Transporter: d.ProxyInstance})
	if err != nil {
		return err
	}

	d.blobClient = containerClient.NewBlockBlobClient(d.blobName)

	return nil
}

func (d *downloadPerfTest) Run(ctx context.Context) error {
	get, err := d.blobClient.Download(ctx, nil)
	if err != nil {
		return err
	}
	downloadedData := &bytes.Buffer{}
	reader := get.Body(azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(reader)
	if err != nil {
		return err
	}
	return reader.Close()
}

func (d *downloadPerfTest) Cleanup(ctx context.Context) error {
	return nil
}

func (d *downloadPerfTest) GlobalCleanup(ctx context.Context) error {
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, d.blobName, nil)
	if err != nil {
		return err
	}

	_, err = containerClient.Delete(context.Background(), nil)
	return err
}

func (d *downloadPerfTest) GetMetadata() perf.PerfTestOptions {
	return d.PerfTestOptions
}

func NewDownloadTest(options *perf.PerfTestOptions) perf.PerfTest {
	if options == nil {
		options = &perf.PerfTestOptions{}
	}
	options.Name = "BlobDownloadTest"
	return &downloadPerfTest{
		PerfTestOptions: *options,
		blobName:        "downloadtest",
	}
}
