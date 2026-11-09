BeforeAll {
  function Get-ChangedFiles {}
  . (Join-Path $PSScriptRoot .. analysis-helpers.ps1)
}

Describe "Get-GoAnalysisPackageSets" {
  BeforeEach {
    Mock Get-ChangedFiles {
      return @(
        "sdk/service/code/fake/server.go"
        "sdk/service/code/testdata/_metadata.json"
        "sdk/service/dependency/go.mod"
        "sdk/service/dependency/go.sum"
        "sdk/service/metadata/testdata/_metadata.json"
      )
    }
  }

  It "filters analysis and dependency packages by changed file type" {
    $result = Get-GoAnalysisPackageSets `
      -PackageDirectories @(
        "sdk/service/code"
        "sdk/service/dependency"
        "sdk/service/metadata"
      ) `
      -FilterToChangedFiles $true

    $result.AnalysisPackages -join "," | Should -Be "sdk/service/code,sdk/service/dependency"
    $result.DependencyCheckPackages -join "," | Should -Be "sdk/service/dependency"
    Should -Invoke Get-ChangedFiles -Times 1 -Exactly
  }

  It "preserves all packages outside pull request filtering" {
    $result = Get-GoAnalysisPackageSets `
      -PackageDirectories @("sdk/service/one", "sdk/service/two") `
      -FilterToChangedFiles $false

    $result.AnalysisPackages -join "," | Should -Be "sdk/service/one,sdk/service/two"
    $result.DependencyCheckPackages -join "," | Should -Be "sdk/service/one,sdk/service/two"
    Should -Invoke Get-ChangedFiles -Times 0 -Exactly
  }
}
