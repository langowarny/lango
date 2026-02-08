package onboard

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/langowarny/lango/internal/config"
)

type ProviderItem struct {
	ID   string
	Type string
}

type ProvidersListModel struct {
	Providers []ProviderItem
	Cursor    int
	Selected  string // ID of selected provider, or "NEW"
	Exit      bool   // True if user wants to go back
}

func NewProvidersListModel(cfg *config.Config) ProvidersListModel {
	var items []ProviderItem
	if cfg.Providers != nil {
		for id, p := range cfg.Providers {
			items = append(items, ProviderItem{ID: id, Type: p.Type})
		}
	}
	// Sort for consistent display
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	return ProvidersListModel{
		Providers: items,
		Cursor:    0,
	}
}

func (m ProvidersListModel) Init() tea.Cmd {
	return nil
}

func (m ProvidersListModel) Update(msg tea.Msg) (ProvidersListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			// +1 for "Add New Provider" option
			if m.Cursor < len(m.Providers) {
				m.Cursor++
			}
		case "enter":
			if m.Cursor == len(m.Providers) {
				m.Selected = "NEW"
			} else {
				m.Selected = m.Providers[m.Cursor].ID
			}
			return m, nil
		case "esc":
			m.Exit = true
			return m, nil
		}
	}
	return m, nil
}

func (m ProvidersListModel) View() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		MarginBottom(1)

	b.WriteString(titleStyle.Render("Manage Providers"))
	b.WriteString("\n\n")

	// Render existing providers
	for i, p := range m.Providers {
		cursor := "  "
		itemStyle := lipgloss.NewStyle()

		if m.Cursor == i {
			cursor = "▸ "
			itemStyle = itemStyle.Foreground(lipgloss.Color("#04B575")).Bold(true)
		}

		b.WriteString(cursor)
		label := fmt.Sprintf("%s (%s)", p.ID, p.Type)
		b.WriteString(itemStyle.Render(label))
		b.WriteString("\n")
	}

	// Render "Add New" option
	cursor := "  "
	itemStyle := lipgloss.NewStyle()
	if m.Cursor == len(m.Providers) {
		cursor = "▸ "
		itemStyle = itemStyle.Foreground(lipgloss.Color("#04B575")).Bold(true)
	}
	b.WriteString(cursor)
	b.WriteString(itemStyle.Render("+ Add New Provider"))
	b.WriteString("\n\n")

	b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render("Use arrow keys to navigate • enter to select • esc to back"))

	return b.String()
}
