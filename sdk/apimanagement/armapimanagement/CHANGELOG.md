# Release History

## v0.2.0 (released)
### Breaking Changes

- Type of `SchemaContractProperties.Document` has been changed from `*SchemaDocumentProperties` to `map[string]interface{}`
- Struct `SchemaDocumentProperties` has been removed
- Field `Sample` of struct `RepresentationContract` has been removed

### New Content

- New const `APITypeWebsocket`
- New const `SoapAPITypeWebSocket`
- New const `CertificateSourceBuiltIn`
- New const `CertificateSourceKeyVault`
- New const `CertificateStatusCompleted`
- New const `CertificateStatusFailed`
- New const `ProtocolWss`
- New const `CertificateSourceCustom`
- New const `CertificateStatusInProgress`
- New const `CertificateSourceManaged`
- New const `ProtocolWs`
- New function `PossibleCertificateSourceValues() []CertificateSource`
- New function `PossibleCertificateStatusValues() []CertificateStatus`
- New function `CertificateSource.ToPtr() *CertificateSource`
- New function `CertificateStatus.ToPtr() *CertificateStatus`
- New struct `APIContactInformation`
- New struct `APILicenseInformation`
- New struct `ParameterExampleContract`
- New field `License` in struct `APIEntityBaseContract`
- New field `Contact` in struct `APIEntityBaseContract`
- New field `TermsOfServiceURL` in struct `APIEntityBaseContract`
- New field `Zones` in struct `APIManagementServiceUpdateParameters`
- New field `PublicIPAddressID` in struct `APIManagementServiceBaseProperties`
- New field `CertificateSource` in struct `HostnameConfiguration`
- New field `CertificateStatus` in struct `HostnameConfiguration`
- New field `SchemaID` in struct `ParameterContract`
- New field `TypeName` in struct `ParameterContract`
- New field `Examples` in struct `ParameterContract`
- New field `PublicIPAddressID` in struct `AdditionalLocation`

Total 4 breaking change(s), 27 additive change(s).


## v0.1.0 (released)
