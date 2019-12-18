package azidentity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	deviceCodeGrantType = "urn:ietf:params:oauth:grant-type:device_code"
)

// DeviceCodeResult ...
type DeviceCodeResult struct {
	UserCode        string        `json:"user_code"`        // User code returned by the service
	DeviceCode      string        `json:"device_code"`      // Device code returned by the service
	VerificationURL string        `json:"verification_uri"` // Verification URL where the user must navigate to authenticate using the device code and credentials. JMR: URL?
	ExpiresOn       time.Duration // Time when the device code will expire.
	Interval        int64         `json:"interval"` // Polling interval time to check for completion of authentication flow.
	Message         string        `json:"message"`  // User friendly text response that can be used for display purpose.
	ClientID        string        // Identifier of the client requesting device code.
	Scopes          []string      // List of the scopes that would be held by token. JMR: Should be readonly
}

// DeviceCodeCredential ...
type DeviceCodeCredential struct {
	client   *aadIdentityClient
	tenantID string       // Gets the Azure Active Directory tenant (directory) Id of the service principal
	clientID string       // Gets the client (application) ID of the service principal
	callback func(string) // Sends the user a message with a verification URL and device code to sign in to the login server
}

// NewDeviceCodeCredential constructs a new DeviceCodeCredential with the details needed to authenticate against Azure Active Directory with a device code.
// tenantID: The Azure Active Directory tenant (directory) Id of the service principal.
// clientID: The client (application) ID of the service principal.
// callback: The callback function used to send the login message back to the user
// options: allow to configure the management of the requests sent to the Azure Active Directory service.
func NewDeviceCodeCredential(tenantID string, clientID string, callback func(string), options *TokenCredentialOptions) *DeviceCodeCredential {
	return &DeviceCodeCredential{tenantID: tenantID, clientID: clientID, callback: callback, client: newAADIdentityClient(options)}
}

// GetToken obtains a token from the Azure Active Directory service, following the device code authentication
// flow. This function first requests a device code and requests that the user login to continue authenticating.
// This function will keep polling the service for a token meanwhile the user logs.
// - scopes: The list of scopes for which the token will have access.
// - ctx: controlling the request lifetime.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *DeviceCodeCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	dc, err := c.client.requestNewDeviceCode(ctx, c.tenantID, c.clientID, opts.Scopes)
	if err != nil {
		return nil, err
	}
	c.callback(dc.Message)
	for {
		tk, err := c.client.authenticateDeviceCode(ctx, c.tenantID, c.clientID, dc.DeviceCode, opts.Scopes)
		if err == nil {
			return tk, err
		}
		var authFailed *AuthenticationFailedError
		if !errors.As(err, &authFailed) {
			return nil, err
		}
		var authRespErr *AuthenticationResponseError
		if !errors.As(authFailed.Unwrap(), &authRespErr) {
			return nil, err
		}
		switch authRespErr.Message {
		case "authorization_pending":
			time.Sleep(time.Duration(dc.Interval) * time.Second)
		default:
			// Any other error should be returned
			return nil, err
		}
	}
}

// AuthenticationPolicy implements the azcore.Credential interface on ClientSecretCredential.
func (c *DeviceCodeCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}

func createDeviceCodeResult(res *azcore.Response) (*DeviceCodeResult, error) {
	value := &DeviceCodeResult{}
	if err := json.Unmarshal(res.Payload, &value); err != nil {
		return nil, fmt.Errorf("DeviceCodeResult: %w", err)
	}
	return value, nil
}

func (c *aadIdentityClient) createDeviceCodeAccessToken(res *azcore.Response) (*azcore.AccessToken, error) {
	value := struct {
		// these are the only fields that we use
		Token        string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
		ExpiresIn    json.Number `json:"expires_in"`
		ExpiresOn    string      `json:"expires_on"` // the value returned in this field varies between a number and a date string
	}{}
	if err := json.Unmarshal(res.Payload, &value); err != nil {
		return nil, fmt.Errorf("internal AccessToken: %w", err)
	}
	expiresIn, err := value.ExpiresIn.Int64()
	if err != nil {
		return nil, err
	}
	c.updateRefreshToken(value.RefreshToken)
	return &azcore.AccessToken{Token: value.Token, ExpiresOn: time.Now().Add(time.Second * time.Duration(expiresIn)).UTC()}, nil
}

