// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppendBlock(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	resp, err := abClient.Create(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 201)

	appendResp, err := abClient.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	require.NoError(t, err)
	require.Equal(t, appendResp.RawResponse.StatusCode, 201)
	require.Equal(t, *appendResp.BlobAppendOffset, "0")
	require.Equal(t, *appendResp.BlobCommittedBlockCount, int32(1))
	require.NotNil(t, appendResp.ETag)
	require.NotNil(t, appendResp.LastModified)
	require.Equal(t, (*appendResp.LastModified).IsZero(), false)
	require.Nil(t, appendResp.ContentMD5)
	require.NotNil(t, appendResp.RequestID)
	require.NotNil(t, appendResp.Version)
	require.NotNil(t, appendResp.Date)
	require.Equal(t, (*appendResp.Date).IsZero(), false)

	appendResp, err = abClient.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	require.NoError(t, err)
	require.Equal(t, *appendResp.BlobAppendOffset, "1024")
	require.Equal(t, *appendResp.BlobCommittedBlockCount, int32(2))
}

func TestAppendBlockWithMD5(t *testing.T) {
	t.Skip("md5 fails authentication")
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	// set up abClient to test
	abClient := containerClient.NewAppendBlobClient(generateBlobName(testName))
	resp, err := abClient.Create(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, resp.RawResponse.StatusCode, 201)

	// test append block with valid MD5 value
	readerToBody, body := getRandomDataAndReader(1024)
	md5Value := md5.Sum(body)
	_ = body
	contentMD5 := md5Value[:]
	appendBlockOptions := AppendBlockOptions{
		TransactionalContentMD5: contentMD5,
	}
	appendResp, err := abClient.AppendBlock(context.Background(), internal.NopCloser(readerToBody), &appendBlockOptions)
	require.NoError(t, err)
	require.Equal(t, appendResp.RawResponse.StatusCode, 201)
	require.Equal(t, *appendResp.BlobAppendOffset, "0")
	require.Equal(t, *appendResp.BlobCommittedBlockCount, int32(1))
	require.NotNil(t, appendResp.ETag)
	require.NotNil(t, appendResp.LastModified)
	require.Equal(t, (*appendResp.LastModified).IsZero(), false)
	require.EqualValues(t, appendResp.ContentMD5, contentMD5)
	require.NotNil(t, appendResp.RequestID)
	require.NotNil(t, appendResp.Version)
	require.NotNil(t, appendResp.Date)
	require.Equal(t, (*appendResp.Date).IsZero(), false)

	// test append block with bad MD5 value
	readerToBody, body = getRandomDataAndReader(1024)
	_, badMD5 := getRandomDataAndReader(16)
	_ = body
	appendBlockOptions = AppendBlockOptions{
		TransactionalContentMD5: badMD5,
	}
	appendResp, err = abClient.AppendBlock(context.Background(), internal.NopCloser(readerToBody), &appendBlockOptions)
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeMD5Mismatch)
}

