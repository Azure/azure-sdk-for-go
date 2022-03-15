#Requires -Version 7.0

Param(
    [string] $serviceDirectory,
)

Push-Location sdk/$serviceDirectory

if (Test-Path -Path testdata/perf) {
    Write-Host "##[command] Building and vetting performance tests in sdk/$serviceDirectory/testdata/perf"

    Write-Host "##[command] Executing 'go build .' in sdk/$serviceDirectory/testdata/perf"
    go build .
    if ($LASTEXITCODE) {
        exit $LASTEXITCODE
    }

    Write-Host "##[command] Executing 'go vet .' in sdk/$serviceDirectory/testdata/perf"
    go vet .
    if ($LASTEXITCODE) {
        exit $LASTEXITCODE
    }
} else {
    Write-Host "##[command] Did not find performance tests in sdk/$serviceDirectory/testdata/perf"
}

Pop-Location
