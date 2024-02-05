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
