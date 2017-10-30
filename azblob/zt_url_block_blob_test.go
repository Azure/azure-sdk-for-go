package azblob

import (
	"context"
	"encoding/base64"
	"fmt"

	chk "gopkg.in/check.v1" // go get gopkg.in/check.v1
)

type BlockBlobURLSuite struct{}

var _ = chk.Suite(&BlockBlobURLSuite{})

func (b *BlockBlobURLSuite) TestPutGetBlocks(c *chk.C) {
	container := getContainer(c)
	defer delContainer(c, container)

	blob := container.NewBlockBlobURL(generateBlobName())

	blockID := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%6d", 0)))

	putResp, err := blob.PutBlock(context.Background(), blockID, content(1024), LeaseAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(putResp.Response().StatusCode, chk.Equals, 201)
	c.Assert(putResp.ContentMD5(), chk.Not(chk.Equals), "")
	c.Assert(putResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(putResp.Version(), chk.Not(chk.Equals), "")
	c.Assert(putResp.Date().IsZero(), chk.Equals, false)

	blockList, err := blob.GetBlockList(context.Background(), BlockListAll, LeaseAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(blockList.Response().StatusCode, chk.Equals, 200)
	c.Assert(blockList.LastModified().IsZero(), chk.Equals, true)
	c.Assert(blockList.ETag(), chk.Equals, ETagNone)
	c.Assert(blockList.ContentType(), chk.Not(chk.Equals), "")
	c.Assert(blockList.BlobContentLength(), chk.Equals, int64(-1))
	c.Assert(blockList.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(blockList.Version(), chk.Not(chk.Equals), "")
	c.Assert(blockList.Date().IsZero(), chk.Equals, false)
	c.Assert(blockList.CommittedBlocks, chk.HasLen, 0)
	c.Assert(blockList.UncommittedBlocks, chk.HasLen, 1)

	listResp, err := blob.PutBlockList(context.Background(), []string{blockID}, nil, BlobHTTPHeaders{}, BlobAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(listResp.Response().StatusCode, chk.Equals, 201)
	c.Assert(listResp.LastModified().IsZero(), chk.Equals, false)
	c.Assert(listResp.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(listResp.ContentMD5(), chk.Not(chk.Equals), "")
	c.Assert(listResp.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(listResp.Version(), chk.Not(chk.Equals), "")
	c.Assert(listResp.Date().IsZero(), chk.Equals, false)

	blockList, err = blob.GetBlockList(context.Background(), BlockListAll, LeaseAccessConditions{})
	c.Assert(err, chk.IsNil)
	c.Assert(blockList.Response().StatusCode, chk.Equals, 200)
	c.Assert(blockList.LastModified().IsZero(), chk.Equals, false)
	c.Assert(blockList.ETag(), chk.Not(chk.Equals), ETagNone)
	c.Assert(blockList.ContentType(), chk.Not(chk.Equals), "")
	c.Assert(blockList.BlobContentLength(), chk.Equals, int64(1024))
	c.Assert(blockList.RequestID(), chk.Not(chk.Equals), "")
	c.Assert(blockList.Version(), chk.Not(chk.Equals), "")
	c.Assert(blockList.Date().IsZero(), chk.Equals, false)
	c.Assert(blockList.CommittedBlocks, chk.HasLen, 1)
	c.Assert(blockList.UncommittedBlocks, chk.HasLen, 0)
}
