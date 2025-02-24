//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	msalerrors "github.com/AzureAD/microsoft-authentication-library-for-go/apps/errors"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/managedidentity"
)

const (
	arcIMDSEndpoint          = "IMDS_ENDPOINT"
	defaultIdentityClientID  = "DEFAULT_IDENTITY_CLIENT_ID"
	identityEndpoint         = "IDENTITY_ENDPOINT"
	identityHeader           = "IDENTITY_HEADER"
	identityServerThumbprint = "IDENTITY_SERVER_THUMBPRINT"
	headerMetadata           = "Metadata"
	imdsEndpoint             = "http://169.254.169.254/metadata/identity/oauth2/token"
	miResID                  = "mi_res_id"
	msiEndpoint              = "MSI_ENDPOINT"
	msiResID                 = "msi_res_id"
	msiSecret                = "MSI_SECRET"
	imdsAPIVersion           = "2018-02-01"
	azureArcAPIVersion       = "2020-06-01"
	qpClientID               = "client_id"
	serviceFabricAPIVersion  = "2019-07-01-preview"
)

var imdsProbeTimeout = time.Second

type managedIdentityClient struct {
	azClient                      *azcore.Client
	imds, probeIMDS, userAssigned bool
	// chained indicates whether the client is part of a credential chain. If true, the client will return
	// a credentialUnavailableError instead of an AuthenticationFailedError for an unexpected IMDS response.
	chained    bool
	msalClient *managedidentity.Client
}

// setIMDSRetryOptionDefaults sets zero-valued fields to default values appropriate for IMDS
func setIMDSRetryOptionDefaults(o *policy.RetryOptions) {
	if o.MaxRetries == 0 {
		o.MaxRetries = 5
	}
	if o.MaxRetryDelay == 0 {
		o.MaxRetryDelay = 1 * time.Minute
	}
	if o.RetryDelay == 0 {
		o.RetryDelay = 2 * time.Second
	}
	if o.StatusCodes == nil {
		o.StatusCodes = []int{
			// IMDS docs recommend retrying 404, 410, 429 and 5xx
			// https://learn.microsoft.com/entra/identity/managed-identities-azure-resources/how-to-use-vm-token#error-handling
			http.StatusNotFound,                      // 404
			http.StatusGone,                          // 410
			http.StatusTooManyRequests,               // 429
			http.StatusInternalServerError,           // 500
			http.StatusNotImplemented,                // 501
			http.StatusBadGateway,                    // 502
			http.StatusServiceUnavailable,            // 503
			http.StatusGatewayTimeout,                // 504
			http.StatusHTTPVersionNotSupported,       // 505
			http.StatusVariantAlsoNegotiates,         // 506
			http.StatusInsufficientStorage,           // 507
			http.StatusLoopDetected,                  // 508
			http.StatusNotExtended,                   // 510
			http.StatusNetworkAuthenticationRequired, // 511
		}
	}
	if o.TryTimeout == 0 {
		o.TryTimeout = 1 * time.Minute
	}
}

// newManagedIdentityClient creates a new instance of the ManagedIdentityClient with the ManagedIdentityCredentialOptions
// that are passed into it along with a default pipeline.
// options: ManagedIdentityCredentialOptions configure policies for the pipeline and the authority host that
// will be used to retrieve tokens and authenticate
func newManagedIdentityClient(options *ManagedIdentityCredentialOptions) (*managedIdentityClient, error) {
	if options == nil {
		options = &ManagedIdentityCredentialOptions{}
	}
	cp := options.ClientOptions
	c := managedIdentityClient{}
	source, err := managedidentity.GetSource()
	if err != nil {
		return nil, err
	}
	env := string(source)
	if source == managedidentity.DefaultToIMDS {
		env = "IMDS"
		c.imds = true
		c.probeIMDS = options.dac
		setIMDSRetryOptionDefaults(&cp.Retry)
	}

	c.azClient, err = azcore.NewClient(module, version, azruntime.PipelineOptions{
		Tracing: azruntime.TracingOptions{
			Namespace: traceNamespace,
		},
	}, &cp)
	if err != nil {
		return nil, err
	}

	id := managedidentity.SystemAssigned()
	if options.ID != nil {
		c.userAssigned = true
		switch s := options.ID.String(); options.ID.idKind() {
		case miClientID:
			id = managedidentity.UserAssignedClientID(s)
		case miObjectID:
			id = managedidentity.UserAssignedObjectID(s)
		case miResourceID:
			id = managedidentity.UserAssignedResourceID(s)
		}
	}
	msalClient, err := managedidentity.New(id, managedidentity.WithHTTPClient(&c), managedidentity.WithRetryPolicyDisabled())
	if err != nil {
		return nil, err
	}
	c.msalClient = &msalClient

	if log.Should(EventAuthentication) {
		msg := fmt.Sprintf("%s will use %s managed identity", credNameManagedIdentity, env)
		if options.ID != nil {
			kind := "client"
			switch options.ID.(type) {
			case ObjectID:
				kind = "object"
			case ResourceID:
				kind = "resource"
			}
			msg += fmt.Sprintf(" with %s ID %q", kind, options.ID.String())
		}
		log.Write(EventAuthentication, msg)
	}

	return &c, nil
}

