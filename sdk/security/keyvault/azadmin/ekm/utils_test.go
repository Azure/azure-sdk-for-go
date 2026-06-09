// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package ekm_test

import (
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	azcred "github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azadmin/ekm"
	"github.com/stretchr/testify/require"
)

const recordingDirectory = "sdk/security/keyvault/azadmin/testdata"

var (
	credential   azcore.TokenCredential
	hsmURL       string
	ekmProxyHost string
	ekmCACert    []byte

	fakeHsmURL       = fmt.Sprintf("https://%s.managedhsm.azure.net/", recording.SanitizedValue)
	fakeEkmProxyHost = fmt.Sprintf("%s.ekm.example.com:443", recording.SanitizedValue)
	// fakeCACert is the byte sequence used in place of a real CA certificate when
	// no EKM_SERVER_CA_CERTIFICATE is supplied. The Managed HSM accepts these
	// bytes verbatim on create/update; only CheckEkmConnection / GetEkmCertificate
	// parse them, and those tests skip gracefully when the proxy isn't reachable.
	fakeCACert = []byte("azure-sdk-for-go-ekm-test-ca-certificate-bytes")
)

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	if recording.GetRecordMode() == recording.PlaybackMode || recording.GetRecordMode() == recording.RecordingMode {
		proxy, err := recording.StartTestProxy(recordingDirectory, nil)
		if err != nil {
			panic(err)
		}

		defer func() {
			err := recording.StopTestProxy(proxy)
			if err != nil {
				panic(err)
			}
		}()
	}

	var err error
	credential, err = azcred.New(nil)
	if err != nil {
		panic(err)
	}

	hsmURL = recording.GetEnvVariable("AZURE_MANAGEDHSM_URL", fakeHsmURL)
	// EKM_PROXY_HOST should point at an EKM proxy reachable from the Managed HSM
	// (for example "ekm.contoso.com:443"). In playback it is replaced by the
	// sanitized placeholder. The value is normalized to always include an explicit
	// port: the host sanitizer regex matches "<host>:<port>", so a bare hostname
	// would also match — and overwrite — the bare CN value in the same body.
	ekmProxyHost = recording.GetEnvVariable("EKM_PROXY_HOST", fakeEkmProxyHost)
	if _, _, err := net.SplitHostPort(ekmProxyHost); err != nil {
		ekmProxyHost += ":443"
	}
	// EKM_SERVER_CA_CERTIFICATE is the CA certificate chain that issued the EKM
	// proxy's server certificate. PEM and base64-encoded DER are both accepted.
	// In playback we fall back to fakeCACert so recordings remain stable.
	ekmCACert = mustDecodeCACertEnv("EKM_SERVER_CA_CERTIFICATE", fakeCACert)

	if recording.GetRecordMode() == recording.RecordingMode {
		if err := recording.AddGeneralRegexSanitizer(fakeHsmURL, hsmURL, nil); err != nil {
			panic(err)
		}
		if err := recording.AddBodyRegexSanitizer(fakeEkmProxyHost, regexp.QuoteMeta(ekmProxyHost), nil); err != nil {
			panic(err)
		}
		// The Subject Common Name on the proxy's server cert must match the
		// host portion of the dial target (without the port), so it also
		// reveals the real hostname and needs its own sanitizer entry.
		realCN := serverSubjectCommonNameFor(ekmProxyHost)
		fakeCN := serverSubjectCommonNameFor(fakeEkmProxyHost)
		if realCN != fakeCN {
			if err := recording.AddBodyRegexSanitizer(fakeCN, regexp.QuoteMeta(realCN), nil); err != nil {
				panic(err)
			}
		}
		// The CA certificate appears base64-encoded in request bodies under the
		// server_ca_certificates JSON field. Sanitize its on-wire base64 form so
		// real certificate material never lands in recordings. The replacement
		// must match what playback sends (base64(fakeCACert)), not the generic
		// "Sanitized" placeholder, because the test-proxy matches request bodies
		// byte-for-byte and the playback test always sends the fake cert.
		caCertB64 := base64.StdEncoding.EncodeToString(ekmCACert)
		fakeCACertB64 := base64.StdEncoding.EncodeToString(fakeCACert)
		if caCertB64 != fakeCACertB64 {
			if err := recording.AddBodyRegexSanitizer(fakeCACertB64, regexp.QuoteMeta(caCertB64), nil); err != nil {
				panic(err)
			}
		}
	}

	return m.Run()
}

func startRecording(t *testing.T) {
	err := recording.Start(t, recordingDirectory, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
}

func startEKMTest(t *testing.T) *ekm.KeyVaultClient {
	// EKM tests require an external key manager proxy reachable from the
	// Managed HSM. CI doesn't provision one, so live runs are skipped unless
	// EKM_PROXY_HOST is supplied. Recording runs still execute (so recordings
	// can be refreshed against an environment that does have a real proxy)
	// and playback runs always execute against the recordings.
	if recording.GetRecordMode() == recording.LiveMode && os.Getenv("EKM_PROXY_HOST") == "" {
		t.Skip("skipping live EKM test: set EKM_PROXY_HOST to run against a real EKM proxy")
	}
	startRecording(t)
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	opts := &ekm.ClientOptions{ClientOptions: azcore.ClientOptions{Transport: transport}}
	client, err := ekm.NewClient(hsmURL, credential, opts)
	require.NoError(t, err)
	return client
}

type serdeModel interface {
	json.Marshaler
	json.Unmarshaler
}

func testSerde[T serdeModel](t *testing.T, model T) {
	data, err := model.MarshalJSON()
	require.NoError(t, err)
	err = model.UnmarshalJSON(data)
	require.NoError(t, err)

	// testing unmarshal error scenarios
	var data2 []byte
	err = model.UnmarshalJSON(data2)
	require.Error(t, err)

	m := regexp.MustCompile(":.*$")
	modifiedData := m.ReplaceAllString(string(data), ":false}")
	if modifiedData != "{}" {
		data3 := []byte(modifiedData)
		err = model.UnmarshalJSON(data3)
		require.Error(t, err)
	}
}

// mustDecodeCACertEnv returns the DER bytes for the certificate stored in the
// named environment variable, accepting either PEM or base64-encoded DER. If
// the variable is empty the supplied fallback is returned, which is used for
// playback so recordings stay deterministic.
func mustDecodeCACertEnv(name string, fallback []byte) []byte {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback
	}
	if block, _ := pem.Decode([]byte(raw)); block != nil && block.Type == "CERTIFICATE" {
		return block.Bytes
	}
	der, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		panic(fmt.Sprintf("%s must contain a PEM certificate or base64-encoded DER: %v", name, err))
	}
	return der
}

// serverSubjectCommonNameFor returns the host portion of an EKM proxy
// "host[:port]" string. The EKM proxy's server certificate must carry this
// value as its Subject Common Name so the Managed HSM can complete TLS
// verification when dialing the proxy.
func serverSubjectCommonNameFor(hostPort string) string {
	if host, _, err := net.SplitHostPort(hostPort); err == nil {
		return host
	}
	return hostPort
}
