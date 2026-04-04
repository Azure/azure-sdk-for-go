// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Test: validateSuccessfulSessionTokenParsing
func TestVectorSessionToken_SuccessfulParsing(t *testing.T) {
	sessionToken := "1#100#1=20#2=5#3=30"
	token, ok := TryCreateVectorSessionToken(sessionToken)

	require.True(t, ok, "Expected successful parsing")
	require.NotNil(t, token)
	require.Equal(t, int64(1), token.GetVersion())
	require.Equal(t, int64(100), token.GetLSN())
	require.Equal(t, sessionToken, token.ConvertToString())
}

// Test: validateSessionTokenParsingWithInvalidVersion
func TestVectorSessionToken_InvalidVersion(t *testing.T) {
	sessionToken := "foo#100#1=20#2=5#3=30"
	token, ok := TryCreateVectorSessionToken(sessionToken)

	require.False(t, ok, "Expected parsing to fail for invalid version")
	require.Nil(t, token)
}

// Test: validateSessionTokenParsingWithInvalidGlobalLsn
func TestVectorSessionToken_InvalidGlobalLSN(t *testing.T) {
	sessionToken := "1#foo#1=20#2=5#3=30"
	token, ok := TryCreateVectorSessionToken(sessionToken)

	require.False(t, ok, "Expected parsing to fail for invalid global LSN")
	require.Nil(t, token)
}

// Test: validateSessionTokenParsingWithInvalidRegionProgress
func TestVectorSessionToken_InvalidRegionProgress(t *testing.T) {
	sessionToken := "1#100#1=20#2=x#3=30"
	token, ok := TryCreateVectorSessionToken(sessionToken)

	require.False(t, ok, "Expected parsing to fail for invalid region progress")
	require.Nil(t, token)
}

// Test: validateSessionTokenParsingWithInvalidFormat
func TestVectorSessionToken_InvalidFormat(t *testing.T) {
	sessionToken := "1;100#1=20#2=40"
	token, ok := TryCreateVectorSessionToken(sessionToken)

	require.False(t, ok, "Expected parsing to fail for invalid format (wrong separator)")
	require.Nil(t, token)
}

// Test: validateSessionTokenParsingFromEmptyString
func TestVectorSessionToken_EmptyString(t *testing.T) {
	sessionToken := ""
	token, ok := TryCreateVectorSessionToken(sessionToken)

	require.False(t, ok, "Expected parsing to fail for empty string")
	require.Nil(t, token)
}

