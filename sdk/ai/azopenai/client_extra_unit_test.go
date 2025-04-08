// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientFormatURL(t *testing.T) {
	testCases := []struct {
		name       string
		client     Client
		path       string
		deployment *string
		expected   string
	}{
		{
			name: "Azure with deployment",
			client: Client{
				clientData: clientData{
					endpoint: "https://example.openai.azure.com",
					azure:    true,
				},
			},
			path:       "/completions",
			deployment: to.Ptr("gpt-35-turbo"),
			expected:   "https://example.openai.azure.com/openai/deployments/gpt-35-turbo/completions",
		},
		{
			name: "Azure with deployment",
			client: Client{
				clientData: clientData{
					endpoint: "https://example.openai.azure.com",
					azure:    true,
				},
			},
			path:       "/completions",
			deployment: to.Ptr("gpt\\-35-turbo"),
			expected:   "https://example.openai.azure.com/openai/deployments/gpt%5C-35-turbo/completions",
		},
		{
			name: "Azure without deployment",
			client: Client{
				clientData: clientData{
					endpoint: "https://example.openai.azure.com",
					azure:    true,
				},
			},
			path:       "/completions",
			deployment: nil,
			expected:   "https://example.openai.azure.com/openai/completions",
		},
		{
			name: "OpenAI with deployment",
			client: Client{
				clientData: clientData{
					endpoint: "https://api.openai.com/v1",
					azure:    false,
				},
			},
			path:       "/completions",
			deployment: to.Ptr("gpt-35-turbo"),
			expected:   "https://api.openai.com/v1/completions",
		},
		{
			name: "Image generation path",
			client: Client{
				clientData: clientData{
					endpoint: "https://example.openai.azure.com",
					azure:    true,
				},
			},
			path:       "/images/generations:submit",
			deployment: to.Ptr("dall-e-3"),
			expected:   "https://example.openai.azure.com/images/generations:submit",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.client.formatURL(tc.path, tc.deployment)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDeserializeAudioTranscription(t *testing.T) {
	// Test case 1: text/plain response
	t.Run("text/plain response", func(t *testing.T) {
		textContent := "This is a transcription"
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(textContent))),
			Header:     make(http.Header),
			Request: &http.Request{
				Body: http.NoBody,
			},
		}
		resp.Header.Set("Content-Type", "text/plain")

		result, err := deserializeAudioTranscription(resp)
		require.NoError(t, err)
		assert.Equal(t, textContent, *result.Text)
	})

	// Test case 2: JSON response
	t.Run("JSON response", func(t *testing.T) {
		jsonData := `{
            "text": "This is a JSON transcription",
            "duration": 10.5
        }`
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(jsonData))),
			Header:     make(http.Header),
			Request: &http.Request{
				Body: http.NoBody,
			},
		}
		resp.Header.Set("Content-Type", "application/json")

		result, err := deserializeAudioTranscription(resp)
		require.NoError(t, err)
		assert.Equal(t, "This is a JSON transcription", *result.Text)
		assert.Equal(t, float32(10.5), *result.Duration)
	})

	// Test case 3: Malformed JSON response
	t.Run("malformed JSON response", func(t *testing.T) {
		// Create a response with invalid JSON
		invalidJSON := `{"text": "This is incomplete JSON`
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(invalidJSON))),
			Header:     make(http.Header),
			Request: &http.Request{
				Body: http.NoBody,
			},
		}
		resp.Header.Set("Content-Type", "application/json")

		// Call the function
		_, err := deserializeAudioTranscription(resp)

		// Verify that the function returns an error
		require.Error(t, err)
	})
}
