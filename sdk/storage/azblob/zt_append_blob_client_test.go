// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
	"io/ioutil"
	"strings"
	"time"
)

func (s *aztestsSuite) TestAppendBlock(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	abClient := containerClient.NewAppendBlobURL(generateBlobName())

	resp, err := abClient.Create(context.Background(), nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)

	appendResp, err := abClient.AppendBlock(context.Background(), getReaderToRandomBytes(1024), nil)
	c.Assert(err, chk.IsNil)
	c.Assert(appendResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*appendResp.BlobAppendOffset, chk.Equals, "0")
	c.Assert(*appendResp.BlobCommittedBlockCount, chk.Equals, int32(1))
	c.Assert(appendResp.ETag, chk.NotNil)
	c.Assert(appendResp.LastModified, chk.NotNil)
	c.Assert((*appendResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(appendResp.ContentMD5, chk.IsNil)
	c.Assert(appendResp.RequestID, chk.NotNil)
	c.Assert(appendResp.Version, chk.NotNil)
	c.Assert(appendResp.Date, chk.NotNil)
	c.Assert((*appendResp.Date).IsZero(), chk.Equals, false)

	appendResp, err = abClient.AppendBlock(context.Background(), getReaderToRandomBytes(1024), nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*appendResp.BlobAppendOffset, chk.Equals, "1024")
	c.Assert(*appendResp.BlobCommittedBlockCount, chk.Equals, int32(2))
}

func (s *aztestsSuite) TestAppendBlockWithMD5(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	// set up abClient to test
	abClient := containerClient.NewAppendBlobURL(generateBlobName())
	resp, err := abClient.Create(context.Background(), nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 201)

	// test append block with valid MD5 value
	readerToBody, body := getRandomDataAndReader(1024)
	md5Value := md5.Sum(body)
	contentMD5 := md5Value[:]
	appendBlockOptions := AppendBlockOptions{
		TransactionalContentMD5: &contentMD5,
	}
	appendResp, err := abClient.AppendBlock(context.Background(), readerToBody, &appendBlockOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(appendResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*appendResp.BlobAppendOffset, chk.Equals, "0")
	c.Assert(*appendResp.BlobCommittedBlockCount, chk.Equals, int32(1))
	c.Assert(appendResp.ETag, chk.NotNil)
	c.Assert(appendResp.LastModified, chk.NotNil)
	c.Assert((*appendResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(appendResp.ContentMD5, chk.NotNil)
	c.Assert(*appendResp.ContentMD5, chk.DeepEquals, contentMD5)
	c.Assert(appendResp.RequestID, chk.NotNil)
	c.Assert(appendResp.Version, chk.NotNil)
	c.Assert(appendResp.Date, chk.NotNil)
	c.Assert((*appendResp.Date).IsZero(), chk.Equals, false)

	// test append block with bad MD5 value
	readerToBody, body = getRandomDataAndReader(1024)
	_, badMD5 := getRandomDataAndReader(16)
	appendBlockOptions = AppendBlockOptions{
		TransactionalContentMD5: &badMD5,
	}
	appendResp, err = abClient.AppendBlock(context.Background(), readerToBody, &appendBlockOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeMD5Mismatch)
}

func (s *aztestsSuite) TestAppendBlockFromURL(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	//ctx := context.Background()
	contentSize := 8 * 1024 // 8KB
	r, sourceData := getRandomDataAndReader(contentSize)
	contentMD5 := md5.Sum(sourceData)
	srcBlob := containerClient.NewAppendBlobURL(generateName("appendsrc"))
	destBlob := containerClient.NewAppendBlobURL(generateName("appenddest"))

	// Prepare source abClient for copy.
	cResp1, err := srcBlob.Create(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp1.RawResponse.StatusCode, chk.Equals, 201)

	appendResp, err := srcBlob.AppendBlock(context.Background(), r, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(appendResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*appendResp.BlobAppendOffset, chk.Equals, "0")
	c.Assert(*appendResp.BlobCommittedBlockCount, chk.Equals, int32(1))
	c.Assert(appendResp.ETag, chk.NotNil)
	c.Assert(appendResp.LastModified, chk.NotNil)
	c.Assert((*appendResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(appendResp.ContentMD5, chk.IsNil)
	c.Assert(appendResp.RequestID, chk.NotNil)
	c.Assert(appendResp.Version, chk.NotNil)
	c.Assert(appendResp.Date, chk.NotNil)
	c.Assert((*appendResp.Date).IsZero(), chk.Equals, false)

	// Get source abClient URL with SAS for AppendBlockFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Append block from URL.
	cResp2, err := destBlob.Create(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp2.RawResponse.StatusCode, chk.Equals, 201)

	//ctx context.Context, source url.URL, contentLength int64, options *AppendBlockURLOptions)
	offset := int64(0)
	count := int64(CountToEnd)
	appendBlockURLOptions := AppendBlockURLOptions{
		Offset: &offset,
		Count:  &count,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(appendFromURLResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*appendFromURLResp.BlobAppendOffset, chk.Equals, "0")
	c.Assert(*appendFromURLResp.BlobCommittedBlockCount, chk.Equals, int32(1))
	c.Assert(appendFromURLResp.ETag, chk.NotNil)
	c.Assert(appendFromURLResp.LastModified, chk.NotNil)
	c.Assert((*appendFromURLResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(appendFromURLResp.ContentMD5, chk.NotNil)
	c.Assert(*appendFromURLResp.ContentMD5, chk.DeepEquals, contentMD5[:])
	c.Assert(appendFromURLResp.RequestID, chk.NotNil)
	c.Assert(appendFromURLResp.Version, chk.NotNil)
	c.Assert(appendFromURLResp.Date, chk.NotNil)
	c.Assert((*appendFromURLResp.Date).IsZero(), chk.Equals, false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	c.Assert(err, chk.IsNil)

	destData, err := ioutil.ReadAll(downloadResp.RawResponse.Body)
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, sourceData)
	_ = downloadResp.Body(RetryReaderOptions{}).Close()
}

func (s *aztestsSuite) TestAppendBlockFromURLWithMD5(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := 8 * 1024 // 8KB
	r, sourceData := getRandomDataAndReader(contentSize)
	md5Value := md5.Sum(sourceData)
	ctx := context.Background() // Use default Background context
	srcBlob := containerClient.NewAppendBlobURL(generateName("appendsrc"))
	destBlob := containerClient.NewAppendBlobURL(generateName("appenddest"))

	// Prepare source abClient for copy.
	cResp1, err := srcBlob.Create(context.Background(), nil)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp1.RawResponse.StatusCode, chk.Equals, 201)

	appendResp, err := srcBlob.AppendBlock(context.Background(), r, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(appendResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*appendResp.BlobAppendOffset, chk.Equals, "0")
	c.Assert(*appendResp.BlobCommittedBlockCount, chk.Equals, int32(1))
	c.Assert(appendResp.ETag, chk.NotNil)
	c.Assert(appendResp.LastModified, chk.NotNil)
	c.Assert((*appendResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(appendResp.ContentMD5, chk.IsNil)
	c.Assert(appendResp.RequestID, chk.NotNil)
	c.Assert(appendResp.Version, chk.NotNil)
	c.Assert(appendResp.Date, chk.NotNil)
	c.Assert((*appendResp.Date).IsZero(), chk.Equals, false)

	// Get source abClient URL with SAS for AppendBlockFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Append block from URL.
	cResp2, err := destBlob.Create(context.Background(), nil)
	c.Assert(err, chk.IsNil)
	c.Assert(cResp2.RawResponse.StatusCode, chk.Equals, 201)

	offset := int64(0)
	count := int64(contentSize)
	contentMD5 := md5Value[:]
	appendBlockURLOptions := AppendBlockURLOptions{
		Offset:           &offset,
		Count:            &count,
		SourceContentMD5: &contentMD5,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(appendFromURLResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*appendFromURLResp.BlobAppendOffset, chk.Equals, "0")
	c.Assert(*appendFromURLResp.BlobCommittedBlockCount, chk.Equals, int32(1))
	c.Assert(appendFromURLResp.ETag, chk.NotNil)
	c.Assert(appendFromURLResp.LastModified, chk.NotNil)
	c.Assert((*appendFromURLResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(appendFromURLResp.ContentMD5, chk.NotNil)
	c.Assert(*appendFromURLResp.ContentMD5, chk.DeepEquals, contentMD5)
	c.Assert(appendFromURLResp.RequestID, chk.NotNil)
	c.Assert(appendFromURLResp.Version, chk.NotNil)
	c.Assert(appendFromURLResp.Date, chk.NotNil)
	c.Assert((*appendFromURLResp.Date).IsZero(), chk.Equals, false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, sourceData)

	// Test append block from URL with bad MD5 value
	_, badMD5 := getRandomDataAndReader(16)
	appendBlockURLOptions = AppendBlockURLOptions{
		Offset:           &offset,
		Count:            &count,
		SourceContentMD5: &badMD5,
	}
	_, err = destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeMD5Mismatch)
}

func (s *aztestsSuite) TestBlobCreateAppendMetadataNonEmpty(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := getAppendBlobClient(c, containerClient)

	createAppendBlobOptions := CreateAppendBlobOptions{
		Metadata: &basicMetadata,
	}
	_, err := abClient.Create(ctx, &createAppendBlobOptions)
	c.Assert(err, chk.IsNil)

	resp, err := abClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Metadata, chk.NotNil)
	c.Assert(resp.Metadata, chk.DeepEquals, basicMetadata)
}

func (s *aztestsSuite) TestBlobCreateAppendMetadataEmpty(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := getAppendBlobClient(c, containerClient)

	createAppendBlobOptions := CreateAppendBlobOptions{
		Metadata: &map[string]string{},
	}
	_, err := abClient.Create(ctx, &createAppendBlobOptions)
	c.Assert(err, chk.IsNil)

	resp, err := abClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Metadata, chk.IsNil)
}

func (s *aztestsSuite) TestBlobCreateAppendMetadataInvalid(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := getAppendBlobClient(c, containerClient)

	createAppendBlobOptions := CreateAppendBlobOptions{
		Metadata: &map[string]string{"In valid!": "bar"},
	}
	_, err := abClient.Create(ctx, &createAppendBlobOptions)
	c.Assert(err, chk.NotNil)
	c.Assert(strings.Contains(err.Error(), invalidHeaderErrorSubstring), chk.Equals, true)
}

func (s *aztestsSuite) TestBlobCreateAppendHTTPHeaders(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := getAppendBlobClient(c, containerClient)

	createAppendBlobOptions := CreateAppendBlobOptions{
		BlobHTTPHeaders: &basicHeaders,
	}
	_, err := abClient.Create(ctx, &createAppendBlobOptions)
	c.Assert(err, chk.IsNil)

	resp, err := abClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	h := resp.NewHTTPHeaders()
	c.Assert(h, chk.DeepEquals, basicHeaders)
}

func validateAppendBlobPut(c *chk.C, abClient AppendBlobClient) {
	resp, err := abClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Metadata, chk.NotNil)
	c.Assert(resp.Metadata, chk.DeepEquals, basicMetadata)
	c.Assert(resp.NewHTTPHeaders(), chk.DeepEquals, basicHeaders)
}

func (s *aztestsSuite) TestBlobCreateAppendIfModifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	createAppendBlobOptions := CreateAppendBlobOptions{
		BlobHTTPHeaders: &basicHeaders,
		Metadata:        &basicMetadata,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := abClient.Create(ctx, &createAppendBlobOptions)
	c.Assert(err, chk.IsNil)

	validateAppendBlobPut(c, abClient)
}

func (s *aztestsSuite) TestBlobCreateAppendIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	createAppendBlobOptions := CreateAppendBlobOptions{
		BlobHTTPHeaders: &basicHeaders,
		Metadata:        &basicMetadata,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := abClient.Create(ctx, &createAppendBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobCreateAppendIfUnmodifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	createAppendBlobOptions := CreateAppendBlobOptions{
		BlobHTTPHeaders: &basicHeaders,
		Metadata:        &basicMetadata,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := abClient.Create(ctx, &createAppendBlobOptions)
	c.Assert(err, chk.IsNil)

	validateAppendBlobPut(c, abClient)
}

func (s *aztestsSuite) TestBlobCreateAppendIfUnmodifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	createAppendBlobOptions := CreateAppendBlobOptions{
		BlobHTTPHeaders: &basicHeaders,
		Metadata:        &basicMetadata,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := abClient.Create(ctx, &createAppendBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobCreateAppendIfMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	createAppendBlobOptions := CreateAppendBlobOptions{
		BlobHTTPHeaders: &basicHeaders,
		Metadata:        &basicMetadata,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err := abClient.Create(ctx, &createAppendBlobOptions)
	c.Assert(err, chk.IsNil)

	validateAppendBlobPut(c, abClient)
}

func (s *aztestsSuite) TestBlobCreateAppendIfMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	eTag := "garbage"
	createAppendBlobOptions := CreateAppendBlobOptions{
		BlobHTTPHeaders: &basicHeaders,
		Metadata:        &basicMetadata,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := abClient.Create(ctx, &createAppendBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobCreateAppendIfNoneMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	eTag := "garbage"
	createAppendBlobOptions := CreateAppendBlobOptions{
		BlobHTTPHeaders: &basicHeaders,
		Metadata:        &basicMetadata,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := abClient.Create(ctx, &createAppendBlobOptions)
	c.Assert(err, chk.IsNil)

	validateAppendBlobPut(c, abClient)
}

func (s *aztestsSuite) TestBlobCreateAppendIfNoneMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	createAppendBlobOptions := CreateAppendBlobOptions{
		BlobHTTPHeaders: &basicHeaders,
		Metadata:        &basicMetadata,
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := abClient.Create(ctx, &createAppendBlobOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobAppendBlockNilBody(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	_, err := abClient.AppendBlock(ctx, bytes.NewReader(nil), nil)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeInvalidHeaderValue)
}

func (s *aztestsSuite) TestBlobAppendBlockEmptyBody(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	_, err := abClient.AppendBlock(ctx, strings.NewReader(""), nil)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeInvalidHeaderValue)
}

func (s *aztestsSuite) TestBlobAppendBlockNonExistentBlob(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := getAppendBlobClient(c, containerClient)

	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), nil)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeBlobNotFound)
}

func validateBlockAppended(c *chk.C, abClient AppendBlobClient, expectedSize int) {
	resp, err := abClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.ContentLength, chk.Equals, int64(expectedSize))
}

func (s *aztestsSuite) TestBlobAppendBlockIfModifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
	c.Assert(err, chk.IsNil)

	validateBlockAppended(c, abClient, len(blockBlobDefaultData))
}

func (s *aztestsSuite) TestBlobAppendBlockIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)
	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
}

// Ping Pong
func (s *aztestsSuite) TestBlobAppendBlockIfUnmodifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
	c.Assert(err, chk.IsNil)

	validateBlockAppended(c, abClient, len(blockBlobDefaultData))
}

func (s *aztestsSuite) TestBlobAppendBlockIfUnmodifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobAppendBlockIfMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
	c.Assert(err, chk.IsNil)

	validateBlockAppended(c, abClient, len(blockBlobDefaultData))
}

func (s *aztestsSuite) TestBlobAppendBlockIfMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	eTag := "garbage"
	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobAppendBlockIfNoneMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	eTag := "garbage"
	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
	c.Assert(err, chk.IsNil)

	validateBlockAppended(c, abClient, len(blockBlobDefaultData))
}

func (s *aztestsSuite) TestBlobAppendBlockIfNoneMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
}

// TODO: Fix this
//func (s *aztestsSuite) TestBlobAppendBlockIfAppendPositionMatchTrueNegOne(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	abClient, _ := createNewAppendBlob(c, containerClient)
//
//	appendPosition := int64(-1)
//	appendBlockOptions := AppendBlockOptions{
//		AppendPositionAccessConditions: &AppendPositionAccessConditions{
//			AppendPosition: &appendPosition,
//		},
//	}
//	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions) // This will cause the library to set the value of the header to 0
//	c.Assert(err, chk.NotNil)
//
//	validateBlockAppended(c, abClient, len(blockBlobDefaultData))
//}

//func (s *aztestsSuite) TestBlobAppendBlockIfAppendPositionMatchZero(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	abClient, _ := createNewAppendBlob(c, containerClient)
//
//	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), nil) // The position will not match, but the condition should be ignored
//	c.Assert(err, chk.IsNil)
//
//	appendPosition := int64(0)
//	appendBlockOptions := AppendBlockOptions{
//		AppendPositionAccessConditions: &AppendPositionAccessConditions{
//			AppendPosition: &appendPosition,
//		},
//	}
//	_, err = abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
//	c.Assert(err, chk.IsNil)
//
//	validateBlockAppended(c, abClient, 2*len(blockBlobDefaultData))
//}

func (s *aztestsSuite) TestBlobAppendBlockIfAppendPositionMatchTrueNonZero(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), nil)
	c.Assert(err, chk.IsNil)

	appendPosition := int64(len(blockBlobDefaultData))
	appendBlockOptions := AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			AppendPosition: &appendPosition,
		},
	}
	_, err = abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
	c.Assert(err, chk.IsNil)

	validateBlockAppended(c, abClient, len(blockBlobDefaultData)*2)
}

func (s *aztestsSuite) TestBlobAppendBlockIfAppendPositionMatchFalseNegOne(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), nil)
	c.Assert(err, chk.IsNil)

	appendPosition := int64(-1)
	appendBlockOptions := AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			AppendPosition: &appendPosition,
		},
	}
	_, err = abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeInvalidHeaderValue)
}

