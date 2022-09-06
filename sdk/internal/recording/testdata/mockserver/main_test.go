//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// This is a mock server used to test the sanitizers in sanitizer.go

package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	indexHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var unmarshalled map[string]string
	err = json.Unmarshal(data, &unmarshalled)

	require.NoError(t, err)
	require.Equal(t, unmarshalled["Tag"], "Value")
	require.Equal(t, unmarshalled["Tag2"], "Value2")
	require.Equal(t, unmarshalled["Tag3"], "https://storageaccount.table.core.windows.net/")
}
