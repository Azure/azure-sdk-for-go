// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"strconv"
	"strings"
	"sync"
)

type partitionRangeIdToSessionTokens map[int64]string

type sessionContainer struct {
	// containerPathToRid maps named container paths to their resource IDs
	containerPathToRid map[string]string

	// ridToSessionTokens maps container resource IDs to their per-partition session tokens
	ridToSessionTokens map[string]partitionRangeIdToSessionTokens

	mutex sync.RWMutex
}

func newSessionContainer() *sessionContainer {
	return &sessionContainer{
		containerPathToRid: make(map[string]string),
		ridToSessionTokens: make(map[string]partitionRangeIdToSessionTokens),
		mutex:              sync.RWMutex{},
	}
}

func (sc *sessionContainer) GetSessionToken(resourceAddress string) string {
	sc.mutex.RLock()
	defer sc.mutex.RUnlock()

	// Get the container path from the resource address e.g. dbs/testdb/colls/testcoll
	containerPath, err := getContainerPath(resourceAddress)
	if err != nil {
		return ""
	}

	// Resolve container resource ID
	rid, ok := sc.containerPathToRid[containerPath]
	if !ok {
		return ""
	}

	// Get the session tokens for the container
	tokensByPartition, ok := sc.ridToSessionTokens[rid]
	if !ok {
		return ""
	}

	// We don't currently support mapping partition keys to partition range IDs, so
	// return all the session tokens for the container.
	combinedTokens := strings.Builder{}
	for partitionRangeId, token := range tokensByPartition {
		if combinedTokens.Len() > 0 {
			combinedTokens.WriteString(",")
		}
		combinedTokens.WriteString(strconv.FormatInt(partitionRangeId, 10))
		combinedTokens.WriteString(":")
		combinedTokens.WriteString(token)
	}

	return combinedTokens.String()
}

func (sc *sessionContainer) SetSessionToken(resourceAddress string, containerRid string, sessionToken string) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	if resourceAddress == "" || containerRid == "" || sessionToken == "" {
		return
	}

	// Get the container path (possibly no-op if already container path)
	containerPath, err := getContainerPath(resourceAddress)
	if err != nil {
		return
	}

	// Check if the RID changed (e.g. recreated container) and clear any old tokens
	existingRid, ok := sc.containerPathToRid[containerPath]
	if ok && existingRid != containerRid {
		delete(sc.ridToSessionTokens, existingRid)
	}
	sc.containerPathToRid[containerPath] = containerRid

	// Create a new map for the container if needed
	if _, ok := sc.ridToSessionTokens[containerRid]; !ok {
		sc.ridToSessionTokens[containerRid] = make(partitionRangeIdToSessionTokens)
	}

	// Store any returned per-partition session tokens
	for _, sessionToken := range strings.Split(sessionToken, ",") {
		tokenParts := strings.Split(sessionToken, ":")
		if len(tokenParts) != 2 {
			continue
		}

		partitionRangeId, err := strconv.ParseInt(tokenParts[0], 10, 64)
		if err != nil {
			continue
		}

		sc.ridToSessionTokens[containerRid][partitionRangeId] = tokenParts[1]
	}
}

func (sc *sessionContainer) ClearSessionToken(resourceAddress string) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()

	// Get the container path from the resource address e.g. dbs/testdb/colls/testcoll
	containerPath, err := getContainerPath(resourceAddress)
	if err != nil {
		return
	}

	// Resolve container resource ID
	rid, ok := sc.containerPathToRid[containerPath]
	if !ok {
		return
	}

	delete(sc.ridToSessionTokens, rid)
	delete(sc.containerPathToRid, containerPath)
}
