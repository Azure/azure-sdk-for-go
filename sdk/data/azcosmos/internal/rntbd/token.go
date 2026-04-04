// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"unicode/utf8"

	"github.com/google/uuid"
)

// ErrCorruptedFrame indicates a protocol violation in the RNTBD frame.
var ErrCorruptedFrame = errors.New("corrupted RNTBD frame")

// CorruptedFrameError provides details about a corrupted frame.
type CorruptedFrameError struct {
	Message string
}

func (e *CorruptedFrameError) Error() string {
	return fmt.Sprintf("corrupted frame: %s", e.Message)
}

func (e *CorruptedFrameError) Unwrap() error {
	return ErrCorruptedFrame
}

func newCorruptedFrameError(format string, args ...interface{}) error {
	return &CorruptedFrameError{Message: fmt.Sprintf(format, args...)}
}

// -----------------------------------------------------------------------------
// RntbdUUID - Microsoft GUID byte order encoding/decoding
// -----------------------------------------------------------------------------

// EmptyUUID represents a zero UUID.
var EmptyUUID = uuid.UUID{}

// EncodeUUID encodes a UUID in Microsoft GUID byte order.
// Microsoft GUID format differs from standard UUID byte order:
// - Data1 (4 bytes): little-endian
// - Data2 (2 bytes): little-endian
// - Data3 (2 bytes): little-endian
// - Data4 (8 bytes): big-endian (last two shorts, then 4 bytes)
func EncodeUUID(id uuid.UUID, out []byte) {
	if len(out) < 16 {
		panic("output buffer too small for UUID")
	}

	// UUID bytes are in big-endian format (RFC 4122)
	// We need to convert to Microsoft GUID format

	// Data1 (bytes 0-3): write as little-endian uint32
	binary.LittleEndian.PutUint32(out[0:4], binary.BigEndian.Uint32(id[0:4]))

	// Data2 (bytes 4-5): write as little-endian uint16
	binary.LittleEndian.PutUint16(out[4:6], binary.BigEndian.Uint16(id[4:6]))

	// Data3 (bytes 6-7): write as little-endian uint16
	binary.LittleEndian.PutUint16(out[6:8], binary.BigEndian.Uint16(id[6:8]))

	// Data4 (bytes 8-15): keep as-is (big-endian in both formats)
	copy(out[8:16], id[8:16])
}

// DecodeUUID decodes a UUID from Microsoft GUID byte order.
func DecodeUUID(in []byte) (uuid.UUID, error) {
	if len(in) < 16 {
		return EmptyUUID, newCorruptedFrameError("invalid UUID length: %d", len(in))
	}

	var id uuid.UUID

	// Data1 (bytes 0-3): read as little-endian uint32, write as big-endian
	binary.BigEndian.PutUint32(id[0:4], binary.LittleEndian.Uint32(in[0:4]))

	// Data2 (bytes 4-5): read as little-endian uint16, write as big-endian
	binary.BigEndian.PutUint16(id[4:6], binary.LittleEndian.Uint16(in[4:6]))

	// Data3 (bytes 6-7): read as little-endian uint16, write as big-endian
	binary.BigEndian.PutUint16(id[6:8], binary.LittleEndian.Uint16(in[6:8]))

	// Data4 (bytes 8-15): keep as-is
	copy(id[8:16], in[8:16])

	return id, nil
}

// WriteUUID writes a UUID to an io.Writer in Microsoft GUID byte order.
func WriteUUID(id uuid.UUID, w io.Writer) error {
	var buf [16]byte
	EncodeUUID(id, buf[:])
	_, err := w.Write(buf[:])
	return err
}

// ReadUUID reads a UUID from an io.Reader in Microsoft GUID byte order.
func ReadUUID(r io.Reader) (uuid.UUID, error) {
	var buf [16]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return EmptyUUID, err
	}
	return DecodeUUID(buf[:])
}

