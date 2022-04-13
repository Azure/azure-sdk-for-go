//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armfrontdoor

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
)

// MarshalJSON implements the json.Marshaller interface for type BackendPoolListResult.
func (b BackendPoolListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", b.NextLink)
	populate(objectMap, "value", b.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type BackendPoolProperties.
func (b BackendPoolProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "backends", b.Backends)
	populate(objectMap, "healthProbeSettings", b.HealthProbeSettings)
	populate(objectMap, "loadBalancingSettings", b.LoadBalancingSettings)
	populate(objectMap, "resourceState", b.ResourceState)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type BackendPoolUpdateParameters.
func (b BackendPoolUpdateParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "backends", b.Backends)
	populate(objectMap, "healthProbeSettings", b.HealthProbeSettings)
	populate(objectMap, "loadBalancingSettings", b.LoadBalancingSettings)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type CustomRule.
func (c CustomRule) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "action", c.Action)
	populate(objectMap, "enabledState", c.EnabledState)
	populate(objectMap, "matchConditions", c.MatchConditions)
	populate(objectMap, "name", c.Name)
	populate(objectMap, "priority", c.Priority)
	populate(objectMap, "rateLimitDurationInMinutes", c.RateLimitDurationInMinutes)
	populate(objectMap, "rateLimitThreshold", c.RateLimitThreshold)
	populate(objectMap, "ruleType", c.RuleType)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type CustomRuleList.
func (c CustomRuleList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "rules", c.Rules)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type Error.
func (e Error) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "code", e.Code)
	populate(objectMap, "details", e.Details)
	populate(objectMap, "innerError", e.InnerError)
	populate(objectMap, "message", e.Message)
	populate(objectMap, "target", e.Target)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type Experiment.
func (e Experiment) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "id", e.ID)
	populate(objectMap, "location", e.Location)
	populate(objectMap, "name", e.Name)
	populate(objectMap, "properties", e.Properties)
	populate(objectMap, "tags", e.Tags)
	populate(objectMap, "type", e.Type)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ExperimentList.
func (e ExperimentList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", e.NextLink)
	populate(objectMap, "value", e.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ExperimentUpdateModel.
func (e ExperimentUpdateModel) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "properties", e.Properties)
	populate(objectMap, "tags", e.Tags)
	return json.Marshal(objectMap)
}

// GetRouteConfiguration implements the RouteConfigurationClassification interface for type ForwardingConfiguration.
func (f *ForwardingConfiguration) GetRouteConfiguration() *RouteConfiguration {
	return &RouteConfiguration{
		ODataType: f.ODataType,
	}
}

// MarshalJSON implements the json.Marshaller interface for type ForwardingConfiguration.
func (f ForwardingConfiguration) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "backendPool", f.BackendPool)
	populate(objectMap, "cacheConfiguration", f.CacheConfiguration)
	populate(objectMap, "customForwardingPath", f.CustomForwardingPath)
	populate(objectMap, "forwardingProtocol", f.ForwardingProtocol)
	objectMap["@odata.type"] = "#Microsoft.Azure.FrontDoor.Models.FrontdoorForwardingConfiguration"
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type ForwardingConfiguration.
func (f *ForwardingConfiguration) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "backendPool":
			err = unpopulate(val, &f.BackendPool)
			delete(rawMsg, key)
		case "cacheConfiguration":
			err = unpopulate(val, &f.CacheConfiguration)
			delete(rawMsg, key)
		case "customForwardingPath":
			err = unpopulate(val, &f.CustomForwardingPath)
			delete(rawMsg, key)
		case "forwardingProtocol":
			err = unpopulate(val, &f.ForwardingProtocol)
			delete(rawMsg, key)
		case "@odata.type":
			err = unpopulate(val, &f.ODataType)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type FrontDoor.
func (f FrontDoor) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "id", f.ID)
	populate(objectMap, "location", f.Location)
	populate(objectMap, "name", f.Name)
	populate(objectMap, "properties", f.Properties)
	populate(objectMap, "tags", f.Tags)
	populate(objectMap, "type", f.Type)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type FrontendEndpointsListResult.
