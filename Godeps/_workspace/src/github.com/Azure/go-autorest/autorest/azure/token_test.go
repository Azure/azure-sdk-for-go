package azure

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest/mocks"
)

const (
	defaultFormData = "client_id=id&client_secret=secret&grant_type=client_credentials&resource=resource"
)

func TestTokenExpires(t *testing.T) {
	tt := time.Now().Add(5 * time.Second)
	tk := newTokenExpiresAt(tt)

	if tk.Expires().Equal(tt) {
		t.Errorf("azure: Token#Expires miscalculated expiration time -- received %v, expected %v", tk.Expires(), tt)
	}
}

func TestTokenIsExpired(t *testing.T) {
	tk := newTokenExpiresAt(time.Now().Add(-5 * time.Second))

	if !tk.IsExpired() {
		t.Errorf("azure: Token#IsExpired failed to mark a stale token as expired -- now %v, token expires at %v",
			time.Now().UTC(), tk.Expires())
	}
}

func TestTokenIsExpiredUninitialized(t *testing.T) {
	tk := &Token{}

	if !tk.IsExpired() {
		t.Errorf("azure: An uninitialized Token failed to mark itself as expired (expiration time %v)", tk.Expires())
	}
}

func TestTokenIsNoExpired(t *testing.T) {
	tk := newTokenExpiresAt(time.Now().Add(1000 * time.Second))

	if tk.IsExpired() {
		t.Errorf("azure: Token marked a fresh token as expired -- now %v, token expires at %v", time.Now().UTC(), tk.Expires())
	}
}

func TestTokenWillExpireIn(t *testing.T) {
	d := 5 * time.Second
	tk := newTokenExpiresIn(d)

	if !tk.WillExpireIn(d) {
		t.Error("azure: Token#WillExpireIn mismeasured expiration time")
	}
}

func TestTokenWithAuthorization(t *testing.T) {
	tk := newToken()

	req, err := autorest.Prepare(&http.Request{}, tk.WithAuthorization())
	if err != nil {
		t.Errorf("azure: Token#WithAuthorization returned an error (%v)", err)
	} else if req.Header.Get(http.CanonicalHeaderKey("Authorization")) != fmt.Sprintf("Bearer %s", tk.AccessToken) {
		t.Error("azure: Token#WithAuthorization failed to set Authorization header")
	}
}

func TestServicePrincipalTokenSetAutoRefresh(t *testing.T) {
	spt := newServicePrincipalToken()

	if !spt.autoRefresh {
		t.Error("azure: ServicePrincipalToken did not default to automatic token refreshing")
	}

	spt.SetAutoRefresh(false)
	if spt.autoRefresh {
		t.Error("azure: ServicePrincipalToken#SetAutoRefresh did not disable automatic token refreshing")
	}
}

func TestServicePrincipalTokenSetRefreshWithin(t *testing.T) {
	spt := newServicePrincipalToken()

	if spt.refreshWithin != defaultRefresh {
		t.Error("azure: ServicePrincipalToken did not correctly set the default refresh interval")
	}

	spt.SetRefreshWithin(2 * defaultRefresh)
	if spt.refreshWithin != 2*defaultRefresh {
		t.Error("azure: ServicePrincipalToken#SetRefreshWithin did not set the refresh interval")
	}
}

func TestServicePrincipalTokenSetSender(t *testing.T) {
	spt := newServicePrincipalToken()

	var s autorest.Sender
	s = mocks.NewSender()
	spt.SetSender(s)
	if !reflect.DeepEqual(s, spt.sender) {
		t.Error("azure: ServicePrincipalToken#SetSender did not set the sender")
	}
}

func TestServicePrincipalTokenRefreshUsesPOST(t *testing.T) {
	spt := newServicePrincipalToken()

	c := mocks.NewSender()
	s := autorest.DecorateSender(c,
		(func() autorest.SendDecorator {
			return func(s autorest.Sender) autorest.Sender {
				return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
					if r.Method != "POST" {
						t.Errorf("azure: ServicePrincipalToken#Refresh did not correctly set HTTP method -- expected %v, received %v", "POST", r.Method)
					}
					return mocks.NewResponse(), nil
				})
			}
		})())
	spt.SetSender(s)
	spt.Refresh()
}

