//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package fake

import (
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"strings"
	"sync"
)

// ServerFactory is a fake server for instances of the armrecoveryservicesbackup.ClientFactory type.
type ServerFactory struct {
	BMSPrepareDataMoveOperationResultServer          BMSPrepareDataMoveOperationResultServer
	BackupEnginesServer                              BackupEnginesServer
	BackupJobsServer                                 BackupJobsServer
	BackupOperationResultsServer                     BackupOperationResultsServer
	BackupOperationStatusesServer                    BackupOperationStatusesServer
	BackupPoliciesServer                             BackupPoliciesServer
	BackupProtectableItemsServer                     BackupProtectableItemsServer
	BackupProtectedItemsServer                       BackupProtectedItemsServer
	BackupProtectionContainersServer                 BackupProtectionContainersServer
	BackupProtectionIntentServer                     BackupProtectionIntentServer
	BackupResourceEncryptionConfigsServer            BackupResourceEncryptionConfigsServer
	BackupResourceStorageConfigsNonCRRServer         BackupResourceStorageConfigsNonCRRServer
	BackupResourceVaultConfigsServer                 BackupResourceVaultConfigsServer
	BackupStatusServer                               BackupStatusServer
	BackupUsageSummariesServer                       BackupUsageSummariesServer
	BackupWorkloadItemsServer                        BackupWorkloadItemsServer
	BackupsServer                                    BackupsServer
	Server                                           Server
	DeletedProtectionContainersServer                DeletedProtectionContainersServer
	ExportJobsOperationResultsServer                 ExportJobsOperationResultsServer
	FeatureSupportServer                             FeatureSupportServer
	FetchTieringCostServer                           FetchTieringCostServer
	GetTieringCostOperationResultServer              GetTieringCostOperationResultServer
	ItemLevelRecoveryConnectionsServer               ItemLevelRecoveryConnectionsServer
	JobCancellationsServer                           JobCancellationsServer
	JobDetailsServer                                 JobDetailsServer
	JobOperationResultsServer                        JobOperationResultsServer
	JobsServer                                       JobsServer
	OperationServer                                  OperationServer
	OperationsServer                                 OperationsServer
	PrivateEndpointServer                            PrivateEndpointServer
	PrivateEndpointConnectionServer                  PrivateEndpointConnectionServer
	ProtectableContainersServer                      ProtectableContainersServer
	ProtectedItemOperationResultsServer              ProtectedItemOperationResultsServer
	ProtectedItemOperationStatusesServer             ProtectedItemOperationStatusesServer
	ProtectedItemsServer                             ProtectedItemsServer
	ProtectionContainerOperationResultsServer        ProtectionContainerOperationResultsServer
	ProtectionContainerRefreshOperationResultsServer ProtectionContainerRefreshOperationResultsServer
	ProtectionContainersServer                       ProtectionContainersServer
	ProtectionIntentServer                           ProtectionIntentServer
	ProtectionPoliciesServer                         ProtectionPoliciesServer
	ProtectionPolicyOperationResultsServer           ProtectionPolicyOperationResultsServer
	ProtectionPolicyOperationStatusesServer          ProtectionPolicyOperationStatusesServer
	RecoveryPointsServer                             RecoveryPointsServer
	RecoveryPointsRecommendedForMoveServer           RecoveryPointsRecommendedForMoveServer
	ResourceGuardProxiesServer                       ResourceGuardProxiesServer
	ResourceGuardProxyServer                         ResourceGuardProxyServer
	RestoresServer                                   RestoresServer
	SecurityPINsServer                               SecurityPINsServer
	TieringCostOperationStatusServer                 TieringCostOperationStatusServer
	ValidateOperationServer                          ValidateOperationServer
	ValidateOperationResultsServer                   ValidateOperationResultsServer
	ValidateOperationStatusesServer                  ValidateOperationStatusesServer
}

// NewServerFactoryTransport creates a new instance of ServerFactoryTransport with the provided implementation.
// The returned ServerFactoryTransport instance is connected to an instance of armrecoveryservicesbackup.ClientFactory via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewServerFactoryTransport(srv *ServerFactory) *ServerFactoryTransport {
	return &ServerFactoryTransport{
		srv: srv,
	}
}

