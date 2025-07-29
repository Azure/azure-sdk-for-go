// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azface_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azface"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ExampleNewClient() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := azface.NewClient("https://<resource-name>.cognitiveservices.azure.com", cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	_ = client
}

func ExampleNewClient_withOptions() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	serviceVersion := azface.ServiceVersionV1_0
	options := &azface.ClientOptions{
		ServiceVersion: &serviceVersion,
	}

	client, err := azface.NewClient("https://<resource-name>.cognitiveservices.azure.com", cred, options)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}
	_ = client
}

func ExampleClient_Detect() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := azface.NewClient("https://<resource-name>.cognitiveservices.azure.com", cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	ctx := context.TODO()
	imageURL := "https://example.com/image.jpg"

	resp, err := client.Detect(ctx, imageURL, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	log.Printf("Detected %d faces", len(resp.Faces))
}

func ExampleClient_Detect_withOptions() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	client, err := azface.NewClient("https://<resource-name>.cognitiveservices.azure.com", cred, nil)
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	ctx := context.TODO()
	imageURL := "https://example.com/image.jpg"

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
	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	for _, face := range resp.Faces {
		if face.FaceID != nil {
			log.Printf("Face ID: %s", *face.FaceID)
		}
		if face.FaceAttributes != nil && face.FaceAttributes.Age != nil {
			log.Printf("Estimated Age: %.1f", *face.FaceAttributes.Age)
		}
	}
}