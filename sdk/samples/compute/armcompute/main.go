package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/compute/armcompute"
	"github.com/Azure/azure-sdk-for-go/sdk/network/armnetwork"
	"github.com/Azure/azure-sdk-for-go/sdk/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/to"
)

var (
	subscriptionId    string
	location          = "westus2"
	resourceGroupName = "sample-resourcegroup"
	vnetName          = "sample-vnet"
	subnetName        = "internal"
	nicName           = "sample-nic"
	vmName            = "sample-vm"
)

func init() {
	subscriptionId = os.Getenv("AZURE_SUBSCRIPTION_ID")
}

func main() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		panic(err)
	}
	conn := armcore.NewDefaultConnection(cred, &armcore.ConnectionOptions{
		Logging: azcore.LogOptions{
			IncludeBody: true,
		},
	})

	ctx := context.Background()
	resourceGroup, err := createResourceGroup(ctx, conn)
	if err != nil {
		log.Fatalf("Cannot create resource group: %+v", err)
	}
	log.Printf("Resource Group ID: %s", *resourceGroup.ID)

	vnet, err := createVirtualNetwork(ctx, conn)
	if err != nil {
		log.Fatalf("Cannot create virtual network: %+v", err)
	}
	log.Printf("Virtual Network ID: %s", *vnet.ID)

	subnet, err := createSubnet(ctx, conn)
	if err != nil {
		log.Fatalf("Cannot create subnet: %+v", err)
	}
	log.Printf("Subnet ID: %s", *subnet.ID)

	nic, err := createNIC(ctx, conn, *subnet.ID)
	if err != nil {
		log.Fatalf("Cannot create network interface: %+v", err)
	}
	log.Printf("Network Interface ID: %s", *nic.ID)

	vm, err := createVirtualMachine(ctx, conn, *nic.ID)
	if err != nil {
		log.Fatalf("Cannot create virtual machine: %+v", err)
	}
	log.Printf("Virtual Machine ID: %s", *vm.ID)
}

func createResourceGroup(ctx context.Context, connection *armcore.Connection) (*armresources.ResourceGroup, error) {
	rgClient := armresources.NewResourceGroupsClient(connection, subscriptionId)

	param := armresources.ResourceGroup{
		Location: to.StringPtr(location),
	}

	resp, err := rgClient.CreateOrUpdate(ctx, resourceGroupName, param, nil)
	if err != nil {
		return nil, err
	}

	return resp.ResourceGroup, nil
}

func updateResourceGroup(ctx context.Context, connection *armcore.Connection) (*armresources.ResourceGroup, error) {
	rgClient := armresources.NewResourceGroupsClient(connection, subscriptionId)

	update := armresources.ResourceGroupPatchable{
		Tags: map[string]*string{
			"new": to.StringPtr("tag"),
		},
	}
	resp, err := rgClient.Update(ctx, resourceGroupName, update, nil)
	if err != nil {
		return nil, err
	}

	return resp.ResourceGroup, nil
}

func listResourceGroups(ctx context.Context, connection *armcore.Connection) ([]*armresources.ResourceGroup, error) {
	rgClient := armresources.NewResourceGroupsClient(connection, subscriptionId)

	pager := rgClient.List(nil)

	var resourceGroups []*armresources.ResourceGroup
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		if resp.ResourceGroupListResult != nil {
			resourceGroups = append(resourceGroups, resp.ResourceGroupListResult.Value...)
		}
	}
	return resourceGroups, pager.Err()
}

func deleteResourceGroup(ctx context.Context, connection *armcore.Connection) error {
	rgClient := armresources.NewResourceGroupsClient(connection, subscriptionId)

	poller, err := rgClient.BeginDelete(ctx, resourceGroupName, nil)
	if err != nil {
		return err
	}
	if _, err := poller.PollUntilDone(ctx, 10*time.Second); err != nil {
		return err
	}

	return nil
}

func createVirtualNetwork(ctx context.Context, connection *armcore.Connection) (*armnetwork.VirtualNetwork, error) {
	vnetClient := armnetwork.NewVirtualNetworksClient(connection, subscriptionId)

	param := armnetwork.VirtualNetwork{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(location),
		},
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.StringPtr("10.0.0.0/16"),
				},
			},
		},
	}
	poller, err := vnetClient.BeginCreateOrUpdate(ctx, resourceGroupName, vnetName, param, nil)
	if err != nil {
		return nil, err
	}
	resp, err := poller.PollUntilDone(ctx, 10*time.Second)
	if err != nil {
		return nil, err
	}

	return resp.VirtualNetwork, nil
}

func deleteVirtualNetwork(ctx context.Context, connection *armcore.Connection) error {
	vnetClient := armnetwork.NewVirtualNetworksClient(connection, subscriptionId)

	poller, err := vnetClient.BeginDelete(ctx, resourceGroupName, vnetName, nil)
	if err != nil {
		return err
	}
	if _, err := poller.PollUntilDone(ctx, 10*time.Second); err != nil {
		return err
	}

	return nil
}

func listVirtualNetwork(ctx context.Context, connection *armcore.Connection) ([]*armnetwork.VirtualNetwork, error) {
	vnetClient := armnetwork.NewVirtualNetworksClient(connection, subscriptionId)

	pager := vnetClient.List(resourceGroupName, nil)

	var virtualNetworks []*armnetwork.VirtualNetwork
	for pager.NextPage(ctx) {
		resp := pager.PageResponse()
		virtualNetworks = append(virtualNetworks, resp.VirtualNetworkListResult.Value...)
	}

	return virtualNetworks, pager.Err()
}

