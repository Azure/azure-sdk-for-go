Param(
    [Parameter(Mandatory=$True)]
    [string] $ServiceDirectory,
    [Parameter(Mandatory=$True)]
    [string] $OutPath,
    [Parameter(Mandatory=$True)]
    [string] $ApiviewUri,
    [Parameter(Mandatory=$True)]
    [string] $ApiKey,
    [Parameter(Mandatory=$True)]
    [string] $ApiLabel,
    [Parameter(Mandatory=$True)]
    [string] $SourceBranch,
    [Parameter(Mandatory=$True)]
    [string] $DefaultBranch,
    [Parameter(Mandatory=$True)]
    [string] $ConfigFileDir,
    [string] $RepoName,
    [string] $BuildId,
    [bool] $MarkPackageAsShipped = $false
)


Write-Host "$PSScriptRoot"
. (Join-Path $PSScriptRoot .. common scripts common.ps1)
$createReviewScript = (Join-Path $PSScriptRoot .. common scripts Create-APIReview.ps1)

$artifactList = @()
foreach ($sdk in (Get-AllPackageInfoFromRepo $ServiceDirectory))
{
    Write-Host "Creating API review artifact for $($sdk.Name)"
    New-Item -ItemType Directory -Path $OutPath/$($sdk.Name) -force
    $fileName = Split-Path -Path $sdk.Name -Leaf
    Compress-Archive -Path $sdk.DirectoryPath -DestinationPath $outPath/$($sdk.Name)/$fileName -force
    Rename-Item $outPath/$($sdk.Name)/$fileName.zip -NewName "$fileName.gosource"

    $artifactList += [PSCustomObject]@{
        name = $sdk.Name
    }
}

&($createReviewScript) -ArtifactList $artifactList  -ArtifactPath $outPath -APIKey $ApiKey -SourceBranch $SourceBranch -DefaultBranch $DefaultBranch -ConfigFileDir $ConfigFileDir -RepoName $RepoName -BuildId $BuildId -MarkPackageAsShipped $MarkPackageAsShipped
```
