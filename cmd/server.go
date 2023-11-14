package cmd

import (
	"MedKick-backend/cmd/server"

	"github.com/spf13/cobra"
)

var runServerCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
}

func init() {
	rootCmd.AddCommand(runServerCmd)
}
