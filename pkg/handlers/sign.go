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
// @Param public_key body cert.SignRequest true "Public key to be signed"
// @Success 201 {object} cert.SignResponse "The signed public key will be returned under the 'signed_key' field"
// @Failure 400 {object} ErrorResponse "Invalid request or failed to parse public key"
// @Failure 404 {object} ErrorResponse "CA not found"
// @Failure 404 {object} ErrorResponse "Requested TTL longer than configured max"
// @Failure 404 {object} ErrorResponse "Requested principals not in valid principal list"
// @Failure 500 {object} ErrorResponse "Failed to sign public key"
// @Router /CA/{id}/Sign [post]
func (a *App) Sign(c echo.Context) error {
	CaID := c.Param("id")
	ca, err := a.Store.GetCAByID(CaID)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{"CA not found"})
	}
	signer, err := a.Store.GetSignerByID(CaID)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{"CA signer not found"})
	}
	// Parse the public key to be signed
	var requestBody = &cert.SignRequest{}

	if err := c.Bind(requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Invalid request"})
	}

	parsedPublicKey, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(requestBody.PublicKey))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Failed to parse public key"})
	}

	if requestBody.TTLMinutes > ca.MaxTTLMinutes {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Requested TTL longer than configured max"})
	}

	if !isSubset(requestBody.Principals, ca.ValidPrincipals) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{"Requested principals not in valid principal list"})
	}

	signedCert, err := cert.SignUserKey(signer, parsedPublicKey, requestBody.Principals, requestBody.TTLMinutes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{"Failed to sign public key"})
	}
	c.Logger().Infof("Signed public key %s for %s", comment, CaID)
	response := cert.SignResponse{
		SignedKey: string(ssh.MarshalAuthorizedKey(signedCert)),
	}

	return c.JSON(http.StatusCreated, response)
}

// IsSubset checks if list1 is a subset of list2
func isSubset(list1, list2 []string) bool {
	// Create a map to store elements of list2
	list2Map := make(map[string]bool)

	// Populate the map with elements from list2
	for _, v := range list2 {
		list2Map[v] = true
	}

	// Check if each element of list1 is in list2Map
	for _, v := range list1 {
		if !list2Map[v] {
			return false // If any element is not in list2, return false
		}
	}

	return true // All elements of list1 are in list2
}
