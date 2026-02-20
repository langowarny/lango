package agent

import (
	"regexp"
	"sort"
)

// PIIDetector detects PII occurrences in text.
type PIIDetector interface {
	Detect(text string) []PIIMatch
}

var _ PIIDetector = (*RegexDetector)(nil)
var _ PIIDetector = (*CompositeDetector)(nil)

// compiledPattern holds a compiled regex and its metadata.
type compiledPattern struct {
	name     string
	category PIICategory
	re       *regexp.Regexp
	validate func(string) bool
}

// RegexDetector detects PII using compiled regex patterns.
type RegexDetector struct {
	patterns []compiledPattern
}

// RegexDetectorConfig configures which patterns the RegexDetector uses.
type RegexDetectorConfig struct {
	DisabledBuiltins []string
	CustomPatterns   map[string]string // name -> regex
	CustomRegex      []string          // legacy unnamed custom patterns

	// Legacy toggles for backward compatibility.
	RedactEmail bool
	RedactPhone bool
}

// NewRegexDetector creates a RegexDetector with the configured patterns.
func NewRegexDetector(cfg RegexDetectorConfig) *RegexDetector {
	disabled := make(map[string]bool, len(cfg.DisabledBuiltins))
	for _, name := range cfg.DisabledBuiltins {
		disabled[name] = true
	}

	var patterns []compiledPattern

	// Add builtin patterns (enabled by default, not disabled).
	for _, bp := range BuiltinPatterns {
		if disabled[bp.Name] {
			continue
		}
		if !bp.EnabledDefault {
			continue
		}

		// Legacy backward compatibility: skip email/phone if explicitly disabled.
		if bp.Name == "email" && !cfg.RedactEmail {
			continue
		}
		if bp.Name == "us_phone" && !cfg.RedactPhone {
			continue
		}

		re := compilePattern(bp.Pattern)
		if re == nil {
			continue
		}
		patterns = append(patterns, compiledPattern{
			name:     bp.Name,
			category: bp.Category,
			re:       re,
			validate: bp.Validate,
		})
	}

	// Add named custom patterns.
	for name, pattern := range cfg.CustomPatterns {
		re := compilePattern(pattern)
		if re == nil {
			continue
		}
		patterns = append(patterns, compiledPattern{
			name:     name,
			category: "custom",
			re:       re,
		})
	}

	// Add legacy unnamed custom regex patterns.
	for i, pattern := range cfg.CustomRegex {
		re := compilePattern(pattern)
		if re == nil {
			continue
		}
		patterns = append(patterns, compiledPattern{
			name:     "custom_" + string(rune('0'+i)),
			category: "custom",
			re:       re,
		})
	}

	return &RegexDetector{patterns: patterns}
}

// Detect finds all PII matches in the given text.
func (d *RegexDetector) Detect(text string) []PIIMatch {
	var matches []PIIMatch

	for _, p := range d.patterns {
		locs := p.re.FindAllStringIndex(text, -1)
		for _, loc := range locs {
			matched := text[loc[0]:loc[1]]

			// Apply optional validation (e.g., Luhn for credit cards).
			if p.validate != nil && !p.validate(matched) {
				continue
			}

			matches = append(matches, PIIMatch{
				PatternName: p.name,
				Category:    p.category,
				Start:       loc[0],
				End:         loc[1],
				Score:       1.0,
			})
		}
	}

	return matches
}

// CompositeDetector chains multiple PIIDetectors and deduplicates overlapping matches.
type CompositeDetector struct {
	detectors []PIIDetector
}

// NewCompositeDetector creates a CompositeDetector from multiple detectors.
func NewCompositeDetector(detectors ...PIIDetector) *CompositeDetector {
	return &CompositeDetector{detectors: detectors}
}

// Detect runs all child detectors and merges results, preferring higher-score
// matches when ranges overlap.
func (c *CompositeDetector) Detect(text string) []PIIMatch {
	var all []PIIMatch
	for _, d := range c.detectors {
		all = append(all, d.Detect(text)...)
	}

	if len(all) == 0 {
		return nil
	}

	// Sort by start position, then by score descending.
	sort.Slice(all, func(i, j int) bool {
		if all[i].Start != all[j].Start {
			return all[i].Start < all[j].Start
		}
		return all[i].Score > all[j].Score
	})

	// Deduplicate overlapping matches, keeping higher-score matches.
	var result []PIIMatch
	lastEnd := -1
	for _, m := range all {
		if m.Start >= lastEnd {
			result = append(result, m)
			if m.End > lastEnd {
				lastEnd = m.End
			}
		}
	}

	return result
}
