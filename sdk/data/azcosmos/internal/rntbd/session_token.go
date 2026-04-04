// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// ISessionToken defines the interface for session tokens.
// Session tokens are immutable after construction.
type ISessionToken interface {
	// IsValid returns true if this session token is valid with respect to the other session token.
	// This is used to decide if the client can accept the server's response based on comparison
	// between client's and server's session tokens.
	IsValid(other ISessionToken) (bool, error)

	// Merge returns a new instance of session token obtained by merging this session token with
	// the given session token. Merge is a commutative operation: a.Merge(b).Equals(b.Merge(a))
	Merge(other ISessionToken) (ISessionToken, error)

	// GetLSN returns the global LSN of this session token.
	GetLSN() int64

	// ConvertToString returns the string representation of this session token.
	ConvertToString() string
}

// PartitionKeyRangeSessionSeparator is the separator between partition key range ID and session token.
const PartitionKeyRangeSessionSeparator = ":"

// VectorSessionToken models a vector clock based session token.
// Session token format: {Version}#{GlobalLSN}#{RegionId1}={LocalLsn1}#{RegionId2}={LocalLsn2}...#{RegionIdN}={LocalLsnN}
//
// 'Version' captures the configuration number of the partition which returned this session token.
// 'Version' is incremented every time the topology of the partition is updated (e.g., Add/Remove/Failover).
//
// The choice of separators '#' and '=' is important. Separators ';' and ',' are used to delimit
// per-partitionKeyRange session tokens.
//
// Instances of this struct are immutable (read only after construction).
type VectorSessionToken struct {
	version          int64
	globalLSN        int64
	localLSNByRegion map[int]int64
	sessionToken     string
}

const (
	segmentSeparator        = '#'
	regionProgressSeparator = '='
)

// TryCreateVectorSessionToken attempts to parse a session token string and returns the parsed token.
// Returns (token, true) on success, (nil, false) on failure.
func TryCreateVectorSessionToken(sessionToken string) (*VectorSessionToken, bool) {
	version, globalLSN, localLSNByRegion, ok := tryParseSessionToken(sessionToken)
	if !ok {
		return nil, false
	}

	return &VectorSessionToken{
		version:          version,
		globalLSN:        globalLSN,
		localLSNByRegion: localLSNByRegion,
		sessionToken:     sessionToken,
	}, true
}

// NewVectorSessionToken creates a new VectorSessionToken with the given values.
// The session token string is computed from the values.
func NewVectorSessionToken(version, globalLSN int64, localLSNByRegion map[int]int64) *VectorSessionToken {
	token := &VectorSessionToken{
		version:          version,
		globalLSN:        globalLSN,
		localLSNByRegion: localLSNByRegion,
	}
	token.sessionToken = token.buildSessionToken()
	return token
}

// GetLSN returns the global LSN.
func (v *VectorSessionToken) GetLSN() int64 {
	return v.globalLSN
}

// GetVersion returns the version number.
func (v *VectorSessionToken) GetVersion() int64 {
	return v.version
}

// GetLocalLSNByRegion returns a copy of the local LSN by region map.
func (v *VectorSessionToken) GetLocalLSNByRegion() map[int]int64 {
	result := make(map[int]int64, len(v.localLSNByRegion))
	for k, val := range v.localLSNByRegion {
		result[k] = val
	}
	return result
}

// Equals checks if two VectorSessionTokens are equal.
func (v *VectorSessionToken) Equals(other *VectorSessionToken) bool {
	if other == nil {
		return false
	}

	if v.version != other.version || v.globalLSN != other.globalLSN {
		return false
	}

	return v.areRegionProgressEqual(other.localLSNByRegion)
}

// IsValid returns true if this session token is valid with respect to the other session token.
func (v *VectorSessionToken) IsValid(other ISessionToken) (bool, error) {
	otherVector, ok := other.(*VectorSessionToken)
	if !ok {
		return false, fmt.Errorf("otherSessionToken must be a VectorSessionToken")
	}

	if otherVector.version < v.version || otherVector.globalLSN < v.globalLSN {
		return false, nil
	}

	if otherVector.version == v.version && len(otherVector.localLSNByRegion) != len(v.localLSNByRegion) {
		return false, NewSessionTokenError(
			fmt.Sprintf("Invalid regions in session token. Expected regions in '%s' but got '%s'",
				v.sessionToken, otherVector.sessionToken))
	}

	for regionID, otherLocalLSN := range otherVector.localLSNByRegion {
		localLSN, exists := v.localLSNByRegion[regionID]
		if !exists {
			// Region mismatch: other session token has progress for a region which is missing in this session token
			// Region mismatch can be ignored only if this session token version is smaller than other session token version
			if v.version == otherVector.version {
				return false, NewSessionTokenError(
					fmt.Sprintf("Invalid regions in session token. Expected regions in '%s' but got '%s'",
						v.sessionToken, otherVector.sessionToken))
			}
			// ignore missing region as other session token version > this session token version
		} else {
			// region is present in both session tokens
			if otherLocalLSN < localLSN {
				return false, nil
			}
		}
	}

	return true, nil
}

