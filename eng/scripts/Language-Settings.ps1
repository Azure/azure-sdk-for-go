$Language = "go"
$packagePattern = "go.mod"
$LanguageDisplayName = "go"

# get version from specific files (*constants.go, *version.go)
function Get-GoModuleVersionInfo($modPath)
{
  $VERSION_LINE_REGEX = ".+\s*=\s*`".*v(?<version>$([AzureEngSemanticVersion]::SEMVER_REGEX))`""

  $versionFiles = Get-ChildItem -Recurse -Path $modPath -Filter *.go

  # for each version file, use regex to search go version num
  foreach ($versionFile in $versionFiles)
  {
    # limit the search to constant and version file
    if (!$versionFile.Name.Contains("constant") -and !$versionFile.Name.Contains("version")) {
      continue
    }
    $content = Get-Content $versionFile -Raw
    # finding where the version number are
    if ($content -match $VERSION_LINE_REGEX) {
        return "$($matches["version"])", $versionFile
    }
  }

  LogWarning "Unable to find version for $modPath"
  return $null
}

function Get-GoModuleProperties($goModPath)
{
  $goModPath = $goModPath -replace "\\", "/"
  if ($goModPath -match "(?<modPath>sdk/(?<serviceDir>(resourcemanager/)?([^/]+/)?(?<modName>[^/]+$)))")
  {
    $modPath = $matches["modPath"]
    $modName = $matches["modName"] # We may need to start readong this from the go.mod file if the path and mod config start to differ
    $serviceDir = $matches["serviceDir"]
    $sdkType = "client"
    if ($modName.StartsWith("arm")) { $sdkType = "mgmt" }

    $modVersion, $versionFile = Get-GoModuleVersionInfo $goModPath

    if (!$modVersion) {
      return $null
    }

    $pkgProp = [PackageProps]::new($modPath, $modVersion, $goModPath, $serviceDir)
    $pkgProp.IsNewSdk = $true
    $pkgProp.SdkType = $sdkType

    $pkgProp | Add-Member -NotePropertyName "VersionFile" -NotePropertyValue $versionFile
    $pkgProp | Add-Member -NotePropertyName "ModuleName" -NotePropertyValue $modName

    return $pkgProp
  }
  return $null
}

# rewrite from artifact-metadata-parsing.ps1 used in RetrievePackages for fetch go single module info
function Get-go-PackageInfoFromPackageFile($pkg, $workingDirectory)
{
    $releaseNotes = ""
    $packageProperties = Get-GoModuleProperties $pkg.Directory

    if ($packageProperties.ChangeLogPath -and $packageProperties.Version)
    {
      $releaseNotes = Get-ChangeLogEntryAsString -ChangeLogLocation $packageProperties.ChangeLogPath `
        -VersionString $packageProperties.Version
    }

    $resultObj = New-Object PSObject -Property @{
      PackageId      = $packageProperties.Name
      PackageVersion = $packageProperties.Version
      ReleaseTag     = "$($packageProperties.Name)/v$($packageProperties.Version)"
      Deployable     = $true
      ReleaseNotes   = $releaseNotes

    }

    return $resultObj
}

function Get-AllPackageInfoFromRepo($serviceDirectory)
{
  $allPackageProps = @()
  $searchPath = Join-Path $RepoRoot "sdk"
  if ($serviceDirectory) {
    $searchPath = Join-Path $searchPath $serviceDirectory
  }

  $pkgFiles = Get-ChildItem -Path $searchPath -Include "go.mod" -Recurse

  foreach ($pkgFile in $pkgFiles)
  {
    $modPropertes = Get-GoModuleProperties $pkgFile.DirectoryName

    if ($modPropertes) {
      $allPackageProps += $modPropertes
    }
  }
  return $allPackageProps
}

function SetPackageVersion ($PackageName, $Version, $ReleaseDate, $PackageProperties, $ReplaceLatestEntryTitle=$true)
{
  if(!$ReleaseDate) {
    $ReleaseDate = Get-Date -Format "yyyy-MM-dd"
  }

  if (!$PackageProperties) {
    $PackageProperties = Get-PkgProperties -PackageName $PackageName
  }

  & "${EngScriptsDir}/Update-ModuleVersion.ps1" `
    -ModulePath $PackageProperties.Name `
    -NewVersionString $Version `
    -ReleaseDate $ReleaseDate `
    -ReplaceLatestEntryTitle $ReplaceLatestEntryTitle
}
