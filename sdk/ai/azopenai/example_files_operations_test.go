//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func ExampleClient_AddUploadPart() {
	azureOpenAIKey := os.Getenv("AOAI_API_KEY")
	azureOpenAIEndpoint := os.Getenv("AOAI_ENDPOINT")

	if azureOpenAIKey == "" || azureOpenAIEndpoint == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

	// Create an upload request
	createUploadResp, err := client.CreateUpload(context.TODO(), azopenai.CreateUploadRequest{
		Bytes:    to.Ptr(int32(10)),
		Filename: to.Ptr("test.txt"),
		MimeType: to.Ptr("text/plain"),
		Purpose:  to.Ptr(azopenai.CreateUploadRequestPurposeAssistants),
	}, nil)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

	// Create parts to upload
	part1 := streaming.NopCloser(strings.NewReader("hello"))
	part2 := streaming.NopCloser(strings.NewReader("world"))

	// Upload the second part
	part2Resp, err := client.AddUploadPart(context.TODO(), *createUploadResp.ID, part2, nil)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

	// Upload the first part
	part1Resp, err := client.AddUploadPart(context.TODO(), *createUploadResp.ID, part1, nil)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

	// Complete the upload by assembling parts in desired sequence (part1 followed by part2)
	// Note: The order in this array determines the final file content order, not the upload order
	uploadResp, err := client.CompleteUpload(context.TODO(), *createUploadResp.ID, azopenai.CompleteUploadRequest{
		PartIDs: []string{*part1Resp.ID, *part2Resp.ID},
	}, nil)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

	// Verify the total size of uploaded parts
	fmt.Println("Total size of uploaded parts:", *uploadResp.Bytes)

	// Delete the uploaded file
	_, err = client.DeleteFile(context.TODO(), *uploadResp.File.ID, nil)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

	// Output:
}

func ExampleClient_UploadFile() {
	azureOpenAIKey := os.Getenv("AOAI_API_KEY")
	azureOpenAIEndpoint := os.Getenv("AOAI_ENDPOINT")

	if azureOpenAIKey == "" || azureOpenAIEndpoint == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

	// Upload a file
	uploadResp, err := client.UploadFile(context.TODO(), streaming.NopCloser(bytes.NewReader([]byte("hello world"))), azopenai.FilePurposeAssistants, nil)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

	// Cleanup: delete the uploaded file
	defer func() {
		_, err := client.DeleteFile(context.TODO(), *uploadResp.ID, nil)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
	}()

	// Get the uploaded file
	getFileResp, err := client.GetFile(context.TODO(), *uploadResp.ID, nil)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}

	// Verify the purpose of the uploaded file
	fmt.Println("Purpose of uploaded file:", *getFileResp.Purpose)

	// Output:
}
