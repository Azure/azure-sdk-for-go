//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armkubernetesconfiguration

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
)

// MarshalJSON implements the json.Marshaller interface for type ComplianceStatus.
func (c ComplianceStatus) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "complianceState", c.ComplianceState)
	populateTimeRFC3339(objectMap, "lastConfigApplied", c.LastConfigApplied)
	populate(objectMap, "message", c.Message)
	populate(objectMap, "messageLevel", c.MessageLevel)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ComplianceStatus.
func (c *ComplianceStatus) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", c, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "complianceState":
			err = unpopulate(val, "ComplianceState", &c.ComplianceState)
			delete(rawMsg, key)
		case "lastConfigApplied":
			err = unpopulateTimeRFC3339(val, "LastConfigApplied", &c.LastConfigApplied)
			delete(rawMsg, key)
		case "message":
			err = unpopulate(val, "Message", &c.Message)
			delete(rawMsg, key)
		case "messageLevel":
			err = unpopulate(val, "MessageLevel", &c.MessageLevel)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", c, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ErrorDetail.
func (e ErrorDetail) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "additionalInfo", e.AdditionalInfo)
	populate(objectMap, "code", e.Code)
	populate(objectMap, "details", e.Details)
	populate(objectMap, "message", e.Message)
	populate(objectMap, "target", e.Target)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ExtensionProperties.
func (e ExtensionProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "aksAssignedIdentity", e.AksAssignedIdentity)
	populate(objectMap, "autoUpgradeMinorVersion", e.AutoUpgradeMinorVersion)
	populate(objectMap, "configurationProtectedSettings", e.ConfigurationProtectedSettings)
	populate(objectMap, "configurationSettings", e.ConfigurationSettings)
	populate(objectMap, "customLocationSettings", e.CustomLocationSettings)
	populate(objectMap, "errorInfo", e.ErrorInfo)
	populate(objectMap, "extensionType", e.ExtensionType)
	populate(objectMap, "installedVersion", e.InstalledVersion)
	populate(objectMap, "packageUri", e.PackageURI)
	populate(objectMap, "provisioningState", e.ProvisioningState)
	populate(objectMap, "releaseTrain", e.ReleaseTrain)
	populate(objectMap, "scope", e.Scope)
	populate(objectMap, "statuses", e.Statuses)
	populate(objectMap, "version", e.Version)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type FluxConfigurationPatch.
func (f FluxConfigurationPatch) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "properties", f.Properties)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type FluxConfigurationPatchProperties.
func (f FluxConfigurationPatchProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "bucket", f.Bucket)
	populate(objectMap, "configurationProtectedSettings", f.ConfigurationProtectedSettings)
	populate(objectMap, "gitRepository", f.GitRepository)
	populate(objectMap, "kustomizations", f.Kustomizations)
	populate(objectMap, "sourceKind", f.SourceKind)
	populate(objectMap, "suspend", f.Suspend)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type FluxConfigurationProperties.
func (f FluxConfigurationProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "bucket", f.Bucket)
	populate(objectMap, "complianceState", f.ComplianceState)
	populate(objectMap, "configurationProtectedSettings", f.ConfigurationProtectedSettings)
	populate(objectMap, "errorMessage", f.ErrorMessage)
	populate(objectMap, "gitRepository", f.GitRepository)
	populate(objectMap, "kustomizations", f.Kustomizations)
	populate(objectMap, "namespace", f.Namespace)
	populate(objectMap, "provisioningState", f.ProvisioningState)
	populate(objectMap, "repositoryPublicKey", f.RepositoryPublicKey)
	populate(objectMap, "scope", f.Scope)
	populate(objectMap, "sourceKind", f.SourceKind)
	populate(objectMap, "sourceSyncedCommitId", f.SourceSyncedCommitID)
	populateTimeRFC3339(objectMap, "sourceUpdatedAt", f.SourceUpdatedAt)
	populateTimeRFC3339(objectMap, "statusUpdatedAt", f.StatusUpdatedAt)
	populate(objectMap, "statuses", f.Statuses)
	populate(objectMap, "suspend", f.Suspend)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type FluxConfigurationProperties.
