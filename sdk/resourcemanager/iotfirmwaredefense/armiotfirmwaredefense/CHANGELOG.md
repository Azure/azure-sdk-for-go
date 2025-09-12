# Release History

## 2.0.0 (2025-09-01)
### Breaking Changes

- Function `*SummariesClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, SummaryName, *SummariesClientGetOptions)` to `(context.Context, string, string, string, SummaryType, *SummariesClientGetOptions)`
- Function `*WorkspacesClient.Update` parameter(s) have been changed from `(context.Context, string, string, WorkspaceUpdateDefinition, *WorkspacesClientUpdateOptions)` to `(context.Context, string, string, WorkspaceUpdate, *WorkspacesClientUpdateOptions)`
- Type of `CryptoKey.KeyType` has been changed from `*string` to `*CryptoKeyType`
- `ProvisioningStateAccepted` from enum `ProvisioningState` has been removed
- `SummaryTypeCVE` from enum `SummaryType` has been removed
- Enum `SummaryName` has been removed
- Function `*FirmwaresClient.GenerateDownloadURL` has been removed
- Function `*FirmwaresClient.GenerateFilesystemDownloadURL` has been removed
- Operation `*WorkspacesClient.Delete` has been changed to LRO, use `*WorkspacesClient.BeginDelete` instead.
- Struct `BinaryHardeningListResult` has been removed
- Struct `CryptoCertificateListResult` has been removed
- Struct `CryptoKeyListResult` has been removed
- Struct `CveListResult` has been removed
- Struct `FirmwareList` has been removed
- Struct `PasswordHashListResult` has been removed
- Struct `SbomComponentListResult` has been removed
- Struct `SummaryListResult` has been removed
- Struct `WorkspaceList` has been removed
- Struct `WorkspaceUpdateDefinition` has been removed
- Field `BinaryHardeningListResult` of struct `BinaryHardeningClientListByFirmwareResponse` has been removed
- Field `Nx`, `Pie`, `Relro` of struct `BinaryHardeningFeatures` has been removed
- Field `Architecture`, `Class`, `Features` of struct `BinaryHardeningResult` has been removed
- Field `Canary`, `Nx`, `Pie`, `Relro`, `Stripped` of struct `BinaryHardeningSummaryResource` has been removed
- Field `KeyAlgorithm`, `KeySize`, `Name`, `Role`, `Usage` of struct `CryptoCertificate` has been removed
- Field `Expired`, `ExpiringSoon`, `PairedKeys`, `SelfSigned`, `ShortKeySize`, `TotalCertificates`, `WeakSignature` of struct `CryptoCertificateSummaryResource` has been removed
- Field `CryptoCertificateListResult` of struct `CryptoCertificatesClientListByFirmwareResponse` has been removed
- Field `KeySize` of struct `CryptoKey` has been removed
- Field `PairedKeys`, `PrivateKeys`, `PublicKeys`, `ShortKeySize`, `TotalKeys` of struct `CryptoKeySummaryResource` has been removed
- Field `CryptoKeyListResult` of struct `CryptoKeysClientListByFirmwareResponse` has been removed
- Field `Name` of struct `CveResult` has been removed
- Field `Critical`, `High`, `Low`, `Medium`, `Unknown` of struct `CveSummary` has been removed
- Field `CveListResult` of struct `CvesClientListByFirmwareResponse` has been removed
- Field `FirmwareList` of struct `FirmwaresClientListByWorkspaceResponse` has been removed
- Field `ID` of struct `PairedKey` has been removed
- Field `PasswordHashListResult` of struct `PasswordHashesClientListByFirmwareResponse` has been removed
- Field `SbomComponentListResult` of struct `SbomComponentsClientListByFirmwareResponse` has been removed
- Field `SummaryListResult` of struct `SummariesClientListByFirmwareResponse` has been removed
- Field `WorkspaceList` of struct `WorkspacesClientListByResourceGroupResponse` has been removed
- Field `WorkspaceList` of struct `WorkspacesClientListBySubscriptionResponse` has been removed

### Features Added

