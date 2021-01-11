package azblob

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	guuid "github.com/google/uuid"
	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
	"io"
	"io/ioutil"
	"math"
	"strings"
	"time"
)

// a single blockID used in tests when only a single ID is needed
var blockID string

func init() {
	u := [64]byte{}
	binary.BigEndian.PutUint32(u[len(guuid.UUID{}):], math.MaxUint32)
	blockID = base64.StdEncoding.EncodeToString(u[:])
}

func (s *aztestsSuite) TestStageGetBlocks(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	bbClient := containerClient.NewBlockBlobClient(generateBlobName())

	data := []string{"Azure ", "Storage ", "Block ", "Blob."}
	base64BlockIDs := make([]string, len(data))

	for index, d := range data {
		base64BlockIDs[index] = blockIDIntToBase64(index)
		putResp, err := bbClient.StageBlock(context.Background(), base64BlockIDs[index], strings.NewReader(d), nil)
		c.Assert(err, chk.IsNil)
		c.Assert(putResp.RawResponse.StatusCode, chk.Equals, 201)
		c.Assert(putResp.ContentMD5, chk.IsNil)
		c.Assert(putResp.RequestID, chk.NotNil)
		c.Assert(putResp.Version, chk.NotNil)
		c.Assert(putResp.Date, chk.NotNil)
		c.Assert((*putResp.Date).IsZero(), chk.Equals, false)
	}

	blockList, err := bbClient.GetBlockList(context.Background(), BlockListTypeAll, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(blockList.RawResponse.StatusCode, chk.Equals, 200)
	c.Assert(blockList.LastModified, chk.IsNil)
	c.Assert(blockList.ETag, chk.IsNil)
	c.Assert(blockList.ContentType, chk.NotNil)
	c.Assert(blockList.BlobContentLength, chk.IsNil)
	c.Assert(blockList.RequestID, chk.NotNil)
	c.Assert(blockList.Version, chk.NotNil)
	c.Assert(blockList.Date, chk.NotNil)
	c.Assert((*blockList.Date).IsZero(), chk.Equals, false)
	c.Assert(blockList.BlockList, chk.NotNil)
	c.Assert(blockList.BlockList.CommittedBlocks, chk.IsNil)
	c.Assert(blockList.BlockList.UncommittedBlocks, chk.NotNil)
	c.Assert(*blockList.BlockList.UncommittedBlocks, chk.HasLen, len(data))

	listResp, err := bbClient.CommitBlockList(context.Background(), base64BlockIDs, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(listResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(listResp.LastModified, chk.NotNil)
	c.Assert((*listResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(listResp.ETag, chk.NotNil)
	c.Assert(listResp.RequestID, chk.NotNil)
	c.Assert(listResp.Version, chk.NotNil)
	c.Assert(listResp.Date, chk.NotNil)
	c.Assert((*listResp.Date).IsZero(), chk.Equals, false)

	blockList, err = bbClient.GetBlockList(context.Background(), BlockListTypeAll, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(blockList.RawResponse.StatusCode, chk.Equals, 200)
	c.Assert(blockList.LastModified, chk.NotNil)
	c.Assert((*blockList.LastModified).IsZero(), chk.Equals, false)
	c.Assert(blockList.ETag, chk.NotNil)
	c.Assert(blockList.ContentType, chk.NotNil)
	c.Assert(*blockList.BlobContentLength, chk.Equals, int64(25))
	c.Assert(blockList.RequestID, chk.NotNil)
	c.Assert(blockList.Version, chk.NotNil)
	c.Assert(blockList.Date, chk.NotNil)
	c.Assert((*blockList.Date).IsZero(), chk.Equals, false)
	c.Assert(blockList.BlockList, chk.NotNil)
	c.Assert(blockList.BlockList.CommittedBlocks, chk.NotNil)
	c.Assert(blockList.BlockList.UncommittedBlocks, chk.IsNil)
	c.Assert(*blockList.BlockList.CommittedBlocks, chk.HasLen, len(data))
}

func (s *aztestsSuite) TestStageBlockFromURL(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)

	ctx := context.Background() // Use default Background context
	srcBlob := containerClient.NewBlockBlobClient(generateBlobName())
	destBlob := containerClient.NewBlockBlobClient(generateBlobName())

	// Prepare source bbClient for copy.
	uploadSrcResp, err := srcBlob.Upload(ctx, rsc, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.RawResponse.StatusCode, chk.Equals, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())

	srcBlobParts.SAS, err = BlobSASSignatureValues{
		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName: srcBlobParts.ContainerName,
		BlobName:      srcBlobParts.BlobName,
		Permissions:   BlobSASPermissions{Read: true}.String(),
	}.NewSASQueryParameters(credential)
	c.Assert(err, chk.IsNil)

	srcBlobURLWithSAS := srcBlobParts.URL()

	// Stage blocks from URL.
	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))

	offset, count := int64(0), int64(contentSize/2)
	options1 := StageBlockFromURLOptions{
		Offset: &offset,
		Count:  &count,
	}

	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
	c.Assert(err, chk.IsNil)
	c.Assert(stageResp1.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(stageResp1.ContentMD5, chk.Not(chk.Equals), "")
	c.Assert(stageResp1.RequestID, chk.Not(chk.Equals), "")
	c.Assert(stageResp1.Version, chk.Not(chk.Equals), "")
	c.Assert(stageResp1.Date.IsZero(), chk.Equals, false)

	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))

	offset2, count2 := int64(contentSize/2), int64(CountToEnd)
	options2 := StageBlockFromURLOptions{
		Offset: &offset2,
		Count:  &count2,
	}

	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
	c.Assert(err, chk.IsNil)
	c.Assert(stageResp2.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(stageResp2.ContentMD5, chk.Not(chk.Equals), "")
	c.Assert(stageResp2.RequestID, chk.Not(chk.Equals), "")
	c.Assert(stageResp2.Version, chk.Not(chk.Equals), "")
	c.Assert(stageResp2.Date.IsZero(), chk.Equals, false)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(blockList.RawResponse.StatusCode, chk.Equals, 200)
	c.Assert(blockList.BlockList, chk.NotNil)
	c.Assert(blockList.BlockList.CommittedBlocks, chk.IsNil)
	c.Assert(blockList.BlockList.UncommittedBlocks, chk.NotNil)
	c.Assert(*blockList.BlockList.UncommittedBlocks, chk.HasLen, 2)

	// Commit block list.
	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(listResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(listResp.LastModified, chk.NotNil)
	c.Assert((*listResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(listResp.ETag, chk.NotNil)
	c.Assert(listResp.RequestID, chk.NotNil)
	c.Assert(listResp.Version, chk.NotNil)
	c.Assert(listResp.Date, chk.NotNil)
	c.Assert((*listResp.Date).IsZero(), chk.Equals, false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, content)
}

func (s *aztestsSuite) TestCopyBlockBlobFromURL(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	const contentSize = 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	contentMD5 := md5.Sum(content)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)
	ctx := context.Background()

	srcBlob := containerClient.NewBlockBlobClient("srcblob")
	destBlob := containerClient.NewBlockBlobClient("destblob")

	// Prepare source bbClient for copy.
	uploadSrcResp, err := srcBlob.Upload(ctx, rsc, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.RawResponse.StatusCode, chk.Equals, 201)

	// Get source blob url with SAS for StageFromURL.
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

	// Invoke copy bbClient from URL.
	sourceContentMD5 := contentMD5[:]
	copyBlockBlobFromURLOptions := CopyBlockBlobFromURLOptions{
		Metadata:         &map[string]string{"foo": "bar"},
		SourceContentMd5: &sourceContentMD5,
	}
	resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 202)
	c.Assert(resp.ETag, chk.NotNil)
	c.Assert(resp.RequestID, chk.NotNil)
	c.Assert(resp.Version, chk.NotNil)
	c.Assert(resp.Date, chk.NotNil)
	c.Assert((*resp.Date).IsZero(), chk.Equals, false)
	c.Assert(resp.CopyID, chk.NotNil)
	c.Assert(*resp.ContentMD5, chk.DeepEquals, sourceContentMD5)
	c.Assert(*resp.CopyStatus, chk.DeepEquals, "success")

	// Make sure the metadata got copied over
	getPropResp, err := destBlob.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	metadata := getPropResp.NewMetadata()
	c.Assert(metadata, chk.NotNil)
	c.Assert(metadata, chk.HasLen, 1)
	c.Assert(metadata, chk.DeepEquals, map[string]string{"foo": "bar"})

	// Check data integrity through downloading.
	downloadResp, err := destBlob.Download(ctx, nil)
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, content)

	// Edge case 1: Provide bad MD5 and make sure the copy fails
	_, badMD5 := getRandomDataAndReader(16)
	copyBlockBlobFromURLOptions1 := CopyBlockBlobFromURLOptions{
		SourceContentMd5: &badMD5,
	}
	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions1)
	c.Assert(err, chk.NotNil)

	// Edge case 2: Not providing any source MD5 should see the CRC getting returned instead
	copyBlockBlobFromURLOptions2 := CopyBlockBlobFromURLOptions{
		SourceContentMd5: &sourceContentMD5,
	}
	resp, err = destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions2)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.RawResponse.StatusCode, chk.Equals, 202)
	c.Assert(*resp.CopyStatus, chk.DeepEquals, "success")
}

