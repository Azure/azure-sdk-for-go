[CmdletBinding()]
param([Parameter(Mandatory)][ValidateSet('Before', 'Install', 'After')][string]$Phase)

Set-StrictMode -Version 4
$ErrorActionPreference = 'Stop'
$script:Sequence = 0
$script:Bytes = 0
$script:UnparsedErrors = 0
$script:WithheldLines = 0
$script:ClientPathsResolved = $false
$script:CacheCandidates = [Collections.Generic.List[hashtable]]::new()
$package = 'azure.sdk.tools.testproxy'
$version = '1.0.0-dev.20260521.2'
$cwd = $env:BUILD_SOURCESDIRECTORY
$out = Join-Path $env:BUILD_ARTIFACTSTAGINGDIRECTORY 'go-nuget-boundary'
[IO.Directory]::CreateDirectory($out) | Out-Null
$destination = Join-Path $out "$Phase.jsonl"
if (Test-Path $destination) { throw 'Diagnostic output already exists; refusing a repeated phase.' }
if ($env:AGENT_OS -ne 'Windows_NT') { throw 'Only the bounded Windows diagnostic is permitted.' }
if ($cwd -ne 'D:\a\_work\1\s' -or $env:BUILD_BINARIESDIRECTORY -ne 'D:\a\_work\1\b') {
    throw 'The diagnostic working directory or destination differs from the approved boundary.'
}
Set-Location $cwd

function Get-Id([string]$Text) {
    [Convert]::ToHexString([Security.Cryptography.SHA256]::HashData([Text.Encoding]::UTF8.GetBytes($Text))).ToLowerInvariant()
}

function Write-Record([string]$Kind, [hashtable]$Data) {
    $script:Sequence++
    $record = @{ phase = $Phase; sequence = $script:Sequence; utc = [DateTime]::UtcNow.ToString('o'); kind = $Kind; data = $Data }
    $line = $record | ConvertTo-Json -Depth 12 -Compress
    $script:Bytes += [Text.Encoding]::UTF8.GetByteCount($line)
    if ($script:Bytes -gt 10485760) { throw 'Sanitized diagnostic output exceeded its byte cap.' }
    [IO.File]::AppendAllText($destination, $line + [Environment]::NewLine)
}

function Get-SourceId([string]$Text) {
    $uri = $null
    if (-not [Uri]::TryCreate($Text, [UriKind]::Absolute, [ref]$uri)) { return 'non-uri-source' }
    if ($uri.Scheme -ne 'https') { return 'non-https-source' }
    $endpoint = Get-Id ($uri.Scheme + '://' + $uri.Host + ':' + $uri.Port + $uri.AbsolutePath)
    if ($uri.Host -eq 'api.nuget.org') { return 'NuGetPublic-' + $endpoint }
    if ($uri.Host -eq 'pkgs.dev.azure.com' -and $uri.AbsolutePath.StartsWith('/azure-sdk/public/')) { return 'AzureSdkPublic-' + $endpoint }
    return 'other-source-' + $endpoint
}

function Write-NativeLine([string]$Line) {
    $found = $false
    $sources = @([regex]::Matches($Line, 'https?://[^\s"''<>]+') | ForEach-Object { Get-SourceId $_.Value } | Sort-Object -Unique)
    $types = @([regex]::Matches($Line, '(?:System|NuGet|Microsoft\.DotNet)\.[A-Za-z0-9_.]+(?:Exception|Error)') | ForEach-Object { $_.Value } | Sort-Object -Unique)
    $action = [regex]::Match($Line, '\b(GET|CACHE|OK|NotFound|Retrying|Unauthorized|Forbidden)\b').Value
    $codes = @([regex]::Matches($Line, '(?i)(?:status(?:\s*code)?[^0-9]{0,64}|HTTP/[0-9.]+\s+)([1-5]\d{2})') | ForEach-Object { $_.Groups[1].Value })
    $hresults = @([regex]::Matches($Line, '0x[0-9a-fA-F]{8}') | ForEach-Object { $_.Value })
    $durations = @([regex]::Matches($Line, '\b(\d+)ms\b') | ForEach-Object { [long]$_.Groups[1].Value })
    $messages = @('Unable to load the service index', 'No such host is known', 'The SSL connection could not be established', 'The remote certificate is invalid', 'Connection refused', 'timed out') | Where-Object { $Line.Contains($_, [StringComparison]::OrdinalIgnoreCase) }
    if ($sources.Count -or $types.Count -or $action -or $codes.Count -or $hresults.Count -or @($messages).Count) {
        Write-Record 'native-event' @{ sources = $sources; action = $action; exceptionTypes = $types; statusCodes = $codes; hresults = $hresults; durationMilliseconds = $durations; knownMessages = @($messages) }
        $found = $true
    }
    if ($Line -match '^\s*(\d+)\.\s+(.+?)\s+\[(Enabled|Disabled)\]\s*$') {
        Write-Record 'resolved-source-state' @{ ordinal = [int]$Matches[1]; keyId = Get-Id $Matches[2]; enabled = ($Matches[3] -eq 'Enabled') }
        $found = $true
    }
    if ($Line -match '^\s*(?:Version:\s*|NuGet Command Line\s*)?(\d+\.\d+\.\d+(?:[.\w+-]*))\s*$') {
        Write-Record 'client-version' @{ version = $Matches[1] }
        $found = $true
    }
    if (-not $found) {
        $script:WithheldLines++
        if ($Line -match '(?i)exception|error|failed') { $script:UnparsedErrors++ }
    }
}