// ServerFactoryTransport connects instances of armrecoveryservicesbackup.ClientFactory to instances of ServerFactory.
// Don't use this type directly, use NewServerFactoryTransport instead.
type ServerFactoryTransport struct {
	srv                                                *ServerFactory
	trMu                                               sync.Mutex
	trBMSPrepareDataMoveOperationResultServer          *BMSPrepareDataMoveOperationResultServerTransport
	trBackupEnginesServer                              *BackupEnginesServerTransport
	trBackupJobsServer                                 *BackupJobsServerTransport
	trBackupOperationResultsServer                     *BackupOperationResultsServerTransport
	trBackupOperationStatusesServer                    *BackupOperationStatusesServerTransport
	trBackupPoliciesServer                             *BackupPoliciesServerTransport
	trBackupProtectableItemsServer                     *BackupProtectableItemsServerTransport
	trBackupProtectedItemsServer                       *BackupProtectedItemsServerTransport
	trBackupProtectionContainersServer                 *BackupProtectionContainersServerTransport
	trBackupProtectionIntentServer                     *BackupProtectionIntentServerTransport
	trBackupResourceEncryptionConfigsServer            *BackupResourceEncryptionConfigsServerTransport
	trBackupResourceStorageConfigsNonCRRServer         *BackupResourceStorageConfigsNonCRRServerTransport
	trBackupResourceVaultConfigsServer                 *BackupResourceVaultConfigsServerTransport
	trBackupStatusServer                               *BackupStatusServerTransport
	trBackupUsageSummariesServer                       *BackupUsageSummariesServerTransport
	trBackupWorkloadItemsServer                        *BackupWorkloadItemsServerTransport
	trBackupsServer                                    *BackupsServerTransport
	trServer                                           *ServerTransport
	trDeletedProtectionContainersServer                *DeletedProtectionContainersServerTransport
	trExportJobsOperationResultsServer                 *ExportJobsOperationResultsServerTransport
	trFeatureSupportServer                             *FeatureSupportServerTransport
	trFetchTieringCostServer                           *FetchTieringCostServerTransport
	trGetTieringCostOperationResultServer              *GetTieringCostOperationResultServerTransport
	trItemLevelRecoveryConnectionsServer               *ItemLevelRecoveryConnectionsServerTransport
	trJobCancellationsServer                           *JobCancellationsServerTransport
	trJobDetailsServer                                 *JobDetailsServerTransport
	trJobOperationResultsServer                        *JobOperationResultsServerTransport
	trJobsServer                                       *JobsServerTransport
	trOperationServer                                  *OperationServerTransport
	trOperationsServer                                 *OperationsServerTransport
	trPrivateEndpointServer                            *PrivateEndpointServerTransport
	trPrivateEndpointConnectionServer                  *PrivateEndpointConnectionServerTransport
	trProtectableContainersServer                      *ProtectableContainersServerTransport
	trProtectedItemOperationResultsServer              *ProtectedItemOperationResultsServerTransport
	trProtectedItemOperationStatusesServer             *ProtectedItemOperationStatusesServerTransport
	trProtectedItemsServer                             *ProtectedItemsServerTransport
	trProtectionContainerOperationResultsServer        *ProtectionContainerOperationResultsServerTransport
	trProtectionContainerRefreshOperationResultsServer *ProtectionContainerRefreshOperationResultsServerTransport
	trProtectionContainersServer                       *ProtectionContainersServerTransport
	trProtectionIntentServer                           *ProtectionIntentServerTransport
	trProtectionPoliciesServer                         *ProtectionPoliciesServerTransport
	trProtectionPolicyOperationResultsServer           *ProtectionPolicyOperationResultsServerTransport
	trProtectionPolicyOperationStatusesServer          *ProtectionPolicyOperationStatusesServerTransport
	trRecoveryPointsServer                             *RecoveryPointsServerTransport
	trRecoveryPointsRecommendedForMoveServer           *RecoveryPointsRecommendedForMoveServerTransport
	trResourceGuardProxiesServer                       *ResourceGuardProxiesServerTransport
	trResourceGuardProxyServer                         *ResourceGuardProxyServerTransport
	trRestoresServer                                   *RestoresServerTransport
	trSecurityPINsServer                               *SecurityPINsServerTransport
	trTieringCostOperationStatusServer                 *TieringCostOperationStatusServerTransport
	trValidateOperationServer                          *ValidateOperationServerTransport
	trValidateOperationResultsServer                   *ValidateOperationResultsServerTransport
	trValidateOperationStatusesServer                  *ValidateOperationStatusesServerTransport
}

