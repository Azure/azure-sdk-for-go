// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/file"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/fileerror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/sas"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azfile/share"
	"github.com/stretchr/testify/require"
)

const (
	fakeFileID      = "12384898975283830"
	fakeShareURL    = "https://account.file.core.windows.net/share"
	fakeFileIDURL   = fakeShareURL + "?fileid=" + fakeFileID
	fakeShareSnap   = "2011-03-09T01:42:34.9360000Z"
	fakeParentID    = "13835128424026341376"
	fakeFileName    = "file.txt"
	fakeHardLinkXML = `<?xml version="1.0" encoding="utf-8"?><HardLinks>` +
		`<HardLink><ParentId>` + fakeParentID + `</ParentId><FileName>` + fakeFileName + `</FileName></HardLink>` +
		`<HardLink><ParentId>0</ParentId><FileName Encoded="true">link%EF%BF%BE</FileName></HardLink>` +
		`</HardLinks>`
)

// fileIDTransport records the requests it handles and answers each with the configured response.
type fileIDTransport struct {
	requests []*http.Request
	status   int
	header   http.Header
	body     string
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
		Body:       io.NopCloser(strings.NewReader(f.body)),
	}, nil
}

func fileIDClientOptions(transport policy.Transporter) *share.ClientOptions {
	return &share.ClientOptions{ClientOptions: policy.ClientOptions{
		Transport: transport,
		Retry:     policy.RetryOptions{MaxRetries: -1},
	}}
}

func newFakeShareClient(t *testing.T, shareURL string, transport policy.Transporter) *share.Client {
	shareClient, err := share.NewClientWithNoCredential(shareURL, fileIDClientOptions(transport))
	require.NoError(t, err)
	return shareClient
}

func TestNewFileClientByFileIDURL(t *testing.T) {
	fileClient := newFakeShareClient(t, fakeShareURL, &fileIDTransport{}).NewFileClientByFileID(fakeFileID)

	require.Equal(t, fakeFileIDURL, fileClient.URL())
	require.Equal(t, fakeFileID, fileClient.FileID())

	urlParts, err := file.ParseURL(fileClient.URL())
	require.NoError(t, err)
	require.Equal(t, "share", urlParts.ShareName)
	require.Equal(t, fakeFileID, urlParts.FileID)
	require.Empty(t, urlParts.DirectoryOrFilePath)
}

func TestNewFileClientByFileIDPreservesShareSnapshot(t *testing.T) {
	shareClient, err := newFakeShareClient(t, fakeShareURL, &fileIDTransport{}).WithSnapshot(fakeShareSnap)
	require.NoError(t, err)

	fileClient := shareClient.NewFileClientByFileID(fakeFileID)
	urlParts, err := file.ParseURL(fileClient.URL())
	require.NoError(t, err)
	require.Equal(t, fakeFileID, urlParts.FileID)
	require.Equal(t, fakeShareSnap, urlParts.ShareSnapshot)
	require.Equal(t, fakeFileID, fileClient.FileID())
}

func TestNewFileClientByFileIDPreservesSAS(t *testing.T) {
	const sasQuery = "se=2015-04-30T02%3A23%3A26Z&sig=Z%2FRHIX5Xcg0Mq2rqI3OlWTjEg2tYkboXr1P9ZUXDtkk%3D&sp=r&sr=s&sv=2015-04-05"
	fileClient := newFakeShareClient(t, fakeShareURL+"?"+sasQuery, &fileIDTransport{}).NewFileClientByFileID(fakeFileID)

	require.Equal(t, fakeFileIDURL+"&"+sasQuery, fileClient.URL())
	require.Equal(t, fakeFileID, fileClient.FileID())
}

