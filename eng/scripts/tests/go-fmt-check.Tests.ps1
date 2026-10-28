BeforeAll {
  $script:formatScript = Join-Path $PSScriptRoot .. go-fmt-check.ps1
  function gofmt {}
}

Describe "Go formatting check" {
  BeforeEach {
    $script:firstPackage = Join-Path $TestDrive "first"
    $script:secondPackage = Join-Path $TestDrive "second"
    New-Item -ItemType Directory -Path $script:firstPackage, $script:secondPackage -Force | Out-Null
    Mock Write-Host {}
    Mock Write-Error {}
  }

  It "succeeds when formatting is clean" {
    Mock gofmt { $global:LASTEXITCODE = 0 }
    & $script:formatScript -Packages $script:firstPackage
    Should -Invoke Write-Error -Times 0 -Exactly
    Should -Invoke Write-Host -Times 1 -Exactly -ParameterFilter {
      $Object -like "Format check succeeded.*"
    }
  }

  It "reports differences from exit code <Code> and checks all packages" -ForEach @(
    @{ Code = 0 }
    @{ Code = 1 }
  ) {
    Mock gofmt {
      $global:LASTEXITCODE = $Code
      "misformatted.go"
      "diff misformatted.go.orig misformatted.go"
    }
    & $script:formatScript -Packages "$script:firstPackage,$script:secondPackage"
    Should -Invoke gofmt -Times 2 -Exactly
    Should -Invoke Write-Error -Times 1 -Exactly
    Should -Invoke Write-Host -Times 2 -Exactly -ParameterFilter {
      $Object -eq " -> misformatted.go"
    }
    Should -Invoke Write-Host -Times 2 -Exactly -ParameterFilter {
      $Object -eq " -> diff misformatted.go.orig misformatted.go"
    }
  }

  It "throws on a tool error and restores the working directory" {
    Mock gofmt { $global:LASTEXITCODE = 2 }
    $originalLocation = Get-Location
    { & $script:formatScript -Packages $script:firstPackage } |
      Should -Throw "*gofmt failed*exit code 2*"
    (Get-Location).Path | Should -Be $originalLocation.Path
    Should -Invoke Write-Host -Times 0 -Exactly -ParameterFilter {
      $Object -like "Format check succeeded.*"
    }
  }
}
