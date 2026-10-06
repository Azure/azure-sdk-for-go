# Verify the template-only Go release completion handoff without executing a release.
# Pester shares BeforeAll variables with It blocks, which the analyzer cannot resolve.
[Diagnostics.CodeAnalysis.SuppressMessageAttribute('PSUseDeclaredVarsMoreThanAssignments', '')]
[CmdletBinding()]
param()

BeforeAll {
  Set-StrictMode -Version 4
  $ErrorActionPreference = 'Stop'
  Import-Module powershell-yaml -ErrorAction Stop
  $root = (Resolve-Path (Join-Path $PSScriptRoot '../../..')).Path
  $entryPath = 'sdk/template/aztemplate/ci.yml'
  $clientPath = 'eng/pipelines/templates/jobs/archetype-sdk-client.yml'
  $releasePath = 'eng/pipelines/templates/jobs/archetype-go-release.yml'
  $entryRaw = [IO.File]::ReadAllText((Join-Path $root $entryPath)).Replace("`r`n", "`n")
  $clientRaw = [IO.File]::ReadAllText((Join-Path $root $clientPath)).Replace("`r`n", "`n")
  $releaseRaw = [IO.File]::ReadAllText((Join-Path $root $releasePath)).Replace("`r`n", "`n")
  $entry = ConvertFrom-Yaml $entryRaw
  $client = ConvertFrom-Yaml $clientRaw
  $release = ConvertFrom-Yaml $releaseRaw
  $releaseStage = @($release.stages | Where-Object { $_['stage'] -eq 'Release' })[0]
  $tagJob = @($releaseStage.jobs | Where-Object { $_['job'] -eq 'TagRepository' })[0]
  $completion = @($tagJob.steps | Where-Object { $_['template'] -eq '/eng/common/pipelines/templates/steps/mark-release-completion.yml' })[0]
  $templateGuard = '${{ if eq(parameters.ServiceDirectory, ''template/aztemplate'') }}'
  $autoGuard = '${{ if and(eq(variables[''Build.Reason''], ''IndividualCI''), eq(variables[''Build.SourceBranch''], ''refs/heads/main'')) }}'
  $bindingGuard = '${{ if and(eq(parameters.ServiceDirectory, ''template/aztemplate''), eq(variables[''Build.Reason''], ''IndividualCI''), eq(variables[''Build.SourceBranch''], ''refs/heads/main'')) }}'

  function Get-BaselineText([string]$Path) {
    # Pinned pre-pilot main: do not compare with a moving origin/main or include any synced changes.
    $lines = @(git -C $root show "b0d63acae6d73aadd7aa0722d96c69ad71f3dd75:$Path")
    if ($LASTEXITCODE -ne 0) { throw "Cannot read baseline $Path" }
    return ($lines -join "`n") + "`n"
  }
}

