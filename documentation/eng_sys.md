# Engineering System Checks

* [Build and Test](#build-and-test)
* [Analyze Stages](#analyze-stages)

## Build and Test

Our build system runs PR changes against the latest two versions of Go on both Windows and Linux.

## Analyze

### Link Verification Check
Verifies all of the links are valid in your README files. This step also checks that locale codes in links are removed. If this is failing first check if you have locale codes (ie. `en-us`) in your links, then check to see if the link works locally.

If you are trying to add a link that will exist in the next PR (ie. you are adding a samples README or migration guide), you can use an `aka.ms` link or use a temporary link (ie: `https://microsoft.com`) and create a follow-up PR to correct the temporary link.

### Lint
Some of the most common linting errors are:
* `errcheck`: An error was returned but it was not checked to be `nil`
* `varcheck`: A variable is unused
* `deadcode`: A struct or method is unused
* `ineffasign`: An ineffectual assignment, the variable is not used after declaration.

For more information about the linters run checkout the [golangci website][golangci_website]

To run this locally, first install the tool with:
```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1
```

```bash
golangci-lint run -c <path_to_root>/eng/.golangci.yml in <path_to_my_package>
```

### Copyright Header Check
Every source file must have the MIT header comment at the top of the file. At the top of each file you need to include the following snippet and a new line before the package definition:
```golang
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package <mypackage>
```

### Format Check
Your package should follow the default formatting, which you can run locally with the command:
```bash
go fmt
```

<!-- LINKS -->
[golangci_website]: https://golangci-lint.run/usage/linters/