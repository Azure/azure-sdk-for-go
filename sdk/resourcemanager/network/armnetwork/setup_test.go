package armnetwork_test

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var pathToPackage = "sdk/resourcemanager/network/armnetwork/testdata"

func TestMain(m *testing.M) {
	// Initialize

	// Run
	exitVal := m.Run()

	// cleanup
	os.Exit(exitVal)
}

func startTest(t *testing.T) func() {
	err := recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)
	return func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	}
}
