// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/file"
	"github.com/stretchr/testify/require"
)

// captureTransport records the last request it sees and returns a fixed 201 response,
// letting tests assert on the outgoing request without hitting the service.
type captureTransport struct {
	req *http.Request
}

func (c *captureTransport) Do(req *http.Request) (*http.Response, error) {
	c.req = req
	return &http.Response{
		Request:    req,
		Status:     "Created",
		StatusCode: http.StatusCreated,
		Header:     http.Header{},
		Body:       http.NoBody,
	}, nil
}

// Rename must percent-encode the source path in the x-ms-rename-source header. Otherwise
// source names containing spaces or non-ASCII characters are sent unencoded and rejected
// by the service with 400 InvalidSourceUri. See issues #23831 and #24369.
func TestFileRenameEncodesSourcePath(t *testing.T) {
	_require := require.New(t)
	ct := &captureTransport{}

	// Source URL whose path contains a non-ASCII character (ö) and a space, percent-encoded.
	srcURL := "https://fake.dfs.core.windows.net/myfs/dir1/l%C3%B6r%20006.jpg"
	fClient, err := file.NewClientWithNoCredential(srcURL, &file.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: ct},
	})
	_require.NoError(err)

	_, err = fClient.Rename(context.Background(), "dir1/renamed.jpg", nil)
	_require.NoError(err)
	_require.NotNil(ct.req)

	// The generated client sets the header via a raw (non-canonical) map key, so read it directly.
	renameSourceVals := ct.req.Header["x-ms-rename-source"]
	_require.NotEmpty(renameSourceVals)
	renameSource := renameSourceVals[0]
	_require.Contains(renameSource, "l%C3%B6r%20006.jpg")
	_require.NotContains(renameSource, " ")
	_require.NotContains(renameSource, "ö")
}
