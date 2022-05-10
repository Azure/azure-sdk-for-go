// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile"
)

type uploadTestOptions struct {
	size int
}

var uploadTestOpts uploadTestOptions = uploadTestOptions{size: 10240}

// uploadTestRegister is called once per process
func uploadTestRegister() {
	flag.IntVar(&uploadTestOpts.size, "size", 10240, "Size in bytes of data to be transferred in upload or download tests.")
}

type uploadTestGlobal struct {
	perf.PerfTestOptions
	shareName         string
	directoryName     string
	fileName          string
	globalShareClient *azfile.ShareClient
}

// NewUploadTest is called once per process
func NewUploadTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	u := &uploadTestGlobal{
		PerfTestOptions: options,
		shareName:       "uploadshare",
		directoryName:   "uploaddir",
		fileName:        "uploadfile",
	}

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	shareClient, err := azfile.NewShareClientFromConnectionString(connStr, u.shareName, nil)
	if err != nil {
		return nil, err
	}
	u.globalShareClient = shareClient
	_, err = u.globalShareClient.Create(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (u *uploadTestGlobal) GlobalCleanup(ctx context.Context) error {
	_, err := u.globalShareClient.Delete(context.Background(), nil)
	return err
}

type uploadPerfTest struct {
	*uploadTestGlobal
	perf.PerfTestOptions
	data       io.ReadSeekCloser
	fileClient *azfile.FileClient
}

// NewPerfTest is called once per goroutine
func (g *uploadTestGlobal) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	u := &uploadPerfTest{
		uploadTestGlobal: g,
		PerfTestOptions:  *options,
	}

	connStr, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	shareClient, err := azfile.NewShareClientFromConnectionString(
		connStr,
		u.uploadTestGlobal.shareName,
		&azfile.ClientOptions{
			Transport: u.PerfTestOptions.Transporter,
		},
	)
	if err != nil {
		return nil, err
	}
	dirClient, err := shareClient.NewDirectoryClient(u.directoryName)
	if err != nil {
		return nil, err
	}

	fClient, err := dirClient.NewFileClient(u.fileName)
	if err != nil {
		return nil, err
	}
	u.fileClient = fClient

	data, err := perf.NewRandomStream(uploadTestOpts.size)
	if err != nil {
		return nil, err
	}
	u.data = data

	return u, nil
}

func (m *uploadPerfTest) Run(ctx context.Context) error {
	_, err := m.data.Seek(0, io.SeekStart) // rewind to the beginning
	if err != nil {
		return err
	}
	_, err = m.fileClient.UploadRange(ctx, 0, m.data, nil)
	return err
}

func (m *uploadPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
