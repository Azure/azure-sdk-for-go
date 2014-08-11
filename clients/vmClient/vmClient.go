package vmClient

import (
	"fmt"
	"time"
	"encoding/xml"
	"encoding/base64"
	"os"
	"io/ioutil"
	"crypto/rand"
	"io"
	"errors"
	"os/user"
	"path"
	"github.com/MSOpenTech/azure-sdk-for-go/clients/locationClient"
	"github.com/MSOpenTech/azure-sdk-for-go/clients/imageClient"
	"github.com/MSOpenTech/azure-sdk-for-go/clients/storageServiceClient"
	azure "github.com/MSOpenTech/azure-sdk-for-go"
)

// REGION PUBLIC METHODS STARTS

func CreateAzureVM(role Role, dnsName, location string) {

	err := locationClient.ResolveLocation(location)
	if err != nil {
		azure.PrintErrorAndExit(err)
	}

	CreateHostedService(dnsName, location)

	fmt.Println("Deploying azure VM configuration ... ")

	vMDeployment := createVMDeploymentConfig(role)
	vMDeploymentBytes, err := xml.Marshal(vMDeployment)
	if err != nil {
		azure.PrintErrorAndExit(err)
	}

	time.Sleep(10)

	requestURL :=  fmt.Sprintf("https://management.core.windows.net/%s/services/hostedservices/%s/deployments", azure.GetPublishSettings().SubscriptionID, role.RoleName)
	azure.SendAzurePostRequest(requestURL, vMDeploymentBytes)
	fmt.Println("done. ")
}

func CreateHostedService(dnsName, location string) {
	fmt.Println("Creating hosted service ... ")
	hostedServiceDeployment := createHostedServiceDeploymentConfig(dnsName, location)
	hostedServiceBytes, err := xml.Marshal(hostedServiceDeployment)
	if err != nil {
		azure.PrintErrorAndExit(err)
	}

	requestURL :=  fmt.Sprintf("https://management.core.windows.net/%s/services/hostedservices", azure.GetPublishSettings().SubscriptionID)
	azure.SendAzurePostRequest(requestURL, hostedServiceBytes)
	fmt.Println("done. ")
}

func CreateAzureVMConfiguration(name, instanceSize, imageName, location string) (Role) {
	fmt.Println("Creating azure VM configuration ... ")

	role := createAzureVMRole(name, instanceSize, imageName, location)
	fmt.Println("done. ")
	return role
}

func AddAzureProvisioningConfig(azureVMConfig Role, userName, password, os string) (Role) {
	fmt.Println("Adding azure provisioning configuration ... ")

	azureVMConfig.ConfigurationSets = createConfigurationSets(azureVMConfig.RoleName, userName, password, os)

	return azureVMConfig
}

func SetAzureVMExtension(azureVMConfiguration Role, name string, publisher string, version string, referenceName string, state string, publicConfigurationValue string, privateConfigurationValue string) (Role) {
	fmt.Printf("Setting azure VM extension: %s  ... \n", name)

	extension := ResourceExtensionReference{}
	extension.Name = name
	extension.Publisher = publisher
	extension.Version = version
	extension.ReferenceName = referenceName
	extension.State = state

	if len(privateConfigurationValue) > 0 {
		privateConfig := ResourceExtensionParameter{}
		privateConfig.Key = "ignored"
		privateConfig.Value = privateConfigurationValue
		privateConfig.Type = "Private"

		extension.ResourceExtensionParameterValues.ResourceExtensionParameterValue = append(extension.ResourceExtensionParameterValues.ResourceExtensionParameterValue, privateConfig)
	}

	if len(publicConfigurationValue) > 0 {
		publicConfig := ResourceExtensionParameter{}
		publicConfig.Key = "ignored"
		publicConfig.Value = publicConfigurationValue
		publicConfig.Type = "Public"

		extension.ResourceExtensionParameterValues.ResourceExtensionParameterValue = append(extension.ResourceExtensionParameterValues.ResourceExtensionParameterValue, publicConfig)
	}

	azureVMConfiguration.ResourceExtensionReferences.ResourceExtensionReference = append(azureVMConfiguration.ResourceExtensionReferences.ResourceExtensionReference, extension)

	return azureVMConfiguration
}