func (s *aztestsSuite) TestBlobSASQueryParamOverrideResponseHeaders(c *chk.C) {
	bsu := getBSU()
	credential, err := getGenericCredential("")
	if err != nil {
		c.Fatal("Invalid credential")
	}
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	const contentSize = 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	contentMD5 := md5.Sum(content)
	_ = contentMD5
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)
	ctx := context.Background()

	bbClient := containerClient.NewBlockBlobClient(generateBlobName())

	uploadSrcResp, err := bbClient.Upload(ctx, rsc, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.RawResponse.StatusCode, chk.Equals, 201)

	// Get blob url with SAS.
	blobParts := NewBlobURLParts(bbClient.URL())

	cacheControlVal := "cache-control-override"
	contentDispositionVal := "content-disposition-override"
	contentEncodingVal := "content-encoding-override"
	contentLanguageVal := "content-language-override"
	contentTypeVal := "content-type-override"

	// Append User Delegation SAS token to URL
	blobParts.SAS, err = BlobSASSignatureValues{
		Protocol:           SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
		ExpiryTime:         time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
		ContainerName:      blobParts.ContainerName,
		BlobName:           blobParts.BlobName,
		Permissions:        BlobSASPermissions{Read: true}.String(),
		CacheControl:       cacheControlVal,
		ContentDisposition: contentDispositionVal,
		ContentEncoding:    contentEncodingVal,
		ContentLanguage:    contentLanguageVal,
		ContentType:        contentTypeVal,
	}.NewSASQueryParameters(credential)
	c.Assert(err, chk.IsNil)

	// Generate new bbClient client
	blobURLWithSAS := blobParts.URL()
	blobRawURL := blobURLWithSAS.String()
	c.Assert(blobRawURL, chk.NotNil)
	blobClientWithSAS, err := NewBlockBlobClient(blobURLWithSAS.String(), azcore.AnonymousCredential(), nil)
	c.Assert(err, chk.IsNil)

	gResp, err := blobClientWithSAS.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*gResp.CacheControl, chk.Equals, cacheControlVal)
	c.Assert(*gResp.ContentDisposition, chk.Equals, contentDispositionVal)
	c.Assert(*gResp.ContentEncoding, chk.Equals, contentEncodingVal)
	c.Assert(*gResp.ContentLanguage, chk.Equals, contentLanguageVal)
	c.Assert(*gResp.ContentType, chk.Equals, contentTypeVal)
}