// Do implements the policy.Transporter interface for ServerFactoryTransport.
func (s *ServerFactoryTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	client := method[:strings.Index(method, ".")]
	var resp *http.Response
	var err error

	switch client {
	case "BMSPrepareDataMoveOperationResultClient":
		initServer(s, &s.trBMSPrepareDataMoveOperationResultServer, func() *BMSPrepareDataMoveOperationResultServerTransport {
			return NewBMSPrepareDataMoveOperationResultServerTransport(&s.srv.BMSPrepareDataMoveOperationResultServer)
		})
		resp, err = s.trBMSPrepareDataMoveOperationResultServer.Do(req)
	case "BackupEnginesClient":
		initServer(s, &s.trBackupEnginesServer, func() *BackupEnginesServerTransport {
			return NewBackupEnginesServerTransport(&s.srv.BackupEnginesServer)
		})
		resp, err = s.trBackupEnginesServer.Do(req)
	case "BackupJobsClient":
		initServer(s, &s.trBackupJobsServer, func() *BackupJobsServerTransport { return NewBackupJobsServerTransport(&s.srv.BackupJobsServer) })
		resp, err = s.trBackupJobsServer.Do(req)
	case "BackupOperationResultsClient":
		initServer(s, &s.trBackupOperationResultsServer, func() *BackupOperationResultsServerTransport {
			return NewBackupOperationResultsServerTransport(&s.srv.BackupOperationResultsServer)
		})
		resp, err = s.trBackupOperationResultsServer.Do(req)
	case "BackupOperationStatusesClient":
		initServer(s, &s.trBackupOperationStatusesServer, func() *BackupOperationStatusesServerTransport {
			return NewBackupOperationStatusesServerTransport(&s.srv.BackupOperationStatusesServer)
		})
		resp, err = s.trBackupOperationStatusesServer.Do(req)
	case "BackupPoliciesClient":
		initServer(s, &s.trBackupPoliciesServer, func() *BackupPoliciesServerTransport {
			return NewBackupPoliciesServerTransport(&s.srv.BackupPoliciesServer)
		})
		resp, err = s.trBackupPoliciesServer.Do(req)
	case "BackupProtectableItemsClient":
		initServer(s, &s.trBackupProtectableItemsServer, func() *BackupProtectableItemsServerTransport {
			return NewBackupProtectableItemsServerTransport(&s.srv.BackupProtectableItemsServer)
		})
		resp, err = s.trBackupProtectableItemsServer.Do(req)
	case "BackupProtectedItemsClient":
		initServer(s, &s.trBackupProtectedItemsServer, func() *BackupProtectedItemsServerTransport {
			return NewBackupProtectedItemsServerTransport(&s.srv.BackupProtectedItemsServer)
		})
		resp, err = s.trBackupProtectedItemsServer.Do(req)
	case "BackupProtectionContainersClient":
		initServer(s, &s.trBackupProtectionContainersServer, func() *BackupProtectionContainersServerTransport {
			return NewBackupProtectionContainersServerTransport(&s.srv.BackupProtectionContainersServer)
		})
		resp, err = s.trBackupProtectionContainersServer.Do(req)
	case "BackupProtectionIntentClient":
		initServer(s, &s.trBackupProtectionIntentServer, func() *BackupProtectionIntentServerTransport {
			return NewBackupProtectionIntentServerTransport(&s.srv.BackupProtectionIntentServer)
		})
		resp, err = s.trBackupProtectionIntentServer.Do(req)
	case "BackupResourceEncryptionConfigsClient":
		initServer(s, &s.trBackupResourceEncryptionConfigsServer, func() *BackupResourceEncryptionConfigsServerTransport {
			return NewBackupResourceEncryptionConfigsServerTransport(&s.srv.BackupResourceEncryptionConfigsServer)
		})
		resp, err = s.trBackupResourceEncryptionConfigsServer.Do(req)
	case "BackupResourceStorageConfigsNonCRRClient":
		initServer(s, &s.trBackupResourceStorageConfigsNonCRRServer, func() *BackupResourceStorageConfigsNonCRRServerTransport {
			return NewBackupResourceStorageConfigsNonCRRServerTransport(&s.srv.BackupResourceStorageConfigsNonCRRServer)
		})
		resp, err = s.trBackupResourceStorageConfigsNonCRRServer.Do(req)
	case "BackupResourceVaultConfigsClient":
		initServer(s, &s.trBackupResourceVaultConfigsServer, func() *BackupResourceVaultConfigsServerTransport {
			return NewBackupResourceVaultConfigsServerTransport(&s.srv.BackupResourceVaultConfigsServer)
		})
		resp, err = s.trBackupResourceVaultConfigsServer.Do(req)
	case "BackupStatusClient":
		initServer(s, &s.trBackupStatusServer, func() *BackupStatusServerTransport { return NewBackupStatusServerTransport(&s.srv.BackupStatusServer) })
		resp, err = s.trBackupStatusServer.Do(req)
	case "BackupUsageSummariesClient":
		initServer(s, &s.trBackupUsageSummariesServer, func() *BackupUsageSummariesServerTransport {
			return NewBackupUsageSummariesServerTransport(&s.srv.BackupUsageSummariesServer)
		})
		resp, err = s.trBackupUsageSummariesServer.Do(req)
	case "BackupWorkloadItemsClient":
		initServer(s, &s.trBackupWorkloadItemsServer, func() *BackupWorkloadItemsServerTransport {
			return NewBackupWorkloadItemsServerTransport(&s.srv.BackupWorkloadItemsServer)
		})
		resp, err = s.trBackupWorkloadItemsServer.Do(req)
	case "BackupsClient":
		initServer(s, &s.trBackupsServer, func() *BackupsServerTransport { return NewBackupsServerTransport(&s.srv.BackupsServer) })
		resp, err = s.trBackupsServer.Do(req)
	case "Client":
		initServer(s, &s.trServer, func() *ServerTransport { return NewServerTransport(&s.srv.Server) })
		resp, err = s.trServer.Do(req)
	case "DeletedProtectionContainersClient":
		initServer(s, &s.trDeletedProtectionContainersServer, func() *DeletedProtectionContainersServerTransport {
			return NewDeletedProtectionContainersServerTransport(&s.srv.DeletedProtectionContainersServer)
		})
		resp, err = s.trDeletedProtectionContainersServer.Do(req)
	case "ExportJobsOperationResultsClient":
		initServer(s, &s.trExportJobsOperationResultsServer, func() *ExportJobsOperationResultsServerTransport {
			return NewExportJobsOperationResultsServerTransport(&s.srv.ExportJobsOperationResultsServer)
		})
		resp, err = s.trExportJobsOperationResultsServer.Do(req)
	case "FeatureSupportClient":
		initServer(s, &s.trFeatureSupportServer, func() *FeatureSupportServerTransport {
			return NewFeatureSupportServerTransport(&s.srv.FeatureSupportServer)
		})
		resp, err = s.trFeatureSupportServer.Do(req)
	case "FetchTieringCostClient":
		initServer(s, &s.trFetchTieringCostServer, func() *FetchTieringCostServerTransport {
			return NewFetchTieringCostServerTransport(&s.srv.FetchTieringCostServer)
		})
		resp, err = s.trFetchTieringCostServer.Do(req)
	case "GetTieringCostOperationResultClient":
		initServer(s, &s.trGetTieringCostOperationResultServer, func() *GetTieringCostOperationResultServerTransport {
			return NewGetTieringCostOperationResultServerTransport(&s.srv.GetTieringCostOperationResultServer)
		})
		resp, err = s.trGetTieringCostOperationResultServer.Do(req)
	case "ItemLevelRecoveryConnectionsClient":
		initServer(s, &s.trItemLevelRecoveryConnectionsServer, func() *ItemLevelRecoveryConnectionsServerTransport {
			return NewItemLevelRecoveryConnectionsServerTransport(&s.srv.ItemLevelRecoveryConnectionsServer)
		})
		resp, err = s.trItemLevelRecoveryConnectionsServer.Do(req)
	case "JobCancellationsClient":
		initServer(s, &s.trJobCancellationsServer, func() *JobCancellationsServerTransport {
			return NewJobCancellationsServerTransport(&s.srv.JobCancellationsServer)
		})
		resp, err = s.trJobCancellationsServer.Do(req)
	case "JobDetailsClient":
		initServer(s, &s.trJobDetailsServer, func() *JobDetailsServerTransport { return NewJobDetailsServerTransport(&s.srv.JobDetailsServer) })
		resp, err = s.trJobDetailsServer.Do(req)
	case "JobOperationResultsClient":
		initServer(s, &s.trJobOperationResultsServer, func() *JobOperationResultsServerTransport {
			return NewJobOperationResultsServerTransport(&s.srv.JobOperationResultsServer)
		})
		resp, err = s.trJobOperationResultsServer.Do(req)
	case "JobsClient":
		initServer(s, &s.trJobsServer, func() *JobsServerTransport { return NewJobsServerTransport(&s.srv.JobsServer) })
		resp, err = s.trJobsServer.Do(req)
	case "OperationClient":
		initServer(s, &s.trOperationServer, func() *OperationServerTransport { return NewOperationServerTransport(&s.srv.OperationServer) })
		resp, err = s.trOperationServer.Do(req)
	case "OperationsClient":
		initServer(s, &s.trOperationsServer, func() *OperationsServerTransport { return NewOperationsServerTransport(&s.srv.OperationsServer) })
		resp, err = s.trOperationsServer.Do(req)
	case "PrivateEndpointClient":
		initServer(s, &s.trPrivateEndpointServer, func() *PrivateEndpointServerTransport {
			return NewPrivateEndpointServerTransport(&s.srv.PrivateEndpointServer)
		})
		resp, err = s.trPrivateEndpointServer.Do(req)
	case "PrivateEndpointConnectionClient":
		initServer(s, &s.trPrivateEndpointConnectionServer, func() *PrivateEndpointConnectionServerTransport {
			return NewPrivateEndpointConnectionServerTransport(&s.srv.PrivateEndpointConnectionServer)
		})
		resp, err = s.trPrivateEndpointConnectionServer.Do(req)
	case "ProtectableContainersClient":
		initServer(s, &s.trProtectableContainersServer, func() *ProtectableContainersServerTransport {
			return NewProtectableContainersServerTransport(&s.srv.ProtectableContainersServer)
		})
		resp, err = s.trProtectableContainersServer.Do(req)
	case "ProtectedItemOperationResultsClient":
		initServer(s, &s.trProtectedItemOperationResultsServer, func() *ProtectedItemOperationResultsServerTransport {
			return NewProtectedItemOperationResultsServerTransport(&s.srv.ProtectedItemOperationResultsServer)
		})
		resp, err = s.trProtectedItemOperationResultsServer.Do(req)
	case "ProtectedItemOperationStatusesClient":
		initServer(s, &s.trProtectedItemOperationStatusesServer, func() *ProtectedItemOperationStatusesServerTransport {
			return NewProtectedItemOperationStatusesServerTransport(&s.srv.ProtectedItemOperationStatusesServer)
		})
		resp, err = s.trProtectedItemOperationStatusesServer.Do(req)
	case "ProtectedItemsClient":
		initServer(s, &s.trProtectedItemsServer, func() *ProtectedItemsServerTransport {
			return NewProtectedItemsServerTransport(&s.srv.ProtectedItemsServer)
		})
		resp, err = s.trProtectedItemsServer.Do(req)
	case "ProtectionContainerOperationResultsClient":
		initServer(s, &s.trProtectionContainerOperationResultsServer, func() *ProtectionContainerOperationResultsServerTransport {
			return NewProtectionContainerOperationResultsServerTransport(&s.srv.ProtectionContainerOperationResultsServer)
		})
		resp, err = s.trProtectionContainerOperationResultsServer.Do(req)
	case "ProtectionContainerRefreshOperationResultsClient":
		initServer(s, &s.trProtectionContainerRefreshOperationResultsServer, func() *ProtectionContainerRefreshOperationResultsServerTransport {
			return NewProtectionContainerRefreshOperationResultsServerTransport(&s.srv.ProtectionContainerRefreshOperationResultsServer)
		})
		resp, err = s.trProtectionContainerRefreshOperationResultsServer.Do(req)
	case "ProtectionContainersClient":
		initServer(s, &s.trProtectionContainersServer, func() *ProtectionContainersServerTransport {
			return NewProtectionContainersServerTransport(&s.srv.ProtectionContainersServer)
		})
		resp, err = s.trProtectionContainersServer.Do(req)
	case "ProtectionIntentClient":
		initServer(s, &s.trProtectionIntentServer, func() *ProtectionIntentServerTransport {
			return NewProtectionIntentServerTransport(&s.srv.ProtectionIntentServer)
		})
		resp, err = s.trProtectionIntentServer.Do(req)
	case "ProtectionPoliciesClient":
		initServer(s, &s.trProtectionPoliciesServer, func() *ProtectionPoliciesServerTransport {
			return NewProtectionPoliciesServerTransport(&s.srv.ProtectionPoliciesServer)
		})
		resp, err = s.trProtectionPoliciesServer.Do(req)
	case "ProtectionPolicyOperationResultsClient":
		initServer(s, &s.trProtectionPolicyOperationResultsServer, func() *ProtectionPolicyOperationResultsServerTransport {
			return NewProtectionPolicyOperationResultsServerTransport(&s.srv.ProtectionPolicyOperationResultsServer)
		})
		resp, err = s.trProtectionPolicyOperationResultsServer.Do(req)
	case "ProtectionPolicyOperationStatusesClient":
		initServer(s, &s.trProtectionPolicyOperationStatusesServer, func() *ProtectionPolicyOperationStatusesServerTransport {
			return NewProtectionPolicyOperationStatusesServerTransport(&s.srv.ProtectionPolicyOperationStatusesServer)
		})
		resp, err = s.trProtectionPolicyOperationStatusesServer.Do(req)
	case "RecoveryPointsClient":
		initServer(s, &s.trRecoveryPointsServer, func() *RecoveryPointsServerTransport {
			return NewRecoveryPointsServerTransport(&s.srv.RecoveryPointsServer)
		})
		resp, err = s.trRecoveryPointsServer.Do(req)
	case "RecoveryPointsRecommendedForMoveClient":
		initServer(s, &s.trRecoveryPointsRecommendedForMoveServer, func() *RecoveryPointsRecommendedForMoveServerTransport {
			return NewRecoveryPointsRecommendedForMoveServerTransport(&s.srv.RecoveryPointsRecommendedForMoveServer)
		})
		resp, err = s.trRecoveryPointsRecommendedForMoveServer.Do(req)
	case "ResourceGuardProxiesClient":
		initServer(s, &s.trResourceGuardProxiesServer, func() *ResourceGuardProxiesServerTransport {
			return NewResourceGuardProxiesServerTransport(&s.srv.ResourceGuardProxiesServer)
		})
		resp, err = s.trResourceGuardProxiesServer.Do(req)
	case "ResourceGuardProxyClient":
		initServer(s, &s.trResourceGuardProxyServer, func() *ResourceGuardProxyServerTransport {
			return NewResourceGuardProxyServerTransport(&s.srv.ResourceGuardProxyServer)
		})
		resp, err = s.trResourceGuardProxyServer.Do(req)
	case "RestoresClient":
		initServer(s, &s.trRestoresServer, func() *RestoresServerTransport { return NewRestoresServerTransport(&s.srv.RestoresServer) })
		resp, err = s.trRestoresServer.Do(req)
	case "SecurityPINsClient":
		initServer(s, &s.trSecurityPINsServer, func() *SecurityPINsServerTransport { return NewSecurityPINsServerTransport(&s.srv.SecurityPINsServer) })
		resp, err = s.trSecurityPINsServer.Do(req)
	case "TieringCostOperationStatusClient":
		initServer(s, &s.trTieringCostOperationStatusServer, func() *TieringCostOperationStatusServerTransport {
			return NewTieringCostOperationStatusServerTransport(&s.srv.TieringCostOperationStatusServer)
		})
		resp, err = s.trTieringCostOperationStatusServer.Do(req)
	case "ValidateOperationClient":
		initServer(s, &s.trValidateOperationServer, func() *ValidateOperationServerTransport {
			return NewValidateOperationServerTransport(&s.srv.ValidateOperationServer)
		})
		resp, err = s.trValidateOperationServer.Do(req)
	case "ValidateOperationResultsClient":
		initServer(s, &s.trValidateOperationResultsServer, func() *ValidateOperationResultsServerTransport {
			return NewValidateOperationResultsServerTransport(&s.srv.ValidateOperationResultsServer)
		})
		resp, err = s.trValidateOperationResultsServer.Do(req)
	case "ValidateOperationStatusesClient":
		initServer(s, &s.trValidateOperationStatusesServer, func() *ValidateOperationStatusesServerTransport {
			return NewValidateOperationStatusesServerTransport(&s.srv.ValidateOperationStatusesServer)
		})
		resp, err = s.trValidateOperationStatusesServer.Do(req)
	default:
		err = fmt.Errorf("unhandled client %s", client)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func initServer[T any](s *ServerFactoryTransport, dst **T, src func() *T) {
	s.trMu.Lock()
	if *dst == nil {
		*dst = src()
	}
	s.trMu.Unlock()
}
