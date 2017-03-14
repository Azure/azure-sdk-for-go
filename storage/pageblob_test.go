package storage

import (
	"bytes"
	"io/ioutil"

	chk "gopkg.in/check.v1"
)

type PageBlobSuite struct{}

var _ = chk.Suite(&PageBlobSuite{})

func (s *PageBlobSuite) TestPutPageBlob(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	size := int64(10 * 1024 * 1024)
	b.Properties.ContentLength = size
	c.Assert(b.PutPageBlob(nil), chk.IsNil)

	// Verify
	err := b.GetProperties(nil)
	c.Assert(err, chk.IsNil)
	c.Assert(b.Properties.ContentLength, chk.Equals, size)
	c.Assert(b.Properties.BlobType, chk.Equals, BlobTypePage)
}

func (s *PageBlobSuite) TestPutPagesUpdate(c *chk.C) {
	cli := getBlobClient(c)
	cnt := cli.GetContainerReference(randContainer())
	b := cnt.GetBlobReference(randName(5))
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	size := int64(10 * 1024 * 1024) // larger than we'll use
	b.Properties.ContentLength = size
	c.Assert(b.PutPageBlob(nil), chk.IsNil)

	chunk1 := []byte(randString(1024))
	chunk2 := []byte(randString(512))

	// Append chunks
	blobRange := BlobRange{
		End: uint64(len(chunk1) - 1),
	}
	c.Assert(b.WriteRange(blobRange, bytes.NewReader(chunk1), nil), chk.IsNil)
	blobRange.Start = uint64(len(chunk1))
	blobRange.End = uint64(len(chunk1) + len(chunk2) - 1)
	c.Assert(b.WriteRange(blobRange, bytes.NewReader(chunk2), nil), chk.IsNil)

	// Verify contents
	options := GetBlobRangeOptions{
		Range: &BlobRange{
			End: uint64(len(chunk1) + len(chunk2) - 1),
		},
	}
	out, err := b.GetRange(&options)
	c.Assert(err, chk.IsNil)
	defer out.Close()
	blobContents, err := ioutil.ReadAll(out)
	c.Assert(err, chk.IsNil)
	c.Assert(blobContents, chk.DeepEquals, append(chunk1, chunk2...))

	// Overwrite first half of chunk1
	chunk0 := []byte(randString(512))
	blobRange.Start = 0
	blobRange.End = uint64(len(chunk0) - 1)
	c.Assert(b.WriteRange(blobRange, bytes.NewReader(chunk0), nil), chk.IsNil)

	// Verify contents
	out, err = b.GetRange(&options)
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
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	size := int64(10 * 1024 * 1024) // larger than we'll use
	b.Properties.ContentLength = size
	c.Assert(b.PutPageBlob(nil), chk.IsNil)

	// Put 0-2047
	chunk := []byte(randString(2048))
	blobRange := BlobRange{
		End: 2047,
	}
	c.Assert(b.WriteRange(blobRange, bytes.NewReader(chunk), nil), chk.IsNil)

	// Clear 512-1023
	blobRange.Start = 512
	blobRange.End = 1023
	c.Assert(b.ClearRange(blobRange, nil), chk.IsNil)

	// Verify contents
	options := GetBlobRangeOptions{
		Range: &BlobRange{
			Start: 0,
			End:   2047,
		},
	}
	out, err := b.GetRange(&options)
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
	c.Assert(cnt.Create(nil), chk.IsNil)
	defer cnt.Delete(nil)

	size := int64(10 * 1024 * 1024) // larger than we'll use
	b.Properties.ContentLength = size
	c.Assert(b.PutPageBlob(nil), chk.IsNil)

	// Get page ranges on empty blob
	out, err := b.GetPageRanges(nil)
	c.Assert(err, chk.IsNil)
	c.Assert(len(out.PageList), chk.Equals, 0)

	// Add 0-512 page
	blobRange := BlobRange{
		End: 511,
	}
	c.Assert(b.WriteRange(blobRange, bytes.NewReader([]byte(randString(512))), nil), chk.IsNil)

	out, err = b.GetPageRanges(nil)
	c.Assert(err, chk.IsNil)
	c.Assert(len(out.PageList), chk.Equals, 1)

	// Add 1024-2048
	blobRange.Start = 1024
	blobRange.End = 2047
	c.Assert(b.WriteRange(blobRange, bytes.NewReader([]byte(randString(1024))), nil), chk.IsNil)

	out, err = b.GetPageRanges(nil)
	c.Assert(err, chk.IsNil)
	c.Assert(len(out.PageList), chk.Equals, 2)
}
