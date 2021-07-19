package aztable

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var AADAuthenticationScope = "https://storage.azure.com/.default"

func Test_TestProxy(t *testing.T) {
	require := require.New(t)
	err := recording.StartRecording(t)
	require.NoError(err)
	defer recording.StopRecording(t)
	if err != nil {
		fmt.Println(err)
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(err)
	testProxyTransport := recording.TestProxyTransport{}
	options := TableClientOptions{
		Scopes:     []string{AADAuthenticationScope},
		HTTPClient: testProxyTransport,
	}
	client, err := NewTableClient("testproxy", "https://seankaneprimx.table.core.windows.net", cred, &options)
	require.NoError(err)

	_, err = client.Create(ctx)
	require.NoError(err)

	client.Delete(ctx)
	require.NoError(err)
}
