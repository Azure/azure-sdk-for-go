package subscriptions

import (
	 original "github.com/Azure/azure-sdk-for-go/service/resources/management/2016-06-01/subscriptions"
)

type (
	 ManagementClient = original.ManagementClient
	 SpendingLimit = original.SpendingLimit
	 State = original.State
	 ListResult = original.ListResult
	 Location = original.Location
	 LocationListResult = original.LocationListResult
	 Policies = original.Policies
	 Subscription = original.Subscription
	 TenantIDDescription = original.TenantIDDescription
	 TenantListResult = original.TenantListResult
	 GroupClient = original.GroupClient
	 TenantsClient = original.TenantsClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 CurrentPeriodOff = original.CurrentPeriodOff
	 Off = original.Off
	 On = original.On
	 Deleted = original.Deleted
	 Disabled = original.Disabled
	 Enabled = original.Enabled
	 PastDue = original.PastDue
	 Warned = original.Warned
)
