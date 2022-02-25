// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package conn

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	endpointKey            = "Endpoint"
	sharedAccessKeyNameKey = "SharedAccessKeyName"
	sharedAccessKeyKey     = "SharedAccessKey"
	entityPathKey          = "EntityPath"
)

type (
	// ParsedConn is the structure of a parsed Service Bus or Event Hub connection string.
	ParsedConn struct {
		Namespace string
		HubName   string
		KeyName   string
		Key       string
	}
)

// newParsedConnection is a constructor for a parsedConn and verifies each of the inputs is non-null.
// namespace is the FQDN of the namespace
func newParsedConnection(namespace, hubName, keyName, key string) *ParsedConn {
	return &ParsedConn{
		Namespace: namespace,
		KeyName:   keyName,
		Key:       key,
		HubName:   hubName,
	}
}

// ParsedConnectionFromStr takes a string connection string from the Azure portal and returns the parsed representation.
// The method will return an error if the Endpoint, SharedAccessKeyName or SharedAccessKey is empty.
func ParsedConnectionFromStr(connStr string) (*ParsedConn, error) {
	var namespace, hubName, keyName, secret string
	splits := strings.Split(connStr, ";")
	for _, split := range splits {
		keyAndValue := strings.Split(split, "=")
		if len(keyAndValue) < 2 {
			return nil, errors.New("failed parsing connection string due to unmatched key value separated by '='")
		}

		// if a key value pair has `=` in the value, recombine them
		key := keyAndValue[0]
		value := strings.Join(keyAndValue[1:], "=")
		switch {
		case strings.EqualFold(endpointKey, key):
			u, err := url.Parse(value)
			if err != nil {
				return nil, errors.New("failed parsing connection string due to an incorrectly formatted Endpoint value")
			}
			namespace = u.Host
		case strings.EqualFold(sharedAccessKeyNameKey, key):
			keyName = value
		case strings.EqualFold(sharedAccessKeyKey, key):
			secret = value
		case strings.EqualFold(entityPathKey, key):
			hubName = value
		}
	}

	parsed := newParsedConnection(namespace, hubName, keyName, secret)
	if namespace == "" {
		return parsed, fmt.Errorf("key %q must not be empty", endpointKey)
	}

	if keyName == "" {
		return parsed, fmt.Errorf("key %q must not be empty", sharedAccessKeyNameKey)
	}

	if secret == "" {
		return parsed, fmt.Errorf("key %q must not be empty", sharedAccessKeyKey)
	}

	return parsed, nil
}
