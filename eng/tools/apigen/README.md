# apigen

apigen describes a Go module's public API in a declaration-only Markdown report for architectural review.

## Usage

Run from this directory, supplying the directory containing the target module's `go.mod`:

```powershell
go run . ..\..\..\sdk\azcore
go run . ..\..\..\sdk\azcore --output ..\..\..\sdk\azcore
```

Generation always writes `api.md` and `api.metadata.yml` to the current working directory by default. Use `--output` or `-o` to choose another output directory. Relative paths are resolved from the current working directory, and the output directory must already exist. The filenames are fixed. Stdout report output (`--output -`) is not supported.

File output also writes `api.metadata.yml` beside the Markdown file. Both files are generated and staged in temporary files before either destination is replaced. Each replacement occurs only after successful write and close; the two replacements are not a single filesystem transaction. New files on Unix-like systems are created with owner-only permissions (`0600` or more restrictive, depending on the process umask). Existing regular-file permissions are preserved; symlinks and other non-regular destinations are rejected. Errors go to stderr with a nonzero exit status. Ctrl+C cancels module loading.

The tool uses the current Go environment and build constraints. Use `go run . --help` for command-line help.

### API review bundles

Use `--output <directory>` to place the two files required by [API Review Hub](https://github.com/Azure/azure-sdk-tools/wiki/API-Review-Hub-Onboarding#5-emit-a-valid-review-bundle) in that directory. The metadata contains:

```yaml
apiHash: 607621e4437753200b71dcdcddaddcabeaf41944273a4600713b004471117189
packageVersion: 1.2.3
packageRelativePath: "sdk/example/example-package"
apigenVersion: 0.1.0
goVersion: 1.27.4
```

- `apiHash` is the lowercase hexadecimal SHA-256 digest of the exact Markdown file bytes, including newlines. No other files, metadata, or version values are included.
- `packageVersion` comes from a module-level string constant named `moduleVersion`, `serviceLibVersion`, or `version` (case-insensitive), in a Go file whose name contains `constant` or `version`, matching the version declaration names and filenames used by `Get-GoModuleVersionInfo` in `eng/scripts/Language-Settings.ps1`. The leading `v` is required in source and omitted from metadata; full semantic versions, including prerelease and build suffixes, are supported. Lookup uses the loaded module packages, including internal packages, but excludes tests, nested modules, dependencies, and files excluded by build constraints. Missing, invalid, or conflicting versions fail file generation without replacing existing reports; repeated identical versions are allowed.
- `packageRelativePath` is the module directory relative to the nearest ancestor containing a `.git` directory or worktree file, with `/` separators. Use `--repo-root <directory>` to override detection, including for source archives without Git metadata. The module must be inside that root; a module at the root uses `.`.
- `apigenVersion` is the generator version defined by `apigenVersion` in `metadata.go`.
- `goVersion` identifies the Go toolchain that compiled apigen, from `runtime.Version()` with the leading `go` removed (for example, `1.27.4`). It is not the target module's Go directive or the Go executable used to load that module. Development toolchain version strings are preserved.

Release checks can reuse the identical hash calculation without loading module source:

```powershell
go run . hash ..\..\..\sdk\azcore\api.md
```

This prints only the digest and a newline; it does not modify either file.

### Generate reports for all SDK modules

From this directory:

```powershell
go build .
.\Generate-AllApiReports.ps1
.\Generate-AllApiReports.ps1 -Concurrency 4
```

The script finds every `go.mod` under the repository's `sdk` directory, including nested modules, and generates or replaces `api.md` and `api.metadata.yml` in each module directory. It continues after generation failures and prints each failed module's absolute directory path in a final summary. Modules without a valid version constant (including non-shipping test modules) fail bundle generation. Exit code 1 means modules failed; exit code 0 includes an all-modules-successful message. Version, path, rendering, and staging failures leave existing reports unchanged.

By default, the script runs the executable beside it: `apigen.exe` on Windows and `apigen` on Linux and macOS. Use `-SdkRoot <directory>` to scan a smaller tree and `-ApigenPath <executable>` to select a different built generator. PowerShell 7 is required.

`-Concurrency <N>` limits the number of simultaneous module-generation tasks and defaults to `1` (sequential). Each task uses a separate apigen process and module output directory; Go supports concurrent use of its build and module caches. Failures do not stop other modules, diagnostics identify their module, and the final failure list is sorted by path. Each process can consume substantial CPU and memory, so increase concurrency conservatively.

## Tests

Run `go test ./...` from this directory. The batch-script tests require PowerShell 7 (`pwsh`) and skip if it is unavailable. The [CI definition](./ci.yml) runs the tests on Windows and Linux when this directory changes.

## Report format

- Each package has a Go code block containing its exported declarations. Its import block includes only packages referenced by public variable declarations, constant types, field types, parameters, results, embedded types, alias targets, and generic constraints, including those in expanded internal aliases. Implementation-only imports are omitted. Referenced types from other public packages are not expanded.
- Each package has five fixed declaration groups: constants and enums, variables, standalone functions, interfaces, and all other types. Each group is sorted alphabetically. Interface aliases belong to the interface group. Enum declarations appear with their constants and values helpers in the first group; other aliases share the final group with structs and other defined types. Associated functions and methods follow their type, sorted within each group.
- In structs, anonymous embedded fields appear first, followed by named fields; each group is sorted alphabetically. Interface members are sorted alphabetically. Fields retain their types, tags, and embedding. Exported embedded fields are shown by name rather than expanding their promoted members. For exported struct types, public members promoted through omitted private embeddings appear in comments after the declaration and its methods. Shadowed and ambiguous selectors are omitted. Each promoted method names the receiving type, as in `// Promoted method on Foo: Read() int` for a method in both method sets or `// Promoted method on *Foo: Reset()` for a pointer-only method. Multi-name struct fields are displayed individually. Source field order is intentionally omitted, so the report does not capture struct layout or positional-literal ordering.
- Constants retain their types and evaluated values, so `iota` and references to unexported constants do not depend on omitted declarations or source ordering.
- Signatures retain type parameters, constraints, and parameter and result names. Method receivers retain their types, including pointer or value distinctions, but omit receiver variable names to avoid diff noise.
- Variables retain any explicit type annotations and initializers, as displayed by pkg.go.dev, except that direct function-literal initializers are omitted and only the variable's explicit type or inferred function signature is shown. For other declarations with no explicit type, a `// Inferred type of Name: Type` comment records the compiler's type for each exported binding. Types from other packages use full import paths in these comments. Independently initialized names are listed separately. Shared multi-result initializers retain their binding order, with unexported names replaced by `_`.
- Aliases to internal types or unexported types in the same package are expanded under their public names, including exported fields and methods. An `Alias of` comment records the target's full package path and type name; the expanded declaration is review notation, not a new defined type.
- References to internal types in expanded declarations use public aliases when available: names in this package are unqualified, and names in other packages are qualified. When no public alias exists, the actual internal type reference is retained. Aliases to public types retain their `type Name = package.Type` declarations.
- Implementation bodies, unexported declarations, and descriptive comments are omitted. Whenever unexported struct fields are omitted, Go's `contains filtered or unexported fields` comment appears at the end of the struct body, after any exported fields. Empty structs and structs with only exported fields have no marker.
- The module section groups the module path, Go version, and a filtered `require` block using `go.mod` syntax. Requirements are deduplicated and sorted by module path, with their selected versions.
- The `require` block includes only external modules imported directly by non-test code anywhere in the module, including internal packages. Standard library imports, imports within this module, test-only dependencies, and dependencies used only transitively are excluded. Direct, non-test usage is determined from source imports, not solely from `go.mod` requirements or `// indirect` annotations.

## Alias normalization

When multiple public aliases identify the same internal type, references prefer an alias in the current package, then the first alias ordered by import path and name. An expanded type's self-references use its own public name.

Embedded struct types retain their field names and embedding. Public-alias selection for an embedded type therefore considers only aliases with the original field name; if none exists, the qualified source type is preserved instead of renaming the field.

Generic aliases and instantiated targets currently produce an explicit error when they require expansion. Unaliased internal types referenced by expanded declarations are preserved, including generic type references.

## Import collection

Imports are collected from displayed declarations, including variable initializers and evaluated constant types. Source aliases and dot imports are preserved; imports used only in omitted declarations, fields, function bodies, or constant expressions are excluded.
Imports used only in array length expressions are also excluded.

When separate source files use the same qualifier for different packages, collection currently returns an explicit error rather than emitting a conflicting import block.

## Declaration grouping

Methods belong to their public receiver type, including non-struct types. Variables always remain in their package-level group, regardless of their types.

String and integer types with exported typed constants are treated as enums, including public aliases such as `log.Event`. The displayed, normalized constant type selects its local enum declaration, so expanded aliases to internal enum types are supported. Group 1 emits unassociated constants first, then enum groups sorted by type name. Each enum group contains its type declaration immediately followed by its own `const (...)` block, then any constructors, values helpers, and methods. Types without exported constants, boolean and floating-point types, and constants whose types remain qualified external references do not form local enum groups.

Constructors named `New` or `NewX...` (with an uppercase letter after `New`) belong to a type only when exactly one result explicitly refers to that local type declaration, including through pointers, slices, or arrays. As an enum-specific exception, `Possible<Type>Values() []<Type>` functions belong to the corresponding enum: the name must match, there must be no parameters, and the sole result must be a slice of that enum. This includes aliases such as `PossibleAccessTierValues() []AccessTier`. Other functions, including `Get...` and `Parse...` helpers, remain at package level regardless of their result types. Generic types are matched by declaration identity, not by their type-parameter names.

An explicitly referenced public alias owns its constructors and matching enum values helper. A function returning an external type does not move when a local alias to that type is added; type identity alone does not establish an association.

## Markdown rendering

The renderer writes module metadata followed by package sections. The root section uses its package name; subpackage headings use module-relative paths so packages such as `policy` and `arm/policy` remain distinct. Without a root package, the report title uses the module path's final component, excluding its major-version suffix.

Constants and enum groups, variables, and standalone functions precede interfaces and other types. Enum declarations are emitted only in their constant groups, not repeated in the final type group. Each type is followed by its associated functions and methods. Nonempty constant and variable groups always use parenthesized blocks, including single declarations; empty groups are omitted.

Go declarations use fixed spacing rather than column alignment, so longer names do not realign neighboring declarations or fields. Nonempty structs and interfaces use multiline bodies. Blank lines inherited from omitted source comments or fields are removed. Multiline literal contents are preserved. Numeric constants use exact Go expressions rather than approximate display values.
Untyped rune constants retain their rune kind, even when their evaluated value cannot be written as a single rune literal.
Array lengths in declaration types are rendered as evaluated decimal integers, so changes to a private constant affect the report only if they change the array length.

An absent Go directive or empty dependency list is omitted from the module block. A dependency without a selected version is an error. The complete document is rendered before output begins, so rendering failures do not leave partial reports; output errors are returned to the caller.

### Ordering and review diffs

Sort keys use case-sensitive lexical ordering, independent of source file names and declaration positions. Packages are ordered by import path; dependencies by module path; imports by standard-library group, path, and alias; and declarations within each group by name. Struct fields are grouped by embedding first, then sorted by field name within each group; unnamed interface members use their type expressions. Shared multi-result variable declarations use their first exported binding as the sort key and retain their binding order. Parameters, results, and type parameters retain their binding order.

Regression tests compare complete Markdown output after file, declaration, field, and interface-member reordering, reversed package traversal, and implementation-only changes. They also verify that adding functions, methods, types, longer-named constants, variables, and fields inserts them at their sorted positions without changing existing declarations, including singleton-to-multiple-entry changes in enum blocks.

Public alias additions can still change the preferred spelling of existing internal type references according to the alias normalization rules above, but do not regroup existing functions based only on type identity.
