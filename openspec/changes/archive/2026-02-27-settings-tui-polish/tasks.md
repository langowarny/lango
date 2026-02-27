## 1. Design System Tokens (internal/cli/tui/styles.go)

- [x] 1.1 Define color palette constants: Primary, Success, Warning, Error, Muted, Foreground, Background, Highlight, Accent, Dim, Separator
- [x] 1.2 Define reusable styles: TitleStyle, SubtitleStyle, SuccessStyle, WarningStyle, ErrorStyle, MutedStyle, HighlightStyle, BoxStyle, ListItemStyle, SelectedItemStyle
- [x] 1.3 Define menu-specific styles: SectionHeaderStyle, SeparatorLineStyle, CursorStyle, ActiveItemStyle, SearchBarStyle, FormTitleBarStyle, FieldDescStyle
- [x] 1.4 Implement Breadcrumb(segments ...string) function with muted prefix segments and primary bold last segment
- [x] 1.5 Implement KeyBadge(key string), HelpEntry(key, label string), and HelpBar(entries ...string) functions
- [x] 1.6 Implement FormatPass, FormatWarn, FormatFail, FormatMuted helper functions

## 2. Breadcrumb Navigation (internal/cli/settings/editor.go)

- [x] 2.1 Add dynamic breadcrumb header in Editor.View() for StepWelcome and StepMenu ("Settings")
- [x] 2.2 Add breadcrumb for StepForm using activeForm.Title ("Settings > {form title}")
- [x] 2.3 Add breadcrumb for StepProvidersList ("Settings > Providers")
- [x] 2.4 Add breadcrumb for StepAuthProvidersList ("Settings > Auth Providers")

## 3. Styled Containers

- [x] 3.1 Wrap welcome screen in RoundedBorder container with Primary border color (editor.go viewWelcome)
- [x] 3.2 Wrap menu body in RoundedBorder container with Muted border color (menu.go View)
- [x] 3.3 Wrap providers list body in RoundedBorder container with Muted border color (providers_list.go View)
- [x] 3.4 Wrap auth providers list body in RoundedBorder container with Muted border color (auth_providers_list.go View)

## 4. Help Bars

- [x] 4.1 Add HelpBar to welcome screen: Enter (Start), Esc (Quit)
- [x] 4.2 Add HelpBar to menu normal mode: Navigate, Select, Search, Back
- [x] 4.3 Add HelpBar to menu search mode: Navigate, Select, Cancel
- [x] 4.4 Add HelpBar to providers list: Navigate, Select, Delete, Back
- [x] 4.5 Add HelpBar to auth providers list: Navigate, Select, Delete, Back

## 5. Verification

- [x] 5.1 Run go build ./... to verify no build errors
- [x] 5.2 Run go test ./... to verify no test failures
