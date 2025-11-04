package main

import (
	"os"

	"github.com/lkn/acme-crypt/internal/crypt"
)

func main() {
	content, err := crypt.GetCurrentAcmeWindowContent()
	if err != nil {
		crypt.WriteToStderr("Failed to get acme window content: %v\n", err)
		os.Exit(1)
	}

	windowName, err := crypt.GetCurrentAcmeWindowName()
	if err != nil {
		crypt.WriteToStderr("Failed to get acme window name: %v\n", err)
		os.Exit(1)
	}

	crypter, err := crypt.GetCrypter()
	if err != nil {
		crypt.WriteToStderr("Failed to initialize crypter: %v\n", err)
		os.Exit(1)
	}

	outputPath := crypt.AddGPGExtension(windowName)

	if err := crypter.Encrypt(content, outputPath); err != nil {
		crypt.WriteToStderr("Failed to encrypt and save file: %v\n", err)
		os.Exit(1)
	}

	crypt.WriteToStderr("Encrypted and saved to: %s\n", outputPath)
}