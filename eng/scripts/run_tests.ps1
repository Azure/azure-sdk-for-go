Param(
    [string] $serviceDir
)

$cwd = Get-Location

# 0. Find all test directories
Write-Host "Finding test directories in 'sdk/$serviceDir'"
$testDirs = & $PSScriptRoot/get_test_dirs.ps1 -serviceDir $serviceDir
# Issues here, not returning any objects
Write-Host "Found test directories $testDirs"

$runTests = $false
# 0b. Verify there are test files with tests to run in at least one of these directories
foreach ($testDir in $testDirs) {
    $testFiles = Get-ChildItem -Path $testDir -Filter *_test.go
    foreach ($testFile in $testFiles) {
        if (Select-String -path $testFile -pattern 'Test' -SimpleMatch) {
            $runTests = $true
        }
    }
}

if (!$runTests) {
    Write-Host "There were no test files found."
    Exit 0
}

go get github.com/jstemmer/go-junit-report
Set-Location $cwd

Write-Host "Proceeding to run tests and add coverage"

# 1. Run tests
foreach ($td in $testDirs) {
    Push-Location $td
    Write-Host "##[command]Executing 'go test -run ""^Test"" -v -coverprofile coverage.txt .' in $td"
    go test -run "^Test" -v -coverprofile coverage.txt . | Tee-Object -FilePath outfile.txt
    if ($LASTEXITCODE) {
        exit $LASTEXITCODE
    }
    Get-Content outfile.txt | go-junit-report > report.xml
    Get-Content report.xml
    Get-Location
    # if no tests were actually run (e.g. examples) delete the coverage file so it's omitted from the coverage report
    if (Select-String -path ./report.xml -pattern '<testsuites></testsuites>' -simplematch -quiet) {
        Write-Host "##[command] Deleting empty coverage file"
        Remove-Item coverage.txt
        Exit 0
    }
}

Pop-Location
