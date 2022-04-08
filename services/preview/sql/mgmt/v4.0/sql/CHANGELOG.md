# Unreleased

## Breaking Changes

### Removed Constants

1. DatabaseState2.DatabaseState2All
1. DatabaseState2.DatabaseState2Deleted
1. DatabaseState2.DatabaseState2Live
1. DatabaseState3.DatabaseState3All
1. DatabaseState3.DatabaseState3Deleted
1. DatabaseState3.DatabaseState3Live
1. DatabaseState4.DatabaseState4All
1. DatabaseState4.DatabaseState4Deleted
1. DatabaseState4.DatabaseState4Live
1. DatabaseState5.DatabaseState5All
1. DatabaseState5.DatabaseState5Deleted
1. DatabaseState5.DatabaseState5Live
1. DatabaseState6.DatabaseState6All
1. DatabaseState6.DatabaseState6Deleted
1. DatabaseState6.DatabaseState6Live
1. LongTermRetentionDatabaseState.LongTermRetentionDatabaseStateAll
1. LongTermRetentionDatabaseState.LongTermRetentionDatabaseStateDeleted
1. LongTermRetentionDatabaseState.LongTermRetentionDatabaseStateLive
1. SecondaryType.Named

### Removed Funcs

1. *BackupLongTermRetentionPoliciesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *BackupLongTermRetentionPolicy.UnmarshalJSON([]byte) error
1. BackupLongTermRetentionPoliciesClient.CreateOrUpdate(context.Context, string, string, string, BackupLongTermRetentionPolicy) (BackupLongTermRetentionPoliciesCreateOrUpdateFuture, error)
1. BackupLongTermRetentionPoliciesClient.CreateOrUpdatePreparer(context.Context, string, string, string, BackupLongTermRetentionPolicy) (*http.Request, error)
1. BackupLongTermRetentionPoliciesClient.CreateOrUpdateResponder(*http.Response) (BackupLongTermRetentionPolicy, error)
1. BackupLongTermRetentionPoliciesClient.CreateOrUpdateSender(*http.Request) (BackupLongTermRetentionPoliciesCreateOrUpdateFuture, error)
1. BackupLongTermRetentionPoliciesClient.Get(context.Context, string, string, string) (BackupLongTermRetentionPolicy, error)
1. BackupLongTermRetentionPoliciesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. BackupLongTermRetentionPoliciesClient.GetResponder(*http.Response) (BackupLongTermRetentionPolicy, error)
1. BackupLongTermRetentionPoliciesClient.GetSender(*http.Request) (*http.Response, error)
1. BackupLongTermRetentionPoliciesClient.ListByDatabase(context.Context, string, string, string) (BackupLongTermRetentionPolicy, error)
1. BackupLongTermRetentionPoliciesClient.ListByDatabasePreparer(context.Context, string, string, string) (*http.Request, error)
1. BackupLongTermRetentionPoliciesClient.ListByDatabaseResponder(*http.Response) (BackupLongTermRetentionPolicy, error)
1. BackupLongTermRetentionPoliciesClient.ListByDatabaseSender(*http.Request) (*http.Response, error)
1. BackupLongTermRetentionPolicy.MarshalJSON() ([]byte, error)
1. NewBackupLongTermRetentionPoliciesClient(string) BackupLongTermRetentionPoliciesClient
1. NewBackupLongTermRetentionPoliciesClientWithBaseURI(string, string) BackupLongTermRetentionPoliciesClient
1. PossibleDatabaseState1Values() []DatabaseState1
1. PossibleDatabaseState2Values() []DatabaseState2
1. PossibleDatabaseState3Values() []DatabaseState3
1. PossibleDatabaseState4Values() []DatabaseState4
1. PossibleDatabaseState5Values() []DatabaseState5
1. PossibleDatabaseState6Values() []DatabaseState6
1. PossibleLongTermRetentionDatabaseStateValues() []LongTermRetentionDatabaseState
1. SystemData.MarshalJSON() ([]byte, error)

### Struct Changes

#### Removed Structs

1. BackupLongTermRetentionPoliciesClient
1. BackupLongTermRetentionPoliciesCreateOrUpdateFuture
1. BackupLongTermRetentionPolicy
1. LongTermRetentionPolicyProperties

#### Removed Struct Fields

1. DatabaseProperties.StorageAccountType

### Signature Changes

#### Const Types

1. All changed type from DatabaseState1 to DatabaseState
1. Deleted changed type from DatabaseState1 to DatabaseState
1. Geo changed type from SecondaryType to BackupStorageRedundancy
1. Live changed type from DatabaseState1 to DatabaseState

#### Funcs

