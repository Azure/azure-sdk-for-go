### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/4fd842fb73656039ec94ce367bcedee25a57bd18/specification/resources/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/4fd842fb73656039ec94ce367bcedee25a57bd18/specification/resources/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources
module-version: 1.2.0-beta.2
package-resources: true
azcore-version: 1.8.0-beta.1
generate-fakes: true
inject-spans: true
```
