// Package onboard implements the lango onboard command.
package onboard

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/langowarny/lango/internal/cli/common"
	"github.com/spf13/cobra"
)

// NewCommand creates the onboard command.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "onboard",
		Short: "Interactive setup wizard for Lango",
		Long: `The onboard command guides you through setting up Lango for the first time.

Steps:
  1. Welcome and mode selection (QuickStart vs Advanced)
  2. API key configuration
  3. Model selection
  4. Channel setup (Telegram, Discord, or Slack)

Configuration is saved to lango.json in the current directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOnboard()
		},
	}

	return cmd
}

func runOnboard() error {
	p := tea.NewProgram(NewWizard())
	model, err := p.Run()
	if err != nil {
		return fmt.Errorf("onboard wizard error: %w", err)
	}

	wizard, ok := model.(*Wizard)
	if !ok {
		return fmt.Errorf("unexpected model type")
	}

	if wizard.Cancelled {
		fmt.Println("\nOnboard cancelled.")
		return nil
	}

	if wizard.Completed {
		if err := wizard.SaveConfig(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		wizard.PrintNextSteps()
	}

	return nil
}

// PrintNextSteps shows hints after successful onboarding.
func (w *Wizard) PrintNextSteps() {
	fmt.Println("\n✓ Lango configuration saved to lango.json")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Set your API key environment variable:")

	if meta, ok := common.GetProviderMetadata(w.state.Current.Agent.Provider); ok && meta.EnvVar != "" {
		fmt.Printf("     export %s=your-key-here\n", meta.EnvVar)
	}

	// Check channels
	if w.state.Current.Channels.Telegram.Enabled {
		fmt.Println("  2. Set your Telegram token:")
		fmt.Println("     export TELEGRAM_BOT_TOKEN=your-token")
	}
	if w.state.Current.Channels.Discord.Enabled {
		fmt.Println("  2. Set your Discord token:")
		fmt.Println("     export DISCORD_BOT_TOKEN=your-token")
	}
	if w.state.Current.Channels.Slack.Enabled {
		fmt.Println("  2. Set your Slack tokens:")
		fmt.Println("     export SLACK_BOT_TOKEN=your-bot-token")
		fmt.Println("     export SLACK_APP_TOKEN=your-app-token")
	}

	fmt.Println("\n  3. Start Lango:")
	fmt.Println("     lango serve")
	fmt.Println("\n  4. (Optional) Run doctor to verify setup:")
	fmt.Println("     lango doctor")

	// Write environment hints to a file for easy sourcing
	_ = os.WriteFile(".lango.env.example", []byte(w.generateEnvExample()), 0644)
	fmt.Println("\n  See .lango.env.example for a template of required environment variables.")
}

func (w *Wizard) generateEnvExample() string {
	example := "# Lango Environment Variables\n"
	example += "# Copy this file to .lango.env and fill in your values\n\n"

	if meta, ok := common.GetProviderMetadata(w.state.Current.Agent.Provider); ok && meta.EnvVar != "" {
		example += fmt.Sprintf("%s=your-%s-api-key\n", meta.EnvVar, w.state.Current.Agent.Provider)
	}

	if w.state.Current.Channels.Telegram.Enabled {
		example += "TELEGRAM_BOT_TOKEN=your-telegram-bot-token\n"
	}
	if w.state.Current.Channels.Discord.Enabled {
		example += "DISCORD_BOT_TOKEN=your-discord-bot-token\n"
	}
	if w.state.Current.Channels.Slack.Enabled {
		example += "SLACK_BOT_TOKEN=xoxb-your-bot-token\n"
		example += "SLACK_APP_TOKEN=xapp-your-app-token\n"
	}

	return example
}
