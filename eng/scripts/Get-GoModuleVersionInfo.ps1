# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

function Get-GoModuleVersionInfo($modPath)
{
  # Only SDK version declarations qualify; service enums can also contain semantic versions.
  $VERSION_CONSTANT_REGEX = "(?m)^\s*(?:const\s+)?(?:moduleVersion|serviceLibVersion|version)\s*(?:string\s*)?=\s*"
  $NO_PREFIX_VERSION_LINE_REGEX = "$VERSION_CONSTANT_REGEX`"(?<bad_version>$([AzureEngSemanticVersion]::SEMVER_REGEX))`""
  $VERSION_LINE_REGEX = "$VERSION_CONSTANT_REGEX`"v(?<version>$([AzureEngSemanticVersion]::SEMVER_REGEX))`""

  $versionFiles = Get-ChildItem -Recurse -Path $modPath -Filter *.go

  foreach ($versionFile in $versionFiles)
  {
    if (!$versionFile.Name.Contains("constant") -and !$versionFile.Name.Contains("version")) {
      continue
    }
    $content = Get-Content $versionFile -Raw

    if ($content -match $VERSION_LINE_REGEX) {
      return "$($matches["version"])", $versionFile.ToString()
    }

    if ($content -match $NO_PREFIX_VERSION_LINE_REGEX) {
      LogError "Version in $versionFile should be 'v$($matches["bad_version"])' not '$($matches["bad_version"])'"
    }
  }

  LogWarning "Unable to find version for $modPath"
  return $null
}
