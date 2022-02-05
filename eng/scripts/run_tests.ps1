#Requires -Version 7.0

Param(
    [string] $serviceDirectory,
    [string] $testTimeout
)

Push-Location sdk/$serviceDirectory
Write-Host "##[command] Executing 'go test -timeout $testTimeout -v -coverprofile coverage.txt ./...' in sdk/$serviceDirectory"

go test -timeout $testTimeout -v -coverprofile coverage.txt ./... | Tee-Object -FilePath outfile.txt
if ($LASTEXITCODE) {
    exit $LASTEXITCODE
}
exit 0

Get-Content outfile.txt | go-junit-report > report.xml

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

    Move-Item ./coverage.xml $repoRoot
    Move-Item ./coverage.html $repoRoot

    # use internal tool to fail if coverage is too low
    Pop-Location

    go run $repoRoot/eng/tools/internal/coverage/coverage.go `
        -config $repoRoot/eng/config.json `
        -serviceDirectory $serviceDirectory `
        -searchDirectory $repoRoot
}
