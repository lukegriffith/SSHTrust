package cert

import "testing"

func TestCaValidation(t *testing.T) {
	table := []struct {
		ca  CaRequest
		res bool
	}{
		// Test cases
		{CaRequest{Name: "TestCA", Type: "rsa", Bits: 2048}, true},          // Valid case
		{CaRequest{Name: "TestCA", Type: "rsa", Bits: 3072}, true},          // Valid case
		{CaRequest{Name: "TestCA", Type: "rsa", Bits: 4096}, true},          // Valid case
		{CaRequest{Name: "TestCA", Type: "ED25519", Bits: 2048}, false},     // Invalid case (wrong Type)
		{CaRequest{Name: "TestCA", Type: "rsa", Bits: 1024}, false},         // Invalid case (wrong Bits)
		{CaRequest{Name: "TestCA", Type: "rsa", Bits: 8192}, false},         // Invalid case (unsupported Bits)
		{CaRequest{Name: "", Type: "rsa", Bits: 2048}, false},               // Invalid case (missing Name)
		{CaRequest{Name: "TestCA", Type: "InvalidType", Bits: 2048}, false}, // Invalid case (wrong Type)
		{CaRequest{Name: "TestCA", Type: "rsa", Bits: 1234}, false},         // Invalid case (wrong Bits)
		{CaRequest{Name: "TestCA", Type: "", Bits: 2048}, false},            // Invalid case (empty Type)
		{CaRequest{Name: "TestCA", Type: "rsa", Bits: 0}, false},            // Invalid case (Bits = 0)
	}

	// Loop through test cases
	for _, tt := range table {
		t.Run(tt.ca.Name, func(t *testing.T) {
			result := tt.ca.Validate()
			if result != tt.res {
				t.Errorf("expected %v, got %v", tt.res, result)
			}
		})
	}
}
