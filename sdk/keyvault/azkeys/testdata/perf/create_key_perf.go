// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
)

type azkeysPerf struct {
	perf.PerfTestOptions
	client   *azkeys.Client
	keyName  string
	vaultURL string
}

func (a *azkeysPerf) GetMetadata() perf.PerfTestOptions {
	return a.PerfTestOptions
}

func (a *azkeysPerf) GlobalSetup(ctx context.Context) error {
	cred, err := azidentity.NewClientSecretCredential(
		os.Getenv("AZKEYS_TENANT_ID"),
		os.Getenv("AZKEYS_CLIENT_ID"),
		os.Getenv("AZKEYS_CLIENT_SECRET"),
		nil,
	)
	if err != nil {
		return err
	}

	options := &azkeys.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: a.ProxyInstance,
		},
	}

	c, err := azkeys.NewClient(a.vaultURL, cred, options)
	if err != nil {
		return err
	}

	_, err = c.CreateRSAKey(context.Background(), a.keyName, nil)
	return err
}

func (a *azkeysPerf) Setup(ctx context.Context) error {
	cred, err := azidentity.NewClientSecretCredential(
		os.Getenv("AZKEYS_TENANT_ID"),
		os.Getenv("AZKEYS_CLIENT_ID"),
		os.Getenv("AZKEYS_CLIENT_SECRET"),
		nil,
	)
	if err != nil {
		return err
	}

	options := &azkeys.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: a.ProxyInstance,
		},
	}

	c, err := azkeys.NewClient(a.vaultURL, cred, options)
	if err != nil {
		return err
	}
	a.client = c
	return nil
}

func (a *azkeysPerf) Run(ctx context.Context) error {
	_, e := a.client.GetKey(ctx, a.keyName, nil)
	return e
}

func (a *azkeysPerf) Cleanup(ctx context.Context) error {
	return nil
}

func (a *azkeysPerf) GlobalCleanup(ctx context.Context) error {
	resp, err := a.client.BeginDeleteKey(ctx, a.keyName, nil)
	if err != nil {
		return err
	}
	_, err = resp.PollUntilDone(ctx, time.Second)
	if err != nil {
		return err
	}
	_, err = a.client.PurgeDeletedKey(ctx, a.keyName, nil)
	return err
}

func NewCreateKeyTest(options *perf.PerfTestOptions) perf.PerfTest {
	if options == nil {
		options = &perf.PerfTestOptions{}
	}
	options.Name = "CreateKeyTest"

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		panic(errors.New("could not find 'AZURE_KEYVAULT_URL' environment variable"))
	}

	return &azkeysPerf{
		PerfTestOptions: *options,
		keyName:         "createKeyTest",
		vaultURL:        vaultURL,
	}
}
