## MODIFIED Requirements

### Requirement: Node constructor accepts SecretsStore
`NewNode()` SHALL accept an optional `*security.SecretsStore` parameter for encrypted node key management. When nil, file-based storage is used.

#### Scenario: Node created with SecretsStore
- **WHEN** `NewNode(cfg, logger, secrets)` is called with a non-nil SecretsStore
- **THEN** the node SHALL use SecretsStore for key storage

#### Scenario: Node created without SecretsStore
- **WHEN** `NewNode(cfg, logger, nil)` is called
- **THEN** the node SHALL fall back to file-based key storage in `cfg.KeyDir`
