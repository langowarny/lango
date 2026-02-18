## REMOVED Requirements

### Requirement: OAuth Login Command
**Reason**: OAuth with AI providers risks account bans. API key authentication is the supported method.
**Migration**: Remove `clientId`, `clientSecret`, and `scopes` from provider configuration. Set `apiKey` using environment variable references (e.g., `${GOOGLE_API_KEY}`).

### Requirement: OAuth Callback Handling
**Reason**: Removed along with OAuth Login Command.
**Migration**: No action needed — callback handling was internal to the login flow.

### Requirement: Secure Token Storage
**Reason**: OAuth tokens are no longer generated or stored.
**Migration**: Existing token files in `~/.lango/tokens/` can be safely deleted.

### Requirement: Automatic Token Refresh
**Reason**: No OAuth tokens to refresh.
**Migration**: Use API keys which do not expire (or rotate manually per provider policy).
