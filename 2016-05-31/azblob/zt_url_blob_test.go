package azblob

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
)

type BlobURLSuite struct{}

var _ = chk.Suite(&BlobURLSuite{})

func getBlob(c *chk.C, container ContainerURL) BlockBlobURL {
	blob := container.NewBlockBlobURL(generateBlobName())
	putResp, err := blob.PutBlob(context.Background(), nil, BlobHTTPHeaders{}, nil, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(putResp.Response().StatusCode, chk.Equals, 201)
	return blob
}

func generateBlobName() string {
	return newUUID().String()
}

func generateBlobNameWithPrefix(prefix string) string {
	return prefix + newUUID().String()
}

func content(n int) *bytes.Reader {
	r, _ := contentAndData(n)
	return r
}

func contentAndData(n int) (*bytes.Reader, []byte) {
	data := make([]byte, n, n)
	for i := 0; i < n; i++ {
		data[i] = byte(i)
	}
	return bytes.NewReader(data), data
}

func (b *BlobURLSuite) TestCreateDelete(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blob := container.NewBlockBlobURL(generateBlobName())

	putResp, err := blob.PutBlob(context.Background(), nil, BlobHTTPHeaders{}, nil, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(putResp.Response().StatusCode, chk.Equals, 201)
	c.Assert(putResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(putResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(putResp.ContentMD5(), chk.Not(chk.Equals), "")
	c.Assert(putResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(putResp.Version(), chk.Not(chk.Equals), "")
	c.Assert(putResp.Date().IsZero(), chk.Equals, false)

	delResp, err := blob.Delete(context.Background(), DeleteSnapshotsOptionNone, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(delResp.Response().StatusCode, chk.Equals, 202)
	c.Assert(delResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(delResp.Version(), chk.Not(chk.Equals), "")
	c.Assert(delResp.Date().IsZero(), chk.Equals, false)
}

func (b *BlobURLSuite) TestGetSetProperties(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blob := getBlob(c, container)

	properties := BlobHTTPHeaders{
		ContentType:     "mytype",
		ContentLanguage: "martian",
	}
	setResp, err := blob.SetProperties(context.Background(), properties, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(setResp.Response().StatusCode, chk.Equals, 200)
	c.Assert(setResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(setResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(setResp.BlobSequenceNumber(), chk.Equals, "")
	c.Assert(setResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(setResp.Version(), chk.Not(chk.Equals), "")
	c.Assert(setResp.Date().IsZero(), chk.Equals, false)

	/*getResp, err := blob.GetProperties(context.Background(), BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(getResp.Response().StatusCode, chk.Equals, 200)
	c.Assert(getResp.ContentType(), chk.Equals, properties.ContentType)
	c.Assert(getResp.ContentLanguage(), chk.Equals, properties.ContentLanguage)
	verifyDateResp(c, getResp.LastModified, false)
	c.Assert(getResp.BlobType(), chk.Not(chk.Equals), "")
	verifyDateResp(c, getResp.CopyCompletionTime, true)
	c.Assert(getResp.CopyStatusDescription(), chk.Equals, "")
	c.Assert(getResp.CopyID(), chk.Equals, "")
	c.Assert(getResp.CopyProgress(), chk.Equals, "")
	c.Assert(getResp.CopySource(), chk.Equals, "")
	c.Assert(getResp.CopyStatus().IsZero(), chk.Equals, true)
	c.Assert(getResp.IsIncrementalCopy(), chk.Equals, "")
	c.Assert(getResp.LeaseDuration().IsZero(), chk.Equals, true)
	c.Assert(getResp.LeaseState(), chk.Equals, LeaseStateAvailable)
	c.Assert(getResp.LeaseStatus(), chk.Equals, LeaseStatusUnlocked)
	c.Assert(getResp.ContentLength(), chk.Not(chk.Equals), "")
	c.Assert(getResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(getResp.ContentMD5(), chk.Equals, "")
	c.Assert(getResp.ContentEncoding(), chk.Equals, "")
	c.Assert(getResp.ContentDisposition(), chk.Equals, "")
	c.Assert(getResp.CacheControl(), chk.Equals, "")
	c.Assert(getResp.BlobSequenceNumber(), chk.Equals, "")
	c.Assert(getResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(getResp.Version(), chk.Not(chk.Equals), "")
	verifyDateResp(c, getResp.Date, false)
	c.Assert(getResp.AcceptRanges(), chk.Equals, "bytes")
	c.Assert(getResp.BlobCommittedBlockCount(), chk.Equals, "")
	*/
}

func (b *BlobURLSuite) TestGetSetMetadata(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blob := getBlob(c, container)

	metadata := Metadata{
		"foo": "foovalue",
		"bar": "barvalue",
	}
	setResp, err := blob.SetMetadata(context.Background(), metadata, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(setResp.Response().StatusCode, chk.Equals, 200)
	c.Assert(setResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(setResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(setResp.Version(), chk.Not(chk.Equals), "")
	c.Assert(setResp.Date().IsZero(), chk.Equals, false)

	getResp, err := blob.GetPropertiesAndMetadata(context.Background(), BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(getResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(getResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(getResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(getResp.Version(), chk.Not(chk.Equals), "")
	c.Assert(getResp.Date().IsZero(), chk.Equals, false)
	md := getResp.NewMetadata()
	c.Assert(md, chk.DeepEquals, metadata)
}

func (b *BlobURLSuite) TestCopy(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	sourceBlob := getBlob(c, container)
	_, err := sourceBlob.PutBlob(context.Background(), content(2048), BlobHTTPHeaders{}, nil, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)

	destBlob := getBlob(c, container)
	copyResp, err := destBlob.StartCopy(context.Background(), sourceBlob.URL(), nil, BlobAccessConditions{}, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(copyResp.Response().StatusCode, chk.Equals, 202)
	c.Assert(copyResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(copyResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(copyResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(copyResp.Version(), chk.Not(chk.Equals), "")
	c.Assert(copyResp.Date().IsZero(), chk.Equals, false)
	c.Assert(copyResp.CopyID(), chk.Not(chk.Equals), "")
	c.Assert(copyResp.CopyStatus(), chk.Not(chk.Equals), "")

	abortResp, err := destBlob.AbortCopy(context.Background(), copyResp.CopyID(), LeaseAccessConditions{})
	// small copy completes before we have time to abort so check for failure case
	c.Assert(err, chk.NotNil)
	c.Assert(abortResp, chk.IsNil)
	se, ok := err.(StorageError)
	c.Assert(ok, chk.Equals, true)
	c.Assert(se.Response().StatusCode, chk.Equals, http.StatusConflict)
}

func (b *BlobURLSuite) TestSnapshot(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blob := getBlob(c, container)

	resp, err := blob.CreateSnapshot(context.Background(), nil, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 201)
	c.Assert(resp.Snapshot().IsZero(), chk.Equals, false)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")
	c.Assert(resp.Date().IsZero(), chk.Equals, false)

	blobs, err := container.ListBlobs(context.Background(), Marker{}, ListBlobsOptions{Details: BlobListingDetails{Snapshots: true}})
	c.Assert(err, chk.IsNil)
	c.Assert(blobs.Blobs.Blob, chk.HasLen, 2)

	_, err = blob.Delete(context.Background(), DeleteSnapshotsOptionOnly, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)

	blobs, err = container.ListBlobs(context.Background(), Marker{}, ListBlobsOptions{Details: BlobListingDetails{Snapshots: true}})
	c.Assert(err, chk.IsNil)
	c.Assert(blobs.Blobs.Blob, chk.HasLen, 1)
}

func (b *BlobURLSuite) TestLeaseAcquireRelease(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blob := getBlob(c, container)

	leaseID := newUUID().String()
	resp, err := blob.AcquireLease(context.Background(), leaseID, 15, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 201)
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.LeaseID(), chk.Equals, leaseID)
	c.Assert(resp.LeaseTime(), chk.Equals, int32(-1))
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")

	resp, err = blob.ReleaseLease(context.Background(), leaseID, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 200)
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.LeaseID(), chk.Equals, "")
	c.Assert(resp.LeaseTime(), chk.Equals, int32(-1))
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")
}

func (b *BlobURLSuite) TestLeaseRenewChangeBreak(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blob := getBlob(c, container)

	leaseID := newUUID().String()
	resp, err := blob.AcquireLease(context.Background(), leaseID, 15, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)

	newID := newUUID().String()
	resp, err = blob.ChangeLease(context.Background(), leaseID, newID, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 200)
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.LeaseID(), chk.Equals, newID)
	c.Assert(resp.LeaseTime(), chk.Equals, int32(-1))
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")

	resp, err = blob.RenewLease(context.Background(), newID, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 200)
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.LeaseID(), chk.Equals, newID)
	c.Assert(resp.LeaseTime(), chk.Equals, int32(-1))
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")

	resp, err = blob.BreakLease(context.Background(), newID, 5, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 202)
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.LeaseID(), chk.Equals, "")
	c.Assert(resp.LeaseTime(), chk.Not(chk.Equals), "")
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")

	resp, err = blob.ReleaseLease(context.Background(), newID, HTTPAccessConditions{})
	c.Assert(err, chk.IsNil)
}

func (b *BlobURLSuite) TestGetBlobRange(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blob := getBlob(c, container)
	contentR, contentD := contentAndData(2048)
	_, err := blob.PutBlob(context.Background(), contentR, BlobHTTPHeaders{}, nil, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)

	resp, err := blob.GetBlob(context.Background(), BlobRange{Offset: 0, Count: 1024}, BlobAccessConditions{}, false)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.ContentLength(), chk.Equals, int64(1024))

	download, err := ioutil.ReadAll(resp.Response().Body)
	c.Assert(err, chk.IsNil)
	c.Assert(download, chk.DeepEquals, contentD[:1024])

	resp, err = blob.GetBlob(context.Background(), BlobRange{Offset: 1024}, BlobAccessConditions{}, false)
	c.Assert(err, chk.IsNil)
	c.Assert(resp.ContentLength(), chk.Equals, int64(1024))

	download, err = ioutil.ReadAll(resp.Response().Body)
	c.Assert(err, chk.IsNil)
	c.Assert(download, chk.DeepEquals, contentD[1024:])

	c.Assert(resp.AcceptRanges(), chk.Equals, "bytes")
	c.Assert(resp.BlobCommittedBlockCount(), chk.Equals, "")
	c.Assert(resp.BlobContentMD5(), chk.Not(chk.Equals), "")
	c.Assert(resp.BlobSequenceNumber(), chk.Equals, "")
	c.Assert(resp.BlobType(), chk.Equals, BlobBlockBlob)
	c.Assert(resp.CacheControl(), chk.Equals, "")
	c.Assert(resp.ContentDisposition(), chk.Equals, "")
	c.Assert(resp.ContentEncoding(), chk.Equals, "")
	c.Assert(resp.ContentMD5(), chk.Equals, "")
	c.Assert(resp.ContentRange(), chk.Equals, "bytes 1024-2047/2048")
	c.Assert(resp.ContentType(), chk.Equals, "application/octet-stream")
	c.Assert(resp.CopyCompletionTime().IsZero(), chk.Equals, true)
	c.Assert(resp.CopyID(), chk.Equals, "")
	c.Assert(resp.CopyProgress(), chk.Equals, "")
	c.Assert(resp.CopySource(), chk.Equals, "")
	c.Assert(resp.CopyStatus(), chk.Equals, CopyStatusNone)
	c.Assert(resp.CopyStatusDescription(), chk.Equals, "")
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.LeaseDuration(), chk.Equals, LeaseDurationNone)
	c.Assert(resp.LeaseState(), chk.Equals, LeaseStateAvailable)
	c.Assert(resp.LeaseStatus(), chk.Equals, LeaseStatusUnlocked)
	c.Assert(resp.NewMetadata(), chk.DeepEquals, Metadata{})
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Response().StatusCode, chk.Equals, 206)
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")
}