function Write-Configuration {
    $paths = [Collections.Generic.List[string]]::new()
    foreach ($machine in @(${env:ProgramFiles(x86)}, $env:ProgramFiles)) {
        if ($machine) {
            $directory = Join-Path $machine 'NuGet\Config'
            if (Test-Path $directory) {
                foreach ($file in (Get-ChildItem $directory -File -Filter '*.config' | Sort-Object Name)) { $paths.Add($file.FullName) }
            }
        }
    }
    foreach ($directory in @((Join-Path $env:APPDATA 'NuGet'), (Join-Path $env:USERPROFILE '.nuget\NuGet'))) {
        $file = Join-Path $directory 'NuGet.Config'
        if (Test-Path $file) { $paths.Add($file) }
        $additional = Join-Path $directory 'config'
        if (Test-Path $additional) {
            foreach ($file in (Get-ChildItem $additional -File -Filter '*.config' | Sort-Object Name)) { $paths.Add($file.FullName) }
        }
    }
    $ancestors = [Collections.Generic.List[string]]::new()
    $directory = [IO.DirectoryInfo]::new($cwd)
    while ($null -ne $directory) { $ancestors.Insert(0, $directory.FullName); $directory = $directory.Parent }
    foreach ($ancestor in $ancestors) {
        $file = Join-Path $ancestor 'NuGet.Config'
        if (Test-Path $file) { $paths.Add($file) }
        $global = Join-Path $ancestor 'global.json'
        if (Test-Path $global) {
            $g = Get-Content $global -Raw | ConvertFrom-Json -AsHashtable
            $sdk = @{}
            if ($g.ContainsKey('sdk')) {
                foreach ($key in @('version', 'rollForward', 'allowPrerelease')) {
                    if ($g.sdk.ContainsKey($key)) {
                        $value = [string]$g.sdk[$key]
                        if ($value -match '^(?:\d[\d.a-zA-Z+-]*|true|false|disable|patch|feature|minor|major|latestPatch|latestFeature|latestMinor|latestMajor)$') {
                            $sdk[$key] = $value
                        } else { $sdk[$key] = 'withheld-unrecognized-value' }
                    }
                }
            }
            Write-Record 'sdk-selection-candidate' @{ pathId = Get-Id $global; sdk = $sdk }
        }
    }
    $order = 0
    foreach ($path in ($paths | Select-Object -Unique)) {
        $order++
        $file = Get-Item $path
        if ($file.Length -gt 1048576) { throw 'Configuration exceeds the bounded capture size.' }
        $xml = [Xml.XmlDocument]::new()
        $xml.XmlResolver = $null
        $xml.LoadXml([IO.File]::ReadAllText($path))
        $sections = @()
        foreach ($section in @('packageSources', 'disabledPackageSources', 'config')) {
            foreach ($node in $xml.SelectNodes("/configuration/$section/*")) {
                if ($node.Name -eq 'clear') { $sections += @{ section = $section; clear = $true }; continue }
                if ($node.Name -ne 'add') { continue }
                $key = $node.GetAttribute('key')
                $value = $node.GetAttribute('value')
                if ($section -eq 'packageSources') {
                    $sections += @{ section = $section; keyId = Get-Id $key; source = Get-SourceId $value }
                } elseif ($section -eq 'disabledPackageSources') {
                    $sections += @{ section = $section; keyId = Get-Id $key; disabled = ($value -eq 'true') }
                } elseif ($key -in @('globalPackagesFolder', 'repositoryPath', 'httpCachePath')) {
                    $sections += @{ section = $section; key = $key; pathId = Get-Id $value }
                    $expanded = [Environment]::ExpandEnvironmentVariables($value)
                    if (-not [IO.Path]::IsPathRooted($expanded)) { $expanded = Join-Path $file.DirectoryName $expanded }
                    if ($key -eq 'globalPackagesFolder') {
                        $script:CacheCandidates.Add(@{ kind = 'config-package-candidate'; path = Join-Path $expanded "$package\$version"; depth = 2 })
                    } elseif ($key -eq 'httpCachePath') {
                        $script:CacheCandidates.Add(@{ kind = 'config-http-candidate'; path = $expanded; depth = 1 })
                    }
                }
            }
        }
        $mapping = @()
        foreach ($source in $xml.SelectNodes('/configuration/packageSourceMapping/packageSource')) {
            $applies = $false
            foreach ($pattern in $source.SelectNodes('package')) {
                if ($package -like $pattern.GetAttribute('pattern')) { $applies = $true }
            }
            $mapping += @{ keyId = Get-Id $source.GetAttribute('key'); appliesToRequestedTool = $applies }
        }
        Write-Record 'config-candidate' @{ candidateOrder = $order; pathId = Get-Id $path; sections = $sections; mapping = $mapping; credentialsPresent = ($null -ne $xml.SelectSingleNode('/configuration/packageSourceCredentials')) }
    }
    Write-Record 'configuration-coverage' @{ status = 'candidates-only-until-client-resolved-paths-confirmed'; count = $order }
}

