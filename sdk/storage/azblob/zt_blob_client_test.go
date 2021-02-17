package azblob

import (
	"bytes"
	"crypto/md5"
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/to"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
)

func (s *aztestsSuite) TestCreateBlobClient(c *chk.C) {
	bsu := getBSU()
	containerClient, containerName := getContainerClient(c, bsu)
	testURL, testName := getBlockBlobClient(c, containerClient)

	parts := NewBlobURLParts(testURL.URL())
	c.Assert(parts.BlobName, chk.Equals, testName)
	c.Assert(parts.ContainerName, chk.Equals, containerName)

	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + ".blob.core.windows.net/" + containerName + "/" + testName
	temp := testURL.URL()
	c.Assert(temp.String(), chk.Equals, correctURL)
}

func (s *aztestsSuite) TestCreateBlobClientWithSnapshotAndSAS(c *chk.C) {
	bsu := getBSU()
	containerClient, containerName := getContainerClient(c, bsu)
	blobClient, blobName := getBlockBlobClient(c, containerClient)

	currentTime := time.Now().UTC()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	sasQueryParams, err := AccountSASSignatureValues{
		Protocol:      SASProtocolHTTPS,
		ExpiryTime:    currentTime.Add(48 * time.Hour),
		Permissions:   AccountSASPermissions{Read: true, List: true}.String(),
		Services:      AccountSASServices{Blob: true}.String(),
		ResourceTypes: AccountSASResourceTypes{Container: true, Object: true}.String(),
	}.NewSASQueryParameters(credential)
	if err != nil {
		c.Fatal(err)
	}

	parts := NewBlobURLParts(blobClient.URL())
	parts.SAS = sasQueryParams
	parts.Snapshot = currentTime.Format(SnapshotTimeFormat)
	testURL := parts.URL()

	// The snapshot format string is taken from the snapshotTimeFormat value in parsing_urls.go. The field is not public, so
	// it is copied here
	correctURL := "https://" + os.Getenv("AZURE_STORAGE_ACCOUNT_NAME") + ".blob.core.windows.net/" + containerName + "/" + blobName +
		"?" + "snapshot=" + currentTime.Format("2006-01-02T15:04:05.0000000Z07:00") + "&" + sasQueryParams.Encode()
	c.Assert(testURL.String(), chk.Equals, correctURL)
}

// func (s *aztestsSuite) TestBlobWithNewPipeline(c *chk.C) {
// 	bsu := getBSU()
// 	containerClient, _ := getContainerClient(c, bsu)
// 	blobClient := containerClient.NewBlockBlobClient(blobPrefix)
//
// 	newBlobClient := blobClient.WithPipeline(newTestPipeline())
//
// 	// exercise the new pipeline
// 	_, err := newBlobClient.GetAccountInfo(ctx)
// 	c.Assert(err, chk.NotNil)
// 	c.Assert(err.Error(), chk.Equals, testPipelineMessage)
// }

func waitForCopy(c *chk.C, copyBlobClient BlockBlobClient, blobCopyResponse BlobStartCopyFromURLResponse) {
	status := *blobCopyResponse.CopyStatus
	// Wait for the copy to finish. If the copy takes longer than a minute, we will fail
	start := time.Now()
	for status != CopyStatusTypeSuccess {
		props, _ := copyBlobClient.GetProperties(ctx, nil)
		status = *props.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			c.Fail()
		}
	}
}

func (s *aztestsSuite) TestBlobStartCopyDestEmpty(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)
	copyblobClient, _ := getBlockBlobClient(c, containerClient)

	blobCopyResponse, err := copyblobClient.StartCopyFromURL(ctx, blobClient.URL(), nil)
	c.Assert(err, chk.IsNil)
	waitForCopy(c, copyblobClient, blobCopyResponse)

	resp, err := copyblobClient.Download(ctx, nil)
	c.Assert(err, chk.IsNil)

	// Read the blob data to verify the copy
	data, err := ioutil.ReadAll(resp.Response().Body)
	c.Assert(*resp.ContentLength(), chk.Equals, int64(len(blockBlobDefaultData)))
	c.Assert(string(data), chk.Equals, blockBlobDefaultData)
	resp.Body(RetryReaderOptions{}).Close()
}

func (s *aztestsSuite) TestBlobStartCopyMetadata(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)
	copyblobClient, _ := getBlockBlobClient(c, containerClient)

	metadata := make(map[string]string)
	metadata["bla"] = "foo"
	options := StartCopyBlobOptions{
		Metadata: &metadata,
	}
	resp, err := copyblobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.IsNil)
	waitForCopy(c, copyblobClient, resp)

	// TODO enabled after generator fix
	//resp2, err := copyblobClient.GetProperties(ctx, nil)
	//c.Assert(err, chk.IsNil)
	//c.Assert(resp2.NewMetadata(), chk.DeepEquals, metadata)
}

func (s *aztestsSuite) TestBlobStartCopyMetadataNil(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)
	copyblobClient, _ := getBlockBlobClient(c, containerClient)

	// Have the destination start with metadata so we ensure the nil metadata passed later takes effect
	_, err := copyblobClient.Upload(ctx, azcore.NopCloser(bytes.NewReader([]byte("data"))), nil)
	c.Assert(err, chk.IsNil)

	resp, err := copyblobClient.StartCopyFromURL(ctx, blobClient.URL(), nil)
	c.Assert(err, chk.IsNil)

	waitForCopy(c, copyblobClient, resp)

	//// TODO enabled after generator fix
	//resp2, err := copyblobClient.GetProperties(ctx, nil)
	//c.Assert(err, chk.IsNil)
	//c.Assert(resp2.NewMetadata(), chk.HasLen, 0)
}

func (s *aztestsSuite) TestBlobStartCopyMetadataEmpty(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)
	copyblobClient, _ := getBlockBlobClient(c, containerClient)

	// Have the destination start with metadata so we ensure the empty metadata passed later takes effect
	_, err := copyblobClient.Upload(ctx, azcore.NopCloser(bytes.NewReader([]byte("data"))), nil)
	c.Assert(err, chk.IsNil)

	metadata := make(map[string]string)
	options := StartCopyBlobOptions{
		Metadata: &metadata,
	}
	resp, err := copyblobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.IsNil)

	waitForCopy(c, copyblobClient, resp)

	// TODO enabled after generator fix
	//resp2, err := copyblobClient.GetProperties(ctx, nil)
	//c.Assert(err, chk.IsNil)
	//c.Assert(resp2.NewMetadata(), chk.HasLen, 0)
}

