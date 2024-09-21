package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lukegriffith/SSHTrust/pkg/cert"
	"io"
	"log"
	"net/http"
)

func CreateCA(name string, bits int, keyType string) error {
	body := cert.CaRequest{
		Name: name,
		Bits: bits,
		Type: keyType,
	}
	jsonValue, _ := json.Marshal(body)

	resp, err := http.Post("http://localhost:8080/CA", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return fmt.Errorf("failed to create CA: %w", err)
	}
	defer resp.Body.Close()

	log.Println("CA created successfully")
	return nil
}

func GetCAPublicKey(id string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/CA/%s", id))
	if err != nil {
		return "", fmt.Errorf("failed to get CA public key: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func SignPublicKey(id, publicKey string) (string, error) {
	body := map[string]string{
		"public_key": publicKey,
	}
	jsonValue, _ := json.Marshal(body)

	resp, err := http.Post(fmt.Sprintf("http://localhost:8080/CA/%s/Sign", id), "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", fmt.Errorf("failed to sign public key: %w", err)
	}
	defer resp.Body.Close()

	signedKey, _ := io.ReadAll(resp.Body)
	return string(signedKey), nil
}

func ListCAs() ([]map[string]interface{}, error) {
	resp, err := http.Get("http://localhost:8080/CA")
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve CA list: %w", err)
	}
	defer resp.Body.Close()

	var cas []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&cas)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA list: %w", err)
	}

	return cas, nil
}
