## Why

The current security model adheres to a Zero Trust policy regarding cryptographic keys: the AI Agent and application code never access raw private keys. However, the current implementation only supports digital signatures. To fully secure the data lifecycle, the system must also support **encryption and decryption** operations (hardware-backed) without exposing keys to the software layer. This ensures that even if the AI or application server is compromised, the private keys remain secure within the hardware module (Secure Enclave).

## What Changes

-   **Update `Signer` Interface**: Rename or extend `Signer` to `CryptoProvider` (interface evolution) to include `Encrypt` and `Decrypt` methods.
-   **Update `RPCSigner`**: Implement the new methods to marshal requests into `encrypt.request` and `decrypt.request` RPC messages.
-   **Update Gateway Server**: Add handlers for the new RPC message types to route them to the connected Host App.
-   **BREAKING**: `Signer` interface expansion requires all implementations to support the new methods (or return "not implemented").

## Capabilities

### New Capabilities
*(None)*

### Modified Capabilities
- `secure-signer`: Update requirements to mandate support for hardware-backed Encrypt and Decrypt operations, not just Signing.

## Impact

-   **Internal Security Package**: `internal/security` interfaces and implementations will be expanded.
-   **Gateway**: `internal/gateway/server.go` RPC routing logic will manage new message types.
-   **Host App Contract**: The JSON-RPC protocol between Gateway and Host App will expand to include encryption/decryption payloads.
