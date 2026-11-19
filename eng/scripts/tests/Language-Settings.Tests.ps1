BeforeAll {
  . (Join-Path $PSScriptRoot .. .. common scripts SemVer.ps1)
  . (Join-Path $PSScriptRoot .. Language-Settings.ps1)
  function LogError($Message) {}
  function LogWarning($Message) {}
}

Describe "Get-GoModuleVersionInfo" {
  BeforeEach {
    $modPath = Join-Path $TestDrive ([guid]::NewGuid().ToString())
    New-Item -ItemType Directory -Path $modPath -Force | Out-Null
    Mock LogError {}
    Mock LogWarning {}
  }

  It "ignores the service version <Value> in constants.go" -TestCases @(
    @{ Value = "8.0.21" }
    @{ Value = "v8.0.21" }
  ) {
    param($Value)
    Set-Content -Path (Join-Path $modPath constants.go) -Value "ServerVersionEight021 ServerVersion = `"$Value`""
    $versionFile = Join-Path $modPath version.go
    Set-Content -Path $versionFile -Value 'moduleVersion = "v2.0.0"'

    $version, $file = Get-GoModuleVersionInfo $modPath

    $version | Should -Be "2.0.0"
    $file | Should -Be $versionFile
    Should -Invoke LogError -Times 0 -Exactly
    Should -Invoke LogWarning -Times 0 -Exactly
  }

  It "recognizes <Declaration> in <FileName>" -TestCases @(
    @{ Declaration = 'moduleVersion = "v1.2.3"'; Expected = "1.2.3"; FileName = "constants.go" }
    @{ Declaration = 'const version string = "v1.2.3-beta.1"'; Expected = "1.2.3-beta.1"; FileName = "version.go" }
    @{ Declaration = 'serviceLibVersion = "v1.6.0-beta.4"'; Expected = "1.6.0-beta.4"; FileName = "version.go" }
  ) {
    param($Declaration, $Expected, $FileName)
    $versionFile = Join-Path $modPath $FileName
    Set-Content -Path $versionFile -Value $Declaration

    $version, $file = Get-GoModuleVersionInfo $modPath

    $version | Should -Be $Expected
    $file | Should -Be $versionFile
    Should -Invoke LogError -Times 0 -Exactly
    Should -Invoke LogWarning -Times 0 -Exactly
  }

  It "still reports a missing v prefix for <Name>" -TestCases @(
    @{ Name = "moduleVersion" }
    @{ Name = "version" }
    @{ Name = "serviceLibVersion" }
  ) {
    param($Name)
    Set-Content -Path (Join-Path $modPath version.go) -Value "const $Name string = `"1.2.3`""

    Get-GoModuleVersionInfo $modPath | Should -BeNullOrEmpty

    Should -Invoke LogError -Times 1 -Exactly -ParameterFilter {
      $Message -like "*should be 'v1.2.3' not '1.2.3'"
    }
    Should -Invoke LogWarning -Times 1 -Exactly
  }

  It "does not treat a commented declaration as the module version" {
    Set-Content -Path (Join-Path $modPath version.go) -Value @'
// moduleVersion = "v9.9.9"
const moduleVersion = "v2.0.0"
'@

    $version, $file = Get-GoModuleVersionInfo $modPath

    $version | Should -Be "2.0.0"
    Should -Invoke LogError -Times 0 -Exactly
  }
}
