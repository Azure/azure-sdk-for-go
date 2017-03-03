package storage

import (
	"fmt"
	"io/ioutil"

	chk "gopkg.in/check.v1"
)

type AppendBlobSuite struct{}

var _ = chk.Suite(&AppendBlobSuite{})

func (s *AppendBlobSuite) TestPutAppendBlob(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	c.Assert(b.PutAppendBlob(nil), chk.IsNil)

	// Verify
	err := b.GetProperties()
	c.Assert(err, chk.IsNil)
	c.Assert(b.Properties.ContentLength, chk.Equals, int64(0))
	c.Assert(b.Properties.BlobType, chk.Equals, BlobTypeAppend)
}

func (s *AppendBlobSuite) TestPutAppendBlobAppendBlocks(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	c.Assert(b.PutAppendBlob(nil), chk.IsNil)

	chunk1 := []byte(randString(1024))
	chunk2 := []byte(randString(512))

	// Append first block
	c.Assert(b.AppendBlock(chunk1, nil), chk.IsNil)

	// Verify contents
	out, err := b.GetRange(fmt.Sprintf("%v-%v", 0, len(chunk1)-1), nil)
	c.Assert(err, chk.IsNil)
	defer out.Close()
	blobContents, err := ioutil.ReadAll(out)
	c.Assert(err, chk.IsNil)
	c.Assert(blobContents, chk.DeepEquals, chunk1)

	// Append second block
	c.Assert(b.AppendBlock(chunk2, nil), chk.IsNil)

	// Verify contents
	out, err = b.GetRange(fmt.Sprintf("%v-%v", 0, len(chunk1)+len(chunk2)-1), nil)
	c.Assert(err, chk.IsNil)
	defer out.Close()
	blobContents, err = ioutil.ReadAll(out)
	c.Assert(err, chk.IsNil)
	c.Assert(blobContents, chk.DeepEquals, append(chunk1, chunk2...))
}

func (s *StorageBlobSuite) TestPutAppendBlobSpecialChars(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	c.Assert(b.PutAppendBlob(nil), chk.IsNil)

	// Verify metadata
	err := b.GetProperties()
	c.Assert(err, chk.IsNil)
	c.Assert(b.Properties.ContentLength, chk.Equals, int64(0))
	c.Assert(b.Properties.BlobType, chk.Equals, BlobTypeAppend)

	chunk1 := []byte(randString(1024))
	chunk2 := []byte(randString(512))

	// Append first block
	c.Assert(b.AppendBlock(chunk1, nil), chk.IsNil)

	// Verify contents
	out, err := b.GetRange(fmt.Sprintf("%v-%v", 0, len(chunk1)-1), nil)
	c.Assert(err, chk.IsNil)
	defer out.Close()
	blobContents, err := ioutil.ReadAll(out)
	c.Assert(err, chk.IsNil)
	c.Assert(blobContents, chk.DeepEquals, chunk1)

	// Append second block
	c.Assert(b.AppendBlock(chunk2, nil), chk.IsNil)

	// Verify contents
	out, err = b.GetRange(fmt.Sprintf("%v-%v", 0, len(chunk1)+len(chunk2)-1), nil)
	c.Assert(err, chk.IsNil)
	defer out.Close()
	blobContents, err = ioutil.ReadAll(out)
	c.Assert(err, chk.IsNil)
	c.Assert(blobContents, chk.DeepEquals, append(chunk1, chunk2...))
}
