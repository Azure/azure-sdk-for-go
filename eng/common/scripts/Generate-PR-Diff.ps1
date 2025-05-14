<#
.SYNOPSIS
Script used to generate the diff.json file for a PR. Explicitly intended to work in a PR context.

.DESCRIPTION
Combines the result of git diff, some parsed details from the diff, and the PR number into a single JSON file. This JSON file is intended for use further along the pipeline.

.PARAMETER ArtifactPath
The folder in which the result will be written.

.PARAMETER TargetPath
The path under which changes will be detected.
#>
[CmdletBinding()]
Param (
  [Parameter(Mandatory = $True)]
  [string] $ArtifactPath,
  [Parameter(Mandatory = $True)]
  [string] $TargetPath,
  [Parameter(Mandatory=$false)]
  [AllowEmptyCollection()]
  [array] $ExcludePaths
)

. (Join-Path $PSScriptRoot "Helpers" "git-helpers.ps1")

function Get-ChangedServices
{
  Param (
    [Parameter(Mandatory = $True)]
    [string[]] $ChangedFiles
  )

  [string[]] $changedServices = $ChangedFiles | Foreach-Object { if ($_ -match "sdk/([^/]+)") { $matches[1] } } | Sort-Object -Unique

  return , $changedServices
}

if (!(Test-Path $ArtifactPath))
{
  New-Item -ItemType Directory -Path $ArtifactPath | Out-Null
}

$ArtifactPath = Resolve-Path $ArtifactPath
$ArtifactName = Join-Path $ArtifactPath "diff.json"

$changedFiles = @()
$changedServices = @()

$changedFiles = @(
  "eng/scripts/Invoke-MgmtTestgen.ps1",
  "eng/scripts/MgmtTestLib.ps1",
  "eng/scripts/build.ps1",
  "eng/tools/generator/cmd/v2/common/cmdProcessor.go",
  "eng/tools/generator/cmd/v2/common/fileProcessor.go",
  "eng/tools/generator/cmd/v2/common/generation.go",
  "sdk/resourcemanager/advisor/armadvisor/build.go",
  "sdk/resourcemanager/agrifood/armagrifood/build.go",
  "sdk/resourcemanager/alertsmanagement/armalertsmanagement/build.go",
  "sdk/resourcemanager/applicationinsights/armapplicationinsights/build.go",
  "sdk/resourcemanager/appplatform/armappplatform/build.go",
  "sdk/resourcemanager/appservice/armappservice/build.go",
  "sdk/resourcemanager/authorization/armauthorization/build.go",
  "sdk/resourcemanager/automanage/armautomanage/build.go",
  "sdk/resourcemanager/azurestackhci/armazurestackhci/build.go",
  "sdk/resourcemanager/billing/armbilling/build.go",
  "sdk/resourcemanager/cdn/armcdn/build.go",
  "sdk/resourcemanager/changeanalysis/armchangeanalysis/build.go",
  "sdk/resourcemanager/confluent/armconfluent/build.go",
  "sdk/resourcemanager/connectedvmware/armconnectedvmware/build.go",
  "sdk/resourcemanager/consumption/armconsumption/build.go",
  "sdk/resourcemanager/customproviders/armcustomproviders/build.go",
  "sdk/resourcemanager/datacatalog/armdatacatalog/build.go",
  "sdk/resourcemanager/dataprotection/armdataprotection/build.go",
  "sdk/resourcemanager/datashare/armdatashare/build.go",
  "sdk/resourcemanager/devops/armdevops/build.go",
  "sdk/resourcemanager/devtestlabs/armdevtestlabs/build.go",
  "sdk/resourcemanager/domainservices/armdomainservices/build.go",
  "sdk/resourcemanager/eventgrid/armeventgrid/build.go",
  "sdk/resourcemanager/frontdoor/armfrontdoor/build.go",
  "sdk/resourcemanager/hybridcontainerservice/armhybridcontainerservice/build.go",
  "sdk/resourcemanager/managednetwork/armmanagednetwork/build.go",
  "sdk/resourcemanager/mobilenetwork/armmobilenetwork/build.go",
  "sdk/resourcemanager/msi/armmsi/build.go",
  "sdk/resourcemanager/mysql/armmysqlflexibleservers/build.go",
  "sdk/resourcemanager/network/armnetwork/build.go",
  "sdk/resourcemanager/networkfunction/armnetworkfunction/build.go",
  "sdk/resourcemanager/paloaltonetworksngfw/armpanngfw/build.go",
  "sdk/resourcemanager/policyinsights/armpolicyinsights/build.go",
  "sdk/resourcemanager/portal/armportal/build.go",
  "sdk/resourcemanager/postgresql/armpostgresqlflexibleservers/build.go",
  "sdk/resourcemanager/powerbiembedded/armpowerbiembedded/build.go",
  "sdk/resourcemanager/purview/armpurview/build.go",
  "sdk/resourcemanager/quota/armquota/build.go",
  "sdk/resourcemanager/resourcehealth/armresourcehealth/build.go",
  "sdk/resourcemanager/resourcemover/armresourcemover/build.go",
  "sdk/resourcemanager/resources/armfeatures/build.go",
  "sdk/resourcemanager/resources/armlocks/build.go",
  "sdk/resourcemanager/resources/armmanagedapplications/build.go",
  "sdk/resourcemanager/resources/armpolicy/build.go",
  "sdk/resourcemanager/saas/armsaas/build.go",
  "sdk/resourcemanager/scvmm/armscvmm/build.go",
  "sdk/resourcemanager/security/armsecurity/build.go",
  "sdk/resourcemanager/solutions/armmanagedapplications/build.go",
  "sdk/resourcemanager/storageimportexport/armstorageimportexport/build.go",
  "sdk/resourcemanager/storsimple1200series/armstorsimple1200series/build.go",
  "sdk/resourcemanager/support/armsupport/build.go",
  "sdk/resourcemanager/synapse/armsynapse/build.go"
)
$deletedFiles = Get-ChangedFiles -DiffPath $TargetPath -DiffFilterType "D"

