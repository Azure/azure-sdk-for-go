#Requires -Version 7.0
param($filter, [switch]$clean, [switch]$vet, [switch]$generate, [switch]$skipBuild)

$startingDirectory = Get-Location
$root = Resolve-Path ($PSScriptRoot + "/../..")
Set-Location $root
$sdks = @{};

foreach ($sdk in (./eng/scripts/get_module_dirs.ps1 -serviceDir 'sdk/...')) {
    $name = $sdk | split-path -leaf
    $sdks[$name] = @{
        'path'      = $sdk;
        'clean'     = $clean;
        'vet'       = $vet;
        'generate'  = $generate;
        'skipBuild' = $skipBuild;
        'root'      = $root;
    }
}

$keys = $sdks.Keys | Sort-Object;
if (![string]::IsNullOrWhiteSpace($filter)) { 
    Write-Host "Using filter: $filter"
    $keys = $keys.Where( { $_ -match $filter }) 
}

$keys | ForEach-Object { $sdks[$_] } | ForEach-Object {
    Push-Location $_.path

    if ($_.clean) {
        Write-Host "##[command]Executing go clean -v ./... in " $_.path
        go clean -v ./...
    }

    if ($_.generate) {
        Write-Host "##[command]Executing autorest.go in " $_.path
        $autorestPath = $_.path + "\autorest.md"
        $autorestVersion = "@autorest/go@4.0.0-preview.23"
        $outputFolder = $_.path
        $root = $_.root
        autorest --use=$autorestVersion --go --track2 --go-sdk-folder=$root --output-folder=$outputFolder --file-prefix="zz_generated_" --clear-output-folder=false $autorestPath
    }
    if (!$_.skipBuild) {
        Write-Host "##[command]Executing go build -v ./... in " $_.path
        go build -x -v ./...
        Write-Host "##[command]Build Complete!"

    }
    if ($_.vet) {
        Write-Host "##[command]Executing go vet ./... in " $_.path
        go vet ./...
    }
    Pop-Location
}

Set-Location $startingDirectory
