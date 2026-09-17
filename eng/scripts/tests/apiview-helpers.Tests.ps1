BeforeAll {
  . (Join-Path $PSScriptRoot .. apiview-helpers.ps1)
}

Describe "New-APIViewArtifacts" {
  It "publishes each package without nesting the sdk directory" {
    $sourceRoot = Join-Path $TestDrive source
    $outputDirectory = Join-Path $TestDrive output
    $directoryToPublish = Join-Path $TestDrive publish
    $firstPackageDirectory = Join-Path $sourceRoot first
    $secondPackageDirectory = Join-Path $sourceRoot second

    New-Item -ItemType Directory -Path $firstPackageDirectory, $secondPackageDirectory -Force | Out-Null
    Set-Content -Path (Join-Path $firstPackageDirectory client.go) -Value "package first"
    Set-Content -Path (Join-Path $secondPackageDirectory client.go) -Value "package second"

    Mock Get-AllPackageInfoFromRepo {
      param($ServiceDirectory)

      if ($ServiceDirectory -eq "sdk/service/first") {
        return [PSCustomObject]@{
          Name = "sdk/service/first"
          DirectoryPath = $firstPackageDirectory
        }
      }
      return [PSCustomObject]@{
        Name = "sdk/service/second"
        DirectoryPath = $secondPackageDirectory
      }
    }

    New-APIViewArtifacts `
      -ServiceDirectory "sdk/service/first" `
      -OutputDirectory $outputDirectory `
      -DirectoryToPublish $directoryToPublish
    New-APIViewArtifacts `
      -ServiceDirectory "sdk/service/second" `
      -OutputDirectory $outputDirectory `
      -DirectoryToPublish $directoryToPublish

    Join-Path $directoryToPublish "sdk/service/first/first.gosource" | Should -Exist
    Join-Path $directoryToPublish "sdk/service/second/second.gosource" | Should -Exist
    Join-Path $directoryToPublish "sdk/sdk" | Should -Not -Exist
    @(Get-ChildItem -Path $directoryToPublish -Filter *.gosource -Recurse).Count | Should -Be 2
  }
}
