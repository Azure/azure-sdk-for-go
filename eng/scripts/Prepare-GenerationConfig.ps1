param(
    [Parameter(Mandatory = $false)]
    [string]$InputPath,
    
    [Parameter(Mandatory = $false)]
    [string]$PrUrl,
    
    [Parameter(Mandatory = $false)]
    [string]$LocalRepoPath,
    
    [Parameter(Mandatory = $false)]
    [string]$WorkspaceRoot = (Get-Location).Path
)

<#
.SYNOPSIS
Prepares the generation config (generatedInput.json) for Go SDK generation.

.DESCRIPTION
This script help to prepare the generation config file. It can handle:
- Local file paths to tspconfig.yaml or main.tsp files
- GitHub PR URLs (with optional local repository path)

.PARAMETER InputPath
Local path to a tspconfig.yaml file or directory containing it.

.PARAMETER PrUrl
GitHub PR URL to extract specification changes from.

.PARAMETER LocalRepoPath
Local path to the spec repository (used with PrUrl).

.PARAMETER WorkspaceRoot
Root directory of the workspace where generatedInput.json will be created.

.EXAMPLE
.\Prepare-GenerationConfig.ps1 -InputPath "C:\specs\azure-rest-api-specs\specification\compute\resource-manager\tspconfig.yaml"

.EXAMPLE
.\Prepare-GenerationConfig.ps1 -PrUrl "https://github.com/Azure/azure-rest-api-specs/pull/12345"

.EXAMPLE
.\Prepare-GenerationConfig.ps1 -PrUrl "https://github.com/Azure/azure-rest-api-specs/pull/12345" -LocalRepoPath "C:\repos\azure-rest-api-specs"
#>

function Write-Info {
    param([string]$Message)
    Write-Host "INFO: $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "WARNING: $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "ERROR: $Message" -ForegroundColor Red
}

function Test-GitHubCli {
    try {
        $null = Get-Command gh -ErrorAction Stop
        return $true
    }
    catch {
        Write-Error "GitHub CLI (gh) is not installed or not in PATH. Please install it from https://cli.github.com/"
        return $false
    }
}

function Extract-PrNumber {
    param([string]$PrUrl)
    
    if ($PrUrl -match "github\.com/.+/pull/(\d+)") {
        return $matches[1]
    }
    
    Write-Error "Invalid GitHub PR URL format. Expected format: https://github.com/owner/repo/pull/number"
    return $null
}

function Find-TspConfigPath {
    param([string]$StartPath)
    
    $currentPath = $StartPath
    
    while ($currentPath -and (Test-Path $currentPath)) {
        $tspConfigPath = Join-Path $currentPath "tspconfig.yaml"
        if (Test-Path $tspConfigPath) {
            return $currentPath
        }
        
        $parent = Split-Path $currentPath -Parent
        if ($parent -eq $currentPath) {
            break
        }
        $currentPath = $parent
    }
    
    return $null
}

function Get-SpecPathFromPr {
    param(
        [string]$PrNumber,
        [string]$LocalRepoPath
    )
    
    $originalLocation = Get-Location
    
    try {
        # Determine repository path
        if ($LocalRepoPath) {
            if (-not (Test-Path $LocalRepoPath)) {
                Write-Error "Local repository path does not exist: $LocalRepoPath"
                return $null
            }
            $repoPath = $LocalRepoPath
            Write-Info "Using local repository: $repoPath"
        }
        else {
            # Clone to temp directory
            $tempDir = Join-Path ([System.IO.Path]::GetTempPath()) "azure-rest-api-specs-$([System.Guid]::NewGuid().ToString('N')[0..7] -join '')"
            Write-Info "Cloning repository to: $tempDir"
            
            $cloneResult = & gh repo clone Azure/azure-rest-api-specs $tempDir 2>&1
            if ($LASTEXITCODE -ne 0) {
                Write-Error "Failed to clone repository: $cloneResult"
                return $null
            }
            $repoPath = $tempDir
        }
        
        # Navigate to repository
        Set-Location $repoPath
        
        # Checkout PR branch
        Write-Info "Checking out PR #$PrNumber"
        $checkoutResult = & gh pr checkout $PrNumber 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Failed to checkout PR: $checkoutResult"
            return $null
        }
        
        # Get changed files
        Write-Info "Getting changed files from PR"
        $changedFiles = & gh pr view $PrNumber --json files --jq '.files[].path' 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Failed to get PR files: $changedFiles"
            return $null
        }
        
        # Filter for .tsp files
        $tspFiles = $changedFiles | Where-Object { $_ -like "*.tsp" }
        
        if (-not $tspFiles) {
            Write-Warning "No .tsp files found in PR changes"
            return $null
        }
        
        Write-Info "Found .tsp files: $($tspFiles -join ', ')"
        
        # Find tspconfig.yaml for each .tsp file
        foreach ($tspFile in $tspFiles) {
            $tspFilePath = Join-Path $repoPath $tspFile
            $tspFileDir = Split-Path $tspFilePath -Parent
            
            $specPath = Find-TspConfigPath $tspFileDir
            if ($specPath) {
                Write-Info "Found tspconfig.yaml in: $specPath"
                return $specPath
            }
        }
        
        Write-Warning "No tspconfig.yaml found for any .tsp files"
        return $null
    }
    finally {
        Set-Location $originalLocation
    }
}

