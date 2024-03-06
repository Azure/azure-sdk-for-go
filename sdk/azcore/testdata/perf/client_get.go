// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

type clientGETTestOptions struct {
	url string
}

var clientGetOpts clientGETTestOptions = clientGETTestOptions{url: ""}

// sleepTestRegister is called once per process
func clientTestRegister() {
	flag.StringVar(&clientGetOpts.url, "url", "", "URL to send a GET request")
}

type globalClientGETTest struct {
	perf.PerfTestOptions
	req policy.Request
}

func newClientGETTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	if clientGetOpts.url == "" {
		fmt.Println("--url/-u flag is required")
		return nil, errors.New("--url/-u flag is required")
	}
	req, err := runtime.NewRequest(ctx, "GET", clientGetOpts.url)
	if err != nil {
		return nil, err
	}
	return &globalClientGETTest{
		PerfTestOptions: options,
		req:             *req,
	}, nil
}

func (g *globalClientGETTest) GlobalCleanup(ctx context.Context) error {
	return nil
}

type clientGETTest struct {
	pipeline runtime.Pipeline
	req      policy.Request
}

func (g *globalClientGETTest) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	pipeline := runtime.NewPipeline("perf", "0.1.0", runtime.PipelineOptions{}, &policy.ClientOptions{
		Transport: options.Transporter,
	})

	return &clientGETTest{
		pipeline: pipeline,
		req:      g.req,
	}, nil
}

func (g *clientGETTest) Run(ctx context.Context) error {
	_, err := g.pipeline.Do(&g.req)
	return err
}

func (s *clientGETTest) Cleanup(ctx context.Context) error {
	return nil
}
