package cmd

import (
	"MedKick-backend/cmd/worker"

	"github.com/spf13/cobra"
)

var runWorkerCmd = &cobra.Command{
	Use: "worker",
	Run: func(cmd *cobra.Command, args []string) {
		worker.Run()
	},
}

func init() {
	rootCmd.AddCommand(runWorkerCmd)
}
