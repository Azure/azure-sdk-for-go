package connectionpolicies

import (
	 original "github.com/Azure/azure-sdk-for-go/service/sql/management/2014-04-01/connectionPolicies"
)

type (
	 ServerConnectionType = original.ServerConnectionType
	 ProxyResource = original.ProxyResource
	 Resource = original.Resource
	 ServerConnectionPolicy = original.ServerConnectionPolicy
	 ServerConnectionPolicyProperties = original.ServerConnectionPolicyProperties
	 ServerConnectionPoliciesClient = original.ServerConnectionPoliciesClient
	 ManagementClient = original.ManagementClient
)

const (
	 Default = original.Default
	 Proxy = original.Proxy
	 Redirect = original.Redirect
	 DefaultBaseURI = original.DefaultBaseURI
)
