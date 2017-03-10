package storage

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"

	chk "gopkg.in/check.v1"
)

type BlockBlobSuite struct{}

var _ = chk.Suite(&BlockBlobSuite{})

func (s *BlockBlobSuite) TestCreateBlockBlobFromReader(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	data := randBytes(8888)
	b.Properties.ContentLength = int64(len(data))
	c.Assert(b.CreateBlockBlobFromReader(bytes.NewReader(data), nil), chk.IsNil)

	resp, err := b.Get(nil)
	c.Assert(err, chk.IsNil)
	gotData, err := ioutil.ReadAll(resp)
	defer resp.Close()

	c.Assert(err, chk.IsNil)
	c.Assert(gotData, chk.DeepEquals, data)
}

func (s *BlockBlobSuite) TestCreateBlockBlobFromReaderWithShortData(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	data := randBytes(8888)
	b.Properties.ContentLength = 9999
	err := b.CreateBlockBlobFromReader(bytes.NewReader(data), nil)
	c.Assert(err, chk.NotNil)

	_, err = b.Get(nil)
	// Upload was incomplete: blob should not have been created.
	c.Assert(err, chk.NotNil)
}

func (s *BlockBlobSuite) TestPutBlock(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	chunk := []byte(randString(1024))
	blockID := base64.StdEncoding.EncodeToString([]byte("lol"))
	c.Assert(b.PutBlock(blockID, chunk, nil), chk.IsNil)
}

func (s *BlockBlobSuite) TestGetBlockList_PutBlockList(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	chunk := []byte(randString(1024))
	blockID := base64.StdEncoding.EncodeToString([]byte("lol"))

	// Put one block
	c.Assert(b.PutBlock(blockID, chunk, nil), chk.IsNil)
	defer b.Delete(nil)

	// Get committed blocks
	committed, err := b.GetBlockList(BlockListTypeCommitted, nil)
	c.Assert(err, chk.IsNil)

	if len(committed.CommittedBlocks) > 0 {
		c.Fatal("There are committed blocks")
	}

	// Get uncommitted blocks
	uncommitted, err := b.GetBlockList(BlockListTypeUncommitted, nil)
	c.Assert(err, chk.IsNil)

	c.Assert(len(uncommitted.UncommittedBlocks), chk.Equals, 1)
	// Commit block list
	c.Assert(b.PutBlockList([]Block{{blockID, BlockStatusUncommitted}}, nil), chk.IsNil)

	// Get all blocks
	all, err := b.GetBlockList(BlockListTypeAll, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(len(all.CommittedBlocks), chk.Equals, 1)
	c.Assert(len(all.UncommittedBlocks), chk.Equals, 0)

	// Verify the block
	thatBlock := all.CommittedBlocks[0]
	c.Assert(thatBlock.Name, chk.Equals, blockID)
	c.Assert(thatBlock.Size, chk.Equals, int64(len(chunk)))
}

func (s *BlockBlobSuite) TestCreateBlockBlob(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	c.Assert(b.CreateBlockBlob(nil), chk.IsNil)

	// Verify
	blocks, err := b.GetBlockList(BlockListTypeAll, nil)
	c.Assert(err, chk.IsNil)
	c.Assert(len(blocks.CommittedBlocks), chk.Equals, 0)
	c.Assert(len(blocks.UncommittedBlocks), chk.Equals, 0)
}

func (s *BlockBlobSuite) TestPutEmptyBlockBlob(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	c.Assert(b.putSingleBlockBlob([]byte{}), chk.IsNil)

	err := b.GetProperties(nil)
	c.Assert(err, chk.IsNil)
	c.Assert(b.Properties.ContentLength, chk.Not(chk.Equals), 0)
}
