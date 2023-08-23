//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package publisher_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type topicVars struct {
	Name     string
	Key      string
	Endpoint string
}

type eventGridVars struct {
	// EG are connection variables for an EventGrid encoded topic.
	EG topicVars
	// CE are connection variables for a CloudEvent encoded topic.
	CE topicVars
}

func newTestVars(t *testing.T) eventGridVars {
	egVars := eventGridVars{
		EG: topicVars{Name: os.Getenv("EVENTGRID_TOPIC_NAME"),
			Key:      os.Getenv("EVENTGRID_TOPIC_KEY"),
			Endpoint: os.Getenv("EVENTGRID_TOPIC_ENDPOINT"),
		},
		CE: topicVars{Name: os.Getenv("EVENTGRID_CE_TOPIC_NAME"),
			Key:      os.Getenv("EVENTGRID_CE_TOPIC_KEY"),
			Endpoint: os.Getenv("EVENTGRID_CE_TOPIC_ENDPOINT"),
		},
	}

	return egVars
}

type dumpFullPolicy struct {
	Prefix string
}

func (p dumpFullPolicy) Do(req *policy.Request) (*http.Response, error) {
	fmt.Printf("\n\n===> BEGIN: REQUEST (%s) <===\n\n", p.Prefix)

	requestBytes, err := httputil.DumpRequestOut(req.Raw(), false)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(requestBytes))
	fmt.Printf("Body: %s\n", string(FormatRequestBytes(req.Raw())))
	fmt.Printf("\n\n===> END: REQUEST (%s)<===\n\n", p.Prefix)

	resp, err := req.Next()

	if err != nil {
		return nil, err
	}

	fmt.Printf("\n\n===> BEGIN: RESPONSE (%s) <===\n\n", p.Prefix)

	responseBytes, err := httputil.DumpResponse(resp, false)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(responseBytes))
	fmt.Printf("Body: %s\n", string(FormatResponseBytes(resp)))

	fmt.Printf("\n\n===> END: RESPONSE (%s) <===\n\n", p.Prefix)
	return resp, err
}

func FormatRequestBytes(req *http.Request) []byte {
	requestBytes, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

	req.Body = io.NopCloser(bytes.NewBuffer(requestBytes))
	return FormatBytes(requestBytes)
}

func FormatResponseBytes(resp *http.Response) []byte {
	requestBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(requestBytes))
	return FormatBytes(requestBytes)
}

func FormatBytes(body []byte) []byte {
	var m *map[string]any
	var l *[]any

	candidates := []any{&m, &l}

	for _, v := range candidates {
		err := json.Unmarshal(body, v)

		if err != nil {
			continue
		}

		if err == nil {
			formattedBytes, err := json.MarshalIndent(v, "  ", "  ")

			if err != nil {
				continue
			}

			return formattedBytes
		}
	}

	// if we can't format it we'll just give it back.
	return body
}