// TODO enabled after azcore fix, the retry policy is retrying header errors
//func (s *aztestsSuite) TestBlobStartCopyMetadataInvalidField(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//	copyblobClient, _ := getBlockBlobClient(c, containerClient)
//
//	metadata := make(map[string]string)
//	metadata["I nvalid."] = "foo"
//	options := StartCopyBlobOptions{
//		Metadata: &metadata,
//	}
//	_, err := copyblobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
//	c.Assert(err, chk.NotNil)
//	c.Assert(strings.Contains(err.Error(), invalidHeaderErrorSubstring), chk.Equals, true)
//}

func (s *aztestsSuite) TestBlobStartCopySourceNonExistant(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := getBlockBlobClient(c, containerClient)
	copyblobClient, _ := getBlockBlobClient(c, containerClient)

	_, err := copyblobClient.StartCopyFromURL(ctx, blobClient.URL(), nil)
	c.Assert(err, chk.NotNil)
	c.Assert(strings.Contains(err.Error(), "not exist"), chk.Equals, true)
}

// TODO enable when container client is fully implemented
//func (s *aztestsSuite) TestBlobStartCopySourcePrivate(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	_, err := containerClient.SetAccessPolicy(ctx, PublicAccessNone, nil, ContainerAccessConditions{})
//	c.Assert(err, chk.IsNil)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	bsu2, err := getAlternateBSU()
//	if err != nil {
//		c.Skip(err.Error())
//		return
//	}
//	copycontainerClient, _ := createNewContainer(c, bsu2)
//	defer deleteContainer(c, copycontainerClient)
//	copyblobClient, _ := getBlockBlobClient(c, copycontainerClient)
//
//	if bsu.String() == bsu2.String() {
//		c.Skip("Test not valid because primary and secondary accounts are the same")
//	}
//	_, err = copyblobClient.StartCopyFromURL(ctx, blobClient.URL(), nil, ModifiedAccessConditions{}, nil)
//	validateStorageError(c, err, ServiceCodeCannotVerifyCopySource)
//}

// TODO enable when container client is fully implemented
//func (s *aztestsSuite) TestBlobStartCopyUsingSASSrc(c *chk.C) {
//	bsu := getBSU()
//	containerClient, containerName := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	_, err := containerClient.SetAccessPolicy(ctx, PublicAccessNone, nil, ContainerAccessConditions{})
//	c.Assert(err, chk.IsNil)
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
//	sasURL := blobClient.URL()
//	sasURL.RawQuery = queryParams.Encode()
//
//	// Create a new container for the destination
//	bsu2, err := getAlternateBSU()
//	if err != nil {
//		c.Skip(err.Error())
//		return
//	}
//	copycontainerClient, _ := createNewContainer(c, bsu2)
//	defer deleteContainer(c, copycontainerClient)
//	copyblobClient, _ := getBlockBlobClient(c, copycontainerClient)
//
//	resp, err := copyblobClient.StartCopyFromURL(ctx, sasURL, nil, ModifiedAccessConditions{}, nil)
//	c.Assert(err, chk.IsNil)
//
//	waitForCopy(c, copyblobClient, resp)
//
//	resp2, err := copyblobClient.Download(ctx, 0, int64(len(blockBlobDefaultData)), nil, false)
//	c.Assert(err, chk.IsNil)
//
//	data, err := ioutil.ReadAll(resp2.Response().Body)
//	c.Assert(resp2.ContentLength(), chk.Equals, int64(len(blockBlobDefaultData)))
//	c.Assert(string(data), chk.Equals, blockBlobDefaultData)
//	resp2.Body(RetryReaderOptions{}).Close()
//}

// TODO enable when container client is fully implemented
//func (s *aztestsSuite) TestBlobStartCopyUsingSASDest(c *chk.C) {
//	bsu := getBSU()
//	containerClient, containerName := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	_, err := containerClient.SetAccessPolicy(ctx, PublicAccessNone, nil, ContainerAccessConditions{})
//	c.Assert(err, chk.IsNil)
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
//	copycontainerClient, copyContainerName := createNewContainer(c, bsu2)
//	defer deleteContainer(c, copycontainerClient)
//	copyblobClient, copyBlobName := getBlockBlobClient(c, copycontainerClient)
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
//	anonURL := bsu2.URL()
//	anonURL.RawQuery = copyQueryParams.Encode()
//	anonPipeline := NewPipeline(NewAnonymousCredential(), PipelineOptions{})
//	anonBSU := NewServiceURL(anonURL, anonPipeline)
//	anoncontainerClient := anonBSU.NewcontainerClient(copyContainerName)
//	anonblobClient := anoncontainerClient.NewBlockBlobClient(copyBlobName)
//
//	// Apply sas to source
//	srcBlobWithSasURL := blobClient.URL()
//	srcBlobWithSasURL.RawQuery = queryParams.Encode()
//
//	resp, err := anonblobClient.StartCopyFromURL(ctx, srcBlobWithSasURL, nil, ModifiedAccessConditions{}, nil)
//	c.Assert(err, chk.IsNil)
//
//	// Allow copy to happen
//	waitForCopy(c, anonblobClient, resp)
//
//	resp2, err := copyblobClient.Download(ctx, 0, int64(len(blockBlobDefaultData)), nil, false)
//	c.Assert(err, chk.IsNil)
//
//	data, err := ioutil.ReadAll(resp2.Response().Body)
//	_, err = resp2.Body(RetryReaderOptions{}).Read(data)
//	c.Assert(resp2.ContentLength(), chk.Equals, int64(len(blockBlobDefaultData)))
//	c.Assert(string(data), chk.Equals, blockBlobDefaultData)
//	resp2.Body(RetryReaderOptions{}).Close()
//}