func TestFileClientFileIDFromConstructorURL(t *testing.T) {
	// a URL that already contains a file ID produces a client that addresses the file by its file ID
	fileClient, err := file.NewClientWithNoCredential(fakeShareURL+"?FileId="+fakeFileID, nil)
	require.NoError(t, err)
	require.Equal(t, fakeFileID, fileClient.FileID())

	cred, err := file.NewSharedKeyCredential("account", "ZmFrZQ==")
	require.NoError(t, err)
	fileClient, err = file.NewClientWithSharedKeyCredential(fakeFileIDURL, cred, nil)
	require.NoError(t, err)
	require.Equal(t, fakeFileID, fileClient.FileID())
}

func TestFileClientFileIDEmptyForPathBasedClient(t *testing.T) {
	fileClient, err := file.NewClientWithNoCredential(fakeShareURL+"/directory/file", nil)
	require.NoError(t, err)
	require.Empty(t, fileClient.FileID())

	fileClient = newFakeShareClient(t, fakeShareURL, &fileIDTransport{}).NewRootDirectoryClient().NewFileClient("file")
	require.Empty(t, fileClient.FileID())
}

func TestFileClientByFileIDGetProperties(t *testing.T) {
	_require := require.New(t)
	transport := &fileIDTransport{header: http.Header{
		"X-Ms-File-Id":        []string{fakeFileID},
		"X-Ms-File-Parent-Id": []string{fakeParentID},
		"X-Ms-File-Name":      []string{fakeFileName},
		"Content-Length":      []string{"1024"},
		"Content-Type":        []string{"application/octet-stream"},
		"X-Ms-Type":           []string{"File"},
	}}
	fileClient := newFakeShareClient(t, fakeShareURL, transport).NewFileClientByFileID(fakeFileID)

	resp, err := fileClient.GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Equal(fakeFileID, *resp.ID)
	_require.Equal(fakeParentID, *resp.ParentID)
	_require.Equal(fakeFileName, *resp.FileName)
	_require.Equal(int64(1024), *resp.ContentLength)

	_require.Len(transport.requests, 1)
	req := transport.requests[0]
	_require.Equal(http.MethodHead, req.Method)
	_require.Equal("/share", req.URL.Path)
	_require.Equal(fakeFileID, req.URL.Query().Get("fileid"))
	_require.Empty(req.URL.Query().Get("comp"))
}

func TestFileClientByFileIDGetPropertiesSharedKey(t *testing.T) {
	_require := require.New(t)
	transport := &fileIDTransport{}
	cred, err := share.NewSharedKeyCredential("account", "ZmFrZQ==")
	_require.NoError(err)
	shareClient, err := share.NewClientWithSharedKeyCredential(fakeShareURL, cred, fileIDClientOptions(transport))
	_require.NoError(err)

	_, err = shareClient.NewFileClientByFileID(fakeFileID).GetProperties(context.Background(), nil)
	_require.NoError(err)
	_require.Len(transport.requests, 1)
	_require.True(strings.HasPrefix(transport.requests[0].Header.Get("Authorization"), "SharedKey account:"))
	_require.Equal(fakeFileID, transport.requests[0].URL.Query().Get("fileid"))
}

func TestFileClientByFileIDGetPropertiesError(t *testing.T) {
	transport := &fileIDTransport{
		status: http.StatusNotFound,
		header: http.Header{"X-Ms-Error-Code": []string{string(fileerror.ResourceNotFound)}},
	}
	fileClient := newFakeShareClient(t, fakeShareURL, transport).NewFileClientByFileID("11111111111111111111")

	_, err := fileClient.GetProperties(context.Background(), nil)
	require.Error(t, err)
	require.True(t, fileerror.HasCode(err, fileerror.ResourceNotFound))
}

