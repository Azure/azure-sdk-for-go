BeforeAll {
  $repoRoot = (Resolve-Path (Join-Path $PSScriptRoot .. .. ..)).Path
  $setHashScript = Join-Path $repoRoot 'eng/scripts/Set-PackageApiReviewHash.ps1'
  $confirmPublishedScript = Join-Path $repoRoot 'eng/scripts/Confirm-GoPackagePublished.ps1'
  $testParserScript = Join-Path $repoRoot 'eng/scripts/Test-GoAPIReviewParser.ps1'
  $buildToolsScript = Join-Path $repoRoot 'eng/scripts/Build-APIReviewTools.ps1'
  $buildHubToolsScript = Join-Path $repoRoot 'eng/scripts/Build-APIReviewHubTools.ps1'
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

  It 'supports shipping profile modules' {
    $profilePath = Join-Path $sourceRoot 'profile/p20200901'
    New-Item -ItemType Directory -Path $profilePath -Force | Out-Null
    Copy-Item -LiteralPath (Join-Path $packagePath 'go.mod'), (Join-Path $packagePath 'version.go') -Destination $profilePath
    Set-Content -LiteralPath $parser -Value ((Get-Content -LiteralPath $parser -Raw).Replace('sdk/service/example', 'profile/p20200901'))
    $bundle = & $bundleScript -SourceRoot $sourceRoot -PackageName 'profile/p20200901' `
      -OutputDirectory $outputDirectory -ParserPath $parser
    $bundle.PackageName | Should -Be 'profile/p20200901'
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
    @{ Diagnostic = 'unhandled node type *ast.StarExpr: type Alias *Value' },
    @{ Diagnostic = 'unhandled declaration []byte' },
    @{ Diagnostic = 'unhandled declaration type *ast.ArrayType for *[2]byte' }
  ) {
    param($Diagnostic)
    Add-Content -LiteralPath $parser -Value "`nWrite-Output '$Diagnostic'"
    { & $bundleScript -SourceRoot $sourceRoot -PackageName 'sdk/service/example' `
        -OutputDirectory $outputDirectory -ParserPath $parser } | Should -Throw '*unsupported syntax*'
    (Test-Path -LiteralPath (Join-Path $outputDirectory 'api.metadata.yml')) | Should -BeFalse
  }

  It 'rejects error and fatal diagnostics before publishing a hash' -TestCases @(
    @{ Level = 3 }, @{ Level = 4 }, @{ Level = 'Error' }, @{ Level = 'Fatal' }
  ) {
    param($Level)
    $contents = (Get-Content -LiteralPath $parser -Raw).Replace(
      "Language = 'Go'",
      "Language = 'Go'`n  Diagnostics = @(@{ Level = '$Level'; TargetID = 'example.Error'; Text = 'omitted promoted methods' })"
    )
    Set-Content -LiteralPath $parser -Value $contents
    { & $bundleScript -SourceRoot $sourceRoot -PackageName 'sdk/service/example' `
        -OutputDirectory $outputDirectory -ParserPath $parser } |
      Should -Throw '*example.Error: omitted promoted methods*'
    (Test-Path -LiteralPath (Join-Path $outputDirectory 'api.md')) | Should -BeFalse
    (Test-Path -LiteralPath (Join-Path $outputDirectory 'api.metadata.yml')) | Should -BeFalse
  }

  It 'allows informational and warning diagnostics' {
    $contents = (Get-Content -LiteralPath $parser -Raw).Replace(
      "Language = 'Go'",
      "Language = 'Go'`n  Diagnostics = @(@{ Level = 1; TargetID = 'example.Alias'; Text = 'alias' }, @{ Level = 2; TargetID = 'example.External'; Text = 'external alias' })"
    )
    Set-Content -LiteralPath $parser -Value $contents
    { & $bundleScript -SourceRoot $sourceRoot -PackageName 'sdk/service/example' `
        -OutputDirectory $outputDirectory -ParserPath $parser } | Should -Not -Throw
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
$updated = (Get-Content -LiteralPath (Join-Path $SourceDirectory 'alias.go') -Raw).Contains('T comparable')
$external = if ($updated) { '=[]rune' } else { '=[]byte' }
$direct = if ($updated) { 'typeDirect[]byte' } else { 'typeDirect=[]byte' }
$constraint = if ($updated) { 'Tcomparable' } else { 'Tany' }
$aliases = @(
  @{ Name = 'ByteSlice'; Definition = '=[]byte'; Methods = @('func(ByteSlice)Len()int', 'func(*ByteSlice)Reset()') },
  @{ Name = 'AnotherSlice'; Definition = '=[]byte'; Methods = @('func(AnotherSlice)Len()int', 'func(*AnotherSlice)Reset()') },
  @{ Name = 'Mapping'; Definition = '=map[string]int'; Methods = @('func(Mapping)Size()int') },
  @{ Name = 'Defined'; Definition = '[]byte'; Methods = @('func(Defined)Local()') },
  @{ Name = 'External'; Definition = $external; Methods = @() }
)
$lines = foreach ($alias in $aliases) {
  @{
    LineId = "parserprobe.$($alias.Name)"
    Tokens = @(@{ Value = "type$($alias.Name)$($alias.Definition)" })
    Children = @($alias.Methods | ForEach-Object { @{ Tokens = @(@{ Value = $_ }) } })
  }
}
$lines += @(
  @{ LineId = 'parserprobe.Direct'; Tokens = @(@{ Value = $direct }) },
  @{ LineId = 'parserprobe.Contract'; Tokens = @(@{ Value = "typeContract[$constraint]interface" }) },
  @{ LineId = 'parserprobe.Responder'; Tokens = @(@{ Value = "typeResponder[$constraint]func(T)" }) }
)
foreach ($name in @('StartTestProxy', 'StopTestProxy', 'ExampleClient')) {
  $lines += @{ LineId = "parserprobe-$name"; Tokens = @(@{ Value = "func$name()" }) }
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

  It 'rejects incomplete API contracts: <Case>' -TestCases @(
    @{ Case = 'alias identity'; Original = "'typeDirect=[]byte'"; Replacement = "'typeDirect[]byte'"; Error = '*expected typeDirect=*' },
    @{ Case = 'generic constraints'; Original = "'Tcomparable'"; Replacement = "'Tany'"; Error = '*expected typeContract*Tcomparable*' },
    @{ Case = 'generic defined types'; Original = '"typeResponder[$constraint]func(T)"'; Replacement = "'typeResponderfunc(T)'"; Error = '*expected typeResponder*Tany*' },
    @{ Case = 'dependency versions'; Original = "'=[]rune'"; Replacement = "'=[]byte'"; Error = '*expected typeExternal=*rune*' },
    @{ Case = 'exported function names'; Original = "'StartTestProxy'"; Replacement = "'Suppressed'"; Error = '*omitted the exported StartTestProxy*' }
  ) {
    param($Case, $Original, $Replacement, $Error)
    Set-Content -LiteralPath $parser -Value ((Get-Content -LiteralPath $parser -Raw).Replace($Original, $Replacement))
    { & $testParserScript -ParserPath $parser -WorkingDirectory $testRoot } | Should -Throw $Error
    @(Get-ChildItem -LiteralPath $testRoot -Filter '.parser-probe-*' -Force).Count | Should -Be 0
  }
}

Describe 'Pinned Go API review tooling' {
  BeforeAll {
    function global:git {
      $Arguments = @($args)
      $global:ToolingGitCalls += ($Arguments -join '|')
      $global:LASTEXITCODE = 0
      if ($Arguments -contains 'fetch' -and $global:ToolingFetchFails) {
        $global:LASTEXITCODE = 1
      }
      elseif ($Arguments[0] -eq 'init') {
        New-Item -ItemType Directory -Path $Arguments[-1] | Out-Null
      }
      elseif ($Arguments -contains 'checkout') {
        $checkout = $Arguments[1]
        New-Item -ItemType Directory -Path (Join-Path $checkout 'src/go'), (Join-Path $checkout 'eng/common/scripts') -Force | Out-Null
        foreach ($script in @('Export-APIViewMarkdown.ps1', 'SemVer.ps1', 'logging.ps1')) {
          Set-Content -LiteralPath (Join-Path $checkout "eng/common/scripts/$script") -Value 'pinned tooling'
        }
      }
      elseif ($Arguments -contains 'rev-parse') {
        return $global:ToolingHead
      }
    }

    function global:go {
      $Arguments = @($args)
      Set-Content -LiteralPath $Arguments[$Arguments.IndexOf('-o') + 1] -Value 'parser binary'
      $global:LASTEXITCODE = 0
    }
  }

  AfterAll {
    Remove-Item Function:\git, Function:\go
    Remove-Variable ToolingGitCalls, ToolingFetchFails, ToolingHead -Scope Global
  }

  BeforeEach {
    $testRoot = Join-Path $TestDrive ([guid]::NewGuid().ToString())
    $languageRoot = Join-Path $testRoot 'language'
    $binDirectory = Join-Path $testRoot 'bin'
    New-Item -ItemType Directory -Path (Join-Path $languageRoot 'eng/scripts') -Force | Out-Null
    $pin = 'a' * 40
    @{ commit = $pin } | ConvertTo-Json | Set-Content -LiteralPath (Join-Path $languageRoot 'eng/api-review-tools.json')
    foreach ($script in @('Export-APIReviewBundle.ps1', 'Get-GoModuleVersionInfo.ps1')) {
      Set-Content -LiteralPath (Join-Path $languageRoot "eng/scripts/$script") -Value 'language tooling'
    }
    Set-Content -LiteralPath (Join-Path $languageRoot 'eng/scripts/Test-GoAPIReviewParser.ps1') `
      -Value 'param($ParserPath, $WorkingDirectory); if (-not (Test-Path -LiteralPath $ParserPath)) { throw "Parser missing" }'
    $global:ToolingGitCalls = @()
    $global:ToolingFetchFails = $false
    $global:ToolingHead = $pin
  }

  It 'uses the same immutable pin with or without an orchestration checkout' -TestCases @(
    @{ UseCheckout = $false }, @{ UseCheckout = $true }
  ) {
    param($UseCheckout)
    $parameters = @{ LanguageToolsSourceDirectory = $languageRoot; BinOutputDirectory = $binDirectory }
    if ($UseCheckout) {
      $parameters.ToolsSourceDirectory = Join-Path $testRoot 'orchestration'
    }
    & $buildToolsScript @parameters
    @($global:ToolingGitCalls | Where-Object { $_ -like "*|fetch|*|origin|$pin" }).Count | Should -Be 1
    (Get-Content -LiteralPath (Join-Path $binDirectory 'Export-APIViewMarkdown.ps1') -Raw).Trim() | Should -Be 'pinned tooling'
    @(Get-ChildItem -LiteralPath $binDirectory -Filter '.tools-source-*' -Force).Count | Should -Be 0
  }

  It 'rejects moving refs before fetching tooling' {
    @{ commit = 'main' } | ConvertTo-Json | Set-Content -LiteralPath (Join-Path $languageRoot 'eng/api-review-tools.json')
    { & $buildToolsScript -LanguageToolsSourceDirectory $languageRoot -BinOutputDirectory $binDirectory } |
      Should -Throw '*immutable 40-character commit*'
    $global:ToolingGitCalls.Count | Should -Be 0
  }

  It 'rejects an incorrect checkout or a failed fetch' -TestCases @(
    @{ FailFetch = $false; Error = '*does not match pinned commit*' },
    @{ FailFetch = $true; Error = '*Failed to fetch API review tooling commit*' }
  ) {
    param($FailFetch, $Error)
    $global:ToolingHead = 'b' * 40
    $global:ToolingFetchFails = $FailFetch
    { & $buildToolsScript -LanguageToolsSourceDirectory $languageRoot -BinOutputDirectory $binDirectory } | Should -Throw $Error
    @(Get-ChildItem -LiteralPath $binDirectory -Filter '.tools-source-*' -Force).Count | Should -Be 0
  }
}

Describe 'Go API Review Hub target tooling' {
  BeforeAll {
    function global:git {
      $arguments = @($args)
      $global:HubGitCalls += ($arguments -join '|')
      $global:LASTEXITCODE = 0
      if ($arguments[0] -eq 'init') {
        New-Item -ItemType Directory -Path $arguments[-1] -Force | Out-Null
      }
      elseif ($arguments -contains 'checkout') {
        $checkout = $arguments[1]
        New-Item -ItemType Directory -Path (Join-Path $checkout 'eng/scripts') -Force | Out-Null
        Set-Content -LiteralPath (Join-Path $checkout 'eng/scripts/Build-APIReviewTools.ps1') -Value @'
param($LanguageToolsSourceDirectory, $BinOutputDirectory)
$global:HubToolsRoot = $LanguageToolsSourceDirectory
'@
      }
      elseif ($arguments -contains 'rev-parse') {
        return $global:HubTargetHead
      }
    }
  }

  AfterAll {
    Remove-Item Function:\git
    Remove-Variable HubGitCalls, HubTargetHead, HubToolsRoot -Scope Global
  }

  BeforeEach {
    $testRoot = Join-Path $TestDrive ([guid]::NewGuid().ToString())
    $sourceRoot = Join-Path $testRoot 'source'
    $toolsDirectory = Join-Path $testRoot 'bin'
    New-Item -ItemType Directory -Path (Join-Path $sourceRoot 'eng/scripts'), $toolsDirectory -Force | Out-Null
    Set-Content -LiteralPath (Join-Path $sourceRoot 'eng/scripts/Build-APIReviewTools.ps1') -Value @'
param($LanguageToolsSourceDirectory, $BinOutputDirectory)
$global:HubToolsRoot = $LanguageToolsSourceDirectory
'@
    $targetCommit = 'a' * 40
    $global:HubTargetHead = $targetCommit
    $global:HubGitCalls = @()
    $global:HubToolsRoot = ''
  }

  It 'builds target tools for both <Kind> and target artifacts' -TestCases @(
    @{ Kind = 'base' }, @{ Kind = 'target' }
  ) {
    param($Kind)
    & $buildHubToolsScript -SourceRoot $sourceRoot -RepositoryFullName 'Azure/azure-sdk-for-go' `
      -TargetCommit $targetCommit -Kind $Kind -ToolsDirectory $toolsDirectory -WorkingDirectory $testRoot
    if ($Kind -eq 'target') {
      $global:HubToolsRoot | Should -Be $sourceRoot
      @($global:HubGitCalls | Where-Object { $_ -like '*|fetch|*' }).Count | Should -Be 0
    }
    else {
      $global:HubToolsRoot | Should -Not -Be $sourceRoot
      @($global:HubGitCalls | Where-Object { $_ -like "*|fetch|*|origin|$targetCommit" }).Count | Should -Be 1
      (Test-Path -LiteralPath $global:HubToolsRoot) | Should -BeFalse
    }
  }

  It 'rejects a <Kind> target checkout that does not match the requested commit' -TestCases @(
    @{ Kind = 'base' }, @{ Kind = 'target' }
  ) {
    param($Kind)
    $global:HubTargetHead = 'b' * 40
    { & $buildHubToolsScript -SourceRoot $sourceRoot -RepositoryFullName 'Azure/azure-sdk-for-go' `
        -TargetCommit $targetCommit -Kind $Kind -ToolsDirectory $toolsDirectory -WorkingDirectory $testRoot } |
      Should -Throw '*does not match commit*'
    $global:HubToolsRoot | Should -BeNullOrEmpty
    @(Get-ChildItem -LiteralPath $testRoot -Filter 'go-api-review-target-*' -Directory).Count | Should -Be 0
  }

  It 'rejects an invalid target commit or repository before fetching' {
    { & $buildHubToolsScript -SourceRoot $sourceRoot -RepositoryFullName 'Azure/azure-sdk-for-go' `
        -TargetCommit 'main' -Kind base -ToolsDirectory $toolsDirectory -WorkingDirectory $testRoot } |
      Should -Throw '*40-character target commit*'
    { & $buildHubToolsScript -SourceRoot $sourceRoot -RepositoryFullName 'invalid/repo/extra' `
        -TargetCommit $targetCommit -Kind base -ToolsDirectory $toolsDirectory -WorkingDirectory $testRoot } |
      Should -Throw '*owner/repository*'
    $global:HubGitCalls.Count | Should -Be 0
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
