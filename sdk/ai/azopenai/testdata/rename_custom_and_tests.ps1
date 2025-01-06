param (
    [switch]$Reverse
)

# Function to rename files
function Rename-GoFiles {
    param (
        [string]$Directory,
        [switch]$Reverse
    )

    $files = if ($Reverse) {
        Get-ChildItem -Path $Directory -Filter "*.go" | Where-Object {
            $_.Name -match "^_(custom.*\.go|.*_test\.go)$"
        }
    } else {
        Get-ChildItem -Path $Directory -Filter "*.go" | Where-Object {
            $_.Name -match "^(custom.*\.go|[^_].*_test\.go)$"
        }
    }

    foreach ($file in $files) {
        $newName = if ($Reverse) {
            # Remove the leading underscore if reversing
            $file.Name -replace "^_", ""
        } else {
            # Add a leading underscore to ignore the file
            "_" + $file.Name
        }

        # Rename the file
        Rename-Item -Path $file.FullName -NewName $newName
    }
}

# Get the current directory
$directory = Get-Location

# Call the function to rename files
Rename-GoFiles -Directory $directory -Reverse:$Reverse