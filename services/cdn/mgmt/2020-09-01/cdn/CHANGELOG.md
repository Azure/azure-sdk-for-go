Generated from https://github.com/Azure/azure-rest-api-specs/tree/3a3a9452f965a227ce43e6b545035b99dd175f23/specification/cdn/resource-manager/readme.md tag: `package-2020-09`

Code generator @microsoft.azure/autorest.go@2.1.168

## Breaking Changes

### Removed Funcs

1. *AFDCustomDomainsCreateFuture.Result(AFDCustomDomainsClient) (AFDDomain, error)
1. *AFDCustomDomainsDeleteFuture.Result(AFDCustomDomainsClient) (autorest.Response, error)
1. *AFDCustomDomainsRefreshValidationTokenFuture.Result(AFDCustomDomainsClient) (ValidationToken, error)
1. *AFDCustomDomainsUpdateFuture.Result(AFDCustomDomainsClient) (AFDDomain, error)
1. *AFDEndpointsCreateFuture.Result(AFDEndpointsClient) (AFDEndpoint, error)
1. *AFDEndpointsDeleteFuture.Result(AFDEndpointsClient) (autorest.Response, error)
1. *AFDEndpointsPurgeContentFuture.Result(AFDEndpointsClient) (autorest.Response, error)
1. *AFDEndpointsUpdateFuture.Result(AFDEndpointsClient) (AFDEndpoint, error)
1. *AFDOriginGroupsCreateFuture.Result(AFDOriginGroupsClient) (AFDOriginGroup, error)
1. *AFDOriginGroupsDeleteFuture.Result(AFDOriginGroupsClient) (autorest.Response, error)
1. *AFDOriginGroupsUpdateFuture.Result(AFDOriginGroupsClient) (AFDOriginGroup, error)
1. *AFDOriginsCreateFuture.Result(AFDOriginsClient) (AFDOrigin, error)
1. *AFDOriginsDeleteFuture.Result(AFDOriginsClient) (autorest.Response, error)
1. *AFDOriginsUpdateFuture.Result(AFDOriginsClient) (AFDOrigin, error)
1. *CustomDomainsCreateFuture.Result(CustomDomainsClient) (CustomDomain, error)
1. *CustomDomainsDeleteFuture.Result(CustomDomainsClient) (CustomDomain, error)
1. *EndpointsCreateFuture.Result(EndpointsClient) (Endpoint, error)
1. *EndpointsDeleteFuture.Result(EndpointsClient) (autorest.Response, error)
1. *EndpointsLoadContentFuture.Result(EndpointsClient) (autorest.Response, error)
1. *EndpointsPurgeContentFuture.Result(EndpointsClient) (autorest.Response, error)
1. *EndpointsStartFuture.Result(EndpointsClient) (Endpoint, error)
1. *EndpointsStopFuture.Result(EndpointsClient) (Endpoint, error)
1. *EndpointsUpdateFuture.Result(EndpointsClient) (Endpoint, error)
1. *OriginGroupsCreateFuture.Result(OriginGroupsClient) (OriginGroup, error)
1. *OriginGroupsDeleteFuture.Result(OriginGroupsClient) (autorest.Response, error)
1. *OriginGroupsUpdateFuture.Result(OriginGroupsClient) (OriginGroup, error)
1. *OriginsCreateFuture.Result(OriginsClient) (Origin, error)
1. *OriginsDeleteFuture.Result(OriginsClient) (autorest.Response, error)
1. *OriginsUpdateFuture.Result(OriginsClient) (Origin, error)
1. *PoliciesCreateOrUpdateFuture.Result(PoliciesClient) (WebApplicationFirewallPolicy, error)
1. *PoliciesUpdateFuture.Result(PoliciesClient) (WebApplicationFirewallPolicy, error)
1. *ProfilesCreateFuture.Result(ProfilesClient) (Profile, error)
1. *ProfilesDeleteFuture.Result(ProfilesClient) (autorest.Response, error)
1. *ProfilesUpdateFuture.Result(ProfilesClient) (Profile, error)
1. *RoutesCreateFuture.Result(RoutesClient) (Route, error)
1. *RoutesDeleteFuture.Result(RoutesClient) (autorest.Response, error)
1. *RoutesUpdateFuture.Result(RoutesClient) (Route, error)
1. *RuleSetsCreateFuture.Result(RuleSetsClient) (RuleSet, error)
1. *RuleSetsDeleteFuture.Result(RuleSetsClient) (autorest.Response, error)
1. *RulesCreateFuture.Result(RulesClient) (Rule, error)
1. *RulesDeleteFuture.Result(RulesClient) (autorest.Response, error)
1. *RulesUpdateFuture.Result(RulesClient) (Rule, error)
1. *SecretsCreateFuture.Result(SecretsClient) (Secret, error)
1. *SecretsDeleteFuture.Result(SecretsClient) (autorest.Response, error)
1. *SecretsUpdateFuture.Result(SecretsClient) (Secret, error)
1. *SecurityPoliciesCreateFuture.Result(SecurityPoliciesClient) (SecurityPolicy, error)
1. *SecurityPoliciesDeleteFuture.Result(SecurityPoliciesClient) (autorest.Response, error)
1. *SecurityPoliciesPatchFuture.Result(SecurityPoliciesClient) (SecurityPolicy, error)

