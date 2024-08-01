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
	DeviceName(port byte) string
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
	CoinDeposited  bool
	soundManager   *SoundManager
	soundMapPort3  map[byte]string
	soundMapPort5  map[byte]string
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
		if si.CoinDeposited {
			result |= 0x01
		}
	case 0x02:
		var result byte
		if ebiten.IsKeyPressed(ebiten.KeyT) {
			result |= 0x04
		}
	default:
		return 0, fmt.Errorf("unsupported hardware port: %02X", addr)
	}

	return result, nil
}

func (si *SpaceInvadersHardware) Out(addr byte, value byte) error {
	switch addr {
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

func (si *SpaceInvadersHardware) DeviceName(port byte) string {
	switch port {
	case 0x01:
		return "INPUT1"
	case 0x02:
		return "INPUT2"
	case 0x03:
		return "SOUND1"
	case 0x05:
		return "SOUND2"
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
