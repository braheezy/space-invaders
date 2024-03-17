package cmd

import (
	"os"

	"github.com/braheezy/8080/internal/emulator"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "8080 <rom>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		emulator.Run(args[0])
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