func TestServicePrincipalTokenRefreshSetsMimeType(t *testing.T) {
	spt := newServicePrincipalToken()

	c := mocks.NewSender()
	s := autorest.DecorateSender(c,
		(func() autorest.SendDecorator {
			return func(s autorest.Sender) autorest.Sender {
				return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
					if r.Header.Get(http.CanonicalHeaderKey("Content-Type")) != "application/x-www-form-urlencoded" {
						t.Errorf("azure: ServicePrincipalToken#Refresh did not correctly set Content-Type -- expected %v, received %v",
							"application/x-form-urlencoded",
							r.Header.Get(http.CanonicalHeaderKey("Content-Type")))
					}
					return mocks.NewResponse(), nil
				})
			}
		})())
	spt.SetSender(s)
	spt.Refresh()
}

func TestServicePrincipalTokenRefreshSetsURL(t *testing.T) {
	spt := newServicePrincipalToken()

	c := mocks.NewSender()
	s := autorest.DecorateSender(c,
		(func() autorest.SendDecorator {
			return func(s autorest.Sender) autorest.Sender {
				return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
					u := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token?api-version=1.0", spt.tenantID)
					if r.URL.String() != u {
						t.Errorf("azure: ServicePrincipalToken#Refresh did not correctly set the URL -- expected %v, received %v",
							u, r.URL)
					}
					return mocks.NewResponse(), nil
				})
			}
		})())
	spt.SetSender(s)
	spt.Refresh()
}

func TestServicePrincipalTokenRefreshSetsBody(t *testing.T) {
	spt := newServicePrincipalToken()

	c := mocks.NewSender()
	s := autorest.DecorateSender(c,
		(func() autorest.SendDecorator {
			return func(s autorest.Sender) autorest.Sender {
				return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
					b, err := ioutil.ReadAll(r.Body)
					if err != nil {
						t.Errorf("azure: Failed to read body of Service Principal token request (%v)", err)
					} else if string(b) != defaultFormData {
						t.Errorf("azure: ServicePrincipalToken#Refresh did not correctly set the HTTP Request Body -- expected %v, received %v",
							defaultFormData, string(b))
					}
					return mocks.NewResponse(), nil
				})
			}
		})())
	spt.SetSender(s)
	spt.Refresh()
}

func TestServicePrincipalTokenRefreshClosesRequestBody(t *testing.T) {
	spt := newServicePrincipalToken()

	resp := mocks.NewResponse()
	c := mocks.NewSender()
	s := autorest.DecorateSender(c,
		(func() autorest.SendDecorator {
			return func(s autorest.Sender) autorest.Sender {
				return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
					return resp, nil
				})
			}
		})())
	spt.SetSender(s)
	spt.Refresh()

	if resp.Body.(*mocks.Body).IsOpen() {
		t.Error("azure: ServicePrincipalToken#Refresh failed to close the HTTP Response Body")
	}
}

func TestServicePrincipalTokenRefreshPropagatesErrors(t *testing.T) {
	spt := newServicePrincipalToken()

	c := mocks.NewSender()
	c.EmitErrors(1)
	spt.SetSender(c)

	err := spt.Refresh()
	if err == nil {
		t.Error("azure: Failed to propagate the request error")
	}
}

func TestServicePrincipalTokenRefreshReturnsErrorIfNotOk(t *testing.T) {
	spt := newServicePrincipalToken()

	c := mocks.NewSender()
	c.EmitStatus("401 NotAuthorized", 401)
	spt.SetSender(c)

	err := spt.Refresh()
	if err == nil {
		t.Error("azure: Failed to return an when receiving a status code other than HTTP 200")
	}
}

func TestServicePrincipalTokenRefreshUnmarshals(t *testing.T) {
	spt := newServicePrincipalToken()

	expiresOn := strconv.Itoa(int(time.Now().Add(3600 * time.Second).Sub(expirationBase).Seconds()))
	j := newTokenJSON(expiresOn, "resource")
	resp := mocks.NewResponseWithContent(j)
	c := mocks.NewSender()
	s := autorest.DecorateSender(c,
		(func() autorest.SendDecorator {
			return func(s autorest.Sender) autorest.Sender {
				return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
					return resp, nil
				})
			}
		})())
	spt.SetSender(s)

	err := spt.Refresh()
	if err != nil {
		t.Errorf("azure: ServicePrincipalToken#Refresh returned an unexpected error (%v)", err)
	} else if spt.AccessToken != "accessToken" ||
		spt.ExpiresIn != "3600" ||
		spt.ExpiresOn != expiresOn ||
		spt.NotBefore != expiresOn ||
		spt.Resource != "resource" ||
		spt.Type != "Bearer" {
		t.Errorf("azure: ServicePrincipalToken#Refresh failed correctly unmarshal the JSON -- expected %v, received %v",
			j, *spt)
	}
}

