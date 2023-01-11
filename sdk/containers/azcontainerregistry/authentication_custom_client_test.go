package azcontainerregistry

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_newAuthenticationClient(t *testing.T) {
	client := newAuthenticationClient("test", nil)
	require.NotNil(t, client)
}
