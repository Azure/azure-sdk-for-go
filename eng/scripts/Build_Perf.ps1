#Requires -Version 7.0

Param(
    [string] $serviceDirectory,
    [bool] $useAzcoreFromMain
)

Push-Location sdk/$serviceDirectory

# Find all 'testdata' directories
$perfDirectories = Get-ChildItem -Path . -Filter testdata -Recurse

if ($perfDirectories.Length -eq 0) {
    Write-Host "##[command] Did not find any performance tests in the directory $(Get-Location)"
    exit 0
}

$failed = $false

foreach ($perfDir in $perfDirectories) {
    Push-Location $perfDir

    if (Test-Path -Path perf) {
        Push-Location perf
        Write-Host "##[command] Building and vetting performance tests in $perfDir/perf"

        if ($useAzcoreFromMain) {
            # using a live azcore might be dragging in updated dependencies
            Write-Host "##[command] Executing 'go mod tidy' in $perfDir/perf"
            go mod tidy
            if ($LASTEXITCODE) {
                $failed = $true
            }
        }

        Write-Host "##[command] Executing 'go build .' in $perfDir/perf"
        go build .
        if ($LASTEXITCODE) {
            $failed = $true
        }

        Write-Host "##[command] Executing 'go vet .' in $perfDir/perf"
        go vet .
        if ($LASTEXITCODE) {
            $failed = $true
        }
        Pop-Location
    }

    Pop-Location
}

Pop-Location

if ($failed) {
    Write-Host "##[command] a failure occurred vetting/building one or more performance tests"
    exit 1
}