// -----------------------------------------------------------------------------
// Token Codec Interface
// -----------------------------------------------------------------------------

// Codec defines the interface for encoding/decoding token values.
type Codec interface {
	// ComputeLength returns the encoded length of the value in bytes.
	ComputeLength(value interface{}) int

	// DefaultValue returns the default value for this codec.
	DefaultValue() interface{}

	// IsValid returns true if the value can be encoded by this codec.
	IsValid(value interface{}) bool

	// Read decodes a value from the reader.
	Read(r io.Reader) (interface{}, error)

	// ReadSlice reads the raw bytes for this token from the reader.
	// This is used for lazy decoding.
	ReadSlice(r io.Reader) ([]byte, error)

	// Write encodes the value to the writer.
	Write(value interface{}, w io.Writer) error
}

// -----------------------------------------------------------------------------
// Codec Implementations
// -----------------------------------------------------------------------------

// byteCodec handles single byte values (also accepts bool).
type byteCodec struct{}

func (c byteCodec) ComputeLength(_ interface{}) int { return 1 }
func (c byteCodec) DefaultValue() interface{}       { return byte(0) }

func (c byteCodec) IsValid(value interface{}) bool {
	switch value.(type) {
	case byte, int8, bool: // byte is alias for uint8
		return true
	case int, int16, int32, int64, uint, uint16, uint32, uint64:
		return true
	}
	return false
}

func (c byteCodec) Read(r io.Reader) (interface{}, error) {
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return nil, err
	}
	return b[0], nil
}

