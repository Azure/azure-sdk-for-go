BeforeAll {
  function CompatibleConvertFrom-Yaml {}
  function go {}
  . (Join-Path $PSScriptRoot .. analysis-helpers.ps1)
}

Describe "Invoke-WithGoAnalysisTools" {
  BeforeEach {
    $script:originalPath = $env:PATH
    $script:originalToolchain = $env:GOTOOLCHAIN
    $script:originalGoRoot = $env:GOROOT
    $env:GOROOT = Join-Path $TestDrive "baseline-go"
    $script:baselineGoRoot = $env:GOROOT
    $script:packagePath = Join-Path $TestDrive "package"
    New-Item -ItemType Directory -Path $script:packagePath -Force | Out-Null
    Set-Content (Join-Path $script:packagePath "ci.yml") "extends: {}"
    Mock CompatibleConvertFrom-Yaml {
      return @{ extends = @{ parameters = @{ AnalysisGoVersion = "1.27.2"; LintVersion = "v2.14.0" } } }
    }
    Mock go {
      $global:LASTEXITCODE = 0
      return (Join-Path $TestDrive "toolchain")
    }
  }

  AfterEach {
    $env:PATH = $script:originalPath
    $env:GOTOOLCHAIN = $script:originalToolchain
    $env:GOROOT = $script:originalGoRoot
  }

  It "uses package tooling for the action and restores the job environment" {
    $result = Invoke-WithGoAnalysisTools $script:packagePath {
      param($settings)
      $env:GOTOOLCHAIN | Should -Be "go1.27.2"
      $env:GOROOT | Should -Be (Join-Path $TestDrive "toolchain")
      $env:PATH.Split([IO.Path]::PathSeparator)[0] | Should -Be (Join-Path $TestDrive "toolchain/bin")
      $settings.LintVersion | Should -Be "v2.14.0"
      return 17
    }
    $result | Should -Be 17
    $env:PATH | Should -Be $script:originalPath
    $env:GOTOOLCHAIN | Should -Be $script:originalToolchain
    $env:GOROOT | Should -Be $script:baselineGoRoot
    Should -Invoke go -Exactly -Times 1
  }

  It "does not change tools for a package without overrides" {
    Mock CompatibleConvertFrom-Yaml { return @{ extends = @{ parameters = @{ ServiceDirectory = "another/package" } } } }
    Invoke-WithGoAnalysisTools $script:packagePath {
      $env:PATH | Should -Be $script:originalPath
      $env:GOTOOLCHAIN | Should -Be $script:originalToolchain
      $env:GOROOT | Should -Be $script:baselineGoRoot
    }
    Should -Invoke go -Exactly -Times 0
  }

  It "does not require a ci file for ordinary packages" {
    Remove-Item (Join-Path $script:packagePath "ci.yml")
    Invoke-WithGoAnalysisTools $script:packagePath {
      param($settings)
      $settings.Count | Should -Be 0
      $env:PATH | Should -Be $script:originalPath
    }
    Should -Invoke CompatibleConvertFrom-Yaml -Exactly -Times 0
    Should -Invoke go -Exactly -Times 0
  }

  It "restores the environment when the check throws" {
    { Invoke-WithGoAnalysisTools $script:packagePath { throw "check failed" } } |
      Should -Throw "*check failed*"
    $env:PATH | Should -Be $script:originalPath
    $env:GOTOOLCHAIN | Should -Be $script:originalToolchain
    $env:GOROOT | Should -Be $script:baselineGoRoot
  }

  It "fails explicitly and restores tools when the selected compiler is unavailable" {
    Mock go { $global:LASTEXITCODE = 1; return $null }
    { Invoke-WithGoAnalysisTools $script:packagePath { throw "must not execute" } } |
      Should -Throw "*Unable to resolve*"
    $env:PATH | Should -Be $script:originalPath
    $env:GOTOOLCHAIN | Should -Be $script:originalToolchain
    $env:GOROOT | Should -Be $script:baselineGoRoot
  }

  It "rejects malformed compiler versions before invoking external commands" {
    Mock CompatibleConvertFrom-Yaml { return @{ extends = @{ parameters = @{ AnalysisGoVersion = "latest" } } } }
    { Invoke-WithGoAnalysisTools $script:packagePath {} } | Should -Throw "*Invalid AnalysisGoVersion*"
    Should -Invoke go -Exactly -Times 0
  }
}