func (f *FluxConfigurationProperties) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", f, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "bucket":
			err = unpopulate(val, "Bucket", &f.Bucket)
			delete(rawMsg, key)
		case "complianceState":
			err = unpopulate(val, "ComplianceState", &f.ComplianceState)
			delete(rawMsg, key)
		case "configurationProtectedSettings":
			err = unpopulate(val, "ConfigurationProtectedSettings", &f.ConfigurationProtectedSettings)
			delete(rawMsg, key)
		case "errorMessage":
			err = unpopulate(val, "ErrorMessage", &f.ErrorMessage)
			delete(rawMsg, key)
		case "gitRepository":
			err = unpopulate(val, "GitRepository", &f.GitRepository)
			delete(rawMsg, key)
		case "kustomizations":
			err = unpopulate(val, "Kustomizations", &f.Kustomizations)
			delete(rawMsg, key)
		case "namespace":
			err = unpopulate(val, "Namespace", &f.Namespace)
			delete(rawMsg, key)
		case "provisioningState":
			err = unpopulate(val, "ProvisioningState", &f.ProvisioningState)
			delete(rawMsg, key)
		case "repositoryPublicKey":
			err = unpopulate(val, "RepositoryPublicKey", &f.RepositoryPublicKey)
			delete(rawMsg, key)
		case "scope":
			err = unpopulate(val, "Scope", &f.Scope)
			delete(rawMsg, key)
		case "sourceKind":
			err = unpopulate(val, "SourceKind", &f.SourceKind)
			delete(rawMsg, key)
		case "sourceSyncedCommitId":
			err = unpopulate(val, "SourceSyncedCommitID", &f.SourceSyncedCommitID)
			delete(rawMsg, key)
		case "sourceUpdatedAt":
			err = unpopulateTimeRFC3339(val, "SourceUpdatedAt", &f.SourceUpdatedAt)
			delete(rawMsg, key)
		case "statusUpdatedAt":
			err = unpopulateTimeRFC3339(val, "StatusUpdatedAt", &f.StatusUpdatedAt)
			delete(rawMsg, key)
		case "statuses":
			err = unpopulate(val, "Statuses", &f.Statuses)
			delete(rawMsg, key)
		case "suspend":
			err = unpopulate(val, "Suspend", &f.Suspend)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", f, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type KustomizationDefinition.
func (k KustomizationDefinition) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "dependsOn", k.DependsOn)
	populate(objectMap, "force", k.Force)
	populate(objectMap, "name", k.Name)
	populate(objectMap, "path", k.Path)
	populate(objectMap, "prune", k.Prune)
	populate(objectMap, "retryIntervalInSeconds", k.RetryIntervalInSeconds)
	populate(objectMap, "syncIntervalInSeconds", k.SyncIntervalInSeconds)
	populate(objectMap, "timeoutInSeconds", k.TimeoutInSeconds)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type KustomizationPatchDefinition.
func (k KustomizationPatchDefinition) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "dependsOn", k.DependsOn)
	populate(objectMap, "force", k.Force)
	populate(objectMap, "path", k.Path)
	populate(objectMap, "prune", k.Prune)
	populate(objectMap, "retryIntervalInSeconds", k.RetryIntervalInSeconds)
	populate(objectMap, "syncIntervalInSeconds", k.SyncIntervalInSeconds)
	populate(objectMap, "timeoutInSeconds", k.TimeoutInSeconds)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ObjectStatusConditionDefinition.
func (o ObjectStatusConditionDefinition) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populateTimeRFC3339(objectMap, "lastTransitionTime", o.LastTransitionTime)
	populate(objectMap, "message", o.Message)
	populate(objectMap, "reason", o.Reason)
	populate(objectMap, "status", o.Status)
	populate(objectMap, "type", o.Type)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ObjectStatusConditionDefinition.
