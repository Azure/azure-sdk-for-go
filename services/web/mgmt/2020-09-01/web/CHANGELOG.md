# Unreleased

## Breaking Changes

### Removed Funcs

1. BaseClient.GenerateGithubAccessTokenForAppserviceCLIAsync(context.Context, AppserviceGithubTokenRequest) (AppserviceGithubToken, error)
1. BaseClient.GenerateGithubAccessTokenForAppserviceCLIAsyncPreparer(context.Context, AppserviceGithubTokenRequest) (*http.Request, error)
1. BaseClient.GenerateGithubAccessTokenForAppserviceCLIAsyncResponder(*http.Response) (AppserviceGithubToken, error)
1. BaseClient.GenerateGithubAccessTokenForAppserviceCLIAsyncSender(*http.Request) (*http.Response, error)

### Struct Changes

#### Removed Struct Fields

1. AppserviceGithubToken.autorest.Response

## Additive Changes

### Struct Changes

#### New Struct Fields

1. SiteConfig.AcrUseManagedIdentityCreds
1. SiteConfig.AcrUserManagedIdentityID
