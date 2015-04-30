package vmutils

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/management"
	"github.com/Azure/azure-sdk-for-go/management/hostedservice"
	"github.com/Azure/azure-sdk-for-go/management/location"
	"github.com/Azure/azure-sdk-for-go/management/osimage"
	storage "github.com/Azure/azure-sdk-for-go/management/storageservice"
	vm "github.com/Azure/azure-sdk-for-go/management/virtualmachine"
	vmimage "github.com/Azure/azure-sdk-for-go/management/virtualmachineimage"

	"github.com/Azure/azure-sdk-for-go/management/testutils"
)

func TestDeployPlatformImage(t *testing.T) {
	client := testutils.GetTestClient(t)
	vmname := GenerateName()
	sa := GetTestStorageAccount(t, client)
	location := sa.StorageServiceProperties.Location

	role := NewVmConfiguration(vmname, "Standard_D3")
	ConfigureDeploymentFromPlatformImage(&role,
		GetLinuxTestImage(t, client).Name,
		fmt.Sprintf("http://%s.blob.core.windows.net/sdktest/%s.vhd", sa.ServiceName, vmname),
		GenerateName())
	ConfigureForLinux(&role, "myvm", "azureuser", GeneratePassword())
	ConfigureWithPublicSSH(&role)

	testRoleConfiguration(t, client, role, location)
}

func TestVMImageList(t *testing.T) {
	client := testutils.GetTestClient(t)
	vmic := vmimage.NewClient(client)
	il, _ := vmic.GetImageList()
	for _, im := range il {
		t.Logf("%s -%s", im.Name, im.Description)
	}
}

func TestDeployPlatformCaptureRedeploy(t *testing.T) {
	client := testutils.GetTestClient(t)
	vmname := GenerateName()
	sa := GetTestStorageAccount(t, client)
	location := sa.StorageServiceProperties.Location

	role := NewVmConfiguration(vmname, "Standard_D3")
	ConfigureDeploymentFromPlatformImage(&role,
		GetLinuxTestImage(t, client).Name,
		fmt.Sprintf("http://%s.blob.core.windows.net/sdktest/%s.vhd", sa.ServiceName, vmname),
		GenerateName())
	ConfigureForLinux(&role, "myvm", "azureuser", GeneratePassword())
	ConfigureWithPublicSSH(&role)

	t.Logf("Deploying VM: %s", vmname)
	createRoleConfiguration(t, client, role, location)

	t.Logf("Wait for deployment to enter running state")
	vmc := vm.NewClient(client)
	status := vm.DeploymentStatusDeploying
	for status != vm.DeploymentStatusRunning {
		deployment, err := vmc.GetDeployment(vmname, vmname)
		if err != nil {
			t.Error(err)
			break
		}
		status = deployment.Status
	}

	t.Logf("Shutting down VM: %s", vmname)
	if err := Await(client, func() (management.OperationId, error) {
		return vmc.ShutdownRole(vmname, vmname, vmname)
	}); err != nil {
		t.Error(err)
	}

	if err := WaitForDeploymentPowerState(client, vmname, vmname, vm.PowerStateStopped); err != nil {
		t.Fatal(err)
	}

	imagename := GenerateName()
	t.Logf("Capturing VMImage: %s", imagename)
	if err := Await(client, func() (management.OperationId, error) {
		return vmc.CaptureRole(vmname, vmname, vmname, imagename, imagename, nil)
	}); err != nil {
		t.Error(err)
	}

	im := GetUserImage(t, client, imagename)
	t.Logf("Found image: %+v", im)

	newvmname := GenerateName()
	role = NewVmConfiguration(newvmname, "Standard_D3")
	ConfigureDeploymentFromPlatformImage(&role,
		im.Name,
		fmt.Sprintf("http://%s.blob.core.windows.net/sdktest/%s.vhd", sa.ServiceName, newvmname),
		GenerateName())
	ConfigureForLinux(&role, newvmname, "azureuser", GeneratePassword())
	ConfigureWithPublicSSH(&role)

	t.Logf("Deploying new VM from freshly captured VM image: %s", newvmname)
	if err := Await(client, func() (management.OperationId, error) {
		return vmc.CreateDeployment(role, vmname)
	}); err != nil {
		t.Error(err)
	}

	deleteHostedService(t, client, vmname)
}

