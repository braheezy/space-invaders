package emulator

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// HardwareIO defines the interface for hardware input/output operations,
// interrupt handling, rendering, and initialization within an 8080-based emulator.
type HardwareIO interface {

	// In reads a byte from the given I/O port address.
	// Returns the byte read, or an error if the operation fails.
	// This will be called when the IN opcode is executed.
	In(addr byte) (byte, error)

	// Out writes a byte to the given I/O port address.
	// Returns an error if the operation fails.
	// This will be called when the OUT opcode is executed.
	Out(addr byte, data byte) error

	// InDeviceName returns the name of the input device for the given I/O port.
	// This is called to print a meaningful name for the device in error messages.
	InDeviceName(port byte) string

	// OutDeviceName returns the name of the output device for the given I/O port.
	// This is called to print a meaningful name for the device in error messages.
	OutDeviceName(port byte) string

	// InterruptConditions returns a slice of interrupt conditions that specify when
	// and how interrupts should be triggered within the emulator.
	InterruptConditions() []Interrupt

	// CyclesPerFrame returns the number of CPU cycles that should be executed per frame.
	CyclesPerFrame() int

	// Draw renders the current state of the hardware to the given ebiten.Image.
	Draw(screen *ebiten.Image)

	// Init initializes the hardware with a reference to RAM.
	// This is called before the emulator executes codes, giving the hardware a chance
	// to do what it needs to do with RAM.
	Init(memory *[65536]byte)

	// Width returns the width of the hardware display in pixels.
	Width() int

	// Height returns the height of the hardware display in pixels.
	Height() int

	// Scale returns the scale factor for rendering the display.
	Scale() int

	// HandleSystemCall processes system calls made by the CPU, typically
	// for input/output operations or other hardware interactions.
	HandleSystemCall(*CPU8080)

	// StartAddress returns the address in memory where execution should begin.
	StartAddress() int

	// ROM holds the data for the read-only memory region.
	// This method is typically used for loading the ROM data into memory.
	// Returns a slice of bytes representing the ROM data.
	ROM() []byte

	// FrameDuration returns the duration of a single frame in milliseconds
	FrameDuration() time.Duration
}

// NullHardware implements HardwareIO with no-op methods for testing purposes.
type NullHardware struct{}

func (nh *NullHardware) In(addr byte) (byte, error) {
	return 0, nil
}
func (nh *NullHardware) Out(addr byte, data byte) error {
	return nil
}
func (nh *NullHardware) InDeviceName(port byte) string {
	return "Null Input Device"
}
func (nh *NullHardware) OutDeviceName(port byte) string {
	return "Null Output Device"
}
func (nh *NullHardware) InterruptConditions() []Interrupt {
	return nil
}
func (nh *NullHardware) CyclesPerFrame() int {
	// Taken from space invaders
	return 33334
}
func (nh *NullHardware) Draw(screen *ebiten.Image) {
	// No-op
}
func (nh *NullHardware) Init(memory *[65536]byte) {
	// No-op
}
func (nh *NullHardware) Width() int {
	return 850
}
func (nh *NullHardware) Height() int {
	return 600
}
func (nh *NullHardware) Scale() int {
	return 1
}
func (nh *NullHardware) HandleSystemCall(cpu *CPU8080) {
	// No-op
}
func (nh *NullHardware) StartAddress() int {
	return 0
}
func (nh *NullHardware) ROM() []byte {
	return []byte{}
}
func (nh *NullHardware) FrameDuration() time.Duration {
	return 17 * time.Millisecond
}
