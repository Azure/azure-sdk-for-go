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
	"strings"
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

const (
	unknown     msiType = 0
	imds        msiType = 1
	appService  msiType = 2
	cloudShell  msiType = 3
	unavailable msiType = 4
)

// ManagedIdentityClient provides the base for authenticating with Managed Identities on Azure VMs and Cloud Shell
// This type initializes a default azcore.Pipeline and IdentityClientOptions.
type managedIdentityClient struct {
	options                IdentityClientOptions
	pipeline               azcore.Pipeline
	imdsEndpoint           *url.URL
	imdsAPIVersion         string
	imdsAvailableTimeoutMS int
	msiType                msiType
	endpoint               *url.URL
}

// NewManagedIdentityClient creates a new instance of the ManagedIdentityClient with the IdentityClientOptions
// that are passed into it along with a default pipeline.
// options: IdentityClientOptions that adds policies for the pipeline and the authority host that
// will be used to retrieve tokens and authenticate
func newManagedIdentityClient(options *ManagedIdentityCredentialOptions) (*managedIdentityClient, error) {
	// TODO: mimic aad client and make string a const parse only once, create in init
	imdsEndpoint, err := url.Parse(imdsEndpoint)
	if err != nil {
		return nil, fmt.Errorf("NewManagedIdentityClient: %w", err)
	}

	if options == nil {
		options = NewDefaultManagedIdentityOptions()
	}
	// document the use of these variables
	// TODO: create a separate pipeline for imds that had its own retry policy
	return &managedIdentityClient{
		// TODO: setDefaultValues should return by val
		options:                options.Options,
		pipeline:               newDefaultPipeline(options.PipelineOptions),
		imdsEndpoint:           imdsEndpoint,
		imdsAPIVersion:         imdsAPIVersion,
		imdsAvailableTimeoutMS: 500,
		msiType:                unknown,
	}, nil
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
	// if msi is unavailable or we were unable to determine the type return a nil access token
	if currentMSI == unavailable || currentMSI == unknown {
		// TODO: add message
		// CP: maybe this should be an auth failed error?
		return nil, &CredentialUnavailableError{Message: "Please make sure you are running in a managed identity environment, such as a VM, Azure Functions, Cloud Shell, etc..."}
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

	resp, err := msg.Do(ctx)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		return c.createAccessToken(resp)
	}

	var authFailed AuthenticationFailedError
	err = json.Unmarshal(resp.Payload, &authFailed)
	if err != nil {
		authFailed.Message = resp.Status
		authFailed.Description = "Failed to unmarshal response: " + err.Error()
	}
	authFailed.Response = resp
	return nil, &authFailed
}

func (c *managedIdentityClient) createAccessToken(res *azcore.Response) (*azcore.AccessToken, error) {
	value := azcore.AccessToken{}
	if err := json.Unmarshal(res.Payload, &value); err != nil {
		return nil, fmt.Errorf("azcore.AccessToken: %w", err)
	}
	// CP: This will change based on the MSI type
	t, err := value.ExpiresIn.Int64()
	if err != nil {
		return nil, err
	}
	// NOTE: look at go-autorest
	value.ExpiresOn = time.Now().Add(time.Second * time.Duration(t)).UTC()
	return &value, nil
}

func (c *managedIdentityClient) createAuthRequest(msiType msiType, clientID string, scopes []string) (*azcore.Request, error) {
	var req *azcore.Request
	var err error
	switch msiType {
	case imds:
		req, err = c.createIMDSAuthRequest(clientID, scopes)
	case appService:
		req, err = c.createAppServiceAuthRequest(clientID, scopes)
	case cloudShell:
		req, err = c.createCloudShellAuthRequest(clientID, scopes)
	default:
		return nil, &CredentialUnavailableError{Message: "Make sure you are running in a valid Managed Identity Environment"}
	}

	return req, err
}

func (c *managedIdentityClient) createIMDSAuthRequest(clientID string, scopes []string) (*azcore.Request, error) {
	request := c.pipeline.NewRequest(http.MethodGet, *c.endpoint)
	request.Header.Set(azcore.HeaderMetadata, "true")
	q := request.URL.Query()
	q.Add("api-version", c.imdsAPIVersion)
	q.Add("resource", strings.Join(scopes, " "))
	request.URL.RawQuery = q.Encode()

	return request, nil
}