func (s *aztestsSuite) TestStageBlockWithMD5(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	bbClient := containerClient.NewBlockBlobClient(generateBlobName())

	// test put block with valid MD5 value
	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)
	md5Value := md5.Sum(content)
	contentMD5 := md5Value[:]

	stageBlockOptions := StageBlockOptions{
		BlockBlobStageBlockOptions: &BlockBlobStageBlockOptions{
			TransactionalContentMd5: &contentMD5,
		},
	}
	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	putResp, err := bbClient.StageBlock(context.Background(), blockID1, rsc, &stageBlockOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(putResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(*(putResp.ContentMD5), chk.DeepEquals, contentMD5)
	c.Assert(putResp.RequestID, chk.NotNil)
	c.Assert(putResp.Version, chk.NotNil)
	c.Assert(putResp.Date, chk.NotNil)
	c.Assert((*putResp.Date).IsZero(), chk.Equals, false)

	// test put block with bad MD5 value
	_, badContent := getRandomDataAndReader(contentSize)
	badMD5Value := md5.Sum(badContent)
	badContentMD5 := badMD5Value[:]

	badStageBlockOptions := StageBlockOptions{
		BlockBlobStageBlockOptions: &BlockBlobStageBlockOptions{
			TransactionalContentMd5: &badContentMD5,
		},
	}
	_, _ = rsc.Seek(0, io.SeekStart)
	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	_, err = bbClient.StageBlock(context.Background(), blockID2, rsc, &badStageBlockOptions)
	c.Assert(err, chk.NotNil)

	// TODO: Fix issue with storage error interface
	//c.Assert(err.(StorageError), chk.Equals, ServiceCodeMd5Mismatch)
}

func (s *aztestsSuite) TestBlobPutBlobHTTPHeaders(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := getBlockBlobClient(c, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		BlobHttpHeaders: &basicHeaders,
	}
	_, err := bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	c.Assert(err, chk.IsNil)

	resp, err := bbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	h := resp.NewHTTPHeaders()
	h.BlobContentMd5 = nil // the service generates a MD5 value, omit before comparing
	c.Assert(h, chk.DeepEquals, basicHeaders)
}

func (s *aztestsSuite) TestBlobPutBlobMetadataNotEmpty(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := getBlockBlobClient(c, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		Metadata: &basicMetadata,
	}
	_, err := bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	c.Assert(err, chk.IsNil)

	resp, err := bbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	actualMetadata := resp.NewMetadata()
	c.Assert(actualMetadata, chk.NotNil)
	c.Assert(actualMetadata, chk.DeepEquals, basicMetadata)
}

func (s *aztestsSuite) TestBlobPutBlobMetadataEmpty(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := getBlockBlobClient(c, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)

	_, err := bbClient.Upload(ctx, rsc, nil)
	c.Assert(err, chk.IsNil)

	resp, err := bbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.NewMetadata(), chk.HasLen, 0)
}

