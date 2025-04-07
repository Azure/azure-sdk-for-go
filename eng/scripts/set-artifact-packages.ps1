param (
    [string] $Packages,
    [string] $PackageInfoFolder
)

$packageSet = "$Packages" -split ","

# retrieve the package info files
$packageProperties = Get-ChildItem -Recurse "$PackageInfoFolder" *.json `
| Foreach-Object { Get-Content -Raw -Path $_.FullName | ConvertFrom-Json }

# filter the package info files to only those that are part of the targeted batch (present in $Packages arg)
# so that we can accurate determine the affected services for the current batch
$changedPackages = $packageProperties | Where-Object { $packageSet -contains $_.ArtifactName }
$changedServicesArray = $changedPackages | ForEach-Object { $_.ServiceDirectory } | Get-Unique
$changedServices = $changedServicesArray -join ","

# remove any package.json files that are not part of the targeted batch
Get-ChildItem -Recurse "$PackageInfoFolder" *.json | ForEach-Object {
    $fileContent = Get-Content -Raw -Path $_.FullName | ConvertFrom-Json
    if ($packageSet -notcontains $fileContent.Name) {
        Remove-Item $_.FullName -Force
        Write-Host "Removed $($_.FullName) as it doesn't belong to the package set that this batch is targeting."
    }
}

$proxyPackages = @()
$nonProxyPackages = @()

# now populate ProxyPackages and NonProxyPackages variables with the package info files that are part of the targeted batch
$changedPackages | ForEach-Object {
    if ($_.CIParameters.UsePipelineProxy) {
        $proxyPackages += $_.ArtifactName
    } else {
        $nonProxyPackages += $_.ArtifactName
    }
}

if ($proxyPackages) {
    Write-Host "##vso[task.setvariable variable=ProxyPackagesPresent;]true"
}
else {
    Write-Host "##vso[task.setvariable variable=ProxyPackagesPresent;]false"
}

if ($nonProxyPackages) {
    Write-Host "##vso[task.setvariable variable=NonProxyPackagesPresent;]true"
}
else {
    Write-Host "##vso[task.setvariable variable=NonProxyPackagesPresent;]false"
}

Write-Host "##vso[task.setvariable variable=ProxyPackages;]$($proxyPackages -join ",")"
Write-Host "##vso[task.setvariable variable=NonProxyPackages;]$($nonProxyPackages -join ",")"
Write-Host "##vso[task.setvariable variable=ChangedServices;]$changedServices"
Write-Host "This run is targeting: $Packages in [$changedServices]"