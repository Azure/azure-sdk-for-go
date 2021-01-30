Generated from https://github.com/Azure/azure-rest-api-specs/tree/138759b8a5987252fd66658078907e1d93969c85/specification/cdn/resource-manager/readme.md tag: `package-2020-09`

Code generator @microsoft.azure/autorest.go@2.1.169


## Breaking Changes

### Removed Constants

1. HealthProbeRequestType.GET
1. HealthProbeRequestType.HEAD
1. SecretType.SecretTypeCustomerCertificate
1. SecretType.SecretTypeManagedCertificate
1. SecretType.SecretTypeURLSigningKey

### Removed Funcs

1. PossibleSecretTypeValues() []SecretType

## Struct Changes

### Removed Struct Fields

1. BaseClient.SubscriptionID1
1. CustomerCertificateParameters.ExpirationDate
1. CustomerCertificateParameters.Subject
1. CustomerCertificateParameters.Thumbprint
1. ManagedCertificateParameters.ExpirationDate
1. ManagedCertificateParameters.Subject
1. ManagedCertificateParameters.Thumbprint
1. RouteProperties.OptimizationType
1. RouteUpdatePropertiesParameters.OptimizationType

## Signature Changes

### Const Types

1. NotSet changed type from HealthProbeRequestType to AfdQueryStringCachingBehavior

### Funcs

1. LogAnalyticsClient.GetLogAnalyticsRankings
	- Params
		- From: context.Context, string, string, []string, []string, float64, date.Time, date.Time, []string
		- To: context.Context, string, string, []string, []string, int32, date.Time, date.Time, []string
1. LogAnalyticsClient.GetLogAnalyticsRankingsPreparer
	- Params
		- From: context.Context, string, string, []string, []string, float64, date.Time, date.Time, []string
		- To: context.Context, string, string, []string, []string, int32, date.Time, date.Time, []string
1. LogAnalyticsClient.GetWafLogAnalyticsRankings
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, float64, []string, []string, []string
		- To: context.Context, string, string, []string, date.Time, date.Time, int32, []string, []string, []string
1. LogAnalyticsClient.GetWafLogAnalyticsRankingsPreparer
	- Params
		- From: context.Context, string, string, []string, date.Time, date.Time, float64, []string, []string, []string
		- To: context.Context, string, string, []string, date.Time, date.Time, int32, []string, []string, []string
1. New
	- Params
		- From: string, string
		- To: string
1. NewAFDCustomDomainsClient
	- Params
		- From: string, string
		- To: string
1. NewAFDCustomDomainsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAFDEndpointsClient
	- Params
		- From: string, string
		- To: string
1. NewAFDEndpointsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAFDOriginGroupsClient
	- Params
		- From: string, string
		- To: string
1. NewAFDOriginGroupsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAFDOriginsClient
	- Params
		- From: string, string
		- To: string
1. NewAFDOriginsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewAFDProfilesClient
	- Params
		- From: string, string
		- To: string
1. NewAFDProfilesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewCustomDomainsClient
	- Params
		- From: string, string
		- To: string
1. NewCustomDomainsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewEdgeNodesClient
	- Params
		- From: string, string
		- To: string
1. NewEdgeNodesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewEndpointsClient
	- Params
		- From: string, string
		- To: string
1. NewEndpointsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewLogAnalyticsClient
	- Params
		- From: string, string
		- To: string
1. NewLogAnalyticsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewManagedRuleSetsClient
	- Params
		- From: string, string
		- To: string
1. NewManagedRuleSetsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewOperationsClient
	- Params
		- From: string, string
		- To: string
1. NewOperationsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewOriginGroupsClient
	- Params
		- From: string, string
		- To: string
1. NewOriginGroupsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewOriginsClient
	- Params
		- From: string, string
		- To: string
1. NewOriginsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewPoliciesClient
	- Params
		- From: string, string
		- To: string
1. NewPoliciesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewProfilesClient
	- Params
		- From: string, string
		- To: string
1. NewProfilesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewResourceUsageClient
	- Params
		- From: string, string
		- To: string
1. NewResourceUsageClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewRoutesClient
	- Params
		- From: string, string
		- To: string
1. NewRoutesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewRuleSetsClient
	- Params
		- From: string, string
		- To: string
1. NewRuleSetsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewRulesClient
	- Params
		- From: string, string
		- To: string
1. NewRulesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSecretsClient
	- Params
		- From: string, string
		- To: string
1. NewSecretsClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewSecurityPoliciesClient
	- Params
		- From: string, string
		- To: string
1. NewSecurityPoliciesClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewValidateClient
	- Params
		- From: string, string
		- To: string
1. NewValidateClientWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. NewWithBaseURI
	- Params
		- From: string, string, string
		- To: string, string
1. SecurityPoliciesClient.Patch
	- Params
		- From: context.Context, string, string, string, SecurityPolicyWebApplicationFirewallParameters
		- To: context.Context, string, string, string, SecurityPolicyProperties
1. SecurityPoliciesClient.PatchPreparer
	- Params
		- From: context.Context, string, string, string, SecurityPolicyWebApplicationFirewallParameters
		- To: context.Context, string, string, string, SecurityPolicyProperties

### Struct Fields

1. CustomerCertificateParameters.Type changed type from SecretType to TypeBasicSecretParameters
1. ManagedCertificateParameters.Type changed type from SecretType to TypeBasicSecretParameters
1. RouteProperties.QueryStringCachingBehavior changed type from QueryStringCachingBehavior to AfdQueryStringCachingBehavior
1. RouteUpdatePropertiesParameters.QueryStringCachingBehavior changed type from QueryStringCachingBehavior to AfdQueryStringCachingBehavior
1. SecurityPolicyProperties.Parameters changed type from *SecurityPolicyWebApplicationFirewallParameters to BasicSecurityPolicyParameters

### New Constants

1. AfdQueryStringCachingBehavior.IgnoreQueryString
1. AfdQueryStringCachingBehavior.UseQueryString
1. HealthProbeRequestType.HealthProbeRequestTypeGET
1. HealthProbeRequestType.HealthProbeRequestTypeHEAD
1. HealthProbeRequestType.HealthProbeRequestTypeNotSet
1. TypeBasicSecretParameters.TypeCustomerCertificate
1. TypeBasicSecretParameters.TypeManagedCertificate

### New Funcs

1. *SecurityPolicyProperties.UnmarshalJSON([]byte) error
1. CustomerCertificateParameters.AsBasicSecretParameters() (BasicSecretParameters, bool)
1. CustomerCertificateParameters.AsCustomerCertificateParameters() (*CustomerCertificateParameters, bool)
1. CustomerCertificateParameters.AsManagedCertificateParameters() (*ManagedCertificateParameters, bool)
1. CustomerCertificateParameters.AsSecretParameters() (*SecretParameters, bool)
1. CustomerCertificateParameters.AsURLSigningKeyParameters() (*URLSigningKeyParameters, bool)
1. CustomerCertificateParameters.MarshalJSON() ([]byte, error)
1. ManagedCertificateParameters.AsBasicSecretParameters() (BasicSecretParameters, bool)
1. ManagedCertificateParameters.AsCustomerCertificateParameters() (*CustomerCertificateParameters, bool)
1. ManagedCertificateParameters.AsManagedCertificateParameters() (*ManagedCertificateParameters, bool)
1. ManagedCertificateParameters.AsSecretParameters() (*SecretParameters, bool)
1. ManagedCertificateParameters.AsURLSigningKeyParameters() (*URLSigningKeyParameters, bool)
1. ManagedCertificateParameters.MarshalJSON() ([]byte, error)
1. PossibleAfdQueryStringCachingBehaviorValues() []AfdQueryStringCachingBehavior
1. SecretParameters.AsCustomerCertificateParameters() (*CustomerCertificateParameters, bool)
1. SecretParameters.AsManagedCertificateParameters() (*ManagedCertificateParameters, bool)
1. URLSigningKeyParameters.AsCustomerCertificateParameters() (*CustomerCertificateParameters, bool)
1. URLSigningKeyParameters.AsManagedCertificateParameters() (*ManagedCertificateParameters, bool)
