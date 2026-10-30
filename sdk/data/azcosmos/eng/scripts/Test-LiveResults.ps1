# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

#Requires -Version 7.0
param(
    [Parameter(Mandatory)]
    [string] $EventsPath,
    [Parameter(Mandatory)]
    [string] $PlanPath,
    [Parameter(Mandatory)]
    [string] $OutputDirectory
)

$ErrorActionPreference = 'Stop'
$plan = @(Get-Content -Raw -LiteralPath $PlanPath | ConvertFrom-Json)
if ($plan.Count -eq 0 -or @($plan.ID | Sort-Object -Unique).Count -ne $plan.Count) {
    throw 'The live scenario manifest is empty or contains duplicate IDs.'
}
foreach ($entry in $plan) {
    if ($entry.ID -notmatch '^[a-z][a-z0-9.-]+$' -or $entry.Profile -notin @('smokeTests', 'coreOperations')) {
        throw 'The live scenario manifest contains an invalid ID or unsupported profile.'
    }
}
$events = @(Get-Content -LiteralPath $EventsPath | ForEach-Object { $_ | ConvertFrom-Json })
$suite = @($events | Where-Object { $_.Test -eq 'TestAzureLive' -and $_.Action -in @('pass', 'fail', 'skip') })
$package = @($events | Where-Object { !$_.Test -and $_.Action -in @('pass', 'fail', 'skip') })
$problems = [System.Collections.Generic.List[string]]::new()
if ($suite.Count -ne 1 -or $suite[0].Action -ne 'pass' -or
    $package.Count -ne 1 -or $package[0].Action -ne 'pass') {
    $problems.Add('The selected live suite/package did not complete successfully.')
}
foreach ($event in $events | Where-Object { $_.Test -and $_.Action -in @('fail', 'skip') }) {
    $problems.Add("Live test $($event.Test) ended as $($event.Action).")
}
foreach ($event in $events | Where-Object { $_.Action -eq 'run' -and $_.Test }) {
    $terminal = @($events | Where-Object { $_.Test -eq $event.Test -and $_.Action -in @('pass', 'fail', 'skip') })
    if ($terminal.Count -ne 1) {
        $problems.Add("Live test $($event.Test) has no unique terminal result.")
    }
}
$results = @(
    foreach ($entry in $plan) {
        $name = 'TestAzureLive/' + $entry.ID
        $result = @($events | Where-Object { $_.Test -eq $name -and $_.Action -in @('pass', 'fail', 'skip') })
        $passed = $result.Count -eq 1 -and $result[0].Action -eq 'pass'
        if (!$passed) {
            $problems.Add("Required selected scenario $($entry.ID) did not pass.")
        }
        [pscustomobject]@{
            scenario = $entry.ID
            profile = $entry.Profile
            backend = 'azureLive'
            outcome = $(if ($passed) { 'passed' } else { 'failed' })
            durationSeconds = $(if ($result.Count -eq 1) { $result[0].Elapsed } else { 0 })
        }
    }
)
New-Item -ItemType Directory -Path $OutputDirectory -Force | Out-Null
$results | ConvertTo-Json -Depth 4 | Set-Content -LiteralPath (Join-Path $OutputDirectory 'summary.json')
$settings = [System.Xml.XmlWriterSettings]::new()
$settings.Indent = $true
$writer = [System.Xml.XmlWriter]::Create((Join-Path $OutputDirectory 'report.xml'), $settings)
try {
    $writer.WriteStartElement('testsuite')
    $writer.WriteAttributeString('name', 'Cosmos Azure Live')
    $writer.WriteAttributeString('tests', [string]$results.Count)
    $writer.WriteAttributeString('failures', [string]@($results | Where-Object outcome -eq 'failed').Count)
    foreach ($result in $results) {
        $writer.WriteStartElement('testcase')
        $writer.WriteAttributeString('name', $result.scenario)
        $writer.WriteAttributeString('classname', 'azcosmos.azureLive.' + $result.profile)
        $writer.WriteAttributeString('time', [string]$result.durationSeconds)
        if ($result.outcome -ne 'passed') {
            $writer.WriteElementString('failure', 'Selected scenario was missing, skipped, or failed.')
        }
        $writer.WriteEndElement()
    }
    $writer.WriteEndElement()
}
finally {
    $writer.Dispose()
}
if ($problems.Count -gt 0) {
    throw ($problems -join [Environment]::NewLine)
}
Write-Host "All $($results.Count) selected Azure Live scenarios passed with no skips."
