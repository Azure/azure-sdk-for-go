//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package synctoken

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Cache contains a collection sync token values.
// Methods on Cache are safe for concurrent use.
// Don't use this type directly, use NewCache() instead.
type Cache struct {
	tokens   map[string]syncTokenInfo
	tokensMu *sync.RWMutex
}

// NewCache creates a new instance of [Cache].
func NewCache() *Cache {
	return &Cache{
		tokens:   map[string]syncTokenInfo{},
		tokensMu: &sync.RWMutex{},
	}
}

// Set adds or updates the cache with the provided sync token.
func (s *Cache) Set(syncToken string) error {
	tokens := strings.TrimSpace(syncToken)
	if tokens == "" {
		return errors.New("syncToken can't be empty")
	}

	s.tokensMu.Lock()
	defer s.tokensMu.Unlock()

	// token format is "<id>=<value>;sn=<sn>" and can contain multiple, comman-delimited values
	for _, token := range strings.Split(tokens, ",") {
		items := strings.Split(token, ";")
		if len(items) != 2 {
			return fmt.Errorf("invalid token %s", token)
		}

		// items[0] contains "<id>=<value>"
		// note that <value> is a base-64 encoded string, so don't try to split on '='
		assignmentIndex := strings.Index(items[0], "=")
		if assignmentIndex < 0 {
			return fmt.Errorf("unexpected token format %s", items[0])
		}
		tokenID := strings.TrimSpace(items[0][:assignmentIndex])
		tokenValue := strings.TrimSpace(items[0][assignmentIndex+1:])

		// items[1] contains "sn=<sn>"
		// parse the version number after the equals sign
		assignmentIndex = strings.Index(items[1], "=")
		if assignmentIndex < 0 {
			return fmt.Errorf("unexpected token version format %s", items[1])
		}
		tokenVersion, err := strconv.ParseInt(strings.TrimSpace(items[1][assignmentIndex+1:]), 10, 64)
		if err != nil {
			return err
		}

		if tk, ok := s.tokens[tokenID]; ok {
			// we already have a sync token for this ID.
			// if the current token is already at this version
			// or newer there's no need to update the map.
			if tk.Version >= tokenVersion {
				continue
			}
		}

		s.tokens[tokenID] = syncTokenInfo{
			Value:   tokenValue,
			Version: tokenVersion,
		}
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
	for k, v := range s.tokens {
		tokens = append(tokens, fmt.Sprintf("%s=%s", k, v.Value))
	}
	return strings.Join(tokens, ",")
}

type syncTokenInfo struct {
	Value   string
	Version int64
}
