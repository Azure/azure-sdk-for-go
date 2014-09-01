package vmClient

import (
	"fmt"
	"time"
	"encoding/xml"
	"encoding/base64"
	"encoding/pem"
	"os"
	"io/ioutil"
	"crypto/rand"
	"crypto/sha1"
	"io"
	"errors"
	"strings"
	"os/user"
	"path"
	"github.com/MSOpenTech/azure-sdk-for-go/clients/locationClient"
	"github.com/MSOpenTech/azure-sdk-for-go/clients/imageClient"
	"github.com/MSOpenTech/azure-sdk-for-go/clients/storageServiceClient"
	azure "github.com/MSOpenTech/azure-sdk-for-go"
)

// REGION PUBLIC METHODS STARTS

func CreateAzureVM(role *Role, dnsName, location string) error {

	err := locationClient.ResolveLocation(location)
	if err != nil {
		return err
	}

	fmt.Println("Creating hosted service... ")
	requestId, err := CreateHostedService(dnsName, location)
	if err != nil {
		return err
	}

	azure.WaitAsyncOperation(requestId)
	fmt.Println("done.")

	if role.UseCertAuth {
		fmt.Println("Uploading cert...")

		err = uploadServiceCert(dnsName, role.CertPath)
		if err != nil {
			return err
		}

		fmt.Println("done.")
	}

	fmt.Println("Deploying azure VM configuration... ")

	vMDeployment := createVMDeploymentConfig(role)
	vMDeploymentBytes, err := xml.Marshal(vMDeployment)
	if err != nil {
		return err
	}

	requestURL :=  fmt.Sprintf("services/hostedservices/%s/deployments", role.RoleName)
	requestId, err = azure.SendAzurePostRequest(requestURL, vMDeploymentBytes)
	if err != nil {
		return err
	}

	azure.WaitAsyncOperation(requestId)

	fmt.Println("done.")
	return nil
}

func CreateHostedService(dnsName, location string) (string, error) {

	err := locationClient.ResolveLocation(location)
	if err != nil {
		return "", err
	}

	hostedServiceDeployment := createHostedServiceDeploymentConfig(dnsName, location)
	hostedServiceBytes, err := xml.Marshal(hostedServiceDeployment)
	if err != nil {
		return "", err
	}

	requestURL :=  "services/hostedservices"
	requestId, azureErr := azure.SendAzurePostRequest(requestURL, hostedServiceBytes)
	if azureErr != nil {
		return "", err
	}

	return requestId, nil
}

func CreateAzureVMConfiguration(name, instanceSize, imageName, location string) (*Role, error) {
	fmt.Println("Creating azure VM configuration... ")

	err := locationClient.ResolveLocation(location)
	if err != nil {
		return nil, err
	}

	role, err := createAzureVMRole(name, instanceSize, imageName, location)
	if err != nil {
		return nil, err
	}

	fmt.Println("done.")
	return role, nil
}

func AddAzureLinuxProvisioningConfig(azureVMConfig *Role, userName, password, certPath string) (*Role, error) {
	fmt.Println("Adding azure provisioning configuration... ")

	configurationSets := ConfigurationSets{}

	provisioningConfig, err := createLinuxProvisioningConfig(azureVMConfig.RoleName, userName, password, certPath)
	if err != nil {
		return nil, err
	}

	configurationSets.ConfigurationSet = append(configurationSets.ConfigurationSet, provisioningConfig)

	networkConfig, networkErr := createNetworkConfig("Linux")
	if networkErr != nil {
		return nil, err
	}

	configurationSets.ConfigurationSet = append(configurationSets.ConfigurationSet, networkConfig)

	azureVMConfig.ConfigurationSets = configurationSets

	if len(certPath) > 0 {
		azureVMConfig.UseCertAuth = true
		azureVMConfig.CertPath = certPath
	}

	fmt.Println("done.")
	return azureVMConfig, nil
}

