package vmutils

import (
	"encoding/xml"
	"testing"

	vm "github.com/MSOpenTech/azure-sdk-for-go/management/virtualmachine"
)

func Test_AddAzureVMExtensionConfiguration(t *testing.T) {

	role := vm.Role{}
	AddAzureVMExtensionConfiguration(&role,
		"nameOfExtension", "nameOfPublisher", "versionOfExtension", "nameOfReference", "state", []byte{}, []byte{})

	data, err := xml.MarshalIndent(role, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if expected := `<Role>
  <ResourceExtensionReferences>
    <ResourceExtensionReference>
      <ReferenceName>nameOfReference</ReferenceName>
      <Publisher>nameOfPublisher</Publisher>
      <Name>nameOfExtension</Name>
      <Version>versionOfExtension</Version>
      <ResourceExtensionParameterValues>
        <ResourceExtensionParameterValue>
          <Key>ignored</Key>
          <Value></Value>
          <Type>Private</Type>
        </ResourceExtensionParameterValue>
        <ResourceExtensionParameterValue>
          <Key>ignored</Key>
          <Value></Value>
          <Type>Public</Type>
        </ResourceExtensionParameterValue>
      </ResourceExtensionParameterValues>
      <State>state</State>
    </ResourceExtensionReference>
  </ResourceExtensionReferences>
</Role>`; string(data) != expected {
		t.Fatalf("Expected %q, but got %q", expected, string(data))
	}
}
