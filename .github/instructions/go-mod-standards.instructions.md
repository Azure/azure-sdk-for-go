---
applyTo: '**/go.mod'
---

- go.mod should only have direct reference to modules within the azure-sdk-for-go (ie: github.com/Azure/azure-sdk-for-go/sdk/<modules>), standard library modules or modules that begin with `golang.org/x/`
- go.mod can have indirect dependencies to 3rd party modules, with the exception of:
  - github.com/stretchr/testify
  
