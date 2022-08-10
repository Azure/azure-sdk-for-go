// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
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

type uploadTestOptions struct {
	size int
}

var uploadTestOpts = uploadTestOptions{size: 10240}

// uploadTestRegister is called once per process
func uploadTestRegister() {
	flag.IntVar(&uploadTestOpts.size, "size", 10240, "Size in bytes of data to be transferred in upload or download tests.")
}

type uploadTestGlobal struct {
	perf.PerfTestOptions
	containerName         string
	blobName              string
	globalContainerClient *container.Client
}

// NewUploadTest is called once per process
func NewUploadTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	u := &uploadTestGlobal{
		PerfTestOptions: options,
		containerName:   "uploadcontainer",
		blobName:        "uploadblob",
	}

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := container.NewClientFromConnectionString(connStr, u.containerName, nil)
	if err != nil {
		return nil, err
	}
	u.globalContainerClient = containerClient
	_, err = u.globalContainerClient.Create(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *uploadTestGlobal) GlobalCleanup(ctx context.Context) error {
	_, err := u.globalContainerClient.Delete(context.Background(), nil)
	return err
}

type uploadPerfTest struct {
	*uploadTestGlobal
	perf.PerfTestOptions
	data       io.ReadSeekCloser
	blobClient *blockblob.Client
}

// NewPerfTest is called once per goroutine
func (g *uploadTestGlobal) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	u := &uploadPerfTest{
		uploadTestGlobal: g,
		PerfTestOptions:  *options,
	}

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := container.NewClientFromConnectionString(
		connStr,
		u.uploadTestGlobal.containerName,
		&azblob.ClientOptions{
			Transport: u.PerfTestOptions.Transporter,
		},
	)
	if err != nil {
		return nil, err
	}
	bc := containerClient.NewBlockBlobClient(u.blobName)
	if err != nil {
		return nil, err
	}
	u.blobClient = bc

	data, err := perf.NewRandomStream(uploadTestOpts.size)
	if err != nil {
		return nil, err
	}
	u.data = data

	return u, nil
}

func (m *uploadPerfTest) Run(ctx context.Context) error {
	_, err := m.data.Seek(0, io.SeekStart) // rewind to the beginning
	if err != nil {
		return err
	}
	_, err = m.blobClient.Upload(ctx, m.data, nil)
	return err
}

func (m *uploadPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
