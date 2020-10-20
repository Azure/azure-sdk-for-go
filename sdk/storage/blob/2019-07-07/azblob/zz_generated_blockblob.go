// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package azblob

import (
	"context"
	"encoding/base64"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type blockBlobClient struct {
	*client
}

// Do invokes the Do() method on the pipeline associated with this client.
func (client *blockBlobClient) Do(req *azcore.Request) (*azcore.Response, error) {
	return client.p.Do(req)
}

// CommitBlockList - The Commit Block List operation writes a blob by specifying the list of block IDs that make up the blob. In order to be written as part of a blob, a block must have been successfully written to the server in a prior Put Block operation. You can call Put Block List to update a blob by uploading only those blocks that have changed, then committing the new and existing blocks together. You can do this by specifying whether to commit a block from the committed block list or from the uncommitted block list, or to commit the most recently uploaded version of the block, whichever list it may belong to.
func (client *blockBlobClient) CommitBlockList(ctx context.Context, blocks BlockLookupList, blockBlobCommitBlockListOptions *BlockBlobCommitBlockListOptions, blobHttpHeaders *BlobHttpHeaders, leaseAccessConditions *LeaseAccessConditions, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo, modifiedAccessConditions *ModifiedAccessConditions) (*BlockBlobCommitBlockListResponse, error) {
	req, err := client.CommitBlockListCreateRequest(ctx, blocks, blockBlobCommitBlockListOptions, blobHttpHeaders, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusCreated) {
		return nil, client.CommitBlockListHandleError(resp)
	}
	result, err := client.CommitBlockListHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CommitBlockListCreateRequest creates the CommitBlockList request.
func (client *blockBlobClient) CommitBlockListCreateRequest(ctx context.Context, blocks BlockLookupList, blockBlobCommitBlockListOptions *BlockBlobCommitBlockListOptions, blobHttpHeaders *BlobHttpHeaders, leaseAccessConditions *LeaseAccessConditions, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo, modifiedAccessConditions *ModifiedAccessConditions) (*azcore.Request, error) {
	req, err := azcore.NewRequest(ctx, http.MethodPut, client.u)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("comp", "blocklist")
	if blockBlobCommitBlockListOptions != nil && blockBlobCommitBlockListOptions.Timeout != nil {
		query.Set("timeout", strconv.FormatInt(int64(*blockBlobCommitBlockListOptions.Timeout), 10))
	}
	req.URL.RawQuery = query.Encode()
	if blobHttpHeaders != nil && blobHttpHeaders.BlobCacheControl != nil {
		req.Header.Set("x-ms-blob-cache-control", *blobHttpHeaders.BlobCacheControl)
	}
	if blobHttpHeaders != nil && blobHttpHeaders.BlobContentType != nil {
		req.Header.Set("x-ms-blob-content-type", *blobHttpHeaders.BlobContentType)
	}
	if blobHttpHeaders != nil && blobHttpHeaders.BlobContentEncoding != nil {
		req.Header.Set("x-ms-blob-content-encoding", *blobHttpHeaders.BlobContentEncoding)
	}
	if blobHttpHeaders != nil && blobHttpHeaders.BlobContentLanguage != nil {
		req.Header.Set("x-ms-blob-content-language", *blobHttpHeaders.BlobContentLanguage)
	}
	if blobHttpHeaders != nil && blobHttpHeaders.BlobContentMd5 != nil {
		req.Header.Set("x-ms-blob-content-md5", base64.StdEncoding.EncodeToString(*blobHttpHeaders.BlobContentMd5))
	}
	if blockBlobCommitBlockListOptions != nil && blockBlobCommitBlockListOptions.TransactionalContentMd5 != nil {
		req.Header.Set("Content-MD5", base64.StdEncoding.EncodeToString(*blockBlobCommitBlockListOptions.TransactionalContentMd5))
	}
	if blockBlobCommitBlockListOptions != nil && blockBlobCommitBlockListOptions.TransactionalContentCrc64 != nil {
		req.Header.Set("x-ms-content-crc64", base64.StdEncoding.EncodeToString(*blockBlobCommitBlockListOptions.TransactionalContentCrc64))
	}
	if blockBlobCommitBlockListOptions != nil && blockBlobCommitBlockListOptions.Metadata != nil {
		for k, v := range *blockBlobCommitBlockListOptions.Metadata {
			req.Header.Set("x-ms-meta-"+k, v)
		}
	}
	if leaseAccessConditions != nil && leaseAccessConditions.LeaseId != nil {
		req.Header.Set("x-ms-lease-id", *leaseAccessConditions.LeaseId)
	}
	if blobHttpHeaders != nil && blobHttpHeaders.BlobContentDisposition != nil {
		req.Header.Set("x-ms-blob-content-disposition", *blobHttpHeaders.BlobContentDisposition)
	}
	if cpkInfo != nil && cpkInfo.EncryptionKey != nil {
		req.Header.Set("x-ms-encryption-key", *cpkInfo.EncryptionKey)
	}
	if cpkInfo != nil && cpkInfo.EncryptionKeySha256 != nil {
		req.Header.Set("x-ms-encryption-key-sha256", *cpkInfo.EncryptionKeySha256)
	}
	if cpkInfo != nil && cpkInfo.EncryptionAlgorithm != nil {
		req.Header.Set("x-ms-encryption-algorithm", "AES256")
	}
	if cpkScopeInfo != nil && cpkScopeInfo.EncryptionScope != nil {
		req.Header.Set("x-ms-encryption-scope", *cpkScopeInfo.EncryptionScope)
	}
	if blockBlobCommitBlockListOptions != nil && blockBlobCommitBlockListOptions.Tier != nil {
		req.Header.Set("x-ms-access-tier", string(*blockBlobCommitBlockListOptions.Tier))
	}
	if modifiedAccessConditions != nil && modifiedAccessConditions.IfModifiedSince != nil {
		req.Header.Set("If-Modified-Since", modifiedAccessConditions.IfModifiedSince.Format(time.RFC1123))
	}
	if modifiedAccessConditions != nil && modifiedAccessConditions.IfUnmodifiedSince != nil {
		req.Header.Set("If-Unmodified-Since", modifiedAccessConditions.IfUnmodifiedSince.Format(time.RFC1123))
	}
	if modifiedAccessConditions != nil && modifiedAccessConditions.IfMatch != nil {
		req.Header.Set("If-Match", *modifiedAccessConditions.IfMatch)
	}
	if modifiedAccessConditions != nil && modifiedAccessConditions.IfNoneMatch != nil {
		req.Header.Set("If-None-Match", *modifiedAccessConditions.IfNoneMatch)
	}
	req.Header.Set("x-ms-version", "2019-07-07")
	if blockBlobCommitBlockListOptions != nil && blockBlobCommitBlockListOptions.RequestId != nil {
		req.Header.Set("x-ms-client-request-id", *blockBlobCommitBlockListOptions.RequestId)
	}
	req.Header.Set("Accept", "application/xml")
	return req, req.MarshalAsXML(blocks)
}

// CommitBlockListHandleResponse handles the CommitBlockList response.
func (client *blockBlobClient) CommitBlockListHandleResponse(resp *azcore.Response) (*BlockBlobCommitBlockListResponse, error) {
	result := BlockBlobCommitBlockListResponse{RawResponse: resp.Response}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if val := resp.Header.Get("Last-Modified"); val != "" {
		lastModified, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return nil, err
		}
		result.LastModified = &lastModified
	}
	if val := resp.Header.Get("Content-MD5"); val != "" {
		contentMd5, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return nil, err
		}
		result.ContentMD5 = &contentMd5
	}
	if val := resp.Header.Get("x-ms-content-crc64"); val != "" {
		contentCrc64, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return nil, err
		}
		result.ContentCRC64 = &contentCrc64
	}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return nil, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("x-ms-request-server-encrypted"); val != "" {
		isServerEncrypted, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		result.IsServerEncrypted = &isServerEncrypted
	}
	if val := resp.Header.Get("x-ms-encryption-key-sha256"); val != "" {
		result.EncryptionKeySHA256 = &val
	}
	if val := resp.Header.Get("x-ms-encryption-scope"); val != "" {
		result.EncryptionScope = &val
	}
	return &result, nil
}