1. LongTermRetentionBackupsClient.ListByDatabase
	- Params
		- From: context.Context, string, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByDatabaseComplete
	- Params
		- From: context.Context, string, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByDatabasePreparer
	- Params
		- From: context.Context, string, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByLocation
	- Params
		- From: context.Context, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByLocationComplete
	- Params
		- From: context.Context, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByLocationPreparer
	- Params
		- From: context.Context, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByResourceGroupDatabase
	- Params
		- From: context.Context, string, string, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByResourceGroupDatabaseComplete
	- Params
		- From: context.Context, string, string, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByResourceGroupDatabasePreparer
	- Params
		- From: context.Context, string, string, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByResourceGroupLocation
	- Params
		- From: context.Context, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByResourceGroupLocationComplete
	- Params
		- From: context.Context, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByResourceGroupLocationPreparer
	- Params
		- From: context.Context, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByResourceGroupServer
	- Params
		- From: context.Context, string, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByResourceGroupServerComplete
	- Params
		- From: context.Context, string, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByResourceGroupServerPreparer
	- Params
		- From: context.Context, string, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByServer
	- Params
		- From: context.Context, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByServerComplete
	- Params
		- From: context.Context, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, *bool, DatabaseState
1. LongTermRetentionBackupsClient.ListByServerPreparer
	- Params
		- From: context.Context, string, string, *bool, LongTermRetentionDatabaseState
		- To: context.Context, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByDatabase
	- Params
		- From: context.Context, string, string, string, *bool, DatabaseState1
		- To: context.Context, string, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByDatabaseComplete
	- Params
		- From: context.Context, string, string, string, *bool, DatabaseState1
		- To: context.Context, string, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByDatabasePreparer
	- Params
		- From: context.Context, string, string, string, *bool, DatabaseState1
		- To: context.Context, string, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByInstance
	- Params
		- From: context.Context, string, string, *bool, DatabaseState2
		- To: context.Context, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByInstanceComplete
	- Params
		- From: context.Context, string, string, *bool, DatabaseState2
		- To: context.Context, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByInstancePreparer
	- Params
		- From: context.Context, string, string, *bool, DatabaseState2
		- To: context.Context, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByLocation
	- Params
		- From: context.Context, string, *bool, DatabaseState3
		- To: context.Context, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByLocationComplete
	- Params
		- From: context.Context, string, *bool, DatabaseState3
		- To: context.Context, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByLocationPreparer
	- Params
		- From: context.Context, string, *bool, DatabaseState3
		- To: context.Context, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByResourceGroupDatabase
	- Params
		- From: context.Context, string, string, string, string, *bool, DatabaseState4
		- To: context.Context, string, string, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByResourceGroupDatabaseComplete
	- Params
		- From: context.Context, string, string, string, string, *bool, DatabaseState4
		- To: context.Context, string, string, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByResourceGroupDatabasePreparer
	- Params
		- From: context.Context, string, string, string, string, *bool, DatabaseState4
		- To: context.Context, string, string, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByResourceGroupInstance
	- Params
		- From: context.Context, string, string, string, *bool, DatabaseState5
		- To: context.Context, string, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByResourceGroupInstanceComplete
	- Params
		- From: context.Context, string, string, string, *bool, DatabaseState5
		- To: context.Context, string, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByResourceGroupInstancePreparer
	- Params
		- From: context.Context, string, string, string, *bool, DatabaseState5
		- To: context.Context, string, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByResourceGroupLocation
	- Params
		- From: context.Context, string, string, *bool, DatabaseState6
		- To: context.Context, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByResourceGroupLocationComplete
	- Params
		- From: context.Context, string, string, *bool, DatabaseState6
		- To: context.Context, string, string, *bool, DatabaseState
1. LongTermRetentionManagedInstanceBackupsClient.ListByResourceGroupLocationPreparer
	- Params
		- From: context.Context, string, string, *bool, DatabaseState6
		- To: context.Context, string, string, *bool, DatabaseState

## Additive Changes

### New Constants

1. BackupStorageRedundancy.Local
1. BackupStorageRedundancy.Zone
1. CurrentBackupStorageRedundancy.CurrentBackupStorageRedundancyGeo
1. CurrentBackupStorageRedundancy.CurrentBackupStorageRedundancyLocal
1. CurrentBackupStorageRedundancy.CurrentBackupStorageRedundancyZone
1. IdentityType.SystemAssignedUserAssigned
1. RequestedBackupStorageRedundancy.RequestedBackupStorageRedundancyGeo
1. RequestedBackupStorageRedundancy.RequestedBackupStorageRedundancyLocal
1. RequestedBackupStorageRedundancy.RequestedBackupStorageRedundancyZone
1. SecondaryType.SecondaryTypeGeo
1. SecondaryType.SecondaryTypeNamed
1. TargetBackupStorageRedundancy.TargetBackupStorageRedundancyGeo
1. TargetBackupStorageRedundancy.TargetBackupStorageRedundancyLocal
1. TargetBackupStorageRedundancy.TargetBackupStorageRedundancyZone

