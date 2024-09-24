package cert

import (
	"errors"
	"golang.org/x/crypto/ssh"
)

type CA struct {
	Name            string
	Signer          ssh.Signer
	Bits            int
	MaxTTLMinutes   int
	ValidPrincipals []string
}

func NewCA(name string, signer ssh.Signer, validPrincipals []string, bits, maxTtl int) CA {
	return CA{
		Name:            name,
		Signer:          signer,
		ValidPrincipals: validPrincipals,
		Bits:            bits,
		MaxTTLMinutes:   maxTtl,
	}
}

func (c CA) CreateResponse() *CaResponse {
	return &CaResponse{
		CommonCa: CommonCa{
			Name:            c.Name,
			Type:            KeyType(c.Signer.PublicKey().Type()),
			Bits:            c.Bits,
			MaxTTLMinutes:   c.MaxTTLMinutes,
			ValidPrincipals: c.ValidPrincipals,
		},
		PublicKey: string(ssh.MarshalAuthorizedKey(c.Signer.PublicKey())),
	}
}

type CommonCa struct {
	// Name of CA
	Name string `json:"name"`
	// Type of ca, rsa, ed25519
	Type KeyType `json:"type"`
	// Key length
	Bits int `json:"bits"`
	// Maximum TTL certs can be signed for
	MaxTTLMinutes int `json:"max_ttl_minutes"`
	// List of Valid Principals
	ValidPrincipals []string `json:"valid_principals"`
}

type CaRequest struct {
	CommonCa
}

func (c CaRequest) Validate() (error, bool) {
	if c.Name == "" {
		return errors.New("invalid name"), false
	}
	if c.Type != RSAKey && c.Type != ED25519 {
		return InvalidKeyErr, false
	}
	if !(c.Bits == 2048 || c.Bits == 3072 || c.Bits == 4096) &&
		c.Type != ED25519 {
		return errors.New("invalid key length"), false
	}
	if len(c.ValidPrincipals) < 1 {
		return errors.New("no principals provided"), false
	}
	if c.MaxTTLMinutes == 0 {
		return errors.New("MaxTTL not set"), false
	}
	return nil, true
}

type CaResponse struct {
	CommonCa
	// CA Public Key
	PublicKey string `json:"public_key"`
}
