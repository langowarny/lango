## Context

Currently, Lango agents require an API key to be manually configured in `lango.json`. This is inconvenient and less secure than OAuth-based authentication. We want to support OAuth flows (like `gcloud auth login`) to obtain access tokens for Google (Gemini) and GitHub (Models).

## Goals / Non-Goals

**Goals:**
- Enable `lango login [provider]` for Google and GitHub.
- Support OAuth 2.0 Authorization Code flow with PKCE (if applicable) or standard web flow.
- Securely store and refresh tokens.
- Transparently use tokens in `Supervisor` when initializing providers.

**Non-Goals:**
- Implementing a full multi-user authentication system for the Lango server itself (this is just for the agent's API access).
- Supporting every possible OAuth provider immediately (start with Google/GitHub).

## Decisions

- **Loopback Redirect**: Use a local HTTP server on a random port to receive the OAuth callback. This is the standard pattern for CLI tools.
- **Token Storage**: Store tokens as JSON files in `~/.lango/tokens/`. While system keychains are more secure, file-based storage is portable and sufficient for an initial implementation, matching tools like `gcloud` or `gh` CLI (which often use file storage or helpers).
- **Config Schema**: Add `clientId`, `clientSecret`, and `scopes` to the `provider` configuration in `lango.json` to allow users to bring their own OAuth apps.
- **Provider Initialization**: Modify `Supervisor` to check for a valid token file first. If found, use it. If not, fall back to `apiKey`. If token exists but is expired, attempt refresh using the stored refresh token.

## Risks / Trade-offs

- **Risk**: Storing tokens in plain JSON files.
    - **Mitigation**: Set file permissions to `0600` (user read/write only). Future work can integrate with OS keychains (e.g., `keyring` library).
- **Risk**: Port conflicts for the callback server.
    - **Mitigation**: Use port `0` to let the OS assign a free random port.

## Migration Plan

- Users can continue using `apiKey` without changes.
- To use OAuth, users must update `lango.json` with OAuth client credentials and run `lango login`.
