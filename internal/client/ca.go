package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/lukegriffith/SSHTrust/pkg/cert"
	"github.com/lukegriffith/SSHTrust/pkg/handlers"
)

func CreateCA(body cert.CaRequest) error {
	jsonValue, _ := json.Marshal(body)
	req, err := MakeRequest(POST, "http://localhost:8080/CA", jsonValue, readToken)

	if err != nil {
		return fmt.Errorf("failed to create request %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("failed to create CA: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorMessage handlers.ErrorResponse
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error unmarshalling api error: %v", err)
		}

		err = json.Unmarshal(bodyBytes, &errorMessage)
		return fmt.Errorf("failed to create CA: %v - %s", resp.StatusCode, errorMessage)
	}
	defer resp.Body.Close()
	log.Println("CA created successfully")
	return nil
}

func GetCA(id string) (string, error) {

	req, err := MakeRequest(GET, fmt.Sprintf("http://localhost:8080/CA/%s", id), nil, readToken)

	if err != nil {
		return "", fmt.Errorf("failed to create request %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", fmt.Errorf("failed to get CA public key: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func SignPublicKey(id string, body cert.SignRequest) (*cert.SignResponse, error) {
	jsonValue, _ := json.Marshal(body)

	req, err := MakeRequest(POST, fmt.Sprintf("http://localhost:8080/CA/%s/Sign", id), jsonValue, readToken)

	if err != nil {
		return nil, fmt.Errorf("failed to create request %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to sign public key: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorMessage handlers.ErrorResponse
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling api error: %v", err)
		}

		err = json.Unmarshal(bodyBytes, &errorMessage)
		return nil, fmt.Errorf("failed to sign key: %v - %s", resp.StatusCode, errorMessage)
	}
	defer resp.Body.Close()

	signedKey, _ := io.ReadAll(resp.Body)

	var result cert.SignResponse
	err = json.Unmarshal(signedKey, &result)

	return &result, nil
}

func ListCAs() ([]cert.CaResponse, error) {

	req, err := MakeRequest(GET, "http://localhost:8080/CA", nil, readToken)

	if err != nil {
		return nil, fmt.Errorf("failed to create request %w", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve CA list: %w", err)
	}
	defer resp.Body.Close()

	var cas []cert.CaResponse
	err = json.NewDecoder(resp.Body).Decode(&cas)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA list: %w", err)
	}

	return cas, nil
}
