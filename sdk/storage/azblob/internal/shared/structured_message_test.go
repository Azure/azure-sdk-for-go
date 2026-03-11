//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"bytes"
	"encoding/binary"
	"hash/crc64"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeRoundTripSmallData(t *testing.T) {
	data := []byte("Hello, structured message!")
	result := EncodeStructuredMessage(data, 0)

	require.Equal(t, int64(len(data)), result.OriginalContentLength)
	require.Greater(t, len(result.EncodedData), len(data))

	decoded, err := DecodeStructuredMessage(result.EncodedData)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
	require.Equal(t, SMVersion, decoded.Version)
	require.Equal(t, SMFlagCRC64, decoded.Flags)
	require.Equal(t, uint16(1), decoded.NumSegments)
}

func TestEncodeDecodeRoundTripEmptyData(t *testing.T) {
	data := []byte{}
	result := EncodeStructuredMessage(data, 0)

	require.Equal(t, int64(0), result.OriginalContentLength)

	decoded, err := DecodeStructuredMessage(result.EncodedData)
	require.NoError(t, err)
	require.Equal(t, 0, len(decoded.Data))
	require.Equal(t, uint16(1), decoded.NumSegments)
}

func TestEncodeDecodeRoundTripExactSegmentSize(t *testing.T) {
	segSize := 1024
	data := make([]byte, segSize)
	for i := range data {
		data[i] = byte(i % 256)
	}
	result := EncodeStructuredMessage(data, segSize)

	require.Equal(t, int64(segSize), result.OriginalContentLength)

	decoded, err := DecodeStructuredMessage(result.EncodedData)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
	require.Equal(t, uint16(1), decoded.NumSegments)
}

func TestEncodeDecodeRoundTripMultiSegment(t *testing.T) {
	segSize := 100
	data := make([]byte, 350) // 4 segments: 100 + 100 + 100 + 50
	for i := range data {
		data[i] = byte(i % 256)
	}
	result := EncodeStructuredMessage(data, segSize)

	decoded, err := DecodeStructuredMessage(result.EncodedData)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
	require.Equal(t, uint16(4), decoded.NumSegments)
}

func TestEncodeDecodeRoundTripLargerData(t *testing.T) {
	data := make([]byte, 1024*1024) // 1MB
	for i := range data {
		data[i] = byte(i % 251)
	}

	segSize := 256 * 1024 // 256KB segments => 4 segments
	result := EncodeStructuredMessage(data, segSize)

	require.Equal(t, int64(len(data)), result.OriginalContentLength)

	decoded, err := DecodeStructuredMessage(result.EncodedData)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
	require.Equal(t, uint16(4), decoded.NumSegments)
}

func TestEncodeDecodeRoundTripSingleByte(t *testing.T) {
	data := []byte{0x42}
	result := EncodeStructuredMessage(data, 0)

	decoded, err := DecodeStructuredMessage(result.EncodedData)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
}

func TestEncodeDecodeRoundTripSegmentSizeOne(t *testing.T) {
	data := []byte("ABC")
	result := EncodeStructuredMessage(data, 1)

	decoded, err := DecodeStructuredMessage(result.EncodedData)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
	require.Equal(t, uint16(3), decoded.NumSegments)
}

