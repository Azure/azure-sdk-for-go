package cosmosdb

import (
	 original "github.com/Azure/azure-sdk-for-go/service/cosmos-db/management/2015-04-08/cosmos-db"
)

type (
	 ManagementClient = original.ManagementClient
	 DatabaseAccountsClient = original.DatabaseAccountsClient
	 DatabaseAccountKind = original.DatabaseAccountKind
	 DatabaseAccountOfferType = original.DatabaseAccountOfferType
	 DefaultConsistencyLevel = original.DefaultConsistencyLevel
	 KeyKind = original.KeyKind
	 ConsistencyPolicy = original.ConsistencyPolicy
	 DatabaseAccount = original.DatabaseAccount
	 DatabaseAccountConnectionString = original.DatabaseAccountConnectionString
	 DatabaseAccountCreateUpdateParameters = original.DatabaseAccountCreateUpdateParameters
	 DatabaseAccountCreateUpdateProperties = original.DatabaseAccountCreateUpdateProperties
	 DatabaseAccountListConnectionStringsResult = original.DatabaseAccountListConnectionStringsResult
	 DatabaseAccountListKeysResult = original.DatabaseAccountListKeysResult
	 DatabaseAccountListReadOnlyKeysResult = original.DatabaseAccountListReadOnlyKeysResult
	 DatabaseAccountPatchParameters = original.DatabaseAccountPatchParameters
	 DatabaseAccountProperties = original.DatabaseAccountProperties
	 DatabaseAccountRegenerateKeyParameters = original.DatabaseAccountRegenerateKeyParameters
	 DatabaseAccountsListResult = original.DatabaseAccountsListResult
	 FailoverPolicies = original.FailoverPolicies
	 FailoverPolicy = original.FailoverPolicy
	 Location = original.Location
	 Resource = original.Resource
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 GlobalDocumentDB = original.GlobalDocumentDB
	 MongoDB = original.MongoDB
	 Parse = original.Parse
	 Standard = original.Standard
	 BoundedStaleness = original.BoundedStaleness
	 ConsistentPrefix = original.ConsistentPrefix
	 Eventual = original.Eventual
	 Session = original.Session
	 Strong = original.Strong
	 Primary = original.Primary
	 PrimaryReadonly = original.PrimaryReadonly
	 Secondary = original.Secondary
	 SecondaryReadonly = original.SecondaryReadonly
)
