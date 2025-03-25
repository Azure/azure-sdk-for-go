#!/usr/bin/env pwsh
param (
    [ValidateSet("remote", "local")]
    [string] $What = "remote"
)

Set-StrictMode -Version Latest
$PSNativeCommandUseErrorActionPreference = $true
$ErrorActionPreference = 'Stop'

function deployUsingLocalAddons() {
    $azureSDKToolsRoot = "<Git clone of azure-sdk-tools>"
    $stressTestAddonsFolder = "$azureSDKToolsRoot/tools/stress-cluster/cluster/kubernetes/stress-test-addons"
    $clusterResourceGroup = "<Resource Group for Cluster>"
    $clusterSubscription = "<Azure Subscription>"
    $helmEnv = "pg2"

    if (-not (Get-ChildItem $stressTestAddonsFolder)) {
        Write-Host "Can't find the the new stress test addons folder at $stressTestAddonsFolder"
        return
    }

    pwsh "$azureSDKToolsRoot/eng/common/scripts/stress-testing/deploy-stress-tests.ps1" `
        -LocalAddonsPath "$stressTestAddonsFolder"  `
        -clusterGroup "$clusterResourceGroup" `
        -subscription "$clusterSubscription" `
        -Environment $helmEnv
}

switch ($What) {
    "remote" {
        Set-Location $PSScriptRoot
        
        #deployUsingLocalAddons
        $gitRoot = git rev-parse --show-toplevel
        pwsh "$gitRoot/eng/common/scripts/stress-testing/deploy-stress-tests.ps1" @args
    }
    "local" {
        Set-Location $PSScriptRoot

        if (-not (Get-ChildItem -Hidden ".env")) {
            Write-Host "Can't find the .env file"
            return
        }

        $matrixFile = "$PSScriptRoot/scenarios-matrix.yaml"
        $imageBuildDir = (Select-String -Path $matrixFile -Pattern 'imageBuildDir' | ForEach-Object { $_.Line.Split(':')[1].Trim() }).Trim('"')

        Write-Output "Using '$imageBuildDir' as the build directory"
        docker build --no-cache -f Dockerfile -t tmp-sb-stress $imageBuildDir

        # TODO: unfortunately, this doesn't quite work because the container doesn't have rights to authenticate using the DefaultAzureCredential.
        # docker run -it `
        #     -e "ENV_FILE=/app/.env" `
        #     -v "$PSScriptRoot/.env:/app/.env" `
        #     tmp-sb-stress "./sb-stress finiteSendAndReceive"
    }
    default {
        Write-Host "Invalid option selected"
        exit 1
    }
}

