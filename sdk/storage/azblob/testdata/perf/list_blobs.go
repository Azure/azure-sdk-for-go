// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
)

type listTestOptions struct {
	count int
}

var listTestOpts = listTestOptions{count: 100}

// uploadTestRegister is called once per process
func listTestRegister() {
	flag.IntVar(&listTestOpts.count, "num-blobs", 100, "Number of blobs to list.")
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
	var (
		containerClient *container.Client
		err             error
	)
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if ok {
		containerClient, err = container.NewClientFromConnectionString(connStr, l.containerName, nil)
		if err != nil {
			return nil, err
		}
	} else {
		storageSuffix, ok := os.LookupEnv("PRIMARY_STORAGE_ACCOUNT_BLOB_ENDPOINT_SUFFIX")
		if !ok {
			return nil, fmt.Errorf("the environment variable 'PRIMARY_STORAGE_ACCOUNT_BLOB_ENDPOINT_SUFFIX' could not be found")
		}
		accountName, ok := os.LookupEnv("BLOB_STORAGE_ACCOUNT_NAME")
		if !ok {
			return nil, fmt.Errorf("the environment variable 'BLOB_STORAGE_ACCOUNT_NAME' could not be found")
		}
		blobUrl := fmt.Sprintf("https://%s.%s/%s", accountName, storageSuffix, l.containerName)

		credential, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, err
		}

		containerClient, err = container.NewClient(blobUrl, credential, nil)
		if err != nil {
			return nil, err
		}
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
	var (
		containerClient *container.Client
		err             error
	)

	if ok {
		containerClient, err = container.NewClientFromConnectionString(connStr, l.containerName, nil)
		if err != nil {
			return err
		}
	} else {
		storageSuffix, ok := os.LookupEnv("PRIMARY_STORAGE_ACCOUNT_BLOB_ENDPOINT_SUFFIX")
		if !ok {
			return fmt.Errorf("the environment variable 'PRIMARY_STORAGE_ACCOUNT_BLOB_ENDPOINT_SUFFIX' could not be found")
		}
		accountName, ok := os.LookupEnv("BLOB_STORAGE_ACCOUNT_NAME")
		if !ok {
			return fmt.Errorf("the environment variable 'BLOB_STORAGE_ACCOUNT_NAME' could not be found")
		}
		blobUrl := fmt.Sprintf("https://%s.%s/%s", accountName, storageSuffix, l.containerName)

		credential, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return err
		}

		containerClient, err = container.NewClient(blobUrl, credential, nil)
		if err != nil {
			return err
		}
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

	var (
		containerClient *container.Client
		err             error
	)

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")

	if ok {
		containerClient, err = container.NewClientFromConnectionString(connStr, u.listTestGlobal.containerName,
			&container.ClientOptions{
				ClientOptions: azcore.ClientOptions{
					Transport: g.PerfTestOptions.Transporter,
				},
			})
		if err != nil {
			return nil, err
		}
	} else {
		storageSuffix, ok := os.LookupEnv("PRIMARY_STORAGE_ACCOUNT_BLOB_ENDPOINT_SUFFIX")
		if !ok {
			return nil, fmt.Errorf("the environment variable 'PRIMARY_STORAGE_ACCOUNT_BLOB_ENDPOINT_SUFFIX' could not be found")
		}
		accountName, ok := os.LookupEnv("BLOB_STORAGE_ACCOUNT_NAME")
		if !ok {
			return nil, fmt.Errorf("the environment variable 'BLOB_STORAGE_ACCOUNT_NAME' could not be found")
		}
		blobUrl := fmt.Sprintf("https://%s.%s/%s", accountName, storageSuffix, u.listTestGlobal.containerName)

		credential, err := azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return nil, err
		}

		containerClient, err = container.NewClient(blobUrl, credential, &container.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: g.PerfTestOptions.Transporter,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	u.containerClient = containerClient

	return u, nil
}

func (m *listPerfTest) Run(ctx context.Context) error {
	c := int32(listTestOpts.count)
	pager := m.containerClient.NewListBlobsFlatPager(&container.ListBlobsFlatOptions{
		MaxResults: &c,
	})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *listPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
