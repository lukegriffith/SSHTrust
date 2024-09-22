package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/lukegriffith/SSHTrust/pkg/cert"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testPrivateKey = `
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn
NhAAAAAwEAAQAAAQEA4J9aYG5OhRqpgGkYiJsjQWct3uVsvKdLXOIucqOTobJEX/Px0lrS
9I1vubmLWY0+VRj7QlNJnDhwsXp/jynReG6gTuEIjpHU4TRoR0x/MB/Wif0V77V7LGbWV8
i92MrRmqiM+5KJCgmP+TWN07nrbaIszTd1vMBbX5wSk6lWdDDWSGE+WRhAJeJtPh5Kvaev
w2QMADtxF7+2p6ozalskwnESEpcTEME0bLsO6Tg0FLFJT0YEMs49oYAuKurDlwANjDcUJl
+LBQUKeODP7UY+jU5MonD+5bjQdZu5lGLJIfUMAV3sWRSXZQ+xLuUwpPb0IgLprKmXosm9
0rUsNYO0xwAAA9DFlbF+xZWxfgAAAAdzc2gtcnNhAAABAQDgn1pgbk6FGqmAaRiImyNBZy
3e5Wy8p0tc4i5yo5OhskRf8/HSWtL0jW+5uYtZjT5VGPtCU0mcOHCxen+PKdF4bqBO4QiO
kdThNGhHTH8wH9aJ/RXvtXssZtZXyL3YytGaqIz7kokKCY/5NY3TuettoizNN3W8wFtfnB
KTqVZ0MNZIYT5ZGEAl4m0+Hkq9p6/DZAwAO3EXv7anqjNqWyTCcRISlxMQwTRsuw7pODQU
sUlPRgQyzj2hgC4q6sOXAA2MNxQmX4sFBQp44M/tRj6NTkyicP7luNB1m7mUYskh9QwBXe
xZFJdlD7Eu5TCk9vQiAumsqZeiyb3StSw1g7THAAAAAwEAAQAAAQEAtWm+ElfEbtfjweQf
fmTdinsMnxLoSU1MHo5GOSxHlbZmZqCc0+mqqvx4GaXzF7Zte0kb1KtzrgofahenYjbCCO
Q/8LTqtkqthd1PwxXTO0jbesK+rsUB4BCGWIu2WJslwiUTCDOHHmYus3U/QJrNu9PZHz0F
iBZLNeLVVhDT1C4BS3E6OOLJNIS6FN/BuqNDM0QQW6QLvF7UrSIOMgLPd1fpZtvae1CNUr
K2iyMb7AhU5SIs//F4yFZHdBTTw26jI5rQcEAW00DlifeZBFGWX5qCtZfR7sD+X0xY1r4/
QV0PajAR7lF5w4at/8/qaZQRkvd9RdEv0E4JUvLjdlx1sQAAAIBA8vTVwE2lmp0ztpmCYp
iMCCnxvZSd48+Gxml56LrTD/paLX3eXJ5xWAkQC0yQSbbV/+XaMamKwo1r2rHvj7nKSsCj
qh1+PvCEibmsoLCZskQQfLbQ1fFzTK+fdVC4GC1h/Vp6jL9haxZxU+gRMPkq43++d6CSgl
JebyqP5QB4/gAAAIEA/cGEAKdXNs8hjucg0X0ahY/RzqwfO32k+vuBnMhr2/9LPOQHZOFv
9FMUXHAaX957J818XiKPwwSVktVcVp/8e2a47cfANtrMe7tNzS7OKBt9UxgN05QKJF4k03
RZmTbNJjLZZkHW6cOrB98BWvtS/gQmwBRqVHW0hGZ/9di+zO0AAACBAOKb4adQoXZ+fbR0
/MO7yFiSbgPy5xjca7skAWBu/bMjILjCKGdBcJuQQOaLXxauHl5/sW+mwYvyu0RyvShLPu
6oyWWqY78iQuJ6Ecq/HUDf/T9Nh/BtF2/MJ80HF978NUWahG0uN7CurBwNeXoqsAwZRRER
I9TE6deN/Eiq6MYDAAAAF2x1a2VncmlmZml0aEBsZ21iLmxvY2FsAQID
-----END OPENSSH PRIVATE KEY-----`

var testPublicKey = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIAJI+V4/0d5xJTDvOuvR/2ZqahzceFbz00IDIFBEaKvc test@testserver.com"

func createMockSigner() (ssh.Signer, error) {
	signer, err := ssh.ParsePrivateKey([]byte(testPrivateKey))
	if err != nil {
		return nil, err
	}
	return signer, nil
}

// Test for the Sign handler using the existing MockStore
func TestSignHandler(t *testing.T) {
	e := echo.New()

	signer, err := createMockSigner()
	if err != nil {
		t.Fatalf("Failed to create mock signer: %v", err)
	}

	mockStore := &MockStore{
		caMap: map[string]*cert.CaResponse{
			"test-ca": {
				Name:            "test-ca",
				MaxTTLMinutes:   60,
				ValidPrincipals: []string{"user1", "user2"},
			},
		},
		signers: map[string]ssh.Signer{
			"test-ca": signer,
		},
	}
	app := &App{Store: mockStore}

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
		caID           string
	}{
		{
			name:           "Successful Signing",
			requestBody:    `{"public_key":"` + testPublicKey + `","principals":["user1"],"ttl_minutes":30}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   "signed_key",
			caID:           "test-ca",
		},
		{
			name:           "CA Not Found",
			requestBody:    `{"public_key":"` + testPublicKey + `","principals":["user1"],"ttl_minutes":30}`,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "CA not found",
			caID:           "nonexistent-ca",
		},
		{
			name:           "Invalid Public Key",
			requestBody:    `{"public_key":"invalid-public-key","principals":["user1"],"ttl_minutes":30}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Failed to parse public key",
			caID:           "test-ca",
		},
		{
			name:           "TTL Exceeds Maximum",
			requestBody:    `{"public_key":"` + testPublicKey + `","principals":["user1"],"ttl_minutes":100}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Requested TTL longer than configured max",
			caID:           "test-ca",
		},
		{
			name:           "Invalid Principals",
			requestBody:    `{"public_key":"` + testPublicKey + `","principals":["invalid"],"ttl_minutes":30}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Requested principals not in valid principal list",
			caID:           "test-ca",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/ca/"+tt.caID+"/sign", strings.NewReader(tt.requestBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.caID)

			err := app.Sign(c)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}
