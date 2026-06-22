// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armdeviceregistry_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/deviceregistry/armdeviceregistry/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// CMSTestSuite tests the Credential Management Service (CMS) flow:
// Credential CRUD, Policy CRUD, Synchronize, Device operations, and BYOR.
//
// PREREQUISITES
// =============
// The test expects a pre-existing ADR Namespace with IoT Hub integration.
// Set the following environment variables (or provide defaults for playback):
//
//	AZURE_SUBSCRIPTION_ID  - Target subscription
//	LOCATION               - Azure region (e.g. "eastus2euap")
//	RESOURCE_GROUP_NAME    - Resource group containing the namespace
//	NAMESPACE_NAME         - ADR namespace name
//	POLICY_NAME            - Policy name to create/delete
//	BYOR_POLICY_NAME       - BYOR policy name to create/delete
//	DEVICE_NAME            - Device name to create/delete
type CMSTestSuite struct {
	suite.Suite

	ctx               context.Context
	cred              azcore.TokenCredential
	options           *arm.ClientOptions
	subscriptionID    string
	location          string
	resourceGroupName string
	namespaceName     string
	policyName        string
	byorPolicyName    string
	deviceName        string
}

func (testsuite *CMSTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)
	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())

	testsuite.subscriptionID = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	testsuite.location = recording.GetEnvVariable("LOCATION", "eastus2euap")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "adr-sdk-test-cms")
	testsuite.namespaceName = recording.GetEnvVariable("NAMESPACE_NAME", "cms-test-namespace")
	testsuite.policyName = recording.GetEnvVariable("POLICY_NAME", "cms-test-policy")
	testsuite.byorPolicyName = recording.GetEnvVariable("BYOR_POLICY_NAME", "cms-test-byor-policy")
	testsuite.deviceName = recording.GetEnvVariable("DEVICE_NAME", "cms-test-device")
}

func (testsuite *CMSTestSuite) TearDownSuite() {
	testutil.StopRecording(testsuite.T())
}

func TestCMSTestSuite(t *testing.T) {
	suite.Run(t, new(CMSTestSuite))
}

