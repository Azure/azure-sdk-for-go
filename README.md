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
	"io/ioutil"
    "os"

    "github.com/MSOpenTech/azure-sdk-for-go/management"
    "github.com/MSOpenTech/azure-sdk-for-go/management/virtualmachine"
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

	subscriptionCert, err := ioutil.ReadFile(SUBSCRIPTION_CERTIFICATE_PATH)
	if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    client, err := management.NewClient(SUBSCRIPTION_ID, subscriptionCert)
    if err != nil {
    	fmt.Println(err)
    	os.Exit(1)
    }

    vmClient := virtualmachine.NewClient(client)

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

# Setting up your Development Environment

In order to contribute to this project, it is important that you have properly set up your Go workspace. To do so, ensure that your `GOPATH` environment variable has been set.

Additional instructions for setting up your Go workspace can be found here: http://golang.org/doc/code.html.

## Forking and Submitting Pull Requests

In order to avoid any issues with Go dependencies and import paths, it is recommended that the following steps be used to properly fork and contribute to this project:

1. Clone the project to your `GOPATH`

  ```sh
  go get github.com/MSOpenTech/azure-sdk-for-go
  ```

2. Fork the project on GitHub
3. Add your fork as a remote

  ```sh
  git remote add [user] https://github.com/[user]/azure-sdk-for-go.git
  ```

4. Create a new feature/bug fix branch

  ```sh
  git checkout -b [branch]
  ```

5. Develop and commit your changes
6. If the project's master branch has been updated significantly throughout the course of development of your feature/bug fix, rebase the master branch's changes on to your feature/bug fix branch and resolve any conflicts.

  ```sh
  git checkout [branch]
  git fetch origin
  git rebase origin/master # Resolve any merge conflicts
  git add [files]
  git rebase --continue
  ```

# License
[Apache 2.0](LICENSE-2.0.txt)
