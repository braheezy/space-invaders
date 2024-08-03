package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var debug bool
var startAddress int

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Show debug messages")
	rootCmd.PersistentFlags().IntVarP(&startAddress, "start", "s", 0, "Set program start address (in decimal)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:          "8080",
	Short:        "8080 CPU emulator",
	SilenceUsage: true,
}
