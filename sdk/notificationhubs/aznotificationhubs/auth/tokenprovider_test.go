//go:build go1.20
// +build go1.20

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package auth

import (
	"strings"
	"testing"
)

const (
	validConnectionString = "Endpoint=sb://my-namespace.servicebus.windows.net/;SharedAccessKeyName=key-name;SharedAccessKey=secret="
)

func TestParseConnectionString(t *testing.T) {
	parsedConnection, err := FromConnectionString(validConnectionString)
	if parsedConnection == nil || err != nil {
		t.Fatalf(`FromConnectionString = %q, %v`, parsedConnection, err)
	}

	if !strings.EqualFold(parsedConnection.Endpoint, "sb://my-namespace.servicebus.windows.net/") {
		t.Fatalf(`ParsedConnection.EndPoint = %q`, parsedConnection.Endpoint)
	}

	if !strings.EqualFold(parsedConnection.KeyName, "key-name") {
		t.Fatalf(`ParsedConnection.KeyName = %q`, parsedConnection.KeyName)
	}

	if !strings.EqualFold(parsedConnection.KeyValue, "secret=") {
		t.Fatalf(`ParsedConnection.KeyValue = %q`, parsedConnection.KeyValue)
	}
}