func SetAzureVMExtension(azureVMConfiguration *Role, name string, publisher string, version string, referenceName string, state string, publicConfigurationValue string, privateConfigurationValue string) (*Role) {
	fmt.Printf("Setting azure VM extension: %s... \n", name)

	extension := ResourceExtensionReference{}
	extension.Name = name
	extension.Publisher = publisher
	extension.Version = version
	extension.ReferenceName = referenceName
	extension.State = state

	if len(privateConfigurationValue) > 0 {
		privateConfig := ResourceExtensionParameter{}
		privateConfig.Key = "ignored"
		privateConfig.Value = base64.StdEncoding.EncodeToString([]byte(privateConfigurationValue))
		privateConfig.Type = "Private"

		extension.ResourceExtensionParameterValues.ResourceExtensionParameterValue = append(extension.ResourceExtensionParameterValues.ResourceExtensionParameterValue, privateConfig)
	}

	if len(publicConfigurationValue) > 0 {
		publicConfig := ResourceExtensionParameter{}
		publicConfig.Key = "ignored"
		publicConfig.Value = base64.StdEncoding.EncodeToString([]byte(publicConfigurationValue))
		publicConfig.Type = "Public"

		extension.ResourceExtensionParameterValues.ResourceExtensionParameterValue = append(extension.ResourceExtensionParameterValues.ResourceExtensionParameterValue, publicConfig)
	}

	azureVMConfiguration.ResourceExtensionReferences.ResourceExtensionReference = append(azureVMConfiguration.ResourceExtensionReferences.ResourceExtensionReference, extension)

	return azureVMConfiguration
}

func SetAzureDockerVMExtension(azureVMConfiguration *Role, dockerCertDir string, dockerPort int, version string) (*Role, error) {
	if len(version) == 0 {
		version = "0.3"
	}

	err := addDockerPort(azureVMConfiguration.ConfigurationSets.ConfigurationSet, dockerPort)
	if err != nil {
		return nil, err
	}

	publicConfiguration := createDockerPublicConfig(dockerPort)
	privateConfiguration, err := createDockerPrivateConfig(dockerCertDir)
	if err != nil {
		return nil, err
	}

	azureVMConfiguration = SetAzureVMExtension(azureVMConfiguration, "DockerExtension", "MSOpenTech.Extensions", version, "DockerExtension", "enable", publicConfiguration, privateConfiguration)
	return azureVMConfiguration, nil
}

