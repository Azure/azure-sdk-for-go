// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"sort"
	"strings"
	"sync"
)

// ISessionContainer defines the interface for session token management.
type ISessionContainer interface {
	// GetSessionToken returns the session token for a collection link.
	GetSessionToken(collectionLink string) string

	// ResolveGlobalSessionToken resolves the global session token for a request.
	ResolveGlobalSessionToken(request *SessionContainerRequest) string

	// ResolvePartitionLocalSessionToken resolves the session token for a specific partition key range.
	ResolvePartitionLocalSessionToken(request *SessionContainerRequest, partitionKeyRangeID string) ISessionToken

	// SetSessionToken sets the session token from a request and response headers.
	SetSessionToken(request *SessionContainerRequest, responseHeaders map[string]string)

	// SetSessionTokenFromRID sets the session token using collection RID and name directly.
	SetSessionTokenFromRID(collectionRID string, collectionFullName string, responseHeaders map[string]string)

	// ClearTokenByCollectionFullName clears tokens for a collection by its full name.
	ClearTokenByCollectionFullName(collectionFullName string)

	// ClearTokenByResourceID clears tokens for a collection by its resource ID.
	ClearTokenByResourceID(resourceID string)
}

// SessionContainerRequest represents a request for session token operations.
// This is a simplified version that captures the essential fields needed for session management.
type SessionContainerRequest struct {
	// IsNameBased indicates whether the request uses name-based routing.
	IsNameBased bool

	// ResourceID is the resource ID string (e.g., collection RID).
	ResourceID string

	// ResourceAddress is the resource address/path.
	ResourceAddress string

	// ResourceType is the type of resource being accessed.
	ResourceType ResourceType

	// OperationType is the type of operation being performed.
	OperationType OperationType

	// RequestContext holds additional request context including resolved partition key range.
	RequestContext *SessionRequestContext
}

// SessionRequestContext holds additional context for session resolution.
type SessionRequestContext struct {
	// ResolvedPartitionKeyRange contains information about the resolved partition.
	ResolvedPartitionKeyRange *PartitionKeyRangeInfo
}

// PartitionKeyRangeInfo contains partition key range information.
type PartitionKeyRangeInfo struct {
	// ID is the partition key range ID.
	ID string

	// Parents contains the IDs of parent partition key ranges (for split handling).
	Parents []string
}

// HTTP header constants for session token management.
const (
	HTTPHeaderSessionToken  = "x-ms-session-token"
	HTTPHeaderOwnerFullName = "x-ms-alt-content-path"
	HTTPHeaderOwnerID       = "x-ms-content-path"
)

// SessionContainer caches session tokens for collections.
// It maps collection resource IDs to per-partition-key-range session tokens.
type SessionContainer struct {
	mu sync.RWMutex

	// hostName is the host name of the session container.
	hostName string

	// collectionResourceIDToSessionTokens maps collection resource ID to partition tokens.
	// map[collectionRID]map[partitionKeyRangeID]ISessionToken
	collectionResourceIDToSessionTokens map[string]map[string]ISessionToken

	// collectionNameToCollectionResourceID maps collection full name to collection resource ID.
	collectionNameToCollectionResourceID map[string]string

	// collectionResourceIDToCollectionName maps collection resource ID to collection full name.
	collectionResourceIDToCollectionName map[string]string
}

// NewSessionContainer creates a new SessionContainer.
func NewSessionContainer(hostName string) *SessionContainer {
	return &SessionContainer{
		hostName:                             hostName,
		collectionResourceIDToSessionTokens:  make(map[string]map[string]ISessionToken),
		collectionNameToCollectionResourceID: make(map[string]string),
		collectionResourceIDToCollectionName: make(map[string]string),
	}
}

// GetHostName returns the host name.
func (sc *SessionContainer) GetHostName() string {
	return sc.hostName
}

// GetSessionToken returns the combined session token for a collection link.
func (sc *SessionContainer) GetSessionToken(collectionLink string) string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	collectionName := getCollectionPath(collectionLink)
	collectionRID, exists := sc.collectionNameToCollectionResourceID[collectionName]
	if !exists {
		return ""
	}

	tokenMap, exists := sc.collectionResourceIDToSessionTokens[collectionRID]
	if !exists {
		return ""
	}

	return getCombinedSessionToken(tokenMap)
}

