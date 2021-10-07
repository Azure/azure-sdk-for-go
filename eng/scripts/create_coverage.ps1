#Requires -Version 7.0

Param(
    [string] $serviceDirectory
)

$repoRoot = Resolve-Path "$PSScriptRoot/../../"

Write-Host "repoRoot $repoRoot"

Write-Host $serviceDirectory
Push-Location sdk/$serviceDirectory

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

Move-Item ./coverage.xml $repoRoot
Move-Item ./coverage.html $repoRoot

# use internal tool to fail if coverage is too low
Pop-Location

go run $repoRoot/eng/tools/internal/coverage/coverage.go `
    -config $repoRoot/eng/config.json `
    -serviceDirectory $serviceDirectory `
    -searchDirectory $repoRoot
