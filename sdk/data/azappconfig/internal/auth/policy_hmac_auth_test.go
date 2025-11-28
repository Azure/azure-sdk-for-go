// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseConnectionString(t *testing.T) {
	ep, id, sc, err := ParseConnectionString("Endpoint=xX;Id=yY;Secret=ZmZm")
	require.NoError(t, err)
	require.Equal(t, "xX", ep)
	require.Equal(t, "yY", id)

	require.Len(t, sc, 3)
	require.Equal(t, byte('f'), sc[0])
	require.Equal(t, byte('f'), sc[1])
	require.Equal(t, byte('f'), sc[2])
}

func TestParseConnectionStringMixedOrder(t *testing.T) {
	ep, id, sc, err := ParseConnectionString("Id=yY;Secret=ZmZm;Endpoint=xX")
	require.NoError(t, err)
	require.Equal(t, "xX", ep)
	require.Equal(t, "yY", id)

	require.Len(t, sc, 3)
	require.Equal(t, byte('f'), sc[0])
	require.Equal(t, byte('f'), sc[1])
	require.Equal(t, byte('f'), sc[2])
}

func TestParseConnectionStringExtraProperties(t *testing.T) {
	ep, id, sc, err := ParseConnectionString("A=aA;Endpoint=xX;B=bB;Id=yY;C=cC;Secret=ZmZm;D=dD;")
	require.NoError(t, err)
	require.Equal(t, "xX", ep)
	require.Equal(t, "yY", id)

	require.Len(t, sc, 3)
	require.Equal(t, byte('f'), sc[0])
	require.Equal(t, byte('f'), sc[1])
	require.Equal(t, byte('f'), sc[2])
}

func TestParseConnectionStringMissingEndoint(t *testing.T) {
	_, _, _, err := ParseConnectionString("Id=yY;Secret=ZmZm")
	require.Error(t, err)
	require.ErrorContains(t, err, "missing Endpoint")
}

func TestParseConnectionStringMissingId(t *testing.T) {
	_, _, _, err := ParseConnectionString("Endpoint=xX;Secret=ZmZm")
	require.Error(t, err)
	require.ErrorContains(t, err, "missing Id")
}

func TestParseConnectionStringMissingSecret(t *testing.T) {
	_, _, _, err := ParseConnectionString("Endpoint=xX;Id=yY")
	require.Error(t, err)
	require.ErrorContains(t, err, "missing Secret")
}

func TestParseConnectionStringDuplicateEndoint(t *testing.T) {
	_, _, _, err := ParseConnectionString("Endpoint=xX;Endpoint=xX;Id=yY;Secret=ZmZm")
	require.Error(t, err)
	require.ErrorContains(t, err, "duplicate Endpoint")
}

func TestParseConnectionStringDuplicateId(t *testing.T) {
	_, _, _, err := ParseConnectionString("Endpoint=xX;Id=yY;Id=yY;Secret=ZmZm")
	require.Error(t, err)
	require.ErrorContains(t, err, "duplicate Id")
}

func TestParseConnectionStringDuplicateSecret(t *testing.T) {
	_, _, _, err := ParseConnectionString("Endpoint=xX;Id=yY;Secret=ZmZm;Secret=zZ")
	require.Error(t, err)
	require.ErrorContains(t, err, "duplicate Secret")
}

func TestParseConnectionStringInvalidEncoding(t *testing.T) {
	_, _, _, err := ParseConnectionString("Endpoint=xX;Id=yY;Secret=badencoding")
	require.Error(t, err)
	require.ErrorContains(t, err, "illegal base64 data")
}
