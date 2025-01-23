Write-Host "$PSScriptRoot"
. (Join-Path $PSScriptRoot .. common scripts common.ps1)

<#
.DESCRIPTION
    Create .gosource APIVIew artifact for go
.PARAMETER ServiceDirectory
    Thee name of the ServiceDirectory
.PARAMETER ArtifactOutputDirectory
    Base output Directory path for the generated gosource artifacts
.PARAMETER ArtifactName
    Name of the group of artifacts to be created
#>
function Create-APIViewArtifact {
        Param(
        [Parameter(Mandatory=$True)]
        [string] $ServiceDirectory,
        [Parameter(Mandatory=$True)]
        [string] $ArtifactOutputDirectory,
        [Parameter(Mandatory=$True)]
        [string] $ArtifactName
    )

    $artifactList = @()

    $artifactsDirectoryPath = Join-Path -Path $ArtifactOutputDirectory $ArtifactName
    New-Item -ItemType Directory -Path $artifactsDirectoryPath -force

    foreach ($sdk in (Get-AllPackageInfoFromRepo $ServiceDirectory))
    {
        Write-Host "Creating API review artifact for $($sdk.Name)"
        New-Item -ItemType Directory -Path $artifactsDirectoryPath/$($sdk.Name) -force
        $fileName = Split-Path -Path $sdk.Name -Leaf
        Compress-Archive -Path $sdk.DirectoryPath -DestinationPath $artifactsDirectoryPath/$($sdk.Name)/$fileName -force
        Rename-Item $artifactsDirectoryPath/$($sdk.Name)/$fileName.zip -NewName "$fileName.gosource"

        $artifactList += [PSCustomObject]@{
            name = $sdk.Name
        }
    }
    return $artifactList
}

<#
.DESCRIPTION
    Create new automatic APIView from a CI run
.PARAMETER ServiceDirectory
    Thee name of the ServiceDirectory
.PARAMETER ArtifactOutputDirectory
    Base output Directory path for the generated gosource artifacts
.PARAMETER ArtifactName
    Name of the group of artifacts to be created
.PARAMETER ApiKey
    The APIview ApiKey
.PARAMETER SourceBranch
    SourceBranch
.PARAMETER DefaultBranch
    DefaultBranch
.PARAMETER ConfigFileDir
    Path to the ConfigFileDir as published in the pipeline
.PARAMETER RepoName
    The name of the repository
.PARAMETER BuildId
    The build Id of the pipeline run
.PARAMETER MarkPackageAsShipped
    Indicate weather to mark the package a s shipped
#>
function New-APIView-From-CI {
    Param(
        [Parameter(Mandatory=$True)]
        [string] $ServiceDirectory,
        [Parameter(Mandatory=$True)]
        [string] $ArtifactOutputDirectory,
        [string] $ArtifactName="APIViewArtifacts",
        [Parameter(Mandatory=$True)]
        [string] $ApiKey,
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
    $artifactList = Create-APIViewArtifact -ServiceDirectory $ServiceDirectory -ArtifactOutputDirectory $ArtifactOutputDirectory -ArtifactName $ArtifactName
    $createReviewScript = (Join-Path $PSScriptRoot .. common scripts Create-APIReview.ps1)

    Write-Host "Create Go APIView using generated artifacts"
    Write-Host "ArtifactList: $artifactList"
    Write-Host "ApiKey: $ApiKey"
    Write-Host "SourceBranch: $SourceBranch"
    Write-Host "DefaultBranch: $DefaultBranch"
    Write-Host "ConfigFileDir: $ConfigFileDir"
    Write-Host "RepoName: $RepoName"
    Write-Host "BuildId: $BuildId"
    Write-Host "MarkPackageAsShipped: $MarkPackageAsShipped"
    Write-Host "createReviewScript: $createReviewScript"


    #&($createReviewScript) -ArtifactList $artifactList -ArtifactPath $outPath -APIKey $ApiKey -SourceBranch $SourceBranch -DefaultBranch $DefaultBranch -ConfigFileDir $ConfigFileDir -RepoName $RepoName -BuildId $BuildId -MarkPackageAsShipped $MarkPackageAsShipped
}

<#
.DESCRIPTION
    Call APIView endpoint to detect changes in API surface and create APIView if required
.PARAMETER ArtifactPath
    Path to the generated artifact
.PARAMETER CommitSha
    Commist sha for the build
.PARAMETER ArtifactName
    Name of the group of artifacts to be created
.PARAMETER ApiKey
    The APIview ApiKey
.PARAMETER SourceBranch
    SourceBranch
.PARAMETER DefaultBranch
    DefaultBranch
.PARAMETER ConfigFileDir
    Path to the ConfigFileDir as published in the pipeline
.PARAMETER RepoName
    The name of the repository
.PARAMETER BuildId
    The build Id of the pipeline run
.PARAMETER MarkPackageAsShipped
    Indicate weather to mark the package a s shipped
#>
function New-APIView-From-PR {
    Param(
        [Parameter(Mandatory=$True)]
        [string] $ArtifactPath,
        [Parameter(Mandatory=$True)]
        [string] $CommitSha,
        [Parameter(Mandatory=$True)]
        [string] $BuildId,
        [Parameter(Mandatory=$True)]
        [string] $PullRequestNumber,
        [Parameter(Mandatory=$True)]
        [string] $RepoFullName,
        [Parameter(Mandatory=$True)]
        [string] $APIViewUri,
        [Parameter(Mandatory=$True)]
        [string] $ArtifactName,
        [Parameter(Mandatory=$True)]
        [string] $DevopsProject
    )

    $detectApiChanges = (Join-Path $PSScriptRoot .. common scripts Detect-Api-Changes.ps1)

    Write-Host "Detect API changes and create Go APIView using generated artifacts"
    &($detectApiChanges) -ArtifactPath $ArtifactPath -CommitSha $CommitSha -BuildId $BuildId `
        -PullRequestNumber $PullRequestNumber -RepoFullName $RepoFullName -APIViewUri $APIViewUri `
        -ArtifactName $ArtifactName -DevopsProject $DevopsProject
}


