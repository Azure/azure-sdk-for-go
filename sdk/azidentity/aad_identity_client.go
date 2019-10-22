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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	clientAssertionType = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
)

// aadIdentityClient type provides the base for authenticating with Client Secret Credentials, Client Certificate Credentials
// and Environment Credentials. This type initializes a default azcore.Pipeline and IdentityClientOptions.
// CP: Should this be exported?
type aadIdentityClient struct {
	options  IdentityClientOptions
	pipeline azcore.Pipeline
}

// newAADIdentityClient creates a new instance of the aadIdentityClient with the IdentityClientOptions that are passed into it
// along with a default pipeline
// CP: consider scenario where we want to pass a custom pipeline to aadIdentityClient, what would that func look like? func (c *AADIdentityClient) NewAADIdentityClientCustomPipeline() bool. If this func is public then aadIdentityClient should also be exported
func newAADIdentityClient(options *IdentityClientOptions) aadIdentityClient {
	if options == nil {
		options, _ = newIdentityClientOptions()
	}
	client := aadIdentityClient{options: *options, pipeline: NewDefaultPipeline(options.pipelineOptions)}
	return client
}

// Authenticate creates a client secret authentication request and returns the resulting Access Token or an error in case of authentication failure.
func (c aadIdentityClient) Authenticate(ctx context.Context, tenantID string, clientID string, clientSecret string, scopes []string) (*AccessToken, error) {
	// CP: Need to activate pipeline diagnostics before beginning the auth process...
	msg, err := c.createClientSecretAuthRequest(tenantID, clientID, clientSecret, scopes)
	if err != nil {
		return nil, fmt.Errorf("createClientSecretAuthRequest: %w", err)
	}

	resp, err := msg.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("azcore.Message Do: %w", err)
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return c.createAccessToken(resp)
	}

	// This should have similar checks like in createResponder in azfile, might even be separated into multiple funcs
	return nil, err
}

// AuthenticateCertificate creates a client certificate authentication request and returns an Access Token or an error in case of authentication failure.
func (c aadIdentityClient) AuthenticateCertificate(ctx context.Context, tenantID string, clientID string, clientCertificatePath string, scopes []string) (*AccessToken, error) {
	// CP: Need to activate pipeline diagnostics before beginning the auth process... azcore distributed tracing policy?
	msg, err := c.createClientCertificateAuthRequest(tenantID, clientID, clientCertificatePath, scopes)
	if err != nil {
		return nil, fmt.Errorf("createClientCertificateAuthRequest: %w", err)
	}

	resp, err := msg.Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("azcore.Message Do: %w", err)
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return c.createAccessToken(resp)
	}

	// CP: Check what error to return in this case.
	return nil, fmt.Errorf("%s", authenticationRequestFailedError)
}

func (c aadIdentityClient) createAccessToken(res *azcore.Response) (*AccessToken, error) {
	// CP: what is the best method to initialize this?
	value := &AccessToken{}
	// value := NewAccessToken("", 0)
	jd := json.NewDecoder(res.Body)
	err := jd.Decode(&value)
	if err != nil {
		return nil, fmt.Errorf("AccessToken: %w", err)
	}
	// CP: CHECK THIS
	value.SetExpiresOn()

	return value, nil
}

func (c aadIdentityClient) createClientSecretAuthRequest(tenantID string, clientID string, clientSecret string, scopes []string) (*azcore.Request, error) {
	urlStr := c.options.AuthorityHost.String() + tenantID + "/oauth2/v2.0/token/"
	urlFormat, err := url.Parse(urlStr)
	// CP: FIX THIS
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("scope", strings.Join(scopes, " "))
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))

	msg := c.pipeline.NewRequest(http.MethodPost, *urlFormat)

	msg.Header = http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}
	msg.SetBody(body)

	return msg, nil
}

func (c aadIdentityClient) createClientCertificateAuthRequest(tenantID string, clientID string, clientCertificate string, scopes []string) (*azcore.Request, error) {
	urlStr := c.options.AuthorityHost.String() + tenantID + "/oauth2/v2.0/token/"
	urlFormat, err := url.Parse(urlStr)
	// CP: FIX THIS
	if err != nil {
		return nil, err
	}

	clientAssertion, err := createClientAssertionJWT(clientID, urlStr, clientCertificate)
	if err != nil {
		return nil, fmt.Errorf("createClientAssertionJWT: %w", err)
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("response_type", "token")
	data.Set("client_id", clientID)
	data.Set("client_assertion_type", clientAssertionType)
	data.Set("client_assertion", clientAssertion)
	data.Set("scope", strings.Join(scopes, " "))
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))

	msg := c.pipeline.NewRequest(http.MethodPost, *urlFormat)

	msg.Header = http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}
	msg.SetBody(body)

	return msg, nil
}

func getPrivateKey(cert string) (*rsa.PrivateKey, error) {
	privateKeyFile, err := os.Open(cert)
	if err != nil {
		return nil, fmt.Errorf("Opening certificate file path Open(): %w", err)
	}

	defer privateKeyFile.Close()

	pemFileInfo, err := privateKeyFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("Get file info from Stat(): %w", err)
	}
	var size int64 = pemFileInfo.Size()

	pemBytes := make([]byte, size)
	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("Read PEM file bytes from bufio.NewReader().Read(): %w", err)
	}

	data, rest := pem.Decode([]byte(pemBytes))

	if data.Type != "PRIVATE KEY" {
		for len(rest) > 0 {
			data, rest = pem.Decode(rest)
			if data.Type == "PRIVATE KEY" {
				privateKeyImported, err := x509.ParsePKCS8PrivateKey(data.Bytes)
				if err != nil {
					return nil, fmt.Errorf("ParsePKCS8PrivateKey: %w", err)
				}

				return privateKeyImported.(*rsa.PrivateKey), nil
			}
		}
		return nil, errors.New("Cannot find PRIVATE KEY in file")
	}

	privateKeyImported, err := x509.ParsePKCS8PrivateKey(data.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ParsePKCS8PrivateKey: %w", err)
	}

	return privateKeyImported.(*rsa.PrivateKey), nil
}
