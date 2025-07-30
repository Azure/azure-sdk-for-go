//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azface

import (
	"context"
	"io"
	
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ClientOptions contains the optional parameters for the NewClient method.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient creates a new instance of Client with the specified values.
//   - endpoint - Face service endpoint URL
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - client options, pass nil to accept the default values.
func NewClient(endpoint string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			runtime.NewBearerTokenPolicy(credential, []string{"https://cognitiveservices.azure.com/.default"}, nil),
		},
	}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	return &Client{
		internal:   azcoreClient,
		endpoint:   endpoint,
		apiVersion: VersionsV12Preview1,
	}, nil
}

// NewClientWithKey creates a new instance of Client with the specified endpoint and API key.
//   - endpoint - Face service endpoint URL
//   - apiKey - subscription key for Face service
//   - options - client options, pass nil to accept the default values.
func NewClientWithKey(endpoint string, apiKey string, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	keyCredential := azcore.NewKeyCredential(apiKey)
	keyPolicy := runtime.NewKeyCredentialPolicy(keyCredential, "Ocp-Apim-Subscription-Key", nil)
	
	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{keyPolicy},
	}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	return &Client{
		internal:   azcoreClient,
		endpoint:   endpoint,
		apiVersion: VersionsV12Preview1,
	}, nil
}

// AdministrationClientOptions contains the optional parameters for the NewAdministrationClient method.
type AdministrationClientOptions struct {
	azcore.ClientOptions
}

// NewAdministrationClient creates a new instance of AdministrationClient with the specified values.
//   - endpoint - Face service endpoint URL
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - client options, pass nil to accept the default values.
func NewAdministrationClient(endpoint string, credential azcore.TokenCredential, options *AdministrationClientOptions) (*AdministrationClient, error) {
	if options == nil {
		options = &AdministrationClientOptions{}
	}

	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			runtime.NewBearerTokenPolicy(credential, []string{"https://cognitiveservices.azure.com/.default"}, nil),
		},
	}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	return &AdministrationClient{
		internal:   azcoreClient,
		endpoint:   endpoint,
		apiVersion: VersionsV12Preview1,
	}, nil
}

// NewAdministrationClientWithKey creates a new instance of AdministrationClient with the specified endpoint and API key.
//   - endpoint - Face service endpoint URL
//   - apiKey - subscription key for Face service
//   - options - client options, pass nil to accept the default values.
func NewAdministrationClientWithKey(endpoint string, apiKey string, options *AdministrationClientOptions) (*AdministrationClient, error) {
	if options == nil {
		options = &AdministrationClientOptions{}
	}

	keyCredential := azcore.NewKeyCredential(apiKey)
	keyPolicy := runtime.NewKeyCredentialPolicy(keyCredential, "Ocp-Apim-Subscription-Key", nil)
	
	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{keyPolicy},
	}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	return &AdministrationClient{
		internal:   azcoreClient,
		endpoint:   endpoint,
		apiVersion: VersionsV12Preview1,
	}, nil
}

// SessionClientOptions contains the optional parameters for the NewSessionClient method.
type SessionClientOptions struct {
	azcore.ClientOptions
}

// NewSessionClient creates a new instance of SessionClient with the specified values.
//   - endpoint - Face service endpoint URL
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - client options, pass nil to accept the default values.
func NewSessionClient(endpoint string, credential azcore.TokenCredential, options *SessionClientOptions) (*SessionClient, error) {
	if options == nil {
		options = &SessionClientOptions{}
	}

	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{
			runtime.NewBearerTokenPolicy(credential, []string{"https://cognitiveservices.azure.com/.default"}, nil),
		},
	}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	return &SessionClient{
		internal:   azcoreClient,
		endpoint:   endpoint,
		apiVersion: VersionsV12Preview1,
	}, nil
}

// NewSessionClientWithKey creates a new instance of SessionClient with the specified endpoint and API key.
//   - endpoint - Face service endpoint URL
//   - apiKey - subscription key for Face service
//   - options - client options, pass nil to accept the default values.
func NewSessionClientWithKey(endpoint string, apiKey string, options *SessionClientOptions) (*SessionClient, error) {
	if options == nil {
		options = &SessionClientOptions{}
	}

	keyCredential := azcore.NewKeyCredential(apiKey)
	keyPolicy := runtime.NewKeyCredentialPolicy(keyCredential, "Ocp-Apim-Subscription-Key", nil)
	
	azcoreClient, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{keyPolicy},
	}, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	return &SessionClient{
		internal:   azcoreClient,
		endpoint:   endpoint,
		apiVersion: VersionsV12Preview1,
	}, nil
}

