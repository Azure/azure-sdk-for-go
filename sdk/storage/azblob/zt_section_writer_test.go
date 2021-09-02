// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
)

//nolint
func (s *azblobUnrecordedTestSuite) TestSectionWriter() {
	_assert := assert.New(s.T())
	b := [10]byte{}
	buffer := newBytesWriter(b[:])

	section := newSectionWriter(buffer, 0, 5)
	_assert.Equal(section.count, int64(5))
	_assert.Equal(section.offset, int64(0))
	_assert.Equal(section.position, int64(0))

	count, err := section.Write([]byte{1, 2, 3})
	_assert.Nil(err)
	_assert.Equal(count, 3)
	_assert.Equal(section.position, int64(3))
	_assert.Equal(b, [10]byte{1, 2, 3, 0, 0, 0, 0, 0, 0, 0})

	count, err = section.Write([]byte{4, 5, 6})
	_assert.Contains(err.Error(), "not enough space for all bytes")
	_assert.Equal(count, 2)
	_assert.Equal(section.position, int64(5))
	_assert.Equal(b, [10]byte{1, 2, 3, 4, 5, 0, 0, 0, 0, 0})

	count, err = section.Write([]byte{6, 7, 8})
	_assert.Contains(err.Error(), "end of section reached")
	_assert.Equal(count, 0)
	_assert.Equal(section.position, int64(5))
	_assert.Equal(b, [10]byte{1, 2, 3, 4, 5, 0, 0, 0, 0, 0})

	// Intentionally create a section writer which will attempt to write
	// outside the bounds of the buffer.
	section = newSectionWriter(buffer, 5, 6)
	_assert.Equal(section.count, int64(6))
	_assert.Equal(section.offset, int64(5))
	_assert.Equal(section.position, int64(0))

	count, err = section.Write([]byte{6, 7, 8})
	_assert.Nil(err)
	_assert.Equal(count, 3)
	_assert.Equal(section.position, int64(3))
	_assert.Equal(b, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 0, 0})

	// Attempt to write past the end of the section. Since the underlying
	// buffer rejects the write it gives the same error as in the normal case.
	count, err = section.Write([]byte{9, 10, 11})
	_assert.Contains(err.Error(), "not enough space for all bytes")
	_assert.Equal(count, 2)
	_assert.Equal(section.position, int64(5))
	_assert.Equal(b, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	// Attempt to write past the end of the buffer. In this case the buffer
	// rejects the write completely since it falls completely out of bounds.
	count, err = section.Write([]byte{11, 12, 13})
	_assert.Contains(err.Error(), "offset value is out of range")
	_assert.Equal(count, 0)
	_assert.Equal(section.position, int64(5))
	_assert.Equal(b, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
}

//nolint
func (s *azblobUnrecordedTestSuite) TestSectionWriterCopySrcDestEmpty() {
	_assert := assert.New(s.T())
	input := make([]byte, 0)
	reader := bytes.NewReader(input)

	output := make([]byte, 0)
	buffer := newBytesWriter(output)
	section := newSectionWriter(buffer, 0, 0)

	count, err := io.Copy(section, reader)
	_assert.Nil(err)
	_assert.Equal(count, int64(0))
}

//nolint
func (s *azblobUnrecordedTestSuite) TestSectionWriterCopyDestEmpty() {
	_assert := assert.New(s.T())
	input := make([]byte, 10)
	reader := bytes.NewReader(input)

	output := make([]byte, 0)
	buffer := newBytesWriter(output)
	section := newSectionWriter(buffer, 0, 0)

	count, err := io.Copy(section, reader)
	_assert.Contains(err.Error(), "end of section reached")
	_assert.Equal(count, int64(0))
}
