# Unreleased

## Breaking Changes

### Removed Funcs

1. *AllowedAudiencesValidation.UnmarshalJSON([]byte) error
1. *AppRegistration.UnmarshalJSON([]byte) error
1. *Apple.UnmarshalJSON([]byte) error
1. *AppleRegistration.UnmarshalJSON([]byte) error
1. *AuthPlatform.UnmarshalJSON([]byte) error
1. *AzureActiveDirectory.UnmarshalJSON([]byte) error
1. *AzureActiveDirectoryLogin.UnmarshalJSON([]byte) error
1. *AzureActiveDirectoryRegistration.UnmarshalJSON([]byte) error
1. *AzureActiveDirectoryValidation.UnmarshalJSON([]byte) error
1. *AzureStaticWebApps.UnmarshalJSON([]byte) error
1. *AzureStaticWebAppsRegistration.UnmarshalJSON([]byte) error
1. *BlobStorageTokenStore.UnmarshalJSON([]byte) error
1. *ClientRegistration.UnmarshalJSON([]byte) error
1. *CookieExpiration.UnmarshalJSON([]byte) error
1. *CustomOpenIDConnectProvider.UnmarshalJSON([]byte) error
1. *Facebook.UnmarshalJSON([]byte) error
1. *FileSystemTokenStore.UnmarshalJSON([]byte) error
1. *ForwardProxy.UnmarshalJSON([]byte) error
1. *GitHub.UnmarshalJSON([]byte) error
1. *GlobalValidation.UnmarshalJSON([]byte) error
1. *Google.UnmarshalJSON([]byte) error
1. *HTTPSettings.UnmarshalJSON([]byte) error
1. *HTTPSettingsRoutes.UnmarshalJSON([]byte) error
1. *IdentityProviders.UnmarshalJSON([]byte) error
1. *JwtClaimChecks.UnmarshalJSON([]byte) error
1. *LegacyMicrosoftAccount.UnmarshalJSON([]byte) error
1. *Login.UnmarshalJSON([]byte) error
1. *LoginRoutes.UnmarshalJSON([]byte) error
1. *LoginScopes.UnmarshalJSON([]byte) error
1. *Nonce.UnmarshalJSON([]byte) error
1. *OpenIDConnectClientCredential.UnmarshalJSON([]byte) error
1. *OpenIDConnectConfig.UnmarshalJSON([]byte) error
1. *OpenIDConnectLogin.UnmarshalJSON([]byte) error
1. *OpenIDConnectRegistration.UnmarshalJSON([]byte) error
1. *TokenStore.UnmarshalJSON([]byte) error
1. *Twitter.UnmarshalJSON([]byte) error
1. *TwitterRegistration.UnmarshalJSON([]byte) error
1. AllowedAudiencesValidation.MarshalJSON() ([]byte, error)
1. AppRegistration.MarshalJSON() ([]byte, error)
1. Apple.MarshalJSON() ([]byte, error)
1. AppleRegistration.MarshalJSON() ([]byte, error)
1. AuthPlatform.MarshalJSON() ([]byte, error)
1. AzureActiveDirectory.MarshalJSON() ([]byte, error)
1. AzureActiveDirectoryLogin.MarshalJSON() ([]byte, error)
1. AzureActiveDirectoryRegistration.MarshalJSON() ([]byte, error)
1. AzureActiveDirectoryValidation.MarshalJSON() ([]byte, error)
1. AzureStaticWebApps.MarshalJSON() ([]byte, error)
1. AzureStaticWebAppsRegistration.MarshalJSON() ([]byte, error)
1. BaseClient.GenerateGithubAccessTokenForAppserviceCLIAsync(context.Context, AppserviceGithubTokenRequest) (AppserviceGithubToken, error)
1. BaseClient.GenerateGithubAccessTokenForAppserviceCLIAsyncPreparer(context.Context, AppserviceGithubTokenRequest) (*http.Request, error)
1. BaseClient.GenerateGithubAccessTokenForAppserviceCLIAsyncResponder(*http.Response) (AppserviceGithubToken, error)
1. BaseClient.GenerateGithubAccessTokenForAppserviceCLIAsyncSender(*http.Request) (*http.Response, error)
1. BlobStorageTokenStore.MarshalJSON() ([]byte, error)
1. ClientRegistration.MarshalJSON() ([]byte, error)
1. CookieExpiration.MarshalJSON() ([]byte, error)
1. CustomOpenIDConnectProvider.MarshalJSON() ([]byte, error)
1. Facebook.MarshalJSON() ([]byte, error)
1. FileSystemTokenStore.MarshalJSON() ([]byte, error)
1. ForwardProxy.MarshalJSON() ([]byte, error)
1. GitHub.MarshalJSON() ([]byte, error)
1. GlobalValidation.MarshalJSON() ([]byte, error)
1. Google.MarshalJSON() ([]byte, error)
1. HTTPSettings.MarshalJSON() ([]byte, error)
1. HTTPSettingsRoutes.MarshalJSON() ([]byte, error)
1. IdentityProvidersProperties.MarshalJSON() ([]byte, error)
1. JwtClaimChecks.MarshalJSON() ([]byte, error)
1. LegacyMicrosoftAccount.MarshalJSON() ([]byte, error)
1. Login.MarshalJSON() ([]byte, error)
1. LoginRoutes.MarshalJSON() ([]byte, error)
1. LoginScopes.MarshalJSON() ([]byte, error)
1. Nonce.MarshalJSON() ([]byte, error)
1. OpenIDConnectClientCredential.MarshalJSON() ([]byte, error)
1. OpenIDConnectConfig.MarshalJSON() ([]byte, error)
1. OpenIDConnectLogin.MarshalJSON() ([]byte, error)
1. OpenIDConnectRegistration.MarshalJSON() ([]byte, error)
1. TokenStore.MarshalJSON() ([]byte, error)
1. Twitter.MarshalJSON() ([]byte, error)
1. TwitterRegistration.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. AllowedAudiencesValidationProperties
1. AppRegistrationProperties
1. AppleProperties
1. AppleRegistrationProperties
1. AuthPlatformProperties
1. AzureActiveDirectoryLoginProperties
1. AzureActiveDirectoryProperties
1. AzureActiveDirectoryRegistrationProperties
1. AzureActiveDirectoryValidationProperties
1. AzureStaticWebAppsProperties
1. AzureStaticWebAppsRegistrationProperties
1. BlobStorageTokenStoreProperties
1. ClientRegistrationProperties
1. CookieExpirationProperties
1. CustomOpenIDConnectProviderProperties
1. FacebookProperties
1. FileSystemTokenStoreProperties
1. ForwardProxyProperties
1. GitHubProperties
1. GlobalValidationProperties
1. GoogleProperties
1. HTTPSettingsProperties
1. HTTPSettingsRoutesProperties
1. IdentityProvidersProperties
1. JwtClaimChecksProperties
1. LegacyMicrosoftAccountProperties
1. LoginProperties
1. LoginRoutesProperties
1. LoginScopesProperties
1. NonceProperties
1. OpenIDConnectClientCredentialProperties
1. OpenIDConnectConfigProperties
1. OpenIDConnectLoginProperties
1. OpenIDConnectRegistrationProperties
1. TokenStoreProperties
1. TwitterProperties
1. TwitterRegistrationProperties

