package databasesecurityalertpolicies

import (
	 original "github.com/Azure/azure-sdk-for-go/service/sql/management/2014-04-01/databaseSecurityAlertPolicies"
)

type (
	 ManagementClient = original.ManagementClient
	 DatabaseThreatDetectionPoliciesClient = original.DatabaseThreatDetectionPoliciesClient
	 SecurityAlertPolicyEmailAccountAdmins = original.SecurityAlertPolicyEmailAccountAdmins
	 SecurityAlertPolicyState = original.SecurityAlertPolicyState
	 SecurityAlertPolicyUseServerDefault = original.SecurityAlertPolicyUseServerDefault
	 DatabaseSecurityAlertPolicy = original.DatabaseSecurityAlertPolicy
	 DatabaseSecurityAlertPolicyProperties = original.DatabaseSecurityAlertPolicyProperties
	 ProxyResource = original.ProxyResource
	 Resource = original.Resource
	 TrackedResource = original.TrackedResource
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Disabled = original.Disabled
	 Enabled = original.Enabled
	 SecurityAlertPolicyStateDisabled = original.SecurityAlertPolicyStateDisabled
	 SecurityAlertPolicyStateEnabled = original.SecurityAlertPolicyStateEnabled
	 SecurityAlertPolicyStateNew = original.SecurityAlertPolicyStateNew
	 SecurityAlertPolicyUseServerDefaultDisabled = original.SecurityAlertPolicyUseServerDefaultDisabled
	 SecurityAlertPolicyUseServerDefaultEnabled = original.SecurityAlertPolicyUseServerDefaultEnabled
)