- New value `ProvisioningStateAnalyzing`, `ProvisioningStateExtracting`, `ProvisioningStatePending` added to enum type `ProvisioningState`
- New value `SummaryTypeCommonVulnerabilitiesAndExposures` added to enum type `SummaryType`
- New enum type `CertificateUsage` with values `CertificateUsageCRLSign`, `CertificateUsageClientAuthentication`, `CertificateUsageCodeSigning`, `CertificateUsageContentCommitment`, `CertificateUsageDataEncipherment`, `CertificateUsageDecipherOnly`, `CertificateUsageDigitalSignature`, `CertificateUsageEmailProtection`, `CertificateUsageEncipherOnly`, `CertificateUsageKeyAgreement`, `CertificateUsageKeyCertSign`, `CertificateUsageKeyEncipherment`, `CertificateUsageNonRepudiation`, `CertificateUsageOCSPSigning`, `CertificateUsageServerAuthentication`, `CertificateUsageTimeStamping`
- New enum type `CryptoKeyType` with values `CryptoKeyTypePrivate`, `CryptoKeyTypePublic`
- New enum type `ExecutableClass` with values `ExecutableClassX64`, `ExecutableClassX86`
- New enum type `SKUTier` with values `SKUTierBasic`, `SKUTierFree`, `SKUTierPremium`, `SKUTierStandard`
- New function `*ClientFactory.NewUsageMetricsClient() *UsageMetricsClient`
- New function `NewUsageMetricsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*UsageMetricsClient, error)`
- New function `*UsageMetricsClient.Get(context.Context, string, string, string, *UsageMetricsClientGetOptions) (UsageMetricsClientGetResponse, error)`
- New function `*UsageMetricsClient.NewListByWorkspacePager(string, string, *UsageMetricsClientListByWorkspaceOptions) *runtime.Pager[UsageMetricsClientListByWorkspaceResponse]`
- New struct `BinaryHardeningResourceListResult`
- New struct `CryptoCertificateResourceListResult`
- New struct `CryptoKeyResourceListResult`
- New struct `CveResourceListResult`
- New struct `CvssScore`
- New struct `FirmwareListResult`
- New struct `PasswordHashResourceListResult`
- New struct `SKU`
- New struct `SbomComponentResourceListResult`
- New struct `SummaryResourceListResult`
- New struct `UsageMetric`
- New struct `UsageMetricListResult`
- New struct `UsageMetricProperties`
- New struct `WorkspaceListResult`
- New struct `WorkspaceUpdate`
- New anonymous field `BinaryHardeningResourceListResult` in struct `BinaryHardeningClientListByFirmwareResponse`
- New field `NoExecute`, `PositionIndependentExecutable`, `RelocationReadOnly` in struct `BinaryHardeningFeatures`
- New field `ExecutableArchitecture`, `ExecutableClass`, `ProvisioningState`, `SecurityHardeningFeatures` in struct `BinaryHardeningResult`
- New field `NotExecutableStackCount`, `PositionIndependentExecutableCount`, `ProvisioningState`, `RelocationReadOnlyCount`, `StackCanaryCount`, `StrippedBinaryCount` in struct `BinaryHardeningSummaryResource`
- New field `CertificateKeyAlgorithm`, `CertificateKeySize`, `CertificateName`, `CertificateRole`, `CertificateUsage`, `ProvisioningState` in struct `CryptoCertificate`
- New field `ExpiredCertificateCount`, `ExpiringSoonCertificateCount`, `PairedKeyCount`, `ProvisioningState`, `SelfSignedCertificateCount`, `ShortKeySizeCount`, `TotalCertificateCount`, `WeakSignatureCount` in struct `CryptoCertificateSummaryResource`
- New anonymous field `CryptoCertificateResourceListResult` in struct `CryptoCertificatesClientListByFirmwareResponse`
- New field `CryptoKeySize`, `ProvisioningState` in struct `CryptoKey`
- New field `PairedKeyCount`, `PrivateKeyCount`, `ProvisioningState`, `PublicKeyCount`, `ShortKeySizeCount`, `TotalKeyCount` in struct `CryptoKeySummaryResource`
- New anonymous field `CryptoKeyResourceListResult` in struct `CryptoKeysClientListByFirmwareResponse`
- New field `ComponentID`, `ComponentName`, `ComponentVersion`, `CveName`, `CvssScores`, `EffectiveCvssScore`, `EffectiveCvssVersion`, `ProvisioningState` in struct `CveResult`
- New field `CriticalCveCount`, `HighCveCount`, `LowCveCount`, `MediumCveCount`, `ProvisioningState`, `UnknownCveCount` in struct `CveSummary`
- New anonymous field `CveResourceListResult` in struct `CvesClientListByFirmwareResponse`
- New field `ProvisioningState` in struct `FirmwareSummary`
- New anonymous field `FirmwareListResult` in struct `FirmwaresClientListByWorkspaceResponse`
- New field `PairedKeyID` in struct `PairedKey`
- New field `ProvisioningState` in struct `PasswordHash`
- New anonymous field `PasswordHashResourceListResult` in struct `PasswordHashesClientListByFirmwareResponse`
- New field `ProvisioningState` in struct `SbomComponent`
- New anonymous field `SbomComponentResourceListResult` in struct `SbomComponentsClientListByFirmwareResponse`
- New anonymous field `SummaryResourceListResult` in struct `SummariesClientListByFirmwareResponse`
- New field `ProvisioningState` in struct `SummaryResourceProperties`
- New field `SKU` in struct `Workspace`
- New anonymous field `WorkspaceListResult` in struct `WorkspacesClientListByResourceGroupResponse`
- New anonymous field `WorkspaceListResult` in struct `WorkspacesClientListBySubscriptionResponse`