func (*managedIdentityClient) CloseIdleConnections() {
	// do nothing
}

func (c *managedIdentityClient) Do(r *http.Request) (*http.Response, error) {
	return doForClient(c.azClient, r)
}

// authenticate acquires an access token
func (c *managedIdentityClient) GetToken(ctx context.Context, tro policy.TokenRequestOptions) (azcore.AccessToken, error) {
	// no need to synchronize around this value because it's true only when DefaultAzureCredential constructed the client,
	// and in that case ChainedTokenCredential.GetToken synchronizes goroutines that would execute this block
	if c.probeIMDS {
		// send a malformed request (no Metadata header) to IMDS to determine whether the endpoint is available
		cx, cancel := context.WithTimeout(ctx, imdsProbeTimeout)
		defer cancel()
		cx = policy.WithRetryOptions(cx, policy.RetryOptions{MaxRetries: -1})
		req, err := azruntime.NewRequest(cx, http.MethodGet, imdsEndpoint)
		if err != nil {
			return azcore.AccessToken{}, fmt.Errorf("failed to create IMDS probe request: %s", err)
		}
		if _, err = c.azClient.Pipeline().Do(req); err != nil {
			msg := err.Error()
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				msg = "managed identity timed out. See https://aka.ms/azsdk/go/identity/troubleshoot#dac for more information"
			}
			return azcore.AccessToken{}, newCredentialUnavailableError(credNameManagedIdentity, msg)
		}
		// send normal token requests from now on because something responded
		c.probeIMDS = false
	}

	ar, err := c.msalClient.AcquireToken(ctx, tro.Scopes[0], managedidentity.WithClaims(tro.Claims))
	if err == nil {
		msg := fmt.Sprintf(scopeLogFmt, credNameManagedIdentity, strings.Join(ar.GrantedScopes, ", "))
		log.Write(EventAuthentication, msg)
		return azcore.AccessToken{Token: ar.AccessToken, ExpiresOn: ar.ExpiresOn.UTC()}, err
	}
	if c.imds {
		var ije msalerrors.InvalidJsonErr
		if c.chained && errors.As(err, &ije) {
			// an unmarshaling error implies the response is from something other than IMDS such as a proxy listening at
			// the same address. Return a credentialUnavailableError so credential chains continue to their next credential
			return azcore.AccessToken{}, newCredentialUnavailableError(credNameManagedIdentity, err.Error())
		}
		resp := getResponseFromError(err)
		if resp == nil {
			return azcore.AccessToken{}, newAuthenticationFailedErrorFromMSAL(credNameManagedIdentity, err)
		}
		switch resp.StatusCode {
		case http.StatusBadRequest:
			if c.userAssigned {
				return azcore.AccessToken{}, newAuthenticationFailedError(credNameManagedIdentity, "the requested identity isn't assigned to this resource", resp)
			}
			msg := "failed to authenticate a system assigned identity"
			if body, err := azruntime.Payload(resp); err == nil && len(body) > 0 {
				msg += fmt.Sprintf(". The endpoint responded with %s", body)
			}
			return azcore.AccessToken{}, newCredentialUnavailableError(credNameManagedIdentity, msg)
		case http.StatusForbidden:
			// Docker Desktop runs a proxy that responds 403 to IMDS token requests. If we get that response,
			// we return credentialUnavailableError so credential chains continue to their next credential
			body, err := azruntime.Payload(resp)
			if err == nil && strings.Contains(string(body), "unreachable") {
				return azcore.AccessToken{}, newCredentialUnavailableError(credNameManagedIdentity, fmt.Sprintf("unexpected response %q", string(body)))
			}
		}
	}
	return azcore.AccessToken{Token: ar.AccessToken, ExpiresOn: ar.ExpiresOn.UTC()}, err
}

