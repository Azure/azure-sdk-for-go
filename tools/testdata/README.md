# Test cases for versioner tool

| Test case | Additional changes | Breaking changes | Existing versions | go.mod file existence | Preview |
| :---: | :---: | :---: | :---: | :---: | :---: |
| a | Yes | No | v1 | Yes | No |
| b | Yes | Yes | v1 | Yes | No |
| c | No | No | v1 | Yes | No |
| d | Yes | Yes | v1, v2 | Yes | No |
| e | Yes | No | v1, v2 | Yes | No |
| f | Yes | No | None | Yes | No |
| g | Yes | Yes | v1 | Yes | No |
| h | Yes | Yes | v1 | Yes | No |
| i | No | No | v1 | Yes | No |
| j | No | No | v1, v2 | Yes | No |
| k | Yes | No | v0 | Yes | Yes |
| l | Yes | Yes | v0 | Yes | Yes |
| m | Yes | No | None | Yes | Yes |
| n | No | No | v0 | Yes | Yes |
| o | No | No | v0 | Yes | Yes |

* For i, all files in v1 and stage are identical with each other. In this case, no version should be bumped.
* For j, all files in v2 and stage are identical with each other. In this case, no version should be bumped.
* For k and l, they are preview packages. For preview packages, version numbers should start from v0.0.0 and they do not bump major versions even received breaking changes.
* For n, it is a test case for patch version of preview package.
* For o, all files are identical with the files in stage directory. In this case, no version should be bumped.

