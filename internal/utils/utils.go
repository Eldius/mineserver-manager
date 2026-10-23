// Package utils holds small generic helpers shared across packages.
package utils

import (
	"fmt"
	"golang.org/x/term"
	"os/user"
	"path/filepath"
	"strings"
)

func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

// PasswordPrompt prompts user for a password
// (hiding it's value from console)
func PasswordPrompt() (string, error) {
	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(0)
	if err != nil {
		err = fmt.Errorf("password prompt: %w", err)
		return "", err
	}
	password := string(bytePassword)

	return strings.TrimSpace(password), nil
}

// ExpandPath expands tilde '~' character for home
func ExpandPath(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

// AbsolutePath returns the absolute path
func AbsolutePath(path string) (string, error) {
	path, err := ExpandPath(path)
	if err != nil {
		return "", fmt.Errorf("expanded path: %w", err)
	}
	return filepath.Abs(path)
}
