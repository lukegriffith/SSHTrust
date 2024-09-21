package handlers

import (
	"fmt"
	echo "github.com/labstack/echo/v4"
	"github.com/lukegriffith/SSHTrust/pkg/cert"
	"github.com/lukegriffith/SSHTrust/pkg/certStore"
	"golang.org/x/crypto/ssh"
	"net/http"
)

type App struct {
	Store certStore.CAStore
}

type SignRequest struct {
	PublicKey string `json:"public_key"`
}
type SignResponse struct {
	SignedKey string `json:"signed_key"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// GetCA retrieves a SSH Certificate Authority (CA) by its ID
// @Summary Get a SSH Certificate Authority (CA) by ID
// @Description Retrieve a CA by its ID from the applications store.
// @Tags CAs
// @Produce  json
// @Param id path string true "CA ID"
// @Success 200 {object} cert.CaResponse
// @Failure 404 {object} ErrorResponse "CA not found"
// @Router /CA/{id} [get]
func (a *App) GetCA(c echo.Context) error {
	// Parse the CA ID from the request
	CaID := c.Param("id")
	CA, err := a.Store.GetCAByID(CaID)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{"CA not found"})
	}

	// Return the CA in the response
	return c.JSON(http.StatusOK, CA)
}

// CreateCA creates a new SSH Certificate Authority (CA)
// @Summary Create a new SSH Certificate Authority (CA)
// @Description Create a new SSH CA and store it in the applications store.
// @Tags CAs
// @Accept  json
// @Produce  json
// @Param CA body cert.CaRequest true "New CA"
// @Success 201 {object} cert.CaResponse "The newly created CA"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 500 {object} ErrorResponse "Could not create CA"
// @Router /CA [post]
func (a *App) CreateCA(c echo.Context) error {
	var newCA cert.CaRequest

	c.Logger().Info("Creating new Key")
	// Bind the incoming JSON request to the CA model
	if err := c.Bind(&newCA); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Invalid request"})
	}

	c.Logger().Info("Input is bound")
	// Call the service to create the CA
	createdCA, err := a.Store.CreateCA(newCA)

	c.Logger().Info("CrateCA is called")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{fmt.Sprintf("Could not create CA: %s", err)})
	}

	c.Logger().Info("Created")
	// Return the newly created CA in the response
	return c.JSON(http.StatusCreated, createdCA)
}

// ListCA lists all Certificate Authorities (CAs)
// @Summary List all Certificate Authorities (CAs)
// @Description Retrieve a list of all CAs stored in the in-memory store.
// @Tags CAs
// @Produce  json
// @Success 200 {array} cert.CaResponse "List of all CAs"
// @Router /CA [get]
func (a *App) ListCA(c echo.Context) error {
	caList, _ := a.Store.ListCAs()
	return c.JSON(http.StatusOK, caList)
}

// Sign a public key using a specific CA
// @Summary Sign a public key with a specific CA
// @Description Use the specified CA to sign a provided public key and return the signed key.
// @Tags CAs
// @Accept  json
// @Produce  json
// @Param id path string true "CA ID"
// @Param public_key body SignRequest true "Public key to be signed"
// @Success 201 {object} SignResponse "The signed public key will be returned under the 'signed_key' field"
// @Failure 400 {object} ErrorResponse "Invalid request or failed to parse public key"
// @Failure 404 {object} ErrorResponse "CA not found"
// @Failure 500 {object} ErrorResponse "Failed to sign public key"
// @Router /CA/{id}/Sign [post]
func (a *App) Sign(c echo.Context) error {
	CaID := c.Param("id")
	signer, err := a.Store.GetSignerByID(CaID)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{"CA not found"})
	}
	// Parse the public key to be signed
	var requestBody = &SignRequest{}

	if err := c.Bind(requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Invalid request"})
	}

	parsedPublicKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(requestBody.PublicKey))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Failed to parse public key"})
	}

	// Sign the public key using the CA from the cert package
	// TODO extract valid principals
	signedCert, err := cert.SignUserKey(signer, parsedPublicKey, []string{"testuser"})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failed to sign public key"})
	}

	response := SignResponse{
		SignedKey: string(ssh.MarshalAuthorizedKey(signedCert)),
	}

	return c.JSON(http.StatusCreated, response)
}
