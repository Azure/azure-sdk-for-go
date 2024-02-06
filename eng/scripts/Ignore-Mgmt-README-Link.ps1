<#
.SYNOPSIS
  Add the Go management module link in readme.md to the 'ignore-links.txt'.

.DESCRIPTION
  This script is used to add the module link to 'ignore-links.txt' to skip it from the links verification. The reason to skip it from the validation is the module link is not available at this time during this pipeline run. It will be available after the release tag is created at Azure/azure-sdk-for-go repository.

.PARAMETER servicePath
  The path of the service.

.PARAMETER outputFile
  The path of the 'ignore-links.txt'.

.EXAMPLE
  Ignore-Mgmt-README-Link -servicePath "/home/azure-sdk-for-go/sdk/servicea" -outputFile "/home/azure-sdk-for-go/eng/ignore-links.txt"

#>

param (
    [ValidateNotNullOrEmpty()]
    [string]$serviePath,
    [ValidateNotNullOrEmpty()]
    [string]$outputFile
)

if ($serviePath -notlike "*sdk/resourcemanager*" -or $serviePath -like "*sdk/resourcemanager/internal*") {
    exit
}

if (-not (Test-Path $serviePath/README.md)) {
    Write-Host "$serviePath/README.md does not exist."
    exit
}

if (-not (Test-Path $outputFile)) {
    Write-Host "$outputFile does not exist."
    exit
}

try {
    $lines = Get-Content -Path $serviePath/README.md

    $pattern = "(https://pkg\.go\.dev/github\.com/Azure/azure-sdk-for-go/$serviePath[a-zA-Z0-9_.\-/]*)"
    $matchPkg = $lines | Where-Object { $_ -match $pattern }
    
    foreach ($m in $matchPkg) {
        $link = [regex]::Match($m, $pattern).Groups[1].Value
        Write-Host "Ignore Mgmt README.md Link: $link"
        $link | Out-File -FilePath $outputFile -Append
    }
}
catch {
    Write-Host "Ignore-Mgmt-README-Link failed with exception:`n$_"
    exit 1
}
