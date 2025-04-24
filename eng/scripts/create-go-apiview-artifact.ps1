param (
    [Parameter(Mandatory = $false)]
    [string] $ArtifactStagingDirectory = "C:/repo/azure-sdk-for-go/artifacts",
    [Parameter(Mandatory = $false)]
    [string] $PackageInfoFolder = "C:/repo/azure-sdk-for-go/PackageInfo",
    [Parameter(Mandatory = $false)]
    [string] $ArtifactName = "packages"
)

Write-Host "Creating APIView artifacts for Go SDKs in `"$ArtifactStagingDirectory`""
Write-Host "PackageInfoFolder: `"$PackageInfoFolder`""

. $PSScriptRoot/../common/scripts/common.ps1
. $PSScriptRoot/../scripts/apiview-helpers.ps1

$pkgs = Get-PackagesFromPackageInfo -PackageInfoFolder "$PackageInfoFolder" `
    -IncludeIndirect $false -CustomCompareFunction { param($pkgProp) { return $pkgProp.CIParameters.IsSdkLibrary } }

$apiviewEnabledPackages = @()
if ($pkgs) {
    $apiviewEnabledPackages = $pkgs | ForEach-Object { $_.Name } | Sort-Object -Unique
}
else {
    Write-Host "No sdk packages were found in package set. Skipping APIView generation."
    exit 0
}

$directoryToPublish = Join-Path -Path $ArtifactStagingDirectory $ArtifactName

foreach($packageDirectory in $apiviewEnabledPackages) {
    New-APIViewArtifacts `
        -ServiceDirectory $packageDirectory `
        -OutputDirectory $ArtifactStagingDirectory `
        -DirectoryToPublish $directoryToPublish
}

Copy-Item "$PackageInfoFolder" -Destination $directoryToPublish -Recurse