// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func (s *aztestsSuite) TestCreateBlobClient() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)

	blobName := generateBlobName(s.T().Name())
	blobClient := getBlockBlobClient(blobName, containerClient)

	blobURLParts := NewBlobURLParts(blobClient.URL())
	_assert.Equal(blobURLParts.BlobName, blobName)
	_assert.Equal(blobURLParts.ContainerName, containerName)

	correctURL := "https://" + os.Getenv(AccountNameEnvVar) + "." + DefaultBlobEndpointSuffix + containerName + "/" + blobName
	_assert.Equal(blobClient.URL(), correctURL)
}

func (s *aztestsSuite) TestCreateBlobClientWithSnapshotAndSAS() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	containerName := generateContainerName(s.T().Name())
	containerClient := getContainerClient(containerName, bsu)

	blobName := generateBlobName(s.T().Name())
	blobClient := getBlockBlobClient(blobName, containerClient)

	currentTime, err := time.Parse(time.UnixDate, "Fri Jun 11 20:00:00 UTC 2049")
	_assert.Nil(err)

	credential, err := getGenericCredential("")
	if err != nil {
		s.Fail(err.Error())
	}
	sasQueryParams, err := AccountSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    currentTime,
		Permissions:   AccountSASPermissions{Read: true, List: true}.String(),
		Services:      AccountSASServices{Blob: true}.String(),
		ResourceTypes: AccountSASResourceTypes{Container: true, Object: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		s.Fail(err.Error())
	}

	parts := NewBlobURLParts(blobClient.URL())
	parts.SAS = sasQueryParams
	parts.Snapshot = currentTime.Format(SnapshotTimeFormat)
	blobURLParts := parts.URL()

	// The snapshot format string is taken from the snapshotTimeFormat value in parsing_urls.go. The field is not public, so
	// it is copied here
	correctURL := "https://" + os.Getenv(AccountNameEnvVar) + "." + DefaultBlobEndpointSuffix + containerName + "/" + blobName +
		"?" + "snapshot=" + currentTime.Format("2006-01-02T15:04:05.0000000Z07:00") + "&" + sasQueryParams.Encode()
	_assert.Equal(blobURLParts, correctURL)
}

//// func (s *aztestsSuite) TestBlobWithNewPipeline() {
//// 	bsu := getBSU()
//// 	containerClient, _ := getContainerClient(c, bsu)
//// 	blobClient := containerClient.NewBlockBlobClient(blobPrefix)
////
//// 	newBlobClient := blobClient.WithPipeline(newTestPipeline())
////
//// 	// exercise the new pipeline
//// 	_, err := newBlobClient.GetAccountInfo(ctx)
//// 	_assert(err, chk.NotNil)
//// 	_assert(err.Error(), chk.Equals, testPipelineMessage)
//// }

func waitForCopy(_assert *assert.Assertions, copyBlobClient BlockBlobClient, blobCopyResponse BlobStartCopyFromURLResponse) {
	status := *blobCopyResponse.CopyStatus
	// Wait for the copy to finish. If the copy takes longer than a minute, we will fail
	start := time.Now()
	for status != CopyStatusSuccess {
		props, _ := copyBlobClient.GetProperties(ctx, nil)
		status = *props.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			_assert.Fail("If the copy takes longer than a minute, we will fail")
		}
	}
}

func (s *aztestsSuite) TestBlobStartCopyDestEmpty() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, bsu)

	_, err := containerClient.Create(ctx, nil)
	_assert.Nil(err)
	defer deleteContainer(containerClient)

	blobClient, blobName := createNewBlockBlob(_assert, testName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	blobCopyResponse, err := copyBlobClient.StartCopyFromURL(ctx, blobClient.URL(), nil)
	_assert.Nil(err)
	waitForCopy(_assert, copyBlobClient, blobCopyResponse)

	resp, err := copyBlobClient.Download(ctx, nil)
	_assert.Nil(err)

	// Read the blob data to verify the copy
	data, err := ioutil.ReadAll(resp.RawResponse.Body)
	_assert.Equal(*resp.ContentLength, int64(len(blockBlobDefaultData)))
	_assert.Equal(string(data), blockBlobDefaultData)
	resp.Body(RetryReaderOptions{}).Close()
}

func (s *aztestsSuite) TestBlobStartCopyMetadata() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerName := generateContainerName(testName)
	containerClient := getContainerClient(containerName, bsu)

	_, err := containerClient.Create(ctx, nil)
	_assert.Nil(err)
	defer deleteContainer(containerClient)

	blobClient, blobName := createNewBlockBlob(_assert, testName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	metadata := make(map[string]string)
	metadata["Bla"] = "foo"
	options := StartCopyBlobOptions{
		Metadata: &metadata,
	}
	resp, err := copyBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	_assert.Nil(err)
	waitForCopy(_assert, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.EqualValues(resp2.Metadata, metadata)
}

func (s *aztestsSuite) TestBlobStartCopyMetadataNil() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerClient, _ := createNewContainer(_assert, testName, bsu)
	defer deleteContainer(containerClient)

	blobClient, _ := createNewBlockBlob(_assert, testName, containerClient)

	anotherBlobName := "copy" + generateBlobName(testName)
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	// Have the destination start with metadata so we ensure the nil metadata passed later takes effect
	_, err := copyBlobClient.Upload(ctx, azcore.NopCloser(bytes.NewReader([]byte("data"))), nil)
	_assert.Nil(err)

	resp, err := copyBlobClient.StartCopyFromURL(ctx, blobClient.URL(), nil)
	_assert.Nil(err)

	waitForCopy(_assert, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp2.Metadata, 0)
}

func (s *aztestsSuite) TestBlobStartCopyMetadataEmpty() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerClient, _ := createNewContainer(_assert, testName, bsu)
	defer deleteContainer(containerClient)

	blobClient, _ := createNewBlockBlob(_assert, testName, containerClient)

	anotherBlobName := "copy" + generateBlobName(testName)
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	// Have the destination start with metadata so we ensure the empty metadata passed later takes effect
	_, err := copyBlobClient.Upload(ctx, azcore.NopCloser(bytes.NewReader([]byte("data"))), nil)
	_assert.Nil(err)

	metadata := make(map[string]string)
	options := StartCopyBlobOptions{
		Metadata: &metadata,
	}
	resp, err := copyBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	_assert.Nil(err)

	waitForCopy(_assert, copyBlobClient, resp)

	resp2, err := copyBlobClient.GetProperties(ctx, nil)
	_assert.Nil(err)
	_assert.Len(resp2.Metadata, 0)
}

func (s *aztestsSuite) TestBlobStartCopyMetadataInvalidField() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerClient, _ := createNewContainer(_assert, testName, bsu)
	defer deleteContainer(containerClient)

	blobClient, _ := createNewBlockBlob(_assert, testName, containerClient)

	anotherBlobName := "copy" + generateBlobName(testName)
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	metadata := make(map[string]string)
	metadata["I nvalid."] = "foo"
	options := StartCopyBlobOptions{
		Metadata: &metadata,
	}
	_, err := copyBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	_assert.NotNil(err)
	_assert.Equal(strings.Contains(err.Error(), invalidHeaderErrorSubstring), true)
}

