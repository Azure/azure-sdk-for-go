// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azface_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azface"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		endpoint    string
		expectError bool
	}{
		{
			name:        "ValidHTTPSEndpoint",
			endpoint:    "https://test.cognitiveservices.azure.com",
			expectError: false,
		},
		{
			name:        "ValidEndpointWithoutScheme",
			endpoint:    "test.cognitiveservices.azure.com",
			expectError: false,
		},
		{
			name:        "EmptyEndpoint",
			endpoint:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cred, err := azidentity.NewDefaultAzureCredential(nil)
			require.NoError(t, err)

			client, err := azface.NewClient(tt.endpoint, cred, nil)

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, client)
			} else {
				require.NoError(t, err)
				require.NotNil(t, client)
			}
		})
	}
}

func TestNewClientWithOptions(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)

	serviceVersion := azface.ServiceVersionV1_0
	options := &azface.ClientOptions{
		ServiceVersion: &serviceVersion,
	}

	client, err := azface.NewClient("test.cognitiveservices.azure.com", cred, options)
	require.NoError(t, err)
	require.NotNil(t, client)
}

func TestClientDetect(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)

	client, err := azface.NewClient("test.cognitiveservices.azure.com", cred, nil)
	require.NoError(t, err)

	ctx := context.Background()
	imageURL := "https://example.com/test-image.jpg"

	t.Run("DetectWithNilOptions", func(t *testing.T) {
		resp, err := client.Detect(ctx, imageURL, nil)
		require.NoError(t, err)
		require.NotNil(t, resp.Faces)
		require.Equal(t, 0, len(resp.Faces)) // Basic implementation returns empty array
	})

	t.Run("DetectWithEmptyOptions", func(t *testing.T) {
		options := &azface.DetectOptions{}
		resp, err := client.Detect(ctx, imageURL, options)
		require.NoError(t, err)
		require.NotNil(t, resp.Faces)
		require.Equal(t, 0, len(resp.Faces))
	})

	t.Run("DetectWithAllOptions", func(t *testing.T) {
		detectionModel := azface.DetectionModelDetection03
		recognitionModel := azface.RecognitionModelRecognition04
		returnFaceAttributes := true
		returnFaceID := true

		options := &azface.DetectOptions{
			DetectionModel:       &detectionModel,
			RecognitionModel:     &recognitionModel,
			ReturnFaceAttributes: &returnFaceAttributes,
			ReturnFaceID:         &returnFaceID,
		}

		resp, err := client.Detect(ctx, imageURL, options)
		require.NoError(t, err)
		require.NotNil(t, resp.Faces)
		require.Equal(t, 0, len(resp.Faces))
	})
}