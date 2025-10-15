// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/perf"
)

type mockJsonTestOptions struct {
	count uint
}

var mockJsonOpts mockJsonTestOptions = mockJsonTestOptions{count: 25}

// sleepTestRegister is called once per process
func mockJsonTestRegister() {
	flag.UintVar(&mockJsonOpts.count, "count", 25, "Number of items per page")
}

type globalMockJsonTest struct {
	perf.PerfTestOptions
	body []byte
}

func NewMockJsonTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	list := List{
		Name: to.Ptr("t0123456789abcdef"),
		Container: &ListItemsContainer{
			Items: make([]*ListItems, mockJsonOpts.count),
		},
	}
	now := time.Now()
	for i := range mockJsonOpts.count {
		name := fmt.Sprintf("testItem%d", i)
		hash := md5.Sum([]byte(name))
		list.Container.Items[i] = &ListItems{
			Name: to.Ptr(name),
			Properties: &ListItemProperties{
				ETag:         to.Ptr(azcore.ETag(fmt.Sprint(i))),
				CreationTime: to.Ptr(now),
				LastModified: to.Ptr(now),
				ContentMD5:   hash[:],
			},
		}
	}
	body, err := json.Marshal(&list)
	if err != nil {
		return nil, err
	}
	return &globalMockJsonTest{
		PerfTestOptions: options,
		body:            body,
	}, nil
}

func (g *globalMockJsonTest) GlobalCleanup(ctx context.Context) error {
	return nil
}

type mockJsonTest struct {
	pipeline runtime.Pipeline
}

func (g *globalMockJsonTest) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	pipeline := runtime.NewPipeline("perf", "0.1.0", runtime.PipelineOptions{}, &policy.ClientOptions{
		Transport: g,
	})
	return &mockJsonTest{
		pipeline: pipeline,
	}, nil
}

func (g *globalMockJsonTest) Do(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(g.body)),
	}, nil
}

func (g *mockJsonTest) Run(ctx context.Context) error {
	req, err := runtime.NewRequest(ctx, "GET", "https://contoso.com/containers/t0123456789abcdef?api-version=2025-10-15")
	if err != nil {
		return err
	}
	resp, err := g.pipeline.Do(req)
	if err != nil {
		return nil
	}
	// Make sure we deserialize the response.
	result := List{}
	return runtime.UnmarshalAsJSON(resp, &result)
}

func (s *mockJsonTest) Cleanup(ctx context.Context) error {
	return nil
}
