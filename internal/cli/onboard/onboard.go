// Package onboard implements the lango onboard command.
package onboard

import (
	"context"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/langowarny/lango/internal/bootstrap"
	"github.com/langowarny/lango/internal/config"
	"github.com/langowarny/lango/internal/configstore"
)

// NewCommand creates the onboard command.
func NewCommand() *cobra.Command {
	var profileName string

	cmd := &cobra.Command{
		Use:   "onboard",
		Short: "Interactive setup wizard for Lango",
		Long: `The onboard command guides you through setting up Lango for the first time.

An interactive menu-based editor lets you configure each section independently:
  - Agent:      Provider, Model, Tokens, Fallback settings
  - Server:     Host, Port, HTTP/WebSocket toggles
  - Channels:   Telegram, Discord, Slack tokens
  - Tools:      Exec timeouts, Browser, Filesystem limits
  - Auth:       OIDC providers, JWT settings
  - Security:   PII interceptor, Signer
  - Session:    Session DB, TTL
  - Knowledge:  Learning limits, Skills, Context per layer
  - Providers:  Manage multiple provider configurations

All settings including API keys are saved in an encrypted profile (~/.lango/lango.db).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOnboard(profileName)
		},
	}

	cmd.Flags().StringVar(&profileName, "profile", "default", "Profile name to create or edit")

	return cmd
}

func runOnboard(profileName string) error {
	// 1. Bootstrap: DB + crypto + configstore initialization.
	//    This must happen before BubbleTea starts because passphrase
	//    acquisition requires terminal input.
	boot, err := bootstrap.Run(bootstrap.Options{})
	if err != nil {
		return fmt.Errorf("bootstrap: %w", err)
	}
	defer boot.DBClient.Close()

	ctx := context.Background()

	// 2. Load existing profile or start with defaults.
	initialCfg, isNew, err := loadOrDefault(ctx, boot.ConfigStore, profileName)
	if err != nil {
		return fmt.Errorf("load profile %q: %w", profileName, err)
	}

	// 3. Run TUI wizard with the initial config.
	p := tea.NewProgram(NewWizardWithConfig(initialCfg))
	model, err := p.Run()
	if err != nil {
		return fmt.Errorf("onboard wizard: %w", err)
	}

	wizard, ok := model.(*Wizard)
	if !ok {
		return fmt.Errorf("unexpected model type")
	}

	if wizard.Cancelled {
		fmt.Println("\nOnboard cancelled.")
		return nil
	}

	if !wizard.Completed {
		return nil
	}

	// 4. Save the edited config as an encrypted profile.
	cfg := wizard.Config()
	if err := boot.ConfigStore.Save(ctx, profileName, cfg); err != nil {
		return fmt.Errorf("save profile %q: %w", profileName, err)
	}

	// 5. Activate the profile if it is new.
	if isNew {
		if err := boot.ConfigStore.SetActive(ctx, profileName); err != nil {
			return fmt.Errorf("activate profile %q: %w", profileName, err)
		}
	}

	// 6. Print next steps.
	printNextSteps(cfg, profileName)

	return nil
}

// loadOrDefault loads an existing profile or returns a default config.
// Returns (config, isNew, error).
func loadOrDefault(ctx context.Context, store *configstore.Store, name string) (*config.Config, bool, error) {
	cfg, err := store.Load(ctx, name)
	if err == nil {
		return cfg, false, nil
	}
	if errors.Is(err, configstore.ErrProfileNotFound) {
		return config.DefaultConfig(), true, nil
	}
	return nil, false, err
}

// printNextSteps shows hints after successful onboarding.
func printNextSteps(_ *config.Config, profileName string) {
	fmt.Printf("\n✓ Configuration saved to encrypted profile %q\n", profileName)
	fmt.Println("  Storage: ~/.lango/lango.db")

	fmt.Println("\nNext steps:")
	fmt.Println("  1. Start Lango:")
	fmt.Println("     lango serve")
	fmt.Println("\n  2. (Optional) Run doctor to verify setup:")
	fmt.Println("     lango doctor")
	fmt.Println("\n  Profile management:")
	fmt.Println("     lango config list    — list all profiles")
	fmt.Println("     lango config use     — switch active profile")
}
