package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// TODO: Load from environment or config or random

var (
	JWTSecret []byte
	Users     UserList
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserList interface {
	GetPasswordHash(un string) (string, error)
	Register(u *User) *echo.HTTPError
}

func GenerateHash(s string) (string, error) {
	// Generate a hashed password using bcrypt with a default cost of bcrypt.DefaultCost
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func Register(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	httpErr := Users.Register(u)

	if httpErr != nil {
		return httpErr
	}
	return c.JSON(http.StatusOK, nil)
}

// Login handler
func Login(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	storedPassHash, err := Users.GetPasswordHash(u.Username)
	if err != nil {
		return echo.ErrNotFound
	}

	// Validate user credentials (check if username and password match)
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassHash), []byte(u.Password)); err != nil {
		c.Logger().Warn(err)
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
