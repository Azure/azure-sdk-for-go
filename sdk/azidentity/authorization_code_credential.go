package azidentity

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
)

const credNameAuthCode = "AuthorizationCodeCredential"

type AuthorizationCodeCredentialOptions struct {
	azcore.ClientOptions

	TenantID     string
	ClientID     string
	ClientSecret string
	AuthCode     string
	RedirectURI  string
}

type AuthorizationCodeCredential struct {
	options AuthorizationCodeCredentialOptions
	client  confidentialClient
	account confidential.Account
	syncer  *syncer
}

func NewAuthorizationCodeCredential(options *AuthorizationCodeCredentialOptions) (*AuthorizationCodeCredential, error) {
	if options == nil {
		options = &AuthorizationCodeCredentialOptions{}
	}

	cred, err := confidential.NewCredFromSecret(options.ClientSecret)
	if err != nil {
		return nil, err
	}

	client, err := getConfidentialClient(options.ClientID, options.TenantID, cred, &options.ClientOptions)
	if err != nil {
		return nil, err
	}

	credential := &AuthorizationCodeCredential{
		options: *options,
		client:  client,
	}
	credential.syncer = newSyncer(credNameAuthCode, options.TenantID, []string{}, credential.requestToken, credential.silentAuth)

	return credential, nil
}

func (c *AuthorizationCodeCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return c.syncer.GetToken(ctx, opts)
}

func (c *AuthorizationCodeCredential) silentAuth(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	res, err := c.client.AcquireTokenSilent(ctx, opts.Scopes, confidential.WithSilentAccount(c.account))
	return azcore.AccessToken{Token: res.AccessToken, ExpiresOn: res.ExpiresOn.UTC()}, err
}

func (c *AuthorizationCodeCredential) requestToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	res, err := c.client.AcquireTokenByAuthCode(ctx, c.options.AuthCode, c.options.RedirectURI, opts.Scopes)
	c.account = res.Account
	return azcore.AccessToken{Token: res.AccessToken, ExpiresOn: res.ExpiresOn.UTC()}, err
}

var _ azcore.TokenCredential = (*AuthorizationCodeCredential)(nil)
