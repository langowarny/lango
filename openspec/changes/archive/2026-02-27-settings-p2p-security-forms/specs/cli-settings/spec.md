## MODIFIED Requirements

### Requirement: Configuration Coverage
The settings editor SHALL support editing all configuration sections including P2P networking and advanced security:
21. **P2P Network** -- Enabled, listen addrs, bootstrap peers, relay, mDNS, max peers, handshake timeout, session token TTL, auto-approve, gossip interval, ZK handshake/attestation, signed challenge, min trust score
22. **P2P ZKP** -- Proof cache dir, proving scheme, SRS mode/path, max credential age
23. **P2P Pricing** -- Enabled, per query price, tool-specific prices
24. **P2P Owner Protection** -- Owner name/email/phone, extra terms, block conversations
25. **P2P Sandbox** -- Tool isolation (enabled, timeout, memory), container sandbox (runtime, image, network, rootfs, CPU, pool)
26. **Security Keyring** -- OS keyring enabled
27. **Security DB Encryption** -- SQLCipher enabled, cipher page size
28. **Security KMS** -- Region, key ID, endpoint, fallback, timeout, retries, Azure vault/version, PKCS#11 module/slot/PIN/key label

#### Scenario: Menu categories
- **WHEN** user launches `lango settings`
- **THEN** the menu SHALL display all categories including P2P Network, P2P ZKP, P2P Pricing, P2P Owner Protection, P2P Sandbox, Security Keyring, Security DB Encryption, Security KMS, grouped under "P2P Network" and "Security" sections

### Requirement: Security form signer provider options
The Security form's signer provider dropdown SHALL include options for all supported providers: local, rpc, enclave, aws-kms, gcp-kms, azure-kv, pkcs11.

#### Scenario: KMS providers available in signer dropdown
- **WHEN** user opens the Security form
- **THEN** the signer provider dropdown SHALL include "aws-kms", "gcp-kms", "azure-kv", and "pkcs11" as options

## ADDED Requirements

### Requirement: P2P Network settings form
The settings TUI SHALL provide a "P2P Network" form with 14 fields covering core P2P networking: enabled, listen addresses, bootstrap peers, relay, mDNS, max peers, handshake timeout, session token TTL, auto-approve known peers, gossip interval, ZK handshake, ZK attestation, require signed challenge, and min trust score.

#### Scenario: User enables P2P networking
- **WHEN** user navigates to "P2P Network" and sets Enabled to true
- **THEN** the config's `p2p.enabled` field SHALL be set to true upon save

#### Scenario: User sets listen addresses
- **WHEN** user enters comma-separated multiaddrs in "Listen Addresses"
- **THEN** the config's `p2p.listenAddrs` SHALL contain each address as a separate array element

### Requirement: P2P ZKP settings form
The settings TUI SHALL provide a "P2P ZKP" form with fields for proof cache directory, proving scheme (plonk/groth16), SRS mode (unsafe/file), SRS path, and max credential age.

#### Scenario: User selects groth16 proving scheme
- **WHEN** user selects "groth16" from the proving scheme dropdown
- **THEN** the config's `p2p.zkp.provingScheme` SHALL be set to "groth16"

### Requirement: P2P Pricing settings form
The settings TUI SHALL provide a "P2P Pricing" form with fields for enabled, price per query, and tool-specific prices (as key:value comma-separated text).

#### Scenario: User sets tool prices
- **WHEN** user enters "exec:0.10,browser:0.50" in the Tool Prices field
- **THEN** the config's `p2p.pricing.toolPrices` SHALL be a map with keys "exec" and "browser"

### Requirement: P2P Owner Protection settings form
The settings TUI SHALL provide a "P2P Owner Protection" form with fields for owner name, email, phone, extra terms, and block conversations. The block conversations field SHALL default to checked when the config value is nil.

#### Scenario: User sets block conversations with nil default
- **WHEN** the config's `blockConversations` is nil
- **THEN** the form SHALL display the checkbox as checked (default true)

#### Scenario: User unchecks block conversations
- **WHEN** user unchecks "Block Conversations"
- **THEN** the config's `p2p.ownerProtection.blockConversations` SHALL be a pointer to false

### Requirement: P2P Sandbox settings form
The settings TUI SHALL provide a "P2P Sandbox" form with fields for tool isolation (enabled, timeout, max memory) and container sandbox (enabled, runtime, image, network mode, read-only rootfs, CPU quota, pool size, pool idle timeout). Container-specific fields SHALL only be visible when Container Sandbox is enabled.

#### Scenario: User configures container sandbox
- **WHEN** user enables container sandbox and selects "docker" runtime
- **THEN** the config's `p2p.toolIsolation.container.enabled` SHALL be true and `runtime` SHALL be "docker"

#### Scenario: Container read-only rootfs defaults to true
- **WHEN** the config's `readOnlyRootfs` is nil
- **THEN** the form SHALL display the checkbox as checked (default true)

### Requirement: Security Keyring settings form
The settings TUI SHALL provide a "Security Keyring" form with a single field for OS keyring enabled/disabled.

#### Scenario: User enables keyring
- **WHEN** user checks "OS Keyring Enabled"
- **THEN** the config's `security.keyring.enabled` SHALL be set to true

### Requirement: Security DB Encryption settings form
The settings TUI SHALL provide a "Security DB Encryption" form with fields for SQLCipher encryption enabled and cipher page size.

#### Scenario: User enables DB encryption
- **WHEN** user checks "SQLCipher Encryption" and sets page size to 4096
- **THEN** the config SHALL have `security.dbEncryption.enabled` true and `cipherPageSize` 4096

#### Scenario: Cipher page size validation
- **WHEN** user enters 0 or a negative number for cipher page size
- **THEN** the form SHALL display a validation error "must be a positive integer"

### Requirement: Security KMS settings form
The settings TUI SHALL provide a "Security KMS" form with conditional field visibility based on the selected backend. Cloud KMS fields (region, endpoint) appear for aws-kms/gcp-kms/azure-kv. Azure-specific fields appear for azure-kv. PKCS#11 fields appear for pkcs11. Common fields (key ID, fallback, timeout, retries) appear for all non-local backends.

#### Scenario: User configures AWS KMS
- **WHEN** user selects "aws-kms" and enters region and key ARN
- **THEN** the config's `security.kms.region` and `security.kms.keyId` SHALL contain the entered values

#### Scenario: PKCS#11 PIN is password field
- **WHEN** the KMS form is displayed with pkcs11 backend selected
- **THEN** the PKCS#11 PIN field SHALL use InputPassword type to mask the value

#### Scenario: Local backend hides KMS fields
- **WHEN** user selects "local" as the KMS backend
- **THEN** all KMS-specific fields SHALL be hidden
