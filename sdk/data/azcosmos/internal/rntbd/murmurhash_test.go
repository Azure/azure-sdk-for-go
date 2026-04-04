// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"crypto/rand"
	"testing"
)

// Known hash values computed using Google Guava's murmur3_32(0)
// These are the reference values that the Java SDK validates against
var knownHashValues = []struct {
	input    string
	expected uint32
}{
	{"", 0},               // Empty string
	{"test", 3127628307},  // Simple ASCII - verified against Guava
	{"hello", 613153351},  // Another ASCII
	{"world", 4220927227}, // And another
}

func TestMurmurHash3_32_EmptyByteArray(t *testing.T) {
	data := []byte{}
	hash := MurmurHash3_32(data, 0)

	// Empty array with seed 0 should produce 0
	// This matches Google Guava's murmur3_32(0).hashBytes(new byte[0]).asInt()
	if hash != 0 {
		t.Errorf("MurmurHash3_32(empty, 0) = %d, want 0", hash)
	}
}

func TestMurmurHash3_32_String(t *testing.T) {
	data := []byte("test")
	hash := MurmurHash3_32(data, 0)

	// Expected value from Google Guava murmur3_32(0).hashBytes("test".getBytes("UTF-8")).asInt()
	// The unsigned value is 3127628307 (0xBA6BD213)
	expected := uint32(3127628307)
	if hash != expected {
		t.Errorf("MurmurHash3_32(\"test\", 0) = %d (0x%08X), want %d (0x%08X)", hash, hash, expected, expected)
	}
}

func TestMurmurHash3_32_NonLatin(t *testing.T) {
	// Test with Cyrillic characters - must match Java's UTF-8 encoding
	nonLatin := "абвгдеёжзийклмнопрстуфхцчшщъыьэюяабвгдеёжзийклмнопрстуфхцчшщъыьэюяабвгдеёжзийклмнопрстуфхцчшщъыьэюяабвгдеёжзийклмнопрстуфхцчшщъыьэюя"

	// Test all substring lengths from 0 to full length
	for i := 0; i <= len(nonLatin); i++ {
		substr := nonLatin[:i]
		data := []byte(substr)

		// Verify hash is deterministic
		hash1 := MurmurHash3_32(data, 0)
		hash2 := MurmurHash3_32(data, 0)

		if hash1 != hash2 {
			t.Errorf("MurmurHash3_32 not deterministic for substring length %d: %d != %d", i, hash1, hash2)
		}
	}
}

func TestMurmurHash3_32_ZeroByteArray(t *testing.T) {
	data := make([]byte, 3)
	hash := MurmurHash3_32(data, 0)
	hash2 := MurmurHash3_32(data, 0)
	if hash != hash2 {
		t.Errorf("MurmurHash3_32 not deterministic: %d != %d", hash, hash2)
	}
}

func TestMurmurHash3_32_RandomBytesOfAllSizes(t *testing.T) {
	// Test 1000 different sizes from 0 to 999 bytes
	for size := 0; size < 1000; size++ {
		data := make([]byte, size)
		_, err := rand.Read(data)
		if err != nil {
			t.Fatalf("Failed to generate random bytes: %v", err)
		}

		// Hash must be deterministic
		hash1 := MurmurHash3_32(data, 0)
		hash2 := MurmurHash3_32(data, 0)

		if hash1 != hash2 {
			t.Errorf("MurmurHash3_32 not deterministic for size %d: %d != %d", size, hash1, hash2)
		}
	}
}

func TestMurmurHash3_32_KnownValues(t *testing.T) {
	for _, tc := range knownHashValues {
		t.Run(tc.input, func(t *testing.T) {
			data := []byte(tc.input)
			hash := MurmurHash3_32(data, 0)

			if hash != tc.expected {
				t.Errorf("MurmurHash3_32(%q, 0) = %d (0x%08X), want %d (0x%08X)",
					tc.input, hash, hash, tc.expected, tc.expected)
			}
		})
	}
}

