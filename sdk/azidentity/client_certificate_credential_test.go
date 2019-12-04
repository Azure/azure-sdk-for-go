package azidentity

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/mock"
)

const (
	certificatePath      = "certificate.pem"
	wrongCertificatePath = "wrong_certificate_path.pem"
)

func TestClientCertificateCredential_CreateAuthRequestSuccess(t *testing.T) {
	cred, err := NewClientCertificateCredential(tenantID, clientID, certificatePath, nil)
	if err != nil {
		t.Fatalf("Failed to instantiate credential")
	}
	req, err := cred.client.createClientCertificateAuthRequest(cred.tenantID, cred.clientID, cred.clientCertificate, []string{scope})
	if err != nil {
		t.Fatalf("Unexpectedly received an error: %v", err)
	}
	if req.Request.Header.Get(azcore.HeaderContentType) != azcore.HeaderURLEncoded {
		t.Fatalf("Unexpected value for Content-Type header")
	}
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		t.Fatalf("Unable to read request body")
	}
	bodyStr := string(body)
	reqQueryParams, err := url.ParseQuery(bodyStr)
	if err != nil {
		t.Fatalf("Unable to parse query params in request")
	}
	if reqQueryParams[qpClientID][0] != clientID {
		t.Fatalf("Unexpected client ID in the client_id header")
	}
	if reqQueryParams[qpGrantType][0] != "client_credentials" {
		t.Fatalf("Wrong grant type in request body")
	}
	if reqQueryParams[qpClientAssertionType][0] != clientAssertionType {
		t.Fatalf("Wrong client assertion type assigned to request")
	}
	if reqQueryParams[qpScope][0] != scope {
		t.Fatalf("Unexpected scope in scope header")
	}
	if len(reqQueryParams[qpClientAssertion][0]) == 0 {
		t.Fatalf("Client assertion is not present on the request")
	}
	if req.Request.URL.Host != defaultTestAuthorityHost {
		t.Fatalf("Unexpected default authority host")
	}
	if req.Request.URL.Scheme != "https" {
		t.Fatalf("Wrong request scheme")
	}
}

func TestClientCertificateCredential_GetTokenSuccess(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
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
}

func TestClientCertificateCredential_GetTokenInvalidCredentials(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srvURL := srv.URL()
	_, err := NewClientCertificateCredential(tenantID, clientID, wrongCertificatePath, &TokenCredentialOptions{PipelineOptions: azcore.PipelineOptions{HTTPClient: srv}, AuthorityHost: &srvURL})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}
