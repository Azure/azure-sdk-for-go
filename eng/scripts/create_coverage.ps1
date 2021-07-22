#Requires -Version 7.0

$coverageFiles = [Collections.Generic.List[String]]@()
Get-Childitem -recurse -path . -filter coverage.txt | foreach-object {
  $covFile = $_.FullName
  Write-Host "Adding $covFile to the list of code coverage files"
  $coverageFiles.Add($covFile)
}
gocovmerge $coverageFiles > mergedCoverage.txt
gocov convert ./mergedCoverage.txt > ./coverage.json

# gocov converts rely on standard input
Get-Content ./coverage.json | gocov-xml > ./coverage.xml
Get-Content ./coverage.json | gocov-html > ./coverage.html

# use internal tool to fail if coverage is too low
go run ../tools/internal/coverage/main.go