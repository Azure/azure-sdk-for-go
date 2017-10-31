package azblob

import (
	"context"

	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
)

type AppendBlobURLSuite struct{}

var _ = chk.Suite(&AppendBlobURLSuite{})

func (b *AppendBlobURLSuite) TestAppendBlock(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blob := container.NewAppendBlobURL(generateBlobName())

	resp, err := blob.Create(context.Background(), nil, BlobHTTPHeaders{}, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(resp.StatusCode, chk.Equals, 201)

	appendResp, err := blob.AppendBlock(context.Background(), content(1024), BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(appendResp.Response().StatusCode, chk.Equals, 201)
	c.Assert(appendResp.BlobAppendOffset(), chk.Equals, "0")
	c.Assert(appendResp.BlobCommittedBlockCount(), chk.Equals, "1")
	c.Assert(appendResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(appendResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(appendResp.ContentMD5(), chk.Not(chk.Equals), "")
	c.Assert(appendResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(appendResp.Version(), chk.Not(chk.Equals), "")
	c.Assert(appendResp.Date().IsZero(), chk.Equals, false)

	appendResp, err = blob.AppendBlock(context.Background(), content(1024), BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(appendResp.BlobAppendOffset(), chk.Equals, "1024")
	c.Assert(appendResp.BlobCommittedBlockCount(), chk.Equals, "2")
}
