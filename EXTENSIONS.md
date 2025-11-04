# Extensions Guide

This document describes how to implement new encryption backends for acme-crypt.

## Architecture Overview

The acme-crypt tools use a pluggable architecture based on Go interfaces. New encryption backends can be added by implementing the `Crypter` interface defined in `internal/crypt/crypt.go`.

## Interface Definition

```go
type Crypter interface {
    Decrypt(encryptedData []byte) ([]byte, error)
    Encrypt(plainData []byte) ([]byte, error)
}
```

## Implementing a New Backend

### 1. Create the Backend

Create a new type that implements the `Crypter` interface:

```go
type MyEncryptionBackend struct {
    // Configuration fields
    keyPath string
    config  *MyConfig
}

func NewMyEncryptionBackend() *MyEncryptionBackend {
    return &MyEncryptionBackend{
        // Initialize from environment variables or config
    }
}

func (m *MyEncryptionBackend) Decrypt(encryptedData []byte) ([]byte, error) {
    // Implement decryption logic
    return decryptedData, nil
}

func (m *MyEncryptionBackend) Encrypt(plainData []byte) ([]byte, error) {
    // Implement encryption logic
    return encryptedData, nil
}
```

### 2. Add Backend Selection

Modify the main functions in `cmd/CryptGet/main.go` and `cmd/CryptPut/main.go` to support backend selection:

```go
func createCrypter() crypt.Crypter {
    backend := os.Getenv("ACME_CRYPT_BACKEND")
    switch backend {
    case "gpg", "":
        return crypt.NewGPGCrypter()
    case "age":
        return crypt.NewAgeCrypter()
    case "openssl":
        return crypt.NewOpenSSLCrypter()
    default:
        fmt.Fprintf(os.Stderr, "Unknown backend: %s\n", backend)
        os.Exit(1)
    }
}
```

### 3. Handle File Extensions

Update the extension handling in `GetOutputFilename()` to recognize your backend's file extensions:

```go
func GetOutputFilename(filename string) string {
    extensions := []string{".gpg", ".asc", ".pgp", ".age", ".enc"}
    for _, ext := range extensions {
        if strings.HasSuffix(filename, ext) {
            return strings.TrimSuffix(filename, ext)
        }
    }
    return filename
}
```

## Configuration

### Environment Variables

Backends should use environment variables for configuration:

- `ACME_CRYPT_BACKEND` - Backend selection (e.g., "gpg", "age", "openssl")
- Backend-specific variables (e.g., `ACME_CRYPT_RCPT` for GPG recipient)

### Example Backend Configurations

#### GPG Backend (Current)
```bash
export ACME_CRYPT_BACKEND="gpg"
export ACME_CRYPT_RCPT="user@example.com"
```

#### Hypothetical Age Backend
```bash
export ACME_CRYPT_BACKEND="age"
export ACME_CRYPT_AGE_RECIPIENTS_FILE="$HOME/.config/age/recipients"
```

#### Hypothetical OpenSSL Backend
```bash
export ACME_CRYPT_BACKEND="openssl"
export ACME_CRYPT_CERT_FILE="$HOME/.config/crypt/cert.pem"
export ACME_CRYPT_KEY_FILE="$HOME/.config/crypt/key.pem"
```

## Error Handling

Backends should:
- Return descriptive errors for common failure cases
- Preserve environment variables needed for GUI operations (like `DISPLAY`)
- Handle missing configuration gracefully
- Use stderr for error output (goes to Acme's +Errors buffer)

## Testing Your Backend

1. **Build the tools**: `mk install`
2. **Set environment variables** for your backend
3. **Test decryption**: Create an encrypted file, highlight it in Acme, run `CryptGet`
4. **Test encryption**: Edit decrypted content, run `CryptPut`
5. **Verify roundtrip**: Ensure the encrypted file can be decrypted again

## Example: Age Backend Implementation

```go
package crypt

import (
    "bytes"
    "os"
    "os/exec"
)

type AgeCrypter struct {
    recipientsFile string
}

func NewAgeCrypter() *AgeCrypter {
    recipientsFile := os.Getenv("ACME_CRYPT_AGE_RECIPIENTS_FILE")
    if recipientsFile == "" {
        recipientsFile = os.Getenv("HOME") + "/.config/age/recipients"
    }
    return &AgeCrypter{recipientsFile: recipientsFile}
}

func (a *AgeCrypter) Decrypt(encryptedData []byte) ([]byte, error) {
    cmd := exec.Command("age", "--decrypt")
    cmd.Stdin = bytes.NewReader(encryptedData)
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    if err := cmd.Run(); err != nil {
        return nil, fmt.Errorf("age decrypt failed: %v\nstderr: %s", err, stderr.String())
    }
    
    return stdout.Bytes(), nil
}

func (a *AgeCrypter) Encrypt(plainData []byte) ([]byte, error) {
    cmd := exec.Command("age", "--recipients-file", a.recipientsFile)
    cmd.Stdin = bytes.NewReader(plainData)
    
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    if err := cmd.Run(); err != nil {
        return nil, fmt.Errorf("age encrypt failed: %v\nstderr: %s", err, stderr.String())
    }
    
    return stdout.Bytes(), nil
}
```

## Contributing

When adding new backends:
1. Follow the existing code style and patterns
2. Add appropriate error handling and logging
3. Update documentation and examples
4. Test thoroughly with Acme integration
5. Consider security implications of your implementation