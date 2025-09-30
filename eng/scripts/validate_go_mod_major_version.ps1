Param(
    [string] $serviceDir
)

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

# fetch the module version from the constant
$modVersion, $null = Get-GoModuleVersionInfo sdk/$serviceDir

if (!$modVersion) {
    Write-Host "Could not find module version for module directory $(Get-Location)"
    exit 1
}

$goModFile = Get-ChildItem -Path sdk/$serviceDir -Filter go.mod

if ($goModFile.Length -eq 0) {
    Write-Host "Could not find go.mod file in the directory $(Get-Location)"
    exit 1
}

# fetch the major version suffix from go.mod if present
$majorVersion = Get-GoModuleMajorVersion $goModFile

if ($majorVersion) {
    # go.mod has a major version
    # ensure it matches the major version defined in the constant
    if (-not $modVersion.StartsWith("$majorVersion.")) {
        Write-Host "Mismatched module major versions. go.mod states $majorVersion while const version states $modVersion"
        exit 1
    }
} else {
    # go.mod has no major version
    # ensure the constant is v0 or v1
    if ($modVersion -notmatch "^[01]\.") {
        Write-Host "The module's identity is missing a major version suffix for const version $modVersion"
        exit 1
    }
}
