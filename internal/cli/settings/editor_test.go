package settings

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditor_EscAtWelcome_Quits(t *testing.T) {
	e := NewEditor()
	require.Equal(t, StepWelcome, e.step)

	model, cmd := e.Update(tea.KeyMsg{Type: tea.KeyEsc})
	ed := model.(*Editor)

	assert.Equal(t, StepWelcome, ed.step)
	assert.NotNil(t, cmd, "esc at welcome should return quit cmd")
}

func TestEditor_EscAtMenu_NavigatesToWelcome(t *testing.T) {
	e := NewEditor()
	e.step = StepMenu

	model, cmd := e.Update(tea.KeyMsg{Type: tea.KeyEsc})
	ed := model.(*Editor)

	assert.Equal(t, StepWelcome, ed.step)
	assert.Nil(t, cmd, "esc at menu should not quit, just navigate back")
}

func TestEditor_EscAtMenuWhileSearching_StaysAtMenu(t *testing.T) {
	e := NewEditor()
	e.step = StepMenu

	// Enter search mode by pressing /
	model, _ := e.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	ed := model.(*Editor)
	require.True(t, ed.menu.IsSearching(), "should be in search mode")

	// Press esc — should cancel search, not navigate back
	model, cmd := ed.Update(tea.KeyMsg{Type: tea.KeyEsc})
	ed = model.(*Editor)

	assert.Equal(t, StepMenu, ed.step, "should stay at menu")
	assert.False(t, ed.menu.IsSearching(), "search should be cancelled")
	assert.Nil(t, cmd)
}

func TestEditor_CtrlC_AlwaysQuits(t *testing.T) {
	tests := []struct {
		give string
		step EditorStep
	}{
		{give: "welcome", step: StepWelcome},
		{give: "menu", step: StepMenu},
		{give: "form", step: StepForm},
		{give: "providers_list", step: StepProvidersList},
		{give: "auth_providers_list", step: StepAuthProvidersList},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			e := NewEditor()
			e.step = tt.step

			model, cmd := e.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
			ed := model.(*Editor)

			assert.True(t, ed.Cancelled, "ctrl+c should set Cancelled")
			assert.NotNil(t, cmd, "ctrl+c should return quit cmd")
		})
	}
}
