package virtualmachinedisk

import "github.com/MSOpenTech/azure-sdk-for-go/azure"

//DiskClient is used to manage operations on Azure Disks
type DiskClient struct {
	client *azure.Client
}