// ResolveGlobalSessionToken resolves the global session token for a request.
func (sc *SessionContainer) ResolveGlobalSessionToken(request *SessionContainerRequest) string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	tokenMap := sc.getPartitionKeyRangeIDToTokenMap(request)
	if tokenMap != nil {
		return getCombinedSessionToken(tokenMap)
	}

	return ""
}

// ResolvePartitionLocalSessionToken resolves the session token for a specific partition key range.
func (sc *SessionContainer) ResolvePartitionLocalSessionToken(request *SessionContainerRequest, partitionKeyRangeID string) ISessionToken {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	tokenMap := sc.getPartitionKeyRangeIDToTokenMap(request)
	return resolvePartitionLocalSessionTokenFromMap(request, partitionKeyRangeID, tokenMap)
}

// SetSessionToken sets the session token from a request and response headers.
func (sc *SessionContainer) SetSessionToken(request *SessionContainerRequest, responseHeaders map[string]string) {
	if request == nil {
		return
	}

	token := responseHeaders[HTTPHeaderSessionToken]
	if token == "" {
		return
	}

	// Determine collection name (priority: owner full name > resource address)
	ownerFullName := responseHeaders[HTTPHeaderOwnerFullName]
	if ownerFullName == "" {
		ownerFullName = request.ResourceAddress
	}
	collectionName := getCollectionPath(ownerFullName)

	// Determine resource ID
	var resourceIDString string
	if !request.IsNameBased {
		resourceIDString = request.ResourceID
	} else {
		resourceIDString = responseHeaders[HTTPHeaderOwnerID]
		if resourceIDString == "" {
			resourceIDString = request.ResourceID
		}
	}

	if resourceIDString == "" || collectionName == "" {
		return
	}

	// Don't set session tokens for master queries
	if isReadingFromMaster(request.ResourceType, request.OperationType) {
		return
	}

	sc.setSessionTokenInternal(resourceIDString, collectionName, token)
}

// SetSessionTokenFromRID sets the session token using collection RID and name directly.
func (sc *SessionContainer) SetSessionTokenFromRID(collectionRID string, collectionFullName string, responseHeaders map[string]string) {
	token := responseHeaders[HTTPHeaderSessionToken]
	if token == "" {
		return
	}

	collectionName := getCollectionPath(collectionFullName)
	sc.setSessionTokenInternal(collectionRID, collectionName, token)
}

// ClearTokenByCollectionFullName clears tokens for a collection by its full name.
func (sc *SessionContainer) ClearTokenByCollectionFullName(collectionFullName string) {
	if collectionFullName == "" {
		return
	}

	collectionName := getCollectionPath(collectionFullName)

	sc.mu.Lock()
	defer sc.mu.Unlock()

	rid, exists := sc.collectionNameToCollectionResourceID[collectionName]
	if !exists {
		return
	}

	delete(sc.collectionResourceIDToSessionTokens, rid)
	delete(sc.collectionResourceIDToCollectionName, rid)
	delete(sc.collectionNameToCollectionResourceID, collectionName)
}

// ClearTokenByResourceID clears tokens for a collection by its resource ID.
func (sc *SessionContainer) ClearTokenByResourceID(resourceID string) {
	if resourceID == "" {
		return
	}

	sc.mu.Lock()
	defer sc.mu.Unlock()

	collectionName, exists := sc.collectionResourceIDToCollectionName[resourceID]
	if !exists {
		return
	}

	delete(sc.collectionResourceIDToSessionTokens, resourceID)
	delete(sc.collectionResourceIDToCollectionName, resourceID)
	delete(sc.collectionNameToCollectionResourceID, collectionName)
}

// setSessionTokenInternal sets the session token internally.
func (sc *SessionContainer) setSessionTokenInternal(collectionRID, collectionName, token string) {
	// Parse the token: format is "partitionKeyRangeId:sessionToken"
	parts := strings.SplitN(token, PartitionKeyRangeSessionSeparator, 2)
	if len(parts) != 2 {
		return
	}

	partitionKeyRangeID := parts[0]
	parsedToken, ok := TryCreateVectorSessionToken(parts[1])
	if !ok {
		return
	}

	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Check if this is a known collection
	existingRID, nameExists := sc.collectionNameToCollectionResourceID[collectionName]
	existingName, ridExists := sc.collectionResourceIDToCollectionName[collectionRID]
	isKnownCollection := nameExists && ridExists &&
		existingRID == collectionRID &&
		existingName == collectionName

	if !isKnownCollection {
		// Register the collection mapping
		sc.collectionNameToCollectionResourceID[collectionName] = collectionRID
		sc.collectionResourceIDToCollectionName[collectionRID] = collectionName
	}

	// Add or merge the session token
	sc.addSessionToken(collectionRID, partitionKeyRangeID, parsedToken)
}

