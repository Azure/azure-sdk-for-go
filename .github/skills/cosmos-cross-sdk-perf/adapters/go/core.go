// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	documentBytes = 1024
	seedCount     = 1024
)

type document struct {
	ID      string `json:"id"`
	PK      string `json:"pk"`
	Owner   string `json:"owner"`
	Padding string `json:"padding"`
}

func itemID(owner string, index int) string {
	return fmt.Sprintf("%s-%04d", owner, index)
}

func itemPK(owner string, index int) string {
	return itemID(owner, index)
}

func makeDocument(owner string, index int) ([]byte, error) {
	if err := validateOwner(owner); err != nil {
		return nil, err
	}
	value := document{ID: itemID(owner, index), PK: itemPK(owner, index), Owner: owner}
	paddingLength := expectedPaddingLength(owner, index)
	if paddingLength < 0 {
		return nil, fmt.Errorf("owner %q leaves no room for document padding", owner)
	}
	value.Padding = strings.Repeat("x", paddingLength)
	body, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	if len(body) != documentBytes {
		return nil, fmt.Errorf("document length is %d, expected %d", len(body), documentBytes)
	}
	return body, nil
}

func validateDocument(body []byte, owner string, index int) error {
	var value document
	if err := json.Unmarshal(body, &value); err != nil {
		return fmt.Errorf("decode document: %w", err)
	}
	if value.ID != itemID(owner, index) || value.PK != itemPK(owner, index) || value.Owner != owner {
		return errors.New("document id, pk, or owner does not match the campaign")
	}
	if len(value.Padding) != expectedPaddingLength(owner, index) || strings.Trim(value.Padding, "x") != "" {
		return errors.New("document padding does not match the exact-size campaign payload")
	}
	return nil
}

func expectedPaddingLength(owner string, index int) int {
	const emptyDocumentBytes = len(`{"id":"","pk":"","owner":"","padding":""}`)
	id := itemID(owner, index)
	return documentBytes - emptyDocumentBytes - 2*len(id) - len(owner)
}

func validateOwner(owner string) error {
	if len(owner) < 8 || len(owner) > 64 {
		return errors.New("owner must contain 8 through 64 characters")
	}
	for _, value := range owner {
		if (value >= 'a' && value <= 'z') || (value >= '0' && value <= '9') || value == '-' {
			continue
		}
		return errors.New("owner may contain only lowercase ASCII letters, digits, and hyphens")
	}
	return nil
}

func keyIndex(sequence uint64) int {
	return int(sequence % seedCount)
}

func ownedBy(body []byte, owner string) (bool, error) {
	var value struct {
		Owner string `json:"owner"`
	}
	if err := json.Unmarshal(body, &value); err != nil {
		return false, err
	}
	return value.Owner == owner, nil
}

func rejectSecretArguments(arguments []string) error {
	for index, argument := range arguments {
		lower := strings.ToLower(argument)
		if lower == "--key" || lower == "--account-key" || strings.HasPrefix(lower, "--key=") ||
			strings.HasPrefix(lower, "--account-key=") {
			return fmt.Errorf("secret command-line argument %q is forbidden", arguments[index])
		}
	}
	for _, name := range []string{"AZURE_COSMOS_KEY", "COSMOS_KEY", "COSMOS_ACCOUNT_KEY", "ACCOUNT_KEY"} {
		if os.Getenv(name) != "" {
			return fmt.Errorf("secret environment variable %s is forbidden; use --key-fd", name)
		}
	}
	return nil
}

func readSecretFD(fdValue string) ([]byte, error) {
	if fdValue == "" {
		return nil, errors.New("--key-fd is required")
	}
	fd, err := strconv.Atoi(fdValue)
	if err != nil || fd < 3 {
		return nil, errors.New("--key-fd must be an inherited descriptor number of at least 3")
	}
	file := os.NewFile(uintptr(fd), "cosmos-account-key")
	if file == nil {
		return nil, errors.New("key descriptor is invalid")
	}
	data, err := io.ReadAll(io.LimitReader(file, 16*1024))
	if err != nil {
		return nil, fmt.Errorf("read key descriptor: %w", err)
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil, errors.New("key descriptor was empty")
	}
	return data, nil
}

func clearBytes(value []byte) {
	for index := range value {
		value[index] = 0
	}
}
