BeforeAll {
  $repoRoot = (Resolve-Path (Join-Path $PSScriptRoot .. .. ..)).Path
  $approvalScript = Join-Path $repoRoot 'eng/scripts/Get-GoPackageApprovalStatus.ps1'

  function global:azsdk {
    param([Parameter(ValueFromRemainingArguments = $true)] [object[]] $Arguments)
    if ($Arguments.Count -eq 1 -and $Arguments[0] -eq '--version') {
      $global:LASTEXITCODE = 0
      return '0.6.38'
    }
    $global:AzSdkArguments = @($Arguments)
    $result = $global:AzSdkResponse | ConvertFrom-Json
    $global:LASTEXITCODE = if ($result.result.isApproved) { 0 } else { 1 }
    return $global:AzSdkResponse
  }
}

AfterAll {
  Remove-Item Function:\azsdk -ErrorAction SilentlyContinue
  Remove-Variable AzSdkArguments, AzSdkResponse, LanguageShort -Scope Global -ErrorAction SilentlyContinue
}

Describe 'Go release approval' {
  BeforeEach {
    $testRoot = Join-Path $TestDrive ([guid]::NewGuid().ToString())
    $packageInfoDirectory = Join-Path $testRoot 'PackageInfo/sdk/service'
    New-Item -ItemType Directory -Path $packageInfoDirectory -Force | Out-Null
    $packageInfoFile = Join-Path $packageInfoDirectory 'example.json'
    @{
      Name = 'sdk/service/example'
      ArtifactName = 'sdk/service/example'
      Version = '1.2.3'
      SdkType = 'client'
      IsNewSdk = $true
      ReleaseStatus = '2026-01-01'
    } | ConvertTo-Json | Set-Content -Path $packageInfoFile
    $global:LanguageShort = 'go'
    $global:AzSdkArguments = @()
  }

  It 'accepts APIView approval without an API hash' {
    $global:AzSdkResponse = '{"operation_status":"Succeeded","result":{"isApproved":true,"finalSource":"APIView","reason":"approved"}}'
    & $approvalScript -PackageInfoFiles @($packageInfoFile) -AzSdkExePath 'azsdk'
    ($global:AzSdkArguments -join '|') | Should -Match "\-\-language\|go\|\-\-package-name\|sdk/service/example\|\-\-package-version\|1\.2\.3"
    $global:AzSdkArguments | Should -Not -Contain '--api-hash'
  }

  It 'blocks release when <Reason> is returned' -TestCases @(
    @{ Reason = 'missingApiHash' }
    @{ Reason = 'packageNameNotApproved' }
  ) {
    param($Reason)
    $global:AzSdkResponse = @{
      operation_status = 'Succeeded'
      result = @{ isApproved = $false; finalSource = 'reviewHub'; reason = $Reason }
    } | ConvertTo-Json -Compress
    { & $approvalScript -PackageInfoFiles @($packageInfoFile) -AzSdkExePath 'azsdk' } |
      Should -Throw '*not approved for release*'
    $global:AzSdkArguments | Should -Not -Contain '--api-hash'
  }

  It 'allows a prerelease with an approved package name while API approval is pending' {
    $packageInfo = Get-Content $packageInfoFile -Raw | ConvertFrom-Json
    $packageInfo.Version = '1.2.3-beta.1'
    $packageInfo | ConvertTo-Json | Set-Content -Path $packageInfoFile
    $global:AzSdkResponse = @{
      operation_status = 'Succeeded'
      result = @{
        isApproved = $false
        finalSource = 'reviewHub'
        reason = 'packageNameApproved'
        apiView = @{ isApproved = $false; packageNameApproved = $true; reason = 'packageNameApproved' }
      }
    } | ConvertTo-Json -Compress

    & $approvalScript -PackageInfoFiles @($packageInfoFile) -AzSdkExePath 'azsdk'

    $global:AzSdkArguments | Should -Not -Contain '--api-hash'
  }

  It 'allows a prerelease when only API Review Hub approves the package name' {
    $packageInfo = Get-Content $packageInfoFile -Raw | ConvertFrom-Json
    $packageInfo.Version = '1.2.3-beta.1'
    $packageInfo | ConvertTo-Json | Set-Content -Path $packageInfoFile
    $global:AzSdkResponse = @{
      operation_status = 'Succeeded'
      result = @{
        isApproved = $false
        finalSource = 'reviewHub'
        reason = 'packageNameApproved'
        reviewHub = @{ isPackageNameApproved = $true }
        apiView = @{ isApproved = $false; packageNameApproved = $false; reason = 'packageNamePending' }
      }
    } | ConvertTo-Json -Compress

    & $approvalScript -PackageInfoFiles @($packageInfoFile) -AzSdkExePath 'azsdk'

    $global:AzSdkArguments | Should -Not -Contain '--api-hash'
  }

  It 'blocks a prerelease when neither backend approves the package name' {
    $packageInfo = Get-Content $packageInfoFile -Raw | ConvertFrom-Json
    $packageInfo.Version = '1.2.3-beta.1'
    $packageInfo | ConvertTo-Json | Set-Content -Path $packageInfoFile
    $global:AzSdkResponse = @{
      operation_status = 'Succeeded'
      result = @{
        isApproved = $false
        finalSource = 'reviewHub'
        reason = 'packageNamePending'
        reviewHub = @{ isPackageNameApproved = $false }
        apiView = @{ isApproved = $false; packageNameApproved = $false; reason = 'packageNamePending' }
      }
    } | ConvertTo-Json -Compress

    { & $approvalScript -PackageInfoFiles @($packageInfoFile) -AzSdkExePath 'azsdk' } |
      Should -Throw '*not approved for release*'
  }

  It 'does not allow a stable version with only package-name approval' {
    $global:AzSdkResponse = @{
      operation_status = 'Succeeded'
      result = @{
        isApproved = $false
        finalSource = 'reviewHub'
        reason = 'packageNameApproved'
        apiView = @{ isApproved = $false; packageNameApproved = $true; reason = 'packageNameApproved' }
      }
    } | ConvertTo-Json -Compress

    { & $approvalScript -PackageInfoFiles @($packageInfoFile) -AzSdkExePath 'azsdk' } |
      Should -Throw '*not approved for release*'
  }

  It 'preserves the API approval exemption for management packages' {
    $packageInfo = Get-Content $packageInfoFile -Raw | ConvertFrom-Json
    $packageInfo.SdkType = 'mgmt'
    $packageInfo | ConvertTo-Json | Set-Content -Path $packageInfoFile
    $global:AzSdkResponse = @{
      operation_status = 'Succeeded'
      result = @{
        isApproved = $false
        finalSource = 'reviewHub'
        reason = 'missingApproval'
        reviewHub = @{ isPackageNameApproved = $false }
        apiView = @{ isApproved = $false; packageNameApproved = $false; reason = 'packageNamePending' }
      }
    } | ConvertTo-Json -Compress

    & $approvalScript -PackageInfoFiles @($packageInfoFile) -AzSdkExePath 'azsdk'
  }

  It 'surfaces operational CLI errors instead of treating them as approval denials' {
    $global:AzSdkResponse = @{
      operation_status = 'Failed'
      response_error = 'approval service unavailable'
      result = @{ isApproved = $false; finalSource = 'none'; reason = 'queryFailed' }
    } | ConvertTo-Json -Compress

    { & $approvalScript -PackageInfoFiles @($packageInfoFile) -AzSdkExePath 'azsdk' } |
      Should -Throw '*approval service unavailable*'
  }
}

Describe 'Go release pipeline' {
  It 'checks approval before creating a release tag' {
    $template = Get-Content (Join-Path $repoRoot 'eng/pipelines/templates/jobs/archetype-go-release.yml') -Raw
    $tagJob = $template.IndexOf('- job: TagRepository')
    $tagJob | Should -BeGreaterThan -1
    $approval = $template.IndexOf('displayName: Check Package Approval Status', $tagJob)
    $otherTriggers = $template.IndexOf('displayName: Check Package Approval Status for other release triggers', $tagJob)
    $tag = $template.IndexOf('/steps/create-tags-and-git-release.yml', $tagJob)

    $approval | Should -BeGreaterThan $tagJob
    $otherTriggers | Should -BeGreaterThan $tagJob
    $tag | Should -BeGreaterThan $approval
    $tag | Should -BeGreaterThan $otherTriggers
  }
}
