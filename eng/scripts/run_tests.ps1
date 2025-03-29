#Requires -Version 7.0

Param(
    [Parameter(Mandatory = $true)]
    [string] $TargetDirectories,
    [string] $testTimeout = "10s",
    [bool] $enableRaceDetector = $false
)

$ErrorActionPreference = 'Stop'

$services = $TargetDirectories -split ","

foreach($serviceDirectory in $services) {
    &$PSScriptRoot/test.ps1 $serviceDirectory $testTimeout $enableRaceDetector
    if ($LASTEXITCODE) {
        Write-Host "##[error] a failure occurred testing the directory: $serviceDirectory. Check above details for more information."
        $failed = $true
    }
}

if ($failed) {
    Write-Host "##[error] a failure occurred testing the directories: $TargetDirectories. Check above details for more information."
    exit 1
}
