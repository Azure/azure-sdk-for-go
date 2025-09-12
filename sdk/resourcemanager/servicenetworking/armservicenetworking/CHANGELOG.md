# Release History

## 1.2.0-beta.1 (2025-04-23)
### Features Added

- New value `PolicyTypeIPAccessRules` added to enum type `PolicyType`
- New enum type `IPAccessRuleAction` with values `IPAccessRuleActionAllow`, `IPAccessRuleActionDeny`
- New struct `FrontendUpdateProperties`
- New struct `IPAccessRule`
- New struct `IPAccessRulesPolicy`
- New struct `IPAccessRulesSecurityPolicy`
- New field `SecurityPolicyConfigurations` in struct `FrontendProperties`
- New field `Properties` in struct `FrontendUpdate`
- New field `IPAccessRulesSecurityPolicy` in struct `SecurityPolicyConfigurations`
- New field `IPAccessRulesPolicy` in struct `SecurityPolicyProperties`
- New field `IPAccessRulesPolicy` in struct `SecurityPolicyUpdateProperties`


## 1.1.0 (2025-01-23)
### Features Added

- New enum type `PolicyType` with values `PolicyTypeWAF`
- New function `*ClientFactory.NewSecurityPoliciesInterfaceClient() *SecurityPoliciesInterfaceClient`
- New function `NewSecurityPoliciesInterfaceClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SecurityPoliciesInterfaceClient, error)`
- New function `*SecurityPoliciesInterfaceClient.BeginCreateOrUpdate(context.Context, string, string, string, SecurityPolicy, *SecurityPoliciesInterfaceClientBeginCreateOrUpdateOptions) (*runtime.Poller[SecurityPoliciesInterfaceClientCreateOrUpdateResponse], error)`
- New function `*SecurityPoliciesInterfaceClient.BeginDelete(context.Context, string, string, string, *SecurityPoliciesInterfaceClientBeginDeleteOptions) (*runtime.Poller[SecurityPoliciesInterfaceClientDeleteResponse], error)`
- New function `*SecurityPoliciesInterfaceClient.Get(context.Context, string, string, string, *SecurityPoliciesInterfaceClientGetOptions) (SecurityPoliciesInterfaceClientGetResponse, error)`
- New function `*SecurityPoliciesInterfaceClient.NewListByTrafficControllerPager(string, string, *SecurityPoliciesInterfaceClientListByTrafficControllerOptions) *runtime.Pager[SecurityPoliciesInterfaceClientListByTrafficControllerResponse]`
- New function `*SecurityPoliciesInterfaceClient.Update(context.Context, string, string, string, SecurityPolicyUpdate, *SecurityPoliciesInterfaceClientUpdateOptions) (SecurityPoliciesInterfaceClientUpdateResponse, error)`
- New struct `SecurityPolicy`
- New struct `SecurityPolicyConfigurations`
- New struct `SecurityPolicyListResult`
- New struct `SecurityPolicyProperties`
- New struct `SecurityPolicyUpdate`
- New struct `SecurityPolicyUpdateProperties`
- New struct `TrafficControllerUpdateProperties`
- New struct `WafPolicy`
- New struct `WafSecurityPolicy`
- New field `SecurityPolicies`, `SecurityPolicyConfigurations` in struct `TrafficControllerProperties`
- New field `Properties` in struct `TrafficControllerUpdate`


## 1.1.0-beta.2 (2024-09-26)
### Bugs Fixed

- Fix wrong url according to the spec fix.

## 1.1.0-beta.1 (2024-08-23)
### Features Added

