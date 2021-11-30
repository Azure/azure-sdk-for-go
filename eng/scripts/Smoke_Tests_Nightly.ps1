#Requires -Version 7.0

$repoRoot = Resolve-Path "$PSScriptRoot/../../"

Push-Location $repoRoot/eng/scripts/smoketest/

go run .

Pop-Location

Push-Location $repoRoot/sdk/smoketests

go mod tidy
go build ./...

Pop-Location

# Clean-up the directory created
Remove-Item -Path $repoRoot/sdk/smoketests -Recurse -Force