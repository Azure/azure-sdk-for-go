#Requires -Version 7.0

Param(
    [string] $filter
)

$packages = Get-AllPackageInfoFromRepo $filter
$targetServices = $packages | Where-Object { $_.CIParameters.NonShipping -eq $false } | Select-Object -Property Service -Unique

$failed = $false

foreach($service in $targetServices) {
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