func (s *aztestsSuite) TestBlobPutBlobMetadataInvalid(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := getBlockBlobClient(c, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		Metadata: &map[string]string{"In valid!": "bar"},
	}
	_, err := bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	c.Assert(err, chk.NotNil)
	c.Assert(strings.Contains(err.Error(), invalidHeaderErrorSubstring), chk.Equals, true)
}

func (s *aztestsSuite) TestBlobPutBlobIfModifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)

	bbClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}
	_, err := bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	c.Assert(err, chk.IsNil)
	validateUpload(c, bbClient.BlobClient)
}

func (s *aztestsSuite) TestBlobPutBlobIfModifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}

	_, err := bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	c.Assert(err, chk.NotNil)

	//TODO: Fix issue with storage error interface
	//c.Assert(strings.Contains(err.Error(), string(ServiceCodeConditionNotMet)), chk.Equals, true)
	//validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutBlobIfUnmodifiedSinceTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(10)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err := bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	c.Assert(err, chk.IsNil)

	validateUpload(c, bbClient.BlobClient)
}

func (s *aztestsSuite) TestBlobPutBlobIfUnmodifiedSinceFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := createNewBlockBlob(c, containerClient)

	currentTime := getRelativeTimeGMT(-10)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}
	_, err := bbClient.Upload(ctx, bytes.NewReader(nil), &uploadBlockBlobOptions)
	_ = err
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutBlobIfMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := createNewBlockBlob(c, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}
	_, err = bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	c.Assert(err, chk.IsNil)

	validateUpload(c, bbClient.BlobClient)
}

func (s *aztestsSuite) TestBlobPutBlobIfMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := createNewBlockBlob(c, containerClient)

	_, err := bbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)

	ifNoneMatch := "garbage"
	uploadBlockBlobOptions := UploadBlockBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: &ifNoneMatch,
		},
	}
	_, err = bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutBlobIfNoneMatchTrue(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := createNewBlockBlob(c, containerClient)

	_, err := bbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)

	ifNoneMatch := "garbage"
	uploadBlockBlobOptions := UploadBlockBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: &ifNoneMatch,
		},
	}

	_, err = bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	c.Assert(err, chk.IsNil)

	validateUpload(c, bbClient.BlobClient)
}

