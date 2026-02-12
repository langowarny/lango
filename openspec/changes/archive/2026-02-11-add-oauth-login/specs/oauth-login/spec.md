## ADDED Requirements

### Requirement: OAuth Login Command
The system SHALL provide a CLI command `lango login [provider]` to initiate OAuth authentication.

#### Scenario: User logs in with Google
- **WHEN** user runs `lango login google`
- **THEN** system starts a local web server on a random port
- **AND** opens the system browser to the Google authorization URL
- **AND** listens for the callback with the authorization code
- **AND** exchanges the code for an access token and refresh token

#### Scenario: User logs in with GitHub
- **WHEN** user runs `lango login github`
- **THEN** system performs OAuth flow with GitHub endpoints
- **AND** saves the resulting token

### Requirement: OAuth Callback Handling
The system SHALL handle the OAuth callback on `localhost`.

#### Scenario: Successful callback
- **WHEN** the provider redirects to `http://localhost:<port>/callback?code=<code>&state=<state>`
- **THEN** system validates the `state` parameter against the generated state
- **AND** exchanges `code` for tokens
- **AND** displays a success message to the user in the browser
- **AND** closes the local server

#### Scenario: Error callback
- **WHEN** the provider redirects with `error=access_denied`
- **THEN** system displays an error message
- **AND** CLI command exits with an error

### Requirement: Secure Token Storage
The system SHALL store OAuth tokens securely in the user's home directory.

#### Scenario: Token file creation
- **WHEN** authentication is successful
- **THEN** system saves the token JSON to `~/.lango/tokens/<provider>.json`
- **AND** sets file permissions to `0600` (read/write by owner only)

### Requirement: Automatic Token Refresh
The system SHALL automatically refresh expired access tokens if a refresh token is available.

#### Scenario: Access token expired
- **WHEN** `Supervisor` attempts to initialize a provider
- **AND** the stored access token is expired
- **AND** a refresh token is present
- **THEN** system requests a new access token from the provider
- **AND** updates the token file with the new credentials
- **AND** proceeds with provider initialization
