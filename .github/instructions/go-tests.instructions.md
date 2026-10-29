---
applyTo: '**/*_test.go'
---

- Use github.com/stretchr/testify for assertions rather than hand-rolled `if got != want { t.Errorf(...) }` checks or helpers that reimplement basic assertions testify already provides (a custom `assertEqual`, a home-grown assertion package). Some commonly used functions: `require.Equal`, `require.NoError`, `require.NotNil`.
- Domain-specific helpers that *compose* testify are encouraged, not discouraged — for example a `requireEqualAttributes(t, a, b)` that asserts the fields of a model, or a `requireEqualCloudEvent` that normalizes recorded values before comparing. Mark them with `t.Helper()` so failures report the caller's line.
- Choose between `require` and `assert` based on what follows the assertion:
  - Use `require` when the rest of the test depends on the predicate holding, so a failure stops the test instead of panicking or cascading. Typical cases are checking an error before using the result, checking for nil before a dereference, and checking length before indexing.
  - Use `assert` for independent checks that should each be reported in a single run, for example verifying several unrelated fields on a response.

  ```go
  require.NoError(t, err)
  require.NotNil(t, resp.Value)
  assert.Equal(t, "value", *resp.Value)
  assert.Equal(t, 200, resp.StatusCode)
  ```
- In tests that record and play back, use `recording.Sleep` rather than `time.Sleep` when waiting for eventual consistency or provisioning. `recording.Sleep` sleeps while recording and is a no-op in `PlaybackMode`, so playback runs don't pay a delay that only matters against the live service. Live-only tests, which never play back, may use `time.Sleep` directly.
- Environment variables required for live testing can be found by looking for recording.Getenv() calls, or os.Getenv() calls in the code. You should place these into a .env file at the root of the module, before testing.
