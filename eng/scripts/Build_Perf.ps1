#Requires -Version 7.0

Param(
    [string] $ServiceDirectories,
    [bool] $useAzcoreFromMain
)

$services = $ServiceDirectories -split ","

$failed = $false

foreach($serviceDirectory in $services) {
    $result = &$PSScriptRoot/perf.ps1 $serviceDirectory $useAzcoreFromMain

    if (!$result) {
        $failed = $true
    }
    if ($LASTEXITCODE) {
        Write-Host "##[error] a failure occurred vetting/building one or more performance tests in $serviceDirectory"
        $failed = $true
    }
}

if ($failed) {
    exit 1
}
