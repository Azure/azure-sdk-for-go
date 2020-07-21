Param(
  [string] $serviceDir = ""
)

if($serviceDir){
  $targetDir = "$PSScriptRoot/../../sdk/$serviceDir"
}
else {
  $targetDir = "$PSScriptRoot/../../sdk"
}

$path = Resolve-Path -Path $targetDir
$modDirs = [System.Collections.ArrayList]@()

# find each module directory under $path
Get-Childitem -recurse -path $path -filter go.mod | foreach-object {
  $cdir = $_.Directory
  Write-Host "Adding $cdir to list of module paths"
  $modDirs += $cdir
}

# return the list of module directories
return $modDirs
