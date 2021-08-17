#Requires -Version 7.0

Param(
    [string] $serviceDirectory
)

Write-Host $serviceDirectory
Push-Location $serviceDirectory

$coverageFiles = [Collections.Generic.List[String]]@()
Get-ChildItem -recurse -path . -filter coverage.txt | ForEach-Object {
  $covFile = $_.FullName
  Write-Host "Adding $covFile to the list of code coverage files"
  $coverageFiles.Add($covFile)
}

# merge coverage files
gocovmerge $coverageFiles > mergedCoverage.txt
gocov convert ./mergedCoverage.txt > ./coverage.json

# gocov converts rely on standard input
Get-Content ./coverage.json | gocov-xml > ./coverage.xml
Get-Content ./coverage.json | gocov-html > ./coverage.html

# use internal tool to fail if coverage is too low
Pop-Location
go run ../tools/internal/coverage/main.go  -serviceDirectory $serviceDirectory