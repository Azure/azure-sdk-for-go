# Unreleased

## Breaking Changes

### Removed Constants

1. BackupStorageRedundancy1.BackupStorageRedundancy1Geo
1. BackupStorageRedundancy1.BackupStorageRedundancy1Local
1. BackupStorageRedundancy1.BackupStorageRedundancy1Zone
1. CurrentBackupStorageRedundancy.CurrentBackupStorageRedundancyGeo
1. CurrentBackupStorageRedundancy.CurrentBackupStorageRedundancyLocal
1. CurrentBackupStorageRedundancy.CurrentBackupStorageRedundancyZone
1. RequestedBackupStorageRedundancy.RequestedBackupStorageRedundancyGeo
1. RequestedBackupStorageRedundancy.RequestedBackupStorageRedundancyLocal
1. RequestedBackupStorageRedundancy.RequestedBackupStorageRedundancyZone
1. StorageAccountType1.StorageAccountType1GRS
1. StorageAccountType1.StorageAccountType1LRS
1. StorageAccountType1.StorageAccountType1ZRS
1. TargetBackupStorageRedundancy.TargetBackupStorageRedundancyGeo
1. TargetBackupStorageRedundancy.TargetBackupStorageRedundancyLocal
1. TargetBackupStorageRedundancy.TargetBackupStorageRedundancyZone
1. TransparentDataEncryptionActivityStatus.TransparentDataEncryptionActivityStatusDecrypting
1. TransparentDataEncryptionActivityStatus.TransparentDataEncryptionActivityStatusEncrypting
1. TransparentDataEncryptionStatus.TransparentDataEncryptionStatusDisabled
1. TransparentDataEncryptionStatus.TransparentDataEncryptionStatusEnabled

### Removed Funcs

1. *OperationsHealth.UnmarshalJSON([]byte) error
1. *OperationsHealthListResultIterator.Next() error
1. *OperationsHealthListResultIterator.NextWithContext(context.Context) error
1. *OperationsHealthListResultPage.Next() error
1. *OperationsHealthListResultPage.NextWithContext(context.Context) error
1. *ReplicationLinksUnlinkFuture.UnmarshalJSON([]byte) error
1. *TransparentDataEncryption.UnmarshalJSON([]byte) error
1. *TransparentDataEncryptionActivity.UnmarshalJSON([]byte) error
1. NewOperationsHealthClient(string) OperationsHealthClient
1. NewOperationsHealthClientWithBaseURI(string, string) OperationsHealthClient
1. NewOperationsHealthListResultIterator(OperationsHealthListResultPage) OperationsHealthListResultIterator
1. NewOperationsHealthListResultPage(OperationsHealthListResult, func(context.Context, OperationsHealthListResult) (OperationsHealthListResult, error)) OperationsHealthListResultPage
1. NewTransparentDataEncryptionActivitiesClient(string) TransparentDataEncryptionActivitiesClient
1. NewTransparentDataEncryptionActivitiesClientWithBaseURI(string, string) TransparentDataEncryptionActivitiesClient
1. OperationsHealth.MarshalJSON() ([]byte, error)
1. OperationsHealthClient.ListByLocation(context.Context, string) (OperationsHealthListResultPage, error)
1. OperationsHealthClient.ListByLocationComplete(context.Context, string) (OperationsHealthListResultIterator, error)
1. OperationsHealthClient.ListByLocationPreparer(context.Context, string) (*http.Request, error)
1. OperationsHealthClient.ListByLocationResponder(*http.Response) (OperationsHealthListResult, error)
1. OperationsHealthClient.ListByLocationSender(*http.Request) (*http.Response, error)
1. OperationsHealthListResult.IsEmpty() bool
1. OperationsHealthListResult.MarshalJSON() ([]byte, error)
1. OperationsHealthListResultIterator.NotDone() bool
1. OperationsHealthListResultIterator.Response() OperationsHealthListResult
1. OperationsHealthListResultIterator.Value() OperationsHealth
1. OperationsHealthListResultPage.NotDone() bool
1. OperationsHealthListResultPage.Response() OperationsHealthListResult
1. OperationsHealthListResultPage.Values() []OperationsHealth
1. OperationsHealthProperties.MarshalJSON() ([]byte, error)
1. PossibleBackupStorageRedundancy1Values() []BackupStorageRedundancy1
1. PossibleCurrentBackupStorageRedundancyValues() []CurrentBackupStorageRedundancy
1. PossibleRequestedBackupStorageRedundancyValues() []RequestedBackupStorageRedundancy
1. PossibleStorageAccountType1Values() []StorageAccountType1
1. PossibleTargetBackupStorageRedundancyValues() []TargetBackupStorageRedundancy
1. PossibleTransparentDataEncryptionActivityStatusValues() []TransparentDataEncryptionActivityStatus
1. PossibleTransparentDataEncryptionStatusValues() []TransparentDataEncryptionStatus
1. ReplicationLinksClient.Unlink(context.Context, string, string, string, string, UnlinkParameters) (ReplicationLinksUnlinkFuture, error)
1. ReplicationLinksClient.UnlinkPreparer(context.Context, string, string, string, string, UnlinkParameters) (*http.Request, error)
1. ReplicationLinksClient.UnlinkResponder(*http.Response) (autorest.Response, error)
1. ReplicationLinksClient.UnlinkSender(*http.Request) (ReplicationLinksUnlinkFuture, error)
1. ResourceIdentityWithUserAssignedIdentities.MarshalJSON() ([]byte, error)
1. TransparentDataEncryption.MarshalJSON() ([]byte, error)
1. TransparentDataEncryptionActivitiesClient.ListByConfiguration(context.Context, string, string, string) (TransparentDataEncryptionActivityListResult, error)
1. TransparentDataEncryptionActivitiesClient.ListByConfigurationPreparer(context.Context, string, string, string) (*http.Request, error)
1. TransparentDataEncryptionActivitiesClient.ListByConfigurationResponder(*http.Response) (TransparentDataEncryptionActivityListResult, error)
1. TransparentDataEncryptionActivitiesClient.ListByConfigurationSender(*http.Request) (*http.Response, error)
1. TransparentDataEncryptionActivity.MarshalJSON() ([]byte, error)
1. TransparentDataEncryptionActivityProperties.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. OperationsHealth
1. OperationsHealthClient
1. OperationsHealthListResult
1. OperationsHealthListResultIterator
1. OperationsHealthListResultPage
1. OperationsHealthProperties
1. ReplicationLinksUnlinkFuture
1. ResourceIdentityWithUserAssignedIdentities
1. TransparentDataEncryption
1. TransparentDataEncryptionActivitiesClient
1. TransparentDataEncryptionActivity
1. TransparentDataEncryptionActivityListResult
1. TransparentDataEncryptionActivityProperties
1. UnlinkParameters

