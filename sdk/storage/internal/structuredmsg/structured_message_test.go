//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package structuredmsg

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc64"
	"io"
	"net"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/internal"
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

	// internal.ValidateSeekableStreamAt0AndGetCount uses Seek(0, SeekCurrent), Seek(0, SeekEnd), Seek(0, SeekStart)
	count, err := internal.ValidateSeekableStreamAt0AndGetCount(enc)
	require.NoError(t, err)
	require.Equal(t, enc.EncodedLength(), count)

	// After validation, encoder should still be at position 0 and readable
	encoded, err := io.ReadAll(enc)
	require.NoError(t, err)
	require.Equal(t, int(count), len(encoded))
}

// --- Decoder error tests (matching .NET StructuredMessageDecodingStreamTests) ---

func TestDecodeBadVersion(t *testing.T) {
	smData := makeCorruptedSM([]byte("test data for version check"), func(d []byte) {
		d[0] = 0xFF
	})
	_, err := SMDecode(smData)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported structured message version: 255")
}

func TestDecodeBadSegmentCRC(t *testing.T) {
	data := make([]byte, 100)
	for i := range data {
		data[i] = byte(i)
	}
	smData := makeCorruptedSM(data, func(d []byte) {
		d[SMHeaderSize+SMSegmentHeaderSize+10] ^= 0xFF
	})
	_, err := SMDecode(smData)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CRC64 mismatch")
	require.Contains(t, err.Error(), "segment 1")
}

func TestDecodeBadMessageCRC(t *testing.T) {
	data := make([]byte, 50)
	for i := range data {
		data[i] = byte(i)
	}
	smData := makeCorruptedSM(data, func(d []byte) {
		d[len(d)-1] ^= 0xFF
	})
	_, err := SMDecode(smData)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CRC64 mismatch")
}

func TestDecodeWrongMessageLength(t *testing.T) {
	smData := makeCorruptedSM([]byte("test message length"), func(d []byte) {
		binary.LittleEndian.PutUint64(d[1:9], 123456789)
	})
	_, err := SMDecode(smData)
	require.Error(t, err)
	require.Contains(t, err.Error(), "length mismatch")
}

func TestDecodeWrongSegmentCountTooMany(t *testing.T) {
	data := []byte("test segment count")
	result := SMEncode(data, 0)
	smData := make([]byte, len(result.EncodedData))
	copy(smData, result.EncodedData)

	binary.LittleEndian.PutUint16(smData[11:13], 2)
	_, err := SMDecode(smData)
	require.Error(t, err)
}

func TestDecodeWrongSegmentCountTooFew(t *testing.T) {
	data := make([]byte, 200)
	for i := range data {
		data[i] = byte(i)
	}
	result := SMEncode(data, 50) // 4 segments
	smData := make([]byte, len(result.EncodedData))
	copy(smData, result.EncodedData)

	binary.LittleEndian.PutUint16(smData[11:13], 3)
	_, err := SMDecode(smData)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CRC64 mismatch")
}

func TestDecodeWrongSegmentNumber(t *testing.T) {
	data := []byte("test segment number check")
	result := SMEncode(data, 0)
	smData := make([]byte, len(result.EncodedData))
	copy(smData, result.EncodedData)

	binary.LittleEndian.PutUint16(smData[SMHeaderSize:SMHeaderSize+2], 123)
	_, err := SMDecode(smData)
	require.Error(t, err)
	require.Contains(t, err.Error(), "segment number mismatch")
	require.Contains(t, err.Error(), "expected 1, got 123")
}

func TestDecodeTruncatedStream(t *testing.T) {
	data := []byte("test truncation handling")
	result := SMEncode(data, 0)
	truncated := result.EncodedData[:len(result.EncodedData)-4]

	_, err := SMDecode(truncated)
	require.Error(t, err)
	require.Contains(t, err.Error(), "length mismatch")
}

