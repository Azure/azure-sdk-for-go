package vmutils

import (
	"encoding/xml"
	"testing"
)

func TestNewLinuxVmPlatformImage(t *testing.T) {
	role, _ := NewVmConfiguration("myplatformvm", "Standard_D3")
	ConfigureDeploymentFromPlatformImage(role,
		"b39f27a8b8c64d52b05eac6a62ebad85__Ubuntu-14_04_2_LTS-amd64-server-20150309-en-us-30GB",
		"http://mystorageacct.blob.core.windows.net/vhds/mybrandnewvm.vhd", "mydisklabel")
	ConfigureForLinux(role, "myvm", "azureuser", "", "2398yyKJGd78e2389ydfncuirdebhf89yh3IUOBY")

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
		t.Fatalf("Expected marshalled xml to be \"%v\", but got \"%v\"", expected, string(bytes))
	}
}