func (s *aztestsSuite) TestBlobStartCopySourceIfModifiedSinceTrue(c *chk.C) {
	currentTime := getRelativeTimeGMT(-10)

	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	accessConditions := SourceModifiedAccessConditions{
		SourceIfModifiedSince: &currentTime,
	}
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &accessConditions,
	}
	destBlobClient, _ := getBlockBlobClient(c, containerClient)
	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.IsNil)

	_, err = destBlobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobStartCopySourceIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(100)
	accessConditions := SourceModifiedAccessConditions{
		SourceIfModifiedSince: &currentTime,
	}
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &accessConditions,
	}
	destBlobClient, _ := getBlockBlobClient(c, containerClient)
	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobStartCopySourceIfUnmodifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)
	accessConditions := SourceModifiedAccessConditions{
		SourceIfUnmodifiedSince: &currentTime,
	}
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &accessConditions,
	}
	destBlobClient, _ := getBlockBlobClient(c, containerClient)
	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.IsNil)

	_, err = destBlobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobStartCopySourceIfUnmodifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)
	accessConditions := SourceModifiedAccessConditions{
		SourceIfUnmodifiedSince: &currentTime,
	}
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &accessConditions,
	}
	destBlobClient, _ := getBlockBlobClient(c, containerClient)
	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobStartCopySourceIfMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	resp, err := blobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)

	accessConditions := SourceModifiedAccessConditions{
		SourceIfMatch: resp.ETag,
	}
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &accessConditions,
	}
	destBlobClient, _ := getBlockBlobClient(c, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.IsNil)

	_, err = destBlobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobStartCopySourceIfMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	randomEtag := "a"
	accessConditions := SourceModifiedAccessConditions{
		SourceIfMatch: &randomEtag,
	}
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &accessConditions,
	}

	destBlobClient, _ := getBlockBlobClient(c, containerClient)
	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobStartCopySourceIfNoneMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	randomEtag := "a"
	accessConditions := SourceModifiedAccessConditions{
		SourceIfNoneMatch: &randomEtag,
	}
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &accessConditions,
	}
	destBlobClient, _ := getBlockBlobClient(c, containerClient)
	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.IsNil)

	_, err = destBlobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobStartCopySourceIfNoneMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	resp, err := blobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)

	accessConditions := SourceModifiedAccessConditions{
		SourceIfNoneMatch: resp.ETag,
	}
	options := StartCopyBlobOptions{
		SourceModifiedAccessConditions: &accessConditions,
	}

	destBlobClient, _ := getBlockBlobClient(c, containerClient)
	_, err = destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobStartCopyDestIfModifiedSinceTrue(c *chk.C) {
	currentTime := getRelativeTimeGMT(-10)
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	accessConditions := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &accessConditions,
	}
	destBlobClient, _ := createNewBlockBlob(c, containerClient) // The blob must exist to have a last-modified time
	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.IsNil)

	_, err = destBlobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobStartCopyDestIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	destBlobClient, _ := createNewBlockBlob(c, containerClient)
	currentTime := getRelativeTimeGMT(10)
	accessConditions := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &accessConditions,
	}
	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobStartCopyDestIfUnmodifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	destBlobClient, _ := createNewBlockBlob(c, containerClient)
	currentTime := getRelativeTimeGMT(10)
	accessConditions := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &accessConditions,
	}
	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.IsNil)

	_, err = destBlobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobStartCopyDestIfUnmodifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)
	destBlobClient, _ := createNewBlockBlob(c, containerClient)
	accessConditions := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &accessConditions,
	}

	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobStartCopyDestIfMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	destBlobClient, _ := createNewBlockBlob(c, containerClient)
	resp, _ := destBlobClient.GetProperties(ctx, nil)

	accessConditions := ModifiedAccessConditions{
		IfMatch: resp.ETag,
	}
	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &accessConditions,
	}

	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.IsNil)

	resp, err = destBlobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobStartCopyDestIfMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	destBlobClient, _ := createNewBlockBlob(c, containerClient)
	resp, _ := destBlobClient.GetProperties(ctx, nil)

	accessConditions := ModifiedAccessConditions{
		IfMatch: resp.ETag,
	}
	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &accessConditions,
	}
	metadata := make(map[string]string)
	metadata["bla"] = "bla"
	_, err := destBlobClient.SetMetadata(ctx, metadata, nil)
	c.Assert(err, chk.IsNil)

	_, err = destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobStartCopyDestIfNoneMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	destBlobClient, _ := createNewBlockBlob(c, containerClient)
	resp, _ := destBlobClient.GetProperties(ctx, nil)

	accessConditions := ModifiedAccessConditions{
		IfNoneMatch: resp.ETag,
	}
	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &accessConditions,
	}
	destBlobClient.SetMetadata(ctx, nil, nil) // SetMetadata chances the blob's etag

	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.IsNil)

	resp, err = destBlobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobStartCopyDestIfNoneMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	destBlobClient, _ := createNewBlockBlob(c, containerClient)
	resp, _ := destBlobClient.GetProperties(ctx, nil)

	accessConditions := ModifiedAccessConditions{
		IfNoneMatch: resp.ETag,
	}
	options := StartCopyBlobOptions{
		ModifiedAccessConditions: &accessConditions,
	}

	_, err := destBlobClient.StartCopyFromURL(ctx, blobClient.URL(), &options)
	c.Assert(err, chk.NotNil)
}

// TODO enabled after SetAccessPolicy
//func (s *aztestsSuite) TestBlobAbortCopyInProgress(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := getBlockBlobClient(c, containerClient)
//
//	// Create a large blob that takes time to copy
//	blobSize := 8 * 1024 * 1024
//	blobReader, blobData := getRandomDataAndReader(blobSize)
//	_, err := blobClient.Upload(ctx, azcore.NopCloser(blobReader), nil)
//	c.Assert(err, chk.IsNil)
//	containerClient.SetAccessPolicy(ctx, PublicAccessBlob, nil, ContainerAccessConditions{}) // So that we don't have to create a SAS
//
//	// Must copy across accounts so it takes time to copy
//	bsu2, err := getAlternateBSU()
//	if err != nil {
//		c.Skip(err.Error())
//		return
//	}
//
//	copycontainerClient, _ := createNewContainer(c, bsu2)
//	copyblobClient, _ := getBlockBlobClient(c, copycontainerClient)
//
//	defer deleteContainer(c, copycontainerClient)
//
//	resp, err := copyblobClient.StartCopyFromURL(ctx, blobClient.URL(), nil, ModifiedAccessConditions{}, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.CopyStatus(), chk.Equals, CopyStatusPending)
//
//	_, err = copyblobClient.AbortCopyFromURL(ctx, resp.CopyID(), LeaseAccessConditions{})
//	if err != nil {
//		// If the error is nil, the test continues as normal.
//		// If the error is not nil, we want to check if it's because the copy is finished and send a message indicating this.
//		validateStorageError(c, err, ServiceCodeNoPendingCopyOperation)
//		c.Error("The test failed because the copy completed because it was aborted")
//	}
//
//	resp2, _ := copyblobClient.GetProperties(ctx, nil)
//	c.Assert(resp2.CopyStatus(), chk.Equals, CopyStatusAborted)
//}

func (s *aztestsSuite) TestBlobAbortCopyNoCopyStarted(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)

	defer deleteContainer(c, containerClient)

	copyblobClient, _ := getBlockBlobClient(c, containerClient)
	_, err := copyblobClient.AbortCopyFromURL(ctx, "copynotstarted", nil)
	c.Assert(err, chk.NotNil)
}

