# Release History

## v0.10.0-beta.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes


## v0.9.2
### Features Added
* Adding support for Service Fabric environment in `ManagedIdentityCredential`
* Adding an option for using a resource ID instead of client ID in `ManagedIdentityCredential`


## v0.9.1
### Features Added
* Add LICENSE.txt and bump version information


## v0.9.0
### Features Added
* Add support for authenticating in Azure Stack environments
* Enable user assigned identities for the IMDS scenario in `ManagedIdentityCredential`
* Add scope to resource conversion in `GetToken()` on `ManagedIdentityCredential`


## v0.8.0
### Features Added
* Updating documentation


## v0.7.1
### Features Added
* Adding port option to `InteractiveBrowserCredential`


## v0.7.0
### Features Added
* Add `redirectURI` parameter back to authentication code flow


## v0.6.1
### Features Added
* Updating query parameter in `ManagedIdentityCredential` and updating datetime string for parsing managed identity access tokens.


## v0.6.0
### Features Added
* Remove `RedirectURL` parameter from auth code flow to align with the MSAL implementation which relies on the native client redirect URL.


## v0.5.0
### Features Added
* Flattening credential options


## v0.4.3
### Features Added
* Adding Azure Arc support in `ManagedIdentityCredential`


## v0.4.2
### Features Added
* Typo fixes


## v0.4.1
### Features Added
* Ensure authority hosts are only HTTPs


## v0.4.0
### Features Added
* Adding options structs for credentials


## v0.3.0
### Features Added
* Update `DeviceCodeCredential` callback


## v0.2.2
### Features Added
* Add `AuthorizationCodeCredential`


## v0.2.1
### Features Added
* Add `InteractiveBrowserCredential`


## v0.2.0
### Features Added
* Refactor `azidentity` on top of `azcore` refactor
* Updated policies to conform to `azcore.Policy` interface changes.
* Updated non-retriable errors to conform to `azcore.NonRetriableError`.
* Fixed calls to `Request.SetBody()` to include content type.
* Switched endpoints to string types and removed extra parsing code.


## v0.1.1
### Features Added
* Add `AzureCLICredential` to `DefaultAzureCredential` chain


## v0.1.0
### Features Added
* Initial Release. Azure Identity library that provides Azure Active Directory token authentication support for the SDK.
