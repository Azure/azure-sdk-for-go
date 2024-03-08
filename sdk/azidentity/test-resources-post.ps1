# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

# IMPORTANT: Do not invoke this file directly. Please instead run eng/common/TestResources/New-TestResources.ps1 from the repository root.

param (
    [hashtable] $DeploymentOutputs
)

$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $true

Write-Host "Building container"
$image = "azidentity-managed-id-test"
Set-Content -Path "$PSScriptRoot/Dockerfile" -Value @"
FROM mcr.microsoft.com/oss/go/microsoft/golang:latest as builder
ENV GOARCH=amd64 GOWORK=off
COPY . /azidentity
WORKDIR /azidentity/testdata/managed-id-test
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

az login --service-principal -u $DeploymentOutputs['AZIDENTITY_CLIENT_ID'] -p $DeploymentOutputs['AZIDENTITY_CLIENT_SECRET'] --tenant $DeploymentOutputs['AZIDENTITY_TENANT_ID']
az account set --subscription $DeploymentOutputs['AZIDENTITY_SUBSCRIPTION_ID']

# Azure Functions deployment: copy the Windows binary from the Docker image, deploy it in a zip
Write-Host "Deploying to Azure Functions"
$container = docker create $image
docker cp ${container}:managed-id-test.exe "$PSScriptRoot/testdata/managed-id-test/"
docker rm -v $container
Compress-Archive -Path "$PSScriptRoot/testdata/managed-id-test/*" -DestinationPath func.zip -Force
az functionapp deploy -g $DeploymentOutputs['AZIDENTITY_RESOURCE_GROUP'] -n $DeploymentOutputs['AZIDENTITY_FUNCTION_NAME'] --src-path func.zip --type zip

az logout
