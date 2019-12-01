package azidentity

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// func Test_CreateAuthRequest(t *testing.T) {
// 	srv, close := mock.NewServer()
// 	defer close()
// 	srv.SetResponse(mock.WithStatusCode(http.StatusOK))
// 	opts := azcore.PipelineOptions{HTTPClient: srv}
// 	var srvURL *url.URL
// 	*srvURL = srv.URL()
// 	cred := NewClientSecretCredential("expected_tenant", "expected_client", "secred", &IdentityClientOptions{PipelineOptions: opts, AuthorityHost: srvURL})
// 	tk, err := cred.GetToken(context.Background(), azcore.TokenRequestOptions{Scopes: []string{"https://storage.azure.com/.default"}})
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	if len(tk.Token) == 0 {
// 		t.Fatalf("Unexpected error")
// 	}
// }

func TestTemp(t *testing.T) {
	// srv, close := mock.NewServer()
	// defer close()
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Success!")
	}

	req := httptest.NewRequest("GET", "http://127.0.0.1:54814", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))

	bod := string(body)
	fmt.Println(bod)
}

func Test_CreateAuthRequest_Pass(t *testing.T) {
	cred := NewClientSecretCredential("expected_tenant", "expected_client", "secret", nil)
	req, err := cred.client.createClientSecretAuthRequest(cred.tenantID, cred.clientID, cred.clientSecret, []string{"http://storage.azure.com/.default"})
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
