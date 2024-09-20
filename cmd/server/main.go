package main

import (
	"encoding/json"
	"fmt"
	"github.com/lukegriffith/ssh-key-server/pkg/cert" // Import your cert package
	"golang.org/x/crypto/ssh"
	"log"
	"net/http"
	"sync"
)

// In-memory storage for CAs
var caStore = struct {
	sync.RWMutex
	cas map[string]ssh.Signer
}{cas: make(map[string]ssh.Signer)}

// Endpoint 1: Create a new CA by name and store it in memory
func createCAHandler(w http.ResponseWriter, r *http.Request) {
	caName := r.URL.Query().Get("name")
	if caName == "" {
		http.Error(w, "Missing CA name", http.StatusBadRequest)
		return
	}

	// Use the existing GenerateSSHKey function from the cert package
	signer, err := cert.GenerateSSHKey()
	if err != nil {
		http.Error(w, "Failed to generate CA keypair", http.StatusInternalServerError)
		return
	}

	// Store the CA in memory
	caStore.Lock()
	caStore.cas[caName] = signer
	caStore.Unlock()

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "CA '%s' created successfully", caName)
}

// Endpoint 2: Obtain the CA's public key
func getCAPublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	caName := r.URL.Query().Get("name")
	if caName == "" {
		http.Error(w, "Missing CA name", http.StatusBadRequest)
		return
	}

	// Get the CA from the in-memory store
	caStore.RLock()
	signer, exists := caStore.cas[caName]
	caStore.RUnlock()

	if !exists {
		http.Error(w, "CA not found", http.StatusNotFound)
		return
	}

	// Marshal the public key into a format that can be returned
	publicKey := ssh.MarshalAuthorizedKey(signer.PublicKey())
	w.WriteHeader(http.StatusOK)
	w.Write(publicKey)
}

// Endpoint 3: Sign a public key with the CA
func signPublicKeyHandler(w http.ResponseWriter, r *http.Request) {
	caName := r.URL.Query().Get("ca")
	if caName == "" {
		http.Error(w, "Missing CA name", http.StatusBadRequest)
		return
	}

	// Get the CA from the in-memory store
	caStore.RLock()
	signer, exists := caStore.cas[caName]
	caStore.RUnlock()

	if !exists {
		http.Error(w, "CA not found", http.StatusNotFound)
		return
	}

	// Parse the public key to be signed
	var requestBody struct {
		PublicKey string `json:"public_key"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil || requestBody.PublicKey == "" {
		http.Error(w, "Invalid public key format", http.StatusBadRequest)
		return
	}

	parsedPublicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(requestBody.PublicKey))
	if err != nil {
		http.Error(w, "Failed to parse public key", http.StatusBadRequest)
		return
	}

	// Sign the public key using the CA from the cert package
	signedCert, err := cert.SignUserKey(signer, parsedPublicKey)
	if err != nil {
		http.Error(w, "Failed to sign public key", http.StatusInternalServerError)
		return
	}

	// Return the signed certificate
	signedCertBytes := ssh.MarshalAuthorizedKey(signedCert)
	w.WriteHeader(http.StatusOK)
	w.Write(signedCertBytes)
}

// HTTP server setup
func main() {
	http.HandleFunc("/create-ca", createCAHandler)
	http.HandleFunc("/get-ca-public-key", getCAPublicKeyHandler)
	http.HandleFunc("/sign-public-key", signPublicKeyHandler)

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