func TestAppendBlockFromURL(t *testing.T) {
	recording.LiveOnly(t) // Live only because of random data and random name
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	//ctx := context.Background()
	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	contentMD5 := md5.Sum(sourceData)
	srcBlob := containerClient.NewAppendBlobClient(generateName("appendsrc"))
	destBlob := containerClient.NewAppendBlobClient(generateName("appenddest"))

	// Prepare source abClient for copy.
	cResp1, err := srcBlob.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp1.RawResponse.StatusCode, 201)

	appendResp, err := srcBlob.AppendBlock(context.Background(), internal.NopCloser(r), nil)
	require.NoError(t, err)
	require.NoError(t, err)
	require.Equal(t, appendResp.RawResponse.StatusCode, 201)
	require.Equal(t, *appendResp.BlobAppendOffset, "0")
	require.Equal(t, *appendResp.BlobCommittedBlockCount, int32(1))
	require.NotNil(t, appendResp.ETag)
	require.NotNil(t, appendResp.LastModified)
	require.Equal(t, (*appendResp.LastModified).IsZero(), false)
	require.Nil(t, appendResp.ContentMD5)
	require.NotNil(t, appendResp.RequestID)
	require.NotNil(t, appendResp.Version)
	require.NotNil(t, appendResp.Date)
	require.Equal(t, (*appendResp.Date).IsZero(), false)

	// Get source abClient URL with SAS for AppendBlockFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	credential, err := getCredential(testAccountDefault)
	require.NoError(t, err)

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	require.NoError(t, err)

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Append block from URL.
	cResp2, err := destBlob.Create(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, cResp2.RawResponse.StatusCode, 201)

	//ctx context.Context, source url.URL, contentLength int64, options *AppendBlockURLOptions)
	offset := int64(0)
	count := int64(CountToEnd)
	appendBlockURLOptions := AppendBlockURLOptions{
		Offset: &offset,
		Count:  &count,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	require.NoError(t, err)
	require.Equal(t, appendFromURLResp.RawResponse.StatusCode, 201)
	require.Equal(t, *appendFromURLResp.BlobAppendOffset, "0")
	require.Equal(t, *appendFromURLResp.BlobCommittedBlockCount, int32(1))
	require.NotNil(t, appendFromURLResp.ETag)
	require.NotNil(t, appendFromURLResp.LastModified)
	require.Equal(t, (*appendFromURLResp.LastModified).IsZero(), false)
	require.NotNil(t, appendFromURLResp.ContentMD5)
	require.EqualValues(t, appendFromURLResp.ContentMD5, contentMD5[:])
	require.NotNil(t, appendFromURLResp.RequestID)
	require.NotNil(t, appendFromURLResp.Version)
	require.NotNil(t, appendFromURLResp.Date)
	require.Equal(t, (*appendFromURLResp.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	require.NoError(t, err)

	destData, err := ioutil.ReadAll(downloadResp.RawResponse.Body)
	require.NoError(t, err)
	require.Equal(t, destData, sourceData)
	err = downloadResp.Body(nil).Close()
	require.NoError(t, err)
}

//nolint
func TestAppendBlockFromURLWithMD5(t *testing.T) {
	t.Skip("md5 fails authentication")
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, sourceData := getRandomDataAndReader(contentSize)
	md5Value := md5.Sum(sourceData)
	ctx := context.Background() // Use default Background context
	srcBlob := containerClient.NewAppendBlobClient(generateName("appendsrc"))
	destBlob := containerClient.NewAppendBlobClient(generateName("appenddest"))

	// Prepare source abClient for copy.
	cResp1, err := srcBlob.Create(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, cResp1.RawResponse.StatusCode, 201)

	appendResp, err := srcBlob.AppendBlock(context.Background(), internal.NopCloser(r), nil)
	require.NoError(t, err)
	require.Equal(t, appendResp.RawResponse.StatusCode, 201)
	require.Equal(t, *appendResp.BlobAppendOffset, "0")
	require.Equal(t, *appendResp.BlobCommittedBlockCount, int32(1))
	require.NotNil(t, appendResp.ETag)
	require.NotNil(t, appendResp.LastModified)
	require.Equal(t, (*appendResp.LastModified).IsZero(), false)
	require.Nil(t, appendResp.ContentMD5)
	require.NotNil(t, appendResp.RequestID)
	require.NotNil(t, appendResp.Version)
	require.NotNil(t, appendResp.Date)
	require.Equal(t, (*appendResp.Date).IsZero(), false)

	// Get source abClient URL with SAS for AppendBlockFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	credential, err := getGenericCredential(nil, testAccountDefault)
	require.NoError(t, err)

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	require.NoError(t, err)

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Append block from URL.
	cResp2, err := destBlob.Create(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, cResp2.RawResponse.StatusCode, 201)

	offset := int64(0)
	count := int64(contentSize)
	contentMD5 := md5Value[:]
	appendBlockURLOptions := AppendBlockURLOptions{
		Offset:           &offset,
		Count:            &count,
		SourceContentMD5: contentMD5,
	}
	appendFromURLResp, err := destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	require.NoError(t, err)
	require.Equal(t, appendFromURLResp.RawResponse.StatusCode, 201)
	require.Equal(t, *appendFromURLResp.BlobAppendOffset, "0")
	require.Equal(t, *appendFromURLResp.BlobCommittedBlockCount, int32(1))
	require.NotNil(t, appendFromURLResp.ETag)
	require.NotNil(t, appendFromURLResp.LastModified)
	require.Equal(t, (*appendFromURLResp.LastModified).IsZero(), false)
	require.NotNil(t, appendFromURLResp.ContentMD5)
	require.EqualValues(t, appendFromURLResp.ContentMD5, contentMD5)
	require.NotNil(t, appendFromURLResp.RequestID)
	require.NotNil(t, appendFromURLResp.Version)
	require.NotNil(t, appendFromURLResp.Date)
	require.Equal(t, (*appendFromURLResp.Date).IsZero(), false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	require.NoError(t, err)
	destData, err := ioutil.ReadAll(downloadResp.Body(nil))
	require.NoError(t, err)
	require.EqualValues(t, destData, sourceData)

	// Test append block from URL with bad MD5 value
	_, badMD5 := getRandomDataAndReader(16)
	appendBlockURLOptions = AppendBlockURLOptions{
		Offset:           &offset,
		Count:            &count,
		SourceContentMD5: badMD5,
	}
	_, err = destBlob.AppendBlockFromURL(ctx, srcBlobURLWithSAS, &appendBlockURLOptions)
	require.Error(t, err)
	validateStorageError(_assert, err, StorageErrorCodeMD5Mismatch)
}

func TestBlobCreateAppendMetadataNonEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	_, err = abClient.Create(ctx, &CreateAppendBlobOptions{
		Metadata: basicMetadata,
	})
	require.NoError(t, err)

	resp, err := abClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Metadata)
	require.EqualValues(t, resp.Metadata, basicMetadata)
}

func TestBlobCreateAppendMetadataEmpty(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := CreateAppendBlobOptions{
		Metadata: map[string]string{},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	require.NoError(t, err)

	resp, err := abClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, resp.Metadata)
}

func TestBlobCreateAppendMetadataInvalid(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := CreateAppendBlobOptions{
		Metadata: map[string]string{"In valid!": "bar"},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	require.Error(t, err)
	require.Contains(t, err.Error(), invalidHeaderErrorSubstring)
}

func TestBlobCreateAppendHTTPHeaders(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	require.NoError(t, err)

	resp, err := abClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	h := resp.GetHTTPHeaders()
	require.EqualValues(t, h, basicHeaders)
}

func validateAppendBlobPut(t *testing.T, abClient AppendBlobClient) {
	resp, err := abClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, resp.Metadata)
	require.EqualValues(t, resp.Metadata, basicMetadata)
	require.EqualValues(t, resp.GetHTTPHeaders(), basicHeaders)
}

func TestBlobCreateAppendIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	require.NoError(t, err)
	require.Equal(t, appendBlobCreateResp.RawResponse.StatusCode, 201)
	require.NotNil(t, appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	require.NoError(t, err)

	validateAppendBlobPut(t, abClient)
}

func TestBlobCreateAppendIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	require.NoError(t, err)
	require.Equal(t, appendBlobCreateResp.RawResponse.StatusCode, 201)
	require.NotNil(t, appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobCreateAppendIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	require.NoError(t, err)
	require.Equal(t, appendBlobCreateResp.RawResponse.StatusCode, 201)
	require.NotNil(t, appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	require.NoError(t, err)

	validateAppendBlobPut(t, abClient)
}

func TestBlobCreateAppendIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	require.NoError(t, err)
	require.Equal(t, appendBlobCreateResp.RawResponse.StatusCode, 201)
	require.NotNil(t, appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobCreateAppendIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	resp, err := abClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	require.NoError(t, err)

	validateAppendBlobPut(t, abClient)
}

func TestBlobCreateAppendIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: to.StringPtr("garbage"),
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobCreateAppendIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	eTag := "garbage"
	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	require.NoError(t, err)

	validateAppendBlobPut(t, abClient)
}

func TestBlobCreateAppendIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	resp, err := abClient.GetProperties(ctx, nil)
	require.NoError(t, err)

	createAppendBlobOptions := CreateAppendBlobOptions{
		HTTPHeaders: &basicHeaders,
		Metadata:    basicMetadata,
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = abClient.Create(ctx, &createAppendBlobOptions)
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobAppendBlockNilBody(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(bytes.NewReader(nil)), nil)
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidHeaderValue)
}

