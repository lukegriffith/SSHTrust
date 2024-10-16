package cmd

import (
	"github.com/lukegriffith/SSHTrust/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		noAuth, _ := cmd.Flags().GetBool("no-auth")
		e := server.SetupServer(noAuth)
		e.Logger.Printf("SSHTrust Started on %s", server.Port)
		if noAuth {
			e.Logger.Printf("No auth enabled %t", noAuth)
		}
		if err := e.Start(server.Port); err != nil {
			e.Logger.Fatal(err)
		}
		// Add server starting logic here
	},
}

func init() {
	serveCmd.Flags().Bool("no-auth", false, "Enable user auth")
	rootCmd.AddCommand(serveCmd)

}
