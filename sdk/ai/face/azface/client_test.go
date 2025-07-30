// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azface_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/face/azface"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

const sampleImageURL = "https://raw.githubusercontent.com/Azure-Samples/cognitive-services-sample-data-files/master/Face/images/detection1.jpg"

// downloadTestImage downloads a test image for face detection
func downloadTestImage(t *testing.T) io.ReadSeekCloser {
	resp, err := http.Get(sampleImageURL)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	
	// Read the image into memory and return as ReadSeekCloser
	imageData, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	require.NoError(t, err)
	
	return io.NopCloser(strings.NewReader(string(imageData))).(io.ReadSeekCloser)
}

func TestClient_DetectFromURL(t *testing.T) {
	client := newClientForTestWithRecording(t)
	
	// Test detect with basic options
	resp, err := client.DetectFromURL(context.Background(), sampleImageURL, &azface.ClientDetectFromURLOptions{
		DetectionModel: to.Ptr(azface.FaceDetectionModelDetection03),
		RecognitionModel: to.Ptr(azface.FaceRecognitionModelRecognition04),
		ReturnFaceAttributes: []azface.FaceAttributeType{
			azface.FaceAttributeTypeAge,
			azface.FaceAttributeTypeGlasses,
		},
		ReturnFaceID: to.Ptr(true),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Value)
	
	// Verify that we detected at least one face
	require.GreaterOrEqual(t, len(resp.Value), 1)
	
	face := resp.Value[0]
	require.NotNil(t, face.FaceID)
	require.NotNil(t, face.FaceRectangle)
	require.NotNil(t, face.FaceAttributes)
	
	// Check face rectangle has valid coordinates
	require.Greater(t, *face.FaceRectangle.Width, int32(0))
	require.Greater(t, *face.FaceRectangle.Height, int32(0))
	
	// Check that requested attributes are present
	require.NotNil(t, face.FaceAttributes.Age)
	require.NotNil(t, face.FaceAttributes.Glasses)
}

func TestClient_Detect(t *testing.T) {
	client := newClientForTestWithRecording(t)
	
	imageData := downloadTestImage(t)
	defer imageData.Close()
	
	// Test detect with image data
	resp, err := client.Detect(context.Background(), imageData, &azface.ClientDetectOptions{
		DetectionModel: to.Ptr(azface.FaceDetectionModelDetection03),
		RecognitionModel: to.Ptr(azface.FaceRecognitionModelRecognition04),
		ReturnFaceAttributes: []azface.FaceAttributeType{
			azface.FaceAttributeTypeAge,
			azface.FaceAttributeTypeSmile,
		},
		ReturnFaceID: to.Ptr(true),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Value)
	
	// Verify that we detected at least one face
	require.GreaterOrEqual(t, len(resp.Value), 1)
	
	face := resp.Value[0]
	require.NotNil(t, face.FaceID)
	require.NotNil(t, face.FaceRectangle)
	require.NotNil(t, face.FaceAttributes)
	
	// Check that requested attributes are present
	require.NotNil(t, face.FaceAttributes.Age)
	require.NotNil(t, face.FaceAttributes.Smile)
}

func TestClient_FindSimilar(t *testing.T) {
	client := newClientForTestWithRecording(t)
	
	// First, detect faces to get face IDs
	resp1, err := client.DetectFromURL(context.Background(), sampleImageURL, &azface.ClientDetectFromURLOptions{
		DetectionModel: to.Ptr(azface.FaceDetectionModelDetection03),
		RecognitionModel: to.Ptr(azface.FaceRecognitionModelRecognition04),
		ReturnFaceID: to.Ptr(true),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp1.Value)
	
	queryFaceID := *resp1.Value[0].FaceID
	
	// Use the same face as candidate for simplicity in this test
	candidateFaceIDs := []string{queryFaceID}
	
	// Test FindSimilar
	resp2, err := client.FindSimilar(context.Background(), queryFaceID, candidateFaceIDs, &azface.ClientFindSimilarOptions{
		MaxNumOfCandidatesReturned: to.Ptr(int32(1)),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp2.FindSimilarResultArray)
	
	// Should find the exact same face with high confidence
	similar := resp2.FindSimilarResultArray[0]
	require.NotNil(t, similar.FaceID)
	require.Equal(t, queryFaceID, *similar.FaceID)
	require.NotNil(t, similar.Confidence)
	require.Greater(t, *similar.Confidence, float64(0.8)) // High confidence for same face
}

func TestClient_Group(t *testing.T) {
	client := newClientForTestWithRecording(t)
	
	// First, detect faces to get face IDs
	resp1, err := client.DetectFromURL(context.Background(), sampleImageURL, &azface.ClientDetectFromURLOptions{
		DetectionModel: to.Ptr(azface.FaceDetectionModelDetection03),
		RecognitionModel: to.Ptr(azface.FaceRecognitionModelRecognition04),
		ReturnFaceID: to.Ptr(true),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp1.Value)
	
	// Extract face IDs
	var faceIDs []string
	for _, face := range resp1.Value {
		if face.FaceID != nil {
			faceIDs = append(faceIDs, *face.FaceID)
		}
	}
	require.NotEmpty(t, faceIDs)
	
	// Test Group operation
	resp2, err := client.Group(context.Background(), faceIDs, nil)
	require.NoError(t, err)
	require.NotNil(t, resp2.Groups)
	
	// Verify the structure of the response
	totalFaces := 0
	for _, group := range resp2.Groups {
		totalFaces += len(group)
	}
	if resp2.MessyGroup != nil {
		totalFaces += len(resp2.MessyGroup)
	}
	
	// All input faces should be accounted for in the groups
	require.Equal(t, len(faceIDs), totalFaces)
}

func TestClient_VerifyFaceToFace(t *testing.T) {
	client := newClientForTestWithRecording(t)
	
	// First, detect faces to get face IDs
	resp1, err := client.DetectFromURL(context.Background(), sampleImageURL, &azface.ClientDetectFromURLOptions{
		DetectionModel: to.Ptr(azface.FaceDetectionModelDetection03),
		RecognitionModel: to.Ptr(azface.FaceRecognitionModelRecognition04),
		ReturnFaceID: to.Ptr(true),
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(resp1.Value), 1)
	
	faceID1 := *resp1.Value[0].FaceID
	
	// For this test, compare the face with itself (should have high confidence)
	resp2, err := client.VerifyFaceToFace(context.Background(), faceID1, faceID1, nil)
	require.NoError(t, err)
	require.NotNil(t, resp2.IsIdentical)
	require.NotNil(t, resp2.Confidence)
	
	// Verifying a face against itself should be identical with high confidence
	require.True(t, *resp2.IsIdentical)
	require.Greater(t, *resp2.Confidence, float64(0.8))
}