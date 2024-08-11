package invaders

import (
	"embed"
	"fmt"
	"image"
	"time"

	"github.com/braheezy/space-invaders/internal/emulator"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/sounds/*
var soundFiles embed.FS

//go:embed assets/invaders.rom
var romFile embed.FS

//go:embed assets/SpaceInvadersArcColorUseCV.png
var cvColorOverlay embed.FS

type ColorScheme int

const (
	BlackAndWhite ColorScheme = iota
	TV
	CV
)

var ColorSchemeNames = []string{
	"BW",
	"TV",
	"CV",
}

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

	// video is the current image to display, the graphics for the current frame
	video  *ebiten.Image
	pixels []byte

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

	ColorScheme    ColorScheme
	cvColorOverlay image.Image

	// DIP switch settings
	// 3, 4, 5, or 6 ships
	ShipsSetting int
	// true = extra ship at 1000, false = extra ship at 1500
	ExtraShipAt1000 bool
	// true = show coin info, false = don't show
	ShowCoinInfoOnDemo bool
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

	cvColorImageFile, _ := cvColorOverlay.Open("assets/SpaceInvadersArcColorUseCV.png")
	img, _, _ := image.Decode(cvColorImageFile)

	return &SpaceInvadersHardware{
		cyclesPerFrame: 33334,
		soundManager:   soundManager,
		soundMapPort3:  soundMapPort3,
		soundMapPort5:  soundMapPort5,
		rom:            romData,
		ColorScheme:    BlackAndWhite,
		cvColorOverlay: img,
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
			// Set bit 3
			result |= 0x08
		}

		// Show coin info on demo screen
		if si.ShowCoinInfoOnDemo {
			// Set bit 7
			result |= 0x80
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

	si.video = ebiten.NewImage(videoWidth, videoHeight)
	si.pixels = make([]byte, videoWidth*videoHeight*4)
}

func (si *SpaceInvadersHardware) Draw(screen *ebiten.Image) {
	// Iterate through each byte in the video RAM
	for i, byteValue := range si.videoRAM {
		originalX := (i % 32) * 8
		originalY := i / 32

		for bit := 0; bit < 8; bit++ {
			pixelOn := byteValue&(1<<bit) != 0
			x := originalX + bit
			y := originalY

			rotatedX := y
			rotatedY := videoHeight - 1 - x

			// Calculate the pixel's index in the array
			index := (rotatedY*videoWidth + rotatedX) * 4

			if pixelOn {
				switch si.ColorScheme {
				case BlackAndWhite:
					si.pixels[index] = 0xFF   // R
					si.pixels[index+1] = 0xFF // G
					si.pixels[index+2] = 0xFF // B
					si.pixels[index+3] = 0xFF // A

				case TV:
					if rotatedY >= 16 && rotatedY < 32 {
						// Red region
						si.pixels[index] = 0xFF   // R
						si.pixels[index+1] = 0x00 // G
						si.pixels[index+2] = 0x00 // B
						si.pixels[index+3] = 0xFF // A
					} else if rotatedY >= (videoHeight - 72) {
						if !(rotatedY >= (videoHeight-16) && rotatedX < 25) && // Bottom left cutout
							!(rotatedY >= (videoHeight-16) && rotatedX >= (videoWidth-88)) { // Bottom right cutout
							// Green region
							si.pixels[index] = 0x00   // R
							si.pixels[index+1] = 0xFF // G
							si.pixels[index+2] = 0x00 // B
							si.pixels[index+3] = 0xFF // A
						} else {
							// Cutout region or area outside the green overlay
							si.pixels[index] = 0xFF   // R
							si.pixels[index+1] = 0xFF // G
							si.pixels[index+2] = 0xFF // B
							si.pixels[index+3] = 0xFF // A
						}
					} else {
						// Default to white if not in any special region
						si.pixels[index] = 0xFF   // R
						si.pixels[index+1] = 0xFF // G
						si.pixels[index+2] = 0xFF // B
						si.pixels[index+3] = 0xFF // A
					}
				case CV:
					if si.cvColorOverlay == nil {
						si.pixels[index] = 0x00   // R
						si.pixels[index+1] = 0x00 // G
						si.pixels[index+2] = 0x00 // B
						si.pixels[index+3] = 0xFF // A
					} else {
						if pixelOn {

							// Sample the color from the CV overlay
							r, g, b, a := si.sampleCVColor(rotatedX, rotatedY)

							si.pixels[index] = r   // R
							si.pixels[index+1] = g // G
							si.pixels[index+2] = b // B
							si.pixels[index+3] = a // A
						} else {
							si.pixels[index] = 0x00   // R
							si.pixels[index+1] = 0x00 // G
							si.pixels[index+2] = 0x00 // B
							si.pixels[index+3] = 0xFF // A
						}
					}
				}

			} else {
				si.pixels[index] = 0x00   // R
				si.pixels[index+1] = 0x00 // G
				si.pixels[index+2] = 0x00 // B
				si.pixels[index+3] = 0xFF // A
			}
		}
	}

	// Write the pixel data to the image
	si.video.WritePixels(si.pixels)

	// Scale and draw the offscreen image to the main screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(displayScale), float64(displayScale))
	screen.DrawImage(si.video, op)
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

// Helper function to sample color from the CV overlay
// ! Because the overlay reference has game sprites also, certain samples will result
// ! in a different color.
// ! All sampled RGB values for all color are either 0, 128 (the base color), or 255 (illuminated when a
// ! sprite is drawn.
func (si *SpaceInvadersHardware) sampleCVColor(x, y int) (r, g, b, a uint8) {
	// Get the color from the overlay
	overlayColor := si.cvColorOverlay.At(x, y)
	r32, g32, b32, _ := overlayColor.RGBA()

	// Convert to 8-bit color values
	r8 := uint8(r32 >> 8)
	g8 := uint8(g32 >> 8)
	b8 := uint8(b32 >> 8)

	// Determine if the color needs to be doubled
	if r8 == 128 {
		r8 = 255
	}
	if g8 == 128 {
		g8 = 255
	}
	if b8 == 128 {
		b8 = 255
	}

	// Use the calculated color for "on" pixels
	return r8, g8, b8, 0xFF
}
