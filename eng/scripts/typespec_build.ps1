#Requires -Version 7.0
param(
    [string]$filter,
    [switch]$clean,
    [switch]$vet,
    [switch]$generate,
    [switch]$skipBuild,
    [switch]$cleanGenerated,
    [switch]$format,
    [switch]$tidy,
    [string]$emitterOptions = "module-version=0.1.0", # default options
    [string]$localSpecRepo
)

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

function Process-TypeSpec ()
{
    $currentDirectory = Get-Location
    
    # Check for tsp-location.yaml first
    if (-not (Test-Path "tsp-location.yaml")) {
        Write-Host "##[error]This is not a TypeSpec-based SDK. No tsp-location.yaml file found in $currentDirectory"
        return
    }

    if ($clean)
    {
        Write-Host "##[command]Executing go clean -v ./... in " $currentDirectory
        go clean -v ./...
        if ($LASTEXITCODE) { exit $LASTEXITCODE }
    }

    if ($cleanGenerated)
    {
        Write-Host "##[command]Cleaning auto-generated files in" $currentDirectory
        (Get-ChildItem -recurse "*.go" | Where-Object { $_.Name -notlike '*_test.go' } | Select-String -Pattern "Code generated by Microsoft" | Select-Object -ExpandProperty path) | Remove-Item -Force
    }

    if ($generate)
    {
        # Execute tsp-client update in current directory
        Write-Host "##[command]Executing tsp-client update in " $currentDirectory
        if ($localSpecRepo) {
            Write-Host "tsp-client update --debug --emitter-options"$emitterOptions "--local-spec-repo" $localSpecRepo
            tsp-client update --debug --emitter-options $emitterOptions --local-spec-repo $localSpecRepo
        } else {
            Write-Host "tsp-client update --debug --emitter-options"$emitterOptions
            tsp-client update --debug --emitter-options $emitterOptions
        }
        if ($LASTEXITCODE) {
            Write-Host "##[error]Error running tsp-client update"
            exit $LASTEXITCODE
        }
    }

    if ($format)
    {
        Write-Host "##[command]Executing gofmt -s -w . in " $currentDirectory
        gofmt -s -w .
        if ($LASTEXITCODE) { exit $LASTEXITCODE }
    }

    if ($tidy)
    {
        Write-Host "##[command]Executing go mod tidy in " $currentDirectory
        go mod tidy
        if ($LASTEXITCODE) { exit $LASTEXITCODE }
    }

    if (!$skipBuild)
    {
        Write-Host "##[command]Executing go build -v ./... in " $currentDirectory
        go build -v ./...
        Write-Host "##[command]Build Complete!"
        if ($LASTEXITCODE) { exit $LASTEXITCODE }
    }

    if ($vet)
    {
        Write-Host "##[command]Executing go vet ./... in " $currentDirectory
        go vet ./...
    }
}

try
{
    $startingDirectory = Get-Location

    $sdks = Get-AllPackageInfoFromRepo $filter

    foreach ($sdk in $sdks)
    {
        Push-Location $sdk.DirectoryPath
        Process-TypeSpec
        Pop-Location
    }

    if ($sdks.Count -eq 0 -and $filter -and (Test-Path -Path (Join-Path $RepoRoot "sdk" $filter)))
    {
        Write-Host "Cannot find go module under $filter, try to build in $(Join-Path $RepoRoot "sdk" $filter)"
        Push-Location (Join-Path $RepoRoot "sdk" $filter)
        Process-TypeSpec
        Pop-Location
    }
}
finally
{
    Set-Location $startingDirectory
}