package vmutils

import (
	"encoding/xml"
	"testing"

	vmdisk "github.com/Azure/azure-sdk-for-go/management/virtualmachinedisk"
)

func TestNewLinuxVmRemoteImage(t *testing.T) {
	role := NewVmConfiguration("myvm", "Standard_D3")
	ConfigureDeploymentFromRemoteImage(&role,
		"http://remote.host/some.vhd?sv=12&sig=ukhfiuwef78687", "Linux",
		"myvm-os-disk", "http://mystorageacct.blob.core.windows.net/vhds/mybrandnewvm.vhd",
		"OSDisk")
	ConfigureForLinux(&role, "myvm", "azureuser", "P@ssword", "2398yyKJGd78e2389ydfncuirowebhf89yh3IUOBY")
	ConfigureWithPublicSSH(&role)

	bytes, err := xml.MarshalIndent(role, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	expected := `<Role>
  <RoleName>myvm</RoleName>
  <RoleType>PersistentVMRole</RoleType>
  <ConfigurationSets>
    <ConfigurationSet>
      <ConfigurationSetType>LinuxProvisioningConfiguration</ConfigurationSetType>
      <HostName>myvm</HostName>
      <UserName>azureuser</UserName>
      <UserPassword>P@ssword</UserPassword>
      <DisableSshPasswordAuthentication>false</DisableSshPasswordAuthentication>
      <SSH>
        <PublicKeys>
          <PublicKey>
            <Fingerprint>2398yyKJGd78e2389ydfncuirowebhf89yh3IUOBY</Fingerprint>
            <Path>/home/azureuser/.ssh/authorized_keys</Path>
          </PublicKey>
        </PublicKeys>
      </SSH>
    </ConfigurationSet>
    <ConfigurationSet>
      <ConfigurationSetType>NetworkConfiguration</ConfigurationSetType>
      <InputEndpoints>
        <InputEndpoint>
          <LocalPort>22</LocalPort>
          <Name>SSH</Name>
          <Port>22</Port>
          <Protocol>TCP</Protocol>
        </InputEndpoint>
      </InputEndpoints>
    </ConfigurationSet>
  </ConfigurationSets>
  <OSVirtualHardDisk>
    <DiskName>myvm-os-disk</DiskName>
    <DiskLabel>OSDisk</DiskLabel>
    <MediaLink>http://mystorageacct.blob.core.windows.net/vhds/mybrandnewvm.vhd</MediaLink>
    <OS>Linux</OS>
    <RemoteSourceImageLink>http://remote.host/some.vhd?sv=12&amp;sig=ukhfiuwef78687</RemoteSourceImageLink>
  </OSVirtualHardDisk>
  <RoleSize>Standard_D3</RoleSize>
  <ProvisionGuestAgent>true</ProvisionGuestAgent>
</Role>`

	if string(bytes) != expected {
		t.Fatalf("Expected marshalled xml to be %q, but got %q", expected, string(bytes))
	}
}

func TestNewLinuxVmPlatformImage(t *testing.T) {
	role := NewVmConfiguration("myplatformvm", "Standard_D3")
	ConfigureDeploymentFromPlatformImage(&role,
		"b39f27a8b8c64d52b05eac6a62ebad85__Ubuntu-14_04_2_LTS-amd64-server-20150309-en-us-30GB",
		"http://mystorageacct.blob.core.windows.net/vhds/mybrandnewvm.vhd", "mydisklabel")
	ConfigureForLinux(&role, "myvm", "azureuser", "", "2398yyKJGd78e2389ydfncuirdebhf89yh3IUOBY")

	bytes, err := xml.MarshalIndent(role, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	expected := `<Role>
  <RoleName>myplatformvm</RoleName>
  <RoleType>PersistentVMRole</RoleType>
  <ConfigurationSets>
    <ConfigurationSet>
      <ConfigurationSetType>LinuxProvisioningConfiguration</ConfigurationSetType>
      <HostName>myvm</HostName>
      <UserName>azureuser</UserName>
      <SSH>
        <PublicKeys>
          <PublicKey>
            <Fingerprint>2398yyKJGd78e2389ydfncuirdebhf89yh3IUOBY</Fingerprint>
            <Path>/home/azureuser/.ssh/authorized_keys</Path>
          </PublicKey>
        </PublicKeys>
      </SSH>
    </ConfigurationSet>
  </ConfigurationSets>
  <OSVirtualHardDisk>
    <MediaLink>http://mystorageacct.blob.core.windows.net/vhds/mybrandnewvm.vhd</MediaLink>
    <SourceImageName>b39f27a8b8c64d52b05eac6a62ebad85__Ubuntu-14_04_2_LTS-amd64-server-20150309-en-us-30GB</SourceImageName>
  </OSVirtualHardDisk>
  <RoleSize>Standard_D3</RoleSize>
  <ProvisionGuestAgent>true</ProvisionGuestAgent>
</Role>`

	if string(bytes) != expected {
		t.Fatalf("Expected marshalled xml to be %q, but got %q", expected, string(bytes))
	}
}

func TestNewVmFromVMImage(t *testing.T) {
	role := NewVmConfiguration("restoredbackup", "Standard_D1")
	ConfigureDeploymentFromVMImage(&role, "myvm-backup-20150209",
		"http://mystorageacct.blob.core.windows.net/vhds/myoldnewvm.vhd")

	bytes, err := xml.MarshalIndent(role, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	expected := `<Role>
  <RoleName>restoredbackup</RoleName>
  <RoleType>PersistentVMRole</RoleType>
  <VMImageName>myvm-backup-20150209</VMImageName>
  <MediaLocation>http://mystorageacct.blob.core.windows.net/vhds/myoldnewvm.vhd</MediaLocation>
  <RoleSize>Standard_D1</RoleSize>
</Role>`

	if string(bytes) != expected {
		t.Fatalf("Expected marshalled xml to be %q, but got %q", expected, string(bytes))
	}
}

func TestNewVmFromExistingDisk(t *testing.T) {
	role := NewVmConfiguration("blobvm", "Standard_D14")
	ConfigureDeploymentFromExistingOSDisk(&role, "myvm-backup-20150209", "OSDisk")
	ConfigureForWindows(&role, "WINVM", "azuser", "P2ssw@rd", true, "")
	ConfigureWindowsToJoinDomain(&role, "user@domain.com", "youReN3verG0nnaGu3ss", "redmond.corp.contoso.com", "")
	ConfigureWithNewDataDisk(&role, "my-brand-new-disk", "http://account.blob.core.windows.net/vhds/newdatadisk.vhd",
		30, vmdisk.HostCachingTypeReadWrite)
	ConfigureWithExistingDataDisk(&role, "data-disk", vmdisk.HostCachingTypeReadOnly)

	bytes, err := xml.MarshalIndent(role, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	expected := `<Role>
  <RoleName>blobvm</RoleName>
  <RoleType>PersistentVMRole</RoleType>
  <ConfigurationSets>
    <ConfigurationSet>
      <ConfigurationSetType>WindowsProvisioningConfiguration</ConfigurationSetType>
      <ComputerName>WINVM</ComputerName>
      <AdminPassword>P2ssw@rd</AdminPassword>
      <EnableAutomaticUpdates>true</EnableAutomaticUpdates>
      <DomainJoin>
        <Credentials>
          <Domain></Domain>
          <Username>user@domain.com</Username>
          <Password>youReN3verG0nnaGu3ss</Password>
        </Credentials>
        <JoinDomain>redmond.corp.contoso.com</JoinDomain>
      </DomainJoin>
      <AdminUsername>azuser</AdminUsername>
    </ConfigurationSet>
  </ConfigurationSets>
  <DataVirtualHardDisks>
    <DataVirtualHardDisk>
      <DiskLabel>my-brand-new-disk</DiskLabel>
      <HostCaching>ReadWrite</HostCaching>
      <MediaLink>http://account.blob.core.windows.net/vhds/newdatadisk.vhd</MediaLink>
      <LogicalDiskSizeInGB>30</LogicalDiskSizeInGB>
    </DataVirtualHardDisk>
    <DataVirtualHardDisk>
      <DiskName>data-disk</DiskName>
      <HostCaching>ReadOnly</HostCaching>
      <Lun>1</Lun>
    </DataVirtualHardDisk>
  </DataVirtualHardDisks>
  <OSVirtualHardDisk>
    <DiskName>myvm-backup-20150209</DiskName>
    <DiskLabel>OSDisk</DiskLabel>
  </OSVirtualHardDisk>
  <RoleSize>Standard_D14</RoleSize>
  <ProvisionGuestAgent>true</ProvisionGuestAgent>
</Role>`

	if string(bytes) != expected {
		t.Fatalf("Expected marshalled xml to be %q, but got %q", expected, string(bytes))
	}
}
