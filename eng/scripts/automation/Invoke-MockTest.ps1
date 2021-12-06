#Requires -Version 7.0
param(
    [string]$inputJsonFile = "input.json",
    [string]$outputJsonFile = "output.json"
)

. (Join-Path $PSScriptRoot .. MgmtTestLib.ps1)

$inputJson = Get-Content $inputJsonFile | Out-String | ConvertFrom-Json
$packageFolder = $inputJson.packageFolder
$packageFolder = $packageFolder -replace "\\", "/"

Write-Host "##[command]Generate example and Mock Test " $packageFolder
Set-Location $packageFolder
Invoke-MgmtTestgen -sdkDirectory $packageFolder -autorestPath $packageFolder/autorest.md -generateExample -generateMockTest

Write-output "Run Mock Test"
$sdk = Get-GoModuleProperties $packageFolder
ExecuteSingleTest $sdk $false

TestAndGenerateReport $packageFolder
$testoutputFile = Join-Path $packageFolder output.txt
$all = (Select-String -Path $testoutputFile -Pattern "=== RUN").Matches.length
$pass = (Select-String -Path $testoutputFile -Pattern "--- PASS").Matches.length
$fail = (Select-String -Path $testoutputFile -Pattern "--- FAIL").Matches.length
$coverage = (Select-String -Path $testoutputFile -Pattern "coverage: (?<coverage>.*)% of statements").Matches[0].Groups["coverage"].Value

$outputJson = [PSCustomObject]@{
    total = $all
    success = $pass
    fail = $fail
    apiCoverage = $coverage
}

$outputJson | ConvertTo-Json -depth 100 | Out-File $outputJsonFile