func (s *aztestsSuite) TestBlobPutBlobIfNoneMatchFalse(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	bbClient, _ := createNewBlockBlob(c, containerClient)

	resp, err := bbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)

	content := make([]byte, 0)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)

	uploadBlockBlobOptions := UploadBlockBlobOptions{
		ModifiedAccessConditions: &ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}
	_, err = bbClient.Upload(ctx, rsc, &uploadBlockBlobOptions)
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func setupPutBlockListTest(c *chk.C) (containerClient ContainerClient, bbClient BlockBlobClient, id string) {
	bsu := getBSU()
	containerClient, _ = createNewContainer(c, bsu)
	bbClient, _ = getBlockBlobClient(c, containerClient)

	_, err := bbClient.StageBlock(ctx, blockID, strings.NewReader(blockBlobDefaultData), nil)
	c.Assert(err, chk.IsNil)
	return containerClient, bbClient, blockID
}

func (s *aztestsSuite) TestBlobPutBlockListHTTPHeadersEmpty(c *chk.C) {
	contClient, bbClient, blockId := setupPutBlockListTest(c)
	defer deleteContainer(c, contClient)

	commitBlockListOptions := CommitBlockListOptions{
		BlobHTTPHeaders: &BlobHttpHeaders{BlobContentDisposition: &blobContentDisposition},
	}
	_, err := bbClient.CommitBlockList(ctx, []string{blockId}, &commitBlockListOptions)
	c.Assert(err, chk.IsNil)

	_, err = bbClient.CommitBlockList(ctx, []string{blockId}, nil)
	c.Assert(err, chk.IsNil)

	resp, err := bbClient.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.ContentDisposition, chk.IsNil)
}

func validateBlobCommitted(c *chk.C, bbClient BlockBlobClient) {
	resp, err := bbClient.GetBlockList(ctx, BlockListTypeAll, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*(resp.BlockList.CommittedBlocks), chk.HasLen, 1)
}

func (s *aztestsSuite) TestBlobPutBlockListIfModifiedSinceTrue(c *chk.C) {
	contClient, bbClient, blockId := setupPutBlockListTest(c)
	defer deleteContainer(c, contClient)
	_, err := bbClient.CommitBlockList(ctx, []string{blockId}, nil) // The bbClient must actually exist to have a modifed time
	c.Assert(err, chk.IsNil)

	currentTime := getRelativeTimeGMT(-10)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(ctx, []string{blockId}, &commitBlockListOptions)
	c.Assert(err, chk.IsNil)

	validateBlobCommitted(c, bbClient)
}

func (s *aztestsSuite) TestBlobPutBlockListIfModifiedSinceFalse(c *chk.C) {
	contClient, bbClient, blockId := setupPutBlockListTest(c)
	defer deleteContainer(c, contClient)

	currentTime := getRelativeTimeGMT(10)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	}
	_, err := bbClient.CommitBlockList(ctx, []string{blockId}, &commitBlockListOptions)
	_ = err
	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutBlockListIfUnmodifiedSinceTrue(c *chk.C) {
	contClient, bbClient, blockId := setupPutBlockListTest(c)
	defer deleteContainer(c, contClient)
	_, err := bbClient.CommitBlockList(ctx, []string{blockId}, nil) // The bbClient must actually exist to have a modifed time
	c.Assert(err, chk.IsNil)

	currentTime := getRelativeTimeGMT(10)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(ctx, []string{blockId}, &commitBlockListOptions)
	c.Assert(err, chk.IsNil)

	validateBlobCommitted(c, bbClient)
}

func (s *aztestsSuite) TestBlobPutBlockListIfUnmodifiedSinceFalse(c *chk.C) {
	contClient, bbClient, blockId := setupPutBlockListTest(c)
	_, err := bbClient.CommitBlockList(ctx, []string{blockId}, nil) // The bbClient must actually exist to have a modifed time
	c.Assert(err, chk.IsNil)
	defer deleteContainer(c, contClient)

	currentTime := getRelativeTimeGMT(-10)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfUnmodifiedSince: &currentTime}},
	}
	_, err = bbClient.CommitBlockList(ctx, []string{blockId}, &commitBlockListOptions)

	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutBlockListIfMatchTrue(c *chk.C) {
	contClient, bbClient, blockId := setupPutBlockListTest(c)
	defer deleteContainer(c, contClient)
	resp, err := bbClient.CommitBlockList(ctx, []string{blockId}, nil) // The bbClient must actually exist to have a modifed time
	c.Assert(err, chk.IsNil)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: resp.ETag}},
	}
	_, err = bbClient.CommitBlockList(ctx, []string{blockId}, &commitBlockListOptions)
	c.Assert(err, chk.IsNil)

	validateBlobCommitted(c, bbClient)
}