## 2.0.0-beta.1 (2025-06-02)
### Breaking Changes

- Function `*SummariesClient.Get` parameter(s) have been changed from `(context.Context, string, string, string, SummaryName, *SummariesClientGetOptions)` to `(context.Context, string, string, string, SummaryType, *SummariesClientGetOptions)`
- Function `*WorkspacesClient.Update` parameter(s) have been changed from `(context.Context, string, string, WorkspaceUpdateDefinition, *WorkspacesClientUpdateOptions)` to `(context.Context, string, string, WorkspaceUpdate, *WorkspacesClientUpdateOptions)`
- Type of `CryptoKey.KeyType` has been changed from `*string` to `*CryptoKeyType`
- Type of `SummaryResource.Name` has been changed from `*string` to `*SummaryType`
- `ProvisioningStateAccepted` from enum `ProvisioningState` has been removed
- `SummaryTypeCVE` from enum `SummaryType` has been removed
- Enum `SummaryName` has been removed
- Function `*FirmwaresClient.GenerateDownloadURL` has been removed
- Function `*FirmwaresClient.GenerateFilesystemDownloadURL` has been removed
- Struct `BinaryHardeningListResult` has been removed
- Struct `CryptoCertificateListResult` has been removed
- Struct `CryptoKeyListResult` has been removed
- Struct `CveComponent` has been removed
- Struct `CveListResult` has been removed
- Struct `FirmwareList` has been removed
- Struct `PasswordHashListResult` has been removed
- Struct `SbomComponentListResult` has been removed
- Struct `SummaryListResult` has been removed
- Struct `WorkspaceList` has been removed
- Struct `WorkspaceUpdateDefinition` has been removed
- Field `BinaryHardeningListResult` of struct `BinaryHardeningClientListByFirmwareResponse` has been removed
- Field `Nx`, `Pie`, `Relro` of struct `BinaryHardeningFeatures` has been removed
- Field `Architecture`, `Class`, `Features` of struct `BinaryHardeningResult` has been removed
- Field `Canary`, `Nx`, `Pie`, `Relro`, `Stripped` of struct `BinaryHardeningSummaryResource` has been removed
- Field `KeyAlgorithm`, `KeySize`, `Name`, `Role`, `Usage` of struct `CryptoCertificate` has been removed
- Field `Expired`, `ExpiringSoon`, `PairedKeys`, `SelfSigned`, `ShortKeySize`, `TotalCertificates`, `WeakSignature` of struct `CryptoCertificateSummaryResource` has been removed
- Field `CryptoCertificateListResult` of struct `CryptoCertificatesClientListByFirmwareResponse` has been removed
- Field `KeySize` of struct `CryptoKey` has been removed
- Field `PairedKeys`, `PrivateKeys`, `PublicKeys`, `ShortKeySize`, `TotalKeys` of struct `CryptoKeySummaryResource` has been removed
- Field `CryptoKeyListResult` of struct `CryptoKeysClientListByFirmwareResponse` has been removed
- Field `Component`, `CvssScore`, `CvssV2Score`, `CvssV3Score`, `CvssVersion`, `Name` of struct `CveResult` has been removed
- Field `Critical`, `High`, `Low`, `Medium`, `Unknown` of struct `CveSummary` has been removed
- Field `CveListResult` of struct `CvesClientListByFirmwareResponse` has been removed
- Field `FirmwareList` of struct `FirmwaresClientListByWorkspaceResponse` has been removed
- Field `ID` of struct `PairedKey` has been removed
- Field `PasswordHashListResult` of struct `PasswordHashesClientListByFirmwareResponse` has been removed
- Field `SbomComponentListResult` of struct `SbomComponentsClientListByFirmwareResponse` has been removed
- Field `SummaryListResult` of struct `SummariesClientListByFirmwareResponse` has been removed
- Field `WorkspaceList` of struct `WorkspacesClientListByResourceGroupResponse` has been removed
- Field `WorkspaceList` of struct `WorkspacesClientListBySubscriptionResponse` has been removed

