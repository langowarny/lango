## Requirements

### Requirement: ProviderConfig type strengthening
The `ProviderConfig.Type` field SHALL use `types.ProviderType` instead of raw `string`.

#### Scenario: Config deserialization with typed provider
- **WHEN** config is loaded via mapstructure/viper
- **THEN** `ProviderConfig.Type` SHALL deserialize correctly as `types.ProviderType`

#### Scenario: Provider validation
- **WHEN** a `ProviderConfig` is created with an unknown provider type
- **THEN** `config.Type.Valid()` SHALL return `false`

### Requirement: AgentConfig fields
`AgentConfig` SHALL include `MaxTurns int`, `ErrorCorrectionEnabled *bool`, and `MaxDelegationRounds int` fields with mapstructure/json tags.

#### Scenario: Zero-value defaults
- **WHEN** config omits `maxTurns`, `errorCorrectionEnabled`, and `maxDelegationRounds`
- **THEN** the zero values (0, nil, 0) SHALL be interpreted as defaults (25, true, 10) by the wiring layer

### Requirement: ObservationalMemoryConfig fields
`ObservationalMemoryConfig` SHALL include `MemoryTokenBudget int` and `ReflectionConsolidationThreshold int` fields with mapstructure/json tags.

#### Scenario: Zero-value defaults
- **WHEN** config omits `memoryTokenBudget` and `reflectionConsolidationThreshold`
- **THEN** the zero values SHALL be interpreted as defaults (4000, 5) by the wiring layer
