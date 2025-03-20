param(
    [Parameter(Mandatory = $true)]
    [string]$PackageInfoFolder,
    [Parameter(Mandatory = $true)]
    [bool]$IncludeIndirect = $false
)


$packageProps = Get-ChildItem -Recurse -Filter "*.json" -Path $PackageInfoFolder `
    | ForEach-Object { Get-Content -Raw $_.FullName | ConvertFrom-Json }
if (-not $IncludeIndirect) {
    $packageProps = $packageProps | Where-Object { $_.IncludedForValidation -eq $false }
}
$shippingPackageProps = $packageProps | Where-Object { $_.CIParameters.LicenseCheck -eq $true }
$changedServicesArray = $shippingPackageProps | ForEach-Object { $_.ServiceDirectory } | Get-Unique

foreach ($changedService in $changedServicesArray) {
    Push-Location
    Write-Host "Ensuring $TargetDirectory/LICENSE.txt file exists"
    if (Test-Path LICENSE.txt) {
        $patternMatches = Get-Content ./LICENSE.txt | Select-String -Pattern 'Copyright (\d{4}|\(c\)) Microsoft'
        if ($patternMatches.Length -eq 0) {
            Write-Host "LICENSE.txt file is invalid"
            $failed = $true
            Pop-Location
        }
    } else {
        Write-Host "Could not find a LICENSE.txt file"
        Pop-Location
        $failed = $true
    }

    Pop-Location
}

if ($failed) {
    Write-Host "License check failed, check output above."
    exit 1
}