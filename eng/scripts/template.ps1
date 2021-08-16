#Requires -Version 7.0
param(
    # The name of the RP which is the directory name in azure-rest-api-specs/specification
    [Parameter(Mandatory=$true)]
    [string]$rpName,
    # The name of the package
    [Parameter(Mandatory=$true)]
    [string]$packageName,
    # The friendly package title of this package which will be the title of README
    [Parameter(Mandatory=$true)]
    [string]$packageTitle, 
    # The commit hash of the azure-rest-api-specs on its main branch to generate the SDK
    [Parameter(Mandatory=$true)]
    [string]$commitHash
)

function ApplyTemplate([System.IO.FileSystemInfo]$templateFile, [string]$targetDirectory) {
    Write-Host "##[command]Copying " $templateFile.Name "to " $targetDirectory
    if ($templateFile.Extension -ne ".tpl") {
        Write-Host "skip " $templateFile
        return
    }
    $content = Get-Content -Path $templateFile
    for ($i = 0; $i -lt $content.Count; $i++) {
        $content[$i] = $content[$i].Replace("{{rpName}}", $rpName).Replace("{{commitID}}", $commitHash).Replace("{{packageName}}", $packageName).Replace("{{PackageTitle}}", $packageTitle)
    }
    $targetFile = Join-Path $targetDirectory $templateFile.Name.Replace(".tpl", "")
    New-Item $targetFile > $null
    Set-Content $targetFile $content
}

$startignDirectory = Get-Location
$root = Resolve-Path ($PSScriptRoot + "/../..")
Set-Location $root
$targetDirectory = Join-Path $root "sdk" $rpName $packageName
$templateDirectory = Join-Path $root "eng/template"
try {
    New-Item -Path $targetDirectory -ItemType "directory" -Force > $null
    Get-ChildItem $templateDirectory | ForEach-Object {
        ApplyTemplate $_ $targetDirectory
    }
}
finally {
    Set-Location $startignDirectory
}