func TestEncodeMessageFormat(t *testing.T) {
	data := []byte("ABCDEFGHIJ") // 10 bytes
	segSize := 5                 // 2 segments of 5 bytes each
	result := EncodeStructuredMessage(data, segSize)

	smData := result.EncodedData

	// Verify Message Header
	require.Equal(t, SMVersion, smData[0])

	msgLen := binary.LittleEndian.Uint64(smData[1:9])
	require.Equal(t, uint64(len(smData)), msgLen)

	flags := binary.LittleEndian.Uint16(smData[9:11])
	require.Equal(t, SMFlagCRC64, flags)

	numSegments := binary.LittleEndian.Uint16(smData[11:13])
	require.Equal(t, uint16(2), numSegments)

	offset := SMHeaderSize

	// Segment 1
	segNum1 := binary.LittleEndian.Uint16(smData[offset : offset+2])
	require.Equal(t, uint16(1), segNum1)
	segLen1 := int64(binary.LittleEndian.Uint64(smData[offset+2 : offset+10]))
	require.Equal(t, int64(5), segLen1)
	offset += SMSegmentHeaderSize

	seg1Data := smData[offset : offset+5]
	require.Equal(t, []byte("ABCDE"), seg1Data)
	offset += 5

	seg1CRC := binary.LittleEndian.Uint64(smData[offset : offset+8])
	expectedSeg1CRC := crc64.Checksum([]byte("ABCDE"), CRC64Table)
	require.Equal(t, expectedSeg1CRC, seg1CRC)
	offset += 8

	// Segment 2
	segNum2 := binary.LittleEndian.Uint16(smData[offset : offset+2])
	require.Equal(t, uint16(2), segNum2)
	segLen2 := int64(binary.LittleEndian.Uint64(smData[offset+2 : offset+10]))
	require.Equal(t, int64(5), segLen2)
	offset += SMSegmentHeaderSize

	seg2Data := smData[offset : offset+5]
	require.Equal(t, []byte("FGHIJ"), seg2Data)
	offset += 5

	seg2CRC := binary.LittleEndian.Uint64(smData[offset : offset+8])
	expectedSeg2CRC := crc64.Checksum([]byte("FGHIJ"), CRC64Table)
	require.Equal(t, expectedSeg2CRC, seg2CRC)
	offset += 8

	// Message Trailer CRC64
	msgCRC := binary.LittleEndian.Uint64(smData[offset : offset+8])
	expectedMsgCRC := crc64.Checksum(data, CRC64Table)
	require.Equal(t, expectedMsgCRC, msgCRC)
	offset += 8

	require.Equal(t, len(smData), offset)
}

func TestEncodeDefaultSegmentSize(t *testing.T) {
	data := make([]byte, 100)
	result := EncodeStructuredMessage(data, 0)

	// With default 4MB segment size, 100 bytes should be 1 segment
	decoded, err := DecodeStructuredMessage(result.EncodedData)
	require.NoError(t, err)
	require.Equal(t, uint16(1), decoded.NumSegments)
}

func TestEncodeMessageLength(t *testing.T) {
	data := []byte("ABCDEFGHIJ") // 10 bytes
	segSize := 5                 // 2 segments

	// Expected length:
	// Header: 13
	// Segment 1: 10 (header) + 5 (data) + 8 (CRC) = 23
	// Segment 2: 10 (header) + 5 (data) + 8 (CRC) = 23
	// Trailer: 8
	// Total: 13 + 23 + 23 + 8 = 67

	result := EncodeStructuredMessage(data, segSize)
	require.Equal(t, 67, len(result.EncodedData))
}

func TestEncodeCRC64MatchesSharedTable(t *testing.T) {
	data := []byte("CRC validation test data")
	expectedCRC := crc64.Checksum(data, CRC64Table)

	result := EncodeStructuredMessage(data, 0)
	smData := result.EncodedData

	// Trailer CRC is last 8 bytes
	trailerCRC := binary.LittleEndian.Uint64(smData[len(smData)-8:])
	require.Equal(t, expectedCRC, trailerCRC)
}

func TestDecodeInvalid(t *testing.T) {
	badInputs := []struct {
		name    string
		data    []byte
		errText string
	}{
		{
			name:    "TruncatedHeader",
			data:    []byte{1, 2, 3},
			errText: "too short for header",
		},
		{
			name:    "WrongVersion",
			data:    makeCorruptedSM([]byte("test"), func(d []byte) { d[0] = 99 }),
			errText: "unsupported structured message version",
		},
		{
			name:    "LengthMismatch",
			data:    makeCorruptedSM([]byte("test"), func(d []byte) { binary.LittleEndian.PutUint64(d[1:9], 999) }),
			errText: "length mismatch",
		},
		{
			name:    "CorruptedSegmentCRC",
			data:    makeCorruptedSM([]byte("Hello, world!"), func(d []byte) { d[36] ^= 0xFF }),
			errText: "CRC64 mismatch",
		},
		{
			name:    "CorruptedData",
			data:    makeCorruptedSM([]byte("Hello, world!"), func(d []byte) { d[25] ^= 0xFF }),
			errText: "CRC64 mismatch",
		},
		{
			name:    "CorruptedTrailerCRC",
			data:    makeCorruptedSM([]byte("Hello, world!"), func(d []byte) { d[len(d)-1] ^= 0xFF }),
			errText: "", // could be segment or trailer mismatch
		},
	}

	for _, tt := range badInputs {
		_, err := DecodeStructuredMessage(tt.data)
		require.Error(t, err)
		if tt.errText != "" {
			require.Contains(t, err.Error(), tt.errText)
		}
	}
}