func (s *aztestsSuite) TestBlobPutBlockListIfMatchFalse(c *chk.C) {
	contClient, bbClient, blockId := setupPutBlockListTest(c)
	defer deleteContainer(c, contClient)
	_, err := bbClient.CommitBlockList(ctx, []string{blockId}, nil) // The bbClient must actually exist to have a modifed time
	c.Assert(err, chk.IsNil)

	eTag := "garbage"
	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(ctx, []string{blockId}, &commitBlockListOptions)

	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeConditionNotMet)

}

func (s *aztestsSuite) TestBlobPutBlockListIfNoneMatchTrue(c *chk.C) {
	contClient, bbClient, blockId := setupPutBlockListTest(c)
	defer deleteContainer(c, contClient)
	_, err := bbClient.CommitBlockList(ctx, []string{blockId}, nil) // The bbClient must actually exist to have a modifed time
	c.Assert(err, chk.IsNil)

	eTag := "garbage"
	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: &eTag}},
	}
	_, err = bbClient.CommitBlockList(ctx, []string{blockId}, &commitBlockListOptions)
	c.Assert(err, chk.IsNil)

	validateBlobCommitted(c, bbClient)
}

func (s *aztestsSuite) TestBlobPutBlockListIfNoneMatchFalse(c *chk.C) {
	contClient, bbClient, blockId := setupPutBlockListTest(c)
	defer deleteContainer(c, contClient)
	resp, err := bbClient.CommitBlockList(ctx, []string{blockId}, nil) // The bbClient must actually exist to have a modifed time
	c.Assert(err, chk.IsNil)

	commitBlockListOptions := CommitBlockListOptions{
		BlobAccessConditions: BlobAccessConditions{ModifiedAccessConditions: &ModifiedAccessConditions{IfNoneMatch: resp.ETag}},
	}
	_, err = bbClient.CommitBlockList(ctx, []string{blockId}, &commitBlockListOptions)

	// TODO: Fix issue with storage error interface
	//validateStorageError(c, err, ServiceCodeConditionNotMet)
}

func (s *aztestsSuite) TestBlobPutBlockListValidateData(c *chk.C) {
	contClient, bbClient, blockId := setupPutBlockListTest(c)
	defer deleteContainer(c, contClient)

	_, err := bbClient.CommitBlockList(ctx, []string{blockId}, nil)

	resp, err := bbClient.Download(ctx, nil)
	c.Assert(err, chk.IsNil)
	data, _ := ioutil.ReadAll(resp.Response().Body)
	c.Assert(string(data), chk.Equals, blockBlobDefaultData)
}

func (s *aztestsSuite) TestBlobPutBlockListModifyBlob(c *chk.C) {
	contClient, bbClient, blockId := setupPutBlockListTest(c)
	defer deleteContainer(c, contClient)

	_, err := bbClient.CommitBlockList(ctx, []string{blockId}, nil)
	c.Assert(err, chk.IsNil)

	_, err = bbClient.StageBlock(ctx, "0001", bytes.NewReader([]byte("new data")), nil)
	c.Assert(err, chk.IsNil)
	_, err = bbClient.StageBlock(ctx, "0010", bytes.NewReader([]byte("new data")), nil)
	c.Assert(err, chk.IsNil)
	_, err = bbClient.StageBlock(ctx, "0011", bytes.NewReader([]byte("new data")), nil)
	c.Assert(err, chk.IsNil)
	_, err = bbClient.StageBlock(ctx, "0100", bytes.NewReader([]byte("new data")), nil)
	c.Assert(err, chk.IsNil)

	_, err = bbClient.CommitBlockList(ctx, []string{"0001", "0011"}, nil)
	c.Assert(err, chk.IsNil)

	resp, err := bbClient.GetBlockList(ctx, BlockListTypeAll, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*(resp.BlockList.CommittedBlocks), chk.HasLen, 2)
	committed := *(resp.BlockList.CommittedBlocks)
	c.Assert(*(committed[0].Name), chk.Equals, "0001")
	c.Assert(*(committed[1].Name), chk.Equals, "0011")
	c.Assert(resp.BlockList.UncommittedBlocks, chk.IsNil)
}

