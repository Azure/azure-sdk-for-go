// // Copyright (c) Microsoft Corporation. All rights reserved.
// // Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/spf13/pflag"
)

type listTestOptions struct {
	count int32
}

var listTestOpts listTestOptions = listTestOptions{count: 100}

// uploadTestRegister is called once per process
func listTestRegister() {
	pflag.Int32Var(&listTestOpts.count, "num-blobs", 100, "Number of blobs to list.")
}

type listTestGlobal struct {
	perf.PerfTestOptions
	containerName string
	blobName      string
}

// NewListTest is called once per process
func NewListTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	l := &listTestGlobal{
		PerfTestOptions: options,
		containerName:   "listcontainer",
		blobName:        "listblob",
	}
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, l.containerName, nil)
	if err != nil {
		return nil, err
	}
	_, err = containerClient.Create(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 100; i++ {
		blobClient := containerClient.NewBlockBlobClient(fmt.Sprintf("%s%d", l.blobName, i))
		_, err = blobClient.Upload(
			context.Background(),
			NopCloser(bytes.NewReader([]byte(""))),
			nil,
		)
		if err != nil {
			return nil, err
		}
	}

	return l, nil
}

func (l *listTestGlobal) GlobalCleanup(ctx context.Context) error {
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, l.containerName, nil)
	if err != nil {
		return err
	}

	_, err = containerClient.Delete(context.Background(), nil)
	return err
}

type listPerfTest struct {
	*listTestGlobal
	perf.PerfTestOptions
	containerClient azblob.ContainerClient
}

// NewPerfTest is called once per goroutine
func (g *listTestGlobal) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	u := &listPerfTest{
		listTestGlobal:  g,
		PerfTestOptions: *options,
	}

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(
		connStr,
		u.listTestGlobal.containerName,
		&azblob.ClientOptions{
			Transporter: u.PerfTestOptions.Transporter,
		},
	)
	if err != nil {
		return nil, err
	}
	u.containerClient = containerClient

	return u, nil
}

func (m *listPerfTest) Run(ctx context.Context) error {
	pager := m.containerClient.ListBlobsFlat(&azblob.ContainerListBlobFlatSegmentOptions{
		Maxresults: &listTestOpts.count,
	})
	for pager.NextPage(context.Background()) {
	}
	return pager.Err()
}

func (m *listPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
