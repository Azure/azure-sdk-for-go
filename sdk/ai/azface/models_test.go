// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azface_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azface"
	"github.com/stretchr/testify/require"
)

func TestFaceModel(t *testing.T) {
	faceID := "test-face-id"
	height := int32(100)
	left := int32(50)
	top := int32(75)
	width := int32(90)
	age := 25.5
	gender := "male"
	smile := 0.8

	face := azface.Face{
		FaceID: &faceID,
		FaceRectangle: &azface.FaceRectangle{
			Height: &height,
			Left:   &left,
			Top:    &top,
			Width:  &width,
		},
		FaceAttributes: &azface.FaceAttributes{
			Age:    &age,
			Gender: &gender,
			Smile:  &smile,
		},
	}

	require.Equal(t, faceID, *face.FaceID)
	require.Equal(t, height, *face.FaceRectangle.Height)
	require.Equal(t, left, *face.FaceRectangle.Left)
	require.Equal(t, top, *face.FaceRectangle.Top)
	require.Equal(t, width, *face.FaceRectangle.Width)
	require.Equal(t, age, *face.FaceAttributes.Age)
	require.Equal(t, gender, *face.FaceAttributes.Gender)
	require.Equal(t, smile, *face.FaceAttributes.Smile)
}

func TestEmotionScores(t *testing.T) {
	anger := 0.1
	contempt := 0.05
	disgust := 0.02
	fear := 0.08
	happiness := 0.75
	neutral := 0.0
	sadness := 0.0
	surprise := 0.0

	emotion := azface.EmotionScores{
		Anger:     &anger,
		Contempt:  &contempt,
		Disgust:   &disgust,
		Fear:      &fear,
		Happiness: &happiness,
		Neutral:   &neutral,
		Sadness:   &sadness,
		Surprise:  &surprise,
	}

	require.Equal(t, anger, *emotion.Anger)
	require.Equal(t, contempt, *emotion.Contempt)
	require.Equal(t, disgust, *emotion.Disgust)
	require.Equal(t, fear, *emotion.Fear)
	require.Equal(t, happiness, *emotion.Happiness)
	require.Equal(t, neutral, *emotion.Neutral)
	require.Equal(t, sadness, *emotion.Sadness)
	require.Equal(t, surprise, *emotion.Surprise)
}

func TestDetectOptions(t *testing.T) {
	detectionModel := azface.DetectionModelDetection03
	recognitionModel := azface.RecognitionModelRecognition04
	returnFaceAttributes := true
	returnFaceID := false

	options := azface.DetectOptions{
		DetectionModel:       &detectionModel,
		RecognitionModel:     &recognitionModel,
		ReturnFaceAttributes: &returnFaceAttributes,
		ReturnFaceID:         &returnFaceID,
	}

	require.Equal(t, detectionModel, *options.DetectionModel)
	require.Equal(t, recognitionModel, *options.RecognitionModel)
	require.Equal(t, returnFaceAttributes, *options.ReturnFaceAttributes)
	require.Equal(t, returnFaceID, *options.ReturnFaceID)
}

func TestClientOptions(t *testing.T) {
	serviceVersion := azface.ServiceVersionV1_0
	options := azface.ClientOptions{
		ServiceVersion: &serviceVersion,
	}

	require.Equal(t, serviceVersion, *options.ServiceVersion)
}