// Test: validateSessionTokenComparison (multiple sub-tests)
func TestVectorSessionToken_Comparison(t *testing.T) {
	t.Run("DifferentVersionsAndRegions", func(t *testing.T) {
		token1, ok1 := TryCreateVectorSessionToken("1#100#1=20#2=5#3=30")
		token2, ok2 := TryCreateVectorSessionToken("2#105#4=10#2=5#3=30")

		require.True(t, ok1)
		require.True(t, ok2)

		// Different tokens are not equal
		require.False(t, token1.Equals(token2))
		require.False(t, token2.Equals(token1))

		// token1 is valid with respect to token2 (token2 has higher version and LSN)
		valid, err := token1.IsValid(token2)
		require.NoError(t, err)
		require.True(t, valid)

		// token2 is NOT valid with respect to token1 (lower version cannot satisfy higher)
		valid, err = token2.IsValid(token1)
		require.NoError(t, err)
		require.False(t, valid)

		// Merge should produce token with higher version/LSN and combined regions
		merged, err := token1.Merge(token2)
		require.NoError(t, err)

		// Expected merge result: version 2, globalLSN 105, regions from higher version (4, 2, 3)
		expectedMerged, ok := TryCreateVectorSessionToken("2#105#2=5#3=30#4=10")
		require.True(t, ok)
		require.True(t, expectedMerged.Equals(merged.(*VectorSessionToken)))
	})

	t.Run("SameVersionDifferentProgress", func(t *testing.T) {
		token1, ok1 := TryCreateVectorSessionToken("1#100#1=20#2=5#3=30")
		token2, ok2 := TryCreateVectorSessionToken("1#100#1=10#2=8#3=30")

		require.True(t, ok1)
		require.True(t, ok2)

		// Not equal (different local LSNs)
		require.False(t, token1.Equals(token2))
		require.False(t, token2.Equals(token1))

		// Neither is valid with respect to the other (same version, some regions have lower LSN)
		valid, err := token1.IsValid(token2)
		require.NoError(t, err)
		require.False(t, valid, "token1 should not be valid wrt token2 (region 1: 20 > 10)")

		valid, err = token2.IsValid(token1)
		require.NoError(t, err)
		require.False(t, valid, "token2 should not be valid wrt token1 (region 2: 8 > 5)")

		// Merge should take max of each region
		merged, err := token1.Merge(token2)
		require.NoError(t, err)

		expectedMerged, ok := TryCreateVectorSessionToken("1#100#1=20#2=8#3=30")
		require.True(t, ok)
		require.True(t, expectedMerged.Equals(merged.(*VectorSessionToken)))
	})

	t.Run("SameVersionHigherGlobalLSN", func(t *testing.T) {
		token1, ok1 := TryCreateVectorSessionToken("1#100#1=20#2=5#3=30")
		token2, ok2 := TryCreateVectorSessionToken("1#102#1=100#2=8#3=30")

		require.True(t, ok1)
		require.True(t, ok2)

		// Not equal
		require.False(t, token1.Equals(token2))
		require.False(t, token2.Equals(token1))

		// token1 is valid wrt token2 (token2 has higher globalLSN and all local LSNs >= token1)
		valid, err := token1.IsValid(token2)
		require.NoError(t, err)
		require.True(t, valid)

		// token2 is NOT valid wrt token1 (token1 has lower globalLSN)
		valid, err = token2.IsValid(token1)
		require.NoError(t, err)
		require.False(t, valid)

		// Merge
		merged, err := token1.Merge(token2)
		require.NoError(t, err)

		expectedMerged, ok := TryCreateVectorSessionToken("1#102#1=100#2=8#3=30")
		require.True(t, ok)
		require.True(t, expectedMerged.Equals(merged.(*VectorSessionToken)))
	})

	t.Run("SameVersionDifferentRegionCount_MergeFails", func(t *testing.T) {
		token1, ok1 := TryCreateVectorSessionToken("1#101#1=20#2=5#3=30")
		token2, ok2 := TryCreateVectorSessionToken("1#100#1=20#2=5#3=30#4=40")

		require.True(t, ok1)
		require.True(t, ok2)

		// Merge should fail when same version but different region count
		_, err := token1.Merge(token2)
		require.Error(t, err)
		require.True(t, IsSessionTokenError(err), "Expected SessionTokenError")
	})

	t.Run("SameVersionDifferentRegionCount_IsValidFails", func(t *testing.T) {
		token1, ok1 := TryCreateVectorSessionToken("1#101#1=20#2=5#3=30")
		token2, ok2 := TryCreateVectorSessionToken("1#100#1=20#2=5#3=30#4=40")

		require.True(t, ok1)
		require.True(t, ok2)

		// IsValid should fail when same version but different region count
		_, err := token2.IsValid(token1)
		require.Error(t, err)
		require.True(t, IsSessionTokenError(err), "Expected SessionTokenError")
	})
}

// Test: session token without region progress (only version and globalLSN)
func TestVectorSessionToken_NoRegionProgress(t *testing.T) {
	sessionToken := "1#100"
	token, ok := TryCreateVectorSessionToken(sessionToken)

	require.True(t, ok, "Expected successful parsing of token without region progress")
	require.NotNil(t, token)
	require.Equal(t, int64(1), token.GetVersion())
	require.Equal(t, int64(100), token.GetLSN())
	require.Equal(t, 0, len(token.GetLocalLSNByRegion()))
}

// Test: NewVectorSessionToken constructs token correctly
func TestNewVectorSessionToken(t *testing.T) {
	localLSNByRegion := map[int]int64{
		1: 20,
		2: 5,
		3: 30,
	}

	token := NewVectorSessionToken(1, 100, localLSNByRegion)

	require.Equal(t, int64(1), token.GetVersion())
	require.Equal(t, int64(100), token.GetLSN())
	require.Equal(t, localLSNByRegion, token.GetLocalLSNByRegion())

	// Verify the string representation is correct (sorted by region ID)
	tokenStr := token.ConvertToString()
	require.Equal(t, "1#100#1=20#2=5#3=30", tokenStr)
}

