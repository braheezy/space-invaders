package cmd

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// HelpSection represents a static section in the menu to display game controls.
type HelpSection struct {
	name     string
	controls []string
}

func (hs *HelpSection) Name() string {
	return hs.name
}

func (hs *HelpSection) Value() interface{} {
	// HelpSection has no modifiable value
	return nil
}

func (hs *HelpSection) SetValue(val interface{}) error {
	// HelpSection cannot be modified
	return nil
}

func (hs *HelpSection) Render(screen *ebiten.Image, x, y int, selected bool) {
	// Colors for the bindings and descriptions
	// Yellow for the keybinding
	bindingColor := color.RGBA{255, 255, 0, 255}
	// White for the description
	descriptionColor := color.RGBA{255, 255, 255, 255}

	// Render the help section title
	title := "Help - Game Controls"
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(color.RGBA{196, 167, 231, 255})
	text.Draw(screen, title, loadedFont, op)

	// Adjust the starting y-coordinate for the controls
	y += 40

	// Render each control in the list
	for i, control := range hs.controls {
		// Split the control into binding and description
		parts := strings.SplitN(control, " - ", 2)
		if len(parts) != 2 {
			// Skip malformed entries
			continue
		}
		binding := parts[0]
		description := parts[1]

		// Render the keybinding
		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(x), float64(y+(i*30)))
		op.ColorScale.ScaleWithColor(bindingColor)
		text.Draw(screen, binding, loadedFont, op)

		// Render the description in the second column
		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(x+200), float64(y+(i*30)))
		op.ColorScale.ScaleWithColor(descriptionColor)
		text.Draw(screen, description, loadedFont, op)
	}
}
