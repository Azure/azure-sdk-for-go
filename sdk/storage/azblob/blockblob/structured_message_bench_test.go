// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

// benchmarkEnv holds pre-created clients for benchmark iterations.
type benchmarkEnv struct {
	containerClient *container.Client
	containerName   string
	svcClient       *azblob.Client
}

func setupBenchEnv(b *testing.B) *benchmarkEnv {
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")
	if accountName == "" || accountKey == "" {
		b.Skip("Set AZURE_STORAGE_ACCOUNT_NAME and AZURE_STORAGE_ACCOUNT_KEY to run live benchmarks")
	}

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		b.Fatal(err)
	}

	svcClient, err := azblob.NewClientWithSharedKeyCredential(
		fmt.Sprintf("https://%s.blob.core.windows.net", accountName), cred, nil)
	if err != nil {
		b.Fatal(err)
	}

	// Use timestamp + sanitized test name for unique container names across sub-benchmarks.
	// Container names must be lowercase, 3-63 chars, only alphanumeric and hyphens.
	sanitized := strings.ToLower(b.Name())
	sanitized = strings.ReplaceAll(sanitized, "/", "-")
	sanitized = strings.ReplaceAll(sanitized, "_", "-")
	if len(sanitized) > 40 {
		sanitized = sanitized[:40]
	}
	containerName := fmt.Sprintf("%s-%d", sanitized, time.Now().UnixNano()%1000000)
	_, err = svcClient.CreateContainer(context.Background(), containerName, nil)
	if err != nil {
		b.Fatal(err)
	}

	env := &benchmarkEnv{
		containerClient: svcClient.ServiceClient().NewContainerClient(containerName),
		containerName:   containerName,
		svcClient:       svcClient,
	}

	b.Cleanup(func() {
		// Delete all blobs then container
		pager := env.containerClient.NewListBlobsFlatPager(nil)
		for pager.More() {
			resp, err := pager.NextPage(context.Background())
			if err != nil {
				break
			}
			for _, blobItem := range resp.Segment.BlobItems {
				_, _ = env.containerClient.NewBlobClient(*blobItem.Name).Delete(context.Background(), nil)
			}
		}
		_, _ = svcClient.DeleteContainer(context.Background(), containerName, nil)
	})

	return env
}

func makeBenchData(size int) []byte {
	data := make([]byte, size)
	_, _ = rand.Read(data)
	return data
}

// ── Upload benchmarks: NoValidation vs ComputeCRC64 (2-pass) vs SM CRC64 (single-pass) ──

