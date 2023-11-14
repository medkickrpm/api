package cmd

import (
	"MedKick-backend/cmd/migration"

	"github.com/spf13/cobra"
)

var runBootstrapCmd = &cobra.Command{
	Use: "migration",
	Run: func(cmd *cobra.Command, args []string) {
		migration.Run()
	},
}

func init() {
	rootCmd.AddCommand(runBootstrapCmd)
}
