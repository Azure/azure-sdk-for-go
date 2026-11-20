#Requires -Version 7.0

# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

[CmdletBinding()]
param(
    [string] $SdkRoot = (Join-Path $PSScriptRoot '..\..\..\sdk'),
    [string] $ApigenPath = (Join-Path $PSScriptRoot 'apigen.exe'),
    [ValidateRange(1, [int]::MaxValue)]
    [int] $Concurrency = 1
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$PSNativeCommandUseErrorActionPreference = $false

$root = Get-Item -LiteralPath $SdkRoot
if (-not $root.PSIsContainer) {
    throw "SDK root is not a directory: $SdkRoot"
}
$generator = Get-Item -LiteralPath $ApigenPath
if ($generator.PSIsContainer) {
    throw "apigen executable is not a file: $ApigenPath"
}

$modules = @(Get-ChildItem -LiteralPath $root.FullName -Filter 'go.mod' -File -Recurse |
    Sort-Object -Property FullName)
if ($modules.Count -eq 0) {
    throw "No Go modules found under $($root.FullName)"
}

$failedModules = [System.Collections.Generic.List[string]]::new()
$generatorPath = $generator.FullName
$completed = 0
Write-Output "Generating API reports for $($modules.Count) modules under $($root.FullName) with up to $Concurrency concurrent tasks"

try {
    $modules | ForEach-Object -Parallel {
        Set-StrictMode -Version Latest
        $ErrorActionPreference = 'Stop'
        $PSNativeCommandUseErrorActionPreference = $false

        $moduleDirectory = $_.DirectoryName
        $diagnostics = @()
        $succeeded = $false
        try {
            $diagnostics = @(& $using:generatorPath $moduleDirectory -o $moduleDirectory 2>&1)
            $generationExitCode = $LASTEXITCODE
            $reportPath = Join-Path $moduleDirectory 'api.md'
            $metadataPath = Join-Path $moduleDirectory 'api.metadata.yml'
            $reportExists = Test-Path -LiteralPath $reportPath -PathType Leaf
            $metadataExists = Test-Path -LiteralPath $metadataPath -PathType Leaf
            $succeeded = $generationExitCode -eq 0 -and $reportExists -and $metadataExists
            if ($generationExitCode -eq 0 -and (-not $reportExists -or -not $metadataExists)) {
                $diagnostics += "apigen returned success but did not create both $reportPath and $metadataPath"
            }
        }
        catch {
            $succeeded = $false
            $diagnostics += $_
        }
        [pscustomobject]@{
            Directory = $moduleDirectory
            Succeeded = $succeeded
            Diagnostics = $diagnostics
        }
    } -ThrottleLimit $Concurrency | ForEach-Object {
        $completed++
        foreach ($diagnostic in $_.Diagnostics) {
            [Console]::Error.WriteLine("[$($_.Directory)] $diagnostic")
        }
        if (-not $_.Succeeded) {
            $failedModules.Add($_.Directory)
        }
        Write-Progress -Activity 'Generating API reports' -Status $_.Directory `
            -PercentComplete (($completed * 100) / $modules.Count)
    }
}
finally {
    Write-Progress -Activity 'Generating API reports' -Completed
}

if ($failedModules.Count -gt 0) {
    Write-Output "Failed to generate API review bundles for $($failedModules.Count) of $($modules.Count) modules:"
    $failedModules | Sort-Object | Write-Output
    exit 1
}

Write-Output "Successfully generated API review bundles for all $($modules.Count) modules."
exit 0
