// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package epk

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

// EffectivePartitionKey holds the computed EPK hash for a partition key value.
// The EPK is a hex-encoded string that determines which physical partition range
// a logical partition key maps to.
type EffectivePartitionKey struct {
	// EPK is the hex-encoded effective partition key hash string, comparable
	// against partitionKeyRange.minInclusive / maxExclusive boundaries.
	EPK string
}

// Partition key component type markers used in hash input encoding.
const (
	componentUndefined byte = 0x00
	componentNull      byte = 0x01
	componentFalse     byte = 0x02
	componentTrue      byte = 0x03
	componentNumber    byte = 0x05
	componentString    byte = 0x08
)

// UndefinedMarker is a sentinel type representing the "undefined" partition key value.
type UndefinedMarker struct{}

// --- MurmurHash3 ---

// murmurhash3_32 computes a 32-bit MurmurHash3 hash.
func murmurhash3_32(data []byte, seed uint32) uint32 {
	const (
		c1 uint32 = 0xcc9e2d51
		c2 uint32 = 0x1b873593
	)
	h1 := seed
	length := uint32(len(data))
	roundedEnd := length & 0xFFFFFFFC

	for i := uint32(0); i < roundedEnd; i += 4 {
		k1 := uint32(data[i]) | uint32(data[i+1])<<8 | uint32(data[i+2])<<16 | uint32(data[i+3])<<24
		k1 *= c1
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= c2

		h1 ^= k1
		h1 = (h1 << 13) | (h1 >> 19)
		h1 = h1*5 + 0xe6546b64
	}

	// tail
	k1 := uint32(0)
	switch length & 3 {
	case 3:
		k1 ^= uint32(data[roundedEnd+2]) << 16
		fallthrough
	case 2:
		k1 ^= uint32(data[roundedEnd+1]) << 8
		fallthrough
	case 1:
		k1 ^= uint32(data[roundedEnd])
		k1 *= c1
		k1 = (k1 << 15) | (k1 >> 17)
		k1 *= c2
		h1 ^= k1
	}

	// finalization
	h1 ^= length
	h1 ^= h1 >> 16
	h1 *= 0x85ebca6b
	h1 ^= h1 >> 13
	h1 *= 0xc2b2ae35
	h1 ^= h1 >> 16

	return h1
}

func rotateLeft64(val uint64, shift uint) uint64 {
	return (val << shift) | (val >> (64 - shift))
}

func mix64(v uint64) uint64 {
	v ^= v >> 33
	v *= 0xff51afd7ed558ccd
	v ^= v >> 33
	v *= 0xc4ceb9fe1a85ec53
	v ^= v >> 33
	return v
}

