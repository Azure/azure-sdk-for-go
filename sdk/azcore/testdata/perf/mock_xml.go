// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/xml"
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

type mockXmlTestOptions struct {
	count uint
}

var mockXmlOpts mockXmlTestOptions = mockXmlTestOptions{count: 25}

// sleepTestRegister is called once per process
func mockXmlTestRegister() {
	flag.UintVar(&mockXmlOpts.count, "count", 25, "Number of items per page")
}

type globalMockXmlTest struct {
	perf.PerfTestOptions
	body []byte
}

func NewMockXmlTest(ctx context.Context, options perf.PerfTestOptions) (perf.GlobalPerfTest, error) {
	list := List{
		Name: to.Ptr("t0123456789abcdef"),
		Container: &ListItemsContainer{
			Items: make([]*ListItems, mockXmlOpts.count),
		},
	}
	now := time.Now()
	for i := range mockXmlOpts.count {
		name := fmt.Sprintf("testBlob%d", i)
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
	body, err := xml.Marshal(&list)
	if err != nil {
		return nil, err
	}
	return &globalMockXmlTest{
		PerfTestOptions: options,
		body:            body,
	}, nil
}

func (g *globalMockXmlTest) GlobalCleanup(ctx context.Context) error {
	return nil
}

type mockXmlTest struct {
	pipeline runtime.Pipeline
}

func (g *globalMockXmlTest) NewPerfTest(ctx context.Context, options *perf.PerfTestOptions) (perf.PerfTest, error) {
	pipeline := runtime.NewPipeline("perf", "0.1.0", runtime.PipelineOptions{}, &policy.ClientOptions{
		Transport: g,
	})
	return &mockXmlTest{
		pipeline: pipeline,
	}, nil
}

func (g *globalMockXmlTest) Do(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(g.body)),
	}, nil
}

func (g *mockXmlTest) Run(ctx context.Context) error {
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
	return runtime.UnmarshalAsXML(resp, &result)
}

func (s *mockXmlTest) Cleanup(ctx context.Context) error {
	return nil
}
