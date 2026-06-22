// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"hash/crc64"
	"io"
	"testing"
)

// data sizes used across benchmarks
var benchSizes = []struct {
	name string
	size int
}{
	{"1KB", 1 * 1024},
	{"64KB", 64 * 1024},
	{"1MB", 1 * 1024 * 1024},
	{"4MB", 4 * 1024 * 1024},
	{"16MB", 16 * 1024 * 1024},
	{"100MB", 100 * 1024 * 1024},
}

func makeBenchData(size int) []byte {
	data := make([]byte, size)
	_, _ = rand.Read(data)
	return data
}

// ── CRC64 baseline ──
// Measures raw CRC64 computation cost — the theoretical minimum overhead for any
// integrity check. SM benchmarks should be compared against this baseline.

func BenchmarkCRC64Checksum(b *testing.B) {
	for _, sz := range benchSizes {
		data := makeBenchData(sz.size)
		b.Run(sz.name, func(b *testing.B) {
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				crc64.Checksum(data, CRC64Table)
			}
		})
	}
}

// ── In-memory encode/decode ──

func BenchmarkSMEncode(b *testing.B) {
	for _, sz := range benchSizes {
		data := makeBenchData(sz.size)
		b.Run(sz.name, func(b *testing.B) {
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				SMEncode(data, 0)
			}
		})
	}
}

func BenchmarkSMDecode(b *testing.B) {
	for _, sz := range benchSizes {
		data := makeBenchData(sz.size)
		encoded := SMEncode(data, 0)
		b.Run(sz.name, func(b *testing.B) {
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = SMDecode(encoded.EncodedData)
			}
		})
	}
}

// ── Streaming encoder ──

func BenchmarkSMEncoderRead(b *testing.B) {
	for _, sz := range benchSizes {
		data := makeBenchData(sz.size)
		b.Run(sz.name, func(b *testing.B) {
			b.SetBytes(int64(sz.size))
			buf := make([]byte, 32*1024) // 32KB read buffer (typical io.Copy size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				reader := bytes.NewReader(data)
				encoder := NewSMEncoder(reader, int64(sz.size), 0)
				for {
					_, err := encoder.Read(buf)
					if err == io.EOF {
						break
					}
				}
			}
		})
	}
}

// ── Streaming decoder ──

func BenchmarkSMDecoderRead(b *testing.B) {
	for _, sz := range benchSizes {
		data := makeBenchData(sz.size)
		encoded := SMEncode(data, 0)
		b.Run(sz.name, func(b *testing.B) {
			b.SetBytes(int64(sz.size))
			buf := make([]byte, 32*1024)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				body := io.NopCloser(bytes.NewReader(encoded.EncodedData))
				decoder := NewSMDecoder(body)
				for {
					_, err := decoder.Read(buf)
					if err == io.EOF {
						break
					}
				}
			}
		})
	}
}

// ── Encoder seek/retry ──
// Measures the cost of seeking back to 0 and re-reading (simulates retry).

func BenchmarkSMEncoderSeekAndReread(b *testing.B) {
	sizes := []struct {
		name string
		size int
	}{
		{"1KB", 1 * 1024},
		{"1MB", 1 * 1024 * 1024},
		{"4MB", 4 * 1024 * 1024},
	}
	for _, sz := range sizes {
		data := makeBenchData(sz.size)
		b.Run(sz.name, func(b *testing.B) {
			b.SetBytes(int64(sz.size))
			buf := make([]byte, 32*1024)
			reader := bytes.NewReader(data)
			encoder := NewSMEncoder(reader, int64(sz.size), 0)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = encoder.Seek(0, io.SeekStart)
				for {
					_, err := encoder.Read(buf)
					if err == io.EOF {
						break
					}
				}
			}
		})
	}
}

// ── Encode + Decode round-trip ──

func BenchmarkSMRoundTrip(b *testing.B) {
	for _, sz := range benchSizes {
		data := makeBenchData(sz.size)
		b.Run(sz.name, func(b *testing.B) {
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := SMEncode(data, 0)
				_, _ = SMDecode(result.EncodedData)
			}
		})
	}
}