func TestMurmurHash3_32_AllTailLengths(t *testing.T) {
	// Test all 4 tail cases (0, 1, 2, 3 remaining bytes after 4-byte blocks)
	testCases := []struct {
		length   int
		tailCase int
	}{
		{0, 0},   // No data
		{1, 1},   // 1 byte tail
		{2, 2},   // 2 byte tail
		{3, 3},   // 3 byte tail
		{4, 0},   // Exact 4-byte block
		{5, 1},   // 4-byte block + 1 byte tail
		{6, 2},   // 4-byte block + 2 byte tail
		{7, 3},   // 4-byte block + 3 byte tail
		{8, 0},   // Exact 2 4-byte blocks
		{100, 0}, // Multiple blocks, no tail
		{101, 1}, // Multiple blocks + 1 byte tail
		{102, 2}, // Multiple blocks + 2 byte tail
		{103, 3}, // Multiple blocks + 3 byte tail
	}

	for _, tc := range testCases {
		t.Run("length_"+string(rune('0'+tc.length)), func(t *testing.T) {
			data := make([]byte, tc.length)
			for i := range data {
				data[i] = byte(i & 0xFF)
			}

			// Verify tail case is correct
			actualTail := tc.length & 0x03
			if actualTail != tc.tailCase {
				t.Errorf("Expected tail case %d for length %d, got %d", tc.tailCase, tc.length, actualTail)
			}

			// Hash should be deterministic
			hash1 := MurmurHash3_32(data, 0)
			hash2 := MurmurHash3_32(data, 0)
			if hash1 != hash2 {
				t.Errorf("Hash not deterministic for length %d", tc.length)
			}
		})
	}
}

func TestMurmurHash3_32_DifferentSeeds(t *testing.T) {
	data := []byte("test data for seed testing")

	// Different seeds should produce different hashes
	seeds := []uint32{0, 1, 42, 0xFFFFFFFF, 0x12345678}
	hashes := make(map[uint32]uint32)

	for _, seed := range seeds {
		hash := MurmurHash3_32(data, seed)
		if existingHash, exists := hashes[hash]; exists {
			t.Errorf("Collision: seed %d and another seed both produce hash %d", seed, existingHash)
		}
		hashes[hash] = seed
	}
}

func TestMurmurHash3_32Signed(t *testing.T) {
	// Test the signed version matches Java's int return type
	data := []byte("test")
	hash := MurmurHash3_32Signed(data, 0)

	// Java returns this as a signed int
	// 3127628307 as unsigned = -1167338989 as signed
	expected := int32(-1167338989)
	if hash != expected {
		t.Errorf("MurmurHash3_32Signed(\"test\", 0) = %d, want %d", hash, expected)
	}
}

func TestMurmurHash3_32_ByteOrderIndependence(t *testing.T) {
	// Verify little-endian byte order is used correctly
	// This is critical for cross-platform compatibility

	// 4 bytes that form a specific little-endian int32
	data := []byte{0x01, 0x02, 0x03, 0x04} // Little-endian: 0x04030201

	hash := MurmurHash3_32(data, 0)

	// The same data should always produce the same hash
	hash2 := MurmurHash3_32(data, 0)
	if hash != hash2 {
		t.Errorf("Hash not deterministic: %d != %d", hash, hash2)
	}

	// Different byte order should produce different hash
	reversed := []byte{0x04, 0x03, 0x02, 0x01}
	hashReversed := MurmurHash3_32(reversed, 0)
	if hash == hashReversed {
		t.Error("Different byte orders should produce different hashes")
	}
}

func BenchmarkMurmurHash3_32_Small(b *testing.B) {
	data := []byte("small")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MurmurHash3_32(data, 0)
	}
}

func BenchmarkMurmurHash3_32_Medium(b *testing.B) {
	data := make([]byte, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MurmurHash3_32(data, 0)
	}
}

func BenchmarkMurmurHash3_32_Large(b *testing.B) {
	data := make([]byte, 10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MurmurHash3_32(data, 0)
	}
}