### Features Added

- New value `ProvisioningStateAnalyzing`, `ProvisioningStateExtracting`, `ProvisioningStatePending` added to enum type `ProvisioningState`
- New value `SummaryTypeCommonVulnerabilitiesAndExposures` added to enum type `SummaryType`
- New enum type `CertificateUsage` with values `CertificateUsageCRLSign`, `CertificateUsageClientAuthentication`, `CertificateUsageCodeSigning`, `CertificateUsageContentCommitment`, `CertificateUsageDataEncipherment`, `CertificateUsageDecipherOnly`, `CertificateUsageDigitalSignature`, `CertificateUsageEmailProtection`, `CertificateUsageEncipherOnly`, `CertificateUsageKeyAgreement`, `CertificateUsageKeyCertSign`, `CertificateUsageKeyEncipherment`, `CertificateUsageNonRepudiation`, `CertificateUsageOCSPSigning`, `CertificateUsageServerAuthentication`, `CertificateUsageTimeStamping`
- New enum type `CryptoKeyType` with values `CryptoKeyTypePrivate`, `CryptoKeyTypePublic`
- New enum type `ExecutableClass` with values `ExecutableClassX64`, `ExecutableClassX86`
- New enum type `SKUTier` with values `SKUTierBasic`, `SKUTierFree`, `SKUTierPremium`, `SKUTierStandard`
- New function `*ClientFactory.NewUsageMetricsClient() *UsageMetricsClient`
- New function `NewUsageMetricsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*UsageMetricsClient, error)`
- New function `*UsageMetricsClient.Get(context.Context, string, string, string, *UsageMetricsClientGetOptions) (UsageMetricsClientGetResponse, error)`
- New function `*UsageMetricsClient.NewListByWorkspacePager(string, string, *UsageMetricsClientListByWorkspaceOptions) *runtime.Pager[UsageMetricsClientListByWorkspaceResponse]`
- New struct `BinaryHardeningResourceListResult`
- New struct `CryptoCertificateResourceListResult`
- New struct `CryptoKeyResourceListResult`
- New struct `CveResourceListResult`
- New struct `CvssScore`
- New struct `FirmwareListResult`
- New struct `PasswordHashResourceListResult`
- New struct `SKU`
- New struct `SbomComponentResourceListResult`
- New struct `SummaryResourceListResult`
- New struct `UsageMetric`
- New struct `UsageMetricListResult`
- New struct `UsageMetricProperties`
- New struct `WorkspaceListResult`
- New struct `WorkspaceUpdate`
- New anonymous field `BinaryHardeningResourceListResult` in struct `BinaryHardeningClientListByFirmwareResponse`
- New field `NoExecute`, `PositionIndependentExecutable`, `RelocationReadOnly` in struct `BinaryHardeningFeatures`
- New field `ExecutableArchitecture`, `ExecutableClass`, `ProvisioningState`, `SecurityHardeningFeatures` in struct `BinaryHardeningResult`
- New field `NotExecutableStackCount`, `PositionIndependentExecutableCount`, `ProvisioningState`, `RelocationReadOnlyCount`, `StackCanaryCount`, `StrippedBinaryCount` in struct `BinaryHardeningSummaryResource`
- New field `CertificateKeyAlgorithm`, `CertificateKeySize`, `CertificateName`, `CertificateRole`, `CertificateUsage`, `ProvisioningState` in struct `CryptoCertificate`
- New field `ExpiredCertificateCount`, `ExpiringSoonCertificateCount`, `PairedKeyCount`, `ProvisioningState`, `SelfSignedCertificateCount`, `ShortKeySizeCount`, `TotalCertificateCount`, `WeakSignatureCount` in struct `CryptoCertificateSummaryResource`
- New anonymous field `CryptoCertificateResourceListResult` in struct `CryptoCertificatesClientListByFirmwareResponse`
- New field `CryptoKeySize`, `ProvisioningState` in struct `CryptoKey`
- New field `PairedKeyCount`, `PrivateKeyCount`, `ProvisioningState`, `PublicKeyCount`, `ShortKeySizeCount`, `TotalKeyCount` in struct `CryptoKeySummaryResource`
- New anonymous field `CryptoKeyResourceListResult` in struct `CryptoKeysClientListByFirmwareResponse`
- New field `ComponentID`, `ComponentName`, `ComponentVersion`, `CveName`, `CvssScores`, `EffectiveCvssScore`, `EffectiveCvssVersion`, `ProvisioningState` in struct `CveResult`
- New field `CriticalCveCount`, `HighCveCount`, `LowCveCount`, `MediumCveCount`, `ProvisioningState`, `UnknownCveCount` in struct `CveSummary`
- New anonymous field `CveResourceListResult` in struct `CvesClientListByFirmwareResponse`
- New field `ProvisioningState` in struct `FirmwareSummary`
- New anonymous field `FirmwareListResult` in struct `FirmwaresClientListByWorkspaceResponse`
- New field `PairedKeyID` in struct `PairedKey`
- New field `ProvisioningState` in struct `PasswordHash`
- New anonymous field `PasswordHashResourceListResult` in struct `PasswordHashesClientListByFirmwareResponse`
- New field `ProvisioningState` in struct `SbomComponent`
- New anonymous field `SbomComponentResourceListResult` in struct `SbomComponentsClientListByFirmwareResponse`
- New anonymous field `SummaryResourceListResult` in struct `SummariesClientListByFirmwareResponse`
- New field `ProvisioningState` in struct `SummaryResourceProperties`
- New field `SKU` in struct `Workspace`
- New anonymous field `WorkspaceListResult` in struct `WorkspacesClientListByResourceGroupResponse`
- New anonymous field `WorkspaceListResult` in struct `WorkspacesClientListBySubscriptionResponse`


