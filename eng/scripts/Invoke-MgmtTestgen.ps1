#Requires -Version 7.0
param(
    [string]$filter,
    [switch]$clean,
    [switch]$vet,
    [switch]$generateExample,
    [switch]$generateMockTest,
    [switch]$skipBuild,
    [switch]$cleanGenerated,
    [switch]$format,
    [switch]$tidy,
    [string]$config = "autorest.md",
    [string]$autorestVersion = "3.7.3",
    [string]$goExtension = "@autorest/go@4.0.0-preview.36",
    [string]$testExtension = "@autorest/gotest@1.3.0",
    [string]$outputFolder
)

. (Join-Path $PSScriptRoot .. common scripts common.ps1)
. (Join-Path $PSScriptRoot MgmtTestLib.ps1)

try
{
    $startingDirectory = Get-Location

    $sdks = Get-AllPackageInfoFromRepo $filter

    foreach ($sdk in $sdks)
    {
        Push-Location $sdk.DirectoryPath
        Invoke-MgmtTestgen -sdkDirectory $sdk.DirectoryPath @psBoundParameters
        Pop-Location
    }
}
finally
{
    Set-Location $startingDirectory
}