func (f FrontendEndpointsListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", f.NextLink)
	populate(objectMap, "value", f.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type HealthProbeSettingsListResult.
func (h HealthProbeSettingsListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", h.NextLink)
	populate(objectMap, "value", h.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type LatencyScorecard.
func (l LatencyScorecard) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "id", l.ID)
	populate(objectMap, "location", l.Location)
	populate(objectMap, "name", l.Name)
	populate(objectMap, "properties", l.Properties)
	populate(objectMap, "tags", l.Tags)
	populate(objectMap, "type", l.Type)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type LatencyScorecardProperties.
func (l LatencyScorecardProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "country", l.Country)
	populate(objectMap, "description", l.Description)
	populateTimeRFC3339(objectMap, "endDateTimeUTC", l.EndDateTimeUTC)
	populate(objectMap, "endpointA", l.EndpointA)
	populate(objectMap, "endpointB", l.EndpointB)
	populate(objectMap, "id", l.ID)
	populate(objectMap, "latencyMetrics", l.LatencyMetrics)
	populate(objectMap, "name", l.Name)
	populateTimeRFC3339(objectMap, "startDateTimeUTC", l.StartDateTimeUTC)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type LatencyScorecardProperties.
func (l *LatencyScorecardProperties) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "country":
			err = unpopulate(val, &l.Country)
			delete(rawMsg, key)
		case "description":
			err = unpopulate(val, &l.Description)
			delete(rawMsg, key)
		case "endDateTimeUTC":
			err = unpopulateTimeRFC3339(val, &l.EndDateTimeUTC)
			delete(rawMsg, key)
		case "endpointA":
			err = unpopulate(val, &l.EndpointA)
			delete(rawMsg, key)
		case "endpointB":
			err = unpopulate(val, &l.EndpointB)
			delete(rawMsg, key)
		case "id":
			err = unpopulate(val, &l.ID)
			delete(rawMsg, key)
		case "latencyMetrics":
			err = unpopulate(val, &l.LatencyMetrics)
			delete(rawMsg, key)
		case "name":
			err = unpopulate(val, &l.Name)
			delete(rawMsg, key)
		case "startDateTimeUTC":
			err = unpopulateTimeRFC3339(val, &l.StartDateTimeUTC)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type ListResult.
func (l ListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", l.NextLink)
	populate(objectMap, "value", l.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type LoadBalancingSettingsListResult.
func (l LoadBalancingSettingsListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", l.NextLink)
	populate(objectMap, "value", l.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ManagedRuleGroupDefinition.
func (m ManagedRuleGroupDefinition) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "description", m.Description)
	populate(objectMap, "ruleGroupName", m.RuleGroupName)
	populate(objectMap, "rules", m.Rules)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ManagedRuleGroupOverride.
func (m ManagedRuleGroupOverride) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "exclusions", m.Exclusions)
	populate(objectMap, "ruleGroupName", m.RuleGroupName)
	populate(objectMap, "rules", m.Rules)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ManagedRuleOverride.
func (m ManagedRuleOverride) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "action", m.Action)
	populate(objectMap, "enabledState", m.EnabledState)
	populate(objectMap, "exclusions", m.Exclusions)
	populate(objectMap, "ruleId", m.RuleID)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ManagedRuleSet.
func (m ManagedRuleSet) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "exclusions", m.Exclusions)
	populate(objectMap, "ruleGroupOverrides", m.RuleGroupOverrides)
	populate(objectMap, "ruleSetAction", m.RuleSetAction)
	populate(objectMap, "ruleSetType", m.RuleSetType)
	populate(objectMap, "ruleSetVersion", m.RuleSetVersion)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ManagedRuleSetDefinition.
func (m ManagedRuleSetDefinition) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "id", m.ID)
	populate(objectMap, "location", m.Location)
	populate(objectMap, "name", m.Name)
	populate(objectMap, "properties", m.Properties)
	populate(objectMap, "tags", m.Tags)
	populate(objectMap, "type", m.Type)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ManagedRuleSetDefinitionList.
func (m ManagedRuleSetDefinitionList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", m.NextLink)
	populate(objectMap, "value", m.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ManagedRuleSetDefinitionProperties.
func (m ManagedRuleSetDefinitionProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "provisioningState", m.ProvisioningState)
	populate(objectMap, "ruleGroups", m.RuleGroups)
	populate(objectMap, "ruleSetId", m.RuleSetID)
	populate(objectMap, "ruleSetType", m.RuleSetType)
	populate(objectMap, "ruleSetVersion", m.RuleSetVersion)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ManagedRuleSetList.
func (m ManagedRuleSetList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "managedRuleSets", m.ManagedRuleSets)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type MatchCondition.
func (m MatchCondition) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "matchValue", m.MatchValue)
	populate(objectMap, "matchVariable", m.MatchVariable)
	populate(objectMap, "negateCondition", m.NegateCondition)
	populate(objectMap, "operator", m.Operator)
	populate(objectMap, "selector", m.Selector)
	populate(objectMap, "transforms", m.Transforms)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PreconfiguredEndpoint.
func (p PreconfiguredEndpoint) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "id", p.ID)
	populate(objectMap, "location", p.Location)
	populate(objectMap, "name", p.Name)
	populate(objectMap, "properties", p.Properties)
	populate(objectMap, "tags", p.Tags)
	populate(objectMap, "type", p.Type)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PreconfiguredEndpointList.
func (p PreconfiguredEndpointList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", p.NextLink)
	populate(objectMap, "value", p.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type Profile.
func (p Profile) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "etag", p.Etag)
	populate(objectMap, "id", p.ID)
	populate(objectMap, "location", p.Location)
	populate(objectMap, "name", p.Name)
	populate(objectMap, "properties", p.Properties)
	populate(objectMap, "tags", p.Tags)
	populate(objectMap, "type", p.Type)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ProfileList.
func (p ProfileList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", p.NextLink)
	populate(objectMap, "value", p.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type ProfileUpdateModel.
func (p ProfileUpdateModel) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "properties", p.Properties)
	populate(objectMap, "tags", p.Tags)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type Properties.
func (p Properties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "backendPools", p.BackendPools)
	populate(objectMap, "backendPoolsSettings", p.BackendPoolsSettings)
	populate(objectMap, "cname", p.Cname)
	populate(objectMap, "enabledState", p.EnabledState)
	populate(objectMap, "friendlyName", p.FriendlyName)
	populate(objectMap, "frontdoorId", p.FrontdoorID)
	populate(objectMap, "frontendEndpoints", p.FrontendEndpoints)
	populate(objectMap, "healthProbeSettings", p.HealthProbeSettings)
	populate(objectMap, "loadBalancingSettings", p.LoadBalancingSettings)
	populate(objectMap, "provisioningState", p.ProvisioningState)
	populate(objectMap, "resourceState", p.ResourceState)
	populate(objectMap, "routingRules", p.RoutingRules)
	populate(objectMap, "rulesEngines", p.RulesEngines)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type PurgeParameters.
func (p PurgeParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "contentPaths", p.ContentPaths)
	return json.Marshal(objectMap)
}

// GetRouteConfiguration implements the RouteConfigurationClassification interface for type RedirectConfiguration.
func (r *RedirectConfiguration) GetRouteConfiguration() *RouteConfiguration {
	return &RouteConfiguration{
		ODataType: r.ODataType,
	}
}

// MarshalJSON implements the json.Marshaller interface for type RedirectConfiguration.
func (r RedirectConfiguration) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "customFragment", r.CustomFragment)
	populate(objectMap, "customHost", r.CustomHost)
	populate(objectMap, "customPath", r.CustomPath)
	populate(objectMap, "customQueryString", r.CustomQueryString)
	objectMap["@odata.type"] = "#Microsoft.Azure.FrontDoor.Models.FrontdoorRedirectConfiguration"
	populate(objectMap, "redirectProtocol", r.RedirectProtocol)
	populate(objectMap, "redirectType", r.RedirectType)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type RedirectConfiguration.
func (r *RedirectConfiguration) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "customFragment":
			err = unpopulate(val, &r.CustomFragment)
			delete(rawMsg, key)
		case "customHost":
			err = unpopulate(val, &r.CustomHost)
			delete(rawMsg, key)
		case "customPath":
			err = unpopulate(val, &r.CustomPath)
			delete(rawMsg, key)
		case "customQueryString":
			err = unpopulate(val, &r.CustomQueryString)
			delete(rawMsg, key)
		case "@odata.type":
			err = unpopulate(val, &r.ODataType)
			delete(rawMsg, key)
		case "redirectProtocol":
			err = unpopulate(val, &r.RedirectProtocol)
			delete(rawMsg, key)
		case "redirectType":
			err = unpopulate(val, &r.RedirectType)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type Resource.
func (r Resource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "id", r.ID)
	populate(objectMap, "location", r.Location)
	populate(objectMap, "name", r.Name)
	populate(objectMap, "tags", r.Tags)
	populate(objectMap, "type", r.Type)
	return json.Marshal(objectMap)
}

// GetRouteConfiguration implements the RouteConfigurationClassification interface for type RouteConfiguration.
func (r *RouteConfiguration) GetRouteConfiguration() *RouteConfiguration { return r }

// MarshalJSON implements the json.Marshaller interface for type RoutingRuleListResult.
func (r RoutingRuleListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", r.NextLink)
	populate(objectMap, "value", r.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type RoutingRuleProperties.
func (r RoutingRuleProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "acceptedProtocols", r.AcceptedProtocols)
	populate(objectMap, "enabledState", r.EnabledState)
	populate(objectMap, "frontendEndpoints", r.FrontendEndpoints)
	populate(objectMap, "patternsToMatch", r.PatternsToMatch)
	populate(objectMap, "resourceState", r.ResourceState)
	populate(objectMap, "routeConfiguration", r.RouteConfiguration)
	populate(objectMap, "rulesEngine", r.RulesEngine)
	populate(objectMap, "webApplicationFirewallPolicyLink", r.WebApplicationFirewallPolicyLink)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type RoutingRuleProperties.
func (r *RoutingRuleProperties) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "acceptedProtocols":
			err = unpopulate(val, &r.AcceptedProtocols)
			delete(rawMsg, key)
		case "enabledState":
			err = unpopulate(val, &r.EnabledState)
			delete(rawMsg, key)
		case "frontendEndpoints":
			err = unpopulate(val, &r.FrontendEndpoints)
			delete(rawMsg, key)
		case "patternsToMatch":
			err = unpopulate(val, &r.PatternsToMatch)
			delete(rawMsg, key)
		case "resourceState":
			err = unpopulate(val, &r.ResourceState)
			delete(rawMsg, key)
		case "routeConfiguration":
			r.RouteConfiguration, err = unmarshalRouteConfigurationClassification(val)
			delete(rawMsg, key)
		case "rulesEngine":
			err = unpopulate(val, &r.RulesEngine)
			delete(rawMsg, key)
		case "webApplicationFirewallPolicyLink":
			err = unpopulate(val, &r.WebApplicationFirewallPolicyLink)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type RoutingRuleUpdateParameters.
func (r RoutingRuleUpdateParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "acceptedProtocols", r.AcceptedProtocols)
	populate(objectMap, "enabledState", r.EnabledState)
	populate(objectMap, "frontendEndpoints", r.FrontendEndpoints)
	populate(objectMap, "patternsToMatch", r.PatternsToMatch)
	populate(objectMap, "routeConfiguration", r.RouteConfiguration)
	populate(objectMap, "rulesEngine", r.RulesEngine)
	populate(objectMap, "webApplicationFirewallPolicyLink", r.WebApplicationFirewallPolicyLink)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type RoutingRuleUpdateParameters.
func (r *RoutingRuleUpdateParameters) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "acceptedProtocols":
			err = unpopulate(val, &r.AcceptedProtocols)
			delete(rawMsg, key)
		case "enabledState":
			err = unpopulate(val, &r.EnabledState)
			delete(rawMsg, key)
		case "frontendEndpoints":
			err = unpopulate(val, &r.FrontendEndpoints)
			delete(rawMsg, key)
		case "patternsToMatch":
			err = unpopulate(val, &r.PatternsToMatch)
			delete(rawMsg, key)
		case "routeConfiguration":
			r.RouteConfiguration, err = unmarshalRouteConfigurationClassification(val)
			delete(rawMsg, key)
		case "rulesEngine":
			err = unpopulate(val, &r.RulesEngine)
			delete(rawMsg, key)
		case "webApplicationFirewallPolicyLink":
			err = unpopulate(val, &r.WebApplicationFirewallPolicyLink)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type RulesEngineAction.
func (r RulesEngineAction) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "requestHeaderActions", r.RequestHeaderActions)
	populate(objectMap, "responseHeaderActions", r.ResponseHeaderActions)
	populate(objectMap, "routeConfigurationOverride", r.RouteConfigurationOverride)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type RulesEngineAction.
func (r *RulesEngineAction) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "requestHeaderActions":
			err = unpopulate(val, &r.RequestHeaderActions)
			delete(rawMsg, key)
		case "responseHeaderActions":
			err = unpopulate(val, &r.ResponseHeaderActions)
			delete(rawMsg, key)
		case "routeConfigurationOverride":
			r.RouteConfigurationOverride, err = unmarshalRouteConfigurationClassification(val)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// MarshalJSON implements the json.Marshaller interface for type RulesEngineListResult.
func (r RulesEngineListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", r.NextLink)
	populate(objectMap, "value", r.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type RulesEngineMatchCondition.
func (r RulesEngineMatchCondition) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "negateCondition", r.NegateCondition)
	populate(objectMap, "rulesEngineMatchValue", r.RulesEngineMatchValue)
	populate(objectMap, "rulesEngineMatchVariable", r.RulesEngineMatchVariable)
	populate(objectMap, "rulesEngineOperator", r.RulesEngineOperator)
	populate(objectMap, "selector", r.Selector)
	populate(objectMap, "transforms", r.Transforms)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type RulesEngineProperties.
func (r RulesEngineProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "resourceState", r.ResourceState)
	populate(objectMap, "rules", r.Rules)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type RulesEngineRule.
func (r RulesEngineRule) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "action", r.Action)
	populate(objectMap, "matchConditions", r.MatchConditions)
	populate(objectMap, "matchProcessingBehavior", r.MatchProcessingBehavior)
	populate(objectMap, "name", r.Name)
	populate(objectMap, "priority", r.Priority)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type RulesEngineUpdateParameters.
func (r RulesEngineUpdateParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "rules", r.Rules)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type TagsObject.
func (t TagsObject) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "tags", t.Tags)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type Timeseries.
func (t Timeseries) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "id", t.ID)
	populate(objectMap, "location", t.Location)
	populate(objectMap, "name", t.Name)
	populate(objectMap, "properties", t.Properties)
	populate(objectMap, "tags", t.Tags)
	populate(objectMap, "type", t.Type)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type TimeseriesProperties.
