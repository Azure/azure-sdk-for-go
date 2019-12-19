// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	authenticationResponseInvalidFormatError = "Invalid response, the authentication response was not in the expected format."
	msiEndpointInvalidURIError               = "The environment variable MSI_ENDPOINT contains an invalid URL."
)

const (
	imdsEndpoint = "http://169.254.169.254/metadata/identity/oauth2/token"
)

const (
	msiEndpointEnvironemntVariable = "MSI_ENDPOINT"
	msiSecretEnvironemntVariable   = "MSI_SECRET"
	appServiceMsiAPIVersion        = "2017-09-01"
	imdsAPIVersion                 = "2018-02-01"
)

type msiType int

const ( // todo rename
	unknown     msiType = 0
	imds        msiType = 1
	appService  msiType = 2
	cloudShell  msiType = 3
	unavailable msiType = 4
)

// ManagedIdentityClient provides the base for authenticating with Managed Identities on Azure VMs and Cloud Shell
// This type initializes a default azcore.Pipeline and IdentityClientOptions.
type managedIdentityClient struct {
	lock                   sync.RWMutex
	options                TokenCredentialOptions
	pipeline               azcore.Pipeline
	imdsAPIVersion         string
	imdsAvailableTimeoutMS int // todo CHANGE TO duration
	msiType                msiType
	endpoint               *url.URL
	refreshToken           string
}

var (
	imdsURL        *url.URL // these are initialized in the init func and are R/O afterwards
	defaultMSIOpts *ManagedIdentityCredentialOptions
)

func init() {
	// The error check is handled in managed_identity_client_test.go
	imdsURL, _ = url.Parse(imdsEndpoint)
	defaultMSIOpts = newDefaultManagedIdentityOptions() // todo move to var initialization
}

func newDefaultManagedIdentityOptions() *ManagedIdentityCredentialOptions {
	return &ManagedIdentityCredentialOptions{
		HTTPClient: azcore.DefaultHTTPClientTransport(),
	}
}

// NewManagedIdentityClient creates a new instance of the ManagedIdentityClient with the IdentityClientOptions
// that are passed into it along with a default pipeline.
// options: IdentityClientOptions that adds policies for the pipeline and the authority host that
// will be used to retrieve tokens and authenticate
func newManagedIdentityClient(options *ManagedIdentityCredentialOptions) *managedIdentityClient {
	options = options.setDefaultValues()
	// TODO document the use of these variables
	return &managedIdentityClient{
		pipeline:               newDefaultMSIPipeline(*options),
		imdsAPIVersion:         imdsAPIVersion,
		imdsAvailableTimeoutMS: 500,
		msiType:                unknown,
	}
}

// Authenticate creates an authentication request for a Managed Identity and returns the resulting Access Token or
// an error in case of authentication failure.
// ctx: The current request context
// clientID: The client (application) ID of the service principal
// scopes: The scopes required for the token
func (c *managedIdentityClient) authenticate(ctx context.Context, clientID string, scopes []string) (*azcore.AccessToken, error) {
	currentMSI, err := c.getMSIType(ctx)
	if err != nil {
		return nil, err
	}
	// This condition should never be true since getMSIType returns an error in these cases
	// if msi is unavailable or we were unable to determine the type return a nil access token
	if currentMSI == unavailable || currentMSI == unknown {
		return nil, &CredentialUnavailableError{CredentialType: "Managed Identity Credential", Message: "Please make sure you are running in a managed identity environment, such as a VM, Azure Functions, Cloud Shell, etc..."}
	}

	AT, err := c.sendAuthRequest(ctx, currentMSI, clientID, scopes)
	if err != nil {
		return nil, err
	}
	return AT, nil
}

