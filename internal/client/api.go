package client

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// Helper function for making HTTP GET requests
func MakeGetRequest(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("request failed with status " + resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