// TODO enabled after metadata fix
//func (s *aztestsSuite) TestBlobSnapshotMetadata(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.CreateSnapshot(ctx, basicMetadata, nil)
//	c.Assert(err, chk.IsNil)
//
//	// Since metadata is specified on the snapshot, the snapshot should have its own metadata different from the (empty) metadata on the source
//	snapshotURL := blobClient.WithSnapshot(resp.Snapshot())
//	resp2, err := snapshotURL.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp2.NewMetadata(), chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotMetadataEmpty(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	c.Assert(err, chk.IsNil)
//
//	resp, err := blobClient.CreateSnapshot(ctx, Metadata{}, nil)
//	c.Assert(err, chk.IsNil)
//
//	// In this case, because no metadata was specified, it should copy the basicMetadata from the source
//	snapshotURL := blobClient.WithSnapshot(resp.Snapshot())
//	resp2, err := snapshotURL.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp2.NewMetadata(), chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotMetadataNil(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	c.Assert(err, chk.IsNil)
//
//	resp, err := blobClient.CreateSnapshot(ctx, nil)
//	c.Assert(err, chk.IsNil)
//
//	snapshotURL := blobClient.WithSnapshot(resp.Snapshot())
//	resp2, err := snapshotURL.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp2.NewMetadata(), chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobSnapshotMetadataInvalid(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.CreateSnapshot(ctx, Metadata{"Invalid Field!": "value"}, nil)
//	c.Assert(err, chk.NotNil)
//	c.Assert(strings.Contains(err.Error(), invalidHeaderErrorSubstring), chk.Equals, true)
//}
//
func (s *aztestsSuite) TestBlobSnapshotBlobNotExist(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := getBlockBlobClient(c, containerClient)

	_, err := blobClient.CreateSnapshot(ctx, nil)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobSnapshotOfSnapshot(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	snapshotURL := blobClient.WithSnapshot(time.Now().UTC().Format(SnapshotTimeFormat))
	// The library allows the server to handle the snapshot of snapshot error
	_, err := snapshotURL.CreateSnapshot(ctx, nil)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobSnapshotIfModifiedSinceTrue(c *chk.C) {
	currentTime := getRelativeTimeGMT(-10)

	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	resp, err := blobClient.CreateSnapshot(ctx, &options)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.Snapshot != "", chk.Equals, true) // i.e. The snapshot time is not zero. If the service gives us back a snapshot time, it successfully created a snapshot
}

func (s *aztestsSuite) TestBlobSnapshotIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)
	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	_, err := blobClient.CreateSnapshot(ctx, &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobSnapshotIfUnmodifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)
	access := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	resp, err := blobClient.CreateSnapshot(ctx, &options)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.Snapshot == "", chk.Equals, false)
}

func (s *aztestsSuite) TestBlobSnapshotIfUnmodifiedSinceFalse(c *chk.C) {
	currentTime := getRelativeTimeGMT(-10)

	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)
	access := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	_, err := blobClient.CreateSnapshot(ctx, &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobSnapshotIfMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	resp, err := blobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)

	access := ModifiedAccessConditions{
		IfMatch: resp.ETag,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	resp2, err := blobClient.CreateSnapshot(ctx, &options)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp2.Snapshot == "", chk.Equals, false)
}

func (s *aztestsSuite) TestBlobSnapshotIfMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	randomEtag := "garbage"
	access := ModifiedAccessConditions{
		IfMatch: &randomEtag,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	_, err := blobClient.CreateSnapshot(ctx, &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobSnapshotIfNoneMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	randomEtag := "garbage"
	access := ModifiedAccessConditions{
		IfNoneMatch: &randomEtag,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	resp, err := blobClient.CreateSnapshot(ctx, &options)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.Snapshot == "", chk.Equals, false)
}

func (s *aztestsSuite) TestBlobSnapshotIfNoneMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	resp, err := blobClient.GetProperties(ctx, nil)
	access := ModifiedAccessConditions{
		IfNoneMatch: resp.ETag,
	}
	options := CreateBlobSnapshotOptions{
		ModifiedAccessConditions: &access,
	}
	_, err = blobClient.CreateSnapshot(ctx, &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobDownloadDataNonExistentBlob(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := getBlockBlobClient(c, containerClient)

	_, err := blobClient.Download(ctx, nil)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobDownloadDataNegativeOffset(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	offset := int64(-1)
	options := DownloadBlobOptions{
		Offset: &offset,
	}
	_, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobDownloadDataOffsetOutOfRange(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	offset := int64(len(blockBlobDefaultData))
	options := DownloadBlobOptions{
		Offset: &offset,
	}
	_, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobDownloadDataCountNegative(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	count := int64(-2)
	options := DownloadBlobOptions{
		Count: &count,
	}
	_, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.IsNil)
}

func (s *aztestsSuite) TestBlobDownloadDataCountZero(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	count := int64(0)
	options := DownloadBlobOptions{
		Count: &count,
	}
	resp, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.IsNil)

	// Specifying a count of 0 results in the value being ignored
	data, err := ioutil.ReadAll(resp.Response().Body)
	c.Assert(err, chk.IsNil)
	c.Assert(string(data), chk.Equals, blockBlobDefaultData)
}

func (s *aztestsSuite) TestBlobDownloadDataCountExact(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	count := int64(len(blockBlobDefaultData))
	options := DownloadBlobOptions{
		Count: &count,
	}
	resp, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.IsNil)

	data, err := ioutil.ReadAll(resp.Response().Body)
	c.Assert(err, chk.IsNil)
	c.Assert(string(data), chk.Equals, blockBlobDefaultData)
}

func (s *aztestsSuite) TestBlobDownloadDataCountOutOfRange(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	count := int64(len(blockBlobDefaultData)) * 2
	options := DownloadBlobOptions{
		Count: &count,
	}
	resp, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.IsNil)

	data, err := ioutil.ReadAll(resp.Response().Body)
	c.Assert(err, chk.IsNil)
	c.Assert(string(data), chk.Equals, blockBlobDefaultData)
}

func (s *aztestsSuite) TestBlobDownloadDataEmptyRangeStruct(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	count := int64(0)
	offset := int64(0)
	options := DownloadBlobOptions{
		Count:  &count,
		Offset: &offset,
	}
	resp, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.IsNil)

	data, err := ioutil.ReadAll(resp.Response().Body)
	c.Assert(err, chk.IsNil)
	c.Assert(string(data), chk.Equals, blockBlobDefaultData)
}

func (s *aztestsSuite) TestBlobDownloadDataContentMD5(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	offset := int64(10)
	count := int64(3)
	getMD5 := true
	options := DownloadBlobOptions{
		Count:              &count,
		Offset:             &offset,
		RangeGetContentMd5: &getMD5,
	}
	resp, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.IsNil)
	mdf := md5.Sum([]byte(blockBlobDefaultData)[10:13])
	c.Assert(*resp.ContentMD5(), chk.DeepEquals, mdf[:])
}

