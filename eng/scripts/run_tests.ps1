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
go get github.com/jstemmer/go-junit-report
go get github.com/axw/gocov/gocov
go get github.com/AlekSi/gocov-xml
go get github.com/matm/gocov-html
go get -u github.com/wadey/gocovmerge
Set-Location $cwd

Write-Host "Proceeding to run tests and add coverage"

# 1. Run tests
foreach ($td in $testDirs) {
    Push-Location $td
    $temp = Get-Location
    Write-Host "Currently in $temp"
    Write-Host "##[command]Executing 'go test -run ""^Test"" -v -coverprofile coverage.txt .' in $td"
    go test -run "^Test" -v -coverprofile coverage.txt . | Tee-Object -FilePath outfile.txt
    if ($LASTEXITCODE) {
        exit $LASTEXITCODE
    }
    Get-Content outfile.txt | go-junit-report > report.xml
    # if no tests were actually run (e.g. examples) delete the coverage file so it's omitted from the coverage report
    if (Select-String -path ./report.xml -pattern '<testsuites></testsuites>' -simplematch -quiet) {
        Write-Host "##[command] Deleting empty coverage file"
        Remove-Item coverage.txt
        Exit 0
    }
}

Set-Location $cwd

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