// murmurhash3_128 computes a 128-bit MurmurHash3 hash, returning (low, high).
func murmurhash3_128(data []byte, seedLow, seedHigh uint64) (uint64, uint64) {
	const (
		c1 uint64 = 0x87c37b91114253d5
		c2 uint64 = 0x4cf5ad432745937f
	)
	h1 := seedLow
	h2 := seedHigh

	// body - process 16-byte blocks
	pos := 0
	for pos < len(data)-15 {
		k1 := binary.LittleEndian.Uint64(data[pos : pos+8])
		k2 := binary.LittleEndian.Uint64(data[pos+8 : pos+16])

		k1 *= c1
		k1 = rotateLeft64(k1, 31)
		k1 *= c2
		h1 ^= k1
		h1 = rotateLeft64(h1, 27)
		h1 += h2
		h1 = h1*5 + 0x52dce729

		k2 *= c2
		k2 = rotateLeft64(k2, 33)
		k2 *= c1
		h2 ^= k2
		h2 = rotateLeft64(h2, 31)
		h2 += h1
		h2 = h2*5 + 0x38495ab5

		pos += 16
	}

	// tail
	k1 := uint64(0)
	k2 := uint64(0)
	n := len(data) & 15

	if n >= 15 {
		k2 ^= uint64(data[pos+14]) << 48
	}
	if n >= 14 {
		k2 ^= uint64(data[pos+13]) << 40
	}
	if n >= 13 {
		k2 ^= uint64(data[pos+12]) << 32
	}
	if n >= 12 {
		k2 ^= uint64(data[pos+11]) << 24
	}
	if n >= 11 {
		k2 ^= uint64(data[pos+10]) << 16
	}
	if n >= 10 {
		k2 ^= uint64(data[pos+9]) << 8
	}
	if n >= 9 {
		k2 ^= uint64(data[pos+8])
	}
	k2 *= c2
	k2 = rotateLeft64(k2, 33)
	k2 *= c1
	h2 ^= k2

	if n >= 8 {
		k1 ^= uint64(data[pos+7]) << 56
	}
	if n >= 7 {
		k1 ^= uint64(data[pos+6]) << 48
	}
	if n >= 6 {
		k1 ^= uint64(data[pos+5]) << 40
	}
	if n >= 5 {
		k1 ^= uint64(data[pos+4]) << 32
	}
	if n >= 4 {
		k1 ^= uint64(data[pos+3]) << 24
	}
	if n >= 3 {
		k1 ^= uint64(data[pos+2]) << 16
	}
	if n >= 2 {
		k1 ^= uint64(data[pos+1]) << 8
	}
	if n >= 1 {
		k1 ^= uint64(data[pos])
	}
	k1 *= c1
	k1 = rotateLeft64(k1, 31)
	k1 *= c2
	h1 ^= k1

	// finalization
	length := uint64(len(data))
	h1 ^= length
	h2 ^= length
	h1 += h2
	h2 += h1
	h1 = mix64(h1)
	h2 = mix64(h2)
	h1 += h2
	h2 += h1

	return h1, h2
}

// --- Hashing helpers ---

// writeForHashing writes a PK component for V1 hashing (string suffix = 0x00).
func writeForHashing(value interface{}, buf *[]byte) {
	writeForHashingCore(value, 0x00, buf)
}

// writeForHashingV2 writes a PK component for V2 hashing (string suffix = 0xFF).
func writeForHashingV2(value interface{}, buf *[]byte) {
	writeForHashingCore(value, 0xFF, buf)
}

func writeForHashingCore(value interface{}, stringSuffix byte, buf *[]byte) {
	switch v := value.(type) {
	case bool:
		if v {
			*buf = append(*buf, componentTrue)
		} else {
			*buf = append(*buf, componentFalse)
		}
	case nil:
		*buf = append(*buf, componentNull)
	case float64:
		*buf = append(*buf, componentNumber)
		var b [8]byte
		binary.LittleEndian.PutUint64(b[:], math.Float64bits(v))
		*buf = append(*buf, b[:]...)
	case string:
		*buf = append(*buf, componentString)
		*buf = append(*buf, []byte(v)...)
		*buf = append(*buf, stringSuffix)
	case UndefinedMarker:
		*buf = append(*buf, componentUndefined)
	}
}

// --- EPK computation ---

// ComputeV1 computes the V1 EPK.
// For each component: murmurhash3_32 of the full value,
// formatted as 12 zero bytes + 4-byte big-endian hash.
func ComputeV1(values []interface{}) string {
	var sb strings.Builder
	for _, v := range values {
		var hashBuf []byte
		writeForHashing(v, &hashBuf)
		hash := murmurhash3_32(hashBuf, 0)

		sb.WriteString("000000000000000000000000")
		fmt.Fprintf(&sb, "%08X", hash)
	}
	return sb.String()
}

// computeV2Hash computes the raw V2 EPK for Hash partitioning (without top-bit masking).
// Do not use for partition routing — use ComputeV2HashForRouting instead.
func computeV2Hash(values []interface{}) string {
	var hashBuf []byte
	for _, comp := range values {
		writeForHashingV2(comp, &hashBuf)
	}

	low, high := murmurhash3_128(hashBuf, 0, 0)
	return hash128ToEPK(low, high)
}

// computeV2MultiHash computes the raw V2 EPK for MultiHash partitioning (without top-bit masking).
// Do not use for partition routing — use ComputeV2MultiHashForRouting instead.
func computeV2MultiHash(values []interface{}) string {
	var sb strings.Builder
	for _, comp := range values {
		var hashBuf []byte
		writeForHashingV2(comp, &hashBuf)

		low, high := murmurhash3_128(hashBuf, 0, 0)
		sb.WriteString(hash128ToEPK(low, high))
	}
	return sb.String()
}

