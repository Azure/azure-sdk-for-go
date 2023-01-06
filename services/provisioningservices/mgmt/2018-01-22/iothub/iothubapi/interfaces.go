// Deprecated: Please note, this package has been deprecated. A replacement package is available [github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/deviceprovisioningservices/armdeviceprovisioningservices](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/deviceprovisioningservices/armdeviceprovisioningservices). We strongly encourage you to upgrade to continue receiving updates. See [Migration Guide](https://aka.ms/azsdk/golang/t2/migration) for guidance on upgrading. Refer to our [deprecation policy](https://azure.github.io/azure-sdk/policies_support.html) for more details.
package iothubapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/provisioningservices/mgmt/2018-01-22/iothub"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
)

// OperationsClientAPI contains the set of methods on the OperationsClient type.
type OperationsClientAPI interface {
	List(ctx context.Context) (result iothub.OperationListResultPage, err error)
	ListComplete(ctx context.Context) (result iothub.OperationListResultIterator, err error)
}

var _ OperationsClientAPI = (*iothub.OperationsClient)(nil)

// DpsCertificateClientAPI contains the set of methods on the DpsCertificateClient type.
type DpsCertificateClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, provisioningServiceName string, certificateName string, certificateDescription iothub.CertificateBodyDescription, ifMatch string) (result iothub.CertificateResponse, err error)
	Delete(ctx context.Context, resourceGroupName string, ifMatch string, provisioningServiceName string, certificateName string, certificatename string, certificaterawBytes []byte, certificateisVerified *bool, certificatepurpose iothub.CertificatePurpose, certificatecreated *date.Time, certificatelastUpdated *date.Time, certificatehasPrivateKey *bool, certificatenonce string) (result autorest.Response, err error)
	GenerateVerificationCode(ctx context.Context, certificateName string, ifMatch string, resourceGroupName string, provisioningServiceName string, certificatename string, certificaterawBytes []byte, certificateisVerified *bool, certificatepurpose iothub.CertificatePurpose, certificatecreated *date.Time, certificatelastUpdated *date.Time, certificatehasPrivateKey *bool, certificatenonce string) (result iothub.VerificationCodeResponse, err error)
	Get(ctx context.Context, certificateName string, resourceGroupName string, provisioningServiceName string, ifMatch string) (result iothub.CertificateResponse, err error)
	List(ctx context.Context, resourceGroupName string, provisioningServiceName string) (result iothub.CertificateListDescription, err error)
	VerifyCertificate(ctx context.Context, certificateName string, ifMatch string, request iothub.VerificationCodeRequest, resourceGroupName string, provisioningServiceName string, certificatename string, certificaterawBytes []byte, certificateisVerified *bool, certificatepurpose iothub.CertificatePurpose, certificatecreated *date.Time, certificatelastUpdated *date.Time, certificatehasPrivateKey *bool, certificatenonce string) (result iothub.CertificateResponse, err error)
}

var _ DpsCertificateClientAPI = (*iothub.DpsCertificateClient)(nil)

// IotDpsResourceClientAPI contains the set of methods on the IotDpsResourceClient type.
type IotDpsResourceClientAPI interface {
	CheckProvisioningServiceNameAvailability(ctx context.Context, arguments iothub.OperationInputs) (result iothub.NameAvailabilityInfo, err error)
	CreateOrUpdate(ctx context.Context, resourceGroupName string, provisioningServiceName string, iotDpsDescription iothub.ProvisioningServiceDescription) (result iothub.IotDpsResourceCreateOrUpdateFuture, err error)
	Delete(ctx context.Context, provisioningServiceName string, resourceGroupName string) (result iothub.IotDpsResourceDeleteFuture, err error)
	Get(ctx context.Context, provisioningServiceName string, resourceGroupName string) (result iothub.ProvisioningServiceDescription, err error)
	GetOperationResult(ctx context.Context, operationID string, resourceGroupName string, provisioningServiceName string, asyncinfo string) (result iothub.AsyncOperationResult, err error)
	ListByResourceGroup(ctx context.Context, resourceGroupName string) (result iothub.ProvisioningServiceDescriptionListResultPage, err error)
	ListByResourceGroupComplete(ctx context.Context, resourceGroupName string) (result iothub.ProvisioningServiceDescriptionListResultIterator, err error)
	ListBySubscription(ctx context.Context) (result iothub.ProvisioningServiceDescriptionListResultPage, err error)
	ListBySubscriptionComplete(ctx context.Context) (result iothub.ProvisioningServiceDescriptionListResultIterator, err error)
	ListKeys(ctx context.Context, provisioningServiceName string, resourceGroupName string) (result iothub.SharedAccessSignatureAuthorizationRuleListResultPage, err error)
	ListKeysComplete(ctx context.Context, provisioningServiceName string, resourceGroupName string) (result iothub.SharedAccessSignatureAuthorizationRuleListResultIterator, err error)
	ListKeysForKeyName(ctx context.Context, provisioningServiceName string, keyName string, resourceGroupName string) (result iothub.SharedAccessSignatureAuthorizationRuleAccessRightsDescription, err error)
	ListValidSkus(ctx context.Context, provisioningServiceName string, resourceGroupName string) (result iothub.IotDpsSkuDefinitionListResultPage, err error)
	ListValidSkusComplete(ctx context.Context, provisioningServiceName string, resourceGroupName string) (result iothub.IotDpsSkuDefinitionListResultIterator, err error)
	Update(ctx context.Context, resourceGroupName string, provisioningServiceName string, provisioningServiceTags iothub.TagsResource) (result iothub.IotDpsResourceUpdateFuture, err error)
}

var _ IotDpsResourceClientAPI = (*iothub.IotDpsResourceClient)(nil)