## Struct Changes

### Removed Struct Fields

1. AFDCustomDomainsCreateFuture.azure.Future
1. AFDCustomDomainsDeleteFuture.azure.Future
1. AFDCustomDomainsRefreshValidationTokenFuture.azure.Future
1. AFDCustomDomainsUpdateFuture.azure.Future
1. AFDEndpointsCreateFuture.azure.Future
1. AFDEndpointsDeleteFuture.azure.Future
1. AFDEndpointsPurgeContentFuture.azure.Future
1. AFDEndpointsUpdateFuture.azure.Future
1. AFDOriginGroupsCreateFuture.azure.Future
1. AFDOriginGroupsDeleteFuture.azure.Future
1. AFDOriginGroupsUpdateFuture.azure.Future
1. AFDOriginsCreateFuture.azure.Future
1. AFDOriginsDeleteFuture.azure.Future
1. AFDOriginsUpdateFuture.azure.Future
1. CustomDomainsCreateFuture.azure.Future
1. CustomDomainsDeleteFuture.azure.Future
1. EndpointsCreateFuture.azure.Future
1. EndpointsDeleteFuture.azure.Future
1. EndpointsLoadContentFuture.azure.Future
1. EndpointsPurgeContentFuture.azure.Future
1. EndpointsStartFuture.azure.Future
1. EndpointsStopFuture.azure.Future
1. EndpointsUpdateFuture.azure.Future
1. OriginGroupsCreateFuture.azure.Future
1. OriginGroupsDeleteFuture.azure.Future
1. OriginGroupsUpdateFuture.azure.Future
1. OriginsCreateFuture.azure.Future
1. OriginsDeleteFuture.azure.Future
1. OriginsUpdateFuture.azure.Future
1. PoliciesCreateOrUpdateFuture.azure.Future
1. PoliciesUpdateFuture.azure.Future
1. ProfilesCreateFuture.azure.Future
1. ProfilesDeleteFuture.azure.Future
1. ProfilesUpdateFuture.azure.Future
1. RoutesCreateFuture.azure.Future
1. RoutesDeleteFuture.azure.Future
1. RoutesUpdateFuture.azure.Future
1. RuleSetsCreateFuture.azure.Future
1. RuleSetsDeleteFuture.azure.Future
1. RulesCreateFuture.azure.Future
1. RulesDeleteFuture.azure.Future
1. RulesUpdateFuture.azure.Future
1. SecretsCreateFuture.azure.Future
1. SecretsDeleteFuture.azure.Future
1. SecretsUpdateFuture.azure.Future
1. SecurityPoliciesCreateFuture.azure.Future
1. SecurityPoliciesDeleteFuture.azure.Future
1. SecurityPoliciesPatchFuture.azure.Future

## Struct Changes

### New Struct Fields

