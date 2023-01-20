// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package exported

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// ConnectionStringProperties are the properties of a connection string
// as returned by [NewConnectionStringProperties].
type ConnectionStringProperties struct {
	// Endpoint represents the Endpoint value in the connection string.
	// Ex: sb://example.servicebus.windows.net
	Endpoint string

	// EventHubName is EntityPath value in the connection string.
	EventHubName string

	// FullyQualifiedNamespace is the Endpoint value without the protocol scheme.
	// Ex: example.servicebus.windows.net
	FullyQualifiedNamespace string

	// SharedAccessKey is the SharedAccessKey value in the connection string.
	SharedAccessKey string

	// SharedAccessKeyName is the SharedAccessKeyName value in the connection string.
	SharedAccessKeyName string

	// SharedAccessSignature is the SharedAccessSignature value in the connection string.
	SharedAccessSignature string
}

// NewConnectionStringProperties takes a connection string from the Azure portal and returns the
// parsed representation The method will return an error if the Endpoint, SharedAccessKeyName
// or SharedAccessKey is empty.
func NewConnectionStringProperties(connStr string) (ConnectionStringProperties, error) {
	const (
		endpointKey              = "Endpoint"
		sharedAccessKeyNameKey   = "SharedAccessKeyName"
		sharedAccessKeyKey       = "SharedAccessKey"
		entityPathKey            = "EntityPath"
		sharedAccessSignatureKey = "SharedAccessSignature"
	)

	// We can parse two types of connection strings.
	// 1. Connection strings generated from the portal (or elsewhere) that contain an embedded key and keyname.
	// 2. A specially formatted connection string with an embedded SharedAccessSignature:
	//   Endpoint=sb://<sb>.servicebus.windows.net;SharedAccessSignature=SharedAccessSignature sr=<sb>.servicebus.windows.net&sig=<base64-sig>&se=<expiry>&skn=<keyname>"
	//var namespace, hubName, keyName, secret, sas string

	csp := ConnectionStringProperties{}

	splits := strings.Split(connStr, ";")

	for _, split := range splits {
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
			csp.SharedAccessKeyName = value
			// keyName = value
		case strings.EqualFold(sharedAccessKeyKey, key):
			csp.SharedAccessKey = value
			// secret = value
		case strings.EqualFold(entityPathKey, key):
			csp.EventHubName = value
			//hubName = value
		case strings.EqualFold(sharedAccessSignatureKey, key):
			csp.SharedAccessSignature = value
			// sas = value
		}
	}

	if csp.FullyQualifiedNamespace == "" {
		return ConnectionStringProperties{}, fmt.Errorf("key %q must not be empty", endpointKey)
	}

	if csp.SharedAccessSignature == "" && csp.SharedAccessKeyName == "" {
		return ConnectionStringProperties{}, fmt.Errorf("key %q must not be empty", sharedAccessKeyNameKey)
	}

	if csp.SharedAccessKey == "" && csp.SharedAccessSignature == "" {
		return ConnectionStringProperties{}, fmt.Errorf("key %q or %q cannot both be empty", sharedAccessKeyKey, sharedAccessSignatureKey)
	}

	return csp, nil
}
