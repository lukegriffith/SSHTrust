package cmd

import (
	"fmt"
	"github.com/lukegriffith/SSHTrust/internal/client" // Update with the correct import path
	"log"

	"github.com/spf13/cobra"
)

var caListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Certificate Authorities",
	Run: func(cmd *cobra.Command, args []string) {
		// Call the client library to fetch the list of CAs
		cas, err := client.ListCAs()
		if err != nil {
			log.Fatalf("Error retrieving CA list: %v", err)
		}

		// Display the list of CAs
		if len(cas) == 0 {
			fmt.Println("No Certificate Authorities found.")
		} else {
			fmt.Println("Certificate Authorities:")
			for _, ca := range cas {
				fmt.Printf("Name: %s, Type: %s, Bits: %v\n", ca["name"], ca["type"], ca["bits"])
			}
		}
	},
}

func init() {
	// Register the list command under the ca command
	caCmd.AddCommand(caListCmd)
}
