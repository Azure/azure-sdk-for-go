# Test cases for versioner tool

| Test case | Additional changes | Breaking changes | Existing versions | go.mod file existence |
| :--- | :--- | :--- | :--- | :--- |
| a | Yes | No | v1 | Yes |
| b | Yes | Yes | v1 | Yes |
| c | No | No | v1 | Yes |
| d | Yes | Yes | v1, v2 | Yes |
| e | Yes | No | v1, v2 | Yes |
| f | Yes | No | None | Yes |
| g | Yes | Yes | v1 | Yes |
| h | Yes | Yes | v1 | Yes |
| i | No | No | v1 | Yes |

* For i, checksum of every file is the same. In this case, no version should be bumped.