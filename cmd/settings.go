package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed fonts/PressStart2P-Regular.ttf
var pressStart2P []byte

//go:embed images/ship.png
var shipPNG []byte

var (
	shipImage *ebiten.Image
)

var (
	mplusFaceSource *text.GoTextFaceSource
	loadedFont      *text.GoTextFace
)

func NewDefaultSettings() []Setting {
	return []Setting{
		&OnOffSetting{name: "Show coin info on demo screen", value: true},
		&OnOffSetting{name: "Extra ship at 1000 instead of 1500", value: false},
		&OnOffSetting{name: "Limit to 60 FPS", value: false},
		&RangeSetting{name: "Ship Count", value: 3, minVal: 3, maxVal: 6},
	}
}

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

	// Load the ship image
	img, _, err := image.Decode(bytes.NewReader(shipPNG))
	if err != nil {
		log.Fatal(err)
	}
	shipImage = ebiten.NewImageFromImage(img)
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

	// Measure the width of the setting name
	nameWidth, _ := text.Measure(s.name, loadedFont, 1.0)

	// Calculate the position for the ON/OFF status based on the name width
	statusX := float64(x) + nameWidth + 20 // Add some padding after the name

	// Draw the ON/OFF status next to the setting name
	statusOp := &text.DrawOptions{}
	statusOp.GeoM.Translate(statusX, float64(y)) // Position for status text
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
	if v, ok := val.(float64); ok {
		// Since JSON numbers are decoded as float64, you might need this check
		if int(v) >= s.minVal && int(v) <= s.maxVal {
			s.value = int(v)
			return nil
		}
	} else if v, ok := val.(int); ok {
		// Direct int comparison
		if v >= s.minVal && v <= s.maxVal {
			s.value = v
			return nil
		}
	}
	return fmt.Errorf("invalid value or out of range")
}

func (s *RangeSetting) Render(screen *ebiten.Image, x, y int, selected bool) {
	// Render the setting name
	nameOp := &text.DrawOptions{}
	nameOp.GeoM.Translate(float64(x), float64(y))
	nameOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, s.name, loadedFont, nameOp)

	// Calculate the position and size for rendering the ships
	shipX := x + 200 // Start drawing ships to the right of the setting name
	shipY := y
	shipScale := 0.10 // Scale factor for the ship image
	shipWidth := float64(shipImage.Bounds().Dx()) * shipScale
	shipSpacing := shipWidth + 10 // Space between ships

	// Render each ship based on the selected value
	for i := 0; i < s.value; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(shipScale, shipScale) // Apply scaling to the image
		op.GeoM.Translate(float64(shipX)+(float64(i)*shipSpacing), float64(shipY))
		screen.DrawImage(shipImage, op)
	}

	// If selected, indicate that this setting is active
	if selected {
		arrowOp := &text.DrawOptions{}
		arrowOp.GeoM.Translate(float64(x-20), float64(y)) // Arrow to indicate selection
		arrowOp.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, ">", loadedFont, arrowOp)
	}
}
