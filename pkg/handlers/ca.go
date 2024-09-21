package handlers

import (
	"fmt"
	echo "github.com/labstack/echo/v4"
	"github.com/lukegriffith/SSHTrust/pkg/cert"
	"github.com/lukegriffith/SSHTrust/pkg/certStore"
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

	// Bind the incoming JSON request to the CA model
	if err := c.Bind(&newCA); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Invalid request"})
	}

	// Call the service to create the CA
	createdCA, err := a.Store.CreateCA(newCA)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{fmt.Sprintf("Could not create CA: %s", err)})
	}

	c.Logger().Infof("created key %s", newCA.Name)
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
