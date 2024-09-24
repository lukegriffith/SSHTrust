package cert

import "testing"

func TestCaValidation(t *testing.T) {
	table := []struct {
		name string
		ca   CaRequest
		res  bool
	}{
		// Test cases with descriptive names
		{"Valid RSA 2048 bits", CaRequest{CommonCa{Name: "TestCA", Type: "ssh-rsa", Bits: 2048, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}}, true},
		{"Valid RSA 3072 bits", CaRequest{CommonCa{Name: "TestCA", Type: "ssh-rsa", Bits: 3072, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}}, true},
		{"Valid RSA 4096 bits", CaRequest{CommonCa{Name: "TestCA", Type: "ssh-rsa", Bits: 4096, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}}, true},
		{"Valid ED25519", CaRequest{CommonCa{Name: "TestCA", Type: "ssh-ed25519", Bits: 0, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}}, true},
		{"Invalid RSA 1024 bits", CaRequest{CommonCa{Name: "TestCA", Type: "ssh-rsa", Bits: 1024, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}}, false},
		{"Invalid RSA 8192 bits", CaRequest{CommonCa{Name: "TestCA", Type: "ssh-rsa", Bits: 8192, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}}, false},
		{"Invalid missing Name", CaRequest{CommonCa{Name: "", Type: "ssh-rsa", Bits: 2048, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}}, false},
		{"Invalid Type InvalidType", CaRequest{CommonCa{Name: "TestCA", Type: "InvalidType", Bits: 2048, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}}, false},
		{"Invalid RSA 1234 bits", CaRequest{CommonCa{Name: "TestCA", Type: "ssh-rsa", Bits: 1234, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}}, false},
		{"Invalid empty Type", CaRequest{CommonCa{Name: "TestCA", Type: "", Bits: 2048, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}}, false},
		{"Invalid RSA 0 bits", CaRequest{CommonCa{Name: "TestCA", Type: "ssh-rsa", Bits: 0, ValidPrincipals: []string{"testuser"}, MaxTTLMinutes: 3600}}, false},
	}

	// Loop through test cases
	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			err, result := tt.ca.Validate()
			if result != tt.res {
				t.Errorf("expected %v, got %v, %s", tt.res, result, err)
			}
		})
	}
}
