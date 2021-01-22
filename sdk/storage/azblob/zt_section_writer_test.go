package azblob

import (
	"bytes"
	"io"

	chk "gopkg.in/check.v1"
)

func (s *aztestsSuite) TestSectionWriter(c *chk.C) {
	b := [10]byte{}
	buffer := newBytesWriter(b[:])

	section := newSectionWriter(buffer, 0, 5)
	c.Assert(section.count, chk.Equals, int64(5))
	c.Assert(section.offset, chk.Equals, int64(0))
	c.Assert(section.position, chk.Equals, int64(0))

	count, err := section.Write([]byte{1, 2, 3})
	c.Assert(err, chk.IsNil)
	c.Assert(count, chk.Equals, 3)
	c.Assert(section.position, chk.Equals, int64(3))
	c.Assert(b, chk.Equals, [10]byte{1, 2, 3, 0, 0, 0, 0, 0, 0, 0})

	count, err = section.Write([]byte{4, 5, 6})
	c.Assert(err, chk.ErrorMatches, "not enough space for all bytes")
	c.Assert(count, chk.Equals, 2)
	c.Assert(section.position, chk.Equals, int64(5))
	c.Assert(b, chk.Equals, [10]byte{1, 2, 3, 4, 5, 0, 0, 0, 0, 0})

	count, err = section.Write([]byte{6, 7, 8})
	c.Assert(err, chk.ErrorMatches, "end of section reached")
	c.Assert(count, chk.Equals, 0)
	c.Assert(section.position, chk.Equals, int64(5))
	c.Assert(b, chk.Equals, [10]byte{1, 2, 3, 4, 5, 0, 0, 0, 0, 0})

	// Intentionally create a section writer which will attempt to write
	// outside the bounds of the buffer.
	section = newSectionWriter(buffer, 5, 6)
	c.Assert(section.count, chk.Equals, int64(6))
	c.Assert(section.offset, chk.Equals, int64(5))
	c.Assert(section.position, chk.Equals, int64(0))

	count, err = section.Write([]byte{6, 7, 8})
	c.Assert(err, chk.IsNil)
	c.Assert(count, chk.Equals, 3)
	c.Assert(section.position, chk.Equals, int64(3))
	c.Assert(b, chk.Equals, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 0, 0})

	// Attempt to write past the end of the section. Since the underlying
	// buffer rejects the write it gives the same error as in the normal case.
	count, err = section.Write([]byte{9, 10, 11})
	c.Assert(err, chk.ErrorMatches, "not enough space for all bytes")
	c.Assert(count, chk.Equals, 2)
	c.Assert(section.position, chk.Equals, int64(5))
	c.Assert(b, chk.Equals, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	// Attempt to write past the end of the buffer. In this case the buffer
	// rejects the write completely since it falls completely out of bounds.
	count, err = section.Write([]byte{11, 12, 13})
	c.Assert(err, chk.ErrorMatches, "offset value is out of range")
	c.Assert(count, chk.Equals, 0)
	c.Assert(section.position, chk.Equals, int64(5))
	c.Assert(b, chk.Equals, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
}

func (s *aztestsSuite) TestSectionWriterCopySrcDestEmpty(c *chk.C) {
	input := make([]byte, 0)
	reader := bytes.NewReader(input)

	output := make([]byte, 0)
	buffer := newBytesWriter(output)
	section := newSectionWriter(buffer, 0, 0)

	count, err := io.Copy(section, reader)
	c.Assert(err, chk.IsNil)
	c.Assert(count, chk.Equals, int64(0))
}

func (s *aztestsSuite) TestSectionWriterCopyDestEmpty(c *chk.C) {
	input := make([]byte, 10)
	reader := bytes.NewReader(input)

	output := make([]byte, 0)
	buffer := newBytesWriter(output)
	section := newSectionWriter(buffer, 0, 0)

	count, err := io.Copy(section, reader)
	c.Assert(err, chk.ErrorMatches, "end of section reached")
	c.Assert(count, chk.Equals, int64(0))
}
