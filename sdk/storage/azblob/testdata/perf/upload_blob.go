// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
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
	blobPrefix            string
	globalContainerClient *container.Client

	// size is the per-iteration upload size in bytes.
	size int64

	// payload is populated only for --upload-method=buffer (which requires a
	// []byte). For stream/single paths it stays nil and randomSeed is used to
	// avoid materializing arbitrarily large payloads in memory.
	payload []byte

	// randomSeed is a small (randomSeedSize) random buffer tiled by
	// randomStream to produce the per-iteration upload bytes for stream/single
	// methods. Shared (read-only) across goroutines.
	randomSeed []byte
}

// NewUploadTest is called once per process
func NewUploadTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	u := &uploadTestGlobal{
		PerfTestOptions: options,
		// Suffix with a unique timestamp so concurrent runs and --no-cleanup
		// leftovers from prior runs do not collide on container creation.
		containerName: fmt.Sprintf("uploadcontainer-%d", time.Now().UnixNano()),
		blobPrefix:    "uploadblob",
		size:          int64(uploadTestOpts.size),
	}

	// Validate memory budget before doing any I/O so the user gets a clear
	// error instead of an OOM kill or a network call followed by a panic.
	if uploadMethod == "buffer" {
		if err := checkBufferMemoryBudget("--upload-method buffer", u.size); err != nil {
			return nil, err
		}
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

	// Only the buffer method requires the full payload to be materialized in
	// RAM (the API signature is []byte). For stream/single we keep a small
	// random seed and tile it via randomStream so --size can be arbitrarily
	// large without OOMing the process.
	if uploadMethod == "buffer" {
		payload, err := generateRandomBytes(uploadTestOpts.size)
		if err != nil {
			return nil, err
		}
		u.payload = payload
	} else {
		seed, err := generateRandomBytes(randomSeedSize)
		if err != nil {
			return nil, err
		}
		u.randomSeed = seed
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
	blobClient *blockblob.Client
	blobName   string
}

// NewPerfTest is called once per goroutine
func (g *uploadTestGlobal) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	u := &uploadPerfTest{
		uploadTestGlobal: g,
		PerfTestOptions:  *options,
		// Each goroutine targets a unique blob to avoid same-blob write
		// contention/throttling at the service.
		blobName: fmt.Sprintf("%s-%s", g.blobPrefix, options.Name),
	}

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	containerClient, err := container.NewClientFromConnectionString(
		connStr,
		u.uploadTestGlobal.containerName,
		&container.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: u.PerfTestOptions.Transporter,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	u.blobClient = containerClient.NewBlockBlobClient(u.blobName)

	return u, nil
}

func (m *uploadPerfTest) Run(ctx context.Context) error {
	switch uploadMethod {
	case "buffer":
		_, err := m.blobClient.UploadBuffer(ctx, m.payload, &blockblob.UploadBufferOptions{
			BlockSize:   commonBlockSize,
			Concurrency: uint16(commonConcurrency),
		})
		return err
	case "stream":
		// Fresh randomStream per iteration so the read offset starts at 0.
		// Tiles the shared 1 MiB random seed, so peak memory is ~1 MiB +
		// SDK chunk buffers regardless of --size.
		stream := newRandomStream(m.randomSeed, m.size)
		_, err := m.blobClient.UploadStream(ctx, stream, &blockblob.UploadStreamOptions{
			BlockSize:   commonBlockSize,
			Concurrency: int(commonConcurrency),
		})
		return err
	case "single", "":
		// Single REST PUT. Uses a tiled randomStream so the payload is not
		// materialized in memory; Azure Storage caps a single PUT Blob at
		// 5000 MiB so very large --size values should use buffer or stream.
		reader := newRandomStream(m.randomSeed, m.size)
		_, err := m.blobClient.Upload(ctx, reader, nil)
		return err
	default:
		return fmt.Errorf("unknown --upload-method %q (expected single|buffer|stream)", uploadMethod)
	}
}

func (m *uploadPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
