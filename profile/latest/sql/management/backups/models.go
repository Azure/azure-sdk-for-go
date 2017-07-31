package backups

import (
	 original "github.com/Azure/azure-sdk-for-go/service/sql/management/2014-04-01/backups"
)

type (
	 RestorableDroppedDatabasesClient = original.RestorableDroppedDatabasesClient
	 RestorePointsClient = original.RestorePointsClient
	 ManagementClient = original.ManagementClient
	 RestorePointType = original.RestorePointType
	 ProxyResource = original.ProxyResource
	 RecoverableDatabase = original.RecoverableDatabase
	 RecoverableDatabaseListResult = original.RecoverableDatabaseListResult
	 RecoverableDatabaseProperties = original.RecoverableDatabaseProperties
	 Resource = original.Resource
	 RestorableDroppedDatabase = original.RestorableDroppedDatabase
	 RestorableDroppedDatabaseListResult = original.RestorableDroppedDatabaseListResult
	 RestorableDroppedDatabaseProperties = original.RestorableDroppedDatabaseProperties
	 RestorePoint = original.RestorePoint
	 RestorePointListResult = original.RestorePointListResult
	 RestorePointProperties = original.RestorePointProperties
	 TrackedResource = original.TrackedResource
	 RecoverableDatabasesClient = original.RecoverableDatabasesClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 CONTINUOUS = original.CONTINUOUS
	 DISCRETE = original.DISCRETE
)
