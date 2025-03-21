param(
    $Packages
)

function Invoke-CopyrightCheck {
    param(
        [string]$Path
    )

    Write-Host "##[command]Check source file formatting in $Path"
    Push-Location $Path
    $missing = gofmt -s -l -d .
    Pop-Location

    return $missing
}

$pkgs = $Packages.Split(",")

$missing = @()
foreach($pkg in $pkgs) {
    $missing += Invoke-CopyrightCheck -Path $pkg
}

if ($missing) {
    Write-Error "The following go files are not formatted correctly. Please go fmt the following files: "
    foreach ($file in $missing) {
        Write-Host " -> $file"
    }
}
else {
    Write-Host "Format check succeeded. All go files within directories [$Packages] are properly formatted."
}
