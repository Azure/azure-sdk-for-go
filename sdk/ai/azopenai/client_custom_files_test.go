// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/stretchr/testify/require"
)

func TestWriteMultipart(t *testing.T) {
	req, err := runtime.NewRequest(context.Background(), http.MethodPost, "https://example.com")
	require.NoError(t, err)

	fileContents := []byte("test file content")
	filename := "test.txt"
	purpose := FilePurpose("test-purpose")

	err = writeMultipart(req, fileContents, filename, purpose)
	require.NoError(t, err)

	body := new(bytes.Buffer)
	_, err = body.ReadFrom(req.Body())
	require.NoError(t, err)

	require.Contains(t, body.String(), "test file content")
	require.Contains(t, body.String(), "test-purpose")
	require.Contains(t, body.String(), filename)
}

func TestWriteMultipart_EmptyFileContents(t *testing.T) {
	req, err := runtime.NewRequest(context.Background(), http.MethodPost, "https://example.com")
	require.NoError(t, err)

	fileContents := []byte("")
	filename := "empty.txt"
	purpose := FilePurpose("test-purpose")

	err = writeMultipart(req, fileContents, filename, purpose)
	require.NoError(t, err)

	body := new(bytes.Buffer)
	_, err = body.ReadFrom(req.Body())
	require.NoError(t, err)

	require.Contains(t, body.String(), "test-purpose")
	require.Contains(t, body.String(), filename)
}

func TestCreateFormFile(t *testing.T) {
	tests := []struct {
		fieldname   string
		filename    string
		expectedErr bool
		expectedCT  string
	}{
		{"file", "test.txt", false, "text/plain"},
		{"file", "test.pdf", false, "application/pdf"},
		{"file", "test.unknown", false, "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			fileWriter, err := createFormFile(writer, tt.fieldname, tt.filename)
			if tt.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, fileWriter)

				// Close the writer to finalize the multipart message
				require.NoError(t, writer.Close())

				// Parse the multipart message
				reader := multipart.NewReader(body, writer.Boundary())
				part, err := reader.NextPart()
				require.NoError(t, err)
				partHeader := part.Header.Get("Content-Type")
				require.Equal(t, tt.expectedCT, partHeader)
			}
		})
	}
}

// Custom writer that always returns an error
type errorWriter struct{}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, io.ErrClosedPipe // Simulating a failure
}

func TestCreateFormFile_InvalidInputs(t *testing.T) {
	t.Run("NilWriterShouldPanic", func(t *testing.T) {
		require.Panics(t, func() {
			_, _ = createFormFile(nil, "field", "file.txt")
		}, "Expected panic when writer is nil")
	})

	t.Run("EmptyFieldAndFilename", func(t *testing.T) {
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		_, err := createFormFile(writer, "", "")
		require.NoError(t, err, "Expected no error for empty fieldname and filename")
	})

	t.Run("UnsupportedFileExtension", func(t *testing.T) {
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		_, err := createFormFile(writer, "field", "unknownfile.unknown")
		require.NoError(t, err, "Expected no error for unsupported file extension")
	})

	t.Run("InvalidWriterShouldReturnError", func(t *testing.T) {
		invalidWriter := multipart.NewWriter(&errorWriter{}) // Writer that always fails

		_, err := createFormFile(invalidWriter, "field", "file.txt")
		require.Error(t, err, "Expected error when using a failing writer")
		require.ErrorContains(t, err, "io: read/write on closed pipe", "Expected specific error message for closed pipe")
	})
}

func TestUploadFileCreateRequest(t *testing.T) {
	// Create a mock client struct that implements the formatURL method
	mockClient := &mockOpenAIClient{}
	// Initialize needed Client fields
	mockClient.endpoint = "https://test.openai.azure.com"

	// Test cases
	t.Run("SuccessWithDefaultFilename", func(t *testing.T) {
		// Setup
		fileContent := "test file content"
		file := strings.NewReader(fileContent)
		purpose := FilePurpose("test-purpose")
		ctx := context.Background()

		// Execute
		req, err := mockClient.uploadFileCreateRequest(ctx, file, purpose, nil)

		// Verify
		require.NoError(t, err)
		require.NotNil(t, req)
		require.Equal(t, http.MethodPost, req.Raw().Method)
		require.Equal(t, []string{"application/json"}, req.Raw().Header["Accept"])

		// Check body
		body := new(bytes.Buffer)
		_, err = body.ReadFrom(req.Body())
		require.NoError(t, err)

		// Validate multipart body has expected content
		bodyStr := body.String()
		require.Contains(t, bodyStr, "test file content")
		require.Contains(t, bodyStr, "file.txt")
		require.Contains(t, bodyStr, "test-purpose")
	})

	t.Run("SuccessWithCustomFilename", func(t *testing.T) {
		// Setup
		fileContent := "test file content"
		file := strings.NewReader(fileContent)
		purpose := FilePurpose("test-purpose")
		customFilename := "custom.txt"
		options := &UploadFileOptions{
			Filename: &customFilename,
		}
		ctx := context.Background()

		// Execute
		req, err := mockClient.uploadFileCreateRequest(ctx, file, purpose, options)

		// Verify
		require.NoError(t, err)
		require.NotNil(t, req)

		// Check body
		body := new(bytes.Buffer)
		_, err = body.ReadFrom(req.Body())
		require.NoError(t, err)

		// Validate multipart body contains custom filename
		bodyStr := body.String()
		require.Contains(t, bodyStr, "custom.txt")
	})

	t.Run("ErrorReadingFile", func(t *testing.T) {
		// Setup - create a ReadSeeker that returns an error
		errorFile := &errorReadSeeker{err: io.ErrClosedPipe}
		purpose := FilePurpose("test-purpose")
		ctx := context.Background()

		// Execute
		req, err := mockClient.uploadFileCreateRequest(ctx, errorFile, purpose, nil)

		// Verify
		require.Error(t, err)
		require.Nil(t, req)
		require.ErrorIs(t, err, io.ErrClosedPipe)
	})
}

// Mock client for testing
type mockOpenAIClient struct {
	Client // Embed the real Client
}

// Mock ReadSeeker that returns an error
type errorReadSeeker struct {
	err error
}

func (e *errorReadSeeker) Read(p []byte) (n int, err error) {
	return 0, e.err
}

func (e *errorReadSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, e.err
}