func (s *aztestsSuite) TestSetTierOnBlobUpload(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)
	for _, tier := range []AccessTier{AccessTierArchive, AccessTierCool, AccessTierHot} {
		bbClient, _ := getBlockBlobClient(c, containerClient)

		uploadBlockBlobOptions := UploadBlockBlobOptions{
			BlobHttpHeaders: &basicHeaders,
			Tier:            &tier,
		}
		_, err := bbClient.Upload(ctx, strings.NewReader(blockBlobDefaultData), &uploadBlockBlobOptions)
		c.Assert(err, chk.IsNil)

		resp, err := bbClient.GetProperties(ctx, nil)
		c.Assert(err, chk.IsNil)
		c.Assert(*resp.AccessTier, chk.Equals, string(tier))
	}
}

func (s *aztestsSuite) TestBlobSetTierOnCommit(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	for _, tier := range []AccessTier{AccessTierCool, AccessTierHot} {
		bbClient, _ := getBlockBlobClient(c, containerClient)

		_, err := bbClient.StageBlock(ctx, blockID, strings.NewReader(blockBlobDefaultData), nil)
		c.Assert(err, chk.IsNil)

		commitBlockListOptions := CommitBlockListOptions{
			Tier: &tier,
		}
		_, err = bbClient.CommitBlockList(ctx, []string{blockID}, &commitBlockListOptions)
		c.Assert(err, chk.IsNil)

		resp, err := bbClient.GetBlockList(ctx, BlockListTypeCommitted, nil)
		c.Assert(err, chk.IsNil)
		c.Assert(resp.BlockList, chk.NotNil)
		c.Assert(resp.BlockList.CommittedBlocks, chk.NotNil)
		c.Assert(resp.BlockList.UncommittedBlocks, chk.IsNil)
		c.Assert(*(resp.BlockList.CommittedBlocks), chk.HasLen, 1)

	}
}

func (s *aztestsSuite) TestSetTierOnCopyBlockBlobFromURL(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	const contentSize = 1 * 1024 // 1 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)
	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient(generateBlobName())

	tier := AccessTierCool
	uploadSrcResp, err := srcBlob.Upload(ctx, rsc, &UploadBlockBlobOptions{Tier: &tier})
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.RawResponse.StatusCode, chk.Equals, 201)

	// Get source blob url with SAS for StageFromURL.
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
	c.Assert(err, chk.IsNil)

	srcBlobParts := NewBlobURLParts(srcBlob.URL())
	srcBlobParts.SAS = sasQueryParams
	srcBlobURLWithSAS := srcBlobParts.URL()

	for _, tier := range []AccessTier{AccessTierArchive, AccessTierCool, AccessTierHot} {
		destBlob := containerClient.NewBlockBlobClient(generateBlobName())

		copyBlockBlobFromURLOptions := CopyBlockBlobFromURLOptions{
			Tier:     &tier,
			Metadata: &map[string]string{"foo": "bar"},
		}
		resp, err := destBlob.CopyFromURL(ctx, srcBlobURLWithSAS, &copyBlockBlobFromURLOptions)
		c.Assert(err, chk.IsNil)
		c.Assert(resp.RawResponse.StatusCode, chk.Equals, 202)
		c.Assert(*resp.CopyStatus, chk.DeepEquals, "success")

		destBlobPropResp, err := destBlob.GetProperties(ctx, nil)
		c.Assert(err, chk.IsNil)
		c.Assert(*destBlobPropResp.AccessTier, chk.Equals, string(tier))

	}
}

