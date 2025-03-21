param(
    $ServiceDirectories
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

$services = $ServiceDirectories.Split(",")

$missing = @()
foreach($service in $services) {
    $missing += Invoke-CopyrightCheck -Path $service
}

if ($missing) {
    Write-Error "The following go files are not formatted correctly. Please go fmt the following files: "
    foreach ($file in $missing) {
        Write-Host " -> $file"
    }
}
else {
    Write-Host "Format check succeeded. All go files within service directories [$ServiceDirectories] are properly formatted."
}