// CommitBlockListHandleError handles the CommitBlockList error response.
func (client *blockBlobClient) CommitBlockListHandleError(resp *azcore.Response) error {
	var err StorageError
	if err := resp.UnmarshalAsXML(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// GetBlockList - The Get Block List operation retrieves the list of blocks that have been uploaded as part of a block blob
func (client *blockBlobClient) GetBlockList(ctx context.Context, listType BlockListType, blockBlobGetBlockListOptions *BlockBlobGetBlockListOptions, leaseAccessConditions *LeaseAccessConditions) (*BlockListResponse, error) {
	req, err := client.GetBlockListCreateRequest(ctx, listType, blockBlobGetBlockListOptions, leaseAccessConditions)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.GetBlockListHandleError(resp)
	}
	result, err := client.GetBlockListHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetBlockListCreateRequest creates the GetBlockList request.
func (client *blockBlobClient) GetBlockListCreateRequest(ctx context.Context, listType BlockListType, blockBlobGetBlockListOptions *BlockBlobGetBlockListOptions, leaseAccessConditions *LeaseAccessConditions) (*azcore.Request, error) {
	req, err := azcore.NewRequest(ctx, http.MethodGet, client.u)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("comp", "blocklist")
	if blockBlobGetBlockListOptions != nil && blockBlobGetBlockListOptions.Snapshot != nil {
		query.Set("snapshot", *blockBlobGetBlockListOptions.Snapshot)
	}
	query.Set("blocklisttype", string(listType))
	if blockBlobGetBlockListOptions != nil && blockBlobGetBlockListOptions.Timeout != nil {
		query.Set("timeout", strconv.FormatInt(int64(*blockBlobGetBlockListOptions.Timeout), 10))
	}
	req.URL.RawQuery = query.Encode()
	if leaseAccessConditions != nil && leaseAccessConditions.LeaseId != nil {
		req.Header.Set("x-ms-lease-id", *leaseAccessConditions.LeaseId)
	}
	req.Header.Set("x-ms-version", "2019-07-07")
	if blockBlobGetBlockListOptions != nil && blockBlobGetBlockListOptions.RequestId != nil {
		req.Header.Set("x-ms-client-request-id", *blockBlobGetBlockListOptions.RequestId)
	}
	req.Header.Set("Accept", "application/xml")
	return req, nil
}

// GetBlockListHandleResponse handles the GetBlockList response.
func (client *blockBlobClient) GetBlockListHandleResponse(resp *azcore.Response) (*BlockListResponse, error) {
	result := BlockListResponse{RawResponse: resp.Response}
	if val := resp.Header.Get("Last-Modified"); val != "" {
		lastModified, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return nil, err
		}
		result.LastModified = &lastModified
	}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if val := resp.Header.Get("Content-Type"); val != "" {
		result.ContentType = &val
	}
	if val := resp.Header.Get("x-ms-blob-content-length"); val != "" {
		blobContentLength, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		result.BlobContentLength = &blobContentLength
	}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return nil, err
		}
		result.Date = &date
	}
	return &result, resp.UnmarshalAsXML(&result.BlockList)
}