function Write-CacheMetadata {
    $roots = @()
    $http = if ($env:NUGET_HTTP_CACHE_PATH) { $env:NUGET_HTTP_CACHE_PATH } else { Join-Path $env:LOCALAPPDATA 'NuGet\v3-cache' }
    $packages = if ($env:NUGET_PACKAGES) { $env:NUGET_PACKAGES } else { Join-Path $env:USERPROFILE '.nuget\packages' }
    $roots += @{ kind = 'http-cache'; path = $http; depth = 1 }
    $roots += @{ kind = 'requested-package'; path = Join-Path $packages "$package\$version"; depth = 2 }
    $roots += @{ kind = 'tool-destination'; path = Join-Path $env:BUILD_BINARIESDIRECTORY 'test-proxy'; depth = 2 }
    $roots += @($script:CacheCandidates.ToArray())
    $count = 0
    foreach ($root in $roots) {
        $exists = Test-Path $root.path
        Write-Record 'cache-root' @{ cacheKind = $root.kind; pathId = Get-Id $root.path; exists = $exists }
        if (-not $exists) { continue }
        Get-ChildItem $root.path -Depth $root.depth -File | ForEach-Object {
            $count++
            if ($count -gt 2000) { throw 'Cache metadata enumeration exceeded its entry cap.' }
            Write-Record 'cache-entry' @{ cacheKind = $root.kind; pathId = Get-Id $_.FullName; length = $_.Length; lastWriteUtc = $_.LastWriteTimeUtc.ToString('o') }
        }
    }
}