func (c *managedIdentityClient) sendAuthRequest(ctx context.Context, msiType msiType, clientID string, scopes []string) (*azcore.AccessToken, error) {
	msg, err := c.createAuthRequest(msiType, clientID, scopes)
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

func (c *managedIdentityClient) createAccessToken(res *azcore.Response) (*azcore.AccessToken, error) {
	value := struct {
		// these are the only fields that we use
		Token        string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
		ExpiresIn    json.Number `json:"expires_in"`
		ExpiresOn    string      `json:"expires_on"` // the value returned in this field varies between a number and a date string
	}{}
	if err := json.Unmarshal(res.Payload, &value); err != nil {
		return nil, fmt.Errorf("internalAccessToken: %w", err)
	}
	if value.ExpiresIn != "" {
		expiresIn, err := value.ExpiresIn.Int64()
		if err != nil {
			return nil, err
		}
		return &azcore.AccessToken{Token: value.Token, ExpiresOn: time.Now().Add(time.Second * time.Duration(expiresIn)).UTC()}, nil
	}
	if expiresOn, err := strconv.Atoi(value.ExpiresOn); err == nil {
		return &azcore.AccessToken{Token: value.Token, ExpiresOn: time.Now().Add(time.Second * time.Duration(expiresOn)).UTC()}, nil
	}
	// this is the case when expires_on is a time string
	// this is the format of the string coming from the service
	if expiresOn, err := time.Parse("01/02/2006 15:04:05 PM +00:00", value.ExpiresOn); err == nil {
		eo := expiresOn.UTC()
		return &azcore.AccessToken{Token: value.Token, ExpiresOn: eo}, nil
	} else {
		return nil, err
	}
}

func (c *managedIdentityClient) createAuthRequest(msiType msiType, clientID string, scopes []string) (*azcore.Request, error) {
	switch msiType {
	case imds:
		return c.createIMDSAuthRequest(scopes), nil
	case appService:
		return c.createAppServiceAuthRequest(clientID, scopes), nil
	case cloudShell:
		return c.createCloudShellAuthRequest(clientID, scopes)
	default:
		errorMsg := ""
		switch msiType {
		case unavailable:
			errorMsg = "unavailable"
		default:
			errorMsg = "unknown"
		}

		return nil, &CredentialUnavailableError{CredentialType: "Managed Identity Credential", Message: "Make sure you are running in a valid Managed Identity Environment. Status: " + errorMsg}
	}
}

func (c *managedIdentityClient) createIMDSAuthRequest(scopes []string) *azcore.Request {
	request := azcore.NewRequest(http.MethodGet, *c.endpoint)
	request.Header.Set(azcore.HeaderMetadata, "true")
	q := request.URL.Query()
	q.Add("api-version", c.imdsAPIVersion)
	q.Add("resource", strings.Join(scopes, " "))
	request.URL.RawQuery = q.Encode()

	return request
}

func (c *managedIdentityClient) createAppServiceAuthRequest(clientID string, scopes []string) *azcore.Request {
	request := azcore.NewRequest(http.MethodGet, *c.endpoint)
	request.Header.Set("secret", os.Getenv(msiSecretEnvironemntVariable))
	q := request.URL.Query()
	q.Add("api-version", appServiceMsiAPIVersion)
	q.Add("resource", strings.Join(scopes, " "))
	if clientID != "" {
		q.Add("client_id", clientID)
	}
	request.URL.RawQuery = q.Encode()

	return request
}

func (c *managedIdentityClient) createCloudShellAuthRequest(clientID string, scopes []string) (*azcore.Request, error) {
	request := azcore.NewRequest(http.MethodPost, *c.endpoint)
	request.Header.Set(azcore.HeaderContentType, azcore.HeaderURLEncoded)
	request.Header.Set(azcore.HeaderMetadata, "true")
	data := url.Values{}
	data.Set("resource", strings.Join(scopes, " "))
	if clientID != "" {
		data.Set("client_id", clientID)
	}
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))
	err := request.SetBody(body)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (c *managedIdentityClient) getMSIType(ctx context.Context) (msiType, error) {
	if c.msiType == unknown { // if we haven't already determined the msi type
		if endpointEnvVar := os.Getenv(msiEndpointEnvironemntVariable); endpointEnvVar != "" { // if the env var MSI_ENDPOINT is set
			endpoint, err := url.Parse(endpointEnvVar)
			if err != nil {
				return unknown, err
			}
			c.endpoint = endpoint
			if secretEnvVar := os.Getenv(msiSecretEnvironemntVariable); secretEnvVar != "" { // if BOTH the env vars MSI_ENDPOINT and MSI_SECRET are set the MsiType is AppService
				c.msiType = appService
			} else { // if ONLY the env var MSI_ENDPOINT is set the MsiType is CloudShell
				c.msiType = cloudShell
			}
		} else if c.imdsAvailable(ctx) { // if MSI_ENDPOINT is NOT set AND the IMDS endpoint is available the MsiType is Imds. This will timeout after 500 milliseconds
			c.endpoint = imdsURL
			c.msiType = imds
		} else { // if MSI_ENDPOINT is NOT set and IMDS enpoint is not available ManagedIdentity is not available
			c.msiType = unavailable
			return c.msiType, &CredentialUnavailableError{CredentialType: "Managed Identity Credential", Message: "Make sure you are running in a valid Managed Identity Environment"}
		}
	}
	return c.msiType, nil
}

func (c *managedIdentityClient) imdsAvailable(ctx context.Context) bool {
	tempCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	request := azcore.NewRequest(http.MethodGet, *imdsURL)
	q := request.URL.Query()
	q.Add("api-version", c.imdsAPIVersion)
	request.URL.RawQuery = q.Encode()
	_, err := c.pipeline.Do(tempCtx, request)
	return err == nil
}
