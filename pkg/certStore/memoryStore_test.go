package certStore

import (
	"testing"

	"github.com/lukegriffith/SSHTrust/pkg/cert"
	"github.com/stretchr/testify/assert"
)

// Test for CreateCA success case
func TestCreateCASuccess(t *testing.T) {
	store := NewInMemoryCaStore()
	// Using the actual CA struct instead of mockCA
	mockRequest := cert.CaRequest{Name: "test-ca", Type: "ssh-rsa", Bits: 2048, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}

	ca, err := store.CreateCA(mockRequest)
	assert.NoError(t, err, "Expected no error when creating CA")
	assert.NotNil(t, ca, "CA should not be nil")
	assert.Equal(t, "test-ca", ca.Name, "CA name should match")
	assert.Equal(t, []string{"testuser"}, ca.ValidPrincipals, "Principals should match")
	assert.Equal(t, 3600, ca.MaxTTLMinutes, "TTL should match")
}

// Test for CreateCA failure case (duplicate CA)
func TestCreateCADuplicate(t *testing.T) {
	store := &InMemortCaStore{
		cas: make(map[string]cert.CA),
	}
	mockRequest := cert.CaRequest{Name: "test-ca", Type: "ssh-rsa", Bits: 2048, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}

	// First CA creation should succeed
	_, err := store.CreateCA(mockRequest)
	assert.NoError(t, err, "Expected no error on first CA creation")

	// Second CA creation with the same name should fail
	_, err = store.CreateCA(mockRequest)
	assert.Error(t, err, "Expected error when creating duplicate CA")
	assert.Equal(t, "CA already exists", err.Error(), "Expected 'CA already exists' error")
}

// Test for GetCAByID success case
func TestGetCAByIDSuccess(t *testing.T) {
	store := &InMemortCaStore{
		cas: make(map[string]cert.CA),
	}

	// Create CA and add it to the store
	mockRequest := cert.CaRequest{Name: "test-ca", Type: "ssh-rsa", Bits: 2048, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}
	ca, _ := store.CreateCA(mockRequest)

	// Retrieve CA by ID
	retrievedCA, err := store.GetCAByID("test-ca")
	assert.NoError(t, err, "Expected no error when retrieving CA by ID")
	assert.NotNil(t, retrievedCA, "Retrieved CA should not be nil")
	assert.Equal(t, ca.Name, retrievedCA.Name, "Expected CA names to match")
}

// Test for GetCAByID failure case (non-existent CA)
func TestGetCAByIDNotFound(t *testing.T) {
	store := &InMemortCaStore{
		cas: make(map[string]cert.CA),
	}

	// Attempt to retrieve non-existent CA
	retrievedCA, err := store.GetCAByID("non-existent-ca")
	assert.Error(t, err, "Expected error when retrieving non-existent CA")
	assert.Nil(t, retrievedCA, "Retrieved CA should be nil")
	assert.Equal(t, "unable to find CA by ID", err.Error(), "Expected 'unable to find CA by ID' error")
}

// Test for GetSignerByID success case
func TestGetSignerByIDSuccess(t *testing.T) {
	store := &InMemortCaStore{
		cas: make(map[string]cert.CA),
	}

	// Create CA and add it to the store
	mockRequest := cert.CaRequest{Name: "test-ca", Type: "ssh-rsa", Bits: 2048, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}
	store.CreateCA(mockRequest)

	// Retrieve Signer by ID
	signer, err := store.GetSignerByID("test-ca")
	assert.NoError(t, err, "Expected no error when retrieving signer by ID")
	assert.NotNil(t, signer, "Signer should not be nil")
}

// Test for GetSignerByID failure case (non-existent CA)
func TestGetSignerByIDNotFound(t *testing.T) {
	store := &InMemortCaStore{
		cas: make(map[string]cert.CA),
	}

	// Attempt to retrieve signer for non-existent CA
	signer, err := store.GetSignerByID("non-existent-ca")
	assert.Error(t, err, "Expected error when retrieving signer for non-existent CA")
	assert.Nil(t, signer, "Signer should be nil")
	assert.Equal(t, "Unable to find CA by ID", err.Error(), "Expected 'Unable to find CA by ID' error")
}

// Test for ListCAs
func TestListCAs(t *testing.T) {
	store := &InMemortCaStore{
		cas: make(map[string]cert.CA),
	}

	// Add two CAs to the store
	mockRequest1 := cert.CaRequest{Name: "test-ca1", Type: "ssh-rsa", Bits: 2048, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}
	mockRequest2 := cert.CaRequest{Name: "test-ca2", Type: "ssh-rsa", Bits: 2048, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}
	store.CreateCA(mockRequest1)
	store.CreateCA(mockRequest2)

	// List all CAs
	cas, err := store.ListCAs()
	assert.NoError(t, err, "Expected no error when listing CAs")
	assert.Len(t, cas, 2, "Expected two CAs to be returned")

	caMap := map[string]interface{}{
		cas[0].Name: nil,
		cas[1].Name: nil,
	}

	// Check that both "test-ca1" and "test-ca2" are present
	_, ca1Exists := caMap["test-ca1"]
	_, ca2Exists := caMap["test-ca2"]

	assert.True(t, ca1Exists, "Expected test-ca1 to be present")
	assert.True(t, ca2Exists, "Expected test-ca2 to be present")
}
