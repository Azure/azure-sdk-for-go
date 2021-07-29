package aztable

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var AADAuthenticationScope = "https://storage.azure.com/.default"

func Test_TestProxyPolicy(t *testing.T) {
	require := require.New(t)
	err := recording.StartRecording(t, nil)
	require.NoError(err)
	defer recording.StopRecording(t, nil)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(err)
	recordingPolicy := recording.NewRecordingPolicy(&recording.RecordingOptions{UseHTTPS: true})
	options := TableClientOptions{
		Scopes:         []string{AADAuthenticationScope},
		PerCallOptions: []azcore.Policy{recordingPolicy},
	}
	client, err := NewTableClient("testproxy", "https://seankaneprim.table.core.windows.net", cred, &options)
	require.NoError(err)

	_, err = client.Create(ctx)
	require.NoError(err)

	_, err = client.Delete(ctx)
	require.NoError(err)
}