func (s *aztestsSuite) TestSetTierOnStageBlockFromURL(c *chk.C) {
	bsu := getBSU()
	containerClient, _ := createNewContainer(c, bsu)
	defer deleteContainer(c, containerClient)

	contentSize := 8 * 1024 // 8 KB
	content := make([]byte, contentSize)
	body := bytes.NewReader(content)
	rsc := azcore.NopCloser(body)
	ctx := context.Background()
	srcBlob := containerClient.NewBlockBlobClient("src" + generateBlobName())
	destBlob := containerClient.NewBlockBlobClient("dst" + generateBlobName())
	tier := AccessTierCool
	uploadSrcResp, err := srcBlob.Upload(ctx, rsc, &UploadBlockBlobOptions{Tier: &tier})
	c.Assert(err, chk.IsNil)
	c.Assert(uploadSrcResp.RawResponse.StatusCode, chk.Equals, 201)

	// Get source blob url with SAS for StageFromURL.
	srcBlobParts := NewBlobURLParts(srcBlob.URL())
	credential, err := getGenericCredential("")
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

	// Stage blocks from URL.
	blockID1 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))
	offset1, count1 := int64(0), int64(4*1024)
	options1 := StageBlockFromURLOptions{
		Offset: &offset1,
		Count:  &count1,
	}
	stageResp1, err := destBlob.StageBlockFromURL(ctx, blockID1, srcBlobURLWithSAS, 0, &options1)
	c.Assert(err, chk.IsNil)
	c.Assert(stageResp1.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(stageResp1.ContentMD5, chk.Not(chk.Equals), "")
	c.Assert(stageResp1.RequestID, chk.Not(chk.Equals), "")
	c.Assert(stageResp1.Version, chk.Not(chk.Equals), "")
	c.Assert(stageResp1.Date.IsZero(), chk.Equals, false)

	blockID2 := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 1)))
	offset2, count2 := int64(4*1024), int64(CountToEnd)
	options2 := StageBlockFromURLOptions{
		Offset: &offset2,
		Count:  &count2,
	}
	stageResp2, err := destBlob.StageBlockFromURL(ctx, blockID2, srcBlobURLWithSAS, 0, &options2)
	c.Assert(err, chk.IsNil)
	c.Assert(stageResp2.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(stageResp2.ContentMD5, chk.Not(chk.Equals), "")
	c.Assert(stageResp2.RequestID, chk.Not(chk.Equals), "")
	c.Assert(stageResp2.Version, chk.Not(chk.Equals), "")
	c.Assert(stageResp2.Date.IsZero(), chk.Equals, false)

	// Check block list.
	blockList, err := destBlob.GetBlockList(context.Background(), BlockListTypeAll, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(blockList.RawResponse.StatusCode, chk.Equals, 200)
	c.Assert(blockList.BlockList, chk.NotNil)
	c.Assert(blockList.BlockList.CommittedBlocks, chk.IsNil)
	c.Assert(blockList.BlockList.UncommittedBlocks, chk.NotNil)
	c.Assert(*blockList.BlockList.UncommittedBlocks, chk.HasLen, 2)

	// Commit block list.
	commitBlockListOptions := CommitBlockListOptions{
		Tier: &tier,
	}
	listResp, err := destBlob.CommitBlockList(context.Background(), []string{blockID1, blockID2}, &commitBlockListOptions)
	c.Assert(err, chk.IsNil)
	c.Assert(listResp.RawResponse.StatusCode, chk.Equals, 201)
	c.Assert(listResp.LastModified, chk.NotNil)
	c.Assert((*listResp.LastModified).IsZero(), chk.Equals, false)
	c.Assert(listResp.ETag, chk.NotNil)
	c.Assert(listResp.RequestID, chk.NotNil)
	c.Assert(listResp.Version, chk.NotNil)
	c.Assert(listResp.Date, chk.NotNil)
	c.Assert((*listResp.Date).IsZero(), chk.Equals, false)

	// Check data integrity through downloading.
	downloadResp, err := destBlob.BlobClient.Download(ctx, nil)
	c.Assert(err, chk.IsNil)
	destData, err := ioutil.ReadAll(downloadResp.Body(RetryReaderOptions{}))
	c.Assert(err, chk.IsNil)
	c.Assert(destData, chk.DeepEquals, content)

	// Get properties to validate the tier
	destBlobPropResp, err := destBlob.GetProperties(ctx, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(*destBlobPropResp.AccessTier, chk.Equals, string(tier))
}
