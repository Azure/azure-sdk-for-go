param (
  $PackageName,
  $ServiceDirectory,
  $repoId
)

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

$apiUrl = "https://api.github.com/repos/$repoId"
Write-Host "Using API URL $apiUrl"

# VERIFY CHANGELOG
$PackageProp = Get-PkgProperties -PackageName $PackageName -ServiceDirectory $ServiceDirectory
$changeLogEntries = Get-ChangeLogEntries -ChangeLogLocation $PackageProp.ChangeLogPath
$changeLogEntry = $changeLogEntries[$PackageProp.Version]

if (!$changeLogEntry)
{
  Write-Host "Changelog not existed for package: $PackageName, version: $($PackageProp.Version)."
  Write-Output "##vso[task.setvariable variable=NeedToRelease;isOutput=true]false"
  return
}

if ([System.String]::IsNullOrEmpty($changeLogEntry.ReleaseStatus) -or $changeLogEntry.ReleaseStatus -eq $CHANGELOG_UNRELEASED_STATUS)
{
  Write-Host "Changelog not in release status for package: $PackageName, version: $($PackageProp.Version)."
  Write-Output "##vso[task.setvariable variable=NeedToRelease;isOutput=true]false"
  return
}

# VERIFY TAG
$existingTags = GetExistingTags($apiUrl)
if ($existingTags -contains "$($PackageProp.Name)/v$($PackageProp.Version)")
{
  Write-Host "Package: $PackageName, version: $($PackageProp.Version) already released."
  Write-Output "##vso[task.setvariable variable=NeedToRelease;isOutput=true]false"
}
else
{
  Write-Host "Package: $PackageName, version: $($PackageProp.Version) need to release."
  Write-Output "##vso[task.setvariable variable=NeedToRelease;isOutput=true]true"
}
