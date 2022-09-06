### AutoRest Configuration

> see https://aka.ms/autorest

``` yaml
azure-arm: true
require:
- https://github.com/Azure/azure-rest-api-specs/blob/{{commitID}}/specification/{{rpName}}/resource-manager/readme.md
- https://github.com/Azure/azure-rest-api-specs/blob/{{commitID}}/specification/{{rpName}}/resource-manager/readme.go.md
license-header: MICROSOFT_MIT_NO_VERSION
module-version: {{packageVersion}}
module-name: sdk/profiles/{{packageTitle}}/resourcemanager/{{rpName}}/{{packageName}}
module: github.com/Azure/azure-sdk-for-go/$(module-name)
output-folder: $(go-sdk-folder)/$(module-name)
{{packageConfig}}
```