1. AFDCustomDomainsCreateFuture.Result
1. AFDCustomDomainsCreateFuture.azure.FutureAPI
1. AFDCustomDomainsDeleteFuture.Result
1. AFDCustomDomainsDeleteFuture.azure.FutureAPI
1. AFDCustomDomainsRefreshValidationTokenFuture.Result
1. AFDCustomDomainsRefreshValidationTokenFuture.azure.FutureAPI
1. AFDCustomDomainsUpdateFuture.Result
1. AFDCustomDomainsUpdateFuture.azure.FutureAPI
1. AFDEndpointsCreateFuture.Result
1. AFDEndpointsCreateFuture.azure.FutureAPI
1. AFDEndpointsDeleteFuture.Result
1. AFDEndpointsDeleteFuture.azure.FutureAPI
1. AFDEndpointsPurgeContentFuture.Result
1. AFDEndpointsPurgeContentFuture.azure.FutureAPI
1. AFDEndpointsUpdateFuture.Result
1. AFDEndpointsUpdateFuture.azure.FutureAPI
1. AFDOriginGroupsCreateFuture.Result
1. AFDOriginGroupsCreateFuture.azure.FutureAPI
1. AFDOriginGroupsDeleteFuture.Result
1. AFDOriginGroupsDeleteFuture.azure.FutureAPI
1. AFDOriginGroupsUpdateFuture.Result
1. AFDOriginGroupsUpdateFuture.azure.FutureAPI
1. AFDOriginsCreateFuture.Result
1. AFDOriginsCreateFuture.azure.FutureAPI
1. AFDOriginsDeleteFuture.Result
1. AFDOriginsDeleteFuture.azure.FutureAPI
1. AFDOriginsUpdateFuture.Result
1. AFDOriginsUpdateFuture.azure.FutureAPI
1. CustomDomainsCreateFuture.Result
1. CustomDomainsCreateFuture.azure.FutureAPI
1. CustomDomainsDeleteFuture.Result
1. CustomDomainsDeleteFuture.azure.FutureAPI
1. EndpointsCreateFuture.Result
1. EndpointsCreateFuture.azure.FutureAPI
1. EndpointsDeleteFuture.Result
1. EndpointsDeleteFuture.azure.FutureAPI
1. EndpointsLoadContentFuture.Result
1. EndpointsLoadContentFuture.azure.FutureAPI
1. EndpointsPurgeContentFuture.Result
1. EndpointsPurgeContentFuture.azure.FutureAPI
1. EndpointsStartFuture.Result
1. EndpointsStartFuture.azure.FutureAPI
1. EndpointsStopFuture.Result
1. EndpointsStopFuture.azure.FutureAPI
1. EndpointsUpdateFuture.Result
1. EndpointsUpdateFuture.azure.FutureAPI
1. OriginGroupsCreateFuture.Result
1. OriginGroupsCreateFuture.azure.FutureAPI
1. OriginGroupsDeleteFuture.Result
1. OriginGroupsDeleteFuture.azure.FutureAPI
1. OriginGroupsUpdateFuture.Result
1. OriginGroupsUpdateFuture.azure.FutureAPI
1. OriginsCreateFuture.Result
1. OriginsCreateFuture.azure.FutureAPI
1. OriginsDeleteFuture.Result
1. OriginsDeleteFuture.azure.FutureAPI
1. OriginsUpdateFuture.Result
1. OriginsUpdateFuture.azure.FutureAPI
1. PoliciesCreateOrUpdateFuture.Result
1. PoliciesCreateOrUpdateFuture.azure.FutureAPI
1. PoliciesUpdateFuture.Result
1. PoliciesUpdateFuture.azure.FutureAPI
1. ProfilesCreateFuture.Result
1. ProfilesCreateFuture.azure.FutureAPI
1. ProfilesDeleteFuture.Result
1. ProfilesDeleteFuture.azure.FutureAPI
1. ProfilesUpdateFuture.Result
1. ProfilesUpdateFuture.azure.FutureAPI
1. RoutesCreateFuture.Result
1. RoutesCreateFuture.azure.FutureAPI
1. RoutesDeleteFuture.Result
1. RoutesDeleteFuture.azure.FutureAPI
1. RoutesUpdateFuture.Result
1. RoutesUpdateFuture.azure.FutureAPI
1. RuleSetsCreateFuture.Result
1. RuleSetsCreateFuture.azure.FutureAPI
1. RuleSetsDeleteFuture.Result
1. RuleSetsDeleteFuture.azure.FutureAPI
1. RulesCreateFuture.Result
1. RulesCreateFuture.azure.FutureAPI
1. RulesDeleteFuture.Result
1. RulesDeleteFuture.azure.FutureAPI
1. RulesUpdateFuture.Result
1. RulesUpdateFuture.azure.FutureAPI
1. SecretsCreateFuture.Result
1. SecretsCreateFuture.azure.FutureAPI
1. SecretsDeleteFuture.Result
1. SecretsDeleteFuture.azure.FutureAPI
1. SecretsUpdateFuture.Result
1. SecretsUpdateFuture.azure.FutureAPI
1. SecurityPoliciesCreateFuture.Result
1. SecurityPoliciesCreateFuture.azure.FutureAPI
1. SecurityPoliciesDeleteFuture.Result
1. SecurityPoliciesDeleteFuture.azure.FutureAPI
1. SecurityPoliciesPatchFuture.Result
1. SecurityPoliciesPatchFuture.azure.FutureAPI
