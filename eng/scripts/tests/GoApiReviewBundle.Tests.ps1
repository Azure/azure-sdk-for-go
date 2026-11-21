BeforeAll {
  $repoRoot = (Resolve-Path (Join-Path $PSScriptRoot .. .. ..)).Path
  $exportScript = Join-Path $repoRoot 'eng/scripts/Export-GoApiReviewBundle.ps1'
  $setHashScript = Join-Path $repoRoot 'eng/scripts/Set-GoPackageApiReviewHash.ps1'
}

AfterAll {
  Remove-Variable BundleTestVersion, BundleTestBadHash, BundleTestPath -Scope Global -ErrorAction SilentlyContinue
}

Describe 'Go API review bundles' {
  BeforeEach {
    $testRoot = Join-Path $TestDrive ([guid]::NewGuid().ToString())
    $sourceRoot = Join-Path $testRoot 'source'
    $packagePath = Join-Path $sourceRoot 'sdk/service/example'
    $outputDirectory = Join-Path $testRoot 'output'
    New-Item -ItemType Directory -Path $packagePath -Force | Out-Null
    Set-Content -Path (Join-Path $packagePath 'go.mod') -Value 'module example'
    $apigen = Join-Path $testRoot 'apigen.ps1'
    @'
$output = $args[4]
Set-Content -LiteralPath (Join-Path $output 'api.md') -Value 'public API'
$hash = (Get-FileHash -LiteralPath (Join-Path $output 'api.md') -Algorithm SHA256).Hash.ToLowerInvariant()
if ($global:BundleTestBadHash) { $hash = '0' * 64 }
@(
  "apiHash: $hash"
  "packageVersion: $global:BundleTestVersion"
  "packageRelativePath: `"$global:BundleTestPath`""
) | Set-Content -LiteralPath (Join-Path $output 'api.metadata.yml')
$global:LASTEXITCODE = 0
'@ | Set-Content -Path $apigen
    $global:BundleTestVersion = '1.2.3-beta.1'
    $global:BundleTestBadHash = $false
    $global:BundleTestPath = 'sdk/service/example'
    $packageInfoFile = Join-Path $testRoot 'package.json'
    @{
      Name = 'sdk/service/example'
      Version = '1.2.3-beta.1'
      ReleaseStatus = '2026-01-01'
    } | ConvertTo-Json | Set-Content -Path $packageInfoFile
  }

  It 'adds the generated bundle hash to the matching package info' {
    & $setHashScript -PackageInfoFile $packageInfoFile -SourceRoot $sourceRoot `
      -ApigenPath $apigen -OutputDirectory $outputDirectory

    $info = Get-Content -Raw $packageInfoFile | ConvertFrom-Json
    $info.ApiHash | Should -Be (Get-FileHash (Join-Path $outputDirectory 'api.md') -Algorithm SHA256).Hash.ToLowerInvariant()
  }

  It 'rejects a package version that differs from the generated metadata' {
    $global:BundleTestVersion = '1.2.4'
    { & $setHashScript -PackageInfoFile $packageInfoFile -SourceRoot $sourceRoot `
        -ApigenPath $apigen -OutputDirectory $outputDirectory } | Should -Throw '*does not match source version*'
    (Get-Content -Raw $packageInfoFile | ConvertFrom-Json).PSObject.Properties['ApiHash'] | Should -BeNullOrEmpty
  }

  It 'rejects a metadata hash that does not match the report' {
    $global:BundleTestBadHash = $true
    { & $exportScript -SourceRoot $sourceRoot -PackageName 'sdk/service/example' `
        -OutputDirectory $outputDirectory -ApigenPath $apigen } | Should -Throw '*does not match api.md*'
  }

  It 'rejects metadata for a different package' {
    $global:BundleTestPath = 'sdk/service/other'
    { & $exportScript -SourceRoot $sourceRoot -PackageName 'sdk/service/example' `
        -OutputDirectory $outputDirectory -ApigenPath $apigen } | Should -Throw '*instead of*'
  }

  It 'rejects package paths outside the SDK' {
    { & $exportScript -SourceRoot $sourceRoot -PackageName 'sdk/../example' `
        -OutputDirectory $outputDirectory -ApigenPath $apigen } | Should -Throw '*repository-relative SDK module path*'
  }
}