func (c byteCodec) ReadSlice(r io.Reader) ([]byte, error) {
	b := make([]byte, 1)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (c byteCodec) Write(value interface{}, w io.Writer) error {
	var b byte
	switch v := value.(type) {
	case byte:
		b = v
	case int8:
		b = byte(v)
	case bool:
		if v {
			b = 0x01
		} else {
			b = 0x00
		}
	case int:
		b = byte(v)
	case int16:
		b = byte(v)
	case int32:
		b = byte(v)
	case int64:
		b = byte(v)
	case uint:
		b = byte(v)
	case uint16:
		b = byte(v)
	case uint32:
		b = byte(v)
	case uint64:
		b = byte(v)
	default:
		return fmt.Errorf("invalid byte value type: %T", value)
	}
	_, err := w.Write([]byte{b})
	return err
}

// ushortCodec handles unsigned 16-bit values (little-endian).
type ushortCodec struct{}

func (c ushortCodec) ComputeLength(_ interface{}) int { return 2 }
func (c ushortCodec) DefaultValue() interface{}       { return uint16(0) }

func (c ushortCodec) IsValid(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

func (c ushortCodec) Read(r io.Reader) (interface{}, error) {
	var b [2]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return nil, err
	}
	return binary.LittleEndian.Uint16(b[:]), nil
}

func (c ushortCodec) ReadSlice(r io.Reader) ([]byte, error) {
	b := make([]byte, 2)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (c ushortCodec) Write(value interface{}, w io.Writer) error {
	var v uint16
	switch n := value.(type) {
	case int:
		v = uint16(n)
	case int16:
		v = uint16(n)
	case int32:
		v = uint16(n)
	case int64:
		v = uint16(n)
	case uint:
		v = uint16(n)
	case uint8:
		v = uint16(n)
	case uint16:
		v = n
	case uint32:
		v = uint16(n)
	case uint64:
		v = uint16(n)
	default:
		return fmt.Errorf("invalid ushort value type: %T", value)
	}
	var b [2]byte
	binary.LittleEndian.PutUint16(b[:], v)
	_, err := w.Write(b[:])
	return err
}

// ulongCodec handles unsigned 32-bit values (little-endian), returned as int64.
type ulongCodec struct{}

func (c ulongCodec) ComputeLength(_ interface{}) int { return 4 }
func (c ulongCodec) DefaultValue() interface{}       { return int64(0) }

func (c ulongCodec) IsValid(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

func (c ulongCodec) Read(r io.Reader) (interface{}, error) {
	var b [4]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return nil, err
	}
	return int64(binary.LittleEndian.Uint32(b[:])), nil
}

func (c ulongCodec) ReadSlice(r io.Reader) ([]byte, error) {
	b := make([]byte, 4)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (c ulongCodec) Write(value interface{}, w io.Writer) error {
	var v uint32
	switch n := value.(type) {
	case int:
		v = uint32(n)
	case int32:
		v = uint32(n)
	case int64:
		v = uint32(n)
	case uint:
		v = uint32(n)
	case uint32:
		v = n
	case uint64:
		v = uint32(n)
	default:
		return fmt.Errorf("invalid ulong value type: %T", value)
	}
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], v)
	_, err := w.Write(b[:])
	return err
}

// longCodec handles signed 32-bit values (little-endian).
type longCodec struct{}

func (c longCodec) ComputeLength(_ interface{}) int { return 4 }
func (c longCodec) DefaultValue() interface{}       { return int32(0) }

func (c longCodec) IsValid(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

func (c longCodec) Read(r io.Reader) (interface{}, error) {
	var b [4]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return nil, err
	}
	return int32(binary.LittleEndian.Uint32(b[:])), nil
}

func (c longCodec) ReadSlice(r io.Reader) ([]byte, error) {
	b := make([]byte, 4)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (c longCodec) Write(value interface{}, w io.Writer) error {
	var v int32
	switch n := value.(type) {
	case int:
		v = int32(n)
	case int32:
		v = n
	case int64:
		v = int32(n)
	case uint:
		v = int32(n)
	case uint32:
		v = int32(n)
	case uint64:
		v = int32(n)
	default:
		return fmt.Errorf("invalid long value type: %T", value)
	}
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(v))
	_, err := w.Write(b[:])
	return err
}

// longlongCodec handles 64-bit values (little-endian).
// Used for both ULongLong and LongLong token types.
type longlongCodec struct{}

func (c longlongCodec) ComputeLength(_ interface{}) int { return 8 }
func (c longlongCodec) DefaultValue() interface{}       { return int64(0) }

func (c longlongCodec) IsValid(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

func (c longlongCodec) Read(r io.Reader) (interface{}, error) {
	var b [8]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return nil, err
	}
	return int64(binary.LittleEndian.Uint64(b[:])), nil
}

func (c longlongCodec) ReadSlice(r io.Reader) ([]byte, error) {
	b := make([]byte, 8)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (c longlongCodec) Write(value interface{}, w io.Writer) error {
	var v int64
	switch n := value.(type) {
	case int:
		v = int64(n)
	case int64:
		v = n
	case uint:
		v = int64(n)
	case uint64:
		v = int64(n)
	default:
		return fmt.Errorf("invalid longlong value type: %T", value)
	}
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(v))
	_, err := w.Write(b[:])
	return err
}

// guidCodec handles 16-byte GUID values in Microsoft byte order.
type guidCodec struct{}

func (c guidCodec) ComputeLength(_ interface{}) int { return 16 }
func (c guidCodec) DefaultValue() interface{}       { return EmptyUUID }

func (c guidCodec) IsValid(value interface{}) bool {
	_, ok := value.(uuid.UUID)
	return ok
}

func (c guidCodec) Read(r io.Reader) (interface{}, error) {
	return ReadUUID(r)
}

func (c guidCodec) ReadSlice(r io.Reader) ([]byte, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (c guidCodec) Write(value interface{}, w io.Writer) error {
	id, ok := value.(uuid.UUID)
	if !ok {
		return fmt.Errorf("invalid GUID value type: %T", value)
	}
	return WriteUUID(id, w)
}

// smallStringCodec handles strings with 1-byte length prefix (max 255 bytes).
type smallStringCodec struct{}

func (c smallStringCodec) ComputeLength(value interface{}) int {
	switch v := value.(type) {
	case string:
		return 1 + len(v) // UTF-8 length
	case []byte:
		return 1 + len(v)
	}
	return 1
}

func (c smallStringCodec) DefaultValue() interface{} { return "" }

func (c smallStringCodec) IsValid(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return len(v) <= 0xFF
	case []byte:
		return len(v) <= 0xFF && utf8.Valid(v)
	}
	return false
}

func (c smallStringCodec) Read(r io.Reader) (interface{}, error) {
	var lenBuf [1]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := int(lenBuf[0])

	if length == 0 {
		return "", nil
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	if !utf8.Valid(data) {
		return nil, newCorruptedFrameError("invalid UTF-8 string")
	}

	return string(data), nil
}

func (c smallStringCodec) ReadSlice(r io.Reader) ([]byte, error) {
	var lenBuf [1]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := int(lenBuf[0])

	result := make([]byte, 1+length)
	result[0] = lenBuf[0]
	if length > 0 {
		if _, err := io.ReadFull(r, result[1:]); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (c smallStringCodec) Write(value interface{}, w io.Writer) error {
	var data []byte
	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("invalid string value type: %T", value)
	}

	if len(data) > 0xFF {
		return newCorruptedFrameError("string too long for SmallString: %d bytes", len(data))
	}

	if _, err := w.Write([]byte{byte(len(data))}); err != nil {
		return err
	}
	if len(data) > 0 {
		if _, err := w.Write(data); err != nil {
			return err
		}
	}
	return nil
}

// stringCodec handles strings with 2-byte length prefix (max 64KB).
type stringCodec struct{}

func (c stringCodec) ComputeLength(value interface{}) int {
	switch v := value.(type) {
	case string:
		return 2 + len(v)
	case []byte:
		return 2 + len(v)
	}
	return 2
}

func (c stringCodec) DefaultValue() interface{} { return "" }

func (c stringCodec) IsValid(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return len(v) <= 0xFFFF
	case []byte:
		return len(v) <= 0xFFFF && utf8.Valid(v)
	}
	return false
}

func (c stringCodec) Read(r io.Reader) (interface{}, error) {
	var lenBuf [2]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := int(binary.LittleEndian.Uint16(lenBuf[:]))

	if length == 0 {
		return "", nil
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	if !utf8.Valid(data) {
		return nil, newCorruptedFrameError("invalid UTF-8 string")
	}

	return string(data), nil
}

func (c stringCodec) ReadSlice(r io.Reader) ([]byte, error) {
	var lenBuf [2]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := int(binary.LittleEndian.Uint16(lenBuf[:]))

	result := make([]byte, 2+length)
	copy(result[0:2], lenBuf[:])
	if length > 0 {
		if _, err := io.ReadFull(r, result[2:]); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (c stringCodec) Write(value interface{}, w io.Writer) error {
	var data []byte
	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("invalid string value type: %T", value)
	}

	if len(data) > 0xFFFF {
		return newCorruptedFrameError("string too long for String: %d bytes", len(data))
	}

	var lenBuf [2]byte
	binary.LittleEndian.PutUint16(lenBuf[:], uint16(len(data)))
	if _, err := w.Write(lenBuf[:]); err != nil {
		return err
	}
	if len(data) > 0 {
		if _, err := w.Write(data); err != nil {
			return err
		}
	}
	return nil
}

// ulongStringCodec handles strings with 4-byte length prefix (max 2GB).
type ulongStringCodec struct{}

func (c ulongStringCodec) ComputeLength(value interface{}) int {
	switch v := value.(type) {
	case string:
		return 4 + len(v)
	case []byte:
		return 4 + len(v)
	}
	return 4
}

func (c ulongStringCodec) DefaultValue() interface{} { return "" }

func (c ulongStringCodec) IsValid(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return len(v) <= math.MaxInt32
	case []byte:
		return len(v) <= math.MaxInt32 && utf8.Valid(v)
	}
	return false
}

func (c ulongStringCodec) Read(r io.Reader) (interface{}, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(lenBuf[:])

	if length > math.MaxInt32 {
		return nil, newCorruptedFrameError("string length exceeds maximum: %d", length)
	}

	if length == 0 {
		return "", nil
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	if !utf8.Valid(data) {
		return nil, newCorruptedFrameError("invalid UTF-8 string")
	}

	return string(data), nil
}

func (c ulongStringCodec) ReadSlice(r io.Reader) ([]byte, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(lenBuf[:])

	if length > math.MaxInt32 {
		return nil, newCorruptedFrameError("string length exceeds maximum: %d", length)
	}

	result := make([]byte, 4+length)
	copy(result[0:4], lenBuf[:])
	if length > 0 {
		if _, err := io.ReadFull(r, result[4:]); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (c ulongStringCodec) Write(value interface{}, w io.Writer) error {
	var data []byte
	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return fmt.Errorf("invalid string value type: %T", value)
	}

	if len(data) > math.MaxInt32 {
		return newCorruptedFrameError("string too long for ULongString: %d bytes", len(data))
	}

	var lenBuf [4]byte
	binary.LittleEndian.PutUint32(lenBuf[:], uint32(len(data)))
	if _, err := w.Write(lenBuf[:]); err != nil {
		return err
	}
	if len(data) > 0 {
		if _, err := w.Write(data); err != nil {
			return err
		}
	}
	return nil
}

// smallBytesCodec handles byte arrays with 1-byte length prefix (max 255 bytes).
type smallBytesCodec struct{}

func (c smallBytesCodec) ComputeLength(value interface{}) int {
	if v, ok := value.([]byte); ok {
		return 1 + len(v)
	}
	return 1
}

func (c smallBytesCodec) DefaultValue() interface{} { return []byte{} }

func (c smallBytesCodec) IsValid(value interface{}) bool {
	v, ok := value.([]byte)
	return ok && len(v) <= 0xFF
}

func (c smallBytesCodec) Read(r io.Reader) (interface{}, error) {
	var lenBuf [1]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := int(lenBuf[0])

	if length == 0 {
		return []byte{}, nil
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c smallBytesCodec) ReadSlice(r io.Reader) ([]byte, error) {
	var lenBuf [1]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := int(lenBuf[0])

	result := make([]byte, 1+length)
	result[0] = lenBuf[0]
	if length > 0 {
		if _, err := io.ReadFull(r, result[1:]); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (c smallBytesCodec) Write(value interface{}, w io.Writer) error {
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid bytes value type: %T", value)
	}

	if len(data) > 0xFF {
		return newCorruptedFrameError("bytes too long for SmallBytes: %d bytes", len(data))
	}

	if _, err := w.Write([]byte{byte(len(data))}); err != nil {
		return err
	}
	if len(data) > 0 {
		if _, err := w.Write(data); err != nil {
			return err
		}
	}
	return nil
}

// bytesCodec handles byte arrays with 2-byte length prefix (max 64KB).
type bytesCodec struct{}

func (c bytesCodec) ComputeLength(value interface{}) int {
	if v, ok := value.([]byte); ok {
		return 2 + len(v)
	}
	return 2
}

func (c bytesCodec) DefaultValue() interface{} { return []byte{} }

func (c bytesCodec) IsValid(value interface{}) bool {
	v, ok := value.([]byte)
	return ok && len(v) < 0xFFFF // Note: Java uses < not <=
}

func (c bytesCodec) Read(r io.Reader) (interface{}, error) {
	var lenBuf [2]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := int(binary.LittleEndian.Uint16(lenBuf[:]))

	if length == 0 {
		return []byte{}, nil
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c bytesCodec) ReadSlice(r io.Reader) ([]byte, error) {
	var lenBuf [2]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := int(binary.LittleEndian.Uint16(lenBuf[:]))

	result := make([]byte, 2+length)
	copy(result[0:2], lenBuf[:])
	if length > 0 {
		if _, err := io.ReadFull(r, result[2:]); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (c bytesCodec) Write(value interface{}, w io.Writer) error {
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid bytes value type: %T", value)
	}

	if len(data) > 0xFFFF {
		return newCorruptedFrameError("bytes too long for Bytes: %d bytes", len(data))
	}

	var lenBuf [2]byte
	binary.LittleEndian.PutUint16(lenBuf[:], uint16(len(data)))
	if _, err := w.Write(lenBuf[:]); err != nil {
		return err
	}
	if len(data) > 0 {
		if _, err := w.Write(data); err != nil {
			return err
		}
	}
	return nil
}

// ulongBytesCodec handles byte arrays with 4-byte length prefix (max 2GB).
type ulongBytesCodec struct{}

func (c ulongBytesCodec) ComputeLength(value interface{}) int {
	if v, ok := value.([]byte); ok {
		return 4 + len(v)
	}
	return 4
}

func (c ulongBytesCodec) DefaultValue() interface{} { return []byte{} }

func (c ulongBytesCodec) IsValid(value interface{}) bool {
	_, ok := value.([]byte)
	return ok
}

func (c ulongBytesCodec) Read(r io.Reader) (interface{}, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(lenBuf[:])

	if length > math.MaxInt32 {
		return nil, newCorruptedFrameError("bytes length exceeds maximum: %d", length)
	}

	if length == 0 {
		return []byte{}, nil
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	return data, nil
}

func (c ulongBytesCodec) ReadSlice(r io.Reader) ([]byte, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(r, lenBuf[:]); err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(lenBuf[:])

	if length > math.MaxInt32 {
		return nil, newCorruptedFrameError("bytes length exceeds maximum: %d", length)
	}

	result := make([]byte, 4+length)
	copy(result[0:4], lenBuf[:])
	if length > 0 {
		if _, err := io.ReadFull(r, result[4:]); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (c ulongBytesCodec) Write(value interface{}, w io.Writer) error {
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid bytes value type: %T", value)
	}

	var lenBuf [4]byte
	binary.LittleEndian.PutUint32(lenBuf[:], uint32(len(data)))
	if _, err := w.Write(lenBuf[:]); err != nil {
		return err
	}
	if len(data) > 0 {
		if _, err := w.Write(data); err != nil {
			return err
		}
	}
	return nil
}

// floatCodec handles 32-bit IEEE floating point values (little-endian).
type floatCodec struct{}

func (c floatCodec) ComputeLength(_ interface{}) int { return 4 }
func (c floatCodec) DefaultValue() interface{}       { return float32(0) }

func (c floatCodec) IsValid(value interface{}) bool {
	switch value.(type) {
	case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

func (c floatCodec) Read(r io.Reader) (interface{}, error) {
	var b [4]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return nil, err
	}
	bits := binary.LittleEndian.Uint32(b[:])
	return math.Float32frombits(bits), nil
}

func (c floatCodec) ReadSlice(r io.Reader) ([]byte, error) {
	b := make([]byte, 4)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (c floatCodec) Write(value interface{}, w io.Writer) error {
	var f float32
	switch v := value.(type) {
	case float32:
		f = v
	case float64:
		f = float32(v)
	case int:
		f = float32(v)
	case int64:
		f = float32(v)
	default:
		return fmt.Errorf("invalid float value type: %T", value)
	}
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], math.Float32bits(f))
	_, err := w.Write(b[:])
	return err
}

// doubleCodec handles 64-bit IEEE floating point values (little-endian).
type doubleCodec struct{}

func (c doubleCodec) ComputeLength(_ interface{}) int { return 8 }
func (c doubleCodec) DefaultValue() interface{}       { return float64(0) }

func (c doubleCodec) IsValid(value interface{}) bool {
	switch value.(type) {
	case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	}
	return false
}

func (c doubleCodec) Read(r io.Reader) (interface{}, error) {
	var b [8]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return nil, err
	}
	bits := binary.LittleEndian.Uint64(b[:])
	return math.Float64frombits(bits), nil
}

func (c doubleCodec) ReadSlice(r io.Reader) ([]byte, error) {
	b := make([]byte, 8)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (c doubleCodec) Write(value interface{}, w io.Writer) error {
	var f float64
	switch v := value.(type) {
	case float64:
		f = v
	case float32:
		f = float64(v)
	case int:
		f = float64(v)
	case int64:
		f = float64(v)
	default:
		return fmt.Errorf("invalid double value type: %T", value)
	}
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], math.Float64bits(f))
	_, err := w.Write(b[:])
	return err
}

// noneCodec handles the Invalid token type (no data).
type noneCodec struct{}

func (c noneCodec) ComputeLength(_ interface{}) int        { return 0 }
func (c noneCodec) DefaultValue() interface{}              { return nil }
func (c noneCodec) IsValid(_ interface{}) bool             { return true }
func (c noneCodec) Read(_ io.Reader) (interface{}, error)  { return nil, nil }
func (c noneCodec) ReadSlice(_ io.Reader) ([]byte, error)  { return nil, nil }
func (c noneCodec) Write(_ interface{}, _ io.Writer) error { return nil }

// -----------------------------------------------------------------------------
// Codec Registry
// -----------------------------------------------------------------------------

// Singleton codec instances
var (
	codecByte        Codec = byteCodec{}
	codecUShort      Codec = ushortCodec{}
	codecULong       Codec = ulongCodec{}
	codecLong        Codec = longCodec{}
	codecLongLong    Codec = longlongCodec{}
	codecGuid        Codec = guidCodec{}
	codecSmallString Codec = smallStringCodec{}
	codecString      Codec = stringCodec{}
	codecULongString Codec = ulongStringCodec{}
	codecSmallBytes  Codec = smallBytesCodec{}
	codecBytes       Codec = bytesCodec{}
	codecULongBytes  Codec = ulongBytesCodec{}
	codecFloat       Codec = floatCodec{}
	codecDouble      Codec = doubleCodec{}
	codecNone        Codec = noneCodec{}
)

// GetCodec returns the codec for a given token type.
func GetCodec(tokenType TokenType) Codec {
	switch tokenType {
	case TokenByte:
		return codecByte
	case TokenUShort:
		return codecUShort
	case TokenULong:
		return codecULong
	case TokenLong:
		return codecLong
	case TokenULongLong, TokenLongLong:
		return codecLongLong
	case TokenGuid:
		return codecGuid
	case TokenSmallString:
		return codecSmallString
	case TokenString:
		return codecString
	case TokenULongString:
		return codecULongString
	case TokenSmallBytes:
		return codecSmallBytes
	case TokenBytes:
		return codecBytes
	case TokenULongBytes:
		return codecULongBytes
	case TokenFloat:
		return codecFloat
	case TokenDouble:
		return codecDouble
	case TokenInvalid:
		return codecNone
	default:
		return codecNone
	}
}

// -----------------------------------------------------------------------------
// Token Structure
// -----------------------------------------------------------------------------

// Token represents a single RNTBD protocol token.
// Wire format: [ID (2 bytes LE)] [Type (1 byte)] [Value (variable)]
type Token struct {
	ID       uint16      // Header ID
	Type     TokenType   // Token type determining the codec
	Value    interface{} // Decoded value (or raw []byte slice for lazy decode)
	present  bool        // Whether this token is present in the stream
	rawSlice []byte      // Raw bytes for lazy decoding
}

// TokenHeaderLength is the size of the token header (ID + Type).
const TokenHeaderLength = 3 // 2 bytes ID + 1 byte Type

// IsPresent returns true if this token has a value set.
func (t *Token) IsPresent() bool {
	return t.present
}

// GetValue returns the decoded value, performing lazy decode if needed.
func (t *Token) GetValue() (interface{}, error) {
	if !t.present {
		return GetCodec(t.Type).DefaultValue(), nil
	}

	// If we have a raw slice but no decoded value, decode it now
	if t.rawSlice != nil && t.Value == nil {
		codec := GetCodec(t.Type)
		val, err := codec.Read(bytes.NewReader(t.rawSlice))
		if err != nil {
			return nil, err
		}
		t.Value = val
		t.rawSlice = nil // Clear the slice after decoding
	}

	return t.Value, nil
}

// SetValue sets the token value and marks it as present.
func (t *Token) SetValue(value interface{}) error {
	codec := GetCodec(t.Type)
	if !codec.IsValid(value) {
		return fmt.Errorf("invalid value for token type %s: %T", t.Type, value)
	}
	t.Value = value
	t.present = true
	t.rawSlice = nil
	return nil
}

// ComputeLength returns the total encoded length of this token.
func (t *Token) ComputeLength() int {
	if !t.present {
		return 0
	}
	codec := GetCodec(t.Type)
	return TokenHeaderLength + codec.ComputeLength(t.Value)
}

// Encode writes the token to the writer.
func (t *Token) Encode(w io.Writer) error {
	if !t.present {
		return nil
	}

	// Write ID (2 bytes, little-endian)
	var idBuf [2]byte
	binary.LittleEndian.PutUint16(idBuf[:], t.ID)
	if _, err := w.Write(idBuf[:]); err != nil {
		return err
	}

	// Write Type (1 byte)
	if _, err := w.Write([]byte{byte(t.Type)}); err != nil {
		return err
	}

	// Write Value
	// If we have a raw slice, write it directly
	if t.rawSlice != nil {
		_, err := w.Write(t.rawSlice)
		return err
	}

	// Otherwise encode the value
	codec := GetCodec(t.Type)
	return codec.Write(t.Value, w)
}

// DecodeToken reads a token from the reader.
// It reads the header (ID + Type) and then the value using readSlice for lazy decoding.
func DecodeToken(r io.Reader) (*Token, error) {
	// Read ID (2 bytes, little-endian)
	var idBuf [2]byte
	if _, err := io.ReadFull(r, idBuf[:]); err != nil {
		return nil, err
	}
	id := binary.LittleEndian.Uint16(idBuf[:])

	// Read Type (1 byte)
	var typeBuf [1]byte
	if _, err := io.ReadFull(r, typeBuf[:]); err != nil {
		return nil, err
	}
	tokenType := TokenTypeFromID(typeBuf[0])

	// Read value slice using the codec's ReadSlice
	codec := GetCodec(tokenType)
	rawSlice, err := codec.ReadSlice(r)
	if err != nil {
		return nil, err
	}

	return &Token{
		ID:       id,
		Type:     tokenType,
		present:  true,
		rawSlice: rawSlice,
	}, nil
}

// TokenTypeFromID returns the TokenType for the given byte ID.
func TokenTypeFromID(id byte) TokenType {
	switch TokenType(id) {
	case TokenByte, TokenUShort, TokenULong, TokenLong,
		TokenULongLong, TokenLongLong, TokenGuid,
		TokenSmallString, TokenString, TokenULongString,
		TokenSmallBytes, TokenBytes, TokenULongBytes,
		TokenFloat, TokenDouble:
		return TokenType(id)
	default:
		return TokenInvalid
	}
}

// NewToken creates a new token with the given ID, type, and value.
func NewToken(id uint16, tokenType TokenType, value interface{}) (*Token, error) {
	t := &Token{
		ID:   id,
		Type: tokenType,
	}
	if value != nil {
		if err := t.SetValue(value); err != nil {
			return nil, err
		}
	}
	return t, nil
}
