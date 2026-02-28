## ADDED Requirements

### Requirement: BiometricProvider SHALL zero C heap buffers before freeing
The `BiometricProvider` SHALL zero all C heap buffers containing plaintext secrets before calling `free()`. Zeroing MUST use a volatile pointer pattern to prevent compiler optimization from eliding the memory wipe.

#### Scenario: Get zeroes C buffer via secure_free
- **WHEN** `BiometricProvider.Get()` retrieves a secret from the Keychain
- **THEN** the C heap buffer SHALL be zeroed via `secure_free()` (volatile pointer loop + free) before control returns to Go

#### Scenario: Set zeroes CString buffer before freeing
- **WHEN** `BiometricProvider.Set()` stores a secret in the Keychain
- **THEN** the `C.CString` buffer containing the plaintext value SHALL be zeroed with `memset` before `free` is called

### Requirement: BiometricProvider SHALL zero intermediate Go byte slices
The `BiometricProvider.Get()` method SHALL copy Keychain data into a Go `[]byte` via `C.GoBytes`, extract the string, and then zero every byte of the `[]byte` slice before it becomes unreachable.

#### Scenario: Get zeroes Go byte slice after string extraction
- **WHEN** `BiometricProvider.Get()` copies data from C heap to Go heap
- **THEN** it SHALL use `C.GoBytes` (not `C.GoStringN`), extract the string via `string(data)`, and zero the `[]byte` with a range loop

### Requirement: secure_free C helper prevents compiler optimization
The C `secure_free` helper function SHALL cast the pointer to `volatile char *` before zeroing to prevent the compiler from optimizing away the memset as a dead store.

#### Scenario: Volatile pointer prevents optimization
- **WHEN** `secure_free(ptr, len)` is called
- **THEN** it SHALL iterate through the buffer using a `volatile char *` pointer, set each byte to zero, and then call `free(ptr)`

#### Scenario: Null pointer safety
- **WHEN** `secure_free(NULL, 0)` is called
- **THEN** it SHALL return without error (NULL guard)