// addSessionToken adds or merges a session token for a partition.
func (sc *SessionContainer) addSessionToken(collectionRID, partitionKeyRangeID string, newToken ISessionToken) {
	tokenMap, exists := sc.collectionResourceIDToSessionTokens[collectionRID]
	if !exists {
		tokenMap = make(map[string]ISessionToken)
		sc.collectionResourceIDToSessionTokens[collectionRID] = tokenMap
	}

	existingToken, exists := tokenMap[partitionKeyRangeID]
	if !exists {
		tokenMap[partitionKeyRangeID] = newToken
		return
	}

	// Merge with existing token
	mergedToken, err := existingToken.Merge(newToken)
	if err != nil {
		// On merge error, keep the existing token
		return
	}
	tokenMap[partitionKeyRangeID] = mergedToken
}

// getPartitionKeyRangeIDToTokenMap returns the token map for a request.
func (sc *SessionContainer) getPartitionKeyRangeIDToTokenMap(request *SessionContainerRequest) map[string]ISessionToken {
	if !request.IsNameBased {
		if request.ResourceID != "" {
			return sc.collectionResourceIDToSessionTokens[request.ResourceID]
		}
	} else {
		collectionName := getCollectionName(request.ResourceAddress)
		if collectionName != "" {
			if rid, exists := sc.collectionNameToCollectionResourceID[collectionName]; exists {
				return sc.collectionResourceIDToSessionTokens[rid]
			}
		}
	}
	return nil
}

// resolvePartitionLocalSessionTokenFromMap resolves the local session token from a token map.
func resolvePartitionLocalSessionTokenFromMap(request *SessionContainerRequest, partitionKeyRangeID string, tokenMap map[string]ISessionToken) ISessionToken {
	if tokenMap == nil {
		return nil
	}

	// Direct match
	if token, exists := tokenMap[partitionKeyRangeID]; exists {
		return token
	}

	// No direct match — check parent partition key ranges.
	// A partition can have more than one parent (merge). In that case, we merge
	// all matching parent tokens to produce a token with the max LSNs from each.
	if request != nil && request.RequestContext != nil &&
		request.RequestContext.ResolvedPartitionKeyRange != nil &&
		len(request.RequestContext.ResolvedPartitionKeyRange.Parents) > 0 {

		parents := request.RequestContext.ResolvedPartitionKeyRange.Parents
		var parentSessionToken ISessionToken

		// Traverse parents from most recent to oldest
		for i := len(parents) - 1; i >= 0; i-- {
			if token, exists := tokenMap[parents[i]]; exists {
				if parentSessionToken == nil {
					parentSessionToken = token
				} else {
					merged, err := parentSessionToken.Merge(token)
					if err != nil {
						// On merge error, return what we have so far
						return parentSessionToken
					}
					parentSessionToken = merged
				}
			}
		}

		return parentSessionToken
	}

	return nil
}

