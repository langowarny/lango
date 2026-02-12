## 1. Setup and Dependencies

- [x] 1.1 Add `golang.org/x/oauth2` dependency to `go.mod`.
- [x] 1.2 Update `lango.json` configuration schema (`internal/config/types.go`) to include `clientId`, `clientSecret`, and `scopes` in `ProviderConfig`.

## 2. CLI Authentication (`auth` package)

- [x] 2.1 Create `internal/cli/auth` package and scaffold `auth.go`.
- [x] 2.2 Implement `lango login [provider]` command structure.
- [x] 2.3 Implement local HTTP server for OAuth callback handling.
- [x] 2.4 Implement token exchange logic using `oauth2` library.
- [x] 2.5 Implement secure token storage (save JSON to `~/.lango/tokens/`).

## 3. Supervisor Integration

- [x] 3.1 Update `Supervisor.initializeProviders` to check for OAuth tokens.
- [x] 3.2 Implement `getAccessToken` method to load and validate tokens.
- [x] 3.3 Implement automatic token refresh logic using refresh tokens.
- [x] 3.4 Update provider initialization to use the resolved access token instead of `apiKey`.

## 4. Verification and Cleanup

- [ ] 4.1 Verify `lango login google` flow locally.
- [ ] 4.2 Verify `lango login github` flow locally.
- [ ] 4.3 Verify agent can make calls using the OAuth token.
- [ ] 4.4 Ensure sensitive files (tokens) have correct permissions (`0600`).
