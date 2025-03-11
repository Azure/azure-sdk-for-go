//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"bytes"
	"context"
	"fmt"
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
		fmt.Println("Error creating client:", err)
		return
	}

	// Create an upload request
	createUploadResp, err := client.CreateUpload(context.Background(), azopenai.CreateUploadRequest{
		Bytes:    to.Ptr(int32(10)),
		Filename: to.Ptr("test.txt"),
		MimeType: to.Ptr("text/plain"),
		Purpose:  to.Ptr(azopenai.CreateUploadRequestPurposeAssistants),
	}, nil)
	if err != nil {
		fmt.Println("Error creating upload:", err)
		return
	}

	// Create parts to upload
	part1 := streaming.NopCloser(strings.NewReader("hello"))
	part2 := streaming.NopCloser(strings.NewReader("world"))

	// Upload the second part
	part2Resp, err := client.AddUploadPart(context.Background(), *createUploadResp.ID, part2, nil)
	if err != nil {
		fmt.Println("Error uploading part 2:", err)
		return
	}

	// Upload the first part
	part1Resp, err := client.AddUploadPart(context.Background(), *createUploadResp.ID, part1, nil)
	if err != nil {
		fmt.Println("Error uploading part 1:", err)
		return
	}

	// Complete the upload by specifying the order of parts
	uploadResp, err := client.CompleteUpload(context.Background(), *createUploadResp.ID, azopenai.CompleteUploadRequest{
		PartIDs: []string{*part2Resp.ID, *part1Resp.ID},
	}, nil)
	if err != nil {
		fmt.Println("Error completing upload:", err)
		return
	}

	// Verify the total size of uploaded parts
	fmt.Println("Total size of uploaded parts:", *uploadResp.Bytes)

	// Delete the uploaded file
	_, err = client.DeleteFile(context.Background(), *uploadResp.File.ID, nil)
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return
	}

	// Output:
	// Total size of uploaded parts: 10
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
		fmt.Println("Error creating client:", err)
		return
	}

	// Upload a file
	uploadResp, err := client.UploadFile(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("hello world"))), azopenai.FilePurposeAssistants, nil)
	if err != nil {
		fmt.Println("Error uploading file:", err)
		return
	}

	// Cleanup: delete the uploaded file
	defer func() {
		_, err := client.DeleteFile(context.Background(), *uploadResp.ID, nil)
		if err != nil {
			fmt.Println("Error deleting file:", err)
		}
	}()

	// Get the uploaded file
	getFileResp, err := client.GetFile(context.Background(), *uploadResp.ID, nil)
	if err != nil {
		fmt.Println("Error getting file:", err)
		return
	}

	// Verify the purpose of the uploaded file
	fmt.Println("Purpose of uploaded file:", *getFileResp.Purpose)

	// List all files and verify the response is not empty
	filesResp, err := client.ListFiles(context.Background(), nil)
	if err != nil {
		fmt.Println("Error listing files:", err)
		return
	}
	fmt.Println("Number of files:", len(filesResp.Data))

	// Output:
	// Purpose of uploaded file: assistants
	// Number of files: 1
}
