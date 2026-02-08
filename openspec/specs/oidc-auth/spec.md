## ADDED Requirements

### Requirement: OIDC Provider Configuration
The system SHALL allow configuring multiple OIDC providers (Google, GitHub) via `lango.json`, specifying Client ID, Client Secret, and Issuer URL.

#### Scenario: Config Loading
- **WHEN** the application starts with OIDC config
- **THEN** the system initializes the OIDC provider verifiers

### Requirement: Login Flow
The system SHALL provide HTTP endpoints to initiate OIDC login and handle the callback.

#### Scenario: Login Initiation
- **WHEN** a user accesses `/auth/login/google`
- **THEN** they are redirected to Google's consent page with the correct state and scopes

#### Scenario: Callback Handling
- **WHEN** the user returns to `/auth/callback/google` with a valid code
- **THEN** the system exchanges the code for an ID token
- **AND** issues a Lango session cookie or token to the user
