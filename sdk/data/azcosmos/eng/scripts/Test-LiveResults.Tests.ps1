# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

#Requires -Version 7.0
$ErrorActionPreference = 'Stop'
$directory = Join-Path ([System.IO.Path]::GetTempPath()) ('cosmos-live-gate-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $directory | Out-Null
$eventsPath = Join-Path $directory 'events.jsonl'
$planPath = Join-Path $directory 'selection.json'
$reportPath = Join-Path $directory 'report.xml'
$summaryPath = Join-Path $directory 'summary.json'
try {
    foreach ($case in @('passed', 'skipped', 'failed', 'missing', 'empty plan', 'duplicate plan',
        'object plan', 'invalid ID', 'unsupported profile', 'unterminated', 'package failed',
        'package missing', 'package skipped', 'suite missing', 'suite skipped', 'suite not run',
        'scenario not run', 'duplicate run', 'duplicate result', 'nested skipped', 'malformed events')) {
        foreach ($path in @($reportPath, $summaryPath)) {
            if (Test-Path -LiteralPath $path) { Remove-Item -LiteralPath $path }
        }
        $plan = @(@{ ID = 'item.lifecycle'; Profile = 'smokeTests' })
        if ($case -eq 'empty plan') { $plan = @() }
        if ($case -eq 'duplicate plan') { $plan = @($plan[0], $plan[0]) }
        if ($case -eq 'object plan') { $plan = $plan[0] }
        if ($case -eq 'invalid ID') { $plan[0].ID = 'not a scenario ID' }
        if ($case -eq 'unsupported profile') { $plan[0].Profile = 'dynamicTopology' }
        ConvertTo-Json -InputObject $plan | Set-Content -LiteralPath $planPath
        $events = @(
            @{ Action = 'run'; Test = 'TestAzureLive' }
            @{ Action = 'run'; Test = 'TestAzureLive/item.lifecycle' }
        )
        if ($case -eq 'suite not run') { $events = @($events[1]) }
        if ($case -eq 'scenario not run') { $events = @($events[0]) }
        if ($case -eq 'duplicate run') { $events += $events[1] }
        if ($case -ne 'missing' -and $case -ne 'unterminated') {
            $events += @{
                Action = $(switch ($case) { 'skipped' { 'skip' } 'failed' { 'fail' } default { 'pass' } })
                Test = 'TestAzureLive/item.lifecycle'
                Elapsed = 0.1
            }
        }
        if ($case -eq 'duplicate result') {
            $events += @{ Action = 'pass'; Test = 'TestAzureLive/item.lifecycle'; Elapsed = 0.1 }
        }
        if ($case -ne 'suite missing') {
            $events += @{ Action = $(if ($case -eq 'suite skipped') { 'skip' } else { 'pass' }); Test = 'TestAzureLive' }
        }
        if ($case -eq 'unterminated') {
            $events += @{ Action = 'run'; Test = 'TestAzureLive/item.lifecycle/nested' }
        }
        if ($case -eq 'nested skipped') {
            $events += @{ Action = 'run'; Test = 'TestAzureLive/item.lifecycle/nested' }
            $events += @{ Action = 'skip'; Test = 'TestAzureLive/item.lifecycle/nested' }
        }
        if ($case -ne 'package missing') {
            $events += @{ Action = $(switch ($case) { 'package failed' { 'fail' } 'package skipped' { 'skip' } default { 'pass' } }) }
        }
        $events | ForEach-Object { ConvertTo-Json -InputObject $_ -Compress } | Set-Content -LiteralPath $eventsPath
        if ($case -eq 'malformed events') { Add-Content -LiteralPath $eventsPath -Value '{' }
        $rejected = $false
        try {
            & "$PSScriptRoot/Test-LiveResults.ps1" -EventsPath $eventsPath -PlanPath $planPath -OutputDirectory $directory
        }
        catch {
            $rejected = $true
        }
        if ($rejected -eq ($case -eq 'passed')) {
            throw "Unexpected gate outcome for fixture '$case'."
        }
        if ($case -eq 'passed') {
            [xml] $report = Get-Content -Raw -LiteralPath $reportPath
            if ($report.testsuite.tests -ne '1' -or $report.testsuite.failures -ne '0' -or
                $report.testsuite.testcase.name -ne 'item.lifecycle' -or $report.testsuite.testcase.time -ne '0.1') {
                throw 'Passing JUnit report lost the scenario identity.'
            }
            $summary = Get-Content -Raw -LiteralPath $summaryPath | ConvertFrom-Json -NoEnumerate
            if ($summary -isnot [array] -or $summary.Count -ne 1 -or
                $summary[0].scenario -ne 'item.lifecycle' -or $summary[0].profile -ne 'smokeTests' -or
                $summary[0].backend -ne 'azureLive' -or $summary[0].outcome -ne 'passed') {
                throw 'Partial selection must remain an array of only the selected scenario.'
            }
        }
        elseif (Test-Path -LiteralPath $reportPath) {
            [xml] $report = Get-Content -Raw -LiteralPath $reportPath
            if ([int]$report.testsuite.failures -lt 1) {
                throw "Rejected fixture '$case' must not publish a passing JUnit report."
            }
        }
        Write-Host "PASS: $case"
    }
}
finally {
    foreach ($path in @($eventsPath, $planPath, $reportPath, $summaryPath)) {
        if (Test-Path -LiteralPath $path) { Remove-Item -LiteralPath $path }
    }
    Remove-Item -LiteralPath $directory
}