#### Removed Struct Fields

1. DatabaseUpdate.*DatabaseProperties
1. ManagedInstanceProperties.StorageAccountType
1. RestorableDroppedDatabaseProperties.ElasticPoolID
1. TransparentDataEncryptionProperties.Status

### Signature Changes

#### Funcs

1. ElasticPoolsClient.ListByServer
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, string, *int64
1. ElasticPoolsClient.ListByServerComplete
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, string, *int64
1. ElasticPoolsClient.ListByServerPreparer
	- Params
		- From: context.Context, string, string, *int32
		- To: context.Context, string, string, *int64
1. LedgerDigestUploadsClient.CreateOrUpdate
	- Returns
		- From: LedgerDigestUploads, error
		- To: LedgerDigestUploadsCreateOrUpdateFuture, error
1. LedgerDigestUploadsClient.CreateOrUpdateSender
	- Returns
		- From: *http.Response, error
		- To: LedgerDigestUploadsCreateOrUpdateFuture, error
1. LedgerDigestUploadsClient.Disable
	- Returns
		- From: LedgerDigestUploads, error
		- To: LedgerDigestUploadsDisableFuture, error
1. LedgerDigestUploadsClient.DisableSender
	- Returns
		- From: *http.Response, error
		- To: LedgerDigestUploadsDisableFuture, error
1. ReplicationLinksClient.FailoverAllowDataLossResponder
	- Returns
		- From: autorest.Response, error
		- To: ReplicationLink, error
1. ReplicationLinksClient.FailoverResponder
	- Returns
		- From: autorest.Response, error
		- To: ReplicationLink, error
1. ServerConnectionPoliciesClient.CreateOrUpdate
	- Returns
		- From: ServerConnectionPolicy, error
		- To: ServerConnectionPoliciesCreateOrUpdateFuture, error
1. ServerConnectionPoliciesClient.CreateOrUpdateSender
	- Returns
		- From: *http.Response, error
		- To: ServerConnectionPoliciesCreateOrUpdateFuture, error
1. SyncGroupsClient.ListLogs
	- Params
		- From: context.Context, string, string, string, string, string, string, string, string
		- To: context.Context, string, string, string, string, string, string, SyncGroupsType, string
1. SyncGroupsClient.ListLogsComplete
	- Params
		- From: context.Context, string, string, string, string, string, string, string, string
		- To: context.Context, string, string, string, string, string, string, SyncGroupsType, string
1. SyncGroupsClient.ListLogsPreparer
	- Params
		- From: context.Context, string, string, string, string, string, string, string, string
		- To: context.Context, string, string, string, string, string, string, SyncGroupsType, string
1. TransparentDataEncryptionsClient.CreateOrUpdate
	- Params
		- From: context.Context, string, string, string, TransparentDataEncryption
		- To: context.Context, string, string, string, LogicalDatabaseTransparentDataEncryption
	- Returns
		- From: TransparentDataEncryption, error
		- To: LogicalDatabaseTransparentDataEncryption, error
1. TransparentDataEncryptionsClient.CreateOrUpdatePreparer
	- Params
		- From: context.Context, string, string, string, TransparentDataEncryption
		- To: context.Context, string, string, string, LogicalDatabaseTransparentDataEncryption
1. TransparentDataEncryptionsClient.CreateOrUpdateResponder
	- Returns
		- From: TransparentDataEncryption, error
		- To: LogicalDatabaseTransparentDataEncryption, error
1. TransparentDataEncryptionsClient.Get
	- Returns
		- From: TransparentDataEncryption, error
		- To: LogicalDatabaseTransparentDataEncryption, error
1. TransparentDataEncryptionsClient.GetResponder
	- Returns
		- From: TransparentDataEncryption, error
		- To: LogicalDatabaseTransparentDataEncryption, error

#### Struct Fields

