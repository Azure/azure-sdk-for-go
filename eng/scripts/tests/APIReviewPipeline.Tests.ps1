BeforeAll {
  $repoRoot = (Resolve-Path (Join-Path $PSScriptRoot .. .. ..)).Path
  $revisionScript = Join-Path $repoRoot 'eng/common/scripts/Create-APIViewRevision.ps1'
  $approvalScript = Join-Path $repoRoot 'eng/common/scripts/Get-PackageApprovalStatus.ps1'
  $markReleasedScript = Join-Path $repoRoot 'eng/common/scripts/Mark-PackageReleased.ps1'

  function global:azsdk {
    param([Parameter(ValueFromRemainingArguments = $true)] [object[]] $Arguments)
    $global:LASTEXITCODE = 0
    if ($Arguments.Count -eq 1 -and $Arguments[0] -eq '--version') {
      return '0.6.38'
    }
    $global:AzSdkArguments = @($Arguments)
    return $global:AzSdkResponse
  }

  function global:az {
    $global:LASTEXITCODE = 0
    return '{"accessToken":"test-token"}'
  }

  function global:Invoke-WebRequest {
    param([string] $Method, [string] $Uri, [object] $Body, [hashtable] $Headers, [int] $MaximumRetryCount)
    $global:ApiViewRequests += [PSCustomObject]@{
      Uri = $Uri
      Body = $Body.ReadAsStringAsync().GetAwaiter().GetResult()
    }
    return [PSCustomObject]@{ Content = 'created'; StatusCode = 200 }
  }
}

AfterAll {
  Remove-Item Function:\azsdk, Function:\az, Function:\Invoke-WebRequest
  Remove-Variable AzSdkArguments, AzSdkResponse, ApiViewRequests, LanguageShort -Scope Global
}

Describe 'Go API review release integration' {
  BeforeEach {
    $testRoot = Join-Path $TestDrive ([guid]::NewGuid().ToString())
    $artifactRoot = Join-Path $testRoot 'artifacts'
    $artifactDirectory = Join-Path $artifactRoot 'sdk/service/example'
    $packageInfoDirectory = Join-Path $testRoot 'PackageInfo/sdk/service'
    New-Item -ItemType Directory -Path $artifactDirectory, $packageInfoDirectory -Force | Out-Null
    Set-Content -Path (Join-Path $artifactDirectory 'example.gosource') -Value 'source archive'
    $packageInfoFile = Join-Path $packageInfoDirectory 'example.json'
    $apiHash = 'a' * 64
    @{
      Name = 'sdk/service/example'
      ArtifactName = 'sdk/service/example'
      Version = '1.2.3'
      SdkType = 'client'
      ReleaseStatus = '2026-01-01'
      ApiHash = $apiHash
    } | ConvertTo-Json | Set-Content -Path $packageInfoFile
    $global:LanguageShort = 'go'
    $global:AzSdkArguments = @()
    $global:ApiViewRequests = @()
  }

  It 'creates an APIView revision from the existing Go source artifact' {
    & $revisionScript -ArtifactPath $artifactRoot -PackageInfoFiles @($packageInfoFile) `
      -SourceBranch 'main' -DefaultBranch 'main' -RepoName 'Azure/azure-sdk-for-go' -BuildId '123'

    $global:ApiViewRequests.Count | Should -Be 1
    $global:ApiViewRequests[0].Uri | Should -Be 'https://apiview.dev/autoreview/upload'
    $global:ApiViewRequests[0].Body | Should -Match '(?s)name="?packageVersion"?.*?1\.2\.3'
    $global:ApiViewRequests[0].Body | Should -Match '(?s)name="?setReleaseTag"?.*?False'
  }

  It 'creates revisions for multiple Go package-info files' {
    $secondArtifactDirectory = Join-Path $artifactRoot 'sdk/service/second'
    New-Item -ItemType Directory -Path $secondArtifactDirectory -Force | Out-Null
    Set-Content -Path (Join-Path $secondArtifactDirectory 'second.gosource') -Value 'source archive'
    $secondPackageInfoFile = Join-Path $packageInfoDirectory 'second.json'
    @{
      Name = 'sdk/service/second'
      ArtifactName = 'sdk/service/second'
      Version = '2.0.0'
      SdkType = 'client'
      ReleaseStatus = '2026-01-01'
    } | ConvertTo-Json | Set-Content -Path $secondPackageInfoFile

    & $revisionScript -ArtifactPath $artifactRoot -PackageInfoFiles @($packageInfoFile, $secondPackageInfoFile) `
      -SourceBranch 'main' -DefaultBranch 'main' -RepoName 'Azure/azure-sdk-for-go' -BuildId '123'

    $global:ApiViewRequests.Count | Should -Be 2
    $global:ApiViewRequests[0].Body | Should -Match 'example\.gosource'
    $global:ApiViewRequests[1].Body | Should -Match 'second\.gosource'
  }

  It 'checks approval and records the same Go package version and API hash' {
    $global:AzSdkResponse = '{"operation_status":"Succeeded","result":{"isApproved":true,"finalSource":"reviewHub","reason":"approved"}}'
    & $approvalScript -PackageInfoFiles @($packageInfoFile) -AzSdkExePath 'azsdk'
    ($global:AzSdkArguments -join '|') | Should -Match "\-\-language\|go\|\-\-package-name\|sdk/service/example\|\-\-package-version\|1\.2\.3"
    ($global:AzSdkArguments -join '|') | Should -Match "\-\-api-hash\|$apiHash"

    $global:AzSdkResponse = '{"operation_status":"Succeeded","api_review_hub":{"isReleased":true},"api_view":{"isReleased":true}}'
    & $markReleasedScript -PackageInfoFiles @($packageInfoFile) -AzSdkExePath 'azsdk'
    ($global:AzSdkArguments -join '|') | Should -Match "\-\-language\|go\|\-\-package-name\|sdk/service/example\|\-\-package-version\|1\.2\.3"
    ($global:AzSdkArguments -join '|') | Should -Match "\-\-api-hash\|$apiHash"
  }

  It 'blocks release when the overall approval is denied' {
    $global:AzSdkResponse = '{"operation_status":"Succeeded","result":{"isApproved":false,"finalSource":"reviewHub","reason":"packageNameNotApproved"}}'
    { & $approvalScript -PackageInfoFiles @($packageInfoFile) -AzSdkExePath 'azsdk' } |
      Should -Throw '*Package approval failed*'
  }
}
