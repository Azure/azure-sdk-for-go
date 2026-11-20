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

function Read-ParserProbe([string] $Executable, [string] $SourceDirectory, [string] $OutputDirectory) {
  New-Item -ItemType Directory -Path $OutputDirectory -Force | Out-Null
  & $Executable $SourceDirectory $OutputDirectory | Out-Host
  if ($LASTEXITCODE -ne 0) {
    throw "Go APIView parser compatibility check failed with exit code $LASTEXITCODE."
  }
  $tokenFiles = @(Get-ChildItem -LiteralPath $OutputDirectory -Filter '*.json' -File)
  if ($tokenFiles.Count -ne 1) {
    throw "Go APIView parser compatibility check produced $($tokenFiles.Count) token files; expected one."
  }
  $review = Get-Content -LiteralPath $tokenFiles[0].FullName -Raw | ConvertFrom-Json
  if (-not $review.PSObject.Properties['ReviewLines']) {
    throw 'Go APIView parser compatibility check returned no review lines.'
  }
  return $review
}

function Assert-ProbeDeclaration($Review, [string] $Name, [string] $Text) {
  $line = Find-ReviewLine $Review.ReviewLines "parserprobe.$Name"
  if ($null -eq $line -or (Get-ReviewLineText $line) -cne $Text) {
    throw "Go APIView parser compatibility check expected $Text for $Name."
  }
}

$probeDirectory = Join-Path $WorkingDirectory ('.parser-probe-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $probeDirectory | Out-Null
$originalModCache = $env:GOMODCACHE
try {
  $sourceDirectory = Join-Path $probeDirectory 'sdk/parserprobe'
  $typesDirectory = Join-Path $sourceDirectory 'internal/types'
  $outputDirectory = Join-Path $probeDirectory 'output'
  New-Item -ItemType Directory -Path $typesDirectory, $outputDirectory -Force | Out-Null
  $moduleFile = @'
module example.com/sdk/parserprobe

go 1.21
require example.com/sdk/dependency v1.0.0
'@
  [IO.File]::WriteAllText((Join-Path $sourceDirectory 'go.mod'), $moduleFile)
  $source = @'
package parserprobe
import (
  "example.com/sdk/parserprobe/internal/types"
  "example.com/sdk/dependency"
)
type ByteSlice = types.Data
type AnotherSlice = types.Data
type Mapping = types.StringMap
type Defined types.Data
func (Defined) Local() {}
type Direct = []byte
type External = dependency.Data
type Contract[T any] interface { Handle(T) }
func StartTestProxy() {}
func StopTestProxy() {}
func ExampleClient() {}
'@
  [IO.File]::WriteAllText((Join-Path $sourceDirectory 'alias.go'), $source)
  [IO.File]::WriteAllText((Join-Path $typesDirectory 'types.go'), @'
package types
type Data []byte
func (Data) Len() int { return 0 }
func (*Data) Reset() {}
type StringMap map[string]int
func (StringMap) Size() int { return 0 }
'@)
  $env:GOMODCACHE = Join-Path $probeDirectory 'modcache'
  foreach ($dependency in @(
    @{ Path = 'modcache/example.com/sdk/dependency@v1.0.0'; Type = '[]byte' },
    @{ Path = 'modcache/example.com/sdk/dependency@v1.1.0'; Type = '[]rune' },
    @{ Path = 'sdk/dependency'; Type = '[]bool' }
  )) {
    $directory = Join-Path $probeDirectory $dependency.Path
    New-Item -ItemType Directory -Path $directory -Force | Out-Null
    [IO.File]::WriteAllText((Join-Path $directory 'go.mod'), "module example.com/sdk/dependency`n`ngo 1.21`n")
    [IO.File]::WriteAllText((Join-Path $directory 'types.go'), "package dependency`ntype Data $($dependency.Type)`n")
  }

  $review = Read-ParserProbe $ParserPath $sourceDirectory $outputDirectory
  $expectedTypes = @(
    @{ Name = 'ByteSlice'; Definition = '=[]byte'; Methods = @('func(ByteSlice)Len()', 'func(*ByteSlice)Reset()'); Absent = @() },
    @{ Name = 'AnotherSlice'; Definition = '=[]byte'; Methods = @('func(AnotherSlice)Len()', 'func(*AnotherSlice)Reset()'); Absent = @() },
    @{ Name = 'Mapping'; Definition = '=map[string]int'; Methods = @('func(Mapping)Size()'); Absent = @() },
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
  Assert-ProbeDeclaration $review 'Direct' 'typeDirect=[]byte'
  Assert-ProbeDeclaration $review 'External' 'typeExternal=[]byte'
  Assert-ProbeDeclaration $review 'Contract' 'typeContract[Tany]interface'
  foreach ($name in @('StartTestProxy', 'StopTestProxy', 'ExampleClient')) {
    if ($null -eq (Find-ReviewLine $review.ReviewLines "parserprobe-$name")) {
      throw "Go APIView parser omitted the exported $name function."
    }
  }

  [IO.File]::WriteAllText((Join-Path $sourceDirectory 'go.mod'), $moduleFile.Replace('v1.0.0', 'v1.1.0'))
  [IO.File]::WriteAllText((Join-Path $sourceDirectory 'alias.go'),
    $source.Replace('T any', 'T comparable').Replace('type Direct = []byte', 'type Direct []byte'))
  $changed = Read-ParserProbe $ParserPath $sourceDirectory (Join-Path $probeDirectory 'changed')
  Assert-ProbeDeclaration $changed 'Direct' 'typeDirect[]byte'
  Assert-ProbeDeclaration $changed 'External' 'typeExternal=[]rune'
  Assert-ProbeDeclaration $changed 'Contract' 'typeContract[Tcomparable]interface'
}
finally {
  $env:GOMODCACHE = $originalModCache
  Remove-Item -LiteralPath $probeDirectory -Recurse -Force
}