func TestDecodeTruncatedSegmentFooter(t *testing.T) {
	data := []byte("test footer truncation")
	result := SMEncode(data, 0)
	smData := make([]byte, len(result.EncodedData))
	copy(smData, result.EncodedData)

	truncatedLen := len(smData) - 12
	truncated := smData[:truncatedLen]
	binary.LittleEndian.PutUint64(truncated[1:9], uint64(truncatedLen))
	_, err := SMDecode(truncated)
	require.Error(t, err)
	require.Contains(t, err.Error(), "insufficient data")
}

func TestDecodeVariousReadSizes(t *testing.T) {
	testCases := []struct {
		name    string
		dataLen int
		segSize int
	}{
		{"Small_DefaultSeg", 100, 0},
		{"NonAligned_SmallSeg", 2005, 512},
		{"Aligned_SmallSeg", 2048, 512},
		{"Large_SmallSeg", 8192, 512},
		{"SingleByte_SmallSeg", 1, 512},
		{"ExactSeg", 512, 512},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := make([]byte, tc.dataLen)
			for i := range data {
				data[i] = byte(i % 251)
			}

			result := SMEncode(data, tc.segSize)
			decoded, err := SMDecode(result.EncodedData)
			require.NoError(t, err)
			require.Equal(t, data, decoded.Data)
		})
	}
}

// --- Streaming Encoder -> Decoder Roundtrip Tests (matching .NET StructuredMessageStreamRoundtripTests) ---

func TestStreamingEncoderDecoderRoundtrip(t *testing.T) {
	testCases := []struct {
		name    string
		dataLen int
		segSize int
		readLen int
	}{
		{"2048_DefaultSeg_8KB", 2048, 0, 8192},
		{"2005_512Seg_512Read", 2005, 512, 512},
		{"2048_512Seg_530Read", 2048, 512, 530},
		{"2005_512Seg_3Read", 2005, 512, 3},
		{"100_50Seg_7Read", 100, 50, 7},
		{"1_1Seg_1Read", 1, 1, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := make([]byte, tc.dataLen)
			for i := range data {
				data[i] = byte(i % 251)
			}

			enc := NewSMEncoder(bytes.NewReader(data), int64(len(data)), tc.segSize)
			encodedData, err := io.ReadAll(enc)
			require.NoError(t, err)
			require.Equal(t, int(enc.EncodedLength()), len(encodedData))

			body := io.NopCloser(bytes.NewReader(encodedData))
			dec := NewSMDecoder(body)

			var decoded []byte
			buf := make([]byte, tc.readLen)
			for {
				n, readErr := dec.Read(buf)
				if n > 0 {
					decoded = append(decoded, buf[:n]...)
				}
				if readErr == io.EOF {
					break
				}
				require.NoError(t, readErr)
			}

			require.Equal(t, data, decoded)
		})
	}
}

func TestStreamingEncoderDecoderRoundtripLargeData(t *testing.T) {
	dataLen := 5 * 1024 * 1024
	segSize := 1024 * 1024

	data := make([]byte, dataLen)
	for i := range data {
		data[i] = byte(i % 251)
	}

	enc := NewSMEncoder(bytes.NewReader(data), int64(len(data)), segSize)
	encodedData, err := io.ReadAll(enc)
	require.NoError(t, err)

	body := io.NopCloser(bytes.NewReader(encodedData))
	dec := NewSMDecoder(body)

	decoded, err := io.ReadAll(dec)
	require.NoError(t, err)
	require.Equal(t, data, decoded)

	decResult := dec.DecodeResult()
	require.NotNil(t, decResult)
	require.Equal(t, uint16(5), decResult.NumSegments)
}

// --- Encoder Binary Format Tests (matching .NET StructuredMessageTests) ---

func TestEncodeStreamHeaderBinary(t *testing.T) {
	data := make([]byte, 1024)
	result := SMEncode(data, 0)
	smData := result.EncodedData

	require.Equal(t, byte(1), smData[0])

	msgLen := binary.LittleEndian.Uint64(smData[1:9])
	require.Equal(t, uint64(len(smData)), msgLen)

	flags := binary.LittleEndian.Uint16(smData[9:11])
	require.Equal(t, uint16(1), flags)

	numSegs := binary.LittleEndian.Uint16(smData[11:13])
	require.Equal(t, uint16(1), numSegs)
}

