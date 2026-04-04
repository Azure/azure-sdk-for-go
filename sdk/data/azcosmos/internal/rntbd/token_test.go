// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeUUID(t *testing.T) {
	tests := []struct {
		name string
		uuid uuid.UUID
	}{
		{"empty", EmptyUUID},
		{"random1", uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")},
		{"random2", uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")},
		{"max", uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf [16]byte
			EncodeUUID(tt.uuid, buf[:])

			decoded, err := DecodeUUID(buf[:])
			require.NoError(t, err)
			require.Equal(t, tt.uuid, decoded)
		})
	}
}

func TestUUIDMicrosoftGUIDByteOrder(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

	var buf [16]byte
	EncodeUUID(id, buf[:])

	require.Equal(t, uint32(0x550e8400), binary.LittleEndian.Uint32(buf[0:4]))
	require.Equal(t, uint16(0xe29b), binary.LittleEndian.Uint16(buf[4:6]))
	require.Equal(t, uint16(0x41d4), binary.LittleEndian.Uint16(buf[6:8]))

	require.Equal(t, byte(0xa7), buf[8])
	require.Equal(t, byte(0x16), buf[9])
	require.Equal(t, byte(0x44), buf[10])
	require.Equal(t, byte(0x66), buf[11])
	require.Equal(t, byte(0x55), buf[12])
	require.Equal(t, byte(0x44), buf[13])
	require.Equal(t, byte(0x00), buf[14])
	require.Equal(t, byte(0x00), buf[15])
}

func TestByteCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenByte)

	tests := []struct {
		name  string
		value interface{}
		want  byte
	}{
		{"byte_zero", byte(0), 0},
		{"byte_max", byte(255), 255},
		{"byte_mid", byte(127), 127},
		{"bool_true", true, 1},
		{"bool_false", false, 0},
		{"int_value", int(42), 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 1, buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.want, val)
		})
	}
}

func TestUShortCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenUShort)

	tests := []struct {
		name  string
		value interface{}
		want  uint16
	}{
		{"zero", uint16(0), 0},
		{"max", uint16(65535), 65535},
		{"mid", uint16(32767), 32767},
		{"int_value", int(1000), 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 2, buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.want, val)
		})
	}
}

func TestULongCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenULong)

	tests := []struct {
		name  string
		value interface{}
		want  int64
	}{
		{"zero", uint32(0), 0},
		{"max", uint32(4294967295), 4294967295},
		{"mid", uint32(2147483647), 2147483647},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 4, buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.want, val)
		})
	}
}

func TestLongCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenLong)

	tests := []struct {
		name  string
		value interface{}
		want  int32
	}{
		{"zero", int32(0), 0},
		{"positive", int32(2147483647), 2147483647},
		{"negative", int32(-2147483648), -2147483648},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 4, buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.want, val)
		})
	}
}

func TestLongLongCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenLongLong)

	tests := []struct {
		name  string
		value interface{}
		want  int64
	}{
		{"zero", int64(0), 0},
		{"positive", int64(9223372036854775807), 9223372036854775807},
		{"negative", int64(-9223372036854775808), -9223372036854775808},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 8, buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.want, val)
		})
	}
}

func TestULongLongCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenULongLong)

	tests := []struct {
		name  string
		value interface{}
		want  int64
	}{
		{"zero", int64(0), 0},
		{"max_signed", int64(9223372036854775807), 9223372036854775807},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 8, buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.want, val)
		})
	}
}

func TestGuidCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenGuid)

	tests := []struct {
		name  string
		value uuid.UUID
	}{
		{"empty", EmptyUUID},
		{"random", uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")},
		{"another", uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 16, buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.value, val)
		})
	}
}

func TestSmallStringCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenSmallString)

	tests := []struct {
		name  string
		value string
	}{
		{"empty", ""},
		{"ascii", "Hello, World!"},
		{"unicode", "こんにちは"},
		{"max_length", string(make([]byte, 255))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 1+len(tt.value), buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.value, val)
		})
	}
}

func TestSmallStringCodecTooLong(t *testing.T) {
	codec := GetCodec(TokenSmallString)

	longString := string(make([]byte, 256))
	require.False(t, codec.IsValid(longString))
}

func TestStringCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenString)

	tests := []struct {
		name  string
		value string
	}{
		{"empty", ""},
		{"ascii", "Hello, World!"},
		{"unicode", "こんにちは世界"},
		{"medium", string(make([]byte, 1000))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 2+len(tt.value), buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.value, val)
		})
	}
}

func TestULongStringCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenULongString)

	tests := []struct {
		name  string
		value string
	}{
		{"empty", ""},
		{"ascii", "Hello, World!"},
		{"unicode", "日本語テスト"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 4+len(tt.value), buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.value, val)
		})
	}
}

func TestSmallBytesCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenSmallBytes)

	tests := []struct {
		name  string
		value []byte
	}{
		{"empty", []byte{}},
		{"single", []byte{0x42}},
		{"multiple", []byte{0x01, 0x02, 0x03, 0x04}},
		{"max_length", make([]byte, 255)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 1+len(tt.value), buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.value, val)
		})
	}
}

func TestBytesCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenBytes)

	tests := []struct {
		name  string
		value []byte
	}{
		{"empty", []byte{}},
		{"single", []byte{0x42}},
		{"medium", make([]byte, 1000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 2+len(tt.value), buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			valBytes, ok := val.([]byte)
			require.True(t, ok)
			require.Equal(t, tt.value, valBytes)
		})
	}
}

func TestULongBytesCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenULongBytes)

	tests := []struct {
		name  string
		value []byte
	}{
		{"empty", []byte{}},
		{"single", []byte{0x42}},
		{"medium", make([]byte, 1000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 4+len(tt.value), buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			valBytes, ok := val.([]byte)
			require.True(t, ok)
			require.Equal(t, tt.value, valBytes)
		})
	}
}

func TestFloatCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenFloat)

	tests := []struct {
		name  string
		value float32
	}{
		{"zero", 0.0},
		{"positive", 3.14159},
		{"negative", -2.71828},
		{"max", math.MaxFloat32},
		{"min", -math.MaxFloat32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 4, buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.value, val)
		})
	}
}

func TestDoubleCodecRoundtrip(t *testing.T) {
	codec := GetCodec(TokenDouble)

	tests := []struct {
		name  string
		value float64
	}{
		{"zero", 0.0},
		{"positive", 3.14159265358979},
		{"negative", -2.71828182845904},
		{"max", math.MaxFloat64},
		{"min", -math.MaxFloat64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.True(t, codec.IsValid(tt.value))

			var buf bytes.Buffer
			err := codec.Write(tt.value, &buf)
			require.NoError(t, err)
			require.Equal(t, 8, buf.Len())

			val, err := codec.Read(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.value, val)
		})
	}
}

func TestNoneCodec(t *testing.T) {
	codec := GetCodec(TokenInvalid)

	require.True(t, codec.IsValid(nil))
	require.True(t, codec.IsValid("anything"))
	require.Equal(t, 0, codec.ComputeLength(nil))
	require.Nil(t, codec.DefaultValue())

	var buf bytes.Buffer
	err := codec.Write(nil, &buf)
	require.NoError(t, err)
	require.Equal(t, 0, buf.Len())

	val, err := codec.Read(&buf)
	require.NoError(t, err)
	require.Nil(t, val)
}

func TestTokenRoundtrip(t *testing.T) {
	tests := []struct {
		name      string
		id        uint16
		tokenType TokenType
		value     interface{}
	}{
		{"byte_token", 0x0001, TokenByte, byte(42)},
		{"ushort_token", 0x0002, TokenUShort, uint16(1000)},
		{"ulong_token", 0x0003, TokenULong, uint32(100000)},
		{"string_token", 0x0004, TokenSmallString, "hello"},
		{"guid_token", 0x0005, TokenGuid, uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewToken(tt.id, tt.tokenType, tt.value)
			require.NoError(t, err)
			require.True(t, token.IsPresent())

			var buf bytes.Buffer
			err = token.Encode(&buf)
			require.NoError(t, err)

			decoded, err := DecodeToken(&buf)
			require.NoError(t, err)
			require.Equal(t, tt.id, decoded.ID)
			require.Equal(t, tt.tokenType, decoded.Type)

			val, err := decoded.GetValue()
			require.NoError(t, err)

			switch expected := tt.value.(type) {
			case byte:
				require.Equal(t, expected, val)
			case uint16:
				require.Equal(t, expected, val)
			case uint32:
				require.Equal(t, int64(expected), val)
			case string:
				require.Equal(t, expected, val)
			case uuid.UUID:
				require.Equal(t, expected, val)
			}
		})
	}
}

func TestTokenNotPresent(t *testing.T) {
	token := &Token{
		ID:   0x0001,
		Type: TokenByte,
	}

	require.False(t, token.IsPresent())
	require.Equal(t, 0, token.ComputeLength())

	var buf bytes.Buffer
	err := token.Encode(&buf)
	require.NoError(t, err)
	require.Equal(t, 0, buf.Len())

	val, err := token.GetValue()
	require.NoError(t, err)
	require.Equal(t, byte(0), val)
}

