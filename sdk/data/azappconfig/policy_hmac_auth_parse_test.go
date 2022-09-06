//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azappconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHmacAuthParseConnectionString(t *testing.T) {
	ep, id, sc, err := parseConnectionString("Endpoint=xX;Id=yY;Secret=ZmZm")
	require.NoError(t, err)
	require.Equal(t, "xX", ep)
	require.Equal(t, "yY", id)

	require.Len(t, sc, 3)
	require.Equal(t, byte('f'), sc[0])
	require.Equal(t, byte('f'), sc[1])
	require.Equal(t, byte('f'), sc[2])
}

func TestHmacAuthParseConnectionStringMixedOrder(t *testing.T) {
	ep, id, sc, err := parseConnectionString("Id=yY;Secret=ZmZm;Endpoint=xX")
	require.NoError(t, err)
	require.Equal(t, "xX", ep)
	require.Equal(t, "yY", id)

	require.Len(t, sc, 3)
	require.Equal(t, byte('f'), sc[0])
	require.Equal(t, byte('f'), sc[1])
	require.Equal(t, byte('f'), sc[2])
}

func TestHmacAuthParseConnectionStringExtraProperties(t *testing.T) {
	ep, id, sc, err := parseConnectionString("A=aA;Endpoint=xX;B=bB;Id=yY;C=cC;Secret=ZmZm;D=dD;")
	require.NoError(t, err)
	require.Equal(t, "xX", ep)
	require.Equal(t, "yY", id)

	require.Len(t, sc, 3)
	require.Equal(t, byte('f'), sc[0])
	require.Equal(t, byte('f'), sc[1])
	require.Equal(t, byte('f'), sc[2])
}

func TestHmacAuthParseConnectionStringMissingEndoint(t *testing.T) {
	_, _, _, err := parseConnectionString("Id=yY;Secret=ZmZm")
	require.Error(t, err)
}

func TestHmacAuthParseConnectionStringMissingId(t *testing.T) {
	_, _, _, err := parseConnectionString("Endpoint=xX;Secret=ZmZm")
	require.Error(t, err)
}

func TestHmacAuthParseConnectionStringMissingSecret(t *testing.T) {
	_, _, _, err := parseConnectionString("Endpoint=xX;Id=yY")
	require.Error(t, err)
}

func TestHmacAuthParseConnectionStringDuplicateEndoint(t *testing.T) {
	_, _, _, err := parseConnectionString("Endpoint=xX;Endpoint=xX;Id=yY;Secret=ZmZm")
	require.Error(t, err)
}

func TestHmacAuthParseConnectionStringDuplicateId(t *testing.T) {
	_, _, _, err := parseConnectionString("Endpoint=xX;Id=yY;Id=yY;Secret=ZmZm")
	require.Error(t, err)
}

func TestHmacAuthParseConnectionStringDuplicateSecret(t *testing.T) {
	_, _, _, err := parseConnectionString("Endpoint=xX;Id=yY;Secret=ZmZm;Secret=zZ")
	require.Error(t, err)
}
