//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func TestNewClient(t *testing.T) {
	type args struct {
		endpoint     string
		credential   azcore.TokenCredential
		deploymentID string
		options      *ClientOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.endpoint, tt.args.credential, tt.args.deploymentID, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClientWithKeyCredential(t *testing.T) {
	type args struct {
		endpoint     string
		credential   KeyCredential
		deploymentID string
		options      *ClientOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClientWithKeyCredential(tt.args.endpoint, tt.args.credential, tt.args.deploymentID, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientWithKeyCredential() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientWithKeyCredential() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetCompletionsStream(t *testing.T) {
	body := CompletionsOptions{
		Prompt:      []*string{to.Ptr("What is Azure OpenAI?")},
		MaxTokens:   to.Ptr(int32(2048 - 127)),
		Temperature: to.Ptr(float32(0.0)),
	}
	cred := KeyCredential{APIKey: apiKey}
	deploymentID := "text-davinci-003"
	client, err := NewClientWithKeyCredential(endpoint, cred, deploymentID, newClientOptionsForTest(t))
	if err != nil {
		t.Errorf("NewClientWithKeyCredential() error = %v", err)
		return
	}
	response, err := client.GetCompletionsStream(context.TODO(), body, nil)
	if err != nil {
		t.Errorf("Client.GetCompletionsStream() error = %v", err)
		return
	}
	reader := response.Events
	defer reader.Close()

	var sb strings.Builder
	var eventCount int
	for {
		event, err := reader.Read()
		if err == io.EOF {
			break
		}
		eventCount++
		if err != nil {
			t.Errorf("reader.Read() error = %v", err)
			return
		}
		sb.WriteString(*event.Choices[0].Text)
	}
	got := sb.String()
	const want = "\n\nAzure OpenAI is a platform from Microsoft that provides access to OpenAI's artificial intelligence (AI) technologies. It enables developers to build, train, and deploy AI models in the cloud. Azure OpenAI provides access to OpenAI's powerful AI technologies, such as GPT-3, which can be used to create natural language processing (NLP) applications, computer vision models, and reinforcement learning models."
	if got != want {
		i := 0
		for i < len(got) && i < len(want) && got[i] == want[i] {
			i++
		}
		t.Errorf("Client.GetCompletionsStream() text[%d] = %c, want %c", i, got[i], want[i])
	}
	if eventCount != 86 {
		t.Errorf("Client.GetCompletionsStream() got = %v, want %v", eventCount, 1)
	}
}
