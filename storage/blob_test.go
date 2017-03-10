package storage

import (
	"bytes"
	"crypto/rand"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	chk "gopkg.in/check.v1"
)

type StorageBlobSuite struct{}

var _ = chk.Suite(&StorageBlobSuite{})

const (
	testContainerPrefix = "zzzztest-"

	dummyStorageAccount = "golangrocksonazure"
	dummyMiniStorageKey = "YmFy"
)

func getBlobClient(c *chk.C) BlobStorageClient {
	return getBasicClient(c).GetBlobService()
}

func (s *StorageBlobSuite) Test_buildPath(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference("lol")
	b := cnt.GetBlobReference("rofl")
	c.Assert(b.buildPath(), chk.Equals, "/lol/rofl")
}

func (s *StorageBlobSuite) Test_pathForResource(c *chk.C) {
	c.Assert(pathForResource("lol", ""), chk.Equals, "/lol")
	c.Assert(pathForResource("lol", "blob"), chk.Equals, "/lol/blob")
}

func (s *StorageBlobSuite) TestBlobExists(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	c.Assert(cnt.Create(nil), chk.IsNil)
	b := cnt.GetBlobReference(randName(5))
	defer cnt.Delete(nil)

	c.Assert(b.putSingleBlockBlob([]byte("Hello!")), chk.IsNil)
	defer b.Delete(nil)

	ok, err := b.Exists()
	c.Assert(err, chk.IsNil)
	c.Assert(ok, chk.Equals, true)
	b.Name += ".lol"
	ok, err = b.Exists()
	c.Assert(err, chk.IsNil)
	c.Assert(ok, chk.Equals, false)

}

func (s *StorageBlobSuite) TestGetBlobURL(c *chk.C) {
	cli, err := NewBasicClient(dummyStorageAccount, dummyMiniStorageKey)
	c.Assert(err, chk.IsNil)
	blobCli := cli.GetBlobService()

	cnt := blobCli.GetContainerReference("c")
	b := cnt.GetBlobReference("nested/blob")
	c.Assert(b.GetURL(), chk.Equals, "https://golangrocksonazure.blob.core.windows.net/c/nested/blob")

	cnt.Name = ""
	c.Assert(b.GetURL(), chk.Equals, "https://golangrocksonazure.blob.core.windows.net/$root/nested/blob")

	b.Name = "blob"
	c.Assert(b.GetURL(), chk.Equals, "https://golangrocksonazure.blob.core.windows.net/$root/blob")

}

func (s *StorageBlobSuite) TestGetBlobContainerURL(c *chk.C) {
	cli, err := NewBasicClient(dummyStorageAccount, dummyMiniStorageKey)
	c.Assert(err, chk.IsNil)
	blobCli := cli.GetBlobService()

	cnt := blobCli.GetContainerReference("c")
	b := cnt.GetBlobReference("")
	c.Assert(b.GetURL(), chk.Equals, "https://golangrocksonazure.blob.core.windows.net/c")

	cnt.Name = ""
	c.Assert(b.GetURL(), chk.Equals, "https://golangrocksonazure.blob.core.windows.net/$root")
}

func (s *StorageBlobSuite) TestDeleteBlobIfExists(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	c.Assert(b.Delete(nil), chk.NotNil)

	ok, err := b.DeleteIfExists(nil)
	c.Assert(err, chk.IsNil)
	c.Assert(ok, chk.Equals, false)
}

func (s *StorageBlobSuite) TestDeleteBlobWithConditions(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	c.Assert(b.CreateBlockBlob(nil), chk.IsNil)
	err := b.GetProperties(nil)
	c.Assert(err, chk.IsNil)
	oldProps := b.Properties

	// Update metadata, so Etag changes
	c.Assert(b.SetMetadata(nil), chk.IsNil)
	err = b.GetProperties(nil)
	c.Assert(err, chk.IsNil)
	newProps := b.Properties

	// "Delete if matches old Etag" should fail without deleting.
	options := DeleteBlobOptions{
		IfMatch: oldProps.Etag,
	}
	err = b.Delete(&options)
	c.Assert(err, chk.FitsTypeOf, AzureStorageServiceError{})
	c.Assert(err.(AzureStorageServiceError).StatusCode, chk.Equals, http.StatusPreconditionFailed)
	ok, err := b.Exists()
	c.Assert(err, chk.IsNil)
	c.Assert(ok, chk.Equals, true)

	// "Delete if matches new Etag" should succeed.
	options.IfMatch = newProps.Etag
	err = b.Delete(&options)
	c.Assert(err, chk.IsNil)
	ok, err = b.Exists()
	c.Assert(err, chk.IsNil)
	c.Assert(ok, chk.Equals, false)
}

