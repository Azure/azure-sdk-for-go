package azidentity

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	certificatePath      = "testPEM.pem"
	wrongCertificatePath = "wrong_certificate_path.pem"
)

func Test_CertificateCreateAuthRequest_Pass(t *testing.T) {
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, nil)
	if err != nil {
		t.Fatalf("Failed to instantiate credential")
	}
	req, err := cred.client.createClientCertificateAuthRequest(cred.tenantID, cred.clientID, cred.clientCertificate, []string{scope})
	if err != nil {
		t.Fatalf("Unexpectedly received an error: %w", err)
	}

	if req.Request.Header.Get(azcore.HeaderContentType) != azcore.HeaderURLEncoded {
		t.Fatalf("Unexpected value for Content-Type header")
	}

	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		t.Fatalf("Unable to read request body")
	}
	bodyStr := string(body)
	splitStr := strings.SplitN(bodyStr, "&", -1)
	reqQueryParams := make(map[string]string, len(splitStr))
	for _, i := range splitStr {
		f := strings.SplitN(i, "=", -1)
		reqQueryParams[f[0]] = f[1]
	}

	if reqQueryParams[qpClientID] != clientID {
		t.Fatalf("Unexpected client ID in the client_id header")
	}
	// TODO: add certificate specific tests here

	if reqQueryParams[qpScope] != url.QueryEscape(scope) {
		t.Fatalf("Unexpected scope in scope header")
	}

	if req.Request.URL.Host != defaultTestAuthorityHost {
		t.Fatalf("Unexpected default authority host")
	}

	if req.Request.URL.Scheme != "https" {
		t.Fatalf("Wrong request scheme")
	}
}

func Test_CertificateGetToken_Success(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithBody([]byte(`{"access_token": "new_token", "expires_in": 3600}`)))
	srvURL := srv.URL()
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, &TokenCredentialOptions{PipelineOptions: azcore.PipelineOptions{HTTPClient: srv}, AuthorityHost: &srvURL})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	close()
}

func Test_CertificateGetToken_InvalidCredentials(t *testing.T) {
	srv, close := mock.NewServer()
	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srvURL := srv.URL()
	_, err := NewClientCertificateCredential(tenantID, clientID, wrongCertificatePath, &TokenCredentialOptions{PipelineOptions: azcore.PipelineOptions{HTTPClient: srv}, AuthorityHost: &srvURL})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
	close()
}
