// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type listTestOptions struct {
	count int
}

var listTestOpts = listTestOptions{count: 100}

// listTestRegister is called once per process
func listTestRegister() {
	flag.IntVar(&listTestOpts.count, "num-blobs", 100, "Total number of blobs to create in the container and list during the test.")
}

type listTestGlobal struct {
	perf.PerfTestOptions
	containerName string
	blobPrefix    string
}

// NewListTest is called once per process
func NewListTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	l := &listTestGlobal{
		PerfTestOptions: options,
		// Suffix with a unique timestamp so concurrent runs and --no-cleanup
		// leftovers from prior runs do not collide on container creation.
		containerName: fmt.Sprintf("listcontainer-%d", time.Now().UnixNano()),
		blobPrefix:    "listblob",
	}
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := container.NewClientFromConnectionString(connStr, l.containerName, nil)
	if err != nil {
		return nil, err
	}
	_, err = containerClient.Create(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	// Seed the container with the requested number of empty blobs. The earlier
	// implementation hard-coded 100 here, which meant matrix entries asking
	// for thousands of blobs silently listed at most 100.
	emptyBody := NopCloser(bytes.NewReader(nil))
	for i := 0; i < listTestOpts.count; i++ {
		blobClient := containerClient.NewBlockBlobClient(fmt.Sprintf("%s%d", l.blobPrefix, i))
		if _, err = blobClient.Upload(context.Background(), emptyBody, nil); err != nil {
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

	containerClient, err := container.NewClientFromConnectionString(connStr, l.containerName, nil)
	if err != nil {
		return err
	}

	_, err = containerClient.Delete(context.Background(), nil)
	return err
}

type listPerfTest struct {
	*listTestGlobal
	perf.PerfTestOptions
	containerClient *container.Client
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

	containerClient, err := container.NewClientFromConnectionString(
		connStr,
		u.listTestGlobal.containerName,
		&container.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: g.PerfTestOptions.Transporter,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	u.containerClient = containerClient

	return u, nil
}

func (m *listPerfTest) Run(ctx context.Context) error {
	opts := &container.ListBlobsFlatOptions{}
	if listPageSize > 0 {
		p := int32(listPageSize)
		opts.MaxResults = &p
	}
	pager := m.containerClient.NewListBlobsFlatPager(opts)
	for pager.More() {
		if _, err := pager.NextPage(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *listPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
