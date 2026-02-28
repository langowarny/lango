## 1. Menu & Routing

- [x] 1.1 Add "P2P Network" section with 5 categories (p2p, p2p_zkp, p2p_pricing, p2p_owner, p2p_sandbox) to `NewMenuModel()` in `menu.go`
- [x] 1.2 Add "Security" section with 3 new categories (security_keyring, security_db, security_kms) alongside existing security/auth entries in `menu.go`
- [x] 1.3 Add 8 new `case` entries in `handleMenuSelection()` in `editor.go`

## 2. Form Builders

- [x] 2.1 Add `derefBool(ptr *bool, defaultVal bool)` helper for `*bool` config fields
- [x] 2.2 Add `formatKeyValueMap(m map[string]string)` helper for map-to-string conversion
- [x] 2.3 Add `NewP2PForm()` -- 14 fields (enabled, listen addrs, bootstrap peers, relay, mDNS, max peers, handshake timeout, session token TTL, auto-approve, gossip interval, ZK handshake, ZK attestation, signed challenge, min trust score)
- [x] 2.4 Add `NewP2PZKPForm()` -- 5 fields (proof cache dir, proving scheme, SRS mode, SRS path, max credential age)
- [x] 2.5 Add `NewP2PPricingForm()` -- 3 fields (enabled, per query price, tool prices)
- [x] 2.6 Add `NewP2POwnerProtectionForm()` -- 5 fields (owner name, email, phone, extra terms, block conversations)
- [x] 2.7 Add `NewP2PSandboxForm()` -- 11 fields with `VisibleWhen` for container sub-fields (tool isolation enabled/timeout/memory, container enabled/runtime/image/network/rootfs/CPU/pool size/pool idle timeout)
- [x] 2.8 Add `NewKeyringForm()` -- 1 field (OS keyring enabled)
- [x] 2.9 Add `NewDBEncryptionForm()` -- 2 fields (SQLCipher enabled, cipher page size)
- [x] 2.10 Add `NewKMSForm()` -- 12 fields with `VisibleWhen` for backend-specific fields (backend, region, key ID, endpoint, fallback, timeout, retries, Azure vault/version, PKCS#11 module/slot/PIN/key label)

## 3. Config Write-back

- [x] 3.1 Add `boolPtr(val bool) *bool` helper in `state_update.go`
- [x] 3.2 Add P2P Network case entries (~14) in `UpdateConfigFromForm()`
- [x] 3.3 Add P2P ZKP case entries (~5) in `UpdateConfigFromForm()`
- [x] 3.4 Add P2P Pricing case entries (~3) in `UpdateConfigFromForm()`
- [x] 3.5 Add P2P Owner Protection case entries (~5) in `UpdateConfigFromForm()`
- [x] 3.6 Add P2P Sandbox case entries (~11) in `UpdateConfigFromForm()`
- [x] 3.7 Add Security Keyring case entry in `UpdateConfigFromForm()`
- [x] 3.8 Add Security DB Encryption case entries (~2) in `UpdateConfigFromForm()`
- [x] 3.9 Add Security KMS case entries (~12) in `UpdateConfigFromForm()`

## 4. Existing Form Update

- [x] 4.1 Expand signer provider options in `NewSecurityForm()` to include aws-kms, gcp-kms, azure-kv, pkcs11

## 5. Tests

- [x] 5.1 Add form field count and key tests for all 8 new forms
- [x] 5.2 Add menu category existence tests for all 8 new categories
- [x] 5.3 Add config round-trip tests for P2P, Sandbox *bool, and KMS fields
- [x] 5.4 Add `derefBool` helper test

## 6. Verification

- [x] 6.1 Run `go build ./...` -- no errors
- [x] 6.2 Run `go test ./internal/cli/settings/...` -- all pass
