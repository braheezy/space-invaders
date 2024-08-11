package cmd

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type MenuScreen struct {
	settings      []Setting
	selectedIndex int
	errorMessage  string
	settingsFile  string
	helpSection   *HelpSection
}

func NewMenuScreen(settingsFile string) *MenuScreen {
	ms := &MenuScreen{
		settingsFile: settingsFile,
	}
	if err := ms.loadSettings(); err != nil {
		ms.errorMessage = fmt.Sprintf("Error loading settings: %v", err)
		ms.initializeDefaultSettings()
	}
	return ms
}

func (ms *MenuScreen) initializeDefaultSettings() {
	// Clone the default settings to ms.settings
	ms.settings = NewDefaultSettings()

	// Initialize the help section separately
	ms.helpSection = &HelpSection{
		name: "Game Controls",
		controls: []string{
			"Arrow Keys/WASD - Move, Navigate menu",
			"Space - Shoot",
			"C - Insert credit",
			"1 - Player 1 Start",
			"2 - Player 2 Start",
			"T - Tilt",
			"Enter - Toggle setting",
			"Tab - Toggle menu",
			"Esc - Quit",
		},
	}
}

func (ms *MenuScreen) GetLimitTPS() bool {
	for _, setting := range ms.settings {
		if onOffSetting, ok := setting.(*OnOffSetting); ok && onOffSetting.name == "Limit to 60 FPS" {
			return onOffSetting.value
		}
	}
	return true
}

func (ms *MenuScreen) GetShipsSetting() int {
	for _, setting := range ms.settings {
		if rangeSetting, ok := setting.(*RangeSetting); ok && rangeSetting.name == "Ship Count" {
			return rangeSetting.value
		}
	}
	return 3 // Default to 3 ships if not found
}

func (ms *MenuScreen) GetExtraShipAt1000() bool {
	for _, setting := range ms.settings {
		if onOffSetting, ok := setting.(*OnOffSetting); ok && onOffSetting.name == "Extra ship at 1000 instead of 1500" {
			return onOffSetting.value
		}
	}
	return false // Default to extra ship at 1500
}

func (ms *MenuScreen) GetShowCoinInfoOnDemo() bool {
	for _, setting := range ms.settings {
		if onOffSetting, ok := setting.(*OnOffSetting); ok && onOffSetting.name == "Show coin info on demo screen" {
			return onOffSetting.value
		}
	}
	return true // Default to showing coin info
}

func (ms *MenuScreen) loadSettings() error {
	file, err := os.Open(ms.settingsFile)
	if os.IsNotExist(err) {
		// If the file does not exist, initialize default settings and return nil
		ms.initializeDefaultSettings()
		return nil
	} else if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	loadedSettings := map[string]interface{}{}
	if err := decoder.Decode(&loadedSettings); err != nil {
		return err
	}

	ms.initializeDefaultSettings()

	for _, setting := range ms.settings {
		if value, ok := loadedSettings[setting.Name()]; ok {
			if err := setting.SetValue(value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ms *MenuScreen) saveSettings() error {
	defaultSettings := NewDefaultSettings()
	settingsMap := make(map[string]interface{})

	for _, setting := range ms.settings {
		for _, defaultSetting := range defaultSettings {
			if setting.Name() == defaultSetting.Name() && setting.Value() != defaultSetting.Value() {
				settingsMap[setting.Name()] = setting.Value()
				break
			}
		}
	}

	if len(settingsMap) == 0 {
		// No settings to save, skip creating or writing to the file
		return nil
	}

	// Now proceed to create and save the file only if there are settings to save
	file, err := os.Create(ms.settingsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(settingsMap)
}

func (ms *MenuScreen) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		ms.selectedIndex = (ms.selectedIndex + 1) % len(ms.settings) // Wrap around to the first setting
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		ms.selectedIndex--
		if ms.selectedIndex < 0 {
			ms.selectedIndex = len(ms.settings) - 1 // Wrap around to the last setting
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ms.toggleSelectedSetting()
	}

	// Handle left/right arrow keys for range settings
	selectedSetting := ms.settings[ms.selectedIndex]
	if rangeSetting, ok := selectedSetting.(*RangeSetting); ok {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			if rangeSetting.value > rangeSetting.minVal {
				rangeSetting.value--
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			if rangeSetting.value < rangeSetting.maxVal {
				rangeSetting.value++
			}
		}
	}
}

func (ms *MenuScreen) toggleSelectedSetting() {

	selectedSetting := ms.settings[ms.selectedIndex]
	switch setting := selectedSetting.(type) {
	case *OnOffSetting:
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			setting.SetValue(!setting.value)
		}
	case *RangeSetting:
		// Cycle through the range
		newValue := setting.value + 1
		if newValue > setting.maxVal {
			newValue = setting.minVal
		}
		setting.SetValue(newValue)
	}

}

func (ms *MenuScreen) Draw(screen *ebiten.Image) {
	// Render the overall menu title
	menuTitle := "Settings Menu"
	titleOp := &text.DrawOptions{}
	titleOp.GeoM.Translate(float64(50), float64(20)) // Position at the top of the screen
	titleOp.ColorScale.ScaleWithColor(color.RGBA{196, 167, 231, 255})
	text.Draw(screen, menuTitle, loadedFont, titleOp)

	// Calculate the start position for drawing the settings
	startX := 50
	startY := 70 // Adjust startY to account for the title
	lineHeight := 30

	// Iterate through each setting and draw it
	for i, setting := range ms.settings {
		y := startY + (i * lineHeight) // Calculate the Y position for each setting
		selected := i == ms.selectedIndex
		setting.Render(screen, startX, y, selected)
	}

	// Draw the help section below the settings
	if ms.helpSection != nil {
		y := startY + (len(ms.settings) * lineHeight) + 50 // Space between settings and help section
		ms.helpSection.Render(screen, startX, y, false)
	}

	// Draw the error message if it exists
	if ms.errorMessage != "" {
		errorColor := color.RGBA{255, 0, 0, 255}
		msgX := 50
		msgY := screen.Bounds().Dy() - 100 // Start drawing a bit higher for word wrapping

		// Word wrap the error message
		lines := wordWrap(ms.errorMessage, 50) // Wrap at 50 characters per line
		for i, line := range lines {
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(msgX), float64(msgY+(i*lineHeight))) // Use consistent line height
			op.ColorScale.ScaleWithColor(errorColor)
			text.Draw(screen, line, loadedFont, op)
		}
	}
}

// Helper function to word wrap text
func wordWrap(text string, maxLineLength int) []string {
	words := strings.Fields(text)
	var lines []string
	var currentLine string

	for _, word := range words {
		if len(currentLine)+len(word)+1 > maxLineLength {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			if len(currentLine) > 0 {
				currentLine += " "
			}
			currentLine += word
		}
	}
	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}
