package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lukegriffith/SSHTrust/pkg/cert"
	"github.com/lukegriffith/SSHTrust/pkg/handlers"
)

func CreateCA(body cert.CaRequest) error {
	jsonValue, _ := json.Marshal(body)
	resp, err := http.Post("http://localhost:8080/CA", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("failed to create CA: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorMessage handlers.ErrorResponse
		bodyBytes, err := ioutil.ReadAll(resp.Body)
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
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/CA/%s", id))
	if err != nil {
		return "", fmt.Errorf("failed to get CA public key: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func SignPublicKey(id string, body cert.SignRequest) (*cert.SignResponse, error) {
	jsonValue, _ := json.Marshal(body)

	resp, err := http.Post(fmt.Sprintf("http://localhost:8080/CA/%s/Sign", id), "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, fmt.Errorf("failed to sign public key: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errorMessage handlers.ErrorResponse
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling api error: %v", err)
		}

		err = json.Unmarshal(bodyBytes, &errorMessage)
		return nil, fmt.Errorf("failed to create CA: %v - %s", resp.StatusCode, errorMessage)
	}
	defer resp.Body.Close()

	signedKey, _ := io.ReadAll(resp.Body)

	var result cert.SignResponse
	err = json.Unmarshal(signedKey, &result)

	return &result, nil
}

func ListCAs() ([]cert.CaResponse, error) {
	resp, err := http.Get("http://localhost:8080/CA")
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
