package vmutils

import (
	"fmt"

	"github.com/MSOpenTech/azure-sdk-for-go/management"
	"github.com/MSOpenTech/azure-sdk-for-go/management/hostedservice"
	"github.com/MSOpenTech/azure-sdk-for-go/management/virtualmachine"
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
	requestId, err := hostedservice.NewClient(client).CreateHostedService(dnsName, location, "", dnsName, "")
	if err != nil {
		panic(err)
	}
	err = client.WaitAsyncOperation(requestId)
	if err != nil {
		panic(err)
	}

	// create virtual machine
	role, err := NewVmConfiguration(dnsName, vmSize)
	if err != nil {
		panic(err)
	}
	ConfigureDeploymentFromPlatformImage(&role,
		vmImage, fmt.Sprintf("http://%s.blob.core.windows.net/sdktest/%s.vhd", storageAccount, dnsName), "")
	ConfigureForLinux(&role, dnsName, userName, userPassword)
	ConfigureWithPublicSSH(&role)

	requestId, err = virtualmachine.NewClient(client).CreateDeployment(role, dnsName)
	if err != nil {
		panic(err)
	}
	err = client.WaitAsyncOperation(requestId)
	if err != nil {
		panic(err)
	}
}