func TestEncodeSegmentHeaderBinary(t *testing.T) {
	data := make([]byte, 10)
	for i := range data {
		data[i] = byte(i)
	}
	result := SMEncode(data, 5) // 2 segments of 5 bytes each
	smData := result.EncodedData

	seg1Num := binary.LittleEndian.Uint16(smData[13:15])
	require.Equal(t, uint16(1), seg1Num)
	seg1Len := binary.LittleEndian.Uint64(smData[15:23])
	require.Equal(t, uint64(5), seg1Len)

	seg2Offset := 13 + 10 + 5 + 8
	seg2Num := binary.LittleEndian.Uint16(smData[seg2Offset : seg2Offset+2])
	require.Equal(t, uint16(2), seg2Num)
	seg2Len := binary.LittleEndian.Uint64(smData[seg2Offset+2 : seg2Offset+10])
	require.Equal(t, uint64(5), seg2Len)
}

func TestEncodeNonAlignedDataSize(t *testing.T) {
	testCases := []struct {
		dataLen int
		segSize int
		expSegs uint16
	}{
		{2005, 512, 4},
		{1, 512, 1},
		{513, 512, 2},
		{1023, 512, 2},
		{10000, 3000, 4},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d_%d", tc.dataLen, tc.segSize), func(t *testing.T) {
			data := make([]byte, tc.dataLen)
			for i := range data {
				data[i] = byte(i % 251)
			}

			result := SMEncode(data, tc.segSize)
			decoded, err := SMDecode(result.EncodedData)
			require.NoError(t, err)
			require.Equal(t, data, decoded.Data)
			require.Equal(t, tc.expSegs, decoded.NumSegments)
		})
	}
}

// --- Decoder via SMDecoder (streaming) error tests ---

func TestDecoderBadVersion(t *testing.T) {
	smData := makeCorruptedSM([]byte("test"), func(d []byte) {
		d[0] = 0xFF
	})
	body := io.NopCloser(bytes.NewReader(smData))
	dec := NewSMDecoder(body)
	_, err := io.ReadAll(dec)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported structured message version")
}

func TestDecoderBadSegmentCRC(t *testing.T) {
	data := make([]byte, 100)
	for i := range data {
		data[i] = byte(i)
	}
	smData := makeCorruptedSM(data, func(d []byte) {
		d[SMHeaderSize+SMSegmentHeaderSize+5] ^= 0xFF
	})
	body := io.NopCloser(bytes.NewReader(smData))
	dec := NewSMDecoder(body)
	_, err := io.ReadAll(dec)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CRC64 mismatch")
}

func TestDecoderBadMessageCRC(t *testing.T) {
	data := make([]byte, 50)
	for i := range data {
		data[i] = byte(i)
	}
	smData := makeCorruptedSM(data, func(d []byte) {
		d[len(d)-1] ^= 0xFF
	})
	body := io.NopCloser(bytes.NewReader(smData))
	dec := NewSMDecoder(body)
	_, err := io.ReadAll(dec)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CRC64 mismatch")
}

func TestDecoderWrongSegmentNumber(t *testing.T) {
	smData := makeCorruptedSM([]byte("seg num test"), func(d []byte) {
		binary.LittleEndian.PutUint16(d[SMHeaderSize:SMHeaderSize+2], 42)
	})
	body := io.NopCloser(bytes.NewReader(smData))
	dec := NewSMDecoder(body)
	_, err := io.ReadAll(dec)
	require.Error(t, err)
	require.Contains(t, err.Error(), "segment number mismatch")
}

func TestDecoderMultiSegmentCRCValidation(t *testing.T) {
	data := make([]byte, 200)
	for i := range data {
		data[i] = byte(i)
	}
	result := SMEncode(data, 50)
	smData := make([]byte, len(result.EncodedData))
	copy(smData, result.EncodedData)

	seg3DataStart := SMHeaderSize + 2*(SMSegmentHeaderSize+50+SMSegmentFooterSize) + SMSegmentHeaderSize
	smData[seg3DataStart+10] ^= 0xFF

	body := io.NopCloser(bytes.NewReader(smData))
	dec := NewSMDecoder(body)
	_, err := io.ReadAll(dec)
	require.Error(t, err)
	require.Contains(t, err.Error(), "segment 3")
	require.Contains(t, err.Error(), "CRC64 mismatch")
}