func SetAzureDockerVMExtension(azureVMConfiguration Role, dockerCertDir string, dockerPort int, version string) (Role) {
	if len(version) == 0 {
		version = "0.3"
	}

	addDockerPort(azureVMConfiguration.ConfigurationSets.ConfigurationSet, dockerPort)
	publicConfiguration := createDockerPublicConfig(dockerPort)
	//privateConfiguration, err := createDockerPrivateConfig(dockerCertDir)
	//if err != nil {
	//	azure.PrintErrorAndExit(err)
	//}

	azureVMConfiguration = SetAzureVMExtension(azureVMConfiguration, "DockerExtension", "MSOpenTech.Extensions", version, "DockerExtension", "enable", publicConfiguration, "")
	return azureVMConfiguration
}

// REGION PUBLIC METHODS ENDS


// REGION PRIVATE METHODS STARTS


func createDockerPublicConfig(dockerPort int) string{
	config := fmt.Sprintf("{ \"dockerport\": \"%v\" }", dockerPort)
	return base64.StdEncoding.EncodeToString([]byte(config))
}

func generateDockerCertificates(dockerCertDir string) {
	password := "Docker123"
	ca := path.Join(dockerCertDir, "ca.pem")
	caKey := path.Join(dockerCertDir, "ca-key.pem")
	serverKey := path.Join(dockerCertDir, "server-key.pem")
	server := path.Join(dockerCertDir, "server.csr")
	serverCert := path.Join(dockerCertDir, "server-cert.pem")
	clientKey := path.Join(dockerCertDir, "client-key.pem")
	client := path.Join(dockerCertDir, "client.csr")
	clientCert := path.Join(dockerCertDir, "client-cert.pem")

	azure.ExecuteCommand(fmt.Sprintf("openssl genrsa -des3 -passout pass:%s -out %s 2048", password, caKey))
	time.Sleep(1000 * time.Millisecond)
	azure.ExecuteCommand(fmt.Sprintf("openssl req -new -x509 -days 365 -passin pass:%s -subj '/C=AU/ST=Some-State/O=Test/CN=\\*' -key %s -out %s", password, caKey, ca))
	time.Sleep(1000 * time.Millisecond)
	azure.ExecuteCommand(fmt.Sprintf("openssl genrsa -des3 -passout pass:%s -out %s 2048", password, serverKey))
	time.Sleep(1000 * time.Millisecond)
	azure.ExecuteCommand(fmt.Sprintf("openssl req -new -passin pass:%s -subj '/C=AU/ST=Some-State/O=Test/CN=\\*' -key %s -out %s", password, serverKey, server))
	time.Sleep(1000 * time.Millisecond)
	azure.ExecuteCommand(fmt.Sprintf("openssl x509 -req -days 365 -passin pass:%s -set_serial 01 -in %s -CA %s -CAkey %s -out %s", password, server, ca, caKey, serverCert))
	time.Sleep(1000 * time.Millisecond)
	azure.ExecuteCommand(fmt.Sprintf("openssl genrsa -des3 -passout pass:%s -out %s 2048", password, clientKey))
	time.Sleep(1000 * time.Millisecond)
	azure.ExecuteCommand(fmt.Sprintf("openssl req -passin pass:%s -subj '/C=AU/ST=Some-State/O=Test/CN=\\*' -new -key %s -out %s", password, clientKey, client))
	time.Sleep(1000 * time.Millisecond)
	exFile := path.Join(dockerCertDir, "extfile.cnf")
	ioutil.WriteFile(exFile, []byte("extendedKeyUsage = clientAuth"), 0644)

	azure.ExecuteCommand(fmt.Sprintf("openssl x509 -req -days 365 -passin pass:%s -set_serial 01 -in %s -CA %s -CAkey %s -out %s -extfile %s", password, client, ca, caKey, clientCert, exFile))
	time.Sleep(1000 * time.Millisecond)
	azure.ExecuteCommand(fmt.Sprintf("openssl rsa -passin pass:%s -passout pass:%s -in %s -out server-key.pem", password, password, serverKey, serverKey))
	time.Sleep(1000 * time.Millisecond)
	azure.ExecuteCommand(fmt.Sprintf("openssl rsa -passin pass:%s -passout pass:%s -in %s -out %s", password, password, clientKey, clientKey))
	time.Sleep(1000 * time.Millisecond)
}

