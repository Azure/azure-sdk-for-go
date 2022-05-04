//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armredis

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
)

// MarshalJSON implements the json.Marshaller interface for type CommonProperties.
func (c CommonProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "enableNonSslPort", c.EnableNonSSLPort)
	populate(objectMap, "minimumTlsVersion", c.MinimumTLSVersion)
	populate(objectMap, "publicNetworkAccess", c.PublicNetworkAccess)
	populate(objectMap, "redisConfiguration", c.RedisConfiguration)
	populate(objectMap, "redisVersion", c.RedisVersion)
	populate(objectMap, "replicasPerMaster", c.ReplicasPerMaster)
	populate(objectMap, "replicasPerPrimary", c.ReplicasPerPrimary)
	populate(objectMap, "shardCount", c.ShardCount)
	populate(objectMap, "tenantSettings", c.TenantSettings)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type CommonPropertiesRedisConfiguration.
func (c CommonPropertiesRedisConfiguration) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "aof-storage-connection-string-0", c.AofStorageConnectionString0)
	populate(objectMap, "aof-storage-connection-string-1", c.AofStorageConnectionString1)
	populate(objectMap, "maxclients", c.Maxclients)
	populate(objectMap, "maxfragmentationmemory-reserved", c.MaxfragmentationmemoryReserved)
	populate(objectMap, "maxmemory-delta", c.MaxmemoryDelta)
	populate(objectMap, "maxmemory-policy", c.MaxmemoryPolicy)
	populate(objectMap, "maxmemory-reserved", c.MaxmemoryReserved)
	populate(objectMap, "preferred-data-archive-auth-method", c.PreferredDataArchiveAuthMethod)
	populate(objectMap, "preferred-data-persistence-auth-method", c.PreferredDataPersistenceAuthMethod)
	populate(objectMap, "rdb-backup-enabled", c.RdbBackupEnabled)
	populate(objectMap, "rdb-backup-frequency", c.RdbBackupFrequency)
	populate(objectMap, "rdb-backup-max-snapshot-count", c.RdbBackupMaxSnapshotCount)
	populate(objectMap, "rdb-storage-connection-string", c.RdbStorageConnectionString)
	populate(objectMap, "zonal-configuration", c.ZonalConfiguration)
	if c.AdditionalProperties != nil {
		for key, val := range c.AdditionalProperties {
			objectMap[key] = val
		}
	}
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type CommonPropertiesRedisConfiguration.
func (c *CommonPropertiesRedisConfiguration) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "aof-storage-connection-string-0":
			err = unpopulate(val, &c.AofStorageConnectionString0)
			delete(rawMsg, key)
		case "aof-storage-connection-string-1":
			err = unpopulate(val, &c.AofStorageConnectionString1)
			delete(rawMsg, key)
		case "maxclients":
			err = unpopulate(val, &c.Maxclients)
			delete(rawMsg, key)
		case "maxfragmentationmemory-reserved":
			err = unpopulate(val, &c.MaxfragmentationmemoryReserved)
			delete(rawMsg, key)
		case "maxmemory-delta":
			err = unpopulate(val, &c.MaxmemoryDelta)
			delete(rawMsg, key)
		case "maxmemory-policy":
			err = unpopulate(val, &c.MaxmemoryPolicy)
			delete(rawMsg, key)
		case "maxmemory-reserved":
			err = unpopulate(val, &c.MaxmemoryReserved)
			delete(rawMsg, key)
		case "preferred-data-archive-auth-method":
			err = unpopulate(val, &c.PreferredDataArchiveAuthMethod)
			delete(rawMsg, key)
		case "preferred-data-persistence-auth-method":
			err = unpopulate(val, &c.PreferredDataPersistenceAuthMethod)
			delete(rawMsg, key)
		case "rdb-backup-enabled":
			err = unpopulate(val, &c.RdbBackupEnabled)
			delete(rawMsg, key)
		case "rdb-backup-frequency":
			err = unpopulate(val, &c.RdbBackupFrequency)
			delete(rawMsg, key)
		case "rdb-backup-max-snapshot-count":
			err = unpopulate(val, &c.RdbBackupMaxSnapshotCount)
			delete(rawMsg, key)
		case "rdb-storage-connection-string":
			err = unpopulate(val, &c.RdbStorageConnectionString)
			delete(rawMsg, key)
		case "zonal-configuration":
			err = unpopulate(val, &c.ZonalConfiguration)
			delete(rawMsg, key)
		default:
			if c.AdditionalProperties == nil {
				c.AdditionalProperties = map[string]interface{}{}
			}
			if val != nil {
				var aux interface{}
				err = json.Unmarshal(val, &aux)
				c.AdditionalProperties[key] = aux
			}
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type CreateParameters.
func (c CreateParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "identity", c.Identity)
	populate(objectMap, "location", c.Location)
	populate(objectMap, "properties", c.Properties)
	populate(objectMap, "tags", c.Tags)
	populate(objectMap, "zones", c.Zones)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type CreateProperties.
