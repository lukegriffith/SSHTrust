package cert

import (
	"crypto/rand"
	"golang.org/x/crypto/ssh"
	"time"
)

// SignUserKey signs a user's public key using the CA private key.
// It returns a signed SSH certificate.
func SignUserKey(caSigner ssh.Signer, userPublicKey ssh.PublicKey, principals []string) (*ssh.Certificate, error) {
	// Create a new certificate with the user's public key
	cert := &ssh.Certificate{
		Key:             userPublicKey,                                // The user's public key
		ValidPrincipals: principals,                                   // Set the principal (can be a username)
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
