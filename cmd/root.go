package cmd

import (
	"fmt"
	"os"

	"github.com/braheezy/space-invaders/internal/emulator"
	"github.com/braheezy/space-invaders/internal/invaders"
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
	Use:   "invaders",
	Short: "Run Space Invaders",
	Run: func(cmd *cobra.Command, args []string) {
		logger := newDefaultLogger()
		if debug {
			logger.SetLevel(log.DebugLevel)
		}

		invadersHardware := invaders.NewSpaceInvadersHardware()

		vm := emulator.NewEmulator(invadersHardware)
		vm.StartInterruptRoutines()
		vm.Logger = logger

		game := NewSpaceInvadersGame(vm)
		vm.Options.LimitTPS = game.menuScreen.GetLimitTPS()

		ebiten.SetWindowTitle("space invaders")
		if vm.Options.LimitTPS {
			ebiten.SetTPS(60)
		} else {
			ebiten.SetTPS(ebiten.SyncWithFPS)
		}
		ebiten.SetWindowSize(vm.Hardware.Width()*vm.Hardware.Scale(), vm.Hardware.Height()*vm.Hardware.Scale())

		if err := ebiten.RunGame(game); err != nil && err != ebiten.Termination {
			game.cpuEmulator.Hardware.Cleanup()
			logger.Fatal(err)
		}
		game.cpuEmulator.Hardware.Cleanup()
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	SilenceUsage: true,
}

// SpaceInvadersGame encapsulates the CPU emulator and the settings menu
type SpaceInvadersGame struct {
	cpuEmulator    *emulator.CPU8080
	inSettingsMenu bool
	menuScreen     *MenuScreen
	tabPressed     bool
}

// NewSpaceInvadersGame creates a new SpaceInvadersGame instance
func NewSpaceInvadersGame(cpuEmulator *emulator.CPU8080) *SpaceInvadersGame {
	return &SpaceInvadersGame{
		cpuEmulator:    cpuEmulator,
		inSettingsMenu: false,
		menuScreen:     NewMenuScreen("settings.json"),
	}
}

// Update fulfills the Game interface for ebiten.
// This runs the emulator for one frame.
func (game *SpaceInvadersGame) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Handle Tab key press to toggle the settings menu
	if ebiten.IsKeyPressed(ebiten.KeyTab) {
		if !game.tabPressed {
			// Toggle only when Tab is first pressed
			game.toggleSettingsMenu()
			game.tabPressed = true
		}
	} else {
		// Reset the tabPressed state when the Tab key is released
		game.tabPressed = false
	}

	if game.inSettingsMenu {
		// Update menu logic
		game.menuScreen.Update()
	} else {
		// Run the CPU emulator
		game.cpuEmulator.Update()
	}

	return nil
}

// Draw fulfills the Game interface for ebiten
func (game *SpaceInvadersGame) Draw(screen *ebiten.Image) {
	if game.inSettingsMenu {
		// Draw the settings menu
		game.menuScreen.Draw(screen)
	} else {
		// Draw the CPU emulator output
		game.cpuEmulator.Draw(screen)
	}
}

// Layout fulfills the Game interface for ebiten
func (game *SpaceInvadersGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

// toggleSettingsMenu toggles the display of the settings menu.
func (game *SpaceInvadersGame) toggleSettingsMenu() {
	game.inSettingsMenu = !game.inSettingsMenu

	if game.inSettingsMenu {
		// Initialize menu screen with a specified settings file path
		game.menuScreen = NewMenuScreen("settings.json")
	} else {
		// Save settings after a change
		if err := game.menuScreen.saveSettings(); err != nil {
			game.menuScreen.errorMessage = fmt.Sprintf("Error saving settings: %v", err)
		}

		// Update game settings from the menu screen
		hardware := game.cpuEmulator.Hardware.(*invaders.SpaceInvadersHardware)
		hardware.ShipsSetting = game.menuScreen.GetShipsSetting()
		hardware.ExtraShipAt1000 = game.menuScreen.GetExtraShipAt1000()
		// Coin info displayed in demo screen 0=ON
		hardware.ShowCoinInfoOnDemo = !game.menuScreen.GetShowCoinInfoOnDemo()
		// Close the menu screen
		game.menuScreen = nil
	}
}
