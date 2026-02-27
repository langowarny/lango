// Package tuicore provides shared TUI form components for CLI commands.
package tuicore

import "github.com/charmbracelet/bubbles/textinput"

// InputType defines the type of input field.
type InputType int

const (
	InputText InputType = iota
	InputInt
	InputBool // Toggled via space
	InputSelect
	InputPassword
)

// Field represents a single configuration field in a form.
type Field struct {
	Key         string
	Label       string
	Description string // Help text shown below the focused field
	Type        InputType
	Value       string
	Placeholder string
	Options     []string // For InputSelect
	Checked     bool     // For InputBool
	Width       int
	Validate    func(string) error

	// VisibleWhen controls conditional visibility. When non-nil, the field is
	// shown only when this function returns true. When nil the field is always visible.
	VisibleWhen func() bool

	// TextInput holds the bubbletea text input model (exported for cross-package use).
	TextInput textinput.Model
	Err       error
}

// IsVisible reports whether this field should be rendered and navigable.
func (f *Field) IsVisible() bool {
	if f.VisibleWhen == nil {
		return true
	}
	return f.VisibleWhen()
}
