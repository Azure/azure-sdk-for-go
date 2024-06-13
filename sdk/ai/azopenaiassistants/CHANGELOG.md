# Release History

## 0.2.0 (TBD)

### Features Added

- Now supports the Assistants V2 API, with support for vector stores.

### Bugs Fixed

### Breaking Changes

- Assistants V1 is no longer supported in this library. For information about how to migrate between V1 and V2, see OpenAI's migration documentation: [(link)](https://platform.openai.com/docs/assistants/migration).
- Types that were suffixed with `Options` have been changed, if their name would conflict with the options for a method. For example: `AssistantsThreadCreationOptions`, the main argument for `CreateThread()`, has been changed to `CreateThreadBody`.

## 0.1.1 (2024-05-07)

### Bugs Fixed

- ThreadRun.RequiredAction was deserialized incorrectly, making it impossible to actually resubmit a tool output. (PR#22834)

## 0.1.0 (2024-03-05)

* Initial release of the `azopenaiassistants` library
