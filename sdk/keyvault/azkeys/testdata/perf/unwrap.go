// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/crypto"
)

type unwrapTestOptions struct{}

var unwrapTestOpts unwrapTestOptions = unwrapTestOptions{}

type unwrapTest struct {
	perf.PerfTestOptions
	keyName      string
	client       *azkeys.Client
	cryptoClient *crypto.Client
	alg          crypto.WrapAlg
	encryptedKey []byte
}

// newUnwrapTest is called once per process
func newUnwrapTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	d := &unwrapTest{
		PerfTestOptions: options,
		keyName:         "livekvtestunwrapperfkey",
		alg:             crypto.WrapAlgRSAOAEP256,
	}

	vaultURL, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		return nil, fmt.Errorf("the environment variable 'AZURE_KEYVAULT_URL' could not be found")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}

	client, err := azkeys.NewClient(vaultURL, cred, &azkeys.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: options.Transporter,
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.CreateRSAKey(ctx, d.keyName, &azkeys.CreateRSAKeyOptions{Size: to.Ptr(int32(2048))})
	if err != nil {
		return nil, err
	}

	cryptoClient, err := crypto.NewClient(*resp.ID, cred, &crypto.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: options.Transporter,
		},
	})
	if err != nil {
		return nil, err
	}

	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return nil, err
	}

	result, err := cryptoClient.WrapKey(ctx, d.alg, b, nil)
	if err != nil {
		return nil, err
	}

	d.encryptedKey = result.EncryptedKey
	d.cryptoClient = cryptoClient
	d.client = client
	return d, nil
}

func (gct *unwrapTest) GlobalCleanup(ctx context.Context) error {
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

type unwrapPerfTest struct {
	cryptoClient *crypto.Client
	alg          crypto.WrapAlg
	encryptedKey []byte
}

// NewPerfTest is called once per goroutine
func (gct *unwrapTest) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	return &unwrapPerfTest{
		cryptoClient: gct.cryptoClient,
		alg:          gct.alg,
		encryptedKey: gct.encryptedKey,
	}, nil
}

func (gcpt *unwrapPerfTest) Run(ctx context.Context) error {
	_, err := gcpt.cryptoClient.UnwrapKey(ctx, gcpt.alg, gcpt.encryptedKey, nil)
	return err
}

func (*unwrapPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