#### Removed Struct Fields

1. AllowedAudiencesValidation.*AllowedAudiencesValidationProperties
1. AllowedAudiencesValidation.ID
1. AllowedAudiencesValidation.Kind
1. AllowedAudiencesValidation.Name
1. AllowedAudiencesValidation.Type
1. AppRegistration.*AppRegistrationProperties
1. AppRegistration.ID
1. AppRegistration.Kind
1. AppRegistration.Name
1. AppRegistration.Type
1. Apple.*AppleProperties
1. Apple.ID
1. Apple.Kind
1. Apple.Name
1. Apple.Type
1. AppleRegistration.*AppleRegistrationProperties
1. AppleRegistration.ID
1. AppleRegistration.Kind
1. AppleRegistration.Name
1. AppleRegistration.Type
1. AppserviceGithubToken.autorest.Response
1. AuthPlatform.*AuthPlatformProperties
1. AuthPlatform.ID
1. AuthPlatform.Kind
1. AuthPlatform.Name
1. AuthPlatform.Type
1. AzureActiveDirectory.*AzureActiveDirectoryProperties
1. AzureActiveDirectory.ID
1. AzureActiveDirectory.Kind
1. AzureActiveDirectory.Name
1. AzureActiveDirectory.Type
1. AzureActiveDirectoryLogin.*AzureActiveDirectoryLoginProperties
1. AzureActiveDirectoryLogin.ID
1. AzureActiveDirectoryLogin.Kind
1. AzureActiveDirectoryLogin.Name
1. AzureActiveDirectoryLogin.Type
1. AzureActiveDirectoryRegistration.*AzureActiveDirectoryRegistrationProperties
1. AzureActiveDirectoryRegistration.ID
1. AzureActiveDirectoryRegistration.Kind
1. AzureActiveDirectoryRegistration.Name
1. AzureActiveDirectoryRegistration.Type
1. AzureActiveDirectoryValidation.*AzureActiveDirectoryValidationProperties
1. AzureActiveDirectoryValidation.ID
1. AzureActiveDirectoryValidation.Kind
1. AzureActiveDirectoryValidation.Name
1. AzureActiveDirectoryValidation.Type
1. AzureStaticWebApps.*AzureStaticWebAppsProperties
1. AzureStaticWebApps.ID
1. AzureStaticWebApps.Kind
1. AzureStaticWebApps.Name
1. AzureStaticWebApps.Type
1. AzureStaticWebAppsRegistration.*AzureStaticWebAppsRegistrationProperties
1. AzureStaticWebAppsRegistration.ID
1. AzureStaticWebAppsRegistration.Kind
1. AzureStaticWebAppsRegistration.Name
1. AzureStaticWebAppsRegistration.Type
1. BlobStorageTokenStore.*BlobStorageTokenStoreProperties
1. BlobStorageTokenStore.ID
1. BlobStorageTokenStore.Kind
1. BlobStorageTokenStore.Name
1. BlobStorageTokenStore.Type
1. ClientRegistration.*ClientRegistrationProperties
1. ClientRegistration.ID
1. ClientRegistration.Kind
1. ClientRegistration.Name
1. ClientRegistration.Type
1. CookieExpiration.*CookieExpirationProperties
1. CookieExpiration.ID
1. CookieExpiration.Kind
1. CookieExpiration.Name
1. CookieExpiration.Type
1. CustomOpenIDConnectProvider.*CustomOpenIDConnectProviderProperties
1. CustomOpenIDConnectProvider.ID
1. CustomOpenIDConnectProvider.Kind
1. CustomOpenIDConnectProvider.Name
1. CustomOpenIDConnectProvider.Type
1. Facebook.*FacebookProperties
1. Facebook.ID
1. Facebook.Kind
1. Facebook.Name
1. Facebook.Type
1. FileSystemTokenStore.*FileSystemTokenStoreProperties
1. FileSystemTokenStore.ID
1. FileSystemTokenStore.Kind
1. FileSystemTokenStore.Name
1. FileSystemTokenStore.Type
1. ForwardProxy.*ForwardProxyProperties
1. ForwardProxy.ID
1. ForwardProxy.Kind
1. ForwardProxy.Name
1. ForwardProxy.Type
1. GitHub.*GitHubProperties
1. GitHub.ID
1. GitHub.Kind
1. GitHub.Name
1. GitHub.Type
1. GlobalValidation.*GlobalValidationProperties
1. GlobalValidation.ID
1. GlobalValidation.Kind
1. GlobalValidation.Name
1. GlobalValidation.Type
1. Google.*GoogleProperties
1. Google.ID
1. Google.Kind
1. Google.Name
1. Google.Type
1. HTTPSettings.*HTTPSettingsProperties
1. HTTPSettings.ID
1. HTTPSettings.Kind
1. HTTPSettings.Name
1. HTTPSettings.Type
1. HTTPSettingsRoutes.*HTTPSettingsRoutesProperties
1. HTTPSettingsRoutes.ID
1. HTTPSettingsRoutes.Kind
1. HTTPSettingsRoutes.Name
1. HTTPSettingsRoutes.Type
1. IdentityProviders.*IdentityProvidersProperties
1. IdentityProviders.ID
1. IdentityProviders.Kind
1. IdentityProviders.Name
1. IdentityProviders.Type
1. JwtClaimChecks.*JwtClaimChecksProperties
1. JwtClaimChecks.ID
1. JwtClaimChecks.Kind
1. JwtClaimChecks.Name
1. JwtClaimChecks.Type
1. LegacyMicrosoftAccount.*LegacyMicrosoftAccountProperties
1. LegacyMicrosoftAccount.ID
1. LegacyMicrosoftAccount.Kind
1. LegacyMicrosoftAccount.Name
1. LegacyMicrosoftAccount.Type
1. Login.*LoginProperties
1. Login.ID
1. Login.Kind
1. Login.Name
1. Login.Type
1. LoginRoutes.*LoginRoutesProperties
1. LoginRoutes.ID
1. LoginRoutes.Kind
1. LoginRoutes.Name
1. LoginRoutes.Type
1. LoginScopes.*LoginScopesProperties
1. LoginScopes.ID
1. LoginScopes.Kind
1. LoginScopes.Name
1. LoginScopes.Type
1. Nonce.*NonceProperties
1. Nonce.ID
1. Nonce.Kind
1. Nonce.Name
1. Nonce.Type
1. OpenIDConnectClientCredential.*OpenIDConnectClientCredentialProperties
1. OpenIDConnectClientCredential.ID
1. OpenIDConnectClientCredential.Kind
1. OpenIDConnectClientCredential.Name
1. OpenIDConnectClientCredential.Type
1. OpenIDConnectConfig.*OpenIDConnectConfigProperties
1. OpenIDConnectConfig.ID
1. OpenIDConnectConfig.Kind
1. OpenIDConnectConfig.Name
1. OpenIDConnectConfig.Type
1. OpenIDConnectLogin.*OpenIDConnectLoginProperties
1. OpenIDConnectLogin.ID
1. OpenIDConnectLogin.Kind
1. OpenIDConnectLogin.Name
1. OpenIDConnectLogin.Type
1. OpenIDConnectRegistration.*OpenIDConnectRegistrationProperties
1. OpenIDConnectRegistration.ID
1. OpenIDConnectRegistration.Kind
1. OpenIDConnectRegistration.Name
1. OpenIDConnectRegistration.Type
1. TokenStore.*TokenStoreProperties
1. TokenStore.ID
1. TokenStore.Kind
1. TokenStore.Name
1. TokenStore.Type
1. Twitter.*TwitterProperties
1. Twitter.ID
1. Twitter.Kind
1. Twitter.Name
1. Twitter.Type
1. TwitterRegistration.*TwitterRegistrationProperties
1. TwitterRegistration.ID
1. TwitterRegistration.Kind
1. TwitterRegistration.Name
1. TwitterRegistration.Type

