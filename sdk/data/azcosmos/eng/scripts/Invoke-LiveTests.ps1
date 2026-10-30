# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

#Requires -Version 7.0
param(
    [Parameter(Mandatory)]
    [string] $OutputDirectory
)

$ErrorActionPreference = 'Stop'
$module = (Resolve-Path "$PSScriptRoot/../..").Path
New-Item -ItemType Directory -Path $OutputDirectory -Force | Out-Null
$OutputDirectory = (Resolve-Path -LiteralPath $OutputDirectory).Path
$events = Join-Path $OutputDirectory 'events.jsonl'
$plan = Join-Path $OutputDirectory 'selection.json'
# Never validate stale results, including when compilation or setup fails before emitting a manifest.
foreach ($file in @($events, $plan, (Join-Path $OutputDirectory 'report.xml'),
    (Join-Path $OutputDirectory 'summary.json'), (Join-Path $OutputDirectory 'metadata.json'))) {
    if (Test-Path -LiteralPath $file) {
        Remove-Item -LiteralPath $file
    }
}

Push-Location $module
try {
    $target = (& go env -json GOOS GOARCH CGO_ENABLED) | ConvertFrom-Json
    $hostOS = [System.Runtime.InteropServices.RuntimeInformation]::IsOSPlatform(
        [System.Runtime.InteropServices.OSPlatform]::Linux) ? 'linux' : (
        [System.Runtime.InteropServices.RuntimeInformation]::IsOSPlatform(
            [System.Runtime.InteropServices.OSPlatform]::OSX) ? 'darwin' : 'unsupported')
    $hostArchitecture = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture
    if ($LASTEXITCODE -ne 0 -or $target.CGO_ENABLED -ne '1' -or
        !(($target.GOOS -eq 'linux' -and $target.GOARCH -eq 'amd64') -or
          ($target.GOOS -eq 'darwin' -and $target.GOARCH -eq 'arm64')) -or
        $target.GOOS -ne $hostOS -or
        ($target.GOARCH -eq 'amd64' -and $hostArchitecture -ne 'X64') -or
        ($target.GOARCH -eq 'arm64' -and $hostArchitecture -ne 'Arm64')) {
        throw 'Azure Live requires an executable native target: cgo Linux amd64 or macOS arm64.'
    }
    if ($env:AZURE_COSMOS_E2E_BACKEND -ne 'azureLive' -or $env:AZURE_RECORD_MODE -ne 'live') {
        throw 'Explicit AZURE_COSMOS_E2E_BACKEND=azureLive and AZURE_RECORD_MODE=live are required.'
    }
    $driverModule = "github.com/Azure/azure-cosmos-driver/$($target.GOOS)/$($target.GOARCH)"
    $driver = (& go list -m -json $driverModule) | ConvertFrom-Json
    if ($LASTEXITCODE -ne 0) { throw 'Cannot identify the native driver dependency.' }
    $commit = & git rev-parse HEAD
    if ($LASTEXITCODE -ne 0) { throw 'Cannot identify the SDK source revision.' }
    $changes = & git status --porcelain --untracked-files=normal -- .
    if ($LASTEXITCODE -ne 0) { throw 'Cannot identify uncommitted SDK changes.' }
    $goVersion = & go env GOVERSION
    if ($LASTEXITCODE -ne 0) { throw 'Cannot identify the Go version.' }
    @{
        sdkCommit = $commit
        sdkDirty = [bool]$changes
        nativeDriver = $driver.Version
        goVersion = $goVersion
        os = $target.GOOS
        architecture = $target.GOARCH
        backend = 'azureLive'
        authentication = 'entra'
        expectedAccountConsistency = 'Session'
        expectedAccountRegions = @('East US')
        expectedAccountWriteMode = 'single'
        buildId = $env:BUILD_BUILDID
    } | ConvertTo-Json | Set-Content -LiteralPath (Join-Path $OutputDirectory 'metadata.json')
    $oldPlan = $env:AZCOSMOS_LIVE_PLAN_FILE
    $env:AZCOSMOS_LIVE_PLAN_FILE = $plan
    try {
        & go test -json -race -count=1 -timeout=12m -run '^TestAzureLive$' ./e2e |
            Tee-Object -FilePath $events | ForEach-Object { Write-Host $_ }
        $testExitCode = $LASTEXITCODE
    }
    finally {
        $env:AZCOSMOS_LIVE_PLAN_FILE = $oldPlan
    }
    & "$PSScriptRoot/Test-LiveResults.ps1" -EventsPath $events -PlanPath $plan -OutputDirectory $OutputDirectory
    if ($testExitCode -ne 0) {
        throw "Live test command failed with exit code $testExitCode."
    }
}
finally {
    Pop-Location
}
