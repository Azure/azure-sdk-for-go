// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
)

// GiB-scale parallelism benchmarks to measure the impact of default concurrency changes.
// These use UploadBuffer/DownloadBuffer which perform parallel chunked transfers,
// where concurrency directly affects throughput.
//
// Run with:
//   AZURE_STORAGE_ACCOUNT_NAME=... AZURE_STORAGE_ACCOUNT_KEY=... \
//   go test ./blockblob/... -bench 'BenchmarkParallelism' -benchtime=3x -run '^$' -timeout 1800s
//
// To compare old vs new defaults:
//   # Old default (concurrency=5):
//   AZURE_STORAGE_USE_LEGACY_DEFAULT_CONCURRENCY=true go test ...
//   # New default (concurrency=NumCPU clamped 8-96):
//   go test ...

func BenchmarkParallelismUploadBuffer(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"256MB", 256 * 1024 * 1024},
		{"512MB", 512 * 1024 * 1024},
		{"1GiB", 1024 * 1024 * 1024},
		{"2GiB", 2 * 1024 * 1024 * 1024},
		{"4GiB", 4 * 1024 * 1024 * 1024},
	}

	for _, sz := range sizes {
		b.Run(sz.name, func(b *testing.B) {
			data := make([]byte, sz.size)
			_, _ = rand.Read(data)

			env := setupBenchEnv(b)
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				blobName := fmt.Sprintf("par-upbuf-%s-%d", sz.name, i)
				bbClient := env.containerClient.NewBlockBlobClient(blobName)
				_, err := bbClient.UploadBuffer(context.Background(), data, &blockblob.UploadBufferOptions{})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParallelismUploadStream(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"256MB", 256 * 1024 * 1024},
		{"512MB", 512 * 1024 * 1024},
		{"1GiB", 1024 * 1024 * 1024},
		{"2GiB", 2 * 1024 * 1024 * 1024},
		{"4GiB", 4 * 1024 * 1024 * 1024},
	}

	for _, sz := range sizes {
		b.Run(sz.name, func(b *testing.B) {
			data := make([]byte, sz.size)
			_, _ = rand.Read(data)

			env := setupBenchEnv(b)
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				blobName := fmt.Sprintf("par-upstr-%s-%d", sz.name, i)
				bbClient := env.containerClient.NewBlockBlobClient(blobName)
				_, err := bbClient.UploadStream(context.Background(),
					bytes.NewReader(data),
					&blockblob.UploadStreamOptions{})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParallelismDownloadBuffer(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"256MB", 256 * 1024 * 1024},
		{"512MB", 512 * 1024 * 1024},
		{"1GiB", 1024 * 1024 * 1024},
		{"2GiB", 2 * 1024 * 1024 * 1024},
		{"4GiB", 4 * 1024 * 1024 * 1024},
	}

	for _, sz := range sizes {
		b.Run(sz.name, func(b *testing.B) {
			data := make([]byte, sz.size)
			_, _ = rand.Read(data)

			env := setupBenchEnv(b)

			// Upload the blob once before benchmarking download
			blobName := fmt.Sprintf("par-dlbuf-%s", sz.name)
			bbClient := env.containerClient.NewBlockBlobClient(blobName)
			_, err := bbClient.UploadBuffer(context.Background(), data, nil)
			if err != nil {
				b.Fatal(err)
			}

			// Get blob properties to know the size for buffer allocation
			props, err := bbClient.GetProperties(context.Background(), nil)
			if err != nil {
				b.Fatal(err)
			}
			blobSize := *props.ContentLength

			b.SetBytes(blobSize)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				buf := make([]byte, blobSize)
				_, err := bbClient.BlobClient().DownloadBuffer(context.Background(), buf, &blob.DownloadBufferOptions{})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