func TestServicePrincipalTokenEnsureFreshRefreshes(t *testing.T) {
	spt := newServicePrincipalToken()
	expireToken(&spt.Token)

	f := false
	c := mocks.NewSender()
	s := autorest.DecorateSender(c,
		(func() autorest.SendDecorator {
			return func(s autorest.Sender) autorest.Sender {
				return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
					f = true
					return mocks.NewResponse(), nil
				})
			}
		})())
	spt.SetSender(s)
	spt.EnsureFresh()
	if !f {
		t.Error("azure: ServicePrincipalToken#EnsureFresh failed to call Refresh for stale token")
	}
}

func TestServicePrincipalTokenEnsureFreshSkipsIfFresh(t *testing.T) {
	spt := newServicePrincipalToken()
	setTokenToExpireIn(&spt.Token, 1000*time.Second)

	f := false
	c := mocks.NewSender()
	s := autorest.DecorateSender(c,
		(func() autorest.SendDecorator {
			return func(s autorest.Sender) autorest.Sender {
				return autorest.SenderFunc(func(r *http.Request) (*http.Response, error) {
					f = true
					return mocks.NewResponse(), nil
				})
			}
		})())
	spt.SetSender(s)
	spt.EnsureFresh()
	if f {
		t.Error("azure: ServicePrincipalToken#EnsureFresh invoked Refresh for fresh token")
	}
}

func TestServicePrincipalTokenWithAuthorization(t *testing.T) {
	spt := newServicePrincipalToken()
	setTokenToExpireIn(&spt.Token, 1000*time.Second)

	req, err := autorest.Prepare(&http.Request{}, spt.WithAuthorization())
	if err != nil {
		t.Errorf("azure: ServicePrincipalToken#WithAuthorization returned an error (%v)", err)
	} else if req.Header.Get(http.CanonicalHeaderKey("Authorization")) != fmt.Sprintf("Bearer %s", spt.AccessToken) {
		t.Error("azure: ServicePrincipalToken#WithAuthorization failed to set Authorization header")
	}
}

func TestServicePrincipalTokenWithAuthorizationReturnsErrorIfCannotRefresh(t *testing.T) {
	spt := newServicePrincipalToken()

	_, err := autorest.Prepare(&http.Request{}, spt.WithAuthorization())
	if err == nil {
		t.Error("azure: ServicePrincipalToken#WithAuthorization failed to return an error when refresh fails")
	}
}

func newToken() *Token {
	return &Token{
		AccessToken: "ASECRETVALUE",
		Resource:    "https://azure.microsoft.com/",
		Type:        "Bearer",
	}
}

func newTokenJSON(expiresOn string, resource string) string {
	return fmt.Sprintf(`{
		"access_token" : "accessToken",
		"expires_in"   : "3600",
		"expires_on"   : "%s",
		"not_before"   : "%s",
		"resource"     : "%s",
		"token_type"   : "Bearer"
		}`,
		expiresOn, expiresOn, resource)
}

func newTokenExpiresIn(expireIn time.Duration) *Token {
	return setTokenToExpireIn(newToken(), expireIn)
}

func newTokenExpiresAt(expireAt time.Time) *Token {
	return setTokenToExpireAt(newToken(), expireAt)
}

func expireToken(t *Token) *Token {
	return setTokenToExpireIn(t, 0)
}

func setTokenToExpireAt(t *Token, expireAt time.Time) *Token {
	t.ExpiresIn = "3600"
	t.ExpiresOn = strconv.Itoa(int(expireAt.Sub(expirationBase).Seconds()))
	t.NotBefore = t.ExpiresOn
	return t
}

func setTokenToExpireIn(t *Token, expireIn time.Duration) *Token {
	return setTokenToExpireAt(t, time.Now().Add(expireIn))
}

func newServicePrincipalToken() *ServicePrincipalToken {
	spt, _ := NewServicePrincipalToken("id", "secret", "tenentId", "resource")
	return spt
}
