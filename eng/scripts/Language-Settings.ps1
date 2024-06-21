$Language = "go"
$packagePattern = "go.mod"
$LanguageDisplayName = "go"

# get version from specific files (*constants.go, *version.go)
function Get-GoModuleVersionInfo($modPath)
{
  $NO_PREFIX_VERSION_LINE_REGEX = ".+\s*=\s*`"(?<bad_version>$([AzureEngSemanticVersion]::SEMVER_REGEX))`""
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

    # This is an easy mistake to make (X.Y.Z instead of vX.Y.Z) so add a very clear error log to make debugging easier
    if ($content -match $NO_PREFIX_VERSION_LINE_REGEX) {
        LogError "Version in $versionFile should be 'v$($matches["bad_version"])' not '$($matches["bad_version"])'"
    }
  }

  LogWarning "Unable to find version for $modPath"
  return $null
}

function Get-GoModuleProperties($goModPath)
{
  $goModPath = $goModPath -replace "\\", "/"
  # We should keep this regex in sync with what is in the azure-sdk repo at https://github.com/Azure/azure-sdk/blob/main/eng/scripts/Query-Azure-Packages.ps1#L238
  # The serviceName named capture group is unused but used in azure-sdk, so it's kept here for parity
  if (!$goModPath.Contains("testdata") -and !$goModPath.Contains("sdk/samples") -and $goModPath -match "(?<modPath>(sdk|profile)/(?<serviceDir>(.*?(?<serviceName>[^/]+)/)?(?<modName>[^/]+$)))")
  {
    $modPath = $matches["modPath"]
    $modName = $matches["modName"] # We may need to start reading this from the go.mod file if the path and mod config start to differ
    $serviceDir = $matches["serviceDir"]
    $sdkType = "client"
    if ($modName.StartsWith("arm") -or $modPath.Contains("resourcemanager")) { $sdkType = "mgmt" }

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

    # We have some cases when processing service directories that non-shipping projects like perfdata
    # we just want to exclude them as opposed to returning a property with invalid data.
    if (!$packageProperties) {
      return $null
    }

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
  $pkgFiles = @()
  if ($serviceDirectory) {
    $searchPath = Join-Path $searchPath $serviceDirectory "go.mod"
    [array]$pkgFiles = @(Get-ChildItem $searchPath)
  } else {
    # If service directory is not passed in, find all modules
    [array]$pkgFiles = Get-ChildItem -Path $searchPath -Include "go.mod" -Recurse
  }

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


function Find-Go-Artifacts-For-Apireview($ArtifactPath, $PackageName)
{
  $artifact = @(Get-ChildItem -Path (Join-Path $ArtifactPath $PackageName) -Filter "*.gosource")
  if ($artifact)
  {
    $packages = @{
      $artifact.FullName = $artifact.FullName
    }
    return $packages
  }
  return $null
}
