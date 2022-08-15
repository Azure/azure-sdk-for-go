param(
    $ServiceDirectory,
    $CodeFileOutDirectory,
    $ParserPath
)

Set-StrictMode -Version 3

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

$sdks = Get-AllPackageInfoFromRepo $ServiceDirectory
if ($sdks)
{
    foreach ($sdk in $sdks)
    {
        $pkgRoot = $sdk.DirectoryPath
        $moduleName = $sdk.ModuleName

        $stagingPath = Join-Path -Path $CodeFileOutDirectory $moduleName
        New-Item $stagingPath -Type Directory

        Compress-Archive -Path $pkgRoot -DestinationPath $stagingPath
        Write-Host "Generating API review file for package $($pkgRoot), review file Path: $($stagingPath)"
        &$ParserPath $pkgRoot $stagingPath
    }
}

