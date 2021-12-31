#Requires -Version 7.0
param(
    [string]$inputJsonFile,
    [string]$outputJsonFile
)

. (Join-Path $PSScriptRoot .. MgmtTestLib.ps1)
Write-output "inputfile:$inputJsonFile, outfile:$outputJsonFile"
$inputJson = Get-Content $inputJsonFile | Out-String | ConvertFrom-Json
$inputJson
$packageFolder = $inputJson.packageFolder
$packageFolder = $packageFolder -replace "\\", "/"

$runLocalMockServer = $false
if ([string]::IsNullOrEmpty($inputJson.mockServerHost)) {
    $runLocalMockServer = $true
}
Write-Host "##[command]Generate example and Mock Test " $packageFolder
Set-Location $packageFolder
Invoke-MgmtTestgen -sdkDirectory $packageFolder -autorestPath $packageFolder/autorest.md -generateExample -generateMockTest

if ($runLocalMockServer -eq $true) {
    Write-Host "Prepare Mock Server"
    PrepareMockServer
    Write-Host "Try Stop mock server"
    StopMockServer
}

Set-Location $packageFolder
Write-output "Run Mock Test"
$sdk = Get-GoModuleProperties $packageFolder
ExecuteSingleTest $sdk $runLocalMockServer

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