func (s *StorageBlobSuite) TestGetBlobProperties(c *chk.C) {
	contents := randString(64)

	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	// Nonexisting blob
	err := b.GetProperties(nil)
	c.Assert(err, chk.NotNil)

	// Put the blob
	c.Assert(b.putSingleBlockBlob([]byte(contents)), chk.IsNil)

	// Get blob properties
	err = b.GetProperties(nil)
	c.Assert(err, chk.IsNil)

	c.Assert(b.Properties.ContentLength, chk.Equals, int64(len(contents)))
	c.Assert(b.Properties.ContentType, chk.Equals, "application/octet-stream")
	c.Assert(b.Properties.BlobType, chk.Equals, BlobTypeBlock)
}

// Ensure it's possible to generate a ListBlobs response with
// metadata, e.g., for a stub server.
func (s *StorageBlobSuite) TestMarshalBlobMetadata(c *chk.C) {
	buf, err := xml.Marshal(Blob{
		Name:       randName(5),
		Properties: BlobProperties{},
		Metadata: map[string]string{
			"lol": "baz < waz",
		},
	})
	c.Assert(err, chk.IsNil)
	c.Assert(string(buf), chk.Matches, `.*<Metadata><Lol>baz &lt; waz</Lol></Metadata>.*`)
}

func (s *StorageBlobSuite) TestGetAndSetMetadata(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	err := b.GetMetadata(nil)
	c.Assert(err, chk.IsNil)
	c.Assert(b.Metadata, chk.HasLen, 0)

	metaPut := BlobMetadata{
		"lol":      "rofl",
		"rofl_baz": "waz qux",
	}
	b.Metadata = metaPut

	err = b.SetMetadata(nil)
	c.Assert(err, chk.IsNil)

	err = b.GetMetadata(nil)
	c.Assert(err, chk.IsNil)
	c.Check(b.Metadata, chk.DeepEquals, metaPut)

	// Case munging
	metaPutUpper := BlobMetadata{
		"Lol":      "different rofl",
		"rofl_BAZ": "different waz qux",
	}
	metaExpectLower := BlobMetadata{
		"lol":      "different rofl",
		"rofl_baz": "different waz qux",
	}

	b.Metadata = metaPutUpper
	err = b.SetMetadata(nil)
	c.Assert(err, chk.IsNil)

	err = b.GetMetadata(nil)
	c.Assert(err, chk.IsNil)
	c.Check(b.Metadata, chk.DeepEquals, metaExpectLower)
}

func (s *StorageBlobSuite) TestSetMetadataWithExtraHeaders(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	b.Metadata = BlobMetadata{
		"lol":      "rofl",
		"rofl_baz": "waz qux",
	}

	options := SetBlobMetadataOptions{
		IfMatch: "incorrect-etag",
	}

	// Set with incorrect If-Match in extra headers should result in error
	err := b.SetMetadata(&options)
	c.Assert(err, chk.NotNil)

	err = b.GetProperties(nil)
	c.Assert(err, chk.IsNil)

	// Set with matching If-Match in extra headers should succeed
	options.IfMatch = b.Properties.Etag
	err = b.SetMetadata(&options)
	c.Assert(err, chk.IsNil)
}

func (s *StorageBlobSuite) TestSetBlobProperties(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	input := BlobProperties{
		CacheControl:    "private, max-age=0, no-cache",
		ContentMD5:      "oBATU+oaDduHWbVZLuzIJw==",
		ContentType:     "application/json",
		ContentEncoding: "gzip",
		ContentLanguage: "de-DE",
	}
	b.Properties = input

	err := b.SetProperties(nil)
	c.Assert(err, chk.IsNil)

	err = b.GetProperties(nil)
	c.Assert(err, chk.IsNil)

	c.Check(b.Properties.CacheControl, chk.Equals, input.CacheControl)
	c.Check(b.Properties.ContentType, chk.Equals, input.ContentType)
	c.Check(b.Properties.ContentMD5, chk.Equals, input.ContentMD5)
	c.Check(b.Properties.ContentEncoding, chk.Equals, input.ContentEncoding)
	c.Check(b.Properties.ContentLanguage, chk.Equals, input.ContentLanguage)
}

