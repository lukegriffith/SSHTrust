package client

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mocking the readToken function
var mockReadToken = func() (string, error) {
	return "mocked_token", nil
}

func TestMakeRequestSuccess(t *testing.T) {
	// Setup
	method := GET
	url := "http://example.com/api"
	body := []byte(`{"key":"value"}`)

	// Execute the function being tested
	req, err := MakeRequest(method, url, body, mockReadToken)

	// Assertions
	assert.NoError(t, err)                                                  // Ensure no error occurred
	assert.NotNil(t, req)                                                   // Ensure the request is not nil
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))     // Check the Content-Type header
	assert.Equal(t, "Bearer mocked_token", req.Header.Get("Authorization")) // Check the Authorization header
	assert.Equal(t, string(method), req.Method)                             // Ensure the method is correct
	assert.Equal(t, url, req.URL.String())                                  // Ensure the URL is correct
}

func TestMakeRequestErrorCreatingRequest(t *testing.T) {
	// Setup invalid URL to cause error in http.NewRequest
	method := GET
	url := string([]byte{0x7f}) // Invalid URL will cause error
	body := []byte(`{"key":"value"}`)

	// Execute the function being tested
	req, err := MakeRequest(method, url, body, mockReadToken)

	// Assertions
	assert.Nil(t, req)                                        // Ensure the request is nil
	assert.Error(t, err)                                      // Ensure an error occurred
	assert.Contains(t, err.Error(), "Error creating request") // Check that the error is related to creating the request
}

func TestMakeRequestErrorReadingToken(t *testing.T) {
	// Mock readToken to return an error
	mockReadToken = func() (string, error) {
		return "", errors.New("token read error")
	}

	// Reset the mock function after the test
	defer func() {
		mockReadToken = func() (string, error) {
			return "mocked_token", nil
		}
	}()

	// Setup
	method := GET
	url := "http://example.com/api"
	body := []byte(`{"key":"value"}`)

	// Execute the function being tested
	req, err := MakeRequest(method, url, body, mockReadToken)

	// Assertions
	assert.Nil(t, req)                                       // Ensure the request is nil
	assert.Error(t, err)                                     // Ensure an error occurred
	assert.Contains(t, err.Error(), "Failed making request") // Check that the error is related to the token
}