// Test segment size overflow protection: when content would produce > 65535 segments,
// NewSMEncoder should auto-increase segment size to keep numSegments within uint16 range.
func TestSMEncoderSegmentSizeOverflowProtection(t *testing.T) {
	// With 1-byte segment size and 70000 bytes of content, numSegments would be 70000.
	// The encoder should auto-increase segment size to ceil(70000/65535) = 2 bytes,
	// resulting in 35000 segments.
	contentLen := int64(70000)
	segmentSize := 1
	content := make([]byte, contentLen)
	for i := range content {
		content[i] = byte(i % 251)
	}

	enc := NewSMEncoder(bytes.NewReader(content), contentLen, segmentSize)

	// Verify the encoder was created successfully and numSegments <= 65535
	require.LessOrEqual(t, enc.numSegments, 65535)
	require.Greater(t, enc.numSegments, 0)

	// Read all encoded data
	encodedData, err := io.ReadAll(enc)
	require.NoError(t, err)
	require.Equal(t, int(enc.EncodedLength()), len(encodedData))

	// Decode and verify round-trip
	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(encodedData)))
	decodedData, err := io.ReadAll(dec)
	require.NoError(t, err)
	require.Equal(t, content, decodedData)
}

func TestSMEncoderSegmentSizeExactlyAtLimit(t *testing.T) {
	// With segment size 1 and exactly 65535 bytes, it should not overflow
	contentLen := int64(65535)
	segmentSize := 1
	content := make([]byte, contentLen)
	for i := range content {
		content[i] = byte(i % 199)
	}

	enc := NewSMEncoder(bytes.NewReader(content), contentLen, segmentSize)
	require.Equal(t, 65535, enc.numSegments)

	encodedData, err := io.ReadAll(enc)
	require.NoError(t, err)

	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(encodedData)))
	decodedData, err := io.ReadAll(dec)
	require.NoError(t, err)
	require.Equal(t, content, decodedData)
}

func TestSMEncoderSegmentSizeJustOverLimit(t *testing.T) {
	// With segment size 1 and 65536 bytes, it should auto-scale to segment size 2
	contentLen := int64(65536)
	segmentSize := 1
	content := make([]byte, contentLen)
	for i := range content {
		content[i] = byte(i % 199)
	}

	enc := NewSMEncoder(bytes.NewReader(content), contentLen, segmentSize)
	require.LessOrEqual(t, enc.numSegments, 65535)
	require.Greater(t, enc.segmentSize, 1)

	encodedData, err := io.ReadAll(enc)
	require.NoError(t, err)

	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(encodedData)))
	decodedData, err := io.ReadAll(dec)
	require.NoError(t, err)
	require.Equal(t, content, decodedData)
}

// segBoundaryReader is a test io.ReadSeeker whose Read returns err exactly when its data is
// exhausted, allowing simulation of a source that returns bytes together with an error at a
// segment boundary.
type segBoundaryReader struct {
	data []byte
	err  error
	off  int
}

func (r *segBoundaryReader) Read(p []byte) (int, error) {
	if r.off >= len(r.data) {
		return 0, r.err
	}
	n := copy(p, r.data[r.off:])
	r.off += n
	if r.off >= len(r.data) {
		return n, r.err
	}
	return n, nil
}

func (r *segBoundaryReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		r.off = int(offset)
	case io.SeekCurrent:
		r.off += int(offset)
	case io.SeekEnd:
		r.off = len(r.data) + int(offset)
	}
	return int64(r.off), nil
}

