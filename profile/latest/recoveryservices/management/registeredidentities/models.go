package registeredidentities

import (
	 original "github.com/Azure/azure-sdk-for-go/service/recoveryservices/management/2016-06-01/registeredidentities"
)

type (
	 ManagementClient = original.ManagementClient
	 AuthType = original.AuthType
	 CertificateRequest = original.CertificateRequest
	 RawCertificateData = original.RawCertificateData
	 ResourceCertificateAndAadDetails = original.ResourceCertificateAndAadDetails
	 ResourceCertificateAndAcsDetails = original.ResourceCertificateAndAcsDetails
	 ResourceCertificateDetails = original.ResourceCertificateDetails
	 VaultCertificateResponse = original.VaultCertificateResponse
	 GroupClient = original.GroupClient
	 VaultCertificatesClient = original.VaultCertificatesClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 AAD = original.AAD
	 AccessControlService = original.AccessControlService
	 ACS = original.ACS
	 AzureActiveDirectory = original.AzureActiveDirectory
	 Invalid = original.Invalid
)
