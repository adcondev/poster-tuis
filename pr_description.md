🧪 Add tests for validateServiceVariantFields

🎯 **What:** The `validateServiceVariantFields` function in `internal/service/service.go` was previously lacking unit tests, which left a gap in coverage for an important part of validating service properties.

📊 **Coverage:** The tests cover the following scenarios:
*   A valid variant object.
*   Invalid characters and empty value in `RegistryName`.
*   Invalid characters, path traversal, slashes, backslashes, and empty value in `ExeName`.
*   Invalid characters (including shell metacharacters and newlines) and empty value in `DisplayName`.

✨ **Result:** Test coverage improved by ensuring that the core logic validating service variants correctly identifies errors before registering or acting on a service variant.
