//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azface

import "github.com/Azure/azure-sdk-for-go/sdk/azcore"

// ClientOptions contains optional settings for the Client.
type ClientOptions struct {
	azcore.ClientOptions
	// ServiceVersion specifies the version of the Face service to use.
	ServiceVersion *ServiceVersion
}

// Face represents a detected face
type Face struct {
	// FaceID - The unique ID of the detected face
	FaceID *string

	// FaceRectangle - The location of the detected face
	FaceRectangle *FaceRectangle

	// FaceAttributes - The attributes of the detected face
	FaceAttributes *FaceAttributes
}

// FaceRectangle represents the location of a face
type FaceRectangle struct {
	// Height - The height of the face rectangle
	Height *int32

	// Left - The left coordinate of the face rectangle
	Left *int32

	// Top - The top coordinate of the face rectangle
	Top *int32

	// Width - The width of the face rectangle
	Width *int32
}

// FaceAttributes represents attributes of a detected face
type FaceAttributes struct {
	// Age - The estimated age of the face
	Age *float64

	// Gender - The estimated gender of the face
	Gender *string

	// Emotion - The emotion scores for the face
	Emotion *EmotionScores

	// Smile - The smile intensity
	Smile *float64
}

// EmotionScores represents emotion scores for a face
type EmotionScores struct {
	// Anger - Anger score
	Anger *float64

	// Contempt - Contempt score
	Contempt *float64

	// Disgust - Disgust score
	Disgust *float64

	// Fear - Fear score
	Fear *float64

	// Happiness - Happiness score
	Happiness *float64

	// Neutral - Neutral score
	Neutral *float64

	// Sadness - Sadness score
	Sadness *float64

	// Surprise - Surprise score
	Surprise *float64
}

// DetectOptions contains options for face detection
type DetectOptions struct {
	// DetectionModel - The detection model to use
	DetectionModel *DetectionModel

	// RecognitionModel - The recognition model to use
	RecognitionModel *RecognitionModel

	// ReturnFaceAttributes - Whether to return face attributes
	ReturnFaceAttributes *bool

	// ReturnFaceID - Whether to return face IDs
	ReturnFaceID *bool
}