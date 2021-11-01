package armcompute_test

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
	"hash/fnv"
	"os"
	"testing"
)

var pathToPackage = "sdk/resourcemanager/compute/armcompute"

func TestMain(m *testing.M) {
	// Initialize

	// Run
	exitVal := m.Run()

	// cleanup
	os.Exit(exitVal)
}

func startTest(t *testing.T) func() {
	err := recording.StartRecording(t, pathToPackage, nil)
	require.NoError(t, err)
	return func() {
		err := recording.StopRecording(t, nil)
		require.NoError(t, err)
	}
}

func createRandomName(t *testing.T, prefix string) (string, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(t.Name()))
	return prefix + fmt.Sprint(h.Sum32()), err
}