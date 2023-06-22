package cmd

import (
	"github.com/spf13/cobra"
	"ticket-api/cmd/bootstrap"
)

var webCmd = &cobra.Command{
	Use:     "app:serve",
	Short:   "An application to start http server for Ticket",
	Long:    "An application to expose REST endpoints to users by starting http server",
	Example: "go run main.go app:serve",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.WebApplication()
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}