func (s *aztestsSuite) TestBlobDownloadDataIfModifiedSinceTrue(c *chk.C) {
	currentTime := getRelativeTimeGMT(-10)

	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := DownloadBlobOptions{
		ModifiedAccessConditions: &access,
	}
	resp, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.ContentLength(), chk.Equals, int64(len(blockBlobDefaultData)))
}

func (s *aztestsSuite) TestBlobDownloadDataIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	access := ModifiedAccessConditions{
		IfModifiedSince: &currentTime,
	}
	options := DownloadBlobOptions{
		ModifiedAccessConditions: &access,
	}
	_, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobDownloadDataIfUnmodifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)
	access := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := DownloadBlobOptions{
		ModifiedAccessConditions: &access,
	}
	resp, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.ContentLength(), chk.Equals, int64(len(blockBlobDefaultData)))
}

func (s *aztestsSuite) TestBlobDownloadDataIfUnmodifiedSinceFalse(c *chk.C) {
	currentTime := getRelativeTimeGMT(-10)

	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)
	access := ModifiedAccessConditions{
		IfUnmodifiedSince: &currentTime,
	}
	options := DownloadBlobOptions{
		ModifiedAccessConditions: &access,
	}
	_, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobDownloadDataIfMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	resp, err := blobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	access := ModifiedAccessConditions{
		IfMatch: resp.ETag,
	}
	options := DownloadBlobOptions{
		ModifiedAccessConditions: &access,
	}

	resp2, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp2.ContentLength(), chk.Equals, int64(len(blockBlobDefaultData)))
}

func (s *aztestsSuite) TestBlobDownloadDataIfMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	resp, err := blobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	access := ModifiedAccessConditions{
		IfMatch: resp.ETag,
	}
	options := DownloadBlobOptions{
		ModifiedAccessConditions: &access,
	}

	blobClient.SetMetadata(ctx, nil, nil)

	_, err = blobClient.Download(ctx, &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobDownloadDataIfNoneMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	resp, err := blobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	access := ModifiedAccessConditions{
		IfNoneMatch: resp.ETag,
	}
	options := DownloadBlobOptions{
		ModifiedAccessConditions: &access,
	}

	blobClient.SetMetadata(ctx, nil, nil)

	resp2, err := blobClient.Download(ctx, &options)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp2.ContentLength(), chk.Equals, int64(len(blockBlobDefaultData)))
}

func (s *aztestsSuite) TestBlobDownloadDataIfNoneMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	resp, err := blobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	access := ModifiedAccessConditions{
		IfNoneMatch: resp.ETag,
	}
	options := DownloadBlobOptions{
		ModifiedAccessConditions: &access,
	}

	_, err = blobClient.Download(ctx, &options)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobDeleteNonExistant(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := getBlockBlobClient(c, containerClient)

	_, err := blobClient.Delete(ctx, nil)
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobDeleteSnapshot(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	resp, err := blobClient.CreateSnapshot(ctx, nil)
	c.Assert(err, chk.IsNil)
	snapshotURL := blobClient.WithSnapshot(*resp.Snapshot)

	_, err = snapshotURL.Delete(ctx, nil)
	c.Assert(err, chk.IsNil)

	// TODO enabled after error handling fix
	//validateBlobDeleted(c, snapshotURL)
}

// TODO enable after container client has list command
//func (s *aztestsSuite) TestBlobDeleteSnapshotsInclude(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.CreateSnapshot(ctx, nil)
//	c.Assert(err, chk.IsNil)
//
//	deleteSnapshots := DeleteSnapshotsOptionTypeInclude
//	_, err = blobClient.Delete(ctx, &DeleteBlobOptions{
//		DeleteSnapshots: &deleteSnapshots,
//	})
//	c.Assert(err, chk.IsNil)
//
//	resp, _ := containerClient.ListBlobsFlatSegment(ctx, Marker{},
//		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
//	c.Assert(resp.Segment.BlobItems, chk.HasLen, 0)
//}

// TODO enable after container client has list command
//func (s *aztestsSuite) TestBlobDeleteSnapshotsOnly(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.CreateSnapshot(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	_, err = blobClient.Delete(ctx, DeleteSnapshotsOptionOnly, nil)
//	c.Assert(err, chk.IsNil)
//
//	resp, _ := containerClient.ListBlobsFlatSegment(ctx, Marker{},
//		ListBlobsSegmentOptions{Details: BlobListingDetails{Snapshots: true}})
//	c.Assert(resp.Segment.BlobItems, chk.HasLen, 1)
//	c.Assert(resp.Segment.BlobItems[0].Snapshot == "", chk.Equals, true)
//}
//
func (s *aztestsSuite) TestBlobDeleteSnapshotsNoneWithSnapshots(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	_, err := blobClient.CreateSnapshot(ctx, nil)
	c.Assert(err, chk.IsNil)
	_, err = blobClient.Delete(ctx, nil)
	c.Assert(err, chk.NotNil)
}

func validateBlobDeleted(c *chk.C, blobClient BlobClient) {
	_, err := blobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.NotNil)

	// TODO cannot be checked right now
	serr := err.(*StorageError) // Delete blob is a HEAD request and does not return a ServiceCode in the body
	c.Assert(strings.Contains(serr.Error(), "not exist"), chk.Equals, true)
}