func createDockerPrivateConfig(dockerCertDir string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		azure.PrintErrorAndExit(err)
	}

	certDir := path.Join(usr.HomeDir, dockerCertDir)

	if _, err := os.Stat(certDir); err == nil {
		fmt.Println("Docker directory exists")
	} else {
		fmt.Println("Docker directory does NOT exists")
		mkdirErr := os.Mkdir(certDir, 0644)
		if mkdirErr != nil {
			azure.PrintErrorAndExit(mkdirErr)
		}
	}

	generateDockerCertificates(certDir)
	caCert, caErr := parseFileToBase64String(path.Join(certDir, "ca.pem"))
	if caErr != nil {
		azure.PrintErrorAndExit(caErr)
	}

	serverCert, serverCertErr := parseFileToBase64String(path.Join(certDir, "server-cert.pem"))
	if serverCertErr != nil {
		azure.PrintErrorAndExit(serverCertErr)
	}

	serverKey, serverKeyErr := parseFileToBase64String(path.Join(certDir, "server-key.pem"))
	if serverKeyErr != nil {
		azure.PrintErrorAndExit(serverKeyErr)
	}

	config := fmt.Sprintf("{ \"ca\": \"%s\", \"server-cert\": \"%s\", \"server-key\": \"%s\" }", caCert, serverCert, serverKey)
	return base64.StdEncoding.EncodeToString([]byte(config)), nil
}

func parseFileToBase64String(path string) (string, error) {
	caCertBytes, caErr := ioutil.ReadFile(path)
	if caErr != nil {
		return "", caErr
	}

	base64Content := base64.StdEncoding.EncodeToString(caCertBytes)
	return base64Content, nil
}

func addDockerPort(configurationSets []ConfigurationSet,  dockerPort int) {
	if len(configurationSets) == 0 {
		azure.PrintErrorAndExit(errors.New("You should set azure VM provisioning config first"))
	}

	for i := 0; i < len(configurationSets); i++ {
		if configurationSets[i].ConfigurationSetType != "NetworkConfiguration" {
			continue
		}

		dockerEndpoint := createEndpoint("docker", "tcp", dockerPort, dockerPort)
		configurationSets[i].InputEndpoints.InputEndpoint = append(configurationSets[i].InputEndpoints.InputEndpoint, dockerEndpoint)
	}
}

func createHostedServiceDeploymentConfig(dnsName, location string) (HostedServiceDeployment) {
	deployment := HostedServiceDeployment{}
	deployment.ServiceName = dnsName
	label := base64.StdEncoding.EncodeToString([]byte(dnsName))
	deployment.Label = label
	deployment.Location = location
	deployment.Xmlns = "http://schemas.microsoft.com/windowsazure"

	return deployment
}

func createVMDeploymentConfig(role Role) (VMDeployment) {
	deployment := VMDeployment{}
	deployment.Name = role.RoleName
	deployment.Xmlns = "http://schemas.microsoft.com/windowsazure"
	deployment.DeploymentSlot = "Production"
	deployment.Label = role.RoleName
	deployment.RoleList.Role = append(deployment.RoleList.Role, role)

	return deployment
}

func createAzureVMRole(name, instanceSize, imageName, location string) (Role){
	config := Role{}
	config.RoleName = name
	config.RoleSize = instanceSize
	config.RoleType = "PersistentVMRole"
	config.ProvisionGuestAgent = true
	config.OSVirtualHardDisk = createOSVirtualHardDisk(name, imageName, location)

	return config
}

