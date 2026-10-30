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
$plan = Get-Content -Raw -LiteralPath $PlanPath | ConvertFrom-Json -NoEnumerate
if ($plan -isnot [array] -or $plan.Count -eq 0 -or @($plan.ID | Sort-Object -Unique).Count -ne $plan.Count) {
    throw 'The live scenario manifest must be a nonempty array with unique IDs.'
}
foreach ($entry in $plan) {
    if ($entry.ID -notmatch '^[a-z][a-z0-9.-]+$' -or $entry.Profile -notin @('smokeTests', 'coreOperations')) {
        throw 'The live scenario manifest contains an invalid ID or unsupported profile.'
    }
}
$events = @(Get-Content -LiteralPath $EventsPath | ForEach-Object { $_ | ConvertFrom-Json })
$suite = @($events | Where-Object { $_.Test -eq 'TestAzureLive' -and $_.Action -in @('pass', 'fail', 'skip') })
$suiteRuns = @($events | Where-Object { $_.Test -eq 'TestAzureLive' -and $_.Action -eq 'run' })
$package = @($events | Where-Object { !$_.Test -and $_.Action -in @('pass', 'fail', 'skip') })
$problems = [System.Collections.Generic.List[string]]::new()
if ($suiteRuns.Count -ne 1 -or $suite.Count -ne 1 -or $suite[0].Action -ne 'pass' -or
    $package.Count -ne 1 -or $package[0].Action -ne 'pass') {
    $problems.Add('The selected live suite/package did not complete successfully.')
}
foreach ($event in $events | Where-Object { $_.Test -and $_.Action -in @('fail', 'skip') }) {
    $problems.Add("Live test $($event.Test) ended as $($event.Action).")
}
foreach ($event in $events | Where-Object { $_.Action -eq 'run' -and $_.Test }) {
    $runs = @($events | Where-Object { $_.Test -eq $event.Test -and $_.Action -eq 'run' })
    $terminal = @($events | Where-Object { $_.Test -eq $event.Test -and $_.Action -in @('pass', 'fail', 'skip') })
    if ($runs.Count -ne 1 -or $terminal.Count -ne 1) {
        $problems.Add("Live test $($event.Test) has no unique run and terminal result.")
    }
}
$executionComplete = $problems.Count -eq 0
$results = @(
    foreach ($entry in $plan) {
        $name = 'TestAzureLive/' + $entry.ID
        $runs = @($events | Where-Object { $_.Test -eq $name -and $_.Action -eq 'run' })
        $result = @($events | Where-Object { $_.Test -eq $name -and $_.Action -in @('pass', 'fail', 'skip') })
        $passed = $executionComplete -and $runs.Count -eq 1 -and $result.Count -eq 1 -and $result[0].Action -eq 'pass'
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
ConvertTo-Json -InputObject $results -Depth 4 | Set-Content -LiteralPath (Join-Path $OutputDirectory 'summary.json')
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
        $writer.WriteAttributeString('time', ([double]$result.durationSeconds).ToString('R', [Globalization.CultureInfo]::InvariantCulture))
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
