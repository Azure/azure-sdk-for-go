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
	blobName        string
	containerClient azblob.ContainerClient
	blobClient      azblob.BlockBlobClient
	data            string
}

func (d *downloadPerfTest) GlobalSetup(ctx context.Context) error {
	d.blobName = "downloadtest"

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZBLOB_CONNECTION_STRING' could not be found")
	}

	t, err := perf.NewProxyTransport(&perf.TransportOptions{TestName: d.GetMetadata().Name})
	if err != nil {
		return err
	}
	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, "downloadtest", &azblob.ClientOptions{Transporter: t})
	if err != nil {
		return err
	}
	d.containerClient = containerClient
	_, err = d.containerClient.Create(context.Background(), nil)
	if err != nil {
		return err
	}

	d.blobClient = d.containerClient.NewBlockBlobClient(d.blobName)

	d.data = "This is all placeholder random data for now. This is all placeholder random data for now. This is all placeholder random data for now."

	_, err = d.blobClient.Upload(context.Background(), NopCloser(bytes.NewReader([]byte(d.data))), nil)

	return err
}

func (d *downloadPerfTest) GlobalCleanup(ctx context.Context) error {
	_, err := d.containerClient.Delete(context.Background(), nil)
	return err
}

func (d *downloadPerfTest) Setup(ctx context.Context) error {
	return nil
}

func (d *downloadPerfTest) Run(ctx context.Context) error {
	_, err := d.blobClient.Download(ctx, nil)
	return err
}

func (d *downloadPerfTest) Cleanup(ctx context.Context) error {
	return nil
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
	}
}
