package azcontainerregistry

import (
	"bytes"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
)

func TestCalculateDigest(t *testing.T) {
	payload := streaming.NopCloser(bytes.NewReader([]byte("test")))
	payload.Seek(3, io.SeekStart)
	result, err := CalculateDigest(payload)
	require.NoError(t, err)
	require.Equal(t, result, "sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")
	pos, err := payload.Seek(0, io.SeekCurrent)
	require.NoError(t, err)
	require.Equal(t, pos, int64(0))
}
