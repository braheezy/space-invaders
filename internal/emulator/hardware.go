package emulator

import (
	"embed"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/*
var soundFiles embed.FS

type HardwareIO interface {
	In(addr byte) (byte, error)
	Out(addr byte, data byte) error
	InDeviceName(port byte) string
	OutDeviceName(port byte) string
	InterruptConditions() []InterruptCondition
	CyclesPerFrame() int
	Draw(*ebiten.Image)
	Init(*[65536]byte)
	Width() int
	Height() int
	Scale() int
}

type SpaceInvadersHardware struct {
	watchdogTimer  byte
	cyclesPerFrame int
	videoRAM       []byte
	soundManager   *SoundManager
	soundMapPort3  map[byte]string
	soundMapPort5  map[byte]string
	shiftAmount    byte
	shiftRegister  uint16
}

const (
	videoWidth   = 224
	videoHeight  = 256
	displayScale = 3
)

func (si *SpaceInvadersHardware) In(addr byte) (byte, error) {
	var result byte

	switch addr {
	case 0x01:
		if ebiten.IsKeyPressed(ebiten.KeyC) {
			result |= 0x01
		}
	case 0x02:
		/*
					Port 2
			 bit 0 = DIP3 00 = 3 ships  10 = 5 ships
			 bit 1 = DIP5 01 = 4 ships  11 = 6 ships
			 bit 2 = Tilt
			 bit 3 = DIP6 0 = extra ship at 1500, 1 = extra ship at 1000
			 bit 4 = P2 shot (1 if pressed)
			 bit 5 = P2 left (1 if pressed)
			 bit 6 = P2 right (1 if pressed)
			 bit 7 = DIP7 Coin info displayed in demo screen 0=ON
		*/
		// TODO: Make DIP setting user configurable
		dip3 := false
		dip5 := false
		if !dip3 && !dip5 {
			// 3 ships
			result |= 0x00
		} else if dip3 && !dip5 {
			// 4 ships
			result |= 0x01
		} else if !dip3 && dip5 {
			// 5 ships
			result |= 0x02
		} else if dip3 && dip5 {
			// 6 ships
			result |= 0x03
		}

		if ebiten.IsKeyPressed(ebiten.KeyT) {
			result |= 0x04 // Tilt
		}
	case 0x03:
		// Read from the shift register
		shiftedValue := si.shiftRegister >> (8 - si.shiftAmount)
		result = byte(shiftedValue & 0xFF)
	default:
		return 0, fmt.Errorf("unsupported hardware port: %02X", addr)
	}

	return result, nil
}

func (si *SpaceInvadersHardware) Out(addr byte, value byte) error {
	switch addr {
	case 0x02:
		// Set the shift offset, using only the lowest 3 bits
		si.shiftAmount = value & 0x07
	case 0x04:
		// Write to the shift register
		si.shiftRegister = (uint16(value) << 8) | (si.shiftRegister >> 8)
	case 0x03:
		si.handleSoundBits(value, si.soundMapPort3)
	case 0x05:
		si.handleSoundBits(value, si.soundMapPort5)
	case 0x06:
		si.watchdogTimer = value
	default:
		return fmt.Errorf("unsupported hardware port: %02X", addr)
	}
	return nil
}

func (si *SpaceInvadersHardware) OutDeviceName(port byte) string {
	switch port {
	case 0x02:
		return "SHFTAMNT"
	case 0x03:
		return "SOUND1"
	case 0x04:
		return "SHFT_DATA"
	case 0x05:
		return "SOUND2"
	case 0x06:
		return "WATCHDOG"
	default:
		return fmt.Sprintf("$%02X", port) // Default to hex representation if unknown
	}
}

func (si *SpaceInvadersHardware) InDeviceName(port byte) string {
	switch port {
	case 0x01:
		return "INPUT1"
	case 0x02:
		return "INPUT2"
	case 0x03:
		return "SHFT_IN"
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
	soundManager, err := NewSoundManagerWithDefaults(soundFiles)
	if err != nil {
		panic(err)
	}

	soundMapPort3 := map[byte]string{
		0x01: "assets/ufo_repeat_low.qoa",
		0x02: "assets/shoot.wav",
		0x04: "assets/player_die.wav",
		0x08: "assets/invader_die.wav",
		0x10: "assets/extra_play.qoa",
		0x20: "assets/SX5.raw", // AMP enable
	}

	soundMapPort5 := map[byte]string{
		0x01: "assets/fleet_move_1.raw",
		0x02: "assets/fleet_move_2.raw",
		0x04: "assets/fleet_move_3.raw",
		0x08: "assets/fleet_move_4.raw",
		0x10: "assets/ufo_hit.qoa",
	}

	return &SpaceInvadersHardware{
		cyclesPerFrame: 33334,
		soundManager:   soundManager,
		soundMapPort3:  soundMapPort3,
		soundMapPort5:  soundMapPort5,
	}
}

func (si *SpaceInvadersHardware) Init(memory *[65536]byte) {
	si.videoRAM = memory[0x2400:0x4000]
}
func (si *SpaceInvadersHardware) Draw(screen *ebiten.Image) {
	img := ebiten.NewImage(videoWidth, videoHeight)

	// Iterate through each byte in the video RAM
	for i, byteValue := range si.videoRAM {
		// Calculate the original coordinates
		originalX := (i % 32) * 8
		originalY := i / 32

		// Iterate through each bit in the byteValue
		for bit := 0; bit < 8; bit++ {
			// Determine if the current bit is "on" (1) or "off" (0)
			pixelOn := byteValue&(1<<bit) != 0

			// Calculate the original coordinates of the pixel
			x := originalX + bit
			y := originalY

			// Transform coordinates for 90-degree counterclockwise rotation
			rotatedX := y
			rotatedY := videoHeight - 1 - x

			// Set the color of the pixel
			if pixelOn {
				img.Set(rotatedX, rotatedY, color.White) // Set pixel to white if "on"
			} else {
				img.Set(rotatedX, rotatedY, color.Black) // Set pixel to black if "off"
			}
		}
	}

	// Scale and draw the offscreen image to the main screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(displayScale), float64(displayScale))
	screen.DrawImage(img, op)
}

func (si *SpaceInvadersHardware) handleSoundBits(value byte, soundMap map[byte]string) {
	for bit, soundFile := range soundMap {
		if value&bit != 0 {
			si.soundManager.Play(soundFile)
		} else {
			si.soundManager.Pause(soundFile)
		}
	}
}

func (si *SpaceInvadersHardware) Width() int {
	return videoWidth
}
func (si *SpaceInvadersHardware) Height() int {
	return videoHeight
}
func (si *SpaceInvadersHardware) Scale() int {
	return displayScale
}
