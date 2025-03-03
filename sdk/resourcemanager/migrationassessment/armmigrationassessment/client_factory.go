// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmigrationassessment

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
)

// ClientFactory is a client factory used to create any client in this module.
// Don't use this type directly, use NewClientFactory instead.
type ClientFactory struct {
	subscriptionID string
	internal       *arm.Client
}

// NewClientFactory creates a new instance of ClientFactory with the specified values.
// The parameter values will be propagated to any client created from this factory.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewClientFactory(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ClientFactory, error) {
	internal, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	return &ClientFactory{
		subscriptionID: subscriptionID,
		internal:       internal,
	}, nil
}

// NewAksAssessmentOperationsClient creates a new instance of AksAssessmentOperationsClient.
func (c *ClientFactory) NewAksAssessmentOperationsClient() *AksAssessmentOperationsClient {
	return &AksAssessmentOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAksClusterOperationsClient creates a new instance of AksClusterOperationsClient.
func (c *ClientFactory) NewAksClusterOperationsClient() *AksClusterOperationsClient {
	return &AksClusterOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAksCostDetailOperationsClient creates a new instance of AksCostDetailOperationsClient.
func (c *ClientFactory) NewAksCostDetailOperationsClient() *AksCostDetailOperationsClient {
	return &AksCostDetailOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAksOptionsOperationsClient creates a new instance of AksOptionsOperationsClient.
func (c *ClientFactory) NewAksOptionsOperationsClient() *AksOptionsOperationsClient {
	return &AksOptionsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAksSummaryOperationsClient creates a new instance of AksSummaryOperationsClient.
func (c *ClientFactory) NewAksSummaryOperationsClient() *AksSummaryOperationsClient {
	return &AksSummaryOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAssessedMachinesOperationsClient creates a new instance of AssessedMachinesOperationsClient.
func (c *ClientFactory) NewAssessedMachinesOperationsClient() *AssessedMachinesOperationsClient {
	return &AssessedMachinesOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAssessedSQLDatabaseV2OperationsClient creates a new instance of AssessedSQLDatabaseV2OperationsClient.
func (c *ClientFactory) NewAssessedSQLDatabaseV2OperationsClient() *AssessedSQLDatabaseV2OperationsClient {
	return &AssessedSQLDatabaseV2OperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAssessedSQLInstanceV2OperationsClient creates a new instance of AssessedSQLInstanceV2OperationsClient.
func (c *ClientFactory) NewAssessedSQLInstanceV2OperationsClient() *AssessedSQLInstanceV2OperationsClient {
	return &AssessedSQLInstanceV2OperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAssessedSQLMachinesOperationsClient creates a new instance of AssessedSQLMachinesOperationsClient.
func (c *ClientFactory) NewAssessedSQLMachinesOperationsClient() *AssessedSQLMachinesOperationsClient {
	return &AssessedSQLMachinesOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAssessedSQLRecommendedEntityOperationsClient creates a new instance of AssessedSQLRecommendedEntityOperationsClient.
func (c *ClientFactory) NewAssessedSQLRecommendedEntityOperationsClient() *AssessedSQLRecommendedEntityOperationsClient {
	return &AssessedSQLRecommendedEntityOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAssessedWebAppV2OperationsClient creates a new instance of AssessedWebAppV2OperationsClient.
func (c *ClientFactory) NewAssessedWebAppV2OperationsClient() *AssessedWebAppV2OperationsClient {
	return &AssessedWebAppV2OperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAssessedWebApplicationOperationsClient creates a new instance of AssessedWebApplicationOperationsClient.
func (c *ClientFactory) NewAssessedWebApplicationOperationsClient() *AssessedWebApplicationOperationsClient {
	return &AssessedWebApplicationOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAssessmentOptionsOperationsClient creates a new instance of AssessmentOptionsOperationsClient.
func (c *ClientFactory) NewAssessmentOptionsOperationsClient() *AssessmentOptionsOperationsClient {
	return &AssessmentOptionsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAssessmentProjectSummaryOperationsClient creates a new instance of AssessmentProjectSummaryOperationsClient.
func (c *ClientFactory) NewAssessmentProjectSummaryOperationsClient() *AssessmentProjectSummaryOperationsClient {
	return &AssessmentProjectSummaryOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAssessmentProjectsOperationsClient creates a new instance of AssessmentProjectsOperationsClient.
func (c *ClientFactory) NewAssessmentProjectsOperationsClient() *AssessmentProjectsOperationsClient {
	return &AssessmentProjectsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAssessmentsOperationsClient creates a new instance of AssessmentsOperationsClient.
func (c *ClientFactory) NewAssessmentsOperationsClient() *AssessmentsOperationsClient {
	return &AssessmentsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAvsAssessedMachinesOperationsClient creates a new instance of AvsAssessedMachinesOperationsClient.
func (c *ClientFactory) NewAvsAssessedMachinesOperationsClient() *AvsAssessedMachinesOperationsClient {
	return &AvsAssessedMachinesOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAvsAssessmentOptionsOperationsClient creates a new instance of AvsAssessmentOptionsOperationsClient.
func (c *ClientFactory) NewAvsAssessmentOptionsOperationsClient() *AvsAssessmentOptionsOperationsClient {
	return &AvsAssessmentOptionsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewAvsAssessmentsOperationsClient creates a new instance of AvsAssessmentsOperationsClient.
func (c *ClientFactory) NewAvsAssessmentsOperationsClient() *AvsAssessmentsOperationsClient {
	return &AvsAssessmentsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewBusinessCaseAvsSummaryOperationsClient creates a new instance of BusinessCaseAvsSummaryOperationsClient.
func (c *ClientFactory) NewBusinessCaseAvsSummaryOperationsClient() *BusinessCaseAvsSummaryOperationsClient {
	return &BusinessCaseAvsSummaryOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewBusinessCaseIaasSummaryOperationsClient creates a new instance of BusinessCaseIaasSummaryOperationsClient.
func (c *ClientFactory) NewBusinessCaseIaasSummaryOperationsClient() *BusinessCaseIaasSummaryOperationsClient {
	return &BusinessCaseIaasSummaryOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewBusinessCaseOperationsClient creates a new instance of BusinessCaseOperationsClient.
func (c *ClientFactory) NewBusinessCaseOperationsClient() *BusinessCaseOperationsClient {
	return &BusinessCaseOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewBusinessCaseOverviewSummaryOperationsClient creates a new instance of BusinessCaseOverviewSummaryOperationsClient.
func (c *ClientFactory) NewBusinessCaseOverviewSummaryOperationsClient() *BusinessCaseOverviewSummaryOperationsClient {
	return &BusinessCaseOverviewSummaryOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewBusinessCasePaasSummaryOperationsClient creates a new instance of BusinessCasePaasSummaryOperationsClient.
func (c *ClientFactory) NewBusinessCasePaasSummaryOperationsClient() *BusinessCasePaasSummaryOperationsClient {
	return &BusinessCasePaasSummaryOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewEvaluatedAvsMachinesOperationsClient creates a new instance of EvaluatedAvsMachinesOperationsClient.
func (c *ClientFactory) NewEvaluatedAvsMachinesOperationsClient() *EvaluatedAvsMachinesOperationsClient {
	return &EvaluatedAvsMachinesOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewEvaluatedMachinesOperationsClient creates a new instance of EvaluatedMachinesOperationsClient.
func (c *ClientFactory) NewEvaluatedMachinesOperationsClient() *EvaluatedMachinesOperationsClient {
	return &EvaluatedMachinesOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewEvaluatedSQLEntitiesOperationsClient creates a new instance of EvaluatedSQLEntitiesOperationsClient.
func (c *ClientFactory) NewEvaluatedSQLEntitiesOperationsClient() *EvaluatedSQLEntitiesOperationsClient {
	return &EvaluatedSQLEntitiesOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewEvaluatedWebAppsOperationsClient creates a new instance of EvaluatedWebAppsOperationsClient.
func (c *ClientFactory) NewEvaluatedWebAppsOperationsClient() *EvaluatedWebAppsOperationsClient {
	return &EvaluatedWebAppsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewGroupsOperationsClient creates a new instance of GroupsOperationsClient.
func (c *ClientFactory) NewGroupsOperationsClient() *GroupsOperationsClient {
	return &GroupsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewHypervCollectorsOperationsClient creates a new instance of HypervCollectorsOperationsClient.
func (c *ClientFactory) NewHypervCollectorsOperationsClient() *HypervCollectorsOperationsClient {
	return &HypervCollectorsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewImportCollectorsOperationsClient creates a new instance of ImportCollectorsOperationsClient.
func (c *ClientFactory) NewImportCollectorsOperationsClient() *ImportCollectorsOperationsClient {
	return &ImportCollectorsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewMachinesOperationsClient creates a new instance of MachinesOperationsClient.
func (c *ClientFactory) NewMachinesOperationsClient() *MachinesOperationsClient {
	return &MachinesOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewOperationsClient creates a new instance of OperationsClient.
func (c *ClientFactory) NewOperationsClient() *OperationsClient {
	return &OperationsClient{
		internal: c.internal,
	}
}

// NewPrivateEndpointConnectionOperationsClient creates a new instance of PrivateEndpointConnectionOperationsClient.
func (c *ClientFactory) NewPrivateEndpointConnectionOperationsClient() *PrivateEndpointConnectionOperationsClient {
	return &PrivateEndpointConnectionOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewPrivateLinkResourceOperationsClient creates a new instance of PrivateLinkResourceOperationsClient.
func (c *ClientFactory) NewPrivateLinkResourceOperationsClient() *PrivateLinkResourceOperationsClient {
	return &PrivateLinkResourceOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewSQLAssessmentOptionsOperationsClient creates a new instance of SQLAssessmentOptionsOperationsClient.
func (c *ClientFactory) NewSQLAssessmentOptionsOperationsClient() *SQLAssessmentOptionsOperationsClient {
	return &SQLAssessmentOptionsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewSQLAssessmentV2OperationsClient creates a new instance of SQLAssessmentV2OperationsClient.
func (c *ClientFactory) NewSQLAssessmentV2OperationsClient() *SQLAssessmentV2OperationsClient {
	return &SQLAssessmentV2OperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewSQLAssessmentV2SummaryOperationsClient creates a new instance of SQLAssessmentV2SummaryOperationsClient.
func (c *ClientFactory) NewSQLAssessmentV2SummaryOperationsClient() *SQLAssessmentV2SummaryOperationsClient {
	return &SQLAssessmentV2SummaryOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewSQLCollectorOperationsClient creates a new instance of SQLCollectorOperationsClient.
func (c *ClientFactory) NewSQLCollectorOperationsClient() *SQLCollectorOperationsClient {
	return &SQLCollectorOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewServerCollectorsOperationsClient creates a new instance of ServerCollectorsOperationsClient.
func (c *ClientFactory) NewServerCollectorsOperationsClient() *ServerCollectorsOperationsClient {
	return &ServerCollectorsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewVmwareCollectorsOperationsClient creates a new instance of VmwareCollectorsOperationsClient.
func (c *ClientFactory) NewVmwareCollectorsOperationsClient() *VmwareCollectorsOperationsClient {
	return &VmwareCollectorsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewWebAppAssessmentOptionsOperationsClient creates a new instance of WebAppAssessmentOptionsOperationsClient.
func (c *ClientFactory) NewWebAppAssessmentOptionsOperationsClient() *WebAppAssessmentOptionsOperationsClient {
	return &WebAppAssessmentOptionsOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewWebAppAssessmentV2OperationsClient creates a new instance of WebAppAssessmentV2OperationsClient.
func (c *ClientFactory) NewWebAppAssessmentV2OperationsClient() *WebAppAssessmentV2OperationsClient {
	return &WebAppAssessmentV2OperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewWebAppAssessmentV2SummaryOperationsClient creates a new instance of WebAppAssessmentV2SummaryOperationsClient.
func (c *ClientFactory) NewWebAppAssessmentV2SummaryOperationsClient() *WebAppAssessmentV2SummaryOperationsClient {
	return &WebAppAssessmentV2SummaryOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewWebAppCollectorOperationsClient creates a new instance of WebAppCollectorOperationsClient.
func (c *ClientFactory) NewWebAppCollectorOperationsClient() *WebAppCollectorOperationsClient {
	return &WebAppCollectorOperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}

// NewWebAppServicePlanV2OperationsClient creates a new instance of WebAppServicePlanV2OperationsClient.
func (c *ClientFactory) NewWebAppServicePlanV2OperationsClient() *WebAppServicePlanV2OperationsClient {
	return &WebAppServicePlanV2OperationsClient{
		subscriptionID: c.subscriptionID,
		internal:       c.internal,
	}
}
