#Requires -Version 7.0
param(
    # filter = service filter
    [string]$filter,
    [switch]$clean,
    [switch]$vet,
    [switch]$generate,
    [switch]$skipBuild,
    [switch]$cleanGenerated,
    [switch]$format,
    [switch]$tidy,
    [switch]$alwaysSetBodyParamRequired,
    [switch]$removeUnreferencedTypes,
    [switch]$factoryGatherAllParams,
    [string]$config = "autorest.md",
    [string]$goExtension = "@autorest/go@4.0.0-preview.71",
    [string]$filePrefix,
    [string]$outputFolder
)

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

function Process-Sdk ()
{
    $currentDirectory = Get-Location

    # Check for autorest.md first
    if (-not (Test-Path "autorest.md")) {
        LogWarning "This is not a Swagger-based SDK. No autorest.md file found in $currentDirectory"
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
        Write-Host "##[command]Executing autorest.go in " $currentDirectory
        $autorestPath = "./" + $config

        if ($outputFolder -eq '')
        {
            $outputFolder = $currentDirectory
        }

        $honorBodyPlacement = "false"
        if (-not $alwaysSetBodyParamRequired)
        {
            $honorBodyPlacement = "true"
        }

        $removeUnreferencedTypesFlag = "false"
        if ($removeUnreferencedTypes)
        {
            $removeUnreferencedTypesFlag = "true"
        }

        $factoryGatherAllParamsFlag = "false"
        if ($factoryGatherAllParams)
        {
            $factoryGatherAllParamsFlag = "true"
        }

        if ($filePrefix)
        {
            Write-Host "autorest --use=$goExtension --go --track2 --output-folder=$outputFolder --file-prefix=$filePrefix --clear-output-folder=false --go.clear-output-folder=false --honor-body-placement=$honorBodyPlacement --remove-unreferenced-types=$removeUnreferencedTypesFlag --factory-gather-all-params=$factoryGatherAllParamsFlag $autorestPath"
            autorest --use=$goExtension --go --track2 --output-folder=$outputFolder --file-prefix=$filePrefix --clear-output-folder=false --go.clear-output-folder=false --honor-body-placement=$honorBodyPlacement --remove-unreferenced-types=$removeUnreferencedTypesFlag --factory-gather-all-params=$factoryGatherAllParamsFlag $autorestPath
        }
        else
        {
            Write-Host "autorest --use=$goExtension --go --track2 --output-folder=$outputFolder --clear-output-folder=false --go.clear-output-folder=false --honor-body-placement=$honorBodyPlacement --remove-unreferenced-types=$removeUnreferencedTypesFlag --factory-gather-all-params=$factoryGatherAllParamsFlag $autorestPath"
            autorest --use=$goExtension --go --track2 --output-folder=$outputFolder --clear-output-folder=false --go.clear-output-folder=false --honor-body-placement=$honorBodyPlacement --remove-unreferenced-types=$removeUnreferencedTypesFlag --factory-gather-all-params=$factoryGatherAllParamsFlag $autorestPath
        }

        if ($LASTEXITCODE)
        {
            Write-Host "##[error]Error running autorest.go"
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
    $sdks = @()

    $sdks = Get-AllPackageInfoFromRepo $filter

    Write-Host "Successfully found $($sdks.Count) go modules to build."
    foreach ($sdk in $sdks)
    {
        Push-Location $sdk.DirectoryPath
        Process-Sdk
        Pop-Location
    }

    if ($sdks.Count -eq 0 -and $filter -and (Test-Path -Path (Join-Path $RepoRoot "sdk" $filter)))
    {
        Write-Host "Cannot find go module under $filter, try to build in $(Join-Path $RepoRoot "sdk" $filter)"
        Push-Location (Join-Path $RepoRoot "sdk" $filter)
        Process-Sdk
        Pop-Location
    }
}
finally
{
    Set-Location $startingDirectory
}
