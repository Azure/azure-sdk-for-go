Param(
    [string] $serviceDir
)

function Get-ModuleDirs ([string] $serviceDir) {
    $modDirs = [Collections.Generic.List[String]]@()

    # find each module directory under $serviceDir
    Get-ChildItem -recurse -path $serviceDir -filter go.mod | ForEach-Object {
        $cdir = $_.Directory
        Write-Host "Adding $cdir to list of module paths"
        $modDirs.Add($cdir)
    }

    # return the list of module directories
    return $modDirs
}

if ($MyInvocation.InvocationName -ne ".") {
    Get-ModuleDirs $serviceDir
}
