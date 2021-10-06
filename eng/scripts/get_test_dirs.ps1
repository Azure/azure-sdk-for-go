Param(
    [string] $serviceDir
)

. (Join-Path $PSScriptRoot .. common scripts common.ps1)

$testDirs = @()

foreach ($sdk in (Get-AllPackageInfoFromRepo $serviceDir))
{
    # find each directory under $serviceDir that contains Go test files
    foreach ($testFile in (Get-ChildItem -recurse -path $sdk.DirectoryPath -filter *_test.go))
    {
        $cdir = $testFile.Directory.FullName
        $tests = Select-String -Path $testFile 'Test' -AllMatches

        if ($tests.Count -gt 0) {
            if ($testDirs -notcontains $cdir) {
                Write-Host "Adding $cdir to list of test directories"
                $testDirs += $cdir
            }
        }
    }
}

# return the list of test directories
return $testDirs