// GetBlockListHandleError handles the GetBlockList error response.
func (client *blockBlobClient) GetBlockListHandleError(resp *azcore.Response) error {
	var err StorageError
	if err := resp.UnmarshalAsXML(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// StageBlock - The Stage Block operation creates a new block to be committed as part of a blob
func (client *blockBlobClient) StageBlock(ctx context.Context, blockId string, contentLength int64, body azcore.ReadSeekCloser, blockBlobStageBlockOptions *BlockBlobStageBlockOptions, leaseAccessConditions *LeaseAccessConditions, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo) (*BlockBlobStageBlockResponse, error) {
	req, err := client.StageBlockCreateRequest(ctx, blockId, contentLength, body, blockBlobStageBlockOptions, leaseAccessConditions, cpkInfo, cpkScopeInfo)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusCreated) {
		return nil, client.StageBlockHandleError(resp)
	}
	result, err := client.StageBlockHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// StageBlockCreateRequest creates the StageBlock request.
func (client *blockBlobClient) StageBlockCreateRequest(ctx context.Context, blockId string, contentLength int64, body azcore.ReadSeekCloser, blockBlobStageBlockOptions *BlockBlobStageBlockOptions, leaseAccessConditions *LeaseAccessConditions, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo) (*azcore.Request, error) {
	req, err := azcore.NewRequest(ctx, http.MethodPut, client.u)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("comp", "block")
	query.Set("blockid", blockId)
	if blockBlobStageBlockOptions != nil && blockBlobStageBlockOptions.Timeout != nil {
		query.Set("timeout", strconv.FormatInt(int64(*blockBlobStageBlockOptions.Timeout), 10))
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Content-Length", strconv.FormatInt(contentLength, 10))
	if blockBlobStageBlockOptions != nil && blockBlobStageBlockOptions.TransactionalContentMd5 != nil {
		req.Header.Set("Content-MD5", base64.StdEncoding.EncodeToString(*blockBlobStageBlockOptions.TransactionalContentMd5))
	}
	if blockBlobStageBlockOptions != nil && blockBlobStageBlockOptions.TransactionalContentCrc64 != nil {
		req.Header.Set("x-ms-content-crc64", base64.StdEncoding.EncodeToString(*blockBlobStageBlockOptions.TransactionalContentCrc64))
	}
	if leaseAccessConditions != nil && leaseAccessConditions.LeaseId != nil {
		req.Header.Set("x-ms-lease-id", *leaseAccessConditions.LeaseId)
	}
	if cpkInfo != nil && cpkInfo.EncryptionKey != nil {
		req.Header.Set("x-ms-encryption-key", *cpkInfo.EncryptionKey)
	}
	if cpkInfo != nil && cpkInfo.EncryptionKeySha256 != nil {
		req.Header.Set("x-ms-encryption-key-sha256", *cpkInfo.EncryptionKeySha256)
	}
	if cpkInfo != nil && cpkInfo.EncryptionAlgorithm != nil {
		req.Header.Set("x-ms-encryption-algorithm", "AES256")
	}
	if cpkScopeInfo != nil && cpkScopeInfo.EncryptionScope != nil {
		req.Header.Set("x-ms-encryption-scope", *cpkScopeInfo.EncryptionScope)
	}
	req.Header.Set("x-ms-version", "2019-07-07")
	if blockBlobStageBlockOptions != nil && blockBlobStageBlockOptions.RequestId != nil {
		req.Header.Set("x-ms-client-request-id", *blockBlobStageBlockOptions.RequestId)
	}
	req.Header.Set("Accept", "application/xml")
	return req, req.SetBody(body, "application/octet-stream")
}

// StageBlockHandleResponse handles the StageBlock response.
func (client *blockBlobClient) StageBlockHandleResponse(resp *azcore.Response) (*BlockBlobStageBlockResponse, error) {
	result := BlockBlobStageBlockResponse{RawResponse: resp.Response}
	if val := resp.Header.Get("Content-MD5"); val != "" {
		contentMd5, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return nil, err
		}
		result.ContentMD5 = &contentMd5
	}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return nil, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("x-ms-content-crc64"); val != "" {
		contentCrc64, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return nil, err
		}
		result.ContentCRC64 = &contentCrc64
	}
	if val := resp.Header.Get("x-ms-request-server-encrypted"); val != "" {
		isServerEncrypted, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		result.IsServerEncrypted = &isServerEncrypted
	}
	if val := resp.Header.Get("x-ms-encryption-key-sha256"); val != "" {
		result.EncryptionKeySHA256 = &val
	}
	if val := resp.Header.Get("x-ms-encryption-scope"); val != "" {
		result.EncryptionScope = &val
	}
	return &result, nil
}

// StageBlockHandleError handles the StageBlock error response.
func (client *blockBlobClient) StageBlockHandleError(resp *azcore.Response) error {
	var err StorageError
	if err := resp.UnmarshalAsXML(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// StageBlockFromURL - The Stage Block operation creates a new block to be committed as part of a blob where the contents are read from a URL.
func (client *blockBlobClient) StageBlockFromURL(ctx context.Context, blockId string, contentLength int64, sourceUrl url.URL, blockBlobStageBlockFromUrlOptions *BlockBlobStageBlockFromURLOptions, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo, leaseAccessConditions *LeaseAccessConditions, sourceModifiedAccessConditions *SourceModifiedAccessConditions) (*BlockBlobStageBlockFromURLResponse, error) {
	req, err := client.StageBlockFromURLCreateRequest(ctx, blockId, contentLength, sourceUrl, blockBlobStageBlockFromUrlOptions, cpkInfo, cpkScopeInfo, leaseAccessConditions, sourceModifiedAccessConditions)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusCreated) {
		return nil, client.StageBlockFromURLHandleError(resp)
	}
	result, err := client.StageBlockFromURLHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// StageBlockFromURLCreateRequest creates the StageBlockFromURL request.
func (client *blockBlobClient) StageBlockFromURLCreateRequest(ctx context.Context, blockId string, contentLength int64, sourceUrl url.URL, blockBlobStageBlockFromUrlOptions *BlockBlobStageBlockFromURLOptions, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo, leaseAccessConditions *LeaseAccessConditions, sourceModifiedAccessConditions *SourceModifiedAccessConditions) (*azcore.Request, error) {
	req, err := azcore.NewRequest(ctx, http.MethodPut, client.u)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("comp", "block")
	query.Set("blockid", blockId)
	if blockBlobStageBlockFromUrlOptions != nil && blockBlobStageBlockFromUrlOptions.Timeout != nil {
		query.Set("timeout", strconv.FormatInt(int64(*blockBlobStageBlockFromUrlOptions.Timeout), 10))
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Content-Length", strconv.FormatInt(contentLength, 10))
	req.Header.Set("x-ms-copy-source", sourceUrl.String())
	if blockBlobStageBlockFromUrlOptions != nil && blockBlobStageBlockFromUrlOptions.SourceRange != nil {
		req.Header.Set("x-ms-source-range", *blockBlobStageBlockFromUrlOptions.SourceRange)
	}
	if blockBlobStageBlockFromUrlOptions != nil && blockBlobStageBlockFromUrlOptions.SourceContentMd5 != nil {
		req.Header.Set("x-ms-source-content-md5", base64.StdEncoding.EncodeToString(*blockBlobStageBlockFromUrlOptions.SourceContentMd5))
	}
	if blockBlobStageBlockFromUrlOptions != nil && blockBlobStageBlockFromUrlOptions.SourceContentcrc64 != nil {
		req.Header.Set("x-ms-source-content-crc64", base64.StdEncoding.EncodeToString(*blockBlobStageBlockFromUrlOptions.SourceContentcrc64))
	}
	if cpkInfo != nil && cpkInfo.EncryptionKey != nil {
		req.Header.Set("x-ms-encryption-key", *cpkInfo.EncryptionKey)
	}
	if cpkInfo != nil && cpkInfo.EncryptionKeySha256 != nil {
		req.Header.Set("x-ms-encryption-key-sha256", *cpkInfo.EncryptionKeySha256)
	}
	if cpkInfo != nil && cpkInfo.EncryptionAlgorithm != nil {
		req.Header.Set("x-ms-encryption-algorithm", "AES256")
	}
	if cpkScopeInfo != nil && cpkScopeInfo.EncryptionScope != nil {
		req.Header.Set("x-ms-encryption-scope", *cpkScopeInfo.EncryptionScope)
	}
	if leaseAccessConditions != nil && leaseAccessConditions.LeaseId != nil {
		req.Header.Set("x-ms-lease-id", *leaseAccessConditions.LeaseId)
	}
	if sourceModifiedAccessConditions != nil && sourceModifiedAccessConditions.SourceIfModifiedSince != nil {
		req.Header.Set("x-ms-source-if-modified-since", sourceModifiedAccessConditions.SourceIfModifiedSince.Format(time.RFC1123))
	}
	if sourceModifiedAccessConditions != nil && sourceModifiedAccessConditions.SourceIfUnmodifiedSince != nil {
		req.Header.Set("x-ms-source-if-unmodified-since", sourceModifiedAccessConditions.SourceIfUnmodifiedSince.Format(time.RFC1123))
	}
	if sourceModifiedAccessConditions != nil && sourceModifiedAccessConditions.SourceIfMatch != nil {
		req.Header.Set("x-ms-source-if-match", *sourceModifiedAccessConditions.SourceIfMatch)
	}
	if sourceModifiedAccessConditions != nil && sourceModifiedAccessConditions.SourceIfNoneMatch != nil {
		req.Header.Set("x-ms-source-if-none-match", *sourceModifiedAccessConditions.SourceIfNoneMatch)
	}
	req.Header.Set("x-ms-version", "2019-07-07")
	if blockBlobStageBlockFromUrlOptions != nil && blockBlobStageBlockFromUrlOptions.RequestId != nil {
		req.Header.Set("x-ms-client-request-id", *blockBlobStageBlockFromUrlOptions.RequestId)
	}
	req.Header.Set("Accept", "application/xml")
	return req, nil
}

// StageBlockFromURLHandleResponse handles the StageBlockFromURL response.
func (client *blockBlobClient) StageBlockFromURLHandleResponse(resp *azcore.Response) (*BlockBlobStageBlockFromURLResponse, error) {
	result := BlockBlobStageBlockFromURLResponse{RawResponse: resp.Response}
	if val := resp.Header.Get("Content-MD5"); val != "" {
		contentMd5, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return nil, err
		}
		result.ContentMD5 = &contentMd5
	}
	if val := resp.Header.Get("x-ms-content-crc64"); val != "" {
		contentCrc64, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return nil, err
		}
		result.ContentCRC64 = &contentCrc64
	}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return nil, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("x-ms-request-server-encrypted"); val != "" {
		isServerEncrypted, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		result.IsServerEncrypted = &isServerEncrypted
	}
	if val := resp.Header.Get("x-ms-encryption-key-sha256"); val != "" {
		result.EncryptionKeySHA256 = &val
	}
	if val := resp.Header.Get("x-ms-encryption-scope"); val != "" {
		result.EncryptionScope = &val
	}
	return &result, nil
}

// StageBlockFromURLHandleError handles the StageBlockFromURL error response.
func (client *blockBlobClient) StageBlockFromURLHandleError(resp *azcore.Response) error {
	var err StorageError
	if err := resp.UnmarshalAsXML(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}

// Upload - The Upload Block Blob operation updates the content of an existing block blob. Updating an existing block blob overwrites any existing metadata on the blob. Partial updates are not supported with Put Blob; the content of the existing blob is overwritten with the content of the new blob. To perform a partial update of the content of a block blob, use the Put Block List operation.
func (client *blockBlobClient) Upload(ctx context.Context, contentLength int64, body azcore.ReadSeekCloser, blockBlobUploadOptions *BlockBlobUploadOptions, blobHttpHeaders *BlobHttpHeaders, leaseAccessConditions *LeaseAccessConditions, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo, modifiedAccessConditions *ModifiedAccessConditions) (*BlockBlobUploadResponse, error) {
	req, err := client.UploadCreateRequest(ctx, contentLength, body, blockBlobUploadOptions, blobHttpHeaders, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusCreated) {
		return nil, client.UploadHandleError(resp)
	}
	result, err := client.UploadHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UploadCreateRequest creates the Upload request.
func (client *blockBlobClient) UploadCreateRequest(ctx context.Context, contentLength int64, body azcore.ReadSeekCloser, blockBlobUploadOptions *BlockBlobUploadOptions, blobHttpHeaders *BlobHttpHeaders, leaseAccessConditions *LeaseAccessConditions, cpkInfo *CpkInfo, cpkScopeInfo *CpkScopeInfo, modifiedAccessConditions *ModifiedAccessConditions) (*azcore.Request, error) {
	req, err := azcore.NewRequest(ctx, http.MethodPut, client.u)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if blockBlobUploadOptions != nil && blockBlobUploadOptions.Timeout != nil {
		query.Set("timeout", strconv.FormatInt(int64(*blockBlobUploadOptions.Timeout), 10))
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("x-ms-blob-type", "BlockBlob")
	if blockBlobUploadOptions != nil && blockBlobUploadOptions.TransactionalContentMd5 != nil {
		req.Header.Set("Content-MD5", base64.StdEncoding.EncodeToString(*blockBlobUploadOptions.TransactionalContentMd5))
	}
	req.Header.Set("Content-Length", strconv.FormatInt(contentLength, 10))
	if blobHttpHeaders != nil && blobHttpHeaders.BlobContentType != nil {
		req.Header.Set("x-ms-blob-content-type", *blobHttpHeaders.BlobContentType)
	}
	if blobHttpHeaders != nil && blobHttpHeaders.BlobContentEncoding != nil {
		req.Header.Set("x-ms-blob-content-encoding", *blobHttpHeaders.BlobContentEncoding)
	}
	if blobHttpHeaders != nil && blobHttpHeaders.BlobContentLanguage != nil {
		req.Header.Set("x-ms-blob-content-language", *blobHttpHeaders.BlobContentLanguage)
	}
	if blobHttpHeaders != nil && blobHttpHeaders.BlobContentMd5 != nil {
		req.Header.Set("x-ms-blob-content-md5", base64.StdEncoding.EncodeToString(*blobHttpHeaders.BlobContentMd5))
	}
	if blobHttpHeaders != nil && blobHttpHeaders.BlobCacheControl != nil {
		req.Header.Set("x-ms-blob-cache-control", *blobHttpHeaders.BlobCacheControl)
	}
	if blockBlobUploadOptions != nil && blockBlobUploadOptions.Metadata != nil {
		for k, v := range *blockBlobUploadOptions.Metadata {
			req.Header.Set("x-ms-meta-"+k, v)
		}
	}
	if leaseAccessConditions != nil && leaseAccessConditions.LeaseId != nil {
		req.Header.Set("x-ms-lease-id", *leaseAccessConditions.LeaseId)
	}
	if blobHttpHeaders != nil && blobHttpHeaders.BlobContentDisposition != nil {
		req.Header.Set("x-ms-blob-content-disposition", *blobHttpHeaders.BlobContentDisposition)
	}
	if cpkInfo != nil && cpkInfo.EncryptionKey != nil {
		req.Header.Set("x-ms-encryption-key", *cpkInfo.EncryptionKey)
	}
	if cpkInfo != nil && cpkInfo.EncryptionKeySha256 != nil {
		req.Header.Set("x-ms-encryption-key-sha256", *cpkInfo.EncryptionKeySha256)
	}
	if cpkInfo != nil && cpkInfo.EncryptionAlgorithm != nil {
		req.Header.Set("x-ms-encryption-algorithm", "AES256")
	}
	if cpkScopeInfo != nil && cpkScopeInfo.EncryptionScope != nil {
		req.Header.Set("x-ms-encryption-scope", *cpkScopeInfo.EncryptionScope)
	}
	if blockBlobUploadOptions != nil && blockBlobUploadOptions.Tier != nil {
		req.Header.Set("x-ms-access-tier", string(*blockBlobUploadOptions.Tier))
	}
	if modifiedAccessConditions != nil && modifiedAccessConditions.IfModifiedSince != nil {
		req.Header.Set("If-Modified-Since", modifiedAccessConditions.IfModifiedSince.Format(time.RFC1123))
	}
	if modifiedAccessConditions != nil && modifiedAccessConditions.IfUnmodifiedSince != nil {
		req.Header.Set("If-Unmodified-Since", modifiedAccessConditions.IfUnmodifiedSince.Format(time.RFC1123))
	}
	if modifiedAccessConditions != nil && modifiedAccessConditions.IfMatch != nil {
		req.Header.Set("If-Match", *modifiedAccessConditions.IfMatch)
	}
	if modifiedAccessConditions != nil && modifiedAccessConditions.IfNoneMatch != nil {
		req.Header.Set("If-None-Match", *modifiedAccessConditions.IfNoneMatch)
	}
	req.Header.Set("x-ms-version", "2019-07-07")
	if blockBlobUploadOptions != nil && blockBlobUploadOptions.RequestId != nil {
		req.Header.Set("x-ms-client-request-id", *blockBlobUploadOptions.RequestId)
	}
	req.Header.Set("Accept", "application/xml")
	return req, req.SetBody(body, "application/octet-stream")
}

// UploadHandleResponse handles the Upload response.
func (client *blockBlobClient) UploadHandleResponse(resp *azcore.Response) (*BlockBlobUploadResponse, error) {
	result := BlockBlobUploadResponse{RawResponse: resp.Response}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if val := resp.Header.Get("Last-Modified"); val != "" {
		lastModified, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return nil, err
		}
		result.LastModified = &lastModified
	}
	if val := resp.Header.Get("Content-MD5"); val != "" {
		contentMd5, err := base64.StdEncoding.DecodeString(val)
		if err != nil {
			return nil, err
		}
		result.ContentMD5 = &contentMd5
	}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return nil, err
		}
		result.Date = &date
	}
	if val := resp.Header.Get("x-ms-request-server-encrypted"); val != "" {
		isServerEncrypted, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}
		result.IsServerEncrypted = &isServerEncrypted
	}
	if val := resp.Header.Get("x-ms-encryption-key-sha256"); val != "" {
		result.EncryptionKeySHA256 = &val
	}
	if val := resp.Header.Get("x-ms-encryption-scope"); val != "" {
		result.EncryptionScope = &val
	}
	return &result, nil
}

// UploadHandleError handles the Upload error response.
func (client *blockBlobClient) UploadHandleError(resp *azcore.Response) error {
	var err StorageError
	if err := resp.UnmarshalAsXML(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}
