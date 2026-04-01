# Release History

## 0.8.0 (2026-04-01)
### Breaking Changes

- Type of `Operation.Origin` has been changed from `*OperationOrigin` to `*Origin`
- Enum `OperationOrigin` has been removed
- Operation `*SQLServerInstancesClient.Update` has been changed to LRO, use `*SQLServerInstancesClient.BeginUpdate` instead.
- Struct `CommonSKU` has been removed
- Struct `ErrorResponse` has been removed
- Struct `ErrorResponseBody` has been removed
- Struct `ProxyResource` has been removed
- Struct `Resource` has been removed
- Struct `TrackedResource` has been removed
- Field `Properties` of struct `Operation` has been removed

### Features Added

- New value `ArcSQLServerLicenseTypeFabricCapacity`, `ArcSQLServerLicenseTypeLicenseOnly`, `ArcSQLServerLicenseTypePayg`, `ArcSQLServerLicenseTypeServerCAL` added to enum type `ArcSQLServerLicenseType`
- New value `ConnectionStatusDiscovered` added to enum type `ConnectionStatus`
- New value `EditionTypeBusinessIntelligence`, `EditionTypeStandardDeveloper`, `EditionTypeUnknown` added to enum type `EditionType`
- New value `HostTypeAWSKubernetesService`, `HostTypeAWSVMWareVirtualMachine`, `HostTypeAzureKubernetesService`, `HostTypeAzureVMWareVirtualMachine`, `HostTypeAzureVirtualMachine`, `HostTypeContainer`, `HostTypeGCPKubernetesService`, `HostTypeGCPVMWareVirtualMachine`, `HostTypeHyperVVirtualMachine` added to enum type `HostType`
- New value `SQLVersionSQLServer2025` added to enum type `SQLVersion`
- New enum type `ActionType` with values `ActionTypeInternal`
- New enum type `ActivationState` with values `ActivationStateActivated`, `ActivationStateDeactivated`
- New enum type `AggregationType` with values `AggregationTypeAverage`, `AggregationTypeCount`, `AggregationTypeMaximum`, `AggregationTypeMinimum`, `AggregationTypeSum`
- New enum type `AlwaysOnRole` with values `AlwaysOnRoleAvailabilityGroupReplica`, `AlwaysOnRoleFailoverClusterInstance`, `AlwaysOnRoleFailoverClusterNode`, `AlwaysOnRoleNone`
- New enum type `ArcSQLServerAvailabilityGroupTypeFilter` with values `ArcSQLServerAvailabilityGroupTypeFilterContained`, `ArcSQLServerAvailabilityGroupTypeFilterDefault`, `ArcSQLServerAvailabilityGroupTypeFilterDistributed`
- New enum type `ArcSQLServerAvailabilityMode` with values `ArcSQLServerAvailabilityModeAsynchronousCommit`, `ArcSQLServerAvailabilityModeSynchronousCommit`
- New enum type `ArcSQLServerFailoverMode` with values `ArcSQLServerFailoverModeAutomatic`, `ArcSQLServerFailoverModeExternal`, `ArcSQLServerFailoverModeManual`, `ArcSQLServerFailoverModeNone`
- New enum type `AssessmentStatus` with values `AssessmentStatusFailure`, `AssessmentStatusSuccess`, `AssessmentStatusWarning`
- New enum type `AutomatedBackupPreference` with values `AutomatedBackupPreferenceNone`, `AutomatedBackupPreferencePrimary`, `AutomatedBackupPreferenceSecondary`, `AutomatedBackupPreferenceSecondaryOnly`
- New enum type `AzureManagedInstanceRole` with values `AzureManagedInstanceRolePrimary`, `AzureManagedInstanceRoleSecondary`
- New enum type `BillingPlan` with values `BillingPlanPaid`, `BillingPlanPayg`
- New enum type `ClusterType` with values `ClusterTypeNone`, `ClusterTypeWsfc`
- New enum type `ConnectionAuth` with values `ConnectionAuthCertificate`, `ConnectionAuthCertificateWindowsKerberos`, `ConnectionAuthCertificateWindowsNegotiate`, `ConnectionAuthCertificateWindowsNtlm`, `ConnectionAuthWindowsKerberos`, `ConnectionAuthWindowsKerberosCertificate`, `ConnectionAuthWindowsNegotiate`, `ConnectionAuthWindowsNegotiateCertificate`, `ConnectionAuthWindowsNtlm`, `ConnectionAuthWindowsNtlmCertificate`
- New enum type `DatabaseCreateMode` with values `DatabaseCreateModeDefault`, `DatabaseCreateModePointInTimeRestore`
- New enum type `DatabaseState` with values `DatabaseStateCopying`, `DatabaseStateEmergency`, `DatabaseStateOffline`, `DatabaseStateOfflineSecondary`, `DatabaseStateOnline`, `DatabaseStateRecovering`, `DatabaseStateRecoveryPending`, `DatabaseStateRestoring`, `DatabaseStateSuspect`
- New enum type `DbFailover` with values `DbFailoverOFF`, `DbFailoverON`
- New enum type `DifferentialBackupHours` with values `DifferentialBackupHoursTwelve`, `DifferentialBackupHoursTwentyFour`
- New enum type `DiscoverySource` with values `DiscoverySourceADS`, `DiscoverySourceAzureArc`, `DiscoverySourceAzureMigrate`, `DiscoverySourceDMSCLI`, `DiscoverySourceDMSPS`, `DiscoverySourceDMSPortal`, `DiscoverySourceDMSSDK`, `DiscoverySourceImport`, `DiscoverySourceOther`, `DiscoverySourceSsma`, `DiscoverySourceSsms`
- New enum type `DtcSupport` with values `DtcSupportNone`, `DtcSupportPERDB`
- New enum type `EncryptionAlgorithm` with values `EncryptionAlgorithmAES`, `EncryptionAlgorithmAESRC4`, `EncryptionAlgorithmNone`, `EncryptionAlgorithmNoneAES`, `EncryptionAlgorithmNoneAESRC4`, `EncryptionAlgorithmNoneRC4`, `EncryptionAlgorithmNoneRC4AES`, `EncryptionAlgorithmRC4`, `EncryptionAlgorithmRC4AES`
- New enum type `ExecutionState` with values `ExecutionStateRunning`, `ExecutionStateWaiting`
- New enum type `FailoverGroupPartnerSyncMode` with values `FailoverGroupPartnerSyncModeAsync`, `FailoverGroupPartnerSyncModeSync`
- New enum type `FailureConditionLevel` with values `FailureConditionLevelFive`, `FailureConditionLevelFour`, `FailureConditionLevelOne`, `FailureConditionLevelThree`, `FailureConditionLevelTwo`
- New enum type `IdentityType` with values `IdentityTypeSystemAssignedManagedIdentity`, `IdentityTypeUserAssignedManagedIdentity`
- New enum type `InitiatedFrom` with values `InitiatedFromADS`, `InitiatedFromAzureArc`, `InitiatedFromDMSCLI`, `InitiatedFromDMSPS`, `InitiatedFromDMSPortal`, `InitiatedFromDMSSDK`, `InitiatedFromOther`, `InitiatedFromSsma`, `InitiatedFromSsms`
- New enum type `InstanceFailoverGroupRole` with values `InstanceFailoverGroupRoleForcePrimaryAllowDataLoss`, `InstanceFailoverGroupRoleForceSecondary`, `InstanceFailoverGroupRolePrimary`, `InstanceFailoverGroupRoleSecondary`
- New enum type `JobStatus` with values `JobStatusFailed`, `JobStatusInProgress`, `JobStatusNotStarted`, `JobStatusSucceeded`
- New enum type `LastExecutionStatus` with values `LastExecutionStatusCompleted`, `LastExecutionStatusFailed`, `LastExecutionStatusFaulted`, `LastExecutionStatusPostponed`, `LastExecutionStatusRescheduled`, `LastExecutionStatusSucceeded`
- New enum type `LicenseCategory` with values `LicenseCategoryCore`
- New enum type `MiLinkAssessmentCategory` with values `MiLinkAssessmentCategoryBoxToMiNetworkConnectivity`, `MiLinkAssessmentCategoryCertificates`, `MiLinkAssessmentCategoryDagCrossValidation`, `MiLinkAssessmentCategoryManagedInstance`, `MiLinkAssessmentCategoryManagedInstanceCrossValidation`, `MiLinkAssessmentCategoryManagedInstanceDatabase`, `MiLinkAssessmentCategoryMiToBoxNetworkConnectivity`, `MiLinkAssessmentCategorySQLInstance`, `MiLinkAssessmentCategorySQLInstanceAg`, `MiLinkAssessmentCategorySQLInstanceDatabase`
- New enum type `MigrationMode` with values `MigrationModeLogShipping`, `MigrationModeLogical`, `MigrationModeMILink`, `MigrationModeOther`, `MigrationModeUnknown`
- New enum type `MigrationStatus` with values `MigrationStatusCancelled`, `MigrationStatusFailed`, `MigrationStatusInProgress`, `MigrationStatusInProgressWithWarnings`, `MigrationStatusNotStarted`, `MigrationStatusSuccessful`, `MigrationStatusUnknown`
- New enum type `Mode` with values `ModeMixed`, `ModeUndefined`, `ModeWindows`
- New enum type `Origin` with values `OriginSystem`, `OriginUser`, `OriginUserSystem`
- New enum type `PrimaryAllowConnections` with values `PrimaryAllowConnectionsALL`, `PrimaryAllowConnectionsReadWrite`
- New enum type `ProvisioningState` with values `ProvisioningStateAccepted`, `ProvisioningStateCanceled`, `ProvisioningStateFailed`, `ProvisioningStateSucceeded`
- New enum type `RecommendationStatus` with values `RecommendationStatusNotReady`, `RecommendationStatusReady`, `RecommendationStatusReadyWithConditions`, `RecommendationStatusUnknown`
- New enum type `RecoveryMode` with values `RecoveryModeBulkLogged`, `RecoveryModeFull`, `RecoveryModeSimple`
- New enum type `ReplicationPartnerType` with values `ReplicationPartnerTypeAzureSQLManagedInstance`, `ReplicationPartnerTypeAzureSqlvm`, `ReplicationPartnerTypeSQLServer`, `ReplicationPartnerTypeUnknown`
- New enum type `ResourceUpdateMode` with values `ResourceUpdateModeSkipResourceUpdate`, `ResourceUpdateModeUpdateAllTargetRecommendationDetails`
- New enum type `Result` with values `ResultFailed`, `ResultNotCompleted`, `ResultSkipped`, `ResultSucceeded`, `ResultTimedOut`
- New enum type `Role` with values `RoleALL`, `RoleNone`, `RolePartner`, `RoleWitness`
- New enum type `SQLServerInstanceBpaColumnType` with values `SQLServerInstanceBpaColumnTypeBool`, `SQLServerInstanceBpaColumnTypeDatetime`, `SQLServerInstanceBpaColumnTypeDouble`, `SQLServerInstanceBpaColumnTypeGUID`, `SQLServerInstanceBpaColumnTypeInt`, `SQLServerInstanceBpaColumnTypeLong`, `SQLServerInstanceBpaColumnTypeString`, `SQLServerInstanceBpaColumnTypeTimespan`
- New enum type `SQLServerInstanceBpaQueryType` with values `SQLServerInstanceBpaQueryTypeBasic`, `SQLServerInstanceBpaQueryTypeHistoricalTrends`
- New enum type `SQLServerInstanceBpaReportType` with values `SQLServerInstanceBpaReportTypeAssessmentDataPoint`, `SQLServerInstanceBpaReportTypeAssessmentSummary`
- New enum type `SQLServerInstanceTargetRecommendationReportSectionType` with values `SQLServerInstanceTargetRecommendationReportSectionTypeFileRequirementsPerDatabase`, `SQLServerInstanceTargetRecommendationReportSectionTypeRequirementsPerDatabase`, `SQLServerInstanceTargetRecommendationReportSectionTypeRequirementsPerInstance`, `SQLServerInstanceTargetRecommendationReportSectionTypeSQLDbTargetRecommendationPerDatabase`, `SQLServerInstanceTargetRecommendationReportSectionTypeSQLMiTargetRecommendationPerInstance`, `SQLServerInstanceTargetRecommendationReportSectionTypeSQLVMTargetRecommendationPerInstance`
- New enum type `SQLServerInstanceTelemetryColumnType` with values `SQLServerInstanceTelemetryColumnTypeBool`, `SQLServerInstanceTelemetryColumnTypeDatetime`, `SQLServerInstanceTelemetryColumnTypeDouble`, `SQLServerInstanceTelemetryColumnTypeGUID`, `SQLServerInstanceTelemetryColumnTypeInt`, `SQLServerInstanceTelemetryColumnTypeLong`, `SQLServerInstanceTelemetryColumnTypeString`, `SQLServerInstanceTelemetryColumnTypeTimespan`
- New enum type `ScopeType` with values `ScopeTypeResourceGroup`, `ScopeTypeSubscription`, `ScopeTypeTenant`
- New enum type `SecondaryAllowConnections` with values `SecondaryAllowConnectionsALL`, `SecondaryAllowConnectionsNO`, `SecondaryAllowConnectionsReadOnly`
- New enum type `SeedingMode` with values `SeedingModeAutomatic`, `SeedingModeManual`
- New enum type `SequencerState` with values `SequencerStateCompleted`, `SequencerStateCreatingSuccessors`, `SequencerStateExecutingAction`, `SequencerStateNotStarted`, `SequencerStateWaitingPredecessors`
- New enum type `ServiceType` with values `ServiceTypeEngine`, `ServiceTypePbirs`, `ServiceTypeSsas`, `ServiceTypeSsis`, `ServiceTypeSsrs`
- New enum type `State` with values `StateActive`, `StateCompleted`, `StateDeleted`, `StateDisabled`, `StateEnabled`, `StateFaulted`, `StateInactive`, `StateSuspended`, `StateTerminated`
- New enum type `TargetType` with values `TargetTypeAzureSQLDatabase`, `TargetTypeAzureSQLManagedInstance`, `TargetTypeAzureSQLVirtualMachine`
- New enum type `Version` with values `VersionSQLServer2012`, `VersionSQLServer2014`, `VersionSQLServer2016`
- New function `*ClientFactory.NewFailoverGroupsClient() *FailoverGroupsClient`
- New function `*ClientFactory.NewSQLServerAvailabilityGroupsClient() *SQLServerAvailabilityGroupsClient`
- New function `*ClientFactory.NewSQLServerDatabasesClient() *SQLServerDatabasesClient`
- New function `*ClientFactory.NewSQLServerEsuLicensesClient() *SQLServerEsuLicensesClient`
- New function `*ClientFactory.NewSQLServerLicensesClient() *SQLServerLicensesClient`
- New function `NewFailoverGroupsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*FailoverGroupsClient, error)`
- New function `*FailoverGroupsClient.BeginCreate(ctx context.Context, resourceGroupName string, sqlManagedInstanceName string, failoverGroupName string, failoverGroupResource FailoverGroupResource, options *FailoverGroupsClientBeginCreateOptions) (*runtime.Poller[FailoverGroupsClientCreateResponse], error)`
- New function `*FailoverGroupsClient.BeginDelete(ctx context.Context, resourceGroupName string, sqlManagedInstanceName string, failoverGroupName string, options *FailoverGroupsClientBeginDeleteOptions) (*runtime.Poller[FailoverGroupsClientDeleteResponse], error)`
- New function `*FailoverGroupsClient.Get(ctx context.Context, resourceGroupName string, sqlManagedInstanceName string, failoverGroupName string, options *FailoverGroupsClientGetOptions) (FailoverGroupsClientGetResponse, error)`
- New function `*FailoverGroupsClient.NewListPager(resourceGroupName string, sqlManagedInstanceName string, options *FailoverGroupsClientListOptions) *runtime.Pager[FailoverGroupsClientListResponse]`
- New function `NewSQLServerAvailabilityGroupsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SQLServerAvailabilityGroupsClient, error)`
- New function `*SQLServerAvailabilityGroupsClient.AddDatabases(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, availabilityGroupName string, databases Databases, options *SQLServerAvailabilityGroupsClientAddDatabasesOptions) (SQLServerAvailabilityGroupsClientAddDatabasesResponse, error)`
- New function `*SQLServerAvailabilityGroupsClient.Create(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, availabilityGroupName string, sqlServerAvailabilityGroupResource SQLServerAvailabilityGroupResource, options *SQLServerAvailabilityGroupsClientCreateOptions) (SQLServerAvailabilityGroupsClientCreateResponse, error)`
- New function `*SQLServerAvailabilityGroupsClient.BeginCreateAvailabilityGroup(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, createAgConfiguration AvailabilityGroupCreateUpdateConfiguration, options *SQLServerAvailabilityGroupsClientBeginCreateAvailabilityGroupOptions) (*runtime.Poller[SQLServerAvailabilityGroupsClientCreateAvailabilityGroupResponse], error)`
- New function `*SQLServerAvailabilityGroupsClient.BeginCreateDistributedAvailabilityGroup(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, createDagConfiguration DistributedAvailabilityGroupCreateUpdateConfiguration, options *SQLServerAvailabilityGroupsClientBeginCreateDistributedAvailabilityGroupOptions) (*runtime.Poller[SQLServerAvailabilityGroupsClientCreateDistributedAvailabilityGroupResponse], error)`
- New function `*SQLServerAvailabilityGroupsClient.BeginCreateManagedInstanceLink(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, createManagedInstanceLinkConfiguration ManagedInstanceLinkCreateUpdateConfiguration, options *SQLServerAvailabilityGroupsClientBeginCreateManagedInstanceLinkOptions) (*runtime.Poller[SQLServerAvailabilityGroupsClientCreateManagedInstanceLinkResponse], error)`
- New function `*SQLServerAvailabilityGroupsClient.BeginDelete(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, availabilityGroupName string, options *SQLServerAvailabilityGroupsClientBeginDeleteOptions) (*runtime.Poller[SQLServerAvailabilityGroupsClientDeleteResponse], error)`
- New function `*SQLServerAvailabilityGroupsClient.BeginDeleteMiLink(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, availabilityGroupName string, options *SQLServerAvailabilityGroupsClientBeginDeleteMiLinkOptions) (*runtime.Poller[SQLServerAvailabilityGroupsClientDeleteMiLinkResponse], error)`
- New function `*SQLServerAvailabilityGroupsClient.DetailView(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, availabilityGroupName string, options *SQLServerAvailabilityGroupsClientDetailViewOptions) (SQLServerAvailabilityGroupsClientDetailViewResponse, error)`
- New function `*SQLServerAvailabilityGroupsClient.Failover(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, availabilityGroupName string, options *SQLServerAvailabilityGroupsClientFailoverOptions) (SQLServerAvailabilityGroupsClientFailoverResponse, error)`
- New function `*SQLServerAvailabilityGroupsClient.BeginFailoverMiLink(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, availabilityGroupName string, managedInstanceResourceID FailoverMiLinkResourceID, options *SQLServerAvailabilityGroupsClientBeginFailoverMiLinkOptions) (*runtime.Poller[SQLServerAvailabilityGroupsClientFailoverMiLinkResponse], error)`
- New function `*SQLServerAvailabilityGroupsClient.ForceFailoverAllowDataLoss(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, availabilityGroupName string, options *SQLServerAvailabilityGroupsClientForceFailoverAllowDataLossOptions) (SQLServerAvailabilityGroupsClientForceFailoverAllowDataLossResponse, error)`
- New function `*SQLServerAvailabilityGroupsClient.Get(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, availabilityGroupName string, options *SQLServerAvailabilityGroupsClientGetOptions) (SQLServerAvailabilityGroupsClientGetResponse, error)`
- New function `*SQLServerAvailabilityGroupsClient.NewListPager(resourceGroupName string, sqlServerInstanceName string, options *SQLServerAvailabilityGroupsClientListOptions) *runtime.Pager[SQLServerAvailabilityGroupsClientListResponse]`
- New function `*SQLServerAvailabilityGroupsClient.RemoveDatabases(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, availabilityGroupName string, databases Databases, options *SQLServerAvailabilityGroupsClientRemoveDatabasesOptions) (SQLServerAvailabilityGroupsClientRemoveDatabasesResponse, error)`
- New function `*SQLServerAvailabilityGroupsClient.BeginUpdate(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, availabilityGroupName string, sqlServerAvailabilityGroupUpdate SQLServerAvailabilityGroupUpdate, options *SQLServerAvailabilityGroupsClientBeginUpdateOptions) (*runtime.Poller[SQLServerAvailabilityGroupsClientUpdateResponse], error)`
- New function `NewSQLServerDatabasesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SQLServerDatabasesClient, error)`
- New function `*SQLServerDatabasesClient.Create(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, databaseName string, sqlServerDatabaseResource SQLServerDatabaseResource, options *SQLServerDatabasesClientCreateOptions) (SQLServerDatabasesClientCreateResponse, error)`
- New function `*SQLServerDatabasesClient.BeginDelete(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, databaseName string, options *SQLServerDatabasesClientBeginDeleteOptions) (*runtime.Poller[SQLServerDatabasesClientDeleteResponse], error)`
- New function `*SQLServerDatabasesClient.Get(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, databaseName string, options *SQLServerDatabasesClientGetOptions) (SQLServerDatabasesClientGetResponse, error)`
- New function `*SQLServerDatabasesClient.NewListPager(resourceGroupName string, sqlServerInstanceName string, options *SQLServerDatabasesClientListOptions) *runtime.Pager[SQLServerDatabasesClientListResponse]`
- New function `*SQLServerDatabasesClient.BeginUpdate(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, databaseName string, sqlServerDatabaseUpdate SQLServerDatabaseUpdate, options *SQLServerDatabasesClientBeginUpdateOptions) (*runtime.Poller[SQLServerDatabasesClientUpdateResponse], error)`
- New function `NewSQLServerEsuLicensesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SQLServerEsuLicensesClient, error)`
- New function `*SQLServerEsuLicensesClient.Create(ctx context.Context, resourceGroupName string, sqlServerEsuLicenseName string, sqlServerEsuLicense SQLServerEsuLicense, options *SQLServerEsuLicensesClientCreateOptions) (SQLServerEsuLicensesClientCreateResponse, error)`
- New function `*SQLServerEsuLicensesClient.Delete(ctx context.Context, resourceGroupName string, sqlServerEsuLicenseName string, options *SQLServerEsuLicensesClientDeleteOptions) (SQLServerEsuLicensesClientDeleteResponse, error)`
- New function `*SQLServerEsuLicensesClient.Get(ctx context.Context, resourceGroupName string, sqlServerEsuLicenseName string, options *SQLServerEsuLicensesClientGetOptions) (SQLServerEsuLicensesClientGetResponse, error)`
- New function `*SQLServerEsuLicensesClient.NewListByResourceGroupPager(resourceGroupName string, options *SQLServerEsuLicensesClientListByResourceGroupOptions) *runtime.Pager[SQLServerEsuLicensesClientListByResourceGroupResponse]`
- New function `*SQLServerEsuLicensesClient.NewListPager(options *SQLServerEsuLicensesClientListOptions) *runtime.Pager[SQLServerEsuLicensesClientListResponse]`
- New function `*SQLServerEsuLicensesClient.Update(ctx context.Context, resourceGroupName string, sqlServerEsuLicenseName string, parameters SQLServerEsuLicenseUpdate, options *SQLServerEsuLicensesClientUpdateOptions) (SQLServerEsuLicensesClientUpdateResponse, error)`
- New function `*SQLServerInstancesClient.NewGetAllAvailabilityGroupsPager(resourceGroupName string, sqlServerInstanceName string, options *SQLServerInstancesClientGetAllAvailabilityGroupsOptions) *runtime.Pager[SQLServerInstancesClientGetAllAvailabilityGroupsResponse]`
- New function `*SQLServerInstancesClient.BeginGetBestPracticesAssessment(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, sqlServerInstanceBpaRequest SQLServerInstanceBpaRequest, options *SQLServerInstancesClientBeginGetBestPracticesAssessmentOptions) (*runtime.Poller[*runtime.Pager[SQLServerInstancesClientGetBestPracticesAssessmentResponse]], error)`
- New function `*SQLServerInstancesClient.BeginGetJobs(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, options *SQLServerInstancesClientBeginGetJobsOptions) (*runtime.Poller[SQLServerInstancesClientGetJobsResponse], error)`
- New function `*SQLServerInstancesClient.GetJobsStatus(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, options *SQLServerInstancesClientGetJobsStatusOptions) (SQLServerInstancesClientGetJobsStatusResponse, error)`
- New function `*SQLServerInstancesClient.BeginGetMigrationReadinessReport(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, options *SQLServerInstancesClientBeginGetMigrationReadinessReportOptions) (*runtime.Poller[SQLServerInstancesClientGetMigrationReadinessReportResponse], error)`
- New function `*SQLServerInstancesClient.BeginGetTargetRecommendationReports(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, options *SQLServerInstancesClientBeginGetTargetRecommendationReportsOptions) (*runtime.Poller[SQLServerInstancesClientGetTargetRecommendationReportsResponse], error)`
- New function `*SQLServerInstancesClient.BeginGetTelemetry(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, sqlServerInstanceTelemetryRequest SQLServerInstanceTelemetryRequest, options *SQLServerInstancesClientBeginGetTelemetryOptions) (*runtime.Poller[*runtime.Pager[SQLServerInstancesClientGetTelemetryResponse]], error)`
- New function `*SQLServerInstancesClient.PostUpgrade(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, options *SQLServerInstancesClientPostUpgradeOptions) (SQLServerInstancesClientPostUpgradeResponse, error)`
- New function `*SQLServerInstancesClient.PreUpgrade(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, options *SQLServerInstancesClientPreUpgradeOptions) (SQLServerInstancesClientPreUpgradeResponse, error)`
- New function `*SQLServerInstancesClient.BeginRunBestPracticeAssessment(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, options *SQLServerInstancesClientBeginRunBestPracticeAssessmentOptions) (*runtime.Poller[SQLServerInstancesClientRunBestPracticeAssessmentResponse], error)`
- New function `*SQLServerInstancesClient.RunBestPracticesAssessment(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, options *SQLServerInstancesClientRunBestPracticesAssessmentOptions) (SQLServerInstancesClientRunBestPracticesAssessmentResponse, error)`
- New function `*SQLServerInstancesClient.BeginRunManagedInstanceLinkAssessment(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, sqlServerInstanceManagedInstanceLinkAssessmentRequest SQLServerInstanceManagedInstanceLinkAssessmentRequest, options *SQLServerInstancesClientBeginRunManagedInstanceLinkAssessmentOptions) (*runtime.Poller[SQLServerInstancesClientRunManagedInstanceLinkAssessmentResponse], error)`
- New function `*SQLServerInstancesClient.RunMigrationAssessment(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, options *SQLServerInstancesClientRunMigrationAssessmentOptions) (SQLServerInstancesClientRunMigrationAssessmentResponse, error)`
- New function `*SQLServerInstancesClient.BeginRunMigrationReadinessAssessment(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, options *SQLServerInstancesClientBeginRunMigrationReadinessAssessmentOptions) (*runtime.Poller[SQLServerInstancesClientRunMigrationReadinessAssessmentResponse], error)`
- New function `*SQLServerInstancesClient.BeginRunTargetRecommendationJob(ctx context.Context, resourceGroupName string, sqlServerInstanceName string, options *SQLServerInstancesClientBeginRunTargetRecommendationJobOptions) (*runtime.Poller[SQLServerInstancesClientRunTargetRecommendationJobResponse], error)`
- New function `NewSQLServerLicensesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SQLServerLicensesClient, error)`
- New function `*SQLServerLicensesClient.Create(ctx context.Context, resourceGroupName string, sqlServerLicenseName string, sqlServerLicense SQLServerLicense, options *SQLServerLicensesClientCreateOptions) (SQLServerLicensesClientCreateResponse, error)`
- New function `*SQLServerLicensesClient.Delete(ctx context.Context, resourceGroupName string, sqlServerLicenseName string, options *SQLServerLicensesClientDeleteOptions) (SQLServerLicensesClientDeleteResponse, error)`
- New function `*SQLServerLicensesClient.Get(ctx context.Context, resourceGroupName string, sqlServerLicenseName string, options *SQLServerLicensesClientGetOptions) (SQLServerLicensesClientGetResponse, error)`
- New function `*SQLServerLicensesClient.NewListByResourceGroupPager(resourceGroupName string, options *SQLServerLicensesClientListByResourceGroupOptions) *runtime.Pager[SQLServerLicensesClientListByResourceGroupResponse]`
- New function `*SQLServerLicensesClient.NewListPager(options *SQLServerLicensesClientListOptions) *runtime.Pager[SQLServerLicensesClientListResponse]`
- New function `*SQLServerLicensesClient.Update(ctx context.Context, resourceGroupName string, sqlServerLicenseName string, parameters SQLServerLicenseUpdate, options *SQLServerLicensesClientUpdateOptions) (SQLServerLicensesClientUpdateResponse, error)`
- New struct `AdditionalMigrationJobAttributes`
- New struct `ArcSQLServerAvailabilityGroupListResult`
- New struct `ArcSQLServerDatabaseListResult`
- New struct `Authentication`
- New struct `AvailabilityGroupConfigure`
- New struct `AvailabilityGroupCreateUpdateConfiguration`
- New struct `AvailabilityGroupCreateUpdateReplicaConfiguration`
- New struct `AvailabilityGroupInfo`
- New struct `AvailabilityGroupRetrievalFilters`
- New struct `AvailabilityGroupState`
- New struct `BackgroundJob`
- New struct `BackupPolicy`
- New struct `BestPracticesAssessment`
- New struct `ClientConnection`
- New struct `CostOptionSelectedValues`
- New struct `CostTypeValues`
- New struct `CronTrigger`
- New struct `DBMEndpoint`
- New struct `DataBaseMigration`
- New struct `DataBaseMigrationAssessment`
- New struct `DatabaseAssessmentsItem`
- New struct `DatabaseMigrationJobsItem`
- New struct `Databases`
- New struct `DiskSizes`
- New struct `DistributedAvailabilityGroupCreateUpdateAvailabilityGroupCertificateConfiguration`
- New struct `DistributedAvailabilityGroupCreateUpdateAvailabilityGroupConfiguration`
- New struct `DistributedAvailabilityGroupCreateUpdateConfiguration`
- New struct `EntraAuthentication`
- New struct `FailoverCluster`
- New struct `FailoverGroupListResult`
- New struct `FailoverGroupProperties`
- New struct `FailoverGroupResource`
- New struct `FailoverGroupSpec`
- New struct `FailoverMiLinkResourceID`
- New struct `HostIPAddressInformation`
- New struct `ImpactedObjectsInfo`
- New struct `ImpactedObjectsSuitabilitySummary`
- New struct `K8SActiveDirectory`
- New struct `K8SActiveDirectoryConnector`
- New struct `K8SNetworkSettings`
- New struct `K8SSecurity`
- New struct `K8SSettings`
- New struct `K8StransparentDataEncryption`
- New struct `ManagedInstanceLinkCreateUpdateConfiguration`
- New struct `MiLinkCreateUpdateConfiguration`
- New struct `Migration`
- New struct `MigrationAssessment`
- New struct `MigrationAssessmentSettings`
- New struct `Monitoring`
- New struct `SKURecommendationResults`
- New struct `SKURecommendationResultsAzureSQLDatabase`
- New struct `SKURecommendationResultsAzureSQLDatabaseTargetSKU`
- New struct `SKURecommendationResultsAzureSQLDatabaseTargetSKUCategory`
- New struct `SKURecommendationResultsAzureSQLManagedInstance`
- New struct `SKURecommendationResultsAzureSQLManagedInstanceTargetSKU`
- New struct `SKURecommendationResultsAzureSQLManagedInstanceTargetSKUCategory`
- New struct `SKURecommendationResultsAzureSQLVirtualMachine`
- New struct `SKURecommendationResultsAzureSQLVirtualMachineTargetSKU`
- New struct `SKURecommendationResultsAzureSQLVirtualMachineTargetSKUCategory`
- New struct `SKURecommendationResultsAzureSQLVirtualMachineTargetSKUVirtualMachineSize`
- New struct `SKURecommendationResultsMonthlyCost`
- New struct `SKURecommendationResultsMonthlyCostOptionItem`
- New struct `SKURecommendationSummary`
- New struct `SKURecommendationSummaryTargetSKU`
- New struct `SKURecommendationSummaryTargetSKUCategory`
- New struct `SQLAvailabilityGroupDatabaseReplicaResourceProperties`
- New struct `SQLAvailabilityGroupIPV4AddressesAndMasksPropertiesItem`
- New struct `SQLAvailabilityGroupReplicaResourceProperties`
- New struct `SQLAvailabilityGroupStaticIPListenerProperties`
- New struct `SQLServerAvailabilityGroupResource`
- New struct `SQLServerAvailabilityGroupResourceProperties`
- New struct `SQLServerAvailabilityGroupResourcePropertiesDatabases`
- New struct `SQLServerAvailabilityGroupResourcePropertiesReplicas`
- New struct `SQLServerAvailabilityGroupUpdate`
- New struct `SQLServerDatabaseResource`
- New struct `SQLServerDatabaseResourceProperties`
- New struct `SQLServerDatabaseResourcePropertiesBackupInformation`
- New struct `SQLServerDatabaseResourcePropertiesDatabaseOptions`
- New struct `SQLServerDatabaseUpdate`
- New struct `SQLServerEsuLicense`
- New struct `SQLServerEsuLicenseListResult`
- New struct `SQLServerEsuLicenseProperties`
- New struct `SQLServerEsuLicenseUpdate`
- New struct `SQLServerEsuLicenseUpdateProperties`
- New struct `SQLServerInstanceBpaColumn`
- New struct `SQLServerInstanceBpaRequest`
- New struct `SQLServerInstanceBpaResponse`
- New struct `SQLServerInstanceJob`
- New struct `SQLServerInstanceJobStatus`
- New struct `SQLServerInstanceJobsRequest`
- New struct `SQLServerInstanceJobsResponse`
- New struct `SQLServerInstanceJobsStatusRequest`
- New struct `SQLServerInstanceJobsStatusResponse`
- New struct `SQLServerInstanceManagedInstanceLinkAssessment`
- New struct `SQLServerInstanceManagedInstanceLinkAssessmentRequest`
- New struct `SQLServerInstanceManagedInstanceLinkAssessmentResponse`
- New struct `SQLServerInstanceMigrationReadinessReportResponse`
- New struct `SQLServerInstanceRunBestPracticesAssessmentResponse`
- New struct `SQLServerInstanceRunMigrationAssessmentResponse`
- New struct `SQLServerInstanceRunMigrationReadinessAssessmentResponse`
- New struct `SQLServerInstanceRunTargetRecommendationJobRequest`
- New struct `SQLServerInstanceRunTargetRecommendationJobResponse`
- New struct `SQLServerInstanceTargetRecommendationReport`
- New struct `SQLServerInstanceTargetRecommendationReportSection`
- New struct `SQLServerInstanceTargetRecommendationReportsRequest`
- New struct `SQLServerInstanceTargetRecommendationReportsResponse`
- New struct `SQLServerInstanceTelemetryColumn`
- New struct `SQLServerInstanceTelemetryRequest`
- New struct `SQLServerInstanceTelemetryResponse`
- New struct `SQLServerInstanceUpdateProperties`
- New struct `SQLServerInstancesClientGetBestPracticesAssessmentResponse`
- New struct `SQLServerInstancesClientGetTelemetryResponse`
- New struct `SQLServerLicense`
- New struct `SQLServerLicenseListResult`
- New struct `SQLServerLicenseProperties`
- New struct `SQLServerLicenseUpdate`
- New struct `SQLServerLicenseUpdateProperties`
- New struct `Schedule`
- New struct `SequencerAction`
- New struct `ServerAssessmentsItem`
- New struct `ServerAssessmentsPropertiesItemsItem`
- New struct `TargetReadiness`
- New field `ActionType` in struct `Operation`
- New field `Security`, `Settings` in struct `SQLManagedInstanceK8SSpec`
- New field `AlwaysOnRole`, `Authentication`, `BackupPolicy`, `BestPracticesAssessment`, `ClientConnection`, `Cores`, `DatabaseMirroringEndpoint`, `DbMasterKeyExists`, `DiscoverySource`, `FailoverCluster`, `IsDigiCertPkiCertTrustConfigured`, `IsHadrEnabled`, `IsMicrosoftPkiCertTrustConfigured`, `LastInventoryUploadTime`, `LastUsageUploadTime`, `MaxServerMemoryMB`, `Migration`, `Monitoring`, `ServiceType`, `TraceFlags`, `UpgradeLockedUntil`, `VMID` in struct `SQLServerInstanceProperties`
- New field `Properties` in struct `SQLServerInstanceUpdate`


## 0.7.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.


## 0.6.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.

## 0.6.0 (2023-03-27)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.5.0 (2022-05-17)

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/azurearcdata/armazurearcdata` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html) since version 0.5.0, which contains breaking changes.

To migrate the existing applications to the latest version, please refer to [Migration Guide](https://aka.ms/azsdk/go/mgmt/migration).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).