# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

[CmdletBinding()]
param(
  [Parameter(Mandatory = $true)]
  [string] $ParserPath,

  [Parameter(Mandatory = $true)]
  [string] $WorkingDirectory
)

Set-StrictMode -Version 3
$ErrorActionPreference = 'Stop'

function Find-ReviewLine([object[]] $Lines, [string] $LineId) {
  foreach ($line in $Lines) {
    if ($line.PSObject.Properties['LineId'] -and $line.LineId -ceq $LineId) {
      return $line
    }
    if ($line.PSObject.Properties['Children'] -and $line.Children) {
      $match = Find-ReviewLine $line.Children $LineId
      if ($null -ne $match) {
        return $match
      }
    }
  }
  return $null
}

function Get-ReviewLineText($Line) {
  return (@($Line.Tokens | ForEach-Object { $_.Value }) -join '')
}

$probeDirectory = Join-Path $WorkingDirectory ('.parser-probe-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $probeDirectory | Out-Null
try {
  $sourceDirectory = Join-Path $probeDirectory 'source'
  $typesDirectory = Join-Path $sourceDirectory 'internal/types'
  $outputDirectory = Join-Path $probeDirectory 'output'
  New-Item -ItemType Directory -Path $typesDirectory, $outputDirectory -Force | Out-Null
  [IO.File]::WriteAllText((Join-Path $sourceDirectory 'go.mod'), @'
module example.com/parserprobe

go 1.21
'@)
  [IO.File]::WriteAllText((Join-Path $sourceDirectory 'alias.go'), @'
package parserprobe
import "example.com/parserprobe/internal/types"
type ByteSlice = types.Data
type AnotherSlice = types.Data
type Mapping = types.StringMap
type Defined types.Data
func (Defined) Local() {}
'@)
  [IO.File]::WriteAllText((Join-Path $typesDirectory 'types.go'), @'
package types
type Data []byte
func (Data) Len() int { return 0 }
func (*Data) Reset() {}
type StringMap map[string]int
func (StringMap) Size() int { return 0 }
'@)

  & $ParserPath $sourceDirectory $outputDirectory | Out-Host
  if ($LASTEXITCODE -ne 0) {
    throw "Go APIView parser compatibility check failed with exit code $LASTEXITCODE."
  }
  $tokenFiles = @(Get-ChildItem -LiteralPath $outputDirectory -Filter '*.json' -File)
  if ($tokenFiles.Count -ne 1) {
    throw "Go APIView parser compatibility check produced $($tokenFiles.Count) token files; expected one."
  }
  $review = Get-Content -LiteralPath $tokenFiles[0].FullName -Raw | ConvertFrom-Json
  if (-not $review.PSObject.Properties['ReviewLines']) {
    throw 'Go APIView parser compatibility check returned no review lines.'
  }
  $expectedTypes = @(
    @{ Name = 'ByteSlice'; Definition = '[]byte'; Methods = @('func(ByteSlice)Len()', 'func(*ByteSlice)Reset()'); Absent = @() },
    @{ Name = 'AnotherSlice'; Definition = '[]byte'; Methods = @('func(AnotherSlice)Len()', 'func(*AnotherSlice)Reset()'); Absent = @() },
    @{ Name = 'Mapping'; Definition = 'map[string]int'; Methods = @('func(Mapping)Size()'); Absent = @() },
    @{ Name = 'Defined'; Definition = '[]byte'; Methods = @('func(Defined)Local()'); Absent = @('func(Defined)Len()', 'func(*Defined)Reset()') }
  )
  foreach ($expected in $expectedTypes) {
    $line = Find-ReviewLine $review.ReviewLines "parserprobe.$($expected.Name)"
    if ($null -eq $line -or (Get-ReviewLineText $line) -cne "type$($expected.Name)$($expected.Definition)") {
      throw "Go APIView parser does not resolve the $($expected.Name) alias. Use a parser that supports aliased composite types."
    }
    foreach ($method in $expected.Methods) {
      $methodLines = @()
      if ($line.PSObject.Properties['Children']) {
        $methodLines = @($line.Children | Where-Object { (Get-ReviewLineText $_).StartsWith($method, [StringComparison]::Ordinal) })
      }
      if ($methodLines.Count -ne 1) {
        throw "Go APIView parser does not resolve $method on an aliased type. Use a parser that supports aliased methods."
      }
    }
    foreach ($method in $expected.Absent) {
      if ($line.PSObject.Properties['Children'] -and
          @($line.Children | Where-Object { (Get-ReviewLineText $_).StartsWith($method, [StringComparison]::Ordinal) }).Count -gt 0) {
        throw "Go APIView parser incorrectly inherited $method on a defined type."
      }
    }
  }
}
finally {
  Remove-Item -LiteralPath $probeDirectory -Recurse -Force
}
