# Unreleased

## Breaking Changes

### Removed Funcs

1. *AllowedPrincipals.UnmarshalJSON([]byte) error
1. *AppRegistration.UnmarshalJSON([]byte) error
1. *Apple.UnmarshalJSON([]byte) error
1. *AzureActiveDirectoryLogin.UnmarshalJSON([]byte) error
1. *AzureActiveDirectoryRegistration.UnmarshalJSON([]byte) error
1. *AzureActiveDirectoryValidation.UnmarshalJSON([]byte) error
1. *AzureStaticWebApps.UnmarshalJSON([]byte) error
1. *BlobStorageTokenStore.UnmarshalJSON([]byte) error
1. *CustomOpenIDConnectProvider.UnmarshalJSON([]byte) error
1. *GitHub.UnmarshalJSON([]byte) error
1. *Google.UnmarshalJSON([]byte) error
1. *LegacyMicrosoftAccount.UnmarshalJSON([]byte) error
1. *Twitter.UnmarshalJSON([]byte) error
1. AllowedPrincipals.MarshalJSON() ([]byte, error)
1. AppRegistration.MarshalJSON() ([]byte, error)
1. Apple.MarshalJSON() ([]byte, error)
1. AzureActiveDirectoryLogin.MarshalJSON() ([]byte, error)
1. AzureActiveDirectoryRegistration.MarshalJSON() ([]byte, error)
1. AzureActiveDirectoryValidation.MarshalJSON() ([]byte, error)
1. AzureStaticWebApps.MarshalJSON() ([]byte, error)
1. BaseClient.GenerateGithubAccessTokenForAppserviceCLIAsync(context.Context, AppserviceGithubTokenRequest) (AppserviceGithubToken, error)
1. BaseClient.GenerateGithubAccessTokenForAppserviceCLIAsyncPreparer(context.Context, AppserviceGithubTokenRequest) (*http.Request, error)
1. BaseClient.GenerateGithubAccessTokenForAppserviceCLIAsyncResponder(*http.Response) (AppserviceGithubToken, error)
1. BaseClient.GenerateGithubAccessTokenForAppserviceCLIAsyncSender(*http.Request) (*http.Response, error)
1. BlobStorageTokenStore.MarshalJSON() ([]byte, error)
1. CustomOpenIDConnectProvider.MarshalJSON() ([]byte, error)
1. GitHub.MarshalJSON() ([]byte, error)
1. Google.MarshalJSON() ([]byte, error)
1. LegacyMicrosoftAccount.MarshalJSON() ([]byte, error)
1. Twitter.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. AllowedPrincipalsProperties
1. AppRegistrationProperties
1. AppleProperties
1. AzureActiveDirectoryLoginProperties
1. AzureActiveDirectoryRegistrationProperties
1. AzureActiveDirectoryValidationProperties
1. AzureStaticWebAppsProperties
1. BlobStorageTokenStoreProperties
1. CustomOpenIDConnectProviderProperties
1. GitHubProperties
1. GoogleProperties
1. LegacyMicrosoftAccountProperties
1. TwitterProperties

#### Removed Struct Fields

1. AllowedPrincipals.*AllowedPrincipalsProperties
1. AllowedPrincipals.ID
1. AllowedPrincipals.Kind
1. AllowedPrincipals.Name
1. AllowedPrincipals.Type
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
1. AppserviceGithubToken.autorest.Response
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
1. BlobStorageTokenStore.*BlobStorageTokenStoreProperties
1. BlobStorageTokenStore.ID
1. BlobStorageTokenStore.Kind
1. BlobStorageTokenStore.Name
1. BlobStorageTokenStore.Type
1. CustomOpenIDConnectProvider.*CustomOpenIDConnectProviderProperties
1. CustomOpenIDConnectProvider.ID
1. CustomOpenIDConnectProvider.Kind
1. CustomOpenIDConnectProvider.Name
1. CustomOpenIDConnectProvider.Type
1. GitHub.*GitHubProperties
1. GitHub.ID
1. GitHub.Kind
1. GitHub.Name
1. GitHub.Type
1. Google.*GoogleProperties
1. Google.ID
1. Google.Kind
1. Google.Name
1. Google.Type
1. LegacyMicrosoftAccount.*LegacyMicrosoftAccountProperties
1. LegacyMicrosoftAccount.ID
1. LegacyMicrosoftAccount.Kind
1. LegacyMicrosoftAccount.Name
1. LegacyMicrosoftAccount.Type
1. Twitter.*TwitterProperties
1. Twitter.ID
1. Twitter.Kind
1. Twitter.Name
1. Twitter.Type

### Signature Changes

#### Funcs

1. ProviderClient.GetAvailableStacks
	- Params
		- From: context.Context, string
		- To: context.Context, ProviderOsTypeSelected