// ComputeV2HashForRouting computes the V2 EPK for routing (with top-2-bit masking).
// The service uses "FF" as the maximum exclusive partition key range sentinel,
// so all valid EPK values must have their first byte masked to [0x00, 0x3F].
func ComputeV2HashForRouting(values []interface{}) string {
	return maskTopBitsForRouting(computeV2Hash(values))
}

// ComputeV2MultiHashForRouting computes the V2 MultiHash EPK for routing.
// Each per-component hash has its top 2 bits masked independently.
func ComputeV2MultiHashForRouting(values []interface{}) string {
	var sb strings.Builder
	for _, comp := range values {
		var hashBuf []byte
		writeForHashingV2(comp, &hashBuf)

		low, high := murmurhash3_128(hashBuf, 0, 0)
		sb.WriteString(maskTopBitsForRouting(hash128ToEPK(low, high)))
	}
	return sb.String()
}

// hash128ToEPK converts a 128-bit hash (low, high) to an EPK hex string.
// The byte array is [low_LE, high_LE] reversed, producing big-endian order.
func hash128ToEPK(low, high uint64) string {
	var bytes [16]byte
	binary.LittleEndian.PutUint64(bytes[0:8], low)
	binary.LittleEndian.PutUint64(bytes[8:16], high)

	// Reverse to big-endian order
	for i, j := 0, 15; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}

	return toHexUpper(bytes[:])
}

// maskTopBitsForRouting clears the two most significant bits of the first byte
// in an EPK hex string. The Cosmos DB partition key range space uses "FF" as the
// maximum exclusive sentinel, so all valid EPK values must have their first byte
// in the range [0x00, 0x3F]. This matches the behavior of other Cosmos DB SDKs.
func maskTopBitsForRouting(epkHex string) string {
	if len(epkHex) < 2 {
		return epkHex
	}
	// Parse the first byte (2 hex chars)
	firstByte := hexCharToNibble(epkHex[0])<<4 | hexCharToNibble(epkHex[1])
	firstByte &= 0x3F
	// Replace the first two hex characters
	return fmt.Sprintf("%02X", firstByte) + epkHex[2:]
}

func hexCharToNibble(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	default:
		return 0
	}
}

// CompareEPK performs length-aware comparison of two effective partition key hex strings.
//
// For hierarchical partition key (HPK) containers, the service may return
// partition key ranges with mixed-length boundaries (e.g., 32-char partial
// vs 64-char fully specified, zero-padded). This function treats two EPKs
// as equal when one is a prefix of the other and the remainder is all '0' characters.
//
// Returns -1 if a < b, 0 if a == b, +1 if a > b.
//
// See: https://github.com/Azure/azure-cosmos-dotnet-v3/pull/5260
func CompareEPK(a, b string) int {
	commonLen := len(a)
	if len(b) < commonLen {
		commonLen = len(b)
	}

	// Compare the common prefix
	prefixA := a[:commonLen]
	prefixB := b[:commonLen]
	if prefixA < prefixB {
		return -1
	}
	if prefixA > prefixB {
		return 1
	}

	// Common prefixes are equal — check the tail of the longer string
	var tail string
	if len(a) > len(b) {
		tail = a[commonLen:]
	} else {
		tail = b[commonLen:]
	}

	// If the tail is all zeros (or empty), the EPKs are equal
	allZeros := true
	for i := 0; i < len(tail); i++ {
		if tail[i] != '0' {
			allZeros = false
			break
		}
	}
	if allZeros {
		return 0
	}

	// Non-zero tail: the longer string is greater
	if len(a) > len(b) {
		return 1
	}
	return -1
}

// toHexUpper returns uppercase hex encoding of data with no separators.
func toHexUpper(data []byte) string {
	var sb strings.Builder
	sb.Grow(len(data) * 2)
	for _, b := range data {
		fmt.Fprintf(&sb, "%02X", b)
	}
	return sb.String()
}
