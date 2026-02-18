## Context

The system currently employs a `Signer` interface and `RPCSigner` implementation to support hardware-backed digital signatures via a connection to a macOS Host App (which accesses the Secure Enclave). The user requirement is to extend this zero-trust model to include Encryption and Decryption, ensuring the AI/Application never accesses the private keys directly.

## Goals / Non-Goals

**Goals:**
*   Extend the current signing infrastructure to support encryption and decryption.
*   Maintain the Zero Trust architecture: keys remain non-exportable in the hardware module.
*   Ensure backward compatibility or clean interface evolution for existing signing users.

**Non-Goals:**
*   Software-based key implementation (keys must be hardware-backed).
*   Key management/rotation APIs (out of scope, handled by Host App natively).

## Decisions

### 1. Renaming `Signer` to `CryptoProvider`
**Decision:** Rename the `Signer` interface to `CryptoProvider` or similar, to reflect its broader scope.
**Rationale:** `Signer` implies only signing. Since we are adding `Encrypt` and `Decrypt`, the name should facilitate this.
**Migration:** We will rename `internal/security/signer.go`'s interface and update all references.

### 2. Interface Definition
```go
type CryptoProvider interface {
    Sign(ctx context.Context, keyID string, payload []byte) ([]byte, error)
    Encrypt(ctx context.Context, keyID string, plaintext []byte) ([]byte, error)
    Decrypt(ctx context.Context, keyID string, ciphertext []byte) ([]byte, error)
}
```

### 3. RPC Protocol Schema
We will add new method types to the existing JSON-RPC structure:
- `sign.request` / `sign.response` (Existing)
- `encrypt.request` / `encrypt.response` (New)
- `decrypt.request` / `decrypt.response` (New)

Each request will carry: `id`, `keyId`, and `payload` (plaintext or ciphertext).

## Risks / Trade-offs

*   **Risk**: Host App Compatibility.
    *   **Mitigation**: The Gateway must handle cases where the connected Host App strictly supports only Signing (older version). The `RPCSigner` should return an explicit "Method not supported" or timeout if the Host doesn't reply to Encrypt/Decrypt requests.
*   **Trade-off**: Latency.
    *   **Detail**: RPC calls to the Host App introduce network/IPC latency compared to local OpenSSL.
    *   **Acceptance**: This is an acceptable trade-off for the security guarantee of non-exportable keys.
