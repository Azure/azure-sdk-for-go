// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azface_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azface"
	"github.com/stretchr/testify/require"
)

func TestServiceVersionValues(t *testing.T) {
	versions := azface.PossibleServiceVersionValues()
	require.NotEmpty(t, versions)
	require.Contains(t, versions, azface.ServiceVersionV1_0)
}

func TestDetectionModelValues(t *testing.T) {
	models := azface.PossibleDetectionModelValues()
	require.NotEmpty(t, models)
	require.Contains(t, models, azface.DetectionModelDetection01)
	require.Contains(t, models, azface.DetectionModelDetection02)
	require.Contains(t, models, azface.DetectionModelDetection03)
	require.Len(t, models, 3)
}

func TestRecognitionModelValues(t *testing.T) {
	models := azface.PossibleRecognitionModelValues()
	require.NotEmpty(t, models)
	require.Contains(t, models, azface.RecognitionModelRecognition01)
	require.Contains(t, models, azface.RecognitionModelRecognition02)
	require.Contains(t, models, azface.RecognitionModelRecognition03)
	require.Contains(t, models, azface.RecognitionModelRecognition04)
	require.Len(t, models, 4)
}

func TestDetectionModelConstants(t *testing.T) {
	require.Equal(t, "detection_01", string(azface.DetectionModelDetection01))
	require.Equal(t, "detection_02", string(azface.DetectionModelDetection02))
	require.Equal(t, "detection_03", string(azface.DetectionModelDetection03))
}

func TestRecognitionModelConstants(t *testing.T) {
	require.Equal(t, "recognition_01", string(azface.RecognitionModelRecognition01))
	require.Equal(t, "recognition_02", string(azface.RecognitionModelRecognition02))
	require.Equal(t, "recognition_03", string(azface.RecognitionModelRecognition03))
	require.Equal(t, "recognition_04", string(azface.RecognitionModelRecognition04))
}