// Merge returns a new session token obtained by merging this token with the other token.
// Merge is commutative: a.Merge(b).Equals(b.Merge(a))
func (v *VectorSessionToken) Merge(other ISessionToken) (ISessionToken, error) {
	otherVector, ok := other.(*VectorSessionToken)
	if !ok {
		return nil, fmt.Errorf("obj must be a VectorSessionToken")
	}

	if v.version == otherVector.version && len(v.localLSNByRegion) != len(otherVector.localLSNByRegion) {
		return nil, NewSessionTokenError(
			fmt.Sprintf("Invalid regions in session token. Expected regions in '%s' but got '%s'",
				v.sessionToken, otherVector.sessionToken))
	}

	var sessionTokenWithHigherVersion, sessionTokenWithLowerVersion *VectorSessionToken
	if v.version < otherVector.version {
		sessionTokenWithLowerVersion = v
		sessionTokenWithHigherVersion = otherVector
	} else {
		sessionTokenWithLowerVersion = otherVector
		sessionTokenWithHigherVersion = v
	}

	highestLocalLSNByRegion := make(map[int]int64)

	for regionID, localLSN1 := range sessionTokenWithHigherVersion.localLSNByRegion {
		if localLSN2, exists := sessionTokenWithLowerVersion.localLSNByRegion[regionID]; exists {
			highestLocalLSNByRegion[regionID] = max(localLSN1, localLSN2)
		} else if v.version == otherVector.version {
			return nil, NewSessionTokenError(
				fmt.Sprintf("Invalid regions in session token. Expected regions in '%s' but got '%s'",
					v.sessionToken, otherVector.sessionToken))
		} else {
			highestLocalLSNByRegion[regionID] = localLSN1
		}
	}

	return NewVectorSessionToken(
		max(v.version, otherVector.version),
		max(v.globalLSN, otherVector.globalLSN),
		highestLocalLSNByRegion,
	), nil
}

// ConvertToString returns the string representation of this session token.
func (v *VectorSessionToken) ConvertToString() string {
	return v.sessionToken
}

// buildSessionToken constructs the session token string from the values.
func (v *VectorSessionToken) buildSessionToken() string {
	if len(v.localLSNByRegion) == 0 {
		return fmt.Sprintf("%d%c%d", v.version, segmentSeparator, v.globalLSN)
	}

	// Sort region IDs for deterministic output
	regionIDs := make([]int, 0, len(v.localLSNByRegion))
	for regionID := range v.localLSNByRegion {
		regionIDs = append(regionIDs, regionID)
	}
	sort.Ints(regionIDs)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d%c%d", v.version, segmentSeparator, v.globalLSN))

	for _, regionID := range regionIDs {
		sb.WriteByte(segmentSeparator)
		sb.WriteString(fmt.Sprintf("%d%c%d", regionID, regionProgressSeparator, v.localLSNByRegion[regionID]))
	}

	return sb.String()
}

// areRegionProgressEqual checks if the region progress maps are equal.
func (v *VectorSessionToken) areRegionProgressEqual(other map[int]int64) bool {
	if len(v.localLSNByRegion) != len(other) {
		return false
	}

	for regionID, localLSN1 := range v.localLSNByRegion {
		localLSN2, exists := other[regionID]
		if !exists || localLSN1 != localLSN2 {
			return false
		}
	}

	return true
}

// tryParseSessionToken attempts to parse a session token string.
func tryParseSessionToken(sessionToken string) (version int64, globalLSN int64, localLSNByRegion map[int]int64, ok bool) {
	if sessionToken == "" {
		return 0, 0, nil, false
	}

	segments := strings.Split(sessionToken, string(segmentSeparator))
	if len(segments) < 2 {
		return 0, 0, nil, false
	}

	var err error
	version, err = strconv.ParseInt(segments[0], 10, 64)
	if err != nil {
		return 0, 0, nil, false
	}

	globalLSN, err = strconv.ParseInt(segments[1], 10, 64)
	if err != nil {
		return 0, 0, nil, false
	}

	localLSNByRegion = make(map[int]int64)

	for i := 2; i < len(segments); i++ {
		regionSegment := segments[i]
		regionIDWithLSN := strings.Split(regionSegment, string(regionProgressSeparator))

		if len(regionIDWithLSN) != 2 {
			return 0, 0, nil, false
		}

		regionID, err := strconv.Atoi(regionIDWithLSN[0])
		if err != nil {
			return 0, 0, nil, false
		}

		localLSN, err := strconv.ParseInt(regionIDWithLSN[1], 10, 64)
		if err != nil {
			return 0, 0, nil, false
		}

		localLSNByRegion[regionID] = localLSN
	}

	return version, globalLSN, localLSNByRegion, true
}

// SessionTokenError represents an error during session token operations.
// This is used for invalid region mismatches and other session token validation failures.
type SessionTokenError struct {
	Message string
}

func (e *SessionTokenError) Error() string {
	return e.Message
}

// NewSessionTokenError creates a new SessionTokenError.
func NewSessionTokenError(message string) *SessionTokenError {
	return &SessionTokenError{Message: message}
}

// IsSessionTokenError checks if the error is a SessionTokenError.
func IsSessionTokenError(err error) bool {
	_, ok := err.(*SessionTokenError)
	return ok
}

// ParseSessionToken parses a session token string and returns the ISessionToken.
// Returns an error if parsing fails.
func ParseSessionToken(sessionToken string) (ISessionToken, error) {
	token, ok := TryCreateVectorSessionToken(sessionToken)
	if !ok {
		return nil, fmt.Errorf("invalid session token: %s", sessionToken)
	}
	return token, nil
}
