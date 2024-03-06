#Requires -Version 7.0

Param(
    [string] $serviceDirectory
)

Write-Host $PSScriptRoot

# 1. Every module uses a replace directive to the local version
# 2. Include every module (data & mgmt) in a go.mod file
# 3. Run `go mod tidy` and ensure it succeeds