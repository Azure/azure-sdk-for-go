package firewallrules

import (
	 original "github.com/Azure/azure-sdk-for-go/service/sql/management/2014-04-01/firewallRules"
)

type (
	 ManagementClient = original.ManagementClient
	 GroupClient = original.GroupClient
	 FirewallRule = original.FirewallRule
	 ListResult = original.ListResult
	 Properties = original.Properties
	 ProxyResource = original.ProxyResource
	 Resource = original.Resource
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
)
