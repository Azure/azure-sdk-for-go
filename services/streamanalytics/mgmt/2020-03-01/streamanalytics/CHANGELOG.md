# Change History

## Additive Changes

### New Constants

1. TypeBasicOutputDataSource.TypeBasicOutputDataSourceTypeMicrosoftAzureFunction

### New Funcs

1. *AzureFunctionOutputDataSource.UnmarshalJSON([]byte) error
1. AzureDataLakeStoreOutputDataSource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsAzureDataLakeStoreOutputDataSource() (*AzureDataLakeStoreOutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsAzureSQLDatabaseOutputDataSource() (*AzureSQLDatabaseOutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsAzureSynapseOutputDataSource() (*AzureSynapseOutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsAzureTableOutputDataSource() (*AzureTableOutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsBasicOutputDataSource() (BasicOutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsBlobOutputDataSource() (*BlobOutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsDocumentDbOutputDataSource() (*DocumentDbOutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsEventHubOutputDataSource() (*EventHubOutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsEventHubV2OutputDataSource() (*EventHubV2OutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsOutputDataSource() (*OutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsPowerBIOutputDataSource() (*PowerBIOutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsServiceBusQueueOutputDataSource() (*ServiceBusQueueOutputDataSource, bool)
1. AzureFunctionOutputDataSource.AsServiceBusTopicOutputDataSource() (*ServiceBusTopicOutputDataSource, bool)
1. AzureFunctionOutputDataSource.MarshalJSON() ([]byte, error)
1. AzureSQLDatabaseOutputDataSource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)
1. AzureSynapseOutputDataSource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)
1. AzureTableOutputDataSource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)
1. BlobOutputDataSource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)
1. DocumentDbOutputDataSource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)
1. EventHubOutputDataSource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)
1. EventHubV2OutputDataSource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)
1. OutputDataSource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)
1. PowerBIOutputDataSource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)
1. ServiceBusQueueOutputDataSource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)
1. ServiceBusTopicOutputDataSource.AsAzureFunctionOutputDataSource() (*AzureFunctionOutputDataSource, bool)

### Struct Changes

#### New Structs

1. AzureFunctionOutputDataSource
1. AzureFunctionOutputDataSourceProperties

#### New Struct Fields

1. BlobDataSourceProperties.AuthenticationMode
1. BlobReferenceInputDataSourceProperties.AuthenticationMode
1. BlobStreamInputDataSourceProperties.AuthenticationMode
