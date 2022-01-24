package perf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommaIze(t *testing.T) {
	require.Equal(t, "0", commaIze(0))
	require.Equal(t, "1", commaIze(1))
	require.Equal(t, "10", commaIze(10))
	require.Equal(t, "100", commaIze(100))
	require.Equal(t, "1,000", commaIze(1000))
	require.Equal(t, "10,000", commaIze(10000))
	require.Equal(t, "100,000", commaIze(100000))
	require.Equal(t, "1,000,000", commaIze(1000000))
	require.Equal(t, "1,000,000,000", commaIze(1000000000))
}
