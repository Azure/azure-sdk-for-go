# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)]
  [string] $SourceRoot,
  [Parameter(Mandatory = $true)]
  [string] $PackageName,
  [Parameter(Mandatory = $true)]
  [string] $OutputDirectory,
  [Parameter(Mandatory = $true)]
  [string] $ApigenPath
)

Set-StrictMode -Version 3
$ErrorActionPreference = 'Stop'

$packageName = $PackageName.Replace('\', '/')
$parts = $packageName.Split('/')
if ($packageName -cnotmatch '^(?:sdk|profile)(?:/[A-Za-z0-9._-]+)+$' -or
    @($parts | Where-Object { $_ -in @('.', '..') }).Count -gt 0) {
  throw "PackageName must be a repository-relative SDK module path: $PackageName"
}

$sourceRootPath = [IO.Path]::GetFullPath($SourceRoot)
$packagePath = $sourceRootPath
foreach ($part in $parts) {
  $packagePath = Join-Path $packagePath $part
}
if (-not (Test-Path -LiteralPath (Join-Path $packagePath 'go.mod') -PathType Leaf)) {
  throw "Go module not found for $PackageName at $packagePath."
}
if (-not (Test-Path -LiteralPath $ApigenPath -PathType Leaf)) {
  throw "apigen executable not found: $ApigenPath"
}

New-Item -ItemType Directory -Path $OutputDirectory -Force | Out-Null
& $ApigenPath $packagePath --repo-root $sourceRootPath --output $OutputDirectory | Out-Host
if ($LASTEXITCODE -ne 0) {
  throw "apigen failed for $PackageName with exit code $LASTEXITCODE."
}

$apiPath = Join-Path $OutputDirectory 'api.md'
$metadataPath = Join-Path $OutputDirectory 'api.metadata.yml'
if (-not (Test-Path -LiteralPath $apiPath -PathType Leaf) -or
    -not (Test-Path -LiteralPath $metadataPath -PathType Leaf)) {
  throw "apigen did not produce both API review files for $PackageName."
}
$metadata = Get-Content -LiteralPath $metadataPath -Raw
$hashMatch = [regex]::Match($metadata, '(?m)^apiHash: ([0-9a-f]{64})\r?$')
$versionMatch = [regex]::Match($metadata, '(?m)^packageVersion: (\S+)\r?$')
$pathMatch = [regex]::Match($metadata, '(?m)^packageRelativePath: (".*")\r?$')
if (-not $hashMatch.Success -or -not $versionMatch.Success -or -not $pathMatch.Success) {
  throw "apigen produced invalid metadata for $PackageName."
}
$metadataPathValue = $pathMatch.Groups[1].Value | ConvertFrom-Json -ErrorAction Stop
if ($metadataPathValue -cne $packageName) {
  throw "apigen returned package path '$metadataPathValue' instead of '$packageName'."
}
$hash = $hashMatch.Groups[1].Value
if ((Get-FileHash -LiteralPath $apiPath -Algorithm SHA256).Hash.ToLowerInvariant() -cne $hash) {
  throw "apigen returned an API hash that does not match api.md for $PackageName."
}

[PSCustomObject]@{
  ApiHash = $hash
  PackageVersion = $versionMatch.Groups[1].Value
}
