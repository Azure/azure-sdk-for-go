# Unreleased

## Breaking Changes

### Removed Constants

1. Kind.KindAggregations
1. Kind.KindCasesAggregation
1. KindBasicEntity.KindBasicEntityKindAccount
1. KindBasicEntity.KindBasicEntityKindAzureResource
1. KindBasicEntity.KindBasicEntityKindBookmark
1. KindBasicEntity.KindBasicEntityKindCloudApplication
1. KindBasicEntity.KindBasicEntityKindDNSResolution
1. KindBasicEntity.KindBasicEntityKindEntity
1. KindBasicEntity.KindBasicEntityKindFile
1. KindBasicEntity.KindBasicEntityKindFileHash
1. KindBasicEntity.KindBasicEntityKindHost
1. KindBasicEntity.KindBasicEntityKindIP
1. KindBasicEntity.KindBasicEntityKindIoTDevice
1. KindBasicEntity.KindBasicEntityKindMailCluster
1. KindBasicEntity.KindBasicEntityKindMailMessage
1. KindBasicEntity.KindBasicEntityKindMailbox
1. KindBasicEntity.KindBasicEntityKindMalware
1. KindBasicEntity.KindBasicEntityKindProcess
1. KindBasicEntity.KindBasicEntityKindRegistryKey
1. KindBasicEntity.KindBasicEntityKindRegistryValue
1. KindBasicEntity.KindBasicEntityKindSecurityAlert
1. KindBasicEntity.KindBasicEntityKindSecurityGroup
1. KindBasicEntity.KindBasicEntityKindSubmissionMail
1. KindBasicEntity.KindBasicEntityKindURL

### Removed Funcs

1. PossibleKindBasicEntityValues() []KindBasicEntity

### Signature Changes

#### Funcs

1. ThreatIntelligenceIndicatorsClient.List
	- Params
		- From: context.Context, string, string, string, string, *int32, string, string
		- To: context.Context, string, string, string, string, string, *int32, string
1. ThreatIntelligenceIndicatorsClient.ListComplete
	- Params
		- From: context.Context, string, string, string, string, *int32, string, string
		- To: context.Context, string, string, string, string, string, *int32, string
1. ThreatIntelligenceIndicatorsClient.ListPreparer
	- Params
		- From: context.Context, string, string, string, string, *int32, string, string
		- To: context.Context, string, string, string, string, string, *int32, string

#### Struct Fields

1. AccountEntity.Kind changed type from KindBasicEntity to Kind
1. Aggregations.Kind changed type from Kind to KindBasicAggregations
1. AzureResourceEntity.Kind changed type from KindBasicEntity to Kind
1. CasesAggregation.Kind changed type from Kind to KindBasicAggregations
1. CloudApplicationEntity.Kind changed type from KindBasicEntity to Kind
1. DNSEntity.Kind changed type from KindBasicEntity to Kind
1. Entity.Kind changed type from KindBasicEntity to Kind
1. FileEntity.Kind changed type from KindBasicEntity to Kind
1. FileHashEntity.Kind changed type from KindBasicEntity to Kind
1. HostEntity.Kind changed type from KindBasicEntity to Kind
1. HuntingBookmark.Kind changed type from KindBasicEntity to Kind
1. IPEntity.Kind changed type from KindBasicEntity to Kind
1. IncidentInfo.Severity changed type from CaseSeverity to IncidentSeverity
1. IoTDeviceEntity.Kind changed type from KindBasicEntity to Kind
1. MailClusterEntity.Kind changed type from KindBasicEntity to Kind
1. MailMessageEntity.Kind changed type from KindBasicEntity to Kind
1. MailboxEntity.Kind changed type from KindBasicEntity to Kind
1. MalwareEntity.Kind changed type from KindBasicEntity to Kind
1. ProcessEntity.Kind changed type from KindBasicEntity to Kind
1. RegistryKeyEntity.Kind changed type from KindBasicEntity to Kind
1. RegistryValueEntity.Kind changed type from KindBasicEntity to Kind
1. SecurityAlert.Kind changed type from KindBasicEntity to Kind
1. SecurityGroupEntity.Kind changed type from KindBasicEntity to Kind
1. SubmissionMailEntity.Kind changed type from KindBasicEntity to Kind
1. URLEntity.Kind changed type from KindBasicEntity to Kind

