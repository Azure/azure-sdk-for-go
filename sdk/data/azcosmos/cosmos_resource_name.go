// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (c *Client) validateResourceName(resourceType, name string) error {
	if c.disableResourceNameValidation {
		return nil
	}
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, "/\\?#\t\n\r") {
		return fmt.Errorf("invalid %s ID: resource names must be non-empty, cannot equal '.' or '..', and cannot contain '/', '\\', '?', '#', tab, line feed, or carriage return", resourceType)
	}
	return nil
}

// validateItemBodyID inspects only the top-level, case-sensitive id property.
// Other properties, including nested id properties, remain opaque.
func (c *Client) validateItemBodyID(body []byte) error {
	if c.disableResourceNameValidation {
		return nil
	}
	if !json.Valid(body) {
		return errors.New("item body must be valid JSON")
	}
	decoder := json.NewDecoder(bytes.NewReader(body))
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	if token != json.Delim('{') {
		return errors.New("item body must be a JSON object")
	}
	seenID := false
	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			return err
		}
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			return err
		}
		if key != "id" {
			continue
		}
		if seenID {
			return errors.New("item body must not contain duplicate id properties")
		}
		seenID = true
		var id string
		if err := json.Unmarshal(value, &id); err != nil {
			return fmt.Errorf("item id must be a string: %w", err)
		}
		if err := c.validateResourceName("item", id); err != nil {
			return err
		}
	}
	return nil
}
