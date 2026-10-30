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
    foreach ($case in @('passed', 'skipped', 'failed', 'missing', 'empty plan', 'duplicate plan', 'unterminated', 'package failed', 'suite missing')) {
        $plan = @(@{ ID = 'item.lifecycle'; Profile = 'smokeTests' })
        if ($case -eq 'empty plan') { $plan = @() }
        if ($case -eq 'duplicate plan') { $plan = @($plan[0], $plan[0]) }
        ConvertTo-Json -InputObject $plan | Set-Content -LiteralPath $planPath
        $events = @(
            @{ Action = 'run'; Test = 'TestAzureLive' }
            @{ Action = 'run'; Test = 'TestAzureLive/item.lifecycle' }
        )
        if ($case -ne 'missing' -and $case -ne 'unterminated') {
            $events += @{
                Action = $(switch ($case) { 'skipped' { 'skip' } 'failed' { 'fail' } default { 'pass' } })
                Test = 'TestAzureLive/item.lifecycle'
                Elapsed = 0.1
            }
        }
        if ($case -ne 'suite missing') {
            $events += @{ Action = 'pass'; Test = 'TestAzureLive' }
        }
        if ($case -eq 'unterminated') {
            $events += @{ Action = 'run'; Test = 'TestAzureLive/item.lifecycle/nested' }
        }
        $events += @{ Action = $(if ($case -eq 'package failed') { 'fail' } else { 'pass' }) }
        $events | ForEach-Object { ConvertTo-Json -InputObject $_ -Compress } | Set-Content -LiteralPath $eventsPath
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
                $report.testsuite.testcase.name -ne 'item.lifecycle') {
                throw 'Passing JUnit report lost the scenario identity.'
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