func (s *aztestsSuite) TestBlobStartCopySourceNonExistent() {
	_assert := assert.New(s.T())
	context := getTestContext(s.T().Name())
	bsu := getBSU(&ClientOptions{
		HTTPClient: context.recording,
		Retry:      azcore.RetryOptions{MaxRetries: -1}})
	testName := s.T().Name()
	containerClient, _ := createNewContainer(_assert, testName, bsu)
	defer deleteContainer(containerClient)

	blobName := generateBlobName(testName)
	blobClient := getBlockBlobClient(blobName, containerClient)

	anotherBlobName := "copy" + blobName
	copyBlobClient := getBlockBlobClient(anotherBlobName, containerClient)

	_, err := copyBlobClient.StartCopyFromURL(ctx, blobClient.URL(), nil)
	_assert.NotNil(err)
	_assert.Equal(strings.Contains(err.Error(), "not exist"), true)
}

//func (s *aztestsSuite) TestBlobStartCopySourcePrivate() {
//	_assert := assert.New(s.T())
//	context := getTestContext(s.T().Name())
//	bsu := getBSU(&ClientOptions{
//		HTTPClient: context.recording,
//		Retry: azcore.RetryOptions{MaxRetries: -1}})
//	testName := s.T().Name()
//	containerClient, _ := createNewContainer(_assert, testName, bsu)
//	defer deleteContainer(containerClient)
//
//	_, err := containerClient.SetAccessPolicy(ctx, nil)
//	_assert.Nil(err)
//
//	blobClient, _ := createNewBlockBlob(_assert, testName, containerClient)
//
//	bsu2, err := getAlternateBSU()

//	if err != nil {
//		s.T().Skip(err.Error())
//		return
//	}
//
//	copyContainerClient, _ := createNewContainer(_assert, "cpyc" + testName, bsu2)
//	defer deleteContainer(copyContainerClient)
//	copyBlobName := "copyb" + generateBlobName(testName)
//	copyBlobClient := getBlockBlobClient(copyBlobName, copyContainerClient)
//
//	if bsu.URL() == bsu2.URL() {
//		s.T().Skip("Test not valid because primary and secondary accounts are the same")
//	}
//	_, err = copyBlobClient.StartCopyFromURL(ctx, blobClient.URL(), nil)
//	validateStorageError(_assert, err, StorageErrorCodeCannotVerifyCopySource)
//}

