package storage

import (
	"fmt"
	"io/ioutil"

	chk "gopkg.in/check.v1"
)

type PageBlobSuite struct{}

var _ = chk.Suite(&PageBlobSuite{})

func (s *PageBlobSuite) TestPutPageBlob(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	size := int64(10 * 1024 * 1024)
	c.Assert(b.PutPageBlob(size, nil), chk.IsNil)

	// Verify
	err := b.GetProperties()
	c.Assert(err, chk.IsNil)
	c.Assert(b.Properties.ContentLength, chk.Equals, size)
	c.Assert(b.Properties.BlobType, chk.Equals, BlobTypePage)
}

func (s *PageBlobSuite) TestPutPagesUpdate(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	size := int64(10 * 1024 * 1024) // larger than we'll use
	c.Assert(b.PutPageBlob(size, nil), chk.IsNil)

	chunk1 := []byte(randString(1024))
	chunk2 := []byte(randString(512))

	// Append chunks
	c.Assert(b.PutPage(0, int64(len(chunk1)-1), PageWriteTypeUpdate, chunk1, nil), chk.IsNil)
	c.Assert(b.PutPage(int64(len(chunk1)), int64(len(chunk1)+len(chunk2)-1), PageWriteTypeUpdate, chunk2, nil), chk.IsNil)

	// Verify contents
	out, err := b.GetRange(fmt.Sprintf("%v-%v", 0, len(chunk1)+len(chunk2)-1), nil)
	c.Assert(err, chk.IsNil)
	defer out.Close()
	blobContents, err := ioutil.ReadAll(out)
	c.Assert(err, chk.IsNil)
	c.Assert(blobContents, chk.DeepEquals, append(chunk1, chunk2...))

	// Overwrite first half of chunk1
	chunk0 := []byte(randString(512))
	c.Assert(b.PutPage(0, int64(len(chunk0)-1), PageWriteTypeUpdate, chunk0, nil), chk.IsNil)

	// Verify contents
	out, err = b.GetRange(fmt.Sprintf("%v-%v", 0, len(chunk1)+len(chunk2)-1), nil)
	c.Assert(err, chk.IsNil)
	defer out.Close()
	blobContents, err = ioutil.ReadAll(out)
	c.Assert(err, chk.IsNil)
	c.Assert(blobContents, chk.DeepEquals, append(append(chunk0, chunk1[512:]...), chunk2...))
}

func (s *PageBlobSuite) TestPutPagesClear(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	size := int64(10 * 1024 * 1024) // larger than we'll use
	c.Assert(b.PutPageBlob(size, nil), chk.IsNil)

	// Put 0-2047
	chunk := []byte(randString(2048))
	c.Assert(b.PutPage(0, 2047, PageWriteTypeUpdate, chunk, nil), chk.IsNil)

	// Clear 512-1023
	c.Assert(b.PutPage(512, 1023, PageWriteTypeClear, nil, nil), chk.IsNil)

	// Verify contents
	out, err := b.GetRange("0-2047", nil)
	c.Assert(err, chk.IsNil)
	contents, err := ioutil.ReadAll(out)
	c.Assert(err, chk.IsNil)
	defer out.Close()
	c.Assert(contents, chk.DeepEquals, append(append(chunk[:512], make([]byte, 512)...), chunk[1024:]...))
}

func (s *PageBlobSuite) TestGetPageRanges(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(), chk.IsNil)
	defer cnt.Delete()

	size := int64(10 * 1024 * 1024) // larger than we'll use
	c.Assert(b.PutPageBlob(size, nil), chk.IsNil)

	// Get page ranges on empty blob
	out, err := b.GetPageRanges()
	c.Assert(err, chk.IsNil)
	c.Assert(len(out.PageList), chk.Equals, 0)

	// Add 0-512 page
	c.Assert(b.PutPage(0, 511, PageWriteTypeUpdate, []byte(randString(512)), nil), chk.IsNil)

	out, err = b.GetPageRanges()
	c.Assert(err, chk.IsNil)
	c.Assert(len(out.PageList), chk.Equals, 1)

	// Add 1024-2048
	c.Assert(b.PutPage(1024, 2047, PageWriteTypeUpdate, []byte(randString(1024)), nil), chk.IsNil)

	out, err = b.GetPageRanges()
	c.Assert(err, chk.IsNil)
	c.Assert(len(out.PageList), chk.Equals, 2)
}
