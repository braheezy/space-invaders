package emulator

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type HardwareIO interface {
	In(addr byte) byte
	Out(addr byte, data byte)
	DeviceName(port byte) string
	InterruptConditions() []InterruptCondition
	CyclesPerFrame() int
	Draw(*ebiten.Image)
	Init(*[65536]byte)
}

type SpaceInvadersHardware struct {
	watchdogTimer  byte
	cyclesPerFrame int
	DisplayScale   int
	videoRAM       []byte
}

func (si *SpaceInvadersHardware) In(addr byte) byte {
	switch addr {
	case 0x02:
		var result byte
		if ebiten.IsKeyPressed(ebiten.KeyT) {
			result |= 0x04
		}
		return result
	default:
		return 0xFF // Return 0xFF if the address is unknown
	}
}

func (si *SpaceInvadersHardware) Out(addr byte, value byte) {
	if addr == 0x03 {
		si.watchdogTimer = value
	}
}

func (si *SpaceInvadersHardware) DeviceName(port byte) string {
	switch port {
	case 0x02:
		return "INPUT2"
	case 0x06:
		return "WATCHDOG"
	default:
		return fmt.Sprintf("$%02X", port) // Default to hex representation if unknown
	}
}

func (si *SpaceInvadersHardware) InterruptConditions() []InterruptCondition {
	return []InterruptCondition{
		{
			// Mid screen interrupt
			Cycle: si.cyclesPerFrame / 2,
			Action: func(vm *CPU8080) {
				// RST 8
				vm.interruptRequest <- 0xCF
			},
		},
		{
			// VBLANK interrupt
			Cycle: si.cyclesPerFrame,
			Action: func(vm *CPU8080) {
				// RST 10
				vm.interruptRequest <- 0xD7
			},
		},
	}
}

func (si *SpaceInvadersHardware) CyclesPerFrame() int {
	return si.cyclesPerFrame // The constant value for Space Invaders
}

func NewSpaceInvadersHardware() *SpaceInvadersHardware {
	return &SpaceInvadersHardware{
		cyclesPerFrame: 33334,
		DisplayScale:   5,
	}
}

func (si *SpaceInvadersHardware) Init(memory *[65536]byte) {
	si.videoRAM = (*memory)[0x2400 : 0x3FFF+1]
}

func (si *SpaceInvadersHardware) Draw(screen *ebiten.Image) {
	scale := si.DisplayScale // Assume this is set to an appropriate scaling factor

	// Create a single pixel image for reuse in drawing each "on" pixel
	pixelImg := ebiten.NewImage(1, 1) // 1x1 pixel image

	// Iterate through each byte in the video RAM
	for i, byteValue := range si.videoRAM {
		// Calculate the corresponding x and y position on the rotated screen
		x := i / 0x20
		yBase := i % 0x20 * 8 // Start y position for this byte

		// Iterate through each bit in the byteValue
		for bit := 0; bit < 8; bit++ {
			// Determine if the current bit is "on" (1) or "off" (0)
			pixelOn := byteValue&(1<<bit) != 0
			y := yBase + (7 - bit) // Calculate y position, adjusting for bit position

			// Set the color of the pixelImg based on the pixel state
			if pixelOn {
				pixelImg.Fill(color.White) // Set pixel to white if "on"
			} else {
				continue // Skip drawing for "off" pixels to enhance performance
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(scale), float64(scale))         // Scale up the pixel
			op.GeoM.Translate(float64(x*scale), float64(y*scale)) // Move the pixel to its proper location

			// Draw the scaled pixel image to the screen
			screen.DrawImage(pixelImg, op)
		}
	}
}
