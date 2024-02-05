param (
    [string]$serviePath,
    [string]$outputFile
)

if ($serviePath -notlike "*sdk/resourcemanager*" -or $serviePath -like "*sdk/resourcemanager/internal*") {
    return
}

$lines = Get-Content -Path $serviePath/README.md

$pattern = "(https://pkg\.go\.dev/github\.com/Azure/azure-sdk-for-go/$serviePath[a-zA-Z0-9_.\-/]*)"
$matchPkg = $lines | Where-Object { $_ -match $pattern }

foreach ($m in $matchPkg) {
    $link = [regex]::Match($m, $pattern).Groups[1].Value
    Write-Host "Ignore Mgmt README.md Link: $link"
    $link | Out-File -FilePath $outputFile -Append
}
