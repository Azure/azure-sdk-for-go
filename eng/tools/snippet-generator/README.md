# Snippet-Generator

The CLI tool `snippet-generator` enables you to embed your sample code from your test files to a markdown file.

## Usage

```shell
snippet-generator <path-to-discover> [true|false]
```

Here `path-to-discover` is the root directory where the tool discovers code snippets and the markdown files to replace. If you pass `true` to this tool, the tool will be in "strict mode", in this mode, the tool will require all the snippets defined in the test source files to be used up, otherwise it will report errors.

## Magic comments