1. CopyLongTermRetentionBackupParametersProperties.TargetBackupStorageRedundancy changed type from TargetBackupStorageRedundancy to BackupStorageRedundancy
1. DatabaseProperties.CurrentBackupStorageRedundancy changed type from CurrentBackupStorageRedundancy to BackupStorageRedundancy
1. DatabaseProperties.RequestedBackupStorageRedundancy changed type from RequestedBackupStorageRedundancy to BackupStorageRedundancy
1. ReplicationLinksFailoverAllowDataLossFuture.Result changed type from func(ReplicationLinksClient) (autorest.Response, error) to func(ReplicationLinksClient) (ReplicationLink, error)
1. ReplicationLinksFailoverFuture.Result changed type from func(ReplicationLinksClient) (autorest.Response, error) to func(ReplicationLinksClient) (ReplicationLink, error)
1. RestorableDroppedDatabaseProperties.BackupStorageRedundancy changed type from BackupStorageRedundancy1 to BackupStorageRedundancy
1. StorageCapability.StorageAccountType changed type from StorageAccountType1 to StorageAccountType
1. UpdateLongTermRetentionBackupParametersProperties.RequestedBackupStorageRedundancy changed type from RequestedBackupStorageRedundancy to BackupStorageRedundancy

## Additive Changes

### New Constants

1. BackupStorageRedundancy.BackupStorageRedundancyGeoZone
1. DatabaseIdentityType.DatabaseIdentityTypeNone
1. DatabaseIdentityType.DatabaseIdentityTypeUserAssigned
1. DatabaseStatus.DatabaseStatusStarting
1. DatabaseStatus.DatabaseStatusStopped
1. DatabaseStatus.DatabaseStatusStopping
1. IdentityType.IdentityTypeSystemAssignedUserAssigned
1. ProvisioningState1.ProvisioningState1Accepted
1. ProvisioningState1.ProvisioningState1Canceled
1. ProvisioningState1.ProvisioningState1Created
1. ProvisioningState1.ProvisioningState1Deleted
1. ProvisioningState1.ProvisioningState1NotSpecified
1. ProvisioningState1.ProvisioningState1Registering
1. ProvisioningState1.ProvisioningState1Running
1. ProvisioningState1.ProvisioningState1TimedOut
1. ProvisioningState1.ProvisioningState1Unrecognized
1. ReplicationMode.ReplicationModeAsync
1. ReplicationMode.ReplicationModeSync
1. ServicePrincipalType.ServicePrincipalTypeNone
1. ServicePrincipalType.ServicePrincipalTypeSystemAssigned
1. SyncGroupsType.SyncGroupsTypeAll
1. SyncGroupsType.SyncGroupsTypeError
1. SyncGroupsType.SyncGroupsTypeSuccess
1. SyncGroupsType.SyncGroupsTypeWarning

### New Funcs