func TestBlobAppendBlockEmptyBody(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader("")), nil)
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeInvalidHeaderValue)
}

func TestBlobAppendBlockNonExistentBlob(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeBlobNotFound)
}

func validateBlockAppended(t *testing.T, abClient AppendBlobClient, expectedSize int) {
	resp, err := abClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *resp.ContentLength, int64(expectedSize))
}

func TestBlobAppendBlockIfModifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	require.NoError(t, err)
	require.Equal(t, appendBlobCreateResp.RawResponse.StatusCode, 201)
	require.NotNil(t, appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
	require.NoError(t, err)

	validateBlockAppended(t, abClient, len(blockBlobDefaultData))
}

func TestBlobAppendBlockIfModifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	require.NoError(t, err)
	require.Equal(t, appendBlobCreateResp.RawResponse.StatusCode, 201)
	require.NotNil(t, appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobAppendBlockIfUnmodifiedSinceTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	require.NoError(t, err)
	require.Equal(t, appendBlobCreateResp.RawResponse.StatusCode, 201)
	require.NotNil(t, appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, 10)

	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
	require.NoError(t, err)

	validateBlockAppended(t, abClient, len(blockBlobDefaultData))
}

func TestBlobAppendBlockIfUnmodifiedSinceFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := getAppendBlobClient(abName, containerClient)

	appendBlobCreateResp, err := abClient.Create(ctx, nil)

	require.NoError(t, err)
	require.Equal(t, appendBlobCreateResp.RawResponse.StatusCode, 201)
	require.NotNil(t, appendBlobCreateResp.Date)

	currentTime := getRelativeTimeFromAnchor(appendBlobCreateResp.Date, -10)

	appendBlockOptions := AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
	require.Error(t, err)

	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobAppendBlockIfMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	})
	require.NoError(t, err)

	validateBlockAppended(t, abClient, len(blockBlobDefaultData))
}

func TestBlobAppendBlockIfMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfMatch: to.StringPtr("garbage"),
			},
		},
	})
	require.Error(t, err)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

func TestBlobAppendBlockIfNoneMatchTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: to.StringPtr("garbage"),
			},
		},
	})
	require.NoError(t, err)
	validateBlockAppended(t, abClient, len(blockBlobDefaultData))
}

func TestBlobAppendBlockIfNoneMatchFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	resp, _ := abClient.GetProperties(ctx, nil)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		BlobAccessConditions: &BlobAccessConditions{
			ModifiedAccessConditions: &ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	})
	require.Error(t, err)
	validateStorageError(_assert, err, StorageErrorCodeConditionNotMet)
}

//// TODO: Fix this
////func (s *azblobTestSuite) TestBlobAppendBlockIfAppendPositionMatchTrueNegOne() {
////	bsu := getServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(_assert, containerClient)
////	abClient, _ := createNewAppendBlob(c, containerClient)
////
////	appendPosition := int64(-1)
////	appendBlockOptions := AppendBlockOptions{
////		AppendPositionAccessConditions: &AppendPositionAccessConditions{
////			AppendPosition: &appendPosition,
////		},
////	}
////	_, err := abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions) // This will cause the library to set the value of the header to 0
////	_assert.Error(err)
////
////	validateBlockAppended(c, abClient, len(blockBlobDefaultData))
////}
//
////func (s *azblobTestSuite) TestBlobAppendBlockIfAppendPositionMatchZero() {
////	bsu := getServiceClient()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(_assert, containerClient)
////	abClient, _ := createNewAppendBlob(c, containerClient)
////
////	_, err := abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil) // The position will not match, but the condition should be ignored
////	_assert.NoError(err)
////
////	appendPosition := int64(0)
////	appendBlockOptions := AppendBlockOptions{
////		AppendPositionAccessConditions: &AppendPositionAccessConditions{
////			AppendPosition: &appendPosition,
////		},
////	}
////	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &appendBlockOptions)
////	_assert.NoError(err)
////
////	validateBlockAppended(c, abClient, 2*len(blockBlobDefaultData))
////}

func TestBlobAppendBlockIfAppendPositionMatchTrueNonZero(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			AppendPosition: to.Int64Ptr(int64(len(blockBlobDefaultData))),
		},
	})
	require.NoError(t, err)

	validateBlockAppended(t, abClient, len(blockBlobDefaultData)*2)
}

func TestBlobAppendBlockIfAppendPositionMatchFalseNegOne(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), nil)
	require.NoError(t, err)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			AppendPosition: to.Int64Ptr(-1),
		},
	})
	require.Error(t, err)
	validateStorageError(_assert, err, StorageErrorCodeInvalidHeaderValue)
}

func TestBlobAppendBlockIfAppendPositionMatchFalseNonZero(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			AppendPosition: to.Int64Ptr(12),
		},
	})
	require.Error(t, err)
	validateStorageError(_assert, err, StorageErrorCodeAppendPositionConditionNotMet)
}

func TestBlobAppendBlockIfMaxSizeTrue(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			MaxSize: to.Int64Ptr(int64(len(blockBlobDefaultData) + 1)),
		},
	})
	require.NoError(t, err)
	validateBlockAppended(t, abClient, len(blockBlobDefaultData))
}

func TestBlobAppendBlockIfMaxSizeFalse(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.AppendBlock(ctx, internal.NopCloser(strings.NewReader(blockBlobDefaultData)), &AppendBlockOptions{
		AppendPositionAccessConditions: &AppendPositionAccessConditions{
			MaxSize: to.Int64Ptr(int64(len(blockBlobDefaultData) - 1)),
		},
	})
	require.Error(t, err)
	validateStorageError(_assert, err, StorageErrorCodeMaxBlobSizeConditionNotMet)
}

