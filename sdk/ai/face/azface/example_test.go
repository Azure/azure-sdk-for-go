// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azface_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/face/azface"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleNewClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := azface.NewClient("https://my-face-service.cognitiveservices.azure.com/", cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client // Use the client
}

func ExampleNewClientWithKey() {
	client, err := azface.NewClientWithKey(
		"https://my-face-service.cognitiveservices.azure.com/", 
		"my-subscription-key", 
		nil,
	)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = client // Use the client
}

func ExampleClient_DetectFromURL() {
	client, err := azface.NewClientWithKey(
		"https://my-face-service.cognitiveservices.azure.com/", 
		"my-subscription-key", 
		nil,
	)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	imageURL := "https://example.com/face.jpg"
	
	resp, err := client.DetectFromURL(context.TODO(), imageURL, &azface.ClientDetectFromURLOptions{
		DetectionModel: to.Ptr(azface.FaceDetectionModelDetection03),
		RecognitionModel: to.Ptr(azface.FaceRecognitionModelRecognition04),
		ReturnFaceAttributes: []azface.FaceAttributeType{
			azface.FaceAttributeTypeAge,
			azface.FaceAttributeTypeGlasses,
		},
		ReturnFaceID: to.Ptr(true),
	})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	for _, face := range resp.Value {
		if face.FaceAttributes != nil && face.FaceAttributes.Age != nil {
			log.Printf("Detected face with age: %f", *face.FaceAttributes.Age)
		}
	}
}

func ExampleClient_FindSimilar() {
	client, err := azface.NewClientWithKey(
		"https://my-face-service.cognitiveservices.azure.com/", 
		"my-subscription-key", 
		nil,
	)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	queryFaceID := "12345678-1234-1234-1234-123456789012"
	candidateFaceIDs := []string{
		"87654321-4321-4321-4321-210987654321",
		"11111111-2222-3333-4444-555555555555",
	}

	resp, err := client.FindSimilar(context.TODO(), queryFaceID, candidateFaceIDs, &azface.ClientFindSimilarOptions{
		MaxNumOfCandidatesReturned: to.Ptr(int32(2)),
	})
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	for _, similar := range resp.FindSimilarResultArray {
		if similar.Confidence != nil {
			log.Printf("Found similar face with confidence: %f", *similar.Confidence)
		}
	}
}

func ExampleClient_Group() {
	client, err := azface.NewClientWithKey(
		"https://my-face-service.cognitiveservices.azure.com/", 
		"my-subscription-key", 
		nil,
	)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	faceIDs := []string{
		"12345678-1234-1234-1234-123456789012",
		"87654321-4321-4321-4321-210987654321",
		"11111111-2222-3333-4444-555555555555",
	}

	resp, err := client.Group(context.TODO(), faceIDs, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	log.Printf("Found %d groups", len(resp.Groups))
	for i, group := range resp.Groups {
		log.Printf("Group %d has %d faces", i, len(group))
	}
}

func ExampleNewAdministrationClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	adminClient, err := azface.NewAdministrationClient("https://my-face-service.cognitiveservices.azure.com/", cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	_ = adminClient // Use the admin client
}