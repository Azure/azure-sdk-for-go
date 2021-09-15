//go:build go1.16
// +build go1.16

package azsecrets

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
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

func createClient(t *testing.T) (*Client, error) {
	vaultUrl, ok := os.LookupEnv("AZURE_KEYVAULT_URL")
	if !ok {
		t.Fatal("Could not find environment variable AZURE_KEYVAULT_URL")
	}

	p := NewRecordingPolicy(t, &recording.RecordingOptions{UseHTTPS: true})
	client, err := recording.GetHTTPClient(t)
	require.NoError(t, err)

	options := &ClientOptions{
		PerCallPolicies: []policy.Policy{p},
		HTTPClient:      client,
	}
	_ = options

	cred, err := azidentity.NewClientSecretCredential(
		os.Getenv("KEYVAULT_TENANT_ID"),
		os.Getenv("KEYVAULT_CLIENT_ID"),
		os.Getenv("KEYVAULT_CLIENT_SECRET"),
		nil,
	)
	require.NoError(t, err)

	return NewClient(vaultUrl, cred, options)
}

func TestSetGetSecret(t *testing.T) {
	recording.StartRecording(t, pathToPackage, nil)
	defer recording.StopRecording(t, nil)

	client, err := createClient(t)
	require.NoError(t, err)

	secretValue := "mySecretValue"
	resp, err := client.SetSecret(context.Background(), "mySecret", secretValue, nil)
	require.NoError(t, err)
	require.Equal(t, *resp.Value, secretValue)

	secretVersion := strings.Split(*resp.ID, "/")

	getResp, err := client.GetSecret(context.Background(), "mySecret", secretVersion[len(secretVersion)-1], nil)
	require.NoError(t, err)
	require.Equal(t, *getResp.Value, secretValue)
}
