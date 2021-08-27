$Language = "go"
$packagePattern = "go.mod"
# rewrite from ChangeLog-Operations.ps1 used in Get-ChangeLogEntriesFromContent for go uses vx.x.x as version number
$RELEASE_TITLE_REGEX = "(?<releaseNoteTitle>^\#+\s+(?<version>v$([AzureEngSemanticVersion]::SEMVER_REGEX))(\s+(?<releaseStatus>\(.+\))))"
# constants for go version fetch
$GO_VERSION_REGEX = ".+\s*=\s*`".*v?(?<version>$([AzureEngSemanticVersion]::SEMVER_REGEX))`""
$VERSION_FILE_SUFFIXS = "*constants.go", "*version.go"

# rewrite from artifact-metadata-parsing.ps1 used in RetrievePackages for fetch go single module info
function Get-go-PackageInfoFromPackageFile ($pkg, $workingDirectory)
{
    $workFolder = $pkg.Directory
    $releaseNotes = ""
    $namespaceName = $workFolder.Name
    $rpName = $workFolder.Parent.Name

    $pkgId = "sdk/$rpName/$namespaceName"
    $pkgVersion = Get-Version $workFolder

    $changeLogLoc = @(Get-ChildItem -Path $workFolder -Recurse -Include "CHANGELOG.md")[0]
    if ($changeLogLoc)
    {
        $releaseNotes = Get-ChangeLogEntryAsString -ChangeLogLocation $changeLogLoc -VersionString v$pkgVersion
    }

    $resultObj = New-Object PSObject -Property @{
        PackageId      = $pkgId
        PackageVersion = $pkgVersion
        ReleaseTag     = "$pkgId/v$pkgVersion"
        Deployable     = $true
        ReleaseNotes   = $releaseNotes
    }

    return $resultObj
}

# get version from specific files (*constants.go, *version.go)
function Get-Version ($pkgPath)
{
    # find any file with surfix
    $versionFiles = [Collections.Generic.List[String]]@()
    foreach ($versionFileSuffix in $VERSION_FILE_SUFFIXS)
    {
        Get-ChildItem -Recurse -Path $pkgPath -Filter $versionFileSuffix | ForEach-Object {
            Write-Host "Adding $_ to list of version files"
            $versionFiles.Add($_)
        }
    }
    
    # for each version file, use regex to search go version num
    foreach ($versionFile in $versionFiles)
    {
        $content = Get-Content $versionFile
        $content = $content.Split("`n")
        try
        {
            # walk the document, finding where the version number are
            foreach ($line in $content)
            {
                if ($line -match $GO_VERSION_REGEX)
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
    }

    Write-Host "Cannot find release version."
    exit 1
}

# rewrite from artifact-metadata-parsing.ps1 as all 0.x.x version should be considered prerelease in go sdk
function CreateReleases($pkgList, $releaseApiUrl, $releaseSha)
{
    foreach ($pkgInfo in $pkgList)
    {
        Write-Host "Creating release $($pkgInfo.Tag)"
  
        $releaseNotes = ""
        if ($null -ne $pkgInfo.ReleaseNotes)
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