# Releasing Packages

After going through a minimal architects board review and preparing your package for release, verify you are ready for release by following the "Release Checklist", and finally release your package by following the "Release Process"

## Release Checklist

- [] Verify there are no replace directives in the go.mod file
- [] Verify the package has a LICENSE file
- [] Verify documentation is present and accurate for all public methods and types. Reference the [content guidelines](https://review.docs.microsoft.com/help/contribute-ref/contribute-ref-how-to-document-sdk?branch=master#api-reference) for best practices. You can start the `godoc` server by running `godoc -http=:6060` from the module home and navigating to `localhost:6060` in the browser.
- [] Verify there are no broken links
- [] Verify all links are non-localized (no "en-us" in links)
- [] Check the package manager link goes to the correct package
- [] Verify Samples
- [] Verify samples are visible in the [sample browser](https://docs.microsoft.com/samples/browse/)
- [] Verify release notes follow [general guidelines](https://azure.github.io/azure-sdk/policies_releasenotes.html)
- [] Verify troubleshooting section of README contains information about how to enable logging
- [] Verify CHANGELOG follows [current guidance](https://azure.github.io/azure-sdk/policies_releases.html#changelog-guidance)
- [] Verify all champion scenarios have a getting started scenario

## Release Process

1. Complete all steps of the Release Checklist shown above
2. Mark the package as 'in-release' by running the `./eng/common/scripts/Prepare-Release.ps1` script and following the prompts. The script may update the version and/or `CHANGELOG.md` of the package. If changes are made, these changes need to be committed and merged before continuing with the release process.
3. Run the pipeline from the `internal` Azure Devops. This will require you to approve the release after both the live and recorded test pipelines pass.
4. Validate the package was released properly by running `go get <your-package>@<your-version>` (ie. `go get github.com/Azure/azure-sdk-for-go/sdk/azcore@v0.20.0`) and validating that pkg.go.dev has updated with the latest version.