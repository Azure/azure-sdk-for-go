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
	blobName        string
	containerClient azblob.ContainerClient
	blobClient      azblob.BlockBlobClient
	data            string
}

func (m *downloadPerfTest) GlobalSetup(ctx context.Context) error {
	m.blobName = "downloadtest"

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZBLOB_CONNECTION_STRING' could not be found")
	}

	t, err := perf.NewProxyTransport(&perf.TransportOptions{TestName: m.GetMetadata()})
	if err != nil {
		return err
	}
	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, "downloadtest", &azblob.ClientOptions{Transporter: t})
	if err != nil {
		return err
	}
	m.containerClient = containerClient
	_, err = m.containerClient.Create(context.Background(), nil)
	if err != nil {
		return err
	}

	m.blobClient = m.containerClient.NewBlockBlobClient(m.blobName)

	m.data = "This is all placeholder random data for now. This is all placeholder random data for now. This is all placeholder random data for now."

	_, err = m.blobClient.Upload(context.Background(), NopCloser(bytes.NewReader([]byte(m.data))), nil)

	return err
}

func (m *downloadPerfTest) GlobalTearDown(ctx context.Context) error {
	_, err := m.containerClient.Delete(context.Background(), nil)
	return err
}

func (m *downloadPerfTest) Setup(ctx context.Context) error {
	return nil
}

func (m *downloadPerfTest) Run(ctx context.Context) error {
	_, err := m.blobClient.Download(ctx, nil)
	return err
}

func (m *downloadPerfTest) TearDown(ctx context.Context) error {
	return nil
}

func (m *downloadPerfTest) GetMetadata() string {
	return "BlobDownloadTest"
}
