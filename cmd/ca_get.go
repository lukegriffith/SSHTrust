package cmd

import (
	"fmt"
	"github.com/lukegriffith/SSHTrust/internal/client" // Update with the correct import path
	"log"

	"github.com/spf13/cobra"
)

var caGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Retrieve the public key of a specific Certificate Authority",
	Args:  cobra.ExactArgs(1), // Ensure exactly 1 argument is passed (the CA ID)
	Run: func(cmd *cobra.Command, args []string) {
		// Extract the CA ID from the command arguments
		id := args[0]

		// Call the client library to retrieve the CA's public key
		publicKey, err := client.GetCAPublicKey(id)
		if err != nil {
			log.Fatalf("Error retrieving CA public key: %v", err)
		}

		// Display the public key
		fmt.Println(publicKey)
	},
}

func init() {
	// Register the get command under the ca command
	caCmd.AddCommand(caGetCmd)
}
