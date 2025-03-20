param (
    [string] $Packages,
    [string] $PackageInfoFolder
)

. $PSScriptRoot/generate-dependency-functions.ps1

$packageSet = "$Packages" -split ","

# retrieve the package info files
$packageProperties = Get-ChildItem -Recurse "$PackageInfoFolder" *.json `
| Foreach-Object { Get-Content -Raw -Path $_.FullName | ConvertFrom-Json }

# filter the package info files to only those that are part of the targeted batch (present in $Packages arg)
# so that we can accurate determine the affected services for the current batch
$changedServicesArray = $packageProperties | Where-Object { $packageSet -contains $_.ArtifactName } `
    | ForEach-Object { $_.ServiceDirectory } | Get-Unique
$changedServices = $changedServicesArray -join ","

# todo: include the removal of the package info files that are not part of the targeted batch

Write-Host "##vso[task.setvariable variable=ChangedServices;]$changedServices"
Write-Host "This run is targeting: $Packages in [$changedServices]"