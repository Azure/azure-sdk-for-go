package azure

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/Godeps/_workspace/src/github.com/Azure/go-autorest/autorest"
)

const (
	defaultRefresh = 5 * time.Minute
	oauthURL       = "https://login.microsoftonline.com/{tenantID}/oauth2/{requestType}?api-version=1.0"
	tokenBaseDate  = "1970-01-01T00:00:00Z"

	// AzureResourceManagerScope is the OAuth scope for the Azure Resource Manager.
	AzureResourceManagerScope = "https://management.azure.com/"
)

var expirationBase time.Time

func init() {
	expirationBase, _ = time.Parse(time.RFC3339, tokenBaseDate)
}

// Token encapsulates the access token used to authorize Azure requests.
type Token struct {
	AccessToken string `json:"access_token"`

	ExpiresIn string `json:"expires_in"`
	ExpiresOn string `json:"expires_on"`
	NotBefore string `json:"not_before"`

	Resource string `json:"resource"`
	Type     string `json:"token_type"`
}

// Expires returns the time.Time when the Token expires.
func (t Token) Expires() time.Time {
	s, err := strconv.Atoi(t.ExpiresOn)
	if err != nil {
		s = -3600
	}
	return expirationBase.Add(time.Duration(s) * time.Second).UTC()
}

// IsExpired returns true if the Token is expired, false otherwise.
func (t Token) IsExpired() bool {
	return t.WillExpireIn(0)
}

// WillExpireIn returns true if the Token will expire after the passed time.Duration interval
// from now, false otherwise.
func (t Token) WillExpireIn(d time.Duration) bool {
	return !t.Expires().After(time.Now().Add(d))
}

// WithAuthorization returns a PrepareDecorator that adds an HTTP Authorization header whose
// value is "Bearer " followed by the AccessToken of the Token.
func (t *Token) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			return (autorest.WithBearerAuthorization(t.AccessToken)(p)).Prepare(r)
		})
	}
}

// ServicePrincipalToken encapsulates a Token created for a Service Principal.
type ServicePrincipalToken struct {
	Token

	clientID      string
	clientSecret  string
	resource      string
	tenantID      string
	autoRefresh   bool
	refreshWithin time.Duration
	sender        autorest.Sender
}

// NewServicePrincipalToken creates a ServicePrincipalToken from the supplied Service Principal
// credentials scoped to the named resource.
func NewServicePrincipalToken(id string, secret string, tenantID string, resource string) (*ServicePrincipalToken, error) {
	spt := &ServicePrincipalToken{
		clientID:      id,
		clientSecret:  secret,
		resource:      resource,
		tenantID:      tenantID,
		autoRefresh:   true,
		refreshWithin: defaultRefresh,
		sender:        &http.Client{}}
	return spt, nil
}

// EnsureFresh will refresh the token if it will expire within the refresh window (as set by
// RefreshWithin).
func (spt *ServicePrincipalToken) EnsureFresh() error {
	if spt.WillExpireIn(spt.refreshWithin) {
		return spt.Refresh()
	}
	return nil
}

// Refresh obtains a fresh token for the Service Principal.
func (spt *ServicePrincipalToken) Refresh() error {
	p := map[string]interface{}{
		"tenantID":    spt.tenantID,
		"requestType": "token",
	}

	v := url.Values{}
	v.Set("client_id", spt.clientID)
	v.Set("client_secret", spt.clientSecret)
	v.Set("grant_type", "client_credentials")
	v.Set("resource", spt.resource)

	req, _ := autorest.Prepare(&http.Request{},
		autorest.AsPost(),
		autorest.AsFormURLEncoded(),
		autorest.WithBaseURL(oauthURL),
		autorest.WithPathParameters(p),
		autorest.WithFormData(v))

	resp, err := autorest.SendWithSender(spt.sender, req)
	if err != nil {
		return autorest.NewErrorWithError(err,
			"azure.ServicePrincipalToken", "Refresh", "Failure sending request for Service Principal %s",
			spt.clientID)
	}

	err = autorest.Respond(resp,
		autorest.WithErrorUnlessOK(),
		autorest.ByUnmarshallingJSON(spt),
		autorest.ByClosing())
	if err != nil {
		return autorest.NewErrorWithError(err,
			"azure.ServicePrincipalToken", "Refresh", "Failure handling response to Service Principal %s request",
			spt.clientID)
	}

	return nil
}

// SetAutoRefresh enables or disables automatic refreshing of stale tokens.
func (spt *ServicePrincipalToken) SetAutoRefresh(autoRefresh bool) {
	spt.autoRefresh = autoRefresh
}

// SetRefreshWithin sets the interval within which if the token will expire, EnsureFresh will
// refresh the token.
func (spt *ServicePrincipalToken) SetRefreshWithin(d time.Duration) {
	spt.refreshWithin = d
	return
}

// SetSender sets the autorest.Sender used when obtaining the Service Principal token. An
// undecorated http.Client is used by default.
func (spt *ServicePrincipalToken) SetSender(s autorest.Sender) {
	spt.sender = s
}

// WithAuthorization returns a PrepareDecorator that adds an HTTP Authorization header whose
// value is "Bearer " followed by the AccessToken of the ServicePrincipalToken.
//
// By default, the token will automatically refresh if nearly expired (as determined by the
// RefreshWithin interval). Use the AutoRefresh method to enable or disable automatically refreshing
// tokens.
func (spt *ServicePrincipalToken) WithAuthorization() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			if spt.autoRefresh {
				err := spt.EnsureFresh()
				if err != nil {
					return r, autorest.NewErrorWithError(err,
						"azure.ServicePrincipalToken", "WithAuthorization", "Failed to refresh Service Principal Token for request to %s",
						r.URL)
				}
			}
			return (autorest.WithBearerAuthorization(spt.AccessToken)(p)).Prepare(r)
		})
	}
}
