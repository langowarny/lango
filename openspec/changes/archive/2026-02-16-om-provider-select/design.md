## Context

The Observational Memory (OM) onboard form currently uses `InputText` for its provider field, while the Agent form uses `InputSelect` with dynamically populated provider options via `buildProviderOptions(cfg)`. The Agent form's fallback_provider field follows the same pattern: empty string first option + registered providers. This inconsistency means OM users must type provider names manually.

## Goals / Non-Goals

**Goals:**
- Make OM provider field consistent with Agent form's provider selection pattern
- Reuse existing `buildProviderOptions(cfg)` helper for option generation
- Include empty string option as first choice (= "use agent default")
- Maintain test coverage for the changed field type

**Non-Goals:**
- Cross-field dependency (dynamically filtering models based on selected provider)
- Changing `om_model` field type (remains InputText, consistent with Agent form's model field)
- Modifying `buildProviderOptions` behavior

## Decisions

1. **Reuse `buildProviderOptions(cfg)` instead of hardcoding options**: The function already handles the fallback logic (defaults to anthropic/openai/gemini/ollama when no providers registered). Reusing it ensures OM stays in sync with Agent form automatically.

2. **Prepend empty string option**: Following the `fallback_provider` pattern (`append([]string{""}, providerOpts...)`), the empty string represents "use agent default", which is the current behavior when provider is left blank.

3. **Keep `om_model` as InputText**: Models vary per provider and the form system lacks cross-field dependency support. Agent form also uses InputText for model, so this maintains consistency.

## Risks / Trade-offs

- [Minimal risk] If a user had previously saved a provider value not in the registered list, the select dropdown won't show it as selected → The value is preserved in config but may appear blank in the dropdown. This matches existing Agent form behavior and is acceptable.