func BenchmarkUpload(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"1KB", 1 * 1024},
		{"1MB", 1 * 1024 * 1024},
		{"4MB", 4 * 1024 * 1024},
		{"16MB", 16 * 1024 * 1024},
		{"32MB", 32 * 1024 * 1024},
	}

	for _, sz := range sizes {
		data := makeBenchData(sz.size)

		b.Run(fmt.Sprintf("NoValidation/%s", sz.name), func(b *testing.B) {
			env := setupBenchEnv(b)
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				blobName := fmt.Sprintf("noval-%s-%d", sz.name, i)
				bbClient := env.containerClient.NewBlockBlobClient(blobName)
				_, err := bbClient.Upload(context.Background(),
					streaming.NopCloser(bytes.NewReader(data)),
					&blockblob.UploadOptions{})
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		// ComputeCRC64: reads entire body to compute checksum, then re-reads to send (2-pass)
		b.Run(fmt.Sprintf("ComputeCRC64/%s", sz.name), func(b *testing.B) {
			env := setupBenchEnv(b)
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				blobName := fmt.Sprintf("crc64-%s-%d", sz.name, i)
				bbClient := env.containerClient.NewBlockBlobClient(blobName)
				_, err := bbClient.Upload(context.Background(),
					streaming.NopCloser(bytes.NewReader(data)),
					&blockblob.UploadOptions{
						TransactionalValidation: blob.TransferValidationTypeComputeCRC64(),
					})
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		// SM CRC64: single-pass streaming CRC computation embedded in the request body
		b.Run(fmt.Sprintf("SMCRC64/%s", sz.name), func(b *testing.B) {
			env := setupBenchEnv(b)
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				blobName := fmt.Sprintf("sm-%s-%d", sz.name, i)
				bbClient := env.containerClient.NewBlockBlobClient(blobName)
				_, err := bbClient.Upload(context.Background(),
					streaming.NopCloser(bytes.NewReader(data)),
					&blockblob.UploadOptions{
						TransactionalValidation: blob.TransferValidationTypeComputeStructuredMessageCRC64(0),
					})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// ── Download benchmarks: NoValidation vs SM CRC64 ──
// Note: there is no ComputeCRC64 equivalent for downloads — the service computes CRC.
// SM adds client-side CRC verification + SM framing overhead on the response body.

func BenchmarkDownload(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"1KB", 1 * 1024},
		{"1MB", 1 * 1024 * 1024},
		{"4MB", 4 * 1024 * 1024},
		{"16MB", 16 * 1024 * 1024},
		{"64MB", 64 * 1024 * 1024},
		{"128MB", 128 * 1024 * 1024},
	}

	for _, sz := range sizes {
		data := makeBenchData(sz.size)

		b.Run(fmt.Sprintf("NoSM/%s", sz.name), func(b *testing.B) {
			env := setupBenchEnv(b)
			// Upload once, then benchmark download
			blobName := fmt.Sprintf("dl-nosm-%s", sz.name)
			bbClient := env.containerClient.NewBlockBlobClient(blobName)
			_, err := bbClient.UploadBuffer(context.Background(), data, nil)
			if err != nil {
				b.Fatal(err)
			}

			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				resp, err := bbClient.DownloadStream(context.Background(), nil)
				if err != nil {
					b.Fatal(err)
				}
				_, _ = io.Copy(io.Discard, resp.Body)
			}
		})

		b.Run(fmt.Sprintf("WithSM/%s", sz.name), func(b *testing.B) {
			env := setupBenchEnv(b)
			// Upload once, then benchmark download with SM
			blobName := fmt.Sprintf("dl-sm-%s", sz.name)
			bbClient := env.containerClient.NewBlockBlobClient(blobName)
			_, err := bbClient.UploadBuffer(context.Background(), data, nil)
			if err != nil {
				b.Fatal(err)
			}

			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				resp, err := bbClient.DownloadStream(context.Background(), &blob.DownloadStreamOptions{
					TransactionalValidation: blob.TransferValidationTypeComputeStructuredMessageCRC64(0),
				})
				if err != nil {
					b.Fatal(err)
				}
				_, _ = io.Copy(io.Discard, resp.Body)
			}
		})
	}
}

// ── UploadBuffer multi-block benchmark: NoValidation vs SM CRC64 ──
// Note: ComputeCRC64 is not supported for multi-block uploads (UnsupportedChecksum error).
// SM CRC64 is the only validation option that works with multi-block uploads because it
// computes per-block checksums on-the-fly.

func BenchmarkUploadBuffer(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"16MB", 16 * 1024 * 1024},
		{"64MB", 64 * 1024 * 1024},
		{"128MB", 128 * 1024 * 1024},
		{"256MB", 256 * 1024 * 1024},
	}

	for _, sz := range sizes {
		data := makeBenchData(sz.size)

		b.Run(fmt.Sprintf("NoSM/%s", sz.name), func(b *testing.B) {
			env := setupBenchEnv(b)
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				blobName := fmt.Sprintf("buf-nosm-%s-%d", sz.name, i)
				bbClient := env.containerClient.NewBlockBlobClient(blobName)
				_, err := bbClient.UploadBuffer(context.Background(), data, &blockblob.UploadBufferOptions{})
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(fmt.Sprintf("WithSM/%s", sz.name), func(b *testing.B) {
			env := setupBenchEnv(b)
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				blobName := fmt.Sprintf("buf-sm-%s-%d", sz.name, i)
				bbClient := env.containerClient.NewBlockBlobClient(blobName)
				_, err := bbClient.UploadBuffer(context.Background(), data, &blockblob.UploadBufferOptions{
					TransactionalValidation: blob.TransferValidationTypeComputeStructuredMessageCRC64(0),
				})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// ── UploadStream multi-block benchmark: NoValidation vs SM CRC64 ──
// UploadStream uses a different code path (copyFromReader) than UploadBuffer (uploadFromReader).
// It reads from a non-seekable stream, buffering into blocks internally.

func BenchmarkUploadStream(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"16MB", 16 * 1024 * 1024},
		{"64MB", 64 * 1024 * 1024},
		{"128MB", 128 * 1024 * 1024},
		{"256MB", 256 * 1024 * 1024},
	}

	for _, sz := range sizes {
		data := makeBenchData(sz.size)

		b.Run(fmt.Sprintf("NoSM/%s", sz.name), func(b *testing.B) {
			env := setupBenchEnv(b)
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				blobName := fmt.Sprintf("str-nosm-%s-%d", sz.name, i)
				bbClient := env.containerClient.NewBlockBlobClient(blobName)
				_, err := bbClient.UploadStream(context.Background(),
					bytes.NewReader(data),
					&blockblob.UploadStreamOptions{})
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		b.Run(fmt.Sprintf("WithSM/%s", sz.name), func(b *testing.B) {
			env := setupBenchEnv(b)
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				blobName := fmt.Sprintf("str-sm-%s-%d", sz.name, i)
				bbClient := env.containerClient.NewBlockBlobClient(blobName)
				_, err := bbClient.UploadStream(context.Background(),
					bytes.NewReader(data),
					&blockblob.UploadStreamOptions{
						TransactionalValidation: blob.TransferValidationTypeComputeStructuredMessageCRC64(0),
					})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
