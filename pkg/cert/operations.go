package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"golang.org/x/crypto/ssh"
	"os"
)

// GenerateSSHKey generates a new SSH keypair with a 4096-bit RSA private key
func GenerateSSHKey(bits int) (ssh.Signer, error) {
	privateKey, err := generatePrivateKey(bits)
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
func generatePrivateKey(bits int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
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
