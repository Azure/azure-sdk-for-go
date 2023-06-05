//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid_test

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid"
)

type testVars struct {
	Key          string
	Endpoint     string
	Topic        string
	Subscription string

	KeyLogPath string
}

func loadEnv() testVars {
	key := os.Getenv("EVENTGRID_KEY")
	ep := os.Getenv("EVENTGRID_ENDPOINT")
	topic := os.Getenv("EVENTGRID_TOPIC")
	sub := os.Getenv("EVENTGRID_SUBSCRIPTION")

	// Setting this variable will cause the test clients to dump out the pre-master-key
	// for your HTTP connection. This allows you decrypt a packet capture from wireshark.
	//
	// If you want to do this just set SSLKEYLOGFILE_TEST env var to a path on disk and
	// Go will write out the key.
	keyLogFile := os.Getenv("SSLKEYLOGFILE_TEST")

	return testVars{
		Key:          key,
		Endpoint:     ep,
		Topic:        topic,
		Subscription: sub,
		KeyLogPath:   keyLogFile,
	}
}

type clientWrapper struct {
	*azeventgrid.Client
	TestVars   testVars
	keyLogFile *os.File
}

func (c clientWrapper) cleanup() {
	if c.keyLogFile != nil {
		c.keyLogFile.Close()
	}
}

func (c clientWrapper) receiveAllCloudEventsFull(ctx context.Context, numEvents int) ([]*azeventgrid.ReceiveDetails, error) {
	var received []*azeventgrid.ReceiveDetails

	for len(received) < numEvents {
		remaining := int32(numEvents - len(received))

		resp, err := c.ReceiveCloudEvents(ctx, c.TestVars.Topic, c.TestVars.Subscription, &azeventgrid.ReceiveCloudEventsOptions{
			MaxEvents:   to.Ptr(remaining),
			MaxWaitTime: to.Ptr[int32](10),
		})

		if errors.Is(err, context.Canceled) {
			break
		}

		if err != nil {
			return nil, err
		}

		received = append(received, resp.Value...)
	}

	return received, nil
}

func newClientForTest() clientWrapper {
	vars := loadEnv()

	var opts *azeventgrid.ClientOptions
	var keyLogWriter *os.File

	if vars.KeyLogPath != "" {
		tmpKeyLogWriter, err := os.OpenFile(vars.KeyLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)

		if err != nil {
			panic(err)
		}

		keyLogWriter = tmpKeyLogWriter

		tp := http.DefaultTransport.(*http.Transport).Clone()
		tp.TLSClientConfig = &tls.Config{
			KeyLogWriter: keyLogWriter,
		}

		opts = &azeventgrid.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: &http.Client{Transport: tp},
			},
		}
	}

	c, err := azeventgrid.NewClientWithSharedKeyCredential(vars.Endpoint, vars.Key, opts)

	if err != nil {
		panic(err)
	}

	return clientWrapper{
		Client:     c,
		TestVars:   vars,
		keyLogFile: keyLogWriter,
	}
}