func createOSVirtualHardDisk(dnsName, imageName, location string) (OSVirtualHardDisk){
	oSVirtualHardDisk := OSVirtualHardDisk{}

	err := imageClient.ResolveImageName(imageName)
	if err != nil {
		azure.PrintErrorAndExit(err)
	}

	oSVirtualHardDisk.SourceImageName = imageName
	oSVirtualHardDisk.MediaLink = getVHDMediaLink(dnsName, location)

	return oSVirtualHardDisk
}

func getVHDMediaLink(dnsName, location string) (string){

	storageService, err := storageServiceClient.GetStorageServiceByLocation(location)
	if err != nil {
		azure.PrintErrorAndExit(err)
	}

	if storageService == nil {

		uuid, err := newUUID()
		if err != nil {
			azure.PrintErrorAndExit(err)
		}

		serviceName := "portalvhds" + uuid
		storageService = storageServiceClient.CreateStorageService(serviceName, location)
	}

	blobEndpoint, blobErr := storageServiceClient.GetBlobEndpoint(storageService)
	if blobErr != nil {
		azure.PrintErrorAndExit(blobErr)
	}

	vhdMediaLink := blobEndpoint + "vhds/" + dnsName + "-" + time.Now().Local().Format("20060102150405") + ".vhd"
	return vhdMediaLink
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	//return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
	return fmt.Sprintf("%x", uuid[10:]), nil
}

func createConfigurationSets(dnsName, userName, password, os string) (ConfigurationSets){
	configurationSets := ConfigurationSets{}

	provisioningConfig, err := createProvisioningConfig(dnsName, userName, password, os)
	if err != nil {
		azure.PrintErrorAndExit(err)
	}

	configurationSets.ConfigurationSet = append(configurationSets.ConfigurationSet, provisioningConfig)

	networkConfig, networkErr := createNetworkConfig(os)
	if networkErr != nil {
		azure.PrintErrorAndExit(networkErr)
	}

	configurationSets.ConfigurationSet = append(configurationSets.ConfigurationSet, networkConfig)

	return configurationSets
}

func createProvisioningConfig(dnsName string, userName string, userPassword, os string) (ConfigurationSet, error) {
	provisioningConfig := ConfigurationSet{}
	var configurationSetType string

	if os == "Linux" {
		configurationSetType = "LinuxProvisioningConfiguration"

		//!TODO If user specified cert base auth we should not set this value
		//if something {
		provisioningConfig.DisableSshPasswordAuthentication = false
		//}
	} else if os == "Windows" {
		configurationSetType = "WindowsProvisioningConfiguration"
	} else {
		return provisioningConfig, errors.New(fmt.Sprintf("You must specify correct OS param. Valid values are 'Linux' and 'Windows'"))
	}

	provisioningConfig.ConfigurationSetType = configurationSetType
	provisioningConfig.HostName = dnsName
	provisioningConfig.UserName = userName
	provisioningConfig.UserPassword = userPassword

	return provisioningConfig, nil
}

func createNetworkConfig(os string) (ConfigurationSet, error) {
	networkConfig := ConfigurationSet{}
	networkConfig.ConfigurationSetType = "NetworkConfiguration"

	var endpoint InputEndpoint
	if os == "Linux" {
		endpoint = createEndpoint("ssh", "tcp", 22, 22)
	} else if os == "Windows" {
		//!TODO should add rdp endpoint
	} else {
		return networkConfig, errors.New(fmt.Sprintf("You must specify correct OS param. Valid values are 'Linux' and 'Windows'"))
	}

	networkConfig.InputEndpoints.InputEndpoint = append(networkConfig.InputEndpoints.InputEndpoint, endpoint)

	//!TODO
	//networkConfig.SubnetNames

	return networkConfig, nil
}

func createEndpoint(name string, protocol string, extertalPort int, internalPort int) (InputEndpoint) {
	endpoint := InputEndpoint{}
	endpoint.Name = name
	endpoint.Protocol = protocol
	endpoint.Port = extertalPort
	endpoint.LocalPort = internalPort

	return endpoint
}

// REGION PRIVATE METHODS ENDS
