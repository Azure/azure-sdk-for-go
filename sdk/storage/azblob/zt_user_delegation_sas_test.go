// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

//
//import (
//	"bytes"
//	"strings"
//	"time"
//
//	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
//	"github.com/Azure/azure-sdk-for-go/sdk/to"
//	chk "gopkg.in/check.v1"
//)
//
////Creates a container and tests permissions by listing blobs
//func (s *azblobTestSuite) TestUserDelegationSASContainer() {
//	bsu := getServiceClient(nil)
//	containerClient, containerName := getContainerClient(bsu)
//	currentTime := time.Now().UTC()
//	time.Sleep(2 * time.Second)
//
//	svcClient, err := getGenericServiceClientWithOAuth(c, "")
//	_assert.Nil(err)
//
//	// Ensuring currTime <= time of sending delegating request request
//	startTime, expiryTime := to.TimePtr(currentTime), to.TimePtr(currentTime.Add(48*time.Hour))
//	userDelegationCred, err := svcClient.GetUserDelegationCredential(ctx, startTime, expiryTime)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	cSAS, err := BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		StartTime:     currentTime.Add(-1 * time.Second),
//		ExpiryTime:    currentTime.Add(24 * time.Hour),
//		Permissions:   "racwdl",
//		ContainerName: containerName,
//	}.NewSASQueryParameters(userDelegationCred)
//
//	// Create anonymous pipeline
//	//p = azcore.NewPipeline(NewAnonymousCredential(), PipelineOptions{})
//
//	// Create the container
//	_, err = containerClient.Create(ctx, nil)
//	defer containerClient.Delete(ctx, nil)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	// Craft a container URL w/ container UDK SAS
//	cURL := NewBlobURLParts(containerClient.URL())
//	cURL.SAS = cSAS
//	cSASURL, err := NewContainerClient(cURL.URL(), azcore.AnonymousCredential(), nil)
//
//	bblob := cSASURL.NewBlockBlobClient("test")
//	_, err = bblob.Upload(ctx, strings.NewReader("hello world!"), nil)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	resp, err := bblob.Download(ctx, nil)
//	data := &bytes.Buffer{}
//	body := resp.Body(RetryReaderOptions{})
//	if body == nil {
//		c.Fatal("download body was nil")
//	}
//	_, err = data.ReadFrom(body)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//	err = body.Close()
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	_assert(data.String(), chk.Equals, "hello world!")
//	_, err = bblob.Delete(ctx, nil)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//}
//
//// Creates a blob, takes a snapshot, downloads from snapshot, and deletes from the snapshot w/ the token
//func (s *azblobTestSuite) TestUserDelegationSASBlob() {
//	// Accumulate prerequisite details to create storage etc.
//	svcClient, err := getGenericServiceClientWithOAuth(c, "")
//	_assert.Nil(err)
//
//	containerClient, containerName := getContainerClient(svcClient)
//	blobClient, blobName := getBlockBlobClient(c, containerClient)
//	currentTime := time.Now().UTC()
//	time.Sleep(time.Second)
//
//	// Create container & upload sample data
//	_, err = containerClient.Create(ctx, nil)
//	defer containerClient.Delete(ctx, nil)
//	if err != nil {
//		s.T().Fatal(err)
//	}
//
//	// Ensuring currTime <= time of sending delegating request request
//	startTime, expiryTime := to.TimePtr(currentTime), to.TimePtr(currentTime.Add(48*time.Hour))
//	cudk, err := svcClient.GetUserDelegationCredential(ctx, startTime, expiryTime)
//	_assert.Nil(err)
//	_assert(cudk, chk.NotNil)
//
//	// Prepare User Delegation SAS query
//	bSAS, err := BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		StartTime:     currentTime,
//		ExpiryTime:    currentTime.Add(24 * time.Hour),
//		Permissions:   "racwd",
//		ContainerName: containerName,
//		BlobName:      blobName,
//	}.NewSASQueryParameters(cudk)
//	_assert.Nil(err)
//
//	// Append User Delegation SAS token to URL
//	bSASParts := NewBlobURLParts(blobClient.URL())
//	bSASParts.SAS = bSAS
//	blobURLWithSAS := bSASParts.URL()
//	_assert(len(blobURLWithSAS), chk.Not(chk.Equals), 0)
//
//	blobClientWithSAS, err := NewBlockBlobClient(blobURLWithSAS, azcore.AnonymousCredential(), nil)
//	_assert.Nil(err)
//
//	data := "Hello World!"
//	_, err = blobClient.Upload(ctx, azcore.NopCloser(strings.NewReader(data)), nil)
//	_assert.Nil(err)
//
//	// Download data via User Delegation SAS URL; must succeed
//	downloadResponse, err := blobClientWithSAS.Download(ctx, nil)
//	_assert.Nil(err)
//
//	downloadedData := &bytes.Buffer{}
//	reader := downloadResponse.Body(RetryReaderOptions{})
//	_, err = downloadedData.ReadFrom(reader)
//	_assert.Nil(err)
//
//	err = reader.Close()
//	_assert.Nil(err)
//	_assert(data, chk.Equals, downloadedData.String())
//
//	// Delete the item using the User Delegation SAS URL; must succeed
//	_, err = blobClientWithSAS.Delete(ctx, nil)
//	_assert.Nil(err)
//}
