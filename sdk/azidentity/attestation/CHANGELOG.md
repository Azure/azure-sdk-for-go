# Release History

## 0.1.0-beta.1 (Unreleased)

### Features Added

- Blank-import activation for automatic KeyGuard-attested managed identity mTLS PoP requests.
  Initialization only registers availability; capability discovery and native operations remain
  lazy. Ordinary bearer requests and known unsupported hosts retain bearer acquisition.
- Added compiled examples for default Key Vault construction and per-MIC opt-out, and
  offline imported-binary coverage of automatic Key Vault requests through MIC and DAC.