// TestStreamingDecoderValidatesTrailerOnExactBufferFill guards against a decoder that returns the
// final payload byte without draining and validating the trailing segment footer/message trailer in
// the same Read. A bounded reader (e.g. RetryReader) could otherwise report EOF and skip validation
// when the payload exactly fills the caller's buffer.
func TestStreamingDecoderValidatesTrailerOnExactBufferFill(t *testing.T) {
	data := make([]byte, 2048)
	for i := range data {
		data[i] = byte(i % 251)
	}

	// Corrupt only the message trailer CRC (last 8 bytes); the segment data and footer remain valid.
	corrupted := makeCorruptedSM(data, func(d []byte) { d[len(d)-1] ^= 0xFF })

	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(corrupted)))

	// Buffer sized exactly to the payload length so the last payload byte fills it exactly.
	buf := make([]byte, len(data))
	n, err := dec.Read(buf)
	require.Equal(t, len(data), n)
	// The trailing framing must have been drained and validated in this same call.
	require.Error(t, err)
	require.Contains(t, err.Error(), "trailer CRC64 mismatch")
}

// TestStreamingDecoderExactBufferFillValid verifies the happy path for the same exact-fill boundary:
// a valid message returns all payload bytes plus a nil error on the fill read, and a subsequent read
// reports io.EOF with no extra bytes.
func TestStreamingDecoderExactBufferFillValid(t *testing.T) {
	data := make([]byte, 2048)
	for i := range data {
		data[i] = byte(i % 251)
	}

	valid := SMEncode(data, 0).EncodedData
	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(valid)))

	buf := make([]byte, len(data))
	n, err := dec.Read(buf)
	require.NoError(t, err)
	require.Equal(t, len(data), n)
	require.Equal(t, data, buf)

	n, err = dec.Read(make([]byte, 8))
	require.Equal(t, 0, n)
	require.ErrorIs(t, err, io.EOF)
}

// TestEncoderPropagatesNonEOFErrorAtSegmentBoundary verifies that when the inner reader returns data
// together with a non-EOF error and those bytes exactly complete a segment, the encoder propagates
// the error instead of silently emitting a valid trailer.
func TestEncoderPropagatesNonEOFErrorAtSegmentBoundary(t *testing.T) {
	data := []byte("ABCD") // exactly one 4-byte segment
	sentinel := errors.New("read failure at segment boundary")

	enc := NewSMEncoder(&segBoundaryReader{data: data, err: sentinel}, int64(len(data)), len(data))
	_, err := io.ReadAll(enc)
	require.Error(t, err)
	require.ErrorIs(t, err, sentinel)
}

// TestEncoderAcceptsDataWithEOFAtFinalBoundary verifies that a source returning its final bytes
// together with io.EOF (valid io.Reader behavior) still produces a correctly encoded message rather
// than a spurious error.
func TestEncoderAcceptsDataWithEOFAtFinalBoundary(t *testing.T) {
	data := []byte("ABCDEFGHIJ") // 10 bytes across two 5-byte segments
	segSize := 5

	enc := NewSMEncoder(&segBoundaryReader{data: data, err: io.EOF}, int64(len(data)), segSize)
	encoded, err := io.ReadAll(enc)
	require.NoError(t, err)

	decoded, err := SMDecode(encoded)
	require.NoError(t, err)
	require.Equal(t, data, decoded.Data)
}

// TestDecoderRejectsMissingCRC64Flag verifies that a structured message whose CRC64 flag is not set
// is rejected. The body is negotiated with properties=crc64, so a missing flag (which would cause the
// decoder to skip all footer/trailer validation) must be treated as an error by both the streaming
// decoder and the one-shot SMDecode.
func TestDecoderRejectsMissingCRC64Flag(t *testing.T) {
	data := []byte("structured message body that should be validated")

	// Clear the CRC64 flag bits (flags occupy bytes [9:11]) while leaving msgLen and framing intact.
	noFlag := makeCorruptedSM(data, func(d []byte) {
		d[9] = 0
		d[10] = 0
	})

	// Streaming decoder must reject it.
	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(noFlag)))
	_, err := io.ReadAll(dec)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing the required CRC64 flag")

	// One-shot decode must reject it too.
	_, err = SMDecode(noFlag)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing the required CRC64 flag")
}