1. ProviderClient.GetAvailableStacksComplete
	- Params
		- From: context.Context, string
		- To: context.Context, ProviderOsTypeSelected
1. ProviderClient.GetAvailableStacksOnPrem
	- Params
		- From: context.Context, string
		- To: context.Context, ProviderOsTypeSelected
1. ProviderClient.GetAvailableStacksOnPremComplete
	- Params
		- From: context.Context, string
		- To: context.Context, ProviderOsTypeSelected
1. ProviderClient.GetAvailableStacksOnPremPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, ProviderOsTypeSelected
1. ProviderClient.GetAvailableStacksPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, ProviderOsTypeSelected
1. ProviderClient.GetFunctionAppStacks
	- Params
		- From: context.Context, string
		- To: context.Context, ProviderStackOsType
1. ProviderClient.GetFunctionAppStacksComplete
	- Params
		- From: context.Context, string
		- To: context.Context, ProviderStackOsType
1. ProviderClient.GetFunctionAppStacksForLocation
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, ProviderStackOsType
1. ProviderClient.GetFunctionAppStacksForLocationComplete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, ProviderStackOsType
1. ProviderClient.GetFunctionAppStacksForLocationPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, ProviderStackOsType
1. ProviderClient.GetFunctionAppStacksPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, ProviderStackOsType
1. ProviderClient.GetWebAppStacks
	- Params
		- From: context.Context, string
		- To: context.Context, ProviderStackOsType
1. ProviderClient.GetWebAppStacksComplete
	- Params
		- From: context.Context, string
		- To: context.Context, ProviderStackOsType
1. ProviderClient.GetWebAppStacksForLocation
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, ProviderStackOsType
1. ProviderClient.GetWebAppStacksForLocationComplete
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, ProviderStackOsType
1. ProviderClient.GetWebAppStacksForLocationPreparer
	- Params
		- From: context.Context, string, string
		- To: context.Context, string, ProviderStackOsType
1. ProviderClient.GetWebAppStacksPreparer
	- Params
		- From: context.Context, string
		- To: context.Context, ProviderStackOsType

## Additive Changes

### New Constants

1. ProviderOsTypeSelected.ProviderOsTypeSelectedAll
1. ProviderOsTypeSelected.ProviderOsTypeSelectedLinux
1. ProviderOsTypeSelected.ProviderOsTypeSelectedLinuxFunctions
1. ProviderOsTypeSelected.ProviderOsTypeSelectedWindows
1. ProviderOsTypeSelected.ProviderOsTypeSelectedWindowsFunctions
1. ProviderStackOsType.ProviderStackOsTypeAll
1. ProviderStackOsType.ProviderStackOsTypeLinux
1. ProviderStackOsType.ProviderStackOsTypeWindows

### New Funcs

1. PossibleProviderOsTypeSelectedValues() []ProviderOsTypeSelected
1. PossibleProviderStackOsTypeValues() []ProviderStackOsType

### Struct Changes

#### New Struct Fields

1. AllowedPrincipals.Groups
1. AllowedPrincipals.Identities
1. AppRegistration.AppID
1. AppRegistration.AppSecretSettingName
1. Apple.Enabled
1. Apple.Login
1. Apple.Registration
1. AzureActiveDirectoryLogin.DisableWWWAuthenticate
1. AzureActiveDirectoryLogin.LoginParameters
1. AzureActiveDirectoryRegistration.ClientID
1. AzureActiveDirectoryRegistration.ClientSecretCertificateIssuer
1. AzureActiveDirectoryRegistration.ClientSecretCertificateSubjectAlternativeName
1. AzureActiveDirectoryRegistration.ClientSecretCertificateThumbprint
1. AzureActiveDirectoryRegistration.ClientSecretSettingName
1. AzureActiveDirectoryRegistration.OpenIDIssuer
1. AzureActiveDirectoryValidation.AllowedAudiences
1. AzureActiveDirectoryValidation.DefaultAuthorizationPolicy
1. AzureActiveDirectoryValidation.JwtClaimChecks
1. AzureStaticWebApps.Enabled
1. AzureStaticWebApps.Registration
1. BlobStorageTokenStore.SasURLSettingName
1. CustomOpenIDConnectProvider.Enabled
1. CustomOpenIDConnectProvider.Login
1. CustomOpenIDConnectProvider.Registration
1. GitHub.Enabled
1. GitHub.Login
1. GitHub.Registration
1. Google.Enabled
1. Google.Login
1. Google.Registration
1. Google.Validation
1. LegacyMicrosoftAccount.Enabled
1. LegacyMicrosoftAccount.Login
1. LegacyMicrosoftAccount.Registration
1. LegacyMicrosoftAccount.Validation
1. Twitter.Enabled
1. Twitter.Registration
