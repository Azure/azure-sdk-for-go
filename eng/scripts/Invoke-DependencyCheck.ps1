Param(
    [string] $PackageDirectory
)

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

$sdks = Get-AllPackageInfoFromRepo

## Use local git to fetch tags without GitHub token
$existingTags = git tag

## Create depcheck module
$workingPath = Join-Path $RepoRoot "sdk" "depcheck"
if (Test-Path -Path $workingPath)
{
    Remove-Item -Path $workingPath -Recurse -Force
}
New-Item -ItemType Directory -Force -Path $workingPath

## Init go mod
Set-Location $workingPath
Write-Host "##[command]Executing go mod init in " $workingPath
go mod init github.com/Azure/azure-sdk-for-go/sdk/depcheck
if ($LASTEXITCODE) { exit $LASTEXITCODE }

# Find whether latest version is in the `retract` section in `go.mod` to judge whether the package is temporarily deprecated. 
function IsPackageDeprecated($sdk)
{
    $RETRACT_SECTION_REGEX = "retract\s*((?<retract>(.|\s)*))"
    $DEPRECATE_SECTION_REGEX = "\/\/\s*Deprecated:"
    $modContent = Get-Content (Join-Path $sdk.DirectoryPath 'go.mod') -Raw
    if ($modContent -match $DEPRECATE_SECTION_REGEX)
    {
        return $true
    }
    if ($modContent -match $RETRACT_SECTION_REGEX)
    {
        return $($matches["retract"]).Indexof($sdk.Version) -ne -1
    }
    return $false
}

# Get all existed packages
$packagesImport = ""
foreach ($sdk in $sdks)
{
    if ($sdk.Name -like "*internal*" -or (IsPackageDeprecated $sdk))
    {
        continue
    }
    $parsedSemver = [AzureEngSemanticVersion]::ParseVersionString($sdk.Version)

    
    if ($sdk.Name -eq $PackageDirectory -or $existingTags -notcontains "$($sdk.Name)/v$($sdk.Version)")
    {
        ## Add replace for new package
        $modPath = Join-Path $RepoRoot "sdk" "depcheck" "go.mod"
        if ($parsedSemver.Major -gt 1)
        {
            Add-Content $modPath "`nreplace github.com/Azure/azure-sdk-for-go/$($sdk.Name)/v$($parsedSemver.Major) => ../../$($sdk.Name)`n"
        }
        else
        {
            Add-Content $modPath "`nreplace github.com/Azure/azure-sdk-for-go/$($sdk.Name) => ../../$($sdk.Name)`n"
        }
    }

    if ($parsedSemver.Major -gt 1)
    {
        $packagesImport = $packagesImport + "`t_ `"github.com/Azure/azure-sdk-for-go/$($sdk.Name)/v$($parsedSemver.Major)`"`n"
    }
    else
    {
        $packagesImport = $packagesImport + "`t_ `"github.com/Azure/azure-sdk-for-go/$($sdk.Name)`"`n"
    }
}

## Add main.go
$mainPath = Join-Path $RepoRoot "sdk" "depcheck" "main.go"
New-Item -Path $mainPath -ItemType File -Value '' -Force
Add-Content $mainPath "package main

import (
$packagesImport
)

func main() {
}
"

## Run go mod tidy
Write-Host "##[command]Executing go mod tidy in " $workingPath
go mod tidy
if ($LASTEXITCODE) { exit $LASTEXITCODE }

## Run go build and go vet
Write-Host "##[command]Executing go build -v ./... in " $workingPath
go build -v ./...
if ($LASTEXITCODE) { exit $LASTEXITCODE }

Write-Host "##[command]Executing go vet ./... in " $workingPath
go vet ./...
if ($LASTEXITCODE) { exit $LASTEXITCODE }

Write-Host "Checking dependency has completed. All packages are compatible."
