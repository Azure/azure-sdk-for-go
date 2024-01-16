Param(
    [string]$searchPath = $PWD.Path
)

$modFiles = Get-ChildItem -Path $searchPath -Include "go.mod" -Recurse

$tidyErrors = $modFiles | ForEach-Object -Parallel {
    Set-Location $_.Directory
    Write-Host (Get-Location)
    $output = go mod tidy 2>&1
    if ($LASTEXITCODE) {
        return @{ Directory = $_.Directory; Output = $output }
    }
}

if ($tidyErrors) {
    Write-Error "Encountered the following tidy failures:" -ErrorAction 'Continue'
    foreach ($err in $tidyErrors) {
        Write-Host "=== $($err.Directory) ==="
        $err.Output
    }
    exit 1
}
