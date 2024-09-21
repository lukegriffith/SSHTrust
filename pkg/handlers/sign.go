package handlers

import (
	echo "github.com/labstack/echo/v4"
	"github.com/lukegriffith/SSHTrust/pkg/cert"
	"golang.org/x/crypto/ssh"
	"net/http"
)

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

	parsedPublicKey, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(requestBody.PublicKey))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Failed to parse public key"})
	}

	// Sign the public key using the CA from the cert package
	// TODO extract valid principals
	signedCert, err := cert.SignUserKey(signer, parsedPublicKey, []string{"testuser"})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failed to sign public key"})
	}
	c.Logger().Infof("Signed public key %s for %s", comment, CaID)
	response := SignResponse{
		SignedKey: string(ssh.MarshalAuthorizedKey(signedCert)),
	}

	return c.JSON(http.StatusCreated, response)
}
