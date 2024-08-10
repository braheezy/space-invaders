package cmd

// import (
// 	"encoding/json"
// 	"fmt"
// 	"image/color"
// 	"os"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/text/v2"
// )

// // MenuScreen is a stub for the settings menu screen
// type MenuScreen struct {
// 	settings      []Setting
// 	selectedIndex int
// 	errorMessage  string
// 	settingsFile  string
// }

// // NewMenuScreen creates a new settings menu screen
// func NewMenuScreen() *MenuScreen {
// 	ms := &MenuScreen{
// 		settingsFile: "settings.json",
// 	}
// 	if err := ms.loadSettings(); err != nil {
// 		ms.errorMessage = fmt.Sprintf("Error loading settings: %v", err)
// 		ms.initializeDefaultSettings()
// 	}
// 	return ms
// }

// // initializeDefaultSettings initializes the settings with default values.
// func (ms *MenuScreen) initializeDefaultSettings() {
// 	ms.settings = []Setting{
// 		&OnOffSetting{name: "Tilt", value: false},
// 		&OnOffSetting{name: "Coin Info", value: true},
// 		&RangeSetting{name: "Ships", value: 3, minVal: 3, maxVal: 6},
// 		&HelpSection{
// 			name: "Help - Game Controls",
// 			controls: []string{
// 				"Arrow Keys: Move",
// 				"Space: Shoot",
// 				"Tab: Toggle Settings",
// 				"Esc: Quit Game",
// 				// Add more controls as needed
// 			},
// 		},
// 	}
// }

// // loadSettings loads the settings from a JSON file.
// func (ms *MenuScreen) loadSettings() error {
// 	file, err := os.Open(ms.settingsFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	decoder := json.NewDecoder(file)
// 	loadedSettings := map[string]interface{}{}
// 	if err := decoder.Decode(&loadedSettings); err != nil {
// 		return err
// 	}

// 	ms.initializeDefaultSettings() // Initialize default settings first

// 	for _, setting := range ms.settings {
// 		if value, ok := loadedSettings[setting.Name()]; ok {
// 			if err := setting.SetValue(value); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// // saveSettings saves the settings to a JSON file.
// func (ms *MenuScreen) saveSettings() error {
// 	file, err := os.Create(ms.settingsFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	settingsMap := make(map[string]interface{})
// 	for _, setting := range ms.settings {
// 		settingsMap[setting.Name()] = setting.Value()
// 	}

// 	encoder := json.NewEncoder(file)
// 	return encoder.Encode(settingsMap)
// }

// // Update handles updates for the settings menu.
// func (ms *MenuScreen) Update() {
// 	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
// 		ms.selectedIndex = (ms.selectedIndex + 1) % len(ms.settings)
// 	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
// 		ms.selectedIndex = (ms.selectedIndex - 1 + len(ms.settings)) % len(ms.settings)
// 	}

// 	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
// 		ms.toggleSelectedSetting()
// 	}

// 	// Implement mouse interaction logic as needed.
// }

// func (ms *MenuScreen) toggleSelectedSetting() {
// 	selectedSetting := ms.settings[ms.selectedIndex]
// 	switch setting := selectedSetting.(type) {
// 	case *OnOffSetting:
// 		setting.SetValue(!setting.value)
// 	case *RangeSetting:
// 		newValue := setting.value + 1
// 		if newValue > setting.maxVal {
// 			newValue = setting.minVal
// 		}
// 		setting.SetValue(newValue)
// 	}

// 	// Save settings after a change
// 	if err := ms.saveSettings(); err != nil {
// 		ms.errorMessage = fmt.Sprintf("Error saving settings: %v", err)
// 	}
// }

// func (ms *MenuScreen) Draw(screen *ebiten.Image) {
// 	// Iterate through each setting and draw it
// 	for i, setting := range ms.settings {
// 		x, y := 50, 50+(i*30) // Adjust coordinates as needed
// 		selected := i == ms.selectedIndex
// 		setting.Render(screen, x, y, selected)
// 	}

// 	// Draw the error message if it exists
// 	if ms.errorMessage != "" {
// 		// Set error message color (e.g., red)
// 		errorColor := color.RGBA{255, 0, 0, 255}

// 		// Define the position for the error message (bottom of the screen)
// 		msgX := 50
// 		msgY := screen.Bounds().Dy() - 50 // 50 pixels from the bottom

// 		// Draw the error message
// 		op := &text.DrawOptions{}
// 		op.GeoM.Translate(float64(msgX), float64(msgY))
// 		op.ColorScale.ScaleWithColor(errorColor)
// 		text.Draw(screen, ms.errorMessage, loadedFont, op)
// 	}
// }
