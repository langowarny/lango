## 1. Core Implementation

- [x] 1.1 Change om_provider field from InputText to InputSelect in NewObservationalMemoryForm, using buildProviderOptions(cfg) with empty string prepended
- [x] 1.2 Remove Placeholder field from om_provider (not applicable for InputSelect)

## 2. Testing

- [x] 2.1 Add TestNewObservationalMemoryForm_ProviderIsSelect test verifying om_provider is InputSelect with empty string as first option
- [x] 2.2 Verify om_model remains InputText in the test

## 3. Verification

- [x] 3.1 Run go build ./... to confirm no compilation errors
- [x] 3.2 Run go test ./internal/cli/onboard/... to confirm all tests pass
