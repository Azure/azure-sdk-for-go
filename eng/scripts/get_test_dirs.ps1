Param(
  [string] $serviceDir
)

$testDirs = [Collections.Generic.List[String]]@()

# find each module directory under $path
Get-Childitem -recurse -path $serviceDir -filter *_test.go | foreach-object {
  $cdir = $_.Directory
  if (!$testDirs.Contains($cdir)) {
    Write-Host "Adding $cdir to list of test directories"
    $testDirs.Add($cdir)
  }
}

# return the list of test directories
return $testDirs