//func (s *aztestsSuite) TestBlobDeleteIfModifiedSinceTrue(c *chk.C) {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.Delete(ctx, DeleteSnapshotsOptionNone,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfModifiedSince: currentTime}})
//	c.Assert(err, chk.IsNil)
//
//	validateBlobDeleted(c, blobClient)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfModifiedSinceFalse(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	_, err := blobClient.Delete(ctx, DeleteSnapshotsOptionNone,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfModifiedSince: currentTime}})
//	validateStorageError(c, err, ServiceCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfUnmodifiedSinceTrue(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	_, err := blobClient.Delete(ctx, DeleteSnapshotsOptionNone,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfUnmodifiedSince: currentTime}})
//	c.Assert(err, chk.IsNil)
//
//	validateBlobDeleted(c, blobClient)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfUnmodifiedSinceFalse(c *chk.C) {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.Delete(ctx, DeleteSnapshotsOptionNone,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfUnmodifiedSince: currentTime}})
//	validateStorageError(c, err, ServiceCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfMatchTrue(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, _ := blobClient.GetProperties(ctx, nil)
//	etag := resp.ETag()
//
//	_, err := blobClient.Delete(ctx, DeleteSnapshotsOptionNone,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfMatch: etag}})
//	c.Assert(err, chk.IsNil)
//
//	validateBlobDeleted(c, blobClient)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfMatchFalse(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, _ := blobClient.GetProperties(ctx, nil)
//	etag := resp.ETag()
//	blobClient.SetMetadata(ctx, nil, nil)
//
//	_, err := blobClient.Delete(ctx, DeleteSnapshotsOptionNone,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfMatch: etag}})
//
//	validateStorageError(c, err, ServiceCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfNoneMatchTrue(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, _ := blobClient.GetProperties(ctx, nil)
//	etag := resp.ETag()
//	blobClient.SetMetadata(ctx, nil, nil)
//
//	_, err := blobClient.Delete(ctx, DeleteSnapshotsOptionNone,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfNoneMatch: etag}})
//	c.Assert(err, chk.IsNil)
//
//	validateBlobDeleted(c, blobClient)
//}
//
//func (s *aztestsSuite) TestBlobDeleteIfNoneMatchFalse(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, _ := blobClient.GetProperties(ctx, nil)
//	etag := resp.ETag()
//
//	_, err := blobClient.Delete(ctx, DeleteSnapshotsOptionNone,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfNoneMatch: etag}})
//	validateStorageError(c, err, ServiceCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfModifiedSinceTrue(c *chk.C) {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	c.Assert(err, chk.IsNil)
//
//	resp, err := blobClient.GetProperties(ctx,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfModifiedSince: currentTime}})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.NewMetadata(), chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfModifiedSinceFalse(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	c.Assert(err, chk.IsNil)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	_, err = blobClient.GetProperties(ctx,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfModifiedSince: currentTime}})
//	c.Assert(err, chk.NotNil)
//	serr := err.(StorageError)
//	c.Assert(serr.Response().StatusCode, chk.Equals, 304) // No service code returned for a HEAD
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfUnmodifiedSinceTrue(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	c.Assert(err, chk.IsNil)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	resp, err := blobClient.GetProperties(ctx,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfUnmodifiedSince: currentTime}})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.NewMetadata(), chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfUnmodifiedSinceFalse(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(-10)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	c.Assert(err, chk.IsNil)
//
//	_, err = blobClient.GetProperties(ctx,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfUnmodifiedSince: currentTime}})
//	c.Assert(err, chk.NotNil)
//	serr := err.(StorageError)
//	c.Assert(serr.Response().StatusCode, chk.Equals, 412)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfMatchTrue(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	c.Assert(err, chk.IsNil)
//
//	resp2, err := blobClient.GetProperties(ctx,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfMatch: resp.ETag()}})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp2.NewMetadata(), chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsOnMissingBlob(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient := containerClient.newBlobClient("MISSING")
//
//	_, err := blobClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.NotNil)
//	serr := err.(StorageError)
//	c.Assert(serr.Response().StatusCode, chk.Equals, 404)
//	c.Assert(serr.ServiceCode(), chk.Equals, ServiceCodeBlobNotFound)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfMatchFalse(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.GetProperties(ctx,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfMatch: ETag("garbage")}})
//	c.Assert(err, chk.NotNil)
//	serr := err.(StorageError)
//	c.Assert(serr.Response().StatusCode, chk.Equals, 412)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfNoneMatchTrue(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata, nil)
//	c.Assert(err, chk.IsNil)
//
//	resp, err := blobClient.GetProperties(ctx,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfNoneMatch: ETag("garbage")}})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.NewMetadata(), chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobGetPropsAndMetadataIfNoneMatchFalse(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.SetMetadata(ctx, nil, nil)
//	c.Assert(err, chk.IsNil)
//
//	_, err = blobClient.GetProperties(ctx,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfNoneMatch: resp.ETag()}})
//	c.Assert(err, chk.NotNil)
//	serr := err.(StorageError)
//	c.Assert(serr.Response().StatusCode, chk.Equals, 304)
//}
//

func (s *aztestsSuite) TestBlobSetPropertiesBasic(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	contentType := "my_type"
	contentDisposition := "my_disposition"
	cacheControl := "control"
	contentLanguage := "my_language"
	contentEncoding := "my_encoding"
	headers := BlobHttpHeaders{
		BlobContentType:        &contentType,
		BlobContentDisposition: &contentDisposition,
		BlobCacheControl:       &cacheControl,
		BlobContentLanguage:    &contentLanguage,
		BlobContentEncoding:    &contentEncoding}
	_, err := blobClient.SetHTTPHeaders(ctx, headers, nil)
	c.Assert(err, chk.IsNil)

	resp, _ := blobClient.GetProperties(ctx, nil)
	h := resp.NewHTTPHeaders()
	c.Assert(h, chk.DeepEquals, headers)
}

func (s *aztestsSuite) TestBlobSetPropertiesEmptyValue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	contentType := "my_type"
	_, err := blobClient.SetHTTPHeaders(ctx, BlobHttpHeaders{BlobContentType: &contentType}, nil)
	c.Assert(err, chk.IsNil)

	_, err = blobClient.SetHTTPHeaders(ctx, BlobHttpHeaders{}, nil)
	c.Assert(err, chk.IsNil)

	resp, err := blobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.ContentType, chk.IsNil)
}

func validatePropertiesSet(c *chk.C, blobClient BlockBlobClient, disposition string) {
	resp, err := blobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*resp.ContentDisposition, chk.Equals, disposition)
}

func (s *aztestsSuite) TestBlobSetPropertiesIfModifiedSinceTrue(c *chk.C) {
	currentTime := getRelativeTimeGMT(-10)

	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	_, err := blobClient.SetHTTPHeaders(ctx, BlobHttpHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}})
	c.Assert(err, chk.IsNil)

	validatePropertiesSet(c, blobClient, "my_disposition")
}

func (s *aztestsSuite) TestBlobSetPropertiesIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	_, err := blobClient.SetHTTPHeaders(ctx, BlobHttpHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}})
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobSetPropertiesIfUnmodifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	_, err := blobClient.SetHTTPHeaders(ctx, BlobHttpHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}})
	c.Assert(err, chk.IsNil)

	validatePropertiesSet(c, blobClient, "my_disposition")
}

func (s *aztestsSuite) TestBlobSetPropertiesIfUnmodifiedSinceFalse(c *chk.C) {
	currentTime := getRelativeTimeGMT(-10)

	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	_, err := blobClient.SetHTTPHeaders(ctx, BlobHttpHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}})
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobSetPropertiesIfMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	resp, err := blobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)

	_, err = blobClient.SetHTTPHeaders(ctx, BlobHttpHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag}})
	c.Assert(err, chk.IsNil)

	validatePropertiesSet(c, blobClient, "my_disposition")
}

func (s *aztestsSuite) TestBlobSetPropertiesIfMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	_, err := blobClient.SetHTTPHeaders(ctx, BlobHttpHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: to.StringPtr("garbage")}})
	c.Assert(err, chk.NotNil)
}

func (s *aztestsSuite) TestBlobSetPropertiesIfNoneMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	_, err := blobClient.SetHTTPHeaders(ctx, BlobHttpHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: to.StringPtr("garbage")}})
	c.Assert(err, chk.IsNil)

	validatePropertiesSet(c, blobClient, "my_disposition")
}