func (s *aztestsSuite) TestBlobAppendBlockIfAppendPositionMatchFalseNonZero(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	appendPosition := int64(12)
	appendBlockOptions := AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			AppendPosition: &appendPosition,
		},
	}
	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
	c.Assert(err, chk.NotNil)

	// validateStorageError(c, err, StorageErrorCodeAppendPositionConditionNotMet)
}

func (s *aztestsSuite) TestBlobAppendBlockIfMaxSizeTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	maxSize := int64(len(blockBlobDefaultData) + 1)
	appendBlockOptions := AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			MaxSize: &maxSize,
		},
	}
	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
	c.Assert(err, chk.IsNil)

	validateBlockAppended(c, abClient, len(blockBlobDefaultData))
}

func (s *aztestsSuite) TestBlobAppendBlockIfMaxSizeFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	abClient, _ := createNewAppendBlob(c, containerClient)

	maxSize := int64(len(blockBlobDefaultData) - 1)
	appendBlockOptions := AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			MaxSize: &maxSize,
		},
	}
	_, err := abClient.AppendBlock(ctx, strings.NewReader(blockBlobDefaultData), &appendBlockOptions)
	c.Assert(err, chk.NotNil)

	validateStorageError(c, err, StorageErrorCodeMaxBlobSizeConditionNotMet)
}
