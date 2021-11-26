# Unreleased

## Additive Changes

### New Constants

1. EncryptionType.EncryptionTypeDouble
1. EncryptionType.EncryptionTypeSingle
1. MetricAggregationType.MetricAggregationTypeAverage

### New Funcs

1. AccountsClient.ListBySubscription(context.Context) (AccountListPage, error)
1. AccountsClient.ListBySubscriptionComplete(context.Context) (AccountListIterator, error)
1. AccountsClient.ListBySubscriptionPreparer(context.Context) (*http.Request, error)
1. AccountsClient.ListBySubscriptionResponder(*http.Response) (AccountList, error)
1. AccountsClient.ListBySubscriptionSender(*http.Request) (*http.Response, error)
1. PossibleEncryptionTypeValues() []EncryptionType
1. PossibleMetricAggregationTypeValues() []MetricAggregationType

### Struct Changes

#### New Structs

1. LogSpecification

#### New Struct Fields

1. Account.Etag
1. Backup.SystemData
1. BackupPolicy.Etag
1. BackupPolicy.SystemData
1. BackupPolicyProperties.BackupPolicyID
1. CapacityPool.Etag
1. MetricSpecification.EnableRegionalMdmAccount
1. MetricSpecification.InternalMetricName
1. MetricSpecification.IsInternal
1. MetricSpecification.SourceMdmAccount
1. MetricSpecification.SourceMdmNamespace
1. MetricSpecification.SupportedAggregationTypes
1. MetricSpecification.SupportedTimeGrainTypes
1. PoolProperties.EncryptionType
1. ServiceSpecification.LogSpecifications
1. SnapshotPolicy.Etag
1. Volume.Etag
