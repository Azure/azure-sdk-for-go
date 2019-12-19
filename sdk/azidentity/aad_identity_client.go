// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bufio"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	clientAssertionType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
	tokenEndpoint       = "/oauth2/v2.0/token/"
)

const (
	qpClientAssertionType = "client_assertion_type"
	qpClientAssertion     = "client_assertion"
	qpClientID            = "client_id"
	qpClientSecret        = "client_secret"
	qpDeviceCode          = "device_code"
	qpGrantType           = "grant_type"
	qpPassword            = "password"
	qpResponseType        = "response_type"
	qpScope               = "scope"
	qpUsername            = "username"
)

// AADIdentityClient provides the base for authenticating with Client Secret Credentials, Client Certificate Credentials
// and Environment Credentials. This type initializes a default azcore.Pipeline and IdentityClientOptions.
type aadIdentityClient struct {
	lock         sync.RWMutex
	options      TokenCredentialOptions
	pipeline     azcore.Pipeline
	refreshToken string
}

// NewAADIdentityClient creates a new instance of the AADIdentityClient with the IdentityClientOptions
// that are passed into it along with a default pipeline.
// options: IdentityClientOptions that adds policies for the pipeline and the authority host that
// will be used to retrieve tokens and authenticate
func newAADIdentityClient(options *TokenCredentialOptions) *aadIdentityClient {
	options = options.setDefaultValues()
	return &aadIdentityClient{options: *options, pipeline: newDefaultPipeline(*options)}
}

// Authenticate creates a client secret authentication request and returns the resulting Access Token or
// an error in case of authentication failure.
// ctx: The current request context
// tenantID: The Azure Active Directory tenant (directory) Id of the service principal
// clientID: The client (application) ID of the service principal
// clientSecret: A client secret that was generated for the App Registration used to authenticate the client
// scopes: The scopes required for the token
func (c *aadIdentityClient) authenticate(ctx context.Context, tenantID string, clientID string, clientSecret string, scopes []string) (*azcore.AccessToken, error) {
	msg, err := c.createClientSecretAuthRequest(tenantID, clientID, clientSecret, scopes)
	if err != nil {
		return nil, err
	}

	resp, err := c.pipeline.Do(ctx, msg)
	if err != nil {
		return nil, err
	}

	// This should never happen under normal conditions
	if resp == nil {
		return nil, &AuthenticationFailedError{msg: "Something unexpected happened with the request and received a nil response"}
	}

	if resp.HasStatusCode(successStatusCodes[:]...) {
		return c.createAccessToken(resp)
	}

	return nil, &AuthenticationFailedError{inner: newAuthenticationResponseError(resp)}
}

// AuthenticateCertificate creates a client certificate authentication request and returns an Access Token or
// an error in case of authentication failure.
// ctx: The current request context
// tenantID: The Azure Active Directory tenant (directory) Id of the service principal
// clientID: The client (application) ID of the service principal
// clientCertificatePath: The path to the client certificate PEM file
// scopes: The scopes required for the token
func (c *aadIdentityClient) authenticateCertificate(ctx context.Context, tenantID string, clientID string, clientCertificatePath string, scopes []string) (*azcore.AccessToken, error) {
	msg, err := c.createClientCertificateAuthRequest(tenantID, clientID, clientCertificatePath, scopes)
	if err != nil {
		return nil, err
	}

	resp, err := c.pipeline.Do(ctx, msg)
	if err != nil {
		return nil, err
	}

	// This should never happen under normal conditions
	if resp == nil {
		return nil, &AuthenticationFailedError{msg: "Something unexpected happened with the request and received a nil response"}
	}

	if resp.HasStatusCode(successStatusCodes[:]...) {
		return c.createAccessToken(resp)
	}

	return nil, &AuthenticationFailedError{inner: newAuthenticationResponseError(resp)}
}

func (c *aadIdentityClient) createAccessToken(res *azcore.Response) (*azcore.AccessToken, error) {
	value := &internalAccessToken{}
	accessToken := &azcore.AccessToken{}
	if err := json.Unmarshal(res.Payload, &value); err != nil {
		return nil, fmt.Errorf("internalAccessToken: %w", err)
	}
	t, err := value.ExpiresIn.Int64()
	if err != nil {
		return nil, err
	}
	// NOTE: look at go-autorest
	accessToken.Token = value.Token
	accessToken.ExpiresOn = time.Now().Add(time.Second * time.Duration(t)).UTC()
	return accessToken, nil
}

func (c *aadIdentityClient) createClientSecretAuthRequest(tenantID string, clientID string, clientSecret string, scopes []string) (*azcore.Request, error) {
	urlStr := c.options.AuthorityHost.String() + tenantID + tokenEndpoint
	urlFormat, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Set(qpGrantType, "client_credentials")
	data.Set(qpClientID, clientID)
	data.Set(qpClientSecret, clientSecret)
	data.Set(qpScope, strings.Join(scopes, " "))
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))
	msg := azcore.NewRequest(http.MethodPost, *urlFormat)
	msg.Header.Set(azcore.HeaderContentType, azcore.HeaderURLEncoded)
	err = msg.SetBody(body)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (c *aadIdentityClient) createClientCertificateAuthRequest(tenantID string, clientID string, clientCertificate string, scopes []string) (*azcore.Request, error) {
	urlStr := c.options.AuthorityHost.String() + tenantID + tokenEndpoint
	urlFormat, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	clientAssertion, err := createClientAssertionJWT(clientID, urlStr, clientCertificate)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set(qpGrantType, "client_credentials")
	data.Set(qpResponseType, "token")
	data.Set(qpClientID, clientID)
	data.Set(qpClientAssertionType, clientAssertionType)
	data.Set(qpClientAssertion, clientAssertion)
	data.Set(qpScope, strings.Join(scopes, " "))
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))
	msg := azcore.NewRequest(http.MethodPost, *urlFormat)
	msg.Header.Set(azcore.HeaderContentType, azcore.HeaderURLEncoded)

	err = msg.SetBody(body)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// AuthenticateUsernamePassword creates a client username and password authentication request and returns an Access Token or
