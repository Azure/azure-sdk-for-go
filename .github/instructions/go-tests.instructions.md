---
applyTo: '**/*_test.go'
---

- Use github.com/stretchr/testify/require for assertions. Some commonly used functions: `require.Equal`, `require.NoError`.
- Environment variables required for live testing can be found by looking for recording.Getenv() calls, or os.Getenv() calls in the code. You should place these into a .env file at the root of the module, before testing.