// Test: ParseSessionToken utility function
func TestParseSessionToken(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		token, err := ParseSessionToken("1#100#1=20#2=5#3=30")
		require.NoError(t, err)
		require.NotNil(t, token)
		require.Equal(t, int64(100), token.GetLSN())
	})

	t.Run("Invalid", func(t *testing.T) {
		token, err := ParseSessionToken("invalid")
		require.Error(t, err)
		require.Nil(t, token)
	})
}

// Test: Merge is commutative (a.Merge(b) == b.Merge(a))
func TestVectorSessionToken_MergeIsCommutative(t *testing.T) {
	token1, _ := TryCreateVectorSessionToken("1#100#1=20#2=5#3=30")
	token2, _ := TryCreateVectorSessionToken("1#102#1=100#2=8#3=30")

	merged1, err1 := token1.Merge(token2)
	merged2, err2 := token2.Merge(token1)

	require.NoError(t, err1)
	require.NoError(t, err2)

	// Merge results should be equal (commutativity)
	require.True(t, merged1.(*VectorSessionToken).Equals(merged2.(*VectorSessionToken)))
}

// Test: GetLocalLSNByRegion returns a copy
func TestVectorSessionToken_GetLocalLSNByRegionReturnsCopy(t *testing.T) {
	token, _ := TryCreateVectorSessionToken("1#100#1=20#2=5#3=30")

	regions := token.GetLocalLSNByRegion()
	regions[1] = 999 // Modify the returned map

	// Original token should not be affected
	originalRegions := token.GetLocalLSNByRegion()
	require.Equal(t, int64(20), originalRegions[1])
}

// Additional edge cases
func TestVectorSessionToken_EdgeCases(t *testing.T) {
	t.Run("SingleSegment", func(t *testing.T) {
		// Only one segment should fail
		token, ok := TryCreateVectorSessionToken("100")
		require.False(t, ok)
		require.Nil(t, token)
	})

	t.Run("InvalidRegionIDFormat", func(t *testing.T) {
		// Region ID that's not an integer
		token, ok := TryCreateVectorSessionToken("1#100#abc=20")
		require.False(t, ok)
		require.Nil(t, token)
	})

	t.Run("MissingRegionLSN", func(t *testing.T) {
		// Missing LSN value for region
		token, ok := TryCreateVectorSessionToken("1#100#1=")
		require.False(t, ok)
		require.Nil(t, token)
	})

	t.Run("MissingEqualsInRegion", func(t *testing.T) {
		// Missing equals sign in region segment
		token, ok := TryCreateVectorSessionToken("1#100#120")
		require.False(t, ok)
		require.Nil(t, token)
	})

	t.Run("NegativeVersion", func(t *testing.T) {
		// Negative version is actually valid (parsed as int64)
		token, ok := TryCreateVectorSessionToken("-1#100#1=20")
		require.True(t, ok)
		require.Equal(t, int64(-1), token.GetVersion())
	})

	t.Run("LargeNumbers", func(t *testing.T) {
		// Very large numbers
		token, ok := TryCreateVectorSessionToken("9223372036854775807#9223372036854775807#1=9223372036854775807")
		require.True(t, ok)
		require.Equal(t, int64(9223372036854775807), token.GetVersion())
		require.Equal(t, int64(9223372036854775807), token.GetLSN())
	})
}

// Test: IsValid with region mismatch when versions differ (should succeed, not error)
func TestVectorSessionToken_IsValid_RegionMismatchDifferentVersions(t *testing.T) {
	// token1 has regions 1,2,3 at version 1
	// token2 has regions 2,3,4 at version 2 (different regions, but version is higher)
	token1, _ := TryCreateVectorSessionToken("1#100#1=20#2=5#3=30")
	token2, _ := TryCreateVectorSessionToken("2#105#2=10#3=35#4=15")

	// token1.IsValid(token2) should succeed because token2 has higher version
	// Missing regions in token1 are ignored when other has higher version
	valid, err := token1.IsValid(token2)
	require.NoError(t, err)
	require.True(t, valid)
}
