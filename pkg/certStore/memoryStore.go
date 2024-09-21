package certStore

import (
	"sync"

	"errors"

	"github.com/lukegriffith/SSHTrust/pkg/cert"
	"golang.org/x/crypto/ssh"
)

type InMemortCaStore struct {
	sync.RWMutex
	cas map[string]cert.CA
}

func NewInMemoryCaStore() *InMemortCaStore {
	return &InMemortCaStore{
		cas: make(map[string]cert.CA),
	}

}

func (store *InMemortCaStore) GetCAByID(ID string) (*cert.CaResponse, error) {
	store.RLock()
	defer store.RUnlock()
	if value, exists := store.cas[ID]; exists {
		return value.CreateResponse(), nil
	}
	return nil, errors.New("Unable to find CA by ID")
}

func (store *InMemortCaStore) GetSignerByID(ID string) (ssh.Signer, error) {
	store.RLock()
	defer store.RUnlock()
	if value, exists := store.cas[ID]; exists {
		return value.Signer, nil
	}
	return nil, errors.New("Unable to find CA by ID")
}

func (store *InMemortCaStore) CreateCA(CAReq cert.CaRequest) (*cert.CaResponse, error) {

	if !CAReq.Validate() {
		return nil, errors.New("Invalid CA Request")
	}
	store.RLock()
	if _, exists := store.cas[CAReq.Name]; exists {
		return nil, errors.New("CA Already Exists")
	}
	store.RUnlock()
	signer, err := cert.GenerateSSHKey(CAReq.Bits)
	if err != nil {
		return nil, errors.New("Failed to generate CA keypair")
	}

	c := cert.CA{
		Name:   CAReq.Name,
		Signer: signer,
		Bits:   CAReq.Bits,
	}

	store.Lock()
	store.cas[CAReq.Name] = c
	store.Unlock()
	return c.CreateResponse(), nil
}

func (store *InMemortCaStore) ListCAs() ([]*cert.CaResponse, error) {
	keys := []*cert.CaResponse{}
	store.RLock()
	defer store.RUnlock()
	for _, value := range store.cas {
		ca := value.CreateResponse()
		keys = append(keys, ca)
	}
	return keys, nil
}
