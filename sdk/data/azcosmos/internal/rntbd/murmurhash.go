// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

// MurmurHash3_32 computes the MurmurHash3 x86 32-bit hash.
//
// The MurmurHash3 algorithm was created by Austin Appleby and placed in the public domain.
// This Go implementation matches the Java SDK implementation exactly to ensure
// cross-platform compatibility for partition key routing.
//
// See: http://github.com/yonik/java_util for the original Java implementation
func MurmurHash3_32(data []byte, seed uint32) uint32 {
	const c1 uint32 = 0xcc9e2d51
	const c2 uint32 = 0x1b873593

	length := len(data)
	h1 := seed
	roundedEnd := length & 0xfffffffc // round down to 4 byte block

	// Process 4-byte blocks
	for i := 0; i < roundedEnd; i += 4 {
		// little endian load order
		k1 := uint32(data[i]) |
			uint32(data[i+1])<<8 |
			uint32(data[i+2])<<16 |
			uint32(data[i+3])<<24

		k1 *= c1
		k1 = rotl32(k1, 15)
		k1 *= c2

		h1 ^= k1
		h1 = rotl32(h1, 13)
		h1 = h1*5 + 0xe6546b64
	}

	// Process tail (remaining 1-3 bytes)
	var k1 uint32 = 0
	switch length & 0x03 {
	case 3:
		k1 = uint32(data[roundedEnd+2]) << 16
		fallthrough
	case 2:
		k1 |= uint32(data[roundedEnd+1]) << 8
		fallthrough
	case 1:
		k1 |= uint32(data[roundedEnd])
		k1 *= c1
		k1 = rotl32(k1, 15)
		k1 *= c2
		h1 ^= k1
	}

	// Finalization
	h1 ^= uint32(length)

	// fmix32
	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return h1
}

// rotl32 performs a 32-bit rotate left
func rotl32(x uint32, r uint) uint32 {
	return (x << r) | (x >> (32 - r))
}

// MurmurHash3_32Signed returns the hash as a signed int32, matching Java's behavior
func MurmurHash3_32Signed(data []byte, seed int) int32 {
	return int32(MurmurHash3_32(data, uint32(seed)))
}