func (c CreateProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "enableNonSslPort", c.EnableNonSSLPort)
	populate(objectMap, "minimumTlsVersion", c.MinimumTLSVersion)
	populate(objectMap, "publicNetworkAccess", c.PublicNetworkAccess)
	populate(objectMap, "redisConfiguration", c.RedisConfiguration)
	populate(objectMap, "redisVersion", c.RedisVersion)
	populate(objectMap, "replicasPerMaster", c.ReplicasPerMaster)
	populate(objectMap, "replicasPerPrimary", c.ReplicasPerPrimary)
	populate(objectMap, "sku", c.SKU)
	populate(objectMap, "shardCount", c.ShardCount)
	populate(objectMap, "staticIP", c.StaticIP)
	populate(objectMap, "subnetId", c.SubnetID)
	populate(objectMap, "tenantSettings", c.TenantSettings)
	return json.Marshal(objectMap)
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

// MarshalJSON implements the json.Marshaller interface for type ErrorDetailAutoGenerated.
func (e ErrorDetailAutoGenerated) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "additionalInfo", e.AdditionalInfo)
	populate(objectMap, "code", e.Code)
	populate(objectMap, "details", e.Details)
	populate(objectMap, "message", e.Message)
	populate(objectMap, "target", e.Target)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type FirewallRuleListResult.
