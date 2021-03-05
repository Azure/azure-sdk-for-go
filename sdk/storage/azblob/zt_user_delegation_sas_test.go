package azblob

import (
	"bytes"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
	chk "gopkg.in/check.v1"
)

//Creates a container and tests permissions by listing blobs
func (s *aztestsSuite) TestUserDelegationSASContainer(c *chk.C) {
	bsu := getBSU()
	containerClient, containerName := getContainerClient(c, bsu)
	currentTime := time.Now().UTC()
	// Ensuring currTime <= time of sending delegating request request
	keyInfo := KeyInfo{
		Start: to.StringPtr(currentTime.Format(SASTimeFormat)),
		Expiry: to.StringPtr(currentTime.Add(48*time.Hour).Format(SASTimeFormat)),
	}
	time.Sleep(2 * time.Second)

	serviceClient, err := getGenericServiceClientWithOAuth(c, "")
	c.Assert(err, chk.IsNil)

	userDelegationCred, err := serviceClient.GetUserDelegationCredential(ctx, keyInfo)
	if err != nil {
		c.Fatal(err)
	}

	cSAS, err := BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		StartTime:     currentTime.Add(-1 * time.Second),
		ExpiryTime:    currentTime.Add(24 * time.Hour),
		Permissions:   "racwdl",
		ContainerName: containerName,
	}.NewSASQueryParameters(userDelegationCred)

	// Create anonymous pipeline
	//p = azcore.NewPipeline(NewAnonymousCredential(), PipelineOptions{})

	// Create the container
	_, err = containerClient.Create(ctx, nil)
	defer containerClient.Delete(ctx, nil)
	if err != nil {
		c.Fatal(err)
	}

	// Craft a container URL w/ container UDK SAS
	cURL := NewBlobURLParts(containerClient.URL())
	cURL.SAS = cSAS
	cSASURL, err := NewContainerClient(cURL.URL(), NewAnonymousCredential(), nil)

	bblob := cSASURL.NewBlockBlobClient("test")
	_, err = bblob.Upload(ctx, strings.NewReader("hello world!"), nil)
	if err != nil {
		c.Fatal(err)
	}

	resp, err := bblob.Download(ctx, nil)
	data := &bytes.Buffer{}
	body := resp.Body(RetryReaderOptions{})
	if body == nil {
		c.Fatal("download body was nil")
	}
	_, err = data.ReadFrom(body)
	if err != nil {
		c.Fatal(err)
	}
	err = body.Close()
	if err != nil {
		c.Fatal(err)
	}

	c.Assert(data.String(), chk.Equals, "hello world!")
	_, err = bblob.Delete(ctx, nil)
	if err != nil {
		c.Fatal(err)
	}
}

// Creates a blob, takes a snapshot, downloads from snapshot, and deletes from the snapshot w/ the token
func (s *aztestsSuite) TestUserDelegationSASBlob(c *chk.C) {
	// Accumulate prerequisite details to create storage etc.
	serviceClient, err := getGenericServiceClientWithOAuth(c, "")
	c.Assert(err, chk.IsNil)

	containerClient, containerName := getContainerClient(c, serviceClient)
	blobClient, blobName := getBlockBlobClient(c, containerClient)
	currentTime := time.Now().UTC()
	time.Sleep(time.Second)

	// Create container & upload sample data
	_, err = containerClient.Create(ctx, nil)
	defer containerClient.Delete(ctx, nil)
	if err != nil {
		c.Fatal(err)
	}

	// Prepare user delegation key
	keyInfo := KeyInfo{Start: to.StringPtr(currentTime.String()), Expiry: to.StringPtr(currentTime.Add(48*time.Hour).String())}
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
	c.Assert(len(blobURLWithSAS), chk.Not(chk.Equals), 0)

	blobClientWithSAS, err := NewBlockBlobClient(blobURLWithSAS, azcore.AnonymousCredential(), nil)
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
