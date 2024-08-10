package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"image/color"
	"log"

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

	// Draw the arrow if the setting is selected
	if selected {
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x-20), float64(y)) // Position the arrow slightly to the left
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, ">", loadedFont, op)
	}

	// Draw the setting name
	nameOp := &text.DrawOptions{}
	nameOp.GeoM.Translate(float64(x), float64(y)) // Position for setting name
	nameOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, s.name, loadedFont, nameOp)

	// Draw the ON/OFF status next to the setting name
	statusOp := &text.DrawOptions{}
	statusOp.GeoM.Translate(float64(x+200), float64(y)) // Position for status text
	if s.value {
		statusOp.ColorScale.ScaleWithColor(onColor)
		text.Draw(screen, "ON", loadedFont, statusOp)
	} else {
		statusOp.ColorScale.ScaleWithColor(offColor)
		text.Draw(screen, "OFF", loadedFont, statusOp)
	}
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