try {
    if ($Phase -eq 'Before' -and (Test-Path (Join-Path $env:BUILD_BINARIESDIRECTORY 'test-proxy'))) {
        throw 'Tool destination already exists. Preserve it; no install, deletion or replacement is permitted.'
    }
    $dotnet = (Get-Command dotnet -CommandType Application -ErrorAction Stop).Source
    Write-Record 'identity' @{ build = $env:BUILD_BUILDID; job = $env:SYSTEM_JOBID; attempt = $env:SYSTEM_JOBATTEMPT; agent = $env:AGENT_NAME; machine = $env:AGENT_MACHINENAME; cwdId = Get-Id $cwd; executablePathId = Get-Id $dotnet }
    $binaryVersion = (Get-Item $dotnet).VersionInfo.FileVersion
    Write-Record 'dotnet-binary' @{ version = if ($binaryVersion -match '^\d[\d.a-zA-Z+ -]*$') { $binaryVersion } else { 'withheld-or-unavailable' } }
    $sdkDirectory = Join-Path (Split-Path $dotnet) 'sdk'
    if (Test-Path $sdkDirectory) {
        $sdkVersions = @(Get-ChildItem $sdkDirectory -Directory | Where-Object { $_.Name -match '^\d+\.\d+\.\d+[a-zA-Z0-9.+-]*$' } | ForEach-Object { $_.Name })
        Write-Record 'installed-sdk-versions' @{ versions = $sdkVersions }
    }
    foreach ($name in @('ImageVersion', 'ImageOS', 'AGENT_VERSION')) {
        $value = [Environment]::GetEnvironmentVariable($name)
        Write-Record 'image-variable' @{ name = $name; value = if ($null -ne $value -and $value -match '^[a-zA-Z0-9_.-]{1,64}$') { $value } else { 'unavailable-or-withheld' } }
    }
    foreach ($name in @('HTTP_PROXY', 'HTTPS_PROXY')) {
        $value = [Environment]::GetEnvironmentVariable($name)
        $proxy = $null
        $parsed = $value -and [Uri]::TryCreate($value, [UriKind]::Absolute, [ref]$proxy)
        Write-Record 'proxy-selection' @{ name = $name; present = ($null -ne $value); hostId = if ($parsed) { Get-Id $proxy.Host } else { $null }; scheme = if ($parsed -and $proxy.Scheme -in @('http', 'https')) { $proxy.Scheme } else { 'unavailable' }; port = if ($parsed) { $proxy.Port } else { $null } }
    }
    $bypass = @($env:NO_PROXY -split ',' | ForEach-Object { $_.Trim().ToLowerInvariant() })
    foreach ($hostName in @('api.nuget.org', 'pkgs.dev.azure.com')) {
        $matching = @($bypass | Where-Object { $_ -eq $hostName -or ($_ -match '^\.[a-z0-9.-]+$' -and $hostName.EndsWith($_)) })
        Write-Record 'proxy-bypass-candidate' @{ publicHost = $hostName; exactOrSuffixEntryPresent = ($matching.Count -gt 0); otherBypassSyntaxNotEvaluated = $true }
    }
    if ($Phase -eq 'Install') {
        Write-Record 'command' @{ package = $package; version = $version; source = 'AzureSdkPublic'; verbosity = 'diagnostic'; cwd = 'Build.SourcesDirectory'; destination = 'Build.BinariesDirectory/test-proxy' }
        $args = @('tool', 'install', $package, '--tool-path', "$env:BUILD_BINARIESDIRECTORY/test-proxy", '--add-source', 'https://pkgs.dev.azure.com/azure-sdk/public/_packaging/azure-sdk-for-net/nuget/v3/index.json', '--version', $version, '--verbosity', 'diagnostic')
        $timer = [Diagnostics.Stopwatch]::StartNew()
        & $dotnet @args 2>&1 | ForEach-Object { Write-NativeLine ([string]$_) }
        $exit = $LASTEXITCODE
        Write-Record 'native-result' @{ exitCode = $exit; elapsedMilliseconds = $timer.ElapsedMilliseconds; withheldLines = $script:WithheldLines; unparsedErrorLines = $script:UnparsedErrors }
        exit $exit
    }
    foreach ($name in @('DOTNET_ROOT', 'DOTNET_HOST_PATH', 'DOTNET_CLI_HOME', 'DOTNET_MULTILEVEL_LOOKUP', 'DOTNET_ROLL_FORWARD', 'NUGET_PACKAGES', 'NUGET_HTTP_CACHE_PATH', 'NUGET_FALLBACK_PACKAGES', 'NUGET_PLUGIN_PATHS')) {
        $value = [Environment]::GetEnvironmentVariable($name)
        Write-Record 'selection-variable' @{ name = $name; present = ($null -ne $value); valueId = if ($null -ne $value) { Get-Id $value } else { $null } }
    }
    Write-Configuration
    Write-CacheMetadata
    if ($Phase -eq 'After') {
        foreach ($command in @(@('--info'), @('nuget', '--version'), @('nuget', 'config', 'paths', '--working-directory', $cwd), @('nuget', 'list', 'source', '--format', 'Detailed'), @('nuget', 'locals', 'all', '--list'))) {
            Write-Record 'read-command-start' @{ command = @($command[0..([Math]::Min(2, $command.Count - 1))]) }
            & $dotnet @command 2>&1 | ForEach-Object {
                $line = [string]$_
                if ($command.Count -ge 3 -and $command[1] -eq 'config' -and $line -match '^[A-Za-z]:\\') {
                    Write-Record 'client-config-path' @{ pathId = Get-Id $line.Trim() }
                } elseif ($command.Count -ge 3 -and $command[1] -eq 'locals' -and $line -match '^([\w-]+):\s+([A-Za-z]:\\.+)$') {
                    Write-Record 'client-cache-path' @{ cacheKind = $Matches[1]; pathId = Get-Id $Matches[2] }
                } else { Write-NativeLine $line }
            }
            if ($command.Count -ge 3 -and $command[1] -eq 'config' -and $LASTEXITCODE -eq 0) { $script:ClientPathsResolved = $true }
            Write-Record 'read-command-result' @{ command = @($command[0..([Math]::Min(2, $command.Count - 1))]); exitCode = $LASTEXITCODE }
        }
    }
    Write-Record 'coverage' @{ withheldLines = $script:WithheldLines; unparsedErrorLines = $script:UnparsedErrors; cachePayloadsRead = $false; clientPathsResolved = $script:ClientPathsResolved }
}
catch {
    Write-Record 'capture-failed' @{ exceptionType = $_.Exception.GetType().FullName; rawMessageWithheld = $true }
    Write-Error 'Bounded diagnostic capture failed; unsafe details were withheld. No retry is permitted.'
    exit 1
}