if ($changedFiles) {
  $changedServices = Get-ChangedServices -ChangedFiles $changedFiles
}
else {
  # ensure we default this to an empty array if not set
  $changedFiles = @()
}

# ExcludePaths is an object array with the default of [] which evaluates to null.
# If the value is null, set it to empty list to ensure that the empty list is
# stored in the json
if (-not $ExcludePaths) {
  $ExcludePaths = @()
}
if (-not $deletedFiles) {
  $deletedFiles = @()
}
if (-not $changedServices) {
  $changedServices = @()
}


# $changedFiles = @(
#   "eng/scripts/Invoke-MgmtTestgen.ps1",
# "eng/scripts/MgmtTestLib.ps1",
# "eng/scripts/build.ps1",
# "eng/tools/generator/cmd/v2/common/cmdProcessor.go",
# "eng/tools/generator/cmd/v2/common/fileProcessor.go",
# "eng/tools/generator/cmd/v2/common/generation.go",
# "sdk/resourcemanager/advisor/armadvisor/build.go",
# "sdk/resourcemanager/agrifood/armagrifood/build.go",
# "sdk/resourcemanager/alertsmanagement/armalertsmanagement/build.go",
# "sdk/resourcemanager/applicationinsights/armapplicationinsights/build.go",
# "sdk/resourcemanager/appplatform/armappplatform/build.go",
# "sdk/resourcemanager/appservice/armappservice/build.go",
# "sdk/resourcemanager/authorization/armauthorization/build.go",
# "sdk/resourcemanager/automanage/armautomanage/build.go",
# "sdk/resourcemanager/azurestackhci/armazurestackhci/build.go",
# "sdk/resourcemanager/billing/armbilling/build.go",
# "sdk/resourcemanager/cdn/armcdn/build.go",
# "sdk/resourcemanager/changeanalysis/armchangeanalysis/build.go",
# "sdk/resourcemanager/confluent/armconfluent/build.go",
# "sdk/resourcemanager/connectedvmware/armconnectedvmware/build.go",
# "sdk/resourcemanager/consumption/armconsumption/build.go",
# "sdk/resourcemanager/customproviders/armcustomproviders/build.go",
# "sdk/resourcemanager/datacatalog/armdatacatalog/build.go",
# "sdk/resourcemanager/dataprotection/armdataprotection/build.go",
# "sdk/resourcemanager/datashare/armdatashare/build.go",
# "sdk/resourcemanager/devops/armdevops/build.go",
# "sdk/resourcemanager/devtestlabs/armdevtestlabs/build.go",
# "sdk/resourcemanager/domainservices/armdomainservices/build.go",
# "sdk/resourcemanager/eventgrid/armeventgrid/build.go",
# "sdk/resourcemanager/frontdoor/armfrontdoor/build.go",
# "sdk/resourcemanager/hybridcontainerservice/armhybridcontainerservice/build.go",
# "sdk/resourcemanager/managednetwork/armmanagednetwork/build.go",
# "sdk/resourcemanager/mobilenetwork/armmobilenetwork/build.go",
# "sdk/resourcemanager/msi/armmsi/build.go",
# "sdk/resourcemanager/mysql/armmysqlflexibleservers/build.go",
# "sdk/resourcemanager/network/armnetwork/build.go",
# "sdk/resourcemanager/networkfunction/armnetworkfunction/build.go",
# "sdk/resourcemanager/paloaltonetworksngfw/armpanngfw/build.go",
# "sdk/resourcemanager/policyinsights/armpolicyinsights/build.go",
# "sdk/resourcemanager/portal/armportal/build.go",
# "sdk/resourcemanager/postgresql/armpostgresqlflexibleservers/build.go",
# "sdk/resourcemanager/powerbiembedded/armpowerbiembedded/build.go",
# "sdk/resourcemanager/purview/armpurview/build.go",
# "sdk/resourcemanager/quota/armquota/build.go",
# "sdk/resourcemanager/resourcehealth/armresourcehealth/build.go",
# "sdk/resourcemanager/resourcemover/armresourcemover/build.go",
# "sdk/resourcemanager/resources/armfeatures/build.go",
# "sdk/resourcemanager/resources/armlocks/build.go",
# "sdk/resourcemanager/resources/armmanagedapplications/build.go",
# "sdk/resourcemanager/resources/armpolicy/build.go",
# "sdk/resourcemanager/saas/armsaas/build.go",
# "sdk/resourcemanager/scvmm/armscvmm/build.go",
# "sdk/resourcemanager/security/armsecurity/build.go",
# "sdk/resourcemanager/solutions/armmanagedapplications/build.go",
# "sdk/resourcemanager/storageimportexport/armstorageimportexport/build.go",
# "sdk/resourcemanager/storsimple1200series/armstorsimple1200series/build.go",
# "sdk/resourcemanager/support/armsupport/build.go",
# "sdk/resourcemanager/synapse/armsynapse/build.go"
# )
# $changedServices = @(
#   ""
# )
# $excludePaths = @()
# $deletedFiles = @()


$result = [PSCustomObject]@{
  "ChangedFiles"    = $changedFiles
  "ChangedServices" = $changedServices
  "ExcludePaths"    = $ExcludePaths
  "DeletedFiles"    = $deletedFiles
  "PRNumber"        = if ($env:SYSTEM_PULLREQUEST_PULLREQUESTNUMBER) { $env:SYSTEM_PULLREQUEST_PULLREQUESTNUMBER } else { "-1" }
}

$json = $result | ConvertTo-Json
$json | Out-File $ArtifactName

Write-Host "`nGenerated diff.json file at $ArtifactName"
Write-Host "  $($json -replace "`n", "`n  ")"
