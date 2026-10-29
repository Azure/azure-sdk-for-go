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

  function Find-PipelineNode {
    param($Node, [string]$PropertyName, [string]$PropertyValue)
    if ($Node -is [System.Collections.IDictionary]) {
      if ($Node[$PropertyName] -eq $PropertyValue) { $Node }
      foreach ($child in $Node.Values) { Find-PipelineNode -Node $child -PropertyName $PropertyName -PropertyValue $PropertyValue }
    }
    elseif ($Node -is [System.Collections.IEnumerable] -and $Node -isnot [string]) {
      foreach ($child in $Node) { Find-PipelineNode -Node $child -PropertyName $PropertyName -PropertyValue $PropertyValue }
    }
  }
  $releaseCalls = @(Find-PipelineNode -Node $client -PropertyName 'template' -PropertyValue 'archetype-go-release.yml')
  $management = ConvertFrom-Yaml ([IO.File]::ReadAllText((Join-Path $root 'eng/pipelines/mgmt-auto-release.yml')))
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
    $releaseCalls.Count | Should -Be 1
    $releaseCalls[0].parameters.ServiceDirectory | Should -BeExactly '${{ parameters.ServiceDirectory }}'
    $releaseCalls[0].parameters.ReleasePlanId | Should -BeExactly '${{ parameters.ReleasePlanId }}'
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

  # Check current parsed safety contracts without requiring history in a shallow CI checkout.
  It 'keeps template triggers and client release eligibility opt-in' {
    $entry.extends.template | Should -BeExactly '/eng/pipelines/templates/jobs/archetype-sdk-client.yml'
    @($entry.trigger.branches.include) -join '|' | Should -BeExactly 'main|hotfix/*|release/*'
    @($entry.pr.branches.include) -join '|' | Should -BeExactly 'main|feature/*|hotfix/*|release/*'
    foreach ($trigger in @($entry.trigger, $entry.pr)) {
      @($trigger.paths.include) -join '|' | Should -BeExactly 'sdk/template/aztemplate/|eng/'
    }
    @($entry.extends.parameters.TriggeringPaths) -join '|' | Should -BeExactly '/eng/'
    foreach ($name in @('SkipPrValidation', 'IncludeRelease')) {
      $parameter = @($client.parameters | Where-Object name -eq $name)
      $parameter.Count | Should -Be 1
      $parameter[0].type | Should -BeExactly 'boolean'
      $parameter[0].default | Should -BeFalse
    }
    $eligibility = '${{ if and(not(and(eq(parameters.SkipPrValidation, true), eq(variables[''Build.Reason''], ''Manual''))), ne(variables[''Build.Reason''], ''PullRequest''), eq(variables[''System.TeamProject''], ''internal'')) }}'
    $notWeekly = '${{ if not(contains(variables[''Build.DefinitionName''], ''weekly'')) }}'
    $includeRelease = '${{ if or(in(variables[''Build.Reason''], ''Manual'', '''', ''IndividualCI''), eq(parameters.IncludeRelease, true)) }}'
    $guarded = @($client.extends.parameters.stages | Where-Object { $_.ContainsKey($eligibility) })
    $guarded.Count | Should -Be 1
    $calls = @(Find-PipelineNode -Node $guarded[0][$eligibility][0][$notWeekly][0][$includeRelease] -PropertyName 'template' -PropertyValue 'archetype-go-release.yml')
    $calls.Count | Should -Be 1
    $calls[0].parameters.ContainsKey('TestPipeline') | Should -BeFalse
    $calls[0].parameters.DependsOn[0] | Should -BeExactly 'Build'
    @($calls[0].parameters.DependsOn).Count | Should -Be 2
    $live = '${{if and(eq(variables[''System.TeamProject''], ''internal''), eq(parameters.RunLiveTests, ''true''))}}'
    $clouds = $calls[0].parameters.DependsOn[1][$live][0]['${{ each cloud in parameters.CloudConfig }}']
    $selected = '${{ if or(contains(parameters.Clouds, cloud.key), and(contains(variables[''Build.DefinitionName''], ''tests-weekly''), contains(parameters.SupportedClouds, cloud.key))) }}'
    $clouds[0][$selected][0]['${{ if not(contains(parameters.UnsupportedClouds, cloud.key)) }}'][0] | Should -BeExactly '${{ cloud.key }}'
    $build = @(Find-PipelineNode -Node $client -PropertyName 'stage' -PropertyValue 'Build')
    $build.Count | Should -Be 1
    $build[0].condition | Should -BeExactly 'not(and(eq(${{ parameters.SkipPrValidation }}, true), eq(variables[''Build.Reason''], ''Manual'')))'
  }

  It 'keeps release verification, auto-release artifact eligibility and approval gates' {
    $check = @(Find-PipelineNode -Node $release -PropertyName 'stage' -PropertyValue 'CheckRelease')
    $check.Count | Should -Be 1
    $check[0].dependsOn | Should -BeExactly '${{ parameters.DependsOn }}'
    $check[0].condition | Should -BeExactly 'and(succeeded(), ne(variables[''SetDevVersion''], ''true''), ne(variables[''Skip.Release''], ''true''), ne(variables[''Build.Repository.Name''], ''Azure/azure-sdk-for-go-pr''), ne(variables.UseAzcoreFromMain, ''true''))'
    $verify = @(Find-PipelineNode -Node $check[0] -PropertyName 'name' -PropertyValue 'Verify')
    $verify.Count | Should -Be 1
    $verify[0].task | Should -BeExactly 'PowerShell@2'
    $verify[0].inputs.filePath | Should -BeExactly './eng/scripts/Verify-NeedToRelease.ps1'
    $verify[0].env.GH_TOKEN | Should -BeExactly '$(GH_TOKEN)'
    @($releaseStage[$autoGuard].dependsOn) -join '|' | Should -BeExactly 'CheckRelease|AutoReleasePrepare'
    ($releaseStage[$autoGuard].condition -replace '\s', '') | Should -BeExactly 'and(succeeded(),eq(dependencies.CheckRelease.outputs[''CheckReleaseJob.Verify.NeedToRelease''],''true''),eq(dependencies.AutoReleasePrepare.outputs[''ResolveAutoReleasePackages.resolve.HasAutoReleaseArtifacts''],''true''))'
    ($releaseStage['${{ else }}'].condition -replace '\s', '') | Should -BeExactly 'and(succeeded(),eq(dependencies.CheckRelease.outputs[''CheckReleaseJob.Verify.NeedToRelease''],''true''))'
    $producer = @($release.stages | Where-Object { $_.ContainsKey($autoGuard) -and -not $_.ContainsKey('stage') })[0][$autoGuard][0]
    @($producer.parameters.Artifacts).Count | Should -Be 2
    $outsideSdk = '${{ if startsWith(parameters.ServiceDirectory, ''../'') }}'
    $producer.parameters.Artifacts[0][$outsideSdk].name | Should -BeExactly '${{ replace(parameters.ServiceDirectory, ''../'', '''') }}'
    $producer.parameters.Artifacts[0][$outsideSdk].safeName | Should -BeExactly '${{ replace(replace(parameters.ServiceDirectory, ''../'', ''''), ''/'', ''_'') }}'
    $producer.parameters.Artifacts[1]['${{ else }}'].name | Should -BeExactly 'sdk/${{ parameters.ServiceDirectory }}'
    $producer.parameters.Artifacts[1]['${{ else }}'].safeName | Should -BeExactly '${{ replace(parameters.ServiceDirectory, ''/'', ''_'') }}'
    $gates = @(Find-PipelineNode -Node $releaseStage -PropertyName 'deployment' -PropertyValue 'ReleaseGate')
    $gates.Count | Should -Be 1
    $gate = $gates[0]
    $gate.ContainsKey('environment') | Should -BeFalse
    $gate[$autoGuard].environment | Should -BeExactly 'none'
    $gate['${{ else }}'].environment | Should -BeExactly 'package-publish'
    $gate.templateContext.type | Should -BeExactly 'releaseJob'
    $gate.templateContext.isProduction | Should -BeTrue
    $gate.pool.name | Should -BeExactly 'azsdk-pool'
  }

  It 'keeps Go publication through verified release tags before completion and version updates' {
    # Go publishes source tags, not signed package binaries; preserve that publication contract.
    $tagJob.condition | Should -BeExactly 'and(succeeded(), ne(variables[''Skip.TagRepository''], ''true''))'
    @($tagJob.steps | Where-Object { $_['checkout'] -eq 'self' }).Count | Should -Be 1
    @($tagJob.steps | Where-Object { $_['download'] -eq 'current' } | ForEach-Object { $_.artifact }) -join '|' | Should -BeExactly 'PackageInfo|packages'
    $changelog = @(Find-PipelineNode -Node $tagJob -PropertyName 'template' -PropertyValue '/eng/common/pipelines/templates/steps/verify-changelog.yml')
    $changelog.Count | Should -Be 1
    $changelog[0].parameters.ForRelease | Should -BeTrue
    $changelog[0].parameters.ServiceDirectory | Should -BeExactly '${{parameters.ServiceDirectory}}'
    $spec = @(Find-PipelineNode -Node $tagJob -PropertyName 'template' -PropertyValue '/eng/common/pipelines/templates/steps/verify-restapi-spec-location.yml')
    $spec.Count | Should -Be 1
    $spec[0].parameters.ArtifactLocation | Should -BeExactly '$(Pipeline.Workspace)'
    foreach ($file in @('./eng/scripts/validate_go_mod.ps1', './eng/scripts/validate_go_mod_major_version.ps1')) {
      $validation = @($tagJob.steps | Where-Object { $_['task'] -eq 'PowerShell@2' -and $_.inputs.filePath -eq $file })
      $validation.Count | Should -Be 1
      $validation[0].inputs.arguments | Should -BeExactly '${{ parameters.ServiceDirectory }}'
      $validation[0].inputs.pwsh | Should -BeTrue
    }
    $publishTemplate = '/eng/common/pipelines/templates/steps/create-tags-and-git-release.yml'
    $publication = @(Find-PipelineNode -Node $tagJob -PropertyName 'template' -PropertyValue $publishTemplate)
    $publication.Count | Should -Be 1
    $publication[0].parameters.ArtifactLocation | Should -BeExactly '$(Build.SourcesDirectory)/sdk/${{ parameters.ServiceDirectory }}'
    $publication[0].parameters.PackageFilter | Should -BeExactly 'sdk/${{ parameters.ServiceDirectory }}'
    $publication[0].parameters.ReleaseSha | Should -BeExactly '$(Build.SourceVersion)'
    $publication[0].parameters.RepoId | Should -BeExactly 'Azure/azure-sdk-for-go'
    $publication[0].parameters.WorkingDirectory | Should -BeExactly '$(System.DefaultWorkingDirectory)'
    $publication[0].parameters.AuthToken | Should -BeExactly ''
    $templates = @($tagJob.steps | Where-Object { $_.ContainsKey('template') } | ForEach-Object { $_.template })
    $templates -join '|' | Should -BeExactly '/eng/common/pipelines/templates/steps/retain-run.yml|/eng/common/pipelines/templates/steps/verify-changelog.yml|/eng/common/pipelines/templates/steps/verify-restapi-spec-location.yml|/eng/common/pipelines/templates/steps/create-tags-and-git-release.yml|/eng/common/pipelines/templates/steps/set-default-branch.yml|/eng/common/pipelines/templates/steps/mark-release-completion.yml'
    $review = @(Find-PipelineNode -Node $tagJob -PropertyName 'task' -PropertyValue 'AzureCLI@2')
    $review.Count | Should -Be 1
    $review[0].condition | Should -BeExactly 'and(succeeded(), ne(variables[''Skip.CreateApiReview''], ''true''))'
    $review[0].inputs.azureSubscription | Should -BeExactly 'APIView prod deployment'
    $review[0].inputs.inlineScript | Should -Match '-MarkPackageAsShipped \$true'
    @($tagJob.steps)[-1].template | Should -BeExactly '/eng/common/pipelines/templates/steps/mark-release-completion.yml'
    $versionGuard = '${{ if not(and(startsWith(parameters.ServiceDirectory, ''resourcemanager''), ne(parameters.ServiceDirectory, ''resourcemanager/internal''))) }}'
    $guarded = @($releaseStage.jobs | Where-Object { $_.ContainsKey($versionGuard) })
    $guarded.Count | Should -Be 1
    $updates = @(Find-PipelineNode -Node $guarded[0][$versionGuard] -PropertyName 'job' -PropertyValue 'UpdatePackageVersion')
    $updates.Count | Should -Be 1
    $updates[0].dependsOn | Should -BeExactly 'TagRepository'
    $updates[0].condition | Should -BeExactly 'and(succeeded(), ne(variables[''Skip.UpdatePackageVersion''], ''true''))'
    $pullRequests = @(Find-PipelineNode -Node $updates[0] -PropertyName 'template' -PropertyValue '/eng/common/pipelines/templates/steps/create-pull-request.yml')
    $pullRequests.Count | Should -Be 1
    $pullRequests[0].parameters.CloseAfterOpenForTesting | Should -BeExactly '${{ parameters.TestPipeline }}'
  }

  It 'keeps the separate management auto-release entry independent of pilot completion' {
    $management.trigger | Should -BeExactly 'none'
    $management.pr | Should -BeExactly 'none'
    $management.extends.template | Should -BeExactly '/eng/pipelines/templates/stages/1es-redirect.yml'
    $management.ContainsKey('parameters') | Should -BeFalse
    @(Find-PipelineNode -Node $management -PropertyName 'template' -PropertyValue 'archetype-go-release.yml').Count | Should -Be 0
    @(Find-PipelineNode -Node $management -PropertyName 'template' -PropertyValue '/eng/common/pipelines/templates/steps/mark-release-completion.yml').Count | Should -Be 0
    @(Find-PipelineNode -Node $management -PropertyName 'name' -PropertyValue 'AutoReleaseSdkPullRequestUrl').Count | Should -Be 0
    $shell = @(Find-PipelineNode -Node $management -PropertyName 'task' -PropertyValue 'ShellScript@2')
    $shell.Count | Should -Be 1
    $shell[0].inputs.scriptPath | Should -BeExactly 'eng/scripts/mgmt-auto-release.sh'
    $shell[0].inputs.cwd | Should -BeExactly '$(System.DefaultWorkingDirectory)'
    $shell[0].env.GH_TOKEN | Should -BeExactly '$(GH_TOKEN)'
  }
}
