package perf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommaIze(t *testing.T) {
	type testStruct struct {
		s string
		i int
	}
	for _, pair := range []testStruct{
		{"0", 0},
		{"100", 100},
		{"1,000", 1000},
		{"10,000", 10000},
		{"100,000", 100000},
		{"1,000,000", 1000000},
		{"1,000,000,000", 1000000000},
		{"1,234,567,890", 1234567890},
		{"987,654,321", 987654321},
	} {
		require.Equal(t, pair.s, commaIze(pair.i))
	}
}