func TestDeployFromVmImage(t *testing.T) {
	client := testutils.GetTestClient(t)
	vmname := GenerateName()
	sa := GetTestStorageAccount(t, client)
	location := sa.StorageServiceProperties.Location

	im := GetVMImage(t, client, func(im vmimage.VMImage) bool {
		return im.Name == "fb83b3509582419d99629ce476bcb5c8__SQL-Server-2014-RTM-12.0.2430.0-OLTP-ENU-Win2012R2-cy14su11"
	})

	role := NewVmConfiguration(vmname, "Standard_D4")
	ConfigureDeploymentFromVMImage(&role, im.Name,
		fmt.Sprintf("http://%s.blob.core.windows.net/%s", sa.ServiceName, vmname))
	ConfigureForWindows(&role, vmname, "azureuser", GeneratePassword(), true, "")
	ConfigureWithPublicSSH(&role)

	testRoleConfiguration(t, client, role, location)
}

func TestRoleStateOperations(t *testing.T) {
	client := testutils.GetTestClient(t)
	vmname := GenerateName()
	sa := GetTestStorageAccount(t, client)
	location := sa.StorageServiceProperties.Location

	role := NewVmConfiguration(vmname, "Standard_D3")
	ConfigureDeploymentFromPlatformImage(&role,
		GetLinuxTestImage(t, client).Name,
		fmt.Sprintf("http://%s.blob.core.windows.net/sdktest/%s.vhd", sa.ServiceName, vmname),
		GenerateName())
	ConfigureForLinux(&role, "myvm", "azureuser", GeneratePassword())

	createRoleConfiguration(t, client, role, location)

	vmc := vm.NewClient(client)
	if err := Await(client, func() (management.OperationId, error) {
		return vmc.ShutdownRole(vmname, vmname, vmname)
	}); err != nil {
		t.Error(err)
	}
	if err := Await(client, func() (management.OperationId, error) {
		return vmc.StartRole(vmname, vmname, vmname)
	}); err != nil {
		t.Error(err)
	}
	if err := Await(client, func() (management.OperationId, error) {
		return vmc.RestartRole(vmname, vmname, vmname)
	}); err != nil {
		t.Error(err)
	}

	deleteHostedService(t, client, vmname)
}

func testRoleConfiguration(t *testing.T, client management.Client, role vm.Role, location string) {
	createRoleConfiguration(t, client, role, location)

	deleteHostedService(t, client, role.RoleName)
}

func createRoleConfiguration(t *testing.T, client management.Client, role vm.Role, location string) {
	vmc := vm.NewClient(client)
	hsc := hostedservice.NewClient(client)
	vmname := role.RoleName

	if err := Await(client, func() (management.OperationId, error) {
		return hsc.CreateHostedService(vmname, location, "", vmname, "")
	}); err != nil {
		t.Error(err)
	}

	if err := Await(client, func() (management.OperationId, error) {
		return vmc.CreateDeployment(role, vmname)
	}); err != nil {
		t.Error(err)
	}
}

func deleteHostedService(t *testing.T, client management.Client, vmname string) {
	t.Logf("Deleting hosted service: %s", vmname)
	if err := Await(client, func() (management.OperationId, error) {
		return hostedservice.NewClient(client).DeleteHostedService(vmname, true)
	}); err != nil {
		t.Error(err)
	}
}

// === utility funcs ===

