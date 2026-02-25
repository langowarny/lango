## ADDED Requirements

### Requirement: Validate overly permissive ACL rules
The firewall SHALL provide a `ValidateRule()` function that rejects allow rules with wildcard peer (`"*"`) combined with wildcard tools (empty list or containing `"*"`). Deny rules SHALL always pass validation.

#### Scenario: Wildcard peer with empty tools (allow)
- **WHEN** `ValidateRule` is called with `{PeerDID: "*", Action: "allow", Tools: []}`
- **THEN** it SHALL return an error "overly permissive rule: allow all peers with all tools is prohibited"

#### Scenario: Wildcard peer with wildcard tool (allow)
- **WHEN** `ValidateRule` is called with `{PeerDID: "*", Action: "allow", Tools: ["*"]}`
- **THEN** it SHALL return an error

#### Scenario: Wildcard peer with specific tools (allow)
- **WHEN** `ValidateRule` is called with `{PeerDID: "*", Action: "allow", Tools: ["echo"]}`
- **THEN** it SHALL return nil (allowed)

#### Scenario: Specific peer with wildcard tools (allow)
- **WHEN** `ValidateRule` is called with `{PeerDID: "did:key:abc", Action: "allow", Tools: ["*"]}`
- **THEN** it SHALL return nil (allowed)

#### Scenario: Wildcard deny rule
- **WHEN** `ValidateRule` is called with `{PeerDID: "*", Action: "deny", Tools: ["*"]}`
- **THEN** it SHALL return nil (deny rules always safe)

## MODIFIED Requirements

### Requirement: AddRule validates before adding
`AddRule()` SHALL validate the rule using `ValidateRule()` before adding it. If validation fails, it SHALL return the error without adding the rule. The return type changes from void to `error`.

#### Scenario: AddRule rejects overly permissive rule
- **WHEN** `AddRule` is called with a wildcard allow-all rule
- **THEN** it SHALL return an error and NOT add the rule to the firewall

#### Scenario: AddRule accepts valid rule
- **WHEN** `AddRule` is called with a specific peer allow rule
- **THEN** it SHALL add the rule and return nil

### Requirement: Initial rules backward compatibility
When constructing a Firewall with `New()`, overly permissive initial rules SHALL be loaded with a warning log (not rejected). This preserves backward compatibility with existing configurations while alerting operators.

#### Scenario: Overly permissive initial rule
- **WHEN** `New()` is called with a wildcard allow-all rule in the initial rules slice
- **THEN** the rule SHALL be loaded (backward compat) and a warning SHALL be logged