// ClientDetectOptions contains the optional parameters for the Client.Detect method.
type ClientDetectOptions struct {
	// The detection model for the face. Defaults to Detection_01.
	DetectionModel *FaceDetectionModel

	// The recognition model for the face. Defaults to Recognition_04.
	RecognitionModel *FaceRecognitionModel

	// Analyze and return the one or more specified face attributes.
	ReturnFaceAttributes []FaceAttributeType

	// Return faceIds of the detected faces. Defaults to true.
	ReturnFaceID *bool

	// Return face landmarks of the detected faces. Defaults to false.
	ReturnFaceLandmarks *bool

	// Return recognition model name. Defaults to false.
	ReturnRecognitionModel *bool
}

// ClientDetectFromURLOptions contains the optional parameters for the Client.DetectFromURL method.
type ClientDetectFromURLOptions struct {
	// The detection model for the face. Defaults to Detection_01.
	DetectionModel *FaceDetectionModel

	// The recognition model for the face. Defaults to Recognition_04.
	RecognitionModel *FaceRecognitionModel

	// Analyze and return the one or more specified face attributes.
	ReturnFaceAttributes []FaceAttributeType

	// Return faceIds of the detected faces. Defaults to true.
	ReturnFaceID *bool

	// Return face landmarks of the detected faces. Defaults to false.
	ReturnFaceLandmarks *bool

	// Return recognition model name. Defaults to false.
	ReturnRecognitionModel *bool
}

// ClientDetectResponse contains the response from method Client.Detect.
type ClientDetectResponse struct {
	// Array of detected faces.
	Value []DetectionResult
}

// ClientDetectFromURLResponse contains the response from method Client.DetectFromURL.
type ClientDetectFromURLResponse struct {
	// Array of detected faces.
	Value []DetectionResult
}

// Detect - Detect human faces in an image, return face rectangles, and optionally with faceIds, landmarks, and attributes.
//
// Please refer to https://learn.microsoft.com/rest/api/face/face-detection-operations/detect for more details.
// If the operation fails it returns an *azcore.ResponseError type.
//   - imageContent - The input image content.
//   - options - ClientDetectOptions contains the optional parameters for the Client.Detect method.
func (client *Client) Detect(ctx context.Context, imageContent io.ReadSeekCloser, options *ClientDetectOptions) (ClientDetectResponse, error) {
	if options == nil {
		options = &ClientDetectOptions{}
	}

	internalOptions := &clientdetectOptions{
		DetectionModel:         options.DetectionModel,
		RecognitionModel:       options.RecognitionModel,
		ReturnFaceAttributes:   options.ReturnFaceAttributes,
		ReturnFaceID:           options.ReturnFaceID,
		ReturnFaceLandmarks:    options.ReturnFaceLandmarks,
		ReturnRecognitionModel: options.ReturnRecognitionModel,
	}

	result, err := client.detect(ctx, imageContent, internalOptions)
	if err != nil {
		return ClientDetectResponse{}, err
	}

	return ClientDetectResponse{Value: result.DetectionResultArray}, nil
}

// DetectFromURL - Detect human faces in an image, return face rectangles, and optionally with faceIds, landmarks, and attributes.
//
// Please refer to https://learn.microsoft.com/rest/api/face/face-detection-operations/detect-from-url for more details.
// If the operation fails it returns an *azcore.ResponseError type.
//   - urlParam - URL of input image.
//   - options - ClientDetectFromURLOptions contains the optional parameters for the Client.DetectFromURL method.
func (client *Client) DetectFromURL(ctx context.Context, urlParam string, options *ClientDetectFromURLOptions) (ClientDetectFromURLResponse, error) {
	if options == nil {
		options = &ClientDetectFromURLOptions{}
	}

	internalOptions := &clientdetectFromURLOptions{
		DetectionModel:         options.DetectionModel,
		RecognitionModel:       options.RecognitionModel,
		ReturnFaceAttributes:   options.ReturnFaceAttributes,
		ReturnFaceID:           options.ReturnFaceID,
		ReturnFaceLandmarks:    options.ReturnFaceLandmarks,
		ReturnRecognitionModel: options.ReturnRecognitionModel,
	}

	result, err := client.detectFromURL(ctx, urlParam, internalOptions)
	if err != nil {
		return ClientDetectFromURLResponse{}, err
	}

	return ClientDetectFromURLResponse{Value: result.DetectionResultArray}, nil
}