func GetTestStorageAccount(t *testing.T, client management.Client) *storage.StorageService {
	t.Log("Retrieving storage account")
	sc := storage.NewClient(client)
	var sa *storage.StorageService
	ssl, err := sc.GetStorageServiceList()
	if err != nil {
		t.Fatal(err)
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	if len(ssl.StorageServices) == 0 {
		t.Log("No storage accounts found, creating a new one")
		lc := location.NewClient(client)
		ll, err := lc.GetLocationList()
		if err != nil {
			t.Fatal(err)
		}
		loc := ll.Locations[rnd.Intn(len(ll.Locations))].Name

		t.Logf("Location for new storage account: %s", loc)
		name := GenerateName()
		sa, err = sc.Create(storage.StorageAccountCreateParameters{
			ServiceName: name,
			Label:       base64.StdEncoding.EncodeToString([]byte(name)),
			Location:    loc,
			AccountType: storage.AccountTypeStandardLRS})
		if err != nil {
			t.Fatal(err)
		}
	} else {

		sa = &ssl.StorageServices[rnd.Intn(len(ssl.StorageServices))]
	}

	t.Logf("Selected storage account '%s' in location '%s'",
		sa.ServiceName, sa.StorageServiceProperties.Location)

	return sa
}

func GetLinuxTestImage(t *testing.T, client management.Client) osimage.OSImage {
	return GetOSImage(t, client, func(im osimage.OSImage) bool {
		return im.Category == "Public" && im.ImageFamily == "Ubuntu Server 14.04 LTS"
	})
}

func GetUserImage(t *testing.T, client management.Client, name string) osimage.OSImage {
	return GetOSImage(t, client, func(im osimage.OSImage) bool {
		return im.Category == "User" && im.Name == name
	})
}

func GetOSImage(t *testing.T, client management.Client, filter func(osimage.OSImage) bool) osimage.OSImage {
	t.Log("Selecting OS image")
	osc := osimage.NewClient(client)
	allimages, err := osc.GetImageList()
	if err != nil {
		t.Fatal(err)
	}
	filtered := []osimage.OSImage{}
	for _, im := range allimages {
		if filter(im) {
			filtered = append(filtered, im)
		}
	}
	if len(filtered) == 0 {
		t.Fatal("Filter too restrictive, no images left?")
	}

	image := filtered[0]
	for _, im := range filtered {
		if im.PublishedDate > image.PublishedDate {
			image = im
		}
	}

	t.Logf("Selecting image '%s'", image.Name)
	return image
}

func GetVMImage(t *testing.T, client management.Client, filter func(vmimage.VMImage) bool) vmimage.VMImage {
	t.Log("Selecting VM image")
	allimages, err := vmimage.NewClient(client).GetImageList()
	if err != nil {
		t.Fatal(err)
	}
	filtered := []vmimage.VMImage{}
	for _, im := range allimages {
		if filter(im) {
			filtered = append(filtered, im)
		}
	}
	if len(filtered) == 0 {
		t.Fatal("Filter too restrictive, no images left?")
	}

	image := filtered[0]
	for _, im := range filtered {
		if im.PublishedDate > image.PublishedDate {
			image = im
		}
	}

	t.Logf("Selecting image '%s'", image.Name)
	return image
}

func GenerateName() string {
	from := "1234567890abcdefghijklmnopqrstuvwxyz"
	return "sdk" + GenerateString(12, from)
}

func GeneratePassword() string {
	pw := GenerateString(20, "1234567890") +
		GenerateString(20, "abcdefghijklmnopqrstuvwxyz") +
		GenerateString(20, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := rnd.Intn(len(pw)-2) + 1

	pw = string(append([]uint8(pw[i:]), pw[:i-1]...))

	return pw
}

func GenerateString(length int, from string) string {
	str := ""
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(str) < length {
		str += string(from[rnd.Intn(len(from))])
	}
	return str
}

type asyncFunc func() (operationId management.OperationId, err error)

func Await(client management.Client, async asyncFunc) error {
	reqId, err := async()
	if err != nil {
		return err
	}
	return client.WaitAsyncOperation(reqId)
}