func TestTokenTypeFromID(t *testing.T) {
	tests := []struct {
		id       byte
		expected TokenType
	}{
		{0x00, TokenByte},
		{0x01, TokenUShort},
		{0x02, TokenULong},
		{0x03, TokenLong},
		{0x04, TokenULongLong},
		{0x05, TokenLongLong},
		{0x06, TokenGuid},
		{0x07, TokenSmallString},
		{0x08, TokenString},
		{0x09, TokenULongString},
		{0x0A, TokenSmallBytes},
		{0x0B, TokenBytes},
		{0x0C, TokenULongBytes},
		{0x0D, TokenFloat},
		{0x0E, TokenDouble},
		{0xFF, TokenInvalid},
		{0x10, TokenInvalid},
		{0x99, TokenInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.expected.String(), func(t *testing.T) {
			result := TokenTypeFromID(tt.id)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestCodecComputeLength(t *testing.T) {
	tests := []struct {
		name      string
		tokenType TokenType
		value     interface{}
		expected  int
	}{
		{"byte", TokenByte, byte(1), 1},
		{"ushort", TokenUShort, uint16(1), 2},
		{"ulong", TokenULong, uint32(1), 4},
		{"long", TokenLong, int32(1), 4},
		{"longlong", TokenLongLong, int64(1), 8},
		{"ulonglong", TokenULongLong, int64(1), 8},
		{"guid", TokenGuid, EmptyUUID, 16},
		{"smallstring_empty", TokenSmallString, "", 1},
		{"smallstring_hello", TokenSmallString, "hello", 6},
		{"string_empty", TokenString, "", 2},
		{"string_hello", TokenString, "hello", 7},
		{"ulongstring_empty", TokenULongString, "", 4},
		{"ulongstring_hello", TokenULongString, "hello", 9},
		{"smallbytes_empty", TokenSmallBytes, []byte{}, 1},
		{"smallbytes_data", TokenSmallBytes, []byte{1, 2, 3}, 4},
		{"bytes_empty", TokenBytes, []byte{}, 2},
		{"bytes_data", TokenBytes, []byte{1, 2, 3}, 5},
		{"ulongbytes_empty", TokenULongBytes, []byte{}, 4},
		{"ulongbytes_data", TokenULongBytes, []byte{1, 2, 3}, 7},
		{"float", TokenFloat, float32(1.0), 4},
		{"double", TokenDouble, float64(1.0), 8},
		{"invalid", TokenInvalid, nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codec := GetCodec(tt.tokenType)
			require.Equal(t, tt.expected, codec.ComputeLength(tt.value))
		})
	}
}

func TestLittleEndianEncoding(t *testing.T) {
	t.Run("ushort_0x1234", func(t *testing.T) {
		codec := GetCodec(TokenUShort)
		var buf bytes.Buffer
		err := codec.Write(uint16(0x1234), &buf)
		require.NoError(t, err)
		require.Equal(t, []byte{0x34, 0x12}, buf.Bytes())
	})

	t.Run("ulong_0x12345678", func(t *testing.T) {
		codec := GetCodec(TokenULong)
		var buf bytes.Buffer
		err := codec.Write(uint32(0x12345678), &buf)
		require.NoError(t, err)
		require.Equal(t, []byte{0x78, 0x56, 0x34, 0x12}, buf.Bytes())
	})

	t.Run("longlong_0x123456789ABCDEF0", func(t *testing.T) {
		codec := GetCodec(TokenLongLong)
		var buf bytes.Buffer
		err := codec.Write(int64(0x123456789ABCDEF0), &buf)
		require.NoError(t, err)
		require.Equal(t, []byte{0xF0, 0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12}, buf.Bytes())
	})
}

func TestReadSlice(t *testing.T) {
	tests := []struct {
		name      string
		tokenType TokenType
		value     interface{}
		sliceLen  int
	}{
		{"byte", TokenByte, byte(42), 1},
		{"ushort", TokenUShort, uint16(1000), 2},
		{"ulong", TokenULong, uint32(100000), 4},
		{"guid", TokenGuid, uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"), 16},
		{"smallstring", TokenSmallString, "hello", 6},
		{"string", TokenString, "hello", 7},
		{"smallbytes", TokenSmallBytes, []byte{1, 2, 3}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codec := GetCodec(tt.tokenType)

			var writeBuf bytes.Buffer
			err := codec.Write(tt.value, &writeBuf)
			require.NoError(t, err)

			slice, err := codec.ReadSlice(bytes.NewReader(writeBuf.Bytes()))
			require.NoError(t, err)
			require.Equal(t, tt.sliceLen, len(slice))
		})
	}
}
