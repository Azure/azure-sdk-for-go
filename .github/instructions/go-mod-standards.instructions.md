---
applyTo: '**/go.mod'
---

go.mod should only have direct reference to:
- modules within the azure-sdk-for-go (ie: github.com/Azure/azure-sdk-for-go/sdk/<modules>)
- Go standard library modules or modules that begin with `golang.org/x/`  
- 3rd party modules, but within this list:
  - github.com/stretchr/testify
  - github.com/joho/godotenv
  - go.opentelemetry.io/otel/*
  - github.com/golang/mock

azidentity can directly reference additional 3rd party modules:
- github.com/AzureAD/microsoft-authentication-extensions-for-go/*
- github.com/golang-jwt/jwt/v5
- github.com/google/uuid # this is only used by a test

azservicebus and azeventhubs can directly reference additional 3rd party modules:
- github.com/coder/websocket
- github.com/Azure/go-amqp

azopenai can directly reference additional 3rd party modules:
- github.com/openai/openai-go/v3

