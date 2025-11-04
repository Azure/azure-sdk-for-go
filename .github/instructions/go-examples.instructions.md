---
applyTo: '**/example*_test.go'
---

- Error handling code should look like this (include all comments):
  ```go
  if err != nil {
    // TODO: Update the following line with your application specific error handling logic
    log.Fatalf(\"ERROR: %s\", err)
  }
  ```
- Don't make examples executable (ie, don't include `//Output:`).
- Use `context.TODO()` in places that require a context, NOT `context.Background()`. This lets the user know they need to provide a context.
