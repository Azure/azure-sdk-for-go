param(
    $ServiceDirectories
)

function Invoke-CopyrightCheck {
    param(
        [string]$Path
    )

    Write-Host "##[command]Executing copyright check in $Path"
    $missing = Get-ChildItem -Path $Path -Recurse -Include *.go -File | ForEach-Object {
        $results = Select-String -Path $_.FullName -Pattern 'Copyright (\d{4}|\(c\)) Microsoft'
        if (-not $results) { $_.FullName }
    }

    return $missing
}

$services = $ServiceDirectories.Split(",")

$missing = @()
foreach($service in $services) {
    $missing += Invoke-CopyrightCheck -Path $service
}

if ($missing) {
    Write-Error "Some go files are missing the copyright header. Please add the copyright header to the following files: "
    foreach ($file in $missing) {
        Write-Host " -> $file"
    }
}
else {
    Write-Host "Copyright check succeeded. All go files within service directories [$ServiceDirectories] have the copyright header."
}
