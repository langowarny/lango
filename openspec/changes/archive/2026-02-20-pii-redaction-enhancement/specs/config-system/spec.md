## ADDED Requirements

### Requirement: InterceptorConfig PII pattern fields
InterceptorConfig SHALL include PIIDisabledPatterns ([]string), PIICustomPatterns (map[string]string), and Presidio (PresidioConfig) fields with appropriate mapstructure and json tags.

#### Scenario: Disabled patterns config
- **WHEN** config JSON contains "piiDisabledPatterns": ["passport", "ipv4"]
- **THEN** InterceptorConfig.PIIDisabledPatterns SHALL be ["passport", "ipv4"]

#### Scenario: Custom patterns config
- **WHEN** config JSON contains "piiCustomPatterns": {"my_id": "\\bID-\\d+\\b"}
- **THEN** InterceptorConfig.PIICustomPatterns SHALL contain the mapping

### Requirement: PresidioConfig type
A new PresidioConfig struct SHALL define Enabled (bool), URL (string, default "http://localhost:5002"), ScoreThreshold (float64, default 0.7), and Language (string, default "en").

#### Scenario: Presidio config loading
- **WHEN** config JSON contains presidio block with enabled=true, url, scoreThreshold, language
- **THEN** InterceptorConfig.Presidio SHALL be populated

#### Scenario: Default values
- **WHEN** no Presidio config is specified
- **THEN** URL SHALL default to "http://localhost:5002", ScoreThreshold to 0.7, Language to "en"
