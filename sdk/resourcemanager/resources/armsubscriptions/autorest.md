### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/4f4073bdb028bc84bc3e6405c1cbaf8e89b83caf/specification/resources/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/4f4073bdb028bc84bc3e6405c1cbaf8e89b83caf/specification/resources/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module: github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions
module-version: 1.3.0-beta.2
package-subscriptions: true
tag: package-subscriptions-2022-12
azcore-version: 1.8.0-beta.1
generate-fakes: true
inject-spans: true
```