### New Funcs

1. *CopyLongTermRetentionBackupParameters.UnmarshalJSON([]byte) error
1. *LongTermRetentionBackupOperationResult.UnmarshalJSON([]byte) error
1. *LongTermRetentionBackupsCopyByResourceGroupFuture.UnmarshalJSON([]byte) error
1. *LongTermRetentionBackupsCopyFuture.UnmarshalJSON([]byte) error
1. *LongTermRetentionBackupsUpdateByResourceGroupFuture.UnmarshalJSON([]byte) error
1. *LongTermRetentionBackupsUpdateFuture.UnmarshalJSON([]byte) error
1. *LongTermRetentionPoliciesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *LongTermRetentionPolicy.UnmarshalJSON([]byte) error
1. *LongTermRetentionPolicyListResultIterator.Next() error
1. *LongTermRetentionPolicyListResultIterator.NextWithContext(context.Context) error
1. *LongTermRetentionPolicyListResultPage.Next() error
1. *LongTermRetentionPolicyListResultPage.NextWithContext(context.Context) error
1. *UpdateLongTermRetentionBackupParameters.UnmarshalJSON([]byte) error
1. CopyLongTermRetentionBackupParameters.MarshalJSON() ([]byte, error)
1. LongTermRetentionBackupOperationResult.MarshalJSON() ([]byte, error)
1. LongTermRetentionBackupsClient.Copy(context.Context, string, string, string, string, CopyLongTermRetentionBackupParameters) (LongTermRetentionBackupsCopyFuture, error)
1. LongTermRetentionBackupsClient.CopyByResourceGroup(context.Context, string, string, string, string, string, CopyLongTermRetentionBackupParameters) (LongTermRetentionBackupsCopyByResourceGroupFuture, error)
1. LongTermRetentionBackupsClient.CopyByResourceGroupPreparer(context.Context, string, string, string, string, string, CopyLongTermRetentionBackupParameters) (*http.Request, error)
1. LongTermRetentionBackupsClient.CopyByResourceGroupResponder(*http.Response) (LongTermRetentionBackupOperationResult, error)
1. LongTermRetentionBackupsClient.CopyByResourceGroupSender(*http.Request) (LongTermRetentionBackupsCopyByResourceGroupFuture, error)
1. LongTermRetentionBackupsClient.CopyPreparer(context.Context, string, string, string, string, CopyLongTermRetentionBackupParameters) (*http.Request, error)
1. LongTermRetentionBackupsClient.CopyResponder(*http.Response) (LongTermRetentionBackupOperationResult, error)
1. LongTermRetentionBackupsClient.CopySender(*http.Request) (LongTermRetentionBackupsCopyFuture, error)
1. LongTermRetentionBackupsClient.Update(context.Context, string, string, string, string, UpdateLongTermRetentionBackupParameters) (LongTermRetentionBackupsUpdateFuture, error)
1. LongTermRetentionBackupsClient.UpdateByResourceGroup(context.Context, string, string, string, string, string, UpdateLongTermRetentionBackupParameters) (LongTermRetentionBackupsUpdateByResourceGroupFuture, error)
1. LongTermRetentionBackupsClient.UpdateByResourceGroupPreparer(context.Context, string, string, string, string, string, UpdateLongTermRetentionBackupParameters) (*http.Request, error)
1. LongTermRetentionBackupsClient.UpdateByResourceGroupResponder(*http.Response) (LongTermRetentionBackupOperationResult, error)
1. LongTermRetentionBackupsClient.UpdateByResourceGroupSender(*http.Request) (LongTermRetentionBackupsUpdateByResourceGroupFuture, error)
1. LongTermRetentionBackupsClient.UpdatePreparer(context.Context, string, string, string, string, UpdateLongTermRetentionBackupParameters) (*http.Request, error)
1. LongTermRetentionBackupsClient.UpdateResponder(*http.Response) (LongTermRetentionBackupOperationResult, error)
1. LongTermRetentionBackupsClient.UpdateSender(*http.Request) (LongTermRetentionBackupsUpdateFuture, error)
1. LongTermRetentionOperationResultProperties.MarshalJSON() ([]byte, error)
1. LongTermRetentionPoliciesClient.CreateOrUpdate(context.Context, string, string, string, LongTermRetentionPolicy) (LongTermRetentionPoliciesCreateOrUpdateFuture, error)
1. LongTermRetentionPoliciesClient.CreateOrUpdatePreparer(context.Context, string, string, string, LongTermRetentionPolicy) (*http.Request, error)
1. LongTermRetentionPoliciesClient.CreateOrUpdateResponder(*http.Response) (LongTermRetentionPolicy, error)
1. LongTermRetentionPoliciesClient.CreateOrUpdateSender(*http.Request) (LongTermRetentionPoliciesCreateOrUpdateFuture, error)
1. LongTermRetentionPoliciesClient.Get(context.Context, string, string, string) (LongTermRetentionPolicy, error)
1. LongTermRetentionPoliciesClient.GetPreparer(context.Context, string, string, string) (*http.Request, error)
1. LongTermRetentionPoliciesClient.GetResponder(*http.Response) (LongTermRetentionPolicy, error)
1. LongTermRetentionPoliciesClient.GetSender(*http.Request) (*http.Response, error)
1. LongTermRetentionPoliciesClient.ListByDatabase(context.Context, string, string, string) (LongTermRetentionPolicyListResultPage, error)
1. LongTermRetentionPoliciesClient.ListByDatabaseComplete(context.Context, string, string, string) (LongTermRetentionPolicyListResultIterator, error)
1. LongTermRetentionPoliciesClient.ListByDatabasePreparer(context.Context, string, string, string) (*http.Request, error)
1. LongTermRetentionPoliciesClient.ListByDatabaseResponder(*http.Response) (LongTermRetentionPolicyListResult, error)
1. LongTermRetentionPoliciesClient.ListByDatabaseSender(*http.Request) (*http.Response, error)
1. LongTermRetentionPolicy.MarshalJSON() ([]byte, error)
1. LongTermRetentionPolicyListResult.IsEmpty() bool
1. LongTermRetentionPolicyListResult.MarshalJSON() ([]byte, error)
1. LongTermRetentionPolicyListResultIterator.NotDone() bool
1. LongTermRetentionPolicyListResultIterator.Response() LongTermRetentionPolicyListResult
1. LongTermRetentionPolicyListResultIterator.Value() LongTermRetentionPolicy
1. LongTermRetentionPolicyListResultPage.NotDone() bool
1. LongTermRetentionPolicyListResultPage.Response() LongTermRetentionPolicyListResult
1. LongTermRetentionPolicyListResultPage.Values() []LongTermRetentionPolicy
1. NewLongTermRetentionPoliciesClient(string) LongTermRetentionPoliciesClient
1. NewLongTermRetentionPoliciesClientWithBaseURI(string, string) LongTermRetentionPoliciesClient
1. NewLongTermRetentionPolicyListResultIterator(LongTermRetentionPolicyListResultPage) LongTermRetentionPolicyListResultIterator
1. NewLongTermRetentionPolicyListResultPage(LongTermRetentionPolicyListResult, func(context.Context, LongTermRetentionPolicyListResult) (LongTermRetentionPolicyListResult, error)) LongTermRetentionPolicyListResultPage
1. PossibleBackupStorageRedundancyValues() []BackupStorageRedundancy
1. PossibleCurrentBackupStorageRedundancyValues() []CurrentBackupStorageRedundancy
1. PossibleDatabaseStateValues() []DatabaseState
1. PossibleRequestedBackupStorageRedundancyValues() []RequestedBackupStorageRedundancy
1. PossibleTargetBackupStorageRedundancyValues() []TargetBackupStorageRedundancy
1. UpdateLongTermRetentionBackupParameters.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. CopyLongTermRetentionBackupParameters
1. CopyLongTermRetentionBackupParametersProperties
1. LongTermRetentionBackupOperationResult
1. LongTermRetentionBackupsCopyByResourceGroupFuture
1. LongTermRetentionBackupsCopyFuture
1. LongTermRetentionBackupsUpdateByResourceGroupFuture
1. LongTermRetentionBackupsUpdateFuture
1. LongTermRetentionOperationResultProperties
1. LongTermRetentionPoliciesClient
1. LongTermRetentionPoliciesCreateOrUpdateFuture
1. LongTermRetentionPolicy
1. LongTermRetentionPolicyListResult
1. LongTermRetentionPolicyListResultIterator
1. LongTermRetentionPolicyListResultPage
1. UpdateLongTermRetentionBackupParameters
1. UpdateLongTermRetentionBackupParametersProperties

#### New Struct Fields

1. DatabaseProperties.CurrentBackupStorageRedundancy
1. DatabaseProperties.RequestedBackupStorageRedundancy
1. LongTermRetentionBackupProperties.BackupStorageRedundancy
1. LongTermRetentionBackupProperties.RequestedBackupStorageRedundancy
