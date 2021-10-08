#Requires -Version 7.0

Param(
    [string] $serviceDirectory
)

Write-Host "Finding test directories in 'sdk/$serviceDirectory'"
$testDirs = & $PSScriptRoot/get_test_dirs.ps1 -serviceDir $serviceDirectory
Write-Host $testDirs

foreach ($td in $testDirs) {
    Push-Location $td
    Write-Host "##[command]Executing 'go test -run "^Test" -v -coverprofile coverage.txt .' in $td"
    go test -run "^Test" -v -coverprofile coverage.txt . | Tee-Object -FilePath outfile.txt
    if ($LASTEXITCODE) {
        exit $LASTEXITCODE
    }
    Get-Content outfile.txt | go-junit-report > report.xml

    # if no tests were actually run (e.g. examples) delete the coverage file so it's omitted from the coverage report
    if (Select-String -path ./report.xml -pattern '<testsuites></testsuites>' -simplematch -quiet) {
        Write-Host "##[command]Deleting empty coverage file"
        Remove-Item coverage.txt
    }
}
