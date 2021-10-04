# Release History

## 0.12.0 (Unreleased)
### Breaking Changes
* Removed `NewAuthenticationPolicy()` from credentials. Clients should instead use azcore's
 `runtime.NewBearerTokenPolicy()` to construct a bearer token authorization policy.
* The `AuthorityHost` field in credential options structs is now a custom type,
  `AuthorityHost`, with underlying type `string`
* `NewChainedTokenCredential` has a new signature to accommodate a placeholder
  options struct:
  ```go
  // before
  cred, err := NewChainedTokenCredential(credA, credB)

  // after
  cred, err := NewChainedTokenCredential([]azcore.TokenCredential{credA, credB}, nil)
  ```
* Removed `ExcludeAzureCLICredential`, `ExcludeEnvironmentCredential`, and `ExcludeMSICredential`
  from `DefaultAzureCredentialOptions`
* `NewClientCertificateCredential` requires the bytes of a certificate instead of
  a path to a certificate file:
  ```go
  // before
  cred, err := NewClientCertificateCredential("tenant", "client-id", "/cert.pem", nil)

  // after
  certData, err := os.ReadFile("/cert.pem")
  cred, err := NewClientCertificateCredential("tenant", "client-id", certData, nil)
  ```
* Removed `InteractiveBrowserCredentialOptions.ClientSecret` and `.Port`
* Removed `AADAuthenticationFailedError`
* Removed `id` parameter of `NewManagedIdentityCredential()`. User assigned identities are now
  specified by `ManagedIdentityCredentialOptions.ID`:
  ```go
  // before
  cred, err := NewManagedIdentityCredential("client-id", nil)
  // or, for a resource ID
  opts := &ManagedIdentityCredentialOptions{ID: ResourceID}
  cred, err := NewManagedIdentityCredential("/subscriptions/...", opts)

  // after
  clientID := ClientID("7cf7db0d-...")
  opts := &ManagedIdentityCredentialOptions{ID: clientID}
  // or, for a resource ID
  resID: ResourceID("/subscriptions/...")
  opts := &ManagedIdentityCredentialOptions{ID: resID}
  cred, err := NewManagedIdentityCredential(opts)
  ```

### Features Added
* Added connection configuration options to `DefaultAzureCredentialOptions`
* `AuthenticationFailedError.RawResponse()` returns the HTTP response motivating the error,
  if available


## 0.11.0 (2021-09-08)
### Breaking Changes
* Unexported `AzureCLICredentialOptions.TokenProvider` and its type,
  `AzureCLITokenProvider`

### Bug Fixes
* `ManagedIdentityCredential.GetToken` returns `CredentialUnavailableError`
  when IMDS has no assigned identity, signaling `DefaultAzureCredential` to
  try other credentials


## 0.10.0 (2021-08-30)
### Breaking Changes
* Update based on `azcore` refactor [#15383](https://github.com/Azure/azure-sdk-for-go/pull/15383)

## 0.9.3 (2021-08-20)

### Bugs Fixed
* `ManagedIdentityCredential.GetToken` no longer mutates its `opts.Scopes`

### Other Changes
* Bumps version of `azcore` to `v0.18.1`


## 0.9.2 (2021-07-23)
### Features Added
* Adding support for Service Fabric environment in `ManagedIdentityCredential`
* Adding an option for using a resource ID instead of client ID in `ManagedIdentityCredential`


## 0.9.1 (2021-05-24)
### Features Added
* Add LICENSE.txt and bump version information


## 0.9.0 (2021-05-21)
### Features Added
* Add support for authenticating in Azure Stack environments
* Enable user assigned identities for the IMDS scenario in `ManagedIdentityCredential`
* Add scope to resource conversion in `GetToken()` on `ManagedIdentityCredential`


## 0.8.0 (2021-01-20)
### Features Added
* Updating documentation


## 0.7.1 (2021-01-04)
### Features Added
* Adding port option to `InteractiveBrowserCredential`


## 0.7.0 (2020-12-11)
### Features Added
* Add `redirectURI` parameter back to authentication code flow


## 0.6.1 (2020-12-09)
### Features Added
* Updating query parameter in `ManagedIdentityCredential` and updating datetime string for parsing managed identity access tokens.


## 0.6.0 (2020-11-16)
### Features Added
* Remove `RedirectURL` parameter from auth code flow to align with the MSAL implementation which relies on the native client redirect URL.


## 0.5.0 (2020-10-30)
### Features Added
* Flattening credential options


## 0.4.3 (2020-10-21)
### Features Added
* Adding Azure Arc support in `ManagedIdentityCredential`


## 0.4.2 (2020-10-16)
### Features Added
* Typo fixes


## 0.4.1 (2020-10-16)
### Features Added
* Ensure authority hosts are only HTTPs


## 0.4.0 (2020-10-16)
### Features Added
* Adding options structs for credentials


## 0.3.0 (2020-10-09)
### Features Added
* Update `DeviceCodeCredential` callback


## 0.2.2 (2020-10-09)
### Features Added
* Add `AuthorizationCodeCredential`


## 0.2.1 (2020-10-06)
### Features Added
* Add `InteractiveBrowserCredential`


## 0.2.0 (2020-09-11)
### Features Added
* Refactor `azidentity` on top of `azcore` refactor
* Updated policies to conform to `policy.Policy` interface changes.
* Updated non-retriable errors to conform to `azcore.NonRetriableError`.
* Fixed calls to `Request.SetBody()` to include content type.
* Switched endpoints to string types and removed extra parsing code.


## 0.1.1 (2020-09-02)
### Features Added
* Add `AzureCLICredential` to `DefaultAzureCredential` chain


## 0.1.0 (2020-07-23)
### Features Added
* Initial Release. Azure Identity library that provides Azure Active Directory token authentication support for the SDK.
