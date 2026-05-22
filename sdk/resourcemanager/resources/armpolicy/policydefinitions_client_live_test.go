// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armpolicy_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armpolicy"
	"github.com/stretchr/testify/suite"
)

type PolicyDefinitionsClientTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	location          string
	resourceGroupName string
	subscriptionID    string
}

func (testsuite *PolicyDefinitionsClientTestSuite) SetupSuite() {
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus")
	testsuite.subscriptionID = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testutil.StartRecording(testsuite.T(), pathToPackage)
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name

}

func (testsuite *PolicyDefinitionsClientTestSuite) TearDownSuite() {
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionID, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestPolicyDefinitionsClient(t *testing.T) {
	suite.Run(t, new(PolicyDefinitionsClientTestSuite))
}

func (testsuite *PolicyDefinitionsClientTestSuite) TestPolicyDefinitionsCRUD() {
	// create policy definition
	fmt.Println("Call operation: PolicyDefinitions_CreateOrUpdate")
	policyDefinitionName := "go-test-definition"
	policyDefinitionsClient, err := armpolicy.NewDefinitionsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = policyDefinitionsClient.CreateOrUpdate(
		testsuite.ctx,
		policyDefinitionName,
		armpolicy.Definition{
			Properties: &armpolicy.DefinitionProperties{
				PolicyType:  to.Ptr(armpolicy.PolicyTypeCustom),
				Description: to.Ptr("test case"),
				Parameters: map[string]*armpolicy.ParameterDefinitionsValue{
					"prefix": {
						Type: to.Ptr(armpolicy.ParameterTypeString),
						Metadata: &armpolicy.ParameterDefinitionsValueMetadata{
							Description: to.Ptr("prefix description"),
							DisplayName: to.Ptr("test case prefix"),
						},
					},
					"suffix": {
						Type: to.Ptr(armpolicy.ParameterTypeString),
						Metadata: &armpolicy.ParameterDefinitionsValueMetadata{
							Description: to.Ptr("suffix description"),
							DisplayName: to.Ptr("test case suffix"),
						},
					},
				},
				PolicyRule: map[string]interface{}{
					"if": map[string]interface{}{
						"not": map[string]interface{}{
							"field": "name",
							"like":  "[concat(parameters('prefix'), '*', parameters('suffix'))]",
						},
					},
					"then": map[string]interface{}{
						"effect": "deny",
					},
				},
			},
		},
		nil,
	)
	testsuite.Require().NoError(err)

	// get policy definition
	fmt.Println("Call operation: PolicyDefinitions_Get")
	_, err = policyDefinitionsClient.Get(testsuite.ctx, policyDefinitionName, nil)
	testsuite.Require().NoError(err)

	// list policy definition
	fmt.Println("Call operation: PolicyDefinitions_List")
	list := policyDefinitionsClient.NewListPager(nil)
	testsuite.Require().True(list.More())

	// list policy definition
	fmt.Println("Call operation: PolicyDefinitions_ListBuiltIn")
	listBuiltIn := policyDefinitionsClient.NewListBuiltInPager(nil)
	testsuite.Require().True(listBuiltIn.More())

	// delete policy definition
	fmt.Println("Call operation: PolicyDefinitions_Delete")
	_, err = policyDefinitionsClient.Delete(testsuite.ctx, policyDefinitionName, nil)
	testsuite.Require().NoError(err)
}
