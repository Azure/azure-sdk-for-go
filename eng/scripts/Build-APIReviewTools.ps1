# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)]
  [string] $ToolsSourceDirectory,

  [Parameter(Mandatory = $true)]
  [string] $LanguageToolsSourceDirectory,

  [Parameter(Mandatory = $true)]
  [string] $BinOutputDirectory
)

$ErrorActionPreference = 'Stop'

$parserDirectory = [IO.Path]::Combine($ToolsSourceDirectory, 'src', 'go')
if (-not (Test-Path -LiteralPath (Join-Path $parserDirectory 'go.mod') -PathType Leaf)) {
  & git -C $ToolsSourceDirectory sparse-checkout add ([IO.Path]::Combine('src', 'go'))
  if ($LASTEXITCODE -ne 0) {
    throw "Failed to include Go APIView parser sources from $ToolsSourceDirectory."
  }
}

if (-not (Test-Path -LiteralPath (Join-Path $parserDirectory 'go.mod') -PathType Leaf)) {
  throw "Go APIView parser sources are missing from $parserDirectory."
}

New-Item -ItemType Directory -Path $BinOutputDirectory -Force | Out-Null
$parserName = if ($IsWindows) { 'apiviewgo.exe' } else { 'apiviewgo' }
$parserPath = Join-Path $BinOutputDirectory $parserName
Push-Location $parserDirectory
try {
  & go build -mod=readonly -o $parserPath .
  if ($LASTEXITCODE -ne 0) {
    throw "Failed to build the Go APIView parser."
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
  [IO.Path]::Combine($LanguageToolsSourceDirectory, 'eng', 'common', 'scripts', 'Export-APIViewMarkdown.ps1'),
  [IO.Path]::Combine($LanguageToolsSourceDirectory, 'eng', 'common', 'scripts', 'SemVer.ps1'),
  [IO.Path]::Combine($LanguageToolsSourceDirectory, 'eng', 'common', 'scripts', 'logging.ps1')
)
Copy-Item -LiteralPath $scripts -Destination $BinOutputDirectory -Force
