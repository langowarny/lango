## ADDED Requirements

### Requirement: ResponseStatus enum type
The system SHALL define `ResponseStatus` as a typed string enum in `protocol/messages.go` with constants `ResponseStatusOK`, `ResponseStatusError`, `ResponseStatusDenied`, `ResponseStatusPaymentRequired` and a `Valid()` method.

#### Scenario: Response.Status uses typed enum
- **WHEN** the protocol handler constructs a `Response`
- **THEN** it SHALL set `Status` using `ResponseStatus` constants, never raw strings

#### Scenario: JSON wire format preserved
- **WHEN** a `Response` with `ResponseStatus` is serialized to JSON
- **THEN** the `status` field SHALL contain the plain string value (e.g., `"ok"`)

### Requirement: ACLAction enum type
The system SHALL define `ACLAction` as a typed string enum in `firewall/firewall.go` with constants `ACLActionAllow`, `ACLActionDeny` and a `Valid()` method.

#### Scenario: ACLRule.Action uses typed enum
- **WHEN** an `ACLRule` is constructed
- **THEN** the `Action` field SHALL be `ACLAction` type, not raw string

### Requirement: WildcardAll constant
The system SHALL define `WildcardAll = "*"` in `firewall/firewall.go`.

#### Scenario: Wildcard comparisons use constant
- **WHEN** firewall code checks for wildcard peer or tool patterns
- **THEN** it SHALL compare against `WildcardAll`, not the literal `"*"`

### Requirement: ProofScheme enum type
The system SHALL define `ProofScheme` as a typed string enum in `zkp/zkp.go` with constants `SchemePlonk`, `SchemeGroth16` and a `Valid()` method.

#### Scenario: ZKP config and proof use typed scheme
- **WHEN** `Config.Scheme`, `ProverService.scheme`, or `Proof.Scheme` stores a proving scheme
- **THEN** it SHALL use the `ProofScheme` type

### Requirement: SRSMode enum type
The system SHALL define `SRSMode` as a typed string enum in `zkp/zkp.go` with constants `SRSModeUnsafe`, `SRSModeFile` and a `Valid()` method.

#### Scenario: ZKP config uses typed SRS mode
- **WHEN** `Config.SRSMode` or `ProverService.srsMode` stores the SRS mode
- **THEN** it SHALL use the `SRSMode` type

### Requirement: KMSProviderName enum type
The system SHALL define `KMSProviderName` as a typed string enum in `security/kms_factory.go` with constants `KMSProviderAWS`, `KMSProviderGCP`, `KMSProviderAzure`, `KMSProviderPKCS11` and a `Valid()` method.

#### Scenario: NewKMSProvider accepts typed name
- **WHEN** `NewKMSProvider` is called
- **THEN** the `providerName` parameter SHALL be `KMSProviderName` type

### Requirement: ChainID type and constants
The system SHALL define `ChainID` as a typed `int64` in `wallet/wallet.go` with constants `ChainEthereumMainnet` (1), `ChainBase` (8453), `ChainBaseSepolia` (84532), `ChainSepolia` (11155111).

#### Scenario: NetworkName uses typed constants
- **WHEN** `NetworkName()` switches on a chain ID
- **THEN** it SHALL compare against `ChainID` constants

### Requirement: CurrencyUSDC constant
The system SHALL define `CurrencyUSDC = "USDC"` in `wallet/wallet.go`.

#### Scenario: All USDC references use constant
- **WHEN** any package references the USDC currency ticker
- **THEN** it SHALL use `wallet.CurrencyUSDC` instead of the string literal `"USDC"`
