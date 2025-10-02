//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package synctoken

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azappconfig/v2/internal/exported"
)

// Cache contains a collection of sync token values.
// Methods on Cache are safe for concurrent use.
// Don't use this type directly, use NewCache() instead.
type Cache struct {
	tokens   map[string]exported.SyncTokenValues
	tokensMu *sync.RWMutex
}

// NewCache creates a new instance of [Cache].
func NewCache() *Cache {
	return &Cache{
		tokens:   map[string]exported.SyncTokenValues{},
		tokensMu: &sync.RWMutex{},
	}
}

// Set adds or updates the cache with the provided sync token.
func (s *Cache) Set(syncToken exported.SyncToken) error {
	tokens, err := exported.ParseSyncToken(syncToken)
	if err != nil {
		return err
	}

	s.tokensMu.Lock()
	defer s.tokensMu.Unlock()

	for _, token := range tokens {
		if tk, ok := s.tokens[token.ID]; ok {
			// we already have a sync token for this ID.
			// if the current token is already at this version
			// or newer don't update the map.
			if tk.Version >= token.Version {
				continue
			}
		}

		s.tokens[token.ID] = token
	}

	return nil
}

// Get returns a sync token representing the current state of the cache.
// Format is "<id1>=<value1>,<id2>=<value2>,..."
func (s *Cache) Get() string {
	s.tokensMu.RLock()
	defer s.tokensMu.RUnlock()

	if len(s.tokens) == 0 {
		return ""
	}
	tokens := []string{}
	for _, token := range s.tokens {
		tokens = append(tokens, fmt.Sprintf("%s=%s", token.ID, token.Value))
	}
	return strings.Join(tokens, ",")
}
