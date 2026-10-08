# Release History

## 0.2.0 (Unreleased)

### Features Added
* `autorefresh` package: a cache for a single expiring value with background refresh, shared acquisitions and conditional invalidation, used by STG105 session authentication and data locality
* `locality` package: the pipeline policy that routes a request to an STG105 layout endpoint while keeping the account Host header

### Other Changes
* Added a dependency on `azcore`

## 0.1.0 (2026-10-01)

### Features Added
* Initial release of the shared internal module for Azure Storage
* Structured message (XSM/1.0) encoder and decoder, consolidated from azblob, azfile, and azdatalake
* `sasurl` package for appending a SAS query string to a resource or account URL that may already carry a query string
