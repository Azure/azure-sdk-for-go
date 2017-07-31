package locks

import (
	 original "github.com/Azure/azure-sdk-for-go/service/resources/management/2016-09-01/locks"
)

type (
	 ManagementLocksClient = original.ManagementLocksClient
	 LockLevel = original.LockLevel
	 ManagementLockListResult = original.ManagementLockListResult
	 ManagementLockObject = original.ManagementLockObject
	 ManagementLockOwner = original.ManagementLockOwner
	 ManagementLockProperties = original.ManagementLockProperties
	 ManagementClient = original.ManagementClient
)

const (
	 CanNotDelete = original.CanNotDelete
	 NotSpecified = original.NotSpecified
	 ReadOnly = original.ReadOnly
	 DefaultBaseURI = original.DefaultBaseURI
)
