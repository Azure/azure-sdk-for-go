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
	result := SMEncode(data, 0)

	require.Equal(t, int64(len(data)), result.OriginalContentLength)
	require.Greater(t, len(result.EncodedData), len(data))

	decoded, err := SMDecode(result.EncodedData)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
	require.Equal(t, SMVersion, decoded.Version)
	require.Equal(t, SMFlagCRC64, decoded.Flags)
	require.Equal(t, uint16(1), decoded.NumSegments)
}

func TestEncodeDecodeRoundTripEmptyData(t *testing.T) {
	data := []byte{}
	result := SMEncode(data, 0)

	require.Equal(t, int64(0), result.OriginalContentLength)

	decoded, err := SMDecode(result.EncodedData)
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
	result := SMEncode(data, segSize)

	require.Equal(t, int64(segSize), result.OriginalContentLength)

	decoded, err := SMDecode(result.EncodedData)
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
	result := SMEncode(data, segSize)

	decoded, err := SMDecode(result.EncodedData)
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
	result := SMEncode(data, segSize)

	require.Equal(t, int64(len(data)), result.OriginalContentLength)

	decoded, err := SMDecode(result.EncodedData)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
	require.Equal(t, uint16(4), decoded.NumSegments)
}

func TestEncodeDecodeRoundTripSingleByte(t *testing.T) {
	data := []byte{0x42}
	result := SMEncode(data, 0)

	decoded, err := SMDecode(result.EncodedData)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
}

func TestEncodeDecodeRoundTripSegmentSizeOne(t *testing.T) {
	data := []byte("ABC")
	result := SMEncode(data, 1)

	decoded, err := SMDecode(result.EncodedData)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
	require.Equal(t, uint16(3), decoded.NumSegments)
}

func TestEncodeMessageFormat(t *testing.T) {
	data := []byte("ABCDEFGHIJ") // 10 bytes
	segSize := 5                 // 2 segments of 5 bytes each
	result := SMEncode(data, segSize)

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
	result := SMEncode(data, 0)

	// With default 4MB segment size, 100 bytes should be 1 segment
	decoded, err := SMDecode(result.EncodedData)
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

	result := SMEncode(data, segSize)
	require.Equal(t, 67, len(result.EncodedData))
}

func TestEncodeCRC64MatchesSharedTable(t *testing.T) {
	data := []byte("CRC validation test data")
	expectedCRC := crc64.Checksum(data, CRC64Table)

	result := SMEncode(data, 0)
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
		t.Run(tt.name, func(t *testing.T) {
			_, err := SMDecode(tt.data)
			require.Error(t, err)
			if tt.errText != "" {
				require.Contains(t, err.Error(), tt.errText)
			}
		})
	}
}

func TestDecodeSegmentLengthExceedsSupportedSize(t *testing.T) {
	data := makeCorruptedSM([]byte("test"), func(d []byte) {
		binary.LittleEndian.PutUint64(d[15:23], ^uint64(0))
	})

	_, err := SMDecode(data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "exceeds supported size")
}

// makeCorruptedSM encodes data then applies a corruption function on the result.
func makeCorruptedSM(data []byte, corrupt func([]byte)) []byte {
	result := SMEncode(data, 0)
	smData := make([]byte, len(result.EncodedData))
	copy(smData, result.EncodedData)
	corrupt(smData)
	return smData
}

func TestEncoderReadSeekClose(t *testing.T) {
	data := []byte("encoder test data")
	enc := NewSMEncoder(bytes.NewReader(data), int64(len(data)), 0)

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
	decoded, err := SMDecode(buf)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
}

func TestEncoderAsReadSeekCloser(t *testing.T) {
	data := []byte("interface compliance test")
	enc := NewSMEncoder(bytes.NewReader(data), int64(len(data)), 0)

	var _ io.ReadSeekCloser = enc

	allData, err := io.ReadAll(enc)
	require.NoError(t, err)
	require.Equal(t, int(enc.EncodedLength()), len(allData))
}

func TestDecoderReadClose(t *testing.T) {
	data := []byte("decoder test with some content here")
	result := SMEncode(data, 10)

	body := io.NopCloser(bytes.NewReader(result.EncodedData))
	dec := NewSMDecoder(body)

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
	dec := NewSMDecoder(body)

	_, err := io.ReadAll(dec)
	require.Error(t, err)
}

