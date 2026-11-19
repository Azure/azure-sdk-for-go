# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)]
  [string] $PackageInfoFile,

  [Parameter(Mandatory = $true)]
  [string] $SourceRoot,

  [Parameter(Mandatory = $true)]
  [string] $ToolsDirectory,

  [Parameter(Mandatory = $true)]
  [string] $ParserPath,

  [Parameter(Mandatory = $true)]
  [string] $OutputDirectory
)

Set-StrictMode -Version 3
$ErrorActionPreference = 'Stop'

$packageInfo = Get-Content -LiteralPath $PackageInfoFile -Raw | ConvertFrom-Json
foreach ($field in @('Name', 'Version', 'ReleaseStatus')) {
  if (-not $packageInfo.PSObject.Properties[$field] -or
      [string]::IsNullOrWhiteSpace([string] $packageInfo.$field)) {
    throw "Package-info file $PackageInfoFile must contain a non-empty $field."
  }
}
if ($packageInfo.ReleaseStatus -eq 'Unreleased') {
  throw "Package $($packageInfo.Name) is not marked for release."
}

$bundle = & (Join-Path $ToolsDirectory 'Export-APIReviewBundle.ps1') `
  -SourceRoot $SourceRoot `
  -PackageName $packageInfo.Name `
  -OutputDirectory $OutputDirectory `
  -ParserPath $ParserPath
if ($bundle.PackageVersion -cne $packageInfo.Version) {
  throw "Package-info version $($packageInfo.Version) does not match source version $($bundle.PackageVersion) for $($packageInfo.Name)."
}
if ($bundle.ApiHash -cnotmatch '^[0-9a-f]{64}$') {
  throw "The generated API hash is invalid for $($packageInfo.Name)."
}

$packageInfo | Add-Member -NotePropertyName ApiHash -NotePropertyValue $bundle.ApiHash -Force
$packageInfo | ConvertTo-Json -Depth 100 | Set-Content -LiteralPath $PackageInfoFile -Encoding utf8
Write-Host "Prepared API hash $($bundle.ApiHash) for $($packageInfo.Name) $($packageInfo.Version)."