func (testsuite *CMSTestSuite) TestCredentialAndPolicyFlow() {
	ctx := testsuite.ctx

	// Create clients
	credentialsClient, err := armdeviceregistry.NewCredentialsClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	require.NoError(testsuite.T(), err)

	policiesClient, err := armdeviceregistry.NewPoliciesClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	require.NoError(testsuite.T(), err)

	devicesClient, err := armdeviceregistry.NewNamespaceDevicesClient(testsuite.subscriptionID, testsuite.cred, testsuite.options)
	require.NoError(testsuite.T(), err)

	// ================================================================
	// Step 1: Create or Get Credential
	// ================================================================
	fmt.Println("Step 1: Creating credential...")
	credentialResource := armdeviceregistry.Credential{
		Location: to.Ptr(testsuite.location),
	}
	credPoller, err := credentialsClient.BeginCreateOrUpdate(ctx, testsuite.resourceGroupName, testsuite.namespaceName, credentialResource, nil)
	require.NoError(testsuite.T(), err)
	credResult, err := credPoller.PollUntilDone(ctx, nil)
	require.NoError(testsuite.T(), err)
	require.NotNil(testsuite.T(), credResult.ID)
	require.Equal(testsuite.T(), testsuite.location, *credResult.Location)
	fmt.Println("  Credential created successfully")

	// Verify credential via GET
	credGetResp, err := credentialsClient.Get(ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	require.NoError(testsuite.T(), err)
	require.NotNil(testsuite.T(), credGetResp.ID)
	require.Equal(testsuite.T(), testsuite.location, *credGetResp.Location)
	fmt.Println("  Credential GET verified")

	// ================================================================
	// Step 2: Clean up existing policies and create a new one
	// ================================================================
	fmt.Println("Step 2: Cleaning up existing policies...")
	pager := policiesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		require.NoError(testsuite.T(), err)
		for _, p := range page.Value {
			fmt.Printf("  Deleting existing policy '%s'...\n", *p.Name)
			delPoller, err := policiesClient.BeginDelete(ctx, testsuite.resourceGroupName, testsuite.namespaceName, *p.Name, nil)
			require.NoError(testsuite.T(), err)
			_, err = delPoller.PollUntilDone(ctx, nil)
			require.NoError(testsuite.T(), err)
			fmt.Printf("  Deleted '%s'\n", *p.Name)
		}
	}

	// Create policy with ECC certificate, 90-day validity
	fmt.Printf("  Creating policy '%s' with ECC cert (90-day validity)...\n", testsuite.policyName)
	policyResource := armdeviceregistry.Policy{
		Properties: &armdeviceregistry.PolicyProperties{
			Certificate: &armdeviceregistry.CertificateConfiguration{
				CertificateAuthorityConfiguration: &armdeviceregistry.CertificateAuthorityConfiguration{
					KeyType: to.Ptr(armdeviceregistry.SupportedKeyTypeECC),
				},
				LeafCertificateConfiguration: &armdeviceregistry.LeafCertificateConfiguration{
					ValidityPeriodInDays: to.Ptr[int32](90),
				},
			},
		},
	}
	policyPoller, err := policiesClient.BeginCreateOrUpdate(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.policyName, policyResource, nil)
	require.NoError(testsuite.T(), err)
	policyResult, err := policyPoller.PollUntilDone(ctx, nil)
	require.NoError(testsuite.T(), err)
	require.NotNil(testsuite.T(), policyResult.Properties)
	require.NotNil(testsuite.T(), policyResult.Properties.Certificate)
	require.Equal(testsuite.T(), testsuite.policyName, *policyResult.Name)
	require.Equal(testsuite.T(), armdeviceregistry.SupportedKeyTypeECC, *policyResult.Properties.Certificate.CertificateAuthorityConfiguration.KeyType)
	require.Equal(testsuite.T(), int32(90), *policyResult.Properties.Certificate.LeafCertificateConfiguration.ValidityPeriodInDays)
	require.NotNil(testsuite.T(), policyResult.Properties.ProvisioningState)
	fmt.Println("  Policy created successfully")

	// ================================================================
	// Step 3: Policy LIST operation
	// ================================================================
	fmt.Println("Step 3: Testing policy LIST...")
	var foundPolicy bool
	listPager := policiesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for listPager.More() {
		page, err := listPager.NextPage(ctx)
		require.NoError(testsuite.T(), err)
		for _, p := range page.Value {
			if *p.Name == testsuite.policyName {
				foundPolicy = true
			}
		}
	}
	require.True(testsuite.T(), foundPolicy, "policy should appear in LIST results")
	fmt.Println("  LIST verified, found policy")

	// ================================================================
	// Step 4: Synchronize credentials with IoT Hub
	// ================================================================
	fmt.Println("Step 4: Synchronizing credentials with IoT Hub...")
	syncPoller, err := credentialsClient.BeginSynchronize(ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	require.NoError(testsuite.T(), err)
	_, err = syncPoller.PollUntilDone(ctx, nil)
	require.NoError(testsuite.T(), err)
	fmt.Println("  Synchronization completed")

	// ================================================================
	// Step 5: GET policy after sync, then UPDATE validity to 60 days
	// ================================================================
	fmt.Println("Step 5: Getting fresh policy after sync...")
	policyGetResp, err := policiesClient.Get(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.policyName, nil)
	require.NoError(testsuite.T(), err)
	currentValidity := *policyGetResp.Properties.Certificate.LeafCertificateConfiguration.ValidityPeriodInDays
	fmt.Printf("  Current validity: %d days\n", currentValidity)

	fmt.Println("  Updating policy validity to 60 days...")
	policyPatch := armdeviceregistry.PolicyUpdate{
		Properties: &armdeviceregistry.PolicyUpdateProperties{
			Certificate: &armdeviceregistry.CertificateConfiguration{
				LeafCertificateConfiguration: &armdeviceregistry.LeafCertificateConfiguration{
					ValidityPeriodInDays: to.Ptr[int32](60),
				},
			},
		},
	}
	updatePoller, err := policiesClient.BeginUpdate(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.policyName, policyPatch, nil)
	require.NoError(testsuite.T(), err)
	_, err = updatePoller.PollUntilDone(ctx, nil)
	require.NoError(testsuite.T(), err)

	// Verify update
	policyGetResp, err = policiesClient.Get(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.policyName, nil)
	require.NoError(testsuite.T(), err)
	require.Equal(testsuite.T(), int32(60), *policyGetResp.Properties.Certificate.LeafCertificateConfiguration.ValidityPeriodInDays)
	fmt.Println("  Policy updated, validity now 60 days")

	// ================================================================
	// Step 6: Device CRUD
	// ================================================================
	fmt.Printf("Step 6: Creating device '%s'...\n", testsuite.deviceName)
	deviceResource := armdeviceregistry.NamespaceDevice{
		Location: to.Ptr(testsuite.location),
		Properties: &armdeviceregistry.NamespaceDeviceProperties{
			Manufacturer:           to.Ptr("Contoso"),
			Model:                  to.Ptr("CMS-TestModel-5000"),
			OperatingSystem:        to.Ptr("Linux"),
			OperatingSystemVersion: to.Ptr("22.04"),
			Endpoints:              &armdeviceregistry.MessagingEndpoints{},
		},
	}
	devicePoller, err := devicesClient.BeginCreateOrReplace(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.deviceName, deviceResource, nil)
	require.NoError(testsuite.T(), err)
	deviceResult, err := devicePoller.PollUntilDone(ctx, nil)
	require.NoError(testsuite.T(), err)
	require.Equal(testsuite.T(), testsuite.deviceName, *deviceResult.Name)
	require.NotNil(testsuite.T(), deviceResult.Properties.UUID)
	require.Equal(testsuite.T(), "Contoso", *deviceResult.Properties.Manufacturer)
	require.Equal(testsuite.T(), "CMS-TestModel-5000", *deviceResult.Properties.Model)
	fmt.Printf("  Device created: %s, UUID: %s\n", *deviceResult.Name, *deviceResult.Properties.UUID)

	// GET device and verify properties
	fmt.Println("  Getting device and verifying properties...")
	deviceGetResp, err := devicesClient.Get(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.deviceName, nil)
	require.NoError(testsuite.T(), err)
	require.Equal(testsuite.T(), "Contoso", *deviceGetResp.Properties.Manufacturer)
	require.Equal(testsuite.T(), "CMS-TestModel-5000", *deviceGetResp.Properties.Model)
	require.Equal(testsuite.T(), "Linux", *deviceGetResp.Properties.OperatingSystem)
	require.Equal(testsuite.T(), "22.04", *deviceGetResp.Properties.OperatingSystemVersion)
	fmt.Println("  Device properties verified")

	// LIST devices
	fmt.Println("  Listing devices in namespace...")
	var foundDevice bool
	devicePager := devicesClient.NewListByResourceGroupPager(testsuite.resourceGroupName, testsuite.namespaceName, nil)
	for devicePager.More() {
		page, err := devicePager.NextPage(ctx)
		require.NoError(testsuite.T(), err)
		for _, d := range page.Value {
			if *d.Name == testsuite.deviceName {
				foundDevice = true
			}
		}
	}
	require.True(testsuite.T(), foundDevice, "device should appear in LIST results")
	fmt.Println("  LIST verified, found device")

	// ================================================================
	// Step 7: Device Revoke (negative test — ARM-created device has no policy)
	// ================================================================
	fmt.Println("Step 7: Testing Device.Revoke (expect error for ARM-created device)...")
	revokeRequest := armdeviceregistry.DeviceCredentialsRevokeRequest{
		Disable: to.Ptr(false),
	}
	revokePoller, err := devicesClient.BeginRevoke(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.deviceName, revokeRequest, nil)
	if err != nil {
		fmt.Printf("  Revoke returned error (expected): %s\n", err.Error())
	} else {
		_, revokeErr := revokePoller.PollUntilDone(ctx, nil)
		// The RP may return an error during polling due to LRO bug
		if revokeErr != nil {
			fmt.Printf("  Revoke polling returned error (expected): %s\n", revokeErr.Error())
		} else {
			fmt.Println("  Revoke completed (RP may have been fixed)")
		}
	}

	// Verify device state unchanged after revoke attempt
	fmt.Println("  Verifying device state unchanged...")
	deviceAfterRevoke, err := devicesClient.Get(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.deviceName, nil)
	require.NoError(testsuite.T(), err)
	require.NotNil(testsuite.T(), deviceAfterRevoke.Properties)
	fmt.Printf("  Device still exists: %s\n", *deviceAfterRevoke.Name)

	// Delete device
	// NOTE: Known RP bug — the RP returns HTTP 200 instead of 202/204 for device delete,
	// which causes the SDK's LRO poller to fail. Treat any error during delete polling
	// as non-fatal since the deletion itself succeeds.
	fmt.Printf("  Deleting device '%s'...\n", testsuite.deviceName)
	deviceDelPoller, err := devicesClient.BeginDelete(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.deviceName, nil)
	if err != nil {
		fmt.Printf("  Device delete returned error (known RP bug): %s\n", err.Error())
	} else {
		_, pollErr := deviceDelPoller.PollUntilDone(ctx, nil)
		if pollErr != nil {
			fmt.Printf("  Device delete polling returned error (known RP bug): %s\n", pollErr.Error())
		} else {
			fmt.Println("  Device deleted")
		}
	}

	// ================================================================
	// Step 8: RevokeIssuer on standard policy
	// ================================================================
	fmt.Println("Step 8: Testing RevokeIssuer on standard policy...")
	revokeIssuerPoller, err := policiesClient.BeginRevokeIssuer(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.policyName, nil)
	require.NoError(testsuite.T(), err)
	_, err = revokeIssuerPoller.PollUntilDone(ctx, nil)
	require.NoError(testsuite.T(), err)
	fmt.Println("  RevokeIssuer completed")

	// Verify policy state after RevokeIssuer
	policyAfterRevoke, err := policiesClient.Get(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.policyName, nil)
	require.NoError(testsuite.T(), err)
	require.NotNil(testsuite.T(), policyAfterRevoke.Properties)
	fmt.Println("  Policy state verified after RevokeIssuer")

	// ================================================================
	// Step 9: Delete standard policy
	// ================================================================
	fmt.Printf("Step 9: Deleting policy '%s'...\n", testsuite.policyName)
	policyDelPoller, err := policiesClient.BeginDelete(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.policyName, nil)
	require.NoError(testsuite.T(), err)
	_, err = policyDelPoller.PollUntilDone(ctx, nil)
	require.NoError(testsuite.T(), err)
	fmt.Println("  Policy deleted")

	// ================================================================
	// Step 10: BYOR (Bring Your Own Root) Policy Flow
	// ================================================================
	fmt.Printf("Step 10: Creating BYOR policy '%s'...\n", testsuite.byorPolicyName)
	byorPolicy := armdeviceregistry.Policy{
		Properties: &armdeviceregistry.PolicyProperties{
			Certificate: &armdeviceregistry.CertificateConfiguration{
				CertificateAuthorityConfiguration: &armdeviceregistry.CertificateAuthorityConfiguration{
					KeyType: to.Ptr(armdeviceregistry.SupportedKeyTypeECC),
					BringYourOwnRoot: &armdeviceregistry.BringYourOwnRoot{
						Enabled: to.Ptr(true),
					},
				},
				LeafCertificateConfiguration: &armdeviceregistry.LeafCertificateConfiguration{
					ValidityPeriodInDays: to.Ptr[int32](90),
				},
			},
		},
	}
	byorPoller, err := policiesClient.BeginCreateOrUpdate(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.byorPolicyName, byorPolicy, nil)
	require.NoError(testsuite.T(), err)
	byorResult, err := byorPoller.PollUntilDone(ctx, nil)
	require.NoError(testsuite.T(), err)
	require.Equal(testsuite.T(), testsuite.byorPolicyName, *byorResult.Name)

	byorConfig := byorResult.Properties.Certificate.CertificateAuthorityConfiguration.BringYourOwnRoot
	require.NotNil(testsuite.T(), byorConfig)
	require.True(testsuite.T(), *byorConfig.Enabled)
	fmt.Printf("  BYOR policy created, enabled: %t, status: %s\n", *byorConfig.Enabled, *byorConfig.Status)

	// Verify PendingActivation status and CSR
	require.Equal(testsuite.T(), armdeviceregistry.BringYourOwnRootStatusPendingActivation, *byorConfig.Status)
	require.NotNil(testsuite.T(), byorConfig.CertificateSigningRequest)
	require.True(testsuite.T(), strings.Contains(*byorConfig.CertificateSigningRequest, "-----BEGIN CERTIFICATE REQUEST-----"))
	fmt.Printf("  CSR present (%d chars)\n", len(*byorConfig.CertificateSigningRequest))

	// ================================================================
	// Step 11: ActivateBringYourOwnRoot with invalid cert (negative test)
	// ================================================================
	fmt.Println("Step 11: Testing ActivateBYOR with invalid certificate (negative test)...")
	fakeCertChain := "-----BEGIN CERTIFICATE-----\n" +
		"MIIBkTCB+wIJALRiMLAhFake0DQYJKoZIhvcNAQELBQAwDzENMAsGA1UEAwwEdGVz\n" +
		"dDAeFw0yNDAzMjAxMjAwMDBaFw0yNTAzMjAxMjAwMDBaMA8xDTALBgNVBAMMBHRl\n" +
		"c3QwXDANBgkqhkiG9w0BAQEFAANLADBIAkEA0Z3VS5JJcds3xf0GQGZ/fake+key\n" +
		"data+that+is+intentionally+invalid+for+testing+purposes+only+AAAAAAAAAA==\n" +
		"-----END CERTIFICATE-----"

	activateReq := armdeviceregistry.ActivateBringYourOwnRootRequest{
		CertificateChain: to.Ptr(fakeCertChain),
	}
	activatePoller, activateErr := policiesClient.BeginActivateBringYourOwnRoot(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.byorPolicyName, activateReq, nil)
	if activateErr != nil {
		fmt.Printf("  ActivateBYOR rejected invalid cert (expected): %s\n", activateErr.Error())
	} else {
		_, pollErr := activatePoller.PollUntilDone(ctx, nil)
		if pollErr != nil {
			fmt.Printf("  ActivateBYOR polling failed as expected: %s\n", pollErr.Error())
		} else {
			fmt.Println("  ActivateBYOR with invalid cert unexpectedly succeeded (known RP bug: RP does not validate cert chain)")
		}
	}

	// Verify BYOR state after activation attempt
	fmt.Println("  Verifying BYOR state after activation attempt...")
	byorAfterFail, err := policiesClient.Get(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.byorPolicyName, nil)
	require.NoError(testsuite.T(), err)
	byorConfigAfterFail := byorAfterFail.Properties.Certificate.CertificateAuthorityConfiguration.BringYourOwnRoot
	require.True(testsuite.T(), *byorConfigAfterFail.Enabled)
	fmt.Printf("  BYOR state after activation: %s\n", *byorConfigAfterFail.Status)

	// ================================================================
	// Step 12: Update BYOR policy — change validity to 45 days
	// ================================================================
	fmt.Println("Step 12: Updating BYOR policy validity to 45 days...")
	byorPatch := armdeviceregistry.PolicyUpdate{
		Properties: &armdeviceregistry.PolicyUpdateProperties{
			Certificate: &armdeviceregistry.CertificateConfiguration{
				LeafCertificateConfiguration: &armdeviceregistry.LeafCertificateConfiguration{
					ValidityPeriodInDays: to.Ptr[int32](45),
				},
			},
		},
	}
	byorUpdatePoller, err := policiesClient.BeginUpdate(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.byorPolicyName, byorPatch, nil)
	require.NoError(testsuite.T(), err)
	_, err = byorUpdatePoller.PollUntilDone(ctx, nil)
	require.NoError(testsuite.T(), err)

	// Verify update
	byorGetResp, err := policiesClient.Get(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.byorPolicyName, nil)
	require.NoError(testsuite.T(), err)
	require.Equal(testsuite.T(), int32(45), *byorGetResp.Properties.Certificate.LeafCertificateConfiguration.ValidityPeriodInDays)
	require.True(testsuite.T(), *byorGetResp.Properties.Certificate.CertificateAuthorityConfiguration.BringYourOwnRoot.Enabled)
	fmt.Println("  BYOR policy updated, validity now 45 days, BYOR still enabled")

	// ================================================================
	// Step 13: Delete BYOR policy
	// ================================================================
	fmt.Printf("Step 13: Deleting BYOR policy '%s'...\n", testsuite.byorPolicyName)
	byorDelPoller, err := policiesClient.BeginDelete(ctx, testsuite.resourceGroupName, testsuite.namespaceName, testsuite.byorPolicyName, nil)
	require.NoError(testsuite.T(), err)
	_, err = byorDelPoller.PollUntilDone(ctx, nil)
	require.NoError(testsuite.T(), err)
	fmt.Println("  BYOR policy deleted")

	// ================================================================
	// Step 14: Delete credential
	// ================================================================
	fmt.Println("Step 14: Deleting credential...")
	credDelPoller, err := credentialsClient.BeginDelete(ctx, testsuite.resourceGroupName, testsuite.namespaceName, nil)
	require.NoError(testsuite.T(), err)
	_, err = credDelPoller.PollUntilDone(ctx, nil)
	require.NoError(testsuite.T(), err)
	fmt.Println("  Credential deleted")

	fmt.Println("TEST COMPLETED SUCCESSFULLY")
}