func TestDecoderDecodeResultBeforeRead(t *testing.T) {
	data := []byte("test")
	result := SMEncode(data, 0)
	body := io.NopCloser(bytes.NewReader(result.EncodedData))
	dec := NewSMDecoder(body)

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

func TestStreamingEncoderMatchesInMemoryEncode(t *testing.T) {
	data := []byte("ABCDEFGHIJ") // 10 bytes
	segSize := 5                 // 2 segments

	// In-memory encode
	inMemResult := SMEncode(data, segSize)

	// Streaming encode
	enc := NewSMEncoder(bytes.NewReader(data), int64(len(data)), segSize)
	streamResult, err := io.ReadAll(enc)
	require.NoError(t, err)

	require.Equal(t, inMemResult.EncodedData, streamResult)
	require.Equal(t, int64(len(data)), enc.OriginalContentLength())
	require.Equal(t, int64(len(inMemResult.EncodedData)), enc.EncodedLength())
}

func TestStreamingEncoderMultiSegment(t *testing.T) {
	data := make([]byte, 350) // 4 segments with segSize=100
	for i := range data {
		data[i] = byte(i % 256)
	}
	enc := NewSMEncoder(bytes.NewReader(data), int64(len(data)), 100)
	encoded, err := io.ReadAll(enc)
	require.NoError(t, err)

	// Verify by decoding
	decoded, err := SMDecode(encoded)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
	require.Equal(t, uint16(4), decoded.NumSegments)
}

func TestStreamingEncoderEmpty(t *testing.T) {
	data := []byte{}
	enc := NewSMEncoder(bytes.NewReader(data), 0, 0)
	encoded, err := io.ReadAll(enc)
	require.NoError(t, err)

	decoded, err := SMDecode(encoded)
	require.NoError(t, err)
	require.Equal(t, 0, len(decoded.Data))
	require.Equal(t, uint16(1), decoded.NumSegments)
}

func TestStreamingEncoderSmallReads(t *testing.T) {
	// Test that the encoder works correctly with very small read buffers
	data := []byte("small buffer test data here!")
	enc := NewSMEncoder(bytes.NewReader(data), int64(len(data)), 10)

	var result []byte
	buf := make([]byte, 3) // intentionally small buffer
	for {
		n, err := enc.Read(buf)
		result = append(result, buf[:n]...)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
	}

	decoded, err := SMDecode(result)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
}

func TestStreamingEncoderToStreamingDecoder(t *testing.T) {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 251)
	}

	// Encode with streaming encoder
	enc := NewSMEncoder(bytes.NewReader(data), int64(len(data)), 256)
	encoded, err := io.ReadAll(enc)
	require.NoError(t, err)

	// Decode with streaming decoder
	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(encoded)))
	rawData, err := io.ReadAll(dec)
	require.NoError(t, err)
	require.Equal(t, data, rawData)

	decResult := dec.DecodeResult()
	require.NotNil(t, decResult)
	require.Equal(t, SMVersion, decResult.Version)
	require.Equal(t, SMFlagCRC64, decResult.Flags)
	require.Equal(t, uint16(4), decResult.NumSegments)
}

func TestStreamingDecoderSmallReads(t *testing.T) {
	data := []byte("streaming decode with tiny buffer")
	smData := SMEncode(data, 10)

	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(smData.EncodedData)))

	var result []byte
	buf := make([]byte, 5) // intentionally small buffer
	for {
		n, err := dec.Read(buf)
		result = append(result, buf[:n]...)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
	}
	require.Equal(t, data, result)
}

func TestStreamingDecoderCRCMismatch(t *testing.T) {
	// Corrupt segment data in an SM payload
	smData := makeCorruptedSM([]byte("Hello, world!"), func(d []byte) { d[25] ^= 0xFF })

	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(smData)))
	_, err := io.ReadAll(dec)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CRC64 mismatch")
}

func TestStreamingDecoderTrailerCRCMismatch(t *testing.T) {
	// Corrupt trailer CRC
	smData := makeCorruptedSM([]byte("test"), func(d []byte) { d[len(d)-1] ^= 0xFF })

	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(smData)))
	_, err := io.ReadAll(dec)
	require.Error(t, err)
}

func TestStreamingDecoderBadVersion(t *testing.T) {
	smData := makeCorruptedSM([]byte("test"), func(d []byte) { d[0] = 99 })

	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(smData)))
	_, err := io.ReadAll(dec)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported structured message version")
}

func TestStreamingEncoderSeekSupport(t *testing.T) {
	data := []byte("seek test")
	enc := NewSMEncoder(bytes.NewReader(data), int64(len(data)), 0)

	// Seek(0, SeekEnd) returns encoded length
	pos, err := enc.Seek(0, io.SeekEnd)
	require.NoError(t, err)
	require.Equal(t, enc.EncodedLength(), pos)

	// Seek(0, SeekCurrent) at initial position returns 0
	pos, err = enc.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	require.Equal(t, int64(0), pos)

	// After reading some data, SeekCurrent returns exact position
	buf := make([]byte, 1)
	_, _ = enc.Read(buf)
	pos, err = enc.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	require.Equal(t, int64(1), pos)

	// Non-zero offset should fail
	_, err = enc.Seek(1, io.SeekStart)
	require.Error(t, err)

	// Non-zero offset with SeekEnd should fail
	_, err = enc.Seek(1, io.SeekEnd)
	require.Error(t, err)

	// Seek(0, SeekStart) resets successfully
	pos, err = enc.Seek(0, io.SeekStart)
	require.NoError(t, err)
	require.Equal(t, int64(0), pos)
}

func TestStreamingEncoderWorksWithValidateSeekableStream(t *testing.T) {
	data := []byte("validate seekable test data")
	enc := NewSMEncoder(bytes.NewReader(data), int64(len(data)), 0)

	// ValidateSeekableStreamAt0AndGetCount uses Seek(0, SeekCurrent), Seek(0, SeekEnd), Seek(0, SeekStart)
	count, err := ValidateSeekableStreamAt0AndGetCount(enc)
	require.NoError(t, err)
	require.Equal(t, enc.EncodedLength(), count)

	// After validation, encoder should still be at position 0 and readable
	encoded, err := io.ReadAll(enc)
	require.NoError(t, err)
	require.Equal(t, int(count), len(encoded))
}
