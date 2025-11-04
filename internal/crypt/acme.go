package crypt

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"9fans.net/go/acme"
)

func CreateAcmeWindow(name string, content []byte) error {
	// Create new acme window
	w, err := acme.New()
	if err != nil {
		return fmt.Errorf("failed to create acme window: %w", err)
	}

	// Set window name
	if err := w.Name("%s", name); err != nil {
		w.Del(true) // Clean up on error
		return fmt.Errorf("failed to set window name: %w", err)
	}

	// Append CryptPut to the window tag
	if _, err := w.Write("tag", []byte(" CryptPut")); err != nil {
		w.Del(true) // Clean up on error
		return fmt.Errorf("failed to append CryptPut to tag: %w", err)
	}

	// Write content to window body
	if _, err := w.Write("body", content); err != nil {
		w.Del(true) // Clean up on error
		return fmt.Errorf("failed to write to acme window body: %w", err)
	}

	return nil
}

func GetCurrentAcmeWindowContent() ([]byte, error) {
	winIDStr := os.Getenv("winid")
	if winIDStr == "" {
		return nil, fmt.Errorf("not running in acme window (winid not set)")
	}

	winID, err := strconv.Atoi(winIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid winid: %w", err)
	}

	// Open the current window by ID
	w, err := acme.Open(winID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open acme window: %w", err)
	}
	defer w.CloseFiles()

	// Read content from window body
	content, err := w.ReadAll("body")
	if err != nil {
		return nil, fmt.Errorf("failed to read acme window body: %w", err)
	}

	return content, nil
}

func GetCurrentAcmeWindowName() (string, error) {
	winIDStr := os.Getenv("winid")
	if winIDStr == "" {
		return "", fmt.Errorf("not running in acme window (winid not set)")
	}

	winID, err := strconv.Atoi(winIDStr)
	if err != nil {
		return "", fmt.Errorf("invalid winid: %w", err)
	}

	// Open the current window by ID
	w, err := acme.Open(winID, nil)
	if err != nil {
		return "", fmt.Errorf("failed to open acme window: %w", err)
	}
	defer w.CloseFiles()

	// Read the window tag
	tag, err := w.ReadAll("tag")
	if err != nil {
		return "", fmt.Errorf("failed to read acme window tag: %w", err)
	}

	tagFields := strings.Fields(string(tag))
	if len(tagFields) == 0 {
		return "", fmt.Errorf("empty tag")
	}

	windowName := tagFields[0]
	if !filepath.IsAbs(windowName) {
		wd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get working directory: %w", err)
		}
		windowName = filepath.Join(wd, windowName)
	}

	return windowName, nil
}

func WriteToStderr(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}