// TestDecodeRejectsTrailingBytes verifies that SMDecode rejects a payload whose declared length
// matches the input but whose parsed segments/trailer end before the end of the buffer, leaving
// unexamined trailing bytes. This guards against declaring fewer segments, placing a valid message
// CRC at the resulting trailer offset, and appending arbitrary bytes.
func TestDecodeRejectsTrailingBytes(t *testing.T) {
	data := []byte("payload that will be followed by junk")
	encoded := SMEncode(data, 0).EncodedData

	// Append arbitrary trailing bytes and bump the header's declared message length to match, so the
	// initial length check passes but the parsed message ends before len(smData).
	junk := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	tampered := make([]byte, 0, len(encoded)+len(junk))
	tampered = append(tampered, encoded...)
	tampered = append(tampered, junk...)
	binary.LittleEndian.PutUint64(tampered[1:9], uint64(len(tampered)))

	_, err := SMDecode(tampered)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected trailing bytes")
}

// TestEncoderPrematureEOFIsUnexpected verifies that when the source ends before the declared content
// length, the encoder surfaces io.ErrUnexpectedEOF rather than a plain io.EOF, so callers such as
// io.ReadAll do not accept a truncated structured message as success.
func TestEncoderPrematureEOFIsUnexpected(t *testing.T) {
	// Declare 10 bytes of content but only supply 4 before EOF.
	data := []byte("ABCD")
	enc := NewSMEncoder(&segBoundaryReader{data: data, err: io.EOF}, 10, 0)

	_, err := io.ReadAll(enc)
	require.Error(t, err)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

// TestDecoderPropagatesNetErrorAtExactSegmentCompletion verifies that when the source returns bytes
// together with a net.Error and those bytes exactly complete a segment, the decoder propagates the
// error (wrapped with %w) instead of silently advancing to the footer state.
func TestDecoderPropagatesNetErrorAtExactSegmentCompletion(t *testing.T) {
	data := []byte("ABCD") // 4 bytes, one segment
	segSize := len(data)

	// Encode a valid structured message to get the framing.
	encoded := SMEncode(data, segSize).EncodedData

	// Build a reader that returns the full encoded bytes but injects a net.Error exactly when the
	// segment data is exhausted. The segBoundaryReader returns (n, err) when its data runs out.
	netErr := &net.DNSError{IsTemporary: true}

	// We need a reader that returns SM framing normally but injects the error at the segment data boundary.
	src := &segDataNetErrorReader{
		encoded: encoded,
		segSize: segSize,
		netErr:  netErr,
	}

	dec := NewSMDecoder(io.NopCloser(src))
	_, err := io.ReadAll(dec)
	require.Error(t, err)
	require.ErrorIs(t, err, netErr)
}

// segDataNetErrorReader serves an encoded SM stream but injects a net.Error exactly when the
// segment data bytes are exhausted (bytes + error at segment completion boundary).
type segDataNetErrorReader struct {
	encoded  []byte
	segSize  int
	netErr   net.Error
	off      int
	injected bool
}

func (r *segDataNetErrorReader) Read(p []byte) (int, error) {
	if r.off >= len(r.encoded) {
		return 0, io.EOF
	}
	// Calculate the end of the first segment's data within the encoded stream.
	// Header (13 bytes) + segment header (6 bytes) + segment data (segSize bytes).
	segDataEnd := SMHeaderSize + SMSegmentHeaderSize + r.segSize
	if !r.injected && r.off < segDataEnd {
		// Serve up to the end of segment data, then inject the error.
		avail := segDataEnd - r.off
		n := copy(p, r.encoded[r.off:r.off+avail])
		r.off += n
		if r.off >= segDataEnd {
			r.injected = true
			return n, r.netErr
		}
		return n, nil
	}
	n := copy(p, r.encoded[r.off:])
	r.off += n
	if r.off >= len(r.encoded) {
		return n, io.EOF
	}
	return n, nil
}

// TestDecoderPrematureEOFDuringSegmentData verifies that when the source returns io.EOF while
// segment data bytes are still expected, the decoder converts it to io.ErrUnexpectedEOF so the
// error is retryable by RetryReader.
func TestDecoderPrematureEOFDuringSegmentData(t *testing.T) {
	data := []byte("ABCDEFGH") // 8 bytes
	segSize := len(data)

	encoded := SMEncode(data, segSize).EncodedData

	// Truncate the encoded stream partway through the segment data (cut after header + seg header + 4 bytes).
	truncateAt := SMHeaderSize + SMSegmentHeaderSize + 4
	truncated := encoded[:truncateAt]

	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(truncated)))
	_, err := io.ReadAll(dec)
	require.Error(t, err)
	require.Contains(t, err.Error(), "segment")
}

