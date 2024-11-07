// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

type getListSecretsTestOptions struct {
	count int
}

var getListSecretsTestOpts getListSecretsTestOptions = getListSecretsTestOptions{
	count: 100,
}

func registerListSecrets() {
	flag.IntVar(&getListSecretsTestOpts.count, "count", 100, "number of secrets to create")
}

type listSecretsTest struct {
	perf.PerfTestOptions
	secretName string
	client     *azsecrets.Client
}

// newListSecretsTest is called once per process
func newListSecretsTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	d := &listSecretsTest{
		PerfTestOptions: options,
		secretName:      "livekvtestlistsecretperfsecret",
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

	for i := 0; i < getListSecretsTestOpts.count; i++ {
		_, err = client.SetSecret(ctx, fmt.Sprintf("%s%d", d.secretName, i), azsecrets.SetSecretParameters{Value: to.Ptr("secret-value")}, nil)
		if err != nil {
			return nil, err
		}
	}

	d.client = client
	return d, nil
}

func (gct *listSecretsTest) GlobalCleanup(ctx context.Context) error {
	for i := 0; i < getListSecretsTestOpts.count; i++ {
		_, err := gct.client.DeleteSecret(ctx, fmt.Sprintf("%s%d", gct.secretName, i), nil)
		if err != nil {
			return err
		}

		_, err = gct.client.PurgeDeletedSecret(ctx, gct.secretName, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

type listSecretsPerfTest struct {
	client     *azsecrets.Client
	secretName string
}

// NewPerfTest is called once per goroutine
func (gct *listSecretsTest) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	return &listSecretsPerfTest{
		client:     gct.client,
		secretName: gct.secretName,
	}, nil
}

func (gcpt *listSecretsPerfTest) Run(ctx context.Context) error {
	pager := gcpt.client.NewListSecretPropertiesPager(nil)
	for pager.More() {
		_, err := pager.NextPage(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (*listSecretsPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
