// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/eng/tools/azgoperf/cmd/recording"
	"github.com/spf13/cobra"
)

var azkeysCmd = &cobra.Command{
	Use:   "azkeys",
	Short: "azkeys perf test",
	Long:  "azkeys perf test longer description",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(c *cobra.Command, args []string) error {
		return RunPerfTest(&azkeysPerf{})
	},
}

func init() {
	rootCmd.AddCommand(azkeysCmd)
	azkeysCmd.Flags().StringVarP(&TestProxy, "testproxy", "x", "", "whether to target http or https proxy (default is neither)")
}

type azkeysPerf struct {
	client  *azkeys.Client
	keyName string
}

func (a *azkeysPerf) GetMetadata() string {
	return "azkeysget"
}

func (a *azkeysPerf) GlobalSetup(ctx context.Context) error {
	vaultURL, ok := os.LookupEnv("AZKEYS_VAULT_URL")
	if !ok {
		return errors.New("could not find 'AZKEYS_VAULT_URL' environment variable")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return err
	}

	options := &azkeys.ClientOptions{}
	if TestProxy == "http" {
		t, err := recording.NewProxyTransport(&recording.TransportOptions{UseHTTPS: false, TestName: a.GetMetadata()})
		if err != nil {
			return err
		}
		options = &azkeys.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: t,
			},
		}
	} else if TestProxy == "https" {
		t, err := recording.NewProxyTransport(&recording.TransportOptions{UseHTTPS: false, TestName: a.GetMetadata()})
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

	_, err = a.client.CreateRSAKey(ctx, a.keyName, nil)
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
	_, err := a.client.BeginDeleteKey(ctx, a.keyName, nil)
	return err
}
