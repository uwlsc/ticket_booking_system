package cmd

import (
	"github.com/spf13/cobra"

	"ticket-api/cmd/bootstrap"
)

var seedCmd = &cobra.Command{
	Use:     "app:seed",
	Short:   "An application to seed users in Ticket application",
	Long:    "An application to seed users in Ticket application",
	Example: "go run main.go app:seed",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.SeedApplication()
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
