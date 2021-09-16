Param(
    [string] $serviceDir
)

$cwd = Get-Location

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

# # 0. Find all test directories
# Write-Host "Finding test directories in 'sdk/$serviceDir'"
# $testDirs = Get-AllPackageInfoFromRepo $serviceDir
# # Issues here, not returning any objects
# Write-Host "Found test directories $testDirs"

$runTests = $false
# 0b. Verify there are test files with tests to run in at least one of these directories
$sdks = Get-AllPackageInfoFromRepo $serviceDir
Write-Host ($sdks | Format-Table | Out-String)
foreach ($sdk in (Get-AllPackageInfoFromRepo $serviceDir))
{
    Write-Host ($sdk | Format-Table | Out-String)
    $testFiles = Get-ChildItem -Path $sdk.Name -Filter *_test.go
    foreach ($testFile in $testFiles) {
        if (Select-String -path $testFile -pattern 'Test' -SimpleMatch) {
            $runTests = $true
        }
    }
}

if (!$runTests) {
    Write-Host "There were no test files found."
    # Exit0
}

go get github.com/jstemmer/go-junit-report
Set-Location $cwd

Write-Host "Proceeding to run tests and add coverage"

# 1. Run tests
foreach ($sdk in (Get-AllPackageInfoFromRepo $serviceDir))
{
    Push-Location $sdk.Name
    Write-Host "##[command]Executing 'go test -run ""^Test"" -v -coverprofile coverage.txt .' in $sdk.Name"
    go test -run "^Test" -v -coverprofile coverage.txt . | Tee-Object -FilePath outfile.txt
    if ($LASTEXITCODE) {
        # Exit$LASTEXITCODE
    }
    Get-Content outfile.txt | go-junit-report > report.xml
    # if no tests were actually run (e.g. examples) delete the coverage file so it's omitted from the coverage report
    if (Select-String -path ./report.xml -pattern '<testsuites></testsuites>' -simplematch -quiet) {
        Write-Host "##[command] Deleting empty coverage file"
        Remove-Item coverage.txt
        # Exit0
    }
}

Pop-Location
