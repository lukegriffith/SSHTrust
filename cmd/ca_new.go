package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/lukegriffith/SSHTrust/internal/client"
	"github.com/lukegriffith/SSHTrust/pkg/cert"
	"github.com/spf13/cobra"
)

var caNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new Certificate Authority",
	Run: func(cmd *cobra.Command, args []string) {
		// Extract the flags
		name, _ := cmd.Flags().GetString("name")
		bits, _ := cmd.Flags().GetInt("bits")
		keyType, _ := cmd.Flags().GetString("type")
		principals, _ := cmd.Flags().GetString("validPrincipals")
		ttl, _ := cmd.Flags().GetInt("ttl")

		// Basic validation
		if name == "" {
			log.Fatal("CA name is required")
		}
		body := cert.CaRequest{
			Name:            name,
			Bits:            bits,
			Type:            cert.KeyType(keyType),
			ValidPrincipals: strings.Split(principals, ","),
			MaxTTLMinutes:   ttl,
		}
		log.Println(body)
		err := client.CreateCA(body)
		if err != nil {
			log.Fatalf("Failed to create CA: %v", err)
		}

		fmt.Printf("CA '%s' created successfully\n", name)
	},
}

func init() {
	// Add flags to the new CA command
	caNewCmd.Flags().StringP("name", "n", "", "Name of the CA (required)")
	caNewCmd.Flags().IntP("bits", "b", 2048, "Key size in bits (optional, default 2048)")
	caNewCmd.Flags().StringP("type", "t", "rsa", "Key type (optional, rsa, ed25519 [default rsa])")
	caNewCmd.Flags().StringP("validPrincipals", "p", "", "comma separated principals (required)")
	caNewCmd.Flags().Int("ttl", 60, "Maximim TTL the CA permits")

	_ = signCmd.MarkFlagRequired("name")
	_ = signCmd.MarkFlagRequired("principals")
	// Register the new CA command under the `ca` command
	caCmd.AddCommand(caNewCmd)
}
