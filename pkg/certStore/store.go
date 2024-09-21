package certStore

import (
	"github.com/lukegriffith/SSHTrust/pkg/cert"
	"golang.org/x/crypto/ssh"
)

type CAStore interface {
	GetCAByID(ID string) (*cert.CaResponse, error)
	GetSignerByID(ID string) (ssh.Signer, error)
	CreateCA(Req cert.CaRequest) (*cert.CaResponse, error)
	ListCAs() ([]*cert.CaResponse, error)
}
