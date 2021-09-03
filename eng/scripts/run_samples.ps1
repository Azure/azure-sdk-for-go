#Requires -Version 7.0

Param(
    [string] $serviceDirectory
)

Write-Host "Running samples for $serviceDirectory"

if (Test-Path -Path "sdk/samples/$serviceDirectory") {
    Set-Location -Path "sdk/samples/$serviceDirectory"

    Get-ChildItem -Recurse -Path . -filter go.mod | ForEach-Object {
        $cwd = $_.Directory
        Write-Host "Running samples in $cwd"

        go run .
        if ($LASTEXITCODE) {
            Write-Host "There was an error running samples."
            exit $LASTEXITCODE
        }
    }
} else {
    Write-Host "Could not find samples in sdk/samples/$serviceDirectory"
}