func (c *managedIdentityClient) createAppServiceAuthRequest(clientID string, scopes []string) (*azcore.Request, error) {
	request := c.pipeline.NewRequest(http.MethodGet, *c.endpoint)
	request.Header.Set("secret", os.Getenv(msiSecretEnvironemntVariable))
	q := request.URL.Query()
	q.Add("api-version", appServiceMsiAPIVersion)
	q.Add("resource", strings.Join(scopes, " "))
	if clientID != "" {
		q.Add("client_id", clientID)
	}
	request.URL.RawQuery = q.Encode()

	return request, nil
}

func (c *managedIdentityClient) createCloudShellAuthRequest(clientID string, scopes []string) (*azcore.Request, error) {
	request := c.pipeline.NewRequest(http.MethodPost, *c.endpoint)
	request.Header.Set(azcore.HeaderContentType, "application/x-www-form-urlencoded")
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
		} else if c.imdsAvailable(ctx) { // if MSI_ENDPOINT is NOT set AND the IMDS endpoint is available the MsiType is Imds
			c.endpoint = c.imdsEndpoint
			c.msiType = imds
		} else { // if MSI_ENDPOINT is NOT set and IMDS enpoint is not available ManagedIdentity is not available
			// CP: should we just fail here? Or is it fine to fail in the func that did the calling?
			c.msiType = unavailable
			// TODO: return a cred unavailable err
			return unknown, &CredentialUnavailableError{Message: "Make sure you are running in a valid Managed Identity Environment"}
		}
	}
	return c.msiType, nil
}

func (c *managedIdentityClient) imdsAvailable(ctx context.Context) bool {
	request := c.pipeline.NewRequest(http.MethodGet, *c.imdsEndpoint)
	q := request.URL.Query()
	q.Add("api-version", c.imdsAPIVersion)
	request.URL.RawQuery = q.Encode()
	// TODO: missing setting timeout and handling it
	// TODO: check for deadline exceeded, return an error
	_, err := request.Do(ctx)
	if err != nil {
		return false
	}

	return true
}

// MISSING: ADD CONSIDERATIONS THAT THIS FUNC INCLUDES IN THE DESERIALIZATION OF THE ACCESSTOKEN
// private static AccessToken Deserialize(JsonElement json)
//         {
//             if (!json.TryGetProperty("access_token", out JsonElement accessTokenProp))
//             {
//                 throw new AuthenticationFailedException(AuthenticationResponseInvalidFormatError);
//             }

//             string accessToken = accessTokenProp.GetString();
//             if (!json.TryGetProperty("expires_on", out JsonElement expiresOnProp))
//             {
//                 throw new AuthenticationFailedException(AuthenticationResponseInvalidFormatError);
//             }

//             DateTimeOffset expiresOn;
//             // if s_msiType is AppService expires_on will be a string formatted datetimeoffset
//             if (s_msiType == MsiType.AppService)
//             {
//                 if (!DateTimeOffset.TryParse(expiresOnProp.GetString(), out expiresOn))
//                 {
//                     throw new AuthenticationFailedException(AuthenticationResponseInvalidFormatError);
//                 }
//             }
//             // otherwise expires_on will be a unix timestamp seconds from epoch
//             else
//             {
//                 // the seconds from epoch may be returned as a Json number or a Json string which is a number
//                 // depending on the environment.  If neither of these are the case we throw an AuthException.
//                 if (!(expiresOnProp.ValueKind == JsonValueKind.Number && expiresOnProp.TryGetInt64(out long expiresOnSec)) &&
//                     !(expiresOnProp.ValueKind == JsonValueKind.String && long.TryParse(expiresOnProp.GetString(), out expiresOnSec)))
//                 {
//                     throw new AuthenticationFailedException(AuthenticationResponseInvalidFormatError);
//                 }

//                 expiresOn = DateTimeOffset.FromUnixTimeSeconds(expiresOnSec);
//             }

//             return new AccessToken(accessToken, expiresOn);
//         }
