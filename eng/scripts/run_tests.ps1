Param(
    [string] $serviceDir
)

$cwd = Get-Location

# 0. Find all test directories
Write-Host "Finding test directories in 'sdk/$serviceDir'"
$testDirs = & $PSScriptRoot/get_test_dirs.ps1 -serviceDir sdk/$serviceDir
# Issues here, not returning any objects
Write-Host "Found test directories $testDirs"
$temp = Get-Location
Write-Host "Currently in $temp"

$runTests = $false
# 0b. Verify there are test files with tests to run in at least one of these directories
$testDirs | ForEach-Object {
    Write-Host "[Loop]: Currently in: $_"
    Get-ChildItem -Path $_ -Filter *_test.go | ForEach-Object {
        Write-Host "[Loop]: $_"
        if (Select-String -path $_ -pattern 'Test' -SimpleMatch) {
            $runTests = $true
        }
    }

    if (!$runTests) {
        Write-Host "There were no test files found."
        Exit 0
    }
}

Write-Host "Downloading coverage tools"
$cwd = Get-Location
Set-Location ./tools/internal/coverage
go get -v
Set-Location $cwd

Write-Host "Proceeding to run tests and add coverage"

# 1. Run tests
foreach ($td in $testDirs) {
    Push-Location $td
    $temp = Get-Location
    Write-Host "Currently in $temp"
    Write-Host "##[command] Executing go test -run ""^Test"" -v -coverprofile coverage.txt $td | go-junit-report -set-exit-code > report.xml"
    go test -run "^Test" -v -coverprofile coverage.txt $td | go-junit-report -set-exit-code > report.xml
    if (!$?) {
        Write-Host "There was an error running the tests"
        Exit $LASTEXITCODE
    } else {
        Write-Host "" "Successfully ran test suite at $temp" ""
    }
    # if no tests were actually run (e.g. examples) delete the coverage file so it's omitted from the coverage report
    if (Select-String -path ./report.xml -pattern '<testsuites></testsuites>' -simplematch -quiet) {
        Write-Host "##[command] Deleting empty coverage file"
        Remove-Item coverage.txt
        Exit 0
    }
}

Set-Location $cwd
$temp = Get-Location
Write-Host "Currently in $temp"

$coverageFiles = [Collections.Generic.List[String]]@()
Get-ChildItem -recurse -path . -filter coverage.txt | ForEach-Object {
    $covFile = $_.FullName
    Write-Host "Adding $covFile to the list of code coverage files"
    $coverageFiles.Add($covFile)
}

# merge coverage files
gocovmerge $coverageFiles > mergedCoverage.txt
gocov convert ./mergedCoverage.txt > ./coverage.json
if (!$?) {
    Write-Host "Issues converting the mergedCoverage.txt to a json file"
    Exit $LASTEXITCODE
}

# gocov converts rely on standard input
Get-Content ./coverage.json | gocov-xml > ./coverage.xml
Get-Content ./coverage.json | gocov-html > ./coverage.html

# use internal tool to fail if coverage is too low
Write-Host "##[command] go run ./tools/internal/coverage/main.go -serviceDirectory $serviceDir"
go run ./tools/internal/coverage/main.go -serviceDirectory $serviceDir
