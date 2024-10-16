package certStore

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/lukegriffith/SSHTrust/pkg/auth"
)

type InMemoryUserList map[string]string

func (ul InMemoryUserList) GetPasswordHash(un string) (string, error) {
	if password, ok := ul[un]; ok {
		return password, nil
	}
	return "", errors.New("Unable to find user")
}

func (ul InMemoryUserList) Register(u *auth.User) *echo.HTTPError {
	if _, ok := ul[u.Username]; ok {
		return echo.ErrBadRequest
	}
	hashedPass, err := auth.GenerateHash(u.Password)
	if err != nil {
		return echo.ErrInternalServerError
	}
	ul[u.Username] = hashedPass

	return nil
}
