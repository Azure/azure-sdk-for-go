# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

<#
.SYNOPSIS
    Checks and installs prerequisites for Go SDK generation.

.DESCRIPTION
    This script verifies that all required tools are installed and properly configured
    for Azure Go SDK generation from API specifications. It checks:
    - Go version 1.23 or later
    - Generator tool v0.1.0
    - Node.js version 20 or later
    - TypeSpec client generator CLI v0.21.0
    - GitHub CLI tool and authentication
    - Git installation

.EXAMPLE
    .\Check-Prerequisites.ps1
    Runs all prerequisite checks and installs missing tools.

.EXAMPLE
    .\Check-Prerequisites.ps1 -CheckOnly
    Only checks prerequisites without installing anything.
#>

param(
    [switch]$CheckOnly = $false,
    [string]$MinGoVersion = "1.23",
    [string]$MinNodeVersion = "22.0.0"
)

$ErrorActionPreference = "Stop"

function Write-StatusMessage {
    param(
        [string]$Message,
        [string]$Status = "INFO"
    )
    
    $color = switch ($Status) {
        "SUCCESS" { "Green" }
        "ERROR" { "Red" }
        "WARNING" { "Yellow" }
        default { "White" }
    }
    
    Write-Host "[$Status] $Message" -ForegroundColor $color
}

function Test-ToolInstalled {
    param(
        [string]$ToolName,
        [string]$Command,
        [string]$VersionFlag = "--version",
        [string]$VersionRegex = "",
        [string]$MinVersion = "",
        [string]$Description = ""
    )
    
    $displayName = if ($Description) { $Description } else { $ToolName }
    Write-StatusMessage "Checking $displayName..."
    
    try {
        # Check if tool is available
        $toolPath = Get-Command $Command -ErrorAction SilentlyContinue
        if (-not $toolPath) {
            Write-StatusMessage "$displayName is not installed" "WARNING"
            return $false
        }
        
        # Get version information
        $versionOutput = & $Command $VersionFlag 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-StatusMessage "$displayName version check failed" "ERROR"
            return $false
        }
        
        # If no version checking is needed, just confirm it's installed
        if (-not $VersionRegex -and -not $MinVersion) {
            Write-StatusMessage "$displayName is installed ✓" "SUCCESS"
            return $true
        }
        
        # Extract version if regex is provided
        if ($VersionRegex) {
            if ($versionOutput -match $VersionRegex) {
                $version = $matches[1]
                
                # Check minimum version if specified
                if ($MinVersion) {
                    if (Compare-Version $version $MinVersion) {
                        Write-StatusMessage "$displayName version $version is installed ✓" "SUCCESS"
                        return $true
                    }
                    else {
                        Write-StatusMessage "$displayName version $version is too old. Minimum required: $MinVersion" "ERROR"
                        return $false
                    }
                }
                else {
                    Write-StatusMessage "$displayName version $version is installed ✓" "SUCCESS"
                    return $true
                }
            }
            else {
                Write-StatusMessage "Could not parse $displayName version: $versionOutput" "ERROR"
                return $false
            }
        }
        else {
            Write-StatusMessage "$displayName is installed ✓" "SUCCESS"
            return $true
        }
    }
    catch {
        Write-StatusMessage "$displayName check failed: $($_.Exception.Message)" "ERROR"
        return $false
    }
}

function Compare-Version {
    param(
        [string]$CurrentVersion,
        [string]$MinimumVersion
    )
    
    try {
        $current = [Version]::Parse($CurrentVersion)
        $minimum = [Version]::Parse($MinimumVersion)
        return $current -ge $minimum
    }
    catch {
        # Fallback to string comparison for complex version formats
        return $CurrentVersion -eq $MinimumVersion -or $CurrentVersion -gt $MinimumVersion
    }
}

function Install-GeneratorTool {
    if ($CheckOnly) {
        Write-StatusMessage "Skipping generator tool installation (CheckOnly mode)" "WARNING"
        return
    }
    
    Write-StatusMessage "Installing generator tool v0.1.0..."
    
    try {
        $installResult = go install github.com/Azure/azure-sdk-for-go/eng/tools/generator@v0.1.0 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-StatusMessage "Generator tool installed successfully ✓" "SUCCESS"
        }
        else {
            Write-StatusMessage "Failed to install generator tool: $installResult" "ERROR"
        }
    }
    catch {
        Write-StatusMessage "Generator tool installation failed: $($_.Exception.Message)" "ERROR"
    }
}