func (s *StorageBlobSuite) TestSnapshotBlob(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	snapshotTime, err := b.Snapshot(nil)
	c.Assert(err, chk.IsNil)
	c.Assert(snapshotTime, chk.NotNil)
}

func (s *StorageBlobSuite) TestSnapshotBlobWithTimeout(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	options := SnapshotOptions{
		Timeout: 0,
	}
	snapshotTime, err := b.Snapshot(&options)
	c.Assert(err, chk.IsNil)
	c.Assert(snapshotTime, chk.NotNil)
}

func (s *StorageBlobSuite) TestSnapshotBlobWithValidLease(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	// generate lease.
	currentLeaseID, err := b.AcquireLease(30, "", nil)
	c.Assert(err, chk.IsNil)

	options := SnapshotOptions{
		LeaseID: currentLeaseID,
	}
	snapshotTime, err := b.Snapshot(&options)
	c.Assert(err, chk.IsNil)
	c.Assert(snapshotTime, chk.NotNil)
}

func (s *StorageBlobSuite) TestSnapshotBlobWithInvalidLease(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	// generate lease.
	leaseID, err := b.AcquireLease(30, "", nil)
	c.Assert(err, chk.IsNil)
	c.Assert(leaseID, chk.Not(chk.Equals), "")

	options := SnapshotOptions{
		LeaseID: "GolangRocksOnAzure",
	}
	snapshotTime, err := b.Snapshot(&options)
	c.Assert(err, chk.NotNil)
	c.Assert(snapshotTime, chk.IsNil)
}

func (s *StorageBlobSuite) TestGetBlobRange(c *chk.C) {
	body := "0123456789"

	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	c.Assert(b.putSingleBlockBlob([]byte(body)), chk.IsNil)
	defer b.Delete(nil)

	cases := []struct {
		options  GetBlobRangeOptions
		expected string
	}{
		{
			options: GetBlobRangeOptions{
				Range: &BlobRange{
					Start: 0,
					End:   uint64(len(body)),
				},
			},
			expected: body,
		},
		{
			options: GetBlobRangeOptions{
				Range: &BlobRange{
					Start: 1,
					End:   3,
				},
			},
			expected: body[1 : 3+1],
		},
		{
			options: GetBlobRangeOptions{
				Range: &BlobRange{
					Start: 3,
					End:   uint64(len(body)),
				},
			},
			expected: body[3:],
		},
	}

	// Read 1-3
	for _, r := range cases {
		resp, err := b.GetRange(&(r.options))
		c.Assert(err, chk.IsNil)
		blobBody, err := ioutil.ReadAll(resp)
		c.Assert(err, chk.IsNil)

		str := string(blobBody)
		c.Assert(str, chk.Equals, r.expected)
	}
}

func (b *Blob) putSingleBlockBlob(chunk []byte) error {
	if len(chunk) > MaxBlobBlockSize {
		return fmt.Errorf("storage: provided chunk (%d bytes) cannot fit into single-block blob (max %d bytes)", len(chunk), MaxBlobBlockSize)
	}

	uri := b.Container.bsc.client.getEndpoint(blobServiceName, b.buildPath(), nil)
	headers := b.Container.bsc.client.getStandardHeaders()
	headers["x-ms-blob-type"] = string(BlobTypeBlock)
	headers["Content-Length"] = strconv.Itoa(len(chunk))

	resp, err := b.Container.bsc.client.exec(http.MethodPut, uri, headers, bytes.NewReader(chunk), b.Container.bsc.auth)
	if err != nil {
		return err
	}
	return checkRespCode(resp.statusCode, []int{http.StatusCreated})
}

func randBytes(n int) []byte {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		panic(err)
	}
	return data
}

func randName(n int) string {
	name := randString(n) + "/" + randString(n)
	return name
}

func randNameWithSpecialChars(n int) string {
	name := randString(n) + "/" + randString(n) + "-._~:?#[]@!$&'()*,;+= " + randString(n)
	return name
}
