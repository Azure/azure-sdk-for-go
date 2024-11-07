// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

type getSecretTestOptions struct{}

var getSecretTestOpts getSecretTestOptions = getSecretTestOptions{}

type getSecretTest struct {
	perf.PerfTestOptions
	secretName string
	client     *azsecrets.Client
}

// newGetSecretTest is called once per process
func newGetSecretTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	d := &getSecretTest{
		PerfTestOptions: options,
		secretName:      "livekvtestgetsecretperfsecret",
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_KEYVAULT_URL' could not be found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, &azsecrets.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: options.Transporter,
		},
	})
	if err != nil {
		return nil, err
	}

	_, err = client.SetSecret(ctx, d.secretName, azsecrets.SetSecretParameters{Value: to.Ptr("secret-value")}, nil)
	if err != nil {
		return nil, err
	}

	d.client = client
	return d, nil
}

func (gct *getSecretTest) GlobalCleanup(ctx context.Context) error {
	_, err := gct.client.DeleteSecret(ctx, gct.secretName, nil)
	if err != nil {
		return err
	}

	_, err = gct.client.PurgeDeletedSecret(ctx, gct.secretName, nil)
	return err
}

type getSecretPerfTest struct {
	client     *azsecrets.Client
	secretName string
}

// NewPerfTest is called once per goroutine
func (gct *getSecretTest) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	return &getSecretPerfTest{
		client:     gct.client,
		secretName: gct.secretName,
	}, nil
}

func (gcpt *getSecretPerfTest) Run(ctx context.Context) error {
	_, err := gcpt.client.GetSecret(ctx, gcpt.secretName, "", nil)
	return err
}

func (*getSecretPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