function Install-TypeSpecTool {
    if ($CheckOnly) {
        Write-StatusMessage "Skipping TypeSpec tool installation (CheckOnly mode)" "WARNING"
        return
    }
    
    Write-StatusMessage "Installing TypeSpec client generator CLI v0.21.0..."
    
    try {
        $installResult = npm install -g "@azure-tools/typespec-client-generator-cli@v0.21.0" 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-StatusMessage "TypeSpec client generator CLI installed successfully ✓" "SUCCESS"
        }
        else {
            Write-StatusMessage "Failed to install TypeSpec client generator CLI: $installResult" "ERROR"
        }
    }
    catch {
        Write-StatusMessage "TypeSpec client generator CLI installation failed: $($_.Exception.Message)" "ERROR"
    }
}

# Main execution
Write-StatusMessage "Starting prerequisite checks for Go SDK generation..." "INFO"
Write-StatusMessage "===============================================" "INFO"

$failedPrerequisites = @()

# Check Go version
if (-not (Test-ToolInstalled -ToolName "Go" -Command "go" -VersionFlag "version" -VersionRegex "go version go(\d+\.\d+)" -MinVersion $MinGoVersion -Description "Go installation and version")) {
    $failedPrerequisites += "Go $MinGoVersion or later (install from https://golang.org/dl/)"
}

# Check and install generator tool
if (-not (Test-ToolInstalled -ToolName "Generator" -Command "generator" -Description "generator tool")) {
    Install-GeneratorTool
    # Re-check after installation
    if (-not (Test-ToolInstalled -ToolName "Generator" -Command "generator" -Description "generator tool")) {
        $failedPrerequisites += "Generator tool (failed to install or not working properly)"
    }
}

# Check Node.js version
if (-not (Test-ToolInstalled -ToolName "Node.js" -Command "node" -VersionRegex "v(\d+\.\d+\.\d+)" -MinVersion $MinNodeVersion -Description "Node.js installation and version")) {
    $failedPrerequisites += "Node.js $MinNodeVersion or later (install from https://nodejs.org/)"
}

# Check and install TypeSpec tool
if (-not (Test-ToolInstalled -ToolName "TypeSpec" -Command "tsp-client" -Description "TypeSpec client generator CLI")) {
    Install-TypeSpecTool
    # Re-check after installation
    if (-not (Test-ToolInstalled -ToolName "TypeSpec" -Command "tsp-client" -Description "TypeSpec client generator CLI")) {
        $failedPrerequisites += "TypeSpec client generator CLI (failed to install or not working properly)"
    }
}

# Check GitHub CLI and authentication
Write-StatusMessage "Checking GitHub CLI installation and authentication..."
try {
    if (-not (Test-ToolInstalled -ToolName "GitHub CLI" -Command "gh" -Description "GitHub CLI")) {
        $failedPrerequisites += "GitHub CLI (install from https://cli.github.com/ and run 'gh auth login')"
    }
    else {
        # Check authentication
        gh auth status 2>&1 | Out-Null
        if ($LASTEXITCODE -eq 0) {
            Write-StatusMessage "GitHub CLI is authenticated ✓" "SUCCESS"
        }
        else {
            Write-StatusMessage "GitHub CLI is not authenticated. Please run 'gh auth login'" "ERROR"
            $failedPrerequisites += "GitHub CLI authentication (run 'gh auth login')"
        }
    }
}
catch {
    Write-StatusMessage "GitHub CLI check failed: $($_.Exception.Message)" "ERROR"
    $failedPrerequisites += "GitHub CLI (install from https://cli.github.com/ and run 'gh auth login')"
}

# Check Git
if (-not (Test-ToolInstalled -ToolName "Git" -Command "git" -Description "Git installation")) {
    $failedPrerequisites += "Git (install from https://git-scm.com/)"
}

Write-StatusMessage "===============================================" "INFO"

if ($failedPrerequisites.Count -eq 0) {
    Write-StatusMessage "All prerequisites are satisfied! ✓" "SUCCESS"
    exit 0
}
else {
    Write-StatusMessage "The following prerequisites are missing or failed to install:" "ERROR"
    foreach ($prerequisite in $failedPrerequisites) {
        Write-StatusMessage "  • $prerequisite" "ERROR"
    }
    Write-StatusMessage "Please install the missing prerequisites and run this script again." "ERROR"
    exit 1
}
