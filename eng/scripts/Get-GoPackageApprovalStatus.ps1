<#
.SYNOPSIS
Checks Go package approval before release.

.DESCRIPTION
Requires API approval for stable versions. Prerelease versions may proceed with package-name approval while API review is pending.

.PARAMETER PackageInfoFiles
Package-info JSON files containing the package name and version.

.PARAMETER AzSdkExePath
The path to the azsdk executable.
#>
[CmdletBinding()]
param (
    [Parameter(Mandatory = $true)]
    [ValidateNotNullOrEmpty()]
    [array] $PackageInfoFiles,

    [ValidateNotNullOrEmpty()]
    [string] $AzSdkExePath = "azsdk"
)

Set-StrictMode -Version 4
$ErrorActionPreference = "Stop"

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

if ([string]::IsNullOrWhiteSpace($LanguageShort) -or $LanguageShort -eq "Unknown") {
    throw "SDK language settings were not loaded."
}

$packageInfoPaths = @($PackageInfoFiles | Where-Object { -not [string]::IsNullOrWhiteSpace($_) })
if ($packageInfoPaths.Count -eq 0) {
    throw "At least one package-info file is required."
}

Confirm-AzSdkCliMinimumVersion $AzSdkExePath ([version] "0.6.38")
$failures = @()
foreach ($packageInfoFile in $packageInfoPaths) {
    try {
        if (-not (Test-Path $packageInfoFile -PathType Leaf)) {
            throw "Package-info file does not exist."
        }

        $packageInfo = Get-Content $packageInfoFile -Raw | ConvertFrom-Json -ErrorAction Stop
        $packageName = if ($packageInfo.PSObject.Properties["Name"]) { [string] $packageInfo.Name } else { "" }
        $packageVersion = if ($packageInfo.PSObject.Properties["Version"]) { [string] $packageInfo.Version } else { "" }
        $releaseStatus = if ($packageInfo.PSObject.Properties["ReleaseStatus"]) { [string] $packageInfo.ReleaseStatus } else { "" }

        if ([string]::IsNullOrWhiteSpace($packageName)) {
            throw "Package-info file does not contain a package Name."
        }
        if ([string]::IsNullOrWhiteSpace($packageVersion)) {
            throw "Package-info file does not contain a package Version."
        }

        $version = [AzureEngSemanticVersion]::ParseVersionString($packageVersion)
        if ($null -eq $version) {
            throw "Package version '$packageVersion' is not a valid semantic version."
        }

        $arguments = @(
            "package",
            "get-approval-status",
            "--language", $LanguageShort,
            "--package-name", $packageName,
            "--package-version", $packageVersion,
            "--output", "json"
        )
        Write-Host "Checking package approval: language=$LanguageShort, package=$packageName, version=$packageVersion, apiHash=not provided"

        try {
            $commandResult = Invoke-AzSdkCliCommand $AzSdkExePath $arguments
            $response = $commandResult.Output | ConvertFrom-Json -ErrorAction Stop
            if (-not $response.PSObject.Properties["result"] -or
                $null -eq $response.result -or
                -not $response.result.PSObject.Properties["isApproved"] -or
                $response.result.isApproved -isnot [bool]) {
                throw "Package approval check returned an invalid response."
            }
            $operationStatus = if ($response.PSObject.Properties["operation_status"]) {
                [string] $response.operation_status
            } else {
                ""
            }
            $responseError = if ($response.PSObject.Properties["response_error"]) {
                [string] $response.response_error
            } else {
                ""
            }
            if ([string]::IsNullOrWhiteSpace($responseError) -and
                $response.PSObject.Properties["response_errors"]) {
                $responseError = @($response.response_errors | Where-Object { -not [string]::IsNullOrWhiteSpace($_) }) -join "; "
            }
            $expectedDenial = $commandResult.ExitCode -eq 1 -and
                $operationStatus -eq "Succeeded" -and
                -not $response.result.isApproved -and
                [string]::IsNullOrWhiteSpace($responseError)

            if ($operationStatus -ne "Succeeded" -or
                ($commandResult.ExitCode -ne 0 -and -not $expectedDenial)) {
                $failureMessage = if (-not [string]::IsNullOrWhiteSpace($responseError)) {
                    $responseError
                } else {
                    "azsdk exited with code $($commandResult.ExitCode)."
                }
                throw "Package approval check failed: $failureMessage"
            }
        }
        catch {
            if ($releaseStatus -eq "Unreleased") {
                Write-Host "$packageName $packageVersion is not marked for release. Ignoring approval check failure: $($_.Exception.Message)"
                continue
            }
            throw
        }

        if ($response.result.isApproved) {
            Write-Host "Package $packageName $packageVersion is approved."
            continue
        }

        if ($releaseStatus -eq "Unreleased") {
            Write-Host "$packageName $packageVersion is not marked for release. Ignoring approval check failure."
            continue
        }

        $requiresApiReview = $packageInfo.SdkType -eq "client" -and
            $packageInfo.PSObject.Properties["IsNewSdk"] -and
            [bool] $packageInfo.IsNewSdk
        if (-not $requiresApiReview) {
            Write-Host "API review is not required for $packageName $packageVersion."
            continue
        }

        $apiViewNameApproved = $response.result.PSObject.Properties["apiView"] -and
            $null -ne $response.result.apiView -and
            $response.result.apiView.PSObject.Properties["packageNameApproved"] -and
            $response.result.apiView.packageNameApproved -is [bool] -and
            $response.result.apiView.packageNameApproved
        $reviewHubNameApproved = $response.result.PSObject.Properties["reviewHub"] -and
            $null -ne $response.result.reviewHub -and
            $response.result.reviewHub.PSObject.Properties["isPackageNameApproved"] -and
            $response.result.reviewHub.isPackageNameApproved -is [bool] -and
            $response.result.reviewHub.isPackageNameApproved
        if ($version.IsPrerelease -and ($apiViewNameApproved -or $reviewHubNameApproved)) {
            Write-Host "Package $packageName $packageVersion has an approved package name. API approval is not required for prerelease versions."
            continue
        }

        throw "Package $packageName $packageVersion is not approved for release."
    }
    catch {
        Write-Error "Package approval failed for ${packageInfoFile}: $($_.Exception.Message)" -ErrorAction Continue
        $failures += "${packageInfoFile}: $($_.Exception.Message)"
    }
}

if ($failures.Count -gt 0) {
    throw "Package approval failed for $($failures.Count) package(s):`n$($failures -join "`n")"
}
