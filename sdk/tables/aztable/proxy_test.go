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
	err := recording.StartRecording(t)
	require.NoError(err)
	defer recording.StopRecording(t)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(err)
	recordingPolicy := recording.NewRecordingPolicy(nil)
	options := TableClientOptions{
		Scopes:         []string{AADAuthenticationScope},
		PerCallOptions: []azcore.Policy{recordingPolicy},
	}
	client, err := NewTableClient("testproxy", "https://seankaneprim.table.core.windows.net", cred, &options)
	require.NoError(err)

	fmt.Println("CALLING CREATE")
	_, err = client.Create(ctx)
	fmt.Println("Create err: ", err.Error())
	require.NoError(err)

	fmt.Println("CALLING DELETE")
	_, err = client.Delete(ctx)
	fmt.Println("Delete err: ", err)
	require.NoError(err)
}

func Test_TestProxyTransport(t *testing.T) {
	require := require.New(t)
	err := recording.StartRecording(t)
	require.NoError(err)
	defer recording.StopRecording(t)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(err)
	testProxyTransport := recording.TestProxyTransport{}
	options := TableClientOptions{
		Scopes:         []string{AADAuthenticationScope},
		HTTPClient: testProxyTransport,
	}
	client, err := NewTableClient("testproxy", "https://seankaneprim.table.core.windows.net", cred, &options)
	require.NoError(err)

	fmt.Println("CALLING CREATE")
	_, err = client.Create(ctx)
	fmt.Println("Create err: ", err.Error())
	require.NoError(err)

	fmt.Println("CALLING DELETE")
	_, err = client.Delete(ctx)
	fmt.Println("Delete err: ", err)
	require.NoError(err)
}
