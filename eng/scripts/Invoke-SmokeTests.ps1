#Requires -Version 7.0

Param(
    [string] $filter
)

. (Join-Path $PSScriptRoot ".." common scripts common.ps1)

$packages = Get-AllPackageInfoFromRepo $filter
$targetServices = $packages | Where-Object { $_.CIParameters.NonShipping -eq $false } | Select-Object -Property ServiceDirectory -Unique

$failed = $false

foreach($service in $targetServices) {
    Write-Host $service
    & $PSScriptRoot/Smoke_Tests_Nightly.ps1 $service

    if ($LASTEXITCODE) {
        Write-Host "Smoke tests failed for $service."
        $failed = $true
    }
}

if ($failed) {
    Write-Host "Smoke tests failed for one or more services."
    exit 1
}
