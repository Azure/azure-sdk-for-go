// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armlogic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/logic/armlogic"
	"github.com/stretchr/testify/suite"
)

type IntegrationAccountsTestSuite struct {
	suite.Suite

	ctx                               context.Context
	cred                              azcore.TokenCredential
	options                           *arm.ClientOptions
	agreementName                     string
	apiName                           string
	assemblyArtifactName              string
	batchConfigurationName            string
	certificateName                   string
	integrationAccountName            string
	integrationServiceEnvironmentName string
	mapName                           string
	partnerName                       string
	schemaName                        string
	sessionName                       string
	location                          string
	resourceGroupName                 string
	subscriptionId                    string
}

func (testsuite *IntegrationAccountsTestSuite) SetupSuite() {
	testutil.StartRecording(testsuite.T(), pathToPackage)

	testsuite.ctx = context.Background()
	testsuite.cred, testsuite.options = testutil.GetCredAndClientOptions(testsuite.T())
	testsuite.agreementName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "agreemen", 14, false)
	testsuite.apiName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "apiname", 13, false)
	testsuite.assemblyArtifactName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "assembly", 14, false)
	testsuite.batchConfigurationName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "batchcon", 14, false)
	testsuite.certificateName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "certific", 14, false)
	testsuite.integrationAccountName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "integrat", 14, false)
	testsuite.integrationServiceEnvironmentName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "integratsenv", 18, false)
	testsuite.mapName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "mapname", 13, false)
	testsuite.partnerName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "partnern", 14, false)
	testsuite.schemaName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "schemana", 14, false)
	testsuite.sessionName, _ = recording.GenerateAlphaNumericID(testsuite.T(), "sessionn", 14, false)
	testsuite.location = recording.GetEnvVariable("LOCATION", "westus")
	testsuite.resourceGroupName = recording.GetEnvVariable("RESOURCE_GROUP_NAME", "scenarioTestTempGroup")
	testsuite.subscriptionId = recording.GetEnvVariable("AZURE_SUBSCRIPTION_ID", "00000000-0000-0000-0000-000000000000")
	resourceGroup, _, err := testutil.CreateResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.location)
	testsuite.Require().NoError(err)
	testsuite.resourceGroupName = *resourceGroup.Name
	testsuite.Prepare()
}

func (testsuite *IntegrationAccountsTestSuite) TearDownSuite() {
	testsuite.Cleanup()
	_, err := testutil.DeleteResourceGroup(testsuite.ctx, testsuite.subscriptionId, testsuite.cred, testsuite.options, testsuite.resourceGroupName)
	testsuite.Require().NoError(err)
	testutil.StopRecording(testsuite.T())
}

func TestIntegrationAccountsTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationAccountsTestSuite))
}

