package cert

import (
	"os"
	"testing"
)

// TestGenerateSSHKey tests the generation of an SSH keypair
func TestGenerateSSHKey(t *testing.T) {
	signer, err := GenerateSSHKey()
	if err != nil {
		t.Fatalf("Failed to generate SSH keypair: %v", err)
	}

	if signer == nil {
		t.Fatal("Expected signer, got nil")
	}

	// Validate that the generated key is a valid SSH public key
	publicKey := signer.PublicKey()
	if publicKey == nil {
		t.Fatalf("Generated key is not a valid SSH public key: %v", err)
	}
}

// TestSavePublicKey tests the saving of the public key to a file
func TestSavePublicKey(t *testing.T) {
	// Generate a test signer (SSH keypair)
	signer, err := GenerateSSHKey()
	if err != nil {
		t.Fatalf("Failed to generate SSH keypair: %v", err)
	}

	// Define a temporary file path to save the public key
	filePath := "test_ssh_ca.pub"

	// Call SavePublicKey and check if it writes the file correctly
	err = SavePublicKey(signer, filePath)
	if err != nil {
		t.Fatalf("Failed to save public key: %v", err)
	}

	// Check if the file was created
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Expected public key file %s to be created, but it doesn't exist", filePath)
	}

	// Clean up by removing the test file
	os.Remove(filePath)
}
