package vmutils

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/management"
	"github.com/Azure/azure-sdk-for-go/management/hostedservice"
	"github.com/Azure/azure-sdk-for-go/management/virtualmachine"
)

func Example() {
	dnsName := "test-vm-from-go"
	storageAccount := "mystorageaccount"
	location := "West US"
	vmSize := "Small"
	vmImage := "b39f27a8b8c64d52b05eac6a62ebad85__Ubuntu-14_04-LTS-amd64-server-20140724-en-us-30GB"
	userName := "testuser"
	userPassword := "Test123"

	client, err := management.ClientFromPublishSettingsFile("path/to/downloaded.publishsettings", "")
	if err != nil {
		panic(err)
	}

	// create hosted service
	operationId, err := hostedservice.NewClient(client).CreateHostedService(dnsName, location, "", dnsName, "")
	if err != nil {
		panic(err)
	}
	err = client.WaitAsyncOperation(operationId)
	if err != nil {
		panic(err)
	}

	// create virtual machine
	role := NewVmConfiguration(dnsName, vmSize)
	ConfigureDeploymentFromPlatformImage(&role,
		vmImage, fmt.Sprintf("http://%s.blob.core.windows.net/sdktest/%s.vhd", storageAccount, dnsName), "")
	ConfigureForLinux(&role, dnsName, userName, userPassword)
	ConfigureWithPublicSSH(&role)

	operationId, err = virtualmachine.NewClient(client).CreateDeployment(role, dnsName)
	if err != nil {
		panic(err)
	}
	err = client.WaitAsyncOperation(operationId)
	if err != nil {
		panic(err)
	}
}
