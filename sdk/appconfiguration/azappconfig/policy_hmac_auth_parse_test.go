//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azappconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHmacAuthParseConnectionString(t *testing.T) {
	ep, id, sc, err := parseConnectionString("Endpoint=xX;Id=yY;Secret=zZ")
	require.Empty(t, err)
	require.Equal(t, ep, "xX")
	require.Equal(t, id, "yY")
	require.Len(t, sc, 2)
	require.Equal(t, sc[0], 'z')
	require.Equal(t, sc[1], 'Z')
}

func TestHmacAuthParseConnectionStringMixedOrder(t *testing.T) {
	ep, id, sc, err := parseConnectionString("Id=yY;Secret=zZ;Endpoint=xX")
	require.Empty(t, err)
	require.Equal(t, ep, "xX")
	require.Equal(t, id, "yY")
	require.Len(t, sc, 2)
	require.Equal(t, sc[0], 'z')
	require.Equal(t, sc[1], 'Z')
}

func TestHmacAuthParseConnectionStringExtraProperties(t *testing.T) {
	ep, id, sc, err := parseConnectionString("A=aA;Endpoint=xX;B=bB;Id=yY;C=cC;Secret=zZ;D=dD;")
	require.Empty(t, err)
	require.Equal(t, ep, "xX")
	require.Equal(t, id, "yY")
	require.Len(t, sc, 2)
	require.Equal(t, sc[0], 'z')
	require.Equal(t, sc[1], 'Z')
}

func TestHmacAuthParseConnectionStringMissingEndoint(t *testing.T) {
	_, _, _, err := parseConnectionString("Id=yY;Secret=zZ")
	require.NotEmpty(t, err)
}

func TestHmacAuthParseConnectionStringMissingId(t *testing.T) {
	_, _, _, err := parseConnectionString("Endpoint=xX;Secret=zZ")
	require.NotEmpty(t, err)
}

func TestHmacAuthParseConnectionStringMissingSecret(t *testing.T) {
	_, _, _, err := parseConnectionString("Endpoint=xX;Id=yY")
	require.NotEmpty(t, err)
}

func TestHmacAuthParseConnectionStringDuplicateEndoint(t *testing.T) {
	_, _, _, err := parseConnectionString("Endpoint=xX;Endpoint=xX;Id=yY;Secret=zZ")
	require.NotEmpty(t, err)
}

func TestHmacAuthParseConnectionStringDuplicateId(t *testing.T) {
	_, _, _, err := parseConnectionString("Endpoint=xX;Id=yY;Id=yY;Secret=zZ")
	require.NotEmpty(t, err)
}

func TestHmacAuthParseConnectionStringDuplicateSecret(t *testing.T) {
	_, _, _, err := parseConnectionString("Endpoint=xX;Id=yY;Secret=zZ;Secret=zZ")
	require.NotEmpty(t, err)
}
