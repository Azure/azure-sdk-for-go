# Snippet-Generator

The CLI tool `snippet-generator` enables you to embed your sample code from your test files to a markdown file.

## Usage

```shell
snippet-generator <path-to-discover> [true|false]
```

Here `path-to-discover` is the root directory where the tool discovers code snippets and the markdown files to replace. If you pass `true` to this tool, the tool will be in "strict mode", in this mode, the tool will require all the snippets defined in the test source files to be used up, otherwise it will report errors.

## Magic comments in go source files

### Snippet Region

Magic comments `// Snippet:Name` and `// EndSnippet` are used to mark a region as a snippet, for example:

```go
package example_test

import (
    "log"
)

// Snippet:Example
func something() {
	x := 1
	log.Println(x)
}
// EndSnippet
```

This will give you a snippet with the name of `Example` and get everything surrounded by the two lines of comments as its content.

### Snippet Ignore

In order to get the go source file compiles, sometimes you have to add some extra lines, you can surround those lines by `// SnippetIgnore` and `// EndSnippetIgnore` to exclude those lines from the content of a snippet.

Snippet Ignore can only be used inside a snippet, otherwise you will get an error.

```go
package example_test

// Snippet:Example
func something() {
	x := 1
	// Do something using x
	// SnippetIgnore
	// This is only a demonstration, but we have to use variable x, otherwise the program will not compile
	_ = x
	// EndSnippetIgnore
}
// EndSnippet
```

## Snippet declaration in markdown files

To reference your code block to some specific snippet, you just need write the code block, leave it empty and put the `Snippet:Name` after the declaration:

```go Snippet:Example
```

By running this tool, the code blocks in markdown files will be automatically fulfilled by the tool.
