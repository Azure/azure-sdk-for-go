# Azure SDK for Go
This project provides various Golang packages that makes it easy to consume and manage
Microsoft Azure Services.

# Installation
- Install Golang: https://golang.org/doc/install
- Get Azure SDK package: 

    go get -d github.com/Azure/azure-sdk-for-go

# Usage

Read Godoc of the package at: http://godoc.org/github.com/Azure/azure-sdk-for-go/management

Download publish settings file from https://manage.windowsazure.com/publishsettings.

# Example: Create a Linux VM

```go
package main

import (
	"encoding/base64"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/management"
	"github.com/Azure/azure-sdk-for-go/management/hostedservice"
	"github.com/Azure/azure-sdk-for-go/management/virtualmachine"
	"github.com/Azure/azure-sdk-for-go/management/vmutils"
)

func main() {
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
	if err := hostedservice.NewClient(client).CreateHostedService(hostedservice.CreateHostedServiceParameters{
		ServiceName: dnsName,
		Location:    location,
		Label:       base64.StdEncoding.EncodeToString([]byte(dnsName))}); err != nil {
		panic(err)
	}

	// create virtual machine
	role := vmutils.NewVMConfiguration(dnsName, vmSize)
	vmutils.ConfigureDeploymentFromPlatformImage(
		&role,
		vmImage,
		fmt.Sprintf("http://%s.blob.core.windows.net/sdktest/%s.vhd", storageAccount, dnsName),
		"")
	vmutils.ConfigureForLinux(&role, dnsName, userName, userPassword)
	vmutils.ConfigureWithPublicSSH(&role)

	operationID, err := virtualmachine.NewClient(client).
		CreateDeployment(role, dnsName, virtualmachine.CreateDeploymentOptions{})
	if err != nil {
		panic(err)
	}
	if err := client.WaitAsyncOperation(operationID); err != nil {
		panic(err)
	}
}
```

# License

This project is published under [Apache 2.0 License](LICENSE).