function Split-SpecPath {
    param([string]$SpecPath)
    
    $normalizedPath = $SpecPath -replace '\\', '/'
    
    if ($normalizedPath -match '^(.+)/specification/(.+)$') {
        return @{
            SpecFolder    = $matches[1] -replace '/', '\'
            ProjectFolder = "specification/$($matches[2])"
        }
    }
    
    Write-Error "Invalid spec path format. Expected path to contain '/specification/' directory."
    return $null
}

function Get-GitInfo {
    param([string]$SpecFolder)
    
    $originalLocation = Get-Location
    
    try {
        Set-Location $SpecFolder
        
        # Get HEAD SHA
        $headSha = & git rev-parse HEAD 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Failed to get git HEAD SHA: $headSha"
            return $null
        }
        
        # Get remote URL
        $remoteOutput = & git remote -v 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-Error "Failed to get git remote: $remoteOutput"
            return $null
        }
        
        Write-Info "Git remote output: $remoteOutput"

        $repoHttpsUrl = $null
        if ($remoteOutput -match 'azure-rest-api-specs-pr') {
            $repoHttpsUrl = "https://github.com/Azure/azure-rest-api-specs-pr"
        }
        elseif ($remoteOutput -match 'azure-rest-api-specs') {
            $repoHttpsUrl = "https://github.com/Azure/azure-rest-api-specs"
        }
        else {
            Write-Error "Unexpected git remote URL. Expected Azure/azure-rest-api-specs or Azure/azure-rest-api-specs-pr"
            return $null
        }
        
        return @{
            HeadSha      = $headSha.Trim()
            RepoHttpsUrl = $repoHttpsUrl
        }
    }
    finally {
        Set-Location $originalLocation
    }
}

function New-GeneratedInputJson {
    param(
        [string]$SpecFolder,
        [string]$HeadSha,
        [string]$RepoHttpsUrl,
        [string]$ProjectFolder,
        [string]$OutputPath
    )
    
    $config = @{
        specFolder                   = $SpecFolder
        headSha                      = $HeadSha
        repoHttpsUrl                 = $RepoHttpsUrl
        relatedTypeSpecProjectFolder = @($ProjectFolder)
    }
    
    $jsonContent = $config | ConvertTo-Json -Depth 10
    
    try {
        $jsonContent | Set-Content -Path $OutputPath -Encoding UTF8
        Write-Info "Generated config file: $OutputPath"
        Write-Info "Config content:"
        Write-Host $jsonContent
        return $true
    }
    catch {
        Write-Error "Failed to write config file: $_"
        return $false
    }
}

# Main execution
Write-Info "Starting generation config preparation..."

# Validate prerequisites for GitHub operations
if ($PrUrl) {
    if (-not (Test-GitHubCli)) {
        exit 1
    }
}

# Determine spec path
$specPath = $null

if ($PrUrl) {
    # GitHub PR scenario
    Write-Info "Processing GitHub PR: $PrUrl"
    
    $prNumber = Extract-PrNumber $PrUrl
    if (-not $prNumber) {
        exit 1
    }
    
    $specPath = Get-SpecPathFromPr $prNumber $LocalRepoPath
}
elseif ($InputPath) {
    # Local file path scenario
    Write-Info "Processing local path: $InputPath"
    
    if (Test-Path $InputPath -PathType Leaf) {
        # It's a file, use its parent directory
        $specPath = Split-Path $InputPath -Parent
    }
    elseif (Test-Path $InputPath -PathType Container) {
        # It's a directory
        $specPath = $InputPath
    }
    else {
        Write-Error "Input path does not exist: $InputPath"
        exit 1
    }
    
    # Verify tspconfig.yaml exists
    $tspConfigPath = Join-Path $specPath "tspconfig.yaml"
    if (-not (Test-Path $tspConfigPath)) {
        Write-Error "tspconfig.yaml not found in: $specPath"
        exit 1
    }
}
else {
    # Try to detect from open editor files (simplified - would need VS Code extension API)
    Write-Warning "No input provided. Please specify either -InputPath or -PrUrl"
    Write-Host "Usage examples:"
    Write-Host "  .\Prepare-GenerationConfig.ps1 -InputPath 'C:\specs\azure-rest-api-specs\specification\storageactions\StorageAction.Management\tspconfig.yaml'"
    Write-Host "  .\Prepare-GenerationConfig.ps1 -PrUrl 'https://github.com/Azure/azure-rest-api-specs/pull/12345'"
    exit 1
}

if (-not $specPath) {
    Write-Error "Failed to determine specification path"
    exit 1
}

Write-Info "Using spec path: $specPath"

# Split spec path
$pathInfo = Split-SpecPath $specPath
if (-not $pathInfo) {
    exit 1
}

Write-Info "Spec folder: $($pathInfo.SpecFolder)"
Write-Info "Project folder: $($pathInfo.ProjectFolder)"

# Get git information
$gitInfo = Get-GitInfo $pathInfo.SpecFolder
if (-not $gitInfo) {
    exit 1
}

Write-Info "HEAD SHA: $($gitInfo.HeadSha)"
Write-Info "Repository URL: $($gitInfo.RepoHttpsUrl)"

# Generate output file
$outputPath = Join-Path $WorkspaceRoot "generatedInput.json"
$success = New-GeneratedInputJson -SpecFolder $pathInfo.SpecFolder -HeadSha $gitInfo.HeadSha -RepoHttpsUrl $gitInfo.RepoHttpsUrl -ProjectFolder $pathInfo.ProjectFolder -OutputPath $outputPath

if ($success) {
    Write-Info "Generation config preparation completed successfully!"
    exit 0
}
else {
    Write-Error "Generation config preparation failed!"
    exit 1
}
