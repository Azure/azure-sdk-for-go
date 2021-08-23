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

Pop-Location
# go run ../tools/internal/coverage/main.go  -serviceDirectory $serviceDirectory

$patternMatches = Get-Content ./coverage.xml | Select-String -Pattern '<coverage line-rate=\"(\d\.\d+)\"'

if ($patternMatches.Length -eq 0) {
  Write-Host "Coverage.xml file did not contain coverage information"
  Exit $1
}

[double] $coverageFloat = $patternMatches.Matches.Groups[1].Captures

Write-Host $coverageFloat

# Read eng/config.json to find appropriate Value

$coverageGoals = Get-Content ./eng/config.json | Out-String | ConvertFrom-Json

Write-Host $coverageGoals

Write-Host $serviceDirectory
Foreach ($pkg in $coverageGoals.Packages) {
  Write-Host $pkg
  if ($pkg.Name -Match $serviceDirectory) {
    $goalCoverage = [double] $pkg.CoverageGoal

    if ($goalCoverage -le $coverageFloat) {
      Write-Host "Coverage is lower than the coverage goal"
      Exit 1
    } else {
      Exit 0
    }
  }
}

Write-Host "Could not find coverage goal, specify the coverage goal for your package in eng/config.json"
Exit $1