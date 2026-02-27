## ADDED Requirements

### Requirement: Settings help lists all category groups
The `lango settings --help` output SHALL display all 6 group sections (Core, Communication, AI & Knowledge, Infrastructure, P2P Network, Security) with their constituent categories.

#### Scenario: User views settings help
- **WHEN** user runs `lango settings --help`
- **THEN** the output lists Core (Providers, Agent, Server, Session), Communication (Channels, Tools, Multi-Agent, A2A Protocol), AI & Knowledge (Knowledge, Skill, Observational Memory, Embedding & RAG, Graph Store, Librarian), Infrastructure (Payment, Cron Scheduler, Background Tasks, Workflow Engine), P2P Network (P2P Network, P2P ZKP, P2P Pricing, P2P Owner Protection, P2P Sandbox), and Security (Security, Auth, Security Keyring, Security DB Encryption, Security KMS)

### Requirement: Settings help mentions keyword search
The `lango settings --help` output SHALL mention the `/` key for keyword search across categories.

#### Scenario: Search feature documented
- **WHEN** user runs `lango settings --help`
- **THEN** the output includes instruction to press `/` to search across all categories by keyword

### Requirement: Doctor help lists all 14 checks
The `lango doctor --help` output SHALL list all 14 diagnostic checks performed.

#### Scenario: User views doctor help
- **WHEN** user runs `lango doctor --help`
- **THEN** the output lists all 14 checks: configuration profile validity, AI provider configuration, API key security, channel token validation, session database, server port, security configuration, companion connectivity, observational memory, output scanning, embedding/RAG, graph store, multi-agent, and A2A protocol

### Requirement: Doctor help documents fix and json flags
The `lango doctor --help` output SHALL describe the `--fix` and `--json` flags in the Long description.

#### Scenario: Flags documented in description
- **WHEN** user runs `lango doctor --help`
- **THEN** the Long description includes usage guidance for `--fix` (automatic repair) and `--json` (machine-readable output)

### Requirement: Onboard help reflects current provider list
The `lango onboard --help` output SHALL list all supported providers including GitHub in step 1.

#### Scenario: GitHub provider listed
- **WHEN** user runs `lango onboard --help`
- **THEN** step 1 lists Anthropic, OpenAI, Gemini, Ollama, and GitHub as provider choices

### Requirement: Onboard help reflects model auto-fetch
The `lango onboard --help` output SHALL mention that models are auto-fetched from the provider in step 2.

#### Scenario: Auto-fetch mentioned
- **WHEN** user runs `lango onboard --help`
- **THEN** step 2 description includes that model selection uses auto-fetched models from the provider

### Requirement: Onboard help reflects approval policy
The `lango onboard --help` output SHALL mention approval policy in step 4.

#### Scenario: Approval policy mentioned
- **WHEN** user runs `lango onboard --help`
- **THEN** step 4 description includes approval policy alongside privacy interceptor and PII redaction