func (c *managedIdentityClient) createAccessToken(res *http.Response) (azcore.AccessToken, error) {
	value := struct {
		// these are the only fields that we use
		Token        string        `json:"access_token,omitempty"`
		RefreshToken string        `json:"refresh_token,omitempty"`
		ExpiresIn    wrappedNumber `json:"expires_in,omitempty"` // this field should always return the number of seconds for which a token is valid
		ExpiresOn    interface{}   `json:"expires_on,omitempty"` // the value returned in this field varies between a number and a date string
	}{}
	if err := azruntime.UnmarshalAsJSON(res, &value); err != nil {
		return azcore.AccessToken{}, newAuthenticationFailedError(credNameManagedIdentity, "Unexpected response content", res)
	}
	if value.ExpiresIn != "" {
		expiresIn, err := json.Number(value.ExpiresIn).Int64()
		if err != nil {
			return azcore.AccessToken{}, err
		}
		return azcore.AccessToken{Token: value.Token, ExpiresOn: time.Now().Add(time.Second * time.Duration(expiresIn)).UTC()}, nil
	}
	switch v := value.ExpiresOn.(type) {
	case float64:
		return azcore.AccessToken{Token: value.Token, ExpiresOn: time.Unix(int64(v), 0).UTC()}, nil
	case string:
		if expiresOn, err := strconv.Atoi(v); err == nil {
			return azcore.AccessToken{Token: value.Token, ExpiresOn: time.Unix(int64(expiresOn), 0).UTC()}, nil
		}
		return azcore.AccessToken{}, newAuthenticationFailedError(credNameManagedIdentity, "unexpected expires_on value: "+v, res)
	default:
		msg := fmt.Sprintf("unsupported type received in expires_on: %T, %v", v, v)
		return azcore.AccessToken{}, newAuthenticationFailedError(credNameManagedIdentity, msg, res)
	}
}

func (c *managedIdentityClient) createAuthRequest(ctx context.Context, id ManagedIDKind, scopes []string) (*policy.Request, error) {
	switch c.msiType {
	case msiTypeIMDS:
		return c.createIMDSAuthRequest(ctx, id, scopes)
	case msiTypeAppService:
		return c.createAppServiceAuthRequest(ctx, id, scopes)
	case msiTypeAzureArc:
		// need to perform preliminary request to retreive the secret key challenge provided by the HIMDS service
		key, err := c.getAzureArcSecretKey(ctx, scopes)
		if err != nil {
			msg := fmt.Sprintf("failed to retreive secret key from the identity endpoint: %v", err)
			return nil, newAuthenticationFailedError(credNameManagedIdentity, msg, nil)
		}
		return c.createAzureArcAuthRequest(ctx, scopes, key)
	case msiTypeAzureML:
		return c.createAzureMLAuthRequest(ctx, id, scopes)
	case msiTypeServiceFabric:
		return c.createServiceFabricAuthRequest(ctx, scopes)
	case msiTypeCloudShell:
		return c.createCloudShellAuthRequest(ctx, scopes)
	default:
		return nil, newCredentialUnavailableError(credNameManagedIdentity, "managed identity isn't supported in this environment")
	}
}

