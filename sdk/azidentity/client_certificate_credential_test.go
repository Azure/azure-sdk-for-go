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

func Test_CertificateCreateAuthRequest_Pass(t *testing.T) {
	cred, err := NewClientCertificateCredential("expected_tenant", "expected_client", "certificate_path", nil)
	if err != nil {
		t.Fatalf("Failed to instantiate credential")
	}
	req, err := cred.client.createClientCertificateAuthRequest(cred.tenantID, cred.clientID, cred.clientCertificate, []string{"http://storage.azure.com/.default"})
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

	if reqQueryParams[qpClientID] != "expected_client" {
		t.Fatalf("Unexpected client ID in the client_id header")
	}

	if reqQueryParams[qpClientSecret] != "secret" {
		t.Fatalf("Unexpected secret in the client_secret header")
	}

	if reqQueryParams[qpScope] != url.QueryEscape("http://storage.azure.com/.default") {
		t.Fatalf("Unexpected scope in scope header")
	}

	if req.Request.URL.Host != "login.microsoftonline.com" {
		t.Fatalf("Unexpected default authority host")
	}

	if req.Request.URL.Scheme != "https" {
		t.Fatalf("Wrong request scheme")
	}
}

func Test_CertificateGetToken_Success(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithBody(`{"access_token": "ey0....", "expires_in": 3600}`))
	tempURL := srv.URL()
	testURL, err := url.Parse(tempURL.String() + "/")
	if err != nil {
		t.Fatalf("Unable to parse url")
	}
	cred := NewClientSecretCredential("expected_tenant", "expected_client", "expected_secret", &IdentityClientOptions{PipelineOptions: azcore.PipelineOptions{HTTPClient: srv}, AuthorityHost: testURL})
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{"https://storage.azure.com/.default"}})
	if err != nil {
		t.Fatalf("Expected an empty error but received: %s", err.Error())
	}
	close()
}

func Test_CertificateGetToken_NilScope(t *testing.T) {
	srv, close := mock.NewServer()
	srv.AppendResponse(mock.WithBody(`{"access_token": "ey0....", "expires_in": 3600}`))
	tempURL := srv.URL()
	testURL, err := url.Parse(tempURL.String() + "/")
	if err != nil {
		t.Fatalf("Unable to parse url")
	}
	cred := NewClientSecretCredential("expected_tenant", "expected_client", "expected_secret", &IdentityClientOptions{PipelineOptions: azcore.PipelineOptions{HTTPClient: srv}, AuthorityHost: testURL})
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
	close()
}

func Test_CertificateGetToken_InvalidCredentials(t *testing.T) {
	srv, close := mock.NewServer()
	srv.SetResponse(mock.WithStatusCode(http.StatusUnauthorized))
	tempURL := srv.URL()
	testURL, err := url.Parse(tempURL.String() + "/")
	if err != nil {
		t.Fatalf("Unable to parse url")
	}
	cred := NewClientSecretCredential("expected_tenant", "expected_client", "wrong_secret", &IdentityClientOptions{PipelineOptions: azcore.PipelineOptions{HTTPClient: srv}, AuthorityHost: testURL})
	_, err = cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{"https://storage.azure.com/.default"}})
	if err == nil {
		t.Fatalf("Expected an error but did not receive one.")
	}
	close()
}
