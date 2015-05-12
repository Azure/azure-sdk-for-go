package vmutils

import (
	"fmt"

	. "github.com/Azure/azure-sdk-for-go/management/virtualmachine"
)

// Adds configuration exposing port 22 externally
func ConfigureWithPublicSSH(role *Role) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	return ConfigureWithExternalPort(role, "SSH", 22, 22, InputEndpointProtocolTcp)
}

// Adds configuration exposing port 3389 externally
func ConfigureWithPublicRDP(role *Role) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	return ConfigureWithExternalPort(role, "RDP", 3389, 3389, InputEndpointProtocolTcp)
}

// Adds configuration exposing port 5986 externally
func ConfigureWithPublicPowerShell(role *Role) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	return ConfigureWithExternalPort(role, "PowerShell", 5986, 5986, InputEndpointProtocolTcp)
}

// Adds a new InputEndpoint to the Role, exposing a port externally
func ConfigureWithExternalPort(role *Role, name string, localport, externalport int, protocol InputEndpointProtocol) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	role.ConfigurationSets = updateOrAddConfig(role.ConfigurationSets, ConfigurationSetTypeNetwork,
		func(config *ConfigurationSet) {
			if config.InputEndpoints == nil {
				config.InputEndpoints = []InputEndpoint{}
			}
			newInputEndpoints := append(config.InputEndpoints, InputEndpoint{
				LocalPort: localport,
				Name:      name,
				Port:      externalport,
				Protocol:  protocol,
			})
			config.InputEndpoints = newInputEndpoints
		})
	return nil
}

// ConfigureWithSecurityGroup associates the Role with a specific network security group
func ConfigureWithSecurityGroup(role *Role, networkSecurityGroup string) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	role.ConfigurationSets = updateOrAddConfig(role.ConfigurationSets, ConfigurationSetTypeNetwork,
		func(config *ConfigurationSet) {
			config.NetworkSecurityGroup = networkSecurityGroup
		})

	return nil
}

// ConfigureWithSubnet associates the Role with a specific subnet
func ConfigureWithSubnet(role *Role, subnet string) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	role.ConfigurationSets = updateOrAddConfig(role.ConfigurationSets, ConfigurationSetTypeNetwork,
		func(config *ConfigurationSet) {
			if config.SubnetNames == nil {
				config.SubnetNames = []string{}
			}

			config.SubnetNames = append(config.SubnetNames, subnet)
		})

	return nil
}
