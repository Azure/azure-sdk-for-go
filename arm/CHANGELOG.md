# CHANGELOG

-----

## `v0.3.0-beta`

- Corrected unintentional struct field renaming and client renaming in v0.2.0-beta

-----

## `v0.2.0-beta`

- Added support for DNS, Redis, and Web site services
- Updated Storage service to API version 2015-06-15
- Updated Network to include routing table support
- Address https://github.com/Azure/azure-sdk-for-go/issues/232
- Address https://github.com/Azure/azure-sdk-for-go/issues/231
- Address https://github.com/Azure/azure-sdk-for-go/issues/230
- Address https://github.com/Azure/azure-sdk-for-go/issues/224
- Address https://github.com/Azure/azure-sdk-for-go/issues/184
- Address https://github.com/Azure/azure-sdk-for-go/issues/183

------

## `v0.1.1-beta`

- Improves the UserAgent string to disambiguate arm packages from others in the SDK
- Improves setting the http.Response into generated results (reduces likelihood of a nil reference)
- Adds gofmt, golint, and govet to Travis CI for the arm packages

##### Fixed Issues

- https://github.com/Azure/azure-sdk-for-go/issues/196
- https://github.com/Azure/azure-sdk-for-go/issues/213

------

## v0.1.0-beta

This release addresses the issues raised against the alpha release and adds more features. Most
notably, to address the challenges of encoding JSON
(see the [comments](https://github.com/Azure/go-autorest#handling-empty-values) in the
[go-autorest](https://github.com/Azure/go-autorest) package) by using pointers for *all* structure
fields (with the exception of enumerations). The
[go-autorest/autorest/to](https://github.com/Azure/go-autorest/tree/master/autorest/to) package
provides helpers to convert to / from pointers. The examples demonstrate their usage.

Additionally, the packages now align with Go coding standards and pass both `golint` and `govet`.
Accomplishing this required renaming various fields and parameters (such as changing Url to URL).

##### Changes

- Changed request / response structures to use pointer fields.
- Changed methods to return `error` instead of `autorest.Error`.
- Re-divided methods to ease asynchronous requests.
- Added paged results support.
- Added a UserAgent string.
- Added changes necessary to pass golint and govet.
- Updated README.md with details on asynchronous requests and paging.
- Saved package dependencies through Godep (for the entire SDK).

##### Fixed Issues:

- https://github.com/Azure/azure-sdk-for-go/issues/205
- https://github.com/Azure/azure-sdk-for-go/issues/206
- https://github.com/Azure/azure-sdk-for-go/issues/211
- https://github.com/Azure/azure-sdk-for-go/issues/212

-----

## v0.1.0-alpha

This release introduces the Azure Resource Manager packages generated from the corresponding
[Swagger API](http://swagger.io) [definitions](https://github.com/Azure/azure-rest-api-specs).
