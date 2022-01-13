// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeysperf

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/azperf/internal/perf"
	"github.com/Azure/azure-sdk-for-go/eng/tools/azperf/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/spf13/cobra"
)

var AzkeysCmd = &cobra.Command{
	Use:   "CreateKeyTest",
	Short: "Create a single RSA key",
	Long:  "Create a RSA key using default options.",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(c *cobra.Command, args []string) error {
		return perf.RunPerfTest(&azkeysPerf{})
	},
}

type azkeysPerf struct {
	client  *azkeys.Client
	keyName string
}

func (a *azkeysPerf) GetMetadata() string {
	return "azkeyscreate"
}

func (a *azkeysPerf) GlobalSetup(ctx context.Context) error {
	err := recording.Start(a.GetMetadata(), nil)
	if err != nil {
		return err
	}
	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		return errors.New("could not find 'AZURE_KEYVAULT_URL' environment variable")
	}

	cred, err := azidentity.NewClientSecretCredential(
		os.Getenv("AZKEYS_TENANT_ID"),
		os.Getenv("AZKEYS_CLIENT_ID"),
		os.Getenv("AZKEYS_CLIENT_SECRET"),
		nil,
	)
	if err != nil {
		return err
	}

	options := &azkeys.ClientOptions{}
	if perf.TestProxy == "http" {
		t, err := recording.NewProxyTransport(&recording.TransportOptions{UseHTTPS: true, TestName: a.GetMetadata()})
		if err != nil {
			return err
		}
		options = &azkeys.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: t,
			},
		}
	} else if perf.TestProxy == "https" {
		t, err := recording.NewProxyTransport(&recording.TransportOptions{UseHTTPS: true, TestName: a.GetMetadata()})
		if err != nil {
			return err
		}
		options = &azkeys.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: t,
			},
		}
	}

	c, err := azkeys.NewClient(vaultURL, cred, options)
	if err != nil {
		return err
	}

	a.client = c
	a.keyName = "myKeyName"

	_, err = a.client.CreateRSAKey(context.Background(), a.keyName, nil)
	return err
}

func (a *azkeysPerf) Setup(ctx context.Context) error {
	return nil
}

func (a *azkeysPerf) Run(ctx context.Context) error {
	_, e := a.client.GetKey(ctx, a.keyName, nil)
	return e
}

func (a *azkeysPerf) TearDown(ctx context.Context) error {
	return nil
}

func (a *azkeysPerf) GlobalTearDown(ctx context.Context) error {
	defer recording.Stop(a.GetMetadata(), nil)
	resp, err := a.client.BeginDeleteKey(ctx, a.keyName, nil)
	if err != nil {return err}
	_, err = resp.PollUntilDone(ctx, time.Second)
	if err != nil {return err}
	_, err = a.client.PurgeDeletedKey(ctx, a.keyName, nil)
	return err
}
