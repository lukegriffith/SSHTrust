package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

// GenerateSSHKey generates a new SSH keypair with a 4096-bit RSA private key
func GenerateSSHKey() (ssh.Signer, error) {
	privateKey, err := generatePrivateKey()
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
func generatePrivateKey() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
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

// SignUserKey signs a user's public key using the CA private key.
// It returns a signed SSH certificate.
func SignUserKey(caSigner ssh.Signer, userPublicKey ssh.PublicKey) (*ssh.Certificate, error) {
	// Create a new certificate with the user's public key
	cert := &ssh.Certificate{
		Key:             userPublicKey,                                // The user's public key
		ValidPrincipals: []string{"testuser"},                         // Set the principal (can be a username)
		ValidAfter:      uint64(time.Now().Unix()),                    // Start time (now)
		ValidBefore:     uint64(time.Now().Add(1 * time.Hour).Unix()), // End time (1-hour TTL)
		CertType:        ssh.UserCert,                                 // Specify that this is a user certificate
		Permissions: ssh.Permissions{
			// Add any custom permissions or extensions as needed
			Extensions: map[string]string{
				"permit-pty": "", // Permit the user to request a PTY (useful for SSH terminal sessions)
			},
		},
	}

	// Sign the certificate using the CA's private key
	err := cert.SignCert(rand.Reader, caSigner)
	if err != nil {
		return nil, err
	}

	return cert, nil
}
