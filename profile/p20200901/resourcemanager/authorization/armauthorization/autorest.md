### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/53b1affe357b3bfbb53721d0a2002382a046d3b0/specification/authorization/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/53b1affe357b3bfbb53721d0a2002382a046d3b0/specification/authorization/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: 0.1.0
module-name: profile/p20200901/resourcemanager/authorization/armauthorization
module: github.com/Azure/azure-sdk-for-go/$(module-name)
output-folder: $(go-sdk-folder)/$(module-name)
tag: profile-hybrid-2020-09-01

```

### Remove moduleName and moduleVersion constant

```yaml
directive:
  - from: constants.go
    where: $
    transform: return $.replace(/const \(\n\s+moduleName.+\n\s+moduleVersion.+\n\)\n/, "");
```

### Add internal import

```yaml
directive:
  - from:
      - "*_client.go"
      - "client.go"
      - "client_factory.go"
    where: $
    transform: return $.replace(/import \(\n/, "import (\n\"github.com/Azure/azure-sdk-for-go/profile/p20200901/internal\"\n");
```

## Change moduleName and moduleVersion in client CTOR

```yaml
directive:
  - from:
      - "*_client.go"
      - "client.go"
      - "client_factory.go"
    where: $
    transform: return $.replace(/moduleVersion/, "internal.ModuleVersion").replace(/moduleName\+"/, "internal.ModuleName+\"/armauthorization");
```