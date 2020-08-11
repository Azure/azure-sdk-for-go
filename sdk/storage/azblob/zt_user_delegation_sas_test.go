package azblob

//
//import (
//	"bytes"
//	"strings"
//	"time"
//
//	chk "gopkg.in/check.v1"
//)
//
////Creates a container and tests permissions by listing blobs
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
//
//// Creates a blob, takes a snapshot, downloads from snapshot, and deletes from the snapshot w/ the token
//func (s *aztestsSuite) TestUserDelegationSASBlob(c *chk.C) {
//	// Accumulate prerequisite details to create storage etc.
//	bsu := getBSU()
//	containerURL, containerName := getContainerClient(c, bsu)
//	blobURL, blobName := getBlockBlobURL(c, containerURL)
//	currentTime := time.Now().UTC()
//	ocred, err := getOAuthCredential("")
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	// Create pipeline to handle requests
//	p := NewPipeline(*ocred, PipelineOptions{})
//
//	// Prepare user delegation key
//	bsu = bsu.WithPipeline(p)
//	keyInfo := NewKeyInfo(currentTime, currentTime.Add(48*time.Hour))
//	budk, err := bsu.GetUserDelegationCredential(ctx, keyInfo, nil, nil) //MUST have TokenCredential
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	// Prepare User Delegation SAS query
//	bSAS, err := BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,
//		StartTime:     currentTime,
//		ExpiryTime:    currentTime.Add(24 * time.Hour),
//		Permissions:   "rd",
//		ContainerName: containerName,
//		BlobName:      blobName,
//	}.NewSASQueryParameters(budk)
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	// Create pipeline
//	p = NewPipeline(NewAnonymousCredential(), PipelineOptions{})
//
//	// Append User Delegation SAS token to URL
//	bSASParts := NewBlobURLParts(blobURL.URL())
//	bSASParts.SAS = bSAS
//	bSASURL := NewBlockBlobURL(bSASParts.URL(), p)
//
//	// Create container & upload sample data
//	_, err = containerURL.Create(ctx, Metadata{}, PublicAccessNone)
//	defer containerURL.Delete(ctx, ContainerAccessConditions{})
//	if err != nil {
//		c.Fatal(err)
//	}
//	data := "Hello World!"
//	_, err = blobURL.Upload(ctx, strings.NewReader(data), BlobHTTPHeaders{ContentType: "text/plain"}, Metadata{}, BlobAccessConditions{})
//	if err != nil {
//		c.Fatal(err)
//	}
//
//	// Download data via User Delegation SAS URL; must succeed
//	downloadResponse, err := bSASURL.Download(ctx, 0, 0, BlobAccessConditions{}, false)
//	if err != nil {
//		c.Fatal(err)
//	}
//	downloadedData := &bytes.Buffer{}
//	reader := downloadResponse.Body(RetryReaderOptions{})
//	_, err = downloadedData.ReadFrom(reader)
//	if err != nil {
//		c.Fatal(err)
//	}
//	err = reader.Close()
//	if err != nil {
//		c.Fatal(err)
//	}
//	c.Assert(data, chk.Equals, downloadedData.String())
//
//	// Delete the item using the User Delegation SAS URL; must succeed
//	_, err = bSASURL.Delete(ctx, DeleteSnapshotsOptionInclude, BlobAccessConditions{})
//	if err != nil {
//		c.Fatal(err)
//	}
//}
