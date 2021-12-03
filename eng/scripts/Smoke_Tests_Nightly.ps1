#Requires -Version 7.0

Param(
    [string] $serviceDirectory
)

$repoRoot = Resolve-Path "$PSScriptRoot/../../"

Push-Location $repoRoot/eng/tools/smoketests

# create a smoketests directory
$smoketestsDir = Join-Path $repoRoot sdk smoketests
Write-Host "Creating a new directory for smoketests at $smoketestsDir"
New-Item -Path $smoketestsDir -ItemType Directory

Push-Location $smoketestsDir
Write-Host "Running 'go mod init' in $pwd"
go mod init github.com/Azure/azure-sdk-for-go/sdk/smoketests
Pop-Location

# Run smoketests script
Write-Host "Running 'go run . -serviceDirectory $serviceDirectory'"
go run . -serviceDirectory $serviceDirectory
if ($LASTEXITCODE) {
    exit $LASTEXITCODE
}

Pop-Location

# Run go mod tidy and go build. If these succeed the smoke tests pass
Push-Location $smoketestsDir
go fmt ./...
Write-Host "Printing content of go.mod file:"
Get-Content go.mod
Write-Host "Printing content of main.go file"
Get-Content main.go
go mod tidy
go build ./...
go run .

Pop-Location

# Clean-up the directory created
Remove-Item -Path $smoketestsDir -Recurse -Force
