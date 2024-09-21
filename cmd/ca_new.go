package cmd

import (
	"fmt"
	"log"

	"github.com/lukegriffith/SSHTrust/internal/client"
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

		// Basic validation
		if name == "" {
			log.Fatal("CA name is required")
		}

		err := client.CreateCA(name, bits, keyType)
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
	caNewCmd.Flags().StringP("type", "t", "rsa", "Key type (optional, default rsa)")

	// Register the new CA command under the `ca` command
	caCmd.AddCommand(caNewCmd)
}

