#Requires -Version 7.0
param($filter, [switch]$vet, [switch]$generate, [switch]$skipBuild, $parallel = 5)

$sdks = @{};

foreach ($sdk in (./eng/scripts/get_module_dirs.ps1 -serviceDir 'sdk/...')) {
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

$keys | ForEach-Object { $sdks[$_] } | ForEach-Object -Parallel {
    Push-Location $_.path

    if (!$skipBuild) {
        Write-Host "##[command]Executing go build -v ./... in " $_.path
        go build -v ./...
    }
    if ($vet) {
        Write-Host "##[command]Executing go vet ./... in " $_.path
        go vet ./...
    }
    if ($generate) {
        Write-Host "##[command]Executing autorest.go in " $_.path
        # TODO
    }
} -ThrottleLimit $parallel
