// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznamespaces

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStringize(t *testing.T) {
	require.Equal(t, "true", mustStringize(t, true))

	require.Equal(t, "1", mustStringize(t, 1))
	require.Equal(t, "1", mustStringize(t, int8(1)))
	require.Equal(t, "1", mustStringize(t, int16(1)))
	require.Equal(t, "1", mustStringize(t, int32(1)))
	require.Equal(t, "1", mustStringize(t, int64(1)))

	require.Equal(t, "1", mustStringize(t, uint8(1)))
	require.Equal(t, "1", mustStringize(t, uint16(1)))
	require.Equal(t, "1", mustStringize(t, uint32(1)))
	require.Equal(t, "1", mustStringize(t, uint64(1)))

	require.Equal(t, "AQID", mustStringize(t, []byte{1, 2, 3}))

	u, err := url.Parse("https://microsoft.com/hello?query=1")
	require.NoError(t, err)
	require.Equal(t, "https://microsoft.com/hello?query=1", mustStringize(t, u))

	require.Equal(t, "hello world", mustStringize(t, "hello world"))

	require.Equal(t, "0001-01-01T00:00:00Z", mustStringize(t, time.Time{}))
	require.Equal(t, "0001-01-01T00:00:00Z", mustStringize(t, &time.Time{}))

	require.Equal(t, "hello", mustStringize(t, stringableType{}))

	type customType struct{}
	_, err = stringize(customType{})
	require.EqualError(t, err, "type aznamespaces.customType cannot be converted to a string")
}

type stringableType struct{}

func (st stringableType) String() string {
	return "hello"
}

func mustStringize(t *testing.T, v any) string {
	str, err := stringize(v)
	require.NoError(t, err)
	return str
}
