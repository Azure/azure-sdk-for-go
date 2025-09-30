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
      return "$($matches["version"])", $versionFile.ToString()
    }

    # This is an easy mistake to make (X.Y.Z instead of vX.Y.Z) so add a very clear error log to make debugging easier
    if ($content -match $NO_PREFIX_VERSION_LINE_REGEX) {
      LogError "Version in $versionFile should be 'v$($matches["bad_version"])' not '$($matches["bad_version"])'"
    }
  }

  LogWarning "Unable to find version for $modPath"
  return $null
}

# returns the major version suffix (e.g. 2) from a module's identity.
# if there's no major version, $null is returned
function Get-GoModuleMajorVersion($modPath)
{
  $content = Get-Content $modPath -Raw

  # module github.com/Azure/azure-sdk-for-go/sdk/azcore/v2
  # we want to return the 2. note that we must enable multi-line
  # mode so that the ^ and $ anchors properly work.
  if ($content -match "(?m)^module\s+\S+/v(?<majorVersion>\d+)\s*$") {
    return "$($matches["majorVersion"])"
  }
  return $null
}

function Get-GoModuleProperties($goModPath)
{
  $goModPath = $goModPath -replace "\\", "/"
  # We should keep this regex in sync with what is in the azure-sdk repo at https://github.com/Azure/azure-sdk/blob/main/eng/scripts/Query-Azure-Packages.ps1#L238
  # The serviceName named capture group is unused but used in azure-sdk, so it's kept here for parity
  if (!$goModPath.Contains("testdata") -and !$goModPath.Contains("sdk/samples") -and $goModPath -match "(?<modPath>(sdk|profile|eng)/(?<serviceDir>(.*?(?<serviceName>[^/]+)/)?(?<modName>[^/]+$)))")
  {
    $modPath = $matches["modPath"]
    $modName = $matches["modName"] # We may need to start reading this from the go.mod file if the path and mod config start to differ
    $serviceDir = $matches["serviceDir"]
    $sdkType = "client"
    if ($modName.StartsWith("arm") -or $modPath.Contains("resourcemanager")) { $sdkType = "mgmt" }
    if ($modPath.Contains("eng/tools")) { $sdkType = "eng" }

    $modVersion, $versionFile = Get-GoModuleVersionInfo $goModPath
    if (!$modVersion -and $sdkType -ne "eng") {
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

function EvaluateCIParam {
  param(
    [HashTable]$ParsedYmlContent,
    [PackageProps]$pkgProp,
    [string]$ParamName,
    [bool]$DefaultValue
  )

  $paramPresent = GetValueSafelyFrom-Yaml $ParsedYmlContent @("extends", "parameters", $ParamName)

  if ($null -ne $paramPresent) {
    $parsedBool = $null

    if ([bool]::TryParse($paramPresent, [ref]$parsedBool)) {
      $pkgProp.CIParameters[$ParamName] = $parsedBool
    }
  }
  else {
    $pkgProp.CIParameters[$ParamName] = $DefaultValue
  }
}


<#
.DESCRIPTION
This function resolves a filter string to a directly invokable directory or directories
that can be assembled by a go binary.
#>
function ResolveSearchPaths {
  param (
      [Parameter(Mandatory=$true)]
      $FilterString
  )

  $resolvedPaths = @()
  $filters = $FilterString.Split(",")

  foreach($filter in $filters) {
    if ($filter.StartsWith("sdk") -or $filter.StartsWith("eng")) {
      $resolvedPaths += (Join-Path $RepoRoot $filter)
    }
    else {
      $resolvedPaths += (Join-Path $RepoRoot "sdk" $filter)
    }
  }

  return ,$resolvedPaths
}


# This parameter can be a straightforward filter string EG "sdk/template/aztemplate,sdk/core/azcore".
# However the Save-Package-Properties call for a service directory context EXPLICITLY passes -ServiceDirectory
# when passing the string through to the function. Due to that, we can't name this the more appropraite "filterString"
# We have to meet that function signature that is called from Get-AllPkgProperties in `Package-Properties.ps1`.
# This $ServiceDirectory argument can actually be a comma separated list of package paths OR the standard service directories,
# but until we make a change over in eng/common/scripts/Package-Properties.ps1 to support that, we will just use the name $ServiceDirectory
# here.
function Get-AllPackageInfoFromRepo($ServiceDirectory)
{
  $allPackageProps = @()
  $pkgFiles = @()

  if ($ServiceDirectory) {
    $searchPaths = ResolveSearchPaths $ServiceDirectory

    foreach ($searchPath in $searchPaths) {
      $pkgFiles += @(Get-ChildItem (Join-Path $searchPath "go.mod"))
    }
  }
  else {
    $searchPath = Join-Path $RepoRoot "sdk"
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

  # populate ci parameters for each package
  foreach ($pkgProp in $allPackageProps) {
    $pkgProp.ArtifactName = $pkgProp.Name
    $ciProps = $pkgProp.GetCIYmlForArtifact()

    # UsePipelineProxy - installs and runs the test proxy in ci.tests.yml, defaults true
    # NonShipping - activate verify changelog in analyze, defaults false
    # IsSdkLibrary - activates Detect API Changes, enables save-package-properties and enables Create API Review steps, defaults true
    # LicenseCheck - activates license check in analyze, defaults true
    if ($ciProps) {
      EvaluateCIParam $ciProps.ParsedYml $pkgProp "UsePipelineProxy" $true
      EvaluateCIParam $ciProps.ParsedYml $pkgProp "NonShipping" $false
      EvaluateCIParam $ciProps.ParsedYml $pkgProp "IsSdkLibrary" $true
      EvaluateCIParam $ciProps.ParsedYml $pkgProp "LicenseCheck" $true
    }
    # if we don't have a ci yml, just set the defaults
    else {
      $pkgProp.CIParameters["UsePipelineProxy"] = $true
      $pkgProp.CIParameters["NonShipping"] = $false
      $pkgProp.CIParameters["IsSdkLibrary"] = $true
      $pkgProp.CIParameters["LicenseCheck"] = $true
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
