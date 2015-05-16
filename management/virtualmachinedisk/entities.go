package virtualmachinedisk

import "github.com/Azure/azure-sdk-for-go/management"

// DiskClient is used to perform operations on Azure Disks
type DiskClient struct {
	client management.Client
}

type HostCachingType string

// Enum values for HostCachingType
const (
	HostCachingTypeNone      = HostCachingType("None")
	HostCachingTypeReadOnly  = HostCachingType("ReadOnly")
	HostCachingTypeReadWrite = HostCachingType("ReadWrite")
)
