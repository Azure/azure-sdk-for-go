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
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates"
)

type getCertificatesTestOptions struct{}

var getCertTestOpts getCertificatesTestOptions = getCertificatesTestOptions{}

type GetCertificateTest struct {
	perf.PerfTestOptions
	certificateName string
	client          *azcertificates.Client
}

// NewGetCertificateTest is called once per process
func NewGetCertificateTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	certName := "livekvtestgetcertperfcert"
	d := &GetCertificateTest{
		PerfTestOptions: options,
		certificateName: certName,
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_KEYVAULT_URL' could not be found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azcertificates.NewClient(vaultURL, cred, nil)
	if err != nil {
		return nil, err
	}

	poller, err := client.BeginCreateCertificate(ctx, d.certificateName, azcertificates.NewDefaultCertificatePolicy(), nil)
	if err != nil {
		return nil, err
	}
	_, err = poller.PollUntilDone(ctx, 500*time.Millisecond)
	if err != nil {
		return nil, err
	}

	d.client = client
	return d, nil
}

func (d *GetCertificateTest) GlobalCleanup(ctx context.Context) error {
	poller, err := d.client.BeginDeleteCertificate(ctx, d.certificateName, nil)
	if err != nil {
		return err
	}

	_, err = poller.PollUntilDone(ctx, 500*time.Millisecond)
	return err
}

type GetCertificatePerfTest struct {
	*GetCertificateTest
	perf.PerfTestOptions
	client          *azcertificates.Client
	certificateName string
}

// NewPerfTest is called once per goroutine
func (g *GetCertificateTest) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	return &GetCertificatePerfTest{
		GetCertificateTest: g,
		PerfTestOptions:    *options,
		client:             g.client,
		certificateName:    g.certificateName,
	}, nil
}

func (d *GetCertificatePerfTest) Run(ctx context.Context) error {
	_, err := d.client.GetCertificate(ctx, d.certificateName, nil)
	return err
}

func (*GetCertificatePerfTest) Cleanup(ctx context.Context) error {
	return nil
}
