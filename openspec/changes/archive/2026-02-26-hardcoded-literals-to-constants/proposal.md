## Why

P2P/A2A/Payment packages accumulated 100+ hardcoded string and integer literals across ~25 files, undermining type safety, making refactoring error-prone, and violating the project's Go style guide. Replacing them with typed enums, named constants, and sentinel errors improves compile-time safety, maintainability, and consistency.

## What Changes

- Add `ChainID` type with named constants for Ethereum/Base/Sepolia chain IDs in the wallet package
- Add `CurrencyUSDC` constant and replace all 20+ occurrences of `"USDC"` across the codebase
- Export `WalletKeyName` constant and use it in x402 signer instead of duplicated string
- Add gas fee constants (`DefaultBaseFeeWei`, `DefaultMaxPriorityFeeWei`, `BaseFeeMultiplier`, `EthAddressLength`) and `BalanceOfSelector` in the payment package
- Add `ResponseStatus` enum type (`ok`/`error`/`denied`/`payment_required`) and 7 sentinel errors in the P2P protocol package, replacing ~40 hardcoded status strings
- Add `ACLAction` enum (`allow`/`deny`), `WildcardAll` constant, and 4 sentinel errors in the firewall package
- Add `ProofScheme` enum (`plonk`/`groth16`) and `SRSMode` enum (`unsafe`/`file`) with `ErrUnsupportedScheme` in the ZKP package
- Add scoring weight constants (`FailureWeight`, `TimeoutWeight`, `BasePenalty`) in the reputation package
- Add `KMSProviderName` enum (`aws-kms`/`gcp-kms`/`azure-kv`/`pkcs11`) in the security package
- Add route/header/tag constants in the A2A server package
- Add `DefaultQuoteExpiry` constant in the paygate package
- Add display truncation constants in CLI payment history

## Capabilities

### New Capabilities

### Modified Capabilities
- `sentinel-errors`: Extended with new sentinel errors in protocol, firewall packages
- `enum-validation`: Extended with new enum types (ResponseStatus, ACLAction, ProofScheme, SRSMode, KMSProviderName, ChainID)

## Impact

- **25 files changed** across wallet, payment, protocol, firewall, zkp, reputation, security, a2a, paygate, x402, app, cli, and tools packages
- `Response.Status` field type changed from `string` to `ResponseStatus` — JSON wire format unchanged (underlying type is `string`)
- `ACLRule.Action` field type changed from `string` to `ACLAction` — JSON wire format unchanged
- `NewKMSProvider` signature changed from `string` to `KMSProviderName` — callsites updated with type casts
- ZKP Config/Proof types changed to use `ProofScheme`/`SRSMode` — callsites updated with type casts
- No breaking changes to external APIs or wire protocols
