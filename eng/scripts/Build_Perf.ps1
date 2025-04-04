#Requires -Version 7.0

Param(
    [string] $TargetDirectories,
    [bool] $useAzcoreFromMain
)

. (Join-Path $PSScriptRoot "Language-Settings.ps1")

$services = ResolveSearchPaths $TargetDirectories

$failed = $false

foreach($serviceDirectory in $services) {
    &$PSScriptRoot/perf.ps1 $serviceDirectory $useAzcoreFromMain

    if ($LASTEXITCODE) {
        Write-Host "##[error] a failure occurred vetting/building one or more performance tests in $serviceDirectory"
        $failed = $true
    }
}

if ($failed) {
    exit 1
}
