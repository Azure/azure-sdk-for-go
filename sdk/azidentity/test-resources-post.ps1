# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

# IMPORTANT: Do not invoke this file directly. Please instead run eng/common/TestResources/New-TestResources.ps1 from the repository root.

param (
  [hashtable] $AdditionalParameters = @{},
  [hashtable] $DeploymentOutputs,

  [Parameter(Mandatory = $true)]
  [ValidateNotNullOrEmpty()]
  [string] $SubscriptionId,

  [Parameter(ParameterSetName = 'Provisioner', Mandatory = $true)]
  [ValidateNotNullOrEmpty()]
  [string] $TenantId,

  [Parameter()]
  [ValidatePattern('^[0-9a-f]{8}(-[0-9a-f]{4}){3}-[0-9a-f]{12}$')]
  [string] $TestApplicationId,

  [Parameter(Mandatory = $true)]
  [ValidateNotNullOrEmpty()]
  [string] $Environment,

  # Captures any arguments from eng/New-TestResources.ps1 not declared here (no parameter errors).
  [Parameter(ValueFromRemainingArguments = $true)]
  $RemainingArguments
)

$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $true

if ($CI) {
  if (!$AdditionalParameters['deployResources']) {
    Write-Host "Skipping post-provisioning script because resources weren't deployed"
    return
  }
  az cloud set -n $Environment
  az login --federated-token $env:ARM_OIDC_TOKEN --service-principal -t $TenantId -u $TestApplicationId
  az account set --subscription $SubscriptionId
}

Write-Host "Building container"
$image = "azidentity-managed-id-test"
Set-Content -Path "$PSScriptRoot/Dockerfile" -Value @"
FROM mcr.microsoft.com/oss/go/microsoft/golang:latest AS builder
ENV GOARCH=amd64 GOWORK=off
COPY . /azidentity
WORKDIR /azidentity/testdata/managed-id-test
RUN go mod tidy
RUN go build -o /build/managed-id-test .
RUN GOOS=windows go build -o /build/managed-id-test.exe .

FROM mcr.microsoft.com/mirror/docker/library/alpine:3.16
RUN apk add gcompat
COPY --from=builder /build/* .
RUN chmod +x managed-id-test
CMD ["./managed-id-test"]
"@
# build from sdk/azidentity because we need that dir in the context (because the test app uses local azidentity)
docker build -t $image "$PSScriptRoot"

$rg = $DeploymentOutputs['AZIDENTITY_RESOURCE_GROUP']

# Azure Functions deployment: copy the Windows binary from the Docker image, deploy it in a zip
Write-Host "Deploying to Azure Functions"
$container = docker create $image
docker cp ${container}:managed-id-test.exe "$PSScriptRoot/testdata/managed-id-test/"
docker rm -v $container
Compress-Archive -Path "$PSScriptRoot/testdata/managed-id-test/*" -DestinationPath func.zip -Force
az functionapp deploy -g $rg -n $DeploymentOutputs['AZIDENTITY_FUNCTION_NAME'] --src-path func.zip --type zip
