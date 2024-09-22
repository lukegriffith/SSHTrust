package cmd

import (
	"fmt"
	"github.com/lukegriffith/SSHTrust/internal/client" // Update with the correct import path
	"log"

	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign [ca_id] [public_key]",
	Short: "Sign a public key using a specific Certificate Authority",
	Args:  cobra.ExactArgs(2), // Ensure exactly 2 arguments: the CA ID and public key
	Run: func(cmd *cobra.Command, args []string) {
		// Extract the CA ID and public key from the command arguments
		caID := args[0]
		publicKey := args[1]

		// Call the client library to sign the public key
		signedKey, err := client.SignPublicKey(caID, publicKey)
		if err != nil {
			log.Fatalf("Error signing public key: %v", err)
		}

		// Display the signed certificate
		fmt.Print(signedKey.SignedKey)
	},
}

func init() {
	// Register the sign command under the root command
	rootCmd.AddCommand(signCmd)
}
