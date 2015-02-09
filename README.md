# Azure SDK for Golang
This project provides a Golang package that makes it easy to consume and manage Microsoft Azure Services.

# Installation
- Install Golang: https://golang.org/doc/install
- Get Azure SDK package: 

```sh
go get github.com/MSOpenTech/azure-sdk-for-go
```
- Install: 

```sh
go install github.com/MSOpenTech/azure-sdk-for-go
```

# Usage

Create linux VM:

```C
package main

import (
    "fmt"
    "os"
    
    "github.com/MSOpenTech/azure-sdk-for-go/azure"
)

func main() {
    dnsName := "test-vm-from-go"
    location := "West US"
    vmSize := "Small"
    vmImage := "b39f27a8b8c64d52b05eac6a62ebad85__Ubuntu-14_04-LTS-amd64-server-20140724-en-us-30GB"
    userName := "testuser"
    userPassword := "Test123"
    sshCert := ""
    sshPort := 22

    client, err := azure.NewClient(&config{SubscriptionId: SUBSCRIPTION_ID, ManagementCertificate: SUBSCRIPTION_CERTIFICATE})
    if err != nil {
    	fmt.Println(err)
    	os.Exit(1)
    }

    vmClient := client.VM()
    
    vmConfig, err := vmClient.CreateAzureVMConfiguration(dnsName, vmSize, vmImage, location)
    if err != nil {
    	fmt.Println(err)
    	os.Exit(1)
    }
    
    vmConfig, err = vmClient.AddAzureLinuxProvisioningConfig(vmConfig, userName, userPassword, sshCert, sshPort)
    if err != nil {
    	fmt.Println(err)
    	os.Exit(1)
    }
    
    err = vmClient.CreateAzureVM(vmConfig, dnsName, location)
    if err != nil {
    	fmt.Println(err)
    	os.Exit(1)
    }
}
```

# License
[Apache 2.0](LICENSE-2.0.txt)
