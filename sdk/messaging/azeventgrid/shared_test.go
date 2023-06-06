//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid_test

import (
	"crypto/tls"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid"
	"github.com/stretchr/testify/require"
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

func requireEqualCloudEvent(t *testing.T, expected *azeventgrid.CloudEvent, actual *azeventgrid.CloudEvent) {
	t.Helper()

	require.NotEmpty(t, actual.ID, "ID is not empty")
	require.NotEmpty(t, actual.SpecVersion, "SpecVersion is not empty")

	expected.ID = actual.ID

	if expected.SpecVersion == nil {
		expected.SpecVersion = actual.SpecVersion
	}

	require.Equal(t, actual, expected)
}