func TestFileClientByFileIDGetFileLinks(t *testing.T) {
	_require := require.New(t)
	transport := &fileIDTransport{
		header: http.Header{
			"Content-Type":             []string{"application/xml"},
			"X-Ms-File-Id":             []string{fakeFileID},
			"X-Ms-File-Parent-Id":      []string{fakeParentID},
			"X-Ms-File-Name":           []string{fakeFileName},
			"X-Ms-Content-Length":      []string{"1024"},
			"X-Ms-Content-Type":        []string{"text/plain"},
			"X-Ms-Cache-Control":       []string{"no-cache"},
			"X-Ms-Content-Disposition": []string{"attachment"},
			"X-Ms-Content-Encoding":    []string{"gzip"},
			"X-Ms-Content-Language":    []string{"en-US"},
			"X-Ms-Content-Md5":         []string{"MDAwMDAwMDAwMDAwMDAwMA=="},
			"X-Ms-Link-Count":          []string{"2"},
			"X-Ms-Type":                []string{"File"},
		},
		body: fakeHardLinkXML,
	}
	fileClient := newFakeShareClient(t, fakeShareURL, transport).NewFileClientByFileID(fakeFileID)

	resp, err := fileClient.GetFileLinks(context.Background(), nil)
	_require.NoError(err)

	// properties of the file are returned with an x-ms- prefix, since the response has a body
	_require.Equal(fakeFileID, *resp.ID)
	_require.Equal(fakeParentID, *resp.ParentID)
	_require.Equal(fakeFileName, *resp.FileName)
	_require.Equal(int64(1024), *resp.FileContentLength)
	_require.Equal("text/plain", *resp.FileContentType)
	_require.Equal("no-cache", *resp.FileCacheControl)
	_require.Equal("attachment", *resp.FileContentDisposition)
	_require.Equal("gzip", *resp.FileContentEncoding)
	_require.Equal("en-US", *resp.FileContentLanguage)
	_require.Equal([]byte("0000000000000000"), resp.FileContentMD5)
	_require.Equal(int64(2), *resp.LinkCount)

	_require.Len(resp.HardLinks, 2)
	_require.Equal(fakeParentID, *resp.HardLinks[0].ParentID)
	_require.Equal(fakeFileName, *resp.HardLinks[0].FileName)
	_require.Equal("0", *resp.HardLinks[1].ParentID)
	_require.Equal("link￾", *resp.HardLinks[1].FileName)

	_require.Len(transport.requests, 1)
	req := transport.requests[0]
	_require.Equal(http.MethodGet, req.Method)
	_require.Equal("/share", req.URL.Path)
	_require.Equal("hardlinks", req.URL.Query().Get("comp"))
	_require.Equal(fakeFileID, req.URL.Query().Get("fileid"))
	_require.Empty(foldedHeaderValues(req.Header, "x-ms-lease-id"))
}

func TestFileClientByFileIDGetFileLinksOptions(t *testing.T) {
	_require := require.New(t)
	transport := &fileIDTransport{body: `<HardLinks />`}
	fileClient := newFakeShareClient(t, fakeShareURL, transport).NewFileClientByFileID(fakeFileID)

	resp, err := fileClient.GetFileLinks(context.Background(), &file.GetFileLinksOptions{
		LeaseAccessConditions: &file.LeaseAccessConditions{LeaseID: to.Ptr("lease-id")},
		ShareSnapshot:         to.Ptr(fakeShareSnap),
	})
	_require.NoError(err)
	_require.Empty(resp.HardLinks)

	_require.Len(transport.requests, 1)
	req := transport.requests[0]
	_require.Equal([]string{"lease-id"}, foldedHeaderValues(req.Header, "x-ms-lease-id"))
	_require.Equal(fakeShareSnap, req.URL.Query().Get("sharesnapshot"))
	_require.Equal(fakeFileID, req.URL.Query().Get("fileid"))
}

func TestFileClientByFileIDGetFileLinksShareSnapshotClient(t *testing.T) {
	_require := require.New(t)
	transport := &fileIDTransport{body: `<HardLinks />`}
	shareClient, err := newFakeShareClient(t, fakeShareURL, transport).WithSnapshot(fakeShareSnap)
	_require.NoError(err)

	_, err = shareClient.NewFileClientByFileID(fakeFileID).GetFileLinks(context.Background(), nil)
	_require.NoError(err)
	_require.Len(transport.requests, 1)
	_require.Equal(fakeShareSnap, transport.requests[0].URL.Query().Get("sharesnapshot"))
	_require.Equal(fakeFileID, transport.requests[0].URL.Query().Get("fileid"))
	_require.Equal("hardlinks", transport.requests[0].URL.Query().Get("comp"))
}