1. *DistributedAvailabilityGroup.UnmarshalJSON([]byte) error
1. *DistributedAvailabilityGroupsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *DistributedAvailabilityGroupsDeleteFuture.UnmarshalJSON([]byte) error
1. *DistributedAvailabilityGroupsListResultIterator.Next() error
1. *DistributedAvailabilityGroupsListResultIterator.NextWithContext(context.Context) error
1. *DistributedAvailabilityGroupsListResultPage.Next() error
1. *DistributedAvailabilityGroupsListResultPage.NextWithContext(context.Context) error
1. *DistributedAvailabilityGroupsUpdateFuture.UnmarshalJSON([]byte) error
1. *EndpointCertificate.UnmarshalJSON([]byte) error
1. *EndpointCertificateListResultIterator.Next() error
1. *EndpointCertificateListResultIterator.NextWithContext(context.Context) error
1. *EndpointCertificateListResultPage.Next() error
1. *EndpointCertificateListResultPage.NextWithContext(context.Context) error
1. *IPv6FirewallRule.UnmarshalJSON([]byte) error
1. *IPv6FirewallRuleListResultIterator.Next() error
1. *IPv6FirewallRuleListResultIterator.NextWithContext(context.Context) error
1. *IPv6FirewallRuleListResultPage.Next() error
1. *IPv6FirewallRuleListResultPage.NextWithContext(context.Context) error
1. *LedgerDigestUploadsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *LedgerDigestUploadsDisableFuture.UnmarshalJSON([]byte) error
1. *LogicalDatabaseTransparentDataEncryption.UnmarshalJSON([]byte) error
1. *LogicalDatabaseTransparentDataEncryptionListResultIterator.Next() error
1. *LogicalDatabaseTransparentDataEncryptionListResultIterator.NextWithContext(context.Context) error
1. *LogicalDatabaseTransparentDataEncryptionListResultPage.Next() error
1. *LogicalDatabaseTransparentDataEncryptionListResultPage.NextWithContext(context.Context) error
1. *ManagedServerDNSAlias.UnmarshalJSON([]byte) error
1. *ManagedServerDNSAliasListResultIterator.Next() error
1. *ManagedServerDNSAliasListResultIterator.NextWithContext(context.Context) error
1. *ManagedServerDNSAliasListResultPage.Next() error
1. *ManagedServerDNSAliasListResultPage.NextWithContext(context.Context) error
1. *ManagedServerDNSAliasesAcquireFuture.UnmarshalJSON([]byte) error
1. *ManagedServerDNSAliasesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *ManagedServerDNSAliasesDeleteFuture.UnmarshalJSON([]byte) error
1. *ServerConnectionPoliciesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *ServerConnectionPolicyListResultIterator.Next() error
1. *ServerConnectionPolicyListResultIterator.NextWithContext(context.Context) error
1. *ServerConnectionPolicyListResultPage.Next() error
1. *ServerConnectionPolicyListResultPage.NextWithContext(context.Context) error
1. *ServerTrustCertificate.UnmarshalJSON([]byte) error
1. *ServerTrustCertificatesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *ServerTrustCertificatesDeleteFuture.UnmarshalJSON([]byte) error
1. *ServerTrustCertificatesListResultIterator.Next() error
1. *ServerTrustCertificatesListResultIterator.NextWithContext(context.Context) error
1. *ServerTrustCertificatesListResultPage.Next() error
1. *ServerTrustCertificatesListResultPage.NextWithContext(context.Context) error
1. DatabaseIdentity.MarshalJSON() ([]byte, error)
1. DatabaseUpdateProperties.MarshalJSON() ([]byte, error)
1. DatabaseUserIdentity.MarshalJSON() ([]byte, error)
1. DistributedAvailabilityGroup.MarshalJSON() ([]byte, error)
1. DistributedAvailabilityGroupProperties.MarshalJSON() ([]byte, error)
1. DistributedAvailabilityGroupsClient.CreateOrUpdate(context.Context, string, string, string, DistributedAvailabilityGroup) (DistributedAvailabilityGroupsCreateOrUpdateFuture, error)
1. DistributedAvailabilityGroupsClient.CreateOrUpdatePreparer(context.Context, string, string, string, DistributedAvailabilityGroup) (*http.Request, error)
1. DistributedAvailabilityGroupsClient.CreateOrUpdateResponder(*http.Response) (DistributedAvailabilityGroup, error)
1. DistributedAvailabilityGroupsClient.CreateOrUpdateSender(*http.Request) (DistributedAvailabilityGroupsCreateOrUpdateFuture, error)
1. DistributedAvailabilityGroupsClient.Delete(context.Context, string, string, string) (DistributedAvailabilityGroupsDeleteFuture, error)
1. DistributedAvailabilityGroupsClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. DistributedAvailabilityGroupsClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. DistributedAvailabilityGroupsClient.DeleteSender(*http.Request) (DistributedAvailabilityGroupsDeleteFuture, error)
1. DistributedAvailabilityGroupsClient.Get(context.Context, string, string, string) (DistributedAvailabilityGroup, error)
1. DistributedAvailabilityGroupsClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. DistributedAvailabilityGroupsClient.GetResponder(*http.Response) (DistributedAvailabilityGroup, error)
1. DistributedAvailabilityGroupsClient.GetSender(*http.Request) (*http.Response, error)
1. DistributedAvailabilityGroupsClient.ListByInstance(context.Context, string, string) (DistributedAvailabilityGroupsListResultPage, error)
1. DistributedAvailabilityGroupsClient.ListByInstanceComplete(context.Context, string, string) (DistributedAvailabilityGroupsListResultIterator, error)
1. DistributedAvailabilityGroupsClient.ListByInstancePreparer(context.Context, string, string) (*http.Request, error)
1. DistributedAvailabilityGroupsClient.ListByInstanceResponder(*http.Response) (DistributedAvailabilityGroupsListResult, error)
1. DistributedAvailabilityGroupsClient.ListByInstanceSender(*http.Request) (*http.Response, error)
1. DistributedAvailabilityGroupsClient.Update(context.Context, string, string, string, DistributedAvailabilityGroup) (DistributedAvailabilityGroupsUpdateFuture, error)
1. DistributedAvailabilityGroupsClient.UpdatePreparer(context.Context, string, string, string, DistributedAvailabilityGroup) (*http.Request, error)
1. DistributedAvailabilityGroupsClient.UpdateResponder(*http.Response) (DistributedAvailabilityGroup, error)
1. DistributedAvailabilityGroupsClient.UpdateSender(*http.Request) (DistributedAvailabilityGroupsUpdateFuture, error)
1. DistributedAvailabilityGroupsListResult.IsEmpty() bool
1. DistributedAvailabilityGroupsListResult.MarshalJSON() ([]byte, error)
1. DistributedAvailabilityGroupsListResultIterator.NotDone() bool
1. DistributedAvailabilityGroupsListResultIterator.Response() DistributedAvailabilityGroupsListResult
1. DistributedAvailabilityGroupsListResultIterator.Value() DistributedAvailabilityGroup
1. DistributedAvailabilityGroupsListResultPage.NotDone() bool
1. DistributedAvailabilityGroupsListResultPage.Response() DistributedAvailabilityGroupsListResult
1. DistributedAvailabilityGroupsListResultPage.Values() []DistributedAvailabilityGroup
1. EndpointCertificate.MarshalJSON() ([]byte, error)
1. EndpointCertificateListResult.IsEmpty() bool
1. EndpointCertificateListResult.MarshalJSON() ([]byte, error)
1. EndpointCertificateListResultIterator.NotDone() bool
1. EndpointCertificateListResultIterator.Response() EndpointCertificateListResult
1. EndpointCertificateListResultIterator.Value() EndpointCertificate
1. EndpointCertificateListResultPage.NotDone() bool
1. EndpointCertificateListResultPage.Response() EndpointCertificateListResult
1. EndpointCertificateListResultPage.Values() []EndpointCertificate
1. EndpointCertificatesClient.Get(context.Context, string, string, string) (EndpointCertificate, error)
1. EndpointCertificatesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. EndpointCertificatesClient.GetResponder(*http.Response) (EndpointCertificate, error)
1. EndpointCertificatesClient.GetSender(*http.Request) (*http.Response, error)
1. EndpointCertificatesClient.ListByInstance(context.Context, string, string) (EndpointCertificateListResultPage, error)
1. EndpointCertificatesClient.ListByInstanceComplete(context.Context, string, string) (EndpointCertificateListResultIterator, error)
1. EndpointCertificatesClient.ListByInstancePreparer(context.Context, string, string) (*http.Request, error)
1. EndpointCertificatesClient.ListByInstanceResponder(*http.Response) (EndpointCertificateListResult, error)
1. EndpointCertificatesClient.ListByInstanceSender(*http.Request) (*http.Response, error)
1. IPv6FirewallRule.MarshalJSON() ([]byte, error)
1. IPv6FirewallRuleListResult.IsEmpty() bool
1. IPv6FirewallRuleListResult.MarshalJSON() ([]byte, error)
1. IPv6FirewallRuleListResultIterator.NotDone() bool
1. IPv6FirewallRuleListResultIterator.Response() IPv6FirewallRuleListResult
1. IPv6FirewallRuleListResultIterator.Value() IPv6FirewallRule
1. IPv6FirewallRuleListResultPage.NotDone() bool
1. IPv6FirewallRuleListResultPage.Response() IPv6FirewallRuleListResult
1. IPv6FirewallRuleListResultPage.Values() []IPv6FirewallRule
1. IPv6FirewallRulesClient.CreateOrUpdate(context.Context, string, string, string, IPv6FirewallRule) (IPv6FirewallRule, error)
1. IPv6FirewallRulesClient.CreateOrUpdatePreparer(context.Context, string, string, string, IPv6FirewallRule) (*http.Request, error)
1. IPv6FirewallRulesClient.CreateOrUpdateResponder(*http.Response) (IPv6FirewallRule, error)
1. IPv6FirewallRulesClient.CreateOrUpdateSender(*http.Request) (*http.Response, error)
1. IPv6FirewallRulesClient.Delete(context.Context, string, string, string) (autorest.Response, error)
1. IPv6FirewallRulesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. IPv6FirewallRulesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. IPv6FirewallRulesClient.DeleteSender(*http.Request) (*http.Response, error)
1. IPv6FirewallRulesClient.Get(context.Context, string, string, string) (IPv6FirewallRule, error)
1. IPv6FirewallRulesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. IPv6FirewallRulesClient.GetResponder(*http.Response) (IPv6FirewallRule, error)
1. IPv6FirewallRulesClient.GetSender(*http.Request) (*http.Response, error)
1. IPv6FirewallRulesClient.ListByServer(context.Context, string, string) (IPv6FirewallRuleListResultPage, error)
1. IPv6FirewallRulesClient.ListByServerComplete(context.Context, string, string) (IPv6FirewallRuleListResultIterator, error)
1. IPv6FirewallRulesClient.ListByServerPreparer(context.Context, string, string) (*http.Request, error)
1. IPv6FirewallRulesClient.ListByServerResponder(*http.Response) (IPv6FirewallRuleListResult, error)
1. IPv6FirewallRulesClient.ListByServerSender(*http.Request) (*http.Response, error)
1. LogicalDatabaseTransparentDataEncryption.MarshalJSON() ([]byte, error)
1. LogicalDatabaseTransparentDataEncryptionListResult.IsEmpty() bool
1. LogicalDatabaseTransparentDataEncryptionListResult.MarshalJSON() ([]byte, error)
1. LogicalDatabaseTransparentDataEncryptionListResultIterator.NotDone() bool
1. LogicalDatabaseTransparentDataEncryptionListResultIterator.Response() LogicalDatabaseTransparentDataEncryptionListResult
1. LogicalDatabaseTransparentDataEncryptionListResultIterator.Value() LogicalDatabaseTransparentDataEncryption
1. LogicalDatabaseTransparentDataEncryptionListResultPage.NotDone() bool
1. LogicalDatabaseTransparentDataEncryptionListResultPage.Response() LogicalDatabaseTransparentDataEncryptionListResult
1. LogicalDatabaseTransparentDataEncryptionListResultPage.Values() []LogicalDatabaseTransparentDataEncryption
1. ManagedDatabaseSensitivityLabelsClient.ListByDatabase(context.Context, string, string, string, string) (SensitivityLabelListResultPage, error)
1. ManagedDatabaseSensitivityLabelsClient.ListByDatabaseComplete(context.Context, string, string, string, string) (SensitivityLabelListResultIterator, error)
1. ManagedDatabaseSensitivityLabelsClient.ListByDatabasePreparer(context.Context, string, string, string, string) (*http.Request, error)
1. ManagedDatabaseSensitivityLabelsClient.ListByDatabaseResponder(*http.Response) (SensitivityLabelListResult, error)
1. ManagedDatabaseSensitivityLabelsClient.ListByDatabaseSender(*http.Request) (*http.Response, error)
1. ManagedServerDNSAlias.MarshalJSON() ([]byte, error)
1. ManagedServerDNSAliasListResult.IsEmpty() bool
1. ManagedServerDNSAliasListResult.MarshalJSON() ([]byte, error)
1. ManagedServerDNSAliasListResultIterator.NotDone() bool
1. ManagedServerDNSAliasListResultIterator.Response() ManagedServerDNSAliasListResult
1. ManagedServerDNSAliasListResultIterator.Value() ManagedServerDNSAlias
1. ManagedServerDNSAliasListResultPage.NotDone() bool
1. ManagedServerDNSAliasListResultPage.Response() ManagedServerDNSAliasListResult
1. ManagedServerDNSAliasListResultPage.Values() []ManagedServerDNSAlias
1. ManagedServerDNSAliasProperties.MarshalJSON() ([]byte, error)
1. ManagedServerDNSAliasesClient.Acquire(context.Context, string, string, string, ManagedServerDNSAliasAcquisition) (ManagedServerDNSAliasesAcquireFuture, error)
1. ManagedServerDNSAliasesClient.AcquirePreparer(context.Context, string, string, string, ManagedServerDNSAliasAcquisition) (*http.Request, error)
1. ManagedServerDNSAliasesClient.AcquireResponder(*http.Response) (ManagedServerDNSAlias, error)
1. ManagedServerDNSAliasesClient.AcquireSender(*http.Request) (ManagedServerDNSAliasesAcquireFuture, error)
1. ManagedServerDNSAliasesClient.CreateOrUpdate(context.Context, string, string, string, ManagedServerDNSAliasCreation) (ManagedServerDNSAliasesCreateOrUpdateFuture, error)
1. ManagedServerDNSAliasesClient.CreateOrUpdatePreparer(context.Context, string, string, string, ManagedServerDNSAliasCreation) (*http.Request, error)
1. ManagedServerDNSAliasesClient.CreateOrUpdateResponder(*http.Response) (ManagedServerDNSAlias, error)
1. ManagedServerDNSAliasesClient.CreateOrUpdateSender(*http.Request) (ManagedServerDNSAliasesCreateOrUpdateFuture, error)
1. ManagedServerDNSAliasesClient.Delete(context.Context, string, string, string) (ManagedServerDNSAliasesDeleteFuture, error)
1. ManagedServerDNSAliasesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. ManagedServerDNSAliasesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ManagedServerDNSAliasesClient.DeleteSender(*http.Request) (ManagedServerDNSAliasesDeleteFuture, error)
1. ManagedServerDNSAliasesClient.Get(context.Context, string, string, string) (ManagedServerDNSAlias, error)
1. ManagedServerDNSAliasesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. ManagedServerDNSAliasesClient.GetResponder(*http.Response) (ManagedServerDNSAlias, error)
1. ManagedServerDNSAliasesClient.GetSender(*http.Request) (*http.Response, error)
1. ManagedServerDNSAliasesClient.ListByManagedInstance(context.Context, string, string) (ManagedServerDNSAliasListResultPage, error)
1. ManagedServerDNSAliasesClient.ListByManagedInstanceComplete(context.Context, string, string) (ManagedServerDNSAliasListResultIterator, error)
1. ManagedServerDNSAliasesClient.ListByManagedInstancePreparer(context.Context, string, string) (*http.Request, error)
1. ManagedServerDNSAliasesClient.ListByManagedInstanceResponder(*http.Response) (ManagedServerDNSAliasListResult, error)
1. ManagedServerDNSAliasesClient.ListByManagedInstanceSender(*http.Request) (*http.Response, error)
1. NewDistributedAvailabilityGroupsClient(string) DistributedAvailabilityGroupsClient
1. NewDistributedAvailabilityGroupsClientWithBaseURI(string, string) DistributedAvailabilityGroupsClient
1. NewDistributedAvailabilityGroupsListResultIterator(DistributedAvailabilityGroupsListResultPage) DistributedAvailabilityGroupsListResultIterator
1. NewDistributedAvailabilityGroupsListResultPage(DistributedAvailabilityGroupsListResult, func(context.Context, DistributedAvailabilityGroupsListResult) (DistributedAvailabilityGroupsListResult, error)) DistributedAvailabilityGroupsListResultPage
1. NewEndpointCertificateListResultIterator(EndpointCertificateListResultPage) EndpointCertificateListResultIterator
1. NewEndpointCertificateListResultPage(EndpointCertificateListResult, func(context.Context, EndpointCertificateListResult) (EndpointCertificateListResult, error)) EndpointCertificateListResultPage
1. NewEndpointCertificatesClient(string) EndpointCertificatesClient
1. NewEndpointCertificatesClientWithBaseURI(string, string) EndpointCertificatesClient
1. NewIPv6FirewallRuleListResultIterator(IPv6FirewallRuleListResultPage) IPv6FirewallRuleListResultIterator
1. NewIPv6FirewallRuleListResultPage(IPv6FirewallRuleListResult, func(context.Context, IPv6FirewallRuleListResult) (IPv6FirewallRuleListResult, error)) IPv6FirewallRuleListResultPage
1. NewIPv6FirewallRulesClient(string) IPv6FirewallRulesClient
1. NewIPv6FirewallRulesClientWithBaseURI(string, string) IPv6FirewallRulesClient
1. NewLogicalDatabaseTransparentDataEncryptionListResultIterator(LogicalDatabaseTransparentDataEncryptionListResultPage) LogicalDatabaseTransparentDataEncryptionListResultIterator
1. NewLogicalDatabaseTransparentDataEncryptionListResultPage(LogicalDatabaseTransparentDataEncryptionListResult, func(context.Context, LogicalDatabaseTransparentDataEncryptionListResult) (LogicalDatabaseTransparentDataEncryptionListResult, error)) LogicalDatabaseTransparentDataEncryptionListResultPage
1. NewManagedServerDNSAliasListResultIterator(ManagedServerDNSAliasListResultPage) ManagedServerDNSAliasListResultIterator
1. NewManagedServerDNSAliasListResultPage(ManagedServerDNSAliasListResult, func(context.Context, ManagedServerDNSAliasListResult) (ManagedServerDNSAliasListResult, error)) ManagedServerDNSAliasListResultPage
1. NewManagedServerDNSAliasesClient(string) ManagedServerDNSAliasesClient
1. NewManagedServerDNSAliasesClientWithBaseURI(string, string) ManagedServerDNSAliasesClient
1. NewServerConnectionPolicyListResultIterator(ServerConnectionPolicyListResultPage) ServerConnectionPolicyListResultIterator
1. NewServerConnectionPolicyListResultPage(ServerConnectionPolicyListResult, func(context.Context, ServerConnectionPolicyListResult) (ServerConnectionPolicyListResult, error)) ServerConnectionPolicyListResultPage
1. NewServerTrustCertificatesClient(string) ServerTrustCertificatesClient
1. NewServerTrustCertificatesClientWithBaseURI(string, string) ServerTrustCertificatesClient
1. NewServerTrustCertificatesListResultIterator(ServerTrustCertificatesListResultPage) ServerTrustCertificatesListResultIterator
1. NewServerTrustCertificatesListResultPage(ServerTrustCertificatesListResult, func(context.Context, ServerTrustCertificatesListResult) (ServerTrustCertificatesListResult, error)) ServerTrustCertificatesListResultPage
1. PossibleDatabaseIdentityTypeValues() []DatabaseIdentityType
1. PossibleReplicationModeValues() []ReplicationMode
1. PossibleServicePrincipalTypeValues() []ServicePrincipalType
1. PossibleSyncGroupsTypeValues() []SyncGroupsType
1. SensitivityLabelsClient.ListByDatabase(context.Context, string, string, string, string) (SensitivityLabelListResultPage, error)
1. SensitivityLabelsClient.ListByDatabaseComplete(context.Context, string, string, string, string) (SensitivityLabelListResultIterator, error)
1. SensitivityLabelsClient.ListByDatabasePreparer(context.Context, string, string, string, string) (*http.Request, error)
1. SensitivityLabelsClient.ListByDatabaseResponder(*http.Response) (SensitivityLabelListResult, error)
1. SensitivityLabelsClient.ListByDatabaseSender(*http.Request) (*http.Response, error)
1. ServerConnectionPoliciesClient.ListByServer(context.Context, string, string) (ServerConnectionPolicyListResultPage, error)
1. ServerConnectionPoliciesClient.ListByServerComplete(context.Context, string, string) (ServerConnectionPolicyListResultIterator, error)
1. ServerConnectionPoliciesClient.ListByServerPreparer(context.Context, string, string) (*http.Request, error)
1. ServerConnectionPoliciesClient.ListByServerResponder(*http.Response) (ServerConnectionPolicyListResult, error)
1. ServerConnectionPoliciesClient.ListByServerSender(*http.Request) (*http.Response, error)
1. ServerConnectionPolicyListResult.IsEmpty() bool
1. ServerConnectionPolicyListResult.MarshalJSON() ([]byte, error)
1. ServerConnectionPolicyListResultIterator.NotDone() bool
1. ServerConnectionPolicyListResultIterator.Response() ServerConnectionPolicyListResult
1. ServerConnectionPolicyListResultIterator.Value() ServerConnectionPolicy
1. ServerConnectionPolicyListResultPage.NotDone() bool
1. ServerConnectionPolicyListResultPage.Response() ServerConnectionPolicyListResult
1. ServerConnectionPolicyListResultPage.Values() []ServerConnectionPolicy
1. ServerTrustCertificate.MarshalJSON() ([]byte, error)
1. ServerTrustCertificateProperties.MarshalJSON() ([]byte, error)
1. ServerTrustCertificatesClient.CreateOrUpdate(context.Context, string, string, string, ServerTrustCertificate) (ServerTrustCertificatesCreateOrUpdateFuture, error)
1. ServerTrustCertificatesClient.CreateOrUpdatePreparer(context.Context, string, string, string, ServerTrustCertificate) (*http.Request, error)
1. ServerTrustCertificatesClient.CreateOrUpdateResponder(*http.Response) (ServerTrustCertificate, error)
1. ServerTrustCertificatesClient.CreateOrUpdateSender(*http.Request) (ServerTrustCertificatesCreateOrUpdateFuture, error)
1. ServerTrustCertificatesClient.Delete(context.Context, string, string, string) (ServerTrustCertificatesDeleteFuture, error)
1. ServerTrustCertificatesClient.DeletePreparer(context.Context, string, string, string) (*http.Request, error)
1. ServerTrustCertificatesClient.DeleteResponder(*http.Response) (autorest.Response, error)
1. ServerTrustCertificatesClient.DeleteSender(*http.Request) (ServerTrustCertificatesDeleteFuture, error)
1. ServerTrustCertificatesClient.Get(context.Context, string, string, string) (ServerTrustCertificate, error)
1. ServerTrustCertificatesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. ServerTrustCertificatesClient.GetResponder(*http.Response) (ServerTrustCertificate, error)
1. ServerTrustCertificatesClient.GetSender(*http.Request) (*http.Response, error)
1. ServerTrustCertificatesClient.ListByInstance(context.Context, string, string) (ServerTrustCertificatesListResultPage, error)
1. ServerTrustCertificatesClient.ListByInstanceComplete(context.Context, string, string) (ServerTrustCertificatesListResultIterator, error)
1. ServerTrustCertificatesClient.ListByInstancePreparer(context.Context, string, string) (*http.Request, error)
1. ServerTrustCertificatesClient.ListByInstanceResponder(*http.Response) (ServerTrustCertificatesListResult, error)
1. ServerTrustCertificatesClient.ListByInstanceSender(*http.Request) (*http.Response, error)
1. ServerTrustCertificatesListResult.IsEmpty() bool
1. ServerTrustCertificatesListResult.MarshalJSON() ([]byte, error)
1. ServerTrustCertificatesListResultIterator.NotDone() bool
1. ServerTrustCertificatesListResultIterator.Response() ServerTrustCertificatesListResult
1. ServerTrustCertificatesListResultIterator.Value() ServerTrustCertificate
1. ServerTrustCertificatesListResultPage.NotDone() bool
1. ServerTrustCertificatesListResultPage.Response() ServerTrustCertificatesListResult
1. ServerTrustCertificatesListResultPage.Values() []ServerTrustCertificate
1. ServicePrincipal.MarshalJSON() ([]byte, error)
1. TransparentDataEncryptionsClient.ListByDatabase(context.Context, string, string, string) (LogicalDatabaseTransparentDataEncryptionListResultPage, error)
1. TransparentDataEncryptionsClient.ListByDatabaseComplete(context.Context, string, string, string) (LogicalDatabaseTransparentDataEncryptionListResultIterator, error)
1. TransparentDataEncryptionsClient.ListByDatabasePreparer(context.Context, string, string, string) (*http.Request, error)
1. TransparentDataEncryptionsClient.ListByDatabaseResponder(*http.Response) (LogicalDatabaseTransparentDataEncryptionListResult, error)
1. TransparentDataEncryptionsClient.ListByDatabaseSender(*http.Request) (*http.Response, error)