// an error in case of authentication failure.
// ctx: The current request context
// tenantID: The Azure Active Directory tenant (directory) Id of the service principal
// clientID: The client (application) ID of the service principal
// username: User's account username
// password: User's account password
// scopes: The scopes required for the token
func (c *aadIdentityClient) authenticateUsernamePassword(ctx context.Context, tenantID string, clientID string, username string, password string, scopes []string) (*azcore.AccessToken, error) {
	msg, err := c.createUsernamePasswordAuthRequest(tenantID, clientID, username, password, scopes)
	if err != nil {
		return nil, err
	}

	resp, err := c.pipeline.Do(ctx, msg)
	if err != nil {
		return nil, err
	}

	// This should never happen under normal conditions
	if resp == nil {
		return nil, &AuthenticationFailedError{msg: "Something unexpected happened with the request and received a nil response"}
	}

	if resp.HasStatusCode(successStatusCodes[:]...) {
		return c.createAccessToken(resp)
	}

	return nil, &AuthenticationFailedError{inner: newAuthenticationResponseError(resp)}
}

func (c *aadIdentityClient) createUsernamePasswordAuthRequest(tenantID string, clientID string, username string, password string, scopes []string) (*azcore.Request, error) {
	urlStr := c.options.AuthorityHost.String() + tenantID + tokenEndpoint
	urlFormat, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set(qpResponseType, "token")
	data.Set(qpGrantType, "password")
	data.Set(qpClientID, clientID)
	data.Set(qpUsername, username)
	data.Set(qpPassword, password)
	data.Set(qpScope, strings.Join(scopes, " "))
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))
	msg := azcore.NewRequest(http.MethodPost, *urlFormat)
	msg.Header.Set(azcore.HeaderContentType, azcore.HeaderURLEncoded)
	err = msg.SetBody(body)
	if err != nil {
		return nil, err
	}
	return msg, nil
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
	c.lock.Lock()
	c.refreshToken = tk
	c.lock.Unlock()
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

	resp, err := c.pipeline.Do(ctx, msg)
	if err != nil {
		return nil, err
	}

	// This should never happen under normal conditions
	if resp == nil {
		return nil, &AuthenticationFailedError{msg: "Something unexpected happened with the request and received a nil response"}
	}

	if resp.HasStatusCode(successStatusCodes[:]...) {
		return c.createDeviceCodeAccessToken(resp)
	}

	return nil, &AuthenticationFailedError{inner: newAuthenticationResponseError(resp)}
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
	msg := azcore.NewRequest(http.MethodPost, *urlFormat)
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

	resp, err := c.pipeline.Do(ctx, msg)
	if err != nil {
		return nil, err
	}
	// This should never happen under normal conditions
	if resp == nil {
		return nil, &AuthenticationFailedError{msg: "Something unexpected happened with the request and received a nil response"}
	}
	if resp.HasStatusCode(successStatusCodes[:]...) {
		return createDeviceCodeResult(resp)
	}
	return nil, &AuthenticationFailedError{inner: newAuthenticationResponseError(resp)}
}

func (c *aadIdentityClient) createDeviceCodeNumberRequest(tenantID string, clientID string, scopes []string) (*azcore.Request, error) {
	if len(tenantID) == 0 { // if the user did not pass in a tenantID then the default value is set
		tenantID = "organizations"
	}
	urlStr := c.options.AuthorityHost.String() + tenantID + "/oauth2/v2.0/devicecode" // endpoint that will return a device code along with the other necessary authentication flow parameters in the DeviceCodeResult struct
	urlFormat, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set(qpClientID, clientID)
	data.Set(qpScope, strings.Join(scopes, " "))
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))
	msg := azcore.NewRequest(http.MethodPost, *urlFormat)
	msg.Header.Set(azcore.HeaderContentType, azcore.HeaderURLEncoded)
	err = msg.SetBody(body)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func getPrivateKey(cert string) (*rsa.PrivateKey, error) {
	privateKeyFile, err := os.Open(cert)
	defer privateKeyFile.Close()
	if err != nil {
		return nil, fmt.Errorf("Opening certificate file path: %w", err)
	}

	pemFileInfo, err := privateKeyFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("Getting certificate file info: %w", err)
	}
	size := pemFileInfo.Size()

	pemBytes := make([]byte, size)
	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("Read PEM file bytes: %w", err)
	}

	data, rest := pem.Decode([]byte(pemBytes))
	const privateKeyBlock = "PRIVATE KEY"
	// NOTE: check types of private keys
	if data.Type != privateKeyBlock {
		for len(rest) > 0 {
			data, rest = pem.Decode(rest)
			if data.Type == privateKeyBlock {
				privateKeyImported, err := x509.ParsePKCS8PrivateKey(data.Bytes)
				if err != nil {
					return nil, fmt.Errorf("ParsePKCS8PrivateKey: %w", err)
				}

				return privateKeyImported.(*rsa.PrivateKey), nil
			}
		}
		return nil, errors.New("Cannot find PRIVATE KEY in file")
	}
	// NOTE: this could be a function local closure
	privateKeyImported, err := x509.ParsePKCS8PrivateKey(data.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ParsePKCS8PrivateKey: %w", err)
	}

	return privateKeyImported.(*rsa.PrivateKey), nil
}
