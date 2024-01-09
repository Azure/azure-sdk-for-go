Param(
    [string]$searchPath = $PWD.Path
)

$modFiles = Get-ChildItem -Path $searchPath -Include "go.mod" -Recurse

$modFiles | ForEach-Object -Parallel {
    Push-Location $_.Directory
    Write-Host (Get-Location)
    go mod tidy
    Pop-Location
    if ($LASTEXITCODE) {
        exit 1
    }
}
