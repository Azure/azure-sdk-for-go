#Requires -Version 7.0

Param(
    [string] $serviceDirectory
)

Push-Location ./sdk/$serviceDirectory

$coverageFiles = [Collections.Generic.List[String]]@()
Get-ChildItem -recurse -path . -filter coverage.txt | ForEach-Object {
  $covFile = $_.FullName
  Write-Host "Adding $covFile to the list of code coverage files"
  $coverageFiles.Add($covFile)
}

Pop-Location

# merge coverage files
gocovmerge $coverageFiles > mergedCoverage.txt
gocov convert ./mergedCoverage.txt > ./coverage.json

# gocov converts rely on standard input
Get-Content ./coverage.json | gocov-xml > ./coverage.xml
Get-Content ./coverage.json | gocov-html > ./coverage.html

$patternMatches = Get-Content ./coverage.xml | Select-String -Pattern '<coverage line-rate=\"(\d\.\d+)\"'
if ($patternMatches.Length -eq 0) {
  Write-Host "Coverage.xml file did not contain coverage information"
  Exit $1
}

[double] $coverageFloat = $patternMatches.Matches.Groups[1].Value

# Read eng/config.json to find appropriate Value
$coverageGoals = Get-Content ./eng/config.json | Out-String | ConvertFrom-Json

Foreach ($pkg in $coverageGoals.Packages) {
  if ($pkg.Name -Match $serviceDirectory) {
    $goalCoverage = [double] $pkg.CoverageGoal

    if ($coverageFloat -le $goalCoverage) {
      Write-Host "Coverage ($coverageFloat) is lower than the coverage goal ($goalCoverage)"
      Exit 1
    } else {
      Exit 0
    }
  }
}

Write-Host "Could not find coverage goal, specify the coverage goal for your package in eng/config.json"
Exit $1