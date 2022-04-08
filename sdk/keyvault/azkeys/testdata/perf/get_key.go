// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

type getCertificatesTestOptions struct{}

var getCertTestOpts decryptTestOptions = decryptTestOptions{}

type GetKeyTest struct {
	perf.PerfTestOptions
	keyName string
	client  *azkeys.Client
}

// NewGetCertificateTest is called once per process
func NewGetCertificateTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	certName := "livekvtestgetcertperfkey"
	d := &GetKeyTest{
		PerfTestOptions: options,
		keyName:         certName,
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_KEYVAULT_URL' could not be found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultURL, cred, nil)
	if err != nil {
		return nil, err
	}

	_, err = client.CreateRSAKey(ctx, d.keyName, nil)
	if err != nil {
		return nil, err
	}

	d.client = client
	return d, nil
}

func (gct *GetKeyTest) GlobalCleanup(ctx context.Context) error {
	poller, err := gct.client.BeginDeleteKey(ctx, gct.keyName, nil)
	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 500*time.Millisecond)
	if err != nil {
		return err
	}

	_, err = gct.client.PurgeDeletedKey(ctx, gct.keyName, nil)
	return err
}

type GetKeyPerfTest struct {
	*GetKeyTest
	perf.PerfTestOptions
	client          *azkeys.Client
	keyName string
}

// NewPerfTest is called once per goroutine
func (gct *GetKeyTest) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	return &GetKeyPerfTest{
		GetKeyTest:      gct,
		PerfTestOptions: *options,
		client:          gct.client,
		keyName: gct.keyName,
	}, nil
}

func (gcpt *GetKeyPerfTest) Run(ctx context.Context) error {
	_, err := gcpt.client.GetKey(ctx, gcpt.keyName, nil)
	return err
}

func (*GetKeyPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