func (testsuite *IntegrationAccountsTestSuite) Prepare() {
	var err error
	// From step IntegrationAccounts_CreateOrUpdate
	fmt.Println("Call operation: IntegrationAccounts_CreateOrUpdate")
	integrationAccountsClient, err := armlogic.NewIntegrationAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = integrationAccountsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, armlogic.IntegrationAccount{
		Location:   to.Ptr(testsuite.location),
		Properties: &armlogic.IntegrationAccountProperties{},
		SKU: &armlogic.IntegrationAccountSKU{
			Name: to.Ptr(armlogic.IntegrationAccountSKUNameStandard),
		},
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Logic/integrationAccounts/{integrationAccountName}
func (testsuite *IntegrationAccountsTestSuite) TestIntegrationAccounts() {
	var err error
	// From step IntegrationAccounts_ListBySubscription
	fmt.Println("Call operation: IntegrationAccounts_ListBySubscription")
	integrationAccountsClient, err := armlogic.NewIntegrationAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	integrationAccountsClientNewListBySubscriptionPager := integrationAccountsClient.NewListBySubscriptionPager(&armlogic.IntegrationAccountsClientListBySubscriptionOptions{Top: nil})
	for integrationAccountsClientNewListBySubscriptionPager.More() {
		_, err := integrationAccountsClientNewListBySubscriptionPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IntegrationAccounts_ListByResourceGroup
	fmt.Println("Call operation: IntegrationAccounts_ListByResourceGroup")
	integrationAccountsClientNewListByResourceGroupPager := integrationAccountsClient.NewListByResourceGroupPager(testsuite.resourceGroupName, &armlogic.IntegrationAccountsClientListByResourceGroupOptions{Top: nil})
	for integrationAccountsClientNewListByResourceGroupPager.More() {
		_, err := integrationAccountsClientNewListByResourceGroupPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IntegrationAccounts_Get
	fmt.Println("Call operation: IntegrationAccounts_Get")
	_, err = integrationAccountsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccounts_Update
	fmt.Println("Call operation: IntegrationAccounts_Update")
	_, err = integrationAccountsClient.Update(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, armlogic.IntegrationAccount{
		Location: to.Ptr(testsuite.location),
		SKU: &armlogic.IntegrationAccountSKU{
			Name: to.Ptr(armlogic.IntegrationAccountSKUNameStandard),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccounts_ListCallbackUrl
	fmt.Println("Call operation: IntegrationAccounts_ListCallbackURL")
	_, err = integrationAccountsClient.ListCallbackURL(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, armlogic.GetCallbackURLParameters{
		KeyType: to.Ptr(armlogic.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccounts_RegenerateAccessKey
	fmt.Println("Call operation: IntegrationAccounts_RegenerateAccessKey")
	_, err = integrationAccountsClient.RegenerateAccessKey(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, armlogic.RegenerateActionParameter{
		KeyType: to.Ptr(armlogic.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Logic/integrationAccounts/{integrationAccountName}/assemblies/{assemblyArtifactName}
func (testsuite *IntegrationAccountsTestSuite) TestIntegrationAccountAssemblies() {
	var err error
	// From step IntegrationAccountAssemblies_CreateOrUpdate
	fmt.Println("Call operation: IntegrationAccountAssemblies_CreateOrUpdate")
	integrationAccountAssembliesClient, err := armlogic.NewIntegrationAccountAssembliesClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = integrationAccountAssembliesClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.assemblyArtifactName, armlogic.AssemblyDefinition{
		Location: to.Ptr(testsuite.location),
		Properties: &armlogic.AssemblyProperties{
			Metadata:     map[string]any{},
			Content:      "Base64 encoded Assembly Content",
			AssemblyName: to.Ptr("System.IdentityModel.Tokens.Jwt"),
			ContentType:  to.Ptr("application/octet-stream"),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountAssemblies_List
	fmt.Println("Call operation: IntegrationAccountAssemblies_List")
	integrationAccountAssembliesClientNewListPager := integrationAccountAssembliesClient.NewListPager(testsuite.resourceGroupName, testsuite.integrationAccountName, nil)
	for integrationAccountAssembliesClientNewListPager.More() {
		_, err := integrationAccountAssembliesClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IntegrationAccountAssemblies_Get
	fmt.Println("Call operation: IntegrationAccountAssemblies_Get")
	_, err = integrationAccountAssembliesClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.assemblyArtifactName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountAssemblies_ListContentCallbackUrl
	fmt.Println("Call operation: IntegrationAccountAssemblies_ListContentCallbackURL")
	_, err = integrationAccountAssembliesClient.ListContentCallbackURL(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.assemblyArtifactName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountAssemblies_Delete
	fmt.Println("Call operation: IntegrationAccountAssemblies_Delete")
	_, err = integrationAccountAssembliesClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.assemblyArtifactName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Logic/integrationAccounts/{integrationAccountName}/agreements/{agreementName}
func (testsuite *IntegrationAccountsTestSuite) TestIntegrationAccountAgreements() {
	var err error
	// From step IntegrationAccountAgreements_CreateOrUpdate
	fmt.Println("Call operation: IntegrationAccountAgreements_CreateOrUpdate")
	integrationAccountAgreementsClient, err := armlogic.NewIntegrationAccountAgreementsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = integrationAccountAgreementsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.agreementName, armlogic.IntegrationAccountAgreement{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"IntegrationAccountAgreement": to.Ptr("<IntegrationAccountAgreementName>"),
		},
		Properties: &armlogic.IntegrationAccountAgreementProperties{
			AgreementType: to.Ptr(armlogic.AgreementTypeAS2),
			Content: &armlogic.AgreementContent{
				AS2: &armlogic.AS2AgreementContent{
					ReceiveAgreement: &armlogic.AS2OneWayAgreement{
						ProtocolSettings: &armlogic.AS2ProtocolSettings{
							AcknowledgementConnectionSettings: &armlogic.AS2AcknowledgementConnectionSettings{
								IgnoreCertificateNameMismatch: to.Ptr(true),
								KeepHTTPConnectionAlive:       to.Ptr(true),
								SupportHTTPStatusCodeContinue: to.Ptr(true),
								UnfoldHTTPHeaders:             to.Ptr(true),
							},
							EnvelopeSettings: &armlogic.AS2EnvelopeSettings{
								AutogenerateFileName:                    to.Ptr(true),
								FileNameTemplate:                        to.Ptr("Test"),
								MessageContentType:                      to.Ptr("text/plain"),
								SuspendMessageOnFileNameGenerationError: to.Ptr(true),
								TransmitFileNameInMimeHeader:            to.Ptr(true),
							},
							ErrorSettings: &armlogic.AS2ErrorSettings{
								ResendIfMDNNotReceived:  to.Ptr(true),
								SuspendDuplicateMessage: to.Ptr(true),
							},
							MdnSettings: &armlogic.AS2MdnSettings{
								DispositionNotificationTo:  to.Ptr("http://tempuri.org"),
								MdnText:                    to.Ptr("Sample"),
								MicHashingAlgorithm:        to.Ptr(armlogic.HashingAlgorithmSHA1),
								NeedMDN:                    to.Ptr(true),
								ReceiptDeliveryURL:         to.Ptr("http://tempuri.org"),
								SendInboundMDNToMessageBox: to.Ptr(true),
								SendMDNAsynchronously:      to.Ptr(true),
								SignMDN:                    to.Ptr(true),
								SignOutboundMDNIfOptional:  to.Ptr(true),
							},
							MessageConnectionSettings: &armlogic.AS2MessageConnectionSettings{
								IgnoreCertificateNameMismatch: to.Ptr(true),
								KeepHTTPConnectionAlive:       to.Ptr(true),
								SupportHTTPStatusCodeContinue: to.Ptr(true),
								UnfoldHTTPHeaders:             to.Ptr(true),
							},
							SecuritySettings: &armlogic.AS2SecuritySettings{
								EnableNRRForInboundDecodedMessages:  to.Ptr(true),
								EnableNRRForInboundEncodedMessages:  to.Ptr(true),
								EnableNRRForInboundMDN:              to.Ptr(true),
								EnableNRRForOutboundDecodedMessages: to.Ptr(true),
								EnableNRRForOutboundEncodedMessages: to.Ptr(true),
								EnableNRRForOutboundMDN:             to.Ptr(true),
								OverrideGroupSigningCertificate:     to.Ptr(false),
							},
							ValidationSettings: &armlogic.AS2ValidationSettings{
								CheckCertificateRevocationListOnReceive: to.Ptr(true),
								CheckCertificateRevocationListOnSend:    to.Ptr(true),
								CheckDuplicateMessage:                   to.Ptr(true),
								CompressMessage:                         to.Ptr(true),
								EncryptMessage:                          to.Ptr(false),
								EncryptionAlgorithm:                     to.Ptr(armlogic.EncryptionAlgorithmAES128),
								InterchangeDuplicatesValidityDays:       to.Ptr[int32](100),
								OverrideMessageProperties:               to.Ptr(true),
								SignMessage:                             to.Ptr(false),
							},
						},
						ReceiverBusinessIdentity: &armlogic.BusinessIdentity{
							Qualifier: to.Ptr("ZZ"),
							Value:     to.Ptr("ZZ"),
						},
						SenderBusinessIdentity: &armlogic.BusinessIdentity{
							Qualifier: to.Ptr("AA"),
							Value:     to.Ptr("AA"),
						},
					},
					SendAgreement: &armlogic.AS2OneWayAgreement{
						ProtocolSettings: &armlogic.AS2ProtocolSettings{
							AcknowledgementConnectionSettings: &armlogic.AS2AcknowledgementConnectionSettings{
								IgnoreCertificateNameMismatch: to.Ptr(true),
								KeepHTTPConnectionAlive:       to.Ptr(true),
								SupportHTTPStatusCodeContinue: to.Ptr(true),
								UnfoldHTTPHeaders:             to.Ptr(true),
							},
							EnvelopeSettings: &armlogic.AS2EnvelopeSettings{
								AutogenerateFileName:                    to.Ptr(true),
								FileNameTemplate:                        to.Ptr("Test"),
								MessageContentType:                      to.Ptr("text/plain"),
								SuspendMessageOnFileNameGenerationError: to.Ptr(true),
								TransmitFileNameInMimeHeader:            to.Ptr(true),
							},
							ErrorSettings: &armlogic.AS2ErrorSettings{
								ResendIfMDNNotReceived:  to.Ptr(true),
								SuspendDuplicateMessage: to.Ptr(true),
							},
							MdnSettings: &armlogic.AS2MdnSettings{
								DispositionNotificationTo:  to.Ptr("http://tempuri.org"),
								MdnText:                    to.Ptr("Sample"),
								MicHashingAlgorithm:        to.Ptr(armlogic.HashingAlgorithmSHA1),
								NeedMDN:                    to.Ptr(true),
								ReceiptDeliveryURL:         to.Ptr("http://tempuri.org"),
								SendInboundMDNToMessageBox: to.Ptr(true),
								SendMDNAsynchronously:      to.Ptr(true),
								SignMDN:                    to.Ptr(true),
								SignOutboundMDNIfOptional:  to.Ptr(true),
							},
							MessageConnectionSettings: &armlogic.AS2MessageConnectionSettings{
								IgnoreCertificateNameMismatch: to.Ptr(true),
								KeepHTTPConnectionAlive:       to.Ptr(true),
								SupportHTTPStatusCodeContinue: to.Ptr(true),
								UnfoldHTTPHeaders:             to.Ptr(true),
							},
							SecuritySettings: &armlogic.AS2SecuritySettings{
								EnableNRRForInboundDecodedMessages:  to.Ptr(true),
								EnableNRRForInboundEncodedMessages:  to.Ptr(true),
								EnableNRRForInboundMDN:              to.Ptr(true),
								EnableNRRForOutboundDecodedMessages: to.Ptr(true),
								EnableNRRForOutboundEncodedMessages: to.Ptr(true),
								EnableNRRForOutboundMDN:             to.Ptr(true),
								OverrideGroupSigningCertificate:     to.Ptr(false),
							},
							ValidationSettings: &armlogic.AS2ValidationSettings{
								CheckCertificateRevocationListOnReceive: to.Ptr(true),
								CheckCertificateRevocationListOnSend:    to.Ptr(true),
								CheckDuplicateMessage:                   to.Ptr(true),
								CompressMessage:                         to.Ptr(true),
								EncryptMessage:                          to.Ptr(false),
								EncryptionAlgorithm:                     to.Ptr(armlogic.EncryptionAlgorithmAES128),
								InterchangeDuplicatesValidityDays:       to.Ptr[int32](100),
								OverrideMessageProperties:               to.Ptr(true),
								SignMessage:                             to.Ptr(false),
							},
						},
						ReceiverBusinessIdentity: &armlogic.BusinessIdentity{
							Qualifier: to.Ptr("AA"),
							Value:     to.Ptr("AA"),
						},
						SenderBusinessIdentity: &armlogic.BusinessIdentity{
							Qualifier: to.Ptr("ZZ"),
							Value:     to.Ptr("ZZ"),
						},
					},
				},
			},
			GuestIdentity: &armlogic.BusinessIdentity{
				Qualifier: to.Ptr("AA"),
				Value:     to.Ptr("AA"),
			},
			GuestPartner: to.Ptr("GuestPartner"),
			HostIdentity: &armlogic.BusinessIdentity{
				Qualifier: to.Ptr("ZZ"),
				Value:     to.Ptr("ZZ"),
			},
			HostPartner: to.Ptr("HostPartner"),
			Metadata:    map[string]any{},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountAgreements_List
	fmt.Println("Call operation: IntegrationAccountAgreements_List")
	integrationAccountAgreementsClientNewListPager := integrationAccountAgreementsClient.NewListPager(testsuite.resourceGroupName, testsuite.integrationAccountName, &armlogic.IntegrationAccountAgreementsClientListOptions{Top: nil,
		Filter: nil,
	})
	for integrationAccountAgreementsClientNewListPager.More() {
		_, err := integrationAccountAgreementsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IntegrationAccountAgreements_Get
	fmt.Println("Call operation: IntegrationAccountAgreements_Get")
	_, err = integrationAccountAgreementsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.agreementName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountAgreements_ListContentCallbackUrl
	fmt.Println("Call operation: IntegrationAccountAgreements_ListContentCallbackURL")
	_, err = integrationAccountAgreementsClient.ListContentCallbackURL(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.agreementName, armlogic.GetCallbackURLParameters{
		KeyType: to.Ptr(armlogic.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountAgreements_Delete
	fmt.Println("Call operation: IntegrationAccountAgreements_Delete")
	_, err = integrationAccountAgreementsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.agreementName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Logic/integrationAccounts/{integrationAccountName}/partners/{partnerName}
func (testsuite *IntegrationAccountsTestSuite) TestIntegrationAccountPartners() {
	var err error
	// From step IntegrationAccountPartners_CreateOrUpdate
	fmt.Println("Call operation: IntegrationAccountPartners_CreateOrUpdate")
	integrationAccountPartnersClient, err := armlogic.NewIntegrationAccountPartnersClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = integrationAccountPartnersClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.partnerName, armlogic.IntegrationAccountPartner{
		Location: to.Ptr(testsuite.location),
		Tags:     map[string]*string{},
		Properties: &armlogic.IntegrationAccountPartnerProperties{
			Content: &armlogic.PartnerContent{
				B2B: &armlogic.B2BPartnerContent{
					BusinessIdentities: []*armlogic.BusinessIdentity{
						{
							Qualifier: to.Ptr("AA"),
							Value:     to.Ptr("ZZ"),
						}},
				},
			},
			Metadata:    map[string]any{},
			PartnerType: to.Ptr(armlogic.PartnerTypeB2B),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountPartners_List
	fmt.Println("Call operation: IntegrationAccountPartners_List")
	integrationAccountPartnersClientNewListPager := integrationAccountPartnersClient.NewListPager(testsuite.resourceGroupName, testsuite.integrationAccountName, &armlogic.IntegrationAccountPartnersClientListOptions{Top: nil,
		Filter: nil,
	})
	for integrationAccountPartnersClientNewListPager.More() {
		_, err := integrationAccountPartnersClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IntegrationAccountPartners_Get
	fmt.Println("Call operation: IntegrationAccountPartners_Get")
	_, err = integrationAccountPartnersClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.partnerName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountPartners_ListContentCallbackUrl
	fmt.Println("Call operation: IntegrationAccountPartners_ListContentCallbackURL")
	_, err = integrationAccountPartnersClient.ListContentCallbackURL(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.partnerName, armlogic.GetCallbackURLParameters{
		KeyType: to.Ptr(armlogic.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountPartners_Delete
	fmt.Println("Call operation: IntegrationAccountPartners_Delete")
	_, err = integrationAccountPartnersClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.partnerName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Logic/integrationAccounts/{integrationAccountName}/maps/{mapName}
func (testsuite *IntegrationAccountsTestSuite) TestIntegrationAccountMaps() {
	var err error
	// From step IntegrationAccountMaps_CreateOrUpdate
	fmt.Println("Call operation: IntegrationAccountMaps_CreateOrUpdate")
	integrationAccountMapsClient, err := armlogic.NewIntegrationAccountMapsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = integrationAccountMapsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.mapName, armlogic.IntegrationAccountMap{
		Location: to.Ptr(testsuite.location),
		Properties: &armlogic.IntegrationAccountMapProperties{
			Content:     to.Ptr("<?xml version=\"1.0\" encoding=\"UTF-16\"?>\r\n<xsl:stylesheet xmlns:xsl=\"http://www.w3.org/1999/XSL/Transform\" xmlns:msxsl=\"urn:schemas-microsoft-com:xslt\" xmlns:var=\"http://schemas.microsoft.com/BizTalk/2003/var\" exclude-result-prefixes=\"msxsl var s0 userCSharp\" version=\"1.0\" xmlns:ns0=\"http://BizTalk_Server_Project4.StringFunctoidsDestinationSchema\" xmlns:s0=\"http://BizTalk_Server_Project4.StringFunctoidsSourceSchema\" xmlns:userCSharp=\"http://schemas.microsoft.com/BizTalk/2003/userCSharp\">\r\n  <xsl:import href=\"http://btsfunctoids.blob.core.windows.net/functoids/functoids.xslt\" />\r\n  <xsl:output omit-xml-declaration=\"yes\" method=\"xml\" version=\"1.0\" />\r\n  <xsl:template match=\"/\">\r\n    <xsl:apply-templates select=\"/s0:Root\" />\r\n  </xsl:template>\r\n  <xsl:template match=\"/s0:Root\">\r\n    <xsl:variable name=\"var:v1\" select=\"userCSharp:StringFind(string(StringFindSource/text()) , &quot;SearchString&quot;)\" />\r\n    <xsl:variable name=\"var:v2\" select=\"userCSharp:StringLeft(string(StringLeftSource/text()) , &quot;2&quot;)\" />\r\n    <xsl:variable name=\"var:v3\" select=\"userCSharp:StringRight(string(StringRightSource/text()) , &quot;2&quot;)\" />\r\n    <xsl:variable name=\"var:v4\" select=\"userCSharp:StringUpperCase(string(UppercaseSource/text()))\" />\r\n    <xsl:variable name=\"var:v5\" select=\"userCSharp:StringLowerCase(string(LowercaseSource/text()))\" />\r\n    <xsl:variable name=\"var:v6\" select=\"userCSharp:StringSize(string(SizeSource/text()))\" />\r\n    <xsl:variable name=\"var:v7\" select=\"userCSharp:StringSubstring(string(StringExtractSource/text()) , &quot;0&quot; , &quot;2&quot;)\" />\r\n    <xsl:variable name=\"var:v8\" select=\"userCSharp:StringConcat(string(StringConcatSource/text()))\" />\r\n    <xsl:variable name=\"var:v9\" select=\"userCSharp:StringTrimLeft(string(StringLeftTrimSource/text()))\" />\r\n    <xsl:variable name=\"var:v10\" select=\"userCSharp:StringTrimRight(string(StringRightTrimSource/text()))\" />\r\n    <ns0:Root>\r\n      <StringFindDestination>\r\n        <xsl:value-of select=\"$var:v1\" />\r\n      </StringFindDestination>\r\n      <StringLeftDestination>\r\n        <xsl:value-of select=\"$var:v2\" />\r\n      </StringLeftDestination>\r\n      <StringRightDestination>\r\n        <xsl:value-of select=\"$var:v3\" />\r\n      </StringRightDestination>\r\n      <UppercaseDestination>\r\n        <xsl:value-of select=\"$var:v4\" />\r\n      </UppercaseDestination>\r\n      <LowercaseDestination>\r\n        <xsl:value-of select=\"$var:v5\" />\r\n      </LowercaseDestination>\r\n      <SizeDestination>\r\n        <xsl:value-of select=\"$var:v6\" />\r\n      </SizeDestination>\r\n      <StringExtractDestination>\r\n        <xsl:value-of select=\"$var:v7\" />\r\n      </StringExtractDestination>\r\n      <StringConcatDestination>\r\n        <xsl:value-of select=\"$var:v8\" />\r\n      </StringConcatDestination>\r\n      <StringLeftTrimDestination>\r\n        <xsl:value-of select=\"$var:v9\" />\r\n      </StringLeftTrimDestination>\r\n      <StringRightTrimDestination>\r\n        <xsl:value-of select=\"$var:v10\" />\r\n      </StringRightTrimDestination>\r\n    </ns0:Root>\r\n  </xsl:template>\r\n</xsl:stylesheet>"),
			ContentType: to.Ptr("application/xml"),
			MapType:     to.Ptr(armlogic.MapTypeXslt),
			Metadata:    map[string]any{},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountMaps_List
	fmt.Println("Call operation: IntegrationAccountMaps_List")
	integrationAccountMapsClientNewListPager := integrationAccountMapsClient.NewListPager(testsuite.resourceGroupName, testsuite.integrationAccountName, &armlogic.IntegrationAccountMapsClientListOptions{Top: nil,
		Filter: nil,
	})
	for integrationAccountMapsClientNewListPager.More() {
		_, err := integrationAccountMapsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IntegrationAccountMaps_Get
	fmt.Println("Call operation: IntegrationAccountMaps_Get")
	_, err = integrationAccountMapsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.mapName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountMaps_ListContentCallbackUrl
	fmt.Println("Call operation: IntegrationAccountMaps_ListContentCallbackURL")
	_, err = integrationAccountMapsClient.ListContentCallbackURL(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.mapName, armlogic.GetCallbackURLParameters{
		KeyType: to.Ptr(armlogic.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountMaps_Delete
	fmt.Println("Call operation: IntegrationAccountMaps_Delete")
	_, err = integrationAccountMapsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.mapName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Logic/integrationAccounts/{integrationAccountName}/schemas/{schemaName}
func (testsuite *IntegrationAccountsTestSuite) TestIntegrationAccountSchemas() {
	var err error
	// From step IntegrationAccountSchemas_CreateOrUpdate
	fmt.Println("Call operation: IntegrationAccountSchemas_CreateOrUpdate")
	integrationAccountSchemasClient, err := armlogic.NewIntegrationAccountSchemasClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = integrationAccountSchemasClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.schemaName, armlogic.IntegrationAccountSchema{
		Location: to.Ptr(testsuite.location),
		Tags: map[string]*string{
			"integrationAccountSchemaName": to.Ptr("IntegrationAccountSchema8120"),
		},
		Properties: &armlogic.IntegrationAccountSchemaProperties{
			Content:     to.Ptr("<?xml version=\"1.0\" encoding=\"utf-16\"?>\r\n<xs:schema xmlns:b=\"http://schemas.microsoft.com/BizTalk/2003\" xmlns=\"http://Inbound_EDI.OrderFile\" targetNamespace=\"http://Inbound_EDI.OrderFile\" xmlns:xs=\"http://www.w3.org/2001/XMLSchema\">\r\n  <xs:annotation>\r\n    <xs:appinfo>\r\n      <b:schemaInfo default_pad_char=\" \" count_positions_by_byte=\"false\" parser_optimization=\"speed\" lookahead_depth=\"3\" suppress_empty_nodes=\"false\" generate_empty_nodes=\"true\" allow_early_termination=\"false\" early_terminate_optional_fields=\"false\" allow_message_breakup_of_infix_root=\"false\" compile_parse_tables=\"false\" standard=\"Flat File\" root_reference=\"OrderFile\" />\r\n      <schemaEditorExtension:schemaInfo namespaceAlias=\"b\" extensionClass=\"Microsoft.BizTalk.FlatFileExtension.FlatFileExtension\" standardName=\"Flat File\" xmlns:schemaEditorExtension=\"http://schemas.microsoft.com/BizTalk/2003/SchemaEditorExtensions\" />\r\n    </xs:appinfo>\r\n  </xs:annotation>\r\n  <xs:element name=\"OrderFile\">\r\n    <xs:annotation>\r\n      <xs:appinfo>\r\n        <b:recordInfo structure=\"delimited\" preserve_delimiter_for_empty_data=\"true\" suppress_trailing_delimiters=\"false\" sequence_number=\"1\" />\r\n      </xs:appinfo>\r\n    </xs:annotation>\r\n    <xs:complexType>\r\n      <xs:sequence>\r\n        <xs:annotation>\r\n          <xs:appinfo>\r\n            <b:groupInfo sequence_number=\"0\" />\r\n          </xs:appinfo>\r\n        </xs:annotation>\r\n        <xs:element name=\"Order\">\r\n          <xs:annotation>\r\n            <xs:appinfo>\r\n              <b:recordInfo sequence_number=\"1\" structure=\"delimited\" preserve_delimiter_for_empty_data=\"true\" suppress_trailing_delimiters=\"false\" child_delimiter_type=\"hex\" child_delimiter=\"0x0D 0x0A\" child_order=\"infix\" />\r\n            </xs:appinfo>\r\n          </xs:annotation>\r\n          <xs:complexType>\r\n            <xs:sequence>\r\n              <xs:annotation>\r\n                <xs:appinfo>\r\n                  <b:groupInfo sequence_number=\"0\" />\r\n                </xs:appinfo>\r\n              </xs:annotation>\r\n              <xs:element name=\"Header\">\r\n                <xs:annotation>\r\n                  <xs:appinfo>\r\n                    <b:recordInfo sequence_number=\"1\" structure=\"delimited\" preserve_delimiter_for_empty_data=\"true\" suppress_trailing_delimiters=\"false\" child_delimiter_type=\"char\" child_delimiter=\"|\" child_order=\"infix\" tag_name=\"HDR|\" />\r\n                  </xs:appinfo>\r\n                </xs:annotation>\r\n                <xs:complexType>\r\n                  <xs:sequence>\r\n                    <xs:annotation>\r\n                      <xs:appinfo>\r\n                        <b:groupInfo sequence_number=\"0\" />\r\n                      </xs:appinfo>\r\n                    </xs:annotation>\r\n                    <xs:element name=\"PODate\" type=\"xs:string\">\r\n                      <xs:annotation>\r\n                        <xs:appinfo>\r\n                          <b:fieldInfo sequence_number=\"1\" justification=\"left\" />\r\n                        </xs:appinfo>\r\n                      </xs:annotation>\r\n                    </xs:element>\r\n                    <xs:element name=\"PONumber\" type=\"xs:string\">\r\n                      <xs:annotation>\r\n                        <xs:appinfo>\r\n                          <b:fieldInfo justification=\"left\" sequence_number=\"2\" />\r\n                        </xs:appinfo>\r\n                      </xs:annotation>\r\n                    </xs:element>\r\n                    <xs:element name=\"CustomerID\" type=\"xs:string\">\r\n                      <xs:annotation>\r\n                        <xs:appinfo>\r\n                          <b:fieldInfo sequence_number=\"3\" justification=\"left\" />\r\n                        </xs:appinfo>\r\n                      </xs:annotation>\r\n                    </xs:element>\r\n                    <xs:element name=\"CustomerContactName\" type=\"xs:string\">\r\n                      <xs:annotation>\r\n                        <xs:appinfo>\r\n                          <b:fieldInfo sequence_number=\"4\" justification=\"left\" />\r\n                        </xs:appinfo>\r\n                      </xs:annotation>\r\n                    </xs:element>\r\n                    <xs:element name=\"CustomerContactPhone\" type=\"xs:string\">\r\n                      <xs:annotation>\r\n                        <xs:appinfo>\r\n                          <b:fieldInfo sequence_number=\"5\" justification=\"left\" />\r\n                        </xs:appinfo>\r\n                      </xs:annotation>\r\n                    </xs:element>\r\n                  </xs:sequence>\r\n                </xs:complexType>\r\n              </xs:element>\r\n              <xs:element minOccurs=\"1\" maxOccurs=\"unbounded\" name=\"LineItems\">\r\n                <xs:annotation>\r\n                  <xs:appinfo>\r\n                    <b:recordInfo sequence_number=\"2\" structure=\"delimited\" preserve_delimiter_for_empty_data=\"true\" suppress_trailing_delimiters=\"false\" child_delimiter_type=\"char\" child_delimiter=\"|\" child_order=\"infix\" tag_name=\"DTL|\" />\r\n                  </xs:appinfo>\r\n                </xs:annotation>\r\n                <xs:complexType>\r\n                  <xs:sequence>\r\n                    <xs:annotation>\r\n                      <xs:appinfo>\r\n                        <b:groupInfo sequence_number=\"0\" />\r\n                      </xs:appinfo>\r\n                    </xs:annotation>\r\n                    <xs:element name=\"PONumber\" type=\"xs:string\">\r\n                      <xs:annotation>\r\n                        <xs:appinfo>\r\n                          <b:fieldInfo sequence_number=\"1\" justification=\"left\" />\r\n                        </xs:appinfo>\r\n                      </xs:annotation>\r\n                    </xs:element>\r\n                    <xs:element name=\"ItemOrdered\" type=\"xs:string\">\r\n                      <xs:annotation>\r\n                        <xs:appinfo>\r\n                          <b:fieldInfo sequence_number=\"2\" justification=\"left\" />\r\n                        </xs:appinfo>\r\n                      </xs:annotation>\r\n                    </xs:element>\r\n                    <xs:element name=\"Quantity\" type=\"xs:string\">\r\n                      <xs:annotation>\r\n                        <xs:appinfo>\r\n                          <b:fieldInfo sequence_number=\"3\" justification=\"left\" />\r\n                        </xs:appinfo>\r\n                      </xs:annotation>\r\n                    </xs:element>\r\n                    <xs:element name=\"UOM\" type=\"xs:string\">\r\n                      <xs:annotation>\r\n                        <xs:appinfo>\r\n                          <b:fieldInfo sequence_number=\"4\" justification=\"left\" />\r\n                        </xs:appinfo>\r\n                      </xs:annotation>\r\n                    </xs:element>\r\n                    <xs:element name=\"Price\" type=\"xs:string\">\r\n                      <xs:annotation>\r\n                        <xs:appinfo>\r\n                          <b:fieldInfo sequence_number=\"5\" justification=\"left\" />\r\n                        </xs:appinfo>\r\n                      </xs:annotation>\r\n                    </xs:element>\r\n                    <xs:element name=\"ExtendedPrice\" type=\"xs:string\">\r\n                      <xs:annotation>\r\n                        <xs:appinfo>\r\n                          <b:fieldInfo sequence_number=\"6\" justification=\"left\" />\r\n                        </xs:appinfo>\r\n                      </xs:annotation>\r\n                    </xs:element>\r\n                    <xs:element name=\"Description\" type=\"xs:string\">\r\n                      <xs:annotation>\r\n                        <xs:appinfo>\r\n                          <b:fieldInfo sequence_number=\"7\" justification=\"left\" />\r\n                        </xs:appinfo>\r\n                      </xs:annotation>\r\n                    </xs:element>\r\n                  </xs:sequence>\r\n                </xs:complexType>\r\n              </xs:element>\r\n            </xs:sequence>\r\n          </xs:complexType>\r\n        </xs:element>\r\n      </xs:sequence>\r\n    </xs:complexType>\r\n  </xs:element>\r\n</xs:schema>"),
			ContentType: to.Ptr("application/xml"),
			Metadata:    map[string]any{},
			SchemaType:  to.Ptr(armlogic.SchemaTypeXML),
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountSchemas_List
	fmt.Println("Call operation: IntegrationAccountSchemas_List")
	integrationAccountSchemasClientNewListPager := integrationAccountSchemasClient.NewListPager(testsuite.resourceGroupName, testsuite.integrationAccountName, &armlogic.IntegrationAccountSchemasClientListOptions{Top: nil,
		Filter: nil,
	})
	for integrationAccountSchemasClientNewListPager.More() {
		_, err := integrationAccountSchemasClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IntegrationAccountSchemas_Get
	fmt.Println("Call operation: IntegrationAccountSchemas_Get")
	_, err = integrationAccountSchemasClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.schemaName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountSchemas_ListContentCallbackUrl
	fmt.Println("Call operation: IntegrationAccountSchemas_ListContentCallbackURL")
	_, err = integrationAccountSchemasClient.ListContentCallbackURL(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.schemaName, armlogic.GetCallbackURLParameters{
		KeyType: to.Ptr(armlogic.KeyTypePrimary),
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountSchemas_Delete
	fmt.Println("Call operation: IntegrationAccountSchemas_Delete")
	_, err = integrationAccountSchemasClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.schemaName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Logic/integrationAccounts/{integrationAccountName}/batchConfigurations/{batchConfigurationName}
func (testsuite *IntegrationAccountsTestSuite) TestIntegrationAccountBatchConfigurations() {
	var err error
	// From step IntegrationAccountBatchConfigurations_CreateOrUpdate
	fmt.Println("Call operation: IntegrationAccountBatchConfigurations_CreateOrUpdate")
	integrationAccountBatchConfigurationsClient, err := armlogic.NewIntegrationAccountBatchConfigurationsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = integrationAccountBatchConfigurationsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.batchConfigurationName, armlogic.BatchConfiguration{
		Location: to.Ptr(testsuite.location),
		Properties: &armlogic.BatchConfigurationProperties{
			BatchGroupName: to.Ptr("DEFAULT"),
			ReleaseCriteria: &armlogic.BatchReleaseCriteria{
				BatchSize:    to.Ptr[int32](234567),
				MessageCount: to.Ptr[int32](10),
				Recurrence: &armlogic.WorkflowTriggerRecurrence{
					Frequency: to.Ptr(armlogic.RecurrenceFrequencyMinute),
					Interval:  to.Ptr[int32](1),
					StartTime: to.Ptr("2017-03-24T11:43:00"),
					TimeZone:  to.Ptr("India Standard Time"),
				},
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountBatchConfigurations_List
	fmt.Println("Call operation: IntegrationAccountBatchConfigurations_List")
	integrationAccountBatchConfigurationsClientNewListPager := integrationAccountBatchConfigurationsClient.NewListPager(testsuite.resourceGroupName, testsuite.integrationAccountName, nil)
	for integrationAccountBatchConfigurationsClientNewListPager.More() {
		_, err := integrationAccountBatchConfigurationsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IntegrationAccountBatchConfigurations_Get
	fmt.Println("Call operation: IntegrationAccountBatchConfigurations_Get")
	_, err = integrationAccountBatchConfigurationsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.batchConfigurationName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountBatchConfigurations_Delete
	fmt.Println("Call operation: IntegrationAccountBatchConfigurations_Delete")
	_, err = integrationAccountBatchConfigurationsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.batchConfigurationName, nil)
	testsuite.Require().NoError(err)
}

// Microsoft.Logic/integrationAccounts/{integrationAccountName}/sessions/{sessionName}
func (testsuite *IntegrationAccountsTestSuite) TestIntegrationAccountSessions() {
	var err error
	// From step IntegrationAccountSessions_CreateOrUpdate
	fmt.Println("Call operation: IntegrationAccountSessions_CreateOrUpdate")
	integrationAccountSessionsClient, err := armlogic.NewIntegrationAccountSessionsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = integrationAccountSessionsClient.CreateOrUpdate(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.sessionName, armlogic.IntegrationAccountSession{
		Properties: &armlogic.IntegrationAccountSessionProperties{
			Content: map[string]any{
				"controlNumber":            "1234",
				"controlNumberChangedTime": "2017-02-21T22:30:11.9923759Z",
			},
		},
	}, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountSessions_List
	fmt.Println("Call operation: IntegrationAccountSessions_List")
	integrationAccountSessionsClientNewListPager := integrationAccountSessionsClient.NewListPager(testsuite.resourceGroupName, testsuite.integrationAccountName, &armlogic.IntegrationAccountSessionsClientListOptions{Top: nil,
		Filter: nil,
	})
	for integrationAccountSessionsClientNewListPager.More() {
		_, err := integrationAccountSessionsClientNewListPager.NextPage(testsuite.ctx)
		testsuite.Require().NoError(err)
		break
	}

	// From step IntegrationAccountSessions_Get
	fmt.Println("Call operation: IntegrationAccountSessions_Get")
	_, err = integrationAccountSessionsClient.Get(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.sessionName, nil)
	testsuite.Require().NoError(err)

	// From step IntegrationAccountSessions_Delete
	fmt.Println("Call operation: IntegrationAccountSessions_Delete")
	_, err = integrationAccountSessionsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, testsuite.sessionName, nil)
	testsuite.Require().NoError(err)
}

func (testsuite *IntegrationAccountsTestSuite) Cleanup() {
	var err error
	// From step IntegrationAccounts_Delete
	fmt.Println("Call operation: IntegrationAccounts_Delete")
	integrationAccountsClient, err := armlogic.NewIntegrationAccountsClient(testsuite.subscriptionId, testsuite.cred, testsuite.options)
	testsuite.Require().NoError(err)
	_, err = integrationAccountsClient.Delete(testsuite.ctx, testsuite.resourceGroupName, testsuite.integrationAccountName, nil)
	testsuite.Require().NoError(err)
}
