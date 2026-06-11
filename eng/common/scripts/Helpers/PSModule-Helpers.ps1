$global:CurrentUserModulePath = ""

# Internal Azure Artifacts NuGet v3 feed used to install all PowerShell modules.
# We deliberately do NOT use PSGallery: pipelines run under network policies
# that block outbound calls to powershellgallery.com.
$global:AzSdkPSResourceRepoUri = "https://pkgs.dev.azure.com/azure-sdk/public/_packaging/azure-sdk-tools/nuget/v3/index.json"
$global:AzSdkPSResourceRepoName = "AzureSdkTools"

function Update-PSModulePathForCI() {
  # Information on PSModulePath taken from docs
  # https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_psmodulepath

  # Information on Az custom module paths on hosted agents taken from
  # https://github.com/microsoft/azure-pipelines-tasks/blob/c9771bc064cd60f47587c68e5c871b7cd13f0f28/Tasks/AzurePowerShellV5/Utility.ps1

  if ($IsWindows) {
    $hostedAgentModulePath = $env:SystemDrive + "\Modules"
    $moduleSeperator = ";"
  }
  else {
    $hostedAgentModulePath = "/usr/share"
    $moduleSeperator = ":"
  }
  $modulePaths = $env:PSModulePath -split $moduleSeperator

  # Remove any hosted agent paths (needed to remove old default azure/azurerm paths which cause conflicts)
  $modulePaths = $modulePaths.Where({ !$_.StartsWith($hostedAgentModulePath) })

  # Add any "az_" paths from the agent which is the lastest set of azure modules
  $AzModuleCachePath = (Get-ChildItem "$hostedAgentModulePath/az_*" -Attributes Directory) -join $moduleSeperator
  if ($AzModuleCachePath -and $env:PSModulePath -notcontains $AzModuleCachePath) {
    $modulePaths += $AzModuleCachePath
  }

  $env:PSModulePath = $modulePaths -join $moduleSeperator

  # Find the path that is under user home directory
  $homeDirectories = $modulePaths.Where({ $_.StartsWith($home) })
  if ($homeDirectories.Count -gt 0) {
    $global:CurrentUserModulePath = $homeDirectories[0]
    if ($homeDirectories.Count -gt 1) {
      Write-Verbose "Found more then one module path starting with $home so selecting the first one $global:CurrentUserModulePath"
    }

    # In some cases the directory might not exist so we need to create it otherwise caching an empty directory will fail
    if (!(Test-Path $global:CurrentUserModulePath)) {
      New-Item $global:CurrentUserModulePath -ItemType Directory > $null
    }
  }
  else {
    Write-Error "Did not find a module path starting with $home to set up a user module path in $env:PSModulePath"
  }
}

function moduleIsInstalled([string]$moduleName, [string]$version) {
  if (-not (Test-Path variable:script:InstalledModules)) {
    $script:InstalledModules = @{}
  }

  if ($script:InstalledModules.ContainsKey("${moduleName}")) {
    $modules = $script:InstalledModules["${moduleName}"]
  }
  else {
    $modules = (Get-Module -ListAvailable $moduleName)
    $script:InstalledModules["${moduleName}"] = $modules
  }

  if ($version -as [Version]) {
    $modules = $modules.Where({ [Version]$_.Version -ge [Version]$version })
    if ($modules.Count -gt 0) {
      Write-Verbose "Using module $($modules[0].Name) with version $($modules[0].Version)."
      return $modules[0]
    }
  }
  return $null
}

function ensureAzSdkPSResourceRepository() {
  Import-Module Microsoft.PowerShell.PSResourceGet -ErrorAction Stop

  $repo = Get-PSResourceRepository -Name $global:AzSdkPSResourceRepoName -ErrorAction SilentlyContinue
  if (-not $repo) {
    Write-Verbose "Registering PSResource repository '$($global:AzSdkPSResourceRepoName)' -> $($global:AzSdkPSResourceRepoUri)"
    Register-PSResourceRepository `
      -Name $global:AzSdkPSResourceRepoName `
      -Uri $global:AzSdkPSResourceRepoUri `
      -Trusted `
      -ErrorAction Stop | Out-Null
    $repo = Get-PSResourceRepository -Name $global:AzSdkPSResourceRepoName -ErrorAction Stop
  }
  elseif (-not $repo.Trusted) {
    Set-PSResourceRepository -Name $global:AzSdkPSResourceRepoName -Trusted | Out-Null
  }

  return $repo
}

function installResource([string]$moduleName, [string]$version) {
  $repo = ensureAzSdkPSResourceRepository

  Write-Verbose "Installing module '$moduleName' version '$version' from '$($repo.Name)' via Install-PSResource"
  # -Repository is a hard scope in PSResourceGet: primary lookup and dependency
  # resolution are constrained to the named repository, so this never reaches
  # PSGallery (unlike PowerShellGet 2.x's Install-Module).
  Install-PSResource `
    -Name $moduleName `
    -Version $version `
    -Repository $repo.Name `
    -Scope CurrentUser `
    -TrustRepository `
    -AcceptLicense `
    -Reinstall:$false `
    -Quiet `
    -ErrorAction Stop

  # Install-PSResource does not emit a PSModuleInfo. Materialize one so callers
  # that pipe to Import-Module continue to work.
  $modules = (Get-Module -ListAvailable $moduleName)
  if ($version -as [Version]) {
    $modules = $modules.Where({ [Version]$_.Version -eq [Version]$version })
  }
  if ($modules.Count -eq 0) {
    throw "Failed to install module $moduleName with version $version"
  }

  if (-not (Test-Path variable:script:InstalledModules)) {
    $script:InstalledModules = @{}
  }
  $script:InstalledModules["${moduleName}"] = $modules
  return $modules[0]
}

function InstallAndImport-ModuleIfNotInstalled([string]$module, [string]$version) {
  if ($null -eq (moduleIsInstalled $module $version)) {
    Install-ModuleIfNotInstalled -WhatIf:$false $module $version | Import-Module
  } elseif (!(Get-Module -Name $module)) {
    Import-Module $module
  }
}

# Manual test at eng/common-tests/psmodule-helpers/Install-Module-Parallel.ps1
# Always installs from the internal Azure Artifacts feed via PSResourceGet
# (Microsoft.PowerShell.PSResourceGet, in-box on PowerShell 7.4+). PSGallery is
# never contacted.
function Install-ModuleIfNotInstalled() {
  [CmdletBinding(SupportsShouldProcess = $true)]
  param(
    [string]$moduleName,
    [string]$version,
    # Retained for back-compat with existing callers; ignored. All installs go
    # through the internal Azure Artifacts feed via PSResourceGet.
    [string]$repositoryUrl
  )

  # Check installed modules before after acquiring lock to avoid a big queue
  $module = moduleIsInstalled -moduleName $moduleName -version $version
  if ($module) { return $module }

  try {
    $mutex = New-Object System.Threading.Mutex($false, "Install-ModuleIfNotInstalled")
    $null = $mutex.WaitOne()

    # Check installed modules again after acquiring lock, in case it has been installed
    $module = moduleIsInstalled -moduleName $moduleName -version $version
    if ($module) { return $module }

    Write-Host "Module '$moduleName' with version '$version' is not installed. Installing from $($global:AzSdkPSResourceRepoUri)."

    $module = installResource -moduleName $moduleName -version $version
    Write-Verbose "Using module '$($module.Name)' with version '$($module.Version)'."
  }
  finally {
    $mutex.ReleaseMutex()
  }

  return $module
}

if ($null -ne $env:SYSTEM_TEAMPROJECTID) {
  Update-PSModulePathForCI
}
