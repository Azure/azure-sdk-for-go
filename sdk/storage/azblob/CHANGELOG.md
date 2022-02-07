# Release History

## 0.3.0 (Unreleased)

### Features Added

### Breaking Changes
* Updated to latest `azcore`.  Public surface area is unchanged.

### Bugs Fixed
* Fixed Issue #16816 : ContainerClient.GetSASToken doesn't allow list permission.
* Fixed Issue #16193 : azblob.GetSASToken wrong singed resource. 
* Fixed Issue #16223 : HttpRange does not expose its fields. 

### Other Changes

## 0.2.0 (2021-11-03)

### Breaking Changes
* Clients now have one constructor per authentication method

## 0.1.0 (2021-09-13)

### Features Added
* This is the initial preview release of the `azblob` library
