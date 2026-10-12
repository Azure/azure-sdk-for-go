// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package directory_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/directory"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
	"github.com/stretchr/testify/require"
)

const (
	fakeFileID    = "12384898975283830"
	fakeParentID  = "13835128424026341376"
	fakeShareURL  = "https://account.file.core.windows.net/share"
	fakeFileIDURL = fakeShareURL + "?fileid=" + fakeFileID
)

// fileIDTransport records the requests it handles and answers each with the configured response.
type fileIDTransport struct {
	requests []*http.Request
	status   int
	header   http.Header
}

func (f *fileIDTransport) Do(req *http.Request) (*http.Response, error) {
	f.requests = append(f.requests, req)
	status := f.status
	if status == 0 {
		status = http.StatusOK
	}
	header := f.header
	if header == nil {
		header = http.Header{}
	}
	return &http.Response{
		Request:    req,
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     header.Clone(),
		Body:       io.NopCloser(strings.NewReader("")),
	}, nil
}

func newFakeShareClient(t *testing.T, transport policy.Transporter) *share.Client {
	shareClient, err := share.NewClientWithNoCredential(fakeShareURL, &share.ClientOptions{ClientOptions: policy.ClientOptions{
		Transport: transport,
		Retry:     policy.RetryOptions{MaxRetries: -1},
	}})
	require.NoError(t, err)
	return shareClient
}

func TestNewDirectoryClientByFileIDURL(t *testing.T) {
	dirClient := newFakeShareClient(t, &fileIDTransport{}).NewDirectoryClientByFileID(fakeFileID)

	require.Equal(t, fakeFileIDURL, dirClient.URL())
	require.Equal(t, fakeFileID, dirClient.FileID())

	urlParts, err := sas.ParseURL(dirClient.URL())
	require.NoError(t, err)
	require.Equal(t, "share", urlParts.ShareName)
	require.Equal(t, fakeFileID, urlParts.FileID)
	require.Empty(t, urlParts.DirectoryOrFilePath)
}

func TestDirectoryClientFileIDFromConstructorURL(t *testing.T) {
	dirClient, err := directory.NewClientWithNoCredential(fakeFileIDURL, nil)
	require.NoError(t, err)
	require.Equal(t, fakeFileID, dirClient.FileID())
}

func TestDirectoryClientFileIDEmptyForPathBasedClient(t *testing.T) {
	dirClient, err := directory.NewClientWithNoCredential(fakeShareURL+"/directory", nil)
	require.NoError(t, err)
	require.Empty(t, dirClient.FileID())

	require.Empty(t, newFakeShareClient(t, &fileIDTransport{}).NewRootDirectoryClient().FileID())
}

func TestDirectoryClientByFileIDGetProperties(t *testing.T) {
	_require := require.New(t)
	transport := &fileIDTransport{header: http.Header{
		"X-Ms-File-Id":        []string{fakeFileID},
		"X-Ms-File-Parent-Id": []string{fakeParentID},
		"X-Ms-File-Name":      []string{"directory"},
	}}
	dirClient := newFakeShareClient(t, transport).NewDirectoryClientByFileID(fakeFileID)

	resp, err := dirClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(fakeFileID, *resp.ID)
	_require.Equal(fakeParentID, *resp.ParentID)
	_require.Equal("directory", *resp.FileName)

	_require.Len(transport.requests, 1)
	req := transport.requests[0]
	_require.Equal("/share", req.URL.Path)
	_require.Equal(fakeFileID, req.URL.Query().Get("fileid"))
	_require.Equal("directory", req.URL.Query().Get("restype"))
}

func TestDirectoryClientByFileIDGetPropertiesError(t *testing.T) {
	transport := &fileIDTransport{
		status: http.StatusNotFound,
		header: http.Header{"X-Ms-Error-Code": []string{string(fileerror.ResourceNotFound)}},
	}
	dirClient := newFakeShareClient(t, transport).NewDirectoryClientByFileID("11111111111111111111")

	_, err := dirClient.GetProperties(context.Background(), nil)
	require.Error(t, err)
	require.True(t, fileerror.HasCode(err, fileerror.ResourceNotFound))
}

func TestDirectoryClientByFileIDUnsupportedOperations(t *testing.T) {
	ctx := context.Background()
	transport := &fileIDTransport{}
	dirClient := newFakeShareClient(t, transport).NewDirectoryClientByFileID(fakeFileID)

	operations := map[string]func() error{
		"Create":        func() error { _, err := dirClient.Create(ctx, nil); return err },
		"Delete":        func() error { _, err := dirClient.Delete(ctx, nil); return err },
		"Rename":        func() error { _, err := dirClient.Rename(ctx, "destination", nil); return err },
		"SetProperties": func() error { _, err := dirClient.SetProperties(ctx, nil); return err },
		"SetMetadata":   func() error { _, err := dirClient.SetMetadata(ctx, nil); return err },
		"ForceCloseHandles": func() error {
			_, err := dirClient.ForceCloseHandles(ctx, "*", nil)
			return err
		},
		"ListHandles": func() error { _, err := dirClient.ListHandles(ctx, nil); return err },
		"NewListFilesAndDirectoriesPager": func() error {
			_, err := dirClient.NewListFilesAndDirectoriesPager(nil).NextPage(ctx)
			return err
		},
	}

	for name, operation := range operations {
		t.Run(name, func(t *testing.T) {
			err := operation()
			require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)
			require.Contains(t, err.Error(), "directory.Client."+name)
		})
	}
	require.Empty(t, transport.requests, "unsupported operations must fail before a request is sent")
}

func TestDirectoryClientByFileIDNavigation(t *testing.T) {
	// the path of a directory addressed by its file ID is not known, so clients navigated to from it are unusable
	ctx := context.Background()
	transport := &fileIDTransport{}
	dirClient := newFakeShareClient(t, transport).NewDirectoryClientByFileID(fakeFileID)

	_, err := dirClient.NewFileClient("file").GetProperties(ctx, nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)
	_, err = dirClient.NewFileClient("file").GetFileLinks(ctx, nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)
	_, err = dirClient.NewSubdirectoryClient("subdirectory").GetProperties(ctx, nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)
	_, err = dirClient.NewSubdirectoryClient("subdirectory").Create(ctx, nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)

	require.Empty(t, transport.requests)
}
