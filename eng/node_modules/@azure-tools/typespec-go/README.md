# @azure-tools/typespec-go

TypeSpec emitter for Go SDKs

## Install

```bash
npm install @azure-tools/typespec-go
```

## Usage

1. Via the command line

```bash
tsp compile . --emit=@azure-tools/typespec-go
```

2. Via the config

```yaml
emit:
  - '@azure-tools/typespec-go'
```

The config can be extended with options as follows:

```yaml
emit:
  - '@azure-tools/typespec-go'
options:
  '@azure-tools/typespec-go':
    option: value
```

## Emitter options

### `emitter-output-dir`

**Type:** `absolutePath`

Defines the emitter output directory. Defaults to `{output-dir}/@azure-tools/typespec-go`
See [Configuring output directory for more info](https://typespec.io/docs/handbook/configuration/configuration/#configuring-output-directory)

### `azcore-version`

**Type:** `string`

Semantic version of azcore without the leading 'v' to use if different from the default version (e.g. 1.2.3).

### `disallow-unknown-fields`

**Type:** `boolean`

When true, unmarshalers will return an error when an unknown field is encountered in the payload. The default is false.

### `file-prefix`

**Type:** `string`

Optional prefix to file names. For example, if you set your file prefix to "zzz*", all generated code files will begin with "zzz*".

### `generate-fakes`

**Type:** `boolean`

When true, enables generation of fake servers. The default is false.

### `head-as-boolean`

**Type:** `boolean`

When true, HEAD requests will return a boolean value based on the HTTP status code. The default is false.

### `inject-spans`

**Type:** `boolean`

Enables generation of spans for distributed tracing. The default is false.

### `module`

**Type:** `string`

The name of the Go module written to go.mod. Omit to skip go.mod generation. When module is specified, module-version must also be specified.

### `module-version`

**Type:** `string`

Semantic version of the Go module without the leading 'v' written to constants.go. (e.g. 1.2.3). When module-version is specified, module must also be specified.

### `rawjson-as-bytes`

**Type:** `boolean`

When true, properties that are untyped (i.e. raw JSON) are exposed as []byte instead of any or map[string]any. The default is false.

### `slice-elements-byval`

**Type:** `boolean`

When true, slice elements will not be pointer-to-type. The default is false.

### `single-client`

**Type:** `boolean`

Indicates package has a single client. This will omit the Client prefix from options and response types. If multiple clients are detected, an error is returned. The default is false.

### `stutter`

**Type:** `string`

Uses the specified value to remove stuttering from types and funcs instead of the built-in algorithm.

### `fix-const-stuttering`

**Type:** `boolean`

When true, fix stuttering for `const` types and values. The default is false.

### `generate-examples`

**Type:** `boolean`

Deprecated. Use generate-samples instead.

### `generate-samples`

**Type:** `boolean`

When true, generate example tests. The default is false.

### `factory-gather-all-params`

**Type:** `boolean`

When true, the `NewClientFactory` constructor gathers all parameters. When false, it only gathers common parameters of clients. The default is true.
