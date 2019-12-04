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
	tenantID                 = "expected_tenant"
	clientID                 = "expected_client"
	secret                   = "secret"
	wrongSecret              = "wrong_secret"
	tokenValue               = "new_token"
	tokenExpiresIn           = "3600"
	scope                    = "http://storage.azure.com/.default"
	defaultTestAuthorityHost = "login.microsoftonline.com"
)

func TestClientSecretCredential_CreateAuthRequestSuccess(t *testing.T) {
	cred := NewClientSecretCredential(tenantID, clientID, secret, nil)
	req, err := cred.client.createClientSecretAuthRequest(cred.tenantID, cred.clientID, cred.clientSecret, []string{scope})
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
	if reqQueryParams[qpClientSecret][0] != secret {
		t.Fatalf("Unexpected secret in the client_secret header")
	}
	if reqQueryParams[qpScope][0] != url.QueryEscape(scope) {
		t.Fatalf("Unexpected scope in scope header")
	}
	if req.Request.URL.Host != defaultTestAuthorityHost {
		t.Fatalf("Unexpected default authority host")
	}
	if req.Request.URL.Scheme != "https" {
		t.Fatalf("Wrong request scheme")
	}
}

func TestClientSecretCredential_GetTokenSuccess(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(`{"access_token": "ey0....", "expires_in": 3600}`)))
	srvURL := srv.URL()
	cred := NewClientSecretCredential(tenantID, clientID, secret, &TokenCredentialOptions{PipelineOptions: azcore.PipelineOptions{HTTPClient: srv}, AuthorityHost: &srvURL})
	_, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestClientSecretCredential_GetTokenInvalidCredentials(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srvURL := srv.URL()
	cred := NewClientSecretCredential(tenantID, clientID, wrongSecret, &TokenCredentialOptions{PipelineOptions: azcore.PipelineOptions{HTTPClient: srv}, AuthorityHost: &srvURL})
	_, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{scope}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
}
