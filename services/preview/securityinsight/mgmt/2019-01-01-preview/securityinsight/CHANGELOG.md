# Unreleased

## Additive Changes

### New Funcs

1. DomainWhoisClient.Get(context.Context, string, string) (EnrichmentDomainWhois, error)
1. DomainWhoisClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. DomainWhoisClient.GetResponder(*http.Response) (EnrichmentDomainWhois, error)
1. DomainWhoisClient.GetSender(*http.Request) (*http.Response, error)
1. IPGeodataClient.Get(context.Context, string, string) (EnrichmentIPGeodata, error)
1. IPGeodataClient.GetPreparer(context.Context, string, string) (*http.Request, error)
1. IPGeodataClient.GetResponder(*http.Response) (EnrichmentIPGeodata, error)
1. IPGeodataClient.GetSender(*http.Request) (*http.Response, error)
1. NewDomainWhoisClient(string) DomainWhoisClient
1. NewDomainWhoisClientWithBaseURI(string, string) DomainWhoisClient
1. NewIPGeodataClient(string) IPGeodataClient
1. NewIPGeodataClientWithBaseURI(string, string) IPGeodataClient

### Struct Changes

#### New Structs

1. DomainWhoisClient
1. EnrichmentDomainWhois
1. EnrichmentDomainWhoisContact
1. EnrichmentDomainWhoisContacts
1. EnrichmentDomainWhoisDetails
1. EnrichmentDomainWhoisRegistrarDetails
1. EnrichmentIPGeodata
1. IPGeodataClient
