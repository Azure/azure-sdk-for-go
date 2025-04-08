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
    exit $GOTESTEXITCODE
}
exit 0