## Additive Changes

### New Constants

1. Kind.KindAccount
1. Kind.KindAzureResource
1. Kind.KindBookmark
1. Kind.KindCloudApplication
1. Kind.KindDNSResolution
1. Kind.KindEntity
1. Kind.KindFile
1. Kind.KindFileHash
1. Kind.KindHost
1. Kind.KindIP
1. Kind.KindIoTDevice
1. Kind.KindMailCluster
1. Kind.KindMailMessage
1. Kind.KindMailbox
1. Kind.KindMalware
1. Kind.KindProcess
1. Kind.KindRegistryKey
1. Kind.KindRegistryValue
1. Kind.KindSecurityAlert
1. Kind.KindSecurityGroup
1. Kind.KindSubmissionMail
1. Kind.KindURL
1. KindBasicAggregations.KindBasicAggregationsKindAggregations
1. KindBasicAggregations.KindBasicAggregationsKindCasesAggregation
1. KindBasicSettings.KindBasicSettingsKindAnomalies
1. SettingKind.SettingKindAnomalies

### New Funcs

1. *Anomalies.UnmarshalJSON([]byte) error
1. Anomalies.AsAnomalies() (*Anomalies, bool)
1. Anomalies.AsBasicSettings() (BasicSettings, bool)
1. Anomalies.AsEntityAnalytics() (*EntityAnalytics, bool)
1. Anomalies.AsEyesOn() (*EyesOn, bool)
1. Anomalies.AsIPSyncer() (*IPSyncer, bool)
1. Anomalies.AsSettings() (*Settings, bool)
1. Anomalies.AsUeba() (*Ueba, bool)
1. Anomalies.MarshalJSON() ([]byte, error)
1. AnomaliesProperties.MarshalJSON() ([]byte, error)
1. DomainWhoisClient.Get(context.Context, string, string) (EnrichmentDomainWhois, error)
1. DomainWhoisClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. DomainWhoisClient.GetResponder(*http.Response) (EnrichmentDomainWhois, error)
1. DomainWhoisClient.GetSender(*http.Request) (*http.Response, error)
1. EntityAnalytics.AsAnomalies() (*Anomalies, bool)
1. EyesOn.AsAnomalies() (*Anomalies, bool)
1. IPGeodataClient.Get(context.Context, string, string) (EnrichmentIPGeodata, error)
1. IPGeodataClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. IPGeodataClient.GetResponder(*http.Response) (EnrichmentIPGeodata, error)
1. IPGeodataClient.GetSender(*http.Request) (*http.Response, error)
1. IPSyncer.AsAnomalies() (*Anomalies, bool)
1. NewDomainWhoisClient(string) DomainWhoisClient
1. NewDomainWhoisClientWithBaseURI(string, string) DomainWhoisClient
1. NewIPGeodataClient(string) IPGeodataClient
1. NewIPGeodataClientWithBaseURI(string, string) IPGeodataClient
1. PossibleKindBasicAggregationsValues() []KindBasicAggregations
1. Settings.AsAnomalies() (*Anomalies, bool)
1. Ueba.AsAnomalies() (*Anomalies, bool)

### Struct Changes

#### New Structs

1. Anomalies
1. AnomaliesProperties
1. DomainWhoisClient
1. EnrichmentDomainWhois
1. EnrichmentDomainWhoisContact
1. EnrichmentDomainWhoisContacts
1. EnrichmentDomainWhoisDetails
1. EnrichmentDomainWhoisRegistrarDetails
1. EnrichmentIPGeodata
1. IPGeodataClient
