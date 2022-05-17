// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile"
)

type listTestOptions struct {
	count int
}

var listTestOpts listTestOptions = listTestOptions{count: 100}

// uploadTestRegister is called once per process
func listTestRegister() {
	flag.IntVar(&listTestOpts.count, "num-blobs", 100, "Number of blobs to list.")
}

type listTestGlobal struct {
	perf.PerfTestOptions
	svcClient  *azfile.ServiceClient
	shareName  string
	shareCount int
	dirName    string
	dirCount   int
	fileName   string
	fileCount  int
}

// NewListTest is called once per process
func NewListTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	l := &listTestGlobal{
		PerfTestOptions: options,
		shareName:       "gosdkperftestshare",
		shareCount:      int(1),
		dirName:         "gosdkperftestdir",
		dirCount:        int(5),
		fileName:        "gosdkperftestfile",
		fileCount:       int(5),
	}
	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	svcClient, err := azfile.NewServiceClientFromConnectionString(connStr, nil)
	if err != nil {
		return nil, err
	}
	l.svcClient = svcClient

	for i := 0; i < l.shareCount; i++ {
		shareClient, err := svcClient.NewShareClient(fmt.Sprintf("%s%d", l.shareName, i))
		if err != nil {
			return nil, err
		}

		for j := 0; j < l.dirCount; j++ {
			dirClient, err := shareClient.NewDirectoryClient(fmt.Sprintf("%s%d%d", l.dirName, i, j))
			if err != nil {
				return nil, err
			}

			for k := 0; k < l.fileCount; k++ {
				fileClient, err := dirClient.NewFileClient(fmt.Sprintf("%s%d%d%d", l.fileName, i, j, k))
				if err != nil {
					return nil, err
				}
				_, err = fileClient.UploadRange(context.Background(), 0, NopCloser(bytes.NewReader([]byte(""))), nil)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return l, nil
}

func (l *listTestGlobal) GlobalCleanup(ctx context.Context) error {
	for i := 0; i < l.shareCount; i++ {
		shareClient, err := l.svcClient.NewShareClient(fmt.Sprintf("%s%d", l.shareName, i))
		if err != nil {
			return err
		}
		_, err = shareClient.Delete(context.TODO(), nil)
		if err != nil {
			return err
		}
	}

	return nil
}

type listPerfTest struct {
	*listTestGlobal
	perf.PerfTestOptions
	shareClient *azfile.ShareClient
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

	containerClient, err := azfile.NewShareClientFromConnectionString(
		connStr,
		u.listTestGlobal.shareName,
		&azfile.ClientOptions{
			Transport: u.PerfTestOptions.Transporter,
		},
	)
	if err != nil {
		return nil, err
	}
	u.shareClient = containerClient

	return u, nil
}

func (l *listPerfTest) Run(ctx context.Context) error {
	pager := l.svcClient.ListShares(nil)
	for pager.More() {
		resp, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}

		if len(resp.ShareItems) != l.shareCount {
			return errors.New("number of shares created not matching with number of shares listed")
		}
	}

	for i := 0; i < l.shareCount; i++ {
		shareClient, err := l.svcClient.NewShareClient(fmt.Sprintf("%s%d", l.shareName, i))
		if err != nil {
			return err
		}

		for j := 0; j < l.dirCount; j++ {
			dirClient, err := shareClient.NewDirectoryClient(fmt.Sprintf("%s%d%d", l.dirName, i, j))
			if err != nil {
				return err
			}

			listFilePager := dirClient.ListFilesAndDirectories(nil)
			for listFilePager.More() {
				resp, err := listFilePager.NextPage(ctx)
				if err != nil {
					return err
				}

				if len(resp.Segment.FileItems) != l.fileCount {
					return errors.New("number of files created not matching with number of files listed")
				}
			}
		}
	}

	return nil
}

func (m *listPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
