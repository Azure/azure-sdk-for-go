// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type downloadTestOptions struct {
	size int
}

var downloadTestOpts = downloadTestOptions{size: 10240}

// downloadTestRegister is called once per process
func downloadTestRegister() {
	flag.IntVar(&downloadTestOpts.size, "size", 10240, "Size in bytes of data to be transferred in upload or download tests. Default is 10240.")
}

type downloadTestGlobal struct {
	perf.PerfTestOptions
	containerName string
	blobName      string
}

// NewDownloadTest is called once per process
func NewDownloadTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	d := &downloadTestGlobal{
		PerfTestOptions: options,
		containerName:   "downloadcontainer",
		blobName:        "downloadblob",
	}

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := container.NewClientFromConnectionString(connStr, d.containerName, nil)
	if err != nil {
		return nil, err
	}
	_, err = containerClient.Create(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	blobClient := containerClient.NewBlockBlobClient(d.blobName)

	data, err := perf.NewRandomStream(downloadTestOpts.size)
	if err != nil {
		return nil, err
	}

	_, err = blobClient.Upload(context.Background(), data, nil)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (d *downloadTestGlobal) GlobalCleanup(ctx context.Context) error {
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := container.NewClientFromConnectionString(connStr, d.containerName, nil)
	if err != nil {
		return err
	}

	_, err = containerClient.Delete(context.Background(), nil)
	return err
}

type downloadPerfTest struct {
	*downloadTestGlobal
	perf.PerfTestOptions
	data       io.ReadSeekCloser
	blobClient *blockblob.Client
}

// NewPerfTest is called once per goroutine
func (g *downloadTestGlobal) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	d := &downloadPerfTest{
		downloadTestGlobal: g,
		PerfTestOptions:    *options,
	}

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := container.NewClientFromConnectionString(connStr, d.downloadTestGlobal.containerName, &azblob.ClientOptions{
		Transport: d.PerfTestOptions.Transporter,
	})
	if err != nil {
		return nil, err
	}
	d.blobClient = containerClient.NewBlockBlobClient(d.blobName)

	data, err := perf.NewRandomStream(downloadTestOpts.size)
	if err != nil {
		return nil, err
	}
	d.data = data

	return d, err
}

func (d *downloadPerfTest) Run(ctx context.Context) error {
	get, err := d.blobClient.Download(ctx, nil)
	if err != nil {
		return err
	}
	downloadedData := &bytes.Buffer{}
	reader := get.Body(nil)
	_, err = downloadedData.ReadFrom(reader)
	if err != nil {
		return err
	}
	return reader.Close()
}

func (*downloadPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
