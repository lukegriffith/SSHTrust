package cert

import "golang.org/x/crypto/ssh"

type CA struct {
	Name   string
	Signer ssh.Signer
	Bits   int
}

func (c CA) CreateResponse() *CaResponse {
	return &CaResponse{
		Name:      c.Name,
		PublicKey: string(ssh.MarshalAuthorizedKey(c.Signer.PublicKey())),
		Type:      c.Signer.PublicKey().Type(),
		Bits:      c.Bits,
	}
}

type CaRequest struct {
	// Name of CA
	Name string `json:"name"`
	// Type of ca, rsa, ed25519
	Type string `json:"type"`
	// Key length
	Bits int `json:"bits"`
}

func (c CaRequest) Validate() bool {
	if c.Name == "" {
		return false
	}
	if c.Type != "rsa" {
		return false
	}
	if !(c.Bits == 2048 || c.Bits == 3072 || c.Bits == 4096) {
		return false
	}
	return true
}

type CaResponse struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Bits      int    `json:"bits"`
	PublicKey string `json:"public_key"`
}