func (t TimeseriesProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "aggregationInterval", t.AggregationInterval)
	populate(objectMap, "country", t.Country)
	populate(objectMap, "endDateTimeUTC", t.EndDateTimeUTC)
	populate(objectMap, "endpoint", t.Endpoint)
	populate(objectMap, "startDateTimeUTC", t.StartDateTimeUTC)
	populate(objectMap, "timeseriesData", t.TimeseriesData)
	populate(objectMap, "timeseriesType", t.TimeseriesType)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type UpdateParameters.
func (u UpdateParameters) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "backendPools", u.BackendPools)
	populate(objectMap, "backendPoolsSettings", u.BackendPoolsSettings)
	populate(objectMap, "enabledState", u.EnabledState)
	populate(objectMap, "friendlyName", u.FriendlyName)
	populate(objectMap, "frontendEndpoints", u.FrontendEndpoints)
	populate(objectMap, "healthProbeSettings", u.HealthProbeSettings)
	populate(objectMap, "loadBalancingSettings", u.LoadBalancingSettings)
	populate(objectMap, "routingRules", u.RoutingRules)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type WebApplicationFirewallPolicy.
func (w WebApplicationFirewallPolicy) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "etag", w.Etag)
	populate(objectMap, "id", w.ID)
	populate(objectMap, "location", w.Location)
	populate(objectMap, "name", w.Name)
	populate(objectMap, "properties", w.Properties)
	populate(objectMap, "sku", w.SKU)
	populate(objectMap, "tags", w.Tags)
	populate(objectMap, "type", w.Type)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type WebApplicationFirewallPolicyList.
func (w WebApplicationFirewallPolicyList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", w.NextLink)
	populate(objectMap, "value", w.Value)
	return json.Marshal(objectMap)
}

// MarshalJSON implements the json.Marshaller interface for type WebApplicationFirewallPolicyProperties.
func (w WebApplicationFirewallPolicyProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "customRules", w.CustomRules)
	populate(objectMap, "frontendEndpointLinks", w.FrontendEndpointLinks)
	populate(objectMap, "managedRules", w.ManagedRules)
	populate(objectMap, "policySettings", w.PolicySettings)
	populate(objectMap, "provisioningState", w.ProvisioningState)
	populate(objectMap, "resourceState", w.ResourceState)
	populate(objectMap, "routingRuleLinks", w.RoutingRuleLinks)
	populate(objectMap, "securityPolicyLinks", w.SecurityPolicyLinks)
	return json.Marshal(objectMap)
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
