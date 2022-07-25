//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

//
////
////import (
////	"bytes"
////	"strings"
////	"time"
////
////	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
////	chk "gopkg.in/check.v1"
////)
////
////func (s *azblobTestSuite) TestSnapshotSAS() {
////	//Generate URLs ----------------------------------------------------------------------------------------------------
////	bsu := getServiceClient(nil)
////	containerClient, containerName := getContainerClient(bsu)
////	blobURL, blobName := getBlockBlobClient(c, containerClient)
////
////	_, err := containerClient.Create(ctx, nil)
////	defer containerClient.Delete(ctx, nil)
////	if err != nil {
////		s.T().Fatal(err)
////	}
////
////	//Create file in container, download from snapshot to test. --------------------------------------------------------
////	blobClient := containerClient.NewBlockBlobClient(blobName)
////	data := "Hello world!"
////
////	contentType := "text/plain"
////	uploadBlockBlobOptions := BlockBlobUploadOptions{
////		HTTPHeaders: &HTTPHeaders{
////			BlobContentType: &contentType,
////		},
////	}
////	_, err = blobClient.Upload(ctx, strings.NewReader(data), &uploadBlockBlobOptions)
////	if err != nil {
////		s.T().Fatal(err)
////	}
////
////	//Create a snapshot & URL
////	createSnapshot, err := blobClient.CreateSnapshot(ctx, nil)
////	if err != nil {
////		s.T().Fatal(err)
////	}
////	_assert(createSnapshot.Snapshot, chk.NotNil)
////
////	//Format snapshot time
////	snapTime, err := time.Parse(SnapshotTimeFormat, *createSnapshot.Snapshot)
////	if err != nil {
////		s.T().Fatal(err)
////	}
////
////	//Get credentials & current time
////	currentTime := time.Now().UTC()
////	credential, err := getGenericCredential("")
////	if err != nil {
////		c.Fatal("Invalid credential")
////	}
////
////	//Create SAS query
////	snapSASQueryParams, err := BlobSASSignatureValues{
////		StartTime:     currentTime,
////		ExpiryTime:    currentTime.Add(48 * time.Hour),
////		SnapshotTime:  snapTime,
////		Permissions:   "racwd",
////		ContainerName: containerName,
////		BlobName:      blobName,
////		Protocol:      SASProtocolHTTPS,
////	}.NewSASQueryParameters(credential)
////	if err != nil {
////		s.T().Fatal(err)
////	}
////	time.Sleep(time.Second * 2)
////
////	//Attach SAS query to block blob URL
////	snapParts := NewBlobURLParts(blobURL.URL())
////	snapParts.SAS = snapSASQueryParams
////	sbUrl, err := NewBlockBlobClient(snapParts.URL(), azcore.AnonymousCredential(), nil)
////
////	//Test the snapshot
////	downloadResponse, err := sbUrl.Download(ctx, nil)
////	if err != nil {
////		s.T().Fatal(err)
////	}
////
////	downloadedData := &bytes.Buffer{}
////	reader := downloadResponse.Body(RetryReaderOptions{})
////	downloadedData.ReadFrom(reader)
////	reader.Close()
////
////	_assert(data, chk.Equals, downloadedData.String())
////
////	//Try to delete snapshot -------------------------------------------------------------------------------------------
////	_, err = sbUrl.Delete(ctx, nil)
////	if err != nil { //This shouldn't fail.
////		s.T().Fatal(err)
////	}
////
////	//Create a normal blob and attempt to use the snapshot SAS against it (assuming failure) ---------------------------
////	//If this succeeds, it means a normal SAS token was created.
////
////	uploadBlockBlobOptions1 := BlockBlobUploadOptions{
////		HTTPHeaders: &HTTPHeaders{
////			BlobContentType: &contentType,
////		},
////	}
////	fsbUrl := containerClient.NewBlockBlobClient("failsnap")
////	_, err = fsbUrl.Upload(ctx, strings.NewReader(data), &uploadBlockBlobOptions1)
////	if err != nil {
////		s.T().Fatal(err) //should succeed to create the blob via normal auth means
////	}
////
////	fsbUrlParts := NewBlobURLParts(fsbUrl.URL())
////	fsbUrlParts.SAS = snapSASQueryParams
////	fsbUrl, err = NewBlockBlobClient(fsbUrlParts.URL(), azcore.AnonymousCredential(), nil) //re-use fsbUrl as we don't need the sharedkey version anymore
////
////	resp, err := fsbUrl.Delete(ctx, nil)
////	if err == nil {
////		c.Fatal(resp) //This SHOULD fail. Otherwise we have a normal SAS token...
////	}
////}
