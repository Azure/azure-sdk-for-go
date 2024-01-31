param (
    [string]$readmePath,
    [string]$outputFile
)

if ($readmePath -notlike "*sdk/resourcemanager*" -or $readmePath -like "*sdk/resourcemanager/internal*") {
    return
}

$lines = Get-Content -Path $readmePath

$pattern = "(https://pkg\.go\.dev/github\.com/Azure/azure-sdk-for-go/sdk/resourcemanager/[a-zA-Z0-9_.\-/]+)"
$matchPkg = $lines | Where-Object { $_ -match $pattern }

foreach ($m in $matchPkg) {
    $link = [regex]::Match($m, $pattern).Groups[1].Value
    $link | Out-File -FilePath $outputFile -Append
}
