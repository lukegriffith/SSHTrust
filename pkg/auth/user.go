package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// TODO: Load from environment or config or random
var JWTSecret = []byte("secret")

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TODO: do this properly, externally, hashed + salted, etc
var Users = map[string]string{
	"admin": "password123", // Sample in-memory user for testing
}

// Login handler
func Login(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	// Validate user credentials (check if username and password match)
	if password, ok := Users[u.Username]; !ok || password != u.Password {
		return echo.ErrUnauthorized
	}

	claims := jwt.MapClaims{
		"authorized": true,
		"user":       u.Username,
		"exp":        time.Now().Add(time.Hour * 72).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(JWTSecret)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
