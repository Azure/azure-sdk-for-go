#Requires -Version 7.0

Param(
    [Parameter(Mandatory = $true)]
    [string] $serviceDirectories,
    [string] $testTimeout = "10s",
    [bool] $enableRaceDetector = $false
)

$ErrorActionPreference = 'Stop'

function Invoke-Test {
    param (
        [string] $serviceDirectory,
        [string] $testTimeout,
        [bool] $enableRaceDetector
    )
    $targetDirectory = $serviceDirectory
    if (-not $serviceDirectory.StartsWith("sdk")) {
        $targetDirectory = Join-Path "sdk" $serviceDirectory
    }
    Push-Location $targetDirectory

    if ($enableRaceDetector) {
        Write-Host "##[command] Executing 'go test -timeout $testTimeout -race -v -coverprofile coverage.txt ./...' in $targetDirectory"
        go test -timeout $testTimeout -race -v -coverprofile coverage.txt ./... | Tee-Object -FilePath outfile.txt
    } else {
        Write-Host "##[command] Executing 'go test -timeout $testTimeout -v -coverprofile coverage.txt ./...' in $targetDirectory"
        go test -timeout $testTimeout -v -coverprofile coverage.txt ./... | Tee-Object -FilePath outfile.txt
    }

    # go test will return a non-zero exit code on test failures so don't skip generating the report in this case
    $GOTESTEXITCODE = $LASTEXITCODE
    Write-Host "Finished test with exit code $GOTESTEXITCODE"

    Get-Content -Raw outfile.txt | go-junit-report > report.xml

    # if no tests were actually run (e.g. examples) delete the coverage file so it's omitted from the coverage report
    if (Select-String -path ./report.xml -pattern '<testsuites></testsuites>' -simplematch -quiet) {
        Write-Host "##[command]Deleting empty coverage file"
        Remove-Item coverage.txt
        Remove-Item outfile.txt
        Remove-Item report.xml

        Pop-Location
    } else {
        # Tests were actually run create a coverage report
        $repoRoot = Resolve-Path "$PSScriptRoot/../../"

        gocov convert ./coverage.txt > ./coverage.json

        # gocov converts rely on standard input
        Get-Content ./coverage.json | gocov-xml > ./coverage.xml
        Get-Content ./coverage.json | gocov-html > ./coverage.html

        Move-Item -Force ./coverage.xml $repoRoot
        Move-Item -Force ./coverage.html $repoRoot

        # use internal tool to fail if coverage is too low
        Pop-Location

        go run $repoRoot/eng/tools/internal/coverage/coverage.go `
            -config $repoRoot/eng/config.json `
            -serviceDirectory $targetDirectory `
            -searchDirectory $repoRoot
    }

    if ($GOTESTEXITCODE) {
        return $GOTESTEXITCODE
    }
    return 0
}


$services = $serviceDirectories -split ","

$failed = $false

foreach($serviceDirectory in $services) {
    $result = Invoke-Test $serviceDirectory $testTimeout $enableRaceDetector

    if ($result -gt 0) {
        Write-Host "An error occured while testing $serviceDirectory."
        $failed = $true
    }
}

if ($failed) {
    Write-Host "##[error] a failure occurred testing the directories: $serviceDirectories. Check above details for more information."
    exit 1
}
