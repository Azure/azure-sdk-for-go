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
	deviceCode         = "device_code"
	deviceCodeResponse = `{"user_code":"BKG99ZAF3","device_code":"BAQABAAEAAACQN9QBRU3jT6bcBQLZNUj7vq0a6sj0GHRFY6Y9IqPV-SsbHst2QKPg3irWUAq9aIgQOwDxrKBCzV_oolXqcYxR4OSiNhCuxNQ1r7ju30vGtGxz2JZ9NkI_btVoY9rvNxSCMhV78pZILlZh7lh952iWUgUY16rxVIndEPz7_K9W-kpHl4nAktRmHc7Y-Y3xVo8gAA","verification_uri":"https://microsoft.com/devicelogin","expires_in":900,"interval":5,"message":"To sign in, use a web browser to open the page https://microsoft.com/devicelogin and enter the code BKG99ZAF3 to authenticate."}`
	deviceCodeScopes   = "user.read offline_access openid profile email"
)

func TestDeviceCodeCredential_CreateAuthRequestSuccess(t *testing.T) {
	handler := func(s string) {}
	cred := NewDeviceCodeCredential(tenantID, clientID, handler, nil)
	req, err := cred.client.createDeviceCodeAuthRequest(cred.tenantID, cred.clientID, deviceCode, []string{deviceCodeScopes})
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
	if reqQueryParams[qpGrantType][0] != deviceCodeGrantType {
		t.Fatalf("Unexpected grant type")
	}
	if reqQueryParams[qpClientID][0] != clientID {
		t.Fatalf("Unexpected client ID in the client_id header")
	}
	if reqQueryParams[qpDeviceCode][0] != deviceCode {
		t.Fatalf("Unexpected username in the username header")
	}
	if reqQueryParams[qpScope][0] != deviceCodeScopes {
		t.Fatalf("Unexpected scope in scope header")
	}
	if req.Request.URL.Host != defaultTestAuthorityHost {
		t.Fatalf("Unexpected default authority host")
	}
	if req.Request.URL.Scheme != "https" {
		t.Fatalf("Wrong request scheme")
	}
}

func TestDeviceCodeCredential_GetTokenSuccess(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(deviceCodeResponse)))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srvURL := srv.URL()
	handler := func(string) {}
	cred := NewDeviceCodeCredential(tenantID, clientID, handler, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	_, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{deviceCodeScopes}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
}

func TestDeviceCodeCredential_GetTokenInvalidCredentials(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srvURL := srv.URL()
	handler := func(string) {}
	cred := NewDeviceCodeCredential(tenantID, clientID, handler, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	_, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{deviceCodeScopes}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
}

func TestDeviceCodeCredential_CreateAuthRequestFail(t *testing.T) {
	srv, close := mock.NewServer()
	defer close()
	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
	srvURL := srv.URL()
	srvURL.Host = "ht @"
	handler := func(string) {}
	cred := NewDeviceCodeCredential(tenantID, clientID, handler, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	_, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{deviceCodeScopes}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one")
	}
}

func TestBearerPolicy_DeviceCodeCredential(t *testing.T) {
	srv, close := mock.NewTLSServer()
	defer close()
	srv.AppendResponse(mock.WithBody([]byte(deviceCodeResponse)))
	srv.AppendResponse(mock.WithBody([]byte(accessTokenRespSuccess)))
	srv.AppendResponse(mock.WithStatusCode(http.StatusOK))
	srvURL := srv.URL()
	handler := func(string) {}
	cred := NewDeviceCodeCredential(tenantID, clientID, handler, &TokenCredentialOptions{HTTPClient: srv, AuthorityHost: &srvURL})
	pipeline := azcore.NewPipeline(
		srv,
		azcore.NewTelemetryPolicy(azcore.TelemetryOptions{}),
		azcore.NewUniqueRequestIDPolicy(),
		azcore.NewRetryPolicy(azcore.RetryOptions{}),
		cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: []string{deviceCodeScopes}}}),
		azcore.NewRequestLogPolicy(azcore.RequestLogOptions{}))
	req := pipeline.NewRequest(http.MethodGet, srv.URL())
	_, err := req.Do(context.Background())
	if err != nil {
		t.Fatalf("Expected an empty error but receive: %v", err)
	}
}
