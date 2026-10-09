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
