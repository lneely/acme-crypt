package crypt

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Crypter interface {
	Decrypt(filePath string) ([]byte, error)
	Encrypt(data []byte, outputPath string) error
}

// Backend registry
var backends = map[string]func() (Crypter, error){
	"gpg": func() (Crypter, error) { return NewGPGCrypter() },
}

// GetCrypter returns a Crypter instance based on ACME_CRYPT_BACKEND environment variable
// Defaults to GPG if not set
func GetCrypter() (Crypter, error) {
	backend := os.Getenv("ACME_CRYPT_BACKEND")
	if backend == "" {
		backend = "gpg" // Default to GPG
	}
	
	constructor, exists := backends[backend]
	if !exists {
		return nil, fmt.Errorf("unsupported backend: %s", backend)
	}
	
	return constructor()
}

// RegisterBackend allows registration of new backend types
func RegisterBackend(name string, constructor func() (Crypter, error)) {
	backends[name] = constructor
}

type GPGCrypter struct {
	recipient string
}

func NewGPGCrypter() (*GPGCrypter, error) {
	recipient := os.Getenv("ACME_CRYPT_RCPT")
	if recipient == "" {
		return nil, fmt.Errorf("ACME_CRYPT_RCPT environment variable not set")
	}
	return &GPGCrypter{recipient: recipient}, nil
}

func (g *GPGCrypter) Decrypt(filePath string) ([]byte, error) {
	cmd := exec.Command("gpg", "--decrypt", "--quiet", "--batch", "--no-tty", filePath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	// Preserve DISPLAY for GUI pinentry
	cmd.Env = append(os.Environ(), "DISPLAY="+os.Getenv("DISPLAY"))
	
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("gpg decrypt failed: %w\nstderr: %s", err, stderr.String())
	}
	
	return stdout.Bytes(), nil
}

func (g *GPGCrypter) Encrypt(data []byte, outputPath string) error {
	// Remove existing file if it exists
	if _, err := os.Stat(outputPath); err == nil {
		if err := os.Remove(outputPath); err != nil {
			return fmt.Errorf("failed to remove existing file %s: %w", outputPath, err)
		}
	}
	
	cmd := exec.Command("gpg", "--encrypt", "--armor", "--recipient", g.recipient, "--batch", "--no-tty", "--output", outputPath)
	cmd.Stdin = bytes.NewReader(data)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	// Preserve DISPLAY for GUI pinentry
	cmd.Env = append(os.Environ(), "DISPLAY="+os.Getenv("DISPLAY"))
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gpg encrypt failed: %w\nstderr: %s", err, stderr.String())
	}
	
	return nil
}

func StripCryptExtension(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".gpg" || ext == ".asc" || ext == ".pgp" {
		return strings.TrimSuffix(filePath, ext)
	}
	return filePath
}

func AddGPGExtension(filePath string) string {
	if !strings.HasSuffix(strings.ToLower(filePath), ".gpg") {
		return filePath + ".gpg"
	}
	return filePath
}