## Additive Changes

### Struct Changes

#### New Struct Fields

1. AllowedAudiencesValidation.AllowedAudiences
1. AppRegistration.AppID
1. AppRegistration.AppSecretSettingName
1. Apple.Enabled
1. Apple.Login
1. Apple.Registration
1. AppleRegistration.ClientID
1. AppleRegistration.ClientSecretSettingName
1. AuthPlatform.ConfigFilePath
1. AuthPlatform.Enabled
1. AuthPlatform.RuntimeVersion
1. AzureActiveDirectory.Enabled
1. AzureActiveDirectory.IsAutoProvisioned
1. AzureActiveDirectory.Login
1. AzureActiveDirectory.Registration
1. AzureActiveDirectory.Validation
1. AzureActiveDirectoryLogin.DisableWWWAuthenticate
1. AzureActiveDirectoryLogin.LoginParameters
1. AzureActiveDirectoryRegistration.ClientID
1. AzureActiveDirectoryRegistration.ClientSecretCertificateIssuer
1. AzureActiveDirectoryRegistration.ClientSecretCertificateSubjectAlternativeName
1. AzureActiveDirectoryRegistration.ClientSecretCertificateThumbprint
1. AzureActiveDirectoryRegistration.ClientSecretSettingName
1. AzureActiveDirectoryRegistration.OpenIDIssuer
1. AzureActiveDirectoryValidation.AllowedAudiences
1. AzureActiveDirectoryValidation.JwtClaimChecks
1. AzureStaticWebApps.Enabled
1. AzureStaticWebApps.Registration
1. AzureStaticWebAppsRegistration.ClientID
1. BlobStorageTokenStore.SasURLSettingName
1. ClientRegistration.ClientID
1. ClientRegistration.ClientSecretSettingName
1. CookieExpiration.Convention
1. CookieExpiration.TimeToExpiration
1. CustomOpenIDConnectProvider.Enabled
1. CustomOpenIDConnectProvider.Login
1. CustomOpenIDConnectProvider.Registration
1. Facebook.Enabled
1. Facebook.GraphAPIVersion
1. Facebook.Login
1. Facebook.Registration
1. FileSystemTokenStore.Directory
1. ForwardProxy.Convention
1. ForwardProxy.CustomHostHeaderName
1. ForwardProxy.CustomProtoHeaderName
1. GitHub.Enabled
1. GitHub.Login
1. GitHub.Registration
1. GlobalValidation.ExcludedPaths
1. GlobalValidation.RedirectToProvider
1. GlobalValidation.RequireAuthentication
1. GlobalValidation.UnauthenticatedClientAction
1. Google.Enabled
1. Google.Login
1. Google.Registration
1. Google.Validation
1. HTTPSettings.ForwardProxy
1. HTTPSettings.RequireHTTPS
1. HTTPSettings.Routes
1. HTTPSettingsRoutes.APIPrefix
1. IdentityProviders.Apple
1. IdentityProviders.AzureActiveDirectory
1. IdentityProviders.AzureStaticWebApps
1. IdentityProviders.CustomOpenIDConnectProviders
1. IdentityProviders.Facebook
1. IdentityProviders.GitHub
1. IdentityProviders.Google
1. IdentityProviders.LegacyMicrosoftAccount
1. IdentityProviders.Twitter
1. JwtClaimChecks.AllowedClientApplications
1. JwtClaimChecks.AllowedGroups
1. LegacyMicrosoftAccount.Enabled
1. LegacyMicrosoftAccount.Login
1. LegacyMicrosoftAccount.Registration
1. LegacyMicrosoftAccount.Validation
1. Login.AllowedExternalRedirectUrls
1. Login.CookieExpiration
1. Login.Nonce
1. Login.PreserveURLFragmentsForLogins
1. Login.Routes
1. Login.TokenStore
1. LoginRoutes.LogoutEndpoint
1. LoginScopes.Scopes
1. Nonce.NonceExpirationInterval
1. Nonce.ValidateNonce
1. OpenIDConnectClientCredential.ClientSecretSettingName
1. OpenIDConnectClientCredential.Method
1. OpenIDConnectConfig.AuthorizationEndpoint
1. OpenIDConnectConfig.CertificationURI
1. OpenIDConnectConfig.Issuer
1. OpenIDConnectConfig.TokenEndpoint
1. OpenIDConnectConfig.WellKnownOpenIDConfiguration
1. OpenIDConnectLogin.NameClaimType
1. OpenIDConnectLogin.Scopes
1. OpenIDConnectRegistration.ClientCredential
1. OpenIDConnectRegistration.ClientID
1. OpenIDConnectRegistration.OpenIDConnectConfiguration
1. TokenStore.AzureBlobStorage
1. TokenStore.Enabled
1. TokenStore.FileSystem
1. TokenStore.TokenRefreshExtensionHours
1. Twitter.Enabled
1. Twitter.Registration
1. TwitterRegistration.ConsumerKey
1. TwitterRegistration.ConsumerSecretSettingName
