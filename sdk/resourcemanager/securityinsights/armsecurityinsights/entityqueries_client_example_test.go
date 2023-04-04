//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armsecurityinsights_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e24bbf6a66cb0a19c072c6f15cee163acbd7acf7/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-09-01-preview/examples/entityQueries/GetEntityQueries.json
func ExampleEntityQueriesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsecurityinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := clientFactory.NewEntityQueriesClient().NewListPager("myRg", "myWorkspace", &armsecurityinsights.EntityQueriesClientListOptions{Kind: to.Ptr(armsecurityinsights.Enum13Expansion)})
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range page.Value {
			// You could use page here. We use blank identifier for just demo purposes.
			_ = v
		}
		// If the HTTP response code is 200 as defined in example definition, your page structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
		// page.EntityQueryList = armsecurityinsights.EntityQueryList{
		// 	Value: []armsecurityinsights.EntityQueryClassification{
		// 		&armsecurityinsights.ExpansionEntityQuery{
		// 			Name: to.Ptr("37ca3555-c135-4a73-a65e-9c1d00323f5d"),
		// 			Type: to.Ptr("Microsoft.SecurityInsights/entityQueries"),
		// 			ID: to.Ptr("/subscriptions/d0cfe6b2-9ac0-4464-9919-dccaee2e48c0/resourceGroups/myRg/providers/Microsoft.OperationalInsights/workspaces/myWorkspace/providers/Microsoft.SecurityInsights/entityQueries/37ca3555-c135-4a73-a65e-9c1d00323f5d"),
		// 			Kind: to.Ptr(armsecurityinsights.EntityQueryKindExpansion),
		// 			Properties: &armsecurityinsights.ExpansionEntityQueriesProperties{
		// 				DataSources: []*string{
		// 					to.Ptr("AzureActivity")},
		// 					DisplayName: to.Ptr("Least active accounts on Azure from this IP"),
		// 					InputEntityType: to.Ptr(armsecurityinsights.EntityTypeIP),
		// 					InputFields: []*string{
		// 						to.Ptr("address")},
		// 						OutputEntityTypes: []*armsecurityinsights.EntityType{
		// 							to.Ptr(armsecurityinsights.EntityTypeAccount)},
		// 							QueryTemplate: to.Ptr("let AccountActivity_byIP = (v_IP_Address:string){\r\n                            AzureActivity\r\n                            | where Caller != '' and CallerIpAddress == v_IP_Address\r\n                            | summarize Account_Aux_StartTime = min(TimeGenerated), Account_Aux_EndTime = max(TimeGenerated), Count = count() by Caller, TenantId\r\n                            | top 10 by Count asc nulls last \r\n                            | extend UPN = iff(Caller contains '@', Caller, ''), Account_AadUserId = iff(Caller !contains '@', Caller,'')\r\n                            | extend Account_Name = split(UPN,'@')[0] , Account_UPNSuffix = split(UPN,'@')[1]\r\n                            | project Account_Name, Account_UPNSuffix, Account_AadUserId, Account_AadTenantId=TenantId, Account_Aux_StartTime , Account_Aux_EndTime};\r\n                            AccountActivity_byIP('<address>')"),
		// 						},
		// 					},
		// 					&armsecurityinsights.ExpansionEntityQuery{
		// 						Name: to.Ptr("97a1d515-abf2-4231-9a35-985f9de0bb91"),
		// 						Type: to.Ptr("Microsoft.SecurityInsights/entityQueries"),
		// 						ID: to.Ptr("/subscriptions/d0cfe6b2-9ac0-4464-9919-dccaee2e48c0/resourceGroups/myRg/providers/Microsoft.OperationalInsights/workspaces/myWorkspace/providers/Microsoft.SecurityInsights/entityQueries/97a1d515-abf2-4231-9a35-985f9de0bb91"),
		// 						Kind: to.Ptr(armsecurityinsights.EntityQueryKindExpansion),
		// 						Properties: &armsecurityinsights.ExpansionEntityQueriesProperties{
		// 							DataSources: []*string{
		// 								to.Ptr("AzureActivity")},
		// 								DisplayName: to.Ptr("Most active accounts on Azure from this IP"),
		// 								InputEntityType: to.Ptr(armsecurityinsights.EntityTypeIP),
		// 								InputFields: []*string{
		// 									to.Ptr("address")},
		// 									OutputEntityTypes: []*armsecurityinsights.EntityType{
		// 										to.Ptr(armsecurityinsights.EntityTypeAccount)},
		// 										QueryTemplate: to.Ptr("let AccountActivity_byIP = (v_IP_Address:string){\r\n                            AzureActivity\r\n                            | where Caller != '' and CallerIpAddress == v_IP_Address\r\n                            | summarize Account_Aux_StartTime = min(TimeGenerated), Account_Aux_EndTime = max(TimeGenerated), Count = count() by Caller, TenantId\r\n                            | top 10 by Count desc nulls last \r\n                            | extend UPN = iff(Caller contains '@', Caller, ''), Account_AadUserId = iff(Caller !contains '@', Caller,'')\r\n                            | extend Account_Name = split(UPN,'@')[0] , Account_UPNSuffix = split(UPN,'@')[1]\r\n                            | project Account_Name, Account_UPNSuffix, Account_AadUserId, Account_AadTenantId=TenantId, Account_Aux_StartTime , Account_Aux_EndTime};\r\n                            AccountActivity_byIP('<address>')"),
		// 									},
		// 							}},
		// 						}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e24bbf6a66cb0a19c072c6f15cee163acbd7acf7/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-09-01-preview/examples/entityQueries/GetActivityEntityQueryById.json
