package aztable

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var AADAuthenticationScope = "https://storage.azure.com/.default"

func createTableClientForRecording(t *testing.T, tableName string, serviceURL string) (*TableClient, error) {
	policy := recording.NewRecordingPolicy(&recording.RecordingOptions{UseHTTPS: true})
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)
	options := &TableClientOptions{
		Scopes:         []string{AADAuthenticationScope},
		PerCallOptions: []azcore.Policy{policy},
	}
	return NewTableClient(tableName, serviceURL, cred, options)
}

func createTableServiceClientForRecording(t *testing.T, serviceURL string) (*TableServiceClient, error) {
	policy := recording.NewRecordingPolicy(&recording.RecordingOptions{UseHTTPS: true})
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)
	options := &TableClientOptions{
		Scopes:         []string{AADAuthenticationScope},
		PerCallOptions: []azcore.Policy{policy},
	}
	return NewTableServiceClient(serviceURL, cred, options)
}

func Test_TestProxyPolicy(t *testing.T) {
	require := require.New(t)
	err := recording.StartRecording(t, nil)
	require.NoError(err)
	defer recording.StopRecording(t, nil)

	client, err := createTableClientForRecording(t, "testproxy", "https://seankaneprim.table.core.windows.net")
	require.NoError(err)

	_, err = client.Create(ctx)
	require.NoError(err)

	_, err = client.Delete(ctx)
	require.NoError(err)
}