### Struct Changes

#### New Structs

1. DatabaseIdentity
1. DatabaseUpdateProperties
1. DatabaseUserIdentity
1. DistributedAvailabilityGroup
1. DistributedAvailabilityGroupProperties
1. DistributedAvailabilityGroupsClient
1. DistributedAvailabilityGroupsCreateOrUpdateFuture
1. DistributedAvailabilityGroupsDeleteFuture
1. DistributedAvailabilityGroupsListResult
1. DistributedAvailabilityGroupsListResultIterator
1. DistributedAvailabilityGroupsListResultPage
1. DistributedAvailabilityGroupsUpdateFuture
1. EndpointCertificate
1. EndpointCertificateListResult
1. EndpointCertificateListResultIterator
1. EndpointCertificateListResultPage
1. EndpointCertificateProperties
1. EndpointCertificatesClient
1. IPv6FirewallRule
1. IPv6FirewallRuleListResult
1. IPv6FirewallRuleListResultIterator
1. IPv6FirewallRuleListResultPage
1. IPv6FirewallRulesClient
1. IPv6ServerFirewallRuleProperties
1. LedgerDigestUploadsCreateOrUpdateFuture
1. LedgerDigestUploadsDisableFuture
1. LogicalDatabaseTransparentDataEncryption
1. LogicalDatabaseTransparentDataEncryptionListResult
1. LogicalDatabaseTransparentDataEncryptionListResultIterator
1. LogicalDatabaseTransparentDataEncryptionListResultPage
1. ManagedServerDNSAlias
1. ManagedServerDNSAliasAcquisition
1. ManagedServerDNSAliasCreation
1. ManagedServerDNSAliasListResult
1. ManagedServerDNSAliasListResultIterator
1. ManagedServerDNSAliasListResultPage
1. ManagedServerDNSAliasProperties
1. ManagedServerDNSAliasesAcquireFuture
1. ManagedServerDNSAliasesClient
1. ManagedServerDNSAliasesCreateOrUpdateFuture
1. ManagedServerDNSAliasesDeleteFuture
1. ServerConnectionPoliciesCreateOrUpdateFuture
1. ServerConnectionPolicyListResult
1. ServerConnectionPolicyListResultIterator
1. ServerConnectionPolicyListResultPage
1. ServerTrustCertificate
1. ServerTrustCertificateProperties
1. ServerTrustCertificatesClient
1. ServerTrustCertificatesCreateOrUpdateFuture
1. ServerTrustCertificatesDeleteFuture
1. ServerTrustCertificatesListResult
1. ServerTrustCertificatesListResultIterator
1. ServerTrustCertificatesListResultPage
1. ServicePrincipal

#### New Struct Fields

1. Database.Identity
1. DatabaseProperties.FederatedClientID
1. DatabaseProperties.SourceResourceID
1. DatabaseUpdate.*DatabaseUpdateProperties
1. DatabaseUpdate.Identity
1. ElasticPoolProperties.HighAvailabilityReplicaCount
1. ElasticPoolUpdateProperties.HighAvailabilityReplicaCount
1. ManagedInstanceProperties.CurrentBackupStorageRedundancy
1. ManagedInstanceProperties.RequestedBackupStorageRedundancy
1. ManagedInstanceProperties.ServicePrincipal
1. TransparentDataEncryptionProperties.State
