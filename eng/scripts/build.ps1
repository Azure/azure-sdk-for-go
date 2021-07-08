function Get-ParentDir() {
    [OutputType([string])]
    param ([string]$path)
    return Split-Path -Path $path -Parent
}

function Get-FilteredChildren {
    param ([string]$Path, [scriptblock]$Filter)    
    return Get-ChildItem -Path $Path -Recurse -Directory | Where-Object -FilterScript $Filter
}

$sdkRoot = Get-ParentDir (Get-ParentDir ($PSScriptRoot))
$skippedDirs = @(
    [String](Join-Path -Path $sdkRoot -ChildPath "vendor"),
    [String](Join-Path -Path $sdkRoot -ChildPath "sdk"),
    [String](Join-Path -Path $sdkRoot -ChildPath "tools\generator")
)

function Test-SkipDir () {
    [OutputType([boolean])]
    param ([string]$path)
    return ($skippedDirs | Where-Object { $path.StartsWith($_) }).Count -eq 0
}

$filteredLocations = Get-FilteredChildren -Path $sdkRoot -Filter { Test-SkipDir $_.FullName };
# foreach ($dir in $filteredLocations) {    
#     Write-Host $dir
# }
$process = Start-Process -FilePath "go" -ArgumentList @("list", $filteredLocations[0]) -PassThru -Wait
Write-Host $filteredLocations[0]
Write-Host $process.StandardError
Write-Host $process.ExitCode