// TestDecoderEOFExactlyCompletingSegmentIsUnexpected verifies that a source returning (n, io.EOF)
// where n exactly completes the segment is treated as an error (footer/trailer must still follow).
func TestDecoderEOFExactlyCompletingSegmentIsUnexpected(t *testing.T) {
	data := []byte("ABCD")
	segSize := len(data)

	encoded := SMEncode(data, segSize).EncodedData

	// Build a reader that returns (n, io.EOF) exactly when segment data is exhausted,
	// before footer/trailer can be read.
	segDataEnd := SMHeaderSize + SMSegmentHeaderSize + segSize
	src := &segBoundaryEOFReader{
		encoded:    encoded,
		segDataEnd: segDataEnd,
	}

	dec := NewSMDecoder(io.NopCloser(src))
	_, err := io.ReadAll(dec)
	require.Error(t, err)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

// segBoundaryEOFReader returns (n, io.EOF) exactly when the segment data is exhausted, simulating
// a source that ends at the segment data boundary before footer/trailer framing.
type segBoundaryEOFReader struct {
	encoded    []byte
	segDataEnd int
	off        int
}

func (r *segBoundaryEOFReader) Read(p []byte) (int, error) {
	if r.off >= len(r.encoded) {
		return 0, io.EOF
	}
	if r.off < r.segDataEnd {
		end := r.segDataEnd
		if end > r.off+len(p) {
			end = r.off + len(p)
		}
		n := copy(p, r.encoded[r.off:end])
		r.off += n
		if r.off >= r.segDataEnd {
			return n, io.EOF
		}
		return n, nil
	}
	n := copy(p, r.encoded[r.off:])
	r.off += n
	if r.off >= len(r.encoded) {
		return n, io.EOF
	}
	return n, nil
}

// TestFillFrameTruncatedFooterIsRetryable verifies that a premature EOF during footer/trailer
// framing reads wraps io.ErrUnexpectedEOF so RetryReader can classify it as retryable.
func TestFillFrameTruncatedFooterIsRetryable(t *testing.T) {
	data := []byte("ABCDEFGH")
	segSize := len(data)

	encoded := SMEncode(data, segSize).EncodedData

	// Truncate after segment data, partway through the segment footer.
	// Header (13) + seg header (6) + seg data (8) + partial footer (4 of 8).
	truncateAt := SMHeaderSize + SMSegmentHeaderSize + segSize + 4
	require.Less(t, truncateAt, len(encoded))
	truncated := encoded[:truncateAt]

	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(truncated)))
	_, err := io.ReadAll(dec)
	require.Error(t, err)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

// TestFillFrameTruncatedTrailerIsRetryable verifies that a premature EOF during message trailer
// reads wraps io.ErrUnexpectedEOF.
func TestFillFrameTruncatedTrailerIsRetryable(t *testing.T) {
	data := []byte("ABCDEFGH")
	segSize := len(data)

	encoded := SMEncode(data, segSize).EncodedData

	// Truncate after segment footer, partway through the message trailer.
	// Header (13) + seg header (6) + seg data (8) + seg footer (8) + partial trailer (4 of 8).
	truncateAt := SMHeaderSize + SMSegmentHeaderSize + segSize + SMSegmentFooterSize + 4
	require.Less(t, truncateAt, len(encoded))
	truncated := encoded[:truncateAt]

	dec := NewSMDecoder(io.NopCloser(bytes.NewReader(truncated)))
	_, err := io.ReadAll(dec)
	require.Error(t, err)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}
