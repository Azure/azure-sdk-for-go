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
Get-Childitem -recurse -path $path -filter go.mod | foreach-object { $modDirs += $_.Directory }

# return the list of module directories
return $modDirs
