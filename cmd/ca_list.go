package cmd

import (
	"fmt"
	"github.com/lukegriffith/SSHTrust/internal/client" // Update with the correct import path
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
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

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Type", "Length"})

		// Display the list of CAs
		if len(cas) == 0 {
			fmt.Println("No Certificate Authorities found.")
		} else {
			for _, ca := range cas {
				table.Append([]string{ca.Name, ca.Type, strconv.Itoa(ca.Bits)})
			}
		}
		table.Render()
	},
}

func init() {
	// Register the list command under the ca command
	caCmd.AddCommand(caListCmd)
}