- New enum type `PolicyType` with values `PolicyTypeWAF`
- New function `*ClientFactory.NewSecurityPoliciesInterfaceClient() *SecurityPoliciesInterfaceClient`
- New function `NewSecurityPoliciesInterfaceClient(string, azcore.TokenCredential, *arm.ClientOptions) (*SecurityPoliciesInterfaceClient, error)`
- New function `*SecurityPoliciesInterfaceClient.BeginCreateOrUpdate(context.Context, string, string, string, SecurityPolicy, *SecurityPoliciesInterfaceClientBeginCreateOrUpdateOptions) (*runtime.Poller[SecurityPoliciesInterfaceClientCreateOrUpdateResponse], error)`
- New function `*SecurityPoliciesInterfaceClient.BeginDelete(context.Context, string, string, string, *SecurityPoliciesInterfaceClientBeginDeleteOptions) (*runtime.Poller[SecurityPoliciesInterfaceClientDeleteResponse], error)`
- New function `*SecurityPoliciesInterfaceClient.Get(context.Context, string, string, string, *SecurityPoliciesInterfaceClientGetOptions) (SecurityPoliciesInterfaceClientGetResponse, error)`
- New function `*SecurityPoliciesInterfaceClient.NewListByTrafficControllerPager(string, string, *SecurityPoliciesInterfaceClientListByTrafficControllerOptions) *runtime.Pager[SecurityPoliciesInterfaceClientListByTrafficControllerResponse]`
- New function `*SecurityPoliciesInterfaceClient.Update(context.Context, string, string, string, SecurityPolicyUpdate, *SecurityPoliciesInterfaceClientUpdateOptions) (SecurityPoliciesInterfaceClientUpdateResponse, error)`
- New struct `SecurityPolicy`
- New struct `SecurityPolicyConfigurations`
- New struct `SecurityPolicyConfigurationsUpdate`
- New struct `SecurityPolicyListResult`
- New struct `SecurityPolicyProperties`
- New struct `SecurityPolicyUpdate`
- New struct `SecurityPolicyUpdateProperties`
- New struct `TrafficControllerUpdateProperties`
- New struct `WafPolicy`
- New struct `WafPolicyUpdate`
- New struct `WafSecurityPolicy`
- New struct `WafSecurityPolicyUpdate`
- New field `SecurityPolicies`, `SecurityPolicyConfigurations` in struct `TrafficControllerProperties`
- New field `Properties` in struct `TrafficControllerUpdate`


## 1.0.0 (2023-11-24)
### Features Added

- Support for test fakes and OpenTelemetry trace spans.

### Other Changes

- Release stable version.


## 0.3.0 (2023-05-26)
### Breaking Changes

- Type of `AssociationProperties.AssociationType` has been changed from `*string` to `*AssociationType`
- Type of `AssociationUpdateProperties.AssociationType` has been changed from `*string` to `*AssociationType`
- Type of `AssociationUpdateProperties.Subnet` has been changed from `*AssociationSubnet` to `*AssociationSubnetUpdate`
- Enum `FrontendIPAddressVersion` has been removed
- Struct `FrontendPropertiesIPAddress` has been removed
- Struct `FrontendUpdateProperties` has been removed
- Field `IPAddressVersion`, `Mode`, `PublicIPAddress` of struct `FrontendProperties` has been removed
- Field `Properties` of struct `FrontendUpdate` has been removed
- Field `Properties` of struct `TrafficControllerUpdate` has been removed

### Features Added

- New enum type `AssociationType` with values `AssociationTypeSubnets`
- New struct `AssociationSubnetUpdate`
- New field `Fqdn` in struct `FrontendProperties`


## 0.2.1 (2023-04-14)
### Bug Fixes

- Fix serialization bug of empty value of `any` type.


## 0.2.0 (2023-03-31)
### Features Added

- New struct `ClientFactory` which is a client factory used to create any client in this module


## 0.1.0 (2023-01-11)
### Other Changes

The package of `github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicenetworking/armservicenetworking` is using our [next generation design principles](https://azure.github.io/azure-sdk/general_introduction.html).

To learn more, please refer to our documentation [Quick Start](https://aka.ms/azsdk/go/mgmt).
