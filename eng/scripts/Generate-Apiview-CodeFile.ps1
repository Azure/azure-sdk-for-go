param(
    $ServiceDirectory,
    $CodeFileOutDirectory,
    $ParserPath
)

Set-StrictMode -Version 3
# download APIView go parser from GitHub release
$modules = Get-ChildItem "go.mod" -Path $ServiceDirectory -Recurse
Write-Host $modules
if ($modules)
{
    foreach ($module in $modules)
    {
        Write-Host "Processing module path $($module)"
        $pkgRoot = Split-Path -Path $module
        Write-Host "Generating API review file for package $($pkgRoot), review file: $($CodeFileOutDirectory)"
        &$ParserPath $pkgRoot $CodeFileOutDirectory
    }
}