## 1.0.0 (2024-03-22)
### Breaking Changes

- Type of `BinaryHardeningFeatures.Canary` has been changed from `*CanaryFlag` to `*bool`
- Type of `BinaryHardeningFeatures.Nx` has been changed from `*NxFlag` to `*bool`
- Type of `BinaryHardeningFeatures.Pie` has been changed from `*PieFlag` to `*bool`
- Type of `BinaryHardeningFeatures.Relro` has been changed from `*RelroFlag` to `*bool`
- Type of `BinaryHardeningFeatures.Stripped` has been changed from `*StrippedFlag` to `*bool`
- Type of `CryptoCertificate.IsExpired` has been changed from `*IsExpired` to `*bool`
- Type of `CryptoCertificate.IsSelfSigned` has been changed from `*IsSelfSigned` to `*bool`
- Type of `CryptoCertificate.IsShortKeySize` has been changed from `*IsShortKeySize` to `*bool`
- Type of `CryptoCertificate.IsWeakSignature` has been changed from `*IsWeakSignature` to `*bool`
- Type of `CryptoKey.IsShortKeySize` has been changed from `*IsShortKeySize` to `*bool`
- Type of `FirmwareProperties.StatusMessages` has been changed from `[]any` to `[]*StatusMessage`
- Enum `CanaryFlag` has been removed
- Enum `IsExpired` has been removed
- Enum `IsSelfSigned` has been removed
- Enum `IsShortKeySize` has been removed
- Enum `IsUpdateAvailable` has been removed
- Enum `IsWeakSignature` has been removed
- Enum `NxFlag` has been removed
- Enum `PieFlag` has been removed
- Enum `RelroFlag` has been removed
- Enum `StrippedFlag` has been removed
- Function `*ClientFactory.NewFirmwareClient` has been removed
- Function `NewFirmwareClient` has been removed
- Function `*FirmwareClient.Create` has been removed
- Function `*FirmwareClient.Delete` has been removed
- Function `*FirmwareClient.GenerateBinaryHardeningDetails` has been removed
- Function `*FirmwareClient.GenerateBinaryHardeningSummary` has been removed
- Function `*FirmwareClient.GenerateComponentDetails` has been removed
- Function `*FirmwareClient.GenerateCryptoCertificateSummary` has been removed
- Function `*FirmwareClient.GenerateCryptoKeySummary` has been removed
- Function `*FirmwareClient.GenerateCveSummary` has been removed
- Function `*FirmwareClient.GenerateDownloadURL` has been removed
- Function `*FirmwareClient.GenerateFilesystemDownloadURL` has been removed
- Function `*FirmwareClient.GenerateSummary` has been removed
- Function `*FirmwareClient.Get` has been removed
- Function `*FirmwareClient.NewListByWorkspacePager` has been removed
- Function `*FirmwareClient.NewListGenerateBinaryHardeningListPager` has been removed
- Function `*FirmwareClient.NewListGenerateComponentListPager` has been removed
- Function `*FirmwareClient.NewListGenerateCryptoCertificateListPager` has been removed
- Function `*FirmwareClient.NewListGenerateCryptoKeyListPager` has been removed
- Function `*FirmwareClient.NewListGenerateCveListPager` has been removed
- Function `*FirmwareClient.NewListGeneratePasswordHashListPager` has been removed
- Function `*FirmwareClient.Update` has been removed
- Struct `BinaryHardening` has been removed
- Struct `BinaryHardeningList` has been removed
- Struct `BinaryHardeningSummary` has been removed
- Struct `Component` has been removed
- Struct `ComponentList` has been removed
- Struct `CryptoCertificateList` has been removed
- Struct `CryptoCertificateSummary` has been removed
- Struct `CryptoKeyList` has been removed
- Struct `CryptoKeySummary` has been removed
- Struct `Cve` has been removed
- Struct `CveList` has been removed
- Struct `PasswordHashList` has been removed
- Field `Undefined` of struct `CveSummary` has been removed
- Field `AdditionalProperties` of struct `PairedKey` has been removed
- Field `UploadURL` of struct `URLToken` has been removed

