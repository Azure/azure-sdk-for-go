package replicationlinks

import (
	 original "github.com/Azure/azure-sdk-for-go/service/sql/management/2014-04-01/replicationLinks"
)

type (
	 ManagementClient = original.ManagementClient
	 ReplicationRole = original.ReplicationRole
	 ReplicationState = original.ReplicationState
	 ListResult = original.ListResult
	 Properties = original.Properties
	 ProxyResource = original.ProxyResource
	 ReplicationLink = original.ReplicationLink
	 Resource = original.Resource
	 GroupClient = original.GroupClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Copy = original.Copy
	 NonReadableSecondary = original.NonReadableSecondary
	 Primary = original.Primary
	 Secondary = original.Secondary
	 Source = original.Source
	 CATCHUP = original.CATCHUP
	 PENDING = original.PENDING
	 SEEDING = original.SEEDING
	 SUSPENDED = original.SUSPENDED
)
