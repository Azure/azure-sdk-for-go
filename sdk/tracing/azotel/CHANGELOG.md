# Release History

## 0.4.1 (Unreleased)

### Features Added

### Breaking Changes

### Bugs Fixed

### Other Changes

* Updated dependencies.
* Added OTel implementation for `Provider.NewPropagator()` to return a `propagation.TraceContext()` propagator for context propagation.
* Added OTel implementation for `tracing.LinkFromContext()` to return a OTel Span Link with the passed in attributes attached to it.
* Added OTel implementation for `Span.AddLink()` to add an OTel Span Link to the current span.
* Added OTel implementation for `Span.SpanContext()` to return the OTel Span Context of the current span.

## 0.4.0 (2023-11-07)

### Other Changes

* Updated to latest release of `azcore` and cleaned up example.

## 0.3.0 (2023-10-16)

### Other Changes

* Updated to latest beta of `azcore`.

## 0.2.0 (2023-07-13)

### Breaking Changes

* The type for parameter `tracerProvider` in function `NewTracingProvider()` has changed to `trace.TracerProvider`.

## 0.1.0 (2023-06-06)

### Features Added

* Initial release
