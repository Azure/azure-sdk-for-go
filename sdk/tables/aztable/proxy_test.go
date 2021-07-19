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

func Test_TestProxy(t *testing.T) {
	require := require.New(t)
	err := recording.StartRecording(t)
	require.Nil(err)
	defer recording.StopRecording(t)
	if err != nil {
		fmt.Println(err)
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.Nil(err)
	testProxyPolicy := recording.TestProxyPolicy{}
	options := TableClientOptions{
		Scopes:         []string{AADAuthenticationScope},
		PerCallOptions: []azcore.Policy{testProxyPolicy},
	}
	client, err := NewTableClient("testproxy", "https://seankaneprimx.table.core.windows.net", cred, &options)
	require.Nil(err)

	_, err = client.Create(ctx)
	require.Nil(err)

	client.Delete(ctx)
	require.Nil(err)
}
