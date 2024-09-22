package cmd

import (
	"github.com/lukegriffith/SSHTrust/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		e := server.SetupServer()
		e.Logger.Printf("SSHTrust Started on %s", server.Port)
		if err := e.Start(server.Port); err != nil {
			e.Logger.Fatal(err)
		}
		// Add server starting logic here
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
