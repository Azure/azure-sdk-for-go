package armcompute_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
	"hash/fnv"
	"net/http"
	"os"
	"testing"
	"time"
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

func authenticateTest(t *testing.T) (azcore.TokenCredential, *arm.ConnectionOptions) {
	p := NewRecordingPolicy(t, &recording.RecordingOptions{UseHTTPS: true})
	client, err := recording.GetHTTPClient(t)
	require.NoError(t, err)

	options := &arm.ConnectionOptions{
		PerCallPolicies: []policy.Policy{p},
		HTTPClient:       client,
	}

	var cred azcore.TokenCredential
	if recording.GetRecordMode() != "playback" {
		tenantId := lookupEnvVar("AZURE_TENANT_ID")
		clientId := lookupEnvVar("AZURE_CLIENT_ID")
		clientSecret := lookupEnvVar("AZURE_CLIENT_SECRET")
		cred, err = azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
		require.NoError(t, err)
	} else {
		cred = &fakeCredential{}
	}

	return cred, options
}

func lookupEnvVar(s string) string {
	ret, ok := os.LookupEnv(s)
	if !ok {
		panic(fmt.Sprintf("Could not find env var: '%s'", s))
	}
	return ret
}

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
	if recording.GetRecordMode() != "live" {
		originalURLHost := req.Raw().URL.Host
		req.Raw().URL.Scheme = "https"
		req.Raw().URL.Host = p.options.Host
		req.Raw().Host = p.options.Host

		req.Raw().Header.Set(recording.UpstreamUriHeader, fmt.Sprintf("%v://%v", p.options.Scheme, originalURLHost))
		req.Raw().Header.Set(recording.ModeHeader, recording.GetRecordMode())
		req.Raw().Header.Set(recording.IdHeader, recording.GetRecordingId(p.t))
	}
	return req.Next()
}

func createRandomName(t *testing.T, prefix string) (string, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(t.Name()))
	return prefix + fmt.Sprint(h.Sum32()), err
}

type fakeCredential struct {
}

func (f *fakeCredential) Do(req *policy.Request) (*http.Response, error) {
	return req.Next()
}

func (f *fakeCredential) NewAuthenticationPolicy(options runtime.AuthenticationOptions) policy.Policy {
	return &fakeCredential{}
}

func (f *fakeCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{
		Token:     "faketoken",
		ExpiresOn: time.Date(2040, time.January, 1, 1, 1, 1, 1, time.UTC),
	}, nil
}
