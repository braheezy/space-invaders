package cmd

import (
	"os"
	"path/filepath"

	"github.com/braheezy/space-invaders/internal/emulator"
	"github.com/charmbracelet/log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

var debug bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Show debug messages")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:  "8080 <rom>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logger := newDefaultLogger()
		if debug {
			logger.SetLevel(log.DebugLevel)
		}

		fileName := filepath.Base(args[0])
		data, err := os.ReadFile(args[0])
		if err != nil {
			logger.Fatal(err)
		}

		vm := emulator.NewCPU8080(&data, emulator.NewSpaceInvadersHardware())
		vm.StartInterruptRoutines()
		vm.Logger = logger

		ebiten.SetWindowTitle(fileName)
		ebiten.SetTPS(60)
		ebiten.SetWindowSize(vm.Hardware.Width()*vm.Hardware.Scale(), vm.Hardware.Height()*vm.Hardware.Scale())

		if err := ebiten.RunGame(vm); err != nil && err != ebiten.Termination {
			logger.Fatal(err)
		}
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	SilenceUsage: true,
}
