# acme-crypt

Go wrapper programs for Acme that interface with encryption/decryption tools to automatically handle encrypted files. **Currently only supports GPG**, though designed with a pluggable interface for future encryption backends.

## Overview

This package provides two commands:
- `CryptGet` - Decrypts an encrypted file and opens it in a new Acme window
- `CryptPut` - Encrypts the contents of an Acme window and saves to disk

## Installation

1. Clone the repository:
   ```bash
   git clone git@github.com:lneely/acme-crypt.git
   cd acme-crypt
   ```

2. Install the binaries:
   ```bash
   mk install
   ```

This installs `CryptGet` and `CryptPut` to `$HOME/bin/`.

## Configuration

Set your GPG recipient in the environment:
```bash
export ACME_CRYPT_RCPT="your@email.com"
```

Ensure your GPG agent is configured with a GUI pinentry program (e.g., pinentry-qt, pinentry-gtk3) since Acme runs in a non-interactive environment.

## Usage

### Decrypting Files (CryptGet)

1. In Acme, highlight an encrypted file (`.gpg`, `.asc`, or `.pgp`)
2. Add `CryptGet` to the window tag
3. Highlight the encrypted filename, hold the middle button then left click on CryptGet to decrypt and open the file.
4. A new window opens with the decrypted content, named after the original file without the encryption extension
5. The new window automatically has `CryptPut` added to its tag.

### Encrypting Files (CryptPut)

#### From Acme Window
1. Edit the decrypted file content in the Acme window created by CryptGet
2. Middle-click on `CryptPut` in the window tag
3. The content is encrypted and saved back to the original encrypted file path

#### From Stdin (New Files)
Create new encrypted files by piping content to CryptPut:
```bash
echo "This is my content" | CryptPut /path/to/file
```
This encrypts the stdin content and saves it to `/path/to/file.gpg`

## Workflow Example

```
1. File: sensitive-data.txt.gpg
2. Highlight "sensitive-data.txt.gpg" in Acme
3. Middle-click CryptGet → Opens window "sensitive-data.txt" with decrypted content
4. Edit the content as needed
5. Middle-click CryptPut → Saves encrypted content back to sensitive-data.txt.gpg
```

## Supported File Extensions

CryptGet automatically strips these extensions when creating the editing window:
- `.gpg`
- `.asc`
- `.pgp`

## Architecture

The program uses a pluggable interface design in `internal/crypt/`:
- `Crypter` interface for encryption/decryption operations
- `GPGCrypter` implementation for GPG backend
- Acme integration utilities for window management

This design allows for easy addition of other encryption tools (OpenSSL, age, etc.) in the future, though currently only GPG is implemented.

## Requirements

- Go 1.19+
- Acme text editor
- GPG with configured keys
- GUI pinentry program (pinentry-qt, pinentry-gtk3, etc.)

## License

GNU General Public License v3