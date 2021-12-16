//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdataprotection_test

import (
	"context"
	"log"

	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dataprotection/armdataprotection"
)

// x-ms-original-file: specification/dataprotection/resource-manager/Microsoft.DataProtection/stable/2021-07-01/examples/BackupInstanceOperations/ListBackupInstances.json
func ExampleBackupInstancesClient_List() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdataprotection.NewBackupInstancesClient("<subscription-id>", cred, nil)
	pager := client.List("<vault-name>",
		"<resource-group-name>",
		nil)
	for pager.NextPage(ctx) {
		if err := pager.Err(); err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range pager.PageResponse().Value {
			log.Printf("BackupInstanceResource.ID: %s\n", *v.ID)
		}
	}
}

// x-ms-original-file: specification/dataprotection/resource-manager/Microsoft.DataProtection/stable/2021-07-01/examples/BackupInstanceOperations/GetBackupInstance.json
func ExampleBackupInstancesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdataprotection.NewBackupInstancesClient("<subscription-id>", cred, nil)
	res, err := client.Get(ctx,
		"<vault-name>",
		"<resource-group-name>",
		"<backup-instance-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("BackupInstanceResource.ID: %s\n", *res.ID)
}

// x-ms-original-file: specification/dataprotection/resource-manager/Microsoft.DataProtection/stable/2021-07-01/examples/BackupInstanceOperations/PutBackupInstance.json
func ExampleBackupInstancesClient_BeginCreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdataprotection.NewBackupInstancesClient("<subscription-id>", cred, nil)
	poller, err := client.BeginCreateOrUpdate(ctx,
		"<vault-name>",
		"<resource-group-name>",
		"<backup-instance-name>",
		armdataprotection.BackupInstanceResource{
			Properties: &armdataprotection.BackupInstance{
				DataSourceInfo: &armdataprotection.Datasource{
					DatasourceType:   to.StringPtr("<datasource-type>"),
					ObjectType:       to.StringPtr("<object-type>"),
					ResourceID:       to.StringPtr("<resource-id>"),
					ResourceLocation: to.StringPtr("<resource-location>"),
					ResourceName:     to.StringPtr("<resource-name>"),
					ResourceType:     to.StringPtr("<resource-type>"),
					ResourceURI:      to.StringPtr("<resource-uri>"),
				},
				DataSourceSetInfo: &armdataprotection.DatasourceSet{
					DatasourceType:   to.StringPtr("<datasource-type>"),
					ObjectType:       to.StringPtr("<object-type>"),
					ResourceID:       to.StringPtr("<resource-id>"),
					ResourceLocation: to.StringPtr("<resource-location>"),
					ResourceName:     to.StringPtr("<resource-name>"),
					ResourceType:     to.StringPtr("<resource-type>"),
					ResourceURI:      to.StringPtr("<resource-uri>"),
				},
				DatasourceAuthCredentials: &armdataprotection.SecretStoreBasedAuthCredentials{
					AuthCredentials: armdataprotection.AuthCredentials{
						ObjectType: to.StringPtr("<object-type>"),
					},
					SecretStoreResource: &armdataprotection.SecretStoreResource{
						SecretStoreType: armdataprotection.SecretStoreTypeAzureKeyVault.ToPtr(),
						URI:             to.StringPtr("<uri>"),
					},
				},
				FriendlyName: to.StringPtr("<friendly-name>"),
				ObjectType:   to.StringPtr("<object-type>"),
				PolicyInfo: &armdataprotection.PolicyInfo{
					PolicyID: to.StringPtr("<policy-id>"),
					PolicyParameters: &armdataprotection.PolicyParameters{
						DataStoreParametersList: []armdataprotection.DataStoreParametersClassification{
							&armdataprotection.AzureOperationalStoreParameters{
								DataStoreParameters: armdataprotection.DataStoreParameters{
									DataStoreType: armdataprotection.DataStoreTypesOperationalStore.ToPtr(),
									ObjectType:    to.StringPtr("<object-type>"),
								},
								ResourceGroupID: to.StringPtr("<resource-group-id>"),
							}},
					},
				},
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("BackupInstanceResource.ID: %s\n", *res.ID)
}

// x-ms-original-file: specification/dataprotection/resource-manager/Microsoft.DataProtection/stable/2021-07-01/examples/BackupInstanceOperations/DeleteBackupInstance.json
func ExampleBackupInstancesClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdataprotection.NewBackupInstancesClient("<subscription-id>", cred, nil)
	poller, err := client.BeginDelete(ctx,
		"<vault-name>",
		"<resource-group-name>",
		"<backup-instance-name>",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/dataprotection/resource-manager/Microsoft.DataProtection/stable/2021-07-01/examples/BackupInstanceOperations/TriggerBackup.json
func ExampleBackupInstancesClient_BeginAdhocBackup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdataprotection.NewBackupInstancesClient("<subscription-id>", cred, nil)
	poller, err := client.BeginAdhocBackup(ctx,
		"<vault-name>",
		"<resource-group-name>",
		"<backup-instance-name>",
		armdataprotection.TriggerBackupRequest{
			BackupRuleOptions: &armdataprotection.AdHocBackupRuleOptions{
				RuleName: to.StringPtr("<rule-name>"),
				TriggerOption: &armdataprotection.AdhocBackupTriggerOption{
					RetentionTagOverride: to.StringPtr("<retention-tag-override>"),
				},
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/dataprotection/resource-manager/Microsoft.DataProtection/stable/2021-07-01/examples/BackupInstanceOperations/ValidateForBackup.json
func ExampleBackupInstancesClient_BeginValidateForBackup() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdataprotection.NewBackupInstancesClient("<subscription-id>", cred, nil)
	poller, err := client.BeginValidateForBackup(ctx,
		"<vault-name>",
		"<resource-group-name>",
		armdataprotection.ValidateForBackupRequest{
			BackupInstance: &armdataprotection.BackupInstance{
				DataSourceInfo: &armdataprotection.Datasource{
					DatasourceType:   to.StringPtr("<datasource-type>"),
					ObjectType:       to.StringPtr("<object-type>"),
					ResourceID:       to.StringPtr("<resource-id>"),
					ResourceLocation: to.StringPtr("<resource-location>"),
					ResourceName:     to.StringPtr("<resource-name>"),
					ResourceType:     to.StringPtr("<resource-type>"),
					ResourceURI:      to.StringPtr("<resource-uri>"),
				},
				DataSourceSetInfo: &armdataprotection.DatasourceSet{
					DatasourceType:   to.StringPtr("<datasource-type>"),
					ObjectType:       to.StringPtr("<object-type>"),
					ResourceID:       to.StringPtr("<resource-id>"),
					ResourceLocation: to.StringPtr("<resource-location>"),
					ResourceName:     to.StringPtr("<resource-name>"),
					ResourceType:     to.StringPtr("<resource-type>"),
					ResourceURI:      to.StringPtr("<resource-uri>"),
				},
				DatasourceAuthCredentials: &armdataprotection.SecretStoreBasedAuthCredentials{
					AuthCredentials: armdataprotection.AuthCredentials{
						ObjectType: to.StringPtr("<object-type>"),
					},
					SecretStoreResource: &armdataprotection.SecretStoreResource{
						SecretStoreType: armdataprotection.SecretStoreTypeAzureKeyVault.ToPtr(),
						URI:             to.StringPtr("<uri>"),
					},
				},
				FriendlyName: to.StringPtr("<friendly-name>"),
				ObjectType:   to.StringPtr("<object-type>"),
				PolicyInfo: &armdataprotection.PolicyInfo{
					PolicyID: to.StringPtr("<policy-id>"),
				},
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/dataprotection/resource-manager/Microsoft.DataProtection/stable/2021-07-01/examples/BackupInstanceOperations/TriggerRehydrate.json
func ExampleBackupInstancesClient_BeginTriggerRehydrate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdataprotection.NewBackupInstancesClient("<subscription-id>", cred, nil)
	poller, err := client.BeginTriggerRehydrate(ctx,
		"<resource-group-name>",
		"<vault-name>",
		"<backup-instance-name>",
		armdataprotection.AzureBackupRehydrationRequest{
			RecoveryPointID:              to.StringPtr("<recovery-point-id>"),
			RehydrationPriority:          armdataprotection.RehydrationPriorityHigh.ToPtr(),
			RehydrationRetentionDuration: to.StringPtr("<rehydration-retention-duration>"),
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/dataprotection/resource-manager/Microsoft.DataProtection/stable/2021-07-01/examples/BackupInstanceOperations/TriggerRestore.json
func ExampleBackupInstancesClient_BeginTriggerRestore() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdataprotection.NewBackupInstancesClient("<subscription-id>", cred, nil)
	poller, err := client.BeginTriggerRestore(ctx,
		"<vault-name>",
		"<resource-group-name>",
		"<backup-instance-name>",
		&armdataprotection.AzureBackupRecoveryPointBasedRestoreRequest{
			AzureBackupRestoreRequest: armdataprotection.AzureBackupRestoreRequest{
				ObjectType: to.StringPtr("<object-type>"),
				RestoreTargetInfo: &armdataprotection.RestoreTargetInfo{
					RestoreTargetInfoBase: armdataprotection.RestoreTargetInfoBase{
						ObjectType:      to.StringPtr("<object-type>"),
						RecoveryOption:  armdataprotection.RecoveryOptionFailIfExists.ToPtr(),
						RestoreLocation: to.StringPtr("<restore-location>"),
					},
					DatasourceAuthCredentials: &armdataprotection.SecretStoreBasedAuthCredentials{
						AuthCredentials: armdataprotection.AuthCredentials{
							ObjectType: to.StringPtr("<object-type>"),
						},
						SecretStoreResource: &armdataprotection.SecretStoreResource{
							SecretStoreType: armdataprotection.SecretStoreTypeAzureKeyVault.ToPtr(),
							URI:             to.StringPtr("<uri>"),
						},
					},
					DatasourceInfo: &armdataprotection.Datasource{
						DatasourceType:   to.StringPtr("<datasource-type>"),
						ObjectType:       to.StringPtr("<object-type>"),
						ResourceID:       to.StringPtr("<resource-id>"),
						ResourceLocation: to.StringPtr("<resource-location>"),
						ResourceName:     to.StringPtr("<resource-name>"),
						ResourceType:     to.StringPtr("<resource-type>"),
						ResourceURI:      to.StringPtr("<resource-uri>"),
					},
					DatasourceSetInfo: &armdataprotection.DatasourceSet{
						DatasourceType:   to.StringPtr("<datasource-type>"),
						ObjectType:       to.StringPtr("<object-type>"),
						ResourceID:       to.StringPtr("<resource-id>"),
						ResourceLocation: to.StringPtr("<resource-location>"),
						ResourceName:     to.StringPtr("<resource-name>"),
						ResourceType:     to.StringPtr("<resource-type>"),
						ResourceURI:      to.StringPtr("<resource-uri>"),
					},
				},
				SourceDataStoreType: armdataprotection.SourceDataStoreTypeVaultStore.ToPtr(),
			},
			RecoveryPointID: to.StringPtr("<recovery-point-id>"),
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}

// x-ms-original-file: specification/dataprotection/resource-manager/Microsoft.DataProtection/stable/2021-07-01/examples/BackupInstanceOperations/ValidateRestore.json
func ExampleBackupInstancesClient_BeginValidateForRestore() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client := armdataprotection.NewBackupInstancesClient("<subscription-id>", cred, nil)
	poller, err := client.BeginValidateForRestore(ctx,
		"<vault-name>",
		"<resource-group-name>",
		"<backup-instance-name>",
		armdataprotection.ValidateRestoreRequestObject{
			RestoreRequestObject: &armdataprotection.AzureBackupRecoveryPointBasedRestoreRequest{
				AzureBackupRestoreRequest: armdataprotection.AzureBackupRestoreRequest{
					ObjectType: to.StringPtr("<object-type>"),
					RestoreTargetInfo: &armdataprotection.RestoreTargetInfo{
						RestoreTargetInfoBase: armdataprotection.RestoreTargetInfoBase{
							ObjectType:      to.StringPtr("<object-type>"),
							RecoveryOption:  armdataprotection.RecoveryOptionFailIfExists.ToPtr(),
							RestoreLocation: to.StringPtr("<restore-location>"),
						},
						DatasourceAuthCredentials: &armdataprotection.SecretStoreBasedAuthCredentials{
							AuthCredentials: armdataprotection.AuthCredentials{
								ObjectType: to.StringPtr("<object-type>"),
							},
							SecretStoreResource: &armdataprotection.SecretStoreResource{
								SecretStoreType: armdataprotection.SecretStoreTypeAzureKeyVault.ToPtr(),
								URI:             to.StringPtr("<uri>"),
							},
						},
						DatasourceInfo: &armdataprotection.Datasource{
							DatasourceType:   to.StringPtr("<datasource-type>"),
							ObjectType:       to.StringPtr("<object-type>"),
							ResourceID:       to.StringPtr("<resource-id>"),
							ResourceLocation: to.StringPtr("<resource-location>"),
							ResourceName:     to.StringPtr("<resource-name>"),
							ResourceType:     to.StringPtr("<resource-type>"),
							ResourceURI:      to.StringPtr("<resource-uri>"),
						},
						DatasourceSetInfo: &armdataprotection.DatasourceSet{
							DatasourceType:   to.StringPtr("<datasource-type>"),
							ObjectType:       to.StringPtr("<object-type>"),
							ResourceID:       to.StringPtr("<resource-id>"),
							ResourceLocation: to.StringPtr("<resource-location>"),
							ResourceName:     to.StringPtr("<resource-name>"),
							ResourceType:     to.StringPtr("<resource-type>"),
							ResourceURI:      to.StringPtr("<resource-uri>"),
						},
					},
					SourceDataStoreType: armdataprotection.SourceDataStoreTypeVaultStore.ToPtr(),
				},
				RecoveryPointID: to.StringPtr("<recovery-point-id>"),
			},
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	_, err = poller.PollUntilDone(ctx, 30*time.Second)
	if err != nil {
		log.Fatal(err)
	}
}
