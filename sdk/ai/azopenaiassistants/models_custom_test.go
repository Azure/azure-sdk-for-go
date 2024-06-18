// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalStringOrObject(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		jsonString, err := json.Marshal("hello world")
		require.NoError(t, err)

		str, obj, err := unmarshalStringOrObject[any](jsonString)
		require.NoError(t, err)
		require.Empty(t, obj)
		require.Equal(t, "hello world", str)
	})

	t.Run("json object", func(t *testing.T) {
		type myType struct {
			Name string
		}
		jsonObj, err := json.Marshal(myType{Name: "hello world"})
		require.NoError(t, err)

		str, obj, err := unmarshalStringOrObject[myType](jsonObj)
		require.NoError(t, err)
		require.Empty(t, str)
		require.Equal(t, myType{Name: "hello world"}, *obj)
	})

	t.Run("errors", func(t *testing.T) {
		str, obj, err := unmarshalStringOrObject[any](nil)
		require.EqualError(t, err, "can't deserialize from an empty slice of bytes")
		require.Empty(t, obj)
		require.Empty(t, str)

		str, obj, err = unmarshalStringOrObject[any]([]byte("\"invalid json string"))
		require.EqualError(t, err, "unexpected end of JSON input")
		require.Empty(t, obj)
		require.Empty(t, str)

		str, obj, err = unmarshalStringOrObject[any]([]byte("{\"InvalidJSONObject\":}"))
		require.EqualError(t, err, "invalid character '}' looking for beginning of value")
		require.Empty(t, obj)
		require.Empty(t, str)
	})
}
