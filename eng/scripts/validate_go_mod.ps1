Param(
    [string] $serviceDir
)

Push-Location sdk/$serviceDir

$goModFiles = Get-ChildItem -Path . -Filter go.mod -Recurse

if ($goModFiles.Length -eq 0) {
    Write-Host "Could not find a go.mod file in the directory"
    exit 1
}

foreach ($goMod in $goModFiles) {
    $patternMatches = Get-Content $goMod.FullName | Select-String -Pattern "replace "
    if ($patternMatches.Length -ne 0) {
        $name = $goMod.FullName
        Write-Host "Found a replace directive in go.mod file at $name"
        exit 1
    }
}

Pop-Location
