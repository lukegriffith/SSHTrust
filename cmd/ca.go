package cmd

import (
	"github.com/spf13/cobra"
)

var caCmd = &cobra.Command{
	Use:   "ca",
	Short: "Manage Certificate Authorities",
}

func init() {
	rootCmd.AddCommand(caCmd)
}