// getCombinedSessionToken combines all tokens in the map into a comma-separated string.
func getCombinedSessionToken(tokens map[string]ISessionToken) string {
	if len(tokens) == 0 {
		return ""
	}

	// Sort partition key range IDs for deterministic output
	keys := make([]string, 0, len(tokens))
	for k := range tokens {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for i, key := range keys {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(key)
		sb.WriteString(PartitionKeyRangeSessionSeparator)
		sb.WriteString(tokens[key].ConvertToString())
	}

	return sb.String()
}

func getCollectionPath(fullPath string) string {
	if fullPath == "" {
		return ""
	}

	path := strings.Trim(fullPath, "/")
	segments := strings.Split(path, "/")

	if len(segments) < 4 {
		return path
	}

	if strings.ToLower(segments[0]) == "dbs" {
		return strings.Join(segments[:4], "/")
	}

	return path
}

// getCollectionName extracts the collection name (same as getCollectionPath for now).
func getCollectionName(resourceAddress string) string {
	return getCollectionPath(resourceAddress)
}

// isReadingFromMaster determines if the operation is a master query.
// Master queries don't use session tokens.
func isReadingFromMaster(resourceType ResourceType, operationType OperationType) bool {
	// Document collection operations: only feed/query variants are master reads.
	// Single-document reads (Read, Head, HeadFeed) go through the partition.
	if resourceType == ResourceCollection {
		return operationType == OperationReadFeed ||
			operationType == OperationQuery ||
			operationType == OperationSQLQuery
	}

	// PartitionKeyRange is master except for GetSplitPoint and AbortSplit,
	// which are partition-level operations.
	if resourceType == ResourcePartitionKeyRange {
		return operationType != OperationGetSplitPoint &&
			operationType != OperationAbortSplit
	}

	// Other master resources
	switch resourceType {
	case ResourceDatabase, ResourceUser, ResourcePermission, ResourceOffer,
		ResourceDatabaseAccount, ResourceTopology, ResourceUserDefinedType:
		return true
	}

	return false
}

// SessionTokenHelper provides helper functions for session token operations.
type SessionTokenHelper struct{}

// Parse parses a session token string.
func (h *SessionTokenHelper) Parse(sessionToken string) (ISessionToken, error) {
	// Handle tokens that may include partition key range ID prefix
	parts := strings.Split(sessionToken, PartitionKeyRangeSessionSeparator)
	tokenPart := parts[len(parts)-1]

	token, ok := TryCreateVectorSessionToken(tokenPart)
	if !ok {
		return nil, NewSessionTokenError("invalid session token: " + sessionToken)
	}
	return token, nil
}

// TryParse attempts to parse a session token string.
func (h *SessionTokenHelper) TryParse(sessionToken string) (ISessionToken, bool) {
	if sessionToken == "" {
		return nil, false
	}

	parts := strings.Split(sessionToken, PartitionKeyRangeSessionSeparator)
	tokenPart := parts[len(parts)-1]

	token, ok := TryCreateVectorSessionToken(tokenPart)
	return token, ok
}

// ResolvePartitionLocalSessionToken resolves a session token from a global session token string.
// If a direct match for the partitionKeyRangeID exists, it is returned immediately.
// Otherwise, parent partition key ranges are checked and merged if found.
func (h *SessionTokenHelper) ResolvePartitionLocalSessionToken(
	request *SessionContainerRequest,
	partitionKeyRangeID string,
	globalSessionToken string,
) (ISessionToken, error) {
	if partitionKeyRangeID == "" || globalSessionToken == "" {
		return nil, nil
	}

	// Parse the global session token (comma-separated list of pkRangeId:token pairs)
	// into a map for easy lookup
	localTokens := strings.Split(globalSessionToken, ",")
	rangeIDToTokenMap := make(map[string]ISessionToken)

	for _, localToken := range localTokens {
		parts := strings.SplitN(localToken, PartitionKeyRangeSessionSeparator, 2)
		if len(parts) != 2 {
			return nil, NewSessionTokenError("invalid session token format: " + localToken)
		}

		rangeID := parts[0]
		tokenString := parts[1]

		parsedToken, ok := TryCreateVectorSessionToken(tokenString)
		if !ok {
			return nil, NewSessionTokenError("invalid session token: " + tokenString)
		}

		rangeIDToTokenMap[rangeID] = parsedToken
	}

	// Check for direct match first - if found, return immediately without merging parents
	if token, exists := rangeIDToTokenMap[partitionKeyRangeID]; exists {
		return token, nil
	}

	// No direct match - check parent partition key ranges
	// A partition can have more than 1 parent (merge). In that case, we apply Merge
	// to generate a token with both parent's max LSNs
	if request != nil && request.RequestContext != nil &&
		request.RequestContext.ResolvedPartitionKeyRange != nil &&
		len(request.RequestContext.ResolvedPartitionKeyRange.Parents) > 0 {

		parents := request.RequestContext.ResolvedPartitionKeyRange.Parents
		var parentSessionToken ISessionToken

		// Traverse parents from most recent to oldest
		for i := len(parents) - 1; i >= 0; i-- {
			parentID := parents[i]
			if token, exists := rangeIDToTokenMap[parentID]; exists {
				if parentSessionToken == nil {
					parentSessionToken = token
				} else {
					merged, err := parentSessionToken.Merge(token)
					if err != nil {
						return nil, err
					}
					parentSessionToken = merged
				}
			}
		}

		return parentSessionToken, nil
	}

	return nil, nil
}

// DefaultSessionTokenHelper is the default instance of SessionTokenHelper.
var DefaultSessionTokenHelper = &SessionTokenHelper{}
