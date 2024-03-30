// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package conn

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type (
	// ConnectionStringProperties are the properties of a connection string
	// as returned by [ParseConnectionString].
	ConnectionStringProperties struct {
		// Endpoint is the Endpoint value in the connection string.
		// Ex: sb://example.servicebus.windows.net
		Endpoint string

		// EntityPath is EntityPath value in the connection string.
		EntityPath *string

		// FullyQualifiedNamespace is the Endpoint value without the protocol scheme.
		// Ex: example.servicebus.windows.net
		FullyQualifiedNamespace string

		// SharedAccessKey is the SharedAccessKey value in the connection string.
		SharedAccessKey *string

		// SharedAccessKeyName is the SharedAccessKeyName value in the connection string.
		SharedAccessKeyName *string

		// SharedAccessSignature is the SharedAccessSignature value in the connection string.
		SharedAccessSignature *string

		// Emulator indicates that the connection string is for an emulator:
		// ex: Endpoint=localhost:6765;SharedAccessKeyName=<< REDACTED >>;SharedAccessKey=<< REDACTED >>;UseDevelopmentEmulator=true
		Emulator bool
	}
)

// ParseConnectionString takes a string connection string from the Azure portal and returns the parsed representation.
// The method will return an error if the Endpoint, SharedAccessKeyName or SharedAccessKey is empty.
func ParseConnectionString(connStr string) (ConnectionStringProperties, error) {
	const (
		endpointKey              = "Endpoint"
		sharedAccessKeyNameKey   = "SharedAccessKeyName"
		sharedAccessKeyKey       = "SharedAccessKey"
		entityPathKey            = "EntityPath"
		sharedAccessSignatureKey = "SharedAccessSignature"
		useEmulator              = "UseDevelopmentEmulator"
	)

	csp := ConnectionStringProperties{}

	splits := strings.Split(connStr, ";")

	for _, split := range splits {
		if split == "" {
			continue
		}

		keyAndValue := strings.SplitN(split, "=", 2)
		if len(keyAndValue) < 2 {
			return ConnectionStringProperties{}, errors.New("failed parsing connection string due to unmatched key value separated by '='")
		}

		// if a key value pair has `=` in the value, recombine them
		key := keyAndValue[0]
		value := strings.Join(keyAndValue[1:], "=")
		switch {
		case strings.EqualFold(endpointKey, key):
			u, err := url.Parse(value)
			if err != nil {
				return ConnectionStringProperties{}, errors.New("failed parsing connection string due to an incorrectly formatted Endpoint value")
			}
			csp.Endpoint = value
			csp.FullyQualifiedNamespace = u.Host
		case strings.EqualFold(sharedAccessKeyNameKey, key):
			csp.SharedAccessKeyName = &value
		case strings.EqualFold(sharedAccessKeyKey, key):
			csp.SharedAccessKey = &value
		case strings.EqualFold(entityPathKey, key):
			csp.EntityPath = &value
		case strings.EqualFold(sharedAccessSignatureKey, key):
			csp.SharedAccessSignature = &value
		case strings.EqualFold(useEmulator, key):
			v, err := strconv.ParseBool(value)

			if err != nil {
				return ConnectionStringProperties{}, err
			}

			csp.Emulator = v
		}
	}

	if csp.Emulator && !strings.HasPrefix(csp.Endpoint, "sb://localhost:") {
		return ConnectionStringProperties{}, fmt.Errorf("UseEmulator=true can only be used with sb://localhost:<port>, not %s", csp.Endpoint)
	}

	if csp.FullyQualifiedNamespace == "" {
		return ConnectionStringProperties{}, fmt.Errorf("key %q must not be empty", endpointKey)
	}

	if csp.SharedAccessSignature == nil && csp.SharedAccessKeyName == nil {
		return ConnectionStringProperties{}, fmt.Errorf("key %q must not be empty", sharedAccessKeyNameKey)
	}

	if csp.SharedAccessKey == nil && csp.SharedAccessSignature == nil {
		return ConnectionStringProperties{}, fmt.Errorf("key %q or %q cannot both be empty", sharedAccessKeyKey, sharedAccessSignatureKey)
	}

	return csp, nil
}
