#Requires -Version 7.0
param(
    # The name of the RP which is the directory name in azure-rest-api-specs/specification
    [Parameter(Mandatory=$true)]
    [string]$rpName,
    # The name of the package
    [Parameter(Mandatory=$true)]
    [string]$packageName,
    # The friendly package title of this package which will be the title of README
    [Parameter(Mandatory=$true)]
    [string]$packageTitle, 
    # The commit hash of the azure-rest-api-specs on its main branch to generate the SDK
    [Parameter(Mandatory=$true)]
    [string]$commitHash
)

function ApplyTemplate([System.IO.FileSystemInfo]$destination) {
    Write-Host "##[command]RP name: " $rpName
    Write-Host "##[command]Package name: " $packageName
    Write-Host "##[command]Package Title: " $packageTitle
    Write-Host "##[command]Commit Hash: " $commitHash
    Write-Host "##[command]Des Directory: " $destination
    
}

$startignDirectory = Get-Location
$root = Resolve-Path ($PSScriptRoot + "/../..")
Set-Location $root
$targetDirectory = Join-Path $root "sdk" $rpName $packageName
$templateDirectory = Join-Path $root "eng/template"
try {
    New-Item -Path $targetDirectory -ItemType "directory" -Force
    Get-ChildItem $templateDirectory | ForEach-Object {
        ApplyTemplate $_
    }
}
finally {
    Set-Location $startignDirectory
}
