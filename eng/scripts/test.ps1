#Requires -Version 7.0

Param(
    [string] $serviceDirectory
)

$repoRoot = Resolve-Path "$PSScriptRoot/../../"

# Load the environment variables
$configFile = Join-Path $repoRoot "eng" "config.json"
$configData = Get-Content -Raw -Path "$configFile" | ConvertFrom-Json

$configData.PSObject.Properties | ForEach-Object {
    Write-Host $_.Name
    Write-Host $_.Value.EnvironmentVariables
    Write-Host $_.Value.EnvironmentVariables.GetType()

    if ($_.Name -eq "EnvironmentVariables") {
        $_.Value | ForEach-Object {
            Write-Host $_
        }
    }
}