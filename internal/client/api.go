package client

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

type MethodType string

var (
	POST MethodType = "POST"
	GET  MethodType = "GET"
)

// Helper function for making HTTP GET requests
func MakeRequest(method MethodType, url string, body []byte, tokenFunc func() (string, error)) (*http.Request, error) {
	req, err := http.NewRequest(string(method), url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("Error creating request %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	token, err := tokenFunc()
	if err != nil {
		log.Println("Unable to find token")
		return req, nil
	}
	bearer := fmt.Sprintf("Bearer %s", token)

	req.Header.Set("Authorization", bearer)

	return req, nil
}
