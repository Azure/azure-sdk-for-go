BeforeAll {
  $repoRoot = (Resolve-Path (Join-Path $PSScriptRoot .. .. ..)).Path
  $approvalScript = Join-Path $repoRoot 'eng/common/scripts/Get-PackageApprovalStatus.ps1'

  function global:azsdk {
    param([Parameter(ValueFromRemainingArguments = $true)] [object[]] $Arguments)
    $global:LASTEXITCODE = 0
    if ($Arguments.Count -eq 1 -and $Arguments[0] -eq '--version') {
      return '0.6.38'
    }
    $global:AzSdkArguments = @($Arguments)
    return $global:AzSdkResponse
  }
}

AfterAll {
  Remove-Item Function:\azsdk
  Remove-Variable AzSdkArguments, AzSdkResponse, LanguageShort -Scope Global
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
}

Describe 'Go release pipeline' {
  It 'checks approval before creating a release tag' {
    $template = Get-Content (Join-Path $repoRoot 'eng/pipelines/templates/jobs/archetype-go-release.yml') -Raw
    $tagJob = $template.IndexOf('- job: TagRepository')
    $tagJob | Should -BeGreaterThan -1
    $approval = $template.IndexOf('/steps/get-package-approval-status.yml', $tagJob)
    $otherTriggers = $template.IndexOf('displayName: Check Package Approval Status for other release triggers', $tagJob)
    $tag = $template.IndexOf('/steps/create-tags-and-git-release.yml', $tagJob)

    $approval | Should -BeGreaterThan $tagJob
    $otherTriggers | Should -BeGreaterThan $tagJob
    $tag | Should -BeGreaterThan $approval
    $tag | Should -BeGreaterThan $otherTriggers
  }
}
