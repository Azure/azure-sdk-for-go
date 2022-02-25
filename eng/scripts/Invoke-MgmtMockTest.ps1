Param(
    [string] $filter
)

. (Join-Path $PSScriptRoot .. common scripts common.ps1)
. (Join-Path $PSScriptRoot MgmtTestLib.ps1)

$env:TEMP = [System.IO.Path]::GetTempPath()
Write-Host "Path tmp: $env:TEMP"

$sdks = Get-AllPackageInfoFromRepo $filter

Write-Host "Prepare mock server"
if ($sdks.Count -eq 0)
{
    Write-Host "No package need to be test"
    exit 0
}
else
{
    PrepareMockServer
    Write-Host "Try Stop mock server"
    StopMockServer
}

foreach ($sdk in $sdks)
{
    if ($sdk.SdkType -eq "mgmt")
    {
        try
        {
            ExecuteSingleTest $sdk
        }
        catch
        {
            Write-Host "##[error]can not finish single test for $sdks :`n$_"
            exit 1
        }
    }
}

Write-Host "Try Stop mock server"
StopMockServer