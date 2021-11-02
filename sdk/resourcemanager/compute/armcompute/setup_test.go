package armcompute_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/require"
	"hash/fnv"
	"net/http"
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

func cleanup(t *testing.T, client *armresources.ResourceGroupsClient, resourceGroupName string) {
	_, err := client.BeginDelete(context.Background(), resourceGroupName, nil)
	require.NoError(t, err)
}

func startTest(t *testing.T) func() {
	err := recording.Start(t, pathToPackage, nil)
	require.NoError(t, err)
	return func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	}
}

func createRandomName(t *testing.T, prefix string) (string, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(t.Name()))
	return prefix + fmt.Sprint(h.Sum32()), err
}

type recordingPolicy struct {
	options recording.RecordingOptions
	t       *testing.T
}

func NewRecordingPolicy(t *testing.T, o *recording.RecordingOptions) policy.Policy {
	if o == nil {
		o = &recording.RecordingOptions{UseHTTPS: true}
	}
	p := &recordingPolicy{options: *o, t: t}
	return p
}

func (p *recordingPolicy) Do(req *policy.Request) (resp *http.Response, err error) {
	if recording.GetRecordMode() != recording.LiveMode {
		originalURLHost := req.Raw().URL.Host
		req.Raw().URL.Scheme = "https"
		//req.Raw().URL.Host = p.options.Host
		//req.Raw().Host = p.options.Host
		scheme := "http"
		if p.options.UseHTTPS {
			scheme = "https"
		}

		req.Raw().Header.Set(recording.UpstreamURIHeader, fmt.Sprintf("%s://%s", scheme, originalURLHost))
		req.Raw().Header.Set(recording.ModeHeader, recording.GetRecordMode())
		req.Raw().Header.Set(recording.IDHeader, recording.GetRecordingId(p.t))
	}
	return req.Next()
}

//func authenticateTest(t *testing.T) (azcore.TokenCredential, *arm.ConnectionOptions) {
//	p := NewRecordingPolicy(t, &recording.RecordingOptions{UseHTTPS: true})
//	client, err := recording.GetHTTPClient(t)
//	require.NoError(t, err)
//
//	options := &arm.ConnectionOptions{
//		PerCallPolicies: []policy.Policy{p},
//		HTTPClient:       client,
//	}
//
//
//	var cred azcore.TokenCredential
//	if recording.GetRecordMode() != "playback" {
//		tenantId := lookupEnvVar("AZURE_TENANT_ID")
//		clientId := lookupEnvVar("AZURE_CLIENT_ID")
//		clientSecret := lookupEnvVar("AZURE_CLIENT_SECRET")
//		cred, err = azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
//		require.NoError(t, err)
//	} else {
//		cred = &fakeCredential{}
//	}
//
//	return cred, options
//}