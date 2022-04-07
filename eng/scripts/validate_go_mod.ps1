Param(
    [string] $serviceDir
)

. (Join-Path $PSScriptRoot . Language-Settings.ps1)

Push-Location sdk/$serviceDir

$goModFiles = Get-ChildItem -Path . -Filter go.mod -Recurse

if ($goModFiles.Length -eq 0) {
    Write-Host "Could not find a go.mod file in the directory $(Get-Location)"
    exit 1
}

$hasError = $false
foreach ($goMod in $goModFiles) {
    Write-Host $goMod.Directory
    $packageProperties = Get-GoModuleProperties $goMod.Directory
    Write-Host "Props $packageProperties"
    if (-Not ($goMod.FullName -Like "*testdata*")) {
        $name = $goMod.FullName
        $patternMatches = Get-Content $name | Select-String -Pattern "replace "
        if ($patternMatches.Length -ne 0) {
            Write-Host "Found a replace directive in go.mod file at $name"
            $hasError = $true
        } else {
            Write-Host "Valid go.mod file at $name"
        }
    }
}

Pop-Location

if ($hasError) {
    exit 1
}
