package main

import (
	"os"
	"path/filepath"

	"github.com/lkn/acme-crypt/internal/crypt"
)

func main() {
	if len(os.Args) != 2 {
		crypt.WriteToStderr("Usage: CryptGet <encrypted-file>\n")
		os.Exit(1)
	}

	encryptedFile := os.Args[1]
	if !filepath.IsAbs(encryptedFile) {
		wd, err := os.Getwd()
		if err != nil {
			crypt.WriteToStderr("Failed to get working directory: %v\n", err)
			os.Exit(1)
		}
		encryptedFile = filepath.Join(wd, encryptedFile)
	}

	if _, err := os.Stat(encryptedFile); os.IsNotExist(err) {
		crypt.WriteToStderr("File does not exist: %s\n", encryptedFile)
		os.Exit(1)
	}

	crypter, err := crypt.GetCrypter()
	if err != nil {
		crypt.WriteToStderr("Failed to initialize crypter: %v\n", err)
		os.Exit(1)
	}

	decryptedContent, err := crypter.Decrypt(encryptedFile)
	if err != nil {
		crypt.WriteToStderr("Failed to decrypt file: %v\n", err)
		os.Exit(1)
	}

	windowName := crypt.StripCryptExtension(encryptedFile)

	if err := crypt.CreateAcmeWindow(windowName, decryptedContent); err != nil {
		crypt.WriteToStderr("Failed to create acme window: %v\n", err)
		os.Exit(1)
	}
}