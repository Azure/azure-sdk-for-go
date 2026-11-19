BeforeAll {
  $repoRoot = (Resolve-Path (Join-Path $PSScriptRoot .. .. ..)).Path
  $setHashScript = Join-Path $repoRoot 'eng/scripts/Set-PackageApiReviewHash.ps1'
  $confirmPublishedScript = Join-Path $repoRoot 'eng/scripts/Confirm-GoPackagePublished.ps1'
  $testParserScript = Join-Path $repoRoot 'eng/scripts/Test-GoAPIReviewParser.ps1'
}

Describe 'Go API review bundles' {
  BeforeEach {
    $testRoot = Join-Path $TestDrive ([guid]::NewGuid().ToString())
    $sourceRoot = Join-Path $testRoot 'source'
    $packagePath = Join-Path $sourceRoot 'sdk/service/example'
    $binDirectory = Join-Path $testRoot 'bin'
    $outputDirectory = Join-Path $testRoot 'bundle'
    New-Item -ItemType Directory -Path $packagePath, $binDirectory -Force | Out-Null
    Set-Content -Path (Join-Path $packagePath 'go.mod') -Value 'module github.com/Azure/azure-sdk-for-go/sdk/service/example'
    Set-Content -Path (Join-Path $packagePath 'version.go') -Value 'const version string = "v1.2.3-beta.1"'

    Copy-Item -LiteralPath @(
      (Join-Path $repoRoot 'eng/scripts/Export-APIReviewBundle.ps1'),
      (Join-Path $repoRoot 'eng/scripts/Get-GoModuleVersionInfo.ps1'),
      (Join-Path $repoRoot 'eng/common/scripts/Export-APIViewMarkdown.ps1'),
      (Join-Path $repoRoot 'eng/common/scripts/SemVer.ps1'),
      (Join-Path $repoRoot 'eng/common/scripts/logging.ps1')
    ) -Destination $binDirectory
    $bundleScript = Join-Path $binDirectory 'Export-APIReviewBundle.ps1'
    $parser = Join-Path $testRoot 'parser.ps1'
    @'
param($PackageDirectory, $OutputDirectory)
$review = [PSCustomObject]@{
  PackageName = 'sdk/service/example'
  Language = 'Go'
  ReviewLines = @(
    [PSCustomObject]@{
      Tokens = @(
        [PSCustomObject]@{ Value = 'package' },
        [PSCustomObject]@{ Value = 'example'; HasSuffixSpace = $false }
      )
      Children = @(
        [PSCustomObject]@{
          Tokens = @(
            [PSCustomObject]@{ Value = 'func' },
            [PSCustomObject]@{ Value = 'New()'; HasSuffixSpace = $false }
          )
        }
      )
    }
  )
}
$review | ConvertTo-Json -Depth 10 | Set-Content -Path (Join-Path $OutputDirectory 'example.json')
$global:LASTEXITCODE = 0
'@ | Set-Content -Path $parser
  }

  It 'renders Go token spacing and emits the source version, path, and API hash' {
    $bundle = & $bundleScript -SourceRoot $sourceRoot -PackageName 'sdk/service/example' `
      -OutputDirectory $outputDirectory -ParserPath $parser

    $apiPath = Join-Path $outputDirectory 'api.md'
    (Get-Content -LiteralPath $apiPath -Raw) | Should -Be (@('```go', 'package example', '    func New()', '```') -join "`n")
    $bundle.PackageVersion | Should -Be '1.2.3-beta.1'
    $bundle.ApiHash | Should -Be (Get-FileHash -LiteralPath $apiPath -Algorithm SHA256).Hash.ToLowerInvariant()
    $metadata = Get-Content -LiteralPath (Join-Path $outputDirectory 'api.metadata.yml') -Raw
    $metadata | Should -Match "apiHash: $($bundle.ApiHash)"
    $metadata | Should -Match 'packageVersion: 1\.2\.3-beta\.1'
    $metadata | Should -Match 'packageRelativePath: sdk/service/example'
  }

  It 'keeps the API hash unchanged when only the package version changes' {
    $first = & $bundleScript -SourceRoot $sourceRoot -PackageName 'sdk/service/example' `
      -OutputDirectory $outputDirectory -ParserPath $parser
    Set-Content -Path (Join-Path $packagePath 'version.go') -Value 'const version string = "v1.2.4"'
    $second = & $bundleScript -SourceRoot $sourceRoot -PackageName 'sdk/service/example' `
      -OutputDirectory (Join-Path $testRoot 'second') -ParserPath $parser

    $second.ApiHash | Should -Be $first.ApiHash
    $second.PackageVersion | Should -Be '1.2.4'
  }

  It 'rejects package traversal and a parser result for a different package' {
    { & $bundleScript -SourceRoot $sourceRoot -PackageName 'sdk/../example' `
        -OutputDirectory $outputDirectory -ParserPath $parser } | Should -Throw '*repository-relative SDK module path*'

    $contents = (Get-Content -LiteralPath $parser -Raw).Replace('sdk/service/example', 'sdk/service/other')
    Set-Content -LiteralPath $parser -Value $contents
    { & $bundleScript -SourceRoot $sourceRoot -PackageName 'sdk/service/example' `
        -OutputDirectory $outputDirectory -ParserPath $parser } | Should -Throw '*different package name*'
  }

  It 'fails when the parser does not produce a token file' {
    Set-Content -LiteralPath $parser -Value '$global:LASTEXITCODE = 0'
    { & $bundleScript -SourceRoot $sourceRoot -PackageName 'sdk/service/example' `
        -OutputDirectory $outputDirectory -ParserPath $parser } | Should -Throw '*produced 0 token files*'
  }

  It 'rejects unsupported syntax reported by the parser' -TestCases @(
    @{ Diagnostic = 'unexpected node type *ast.ArrayType' },
    @{ Diagnostic = 'unhandled node type *ast.StarExpr: type Alias *Value' }
  ) {
    param($Diagnostic)
    Add-Content -LiteralPath $parser -Value "`nWrite-Output '$Diagnostic'"
    { & $bundleScript -SourceRoot $sourceRoot -PackageName 'sdk/service/example' `
        -OutputDirectory $outputDirectory -ParserPath $parser } | Should -Throw '*unsupported syntax*'
    (Test-Path -LiteralPath (Join-Path $outputDirectory 'api.metadata.yml')) | Should -BeFalse
  }

  It 'emits an empty API when the module has no public review lines' {
    @'
param($PackageDirectory, $OutputDirectory)
@{ PackageName = 'sdk/service/example'; Language = 'Go' } |
  ConvertTo-Json | Set-Content -Path (Join-Path $OutputDirectory 'example.json')
$global:LASTEXITCODE = 0
'@ | Set-Content -LiteralPath $parser
    & $bundleScript -SourceRoot $sourceRoot -PackageName 'sdk/service/example' `
      -OutputDirectory $outputDirectory -ParserPath $parser | Out-Null

    (Get-Content -LiteralPath (Join-Path $outputDirectory 'api.md') -Raw) |
      Should -Be (@('```go', '', '```') -join "`n")
  }

  It 'places the bundle hash in release package info only for a matching released version' {
    $packageInfoPath = Join-Path $testRoot 'package.json'
    @{
      Name = 'sdk/service/example'
      Version = '1.2.3-beta.1'
      ReleaseStatus = '2026-01-01'
    } | ConvertTo-Json | Set-Content -Path $packageInfoPath

    & $setHashScript -PackageInfoFile $packageInfoPath -SourceRoot $sourceRoot `
      -ToolsDirectory $binDirectory -ParserPath $parser -OutputDirectory $outputDirectory

    $packageInfo = Get-Content -LiteralPath $packageInfoPath -Raw | ConvertFrom-Json
    $packageInfo.ApiHash | Should -Be (Get-FileHash -LiteralPath (Join-Path $outputDirectory 'api.md') -Algorithm SHA256).Hash.ToLowerInvariant()

    $packageInfo.Version = '1.2.4'
    $packageInfo | ConvertTo-Json | Set-Content -LiteralPath $packageInfoPath
    { & $setHashScript -PackageInfoFile $packageInfoPath -SourceRoot $sourceRoot `
        -ToolsDirectory $binDirectory -ParserPath $parser -OutputDirectory $outputDirectory } | Should -Throw '*does not match source version*'

    $packageInfo.ReleaseStatus = 'Unreleased'
    $packageInfo | ConvertTo-Json | Set-Content -LiteralPath $packageInfoPath
    { & $setHashScript -PackageInfoFile $packageInfoPath -SourceRoot $sourceRoot `
        -ToolsDirectory $binDirectory -ParserPath $parser -OutputDirectory $outputDirectory } | Should -Throw '*not marked for release*'
  }
}

Describe 'Go APIView parser compatibility' {
  BeforeEach {
    $testRoot = Join-Path $TestDrive ([guid]::NewGuid().ToString())
    New-Item -ItemType Directory -Path $testRoot | Out-Null
    $parser = Join-Path $testRoot 'parser.ps1'
    @'
param($SourceDirectory, $OutputDirectory)
$aliases = @(
  @{ Name = 'ByteSlice'; Definition = '[]byte'; Methods = @('func(ByteSlice)Len()int', 'func(*ByteSlice)Reset()') },
  @{ Name = 'AnotherSlice'; Definition = '[]byte'; Methods = @('func(AnotherSlice)Len()int', 'func(*AnotherSlice)Reset()') },
  @{ Name = 'Mapping'; Definition = 'map[string]int'; Methods = @('func(Mapping)Size()int') },
  @{ Name = 'Defined'; Definition = '[]byte'; Methods = @('func(Defined)Local()') }
)
$lines = foreach ($alias in $aliases) {
  @{
    LineId = "parserprobe.$($alias.Name)"
    Tokens = @(@{ Value = "type$($alias.Name)$($alias.Definition)" })
    Children = @($alias.Methods | ForEach-Object { @{ Tokens = @(@{ Value = $_ }) } })
  }
}
@{ ReviewLines = @(@{ LineId = 'parserprobe'; Tokens = @(); Children = @($lines) }) } |
  ConvertTo-Json -Depth 10 | Set-Content -Path (Join-Path $OutputDirectory 'parserprobe.json')
$global:LASTEXITCODE = 0
'@ | Set-Content -Path $parser
  }

  It 'accepts aliases with their underlying types and renamed methods' {
    { & $testParserScript -ParserPath $parser -WorkingDirectory $testRoot } | Should -Not -Throw
    @(Get-ChildItem -LiteralPath $testRoot -Filter '.parser-probe-*' -Force).Count | Should -Be 0
  }

  It 'rejects a silently unresolved map alias' {
    $contents = (Get-Content -LiteralPath $parser -Raw).Replace('map[string]int', 'internal/types.StringMap')
    Set-Content -LiteralPath $parser -Value $contents
    { & $testParserScript -ParserPath $parser -WorkingDirectory $testRoot } | Should -Throw '*does not resolve the Mapping alias*'
  }

  It 'rejects methods that still use the original receiver name' {
    $contents = (Get-Content -LiteralPath $parser -Raw).Replace('func(ByteSlice)Len()', 'func(Data)Len()')
    Set-Content -LiteralPath $parser -Value $contents
    { & $testParserScript -ParserPath $parser -WorkingDirectory $testRoot } | Should -Throw '*does not resolve*aliased type*'
  }

  It 'rejects inherited methods on defined types' {
    $contents = (Get-Content -LiteralPath $parser -Raw).Replace(
      "Methods = @('func(Defined)Local()')",
      "Methods = @('func(Defined)Local()', 'func(Defined)Len()int')"
    )
    Set-Content -LiteralPath $parser -Value $contents
    { & $testParserScript -ParserPath $parser -WorkingDirectory $testRoot } | Should -Throw '*incorrectly inherited*'
  }
}

Describe 'Confirming Go package publication' {
  BeforeAll {
    function global:git {
      param([Parameter(ValueFromRemainingArguments = $true)] [object[]] $Arguments)
      $global:LASTEXITCODE = $global:GitExitCode
      return $global:GitRemoteRefs
    }
  }

  AfterAll {
    Remove-Item Function:\git
    Remove-Variable GitExitCode, GitRemoteRefs -Scope Global
  }

  BeforeEach {
    $packageInfoPath = Join-Path $TestDrive 'published-package.json'
    @{ Name = 'sdk/service/example'; Version = '1.2.3' } |
      ConvertTo-Json | Set-Content -Path $packageInfoPath
    $expectedCommit = 'aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa'
    $tag = 'refs/tags/sdk/service/example/v1.2.3'
    $global:GitExitCode = 0
    $global:GitRemoteRefs = @()
  }

  It 'accepts an annotated tag pointing at the source commit' {
    $global:GitRemoteRefs = @(
      "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb`t$tag",
      "$expectedCommit`t${tag}^{}"
    )
    { & $confirmPublishedScript -PackageInfoFile $packageInfoPath -ExpectedCommit $expectedCommit } | Should -Not -Throw
  }

  It 'rejects a missing tag or a tag for a different commit' {
    $global:GitExitCode = 2
    { & $confirmPublishedScript -PackageInfoFile $packageInfoPath -ExpectedCommit $expectedCommit } | Should -Throw '*was not found*'

    $global:GitExitCode = 0
    $global:GitRemoteRefs = @("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb`t$tag")
    { & $confirmPublishedScript -PackageInfoFile $packageInfoPath -ExpectedCommit $expectedCommit } | Should -Throw '*does not point to the source commit*'
  }

  It 'reports failures querying the remote repository' {
    $global:GitExitCode = 128
    { & $confirmPublishedScript -PackageInfoFile $packageInfoPath -ExpectedCommit $expectedCommit } | Should -Throw '*Failed to query published tag*'
  }
}
