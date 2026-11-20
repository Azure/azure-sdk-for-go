# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)]
  [string] $SourceRoot,

  [Parameter(Mandatory = $true)]
  [string] $RepositoryFullName,

  [Parameter(Mandatory = $true)]
  [string] $TargetCommit,

  [Parameter(Mandatory = $true)]
  [ValidateSet('base', 'target')]
  [string] $Kind,

  [Parameter(Mandatory = $true)]
  [string] $ToolsDirectory,

  [Parameter(Mandatory = $true)]
  [string] $WorkingDirectory
)

Set-StrictMode -Version 3
$ErrorActionPreference = 'Stop'

if ($TargetCommit -notmatch '^[0-9a-fA-F]{40}$') {
  throw 'A pinned 40-character target commit SHA is required.'
}
if ($RepositoryFullName -cnotmatch '^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$') {
  throw 'RepositoryFullName must be in owner/repository form.'
}

$targetRoot = $SourceRoot
$temporaryTarget = $null
try {
  if ($Kind -eq 'base') {
    $temporaryTarget = Join-Path $WorkingDirectory ('go-api-review-target-' + [guid]::NewGuid().ToString('N'))
    $targetRoot = $temporaryTarget
    & git init --quiet $targetRoot
    if ($LASTEXITCODE -ne 0) { throw 'Failed to initialize the target tooling checkout.' }
    & git -C $targetRoot remote add origin "https://github.com/$RepositoryFullName.git"
    if ($LASTEXITCODE -ne 0) { throw 'Failed to configure the target tooling repository.' }
    & git -C $targetRoot sparse-checkout set --cone eng
    if ($LASTEXITCODE -ne 0) { throw 'Failed to configure the target tooling sparse checkout.' }
    & git -C $targetRoot fetch --quiet --depth 1 --filter=blob:none origin $TargetCommit
    if ($LASTEXITCODE -ne 0) { throw "Failed to fetch target tooling commit $TargetCommit." }
    & git -C $targetRoot checkout --quiet --detach FETCH_HEAD
    if ($LASTEXITCODE -ne 0) { throw "Failed to check out target tooling commit $TargetCommit." }
  }
  $actualCommit = & git -C $targetRoot rev-parse HEAD
  if ($LASTEXITCODE -ne 0 -or $actualCommit -ine $TargetCommit) {
    throw "Target tooling checkout does not match commit $TargetCommit."
  }
  & (Join-Path $targetRoot 'eng/scripts/Build-APIReviewTools.ps1') `
    -LanguageToolsSourceDirectory $targetRoot -BinOutputDirectory $ToolsDirectory
}
finally {
  if ($temporaryTarget -and (Test-Path -LiteralPath $temporaryTarget)) {
    Remove-Item -LiteralPath $temporaryTarget -Recurse -Force
  }
}
