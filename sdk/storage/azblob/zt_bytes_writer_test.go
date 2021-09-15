// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"bytes"
	"github.com/stretchr/testify/assert"
)

//nolint
func (s *azblobUnrecordedTestSuite) TestBytesWriterWriteAt() {
	_assert := assert.New(s.T())
	b := make([]byte, 10)
	buffer := newBytesWriter(b)

	count, err := buffer.WriteAt([]byte{1, 2}, 10)
	_assert.Contains(err.Error(), "offset value is out of range")
	_assert.Equal(count, 0)

	count, err = buffer.WriteAt([]byte{1, 2}, -1)
	_assert.Contains(err.Error(), "offset value is out of range")
	_assert.Equal(count, 0)

	count, err = buffer.WriteAt([]byte{1, 2}, 9)
	_assert.Contains(err.Error(), "not enough space for all bytes")
	_assert.Equal(count, 1)
	_assert.Equal(bytes.Compare(b, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}), 0)

	count, err = buffer.WriteAt([]byte{1, 2}, 8)
	_assert.Nil(err)
	_assert.Equal(count, 2)
	_assert.Equal(bytes.Compare(b, []byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 2}), 0)
}
