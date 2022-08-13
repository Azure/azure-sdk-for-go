param(
    $ServiceDirectory,
    $CodeFileOutDirectory,
    $ParserPath
)

Set-StrictMode -Version 3

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

$sdks = Get-AllPackageInfoFromRepo $ServiceDirectory
Write-Host $sdks
if ($sdks)
{
    foreach ($sdk in $sdks)
    {
        $pkgRoot = $sdk.DirectoryPath
        Write-Host "Processing SDK $($pkgRoot)"
        $pkgRoot = Split-Path -Path $module
        Write-Host "Generating API review file for package $($pkgRoot), review file: $($CodeFileOutDirectory)"
        &$ParserPath $pkgRoot $CodeFileOutDirectory
    }
}

