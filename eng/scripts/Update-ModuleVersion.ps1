<#
.SYNOPSIS
Bumps up package versions after release

.PARAMETER ModulePath
The full module path for a package '/sdk/<module>' or `/sdk/<group>/<module>`

.PARAMETER NewVersionString
New version string to use. Must follow SemVer conventions.

.DESCRIPTION
This script bumps up package versions following conventions defined at https://github.com/Azure/azure-sdk/blob/main/docs/policies/releases.md#incrementing-after-release-go
#>

[CmdletBinding()]
Param (
  [Parameter(Mandatory=$True)]
  [string] $ModulePath,
  [string] $NewVersionString,
  [string] $ReleaseDate,
  [boolean] $ReplaceLatestEntryTitle=$false
 )

. (Join-Path $PSScriptRoot ".." common scripts common.ps1)

$moduleProperties = Get-PkgProperties -PackageName $ModulePath

if (!$moduleProperties)
{
  Write-Error "Could not find properties for module $ModulePath!"
  exit 1
}

$incrementVersion = !$NewVersionString;
$semVer = [AzureEngSemanticVersion]::ParseVersionString($moduleProperties.Version);

if ($incrementVersion) {
  if (!$semVer) {
    LogError "Could not parse existing version $($moduleProperties.Version)"
    exit 1
  }

  if ($semVer.PrereleaseLabel -ne "zzz") {
    $semVer.PrereleaseNumber++
  }
  else {
    $semVer.Patch++
  }
}
else {
  $semVer = [AzureEngSemanticVersion]::ParseVersionString($NewVersionString)
  if (!$semVer) {
    LogError "Could not parse given version $NewVersionString"
    exit 1
  }
}

Write-Verbose "New Version: $semVer"

# Update version in version file.
$versionFileContent = Get-Content -Path $moduleProperties.VersionFile -Raw
$newVersionFileContent = $versionFileContent -replace $moduleProperties.Version, $semVer.ToString()
$newVersionFileContent | Set-Content -Path $moduleProperties.VersionFile -NoNewline

$unreleased = $incrementVersion
$updateExisting = !$incrementVersion

if ($ReplaceLatestEntryTitle) {
  $updateExisting = $ReplaceLatestEntryTitle
}

# Update change log entry
& "${RepoRoot}/eng/common/scripts/Update-ChangeLog.ps1" `
  -Version $semVer.ToString() `
  -ChangelogPath $moduleProperties.ChangeLogPath `
  -Unreleased $unreleased `
  -ReplaceLatestEntryTitle $updateExisting `
  -ReleaseDate $releaseDate

exit 0
