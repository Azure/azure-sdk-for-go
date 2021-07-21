Param(
  [string] $serviceDir
)

$testDirs = [Collections.Generic.List[String]]@()

# find each directory under $serviceDir that contains Go test files
Get-Childitem -recurse -path $serviceDir -filter *_test.go | foreach-object {
  $cdir = $_.Directory
  $tests = Select-String -Path $_ 'Test' -AllMatches

  if ($tests.Count -gt 0) {
    if (!$testDirs.Contains($cdir)) {
      Write-Host "Adding $cdir to list of test directories"
      $testDirs.Add($cdir)
    }
  }
}

# return the list of test directories
return $testDirs
