# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

<#
.SYNOPSIS
    Checks and installs prerequisites for Go SDK generation.

.DESCRIPTION
    This script verifies that all required tools are installed and properly configured
    for Azure Go SDK generation from API specifications.

.EXAMPLE
    .\Check-Prerequisites.ps1
    Runs all prerequisite checks and installs missing tools.
#>

param(
    [string]$MinGoVersion = "1.23",
    [string]$MinNodeVersion = "20.0.0",
    [string]$GeneratorVersion = "v0.1.0"
)

$ErrorActionPreference = "Stop"
$WarningPreference = "Stop"
Set-StrictMode -Version Latest

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
        [Parameter(Mandatory = $true)]
        [string]$Description,
        [Parameter(Mandatory = $true)]
        [string]$Command,
        [string]$VersionFlag = "--version",
        [string]$VersionRegex = "",
        [string]$MinVersion = ""
    )
    
    Write-StatusMessage "Checking $Description..."
    
    try {
        # Check if tool is available
        $toolPath = Get-Command $Command -ErrorAction SilentlyContinue
        if (-not $toolPath) {
            Write-StatusMessage "$Description is not installed" "WARNING"
            return $false
        }
        
        # Get version information
        $versionOutput = & $Command $VersionFlag 2>&1
        if ($LASTEXITCODE -ne 0) {
            Write-StatusMessage "$Description version check failed" "ERROR"
            return $false
        }
        
        # If no version checking is needed, just confirm it's installed
        if (-not $VersionRegex -and -not $MinVersion) {
            Write-StatusMessage "$Description is installed ✓" "SUCCESS"
            return $true
        }
        
        # Extract version if regex is provided
        if ($VersionRegex) {
            if ($versionOutput -match $VersionRegex) {
                $version = $matches[1]
                
                # Check minimum version if specified
                if ($MinVersion) {
                    if (Compare-Version $version $MinVersion) {
                        Write-StatusMessage "$Description version $version is installed ✓" "SUCCESS"
                        return $true
                    }
                    else {
                        Write-StatusMessage "$Description version $version is too old. Minimum required: $MinVersion" "ERROR"
                        return $false
                    }
                }
                else {
                    Write-StatusMessage "$Description version $version is installed ✓" "SUCCESS"
                    return $true
                }
            }
            else {
                Write-StatusMessage "Could not parse $Description version: $versionOutput" "ERROR"
                return $false
            }
        }
        else {
            Write-StatusMessage "$Description is installed ✓" "SUCCESS"
            return $true
        }
    }
    catch {
        Write-StatusMessage "$Description check failed: $($_.Exception.Message)" "ERROR"
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
    param(
        [Parameter(Mandatory = $true)]
        [string]$GeneratorVersion
    )

    Write-StatusMessage "Installing generator tool $GeneratorVersion..."

    try {
        $installResult = go install "github.com/Azure/azure-sdk-for-go/eng/tools/generator@$GeneratorVersion" 2>&1
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

function Install-TypeSpecCompiler {
    Write-StatusMessage "Installing TypeSpec compiler..."
    
    try {
        $installResult = npm install -g @typespec/compiler 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-StatusMessage "TypeSpec compiler installed successfully ✓" "SUCCESS"
        }
        else {
            Write-StatusMessage "Failed to install TypeSpec compiler: $installResult" "ERROR"
        }
    }
    catch {
        Write-StatusMessage "TypeSpec compiler installation failed: $($_.Exception.Message)" "ERROR"
    }
}

function Install-TypeSpecTool {
    Write-StatusMessage "Installing TypeSpec client generator CLI..."
    
    try {
        $installResult = npm install -g @azure-tools/typespec-client-generator-cli 2>&1
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
if (-not (Test-ToolInstalled -Description "Go installation and version" -Command "go" -VersionFlag "version" -VersionRegex "go version go(\d+\.\d+)" -MinVersion $MinGoVersion)) {
    $failedPrerequisites += "Go $MinGoVersion or later (install from https://golang.org/dl/)"
}

# Check and install generator tool
if (-not (Test-ToolInstalled -Description "generator tool" -Command "generator")) {
    Install-GeneratorTool -GeneratorVersion $GeneratorVersion
    # Re-check after installation
    if (-not (Test-ToolInstalled -Description "generator tool" -Command "generator")) {
        $failedPrerequisites += "Generator tool (failed to install or not working properly)"
    }
}

# Check Node.js version
if (-not (Test-ToolInstalled -Description "Node.js installation and version" -Command "node" -VersionRegex "v(\d+\.\d+\.\d+)" -MinVersion $MinNodeVersion)) {
    $failedPrerequisites += "Node.js $MinNodeVersion or later (install from https://nodejs.org/)"
}

# Check and install TypeSpec compiler
if (-not (Test-ToolInstalled -Description "TypeSpec compiler" -Command "tsp")) {
    Install-TypeSpecCompiler
    # Re-check after installation
    if (-not (Test-ToolInstalled -Description "TypeSpec compiler" -Command "tsp")) {
        $failedPrerequisites += "TypeSpec compiler (failed to install or not working properly)"
    }
}

# Check and install TypeSpec tool
if (-not (Test-ToolInstalled -Description "TypeSpec client generator CLI" -Command "tsp-client")) {
    Install-TypeSpecTool
    # Re-check after installation
    if (-not (Test-ToolInstalled -Description "TypeSpec client generator CLI" -Command "tsp-client")) {
        $failedPrerequisites += "TypeSpec client generator CLI (failed to install or not working properly)"
    }
}

# Check GitHub CLI and authentication
Write-StatusMessage "Checking GitHub CLI installation and authentication..."
try {
    if (-not (Test-ToolInstalled -Description "GitHub CLI" -Command "gh")) {
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