func (s *aztestsSuite) TestBlobSetPropertiesIfNoneMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	blobClient, _ := createNewBlockBlob(c, containerClient)

	resp, err := blobClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)

	_, err = blobClient.SetHTTPHeaders(ctx, BlobHttpHeaders{BlobContentDisposition: to.StringPtr("my_disposition")},
		&SetBlobHTTPHeadersOptions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag}})
	c.Assert(err, chk.NotNil)
}

// TODO enable after metadata fix
//func (s *aztestsSuite) TestBlobSetMetadataNil(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, Metadata{"not": "nil"}, nil)
//	c.Assert(err, chk.IsNil)
//
//	_, err = blobClient.SetMetadata(ctx, nil, nil)
//	c.Assert(err, chk.IsNil)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.NewMetadata(), chk.HasLen, 0)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataEmpty(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, Metadata{"not": "nil"}, nil)
//	c.Assert(err, chk.IsNil)
//
//	_, err = blobClient.SetMetadata(ctx, Metadata{}, nil)
//	c.Assert(err, chk.IsNil)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.NewMetadata(), chk.HasLen, 0)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataInvalidField(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, Metadata{"Invalid field!": "value"}, nil)
//	c.Assert(err, chk.NotNil)
//	c.Assert(strings.Contains(err.Error(), invalidHeaderErrorSubstring), chk.Equals, true)
//}
//
//func validateMetadataSet(c *chk.C, blobClient BlockblobClient) {
//	resp, err := blobClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.NewMetadata(), chk.DeepEquals, basicMetadata)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfModifiedSinceTrue(c *chk.C) {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfModifiedSince: currentTime}})
//	c.Assert(err, chk.IsNil)
//
//	validateMetadataSet(c, blobClient)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfModifiedSinceFalse(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfModifiedSince: currentTime}})
//	validateStorageError(c, err, ServiceCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfUnmodifiedSinceTrue(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	currentTime := getRelativeTimeGMT(10)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfUnmodifiedSince: currentTime}})
//	c.Assert(err, chk.IsNil)
//
//	validateMetadataSet(c, blobClient)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfUnmodifiedSinceFalse(c *chk.C) {
//	currentTime := getRelativeTimeGMT(-10)
//
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfUnmodifiedSince: currentTime}})
//	validateStorageError(c, err, ServiceCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfMatchTrue(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//
//	_, err = blobClient.SetMetadata(ctx, basicMetadata,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfMatch: resp.ETag()}})
//	c.Assert(err, chk.IsNil)
//
//	validateMetadataSet(c, blobClient)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfMatchFalse(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfMatch: ETag("garbage")}})
//	validateStorageError(c, err, ServiceCodeConditionNotMet)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfNoneMatchTrue(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err := blobClient.SetMetadata(ctx, basicMetadata,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfNoneMatch: ETag("garbage")}})
//	c.Assert(err, chk.IsNil)
//
//	validateMetadataSet(c, blobClient)
//}
//
//func (s *aztestsSuite) TestBlobSetMetadataIfNoneMatchFalse(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//
//	_, err = blobClient.SetMetadata(ctx, basicMetadata,
//		BlobAccessConditions{ModifiedAccessConditions: ModifiedAccessConditions{IfNoneMatch: resp.ETag()}})
//	validateStorageError(c, err, ServiceCodeConditionNotMet)
//}
//
// TODO problem with blob versioning behavior change
//func testBlobsUndeleteImpl(c *chk.C, bsu ServiceURL) error {
//	//containerClient, _ := createNewContainer(c, bsu)
//	//defer deleteContainer(c, containerClient)
//	//blobClient, _ := createNewBlockBlob(c, containerClient)
//	//
//	//_, err := blobClient.Delete(ctx, DeleteSnapshotsOptionNone, nil)
//	//c.Assert(err, chk.IsNil) // This call will not have errors related to slow update of service properties, so we assert.
//	//
//	//_, err = blobClient.Undelete(ctx)
//	//if err != nil { // We want to give the wrapper method a chance to check if it was an error related to the service properties update.
//	//	return err
//	//}
//	//
//	//resp, err := blobClient.GetProperties(ctx, nil)
//	//if err != nil {
//	//	return errors.New(string(err.(StorageError).ServiceCode()))
//	//}
//	//c.Assert(resp.BlobType(), chk.Equals, BlobBlockBlob) // We could check any property. This is just to double check it was undeleted.
//	return nil
//}
////
//func (s *aztestsSuite) TestBlobsUndelete(c *chk.C) {
//	bsu := getBSU()
//
//	runTestRequiringServiceProperties(c, bsu, string(404), enableSoftDelete, testBlobsUndeleteImpl, disableSoftDelete)
//}

// TODO enable after page blob is supported
//func setAndCheckBlobTier(c *chk.C, blobClient BlockBlobClient, tier AccessTier) {
//	_, err := blobClient.SetTier(ctx, tier, nil)
//	c.Assert(err, chk.IsNil)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(*resp.AccessTier, chk.Equals, string(tier))
//}
//
//func (s *aztestsSuite) TestBlobSetTierAllTiers(c *chk.C) {
//	bsu := getBSU()
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	setAndCheckBlobTier(c, blobClient, AccessTierHot)
//	setAndCheckBlobTier(c, blobClient, AccessTierCool)
//	setAndCheckBlobTier(c, blobClient, AccessTierArchive)
//
//	bsu, err = getPremiumBSU()
//	if err != nil {
//		c.Skip(err.Error())
//	}
//
//	containerClient, _ = createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	pageblobClient, _ := createNewPageBlob(c, containerClient)
//
//	setAndCheckBlobTier(c, containerClient, pageblobClient.BlobClient, AccessTierP4)
//	setAndCheckBlobTier(c, containerClient, pageblobClient.BlobClient, AccessTierP6)
//	setAndCheckBlobTier(c, containerClient, pageblobClient.blobClient, AccessTierP10)
//	setAndCheckBlobTier(c, containerClient, pageblobClient.blobClient, AccessTierP20)
//	setAndCheckBlobTier(c, containerClient, pageblobClient.blobClient, AccessTierP30)
//	setAndCheckBlobTier(c, containerClient, pageblobClient.blobClient, AccessTierP40)
//	setAndCheckBlobTier(c, containerClient, pageblobClient.blobClient, AccessTierP50)
//}

