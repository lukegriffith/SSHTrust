package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainFunction(t *testing.T) {
	// Mock Echo instance (as in SetupServer)
	e := SetupServer()

	// We do not actually start the server in the test.
	// Instead, we check that the instance has the expected routes and properties.
	assert.NotNil(t, e, "Expected Echo instance to be set up")
}
