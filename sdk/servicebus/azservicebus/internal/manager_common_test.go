package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructAtomPath(t *testing.T) {
	basePath := ConstructAtomPath("/something", 1, 2)

	// I'm assuming the ordering is non-deterministic since the underlying values are just a map
	assert.Truef(t, basePath == "/something?%24skip=1&%24top=2" || basePath == "/something?%24top=2&%24skip=1", "%s wasn't one of our two variations", basePath)

	basePath = ConstructAtomPath("/something", 0, -1)
	assert.EqualValues(t, "/something", basePath, "Values <= 0 are ignored")

	basePath = ConstructAtomPath("/something", -1, 0)
	assert.EqualValues(t, "/something", basePath, "Values <= 0 are ignored")
}
