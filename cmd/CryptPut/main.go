package main

import (
	"io"
	"os"

	"github.com/lkn/acme-crypt/internal/crypt"
)

func main() {
	var content []byte
	var outputPath string
	var err error

	// Check if we have a file path argument (stdin mode)
	if len(os.Args) > 1 {
		// Read from stdin
		content, err = io.ReadAll(os.Stdin)
		if err != nil {
			crypt.WriteToStderr("Failed to read from stdin: %v\n", err)
			os.Exit(1)
		}

		// Use the provided file path and add GPG extension
		outputPath = crypt.AddGPGExtension(os.Args[1])
	} else {
		// Original acme window mode
		content, err = crypt.GetCurrentAcmeWindowContent()
		if err != nil {
			crypt.WriteToStderr("Failed to get acme window content: %v\n", err)
			os.Exit(1)
		}

		windowName, err := crypt.GetCurrentAcmeWindowName()
		if err != nil {
			crypt.WriteToStderr("Failed to get acme window name: %v\n", err)
			os.Exit(1)
		}

		outputPath = crypt.AddGPGExtension(windowName)
	}

	crypter, err := crypt.GetCrypter()
	if err != nil {
		crypt.WriteToStderr("Failed to initialize crypter: %v\n", err)
		os.Exit(1)
	}

	if err := crypter.Encrypt(content, outputPath); err != nil {
		crypt.WriteToStderr("Failed to encrypt and save file: %v\n", err)
		os.Exit(1)
	}

	crypt.WriteToStderr("Encrypted and saved to: %s\n", outputPath)
}