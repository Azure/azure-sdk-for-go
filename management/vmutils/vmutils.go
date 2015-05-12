// Package vmutils provides convenience methods for building vm role configurations
package vmutils

import (
	"fmt"

	. "github.com/Azure/azure-sdk-for-go/management/virtualmachine"
)

const (
	errParamNotSpecified = "Parameter %s is not specified."
)

// Creates configuration for a new virtual machine Role
func NewVmConfiguration(name string, roleSize string) Role {
	return Role{
		RoleName:            name,
		RoleType:            "PersistentVMRole",
		RoleSize:            roleSize,
		ProvisionGuestAgent: true,
	}
}

// Adds configuration for when deploying a generalized Linux image. If "password" is left empty,
// SSH password security will be disabled by default. Certificates with SSH public keys
// should already be uploaded to the cloud service where the VM will be deployed
// and referenced here only by their thumbprint.
func ConfigureForLinux(role *Role, hostname, user, password string, sshPubkeyCertificateThumbprint ...string) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	role.ConfigurationSets = updateOrAddConfig(role.ConfigurationSets, ConfigurationSetTypeLinuxProvisioning,
		func(config *ConfigurationSet) {
			config.HostName = hostname
			config.UserName = user
			config.UserPassword = password
			if password != "" {
				config.DisableSshPasswordAuthentication = "false"
			}
			if len(sshPubkeyCertificateThumbprint) != 0 {
				config.SSH = &SSH{}
				var keys []PublicKey
				if config.SSH.PublicKeys == nil {
					keys = []PublicKey{}
				} else {
					keys = config.SSH.PublicKeys
				}
				for _, k := range sshPubkeyCertificateThumbprint {
					keys = append(keys,
						PublicKey{
							Fingerprint: k,
							Path:        "/home/" + user + "/.ssh/authorized_keys",
						})
				}
				config.SSH.PublicKeys = keys
			}
		})
	return nil
}

// Adds configuration for when deploying a generalized Windows image.
// timeZone can be left empty. For a complete list of supported time zone entries, you can either
// refer to the values listed in the registry entry "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Time Zones"
// or you can use the tzutil command-line tool to list the valid time.
func ConfigureForWindows(role *Role, hostname, user, password string, enableAutomaticUpdates bool, timeZone string) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	role.ConfigurationSets = updateOrAddConfig(role.ConfigurationSets, ConfigurationSetTypeWindowsProvisioning,
		func(config *ConfigurationSet) {
			config.ComputerName = hostname
			config.AdminUsername = user
			config.AdminPassword = password
			config.EnableAutomaticUpdates = enableAutomaticUpdates
			config.TimeZone = timeZone
		})
	return nil
}

// Adds configuration to join a new Windows vm to a domain. "username" must be in UPN form (user@domain.com),
// "machineOU" can be left empty
func ConfigureWindowsToJoinDomain(role *Role, username, password, domainToJoin, machineOU string) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	winconfig := findConfig(role.ConfigurationSets, ConfigurationSetTypeWindowsProvisioning)
	if winconfig != nil {
		winconfig.DomainJoin = &DomainJoin{
			Credentials:     Credentials{Username: username, Password: password},
			JoinDomain:      domainToJoin,
			MachineObjectOU: machineOU,
		}
	}
	return nil
}
