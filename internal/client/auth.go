package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lukegriffith/SSHTrust/pkg/auth"
	"github.com/lukegriffith/SSHTrust/pkg/handlers"
)

var TokenLocation = ".sshtrust.token"

type TokenResponse struct {
	Token string `json:"token"` // Adjust this if the token key in the response JSON is different
}

func Login(body auth.User) error {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	jsonValue, _ := json.Marshal(body)
	resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorMessage handlers.ErrorResponse
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error unmarshalling api error: %v", err)
		}

		err = json.Unmarshal(bodyBytes, &errorMessage)
		return fmt.Errorf("failed to login: %v - %s", resp.StatusCode, errorMessage)
	}
	defer resp.Body.Close()
	// Read the response body
	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body %w", err)
	}
	var tokenResp TokenResponse
	err = json.Unmarshal(bodyResp, &tokenResp)
	if err != nil {
		return fmt.Errorf("Error parsing JSON: %w", err)
	}

	tokenFilePath := filepath.Join(homeDir, TokenLocation)
	// Write the token to a file
	err = writeTokenToFile(tokenResp.Token, tokenFilePath)
	if err != nil {
		return fmt.Errorf("Error writing token to file %w", err)
	}
	return nil
}

func writeTokenToFile(token, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(token)
	if err != nil {
		return err
	}

	return nil
}

func readToken() (string, error) {
	// Get the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %v", err)
	}

	// Build the full path to the token file
	tokenFilePath := filepath.Join(homeDir, TokenLocation)

	// Read the file contents
	tokenBytes, err := os.ReadFile(tokenFilePath)
	if err != nil {
		return "", fmt.Errorf("could not read token file: %v", err)
	}

	return string(tokenBytes), nil

}
