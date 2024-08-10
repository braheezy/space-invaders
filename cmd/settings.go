package cmd

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed fonts/PressStart2P-Regular.ttf
var pressStart2P []byte

var (
	mplusFaceSource *text.GoTextFaceSource
	loadedFont      *text.GoTextFace
)

func init() {
	// Load the embedded font into a GoTextFaceSource
	s, err := text.NewGoTextFaceSource(bytes.NewReader(pressStart2P))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
	loadedFont = &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   12, // Example font size, adjust as needed
	}
}

// Setting represents a general interface for a setting in the menu.
type Setting interface {
	Name() string
	Value() interface{}
	SetValue(interface{}) error
	Render(screen *ebiten.Image, x, y int, selected bool)
}

// OnOffSetting is a simple on/off toggle setting.
type OnOffSetting struct {
	name  string
	value bool
}

func (s *OnOffSetting) Name() string {
	return s.name
}

func (s *OnOffSetting) Value() interface{} {
	return s.value
}

func (s *OnOffSetting) SetValue(val interface{}) error {
	if v, ok := val.(bool); ok {
		s.value = v
		return nil
	}
	return fmt.Errorf("invalid value type")
}

func (s *OnOffSetting) Render(screen *ebiten.Image, x, y int, selected bool) {
	// Set colors for ON and OFF states
	onColor := color.RGBA{0, 255, 0, 255}  // Green color for "ON"
	offColor := color.RGBA{255, 0, 0, 255} // Red color for "OFF"

	// Draw the setting name
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, s.name, loadedFont, op)

	// Determine the text and color based on the value
	var statusText string
	var statusColor color.RGBA
	if s.value {
		statusText = "ON"
		statusColor = onColor
	} else {
		statusText = "OFF"
		statusColor = offColor
	}

	// Draw the ON/OFF status next to the setting name
	op = &text.DrawOptions{}
	op.GeoM.Translate(float64(x+200), float64(y)) // Adjust x+200 as needed for spacing
	op.ColorScale.ScaleWithColor(statusColor)
	text.Draw(screen, statusText, loadedFont, op)
}

// RangeSetting represents a setting that can take a range of values.
type RangeSetting struct {
	name   string
	value  int
	minVal int
	maxVal int
}

func (s *RangeSetting) Name() string {
	return s.name
}

func (s *RangeSetting) Value() interface{} {
	return s.value
}

func (s *RangeSetting) SetValue(val interface{}) error {
	if v, ok := val.(int); ok && v >= s.minVal && v <= s.maxVal {
		s.value = v
		return nil
	}
	return fmt.Errorf("invalid value or out of range")
}

func (s *RangeSetting) Render(screen *ebiten.Image, x, y int, selected bool) {
	// Render the setting name and value (e.g., "5 ships") at the given coordinates
	// (Implement drawing logic here)
}

// HelpSection represents a static section in the menu to display game controls.
type HelpSection struct {
	name     string
	controls []string
}

func (hs *HelpSection) Name() string {
	return hs.name
}

func (hs *HelpSection) Value() interface{} {
	return nil // HelpSection has no modifiable value
}

func (hs *HelpSection) SetValue(val interface{}) error {
	return nil // HelpSection cannot be modified
}

func (hs *HelpSection) Render(screen *ebiten.Image, x, y int, selected bool) {
	// Colors for the bindings and descriptions
	bindingColor := color.RGBA{255, 255, 0, 255}       // Yellow for the keybinding
	descriptionColor := color.RGBA{255, 255, 255, 255} // White for the description

	// Render the help section title
	title := "Help - Game Controls"
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, title, loadedFont, op)

	// Adjust the starting y-coordinate for the controls
	y += 40 // Adjust spacing as needed

	// Render each control in the list
	for i, control := range hs.controls {
		// Split the control into binding and description
		parts := strings.SplitN(control, " - ", 2)
		if len(parts) != 2 {
			continue // Skip malformed entries
		}
		binding := parts[0]
		description := parts[1]

		// Render the keybinding
		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(x), float64(y+(i*30))) // Adjust line height as needed
		op.ColorScale.ScaleWithColor(bindingColor)
		text.Draw(screen, binding, loadedFont, op)

		// Render the description in the second column
		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(x+200), float64(y+(i*30))) // Adjust x+200 for column spacing
		op.ColorScale.ScaleWithColor(descriptionColor)
		text.Draw(screen, description, loadedFont, op)
	}
}

type MenuScreen struct {
	settings      []Setting
	selectedIndex int
	errorMessage  string
	settingsFile  string
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
		&RangeSetting{name: "Ships", value: 3, minVal: 3, maxVal: 6},
		&HelpSection{
			name: "Help - Game Controls",
			controls: []string{
				"Arrow Keys: Move",
				"Space: Shoot",
				"Tab: Toggle Settings",
				"Esc: Quit Game",
				// Add more controls as needed
			},
		},
	}
}

func (ms *MenuScreen) loadSettings() error {
	file, err := os.Open(ms.settingsFile)
	if err != nil {
		return err
	}
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
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		ms.selectedIndex = (ms.selectedIndex + 1) % len(ms.settings)
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		ms.selectedIndex = (ms.selectedIndex - 1 + len(ms.settings)) % len(ms.settings)
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		ms.toggleSelectedSetting()
	}
}

func (ms *MenuScreen) toggleSelectedSetting() {
	selectedSetting := ms.settings[ms.selectedIndex]
	switch setting := selectedSetting.(type) {
	case *OnOffSetting:
		setting.SetValue(!setting.value)
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
	// Iterate through each setting and draw it
	for i, setting := range ms.settings {
		x, y := 50, 50+(i*30) // Adjust coordinates as needed
		selected := i == ms.selectedIndex
		setting.Render(screen, x, y, selected)
	}

	// Draw the error message if it exists
	if ms.errorMessage != "" {
		// Set error message color (e.g., red)
		errorColor := color.RGBA{255, 0, 0, 255}

		// Define the position for the error message (bottom of the screen)
		msgX := 50
		msgY := screen.Bounds().Dy() - 50 // 50 pixels from the bottom

		// Draw the error message
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(msgX), float64(msgY))
		op.ColorScale.ScaleWithColor(errorColor)
		text.Draw(screen, ms.errorMessage, loadedFont, op)
	}
}
