// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob
//
//import (
//	"bytes"
//
//	chk "gopkg.in/check.v1"
//)
//
//func (s *aztestsSuite) TestBytesWriterWriteAt(c *chk.C) {
//	b := make([]byte, 10)
//	buffer := newBytesWriter(b)
//
//	count, err := buffer.WriteAt([]byte{1, 2}, 10)
//	c.Assert(err, chk.ErrorMatches, "offset value is out of range")
//	c.Assert(count, chk.Equals, 0)
//
//	count, err = buffer.WriteAt([]byte{1, 2}, -1)
//	c.Assert(err, chk.ErrorMatches, "offset value is out of range")
//	c.Assert(count, chk.Equals, 0)
//
//	count, err = buffer.WriteAt([]byte{1, 2}, 9)
//	c.Assert(err, chk.ErrorMatches, "not enough space for all bytes")
//	c.Assert(count, chk.Equals, 1)
//	c.Assert(bytes.Compare(b, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}), chk.Equals, 0)
//
//	count, err = buffer.WriteAt([]byte{1, 2}, 8)
//	c.Assert(err, chk.IsNil)
//	c.Assert(count, chk.Equals, 2)
//	c.Assert(bytes.Compare(b, []byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 2}), chk.Equals, 0)
//}
