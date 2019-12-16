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

func newDeviceCodeInfo(dc DeviceCodeResult) DeviceCodeInfo {
	return DeviceCodeInfo{UserCode: dc.UserCode, DeviceCode: dc.DeviceCode, VerificationURL: dc.VerificationURL,
		ExpiresOn: dc.ExpiresOn, Interval: dc.Interval, Message: dc.Message, ClientID: dc.ClientID, Scopes: dc.Scopes}
}

// DeviceCodeInfo details of the device code to present to a user to allow them to authenticate through the device code authentication flow.
type DeviceCodeInfo struct {
	// JMR: Make all these private and add public getter methods?
	UserCode        string        // User code returned by the service
	DeviceCode      string        // Device code returned by the service
	VerificationURL string        // Verification URL where the user must navigate to authenticate using the device code and credentials. JMR: URL?
	ExpiresOn       time.Duration // Time when the device code will expire.
	Interval        int64         // Polling interval time to check for completion of authentication flow.
	Message         string        // User friendly text response that can be used for display purpose.
	ClientID        string        // Identifier of the client requesting device code.
	Scopes          []string      // List of the scopes that would be held by token. JMR: Should be readonly
}

// DeviceCodeCredential ...
type DeviceCodeCredential struct {
	client   *aadIdentityClient
	tenantID string // Gets the Azure Active Directory tenant (directory) Id of the service principal
	clientID string // Gets the client (application) ID of the service principal
}

// NewDeviceCodeCredential constructs a new ClientSecretCredential with the details needed to authenticate against Azure Active Directory with a client secret.
// tenantID: The Azure Active Directory tenant (directory) Id of the service principal.
// clientID: The client (application) ID of the service principal.
// options: allow to configure the management of the requests sent to the Azure Active Directory service.
func NewDeviceCodeCredential(tenantID string, clientID string, options *TokenCredentialOptions, reqOpts azcore.TokenRequestOptions) *DeviceCodeCredential {
	return &DeviceCodeCredential{tenantID: tenantID, clientID: clientID, client: newAADIdentityClient(options)}
}

// AuthenticateDeviceCode creates a device code authentication request and returns an Access Token or
// an error in case of authentication failure.
// ctx: The current request context
// tenantID: The Azure Active Directory tenant (directory) Id of the service principal
// clientID: The client (application) ID of the service principal
// username: User's account username
// password: User's account password
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
		return nil, &AuthenticationFailedError{Err: errors.New("Something unexpected happened with the request and received a nil response")}
	}

	if resp.HasStatusCode(successStatusCodes[:]...) {
		return c.createAccessToken(resp)
	}

	return nil, &AuthenticationFailedError{Err: newAuthenticationResponseError(resp)}
}

func (c *aadIdentityClient) createDeviceCodeAuthRequest(tenantID string, clientID string, deviceCode string, scopes []string) (*azcore.Request, error) {
	urlStr := c.options.AuthorityHost.String() + tenantID + tokenEndpoint
	urlFormat, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set(qpGrantType, "urn:ietf:params:oauth:grant-type:device_code")
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

// GetToken ...
func (c *DeviceCodeCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	dc, err := RequestNewDeviceCode(ctx, c.tenantID, c.clientID, opts.Scopes)
	if err != nil {
		return nil, err
	}
	fmt.Println(dc)
	// TODO / Problem need to return message here
	for {
		tk, err := c.client.authenticateDeviceCode(ctx, c.tenantID, c.clientID, dc.DeviceCode, opts.Scopes) // Problem 1: we would need scopes to be passed into the constructor
		if err == nil {
			return tk, err // Problem 2: we would need to have a field where we store the token
		}
		var authFailed *AuthenticationResponseError
		if !errors.As(err, &authFailed) {
			return nil, err
		} else {
			switch authFailed.Message {
			case "authorization_pending":
				continue
			case "authorization_declined":
				return nil, err
			case "expired_token":
				return nil, err
			case "bad_verification_code":
				return nil, err
			default:
				// Any other error should be returned
				return nil, err
			}
		}
		time.Sleep(5 * time.Second)
	}
	return nil, err

}

func RequestNewDeviceCode(ctx context.Context, tenantID, clientID string, scopes []string) (*DeviceCodeResult, error) {
	if len(tenantID) == 0 { // if the user did not pass in a tenantID then the default value is set
		tenantID = "organizations"
	}
	p := newDefaultPipeline(TokenCredentialOptions{})                     // create a default pipeline since this request is carried out before credential instantiation
	urlStr := defaultAuthorityHost + tenantID + "/oauth2/v2.0/devicecode" // endpoint that will return a device code along with the other parameters in the DeviceCodeResult struct
	urlFormat, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set(qpClientID, clientID)
	data.Set(qpScope, strings.Join(scopes, " "))
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))
	msg := p.NewRequest(http.MethodPost, *urlFormat)
	msg.Header.Set(azcore.HeaderContentType, azcore.HeaderURLEncoded)
	err = msg.SetBody(body)
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
		return createDeviceCodeAccessToken(resp)
	}
	return nil, &AuthenticationFailedError{Err: newAuthenticationResponseError(resp)}
}

func createDeviceCodeAccessToken(res *azcore.Response) (*DeviceCodeResult, error) {
	value := &DeviceCodeResult{}
	if err := json.Unmarshal(res.Payload, &value); err != nil {
		return nil, fmt.Errorf("internalAccessToken: %w", err)
	}
	return value, nil
}
