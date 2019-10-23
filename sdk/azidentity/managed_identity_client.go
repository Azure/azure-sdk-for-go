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

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	authenticationResponseInvalidFormatError         = "Invalid response, the authentication response was not in the expected format."
	msiEndpointInvalidURIError                       = "The environment variable MSI_ENDPOINT contains an invalid Uri."
	authenticationRequestFailedError                 = "The request to the identity service failed. See inner exception for details."
	msiEndpointEnvironemntVariable                   = "MSI_ENDPOINT"
	msiSecretEnvironemntVariable                     = "MSI_SECRET"
	appServiceMsiAPIVersion                          = "2017-09-01"
	unknown                                  msiType = 0
	imds                                     msiType = 1
	appService                               msiType = 2
	cloudShell                               msiType = 3
	unavailable                              msiType = 4
)

type msiType int

// ManagedIdentityClient ...
type ManagedIdentityClient struct {
	options                IdentityClientOptions
	pipeline               azcore.Pipeline
	sIMDSEndpoint          url.URL
	imdsAPIVersion         string
	imdsAvailableTimeoutMS int
	sMSIType               msiType
	sEndpoint              url.URL
}

// NewManagedIdentityClient ...
func NewManagedIdentityClient(options *IdentityClientOptions) *ManagedIdentityClient {
	if options == nil {
		options, _ = newIdentityClientOptions()
	}

	sIMDSEndpoint, err := url.Parse("http://169.254.169.254/metadata/identity/oauth2/token")
	if err != nil {
		// return nil, fmt.Errorf("NewManagedIdentityClient: %w", err)
		return &ManagedIdentityClient{}
	}
	imdsAPIVersion := "2018-02-01"
	imdsAvailableTimeoutMs := 500
	sMSIType := unknown
	client := ManagedIdentityClient{options: *options, pipeline: NewDefaultPipeline(options.pipelineOptions), sIMDSEndpoint: *sIMDSEndpoint, imdsAPIVersion: imdsAPIVersion, imdsAvailableTimeoutMS: imdsAvailableTimeoutMs, sMSIType: sMSIType}
	return &client
}

// Authenticate ...
func (c *ManagedIdentityClient) Authenticate(ctx context.Context, clientID string, scopes []string) (*AccessToken, error) {
	// TODO: fix variable name
	sType, err := c.getMSIType(ctx)
	if err != nil {
		return nil, fmt.Errorf("Authenticate: %w", err)
	}
	// if msi is unavailable or we were unable to determine the type return a default access token
	if sType == int(unavailable) || sType == int(unknown) {
		// CP: TO DO this needs to return a default accessToken
		return nil, nil
	}

	AT, err := c.sendAuthRequest(ctx, sType, clientID, scopes)
	if err != nil {
		// return nil, errors.New(authenticationRequestFailedError)
		return nil, fmt.Errorf("Failed to authenticate: %w", err)
	}
	return AT, nil
}

