package invaders

import (
	"embed"
	"fmt"
	"image/color"
	"time"

	"github.com/braheezy/space-invaders/internal/emulator"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/sounds/*
var soundFiles embed.FS

//go:embed assets/invaders.rom
var romFile embed.FS

// SpaceInvadersHardware represents the hardware-specific implementation for the Space Invaders game.
type SpaceInvadersHardware struct {
	// watchdogTimer is used to simulate the watchdog timer functionality, which resets the game
	// if the CPU gets stuck in an infinite loop or takes too long to complete an operation.
	watchdogTimer byte

	// cyclesPerFrame defines the number of CPU cycles that should be executed per frame.
	// This value helps in synchronizing the CPU execution with the display refresh rate.
	cyclesPerFrame int

	// videoRAM holds the video memory where the graphical data for the display is stored.
	// This memory is updated by the CPU to reflect changes in the game graphics.
	videoRAM []byte

	// soundManager manages the playback of sound effects for the game.
	soundManager *emulator.SoundManager

	// soundMapPort3 maps the bits in port 3 to their corresponding sound file names.
	// This mapping is used to determine which sound to play when a bit in port 3 is set.
	soundMapPort3 map[byte]string

	// soundMapPort5 maps the bits in port 5 to their corresponding sound file names.
	// This mapping is used to determine which sound to play when a bit in port 5 is set.
	soundMapPort5 map[byte]string

	// shiftAmount specifies the number of positions to shift the data in the shift register.
	// This value is set by writing to port 2.
	shiftAmount byte

	// shiftRegister holds the 16-bit shift register value.
	// This register is used to shift data for certain operations, as specified by the hardware.
	shiftRegister uint16

	// rom contains the read-only memory (ROM) data for the Space Invaders game.
	// This data includes the game code and other static information needed by the emulator.
	rom []byte

	// lastSound1 holds the previous state of the sound control bits for port 3.
	// This value is used to detect changes in the sound control bits and play the corresponding sounds.
	lastSound1 byte

	// lastSound2 holds the previous state of the sound control bits for port 5.
	// This value is used to detect changes in the sound control bits and play the corresponding sounds.
	lastSound2 byte

	// DIP switch settings
	ShipsSetting       int  // 3, 4, 5, or 6 ships
	ExtraShipAt1000    bool // true = extra ship at 1000, false = extra ship at 1500
	ShowCoinInfoOnDemo bool // true = show coin info, false = don't show
}

const (
	videoWidth   = 224
	videoHeight  = 256
	displayScale = 3
	startAddress = 0x0
)

func NewSpaceInvadersHardware() *SpaceInvadersHardware {
	soundManager, err := emulator.NewSoundManager(44100, 1, soundFiles)
	if err != nil {
		panic(err)
	}

	soundMapPort3 := map[byte]string{
		0: "assets/sounds/ufo_repeat_low.qoa",
		1: "assets/sounds/shoot.qoa",
		2: "assets/sounds/player_die.qoa",
		3: "assets/sounds/invader_die.qoa",
		4: "assets/sounds/extra_play.qoa",
	}

	soundMapPort5 := map[byte]string{
		0: "assets/sounds/fleet_move_1.qoa",
		1: "assets/sounds/fleet_move_2.qoa",
		2: "assets/sounds/fleet_move_3.qoa",
		3: "assets/sounds/fleet_move_4.qoa",
		4: "assets/sounds/ufo_hit.qoa",
	}

	romData, _ := romFile.ReadFile("assets/invaders.rom")

	return &SpaceInvadersHardware{
		cyclesPerFrame: 33334,
		soundManager:   soundManager,
		soundMapPort3:  soundMapPort3,
		soundMapPort5:  soundMapPort5,
		rom:            romData,
	}
}

func (si *SpaceInvadersHardware) In(addr byte) (byte, error) {
	var result byte

	switch addr {
	case 0x01:
		/*
					Port 1
			 bit 0 = CREDIT (1 if deposit)
			 bit 1 = 2P start (1 if pressed)
			 bit 2 = 1P start (1 if pressed)
			 bit 3 = Always 1
			 bit 4 = 1P shot (1 if pressed)
			 bit 5 = 1P left (1 if pressed)
			 bit 6 = 1P right (1 if pressed)
			 bit 7 = Not connected
		*/
		// Credit button aka insert coin
		if ebiten.IsKeyPressed(ebiten.KeyC) {
			result |= 0x01
		}
		// Player 2 start
		if ebiten.IsKeyPressed(ebiten.Key2) {
			result |= 0x02
		}
		// Player 1 start
		if ebiten.IsKeyPressed(ebiten.Key1) {
			result |= 0x04
		}
		// Player 1 shoot
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			result |= 0x10
		}
		// Player 1 left
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			result |= 0x20
		}
		// Player 1 right
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			result |= 0x40
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
		switch si.ShipsSetting {
		case 3:
			result |= 0x00
		case 4:
			result |= 0x01
		case 5:
			result |= 0x02
		case 6:
			result |= 0x03
		}

		// Extra ship at 1000 or 1500
		if si.ExtraShipAt1000 {
			result |= 0x08 // Set bit 3 if extra ship is at 1000
		}

		// Show coin info on demo screen
		if si.ShowCoinInfoOnDemo {
			result |= 0x80 // Set bit 7 to display coin info in demo
		}

		// Tilt
		if ebiten.IsKeyPressed(ebiten.KeyT) {
			result |= 0x04
		}
		// Player 2 shoot
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			result |= 0x10
		}
		// Player 2 left
		if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			result |= 0x20
		}
		// Player 2 right
		if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			result |= 0x40
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
		si.handleSoundBits(value, si.soundMapPort3, &si.lastSound1)
	case 0x05:
		si.handleSoundBits(value, si.soundMapPort5, &si.lastSound2)
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
		// Default to hex representation if unknown
		return fmt.Sprintf("$%02X", port)
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
		// Default to hex representation if unknown
		return fmt.Sprintf("$%02X", port)
	}
}

