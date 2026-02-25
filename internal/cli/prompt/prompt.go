package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

// Confirm prompts the user for a yes/no confirmation and returns true for yes.
func Confirm(msg string) (bool, error) {
	fmt.Printf("%s [y/N]: ", msg)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	answer := strings.TrimSpace(strings.ToLower(line))
	return answer == "y" || answer == "yes", nil
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
