#Requires -Version 7.0

Param(
    [string] $TargetDirectories,
    [bool] $useAzcoreFromMain
)

# get access to the common functions so we can use the common implementation of
# "give me a resolved invokable path for this filter string". Will properly handle
# one or multiple service directories in either fully qualified or pure servicedirectory name format
. (Join-Path $PSScriptRoot ".." "common" "scripts" "common.ps1" )
$services = ResolveSearchPaths -FilterString $TargetDirectories

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
