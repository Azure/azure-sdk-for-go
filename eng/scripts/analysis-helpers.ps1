function Get-GoAnalysisPackageSets {
  [CmdletBinding()]
  param(
    [string[]] $PackageDirectories,
    [bool] $FilterToChangedFiles = $false
  )

  $packages = @($PackageDirectories | Where-Object { $_ } | ForEach-Object {
    ($_ -replace "\\", "/").TrimEnd("/")
  })

  if (!$FilterToChangedFiles -or !$packages) {
    return [PSCustomObject]@{
      AnalysisPackages = $packages
      DependencyCheckPackages = $packages
    }
  }

  $changedFiles = @(Get-ChangedFiles -DiffFilterType '' | ForEach-Object { $_ -replace "\\", "/" })
  $analysisPackages = @()
  $dependencyCheckPackages = @()

  foreach ($packageDirectory in $packages) {
    $packagePrefix = "$packageDirectory/"
    $packageChangedFiles = @($changedFiles | Where-Object {
      $_.StartsWith($packagePrefix, [System.StringComparison]::OrdinalIgnoreCase)
    })

    $hasAnalysisChanges = $false
    $hasDependencyChanges = $false
    foreach ($file in $packageChangedFiles) {
      $fileName = Split-Path $file -Leaf
      $isModuleFile = $fileName -in @("go.mod", "go.sum")

      if ($isModuleFile) {
        $hasDependencyChanges = $true
      }
      if ($isModuleFile -or $file.EndsWith(".go", [System.StringComparison]::OrdinalIgnoreCase)) {
        $hasAnalysisChanges = $true
      }
    }

    if ($hasAnalysisChanges) {
      $analysisPackages += $packageDirectory
    }
    if ($hasDependencyChanges) {
      $dependencyCheckPackages += $packageDirectory
    }
  }

  return [PSCustomObject]@{
    AnalysisPackages = $analysisPackages
    DependencyCheckPackages = $dependencyCheckPackages
  }
}

# Package-local compiler overrides also apply to tools whose own go.mod targets an older Go.
function Invoke-WithGoAnalysisTools {
  [CmdletBinding()]
  param(
    [Parameter(Mandatory = $true)][string] $PackageDirectory,
    [Parameter(Mandatory = $true)][scriptblock] $Action
  )

  $packagePath = (Resolve-Path $PackageDirectory -ErrorAction Stop).Path
  $ciPath = Join-Path $packagePath "ci.yml"
  $settings = @{}
  if (Test-Path $ciPath) {
    $ci = CompatibleConvertFrom-Yaml -Content (Get-Content $ciPath -Raw)
    if ($ci.extends.parameters) {
      $settings = $ci.extends.parameters
    }
  }

  $originalPath = $env:PATH
  $originalToolchain = $env:GOTOOLCHAIN
  $originalGoRoot = $env:GOROOT
  try {
    if ($settings.AnalysisGoVersion) {
      if ($settings.AnalysisGoVersion -notmatch '^\d+\.\d+\.\d+$') {
        throw "Invalid AnalysisGoVersion in ${ciPath}: $($settings.AnalysisGoVersion)"
      }
      $env:GOTOOLCHAIN = "go$($settings.AnalysisGoVersion)"
      $goRoot = & go env GOROOT
      if ($LASTEXITCODE -ne 0 -or !$goRoot) {
        throw "Unable to resolve $env:GOTOOLCHAIN for $packagePath"
      }
      $env:GOROOT = $goRoot
      $env:PATH = "$(Join-Path $goRoot 'bin')$([IO.Path]::PathSeparator)$originalPath"
    }
    & $Action $settings
  }
  finally {
    $env:PATH = $originalPath
    $env:GOTOOLCHAIN = $originalToolchain
    $env:GOROOT = $originalGoRoot
  }
}
