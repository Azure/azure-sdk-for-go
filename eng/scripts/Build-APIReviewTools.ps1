# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

[CmdletBinding()]
param(
  [string] $ToolsSourceDirectory,

  [Parameter(Mandatory = $true)]
  [string] $LanguageToolsSourceDirectory,

  [Parameter(Mandatory = $true)]
  [string] $BinOutputDirectory
)

Set-StrictMode -Version 3
$ErrorActionPreference = 'Stop'

$configurationPath = Join-Path $LanguageToolsSourceDirectory 'eng/api-review-tools.json'
$configuration = Get-Content -LiteralPath $configurationPath -Raw | ConvertFrom-Json
if (-not $configuration.PSObject.Properties['commit'] -or $configuration.commit -notmatch '^[0-9a-f]{40}$') {
  throw "API review tooling configuration must contain an immutable 40-character commit: $configurationPath"
}
$toolsCommit = $configuration.commit
$toolsRepository = 'https://github.com/Azure/azure-sdk-tools.git'
if ($ToolsSourceDirectory) {
  & git -C $ToolsSourceDirectory cat-file -e "$toolsCommit^{commit}" 2>$null
  if ($LASTEXITCODE -eq 0) {
    $toolsRepository = [IO.Path]::GetFullPath($ToolsSourceDirectory)
  }
}

$BinOutputDirectory = [IO.Path]::GetFullPath($BinOutputDirectory)
New-Item -ItemType Directory -Path $BinOutputDirectory -Force | Out-Null
$pinnedSource = Join-Path $BinOutputDirectory ('.tools-source-' + [guid]::NewGuid().ToString('N'))
try {
  & git init --quiet $pinnedSource
  if ($LASTEXITCODE -ne 0) {
    throw 'Failed to initialize the pinned API review tooling checkout.'
  }
  & git -C $pinnedSource remote add origin $toolsRepository
  if ($LASTEXITCODE -ne 0) {
    throw 'Failed to configure the API review tooling repository.'
  }
  & git -C $pinnedSource sparse-checkout set --cone 'src/go' 'eng/common/scripts'
  if ($LASTEXITCODE -ne 0) {
    throw 'Failed to configure the API review tooling sparse checkout.'
  }
  & git -C $pinnedSource fetch --quiet --depth 1 --filter=blob:none origin $toolsCommit
  if ($LASTEXITCODE -ne 0) {
    throw "Failed to fetch API review tooling commit $toolsCommit."
  }
  & git -C $pinnedSource checkout --quiet --detach FETCH_HEAD
  if ($LASTEXITCODE -ne 0) {
    throw "Failed to check out API review tooling commit $toolsCommit."
  }
  $actualCommit = & git -C $pinnedSource rev-parse HEAD
  if ($LASTEXITCODE -ne 0 -or $actualCommit -ine $toolsCommit) {
    throw "API review tooling checkout does not match pinned commit $toolsCommit."
  }
  Write-Information "Building Go API review tooling from azure-sdk-tools $toolsCommit." -InformationAction Continue

  $parserName = if ($IsWindows) { 'apiviewgo.exe' } else { 'apiviewgo' }
  $parserPath = Join-Path $BinOutputDirectory $parserName
  Push-Location ([IO.Path]::Combine($pinnedSource, 'src', 'go'))
  try {
    & go build -mod=readonly -o $parserPath .
    if ($LASTEXITCODE -ne 0) {
      throw 'Failed to build the Go APIView parser.'
    }
  }
  finally {
    Pop-Location
  }

  & (Join-Path $LanguageToolsSourceDirectory 'eng/scripts/Test-GoAPIReviewParser.ps1') `
    -ParserPath $parserPath -WorkingDirectory $BinOutputDirectory

  $scripts = @(
    [IO.Path]::Combine($LanguageToolsSourceDirectory, 'eng', 'scripts', 'Export-APIReviewBundle.ps1'),
    [IO.Path]::Combine($LanguageToolsSourceDirectory, 'eng', 'scripts', 'Get-GoModuleVersionInfo.ps1'),
    [IO.Path]::Combine($pinnedSource, 'eng', 'common', 'scripts', 'Export-APIViewMarkdown.ps1'),
    [IO.Path]::Combine($pinnedSource, 'eng', 'common', 'scripts', 'SemVer.ps1'),
    [IO.Path]::Combine($pinnedSource, 'eng', 'common', 'scripts', 'logging.ps1')
  )
  Copy-Item -LiteralPath $scripts -Destination $BinOutputDirectory -Force
}
finally {
  if (Test-Path -LiteralPath $pinnedSource) {
    Remove-Item -LiteralPath $pinnedSource -Recurse -Force
  }
}
