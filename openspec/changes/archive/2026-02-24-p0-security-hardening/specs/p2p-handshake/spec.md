## MODIFIED Requirements

### Requirement: Signature verification
The handshake verifier SHALL perform full ECDSA secp256k1 signature verification by recovering the public key from the signature and comparing it with the claimed public key, instead of accepting any non-empty signature.

#### Scenario: Valid signature accepted
- **WHEN** a challenge response contains a 65-byte ECDSA signature that recovers to a public key matching `resp.PublicKey`
- **THEN** the verifier SHALL accept the response as authenticated

#### Scenario: Invalid signature rejected
- **WHEN** a challenge response contains a signature that recovers to a public key NOT matching `resp.PublicKey`
- **THEN** the verifier SHALL reject the response with "signature public key mismatch" error

#### Scenario: Wrong signature length rejected
- **WHEN** a challenge response contains a signature that is not exactly 65 bytes
- **THEN** the verifier SHALL reject the response with "invalid signature length" error

#### Scenario: Corrupted signature rejected
- **WHEN** a challenge response contains a 65-byte signature that cannot be recovered to a valid public key
- **THEN** the verifier SHALL reject the response with an error

#### Scenario: No proof or signature rejected
- **WHEN** a challenge response contains neither a ZK proof nor a signature
- **THEN** the verifier SHALL reject the response with "no proof or signature in response" error

### Requirement: Constant-time nonce comparison
The handshake verifier SHALL use `hmac.Equal()` for nonce comparison to prevent timing side-channel attacks.

#### Scenario: Nonce mismatch detected securely
- **WHEN** the response nonce does not match the challenge nonce
- **THEN** the verifier SHALL reject the response with "nonce mismatch" error using constant-time comparison
