package handlers

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/lukegriffith/SSHTrust/pkg/cert"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
)

// Mock the CAStore for testing purposes
type MockStore struct {
	caMap   map[string]*cert.CaResponse
	signers map[string]ssh.Signer
}

func (m *MockStore) GetCAByID(ID string) (*cert.CaResponse, error) {
	if ca, exists := m.caMap[ID]; exists {
		return ca, nil
	}
	return nil, errors.New("CA not found")
}

func (m *MockStore) CreateCA(ca cert.CaRequest) (*cert.CaResponse, error) {
	if _, exists := m.caMap[ca.Name]; exists {
		return nil, errors.New("CA already exists")
	}
	resp := &cert.CaResponse{
		Name: ca.Name,
		Type: ca.Type,
		Bits: ca.Bits,
	}
	m.caMap[ca.Name] = resp
	return resp, nil
}

func (m *MockStore) ListCAs() ([]*cert.CaResponse, error) {
	cas := []*cert.CaResponse{}
	for _, ca := range m.caMap {
		cas = append(cas, ca)
	}
	return cas, nil
}

func (m *MockStore) GetSignerByID(ID string) (ssh.Signer, error) {
	// Mock implementation for signers
	return m.signers[ID], nil
}

// Test for GetCA handler
func TestGetCAHandler(t *testing.T) {
	e := echo.New()

	mockStore := &MockStore{
		caMap: map[string]*cert.CaResponse{
			"test-ca": {Name: "test-ca", PublicKey: "test-public-key"},
		},
	}
	app := &App{Store: mockStore}

	req := httptest.NewRequest(http.MethodGet, "/ca/test-ca", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("test-ca")

	// Call the handler
	err := app.GetCA(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "test-ca")
	}
}

// Test for GetCA handler when CA is not found
func TestGetCANotFoundHandler(t *testing.T) {
	e := echo.New()

	mockStore := &MockStore{
		caMap: map[string]*cert.CaResponse{},
	}
	app := &App{Store: mockStore}

	req := httptest.NewRequest(http.MethodGet, "/ca/nonexistent-ca", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nonexistent-ca")

	// Call the handler
	err := app.GetCA(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "CA not found")
	}
}

// Test for CreateCA handler
func TestCreateCAHandler(t *testing.T) {
	e := echo.New()

	mockStore := &MockStore{
		caMap: map[string]*cert.CaResponse{},
	}
	app := &App{Store: mockStore}

	// Prepare a valid CA creation request
	newCA := cert.CA{Name: "new-ca"}
	reqBody, _ := json.Marshal(newCA)
	req := httptest.NewRequest(http.MethodPost, "/ca", strings.NewReader(string(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err := app.CreateCA(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "new-ca")
	}
}

// Test for CreateCA handler with invalid request
func TestCreateCAInvalidRequestHandler(t *testing.T) {
	e := echo.New()

	mockStore := &MockStore{
		caMap: map[string]*cert.CaResponse{},
	}
	app := &App{Store: mockStore}

	// Prepare an invalid JSON request (malformed JSON)
	req := httptest.NewRequest(http.MethodPost, "/ca", strings.NewReader("invalid-json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err := app.CreateCA(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid request")
	}
}

// Test for ListCA handler
func TestListCAHandler(t *testing.T) {
	e := echo.New()

	mockStore := &MockStore{
		caMap: map[string]*cert.CaResponse{
			"ca1": {Name: "ca1", PublicKey: "publickey1"},
			"ca2": {Name: "ca2", PublicKey: "publickey2"},
		},
	}
	app := &App{Store: mockStore}

	req := httptest.NewRequest(http.MethodGet, "/cas", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the handler
	err := app.ListCA(c)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "ca1")
		assert.Contains(t, rec.Body.String(), "ca2")
	}
}
