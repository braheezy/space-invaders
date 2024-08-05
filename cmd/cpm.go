package cmd

import (
	"github.com/braheezy/space-invaders/internal/cpm"
	"github.com/braheezy/space-invaders/internal/emulator"
	"github.com/charmbracelet/log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(cpmCmd)
}

var cpmCmd = &cobra.Command{
	Use:   "cpm",
	Short: "Run CP/M test",
	Run: func(cmd *cobra.Command, args []string) {
		logger := newDefaultLogger()
		if debug {
			logger.SetLevel(log.DebugLevel)
		}
		cpmHardware := cpm.NewCPMHardware()

		// Assuming NewCPMHardware() sets up the CP/M environment
		vm := emulator.NewEmulator(cpmHardware)
		vm.StartInterruptRoutines()
		vm.Logger = logger
		vm.Options.UnlimitedTPS = true

		ebiten.SetWindowTitle("cpm test")
		if vm.Options.UnlimitedTPS {
			ebiten.SetTPS(ebiten.SyncWithFPS)
		} else {
			ebiten.SetTPS(60)
		}
		ebiten.SetWindowSize(vm.Hardware.Width()*vm.Hardware.Scale(), vm.Hardware.Height()*vm.Hardware.Scale())

		if err := ebiten.RunGame(vm); err != nil && err != ebiten.Termination {
			logger.Fatal(err)
		}
	},
}
