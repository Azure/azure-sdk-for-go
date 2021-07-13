package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructAtomPath(t *testing.T) {
	baseUrl := ConstructAtomPath("/something", 1, 2)

	// I'm assuming the ordering is non-deterministic since the underlying values are just a map
	assert.Truef(t, baseUrl == "/something?%24skip=1&%24top=2" || baseUrl == "/something?%24top=2&%24skip=1", "%s wasn't one of our two variations", baseUrl)

	baseUrl = ConstructAtomPath("/something", 0, -1)
	assert.EqualValues(t, "/something", baseUrl, "Values <= 0 are ignored")

	baseUrl = ConstructAtomPath("/something", -1, 0)
	assert.EqualValues(t, "/something", baseUrl, "Values <= 0 are ignored")
}