func TestFileClientByFileIDGetFileLinksError(t *testing.T) {
	transport := &fileIDTransport{
		status: http.StatusNotFound,
		header: http.Header{"X-Ms-Error-Code": []string{string(fileerror.ResourceNotFound)}},
	}
	fileClient := newFakeShareClient(t, fakeShareURL, transport).NewFileClientByFileID("11111111111111111111")

	_, err := fileClient.GetFileLinks(context.Background(), nil)
	require.Error(t, err)
	require.True(t, fileerror.HasCode(err, fileerror.ResourceNotFound))
}

func TestFileClientGetFileLinksPathBasedClient(t *testing.T) {
	// GetFileLinks is only supported when the file is addressed by its file ID; the request is never sent
	transport := &fileIDTransport{}
	fileClient := newFakeShareClient(t, fakeShareURL, transport).NewDirectoryClient("directory").NewFileClient("file")

	_, err := fileClient.GetFileLinks(context.Background(), nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDRequired)
	require.Contains(t, err.Error(), "GetFileLinks")
	require.Empty(t, transport.requests)
}

func TestFileClientByFileIDUnsupportedOperations(t *testing.T) {
	ctx := context.Background()
	transport := &fileIDTransport{}
	fileClient := newFakeShareClient(t, fakeShareURL, transport).NewFileClientByFileID(fakeFileID)

	tmpFile, err := os.CreateTemp(t.TempDir(), "fileid")
	require.NoError(t, err)
	defer func() { _ = tmpFile.Close() }()

	body := func() io.ReadSeekCloser { return streaming.NopCloser(bytes.NewReader([]byte("data"))) }

	operations := map[string]func() error{
		"Create": func() error { _, err := fileClient.Create(ctx, 1024, nil); return err },
		"Delete": func() error { _, err := fileClient.Delete(ctx, nil); return err },
		"Rename": func() error { _, err := fileClient.Rename(ctx, "destination", nil); return err },
		"SetHTTPHeaders": func() error {
			_, err := fileClient.SetHTTPHeaders(ctx, nil)
			return err
		},
		"SetMetadata": func() error { _, err := fileClient.SetMetadata(ctx, nil); return err },
		"StartCopyFromURL": func() error {
			_, err := fileClient.StartCopyFromURL(ctx, fakeShareURL+"/source", nil)
			return err
		},
		"AbortCopy":   func() error { _, err := fileClient.AbortCopy(ctx, "copy-id", nil); return err },
		"Resize":      func() error { _, err := fileClient.Resize(ctx, 1024, nil); return err },
		"UploadRange": func() error { _, err := fileClient.UploadRange(ctx, 0, body(), nil); return err },
		"ClearRange": func() error {
			_, err := fileClient.ClearRange(ctx, file.HTTPRange{Offset: 0, Count: 512}, nil)
			return err
		},
		"UploadRangeFromURL": func() error {
			_, err := fileClient.UploadRangeFromURL(ctx, fakeShareURL+"/source", 0, 0, 512, nil)
			return err
		},
		"GetRangeList": func() error { _, err := fileClient.GetRangeList(ctx, nil); return err },
		"NewListRangesPager": func() error {
			_, err := fileClient.NewListRangesPager(nil).NextPage(ctx)
			return err
		},
		"ForceCloseHandles": func() error {
			_, err := fileClient.ForceCloseHandles(ctx, "*", nil)
			return err
		},
		"ListHandles": func() error { _, err := fileClient.ListHandles(ctx, nil); return err },
		"CreateHardLink": func() error {
			_, err := fileClient.CreateHardLink(ctx, "directory/target", nil)
			return err
		},
		"CreateSymbolicLink": func() error {
			_, err := fileClient.CreateSymbolicLink(ctx, "target", nil)
			return err
		},
		"GetSymbolicLink": func() error { _, err := fileClient.GetSymbolicLink(ctx, nil); return err },
		"GetSASURL": func() error {
			_, err := fileClient.GetSASURL(sas.FilePermissions{Read: true}, time.Now().Add(time.Hour), nil)
			return err
		},
		"UploadBuffer":   func() error { return fileClient.UploadBuffer(ctx, []byte("data"), nil) },
		"UploadFile":     func() error { return fileClient.UploadFile(ctx, tmpFile, nil) },
		"UploadStream":   func() error { return fileClient.UploadStream(ctx, strings.NewReader("data"), nil) },
		"DownloadStream": func() error { _, err := fileClient.DownloadStream(ctx, nil); return err },
		"DownloadBuffer": func() error {
			_, err := fileClient.DownloadBuffer(ctx, make([]byte, 4), nil)
			return err
		},
		"DownloadFile": func() error { _, err := fileClient.DownloadFile(ctx, tmpFile, nil); return err },
	}

	for name, operation := range operations {
		t.Run(name, func(t *testing.T) {
			err := operation()
			require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)
			require.Contains(t, err.Error(), "file.Client."+name)
		})
	}
	require.Empty(t, transport.requests, "unsupported operations must fail before a request is sent")
}

