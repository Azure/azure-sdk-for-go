package virtualmachinedisk

import "github.com/Azure/azure-sdk-for-go/management"

//DiskClient is used to manage operations on Azure Disks
type DiskClient struct {
	client management.Client
}

type HostCachingType string

const (
	// Enum values for HostCachingType
	HostCachingTypeNone      = HostCachingType("None")
	HostCachingTypeReadOnly  = HostCachingType("ReadOnly")
	HostCachingTypeReadWrite = HostCachingType("ReadWrite")
)