//func (s *aztestsSuite) TestBlobStartCopyUsingSASSrc() {
//	bsu := getBSU(nil)
//	containerClient, containerName := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	_, err := containerClient.SetAccessPolicy(ctx, nil)
//	_assert(err, chk.IsNil)
//	blobClient, blobName := createNewBlockBlob(c, containerClient)
//
//	// Create sas values for the source blob
//	credential, err := getGenericCredential("")
//	if err != nil {
//		c.Fatal("Invalid credential")
//	}
//	serviceSASValues := BlobSASSignatureValues{StartTime: time.Now().Add(-1 * time.Hour).UTC(),
//		ExpiryTime: time.Now().Add(time.Hour).UTC(), Permissions: BlobSASPermissions{Read: true, Write: true}.String(),
//		ContainerName: containerName, BlobName: blobName}
//	queryParams, err := serviceSASValues.NewSASQueryParameters(credential)
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	// Create URLs to the destination blob with sas parameters
//	sasURL := NewBlobURLParts(blobClient.URL())
//	sasURL.SAS = queryParams
//
//	// Create a new container for the destination
//	bsu2, err := getAlternateBSU()
//	if err != nil {
//		c.Skip(err.Error())
//		return
//	}
//	copyContainerClient, _ := createNewContainer(c, bsu2)
//	defer deleteContainer(copyContainerClient)
//	copyBlobClient, _ := getBlockBlobClient(c, copyContainerClient)
//
//	resp, err := copyBlobClient.StartCopyFromURL(ctx, sasURL.URL(), nil)
//	_assert(err, chk.IsNil)
//
//	waitForCopy(c, copyBlobClient, resp)
//
//	offset, count := int64(0), int64(len(blockBlobDefaultData))
//	downloadBlobOptions := DownloadBlobOptions{
//		Offset: &offset,
//		Count:  &count,
//	}
//	resp2, err := copyBlobClient.Download(ctx, &downloadBlobOptions)
//	_assert(err, chk.IsNil)
//
//	data, err := ioutil.ReadAll(resp2.RawResponse.Body)
//	_assert(*resp2.ContentLength, chk.Equals, int64(len(blockBlobDefaultData)))
//	_assert(string(data), chk.Equals, blockBlobDefaultData)
//	resp2.Body(RetryReaderOptions{}).Close()
//}
//
//func (s *aztestsSuite) TestBlobStartCopyUsingSASDest() {
//	bsu := getBSU(nil)
//	containerClient, containerName := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	_, err := containerClient.SetAccessPolicy(ctx, nil)
//	_assert(err, chk.IsNil)
//	blobClient, blobName := createNewBlockBlob(c, containerClient)
//	_ = blobClient
//
//	// Generate SAS on the source
//	serviceSASValues := BlobSASSignatureValues{ExpiryTime: time.Now().Add(time.Hour).UTC(),
//		Permissions: BlobSASPermissions{Read: true, Write: true, Create: true}.String(), ContainerName: containerName, BlobName: blobName}
//	credential, err := getGenericCredential("")
//	if err != nil {
//		c.Fatal("Invalid credential")
//	}
//	queryParams, err := serviceSASValues.NewSASQueryParameters(credential)
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	// Create destination container
//	bsu2, err := getAlternateBSU()
//	if err != nil {
//		c.Skip(err.Error())
//		return
//	}
//
//	copyContainerClient, copyContainerName := createNewContainer(c, bsu2)
//	defer deleteContainer(copyContainerClient)
//	copyBlobClient, copyBlobName := getBlockBlobClient(c, copyContainerClient)
//
//	// Generate Sas for the destination
//	credential, err = getGenericCredential("SECONDARY_")
//	if err != nil {
//		c.Fatal("Invalid secondary credential")
//	}
//	copyServiceSASvalues := BlobSASSignatureValues{StartTime: time.Now().Add(-1 * time.Hour).UTC(),
//		ExpiryTime: time.Now().Add(time.Hour).UTC(), Permissions: BlobSASPermissions{Read: true, Write: true}.String(),
//		ContainerName: copyContainerName, BlobName: copyBlobName}
//	copyQueryParams, err := copyServiceSASvalues.NewSASQueryParameters(credential)
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	// Generate anonymous URL to destination with SAS
//	anonURL := NewBlobURLParts(bsu2.URL())
//	anonURL.SAS = copyQueryParams
//	anonymousBSU, err := NewServiceClient(anonURL.URL(), azcore.AnonymousCredential(), nil)
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	anonymousContainerClient := anonymousBSU.NewContainerClient(copyContainerName)
//	anonymousBlobClient := anonymousContainerClient.NewBlockBlobClient(copyBlobName)
//
//	// Apply sas to source
//	srcBlobWithSasURL := NewBlobURLParts(blobClient.URL())
//	srcBlobWithSasURL.SAS = queryParams
//
//	resp, err := anonymousBlobClient.StartCopyFromURL(ctx, srcBlobWithSasURL.URL(), nil)
//	_assert(err, chk.IsNil)
//
//	// Allow copy to happen
//	waitForCopy(c, anonymousBlobClient, resp)
//
//	offset, count := int64(0), int64(len(blockBlobDefaultData))
//	downloadBlobOptions := DownloadBlobOptions{
//		Offset: &offset,
//		Count:  &count,
//	}
//	resp2, err := copyBlobClient.Download(ctx, &downloadBlobOptions)
//	_assert(err, chk.IsNil)
//
//	data, err := ioutil.ReadAll(resp2.RawResponse.Body)
//	_, err = resp2.Body(RetryReaderOptions{}).Read(data)
//	_assert(*resp2.ContentLength, chk.Equals, int64(len(blockBlobDefaultData)))
//	_assert(string(data), chk.Equals, blockBlobDefaultData)
//	resp2.Body(RetryReaderOptions{}).Close()
//}
//
//func (s *aztestsSuite) TestBlobStartCopySourceIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	accessConditions := SourceModifiedAccessConditions{
//		SourceIfModifiedSince: &currentTime,
//	}
//	options := StartCopyBlobOptions{
//		SourceModifiedAccessConditions: &accessConditions,
//	}
//	destBlobClient, _ := getBlockBlobClient(c, containerClient)
//	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.IsNil)
//
//	_, err = destBlobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopySourceIfModifiedSinceFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(100)
//	accessConditions := SourceModifiedAccessConditions{
//		SourceIfModifiedSince: &currentTime,
//	}
//	options := StartCopyBlobOptions{
//		SourceModifiedAccessConditions: &accessConditions,
//	}
//	destBlobClient, _ := getBlockBlobClient(c, containerClient)
//	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopySourceIfUnmodifiedSinceTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//	accessConditions := SourceModifiedAccessConditions{
//		SourceIfUnmodifiedSince: &currentTime,
//	}
//	options := StartCopyBlobOptions{
//		SourceModifiedAccessConditions: &accessConditions,
//	}
//	destBlobClient, _ := getBlockBlobClient(c, containerClient)
//	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.IsNil)
//
//	_, err = destBlobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopySourceIfUnmodifiedSinceFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//	accessConditions := SourceModifiedAccessConditions{
//		SourceIfUnmodifiedSince: &currentTime,
//	}
//	options := StartCopyBlobOptions{
//		SourceModifiedAccessConditions: &accessConditions,
//	}
//	destBlobClient, _ := getBlockBlobClient(c, containerClient)
//	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopySourceIfMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//
//	accessConditions := SourceModifiedAccessConditions{
//		SourceIfMatch: resp.ETag,
//	}
//	options := StartCopyBlobOptions{
//		SourceModifiedAccessConditions: &accessConditions,
//	}
//	destBlobClient, _ := getBlockBlobClient(c, containerClient)
//	_, err = destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.IsNil)
//
//	_, err = destBlobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopySourceIfMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	randomEtag := "a"
//	accessConditions := SourceModifiedAccessConditions{
//		SourceIfMatch: &randomEtag,
//	}
//	options := StartCopyBlobOptions{
//		SourceModifiedAccessConditions: &accessConditions,
//	}
//
//	destBlobClient, _ := getBlockBlobClient(c, containerClient)
//	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopySourceIfNoneMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	randomEtag := "a"
//	accessConditions := SourceModifiedAccessConditions{
//		SourceIfNoneMatch: &randomEtag,
//	}
//	options := StartCopyBlobOptions{
//		SourceModifiedAccessConditions: &accessConditions,
//	}
//	destBlobClient, _ := getBlockBlobClient(c, containerClient)
//	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.IsNil)
//
//	_, err = destBlobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopySourceIfNoneMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//
//	accessConditions := SourceModifiedAccessConditions{
//		SourceIfNoneMatch: resp.ETag,
//	}
//	options := StartCopyBlobOptions{
//		SourceModifiedAccessConditions: &accessConditions,
//	}
//
//	destBlobClient, _ := getBlockBlobClient(c, containerClient)
//	_, err = destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopyDestIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	accessConditions := ModifiedAccessConditions{
//		IfModifiedSince: &currentTime,
//	}
//	options := StartCopyBlobOptions{
//		ModifiedAccessConditions: &accessConditions,
//	}
//	destBlobClient, _ := createNewBlockBlob(c, containerClient) // The blob must exist to have a last-modified time
//	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.IsNil)
//
//	_, err = destBlobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopyDestIfModifiedSinceFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	destBlobClient, _ := createNewBlockBlob(c, containerClient)
//	currentTime := getRelativeTimeGMT(10)
//	accessConditions := ModifiedAccessConditions{
//		IfModifiedSince: &currentTime,
//	}
//	options := StartCopyBlobOptions{
//		ModifiedAccessConditions: &accessConditions,
//	}
//	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopyDestIfUnmodifiedSinceTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	destBlobClient, _ := createNewBlockBlob(c, containerClient)
//	currentTime := getRelativeTimeGMT(10)
//	accessConditions := ModifiedAccessConditions{
//		IfUnmodifiedSince: &currentTime,
//	}
//	options := StartCopyBlobOptions{
//		ModifiedAccessConditions: &accessConditions,
//	}
//	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.IsNil)
//
//	_, err = destBlobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopyDestIfUnmodifiedSinceFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//	destBlobClient, _ := createNewBlockBlob(c, containerClient)
//	accessConditions := ModifiedAccessConditions{
//		IfUnmodifiedSince: &currentTime,
//	}
//	options := StartCopyBlobOptions{
//		ModifiedAccessConditions: &accessConditions,
//	}
//
//	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopyDestIfMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	destBlobClient, _ := createNewBlockBlob(c, containerClient)
//	resp, _ := destBlobClient.GetProperties(ctx, nil)
//
//	accessConditions := ModifiedAccessConditions{
//		IfMatch: resp.ETag,
//	}
//	options := StartCopyBlobOptions{
//		ModifiedAccessConditions: &accessConditions,
//	}
//
//	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.IsNil)
//
//	resp, err = destBlobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopyDestIfMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	destBlobClient, _ := createNewBlockBlob(c, containerClient)
//	resp, _ := destBlobClient.GetProperties(ctx, nil)
//
//	accessConditions := ModifiedAccessConditions{
//		IfMatch: resp.ETag,
//	}
//	options := StartCopyBlobOptions{
//		ModifiedAccessConditions: &accessConditions,
//	}
//	metadata := make(map[string]string)
//	metadata["bla"] = "bla"
//	_, err := destBlobClient.SetMetadata(ctx, metadata, nil)
//	_assert(err, chk.IsNil)
//
//	_, err = destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopyDestIfNoneMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	destBlobClient, _ := createNewBlockBlob(c, containerClient)
//	resp, _ := destBlobClient.GetProperties(ctx, nil)
//
//	accessConditions := ModifiedAccessConditions{
//		IfNoneMatch: resp.ETag,
//	}
//	options := StartCopyBlobOptions{
//		ModifiedAccessConditions: &accessConditions,
//	}
//	_, err := destBlobClient.SetMetadata(ctx, nil, nil) // SetMetadata chances the blob's etag
//	_assert(err, chk.IsNil)
//
//	_, err = destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.IsNil)
//
//	resp, err = destBlobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//}
//
//func (s *aztestsSuite) TestBlobStartCopyDestIfNoneMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	destBlobClient, _ := createNewBlockBlob(c, containerClient)
//	resp, _ := destBlobClient.GetProperties(ctx, nil)
//
//	accessConditions := ModifiedAccessConditions{
//		IfNoneMatch: resp.ETag,
//	}
//	options := StartCopyBlobOptions{
//		ModifiedAccessConditions: &accessConditions,
//	}
//
//	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobAbortCopyInProgress() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := getBlockBlobClient(c, containerClient)
//
//	// Create a large blob that takes time to copy
//	blobSize := 8 * 1024 * 1024
//	blobReader, _ := getRandomDataAndReader(blobSize)
//	_, err := blobClient.Upload(ctx, azcore.NopCloser(blobReader), nil)
//	_assert(err, chk.IsNil)
//
//	access := PublicAccessBlob
//	setAccessPolicyOptions := SetAccessPolicyOptions{
//		ContainerSetAccessPolicyOptions: ContainerSetAccessPolicyOptions{
//			Access: &access,
//		},
//	}
//	_, err = containerClient.SetAccessPolicy(ctx, &setAccessPolicyOptions) // So that we don't have to create a SAS
//	_assert(err, chk.IsNil)
//
//	// Must copy across accounts so it takes time to copy
//	bsu2, err := getAlternateBSU()
//	if err != nil {
//		c.Skip(err.Error())
//		return
//	}
//
//	copyContainerClient, _ := createNewContainer(c, bsu2)
//	copyBlobClient, _ := getBlockBlobClient(c, copyContainerClient)
//
//	defer deleteContainer(copyContainerClient)
//
//	resp, err := copyBlobClient.StartCopyFromURL(ctx, blobClient.URL(), nil)
//	_assert(err, chk.IsNil)
//	_assert(resp.CopyStatus, chk.Equals, CopyStatusPending)
//	_assert(resp.CopyID, chk.NotNil)
//
//	_, err = copyBlobClient.AbortCopyFromURL(ctx, *resp.CopyID, nil)
//	if err != nil {
//		// If the error is nil, the test continues as normal.
//		// If the error is not nil, we want to check if it's because the copy is finished and send a message indicating this.
//		validateStorageError(c, err, StorageErrorCodeNoPendingCopyOperation)
//		c.Error("The test failed because the copy completed because it was aborted")
//	}
//
//	resp2, _ := copyBlobClient.GetProperties(ctx, nil)
//	_assert(resp2.CopyStatus, chk.Equals, CopyStatusAborted)
//}
//
//func (s *aztestsSuite) TestBlobAbortCopyNoCopyStarted() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//
//	defer deleteContainer(containerClient)
//
//	copyBlobClient, _ := getBlockBlobClient(c, containerClient)
//	_, err := copyBlobClient.AbortCopyFromURL(ctx, "copynotstarted", nil)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotMetadata() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
//		Metadata: &basicMetadata,
//	}
//	resp, err := blobClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
//	_assert(err, chk.IsNil)
//	_assert(resp.Snapshot, chk.NotNil)
//
//	// Since metadata is specified on the snapshot, the snapshot should have its own metadata different from the (empty) metadata on the source
//	snapshotURL := blobClient.WithSnapshot(*resp.Snapshot)
//	resp2, err := snapshotURL.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(resp2.Metadata, chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotMetadataEmpty() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	_assert(err, chk.IsNil)
//
//	resp, err := blobClient.CreateSnapshot(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(resp.Snapshot, chk.NotNil)
//
//	// In this case, because no metadata was specified, it should copy the basicMetadata from the source
//	snapshotURL := blobClient.WithSnapshot(*resp.Snapshot)
//	resp2, err := snapshotURL.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(resp2.Metadata, chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotMetadataNil() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	_assert(err, chk.IsNil)
//
//	resp, err := blobClient.CreateSnapshot(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(resp.Snapshot, chk.NotNil)
//
//	snapshotURL := blobClient.WithSnapshot(*resp.Snapshot)
//	resp2, err := snapshotURL.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(resp2.Metadata, chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotMetadataInvalid() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	createBlobSnapshotOptions := CreateBlobSnapshotOptions{
//		Metadata: &map[string]string{"Invalid Field!": "value"},
//	}
//	_, err := blobClient.CreateSnapshot(ctx, &createBlobSnapshotOptions)
//	_assert(err, chk.NotNil)
//	_assert(strings.Contains(err.Error(), invalidHeaderErrorSubstring), chk.Equals, true)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotBlobNotExist() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := getBlockBlobClient(c, containerClient)
//
//	_, err := blobClient.CreateSnapshot(ctx, nil)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotOfSnapshot() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	snapshotURL := blobClient.WithSnapshot(time.Now().UTC().Format(SnapshotTimeFormat))
//	// The library allows the server to handle the snapshot of snapshot error
//	_, err := snapshotURL.CreateSnapshot(ctx, nil)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	access := ModifiedAccessConditions{
//		IfModifiedSince: &currentTime,
//	}
//	options := CreateBlobSnapshotOptions{
//		ModifiedAccessConditions: &access,
//	}
//	resp, err := blobClient.CreateSnapshot(ctx, &options)
//	_assert(err, chk.IsNil)
//	_assert(*resp.Snapshot != "", chk.Equals, true) // i.e. The snapshot time is not zero. If the service gives us back a snapshot time, it successfully created a snapshot
//}
//
//func (s *aztestsSuite) TestBlobSnapshotIfModifiedSinceFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//	access := ModifiedAccessConditions{
//		IfModifiedSince: &currentTime,
//	}
//	options := CreateBlobSnapshotOptions{
//		ModifiedAccessConditions: &access,
//	}
//	_, err := blobClient.CreateSnapshot(ctx, &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotIfUnmodifiedSinceTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//	access := ModifiedAccessConditions{
//		IfUnmodifiedSince: &currentTime,
//	}
//	options := CreateBlobSnapshotOptions{
//		ModifiedAccessConditions: &access,
//	}
//	resp, err := blobClient.CreateSnapshot(ctx, &options)
//	_assert(err, chk.IsNil)
//	_assert(*resp.Snapshot == "", chk.Equals, false)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotIfUnmodifiedSinceFalse() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//	access := ModifiedAccessConditions{
//		IfUnmodifiedSince: &currentTime,
//	}
//	options := CreateBlobSnapshotOptions{
//		ModifiedAccessConditions: &access,
//	}
//	_, err := blobClient.CreateSnapshot(ctx, &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotIfMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//
//	access := ModifiedAccessConditions{
//		IfMatch: resp.ETag,
//	}
//	options := CreateBlobSnapshotOptions{
//		ModifiedAccessConditions: &access,
//	}
//	resp2, err := blobClient.CreateSnapshot(ctx, &options)
//	_assert(err, chk.IsNil)
//	_assert(*resp2.Snapshot == "", chk.Equals, false)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotIfMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	randomEtag := "garbage"
//	access := ModifiedAccessConditions{
//		IfMatch: &randomEtag,
//	}
//	options := CreateBlobSnapshotOptions{
//		ModifiedAccessConditions: &access,
//	}
//	_, err := blobClient.CreateSnapshot(ctx, &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotIfNoneMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	randomEtag := "garbage"
//	access := ModifiedAccessConditions{
//		IfNoneMatch: &randomEtag,
//	}
//	options := CreateBlobSnapshotOptions{
//		ModifiedAccessConditions: &access,
//	}
//	resp, err := blobClient.CreateSnapshot(ctx, &options)
//	_assert(err, chk.IsNil)
//	_assert(*resp.Snapshot == "", chk.Equals, false)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotIfNoneMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//
//	access := ModifiedAccessConditions{
//		IfNoneMatch: resp.ETag,
//	}
//	options := CreateBlobSnapshotOptions{
//		ModifiedAccessConditions: &access,
//	}
//	_, err = blobClient.CreateSnapshot(ctx, &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataNonExistentBlob() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := getBlockBlobClient(c, containerClient)
//
//	_, err := blobClient.Download(ctx, nil)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataNegativeOffset() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	offset := int64(-1)
//	options := DownloadBlobOptions{
//		Offset: &offset,
//	}
//	_, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.IsNil)
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataOffsetOutOfRange() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	offset := int64(len(blockBlobDefaultData))
//	options := DownloadBlobOptions{
//		Offset: &offset,
//	}
//	_, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataCountNegative() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	count := int64(-2)
//	options := DownloadBlobOptions{
//		Count: &count,
//	}
//	_, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.IsNil)
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataCountZero() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	count := int64(0)
//	options := DownloadBlobOptions{
//		Count: &count,
//	}
//	resp, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.IsNil)
//
//	// Specifying a count of 0 results in the value being ignored
//	data, err := ioutil.ReadAll(resp.RawResponse.Body)
//	_assert(err, chk.IsNil)
//	_assert(string(data), chk.Equals, blockBlobDefaultData)
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataCountExact() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	count := int64(len(blockBlobDefaultData))
//	options := DownloadBlobOptions{
//		Count: &count,
//	}
//	resp, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.IsNil)
//
//	data, err := ioutil.ReadAll(resp.RawResponse.Body)
//	_assert(err, chk.IsNil)
//	_assert(string(data), chk.Equals, blockBlobDefaultData)
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataCountOutOfRange() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	count := int64(len(blockBlobDefaultData)) * 2
//	options := DownloadBlobOptions{
//		Count: &count,
//	}
//	resp, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.IsNil)
//
//	data, err := ioutil.ReadAll(resp.RawResponse.Body)
//	_assert(err, chk.IsNil)
//	_assert(string(data), chk.Equals, blockBlobDefaultData)
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataEmptyRangeStruct() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	count := int64(0)
//	offset := int64(0)
//	options := DownloadBlobOptions{
//		Count:  &count,
//		Offset: &offset,
//	}
//	resp, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.IsNil)
//
//	data, err := ioutil.ReadAll(resp.RawResponse.Body)
//	_assert(err, chk.IsNil)
//	_assert(string(data), chk.Equals, blockBlobDefaultData)
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataContentMD5() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	offset := int64(10)
//	count := int64(3)
//	getMD5 := true
//	options := DownloadBlobOptions{
//		Count:              &count,
//		Offset:             &offset,
//		RangeGetContentMD5: &getMD5,
//	}
//	resp, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.IsNil)
//	mdf := md5.Sum([]byte(blockBlobDefaultData)[10:13])
//	_assert(*resp.ContentMD5, chk.DeepEquals, mdf[:])
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	access := ModifiedAccessConditions{
//		IfModifiedSince: &currentTime,
//	}
//	options := DownloadBlobOptions{
//		ModifiedAccessConditions: &access,
//	}
//	resp, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.IsNil)
//	_assert(*resp.ContentLength, chk.Equals, int64(len(blockBlobDefaultData)))
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataIfModifiedSinceFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	access := ModifiedAccessConditions{
//		IfModifiedSince: &currentTime,
//	}
//	options := DownloadBlobOptions{
//		ModifiedAccessConditions: &access,
//	}
//	_, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataIfUnmodifiedSinceTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//	access := ModifiedAccessConditions{
//		IfUnmodifiedSince: &currentTime,
//	}
//	options := DownloadBlobOptions{
//		ModifiedAccessConditions: &access,
//	}
//	resp, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.IsNil)
//	_assert(*resp.ContentLength, chk.Equals, int64(len(blockBlobDefaultData)))
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataIfUnmodifiedSinceFalse() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//	access := ModifiedAccessConditions{
//		IfUnmodifiedSince: &currentTime,
//	}
//	options := DownloadBlobOptions{
//		ModifiedAccessConditions: &access,
//	}
//	_, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataIfMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	access := ModifiedAccessConditions{
//		IfMatch: resp.ETag,
//	}
//	options := DownloadBlobOptions{
//		ModifiedAccessConditions: &access,
//	}
//
//	resp2, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.IsNil)
//	_assert(*resp2.ContentLength, chk.Equals, int64(len(blockBlobDefaultData)))
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataIfMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	access := ModifiedAccessConditions{
//		IfMatch: resp.ETag,
//	}
//	options := DownloadBlobOptions{
//		ModifiedAccessConditions: &access,
//	}
//
//	_, err = blobClient.SetMetadata(ctx, nil, nil)
//	_assert(err, chk.IsNil)
//
//	_, err = blobClient.Download(ctx, &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataIfNoneMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	access := ModifiedAccessConditions{
//		IfNoneMatch: resp.ETag,
//	}
//	options := DownloadBlobOptions{
//		ModifiedAccessConditions: &access,
//	}
//
//	_, err = blobClient.SetMetadata(ctx, nil, nil)
//	_assert(err, chk.IsNil)
//
//	resp2, err := blobClient.Download(ctx, &options)
//	_assert(err, chk.IsNil)
//	_assert(*resp2.ContentLength, chk.Equals, int64(len(blockBlobDefaultData)))
//}
//
//func (s *aztestsSuite) TestBlobDownloadDataIfNoneMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	access := ModifiedAccessConditions{
//		IfNoneMatch: resp.ETag,
//	}
//	options := DownloadBlobOptions{
//		ModifiedAccessConditions: &access,
//	}
//
//	_, err = blobClient.Download(ctx, &options)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobDeleteNonExistant() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := getBlockBlobClient(c, containerClient)
//
//	_, err := blobClient.Delete(ctx, nil)
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobDeleteSnapshot() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.CreateSnapshot(ctx, nil)
//	_assert(err, chk.IsNil)
//	snapshotURL := blobClient.WithSnapshot(*resp.Snapshot)
//
//	_, err = snapshotURL.Delete(ctx, nil)
//	_assert(err, chk.IsNil)
//
//	validateBlobDeleted(c, snapshotURL.BlobClient)
//}
//
////func (s *aztestsSuite) TestBlobDeleteSnapshotsInclude() {
////	bsu := getBSU()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////	blobClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err := blobClient.CreateSnapshot(ctx, nil)
////	_assert(err, chk.IsNil)
////
////	deleteSnapshots := DeleteSnapshotsOptionInclude
////	_, err = blobClient.Delete(ctx, &DeleteBlobOptions{
////		DeleteSnapshots: &deleteSnapshots,
////	})
////	_assert(err, chk.IsNil)
////
////	include := []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots}
////	containerListBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
////		Include: &include,
////	}
////	blobs, errChan := containerClient.ListBlobsFlatSegment(ctx, 3, 0, &containerListBlobFlatSegmentOptions)
////	_assert(<- errChan, chk.IsNil)
////	_assert(<- blobs, chk.HasLen, 0)
////}
//
////func (s *aztestsSuite) TestBlobDeleteSnapshotsOnly() {
////	bsu := getBSU()
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////	blobClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err := blobClient.CreateSnapshot(ctx, nil)
////	_assert(err, chk.IsNil)
////	deleteSnapshot := DeleteSnapshotsOptionOnly
////	deleteBlobOptions := DeleteBlobOptions{
////		DeleteSnapshots: &deleteSnapshot,
////	}
////	_, err = blobClient.Delete(ctx, &deleteBlobOptions)
////	_assert(err, chk.IsNil)
////
////	include := []ListBlobsIncludeItem{ListBlobsIncludeItemSnapshots}
////	containerListBlobFlatSegmentOptions := ContainerListBlobFlatSegmentOptions{
////		Include: &include,
////	}
////	blobs, errChan := containerClient.ListBlobsFlatSegment(ctx, 3, 0, &containerListBlobFlatSegmentOptions)
////	_assert(<- errChan, chk.IsNil)
////	_assert(blobs, chk.HasLen, 1)
////	_assert(*(<-blobs).Snapshot == "", chk.Equals, true)
////}
//
//func (s *aztestsSuite) TestBlobDeleteSnapshotsNoneWithSnapshots() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.CreateSnapshot(ctx, nil)
//	_assert(err, chk.IsNil)
//	_, err = blobClient.Delete(ctx, nil)
//	_assert(err, chk.NotNil)
//}
//
//func validateBlobDeleted(, blobClient BlobClient) {
//	_, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.NotNil)
//
//	var serr *StorageError
//	_assert(errors.As(err, &serr), chk.Equals, true)
//	_assert(serr.ErrorCode, chk.Equals, StorageErrorCodeBlobNotFound)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	deleteBlobOptions := DeleteBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
//	}
//	_, err := blobClient.Delete(ctx, &deleteBlobOptions)
//	_assert(err, chk.IsNil)
//
//	validateBlobDeleted(c, blobClient.BlobClient)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfModifiedSinceFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	deleteBlobOptions := DeleteBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
//	}
//	_, err := blobClient.Delete(ctx, &deleteBlobOptions)
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfUnmodifiedSinceTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	deleteBlobOptions := DeleteBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
//	}
//	_, err := blobClient.Delete(ctx, &deleteBlobOptions)
//	_assert(err, chk.IsNil)
//
//	validateBlobDeleted(c, blobClient.BlobClient)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfUnmodifiedSinceFalse() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	deleteBlobOptions := DeleteBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
//	}
//	_, err := blobClient.Delete(ctx, &deleteBlobOptions)
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, _ := blobClient.GetProperties(ctx, nil)
//
//	deleteBlobOptions := DeleteBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
//	}
//	_, err := blobClient.Delete(ctx, &deleteBlobOptions)
//	_assert(err, chk.IsNil)
//
//	validateBlobDeleted(c, blobClient.BlobClient)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	etag := resp.ETag
//
//	_, err = blobClient.SetMetadata(ctx, nil, nil)
//	_assert(err, chk.IsNil)
//
//	deleteBlobOptions := DeleteBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: etag},
//	}
//	_, err = blobClient.Delete(ctx, &deleteBlobOptions)
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfNoneMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, _ := blobClient.GetProperties(ctx, nil)
//	etag := resp.ETag
//	_, err := blobClient.SetMetadata(ctx, nil, nil)
//	_assert(err, chk.IsNil)
//
//	deleteBlobOptions := DeleteBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: etag},
//	}
//	_, err = blobClient.Delete(ctx, &deleteBlobOptions)
//	_assert(err, chk.IsNil)
//
//	validateBlobDeleted(c, blobClient.BlobClient)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfNoneMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, _ := blobClient.GetProperties(ctx, nil)
//	etag := resp.ETag
//
//	deleteBlobOptions := DeleteBlobOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: etag},
//	}
//	_, err := blobClient.Delete(ctx, &deleteBlobOptions)
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	_assert(err, chk.IsNil)
//
//	getBlobPropertiesOptions := GetBlobPropertiesOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
//	}
//	resp, err := blobClient.GetProperties(ctx, &getBlobPropertiesOptions)
//	_assert(err, chk.IsNil)
//	_assert(resp.Metadata, chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfModifiedSinceFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	_assert(err, chk.IsNil)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	getBlobPropertiesOptions := GetBlobPropertiesOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
//	}
//	_, err = blobClient.GetProperties(ctx, &getBlobPropertiesOptions)
//	_assert(err, chk.NotNil)
//	var serr *StorageError
//	_assert(errors.As(err, &serr), chk.Equals, true)
//	_assert(serr.response.StatusCode, chk.Equals, 304) // No service code returned for a HEAD
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfUnmodifiedSinceTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	_assert(err, chk.IsNil)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	getBlobPropertiesOptions := GetBlobPropertiesOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
//	}
//	resp, err := blobClient.GetProperties(ctx, &getBlobPropertiesOptions)
//	_assert(err, chk.IsNil)
//	_assert(resp.Metadata, chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfUnmodifiedSinceFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	_assert(err, chk.IsNil)
//
//	getBlobPropertiesOptions := GetBlobPropertiesOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
//	}
//	_, err = blobClient.GetProperties(ctx, &getBlobPropertiesOptions)
//	_assert(err, chk.NotNil)
//	var serr *StorageError
//	_assert(errors.As(err, &serr), chk.Equals, true)
//	_assert(serr.response.StatusCode, chk.Equals, 412)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	_assert(err, chk.IsNil)
//
//	getBlobPropertiesOptions := GetBlobPropertiesOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
//	}
//	resp2, err := blobClient.GetProperties(ctx, &getBlobPropertiesOptions)
//	_assert(err, chk.IsNil)
//	_assert(resp2.Metadata, chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsOnMissingBlob() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient := containerClient.NewBlobClient("MISSING")
//
//	_, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.NotNil)
//	var storageError *StorageError
//	_assert(errors.As(err, &storageError), chk.Equals, true)
//	_assert(storageError.response.StatusCode, chk.Equals, 404)
//	//_assert(storageError.ErrorCode(), chk.Equals, StorageErrorCodeBlobNotFound)
//	_assert(storageError.response.Header.Get("x-ms-error-code"), chk.Equals, string(StorageErrorCodeBlobNotFound))
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	eTag := "garbage"
//	getBlobPropertiesOptions := GetBlobPropertiesOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: &eTag},
//	}
//	_, err := blobClient.GetProperties(ctx, &getBlobPropertiesOptions)
//	_assert(err, chk.NotNil)
//	var serr *StorageError
//	_assert(errors.As(err, &serr), chk.Equals, true)
//	_assert(serr.response.StatusCode, chk.Equals, 412)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfNoneMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	_assert(err, chk.IsNil)
//
//	eTag := "garbage"
//	getBlobPropertiesOptions := GetBlobPropertiesOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: &eTag},
//	}
//	resp, err := blobClient.GetProperties(ctx, &getBlobPropertiesOptions)
//	_assert(err, chk.IsNil)
//	_assert(resp.Metadata, chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfNoneMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.SetMetadata(ctx, nil, nil)
//	_assert(err, chk.IsNil)
//
//	getBlobPropertiesOptions := GetBlobPropertiesOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag},
//	}
//	_, err = blobClient.GetProperties(ctx, &getBlobPropertiesOptions)
//	_assert(err, chk.NotNil)
//	var serr *StorageError
//	_assert(errors.As(err, &serr), chk.Equals, true)
//	_assert(serr.response.StatusCode, chk.Equals, 304)
//}
//
//func (s *aztestsSuite) TestBlobSetPropertiesBasic() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	contentType := "my_type"
//	contentDisposition := "my_disposition"
//	cacheControl := "control"
//	contentLanguage := "my_language"
//	contentEncoding := "my_encoding"
//	headers := BlobHTTPHeaders{
//		BlobContentType:        &contentType,
//		BlobContentDisposition: &contentDisposition,
//		BlobCacheControl:       &cacheControl,
//		BlobContentLanguage:    &contentLanguage,
//		BlobContentEncoding:    &contentEncoding}
//	_, err := blobClient.SetHTTPHeaders(ctx, headers, nil)
//	_assert(err, chk.IsNil)
//
//	resp, _ := blobClient.GetProperties(ctx, nil)
//	h := resp.NewHTTPHeaders()
//	_assert(h, chk.DeepEquals, headers)
//}
//
//func (s *aztestsSuite) TestBlobSetPropertiesEmptyValue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	contentType := "my_type"
//	_, err := blobClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentType: &contentType}, nil)
//	_assert(err, chk.IsNil)
//
//	_, err = blobClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{}, nil)
//	_assert(err, chk.IsNil)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(resp.ContentType, chk.IsNil)
//}
//
//func validatePropertiesSet(, blobClient BlockBlobClient, disposition string) {
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(*resp.ContentDisposition, chk.Equals, disposition)
//}
//
//func (s *aztestsSuite) TestBlobSetPropertiesIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
//		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}})
//	_assert(err, chk.IsNil)
//
//	validatePropertiesSet(c, blobClient, "my_disposition")
//}
//
//func (s *aztestsSuite) TestBlobSetPropertiesIfModifiedSinceFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	_, err := blobClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
//		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}})
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobSetPropertiesIfUnmodifiedSinceTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	_, err := blobClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
//		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}})
//	_assert(err, chk.IsNil)
//
//	validatePropertiesSet(c, blobClient, "my_disposition")
//}
//
//func (s *aztestsSuite) TestBlobSetPropertiesIfUnmodifiedSinceFalse() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
//		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}})
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobSetPropertiesIfMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//
//	_, err = blobClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
//		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag}})
//	_assert(err, chk.IsNil)
//
//	validatePropertiesSet(c, blobClient, "my_disposition")
//}
//
//func (s *aztestsSuite) TestBlobSetPropertiesIfMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
//		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: to.StringPtr("garbage")}})
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobSetPropertiesIfNoneMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
//		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: to.StringPtr("garbage")}})
//	_assert(err, chk.IsNil)
//
//	validatePropertiesSet(c, blobClient, "my_disposition")
//}
//
//func (s *aztestsSuite) TestBlobSetPropertiesIfNoneMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//
//	_, err = blobClient.SetHTTPHeaders(ctx, BlobHTTPHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
//		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag}})
//	_assert(err, chk.NotNil)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataNil() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, map[string]string{"not": "nil"}, nil)
//	_assert(err, chk.IsNil)
//
//	_, err = blobClient.SetMetadata(ctx, nil, nil)
//	_assert(err, chk.IsNil)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(resp.Metadata, chk.HasLen, 0)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataEmpty() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, map[string]string{"not": "nil"}, nil)
//	_assert(err, chk.IsNil)
//
//	_, err = blobClient.SetMetadata(ctx, map[string]string{}, nil)
//	_assert(err, chk.IsNil)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(resp.Metadata, chk.HasLen, 0)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataInvalidField() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, map[string]string{"Invalid field!": "value"}, nil)
//	_assert(err, chk.NotNil)
//	_assert(strings.Contains(err.Error(), invalidHeaderErrorSubstring), chk.Equals, true)
//}
//
//func validateMetadataSet(, blobClient BlockBlobClient) {
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(resp.Metadata, chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfModifiedSinceTrue() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	setBlobMetadataOptions := SetBlobMetadataOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
//	}
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
//	_assert(err, chk.IsNil)
//
//	validateMetadataSet(c, blobClient)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfModifiedSinceFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	setBlobMetadataOptions := SetBlobMetadataOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime},
//	}
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfUnmodifiedSinceTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	setBlobMetadataOptions := SetBlobMetadataOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
//	}
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
//	_assert(err, chk.IsNil)
//
//	validateMetadataSet(c, blobClient)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfUnmodifiedSinceFalse() {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	setBlobMetadataOptions := SetBlobMetadataOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
//	}
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//
//	setBlobMetadataOptions := SetBlobMetadataOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag},
//	}
//	_, err = blobClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
//	_assert(err, chk.IsNil)
//
//	validateMetadataSet(c, blobClient)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	eTag := "garbage"
//	setBlobMetadataOptions := SetBlobMetadataOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: &eTag},
//	}
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfNoneMatchTrue() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	eTag := "garbage"
//	setBlobMetadataOptions := SetBlobMetadataOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: &eTag},
//	}
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
//	_assert(err, chk.IsNil)
//
//	validateMetadataSet(c, blobClient)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfNoneMatchFalse() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//
//	setBlobMetadataOptions := SetBlobMetadataOptions{
//		ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag},
//	}
//	_, err = blobClient.SetMetadata(ctx, basicMetadata, &setBlobMetadataOptions)
//	validateStorageError(c, err, StorageErrorCodeConditionNotMet)
//}
//
//func testBlobsUndeleteImpl(, bsu ServiceClient) error {
//	//containerURL, _ := createNewContainer(c, bsu)
//	//defer deleteContainer(containerURL)
//	//blobURL, _ := createNewBlockBlob(c, containerURL)
//	//
//	//_, err := blobURL.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
//	//_assert(err, chk.IsNil) // This call will not have errors related to slow update of service properties, so we assert.
//	//
//	//_, err = blobURL.Undelete(ctx)
//	//if err != nil { // We want to give the wrapper method a chance to check if it was an error related to the service properties update.
//	//	return err
//	//}
//	//
//	//resp, err := blobURL.GetProperties(ctx, BlobAccessConditions{})
//	//if err != nil {
//	//	return errors.New(string(err.(StorageError).ErrorCode()))
//	//}
//	//_assert(resp.BlobType(), chk.Equals, BlobBlockBlob) // We could check any property. This is just to double check it was undeleted.
//	return nil
//}
//
//func (s *aztestsSuite) TestBlobsUndelete() {
//	bsu := getBSU(nil)
//
//	code := 404
//	runTestRequiringServiceProperties(c, bsu, string(rune(code)), enableSoftDelete, testBlobsUndeleteImpl, disableSoftDelete)
//}
//
//func setAndCheckBlobTier(, blobClient BlobClient, tier AccessTier) {
//	_, err := blobClient.SetTier(ctx, tier, nil)
//	_assert(err, chk.IsNil)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	_assert(err, chk.IsNil)
//	_assert(*resp.AccessTier, chk.Equals, string(tier))
//}
//
//func (s *aztestsSuite) TestBlobSetTierAllTiers() {
//	bsu := getBSU(nil)
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	setAndCheckBlobTier(c, blobClient.BlobClient, AccessTierHot)
//	setAndCheckBlobTier(c, blobClient.BlobClient, AccessTierCool)
//	setAndCheckBlobTier(c, blobClient.BlobClient, AccessTierArchive)
//
//	bsu, err := getPremiumBSU()
//	if err != nil {
//		c.Skip(err.Error())
//	}
//
//	containerClient, _ = createNewContainer(c, bsu)
//	defer deleteContainer(containerClient)
//	pageblobClient, _ := createNewPageBlob(c, containerClient)
//
//	setAndCheckBlobTier(c, pageblobClient.BlobClient, AccessTierP4)
//	setAndCheckBlobTier(c, pageblobClient.BlobClient, AccessTierP6)
//	setAndCheckBlobTier(c, pageblobClient.BlobClient, AccessTierP10)
//	setAndCheckBlobTier(c, pageblobClient.BlobClient, AccessTierP20)
//	setAndCheckBlobTier(c, pageblobClient.BlobClient, AccessTierP30)
//	setAndCheckBlobTier(c, pageblobClient.BlobClient, AccessTierP40)
//	setAndCheckBlobTier(c, pageblobClient.BlobClient, AccessTierP50)
//}
//
////func (s *aztestsSuite) TestBlobTierInferred() {
////	bsu, err := getPremiumBSU()
////	if err != nil {
////		c.Skip(err.Error())
////	}
////
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////	blobClient, _ := createNewPageBlob(c, containerClient)
////
////	resp, err := blobClient.GetProperties(ctx, nil)
////	_assert(err, chk.IsNil)
////	_assert(resp.AccessTierInferred(), chk.Equals, "true")
////
////	resp2, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert(err, chk.IsNil)
////	_assert(resp2.Segment.BlobItems[0].Properties.AccessTierInferred, chk.NotNil)
////	_assert(resp2.Segment.BlobItems[0].Properties.AccessTier, chk.Not(chk.Equals), "")
////
////	_, err = blobClient.SetTier(ctx, AccessTierP4, LeaseAccessConditions{})
////	_assert(err, chk.IsNil)
////
////	resp, err = blobClient.GetProperties(ctx, nil)
////	_assert(err, chk.IsNil)
////	_assert(resp.AccessTierInferred(), chk.Equals, "")
////
////	resp2, err = containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert(err, chk.IsNil)
////	_assert(resp2.Segment.BlobItems[0].Properties.AccessTierInferred, chk.IsNil) // AccessTierInferred never returned if false
////}
////
////func (s *aztestsSuite) TestBlobArchiveStatus() {
////	bsu, err := getBlobStorageBSU()
////	if err != nil {
////		c.Skip(err.Error())
////	}
////
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////	blobClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err = blobClient.SetTier(ctx, AccessTierArchive, LeaseAccessConditions{})
////	_assert(err, chk.IsNil)
////	_, err = blobClient.SetTier(ctx, AccessTierCool, LeaseAccessConditions{})
////	_assert(err, chk.IsNil)
////
////	resp, err := blobClient.GetProperties(ctx, nil)
////	_assert(err, chk.IsNil)
////	_assert(resp.ArchiveStatus(), chk.Equals, string(ArchiveStatusRehydratePendingToCool))
////
////	resp2, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert(err, chk.IsNil)
////	_assert(resp2.Segment.BlobItems[0].Properties.ArchiveStatus, chk.Equals, ArchiveStatusRehydratePendingToCool)
////
////	// delete first blob
////	_, err = blobClient.Delete(ctx, DeleteSnapshotsOptionNone, nil)
////	_assert(err, chk.IsNil)
////
////	blobClient, _ = createNewBlockBlob(c, containerClient)
////
////	_, err = blobClient.SetTier(ctx, AccessTierArchive, LeaseAccessConditions{})
////	_assert(err, chk.IsNil)
////	_, err = blobClient.SetTier(ctx, AccessTierHot, LeaseAccessConditions{})
////	_assert(err, chk.IsNil)
////
////	resp, err = blobClient.GetProperties(ctx, nil)
////	_assert(err, chk.IsNil)
////	_assert(resp.ArchiveStatus(), chk.Equals, string(ArchiveStatusRehydratePendingToHot))
////
////	resp2, err = containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
////	_assert(err, chk.IsNil)
////	_assert(resp2.Segment.BlobItems[0].Properties.ArchiveStatus, chk.Equals, ArchiveStatusRehydratePendingToHot)
////}
////
////func (s *aztestsSuite) TestBlobTierInvalidValue() {
////	bsu, err := getBlobStorageBSU()
////	if err != nil {
////		c.Skip(err.Error())
////	}
////
////	containerClient, _ := createNewContainer(c, bsu)
////	defer deleteContainer(containerClient)
////	blobClient, _ := createNewBlockBlob(c, containerClient)
////
////	_, err = blobClient.SetTier(ctx, AccessTierType("garbage"), LeaseAccessConditions{})
////	validateStorageError(c, err, StorageErrorCodeInvalidHeaderValue)
////}
////
//func (s *aztestsSuite) TestblobClientPartsSASQueryTimes() {
//	StartTimesInputs := []string{
//		"2020-04-20",
//		"2020-04-20T07:00Z",
//		"2020-04-20T07:15:00Z",
//		"2020-04-20T07:30:00.1234567Z",
//	}
//	StartTimesExpected := []time.Time{
//		time.Date(2020, time.April, 20, 0, 0, 0, 0, time.UTC),
//		time.Date(2020, time.April, 20, 7, 0, 0, 0, time.UTC),
//		time.Date(2020, time.April, 20, 7, 15, 0, 0, time.UTC),
//		time.Date(2020, time.April, 20, 7, 30, 0, 123456700, time.UTC),
//	}
//	ExpiryTimesInputs := []string{
//		"2020-04-21",
//		"2020-04-20T08:00Z",
//		"2020-04-20T08:15:00Z",
//		"2020-04-20T08:30:00.2345678Z",
//	}
//	ExpiryTimesExpected := []time.Time{
//		time.Date(2020, time.April, 21, 0, 0, 0, 0, time.UTC),
//		time.Date(2020, time.April, 20, 8, 0, 0, 0, time.UTC),
//		time.Date(2020, time.April, 20, 8, 15, 0, 0, time.UTC),
//		time.Date(2020, time.April, 20, 8, 30, 0, 234567800, time.UTC),
//	}
//
//	for i := 0; i < len(StartTimesInputs); i++ {
//		urlString :=
//			"https://myaccount.blob.core.windows.net/mycontainer/mydirectory/myfile.txt?" +
//				"se=" + url.QueryEscape(ExpiryTimesInputs[i]) + "&" +
//				"sig=NotASignature&" +
//				"sp=r&" +
//				"spr=https&" +
//				"sr=b&" +
//				"st=" + url.QueryEscape(StartTimesInputs[i]) + "&" +
//				"sv=2019-10-10"
//
//		parts := NewBlobURLParts(urlString)
//		_assert(parts.Scheme, chk.Equals, "https")
//		_assert(parts.Host, chk.Equals, "myaccount.blob.core.windows.net")
//		_assert(parts.ContainerName, chk.Equals, "mycontainer")
//		_assert(parts.BlobName, chk.Equals, "mydirectory/myfile.txt")
//
//		sas := parts.SAS
//		_assert(sas.StartTime(), chk.Equals, StartTimesExpected[i])
//		_assert(sas.ExpiryTime(), chk.Equals, ExpiryTimesExpected[i])
//
//		_assert(parts.URL(), chk.Equals, urlString)
//	}
//}
//
//func (s *aztestsSuite) TestDownloadBlockBlobUnexpectedEOF() {
//	bsu := getBSU(nil)
//	cURL, _ := createNewContainer(c, bsu)
//	defer deleteContainer(cURL)
//	bURL, _ := createNewBlockBlob(c, cURL) // This uploads for us.
//	resp, err := bURL.Download(ctx, nil)
//	_assert(err, chk.IsNil)
//
//	// Verify that we can inject errors first.
//	reader := resp.Body(InjectErrorInRetryReaderOptions(errors.New("unrecoverable error")))
//
//	_, err = ioutil.ReadAll(reader)
//	_assert(err, chk.NotNil)
//	_assert(err.Error(), chk.Equals, "unrecoverable error")
//
//	// Then inject the retryable error.
//	reader = resp.Body(InjectErrorInRetryReaderOptions(io.ErrUnexpectedEOF))
//
//	buf, err := ioutil.ReadAll(reader)
//	_assert(err, chk.IsNil)
//	_assert(buf, chk.DeepEquals, []byte(blockBlobDefaultData))
//}
//
//func InjectErrorInRetryReaderOptions(err error) RetryReaderOptions {
//	return RetryReaderOptions{
//		MaxRetryRequests:       1,
//		doInjectError:          true,
//		doInjectErrorRound:     0,
//		injectedError:          err,
//		NotifyFailedRead:       nil,
//		TreatEarlyCloseAsError: false,
//	}
//}
