//go:build go1.16
// +build go1.16

package azsecrets

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

var pathToPackage = "sdk/keyvault/azsecrets"

type recordingPolicy struct {
	options recording.RecordingOptions
	t       *testing.T
}

func NewRecordingPolicy(t *testing.T, o *recording.RecordingOptions) policy.Policy {
	if o == nil {
		o = &recording.RecordingOptions{}
	}
	p := &recordingPolicy{options: *o, t: t}
	p.options.Init()
	return p
}

func (p *recordingPolicy) Do(req *policy.Request) (resp *http.Response, err error) {
	originalURLHost := req.Raw().URL.Host
	req.Raw().URL.Scheme = "https"
	req.Raw().URL.Host = p.options.Host
	req.Raw().Host = p.options.Host

	req.Raw().Header.Set(recording.UpstreamUriHeader, fmt.Sprintf("%v://%v", p.options.Scheme, originalURLHost))
	req.Raw().Header.Set(recording.ModeHeader, recording.GetRecordMode())
	req.Raw().Header.Set(recording.IdHeader, recording.GetRecordingId(p.t))

	return req.Next()
}

func createClient(t *testing.T, vaultUrl string) (*Client, error) {
	p := NewRecordingPolicy(t, &recording.RecordingOptions{UseHTTPS: true})
	client, err := recording.GetHTTPClient(t)
	require.NoError(t, err)

	options := &ClientOptions{
		PerCallPolicies: []policy.Policy{p},
		HTTPClient:      client,
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	require.NoError(t, err)

	return NewClient(vaultUrl, cred, options)
}

func TestSetSecret(t *testing.T) {
	recording.StartRecording(t, pathToPackage, nil)
	defer recording.StopRecording(t, nil)

	vaultUrl, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		t.Fatal("Could not find environment variable AZURE_KEYVAULT_URL")
	}
	fmt.Println(vaultUrl)
	client, err := createClient(t, vaultUrl)
	require.NoError(t, err)

	_, err = client.SetSecret(context.Background(), "mySecret", "mySecretValue", nil)
	require.NoError(t, err)
}
