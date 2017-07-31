package vmssnetworkinterface

import (
	 original "github.com/Azure/azure-sdk-for-go/service/network/management/2017-03-01/vmssNetworkInterface"
)

type (
	 ManagementClient = original.ManagementClient
	 IPAllocationMethod = original.IPAllocationMethod
	 IPVersion = original.IPVersion
	 RouteNextHopType = original.RouteNextHopType
	 SecurityRuleAccess = original.SecurityRuleAccess
	 SecurityRuleDirection = original.SecurityRuleDirection
	 SecurityRuleProtocol = original.SecurityRuleProtocol
	 TransportProtocol = original.TransportProtocol
	 ApplicationGatewayBackendAddress = original.ApplicationGatewayBackendAddress
	 ApplicationGatewayBackendAddressPool = original.ApplicationGatewayBackendAddressPool
	 ApplicationGatewayBackendAddressPoolPropertiesFormat = original.ApplicationGatewayBackendAddressPoolPropertiesFormat
	 BackendAddressPool = original.BackendAddressPool
	 BackendAddressPoolPropertiesFormat = original.BackendAddressPoolPropertiesFormat
	 InboundNatRule = original.InboundNatRule
	 InboundNatRulePropertiesFormat = original.InboundNatRulePropertiesFormat
	 IPConfiguration = original.IPConfiguration
	 IPConfigurationPropertiesFormat = original.IPConfigurationPropertiesFormat
	 NetworkInterface = original.NetworkInterface
	 NetworkInterfaceDNSSettings = original.NetworkInterfaceDNSSettings
	 NetworkInterfaceIPConfiguration = original.NetworkInterfaceIPConfiguration
	 NetworkInterfaceIPConfigurationPropertiesFormat = original.NetworkInterfaceIPConfigurationPropertiesFormat
	 NetworkInterfaceListResult = original.NetworkInterfaceListResult
	 NetworkInterfacePropertiesFormat = original.NetworkInterfacePropertiesFormat
	 NetworkSecurityGroup = original.NetworkSecurityGroup
	 NetworkSecurityGroupPropertiesFormat = original.NetworkSecurityGroupPropertiesFormat
	 PublicIPAddress = original.PublicIPAddress
	 PublicIPAddressDNSSettings = original.PublicIPAddressDNSSettings
	 PublicIPAddressPropertiesFormat = original.PublicIPAddressPropertiesFormat
	 Resource = original.Resource
	 ResourceNavigationLink = original.ResourceNavigationLink
	 ResourceNavigationLinkFormat = original.ResourceNavigationLinkFormat
	 Route = original.Route
	 RoutePropertiesFormat = original.RoutePropertiesFormat
	 RouteTable = original.RouteTable
	 RouteTablePropertiesFormat = original.RouteTablePropertiesFormat
	 SecurityRule = original.SecurityRule
	 SecurityRulePropertiesFormat = original.SecurityRulePropertiesFormat
	 Subnet = original.Subnet
	 SubnetPropertiesFormat = original.SubnetPropertiesFormat
	 SubResource = original.SubResource
	 NetworkInterfacesClient = original.NetworkInterfacesClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
	 Dynamic = original.Dynamic
	 Static = original.Static
	 IPv4 = original.IPv4
	 IPv6 = original.IPv6
	 Internet = original.Internet
	 None = original.None
	 VirtualAppliance = original.VirtualAppliance
	 VirtualNetworkGateway = original.VirtualNetworkGateway
	 VnetLocal = original.VnetLocal
	 Allow = original.Allow
	 Deny = original.Deny
	 Inbound = original.Inbound
	 Outbound = original.Outbound
	 Asterisk = original.Asterisk
	 TCP = original.TCP
	 UDP = original.UDP
	 TransportProtocolTCP = original.TransportProtocolTCP
	 TransportProtocolUDP = original.TransportProtocolUDP
)
