package cmd

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"github.com/lukegriffith/SSHTrust/internal/client"
	"github.com/lukegriffith/SSHTrust/pkg/auth"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to server",
	Run: func(cmd *cobra.Command, args []string) {
		// Extract the flags
		userName, _ := cmd.Flags().GetString("username")
		stdin, _ := cmd.Flags().GetBool("stdin")
		var password string
		var err error

		if stdin {
			// Read password from stdin
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter password from stdin: ")
			password, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading password from stdin:", err)
				return
			}
			password = password[:len(password)-1] // Remove the trailing newline

		} else {
			// Prompt for password securely
			fmt.Print("Enter password: ")
			bytePassword, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Println("Error reading password:", err)
				return
			}
			password = string(bytePassword)
			fmt.Println() // Move to the next line after password input
		}

		err = client.Login(auth.User{
			Username: userName,
			Password: password,
		})
		if err != nil {
			fmt.Println("Login Fail %w", err)
			return
		}
		fmt.Println("Login Success")

	},
}

func init() {
	// Add flags to the new CA command
	loginCmd.Flags().StringP("username", "u", "", "Login username (required)")
	loginCmd.Flags().BoolP("stdin", "i", false, "Read password from stdin")
	_ = loginCmd.MarkFlagRequired("username")
	// Register the new CA command under the `ca` command
	rootCmd.AddCommand(loginCmd)
}