// ── Streaming round-trip (encoder → decoder) ──

func BenchmarkSMStreamingRoundTrip(b *testing.B) {
	for _, sz := range benchSizes {
		data := makeBenchData(sz.size)
		b.Run(sz.name, func(b *testing.B) {
			b.SetBytes(int64(sz.size))
			buf := make([]byte, 32*1024)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Encode: read all from encoder into a buffer
				reader := bytes.NewReader(data)
				encoder := NewSMEncoder(reader, int64(sz.size), 0)
				var encoded bytes.Buffer
				encoded.Grow(sz.size + sz.size/10) // ~10% overhead estimate
				for {
					n, err := encoder.Read(buf)
					if n > 0 {
						encoded.Write(buf[:n])
					}
					if err == io.EOF {
						break
					}
				}
				// Decode: stream through decoder
				body := io.NopCloser(bytes.NewReader(encoded.Bytes()))
				decoder := NewSMDecoder(body)
				for {
					_, err := decoder.Read(buf)
					if err == io.EOF {
						break
					}
				}
			}
		})
	}
}

// ── ComputeCRC64 (2-pass) vs SM CRC64 (single-pass) ──
// ComputeCRC64 reads the entire body to compute CRC, then the body is re-read for the HTTP send.
// SM CRC64 computes CRC on-the-fly during a single read pass.

func BenchmarkComputeCRC64VsSMEncoder(b *testing.B) {
	for _, sz := range benchSizes {
		data := makeBenchData(sz.size)

		// Simulate ComputeCRC64: ReadAll + Checksum + re-read (2-pass)
		b.Run(fmt.Sprintf("ComputeCRC64_2pass/%s", sz.name), func(b *testing.B) {
			b.SetBytes(int64(sz.size))
			buf := make([]byte, 32*1024)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Pass 1: read all + compute CRC
				reader := bytes.NewReader(data)
				allData, _ := io.ReadAll(reader)
				crc64.Checksum(allData, CRC64Table)

				// Pass 2: re-read the data (simulates HTTP send)
				reader2 := bytes.NewReader(allData)
				for {
					_, err := reader2.Read(buf)
					if err == io.EOF {
						break
					}
				}
			}
		})

		// SM Encoder: single-pass streaming CRC
		b.Run(fmt.Sprintf("SMEncoder_1pass/%s", sz.name), func(b *testing.B) {
			b.SetBytes(int64(sz.size))
			buf := make([]byte, 32*1024)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				reader := bytes.NewReader(data)
				encoder := NewSMEncoder(reader, int64(sz.size), 0)
				for {
					_, err := encoder.Read(buf)
					if err == io.EOF {
						break
					}
				}
			}
		})
	}
}

// ── Segment size impact ──
// Measures how different segment sizes affect encoding performance for a fixed data size.

func BenchmarkSMEncodeSegmentSizes(b *testing.B) {
	dataSize := 16 * 1024 * 1024 // 16MB
	data := makeBenchData(dataSize)
	segSizes := []struct {
		name string
		size int
	}{
		{"256KB", 256 * 1024},
		{"1MB", 1 * 1024 * 1024},
		{"4MB_default", 4 * 1024 * 1024},
		{"8MB", 8 * 1024 * 1024},
		{"16MB", 16 * 1024 * 1024},
	}
	for _, seg := range segSizes {
		b.Run(fmt.Sprintf("data16MB_seg%s", seg.name), func(b *testing.B) {
			b.SetBytes(int64(dataSize))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				SMEncode(data, seg.size)
			}
		})
	}
}

// ── Overhead measurement ──
// Compares raw io.Copy throughput vs SM encoder throughput to quantify framing overhead.

func BenchmarkBaselineIOCopy(b *testing.B) {
	for _, sz := range benchSizes {
		data := makeBenchData(sz.size)
		b.Run(sz.name, func(b *testing.B) {
			b.SetBytes(int64(sz.size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				reader := bytes.NewReader(data)
				_, _ = io.Copy(io.Discard, reader)
			}
		})
	}
}
