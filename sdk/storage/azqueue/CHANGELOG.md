## Release History

### 2.0.0-beta.2 (Unreleased)

#### Features Added

#### Breaking Changes

#### Bugs Fixed

#### Other Changes

### 2.0.0-beta.1 (2025-07-08)

#### Breaking Changes
* The type of ApproximateMessagesCount has changed from *int32 to *int64 to support large queues with more than int32 max value messages.

### Other Changes
* Updated `azidentity` version to `1.10.1`

### 1.0.1 (2025-04-30)

#### Features Added
* Updated `azidentity` version to `1.9.0`
* Updated `azcore` version to `1.18.0`
* Update transitive dependency `github.com/golang-jwt/jwt`, addressing security vulnerability [CVE-2025-30204](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2025-30204).

#### Bugs Fixed
* Fixed issue where some requests fail with mismatch in string to sign.
* Fixed service SAS creation where expiry time or permissions can be omitted when stored access policy is used.

#### Other Changes
* Integrate `InsecureAllowCredentialWithHTTP` client options.
* Update dependencies.

### 1.0.0 (2023-05-09)

### Features Added

* This is the initial GA release of the `azqueue` library


### 0.1.0 (2023-02-15)

### Features Added

* This is the initial preview release of the `azqueue` library
