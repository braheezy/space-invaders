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
	ms.settings = []Setting{
		&OnOffSetting{name: "Tilt", value: false},
		&OnOffSetting{name: "Coin Info", value: true},
		// &RangeSetting{name: "Ships", value: 3, minVal: 3, maxVal: 6},
	}

	// Initialize the help section separately
	ms.helpSection = &HelpSection{
		name: "Help - Game Controls",
		controls: []string{
			"Arrow Keys - Move",
			"Space - Shoot",
			"Tab - Toggle Settings",
			"Esc - Quit Game",
			// Add more controls as needed
		},
	}
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
	println("here", err)
	defer file.Close()

	decoder := json.NewDecoder(file)
	loadedSettings := map[string]interface{}{}
	if err := decoder.Decode(&loadedSettings); err != nil {
		return err
	}

	ms.initializeDefaultSettings() // Initialize default settings first

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
	file, err := os.Create(ms.settingsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	settingsMap := make(map[string]interface{})
	for _, setting := range ms.settings {
		settingsMap[setting.Name()] = setting.Value()
	}

	encoder := json.NewEncoder(file)
	return encoder.Encode(settingsMap)
}

func (ms *MenuScreen) Update() {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		ms.selectedIndex = (ms.selectedIndex + 1) % len(ms.settings) // Wrap around to the first setting
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		ms.selectedIndex--
		if ms.selectedIndex < 0 {
			ms.selectedIndex = len(ms.settings) - 1 // Wrap around to the last setting
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		ms.toggleSelectedSetting()
	}
}

func (ms *MenuScreen) toggleSelectedSetting() {
	// Use a debounce approach to avoid flickering
	// if ms.selectedIndex < 0 || ms.selectedIndex >= len(ms.settings) {
	// 	return
	// }

	selectedSetting := ms.settings[ms.selectedIndex]
	switch setting := selectedSetting.(type) {
	case *OnOffSetting:
		if ebiten.IsKeyPressed(ebiten.KeyEnter) { // Debounce to prevent flickering
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

	// Save settings after a change
	if err := ms.saveSettings(); err != nil {
		ms.errorMessage = fmt.Sprintf("Error saving settings: %v", err)
	}
}

func (ms *MenuScreen) Draw(screen *ebiten.Image) {
	// Calculate the start position for drawing
	startX := 50
	startY := 50
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
