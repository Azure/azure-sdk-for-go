// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"crypto/sha256"
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

type signTestOptions struct{}

var signTestOpts signTestOptions = signTestOptions{}

type signTest struct {
	perf.PerfTestOptions
	keyName      string
	client       *azkeys.Client
	cryptoClient *crypto.Client
	signAlg      crypto.SignatureAlg
	digest       []byte
}

// newSignTest is called once per process
func newSignTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	d := &signTest{
		PerfTestOptions: options,
		keyName:         "livekvtestsignperfkey",
		signAlg:         crypto.SignatureAlgRS256,
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

	d.cryptoClient = cryptoClient
	d.client = client

	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return nil, err
	}
	hasher := sha256.New()
	d.digest = hasher.Sum(b)
	return d, nil
}

func (gct *signTest) GlobalCleanup(ctx context.Context) error {
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

type signPerfTest struct {
	cryptoClient *crypto.Client
	alg          crypto.SignatureAlg
	digest       []byte
}

// NewPerfTest is called once per goroutine
func (gct *signTest) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	return &signPerfTest{
		alg:    gct.signAlg,
		digest: gct.digest,
	}, nil
}

func (gcpt *signPerfTest) Run(ctx context.Context) error {
	_, err := gcpt.cryptoClient.Sign(ctx, gcpt.alg, gcpt.digest, nil)
	return err
}

func (*signPerfTest) Cleanup(ctx context.Context) error {
	return nil
}
