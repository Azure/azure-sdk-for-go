// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
)

type getCertificatesTestOptions struct{}

var getCertTestOpts getCertificatesTestOptions = getCertificatesTestOptions{}

type getCertificateTest struct {
	perf.PerfTestOptions
	certificateName string
	client          *azcertificates.Client
}

// newGetCertificateTest is called once per process
func newGetCertificateTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	certName := "livekvtestgetcertperfcert"
	d := &getCertificateTest{
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

	client, err := azcertificates.NewClient(vaultURL, cred, &azcertificates.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: options.Transporter,
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = client.CreateCertificate(ctx, d.certificateName, azcertificates.CreateCertificateParameters{}, nil)
	if err != nil {
		return nil, err
	}

	d.client = client
	return d, nil
}

func (gct *getCertificateTest) GlobalCleanup(ctx context.Context) error {
	_, err := gct.client.DeleteCertificate(ctx, gct.certificateName, nil)
	if err != nil {
		return err
	}

	_, err = gct.client.PurgeDeletedCertificate(ctx, gct.certificateName, nil)
	return err
}

type getCertificatePerfTest struct {
	*getCertificateTest
	perf.PerfTestOptions
	client          *azcertificates.Client
	certificateName string
}

// NewPerfTest is called once per goroutine
func (gct *getCertificateTest) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	return &getCertificatePerfTest{
		getCertificateTest: gct,
		PerfTestOptions:    *options,
		client:             gct.client,
		certificateName:    gct.certificateName,
	}, nil
}

func (gcpt *getCertificatePerfTest) Run(ctx context.Context) error {
	_, err := gcpt.client.GetCertificate(ctx, gcpt.certificateName, "", nil)
	return err
}

func (*getCertificatePerfTest) Cleanup(ctx context.Context) error {
	return nil
}