func (c *ManagedIdentityClient) sendAuthRequest(ctx context.Context, msiType int, clientID string, scopes []string) (*AccessToken, error) {
	msg, err := c.createAuthRequest(msiType, clientID, scopes)
	if err != nil {
		return nil, fmt.Errorf("SendAuthRequest: %w", err)
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

func (c *ManagedIdentityClient) createAccessToken(res *azcore.Response) (*AccessToken, error) {
	// CP: what is the best method to initialize this?
	value := &AccessToken{}
	// value := NewAccessToken("", 0)
	jd := json.NewDecoder(res.Body)
	err := jd.Decode(&value)
	if err != nil {
		return nil, fmt.Errorf("Decode: %w", err)
	}
	// CP: CHECK THIS
	value.SetExpiresOn()

	return value, nil
}

func (c *ManagedIdentityClient) createAuthRequest(msiType int, clientID string, scopes []string) (*azcore.Request, error) {
	var req *azcore.Request
	var err error
	// CP: TODO check this switch
	switch msiType {
	case int(imds):
		req, err = c.createImdsAuthRequest(clientID, scopes)
	case int(appService):
		req, err = c.createAppServiceAuthRequest(clientID, scopes)
	case int(cloudShell):
		req, err = c.createCloudShellAuthRequest(clientID, scopes)
	default:
		// TODO: missing this
	}

	return req, err
}

func (c *ManagedIdentityClient) createImdsAuthRequest(clientID string, scopes []string) (*azcore.Request, error) {
	// TODO: check other sdk's handling of resources, it's in ScopeUtilities
	// // covert the scopes to a resource string
	// resource = ScopeUtilities.ScopesToResource(scopes);

	request := c.pipeline.NewRequest(http.MethodGet, c.sEndpoint)
	request.Header = http.Header{"Metadata": []string{"true"}}
	q := request.URL.Query()
	q.Add("api-version", c.imdsAPIVersion)
	q.Add("resource", strings.Join(scopes, " "))
	request.URL.RawQuery = q.Encode()
	// data := url.Values{}
	// data.Set("api-version", c.imdsAPIVersion)
	// data.Set("resource", strings.Join(scopes, " "))
	// if clientID != "" {
	// 	data.Set("client_id", clientID)
	// }
	// dataEncoded := data.Encode()
	// body := azcore.NopCloser(strings.NewReader(dataEncoded))
	// request.SetBody(body)

	return request, nil
}

// CP: TODO need to fix and test this to add query params to the raw query
func (c *ManagedIdentityClient) createAppServiceAuthRequest(clientID string, scopes []string) (*azcore.Request, error) {
	// TODO: check other sdk's handling of resources, it's in ScopeUtilities
	// // covert the scopes to a resource string
	// resource = ScopeUtilities.ScopesToResource(scopes);

	request := c.pipeline.NewRequest(http.MethodGet, c.sEndpoint)
	request.Header = http.Header{"secret": []string{os.Getenv(msiSecretEnvironemntVariable)}}

	data := url.Values{}
	data.Set("api-version", appServiceMsiAPIVersion)
	data.Set("resource", strings.Join(scopes, " "))
	if clientID != "" {
		data.Set("client_id", clientID)
	}
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))
	request.SetBody(body)

	return request, nil
}

// CP: TODO need to fix and test this to add query params to the raw query
func (c *ManagedIdentityClient) createCloudShellAuthRequest(clientID string, scopes []string) (*azcore.Request, error) {
	// TODO: check other sdk's handling of resources, it's in ScopeUtilities
	// // covert the scopes to a resource string
	// resource = ScopeUtilities.ScopesToResource(scopes);

	request := c.pipeline.NewRequest(http.MethodPost, c.sEndpoint)
	request.Header = http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}}
	request.Header = http.Header{"Metadata": []string{"true"}}
	data := url.Values{}
	data.Set("resource", strings.Join(scopes, " "))
	if clientID != "" {
		data.Set("client_id", clientID)
	}
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))
	request.SetBody(body)

	return request, nil
}

func (c *ManagedIdentityClient) getMSIType(ctx context.Context) (int, error) {
	if c.sMSIType == unknown { // if we haven't already determined the msi type
		endpointEnvVar := os.Getenv(msiEndpointEnvironemntVariable)
		secretEnvVar := os.Getenv(msiSecretEnvironemntVariable)

		if endpointEnvVar != "" { // if the env var MSI_ENDPOINT is set
			sEndpoint, err := url.Parse(endpointEnvVar)
			if err != nil {
				return -1, fmt.Errorf("getMSIType: %w", err)
			}
			c.sEndpoint = *sEndpoint
			if secretEnvVar != "" { // if BOTH the env vars MSI_ENDPOINT and MSI_SECRET are set the MsiType is AppService
				c.sMSIType = appService
			} else { // if ONLY the env var MSI_ENDPOINT is set the MsiType is CloudShell
				c.sMSIType = cloudShell
			}
		} else if c.imdsAvailable(ctx) { // if MSI_ENDPOINT is NOT set AND the IMDS endpoint is available the MsiType is Imds
			c.sEndpoint = c.sIMDSEndpoint
			c.sMSIType = imds
		} else { // if MSI_ENDPOINT is NOT set and IMDS enpoint is not available ManagedIdentity is not available
			c.sMSIType = unavailable
		}

	}
	return int(c.sMSIType), nil
}

func (c *ManagedIdentityClient) imdsAvailable(ctx context.Context) bool {
	req := c.pipeline.NewRequest(http.MethodGet, c.sIMDSEndpoint)
	data := url.Values{}
	data.Set("api-version", c.imdsAPIVersion)
	dataEncoded := data.Encode()
	body := azcore.NopCloser(strings.NewReader(dataEncoded))
	req.SetBody(body)

	// TODO: missing setting timeout and handling it
	_, err := req.Do(ctx)
	if err != nil {
		return false
	}
	// TODO: check logic here important
	return true
}

// MISSING THIS
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