### Features Added

- New enum type `SummaryName` with values `SummaryNameBinaryHardening`, `SummaryNameCVE`, `SummaryNameCryptoCertificate`, `SummaryNameCryptoKey`, `SummaryNameFirmware`
- New enum type `SummaryType` with values `SummaryTypeBinaryHardening`, `SummaryTypeCVE`, `SummaryTypeCryptoCertificate`, `SummaryTypeCryptoKey`, `SummaryTypeFirmware`
- New function `NewBinaryHardeningClient(string, azcore.TokenCredential, *arm.ClientOptions) (*BinaryHardeningClient, error)`
- New function `*BinaryHardeningClient.NewListByFirmwarePager(string, string, string, *BinaryHardeningClientListByFirmwareOptions) *runtime.Pager[BinaryHardeningClientListByFirmwareResponse]`
- New function `*BinaryHardeningSummaryResource.GetSummaryResourceProperties() *SummaryResourceProperties`
- New function `*ClientFactory.NewBinaryHardeningClient() *BinaryHardeningClient`
- New function `*ClientFactory.NewCryptoCertificatesClient() *CryptoCertificatesClient`
- New function `*ClientFactory.NewCryptoKeysClient() *CryptoKeysClient`
- New function `*ClientFactory.NewCvesClient() *CvesClient`
- New function `*ClientFactory.NewFirmwaresClient() *FirmwaresClient`
- New function `*ClientFactory.NewPasswordHashesClient() *PasswordHashesClient`
- New function `*ClientFactory.NewSbomComponentsClient() *SbomComponentsClient`
- New function `*ClientFactory.NewSummariesClient() *SummariesClient`
- New function `*CryptoCertificateSummaryResource.GetSummaryResourceProperties() *SummaryResourceProperties`
- New function `NewCryptoCertificatesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CryptoCertificatesClient, error)`
- New function `*CryptoCertificatesClient.NewListByFirmwarePager(string, string, string, *CryptoCertificatesClientListByFirmwareOptions) *runtime.Pager[CryptoCertificatesClientListByFirmwareResponse]`
- New function `*CryptoKeySummaryResource.GetSummaryResourceProperties() *SummaryResourceProperties`
- New function `NewCryptoKeysClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CryptoKeysClient, error)`
- New function `*CryptoKeysClient.NewListByFirmwarePager(string, string, string, *CryptoKeysClientListByFirmwareOptions) *runtime.Pager[CryptoKeysClientListByFirmwareResponse]`
- New function `*CveSummary.GetSummaryResourceProperties() *SummaryResourceProperties`
- New function `NewCvesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*CvesClient, error)`
- New function `*CvesClient.NewListByFirmwarePager(string, string, string, *CvesClientListByFirmwareOptions) *runtime.Pager[CvesClientListByFirmwareResponse]`
- New function `*FirmwareSummary.GetSummaryResourceProperties() *SummaryResourceProperties`
- New function `NewFirmwaresClient(string, azcore.TokenCredential, *arm.ClientOptions) (*FirmwaresClient, error)`
- New function `*FirmwaresClient.Create(context.Context, string, string, string, Firmware, *FirmwaresClientCreateOptions) (FirmwaresClientCreateResponse, error)`
- New function `*FirmwaresClient.Delete(context.Context, string, string, string, *FirmwaresClientDeleteOptions) (FirmwaresClientDeleteResponse, error)`
- New function `*FirmwaresClient.GenerateDownloadURL(context.Context, string, string, string, *FirmwaresClientGenerateDownloadURLOptions) (FirmwaresClientGenerateDownloadURLResponse, error)`
- New function `*FirmwaresClient.GenerateFilesystemDownloadURL(context.Context, string, string, string, *FirmwaresClientGenerateFilesystemDownloadURLOptions) (FirmwaresClientGenerateFilesystemDownloadURLResponse, error)`
- New function `*FirmwaresClient.Get(context.Context, string, string, string, *FirmwaresClientGetOptions) (FirmwaresClientGetResponse, error)`
- New function `*FirmwaresClient.NewListByWorkspacePager(string, string, *FirmwaresClientListByWorkspaceOptions) *runtime.Pager[FirmwaresClientListByWorkspaceResponse]`
- New function `*FirmwaresClient.Update(context.Context, string, string, string, FirmwareUpdateDefinition, *FirmwaresClientUpdateOptions) (FirmwaresClientUpdateResponse, error)`
- New function `NewPasswordHashesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*PasswordHashesClient, error)`
- New function `*PasswordHashesClient.NewListByFirmwarePager(string, string, string, *PasswordHashesClientListByFirmwareOptions) *runtime.Pager[PasswordHashesClientListByFirmwareResponse]`
- New function `NewSbomComponentsClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SbomComponentsClient, error)`
- New function `*SbomComponentsClient.NewListByFirmwarePager(string, string, string, *SbomComponentsClientListByFirmwareOptions) *runtime.Pager[SbomComponentsClientListByFirmwareResponse]`
- New function `NewSummariesClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SummariesClient, error)`
- New function `*SummariesClient.Get(context.Context, string, string, string, SummaryName, *SummariesClientGetOptions) (SummariesClientGetResponse, error)`
- New function `*SummariesClient.NewListByFirmwarePager(string, string, string, *SummariesClientListByFirmwareOptions) *runtime.Pager[SummariesClientListByFirmwareResponse]`
- New function `*SummaryResourceProperties.GetSummaryResourceProperties() *SummaryResourceProperties`
- New struct `BinaryHardeningListResult`
- New struct `BinaryHardeningResource`
- New struct `BinaryHardeningResult`
- New struct `BinaryHardeningSummaryResource`
- New struct `CryptoCertificateListResult`
- New struct `CryptoCertificateResource`
- New struct `CryptoCertificateSummaryResource`
- New struct `CryptoKeyListResult`
- New struct `CryptoKeyResource`
- New struct `CryptoKeySummaryResource`
- New struct `CveComponent`
- New struct `CveListResult`
- New struct `CveResource`
- New struct `CveResult`
- New struct `PasswordHashListResult`
- New struct `PasswordHashResource`
- New struct `SbomComponent`
- New struct `SbomComponentListResult`
- New struct `SbomComponentResource`
- New struct `StatusMessage`
- New struct `SummaryListResult`
- New struct `SummaryResource`
- New field `SummaryType` in struct `CveSummary`
- New field `SummaryType` in struct `FirmwareSummary`


## 0.2.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.1.0 (2023-07-28)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/iotfirmwaredefense/armiotfirmwaredefense` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
