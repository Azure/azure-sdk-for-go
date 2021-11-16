#Requires -Version 7.0
param(
    [string]$filter,
    [switch]$clean,
    [switch]$vet,
    [switch]$generateExample,
    [switch]$generateMockTest,
    [switch]$skipBuild,
    [switch]$cleanGenerated,
    [switch]$format,
    [switch]$tidy,
    [string]$config = "autorest.md",
    [string]$autorestVersion = "3.6.2",
    [string]$goExtension = "@autorest/go@4.0.0-preview.31",
    [string]$testExtension = "@autorest/gotest@1.1.2",
    [string]$outputFolder
)

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

function Invoke-MgmtTestgen ()
{
    $currentDirectory = Get-Location
    if ($clean)
    {
        Write-Host "##[command]Executing go clean -v ./... in " $currentDirectory
        go clean -v ./...
        if ($LASTEXITCODE) { exit $LASTEXITCODE }
    }

    if ($cleanGenerated)
    {
        Write-Host "##[command]Cleaning auto-generated files in" $currentDirectory
        Remove-Item "ze_generated_*"
        Remove-Item "zt_generated_*"
    }

    if ($generateExample -or $generateMockTest)
    {
        Write-Host "##[command]Executing autorest.gotest in " $currentDirectory
        $autorestPath = "./" + $config
        
        if ($outputFolder -eq '')
        {
            $outputFolder = $currentDirectory
        }
        $exampleFlag = "false"
        if ($generateExample)
        {
            $exampleFlag = "true"
        }
        $mockTestFlag = "true"
        if (-not $generateMockTest)
        {
            $mockTestFlag = "false"
        }
        Write-Host "autorest --version=$autorestVersion --use=$goExtension --use=$testExtension --go --track2 --output-folder=$outputFolder --clear-output-folder=false --go.clear-output-folder=false --generate-sdk=false --testmodeler.generate-mock-test=$mockTestFlag --testmodeler.generate-sdk-example=$exampleFlag $autorestPath"
        autorest --version=$autorestVersion --use=$goExtension --use=$testExtension --go --track2 --output-folder=$outputFolder --clear-output-folder=false --go.clear-output-folder=false --generate-sdk=false --testmodeler.generate-mock-test=$mockTestFlag --testmodeler.generate-sdk-example=$exampleFlag $autorestPath
        if ($LASTEXITCODE)
        {
            Write-Host "##[error]Error running autorest.gotest"
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
        Write-Host "##[command]Executing go build -x -v ./... in " $currentDirectory
        go build -x -v ./...
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
        Invoke-MgmtTestgen
        Pop-Location
    }
}
finally
{
    Set-Location $startingDirectory
}