func TestFileClientByFileIDGetSASURLWithSharedKey(t *testing.T) {
	// a file ID addressed request has no path to sign, so a file SAS cannot be generated even with a shared key
	cred, err := share.NewSharedKeyCredential("account", "ZmFrZQ==")
	require.NoError(t, err)
	shareClient, err := share.NewClientWithSharedKeyCredential(fakeShareURL, cred, nil)
	require.NoError(t, err)

	_, err = shareClient.NewFileClientByFileID(fakeFileID).GetSASURL(sas.FilePermissions{Read: true}, time.Now().Add(time.Hour), nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)

	// the share is still known, so a share SAS can be generated as before
	_, err = shareClient.GetSASURL(sas.SharePermissions{Read: true}, time.Now().Add(time.Hour), nil)
	require.NoError(t, err)
}

func TestFileClientByFileIDLeaseOperations(t *testing.T) {
	ctx := context.Background()
	transport := &fileIDTransport{}
	fileClient := newFakeShareClient(t, fakeShareURL, transport).NewFileClientByFileID(fakeFileID)

	leaseClient, err := lease.NewFileClient(fileClient, &lease.FileClientOptions{LeaseID: to.Ptr("c820a799-76d7-4ee2-6e15-546f19325c2c")})
	require.NoError(t, err)

	_, err = leaseClient.Acquire(ctx, nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)
	_, err = leaseClient.Release(ctx, nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)
	_, err = leaseClient.Change(ctx, "6b3e65e5-e1bb-4a3f-8b72-13e9bc9cd3bf", nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)
	_, err = leaseClient.Break(ctx, nil)
	require.ErrorIs(t, err, fileerror.ErrFileIDNotSupported)

	require.Empty(t, transport.requests)
}

func TestFileClientPathBasedOperationsUnaffected(t *testing.T) {
	// path based clients keep sending requests for the operations that file ID addressed clients reject
	transport := &fileIDTransport{}
	fileClient := newFakeShareClient(t, fakeShareURL, transport).NewDirectoryClient("directory").NewFileClient("file")

	_, err := fileClient.GetProperties(context.Background(), nil)
	require.NoError(t, err)
	_, err = fileClient.SetMetadata(context.Background(), nil)
	require.NoError(t, err)

	require.Len(t, transport.requests, 2)
	for _, req := range transport.requests {
		require.Equal(t, "/share/directory/file", req.URL.Path)
		require.Empty(t, req.URL.Query().Get("fileid"))
	}
}
