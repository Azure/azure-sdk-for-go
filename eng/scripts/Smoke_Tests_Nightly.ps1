#Requires -Version 7.0

$repoRoot = Resolve-Path "$PSScriptRoot/../../"

Push-Location $repoRoot/eng/tools/smoketests

# create a smoketests directory
$smoketestsDir = Join-Path $repoRoot sdk smoketests
Write-Host "Creating a new directory for smoketests at $smoketestsDir"
New-Item -Path $smoketestsDir -ItemType Directory

Push-Location $smoketestsDir
go mod init
Pop-Location

# Run smoketests script
go run .

Pop-Location

# Run go mod tidy and go build. If these succeed the smoke tests pass
Push-Location $smoketestsDir
go mod tidy
go build ./...

Pop-Location

# Clean-up the directory created
Remove-Item -Path $smoketestsDir -Recurse -Force