//func (s *aztestsSuite) TestBlobTierInferred(c *chk.C) {
//	bsu, err := getPremiumBSU()
//	if err != nil {
//		c.Skip(err.Error())
//	}
//
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewPageBlob(c, containerClient)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.AccessTierInferred(), chk.Equals, "true")
//
//	resp2, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp2.Segment.BlobItems[0].Properties.AccessTierInferred, chk.NotNil)
//	c.Assert(resp2.Segment.BlobItems[0].Properties.AccessTier, chk.Not(chk.Equals), "")
//
//	_, err = blobClient.SetTier(ctx, AccessTierP4, LeaseAccessConditions{})
//	c.Assert(err, chk.IsNil)
//
//	resp, err = blobClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.AccessTierInferred(), chk.Equals, "")
//
//	resp2, err = containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp2.Segment.BlobItems[0].Properties.AccessTierInferred, chk.IsNil) // AccessTierInferred never returned if false
//}
//
//func (s *aztestsSuite) TestBlobArchiveStatus(c *chk.C) {
//	bsu, err := getBlobStorageBSU()
//	if err != nil {
//		c.Skip(err.Error())
//	}
//
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err = blobClient.SetTier(ctx, AccessTierArchive, LeaseAccessConditions{})
//	c.Assert(err, chk.IsNil)
//	_, err = blobClient.SetTier(ctx, AccessTierCool, LeaseAccessConditions{})
//	c.Assert(err, chk.IsNil)
//
//	resp, err := blobClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.ArchiveStatus(), chk.Equals, string(ArchiveStatusRehydratePendingToCool))
//
//	resp2, err := containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp2.Segment.BlobItems[0].Properties.ArchiveStatus, chk.Equals, ArchiveStatusRehydratePendingToCool)
//
//	// delete first blob
//	_, err = blobClient.Delete(ctx, DeleteSnapshotsOptionNone, nil)
//	c.Assert(err, chk.IsNil)
//
//	blobClient, _ = createNewBlockBlob(c, containerClient)
//
//	_, err = blobClient.SetTier(ctx, AccessTierArchive, LeaseAccessConditions{})
//	c.Assert(err, chk.IsNil)
//	_, err = blobClient.SetTier(ctx, AccessTierHot, LeaseAccessConditions{})
//	c.Assert(err, chk.IsNil)
//
//	resp, err = blobClient.GetProperties(ctx, nil)
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp.ArchiveStatus(), chk.Equals, string(ArchiveStatusRehydratePendingToHot))
//
//	resp2, err = containerClient.ListBlobsFlatSegment(ctx, Marker{}, ListBlobsSegmentOptions{})
//	c.Assert(err, chk.IsNil)
//	c.Assert(resp2.Segment.BlobItems[0].Properties.ArchiveStatus, chk.Equals, ArchiveStatusRehydratePendingToHot)
//}
//
//func (s *aztestsSuite) TestBlobTierInvalidValue(c *chk.C) {
//	bsu, err := getBlobStorageBSU()
//	if err != nil {
//		c.Skip(err.Error())
//	}
//
//	containerClient, _ := createNewContainer(c, bsu)
//	defer deleteContainer(c, containerClient)
//	blobClient, _ := createNewBlockBlob(c, containerClient)
//
//	_, err = blobClient.SetTier(ctx, AccessTierType("garbage"), LeaseAccessConditions{})
//	validateStorageError(c, err, ServiceCodeInvalidHeaderValue)
//}
//
func (s *aztestsSuite) TestblobClientPartsSASQueryTimes(c *chk.C) {
	StartTimesInputs := []string{
		"2020-04-20",
		"2020-04-20T07:00Z",
		"2020-04-20T07:15:00Z",
		"2020-04-20T07:30:00.1234567Z",
	}
	StartTimesExpected := []time.Time{
		time.Date(2020, time.April, 20, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 20, 7, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 20, 7, 15, 0, 0, time.UTC),
		time.Date(2020, time.April, 20, 7, 30, 0, 123456700, time.UTC),
	}
	ExpiryTimesInputs := []string{
		"2020-04-21",
		"2020-04-20T08:00Z",
		"2020-04-20T08:15:00Z",
		"2020-04-20T08:30:00.2345678Z",
	}
	ExpiryTimesExpected := []time.Time{
		time.Date(2020, time.April, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 20, 8, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 20, 8, 15, 0, 0, time.UTC),
		time.Date(2020, time.April, 20, 8, 30, 0, 234567800, time.UTC),
	}

	for i := 0; i < len(StartTimesInputs); i++ {
		urlString :=
			"https://myaccount.blob.core.windows.net/mycontainer/mydirectory/myfile.txt?" +
				"se=" + url.QueryEscape(ExpiryTimesInputs[i]) + "&" +
				"sig=NotASignature&" +
				"sp=r&" +
				"spr=https&" +
				"sr=b&" +
				"st=" + url.QueryEscape(StartTimesInputs[i]) + "&" +
				"sv=2019-10-10"
		url, _ := url.Parse(urlString)

		parts := NewBlobURLParts(*url)
		c.Assert(parts.Scheme, chk.Equals, "https")
		c.Assert(parts.Host, chk.Equals, "myaccount.blob.core.windows.net")
		c.Assert(parts.ContainerName, chk.Equals, "mycontainer")
		c.Assert(parts.BlobName, chk.Equals, "mydirectory/myfile.txt")

		sas := parts.SAS
		c.Assert(sas.StartTime(), chk.Equals, StartTimesExpected[i])
		c.Assert(sas.ExpiryTime(), chk.Equals, ExpiryTimesExpected[i])

		uResult := parts.URL()
		c.Assert(uResult.String(), chk.Equals, urlString)
	}
}

func (s *aztestsSuite) TestDownloadBlockBlobUnexpectedEOF(c *chk.C) {
	bsu := getBSU()
	cURL, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, cURL)
	bURL, _ := createNewBlockBlob(c, cURL) // This uploads for us.
	resp, err := bURL.Download(ctx, nil)
	c.Assert(err, chk.IsNil)

	// Verify that we can inject errors first.
	reader := resp.Body(InjectErrorInRetryReaderOptions(errors.New("unrecoverable error")))

	_, err = ioutil.ReadAll(reader)
	c.Assert(err, chk.NotNil)
	c.Assert(err.Error(), chk.Equals, "unrecoverable error")

	// Then inject the retryable error.
	reader = resp.Body(InjectErrorInRetryReaderOptions(io.ErrUnexpectedEOF))

	buf, err := ioutil.ReadAll(reader)
	c.Assert(err, chk.IsNil)
	c.Assert(buf, chk.DeepEquals, []byte(blockBlobDefaultData))
}

func InjectErrorInRetryReaderOptions(err error) RetryReaderOptions {
	return RetryReaderOptions{
		MaxRetryRequests:       1,
		doInjectError:          true,
		doInjectErrorRound:     0,
		injectedError:          err,
		NotifyFailedRead:       nil,
		TreatEarlyCloseAsError: false,
	}
}
