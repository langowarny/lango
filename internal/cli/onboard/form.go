package onboard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InputType defines the type of input field
type InputType int

const (
	InputText InputType = iota
	InputInt
	InputBool // Toggled via space
	InputSelect
)

// Field represents a single configuration field in a form
type Field struct {
	Key         string
	Label       string
	Type        InputType
	Value       string
	Placeholder string
	Options     []string // For InputSelect
	Checked     bool     // For InputBool
	Width       int
	Validate    func(string) error

	// Internal UI state
	textInput textinput.Model
	err       error
}

// FormModel manages a list of fields
type FormModel struct {
	Title    string
	Fields   []*Field
	Cursor   int
	Focus    bool
	OnSave   func(map[string]interface{})
	OnCancel func()
}

func NewFormModel(title string) FormModel {
	return FormModel{
		Title:  title,
		Fields: []*Field{},
		Cursor: 0,
	}
}

func (m *FormModel) AddField(f *Field) {
	if f.Type == InputText || f.Type == InputInt {
		ti := textinput.New()
		ti.Placeholder = f.Placeholder
		ti.SetValue(f.Value)
		ti.CharLimit = 100
		ti.Width = 30
		if f.Width > 0 {
			ti.Width = f.Width
		}
		f.textInput = ti
	}
	m.Fields = append(m.Fields, f)
}

func (m FormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m FormModel) Update(msg tea.Msg) (FormModel, tea.Cmd) {
	if !m.Focus {
		return m, nil
	}

	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "shift+tab":
			if m.Cursor > 0 {
				m.Cursor--
			} else {
				m.Cursor = len(m.Fields) // Wrap to buttons if we add them later?
				// For now stick to fields
			}
		case "down", "tab":
			if m.Cursor < len(m.Fields)-1 {
				m.Cursor++
			}
		case "enter":
			// For boolean, toggle
			// For others, maybe validate?
			// Usually enter on a form ideally submits, but here we navigate fields
			// keeping standard navigation simple
		case " ":
			field := m.Fields[m.Cursor]
			if field.Type == InputBool {
				field.Checked = !field.Checked
			}
		case "esc":
			if m.OnCancel != nil {
				m.OnCancel()
			}
		}
	}

	// Update specific field logic
	field := m.Fields[m.Cursor]
	if field.Type == InputText || field.Type == InputInt {
		var inputCmd tea.Cmd
		field.textInput, inputCmd = field.textInput.Update(msg)
		field.Value = field.textInput.Value()
		cmd = inputCmd
	}

	// Handle Select Logic (Left/Right to cycle options)
	if field.Type == InputSelect {
		if msg, ok := msg.(tea.KeyMsg); ok {
			switch msg.String() {
			case "right", "l":
				// Cycle forward
				// Need current index logic
				// Simplified: Just cycle blindly for now or find index
				idx := -1
				for i, opt := range field.Options {
					if opt == field.Value {
						idx = i
						break
					}
				}
				if idx < len(field.Options)-1 {
					field.Value = field.Options[idx+1]
				} else if len(field.Options) > 0 {
					field.Value = field.Options[0]
				}
			case "left", "h":
				// Cycle backward
				idx := -1
				for i, opt := range field.Options {
					if opt == field.Value {
						idx = i
						break
					}
				}
				if idx > 0 {
					field.Value = field.Options[idx-1]
				} else if len(field.Options) > 0 {
					field.Value = field.Options[len(field.Options)-1]
				}
			}
		}
	}

	return m, cmd
}

func (m FormModel) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(lipgloss.Color("#7D56F4")).MarginBottom(1)
	b.WriteString(titleStyle.Render(m.Title))
	b.WriteString("\n")

	for i, f := range m.Fields {
		labelStyle := lipgloss.NewStyle().Width(20)
		if i == m.Cursor {
			labelStyle = labelStyle.Foreground(lipgloss.Color("#04B575")).Bold(true)
		}

		b.WriteString(labelStyle.Render(f.Label))

		switch f.Type {
		case InputText, InputInt:
			if i == m.Cursor {
				f.textInput.Focus()
				f.textInput.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575"))
			} else {
				f.textInput.Blur()
				f.textInput.TextStyle = lipgloss.NewStyle()
			}
			b.WriteString(f.textInput.View())

		case InputBool:
			check := "[ ]"
			if f.Checked {
				check = "[x]"
			}
			if i == m.Cursor {
				check = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render(check)
			}
			b.WriteString(check)

		case InputSelect:
			val := f.Value
			if val == "" && len(f.Options) > 0 {
				val = f.Options[0]
			}
			if i == m.Cursor {
				val = fmt.Sprintf("< %s >", val)
				val = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Render(val)
			}
			b.WriteString(val)
		}
		b.WriteString("\n")
	}

	// Help Footer
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render("tab/shift+tab: nav • space: toggle • ←/→: select options • esc: back"))

	return b.String()
}
