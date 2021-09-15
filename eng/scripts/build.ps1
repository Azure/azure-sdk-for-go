#Requires -Version 7.0
param([string]$filter, [switch]$clean, [switch]$vet, [switch]$generate, [switch]$skipBuild, [switch]$cleanGenerated, [switch]$format, [switch]$tidy, [string]$config = "autorest.md", [string]$outputFolder)

. $PSScriptRoot/meta_generation.ps1
. $PSScriptRoot/get_module_dirs.ps1


function Process-Sdk () {
    $currentDirectory = Get-Location
    if ($clean) {
        Write-Host "##[command]Executing go clean -v ./... in " $currentDirectory
        go clean -v ./...
        if ($LASTEXITCODE) { exit $LASTEXITCODE }
    }

    if ($cleanGenerated) {
        Write-Host "##[command]Cleaning auto-generated files in" $currentDirectory
        Remove-Item "zz_generated_*"
    }

    if ($generate) {
        Write-Host "##[command]Executing autorest.go in " $currentDirectory
        $autorestPath = "./" + $config

        if (ShouldGenerate-AutorestConfig $autorestPath) {
            Generate-AutorestConfig $autorestPath
            $removeAutorestFile = $true
        }

        $autorestVersion = "@autorest/go@4.0.0-preview.27"
        if ($outputFolder -eq '') {
            $outputFolder = $currentDirectory
        }
        autorest --use=$autorestVersion --go --track2 --go-sdk-folder=$root --output-folder=$outputFolder --file-prefix="zz_generated_" --clear-output-folder=false $autorestPath
        if ($LASTEXITCODE) {
            Write-Host "##[error]Error running autorest.go"
            exit $LASTEXITCODE
        }
        if ($removeAutorestFile) {
            Remove-Item $autorestPath
        }
    }

    if ($format) {
        Write-Host "##[command]Executing gofmt -s -w . in " $currentDirectory
        gofmt -s -w .
        if ($LASTEXITCODE) {exit $LASTEXITCODE}
    }

    if ($tidy) {
        Write-Host "##[command]Executing go mod tidy in " $currentDirectory
        go mod tidy
        if ($LASTEXITCODE) {exit $LASTEXITCODE}
    }

    if (!$skipBuild) {
        Write-Host "##[command]Executing go build -x -v ./... in " $currentDirectory
        go build -x -v ./...
        Write-Host "##[command]Build Complete!"
        if ($LASTEXITCODE) {exit $LASTEXITCODE}
    }

    if ($vet) {
        Write-Host "##[command]Executing go vet ./... in " $currentDirectory
        go vet ./...
    }

}

$startingDirectory = Get-Location
$root = Resolve-Path ($PSScriptRoot + "/../..")
Set-Location $root
$sdks = @{};

foreach ($sdk in (Get-ModuleDirs 'sdk/')) {
    $name = $sdk | split-path -leaf
    $sdks[$name] = @{
        'path' = $sdk;
    }
}

$keys = $sdks.Keys | Sort-Object;
if (![string]::IsNullOrWhiteSpace($filter)) {
    Write-Host "Using filter: $filter"
    $keys = $keys.Where( { $_ -match $filter })
}

try {
    $keys | ForEach-Object { $sdks[$_] } | ForEach-Object {
        Push-Location $_.path
        Process-Sdk
        Pop-Location
    }
}
finally {
    Set-Location $startingDirectory
}