// makeCorruptedSM encodes data then applies a corruption function on the result.
func makeCorruptedSM(data []byte, corrupt func([]byte)) []byte {
	result := EncodeStructuredMessage(data, 0)
	smData := make([]byte, len(result.EncodedData))
	copy(smData, result.EncodedData)
	corrupt(smData)
	return smData
}

func TestEncoderReadSeekClose(t *testing.T) {
	data := []byte("encoder test data")
	enc := NewStructuredMessageEncoder(data, 0)

	require.Equal(t, int64(len(data)), enc.OriginalContentLength())
	require.Greater(t, enc.EncodedLength(), int64(len(data)))

	// Read all
	buf := make([]byte, enc.EncodedLength())
	n, err := io.ReadFull(enc, buf)
	require.NoError(t, err)
	require.Equal(t, int(enc.EncodedLength()), n)

	// Seek back to start
	pos, err := enc.Seek(0, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(0), pos)

	// Read again and compare
	buf2 := make([]byte, enc.EncodedLength())
	n2, err := io.ReadFull(enc, buf2)
	require.NoError(t, err)
	require.Equal(t, int(enc.EncodedLength()), n2)
	require.Equal(t, buf, buf2)

	require.NoError(t, enc.Close())

	// Decode the output to verify correctness
	decoded, err := DecodeStructuredMessage(buf)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
}

func TestEncoderAsReadSeekCloser(t *testing.T) {
	data := []byte("interface compliance test")
	enc := NewStructuredMessageEncoder(data, 0)

	var _ io.ReadSeekCloser = enc

	allData, err := io.ReadAll(enc)
	require.NoError(t, err)
	require.Equal(t, int(enc.EncodedLength()), len(allData))
}

func TestDecoderReadClose(t *testing.T) {
	data := []byte("decoder test with some content here")
	result := EncodeStructuredMessage(data, 10)

	body := io.NopCloser(bytes.NewReader(result.EncodedData))
	dec := NewStructuredMessageDecoder(body)

	rawData, err := io.ReadAll(dec)
	require.NoError(t, err)
	require.Equal(t, data, rawData)

	decResult := dec.DecodeResult()
	require.NotNil(t, decResult)
	require.Equal(t, SMVersion, decResult.Version)
	require.Equal(t, SMFlagCRC64, decResult.Flags)

	require.NoError(t, dec.Close())
}

func TestDecoderInvalidBody(t *testing.T) {
	body := io.NopCloser(bytes.NewReader([]byte{0xFF, 0x01, 0x02}))
	dec := NewStructuredMessageDecoder(body)

	_, err := io.ReadAll(dec)
	require.Error(t, err)
}

func TestDecoderDecodeResultBeforeRead(t *testing.T) {
	data := []byte("test")
	result := EncodeStructuredMessage(data, 0)
	body := io.NopCloser(bytes.NewReader(result.EncodedData))
	dec := NewStructuredMessageDecoder(body)

	require.Nil(t, dec.DecodeResult())
}

func TestStructuredMessageConstants(t *testing.T) {
	require.Equal(t, uint8(1), SMVersion)
	require.Equal(t, uint16(0x0001), SMFlagCRC64)
	require.Equal(t, 4*1024*1024, SMDefaultSegmentSize)
	require.Equal(t, 13, SMHeaderSize)
	require.Equal(t, 10, SMSegmentHeaderSize)
	require.Equal(t, 8, SMSegmentFooterSize)
	require.Equal(t, 8, SMMessageTrailerSize)
	require.Equal(t, "XSM/1.0; properties=crc64", SMHeaderValue)
}
