// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type listTestOptions struct {
	count       int
	parallelism int
}

var listTestOpts = listTestOptions{count: 100}

// listTestRegister is called once per process
func listTestRegister() {
	flag.IntVar(&listTestOpts.count, "num-blobs", 100, "Total number of blobs to create in the container and list during the test.")
	flag.IntVar(&listTestOpts.parallelism, "num-blobs-parallelism", 0, "Number of parallel workers used to create blobs during list test setup (0=number of CPU cores).")
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
	containerClient, err := newContainerClient(l.containerName, nil)
	if err != nil {
		return nil, err
	}
	_, err = containerClient.Create(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Seed the container with the requested number of empty blobs. The earlier
	// implementation hard-coded 100 here, which meant matrix entries asking
	// for thousands of blobs silently listed at most 100.
	if err = l.seedBlobs(ctx, containerClient, listTestOpts.count, listTestOpts.parallelism); err != nil {
		return nil, err
	}

	return l, nil
}

// seedBlobs creates count empty blobs in the container using a pool of workers.
// Uploading blobs one at a time is prohibitively slow for large blob counts
// (e.g. 500k), so uploads are parallelized across parallelism workers. When
// parallelism is <= 0 it defaults to the number of CPU cores on the machine.
func (l *listTestGlobal) seedBlobs(ctx context.Context, containerClient *container.Client, count, parallelism int) error {
	if count <= 0 {
		return nil
	}
	if parallelism <= 0 {
		parallelism = runtime.NumCPU()
	}
	if parallelism > count {
		parallelism = count
	}

	// Cancel outstanding uploads as soon as any worker reports an error.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	indexes := make(chan int)
	var wg sync.WaitGroup
	var once sync.Once
	var firstErr error

	for w := 0; w < parallelism; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range indexes {
				blobClient := containerClient.NewBlockBlobClient(fmt.Sprintf("%s%d", l.blobPrefix, i))
				if _, err := blobClient.Upload(ctx, NopCloser(bytes.NewReader(nil)), nil); err != nil {
					once.Do(func() {
						firstErr = err
						cancel()
					})
					return
				}
			}
		}()
	}

feed:
	for i := 0; i < count; i++ {
		select {
		case indexes <- i:
		case <-ctx.Done():
			break feed
		}
	}
	close(indexes)
	wg.Wait()

	return firstErr
}

func (l *listTestGlobal) GlobalCleanup(ctx context.Context) error {
	containerClient, err := newContainerClient(l.containerName, nil)
	if err != nil {
		return err
	}

	_, err = containerClient.Delete(ctx, nil)
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

	containerClient, err := newContainerClient(
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