func (si *SpaceInvadersHardware) InterruptConditions() []emulator.Interrupt {
	return []emulator.Interrupt{
		{
			// Mid screen interrupt
			Cycle: si.cyclesPerFrame / 2,
			Action: func(vm *emulator.CPU8080) {
				// RST 8
				vm.InterruptRequest <- 0xCF
			},
		},
		{
			// VBLANK interrupt
			Cycle: si.cyclesPerFrame,
			Action: func(vm *emulator.CPU8080) {
				// RST 10
				vm.InterruptRequest <- 0xD7
			},
		},
	}
}

func (si *SpaceInvadersHardware) CyclesPerFrame() int {
	return si.cyclesPerFrame
}

func (si *SpaceInvadersHardware) Init(memory *[65536]byte) {
	// memory location 0x2400 to 0x3FFF contain the graphic data
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
				img.Set(rotatedX, rotatedY, color.White)
			} else {
				img.Set(rotatedX, rotatedY, color.Black)
			}
		}
	}

	// Scale and draw the offscreen image to the main screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(displayScale), float64(displayScale))
	screen.DrawImage(img, op)
}

// handleSoundBits handles the playing of sound effects based on changes in the sound control bits.
// It checks which bits are set in the provided value, compares them with the previous value, and
// triggers the sound playback if a bit transitions from 0 to 1.
//
// Parameters:
// - value: The current state of the sound control bits.
// - soundMap: A map where the key is the bit position and the value is the corresponding sound file.
// - lastValue: A pointer to the previous state of the sound control bits, used to detect changes.
func (si *SpaceInvadersHardware) handleSoundBits(value byte, soundMap map[byte]string, lastValue *byte) {
	for bit, soundFile := range soundMap {
		// Create a bitmask for the current bit position
		bitMask := 1 << bit
		// Check if the current bit is set in the new value
		currentBitSet := value & byte(bitMask)
		// Check if the current bit was set in the previous value
		lastBitSet := *lastValue & byte(bitMask)
		// If the bit transitioned from 0 to 1...
		if currentBitSet != 0 && lastBitSet == 0 {
			// Play the corresponding sound
			si.soundManager.Play(soundFile)
		}
	}
	*lastValue = value
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
func (si *SpaceInvadersHardware) StartAddress() int {
	return startAddress
}
func (si *SpaceInvadersHardware) HandleSystemCall(*emulator.CPU8080) {
	// no system calls in space invaders!
}
func (si *SpaceInvadersHardware) ROM() []byte {
	return si.rom
}
func (si *SpaceInvadersHardware) FrameDuration() time.Duration {
	// 60 FPS -> 1000ms / 60 = 16.67ms per frame, approximate to 17ms
	return 17 * time.Millisecond
}
func (si *SpaceInvadersHardware) Cleanup() {
	si.soundManager.Cleanup()
}
