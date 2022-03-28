# Unreleased

## Breaking Changes

### Signature Changes

#### Funcs

1. CustomDomainsClient.DisableCustomHTTPS
	- Returns
		- From: CustomDomain, error
		- To: CustomDomainsDisableCustomHTTPSFuture, error
1. CustomDomainsClient.DisableCustomHTTPSSender
	- Returns
		- From: *http.Response, error
		- To: CustomDomainsDisableCustomHTTPSFuture, error
1. CustomDomainsClient.EnableCustomHTTPS
	- Returns
		- From: CustomDomain, error
		- To: CustomDomainsEnableCustomHTTPSFuture, error
1. CustomDomainsClient.EnableCustomHTTPSSender
	- Returns
		- From: *http.Response, error
		- To: CustomDomainsEnableCustomHTTPSFuture, error

## Additive Changes

### New Constants

1. Transform.RemoveNulls
1. Transform.Trim
1. Transform.URLDecode
1. Transform.URLEncode

### New Funcs

1. *CustomDomainProperties.UnmarshalJSON([]byte) error
1. *CustomDomainsDisableCustomHTTPSFuture.UnmarshalJSON([]byte) error
1. *CustomDomainsEnableCustomHTTPSFuture.UnmarshalJSON([]byte) error

### Struct Changes

#### New Structs

1. CustomDomainsDisableCustomHTTPSFuture
1. CustomDomainsEnableCustomHTTPSFuture

#### New Struct Fields

1. CustomDomainProperties.CustomHTTPSParameters
