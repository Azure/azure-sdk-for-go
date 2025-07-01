#Requires -Version 7.0

Param(
    [Parameter(Mandatory = $true)]
    [string] $TargetDirectories,
    [string] $testTimeout = "10s",
    [bool] $enableRaceDetector = $false
)

$ErrorActionPreference = 'Stop'

$directories = $TargetDirectories -split ","

foreach($serviceOrPackageDir in $directories) {
    &$PSScriptRoot/test.ps1 $serviceOrPackageDir $testTimeout $enableRaceDetector
    if ($LASTEXITCODE) {
        Write-Host "##[error] a failure occurred testing the directory: $serviceOrPackageDir. Check above details for more information."
        $failed = $true
    }
}

if ($failed) {
    Write-Host "##[error] a failure occurred testing the directories: $TargetDirectories. Check above details for more information."
    exit 1
}
