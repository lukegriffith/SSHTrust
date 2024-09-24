package cert

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"os"

	"golang.org/x/crypto/ssh"
)

type KeyType string

const (
	RSAKey  KeyType = "rsa"
	ED25519 KeyType = "ed25519"
)

var InvalidKeyErr error = errors.New("unsupported key type")

// GenerateSSHKey generates a new SSH keypair with a 4096-bit RSA private key
func GenerateSSHKey(keyType KeyType, bits int) (ssh.Signer, error) {
	privateKey, err := generatePrivateKey(keyType, bits)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	// Create an SSH signer using the generated private key
	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		return nil, err
	}
	return signer, nil
}

// generatePrivateKey generates a new RSA private key
func generatePrivateKey(keyType KeyType, bits int) (interface{}, error) {
	switch keyType {
	case RSAKey:
		privateKey, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			return nil, err
		}
		return privateKey, nil
	case ED25519:
		_, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		return privateKey, nil
	default:
		return nil, InvalidKeyErr
	}
}

// SavePublicKey saves the public key from the SSH signer to a file
func SavePublicKey(signer ssh.Signer, filePath string) error {
	publicKey := ssh.MarshalAuthorizedKey(signer.PublicKey())
	err := os.WriteFile(filePath, publicKey, 0644)
	if err != nil {
		return err
	}
	return nil
}
