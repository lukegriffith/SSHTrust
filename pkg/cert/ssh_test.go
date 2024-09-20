package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"testing"
)

// Helper function to start a lightweight SSH server in a separate goroutine
func startSSHServer(caPublicKey ssh.PublicKey) (net.Listener, error) {
	// Configure the SSH server to trust the provided CA public key
	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			cert, ok := pubKey.(*ssh.Certificate)
			if !ok {
				log.Printf("Received public key that is not a certificate.")
				return nil, ssh.ErrNoAuth
			}
			// Compare the SignatureKey with the trusted CA public key using bytes.Equal
			if bytes.Equal(ssh.MarshalAuthorizedKey(cert.SignatureKey), ssh.MarshalAuthorizedKey(caPublicKey)) {
				log.Printf("Successful certificate validation for user: %s", conn.User())
				return nil, nil
			}
			log.Printf("Certificate signed by unknown authority for user: %s", conn.User())
			return nil, ssh.ErrNoAuth
		},
	}

	// Generate a private key for the SSH server dynamically
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Create an SSH signer from the private key
	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		return nil, err
	}

	// Add the generated private key as the host key for the server
	config.AddHostKey(signer)

	// Start listening for SSH connections on a random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			// Accept incoming connections
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error accepting connection: %v", err)
				return
			}

			// Handle the incoming SSH connection
			sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
			if err != nil {
				log.Printf("Failed to establish SSH connection: %v", err)
				continue
			}
			log.Printf("New SSH connection established from %s", sshConn.RemoteAddr())

			// Handle requests (this is where session requests will come in)
			go ssh.DiscardRequests(reqs)

			// Handle incoming channels (e.g., session channels)
			for newChannel := range chans {
				if newChannel.ChannelType() == "session" {
					channel, requests, err := newChannel.Accept()
					if err != nil {
						log.Printf("Could not accept channel: %v", err)
						continue
					}

					// Handle the session requests like exec, shell, etc.
					go func(in <-chan *ssh.Request) {
						for req := range in {
							switch req.Type {
							case "exec":
								log.Println("Exec request received")
								// Run the command (in this case, we just respond with success)
								req.Reply(true, nil)
								_, err := io.WriteString(channel, "success\n")
								if err != nil {
									log.Printf("Failed to write to channel: %v", err)
								}

								// Send exit status
								channel.SendRequest("exit-status", false, ssh.Marshal(&struct{ ExitStatus uint32 }{0}))

								// Close the channel
								channel.Close()
							default:
								req.Reply(false, nil)
							}
						}
					}(requests)
				} else {
					newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
					log.Printf("Unknown channel type: %s", newChannel.ChannelType())
				}
			}

			sshConn.Close()
		}
	}()

	return listener, nil
}

// TestGenerateSSHCertificateEndToEnd tests end-to-end SSH certificate signing and server validation
func TestGenerateSSHCertificateEndToEnd(t *testing.T) {
	// Step 1: Generate the SSH CA keypair (this is the key that will sign certificates)
	sshCA, err := GenerateSSHKey()
	if err != nil {
		t.Fatalf("Failed to generate CA keypair: %v", err)
	}

	// Step 2: Save the public key for the CA (to simulate it being known to the server)
	err = SavePublicKey(sshCA, "test_ca.pub")
	if err != nil {
		t.Fatalf("Failed to save CA public key: %v", err)
	}

	// Step 3: Load the saved CA public key for use by the SSH server
	caPublicKeyBytes, err := os.ReadFile("test_ca.pub")
	if err != nil {
		t.Fatalf("Failed to read CA public key file: %v", err)
	}
	caPublicKey, _, _, _, err := ssh.ParseAuthorizedKey(caPublicKeyBytes)
	if err != nil {
		t.Fatalf("Failed to parse CA public key: %v", err)
	}

	// Step 4: Start the SSH server
	listener, err := startSSHServer(caPublicKey)
	if err != nil {
		t.Fatalf("Failed to start SSH server: %v", err)
	}
	defer listener.Close()

	// Step 5: Generate a user keypair (this will simulate a user's public key)
	userPrivateKey, err := GenerateSSHKey()
	if err != nil {
		t.Fatalf("Failed to generate user keypair: %v", err)
	}

	// Step 6: Sign the user's public key with the CA
	signedCert, err := SignUserKey(sshCA, userPrivateKey.PublicKey())
	if err != nil {
		t.Fatalf("Failed to sign user's public key: %v", err)
	}

	// Step 7: Create an SSH client that connects to the server using the signed certificate
	clientConfig := &ssh.ClientConfig{
		User: "testuser",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(&signedCertSigner{
				signer: userPrivateKey, // The user's private key (to sign the request)
				cert:   signedCert,     // The signed certificate
			}),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Ignoring server host key for simplicity in testing
	}

	// Step 8: Connect to the SSH server
	address := listener.Addr().String()
	client, err := ssh.Dial("tcp", address, clientConfig)
	if err != nil {
		t.Fatalf("Failed to connect to SSH server: %v", err)
	}
	defer client.Close()

	// Step 9: Validate the client is accepted by the server
	session, err := client.NewSession()
	if err != nil {
		t.Fatalf("Failed to create SSH session: %v", err)
	}
	defer session.Close()

	// Run a basic command to validate that the server accepts the certificate
	err = session.Run("echo success")
	if err != nil {
		t.Fatalf("Failed to run command on SSH server: %v", err)
	}
}

// signedCertSigner is a custom signer that includes both the private key and the signed certificate
type signedCertSigner struct {
	signer ssh.Signer
	cert   *ssh.Certificate
}

func (s *signedCertSigner) PublicKey() ssh.PublicKey {
	return s.cert
}

func (s *signedCertSigner) Sign(rand io.Reader, data []byte) (*ssh.Signature, error) {
	return s.signer.Sign(rand, data)
}
