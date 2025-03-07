// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"bytes"
	"context"
	"net/http"
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
