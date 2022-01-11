// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/eng/tools/azperf/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/spf13/cobra"
)

var UploadBlobCmd = &cobra.Command{
	Use:   "BlobUploadTest",
	Short: "BlobUploadTest performance test",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(c *cobra.Command, args []string) error {
		return perf.RunPerfTest(&uploadPerfTest{})
	},
}

type uploadPerfTest struct {
	blobName        string
	containerClient azblob.ContainerClient
	blobClient      azblob.BlockBlobClient
	data            string
}

func (m *uploadPerfTest) GlobalSetup(ctx context.Context) error {
	m.blobName = "uploadtest"

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return fmt.Errorf("the environment variable 'AZBLOB_CONNECTION_STRING' could not be found")
	}

	containerClient, err := azblob.NewContainerClientFromConnectionString(connStr, "uploadtest", nil)
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

	return nil
}

func (m *uploadPerfTest) GlobalTearDown(ctx context.Context) error {
	_, err := m.containerClient.Delete(context.Background(), nil)
	return err
}

func (m *uploadPerfTest) Setup(ctx context.Context) error {
	return nil
}

func (m *uploadPerfTest) Run(ctx context.Context) error {
	_, err := m.blobClient.Upload(ctx, NopCloser(bytes.NewReader([]byte(m.data))), nil)
	return err
}

func (m *uploadPerfTest) TearDown(ctx context.Context) error {
	return nil
}

func (m *uploadPerfTest) GetMetadata() string {
	return "BlobUploadTest"
}
