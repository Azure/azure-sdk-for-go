package azblob

import (
	"context"

	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
)

type PageBlobURLSuite struct{}

var _ = chk.Suite(&PageBlobURLSuite{})

func getPageBlob(c *chk.C, container ContainerURL) PageBlobURL {
	blob := container.NewPageBlobURL(generateBlobName())
	resp, err := blob.Create(context.Background(), 4096, 0, nil, BlobHTTPHeaders{}, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.StatusCode, chk.Equals, 201)
	return blob
}

func (b *PageBlobURLSuite) TestPutGetPages(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blob := getPageBlob(c, container)

	pageRange := PageRange{Start: 0, End: 1023}
	putResp, err := blob.PutPages(context.Background(), pageRange, content(1024), BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(putResp.Response().StatusCode, chk.Equals, 201)
	c.Assert(putResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(putResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(putResp.ContentMD5(), chk.Not(chk.Equals), "")
	c.Assert(putResp.BlobSequenceNumber(), chk.Equals, int32(0))
	c.Assert(putResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(putResp.Version(), chk.Not(chk.Equals), "")
	c.Assert(putResp.Date().IsZero(), chk.Equals, false)

	pageList, err := blob.GetPageRanges(context.Background(), pageRange, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(pageList.Response().StatusCode, chk.Equals, 200)
	c.Assert(pageList.LastModified().IsZero(), chk.Equals, false)
	c.Assert(pageList.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(pageList.BlobContentLength(), chk.Equals, int64(4096))
	c.Assert(pageList.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(pageList.Version(), chk.Not(chk.Equals), "")
	c.Assert(pageList.Date().IsZero(), chk.Equals, false)
	c.Assert(pageList.PageRange, chk.HasLen, 1)
	c.Assert(pageList.PageRange[0], chk.DeepEquals, pageRange)
}

func (b *PageBlobURLSuite) TestClearDiffPages(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blob := getPageBlob(c, container)
	_, err := blob.PutPages(context.Background(), PageRange{Start: 0, End: 2047}, content(2048), BlobAccessConditions{})
	c.Assert(err, chk.IsNil)

	snapshotResp, err := blob.CreateSnapshot(context.Background(), nil, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)

	_, err = blob.PutPages(context.Background(), PageRange{Start: 2048, End: 4095}, content(2048), BlobAccessConditions{})
	c.Assert(err, chk.IsNil)

	pageList, err := blob.GetPageRangesDiff(context.Background(), PageRange{Start: 0, End: 4095}, snapshotResp.Snapshot(), BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(pageList.PageRange, chk.HasLen, 1)
	c.Assert(pageList.PageRange[0].Start, chk.Equals, int32(2048))
	c.Assert(pageList.PageRange[0].End, chk.Equals, int32(4095))

	clearResp, err := blob.ClearPages(context.Background(), PageRange{Start: 2048, End: 4095}, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(clearResp.Response().StatusCode, chk.Equals, 201)

	pageList, err = blob.GetPageRangesDiff(context.Background(), PageRange{Start: 0, End: 4095}, snapshotResp.Snapshot(), BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(pageList.PageRange, chk.HasLen, 0)
}

func (b *PageBlobURLSuite) TestIncrementalCopy(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)
	_, err := container.SetPermissions(context.Background(), PublicAccessBlob, nil, ContainerAccessConditions{})
	c.Assert(err, chk.IsNil)

	srcBlob := getPageBlob(c, container)
	_, err = srcBlob.PutPages(context.Background(), PageRange{Start: 0, End: 1023}, content(1024), BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	snapshotResp, err := srcBlob.CreateSnapshot(context.Background(), nil, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)

	dstBlob := container.NewPageBlobURL(generateBlobName())

	resp, err := dstBlob.StartIncrementalCopy(context.Background(), srcBlob.URL(), snapshotResp.Snapshot(), nil, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 202)
	c.Assert(resp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(resp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(resp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(resp.Version(), chk.Not(chk.Equals), "")
	c.Assert(resp.Date().IsZero(), chk.Equals, false)
	c.Assert(resp.CopyID(), chk.Not(chk.Equals), "")
	c.Assert(resp.CopyStatus(), chk.Not(chk.Equals), "")
}

func (b *PageBlobURLSuite) TestResizePageBlob(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blob := getPageBlob(c, container)
	resp, err := blob.Resize(context.Background(), 2048, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 200)

	resp, err = blob.Resize(context.Background(), 8192, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 200)
}

/*
NOTE: We've decided not to expose this from the convenience layer.
func (b *PageBlobURLSuite) TestPageSequenceNumbers(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blob := getPageBlob(c, container)

	resp, err := blob.SetSequenceNumber(context.Background(), SequenceNumberActionIncrement, 0, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 200)

	resp, err = blob.SetSequenceNumber(context.Background(), SequenceNumberActionMax, 7, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 200)

	resp, err = blob.SetSequenceNumber(context.Background(), SequenceNumberActionUpdate, 11, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.Response().StatusCode, chk.Equals, 200)
}
*/
