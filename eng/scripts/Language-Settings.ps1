$Language = "go"
$LanguageDisplayName = "go"

function Get-AllPackageInfoFromRepo ($serviceDirectory)
{
    $allPackageProps = @()
    $searchPath = Join-Path $RepoRoot "sdk"
    if ($serviceDirectory)
    {
        $searchPath = Join-Path $searchPath $serviceDirectory 
    }

    $pkgFiles = Get-ChildItem -Path $searchPath -Include "go.mod" -Recurse

    foreach ($pkgFile in $pkgFiles)
    {
        $pkgPath = $pkgFile.DirectoryName
        $pkgName = $pkgFile.Directory.Name
        $pkgVersion = ""
        $serviceDirectory = $pkgName

        if ($pkgPath -match ".*\\sdk\\(?<serviceDir>.*)")
        {
            $serviceDirectory = $Matches["serviceDir"]
        }
        $versionFile = Join-Path $pkgPath "version.go"

        if (Test-Path $versionFile)
        {
            $versionFileContent = (Get-Content -Path $versionFile -Raw) -replace '"', '\"'
            Push-Location $PSScriptRoot
            $pkgVersion = go run get_pkg_version.go $versionFileContent
            Pop-Location
        }

        $pkgProp = [PackageProps]::new($pkgName, $pkgVersion.Trim('"'), $pkgPath, $serviceDirectory)
        $pkgProp.IsNewSdk = $true
        $pkgProp.SdkType = "client"
        $allPackageProps += $pkgProp
    }
    return $allPackageProps
}