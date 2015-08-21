package main

import (
	"fmt"
	"os"

	"github.com/azure/azure-sdk-for-go/arm/compute"
	"github.com/azure/azure-sdk-for-go/arm/examples/helpers"
	"github.com/azure/go-autorest/autorest/azure"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please provide a resource group and name to use")
		os.Exit(1)
	}
	resourceGroup := os.Args[1]
	name := os.Args[2]

	c, err := helpers.LoadCredentials()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	vmc := compute.NewVirtualMachinesClient(c["subscriptionId"])

	spt, err := helpers.NewServicePrincipalTokenFromCredentials(c, azure.AzureResourceManagerScope)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	vmc.Authorizer = spt

	vm := compute.VirtualMachine{}
	vm, err := vmc.CreateOrUpdate(resourceGroup, name, parameters)
}