func (c *aadIdentityClient) updateRefreshToken(tk string) {
	// c.lock.Lock()
	// c.refreshToken = tk
	// c.lock.Unlock()
}

// authenticateDeviceCode creates a device code authentication request and returns an Access Token or
// an error in case of authentication failure.
// ctx: The current request context
// tenantID: The Azure Active Directory tenant (directory) Id of the service principal
// clientID: The client (application) ID of the service principal
// deviceCode: The device code associated with the request
// scopes: The scopes required for the token
func (c *aadIdentityClient) authenticateDeviceCode(ctx context.Context, tenantID string, clientID string, deviceCode string, scopes []string) (*azcore.AccessToken, error) {
	msg, err := c.createDeviceCodeAuthRequest(tenantID, clientID, deviceCode, scopes)
	if err != nil {
		return nil, err
	}

	resp, err := msg.Do(ctx)
	if err != nil {
		return nil, err
	}

	// This should never happen under normal conditions
	if resp == nil {
		return nil, &AuthenticationFailedError{Message: "Something unexpected happened with the request and received a nil response"}
	}

	if resp.HasStatusCode(successStatusCodes[:]...) {
		return c.createDeviceCodeAccessToken(resp)
	}

	return nil, &AuthenticationFailedError{Err: newAuthenticationResponseError(resp)}
}

func (c *aadIdentityClient) createDeviceCodeAuthRequest(tenantID string, clientID string, deviceCode string, scopes []string) (*azcore.Request, error) {
	if len(tenantID) == 0 { // if the user did not pass in a tenantID then the default value is set
		tenantID = "organizations"
	}
	urlStr := c.options.AuthorityHost.String() + tenantID + tokenEndpoint
	urlFormat, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set(qpGrantType, deviceCodeGrantType)
	data.Set(qpClientID, clientID)
	data.Set(qpDeviceCode, deviceCode)
	data.Set(qpScope, strings.Join(scopes, " "))
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))
	msg := c.pipeline.NewRequest(http.MethodPost, *urlFormat)
	msg.Header.Set(azcore.HeaderContentType, azcore.HeaderURLEncoded)
	err = msg.SetBody(body)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *aadIdentityClient) requestNewDeviceCode(ctx context.Context, tenantID, clientID string, scopes []string) (*DeviceCodeResult, error) {
	msg, err := c.createDeviceCodeNumberRequest(tenantID, clientID, scopes)
	if err != nil {
		return nil, err
	}

	resp, err := msg.Do(ctx)
	if err != nil {
		return nil, err
	}
	// This should never happen under normal conditions
	if resp == nil {
		return nil, &AuthenticationFailedError{Err: errors.New("Something unexpected happened with the request and received a nil response")}
	}
	if resp.HasStatusCode(successStatusCodes[:]...) {
		return createDeviceCodeResult(resp)
	}
	return nil, &AuthenticationFailedError{Err: newAuthenticationResponseError(resp)}
}

func (c *aadIdentityClient) createDeviceCodeNumberRequest(tenantID string, clientID string, scopes []string) (*azcore.Request, error) {
	if len(tenantID) == 0 { // if the user did not pass in a tenantID then the default value is set
		tenantID = "organizations"
	}
	urlStr := c.options.AuthorityHost.String() + tenantID + "/oauth2/v2.0/devicecode" // endpoint that will return a device code along with the other parameters in the DeviceCodeResult struct
	urlFormat, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set(qpClientID, clientID)
	data.Set(qpScope, strings.Join(scopes, " "))
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))
	msg := c.pipeline.NewRequest(http.MethodPost, *urlFormat)
	msg.Header.Set(azcore.HeaderContentType, azcore.HeaderURLEncoded)
	err = msg.SetBody(body)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
