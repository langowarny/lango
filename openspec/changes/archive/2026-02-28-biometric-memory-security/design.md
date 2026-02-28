## Context

The `BiometricProvider` in `internal/keyring/biometric_darwin.go` uses CGo to call
macOS Security framework APIs for Touch ID-protected Keychain access. The current
implementation correctly applies `kSecAccessControlBiometryAny` for access control but
does not zero sensitive plaintext buffers before freeing them. This leaves passphrase
data exposed in freed heap pages, vulnerable to memory dump or core dump analysis.

Current data flow in `Get()`:
1. C `malloc` + `memcpy` from CFData → C heap buffer
2. `C.GoStringN` copies to Go heap → Go string (immutable)
3. `C.free` releases C buffer without zeroing

Current data flow in `Set()`:
1. `C.CString(value)` allocates NUL-terminated C buffer
2. Passed to `keychain_set_biometric`
3. `C.free` releases without zeroing

## Goals / Non-Goals

**Goals:**
- Zero all C heap buffers containing plaintext before freeing
- Zero intermediate Go `[]byte` copies before they become unreachable
- Use `volatile` pointer pattern to prevent compiler optimization of zeroing
- Update documentation to reflect hardware-backed keyring terminology

**Non-Goals:**
- Converting Provider interface from `string` to `[]byte` (scope too large, touches all consumers)
- Adding `kSecAttrAccessGroup` (CLI binary has no bundle ID, making it ineffective)
- Using Secure Enclave directly (only supports EC keys, not symmetric passphrase storage)
- Eliminating the final Go `string` copy (Go strings are immutable by design)

## Decisions

### D1: Volatile-pointer zeroing in C (`secure_free`)

**Choice**: Custom `secure_free(char *ptr, int len)` using `volatile char *` loop.

**Alternatives considered**:
- `memset_s` (C11 Annex K): Not available on all macOS toolchains via CGo
- `explicit_bzero` (BSD): Not portable to CGo compilation context
- `SecureZeroMemory` (Windows): Platform-specific

**Rationale**: The `volatile` cast is the most portable pattern and is recommended by
CERT C (MSC06-C). The compiler cannot optimize away writes through a `volatile` pointer.

### D2: `[]byte` intermediate in Go `Get()`

**Choice**: Use `C.GoBytes` → `string()` → zero the `[]byte`, instead of `C.GoStringN`.

**Rationale**: `C.GoStringN` returns an immutable `string` that cannot be zeroed.
By going through `[]byte` first, we can zero the intermediate copy after extracting
the string. This reduces the window of plaintext exposure from two copies (C + Go string)
to one (Go string only, which is unavoidable without interface changes).

### D3: `memset` before `free` in `Set()`

**Choice**: Call `C.memset(ptr, 0, len+1)` before `C.free` on the `CString` buffer.

**Rationale**: The `CString` buffer holds the plaintext passphrase in C heap. While the
passphrase is also present as a Go string (caller-side), zeroing the C copy removes one
attack surface. Using `memset` here is acceptable since the buffer is freed immediately
after in the same `defer` — no optimization window for the compiler to skip it.

## Risks / Trade-offs

- **[Go string remains in memory]** → The final Go `string` returned from `Get()` cannot be
  zeroed due to Go's immutable string semantics. Mitigation: this is a known Go limitation;
  the C-side and `[]byte` zeroing still reduces attack surface significantly.
- **[volatile loop performance]** → Negligible; `secure_free` runs once per Keychain
  access (user-interactive operation, not hot path).
- **[GC may copy []byte before zeroing]** → Possible but unlikely in practice for
  short-lived buffers. Full mitigation would require `runtime.KeepAlive` or pinning,
  which is overkill for this threat model.
