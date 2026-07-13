// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
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
	size          int
}

// NewDownloadTest is called once per process
func NewDownloadTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	d := &downloadTestGlobal{
		PerfTestOptions: options,
		// Suffix with a unique timestamp so concurrent runs and --no-cleanup
		// leftovers from prior runs do not collide on container creation.
		containerName: fmt.Sprintf("downloadcontainer-%d", time.Now().UnixNano()),
		blobName:      "downloadblob",
		size:          downloadTestOpts.size,
	}

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := container.NewClientFromConnectionString(connStr, d.containerName, nil)
	if err != nil {
		return nil, err
	}
	_, err = containerClient.Create(ctx, nil)
	if err != nil {
		return nil, err
	}

	blobClient := containerClient.NewBlockBlobClient(d.blobName)

	// Seed the test blob without materializing the full payload in RAM. We
	// stream a randomStream (tiles a 1 MiB random seed up to d.size) through
	// UploadStream, which internally stages parallel blocks. This keeps the
	// process memory bounded by the SDK's BlockSize x Concurrency buffer pool
	// regardless of how large d.size is.
	seed, err := generateRandomBytes(randomSeedSize)
	if err != nil {
		return nil, err
	}
	_, err = blobClient.UploadStream(
		ctx,
		newRandomStream(seed, int64(d.size)),
		&blockblob.UploadStreamOptions{
			BlockSize:   commonBlockSize,
			Concurrency: int(commonConcurrency),
		},
	)
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

	_, err = containerClient.Delete(ctx, nil)
	return err
}

type downloadPerfTest struct {
	*downloadTestGlobal
	perf.PerfTestOptions
	blobClient *blob.Client
	// buffer is a per-goroutine, preallocated download target reused across
	// iterations for the "buffer" download method. Avoids per-iteration
	// allocation/copy cost in the measurement.
	buffer []byte
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

	containerClient, err := container.NewClientFromConnectionString(connStr, d.downloadTestGlobal.containerName, &container.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: d.PerfTestOptions.Transporter,
		},
	})
	if err != nil {
		return nil, err
	}
	d.blobClient = containerClient.NewBlockBlobClient(d.blobName).BlobClient()

	if downloadMethod == "buffer" {
		// Guard against allocating per-goroutine buffers that won't fit in
		// available system memory. The total memory cost of the buffer
		// download method is roughly size * parallel; we apply the same 80%
		// budget check used by --upload-method=buffer.
		if err := checkBufferMemoryBudget("--download-method buffer", int64(g.size)*int64(perf.Parallel())); err != nil {
			return nil, err
		}
		d.buffer = make([]byte, g.size)
	}

	return d, nil
}

func (d *downloadPerfTest) Run(ctx context.Context) error {
	switch downloadMethod {
	case "buffer":
		_, err := d.blobClient.DownloadBuffer(ctx, d.buffer, &blob.DownloadBufferOptions{
			BlockSize:   commonBlockSize,
			Concurrency: uint16(commonConcurrency),
		})
		return err
	case "stream", "":
		get, err := d.blobClient.DownloadStream(ctx, nil)
		if err != nil {
			return err
		}
		defer get.Body.Close()
		_, err = io.Copy(io.Discard, get.Body)
		return err
	default:
		return fmt.Errorf("unknown --download-method %q (expected stream|buffer)", downloadMethod)
	}
}

func (*downloadPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
