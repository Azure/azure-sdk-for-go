$Language = "go"
$packagePattern = "go.mod"
$RELEASE_TITLE_REGEX = "(?<releaseNoteTitle>^\#+\s+(?<version>v$([AzureEngSemanticVersion]::SEMVER_REGEX))(\s+(?<releaseStatus>\(.+\))))"
$AUTOREST_VERSION_REGEX = "^module-version:\s+(?<version>$([AzureEngSemanticVersion]::SEMVER_REGEX))"

function Get-go-PackageInfoFromPackageFile ($pkg, $workingDirectory)
{
    $workFolder = $($pkg.Directory)
    $releaseNotes = ""
    $namespaceName = $workFolder.Name
    $rpName = $workFolder.Parent.Name

    $pkgId = "sdk/$rpName/$namespaceName"
    $pkgVersion = Get-Version($workFolder)

    $changeLogLoc = @(Get-ChildItem -Path $workFolder -Recurse -Include "CHANGELOG.md")[0]
    if ($changeLogLoc)
    {
        $releaseNotes = Get-ChangeLogEntryAsString -ChangeLogLocation $changeLogLoc -VersionString v$pkgVersion
    }

    $resultObj = New-Object PSObject -Property @{
        PackageId      = $pkgId
        PackageVersion = $pkgVersion
        ReleaseTag     = "$($pkgId)/v$($pkgVersion)"
        Deployable     = $true
        ReleaseNotes   = $releaseNotes
    }

    return $resultObj
}

# get version from go config auterest.md
function Get-Version ($pkgPath)
{
    $content = Get-Content(Join-Path $pkgPath "autorest.md")
    $content = $content.Split("`n")

    try
    {
        # walk the document, finding where the version specifiers are
        foreach ($line in $content)
        {
            if ($line -match $AUTOREST_VERSION_REGEX)
            {
                return $matches["version"]
            }
        }
    }
    catch
    {
        Write-Error "Error parsing version."
        Write-Error $_
    }
    Write-Host "Cannot find release version."
    exit(1)
}

# rewrite for go as all 0.x.x version should be considered prerelease
function CreateReleases($pkgList, $releaseApiUrl, $releaseSha)
{
    foreach ($pkgInfo in $pkgList)
    {
        Write-Host "Creating release $($pkgInfo.Tag)"
  
        $releaseNotes = ""
        if ($pkgInfo.ReleaseNotes -ne $null)
        {
            $releaseNotes = $pkgInfo.ReleaseNotes
        }
  
        $isPrerelease = $False
  
        $parsedSemver = [AzureEngSemanticVersion]::ParseVersionString($pkgInfo.PackageVersion)
  
        if ($parsedSemver -and ($parsedSemver.IsPrerelease -or ($parsedSemver.Major -eq 0)))
        {
            $isPrerelease = $true
        }
  
        $url = $releaseApiUrl
        $body = ConvertTo-Json @{
            tag_name         = $pkgInfo.Tag
            target_commitish = $releaseSha
            name             = $pkgInfo.Tag
            draft            = $False
            prerelease       = $isPrerelease
            body             = $releaseNotes
        }
  
        $headers = @{
            "Content-Type"  = "application/json"
            "Authorization" = "token $env:GH_TOKEN"
        }

        Invoke-RestMethod -Uri $url -Body $body -Headers $headers -Method "Post" -MaximumRetryCount 3 -RetryIntervalSec 10
    }
}