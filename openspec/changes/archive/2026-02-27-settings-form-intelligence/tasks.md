## 1. Inline Field Descriptions

- [x] 1.1 Add `Description string` field to `tuicore.Field` struct in field.go
- [x] 1.2 Render description below focused field in FormModel.View() in form.go
- [x] 1.3 Add descriptions to all Agent form fields in forms_impl.go
- [x] 1.4 Add descriptions to Server, Channels, Tools, Session form fields
- [x] 1.5 Add descriptions to Security form fields (interceptor, PII, Presidio, signer)
- [x] 1.6 Add descriptions to Knowledge, Skill, Observational Memory form fields
- [x] 1.7 Add descriptions to Embedding & RAG, Graph Store form fields
- [x] 1.8 Add descriptions to Multi-Agent, A2A, Payment form fields
- [x] 1.9 Add descriptions to Cron, Background, Workflow form fields
- [x] 1.10 Add descriptions to Librarian, P2P Network, P2P ZKP form fields
- [x] 1.11 Add descriptions to P2P Pricing, Owner Protection, Sandbox form fields
- [x] 1.12 Add descriptions to Security Keyring, DB Encryption, KMS form fields
- [x] 1.13 Add descriptions to OIDC Provider form fields

## 2. Field Validators

- [x] 2.1 Add Temperature validator (0.0-2.0 range) to Agent form
- [x] 2.2 Add port validator (1-65535) to Server form
- [x] 2.3 Add Max Read Size validator (positive integer) to Tools form
- [x] 2.4 Add Max History Turns validator (positive integer) to Session form
- [x] 2.5 Add Knowledge Max Context validator (positive integer)
- [x] 2.6 Add validators to Observational Memory numeric fields (positive/non-negative)
- [x] 2.7 Add Embedding Dimensions and RAG Max Results validators (non-negative)
- [x] 2.8 Add Graph Max Depth and Max Expansion validators (positive integer)
- [x] 2.9 Add Cron Max Jobs validator (positive integer)
- [x] 2.10 Add Background Yield Time (non-negative) and Max Tasks (positive) validators
- [x] 2.11 Add Workflow Max Steps validator (positive integer)
- [x] 2.12 Add P2P Max Peers validator (positive integer)
- [x] 2.13 Add P2P Min Trust Score validator (0.0-1.0 range)
- [x] 2.14 Add Librarian numeric field validators (threshold, cooldown, max inquiries)
- [x] 2.15 Add Skill Max Bulk Import and Import Concurrency validators (positive integer)
- [x] 2.16 Add Security Approval Timeout validator (non-negative integer)
- [x] 2.17 Add Payment Chain ID validator (integer)

## 3. Auto-Fetch Model Options

- [x] 3.1 Create `model_fetcher.go` with `newProviderFromConfig()` supporting OpenAI, Anthropic, Gemini, Ollama, GitHub
- [x] 3.2 Implement `fetchModelOptions()` with 5s timeout, sorted output, current model inclusion
- [x] 3.3 Wire model fetch into NewAgentForm for primary model field
- [x] 3.4 Wire model fetch into NewAgentForm for fallback model field
- [x] 3.5 Wire model fetch into NewObservationalMemoryForm with agent provider fallback
- [x] 3.6 Wire model fetch into NewEmbeddingForm for embedding model
- [x] 3.7 Wire model fetch into NewLibrarianForm with agent provider fallback

## 4. Unify Embedding Provider

- [x] 4.1 Update NewEmbeddingForm to use single `emb_provider_id` field mapped to `cfg.Embedding.Provider`
- [x] 4.2 Update `emb_provider_id` case in UpdateConfigFromForm to also clear `cfg.Embedding.ProviderID`

## 5. Conditional Field Visibility

- [x] 5.1 Add `VisibleWhen func() bool` field to `tuicore.Field` struct
- [x] 5.2 Add `IsVisible() bool` method to Field
- [x] 5.3 Add `VisibleFields() []*Field` method to FormModel
- [x] 5.4 Update FormModel.Update() to use VisibleFields() for cursor navigation
- [x] 5.5 Update FormModel.View() to iterate VisibleFields() instead of Fields
- [x] 5.6 Add cursor clamping after visibility changes in Update()
- [x] 5.7 Add VisibleWhen closures to Channel token fields (Telegram, Discord, Slack)
- [x] 5.8 Add VisibleWhen closures to Security interceptor sub-fields
- [x] 5.9 Add nested VisibleWhen for Presidio fields (interceptor AND presidio enabled)
- [x] 5.10 Add VisibleWhen to Signer RPC URL and Key ID fields based on provider value
- [x] 5.11 Add VisibleWhen to P2P Container sandbox fields
- [x] 5.12 Add VisibleWhen to KMS Azure and PKCS11 fields based on backend type

## 6. Verification

- [x] 6.1 Run `go build ./...` -- zero errors
- [x] 6.2 Run `go test ./...` -- all tests pass
