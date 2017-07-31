package notificationhubs

import (
	 original "github.com/Azure/azure-sdk-for-go/service/notificationhubs/management/2017-04-01/notificationhubs"
)

type (
	 NamespacesClient = original.NamespacesClient
	 GroupClient = original.GroupClient
	 ManagementClient = original.ManagementClient
	 HubsClient = original.HubsClient
	 AccessRights = original.AccessRights
	 NamespaceType = original.NamespaceType
	 SkuName = original.SkuName
	 AdmCredential = original.AdmCredential
	 AdmCredentialProperties = original.AdmCredentialProperties
	 ApnsCredential = original.ApnsCredential
	 ApnsCredentialProperties = original.ApnsCredentialProperties
	 BaiduCredential = original.BaiduCredential
	 BaiduCredentialProperties = original.BaiduCredentialProperties
	 CheckAvailabilityParameters = original.CheckAvailabilityParameters
	 CheckAvailabilityResult = original.CheckAvailabilityResult
	 CheckNameAvailabilityRequestParameters = original.CheckNameAvailabilityRequestParameters
	 CheckNameAvailabilityResponse = original.CheckNameAvailabilityResponse
	 CreateOrUpdateParameters = original.CreateOrUpdateParameters
	 GcmCredential = original.GcmCredential
	 GcmCredentialProperties = original.GcmCredentialProperties
	 ListResult = original.ListResult
	 MpnsCredential = original.MpnsCredential
	 MpnsCredentialProperties = original.MpnsCredentialProperties
	 NamespaceCreateOrUpdateParameters = original.NamespaceCreateOrUpdateParameters
	 NamespaceListResult = original.NamespaceListResult
	 NamespacePatchParameters = original.NamespacePatchParameters
	 NamespaceProperties = original.NamespaceProperties
	 NamespaceResource = original.NamespaceResource
	 PnsCredentialsProperties = original.PnsCredentialsProperties
	 PnsCredentialsResource = original.PnsCredentialsResource
	 PolicykeyResource = original.PolicykeyResource
	 Properties = original.Properties
	 Resource = original.Resource
	 ResourceListKeys = original.ResourceListKeys
	 ResourceType = original.ResourceType
	 SharedAccessAuthorizationRuleCreateOrUpdateParameters = original.SharedAccessAuthorizationRuleCreateOrUpdateParameters
	 SharedAccessAuthorizationRuleListResult = original.SharedAccessAuthorizationRuleListResult
	 SharedAccessAuthorizationRuleProperties = original.SharedAccessAuthorizationRuleProperties
	 SharedAccessAuthorizationRuleResource = original.SharedAccessAuthorizationRuleResource
	 Sku = original.Sku
	 SubResource = original.SubResource
	 WnsCredential = original.WnsCredential
	 WnsCredentialProperties = original.WnsCredentialProperties
	 NameClient = original.NameClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Listen = original.Listen
	 Manage = original.Manage
	 Send = original.Send
	 Messaging = original.Messaging
	 NotificationHub = original.NotificationHub
	 Basic = original.Basic
	 Free = original.Free
	 Standard = original.Standard
)