func ExampleEntityQueriesClient_Get_getAnActivityEntityQuery() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsecurityinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewEntityQueriesClient().Get(ctx, "myRg", "myWorkspace", "07da3cc8-c8ad-4710-a44e-334cdcb7882b", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armsecurityinsights.EntityQueriesClientGetResponse{
	// 	                            EntityQueryClassification: &armsecurityinsights.ActivityEntityQuery{
	// 		Name: to.Ptr("07da3cc8-c8ad-4710-a44e-334cdcb7882b"),
	// 		Type: to.Ptr("Microsoft.SecurityInsights/entityQueries"),
	// 		ID: to.Ptr("/subscriptions/d0cfe6b2-9ac0-4464-9919-dccaee2e48c0/resourceGroups/myRg/providers/Microsoft.OperationalInsights/workspaces/myWorkspace/providers/Microsoft.SecurityInsights/entityQueries/07da3cc8-c8ad-4710-a44e-334cdcb7882b"),
	// 		Kind: to.Ptr(armsecurityinsights.EntityQueryKindActivity),
	// 		Properties: &armsecurityinsights.ActivityEntityQueriesProperties{
	// 			Description: to.Ptr("Account deleted on host"),
	// 			Content: to.Ptr("On '{{Computer}}' the account '{{TargetAccount}}' was deleted by '{{AddedBy}}'"),
	// 			CreatedTimeUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2019-01-01T13:15:30Z"); return t}()),
	// 			Enabled: to.Ptr(true),
	// 			EntitiesFilter: map[string][]*string{
	// 				"Host_OsFamily": []*string{
	// 					to.Ptr("Windows")},
	// 				},
	// 				InputEntityType: to.Ptr(armsecurityinsights.EntityTypeHost),
	// 				LastModifiedTimeUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2019-01-01T13:15:30Z"); return t}()),
	// 				QueryDefinitions: &armsecurityinsights.ActivityEntityQueriesPropertiesQueryDefinitions{
	// 					Query: to.Ptr("let GetAccountActions = (v_Host_Name:string, v_Host_NTDomain:string, v_Host_DnsDomain:string, v_Host_AzureID:string, v_Host_OMSAgentID:string){\nSecurityEvent\n| where EventID in (4725, 4726, 4767, 4720, 4722, 4723, 4724)\n// parsing for Host to handle variety of conventions coming from data\n| extend Host_HostName = case(\nComputer has '@', tostring(split(Computer, '@')[0]),\nComputer has '\\\\', tostring(split(Computer, '\\\\')[1]),\nComputer has '.', tostring(split(Computer, '.')[0]),\nComputer\n)\n| extend Host_NTDomain = case(\nComputer has '\\\\', tostring(split(Computer, '\\\\')[0]), \nComputer has '.', tostring(split(Computer, '.')[-2]), \nComputer\n)\n| extend Host_DnsDomain = case(\nComputer has '\\\\', tostring(split(Computer, '\\\\')[0]), \nComputer has '.', strcat_array(array_slice(split(Computer,'.'),-2,-1),'.'), \nComputer\n)\n| where (Host_HostName =~ v_Host_Name and Host_NTDomain =~ v_Host_NTDomain) \nor (Host_HostName =~ v_Host_Name and Host_DnsDomain =~ v_Host_DnsDomain) \nor v_Host_AzureID =~ _ResourceId \nor v_Host_OMSAgentID == SourceComputerId\n| project TimeGenerated, EventID, Activity, Computer, TargetAccount, TargetUserName, TargetDomainName, TargetSid, SubjectUserName, SubjectUserSid, _ResourceId, SourceComputerId\n| extend AddedBy = SubjectUserName\n// Future support for Activities\n| extend timestamp = TimeGenerated, HostCustomEntity = Computer, AccountCustomEntity = TargetAccount\n};\nGetAccountActions('{{Host_HostName}}', '{{Host_NTDomain}}', '{{Host_DnsDomain}}', '{{Host_AzureID}}', '{{Host_OMSAgentID}}')\n \n| where EventID == 4726 "),
	// 				},
	// 				RequiredInputFieldsSets: [][]*string{
	// 					[]*string{
	// 						to.Ptr("Host_HostName"),
	// 						to.Ptr("Host_NTDomain")},
	// 						[]*string{
	// 							to.Ptr("Host_HostName"),
	// 							to.Ptr("Host_DnsDomain")},
	// 							[]*string{
	// 								to.Ptr("Host_AzureID")},
	// 								[]*string{
	// 									to.Ptr("Host_OMSAgentID")}},
	// 									Title: to.Ptr("An account was deleted on this host"),
	// 								},
	// 							},
	// 							                        }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e24bbf6a66cb0a19c072c6f15cee163acbd7acf7/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-09-01-preview/examples/entityQueries/GetExpansionEntityQueryById.json
func ExampleEntityQueriesClient_Get_getAnExpansionEntityQuery() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsecurityinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewEntityQueriesClient().Get(ctx, "myRg", "myWorkspace", "07da3cc8-c8ad-4710-a44e-334cdcb7882b", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armsecurityinsights.EntityQueriesClientGetResponse{
	// 	                            EntityQueryClassification: &armsecurityinsights.ExpansionEntityQuery{
	// 		Name: to.Ptr("07da3cc8-c8ad-4710-a44e-334cdcb7882b"),
	// 		Type: to.Ptr("Microsoft.SecurityInsights/entityQueries"),
	// 		ID: to.Ptr("/subscriptions/d0cfe6b2-9ac0-4464-9919-dccaee2e48c0/resourceGroups/myRg/providers/Microsoft.OperationalInsights/workspaces/myWorkspace/providers/Microsoft.SecurityInsights/entityQueries/07da3cc8-c8ad-4710-a44e-334cdcb7882b"),
	// 		Kind: to.Ptr(armsecurityinsights.EntityQueryKindExpansion),
	// 		Properties: &armsecurityinsights.ExpansionEntityQueriesProperties{
	// 			DataSources: []*string{
	// 				to.Ptr("SecurityEvent")},
	// 				DisplayName: to.Ptr("Parent processes running on host"),
	// 				InputEntityType: to.Ptr(armsecurityinsights.EntityTypeHost),
	// 				InputFields: []*string{
	// 					to.Ptr("hostName")},
	// 					OutputEntityTypes: []*armsecurityinsights.EntityType{
	// 						to.Ptr(armsecurityinsights.EntityTypeProcess)},
	// 						QueryTemplate: to.Ptr("let GetParentProcessesOnHost = (v_Host_HostName:string){\r\n                            SecurityEvent \r\n                            | where EventID == 4688 \r\n                            | where isnotempty(ParentProcessName)\r\n                            | where NewProcessName !contains ':\\\\Windows\\\\System32\\\\conhost.exe' and ParentProcessName !contains ':\\\\Windows\\\\System32\\\\conhost.exe'\r\n                            and NewProcessName !contains ':\\\\Windows\\\\Microsoft.NET\\\\Framework64\\\\v2.0.50727\\\\csc.exe' and ParentProcessName !contains ':\\\\Windows\\\\Microsoft.NET\\\\Framework64\\\\v2.0.50727\\\\csc.exe'\r\n                            and NewProcessName !contains ':\\\\Windows\\\\Microsoft.NET\\\\Framework64\\\\v2.0.50727\\\\cvtres.exe' and ParentProcessName !contains ':\\\\Windows\\\\Microsoft.NET\\\\Framework64\\\\v2.0.50727\\\\cvtres.exe'\r\n                            and NewProcessName!contains ':\\\\Program Files\\\\Microsoft Monitoring Agent\\\\Agent\\\\MonitoringHost.exe' and ParentProcessName !contains ':\\\\Program Files\\\\Microsoft Monitoring Agent\\\\Agent\\\\MonitoringHost.exe'\r\n                            and ParentProcessName !contains ':\\\\Windows\\\\CCM\\\\CcmExec.exe'\r\n                            | where(ParentProcessName !contains ':\\\\Windows\\\\System32\\\\svchost.exe' and (NewProcessName !contains ':\\\\Windows\\\\System32\\\\wbem\\\\WmiPrvSE.exe' or NewProcessName !contains ':\\\\Windows\\\\SysWOW64\\\\wbem\\\\WmiPrvSE.exe'))\r\n                            | where(ParentProcessName !contains ':\\\\Windows\\\\System32\\\\services.exe' and NewProcessName !contains ':\\\\Windows\\\\servicing\\\\TrustedInstaller.exe')\r\n                            | where toupper(Computer) contains v_Host_HostName or toupper(WorkstationName) contains v_Host_HostName\r\n                            | summarize min(TimeGenerated), max(TimeGenerated) by Account, Computer, ParentProcessName, NewProcessName, CommandLine, ProcessId\r\n                            | project min_TimeGenerated, max_TimeGenerated, Account, Computer, ParentProcessName, NewProcessName, CommandLine, ProcessId\r\n                            | project-rename Process_Host_UnstructuredName=Computer, Process_Account_UnstructuredName=Account, Process_CommandLine=CommandLine, Process_ProcessId=ProcessId, Process_ImageFile_FullPath=NewProcessName, Process_ParentProcess_ImageFile_FullPath=ParentProcessName\r\n                            | top 10 by min_TimeGenerated asc};\r\n                            GetParentProcessesOnHost(toupper('<hostName>'))"),
	// 					},
	// 				},
	// 				                        }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e24bbf6a66cb0a19c072c6f15cee163acbd7acf7/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-09-01-preview/examples/entityQueries/CreateEntityQueryActivity.json
func ExampleEntityQueriesClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsecurityinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := clientFactory.NewEntityQueriesClient().CreateOrUpdate(ctx, "myRg", "myWorkspace", "07da3cc8-c8ad-4710-a44e-334cdcb7882b", &armsecurityinsights.ActivityCustomEntityQuery{
		Etag: to.Ptr("\"0300bf09-0000-0000-0000-5c37296e0000\""),
		Kind: to.Ptr(armsecurityinsights.CustomEntityQueryKindActivity),
		Properties: &armsecurityinsights.ActivityEntityQueriesProperties{
			Description: to.Ptr("Account deleted on host"),
			Content:     to.Ptr("On '{{Computer}}' the account '{{TargetAccount}}' was deleted by '{{AddedBy}}'"),
			Enabled:     to.Ptr(true),
			EntitiesFilter: map[string][]*string{
				"Host_OsFamily": {
					to.Ptr("Windows")},
			},
			InputEntityType: to.Ptr(armsecurityinsights.EntityTypeHost),
			QueryDefinitions: &armsecurityinsights.ActivityEntityQueriesPropertiesQueryDefinitions{
				Query: to.Ptr("let GetAccountActions = (v_Host_Name:string, v_Host_NTDomain:string, v_Host_DnsDomain:string, v_Host_AzureID:string, v_Host_OMSAgentID:string){\nSecurityEvent\n| where EventID in (4725, 4726, 4767, 4720, 4722, 4723, 4724)\n// parsing for Host to handle variety of conventions coming from data\n| extend Host_HostName = case(\nComputer has '@', tostring(split(Computer, '@')[0]),\nComputer has '\\\\', tostring(split(Computer, '\\\\')[1]),\nComputer has '.', tostring(split(Computer, '.')[0]),\nComputer\n)\n| extend Host_NTDomain = case(\nComputer has '\\\\', tostring(split(Computer, '\\\\')[0]), \nComputer has '.', tostring(split(Computer, '.')[-2]), \nComputer\n)\n| extend Host_DnsDomain = case(\nComputer has '\\\\', tostring(split(Computer, '\\\\')[0]), \nComputer has '.', strcat_array(array_slice(split(Computer,'.'),-2,-1),'.'), \nComputer\n)\n| where (Host_HostName =~ v_Host_Name and Host_NTDomain =~ v_Host_NTDomain) \nor (Host_HostName =~ v_Host_Name and Host_DnsDomain =~ v_Host_DnsDomain) \nor v_Host_AzureID =~ _ResourceId \nor v_Host_OMSAgentID == SourceComputerId\n| project TimeGenerated, EventID, Activity, Computer, TargetAccount, TargetUserName, TargetDomainName, TargetSid, SubjectUserName, SubjectUserSid, _ResourceId, SourceComputerId\n| extend AddedBy = SubjectUserName\n// Future support for Activities\n| extend timestamp = TimeGenerated, HostCustomEntity = Computer, AccountCustomEntity = TargetAccount\n};\nGetAccountActions('{{Host_HostName}}', '{{Host_NTDomain}}', '{{Host_DnsDomain}}', '{{Host_AzureID}}', '{{Host_OMSAgentID}}')\n \n| where EventID == 4726 "),
			},
			RequiredInputFieldsSets: [][]*string{
				{
					to.Ptr("Host_HostName"),
					to.Ptr("Host_NTDomain")},
				{
					to.Ptr("Host_HostName"),
					to.Ptr("Host_DnsDomain")},
				{
					to.Ptr("Host_AzureID")},
				{
					to.Ptr("Host_OMSAgentID")}},
			Title: to.Ptr("An account was deleted on this host"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res = armsecurityinsights.EntityQueriesClientCreateOrUpdateResponse{
	// 	                            EntityQueryClassification: &armsecurityinsights.ActivityEntityQuery{
	// 		Name: to.Ptr("07da3cc8-c8ad-4710-a44e-334cdcb7882b"),
	// 		Type: to.Ptr("Microsoft.SecurityInsights/entityQueries"),
	// 		ID: to.Ptr("/subscriptions/d0cfe6b2-9ac0-4464-9919-dccaee2e48c0/resourceGroups/myRg/providers/Microsoft.OperationalInsights/workspaces/myWorkspace/providers/Microsoft.SecurityInsights/entityQueries/07da3cc8-c8ad-4710-a44e-334cdcb7882b"),
	// 		Etag: to.Ptr("\"0300bf09-0000-0000-0000-5c37296e0000\""),
	// 		Kind: to.Ptr(armsecurityinsights.EntityQueryKindActivity),
	// 		Properties: &armsecurityinsights.ActivityEntityQueriesProperties{
	// 			Description: to.Ptr("Account deleted on host"),
	// 			Content: to.Ptr("On '{{Computer}}' the account '{{TargetAccount}}' was deleted by '{{AddedBy}}'"),
	// 			CreatedTimeUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2019-01-01T13:15:30Z"); return t}()),
	// 			Enabled: to.Ptr(true),
	// 			EntitiesFilter: map[string][]*string{
	// 				"Host_OsFamily": []*string{
	// 					to.Ptr("Windows")},
	// 				},
	// 				InputEntityType: to.Ptr(armsecurityinsights.EntityTypeHost),
	// 				LastModifiedTimeUTC: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2019-01-01T13:15:30Z"); return t}()),
	// 				QueryDefinitions: &armsecurityinsights.ActivityEntityQueriesPropertiesQueryDefinitions{
	// 					Query: to.Ptr("let GetAccountActions = (v_Host_Name:string, v_Host_NTDomain:string, v_Host_DnsDomain:string, v_Host_AzureID:string, v_Host_OMSAgentID:string){\nSecurityEvent\n| where EventID in (4725, 4726, 4767, 4720, 4722, 4723, 4724)\n// parsing for Host to handle variety of conventions coming from data\n| extend Host_HostName = case(\nComputer has '@', tostring(split(Computer, '@')[0]),\nComputer has '\\\\', tostring(split(Computer, '\\\\')[1]),\nComputer has '.', tostring(split(Computer, '.')[0]),\nComputer\n)\n| extend Host_NTDomain = case(\nComputer has '\\\\', tostring(split(Computer, '\\\\')[0]), \nComputer has '.', tostring(split(Computer, '.')[-2]), \nComputer\n)\n| extend Host_DnsDomain = case(\nComputer has '\\\\', tostring(split(Computer, '\\\\')[0]), \nComputer has '.', strcat_array(array_slice(split(Computer,'.'),-2,-1),'.'), \nComputer\n)\n| where (Host_HostName =~ v_Host_Name and Host_NTDomain =~ v_Host_NTDomain) \nor (Host_HostName =~ v_Host_Name and Host_DnsDomain =~ v_Host_DnsDomain) \nor v_Host_AzureID =~ _ResourceId \nor v_Host_OMSAgentID == SourceComputerId\n| project TimeGenerated, EventID, Activity, Computer, TargetAccount, TargetUserName, TargetDomainName, TargetSid, SubjectUserName, SubjectUserSid, _ResourceId, SourceComputerId\n| extend AddedBy = SubjectUserName\n// Future support for Activities\n| extend timestamp = TimeGenerated, HostCustomEntity = Computer, AccountCustomEntity = TargetAccount\n};\nGetAccountActions('{{Host_HostName}}', '{{Host_NTDomain}}', '{{Host_DnsDomain}}', '{{Host_AzureID}}', '{{Host_OMSAgentID}}')\n \n| where EventID == 4726 "),
	// 				},
	// 				RequiredInputFieldsSets: [][]*string{
	// 					[]*string{
	// 						to.Ptr("Host_HostName"),
	// 						to.Ptr("Host_NTDomain")},
	// 						[]*string{
	// 							to.Ptr("Host_HostName"),
	// 							to.Ptr("Host_DnsDomain")},
	// 							[]*string{
	// 								to.Ptr("Host_AzureID")},
	// 								[]*string{
	// 									to.Ptr("Host_OMSAgentID")}},
	// 									Title: to.Ptr("An account was deleted on this host"),
	// 								},
	// 							},
	// 							                        }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/e24bbf6a66cb0a19c072c6f15cee163acbd7acf7/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-09-01-preview/examples/entityQueries/DeleteEntityQuery.json
func ExampleEntityQueriesClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armsecurityinsights.NewClientFactory("<subscription-id>", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = clientFactory.NewEntityQueriesClient().Delete(ctx, "myRg", "myWorkspace", "07da3cc8-c8ad-4710-a44e-334cdcb7882b", nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