Describe 'Go template release completion contract' {
  It 'declares optional string IDs with a zero default at every relay' {
    foreach ($document in @($entry, $client, $release)) {
      $ids = @($document.parameters | Where-Object { $_.name -eq 'ReleasePlanId' })
      $ids.Count | Should -Be 1
      $ids[0].type | Should -BeExactly 'string'
      $ids[0].default | Should -BeOfType [string]
      $ids[0].default | Should -BeExactly '0'
    }
    $entry.extends.parameters.ReleasePlanId | Should -BeExactly '${{ parameters.ReleasePlanId }}'
    $entry.extends.parameters.ServiceDirectory | Should -BeExactly 'template/aztemplate'
    $clientRaw | Should -Match 'ServiceDirectory: \$\{\{ parameters.ServiceDirectory \}\}\n\s+ReleasePlanId: \$\{\{ parameters.ReleasePlanId \}\}'
    @($release.parameters | Where-Object name -eq 'DependsOn')[0].default | Should -BeExactly 'Build'
    @($release.parameters | Where-Object name -eq 'TestPipeline')[0].default | Should -BeFalse
    @($release.parameters | Where-Object name -eq 'ServiceDirectory')[0].default | Should -BeExactly ''
  }

  It 'keeps common completion inputs unchanged and scopes new arguments to the template' {
    $completion.parameters.ConfigFileDir | Should -BeExactly '$(Pipeline.Workspace)/PackageInfo'
    $completion.parameters.PackageArtifactName | Should -BeExactly 'sdk/${{ parameters.ServiceDirectory }}'
    @($completion.parameters.Keys).Count | Should -Be 3
    $pilot = $completion.parameters[$templateGuard]
    @($pilot.Keys).Count | Should -Be 2
    $pilot.ReleasePlanId | Should -BeExactly '${{ parameters.ReleasePlanId }}'
    $pilot[$autoGuard].SdkPullRequest | Should -BeExactly '$(AutoReleaseSdkPullRequestUrl)'
    @($pilot[$autoGuard].Keys).Count | Should -Be 1
  }

  It 'binds the fixed producer output only with the actual direct stage dependency' {
    $binding = @($releaseStage.variables | Where-Object { $_.ContainsKey($bindingGuard) })
    $binding.Count | Should -Be 1
    $binding[0][$bindingGuard][0].name | Should -BeExactly 'AutoReleaseSdkPullRequestUrl'
    $binding[0][$bindingGuard][0].value | Should -BeExactly '$[ stageDependencies.AutoReleasePrepare.ResolveAutoReleasePackages.outputs[''resolve.AutoReleaseSdkPullRequestUrl''] ]'
    $releaseStage[$autoGuard].dependsOn | Should -Contain 'AutoReleasePrepare'
    $releaseStage[$autoGuard].dependsOn | Should -Contain 'CheckRelease'
    $releaseStage['${{ else }}'].dependsOn | Should -BeExactly 'CheckRelease'
    $producer = @($release.stages | Where-Object { $_.ContainsKey($autoGuard) -and -not $_.ContainsKey('stage') })
    $producer.Count | Should -Be 1
    $producer[0][$autoGuard][0].template | Should -BeExactly '/eng/common/pipelines/templates/stages/archetype-auto-release-prepare.yml'
    $producer[0][$autoGuard][0].parameters.DependsOn | Should -BeExactly '${{ parameters.DependsOn }}'
    $tagJob.dependsOn | Should -BeExactly 'ReleaseGate'
    @($tagJob.steps | Where-Object { $_['download'] -eq 'current' } | ForEach-Object { $_.artifact }) | Should -Contain 'PackageInfo'
  }

  # These cases validate the exact parsed guards above, not a general Azure template evaluator.
  It 'selects the expected handoff for <Service>, <Reason>, <Branch>' -ForEach @(
    @{ Service = 'template/aztemplate'; Reason = 'Manual'; Branch = 'refs/heads/main'; Id = $true; Url = $false }
    @{ Service = 'template/aztemplate'; Reason = 'IndividualCI'; Branch = 'refs/heads/main'; Id = $true; Url = $true }
    @{ Service = 'template/aztemplate'; Reason = 'IndividualCI'; Branch = 'refs/heads/release/test'; Id = $true; Url = $false }
    @{ Service = 'template/aztemplate'; Reason = 'BatchedCI'; Branch = 'refs/heads/main'; Id = $true; Url = $false }
    @{ Service = 'template/aztemplate'; Reason = 'Schedule'; Branch = 'refs/heads/main'; Id = $true; Url = $false }
    @{ Service = 'template/aztemplate'; Reason = 'PullRequest'; Branch = 'refs/pull/1/merge'; Id = $true; Url = $false }
    @{ Service = 'azcore'; Reason = 'Manual'; Branch = 'refs/heads/main'; Id = $false; Url = $false }
    @{ Service = 'azcore'; Reason = 'IndividualCI'; Branch = 'refs/heads/main'; Id = $false; Url = $false }
    @{ Service = 'resourcemanager/compute/armcompute'; Reason = 'IndividualCI'; Branch = 'refs/heads/main'; Id = $false; Url = $false }
    @{ Service = 'template/other'; Reason = 'IndividualCI'; Branch = 'refs/heads/main'; Id = $false; Url = $false }
  ) {
    $handoff = @{
      ConfigFileDir = $completion.parameters.ConfigFileDir
      PackageArtifactName = $completion.parameters.PackageArtifactName
    }
    if ($Service -eq 'template/aztemplate') {
      $handoff.ReleasePlanId = $completion.parameters[$templateGuard].ReleasePlanId
      if ($Reason -eq 'IndividualCI' -and $Branch -eq 'refs/heads/main') {
        $handoff.SdkPullRequest = $completion.parameters[$templateGuard][$autoGuard].SdkPullRequest
      }
    }
    $handoff.ContainsKey('ReleasePlanId') | Should -Be $Id
    $handoff.ContainsKey('SdkPullRequest') | Should -Be $Url
    if (-not $Id) { @($handoff.Keys).Count | Should -Be 2 }
    # The client still excludes PR releases; this case does not assert a PR can release.
    $clientRaw | Should -Match "ne\(variables\['Build.Reason'\], 'PullRequest'\)"
  }

  It 'preserves the entire raw baseline after removing only the pilot additions' {
    $entryWithoutPilot = $entryRaw.Replace("  - name: ReleasePlanId`n    type: string`n    default: '0'`n", '').Replace('    ReleasePlanId: ${{ parameters.ReleasePlanId }}' + "`n", '')
    $entryWithoutPilot | Should -BeExactly (Get-BaselineText $entryPath)
    $clientWithoutPilot = $clientRaw.Replace("  - name: ReleasePlanId`n    type: string`n    default: '0'`n", '').Replace('                      ReleasePlanId: ${{ parameters.ReleasePlanId }}' + "`n", '')
    $clientWithoutPilot | Should -BeExactly (Get-BaselineText $clientPath)
    $releaseWithoutPilot = $releaseRaw -replace '(?s)\Aparameters:.*?\nstages:', "parameters:`n  DependsOn: Build`n  TestPipeline: false`n  ServiceDirectory: ''`n`nstages:"
    $releaseWithoutPilot = $releaseWithoutPilot -replace '(?m)^      # This output exists only.*\n      - \$\{\{ if and\(eq\(parameters.ServiceDirectory.*\n        - name: AutoReleaseSdkPullRequestUrl\n          value: .*\n', ''
    $releaseWithoutPilot = $releaseWithoutPilot -replace '(?m)^              # Template-only pilot.*\n              \$\{\{ if eq\(parameters.ServiceDirectory.*\n                ReleasePlanId: .*\n                \$\{\{ if and\(eq\(variables.*\n                  SdkPullRequest: .*\n', ''
    # Exact raw equality protects every gate, condition, environment, artifact and non-template step.
    $releaseWithoutPilot | Should -BeExactly (Get-BaselineText $releasePath)
  }

  It 'leaves the separate management auto-release entry byte-for-byte unchanged' {
    $path = 'eng/pipelines/mgmt-auto-release.yml'
    [IO.File]::ReadAllText((Join-Path $root $path)).Replace("`r`n", "`n") | Should -BeExactly (Get-BaselineText $path)
  }
}
