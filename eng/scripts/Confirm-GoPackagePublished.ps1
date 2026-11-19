# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)]
  [string] $PackageInfoFile,

  [Parameter(Mandatory = $true)]
  [string] $ExpectedCommit
)

Set-StrictMode -Version 3
$ErrorActionPreference = 'Stop'

$packageInfo = Get-Content -LiteralPath $PackageInfoFile -Raw | ConvertFrom-Json
foreach ($field in @('Name', 'Version')) {
  if (-not $packageInfo.PSObject.Properties[$field] -or
      [string]::IsNullOrWhiteSpace([string] $packageInfo.$field)) {
    throw "Package-info file $PackageInfoFile must contain a non-empty $field."
  }
}

$tagRef = "refs/tags/$($packageInfo.Name)/v$($packageInfo.Version)"
$lines = @(& git ls-remote --exit-code 'https://github.com/Azure/azure-sdk-for-go.git' $tagRef "${tagRef}^{}")
if ($LASTEXITCODE -eq 2) {
  throw "The published tag $tagRef was not found."
}
if ($LASTEXITCODE -ne 0) {
  throw "Failed to query published tag $tagRef (git exit code $LASTEXITCODE)."
}

$refs = @{}
foreach ($line in $lines) {
  $parts = $line -split '\s+', 2
  if ($parts.Count -eq 2) {
    $refs[$parts[1]] = $parts[0]
  }
}
$commit = if ($refs.ContainsKey("${tagRef}^{}")) { $refs["${tagRef}^{}"] } else { $refs[$tagRef] }
if ($commit -ne $ExpectedCommit) {
  throw "Published tag $tagRef does not point to the source commit $ExpectedCommit."
}
