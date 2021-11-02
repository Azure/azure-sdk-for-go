Param(
    [string] $serviceDir
)

Push-Location sdk/$serviceDir

$goModFiles = Get-ChildItem -Path . -Filter go.mod -Recurse

if ($goModFiles.Length -eq 0) {
    Write-Host "Could not find a go.mod file in the directory $(pwd)"
    exit 1
}

$hasError = $false
foreach ($goMod in $goModFiles) {
    $patternMatches = Get-Content $goMod.FullName | Select-String -Pattern "replace "
    if ($patternMatches.Length -ne 0) {
        $name = $goMod.FullName
        Write-Host "Found a replace directive in go.mod file at $name"
        $hasError = $true
    }
}

Pop-Location

if ($hasError) {
    exit 1
}