func (o *ObjectStatusConditionDefinition) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", o, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "lastTransitionTime":
			err = unpopulateTimeRFC3339(val, "LastTransitionTime", &o.LastTransitionTime)
			delete(rawMsg, key)
		case "message":
			err = unpopulate(val, "Message", &o.Message)
			delete(rawMsg, key)
		case "reason":
			err = unpopulate(val, "Reason", &o.Reason)
			delete(rawMsg, key)
		case "status":
			err = unpopulate(val, "Status", &o.Status)
			delete(rawMsg, key)
		case "type":
			err = unpopulate(val, "Type", &o.Type)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", o, err)
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ObjectStatusDefinition.
func (o ObjectStatusDefinition) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "appliedBy", o.AppliedBy)
	populate(objectMap, "complianceState", o.ComplianceState)
	populate(objectMap, "helmReleaseProperties", o.HelmReleaseProperties)
	populate(objectMap, "kind", o.Kind)
	populate(objectMap, "name", o.Name)
	populate(objectMap, "namespace", o.Namespace)
	populate(objectMap, "statusConditions", o.StatusConditions)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PatchExtension.
func (p PatchExtension) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "properties", p.Properties)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PatchExtensionProperties.
func (p PatchExtensionProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "autoUpgradeMinorVersion", p.AutoUpgradeMinorVersion)
	populate(objectMap, "configurationProtectedSettings", p.ConfigurationProtectedSettings)
	populate(objectMap, "configurationSettings", p.ConfigurationSettings)
	populate(objectMap, "releaseTrain", p.ReleaseTrain)
	populate(objectMap, "version", p.Version)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type SourceControlConfigurationProperties.
func (s SourceControlConfigurationProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "complianceStatus", s.ComplianceStatus)
	populate(objectMap, "configurationProtectedSettings", s.ConfigurationProtectedSettings)
	populate(objectMap, "enableHelmOperator", s.EnableHelmOperator)
	populate(objectMap, "helmOperatorProperties", s.HelmOperatorProperties)
	populate(objectMap, "operatorInstanceName", s.OperatorInstanceName)
	populate(objectMap, "operatorNamespace", s.OperatorNamespace)
	populate(objectMap, "operatorParams", s.OperatorParams)
	populate(objectMap, "operatorScope", s.OperatorScope)
	populate(objectMap, "operatorType", s.OperatorType)
	populate(objectMap, "provisioningState", s.ProvisioningState)
	populate(objectMap, "repositoryPublicKey", s.RepositoryPublicKey)
	populate(objectMap, "repositoryUrl", s.RepositoryURL)
	populate(objectMap, "sshKnownHostsContents", s.SSHKnownHostsContents)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type SystemData.
func (s SystemData) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populateTimeRFC3339(objectMap, "createdAt", s.CreatedAt)
	populate(objectMap, "createdBy", s.CreatedBy)
	populate(objectMap, "createdByType", s.CreatedByType)
	populateTimeRFC3339(objectMap, "lastModifiedAt", s.LastModifiedAt)
	populate(objectMap, "lastModifiedBy", s.LastModifiedBy)
	populate(objectMap, "lastModifiedByType", s.LastModifiedByType)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type SystemData.
func (s *SystemData) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return fmt.Errorf("unmarshalling type %T: %v", s, err)
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "createdAt":
			err = unpopulateTimeRFC3339(val, "CreatedAt", &s.CreatedAt)
			delete(rawMsg, key)
		case "createdBy":
			err = unpopulate(val, "CreatedBy", &s.CreatedBy)
			delete(rawMsg, key)
		case "createdByType":
			err = unpopulate(val, "CreatedByType", &s.CreatedByType)
			delete(rawMsg, key)
		case "lastModifiedAt":
			err = unpopulateTimeRFC3339(val, "LastModifiedAt", &s.LastModifiedAt)
			delete(rawMsg, key)
		case "lastModifiedBy":
			err = unpopulate(val, "LastModifiedBy", &s.LastModifiedBy)
			delete(rawMsg, key)
		case "lastModifiedByType":
			err = unpopulate(val, "LastModifiedByType", &s.LastModifiedByType)
			delete(rawMsg, key)
		}
		if err != nil {
			return fmt.Errorf("unmarshalling type %T: %v", s, err)
		}
	}
	return nil
}

func populate(m map[string]interface{}, k string, v interface{}) {
	if v == nil {
		return
	} else if azcore.IsNullValue(v) {
		m[k] = nil
	} else if !reflect.ValueOf(v).IsNil() {
		m[k] = v
	}
}

func unpopulate(data json.RawMessage, fn string, v interface{}) error {
	if data == nil {
		return nil
	}
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("struct field %s: %v", fn, err)
	}
	return nil
}
