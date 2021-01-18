package azblob

import (
	"bytes"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	chk "gopkg.in/check.v1"
)

// TODO enable after container is implemented
//Creates a container and tests permissions by listing blobs
//func (s *aztestsSuite) TestUserDelegationSASContainer(c *chk.C) {
//	bsu := getBSU()
//	containerURL, containerName := getContainerClient(c, bsu)
//	currentTime := time.Now().UTC()
//	ocred, err := getOAuthCredential("")
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	// Create pipeline w/ OAuth to handle user delegation key obtaining
//	p := NewPipeline(*ocred, PipelineOptions{})
//
//	bsu = bsu.WithPipeline(p)
//	keyInfo := NewKeyInfo(currentTime, currentTime.Add(48*time.Hour))
//	cudk, err := bsu.GetUserDelegationCredential(ctx, keyInfo, nil, nil)
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	cSAS, err := BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		StartTime:     currentTime,
//		ExpiryTime:    currentTime.Add(24 * time.Hour),
//		Permissions:   "racwdl",
//		ContainerName: containerName,
//	}.NewSASQueryParameters(cudk)
//
//	// Create anonymous pipeline
//	p = NewPipeline(NewAnonymousCredential(), PipelineOptions{})
//
//	// Create the container
//	_, err = containerURL.Create(ctx, Metadata{}, PublicAccessNone)
//	defer containerURL.Delete(ctx, ContainerAccessConditions{})
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	// Craft a container URL w/ container UDK SAS
//	cURL := containerURL.URL()
//	cURL.RawQuery += cSAS.Encode()
//	cSASURL := NewContainerURL(cURL, p)
//
//	bblob := cSASURL.NewBlockBlobURL("test")
//	_, err = bblob.Upload(ctx, strings.NewReader("hello world!"), BlobHTTPHeaders{}, Metadata{}, BlobAccessConditions{})
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	resp, err := bblob.Download(ctx, 0, 0, BlobAccessConditions{}, false)
//	data := &bytes.Buffer{}
//	body := resp.Body(RetryReaderOptions{})
//	if body == nil {
//		c.Fatal("download body was nil")
//	}
//	_, err = data.ReadFrom(body)
//	if err != nil {
//		c.Fatal(err)
//	}
//	err = body.Close()
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	c.Assert(data.String(), chk.Equals, "hello world!")
//	_, err = bblob.Delete(ctx, DeleteSnapshotsOptionNone, BlobAccessConditions{})
//	if err != nil {
//		c.Fatal(err)
//	}
//}

// Creates a blob, takes a snapshot, downloads from snapshot, and deletes from the snapshot w/ the token
func (s *aztestsSuite) TestUserDelegationSASBlob(c *chk.C) {
	// Accumulate prerequisite details to create storage etc.
	serviceClient, err := getGenericServiceClientWithOAuth(c, "")
	c.Assert(err, chk.IsNil)

	containerURL, containerName := getContainerClient(c, serviceClient)
	blobClient, blobName := getBlockBlobClient(c, containerURL)
	currentTime := time.Now().UTC()

	// Create container & upload sample data
	_, err = containerURL.Create(ctx, nil)
	defer containerURL.Delete(ctx, nil)
	if err != nil {
		c.Fatal(err)
	}

	// Prepare user delegation key
	keyInfo := NewKeyInfo(currentTime, currentTime.Add(48*time.Hour))
	cudk, err := serviceClient.GetUserDelegationCredential(ctx, keyInfo)
	c.Assert(err, chk.IsNil)
	c.Assert(cudk, chk.NotNil)

	// Prepare User Delegation SAS query
	bSAS, err := BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		StartTime:     currentTime,
		ExpiryTime:    currentTime.Add(24 * time.Hour),
		Permissions:   "racwd",
		ContainerName: containerName,
		BlobName:      blobName,
	}.NewSASQueryParameters(cudk)
	c.Assert(err, chk.IsNil)

	// Append User Delegation SAS token to URL
	bSASParts := NewBlobURLParts(blobClient.URL())
	bSASParts.SAS = bSAS
	blobURLWithSAS := bSASParts.URL()
	blobRawURL := blobURLWithSAS.String()
	c.Assert(blobRawURL, chk.NotNil)

	blobClientWithSAS, err := NewBlockBlobClient(blobURLWithSAS.String(), azcore.AnonymousCredential(), nil)
	c.Assert(err, chk.IsNil)

	data := "Hello World!"
	_, err = blobClient.Upload(ctx, azcore.NopCloser(strings.NewReader(data)), nil)
	c.Assert(err, chk.IsNil)

	// Download data via User Delegation SAS URL; must succeed
	downloadResponse, err := blobClientWithSAS.Download(ctx, nil)
	c.Assert(err, chk.IsNil)

	downloadedData := &bytes.Buffer{}
	reader := downloadResponse.Body(RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(reader)
	c.Assert(err, chk.IsNil)

	err = reader.Close()
	c.Assert(err, chk.IsNil)
	c.Assert(data, chk.Equals, downloadedData.String())

	// Delete the item using the User Delegation SAS URL; must succeed
	_, err = blobClientWithSAS.Delete(ctx, nil)
	c.Assert(err, chk.IsNil)
}