func TestSealAppendBlob(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	appendResp, err := abClient.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	require.NoError(t, err)
	require.Equal(t, appendResp.RawResponse.StatusCode, 201)
	require.Equal(t, *appendResp.BlobAppendOffset, "0")
	require.Equal(t, *appendResp.BlobCommittedBlockCount, int32(1))

	sealResp, err := abClient.SealAppendBlob(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *sealResp.IsSealed, true)

	appendResp, err = abClient.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	require.Error(t, err)
	validateStorageError(_assert, err, "BlobIsSealed")

	getPropResp, err := abClient.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *getPropResp.IsSealed, true)
}

// TODO: Learn about the behaviour of AppendPosition
// nolint
//func (s *azblobUnrecordedTestSuite) TestSealAppendBlobWithAppendConditions() {
//	_assert := assert.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := getServiceClient(nil, testAccountDefault, nil)
//	if err != nil {
//		s.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := generateContainerName(testName)
//	containerClient := createNewContainer(t, containerName, svcClient)
//	defer deleteContainer(_assert, containerClient)
//
//	abName := generateBlobName(testName)
//	abClient := createNewAppendBlob(_assert, abName, containerClient)
//
//	sealResp, err := abClient.SealAppendBlob(ctx, &SealAppendBlobOptions{
//		AppendPositionAccessConditions: &AppendPositionAccessConditions{
//			AppendPosition: to.Int64Ptr(1),
//		},
//	})
//	_assert.Error(err)
//	_ = sealResp
//
//	sealResp, err = abClient.SealAppendBlob(ctx, &SealAppendBlobOptions{
//		AppendPositionAccessConditions: &AppendPositionAccessConditions{
//			AppendPosition: to.Int64Ptr(0),
//		},
//	})
//}

func TestCopySealedBlob(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	_, err = abClient.SealAppendBlob(ctx, nil)
	require.NoError(t, err)

	copiedBlob1 := getAppendBlobClient("copy1"+abName, containerClient)
	// copy sealed blob will get a sealed blob
	_, err = copiedBlob1.StartCopyFromURL(ctx, abClient.URL(), nil)
	require.NoError(t, err)

	getResp1, err := copiedBlob1.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp1.IsSealed, true)

	_, err = copiedBlob1.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	require.Error(t, err)
	validateStorageError(_assert, err, "BlobIsSealed")

	copiedBlob2 := getAppendBlobClient("copy2"+abName, containerClient)
	_, err = copiedBlob2.StartCopyFromURL(ctx, abClient.URL(), &StartCopyBlobOptions{
		SealBlob: to.BoolPtr(true),
	})
	require.NoError(t, err)

	getResp2, err := copiedBlob2.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp2.IsSealed, true)

	_, err = copiedBlob2.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	require.Error(t, err)
	validateStorageError(_assert, err, "BlobIsSealed")

	copiedBlob3 := getAppendBlobClient("copy3"+abName, containerClient)
	_, err = copiedBlob3.StartCopyFromURL(ctx, abClient.URL(), &StartCopyBlobOptions{
		SealBlob: to.BoolPtr(false),
	})
	require.NoError(t, err)

	getResp3, err := copiedBlob3.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Nil(t, getResp3.IsSealed)

	appendResp3, err := copiedBlob3.AppendBlock(context.Background(), getReaderToGeneratedBytes(1024), nil)
	require.NoError(t, err)
	require.Equal(t, appendResp3.RawResponse.StatusCode, 201)
	require.Equal(t, *appendResp3.BlobAppendOffset, "0")
	require.Equal(t, *appendResp3.BlobCommittedBlockCount, int32(1))
}

func TestCopyUnsealedBlob(t *testing.T) {
	stop := start(t)
	defer stop()

	_assert := assert.New(t)
	testName := t.Name()
	svcClient, err := createServiceClient(t, testAccountDefault)
	require.NoError(t, err)

	containerName := generateContainerName(testName)
	containerClient := createNewContainer(t, containerName, svcClient)
	defer deleteContainer(_assert, containerClient)

	abName := generateBlobName(testName)
	abClient := createNewAppendBlob(_assert, abName, containerClient)

	copiedBlob := getAppendBlobClient("copy"+abName, containerClient)
	_, err = copiedBlob.StartCopyFromURL(ctx, abClient.URL(), &StartCopyBlobOptions{
		SealBlob: to.BoolPtr(true),
	})
	require.NoError(t, err)

	getResp, err := copiedBlob.GetProperties(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, *getResp.IsSealed, true)
}
