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

if (-not (Test-Path Env:TF_BUILD)) {
    $StrictMode = $true

    Write-Host "installing"
    go install github.com/Azure/azure-sdk-for-go/eng/tools/snippet-generator@snippet-generator
} else {
    go install $root/../eng/tools/snippet-generator
}

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
