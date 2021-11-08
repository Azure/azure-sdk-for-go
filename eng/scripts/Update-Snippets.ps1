[CmdletBinding()]
[CmdletBinding()]
param (
    [Parameter(Position=0)]
    [ValidateNotNullOrEmpty()]
    [string] $ServiceDirectory,

    [Parameter()]
    [switch] $StrictMode
)

$root = "$PSScriptRoot/../../sdk"

if ($ServiceDirectory -and ($ServiceDirectory -ne "*")) {
    $root += '/' + $ServiceDirectory
}

$snippetGeneratorPath = Resolve-Path "$PSScriptRoot/../tools/snippet-generator"

Push-Location $snippetGeneratorPath
go install
Pop-Location

Write-Host "Updating snippets in $root"
if ($StrictMode) {
    Resolve-Path "$root" | ForEach-Object {
        & snippet-generator $_ true 
    }
} else {
    Resolve-Path "$root" | ForEach-Object {
        & snippet-generator $_ false
    }
}
