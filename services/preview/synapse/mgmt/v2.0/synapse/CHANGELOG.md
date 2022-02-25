# Change History

## Breaking Changes

### Removed Constants

1. StorageAccountType.StorageAccountTypeZRS

## Additive Changes

### New Constants

1. Kind.KindReadOnlyFollowing
1. PrincipalsModificationKind.PrincipalsModificationKindNone
1. PrincipalsModificationKind.PrincipalsModificationKindReplace
1. PrincipalsModificationKind.PrincipalsModificationKindUnion

### New Funcs

1. *ReadOnlyFollowingDatabase.UnmarshalJSON([]byte) error
1. Database.AsReadOnlyFollowingDatabase() (*ReadOnlyFollowingDatabase, bool)
1. PossiblePrincipalsModificationKindValues() []PrincipalsModificationKind
1. ReadOnlyFollowingDatabase.AsBasicDatabase() (BasicDatabase, bool)
1. ReadOnlyFollowingDatabase.AsDatabase() (*Database, bool)
1. ReadOnlyFollowingDatabase.AsReadOnlyFollowingDatabase() (*ReadOnlyFollowingDatabase, bool)
1. ReadOnlyFollowingDatabase.AsReadWriteDatabase() (*ReadWriteDatabase, bool)
1. ReadOnlyFollowingDatabase.MarshalJSON() ([]byte, error)
1. ReadOnlyFollowingDatabaseProperties.MarshalJSON() ([]byte, error)
1. ReadWriteDatabase.AsReadOnlyFollowingDatabase() (*ReadOnlyFollowingDatabase, bool)
1. SQLPoolResourceProperties.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. ManagedIntegrationRuntimeManagedVirtualNetworkReference
1. ReadOnlyFollowingDatabase
1. ReadOnlyFollowingDatabaseProperties

#### New Struct Fields

1. DynamicExecutorAllocation.MaxExecutors
1. DynamicExecutorAllocation.MinExecutors
1. ExtendedServerBlobAuditingPolicyProperties.IsDevopsAuditEnabled
1. ManagedIntegrationRuntime.*ManagedIntegrationRuntimeManagedVirtualNetworkReference
1. SelfHostedIntegrationRuntimeStatusTypeProperties.NewerVersions
1. SelfHostedIntegrationRuntimeStatusTypeProperties.ServiceRegion
1. ServerBlobAuditingPolicyProperties.IsDevopsAuditEnabled
1. WorkspaceProperties.TrustedServiceBypassEnabled
