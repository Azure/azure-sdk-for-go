//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob_test

//
//import (
//	"bytes"
//	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal"
//	"github.com/stretchr/testify/require"
//	"io"
//)
//
////nolint
//func (s *azblobUnrecordedTestSuite) TestSectionWriter() {
//	_require := require.New(s.T())
//	b := [10]byte{}
//	buffer := internal.NewBytesWriter(b[:])
//
//	section := internal.NewSectionWriter(buffer, 0, 5)
//	_require.Equal(section.Count, int64(5))
//	_require.Equal(section.Offset, int64(0))
//	_require.Equal(section.Position, int64(0))
//
//	count, err := section.Write([]byte{1, 2, 3})
//	_require.Nil(err)
//	_require.Equal(count, 3)
//	_require.Equal(section.Position, int64(3))
//	_require.Equal(b, [10]byte{1, 2, 3, 0, 0, 0, 0, 0, 0, 0})
//
//	count, err = section.Write([]byte{4, 5, 6})
//	_require.Contains(err.Error(), "not enough space for all bytes")
//	_require.Equal(count, 2)
//	_require.Equal(section.Position, int64(5))
//	_require.Equal(b, [10]byte{1, 2, 3, 4, 5, 0, 0, 0, 0, 0})
//
//	count, err = section.Write([]byte{6, 7, 8})
//	_require.Contains(err.Error(), "end of section reached")
//	_require.Equal(count, 0)
//	_require.Equal(section.Position, int64(5))
//	_require.Equal(b, [10]byte{1, 2, 3, 4, 5, 0, 0, 0, 0, 0})
//
//	// Intentionally create a section writer which will attempt to write
//	// outside the bounds of the buffer.
//	section = internal.NewSectionWriter(buffer, 5, 6)
//	_require.Equal(section.Count, int64(6))
//	_require.Equal(section.Offset, int64(5))
//	_require.Equal(section.Position, int64(0))
//
//	count, err = section.Write([]byte{6, 7, 8})
//	_require.Nil(err)
//	_require.Equal(count, 3)
//	_require.Equal(section.Position, int64(3))
//	_require.Equal(b, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 0, 0})
//
//	// Attempt to write past the end of the section. Since the underlying
//	// buffer rejects the write it gives the same error as in the normal case.
//	count, err = section.Write([]byte{9, 10, 11})
//	_require.Contains(err.Error(), "not enough space for all bytes")
//	_require.Equal(count, 2)
//	_require.Equal(section.Position, int64(5))
//	_require.Equal(b, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
//
//	// Attempt to write past the end of the buffer. In this case the buffer
//	// rejects the write completely since it falls completely out of bounds.
//	count, err = section.Write([]byte{11, 12, 13})
//	_require.Contains(err.Error(), "offset value is out of range")
//	_require.Equal(count, 0)
//	_require.Equal(section.Position, int64(5))
//	_require.Equal(b, [10]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
//}
//
////nolint
//func (s *azblobUnrecordedTestSuite) TestSectionWriterCopySrcDestEmpty() {
//	_require := require.New(s.T())
//	input := make([]byte, 0)
//	reader := bytes.NewReader(input)
//
//	output := make([]byte, 0)
//	buffer := internal.NewBytesWriter(output)
//	section := internal.NewSectionWriter(buffer, 0, 0)
//
//	count, err := io.Copy(section, reader)
//	_require.Nil(err)
//	_require.Equal(count, int64(0))
//}
//
////nolint
//func (s *azblobUnrecordedTestSuite) TestSectionWriterCopyDestEmpty() {
//	_require := require.New(s.T())
//	input := make([]byte, 10)
//	reader := bytes.NewReader(input)
//
//	output := make([]byte, 0)
//	buffer := internal.NewBytesWriter(output)
//	section := internal.NewSectionWriter(buffer, 0, 0)
//
//	count, err := io.Copy(section, reader)
//	_require.Contains(err.Error(), "end of section reached")
//	_require.Equal(count, int64(0))
//}
