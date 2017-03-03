package storage

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	blobCopyStatusPending = "pending"
	blobCopyStatusSuccess = "success"
	blobCopyStatusAborted = "aborted"
	blobCopyStatusFailed  = "failed"
)

// Copy starts a blob copy operation and waits for the operation to
// complete. sourceBlob parameter must be a canonical URL to the blob (can be
// obtained using GetBlobURL method.) There is no SLA on blob copy and therefore
// this helper method works faster on smaller files.
//
// See https://msdn.microsoft.com/en-us/library/azure/dd894037.aspx
func (b *Blob) Copy(sourceBlob string) error {
	copyID, err := b.StartCopy(sourceBlob)
	if err != nil {
		return err
	}

	return b.WaitForCopy(copyID)
}

// StartCopy starts a blob copy operation.
// sourceBlob parameter must be a canonical URL to the blob (can be
// obtained using GetBlobURL method.)
//
// See https://msdn.microsoft.com/en-us/library/azure/dd894037.aspx
func (b *Blob) StartCopy(sourceBlob string) (string, error) {
	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), nil)

	headers := b.Container.bsc.client.getStandardHeaders()
	headers["x-ms-copy-source"] = sourceBlob

	resp, err := b.Container.bsc.client.exec(http.MethodPut, uri, headers, nil, b.Container.bsc.auth)
	if err != nil {
		return "", err
	}
	defer readAndCloseBody(resp.body)

	if err := checkRespCode(resp.statusCode, []int{http.StatusAccepted, http.StatusCreated}); err != nil {
		return "", err
	}

	copyID := resp.headers.Get("x-ms-copy-id")
	if copyID == "" {
		return "", errors.New("Got empty copy id header")
	}
	return copyID, nil
}

// AbortCopy aborts a BlobCopy which has already been triggered by the StartBlobCopy function.
// copyID is generated from StartBlobCopy function.
// currentLeaseID is required IF the destination blob has an active lease on it.
// As defined in https://msdn.microsoft.com/en-us/library/azure/jj159098.aspx
func (b *Blob) AbortCopy(copyID, currentLeaseID string, timeout int) error {
	params := url.Values{
		"comp":   {"copy"},
		"copyid": {copyID},
	}
	if timeout > 0 {
		params.Add("timeout", strconv.Itoa(timeout))
	}

	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), params)
	headers := b.Container.bsc.client.getStandardHeaders()
	headers["x-ms-copy-action"] = "abort"

	if currentLeaseID != "" {
		headers[headerLeaseID] = currentLeaseID
	}

	resp, err := b.Container.bsc.client.exec(http.MethodPut, uri, headers, nil, b.Container.bsc.auth)
	if err != nil {
		return err
	}
	readAndCloseBody(resp.body)
	return checkRespCode(resp.statusCode, []int{http.StatusNoContent})
}

// WaitForCopy loops until a BlobCopy operation is completed (or fails with error)
func (b *Blob) WaitForCopy(copyID string) error {
	for {
		err := b.GetProperties()
		if err != nil {
			return err
		}

		if b.Properties.CopyID != copyID {
			return errBlobCopyIDMismatch
		}

		switch b.Properties.CopyStatus {
		case blobCopyStatusSuccess:
			return nil
		case blobCopyStatusPending:
			continue
		case blobCopyStatusAborted:
			return errBlobCopyAborted
		case blobCopyStatusFailed:
			return fmt.Errorf("storage: blob copy failed. Id=%s Description=%s", b.Properties.CopyID, b.Properties.CopyStatusDescription)
		default:
			return fmt.Errorf("storage: unhandled blob copy status: '%s'", b.Properties.CopyStatus)
		}
	}
}
