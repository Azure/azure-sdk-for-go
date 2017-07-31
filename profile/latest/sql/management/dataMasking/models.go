package datamasking

import (
	 original "github.com/Azure/azure-sdk-for-go/service/sql/management/2014-04-01/dataMasking"
)

type (
	 ManagementClient = original.ManagementClient
	 Function = original.Function
	 RuleState = original.RuleState
	 State = original.State
	 Policy = original.Policy
	 PolicyProperties = original.PolicyProperties
	 ProxyResource = original.ProxyResource
	 Resource = original.Resource
	 Rule = original.Rule
	 RuleListResult = original.RuleListResult
	 RuleProperties = original.RuleProperties
	 PoliciesClient = original.PoliciesClient
	 RulesClient = original.RulesClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 CCN = original.CCN
	 Default = original.Default
	 Email = original.Email
	 Number = original.Number
	 SSN = original.SSN
	 Text = original.Text
	 Disabled = original.Disabled
	 Enabled = original.Enabled
	 StateDisabled = original.StateDisabled
	 StateEnabled = original.StateEnabled
)
