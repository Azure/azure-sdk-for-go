Param(
    [string] $ServiceDirectory
)

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

$moduleDirectory = Join-Path $RepoRoot "sdk" $ServiceDirectory
$goModFile = Join-Path $moduleDirectory "go.mod"

if (!(Test-Path $goModFile)) {
    Write-Host "##[command]The file $goModFile doesn't exist"
    exit 1
}

if ((Get-Content $goModFile -raw) -notmatch "github.com/Azure/azure-sdk-for-go/sdk/azcore") {
    # no dependency on azcore, exit
    Write-Host "##[command]No azcore dependency found in " $goModFile
    return
}

# walk up the directory tree until we find the sdk directory, constructing the relative path as we go
$relativePath = ""
for ($parent = $moduleDirectory; !$parent.EndsWith("$([IO.Path]::DirectorySeparatorChar)sdk"); $parent = (Split-Path $parent)) {
    if ($parent -eq $RepoRoot) {
        # we hit the root of the repo, bail to prevent infinite loop
        Write-Host "##[command]Walked to repo root which is unexpected"
        exit 1
    }
    $relativePath += "../"
}

# add a replace directive to go.mod with a relative path to the azcore directory
$replace = "replace github.com/Azure/azure-sdk-for-go/sdk/azcore => $($relativePath)azcore"
Write-Host "##[command]Adding replace statement " $replace
Add-Content -Path $goModFile -Value "`n$($replace)"

## go mod tidy
Write-Host "##[command]Executing go mod tidy in " $moduleDirectory
Set-Location $moduleDirectory
go mod tidy
if ($LASTEXITCODE) { exit $LASTEXITCODE }
