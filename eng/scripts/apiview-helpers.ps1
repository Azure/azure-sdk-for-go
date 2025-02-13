Write-Host "$PSScriptRoot"
. (Join-Path $PSScriptRoot .. common scripts common.ps1)

<#
.DESCRIPTION
    Create .gosource APIVIew artifact for go
.PARAMETER ServiceDirectory
    Thee name of the ServiceDirectory
.PARAMETER OutputDirectory
    Base output Directory path for the generated gosource artifacts
.PARAMETER DirectoryToPublish
    Directory containing all artifacts to be publisehd to the pipeline
#>
function New-APIViewArtifacts {
        Param(
        [Parameter(Mandatory=$True)]
        [string] $ServiceDirectory,
        [Parameter(Mandatory=$True)]
        [string] $OutputDirectory,
        [Parameter(Mandatory=$True)]
        [string] $DirectoryToPublish       
    )

    foreach ($sdk in (Get-AllPackageInfoFromRepo $ServiceDirectory))
    {
        Write-Host "Creating API review artifact for $($sdk.Name)"
        $sdkDirectoryPath = Join-Path -Path $OutputDirectory $sdk.Name
        New-Item -ItemType Directory -Path $sdkDirectoryPath -force
        $fileName = Split-Path -Path $sdk.Name -Leaf
        $compressedArchivePath = Join-Path $sdkDirectoryPath "$fileName.zip"
        Compress-Archive -Path $sdk.DirectoryPath -DestinationPath $compressedArchivePath -force
        Rename-Item $compressedArchivePath -NewName "$fileName.gosource"

        $artifactParentDirectory = $sdk.Name -Split "/" | Select-Object -First 1
        Copy-Item "$OutputDirectory/$artifactParentDirectory" -Destination "$DirectoryToPublish/$artifactParentDirectory" -Recurse
    }
}

<#
.DESCRIPTION
    Create new automatic APIView from a CI run
.PARAMETER ServiceDirectory
    Thee name of the ServiceDirectory
.PARAMETER ArtifactPath
    Directory containing the gosources artifact
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
function New-APIViewFromCI {
    Param(
        [Parameter(Mandatory=$True)]
        [string] $ServiceDirectory,
        [Parameter(Mandatory=$True)]
        [string] $ArtifactPath,
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
    $artifactList = @()
    
    Get-AllPackageInfoFromRepo $ServiceDirectory | ForEach-Object { 
        $artifactList += [PSCustomObject]@{
            name = $sdk.Name
        }
    }

    $createReviewScript = (Join-Path $PSScriptRoot .. common scripts Create-APIReview.ps1)

    Write-Host "Create Go APIView using generated artifacts"
    &($createReviewScript) `
        -ArtifactList $artifactList `
        -ArtifactPath $ArtifactPath `
        -APIKey $ApiKey `
        -SourceBranch $SourceBranch `
        -DefaultBranch $DefaultBranch `
        -ConfigFileDir $ConfigFileDir `
        -RepoName $RepoName `
        -BuildId $BuildId `
        -MarkPackageAsShipped $MarkPackageAsShipped
}