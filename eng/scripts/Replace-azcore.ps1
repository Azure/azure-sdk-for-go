Param(
    [string] $ModuleDirectory
)

$goModFile = Join-Path $ModuleDirectory "go.mod"

if (!(Test-Path $goModFile)) {
    Write-Host "##[command]The file $goModFile doesn't exist"
    exit 1
}

if (!$ModuleDirectory.Contains("$([IO.Path]::DirectorySeparatorChar)sdk$([IO.Path]::DirectorySeparatorChar)")) {
    Write-Host "##[command]Directory $ModuleDirectory doesn't appear to be an SDK module"
    exit 1
}

if ((Get-Content $goModFile -raw) -notmatch "github.com/Azure/azure-sdk-for-go/sdk/azcore") {
    # no dependency on azcore, exit
    Write-Host "##[command]No azcore dependency found in " $goModFile
    return
}

# walk up the directory tree until we find the sdk directory, constructing the relative path as we go
$relativePath = ""
for ($parent = $ModuleDirectory; !$parent.EndsWith("$([IO.Path]::DirectorySeparatorChar)sdk"); $parent = (Split-Path $parent)) {
    $relativePath += "../"
}

# add a replace directive to go.mod with a relative path to the azcore directory
$replace = "replace github.com/Azure/azure-sdk-for-go/sdk/azcore => $($relativePath)azcore"
Write-Host "##[command]Adding replace statement " $replace
Add-Content -Path $goModFile -Value "`n$($replace)"

## go mod tidy
Write-Host "##[command]Executing go mod tidy in " $ModuleDirectory
Set-Location $ModuleDirectory
go mod tidy
if ($LASTEXITCODE) { exit $LASTEXITCODE }
