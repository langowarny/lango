## 1. Critical Bug Fix

- [x] 1.1 Fix cursor bounds panic in `form.go` — change `len(m.Fields)` to no-op when cursor is at 0

## 2. Config Default Fix

- [x] 2.1 Change default `DatabasePath` from `sessions.db` to `data.db` in `config/loader.go`

## 3. Dead Code Removal

- [x] 3.1 Remove DB Passphrase field from `NewSecurityForm()` in `forms_impl.go`

## 4. Session Form Separation

- [x] 4.1 Create `NewSessionForm()` function in `forms_impl.go` with DB Path, TTL, Max History Turns fields
- [x] 4.2 Remove session fields (db_path, ttl, max_history_turns) from `NewSecurityForm()`
- [x] 4.3 Add "Session" category to menu in `menu.go` before "Security"
- [x] 4.4 Add "session" case to `handleMenuSelection()` in `wizard.go`

## 5. Provider Delete

- [x] 5.1 Add `Deleted` field to `ProvidersListModel` struct in `providers_list.go`
- [x] 5.2 Add `d` key binding in `ProvidersListModel.Update()` to set `Deleted` for current provider
- [x] 5.3 Handle `Deleted` in `wizard.go` — remove from state and refresh list
- [x] 5.4 Update help footer to include `d: delete`

## 6. Dynamic Provider Options

- [x] 6.1 Create `buildProviderOptions()` helper in `forms_impl.go`
- [x] 6.2 Use dynamic options for Provider field in `NewAgentForm()`
- [x] 6.3 Use dynamic options for Fallback Provider field in `NewAgentForm()`

## 7. Provider ID UX

- [x] 7.1 Reorder `NewProviderForm()` to show Type selector before ID field
- [x] 7.2 Change ID field label to "Provider Name" and update placeholder

## 8. Test Updates

- [x] 8.1 Update `TestNewSecurityForm_AllFields` to expect 9 fields (without session/passphrase fields)
- [x] 8.2 Add `TestNewSessionForm_AllFields` test for the new session form
- [x] 8.3 Verify `go build ./...` passes
- [x] 8.4 Verify `go test ./...` passes