func GetVMDeployment(cloudserviceName, deploymentName string) (*VMDeployment, error) {
	deployment := new(VMDeployment)

	requestURL :=  fmt.Sprintf("services/hostedservices/%s/deployments/%s", cloudserviceName, deploymentName)
	response, azureErr := azure.SendAzureGetRequest(requestURL)
	if azureErr != nil {
		if strings.Contains(azureErr.Error(), "Code: ResourceNotFound") {
			return nil, nil
		}

		return nil, azureErr
	}

	err := xml.Unmarshal(response, deployment)
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func GetRole(cloudserviceName, deploymentName, roleName string) (*Role, error) {
	role := new(Role)

	requestURL :=  fmt.Sprintf("services/hostedservices/%s/deployments/%s/roles/%s", cloudserviceName, deploymentName, roleName)
	response, azureErr := azure.SendAzureGetRequest(requestURL)
	if azureErr != nil {
		if strings.Contains(azureErr.Error(), "Code: ResourceNotFound") {
			return nil, nil
		}

		return nil, azureErr
	}

	err := xml.Unmarshal(response, role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func StartRole(cloudserviceName, deploymentName, roleName string) (error) {
	startRoleOperation := createStartRoleOperation()

	startRoleOperationBytes, err := xml.Marshal(startRoleOperation)
	if err != nil {
		return err
	}

	requestURL :=  fmt.Sprintf("services/hostedservices/%s/deployments/%s/roleinstances/%s/Operations", cloudserviceName, deploymentName, roleName)
	requestId, azureErr := azure.SendAzurePostRequest(requestURL, startRoleOperationBytes)
	if azureErr != nil {
		return azureErr
	}

	azure.WaitAsyncOperation(requestId)
	return nil
}

// REGION PUBLIC METHODS ENDS


// REGION PRIVATE METHODS STARTS

func createStartRoleOperation() StartRoleOperation {
	startRoleOperation := StartRoleOperation{}
	startRoleOperation.OperationType = "StartRoleOperation"
	startRoleOperation.Xmlns = "http://schemas.microsoft.com/windowsazure"

	return startRoleOperation
}

func createDockerPublicConfig(dockerPort int) string{
	config := fmt.Sprintf("{ \"dockerport\": \"%v\" }", dockerPort)
	return config
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
		return "", err
	}

	certDir := path.Join(usr.HomeDir, dockerCertDir)

	if _, err := os.Stat(certDir); err == nil {
		fmt.Println("Docker directory exists")
	} else {
		fmt.Println("Docker directory does NOT exists")
		mkdirErr := os.Mkdir(certDir, 0644)
		if mkdirErr != nil {
			return "", err
		}
	}

	//generateDockerCertificates(certDir)
	caCert, err := parseFileToBase64String(path.Join(certDir, "ca.pem"))
	if err != nil {
		return "", err
	}

	serverCert, err := parseFileToBase64String(path.Join(certDir, "server-cert.pem"))
	if err != nil {
		return "", err
	}

	serverKey, err := parseFileToBase64String(path.Join(certDir, "server-key.pem"))
	if err != nil {
		return "", err
	}

	config := fmt.Sprintf("{ \"ca\": \"%s\", \"server-cert\": \"%s\", \"server-key\": \"%s\" }", caCert, serverCert, serverKey)
	return config, nil
}

func parseFileToBase64String(path string) (string, error) {
	caCertBytes, caErr := ioutil.ReadFile(path)
	if caErr != nil {
		return "", caErr
	}

	base64Content := base64.StdEncoding.EncodeToString(caCertBytes)
	return base64Content, nil
}

func addDockerPort(configurationSets []ConfigurationSet,  dockerPort int) error {
	if len(configurationSets) == 0 {
		return errors.New("You should set azure VM provisioning config first")
	}

	for i := 0; i < len(configurationSets); i++ {
		if configurationSets[i].ConfigurationSetType != "NetworkConfiguration" {
			continue
		}

		dockerEndpoint := createEndpoint("docker", "tcp", dockerPort, dockerPort)
		configurationSets[i].InputEndpoints.InputEndpoint = append(configurationSets[i].InputEndpoints.InputEndpoint, dockerEndpoint)
	}

	return nil
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

func createVMDeploymentConfig(role *Role) (VMDeployment) {
	deployment := VMDeployment{}
	deployment.Name = role.RoleName
	deployment.Xmlns = "http://schemas.microsoft.com/windowsazure"
	deployment.DeploymentSlot = "Production"
	deployment.Label = role.RoleName
	deployment.RoleList.Role = append(deployment.RoleList.Role, role)

	return deployment
}

func createAzureVMRole(name, instanceSize, imageName, location string) (*Role, error){
	config := new(Role)
	config.RoleName = name
	config.RoleSize = instanceSize
	config.RoleType = "PersistentVMRole"
	config.ProvisionGuestAgent = true
	var err error
	config.OSVirtualHardDisk, err = createOSVirtualHardDisk(name, imageName, location)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func createOSVirtualHardDisk(dnsName, imageName, location string) (OSVirtualHardDisk, error){
	oSVirtualHardDisk := OSVirtualHardDisk{}

	err := imageClient.ResolveImageName(imageName)
	if err != nil {
		return oSVirtualHardDisk, err
	}

	oSVirtualHardDisk.SourceImageName = imageName
	oSVirtualHardDisk.MediaLink, err = getVHDMediaLink(dnsName, location)
	if err != nil {
		return oSVirtualHardDisk, err
	}

	return oSVirtualHardDisk, nil
}

func getVHDMediaLink(dnsName, location string) (string, error){

	storageService, err := storageServiceClient.GetStorageServiceByLocation(location)
	if err != nil {
		return "", err
	}

	if storageService == nil {

		uuid, err := newUUID()
		if err != nil {
			return "", err
		}

		serviceName := "portalvhds" + uuid
		storageService, err = storageServiceClient.CreateStorageService(serviceName, location)
		if err != nil {
			return "", err
		}
	}

	blobEndpoint, err := storageServiceClient.GetBlobEndpoint(storageService)
	if err != nil {
		return "", err
	}

	vhdMediaLink := blobEndpoint + "vhds/" + dnsName + "-" + time.Now().Local().Format("20060102150405") + ".vhd"
	return vhdMediaLink, nil
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

func createLinuxProvisioningConfig(dnsName, userName, userPassword, certPath string) (ConfigurationSet, error) {
	provisioningConfig := ConfigurationSet{}

	disableSshPasswordAuthentication := false
	if len(userPassword) == 0 {
		disableSshPasswordAuthentication = true
		// We need to set dummy password otherwise azure API will throw an error
		userPassword = "P@ssword1"
	}

	provisioningConfig.DisableSshPasswordAuthentication = disableSshPasswordAuthentication
	provisioningConfig.ConfigurationSetType = "LinuxProvisioningConfiguration"
	provisioningConfig.HostName = dnsName
	provisioningConfig.UserName = userName
	provisioningConfig.UserPassword = userPassword

	if len(certPath) > 0 {
		var err error
		provisioningConfig.SSH, err = createSshConfig(certPath, userName)
		if err != nil {
			return provisioningConfig, err
		}
	}

	return provisioningConfig, nil
}

func uploadServiceCert(dnsName, certPath string) (error) {
	certificateConfig, err := createServiceCertDeploymentConf(certPath)
	if err != nil {
		return err
	}

	certificateConfigBytes, err := xml.Marshal(certificateConfig)
	if err != nil {
		return err
	}

	requestURL :=  fmt.Sprintf("services/hostedservices/%s/certificates", dnsName)
	requestId, azureErr := azure.SendAzurePostRequest(requestURL, certificateConfigBytes)
	if azureErr != nil {
		return azureErr
	}

	err = azure.WaitAsyncOperation(requestId)
	return err
}

func createServiceCertDeploymentConf(certPath string) (ServiceCertificate, error) {
	certConfig := ServiceCertificate{}
	certConfig.Xmlns = "http://schemas.microsoft.com/windowsazure"
	data , err := ioutil.ReadFile(certPath)
	if err != nil {
		return certConfig, err
	}

	certData := base64.StdEncoding.EncodeToString(data)
	certConfig.Data = certData
	certConfig.CertificateFormat = "pfx"

	return certConfig, nil
}

func createSshConfig(certPath, userName string) (SSH, error) {
	sshConfig := SSH{}
	publicKey := PublicKey{}

	err := checkServiceCertExtension(certPath)
	if err != nil {
		return sshConfig, err
	}

	fingerprint, err := getServiceCertFingerprint(certPath)
	if err != nil {
		return sshConfig, err
	}

	publicKey.Fingerprint = fingerprint
	publicKey.Path = "/home/" + userName + "/.ssh/authorized_keys"

	sshConfig.PublicKeys.PublicKey = append(sshConfig.PublicKeys.PublicKey, publicKey)
	return sshConfig, nil
}

func getServiceCertFingerprint(certPath string) (string, error) {
	certData, readErr := ioutil.ReadFile(certPath)
	if readErr != nil {
		return "", readErr
	}

	block, rest := pem.Decode(certData)
	if block == nil {
		return "", errors.New(string(rest))
	}

	sha1sum := sha1.Sum(block.Bytes)
	fingerprint := fmt.Sprintf("%X", sha1sum)
	return fingerprint, nil
}

func checkServiceCertExtension(certPath string) (error) {
	certParts := strings.Split(certPath, ".")
	certExt := certParts[len(certParts) - 1]

	acceptedExtension := "pem"
	if certExt != acceptedExtension {
		return errors.New(fmt.Sprintf("Certificate %s is invalid. Please specify %s certificate.", certPath, acceptedExtension))
	}

	return nil
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
