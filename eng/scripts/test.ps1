param (
    [string] $serviceDirectory,
    [string] $testTimeout,
    [bool] $enableRaceDetector
)
. $PSScriptRoot/../common/scripts/common.ps1
# we are passing in a single item here, so the first item is all we need as that is all there will ever be
$targetDirectory = (ResolveSearchPaths $serviceDirectory)[0]
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

    go tool cover -func coverage.txt | Out-File -FilePath coveragefunc.txt
    Move-Item -Force ./coveragefunc.txt $repoRoot

    # use internal tool to fail if coverage is too low
    Pop-Location

    go run $repoRoot/eng/tools/internal/coverage/coverage.go `
        -config $repoRoot/eng/config.json `
        -serviceDirectory $targetDirectory `
        -searchDirectory $repoRoot

    if ($LASTEXITCODE -gt 0) {
        exit $LASTEXITCODE
    }
}

if ($GOTESTEXITCODE) {
    exit $GOTESTEXITCODE
}
exit 0
