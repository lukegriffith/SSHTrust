package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/lukegriffith/SSHTrust/internal/client" // Update with the correct import path
	"github.com/lukegriffith/SSHTrust/pkg/cert"

	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a public key using a specific Certificate Authority",
	Run: func(cmd *cobra.Command, args []string) {
		// Extract the CA ID and public key from the command arguments
		caID, _ := cmd.Flags().GetString("name")
		publicKey, _ := cmd.Flags().GetString("public_key")
		principals, _ := cmd.Flags().GetString("principals")
		ttl, _ := cmd.Flags().GetInt("ttl")

		body := cert.SignRequest{
			PublicKey:  publicKey,
			Principals: strings.Split(principals, ","),
			TTLMinutes: ttl,
		}
		// Call the client library to sign the public key
		signedKey, err := client.SignPublicKey(caID, body)
		if err != nil {
			log.Fatalf("Error signing public key: %v", err)
		}

		// Display the signed certificate
		fmt.Print(signedKey.SignedKey)
	},
}

func init() {
	// Add flags
	signCmd.Flags().StringP("name", "n", "", "Name of the CA")
	signCmd.Flags().StringP("public_key", "k", "", "Public key to be signed")
	signCmd.Flags().StringP("principals", "p", "", "Comma-separated list of principals for the certificate")
	signCmd.Flags().Int("ttl", 60, "Time to live for the certificate in seconds")

	// Optionally, mark flags as required
	_ = signCmd.MarkFlagRequired("name")
	_ = signCmd.MarkFlagRequired("public_key")
	_ = signCmd.MarkFlagRequired("principals")
	// Register the sign command under the root command
	rootCmd.AddCommand(signCmd)
}
