package vmutils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	. "github.com/Azure/azure-sdk-for-go/management/virtualmachine"
)

const (
	dockerPublicConfigVersion = 2
)

func AddAzureVMExtensionConfiguration(role *Role, name, publisher, version, referenceName, state string,
	publicConfigurationValue, privateConfigurationValue []byte) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	extension := ResourceExtensionReference{
		Name:          name,
		Publisher:     publisher,
		Version:       version,
		ReferenceName: referenceName,
		State:         state,
	}

	if len(privateConfigurationValue) == 0 {
		extension.ParameterValues = append(extension.ParameterValues, ResourceExtensionParameter{
			Key:   "ignored",
			Value: base64.StdEncoding.EncodeToString(privateConfigurationValue),
			Type:  "Private",
		})
	}

	if len(publicConfigurationValue) == 0 {
		extension.ParameterValues = append(extension.ParameterValues, ResourceExtensionParameter{
			Key:   "ignored",
			Value: base64.StdEncoding.EncodeToString(publicConfigurationValue),
			Type:  "Public",
		})
	}

	if role.ResourceExtensionReferences == nil {
		role.ResourceExtensionReferences = &[]ResourceExtensionReference{extension}
	} else {
		*role.ResourceExtensionReferences = append(*role.ResourceExtensionReferences, extension)
	}

	return nil
}

// Adds the DockerExtension to the role configuratioon and opens a port "dockerPort"
func AddAzureDockerVMExtensionConfiguration(role *Role, dockerPort int, version string) error {
	if role == nil {
		return fmt.Errorf(errParamNotSpecified, "role")
	}

	ConfigureWithExternalPort(role, "docker", dockerPort, dockerPort, InputEndpointProtocolTcp)

	publicConfiguration, err := createDockerPublicConfig(dockerPort)
	if err != nil {
		return err
	}

	privateConfiguration, err := json.Marshal(dockerPrivateConfig{})
	if err != nil {
		return err
	}

	return AddAzureVMExtensionConfiguration(role,
		"DockerExtension", "MSOpenTech.Extensions",
		version, "DockerExtension", "enable",
		publicConfiguration, privateConfiguration)
}

func createDockerPublicConfig(dockerPort int) ([]byte, error) {
	return json.Marshal(dockerPublicConfig{DockerPort: dockerPort, Version: dockerPublicConfigVersion})
}

type dockerPublicConfig struct {
	DockerPort int `json:"dockerport"`
	Version    int `json:"version"`
}

type dockerPrivateConfig struct{}
