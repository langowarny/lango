package prompt

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

// IsInteractive returns true if the standard input is a terminal
func IsInteractive() bool {
	return term.IsTerminal(int(syscall.Stdin))
}

// Passphrase prompts the user for a passphrase with hidden input
func Passphrase(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add newline after input
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

// PassphraseConfirm prompts for a passphrase and its confirmation
func PassphraseConfirm(prompt, confirmPrompt string) (string, error) {
	pass1, err := Passphrase(prompt)
	if err != nil {
		return "", err
	}

	pass2, err := Passphrase(confirmPrompt)
	if err != nil {
		return "", err
	}

	if pass1 != pass2 {
		return "", fmt.Errorf("passphrases do not match")
	}

	return pass1, nil
}