func (f FirewallRuleListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", f.NextLink)
	populate(objectMap, "value", f.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ImportRDBParameters.
func (i ImportRDBParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "files", i.Files)
	populate(objectMap, "format", i.Format)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type LinkedServerWithPropertiesList.
func (l LinkedServerWithPropertiesList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", l.NextLink)
	populate(objectMap, "value", l.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ListResult.
func (l ListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", l.NextLink)
	populate(objectMap, "value", l.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ManagedServiceIdentity.
func (m ManagedServiceIdentity) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "principalId", m.PrincipalID)
	populate(objectMap, "tenantId", m.TenantID)
	populate(objectMap, "type", m.Type)
	populate(objectMap, "userAssignedIdentities", m.UserAssignedIdentities)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type NotificationListResponse.
func (n NotificationListResponse) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", n.NextLink)
	populate(objectMap, "value", n.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type OperationListResult.
func (o OperationListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", o.NextLink)
	populate(objectMap, "value", o.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type OperationStatus.
func (o OperationStatus) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populateTimeRFC3339(objectMap, "endTime", o.EndTime)
	populate(objectMap, "error", o.Error)
	populate(objectMap, "id", o.ID)
	populate(objectMap, "name", o.Name)
	populate(objectMap, "operations", o.Operations)
	populate(objectMap, "percentComplete", o.PercentComplete)
	populate(objectMap, "properties", o.Properties)
	populateTimeRFC3339(objectMap, "startTime", o.StartTime)
	populate(objectMap, "status", o.Status)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type OperationStatus.
func (o *OperationStatus) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "endTime":
			err = unpopulateTimeRFC3339(val, &o.EndTime)
			delete(rawMsg, key)
		case "error":
			err = unpopulate(val, &o.Error)
			delete(rawMsg, key)
		case "id":
			err = unpopulate(val, &o.ID)
			delete(rawMsg, key)
		case "name":
			err = unpopulate(val, &o.Name)
			delete(rawMsg, key)
		case "operations":
			err = unpopulate(val, &o.Operations)
			delete(rawMsg, key)
		case "percentComplete":
			err = unpopulate(val, &o.PercentComplete)
			delete(rawMsg, key)
		case "properties":
			err = unpopulate(val, &o.Properties)
			delete(rawMsg, key)
		case "startTime":
			err = unpopulateTimeRFC3339(val, &o.StartTime)
			delete(rawMsg, key)
		case "status":
			err = unpopulate(val, &o.Status)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type OperationStatusResult.
func (o OperationStatusResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populateTimeRFC3339(objectMap, "endTime", o.EndTime)
	populate(objectMap, "error", o.Error)
	populate(objectMap, "id", o.ID)
	populate(objectMap, "name", o.Name)
	populate(objectMap, "operations", o.Operations)
	populate(objectMap, "percentComplete", o.PercentComplete)
	populateTimeRFC3339(objectMap, "startTime", o.StartTime)
	populate(objectMap, "status", o.Status)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type OperationStatusResult.
func (o *OperationStatusResult) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "endTime":
			err = unpopulateTimeRFC3339(val, &o.EndTime)
			delete(rawMsg, key)
		case "error":
			err = unpopulate(val, &o.Error)
			delete(rawMsg, key)
		case "id":
			err = unpopulate(val, &o.ID)
			delete(rawMsg, key)
		case "name":
			err = unpopulate(val, &o.Name)
			delete(rawMsg, key)
		case "operations":
			err = unpopulate(val, &o.Operations)
			delete(rawMsg, key)
		case "percentComplete":
			err = unpopulate(val, &o.PercentComplete)
			delete(rawMsg, key)
		case "startTime":
			err = unpopulateTimeRFC3339(val, &o.StartTime)
			delete(rawMsg, key)
		case "status":
			err = unpopulate(val, &o.Status)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type PatchScheduleListResult.
func (p PatchScheduleListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", p.NextLink)
	populate(objectMap, "value", p.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PrivateEndpointConnectionListResult.
func (p PrivateEndpointConnectionListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "value", p.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PrivateLinkResourceListResult.
func (p PrivateLinkResourceListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "value", p.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PrivateLinkResourceProperties.
func (p PrivateLinkResourceProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "groupId", p.GroupID)
	populate(objectMap, "requiredMembers", p.RequiredMembers)
	populate(objectMap, "requiredZoneNames", p.RequiredZoneNames)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type Properties.
func (p Properties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "accessKeys", p.AccessKeys)
	populate(objectMap, "enableNonSslPort", p.EnableNonSSLPort)
	populate(objectMap, "hostName", p.HostName)
	populate(objectMap, "instances", p.Instances)
	populate(objectMap, "linkedServers", p.LinkedServers)
	populate(objectMap, "minimumTlsVersion", p.MinimumTLSVersion)
	populate(objectMap, "port", p.Port)
	populate(objectMap, "privateEndpointConnections", p.PrivateEndpointConnections)
	populate(objectMap, "provisioningState", p.ProvisioningState)
	populate(objectMap, "publicNetworkAccess", p.PublicNetworkAccess)
	populate(objectMap, "redisConfiguration", p.RedisConfiguration)
	populate(objectMap, "redisVersion", p.RedisVersion)
	populate(objectMap, "replicasPerMaster", p.ReplicasPerMaster)
	populate(objectMap, "replicasPerPrimary", p.ReplicasPerPrimary)
	populate(objectMap, "sku", p.SKU)
	populate(objectMap, "sslPort", p.SSLPort)
	populate(objectMap, "shardCount", p.ShardCount)
	populate(objectMap, "staticIP", p.StaticIP)
	populate(objectMap, "subnetId", p.SubnetID)
	populate(objectMap, "tenantSettings", p.TenantSettings)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type RebootParameters.
func (r RebootParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "ports", r.Ports)
	populate(objectMap, "rebootType", r.RebootType)
	populate(objectMap, "shardId", r.ShardID)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ResourceInfo.
func (r ResourceInfo) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "id", r.ID)
	populate(objectMap, "identity", r.Identity)
	populate(objectMap, "location", r.Location)
	populate(objectMap, "name", r.Name)
	populate(objectMap, "properties", r.Properties)
	populate(objectMap, "tags", r.Tags)
	populate(objectMap, "type", r.Type)
	populate(objectMap, "zones", r.Zones)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ScheduleEntries.
func (s ScheduleEntries) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "scheduleEntries", s.ScheduleEntries)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type TrackedResource.
func (t TrackedResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "id", t.ID)
	populate(objectMap, "location", t.Location)
	populate(objectMap, "name", t.Name)
	populate(objectMap, "tags", t.Tags)
	populate(objectMap, "type", t.Type)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type UpdateParameters.
func (u UpdateParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "identity", u.Identity)
	populate(objectMap, "properties", u.Properties)
	populate(objectMap, "tags", u.Tags)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type UpdateProperties.
func (u UpdateProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "enableNonSslPort", u.EnableNonSSLPort)
	populate(objectMap, "minimumTlsVersion", u.MinimumTLSVersion)
	populate(objectMap, "publicNetworkAccess", u.PublicNetworkAccess)
	populate(objectMap, "redisConfiguration", u.RedisConfiguration)
	populate(objectMap, "redisVersion", u.RedisVersion)
	populate(objectMap, "replicasPerMaster", u.ReplicasPerMaster)
	populate(objectMap, "replicasPerPrimary", u.ReplicasPerPrimary)
	populate(objectMap, "sku", u.SKU)
	populate(objectMap, "shardCount", u.ShardCount)
	populate(objectMap, "tenantSettings", u.TenantSettings)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type UpgradeNotification.
func (u UpgradeNotification) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "name", u.Name)
	populateTimeRFC3339(objectMap, "timestamp", u.Timestamp)
	populate(objectMap, "upsellNotification", u.UpsellNotification)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type UpgradeNotification.
func (u *UpgradeNotification) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "name":
			err = unpopulate(val, &u.Name)
			delete(rawMsg, key)
		case "timestamp":
			err = unpopulateTimeRFC3339(val, &u.Timestamp)
			delete(rawMsg, key)
		case "upsellNotification":
			err = unpopulate(val, &u.UpsellNotification)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
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

func unpopulate(data json.RawMessage, v interface{}) error {
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, v)
}