func createSubnet(ctx context.Context, connection *armcore.Connection) (*armnetwork.Subnet, error) {
	subnetClient := armnetwork.NewSubnetsClient(connection, subscriptionId)

	param := armnetwork.Subnet{
		Properties: &armnetwork.SubnetPropertiesFormat{
			AddressPrefix: to.StringPtr("10.0.2.0/24"),
		},
	}
	poller, err := subnetClient.BeginCreateOrUpdate(ctx, resourceGroupName, vnetName, subnetName, param, nil)
	if err != nil {
		return nil, err
	}
	resp, err := poller.PollUntilDone(ctx, 10*time.Second)
	if err != nil {
		return nil, err
	}
	return resp.Subnet, nil
}

func deleteSubnet(ctx context.Context, connection *armcore.Connection) error {
	subnetClient := armnetwork.NewSubnetsClient(connection, subscriptionId)

	poller, err := subnetClient.BeginDelete(ctx, resourceGroupName, vnetName, subnetName, nil)
	if err != nil {
		return err
	}
	if _, err := poller.PollUntilDone(ctx, 10*time.Second); err != nil {
		return err
	}

	return nil
}

func createNIC(ctx context.Context, connection *armcore.Connection, subnetID string) (*armnetwork.NetworkInterface, error) {
	nicClient := armnetwork.NewNetworkInterfacesClient(connection, subscriptionId)

	param := armnetwork.NetworkInterface{
		Resource: armnetwork.Resource{
			Location: to.StringPtr(location),
		},
		Properties: &armnetwork.NetworkInterfacePropertiesFormat{
			IPConfigurations: []*armnetwork.NetworkInterfaceIPConfiguration{
				{
					Name: to.StringPtr("internal"),
					Properties: &armnetwork.NetworkInterfaceIPConfigurationPropertiesFormat{
						PrivateIPAllocationMethod: armnetwork.IPAllocationMethodDynamic.ToPtr(),
						Subnet: &armnetwork.Subnet{
							SubResource: armnetwork.SubResource{
								ID: to.StringPtr(subnetID),
							},
						},
					},
				},
			},
		},
	}
	poller, err := nicClient.BeginCreateOrUpdate(ctx, resourceGroupName, nicName, param, nil)
	if err != nil {
		return nil, err
	}
	resp, err := poller.PollUntilDone(ctx, 10*time.Second)
	if err != nil {
		return nil, err
	}

	return resp.NetworkInterface, nil
}

func deleteNIC(ctx context.Context, connection *armcore.Connection) error {
	nicClient := armnetwork.NewNetworkInterfacesClient(connection, subscriptionId)

	poller, err := nicClient.BeginDelete(ctx, resourceGroupName, nicName, nil)
	if err != nil {
		return err
	}
	if _, err := poller.PollUntilDone(ctx, 10*time.Second); err != nil {
		return err
	}

	return nil
}

func createVirtualMachine(ctx context.Context, connection *armcore.Connection, nicID string) (*armcompute.VirtualMachine, error) {
	vmClient := armcompute.NewVirtualMachinesClient(connection, subscriptionId)

	param := armcompute.VirtualMachine{
		Resource: armcompute.Resource{
			Location: to.StringPtr(location),
		},
		Identity: &armcompute.VirtualMachineIdentity{
			Type: armcompute.ResourceIdentityTypeSystemAssigned.ToPtr(),
		},
		Properties: &armcompute.VirtualMachineProperties{
			HardwareProfile: &armcompute.HardwareProfile{
				VMSize: armcompute.VirtualMachineSizeTypesStandardF2.ToPtr(),
			},
			OSProfile: &armcompute.OSProfile{
				AdminPassword:        to.StringPtr("P@$$w0rd1234!"),
				AdminUsername:        to.StringPtr("adminuser"),
				ComputerName:         to.StringPtr("arcturus"),
				WindowsConfiguration: &armcompute.WindowsConfiguration{},
			},
			NetworkProfile: &armcompute.NetworkProfile{
				NetworkInterfaces: []*armcompute.NetworkInterfaceReference{
					{
						SubResource: armcompute.SubResource{
							ID: to.StringPtr(nicID),
						},
					},
				},
			},
			StorageProfile: &armcompute.StorageProfile{
				ImageReference: &armcompute.ImageReference{
					Offer:     to.StringPtr("WindowsServer"),
					Publisher: to.StringPtr("MicrosoftWindowsServer"),
					SKU:       to.StringPtr("2016-Datacenter"),
					Version:   to.StringPtr("latest"),
				},
				OSDisk: &armcompute.OSDisk{
					CreateOption: armcompute.DiskCreateOptionTypesFromImage.ToPtr(),
					Caching:      armcompute.CachingTypesReadWrite.ToPtr(),
					ManagedDisk: &armcompute.ManagedDiskParameters{
						StorageAccountType: armcompute.StorageAccountTypesStandardLRS.ToPtr(),
					},
					OSType: armcompute.OperatingSystemTypesWindows.ToPtr(),
				},
			},
		},
	}

	poller, err := vmClient.BeginCreateOrUpdate(ctx, resourceGroupName, vmName, param, nil)
	if err != nil {
		return nil, err
	}

	resp, err := poller.PollUntilDone(ctx, 10*time.Second)
	if err != nil {
		return nil, err
	}

	return resp.VirtualMachine, nil
}

func deleteVirtualMachine(ctx context.Context, connection *armcore.Connection) error {
	vmClient := armcompute.NewVirtualMachinesClient(connection, subscriptionId)

	poller, err := vmClient.BeginDelete(ctx, resourceGroupName, vmName, &armcompute.VirtualMachinesBeginDeleteOptions{
		ForceDeletion: to.BoolPtr(true),
	})
	if err != nil {
		return err
	}
	if _, err := poller.PollUntilDone(ctx, 10*time.Second); err != nil {
		return err
	}

	return nil
}
