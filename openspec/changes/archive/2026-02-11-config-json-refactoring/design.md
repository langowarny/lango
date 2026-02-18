## Context

Currently, `lango.json` allows setting the API key in `agent.apiKey` (for simple setups) and `providers.<name>.apiKey` (for multiple providers). With the addition of OAuth, `agent.apiKey` becomes misleading and insufficient. We need to unify credential management.

## Goals / Non-Goals

**Goals:**
- Simplify `lango.json` by removing credential duplication.
- Enforce a strict separation: `agent` configures *behavior* (which provider/model to use), `providers` configures *access* (keys, tokens, endpoints).
- Update `lango.example.json` to show best practices.

**Non-Goals:**
- Changing how internal secrets storage (keychain) works; this is purely about the configuration file structure.

## Decisions

- **Remove `agent.apiKey`**: This field will be deprecated and removed from the `AgentConfig` struct. The application will fail to start if the user relies on it, forcing migration to the new structure.
- **Provider Reference**: `agent.provider` string value must match a key in the `providers` map. If it doesn't, and `providers` is empty, we can't initialize.
- **Implicit Default**: If `providers` has only one entry and `agent.provider` is empty, we *could* infer it, but explicit is better. We will require `agent.provider` to be set.

## Risks / Trade-offs

- **Risk**: Breaking change for existing users.
    - **Mitigation**: Fail fast with a clear error message: "agent.apiKey is no longer supported; please move your API key to the 'providers' section in lango.json".

## Migration Plan

1.  Update `lango.example.json` first.
2.  Update `internal/config` structs.
3.  Update `Supervisor` initialization logic.
4.  Instruct users to update their local `lango.json`.
