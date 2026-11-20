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
  [string] $ParserPath
)

Set-StrictMode -Version 3
$ErrorActionPreference = 'Stop'

. (Join-Path $PSScriptRoot 'SemVer.ps1')
. (Join-Path $PSScriptRoot 'logging.ps1')
. (Join-Path $PSScriptRoot 'Get-GoModuleVersionInfo.ps1')

function Initialize-GoTokenSpacing([object[]] $Lines) {
  foreach ($line in $Lines) {
    foreach ($token in @($line.Tokens)) {
      if (-not $token.PSObject.Properties['HasPrefixSpace']) {
        $token | Add-Member -NotePropertyName HasPrefixSpace -NotePropertyValue $false
      }
      if (-not $token.PSObject.Properties['HasSuffixSpace']) {
        $token | Add-Member -NotePropertyName HasSuffixSpace -NotePropertyValue $true
      }
    }
    if ($line.PSObject.Properties['Children'] -and $line.Children) {
      Initialize-GoTokenSpacing $line.Children
    }
  }
}

$packageRelativePath = $PackageName.Replace('\', '/')
$pathParts = $packageRelativePath.Split('/')
if ($packageRelativePath -cnotmatch '^(?:sdk|profile)(?:/[A-Za-z0-9._-]+)+$' -or
    @($pathParts | Where-Object { $_ -in @('.', '..') }).Count -gt 0) {
  throw "PackageName must be a repository-relative SDK module path: $PackageName"
}

$sourceRootPath = [IO.Path]::GetFullPath($SourceRoot)
if (-not (Test-Path -LiteralPath $sourceRootPath -PathType Container)) {
  throw "Source root does not exist: $sourceRootPath"
}
$packagePath = $sourceRootPath
foreach ($part in $pathParts) {
  $packagePath = Join-Path $packagePath $part
}
if (-not (Test-Path -LiteralPath (Join-Path $packagePath 'go.mod') -PathType Leaf)) {
  throw "Go module not found for $PackageName at $packagePath."
}

$packageVersion, $versionFile = Get-GoModuleVersionInfo $packagePath
if ([string]::IsNullOrWhiteSpace($packageVersion)) {
  throw "Could not determine the version of $PackageName."
}

New-Item -ItemType Directory -Path $OutputDirectory -Force | Out-Null
$tokenDirectory = Join-Path $OutputDirectory ('.apiview-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $tokenDirectory | Out-Null
try {
  $parserOutput = @(& $ParserPath $packagePath $tokenDirectory)
  $parserExitCode = $LASTEXITCODE
  $parserOutput | Out-Host
  if ($parserExitCode -ne 0) {
    throw "Go APIView parser failed for $PackageName with exit code $parserExitCode."
  }
  $unsupported = @($parserOutput | Where-Object { $_ -match '^(?:(?:unexpected|unhandled) node type\b|unhandled declaration\b)' })
  if ($unsupported.Count -gt 0) {
    throw "Go APIView parser reported unsupported syntax for ${PackageName}: $($unsupported -join '; ')"
  }

  $tokenFiles = @(Get-ChildItem -LiteralPath $tokenDirectory -Filter '*.json' -File)
  if ($tokenFiles.Count -ne 1) {
    throw "Go APIView parser produced $($tokenFiles.Count) token files for $PackageName; expected one."
  }
  $review = Get-Content -LiteralPath $tokenFiles[0].FullName -Raw | ConvertFrom-Json
  if (-not $review.PSObject.Properties['PackageName'] -or $review.PackageName -cne $packageRelativePath) {
    throw "Go APIView parser returned a different package name for $PackageName."
  }
  if (-not $review.PSObject.Properties['Language'] -or $review.Language -cne 'Go') {
    throw "Go APIView parser returned an invalid language for $PackageName."
  }
  if ($review.PSObject.Properties['Diagnostics'] -and $review.Diagnostics) {
    $errors = @($review.Diagnostics | Where-Object { $_.Level -in @(3, 4, 'Error', 'Fatal') })
    if ($errors.Count -gt 0) {
      $details = ($errors | ForEach-Object { "$($_.TargetID): $($_.Text)" }) -join '; '
      throw "Go APIView parser reported incomplete API for ${PackageName}: $details"
    }
  }

  if (-not $review.PSObject.Properties['ReviewLines'] -or -not $review.ReviewLines) {
    $review | Add-Member -NotePropertyName ReviewLines -NotePropertyValue @([PSCustomObject]@{ Tokens = @() }) -Force
  }
  Initialize-GoTokenSpacing $review.ReviewLines
  $review | ConvertTo-Json -Depth 100 | Set-Content -LiteralPath $tokenFiles[0].FullName -Encoding utf8

  $apiPath = Join-Path $OutputDirectory 'api.md'
  & (Join-Path $PSScriptRoot 'Export-APIViewMarkdown.ps1') -TokenJsonPath $tokenFiles[0].FullName -OutputPath $apiPath
  if (-not (Test-Path -LiteralPath $apiPath -PathType Leaf)) {
    throw "API review markdown was not generated at $apiPath."
  }

  $apiHash = (Get-FileHash -LiteralPath $apiPath -Algorithm SHA256).Hash.ToLowerInvariant()
  $metadata = "apiHash: $apiHash`npackageVersion: $packageVersion`npackageRelativePath: $packageRelativePath`n"
  [IO.File]::WriteAllText((Join-Path $OutputDirectory 'api.metadata.yml'), $metadata, [Text.UTF8Encoding]::new($false))

  [PSCustomObject]@{
    ApiHash = $apiHash
    PackageName = $packageRelativePath
    PackageVersion = $packageVersion
  }
}
finally {
  Remove-Item -LiteralPath $tokenDirectory -Recurse -Force
}