func (c *managedIdentityClient) createIMDSAuthRequest(ctx context.Context, id ManagedIDKind, scopes []string) (*policy.Request, error) {
	request, err := azruntime.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Raw().Header.Set(headerMetadata, "true")
	q := request.Raw().URL.Query()
	q.Set("api-version", imdsAPIVersion)
	q.Set("resource", strings.Join(scopes, " "))
	if id != nil {
		switch id.idKind() {
		case miClientID:
			q.Set(qpClientID, id.String())
		case miObjectID:
			q.Set("object_id", id.String())
		case miResourceID:
			q.Set(msiResID, id.String())
		}
	}
	request.Raw().URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) createAppServiceAuthRequest(ctx context.Context, id ManagedIDKind, scopes []string) (*policy.Request, error) {
	request, err := azruntime.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Raw().Header.Set("X-IDENTITY-HEADER", os.Getenv(identityHeader))
	q := request.Raw().URL.Query()
	q.Set("api-version", "2019-08-01")
	q.Set("resource", scopes[0])
	if id != nil {
		switch id.idKind() {
		case miClientID:
			q.Set(qpClientID, id.String())
		case miObjectID:
			q.Set("principal_id", id.String())
		case miResourceID:
			q.Set(miResID, id.String())
		}
	}
	request.Raw().URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) createAzureMLAuthRequest(ctx context.Context, id ManagedIDKind, scopes []string) (*policy.Request, error) {
	request, err := azruntime.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Raw().Header.Set("secret", os.Getenv(msiSecret))
	q := request.Raw().URL.Query()
	q.Set("api-version", "2017-09-01")
	q.Set("resource", strings.Join(scopes, " "))
	q.Set("clientid", os.Getenv(defaultIdentityClientID))
	if id != nil {
		switch id.idKind() {
		case miClientID:
			q.Set("clientid", id.String())
		case miObjectID:
			return nil, newAuthenticationFailedError(credNameManagedIdentity, "Azure ML doesn't support specifying a managed identity by object ID", nil)
		case miResourceID:
			return nil, newAuthenticationFailedError(credNameManagedIdentity, "Azure ML doesn't support specifying a managed identity by resource ID", nil)
		}
	}
	request.Raw().URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) createServiceFabricAuthRequest(ctx context.Context, scopes []string) (*policy.Request, error) {
	request, err := azruntime.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	q := request.Raw().URL.Query()
	request.Raw().Header.Set("Accept", "application/json")
	request.Raw().Header.Set("Secret", os.Getenv(identityHeader))
	q.Set("api-version", serviceFabricAPIVersion)
	q.Set("resource", strings.Join(scopes, " "))
	request.Raw().URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) getAzureArcSecretKey(ctx context.Context, resources []string) (string, error) {
	// create the request to retreive the secret key challenge provided by the HIMDS service
	request, err := azruntime.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return "", err
	}
	request.Raw().Header.Set(headerMetadata, "true")
	q := request.Raw().URL.Query()
	q.Set("api-version", azureArcAPIVersion)
	q.Set("resource", strings.Join(resources, " "))
	request.Raw().URL.RawQuery = q.Encode()
	// send the initial request to get the short-lived secret key
	response, err := c.azClient.Pipeline().Do(request)
	if err != nil {
		return "", err
	}
	// the endpoint is expected to return a 401 with the WWW-Authenticate header set to the location
	// of the secret key file. Any other status code indicates an error in the request.
	if response.StatusCode != 401 {
		msg := fmt.Sprintf("expected a 401 response, received %d", response.StatusCode)
		return "", newAuthenticationFailedError(credNameManagedIdentity, msg, response)
	}
	header := response.Header.Get("WWW-Authenticate")
	if len(header) == 0 {
		return "", newAuthenticationFailedError(credNameManagedIdentity, "HIMDS response has no WWW-Authenticate header", nil)
	}
	// the WWW-Authenticate header is expected in the following format: Basic realm=/some/file/path.key
	_, p, found := strings.Cut(header, "=")
	if !found {
		return "", newAuthenticationFailedError(credNameManagedIdentity, "unexpected WWW-Authenticate header from HIMDS: "+header, nil)
	}
	expected, err := arcKeyDirectory()
	if err != nil {
		return "", err
	}
	if filepath.Dir(p) != expected || !strings.HasSuffix(p, ".key") {
		return "", newAuthenticationFailedError(credNameManagedIdentity, "unexpected file path from HIMDS service: "+p, nil)
	}
	f, err := os.Stat(p)
	if err != nil {
		return "", newAuthenticationFailedError(credNameManagedIdentity, fmt.Sprintf("could not stat %q: %v", p, err), nil)
	}
	if s := f.Size(); s > 4096 {
		return "", newAuthenticationFailedError(credNameManagedIdentity, fmt.Sprintf("key is too large (%d bytes)", s), nil)
	}
	key, err := os.ReadFile(p)
	if err != nil {
		return "", newAuthenticationFailedError(credNameManagedIdentity, fmt.Sprintf("could not read %q: %v", p, err), nil)
	}
	return string(key), nil
}

func (c *managedIdentityClient) createAzureArcAuthRequest(ctx context.Context, resources []string, key string) (*policy.Request, error) {
	request, err := azruntime.NewRequest(ctx, http.MethodGet, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Raw().Header.Set(headerMetadata, "true")
	request.Raw().Header.Set("Authorization", fmt.Sprintf("Basic %s", key))
	q := request.Raw().URL.Query()
	q.Set("api-version", azureArcAPIVersion)
	q.Set("resource", strings.Join(resources, " "))
	request.Raw().URL.RawQuery = q.Encode()
	return request, nil
}

func (c *managedIdentityClient) createCloudShellAuthRequest(ctx context.Context, scopes []string) (*policy.Request, error) {
	request, err := azruntime.NewRequest(ctx, http.MethodPost, c.endpoint)
	if err != nil {
		return nil, err
	}
	request.Raw().Header.Set(headerMetadata, "true")
	data := url.Values{}
	data.Set("resource", strings.Join(scopes, " "))
	dataEncoded := data.Encode()
	body := streaming.NopCloser(strings.NewReader(dataEncoded))
	if err := request.SetBody(body, "application/x-www-form-urlencoded"); err != nil {
		return nil, err
	}
	return request, nil
}
