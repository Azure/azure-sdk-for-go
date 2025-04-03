param(
    [Parameter(Mandatory = $true)]
    [string]$PackageInfoFolder
)


. (Join-Path $PSScriptRoot ".." "common" "scripts" "common.ps1")

$licenseCheckEnabledPackages = Get-PackagesFromPackageInfo -PackageInfoFolder $PackageInfoFolder `
  -IncludeIndirect $false -CustomCompareFunction { param($pkgProp) { return $pkgProp.CIParameters.LicenseCheck } }

foreach ($targetPackage in $licenseCheckEnabledPackages) {
    $target = $targetPackage.Name
    Push-Location $target
    Write-Host "Ensuring $target/LICENSE.txt file exists"
    if (Test-Path LICENSE.txt) {
        $patternMatches = Get-Content ./LICENSE.txt | Select-String -Pattern 'Copyright (\d{4}|\(c\)) Microsoft'
        if ($patternMatches.Length -eq 0) {
            Write-Host "LICENSE.txt file is invalid"
            $failed = $true
        }
    } else {
        Write-Host "Could not find a LICENSE.txt file"
        $failed = $true
    }
    Pop-Location
}

if ($failed) {
    Write-Host "License check failed, check output above."
    exit 1
}