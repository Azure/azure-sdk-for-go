package redis

import (
	 original "github.com/Azure/azure-sdk-for-go/service/redis/management/2017-02-01/redis"
)

type (
	 LinkedServerClient = original.LinkedServerClient
	 DayOfWeek = original.DayOfWeek
	 KeyType = original.KeyType
	 RebootType = original.RebootType
	 ReplicationRole = original.ReplicationRole
	 SkuFamily = original.SkuFamily
	 SkuName = original.SkuName
	 AccessKeys = original.AccessKeys
	 CreateParameters = original.CreateParameters
	 CreateProperties = original.CreateProperties
	 ExportRDBParameters = original.ExportRDBParameters
	 ForceRebootResponse = original.ForceRebootResponse
	 ImportRDBParameters = original.ImportRDBParameters
	 LinkedServer = original.LinkedServer
	 LinkedServerCreateParameters = original.LinkedServerCreateParameters
	 LinkedServerCreateProperties = original.LinkedServerCreateProperties
	 LinkedServerList = original.LinkedServerList
	 LinkedServerProperties = original.LinkedServerProperties
	 LinkedServerWithProperties = original.LinkedServerWithProperties
	 LinkedServerWithPropertiesList = original.LinkedServerWithPropertiesList
	 ListResult = original.ListResult
	 PatchSchedule = original.PatchSchedule
	 Properties = original.Properties
	 RebootParameters = original.RebootParameters
	 RegenerateKeyParameters = original.RegenerateKeyParameters
	 Resource = original.Resource
	 ResourceProperties = original.ResourceProperties
	 ResourceType = original.ResourceType
	 ScheduleEntries = original.ScheduleEntries
	 ScheduleEntry = original.ScheduleEntry
	 Sku = original.Sku
	 UpdateParameters = original.UpdateParameters
	 UpdateProperties = original.UpdateProperties
	 PatchSchedulesClient = original.PatchSchedulesClient
	 GroupClient = original.GroupClient
	 ManagementClient = original.ManagementClient
)

const (
	 Everyday = original.Everyday
	 Friday = original.Friday
	 Monday = original.Monday
	 Saturday = original.Saturday
	 Sunday = original.Sunday
	 Thursday = original.Thursday
	 Tuesday = original.Tuesday
	 Wednesday = original.Wednesday
	 Weekend = original.Weekend
	 Primary = original.Primary
	 Secondary = original.Secondary
	 AllNodes = original.AllNodes
	 PrimaryNode = original.PrimaryNode
	 SecondaryNode = original.SecondaryNode
	 ReplicationRolePrimary = original.ReplicationRolePrimary
	 ReplicationRoleSecondary = original.ReplicationRoleSecondary
	 C = original.C
	 P = original.P
	 Basic = original.Basic
	 Premium = original.Premium
	 Standard = original.Standard
	 DefaultBaseURI = original.DefaultBaseURI
)
