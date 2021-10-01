Param(
  [string] $serviceDir
)

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

$modDirs = @()
foreach ($sdk in (Get-AllPackageInfoFromRepo $serviceDir))
{
    $modDirs += $sdk.DirectoryPath
}

return $modDirs
