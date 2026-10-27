param(
    $Packages
)

. (Join-Path $PSScriptRoot .. common scripts common.ps1)
. (Join-Path $PSScriptRoot analysis-helpers.ps1)

function Invoke-CopyrightCheck {
    param(
        [string]$Path
    )

    Write-Host "##[command]Check source file formatting in $Path"
    $missing = Invoke-WithGoAnalysisTools -PackageDirectory $Path -Action {
        Push-Location $Path
        try {
            $output = gofmt -s -l -d .
            if ($LASTEXITCODE -notin @(0, 1)) {
                throw "gofmt failed for $Path (exit code $LASTEXITCODE)"
            }
            $output
        }
        finally {
            Pop-Location
        }
    }

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
